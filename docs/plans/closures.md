# Plan: Closures + Lambda Codegen

Spec: [`0004-TypeSystem.md` — Closures](../specs/0004-TypeSystem.md#closures--type-fn-closure) ([TYPE-FN-CLOSURE]).

Parent: [`production-primitives.md`](production-primitives.md).

## Problem

Lambdas don't work at codegen. The simplest case fails:

```osprey
let f = fn(x: int) => x + 1
print(f(10))
```

```
Compilation failed
Execution failed: failed to generate LLVM IR: line 2:22: undefined variable: x
```

The parameter `x` is undefined inside the lambda body. Both the no-capture case (`fn(x) => x + 1`) and the capture case (`fn makeAdder(n) = fn(x) => x + n`) fail with the same error. The spec at [0004-TypeSystem.md:102-104](../specs/0004-TypeSystem.md#L102-L104) advertises both as working.

Type inference for lambdas DOES bind parameters ([type_inference.go:1570-1574](../../compiler/internal/codegen/type_inference.go#L1570-L1574)) and that pass succeeds. The failure is in codegen.

## Why it's broken (suspected)

[`generateLambdaExpression`](../../compiler/internal/codegen/fiber_generation.go#L371) at fiber_generation.go:371 *appears* to bind parameters into `g.variables` (line 408–412):

```go
for i, param := range lambda.Parameters {
    if i < len(lambdaFunc.Params) {
        g.variables[param.Name] = lambdaFunc.Params[i]
    }
}
```

But the error fires from [expression_generation.go:702](../../compiler/internal/codegen/expression_generation.go#L702) (the codegen-side identifier resolution), which means by the time the body's `x + 1` is generated, `x` is not in `g.variables`. Three plausible causes — the implementation agent must determine which before coding:

1. The `LambdaExpression` AST node never reaches `generateLambdaExpression` because a wrapper dispatcher (e.g. `generateChannelOrUnsupportedExpression` at expression_generation.go:106) takes precedence and routes elsewhere. **Verify**: add a `log.Printf` at the entry of `generateLambdaExpression` and re-run the probe.
2. The lambda is generated *eagerly* at let-binding time but then re-generated when called from `print(f(10))`, and the second generation doesn't see the params. **Verify**: count how many functions named `lambda_*` get added to the module.
3. The param binding is overwritten by the `maps.Copy(savedVars, g.variables)` at fiber_generation.go:404–405, which copies *from `g.variables` into `savedVars`* (a save) but the restore at line 443 (`g.variables = savedVars`) wipes the param bindings before the body has finished generating. **Verify**: trace the order of these operations relative to body codegen.

The "Variable doesn't exist anywhere" branch the error came from is at expression_generation.go:701-702 — read that whole function (`generateIdentifier`) to understand the lookup order.

## Phase 1 — Make the no-capture case work

Goal: `let f = fn(x: int) => x + 1; print(f(10))` prints `11`.

- [ ] **1.1** Add a failing test in `tests/integration/` that compiles and runs the snippet above; assert stdout is `11\n`. (This test must be checked in BEFORE the fix.)
- [ ] **1.2** Diagnose the actual cause (one of the three above, or another). Document the finding inline in `fiber_generation.go` above the function.
- [ ] **1.3** Fix the param-binding order so `g.variables[param.Name]` is live when the body codegen runs and stays live until the body finishes.
- [ ] **1.4** Rerun the test from 1.1; it must pass.
- [ ] **1.5** Add a second test with a `string`-typed param: `let g = fn(s: string) => s + "!"; print(g("hi"))`. The current code hardcodes `types.I64` at fiber_generation.go:379 — generalise to use `g.getLLVMType(...)` based on the param's inferred type.
- [ ] **1.6** Add a third test with two params: `let h = fn(a: int, b: int) => a + b; print(h(3, 4))`.

## Phase 2 — Capture

Goal: `fn makeAdder(n: int) -> (int) -> int = fn(x: int) => x + n` works, and the returned closure remains callable after `makeAdder` returns.

- [ ] **2.1** Add a failing test for `makeAdder(5)(3) == 8` and `let add5 = makeAdder(5); add5(3) == 8`.
- [ ] **2.2** Decide capture strategy. Two viable approaches:
  - **(a) Closure-converted struct.** Walk the lambda body before codegen; collect every free identifier; allocate a struct `{ fn_ptr, captured0, captured1, ... }`; lower call sites to `(*closure->fn_ptr)(closure, args...)`. This is the standard ML/Haskell approach.
  - **(b) Heap-allocated environment per call.** Materialise an environment record at the point of lambda construction; pass it as an implicit first arg. Simpler but performs worse for capture-free lambdas.
  - **Recommendation**: (a). Capture-free lambdas trivially lower to plain function pointers (no struct), so they pay no cost and stay zero-overhead.
- [ ] **2.3** Implement free-variable analysis. A new file `internal/codegen/closure_conversion.go` (under the 500-LOC cap) holds it. Input: `*ast.LambdaExpression`. Output: `[]FreeVar { Name string, Type Type }`. Reuses the existing `typeInferer.env`.
- [ ] **2.4** Implement struct generation. The closure value at runtime is `{ i8* fn_ptr, [N x i8] captured_storage }`. Lambda function gets an extra leading `i8* env` parameter; body accesses captured vars via GEP into `env`.
- [ ] **2.5** Rewrite call sites. When calling a value of function type that was bound from a lambda, prepend the env pointer. Direct calls to top-level `fn` names continue to compile as before (no env).
- [ ] **2.6** Heap-allocate the closure struct via `malloc` (osprey already links it; see `runtime/system_runtime.c`). Reference-count or rely on the existing memory model; document in the file header which.
- [ ] **2.7** Tests from 2.1 must pass.

## Phase 3 — Higher-order use sites

Goal: closures work as arguments to `map`, `filter`, `fold`, `forEach`, and as record fields.

- [ ] **3.1** Test: `[1, 2, 3] |> map(fn(x) => x * 2)` evaluates to `[2, 4, 6]`. (Depends on collection-literal codegen rewiring — Phase 2.2 in [`collections.md`](collections.md). If that's still deferred, use `forEachList` instead.)
- [ ] **3.2** Test: `let prefix = ">>"; ["a", "b"] |> forEachList(fn(s) => print(prefix + s))`.
- [ ] **3.3** Test: a record with a function field — `type Validator = { check: (int) -> bool }; let v = Validator { check: fn(n) => n > 0 }; print(v.check(5))`.

## Phase 4 — Negative tests

- [ ] **4.1** `examples/failscompilation/closure_captures_unbound.ospo` — lambda body references a name that's in no enclosing scope; expected error mentions the missing name AND that closures only capture in-scope bindings.
- [ ] **4.2** `examples/failscompilation/closure_arity.ospo` — `let f = fn(x) => x; f(1, 2)` — expected error names both arities.

## Out of scope

- Mutable captures (Osprey is immutable-by-default; nothing to do).
- Polymorphic lambdas (`fn<T>(x: T) => x`) — lambdas use the surrounding HM scheme; explicit type params on lambdas are a separate workstream.
- Recursive lambdas (`let rec` — would need a fixpoint operator; not in scope here).

## TODO checklist

### Phase 1 — No-capture lambda
- [ ] 1.1 Failing test for `let f = fn(x: int) => x + 1; print(f(10))`
- [ ] 1.2 Diagnose the actual cause; document in `fiber_generation.go`
- [ ] 1.3 Fix param binding order
- [ ] 1.4 Test from 1.1 passes
- [ ] 1.5 Param-type generalisation (no more hardcoded `types.I64` at fiber_generation.go:379)
- [ ] 1.6 Two-param lambda test

### Phase 2 — Capture
- [ ] 2.1 Failing test for `makeAdder`
- [ ] 2.2 Decide capture strategy (recommend closure-converted struct)
- [ ] 2.3 `internal/codegen/closure_conversion.go` — free-variable analysis
- [ ] 2.4 Closure struct codegen + extra env param
- [ ] 2.5 Call-site rewrite
- [ ] 2.6 Heap allocation + memory model decision
- [ ] 2.7 Tests from 2.1 pass

### Phase 3 — Higher-order
- [ ] 3.1 `map`/`forEachList` with lambda
- [ ] 3.2 Lambda as `forEachList` argument capturing outer string
- [ ] 3.3 Function-typed record field

### Phase 4 — Negative tests
- [ ] 4.1 `closure_captures_unbound.ospo`
- [ ] 4.2 `closure_arity.ospo`

### Acceptance
- [ ] All `examples/tested/` examples that use `fn(x) => ...` syntax compile and run.
- [ ] [`production-primitives.md`](production-primitives.md) JSON-parser canary uses a closure in at least one place (parser-combinator style).
