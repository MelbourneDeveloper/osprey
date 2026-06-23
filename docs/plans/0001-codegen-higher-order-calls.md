# Plan 0001 — Higher-Order / Indirect Calls in Codegen

**Subsystem:** `crates/osprey-codegen`
**Status:** Partially implemented
**Spec:** [0005-FunctionCalls.md](../specs/0005-FunctionCalls.md), [0010-LoopConstructsAndFunctionalIterators.md](../specs/0010-LoopConstructsAndFunctionalIterators.md)

## Summary

A call whose callee is a plain identifier (a top-level function, a builtin, or a
function-typed local bound through a closure cell) lowers correctly. A call whose
callee is an **arbitrary expression** — a call result, a record field, or a
chained application — fails loudly. The machinery to call *through a closure
cell* already exists; it is simply not reachable from a non-identifier callee.

## What works today

- Direct calls to top-level functions and builtins — [crates/osprey-codegen/src/expr.rs](../../crates/osprey-codegen/src/expr.rs).
- Calls to a function-typed **local** through its closure cell, via `try_indirect()`.
- Direct lambda application (inline beta-reduction).
- The generic "call through a cell given a signature" path already exists in `closure::cell_call()` — [crates/osprey-codegen/src/closure.rs](../../crates/osprey-codegen/src/closure.rs).

## Where it bails

```rust
// crates/osprey-codegen/src/expr.rs:374
let Expr::Identifier(ident) = function else {
    return Err(CodegenError::unsupported("indirect / higher-order call"));
};
```

Consequently these all fail and must be worked around by binding to a `let` first:

- `makeAdder(5)(3)` — calling a call result.
- `config.processor(x)` — calling a closure stored in a record field.
- `f()()` — chained application.

The same limitation cascades into iterator combinators: `callback_of()` rejects
any callback that is not a name or a literal lambda — [crates/osprey-codegen/src/iter.rs](../../crates/osprey-codegen/src/iter.rs). Finishing this plan unblocks computed/field-access iterator callbacks too.

## Implementation plan

1. **Generalize the callee path.** In the call lowering in
   [expr.rs](../../crates/osprey-codegen/src/expr.rs), when the callee is not an
   `Expr::Identifier`, evaluate it with `gen_expr` to a `Value` holding a closure
   handle instead of erroring.
2. **Recover the signature.** Read the inferred function type of the callee
   expression from the type table (the same inference the let-bound path already
   relies on) and convert it to a `FnSig` via `Codegen::fn_value_sig`.
3. **Dispatch through the existing cell-call.** Feed the handle + `FnSig` into
   `closure::cell_call()` — no new ABI needed.
4. **Route `callback_of()` through the same path.** Replace the bail in
   [iter.rs](../../crates/osprey-codegen/src/iter.rs) with a fallback that
   evaluates the callback expression to a cell and calls it, so `map`/`filter`/
   `forEach`/`fold` accept computed and field-access callbacks.
5. **Keep failing loudly** where the inferred type is still generic — that is the
   separate concern handled by [Plan 0002](0002-codegen-generic-function-values.md).

## Testing

- Extend an existing example under
  [compiler/examples/tested/basics/](../../compiler/examples/tested/basics/)
  (e.g. the function-composition example) with: a curried `makeAdder(5)(3)`, a
  record field holding a lambda then calling it, and an iterator whose callback is
  a field access. Update the matching `.expectedoutput`.
- Add a `failscompilation` case only if a genuinely uncallable expression should
  still be rejected.

## Risks / considerations

- Inferring the signature of an arbitrary call result is the crux; if the type
  table does not yet record a concrete function type at that node, surface a clear
  diagnostic rather than guessing (CLAUDE.md: no placeholders).
- Watch for `Result` auto-unwrap interactions when the callee expression is itself
  `Result`-typed.
- The concurrent `freevars.rs` refactor centralizes capture analysis; reuse it
  rather than re-deriving free variables for the evaluated callee.

## TODO

- [ ] Replace the `Expr::Identifier` guard in `expr.rs:374` with a general
      callee-expression path.
- [ ] Resolve the callee's inferred `FnSig` from the type table.
- [ ] Dispatch non-identifier callees through `closure::cell_call()`.
- [ ] Generalize `callback_of()` in `iter.rs` to accept computed/field callbacks.
- [ ] Emit a precise diagnostic when the callee type is unknown/non-function.
- [ ] Extend a `tested/basics` example with curried, field-stored, and chained
      calls; refresh `.expectedoutput`.
- [ ] `make ci` green.
