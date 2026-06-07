# Plan: Closures + Lambda Codegen

Spec: [`0004-TypeSystem.md` — Closures](../specs/0004-TypeSystem.md#closures--type-fn-closure) ([TYPE-FN-CLOSURE]).

Parent: [`production-primitives.md`](production-primitives.md).
Downstream: gating dependency for [`backend-framework.md`](backend-framework.md) (composable
`Middleware = fn(Handler) -> Handler`).

## Problem (UPDATED — Phase 1 has landed; capture is the live work)

**Phase 1 is done.** Non-capturing lambdas now work in **every** position — `let`, argument, and record
field. Re-confirmed this session by compiling:

- `let f = fn(x: int) => x + 1; print(f(10))` → `11` ✅
- function value passed as an arg, and a non-capturing lambda stored in a record field → work ✅
- `let f = record.fnField; f(x)` (call a function-valued field via a local) → works ✅

**What is still broken: capture.** Any lambda that references an outer binding emits dangling SSA refs:

| Probe | llc error |
|---|---|
| `let p="X"; fn(b) => p + b` (capture outer `let`) | `instruction forward referenced with type 'ptr'` |
| `fn(h) => fn(b) => h(b)` (curried middleware) | `use of undefined value '%h'` |
| `fn adder(n) = fn(x) => x + n` (named fn returns capturing lambda) | `use of undefined value '%n'` |
| `obj.fnField(args)` directly | `UFCS call _.op(...) rewrites to op(_, ...): function not declared: op` |

## Why it's broken (CONFIRMED root cause)

A function VALUE is a **bare LLVM function pointer**, never a closure pair `{fn, env}`
([`getLLVMFunctionType`, function_signatures.go:1272](../../compiler/internal/codegen/function_signatures.go#L1272)).
`generateLambdaExpression` ([fiber_generation.go:402](../../compiler/internal/codegen/fiber_generation.go#L402))
emits a fresh top-level func and binds the lambda's *own params*, but performs **zero capture** — it
saves/restores `g.variables` but never threads outer values into the new function. So a lambda body's
reference to `p`/`h`/`n` resolves against `g.variables`, which still holds the **parent** function's SSA
values — out of scope in the emitted lambda, hence the dangling refs. The `wrap`/`adder` cases additionally
need the env to survive *after the outer function returns* (escaping closures).

Only the **`spawn`** path was fixed, via a **global-spill hack**: `captureVariablesInExpression`
([fiber_generation.go:603](../../compiler/internal/codegen/fiber_generation.go#L603)) collects free vars,
`spillNonConstantCaptures` ([fiber_generation.go:574](../../compiler/internal/codegen/fiber_generation.go#L574))
stores each into a per-spawn **module global**, and the closure reloads them. That's fine for `spawn`'s
fire-once semantics but is **not** a correct general closure model (breaks re-entrancy, recursion, multiple
instances, threads). General lambdas need a proper per-instance heap env (Phase 2).

## Phase 1 — Make the no-capture case work ✅ LANDED

Goal (met): `let f = fn(x: int) => x + 1; print(f(10))` prints `11`; string-typed and multi-param
lambdas, lambdas-as-args, and non-capturing lambdas in record fields all compile and run. Keep/extend the
golden coverage when touching the lambda path in Phase 2 so this never regresses.

## Phase 2 — Capture

Goal: `fn makeAdder(n: int) -> (int) -> int = fn(x: int) => x + n` works, and the returned closure remains callable after `makeAdder` returns.

- [ ] **2.1** Add a failing test for `makeAdder(5)(3) == 8` and `let add5 = makeAdder(5); add5(3) == 8`.
- [ ] **2.2** Capture strategy — **closure-pair fat pointer** (confirmed by this session's grounding).
  A function VALUE becomes a 2-field struct `osprey.closure = { i8* fnptr, i8* env }` (define once in the
  generator). `env` is `null` for non-capturing lambdas / top-level functions, so capture-free lambdas pay
  no allocation. Change `getLLVMFunctionType`
  ([function_signatures.go:1272](../../compiler/internal/codegen/function_signatures.go#L1272)) +
  the string-fallback at [:1208-1216](../../compiler/internal/codegen/function_signatures.go#L1208-L1216).
- [ ] **2.3** Free-variable analysis — **reuse the existing walker**
  `captureVariablesInExpression` ([fiber_generation.go:603](../../compiler/internal/codegen/fiber_generation.go#L603))
  rather than a parallel implementation (no-duplication rule). Two required upgrades to it: (a) **subtract
  the lambda's own parameters** from the captured set; (b) **recurse into** `LambdaExpression`,
  `MatchExpression`, and `FieldAccessExpression` (it currently handles only Identifier/Call/Binary/Unary —
  incomplete for nested lambdas). If a separate file is warranted, `internal/codegen/closure_conversion.go`
  (under 500 LOC), but prefer extending the existing walker.
- [ ] **2.4** In `generateLambdaExpression` ([fiber_generation.go:402](../../compiler/internal/codegen/fiber_generation.go#L402)):
  for each captured var, **`malloc` a heap env struct** `{ captured fields… }` (NOT a module global — it
  must outlive the enclosing frame for escaping closures like `adder`), store current values via the parent
  builder, `bitcast` to `i8*`. Add a hidden leading `i8* %__env` param to the lambda; at entry, bitcast +
  load each captured field into `g.variables` before body codegen (mirrors the spawn reload at
  [fiber_generation.go:162-164](../../compiler/internal/codegen/fiber_generation.go#L162-L164) but
  per-instance). Return the closure-pair `{ bitcast(lambda to i8*), env }`.
- [ ] **2.5** Call sites — in `generateCallExpression`/`resolveFunctionValue`/`validateCallableType`
  ([llvm.go:25](../../compiler/internal/codegen/llvm.go#L25),
  [:102](../../compiler/internal/codegen/llvm.go#L102),
  [:167](../../compiler/internal/codegen/llvm.go#L167)): when the callee is a closure pair, extract field 0
  (fnptr) + field 1 (env) and emit `NewCall(fnptr, env, args...)` ([:61](../../compiler/internal/codegen/llvm.go#L61)).
  **Preserve the direct-call fast path** for named top-level funcs (`resolveMonomorphizedFunction`,
  [llvm.go:146](../../compiler/internal/codegen/llvm.go#L146)) so the bulk of the codebase is unperturbed.
  Simplest ABI: all function-value signatures take a leading `i8*` env that top-level callees ignore.
- [ ] **2.6** Heap-allocate the env via `malloc` (already linked). Document the memory model in the file
  header. **Do NOT touch the effect-handler ABI** ([effects_generation.go:416-438](../../compiler/internal/codegen/effects_generation.go#L416-L438))
  — handlers are non-capturing; wrap as `{fn, null}` only if a handler is ever used as a first-class value.
- [ ] **2.7** Function-typed record fields ([function_signatures.go] layout, field access
  [expression_generation.go:1631](../../compiler/internal/codegen/expression_generation.go#L1631)) become
  the closure-pair struct — so storing a *capturing* lambda in a field now works, and `let f = b.op; f(x)`
  keeps working.
- [ ] **2.8** Tests from 2.1 must pass; add `wrap`/`adder`/curried-middleware cases.
- [ ] **2.9** (cleanup) Route `generateSpawnExpression` through the shared env mechanism and **delete
  `spillNonConstantCaptures`** ([fiber_generation.go:574](../../compiler/internal/codegen/fiber_generation.go#L574)) —
  removes the duplicate, buggy global-spill capture logic.

## Phase 3 — Higher-order use sites

Goal: closures work as arguments to `map`, `filter`, `fold`, `forEach`, and as record fields.

- [ ] **3.1** Test: `[1, 2, 3] |> map(fn(x) => x * 2)` evaluates to `[2, 4, 6]`. (Depends on collection-literal codegen rewiring — Phase 2.2 in [`collections.md`](collections.md). If that's still deferred, use `forEachList` instead.)
- [ ] **3.2** Test: `let prefix = ">>"; ["a", "b"] |> forEachList(fn(s) => print(prefix + s))`.
- [ ] **3.3** Test: a record with a function field — `type Validator = { check: (int) -> bool }; let v = Validator { check: fn(n) => n > 0 }; print(v.check(5))`.

## Phase 4 — Negative tests

- [ ] **4.1** `examples/failscompilation/closure_captures_unbound.ospo` — lambda body references a name that's in no enclosing scope; expected error mentions the missing name AND that closures only capture in-scope bindings.
- [ ] **4.2** `examples/failscompilation/closure_arity.ospo` — `let f = fn(x) => x; f(1, 2)` — expected error names both arities.

## Phase 5 — UFCS vs. function-field-call disambiguation

Goal: `obj.fnField(args)` **calls the function-valued field** instead of erroring `function not declared`.
This rule is **already specified** — [0012-Built-InFunctions.md:60](../specs/0012-Built-InFunctions.md#L60):
"If a record has a field named `f`, field access wins; UFCS is the fallback." The codegen is non-conforming.

Today `obj.fnField(args)` is unconditionally built as a `MethodCallExpression` at parse time (purely because
an `LPAREN` follows the field — [builder_calls.go:69,117](../../compiler/internal/ast/builder_calls.go#L69)),
then blindly UFCS-rewritten to `fnField(obj, args)` in codegen
([expression_generation.go:1724](../../compiler/internal/codegen/expression_generation.go#L1724)) and the
inferer ([type_inference.go:706-727](../../compiler/internal/codegen/type_inference.go#L706-L727)). A helper
`receiverHasField` already exists ([type_inference.go:829-845](../../compiler/internal/codegen/type_inference.go#L829-L845),
via `RecordType.HasField`) but is used **only to improve the error message**, not to change dispatch.

- [ ] **5.1** Failing test: `type B = { op: fn(int) -> int }; let b = B { op: fn(n) => n + 1 }; print(b.op(41))` → `42`.
- [ ] **5.2** Before the UFCS rewrite in BOTH `generateMethodCallExpression`
  ([expression_generation.go:1724](../../compiler/internal/codegen/expression_generation.go#L1724)) and the
  inferer ([type_inference.go:709](../../compiler/internal/codegen/type_inference.go#L709)), check
  `receiverHasField(object, name)`. If true → **field-call**: build a `CallExpression` whose `Function` is a
  `FieldAccessExpression{object, name}` (the closure-pair call path from Phase 2.5), passing only the call's
  own args (do NOT prepend `object`). If false → keep the existing UFCS rewrite. Realises the
  field-access-wins-over-UFCS rule the error message already cites. Depends on Phase 2.

## Out of scope

- Mutable captures (Osprey is immutable-by-default; nothing to do).
- Polymorphic lambdas (`fn<T>(x: T) => x`) — lambdas use the surrounding HM scheme; explicit type params on lambdas are a separate workstream.
- Recursive lambdas (`let rec` — would need a fixpoint operator; not in scope here).

## TODO checklist

### Phase 1 — No-capture lambda ✅ LANDED
- [x] No-capture lambdas work as `let`, arg, and record field (re-confirmed this session)

### Phase 2 — Capture
- [ ] 2.1 Failing test for `makeAdder`, `wrap` (curried middleware), capture-outer-`let`
- [ ] 2.2 Closure-pair fat pointer `{ i8* fnptr, i8* env }` in `getLLVMFunctionType`
- [ ] 2.3 Reuse + extend `captureVariablesInExpression` (exclude params; recurse into Lambda/Match/FieldAccess)
- [ ] 2.4 Heap env (`malloc`) + hidden `i8* %__env` param + entry reload in `generateLambdaExpression`
- [ ] 2.5 Call-site fnptr+env extraction; preserve direct-call fast path
- [ ] 2.6 Memory-model doc; leave effect-handler ABI untouched
- [ ] 2.7 Function-typed record fields hold closure pairs (capturing lambda in a field works)
- [ ] 2.8 `wrap`/`adder`/curried-middleware tests pass
- [ ] 2.9 Cleanup: route spawn through shared env; delete `spillNonConstantCaptures`

### Phase 3 — Higher-order
- [ ] 3.1 `map`/`forEachList` with lambda
- [ ] 3.2 Lambda as `forEachList` argument capturing outer string
- [ ] 3.3 Function-typed record field

### Phase 4 — Negative tests
- [ ] 4.1 `closure_captures_unbound.ospo`
- [ ] 4.2 `closure_arity.ospo`

### Phase 5 — UFCS vs. field-call disambiguation
- [ ] 5.1 Failing test `b.op(41)` → calls the field, prints `42`
- [ ] 5.2 `receiverHasField` gate before UFCS rewrite (codegen + inferer)

### Acceptance
- [ ] All `examples/tested/` examples that use `fn(x) => ...` syntax compile and run.
- [ ] [`production-primitives.md`](production-primitives.md) JSON-parser canary uses a closure in at least one place (parser-combinator style).
