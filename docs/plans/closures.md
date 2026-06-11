# Plan: Closures + Lambda Codegen

Spec: [`0004-TypeSystem.md` — Closures](../specs/0004-TypeSystem.md#closures--type-fn-closure) ([TYPE-FN-CLOSURE]).

Parent: [`production-primitives.md`](production-primitives.md).
Downstream: gating dependency for [`backend-framework.md`](backend-framework.md) (composable
`Middleware = fn(Handler) -> Handler`).

## Problem (UPDATED — re-probed against the Rust compiler)

**Phase 1 is done, and the let-bound capture case now works too.** Probes against
`target/release/osprey`:

- `let f = fn(x: int) => x + 1; print(f(10))` → `11` ✅
- `let n = 10; let addN = fn(x) => x + n; addN(5)` → `15` ✅ (let-bound capturing lambda)
- `let wrap = fn(h) => fn(b) => h(b)` compiles ✅ (let-bound, inlined)
- function value passed as an arg, and a non-capturing lambda stored in a record field → work ✅

**What is still missing: escaping closures** — a lambda *returned or stored as a first-class value*
that captures its maker's state. `fn makeAdder(n: int) = fn(x) => x + n` has no expressible return
type (the checker rejects calling the resulting `any`), and a *capturing* lambda passed as a function
value loses its environment.

## Why (the Rust lowering today)

Lambdas take two paths in `crates/osprey-codegen`:

1. **Let-bound lambdas inline at the call site**, where the captured names are still in scope — this
   is why `addN(5)` and curried `wrap` work. The lambda never becomes a runtime value.
2. **Lambdas used as *values*** (passed to a function-typed parameter) are **lifted** to a top-level
   function by `lift_lambda` ([`crates/osprey-codegen/src/genfn.rs`](../../crates/osprey-codegen/src/genfn.rs)),
   which is explicit that "a free identifier bound in the enclosing scope is not captured (the backend
   lowers no closures)". A function value is a bare code pointer (`i8*`), never a closure pair `{fn, env}`.

Only the **`spawn`** path threads captures, via a **global-spill mechanism**: `free_idents`
([`crates/osprey-codegen/src/fiber.rs`](../../crates/osprey-codegen/src/fiber.rs)) collects free vars and
`spill_captures` stores each into a per-spawn **module global** the emitted thunk reloads. That's fine for
`spawn`'s fire-once semantics but is **not** a general closure model (breaks re-entrancy, recursion,
multiple instances). General escaping lambdas need a per-instance heap env (Phase 2).

## Phase 1 — Make the no-capture case work ✅ LANDED

Goal (met): `let f = fn(x: int) => x + 1; print(f(10))` prints `11`; string-typed and multi-param
lambdas, lambdas-as-args, and non-capturing lambdas in record fields all compile and run. Keep/extend the
golden coverage when touching the lambda path in Phase 2 so this never regresses.

## Phase 2 — Capture

Goal: `fn makeAdder(n: int) -> (int) -> int = fn(x: int) => x + n` works, and the returned closure remains callable after `makeAdder` returns.

- [ ] **2.1** Add a failing test for `makeAdder(5)(3) == 8` and `let add5 = makeAdder(5); add5(3) == 8`.
- [ ] **2.2** Capture strategy — **closure-pair fat pointer**. A function VALUE becomes a 2-field struct
  `osprey.closure = { i8* fnptr, i8* env }` (define once in the generator). `env` is `null` for
  non-capturing lambdas / top-level functions, so capture-free lambdas pay no allocation. Change the
  function-value spelling (`fn_ptr_spelling` in
  [`crates/osprey-codegen/src/genfn.rs`](../../crates/osprey-codegen/src/genfn.rs)) and the indirect-call
  bitcast that shares it.
- [ ] **2.3** Free-variable analysis — **reuse the existing walker** `free_idents`
  ([`crates/osprey-codegen/src/fiber.rs`](../../crates/osprey-codegen/src/fiber.rs)) rather than a parallel
  implementation (no-duplication rule); it already subtracts locally-bound names and recurses the full
  `Expr` surface. Move it somewhere neutral (e.g. a small `freevars` module) if fiber.rs stops being its
  only consumer.
- [ ] **2.4** In `lift_lambda`: for each captured var, **`malloc` a heap env struct** `{ captured fields… }`
  (NOT a module global — it must outlive the enclosing frame for escaping closures like `adder`), store
  current values via the parent builder, `bitcast` to `i8*`. Add a hidden leading `i8* %__env` param to the
  lifted lambda; at entry, load each captured field into scope before body codegen (mirrors the spawn
  thunk's capture reload in fiber.rs, but per-instance). Return the closure-pair
  `{ bitcast(lambda to i8*), env }`.
- [ ] **2.5** Call sites — in the indirect-call path
  ([`crates/osprey-codegen/src/call.rs`](../../crates/osprey-codegen/src/call.rs)): when the callee is a
  closure pair, extract field 0 (fnptr) + field 1 (env) and call `fnptr(env, args...)`. **Preserve the
  direct-call fast path** for named top-level funcs so the bulk of the codebase is unperturbed. Simplest
  ABI: all function-value signatures take a leading `i8*` env that top-level callees ignore.
- [ ] **2.6** Heap-allocate the env via `malloc` (already linked). Document the memory model in the file
  header. **Do NOT touch the effect-handler ABI**
  ([`crates/osprey-codegen/src/effects.rs`](../../crates/osprey-codegen/src/effects.rs)) — handlers are
  non-capturing; wrap as `{fn, null}` only if a handler is ever used as a first-class value.
- [ ] **2.7** Function-typed record fields (ctor layout in
  [`crates/osprey-codegen/src/builder.rs`](../../crates/osprey-codegen/src/builder.rs) /
  [`aggregate.rs`](../../crates/osprey-codegen/src/aggregate.rs)) become the closure-pair struct — so
  storing a *capturing* lambda in a field works, and `let f = b.op; f(x)` keeps working.
- [ ] **2.8** Tests from 2.1 must pass; add `wrap`/`adder`/curried-middleware cases.
- [ ] **2.9** (cleanup) Route `spawn` through the shared env mechanism and **delete `spill_captures`**
  ([`crates/osprey-codegen/src/fiber.rs`](../../crates/osprey-codegen/src/fiber.rs)) — removes the
  spawn-only global-spill capture mechanism.

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

Today `obj.fnField(args)` is rewritten to `fnField(obj, args)` at parse time — the UFCS sugar lives in the
syntax layer ([`crates/osprey-syntax/src/expr.rs`](../../crates/osprey-syntax/src/expr.rs), "UFCS:
`x.f(a, …)` is sugar for `f(x, a, …)`") — so by the time the checker and codegen see the call, the
field-call interpretation is already gone.

- [ ] **5.1** Failing test: `type B = { op: fn(int) -> int }; let b = B { op: fn(n) => n + 1 }; print(b.op(41))` → `42`.
- [ ] **5.2** Gate the UFCS rewrite on "does the receiver's record type have this field?". Since the
  rewrite happens pre-inference, either (a) defer the decision: parse to a neutral `MethodCall` node and
  let `crates/osprey-types` resolve field-call vs UFCS during inference, or (b) keep the rewrite but have
  the checker re-interpret `f(x, …)` as a field call when `x`'s record type declares a function-typed
  field `f`. Option (a) is cleaner. Field access wins; UFCS is the fallback — the rule the spec already
  states. Depends on Phase 2 (the field's value must be a callable closure pair).

## Out of scope

- Mutable captures (Osprey is immutable-by-default; nothing to do).
- Polymorphic lambdas (`fn<T>(x: T) => x`) — lambdas use the surrounding HM scheme; explicit type params on lambdas are a separate workstream.
- Recursive lambdas (`let rec` — would need a fixpoint operator; not in scope here).

## TODO checklist

### Phase 1 — No-capture lambda ✅ LANDED
- [x] No-capture lambdas work as `let`, arg, and record field (re-confirmed this session)

### Phase 2 — Capture
- [ ] 2.1 Failing test for `makeAdder`, `wrap` (curried middleware), capture-outer-`let`
- [ ] 2.2 Closure-pair fat pointer `{ i8* fnptr, i8* env }` in `fn_ptr_spelling` (genfn.rs)
- [ ] 2.3 Reuse `free_idents` (fiber.rs) for free-variable analysis
- [ ] 2.4 Heap env (`malloc`) + hidden `i8* %__env` param + entry reload in `lift_lambda`
- [ ] 2.5 Call-site fnptr+env extraction (call.rs); preserve direct-call fast path
- [ ] 2.6 Memory-model doc; leave effect-handler ABI untouched
- [ ] 2.7 Function-typed record fields hold closure pairs (capturing lambda in a field works)
- [ ] 2.8 `wrap`/`adder`/curried-middleware tests pass
- [ ] 2.9 Cleanup: route spawn through shared env; delete `spill_captures`

### Phase 3 — Higher-order
- [ ] 3.1 `map`/`forEachList` with lambda
- [ ] 3.2 Lambda as `forEachList` argument capturing outer string
- [ ] 3.3 Function-typed record field

### Phase 4 — Negative tests
- [ ] 4.1 `closure_captures_unbound.ospo`
- [ ] 4.2 `closure_arity.ospo`

### Phase 5 — UFCS vs. field-call disambiguation
- [ ] 5.1 Failing test `b.op(41)` → calls the field, prints `42`
- [ ] 5.2 Field-presence gate before the UFCS rewrite (syntax layer + checker)

### Acceptance
- [ ] All `examples/tested/` examples that use `fn(x) => ...` syntax compile and run.
- [ ] [`production-primitives.md`](production-primitives.md) JSON-parser canary uses a closure in at least one place (parser-combinator style).
