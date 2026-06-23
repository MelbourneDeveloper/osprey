# Plan 0008 — Effect `resume` / Continuations

**Subsystem:** `crates/osprey-syntax`, `crates/osprey-ast`, `crates/osprey-types`,
`crates/osprey-codegen`, `compiler/runtime`
**Status:** Partially implemented (handlers work as value substitution)
**Spec:** [0017-AlgebraicEffects.md](../specs/0017-AlgebraicEffects.md)

## Summary

Effects are real and compile-time safe: declarations, `perform`, `handle … in`,
effect annotations, and unhandled-effect rejection all work. What is missing is
the defining feature of *algebraic* effects — `resume`. Today a handler arm
behaves as a value substitution: it returns a value and the performing
computation does **not** continue from the `perform` site. This is the largest and
highest-risk of the partial features; it is the capstone, not a quick win.

## Evidence

- Spec status: *"Continuation/`resume` semantics inside handlers are not yet
  implemented; current handlers act as value substitutions …"* —
  [0017-AlgebraicEffects.md](../specs/0017-AlgebraicEffects.md) §Status.
- Codegen note: *"example handlers never resume, so an arm is an ordinary
  function"* — [crates/osprey-codegen/src/effects.rs](../../crates/osprey-codegen/src/effects.rs).

## What works today

- Effect declarations, `perform X.op(args)`, `handle X arm… in body`.
- Dynamic handler stack with push/pop/lookup, snapshot/restore across fiber
  boundaries — [compiler/runtime/effects_runtime.c](../../compiler/runtime/effects_runtime.c),
  [crates/osprey-codegen/src/effects.rs](../../crates/osprey-codegen/src/effects.rs).
- Compile-time unhandled-effect checking —
  [crates/osprey-types/src/check.rs](../../crates/osprey-types/src/check.rs).
- Working example —
  [compiler/examples/tested/effects/algebraic_effects_comprehensive.osp](../../compiler/examples/tested/effects/algebraic_effects_comprehensive.osp).

## Where it stops

A handler arm is lowered as an ordinary function returning a value; there is no
captured continuation, so `resume v` (continue the performer with `v`) cannot be
expressed. Multi-shot resume (resuming the same continuation more than once) is
likewise impossible.

## Implementation plan

This needs a continuation mechanism. Recommended phased approach:

1. **Surface syntax + AST.** Add `resume <expr>` to the grammar
   ([tree-sitter-osprey/](../../tree-sitter-osprey/)), the AST
   ([crates/osprey-ast/src/lib.rs](../../crates/osprey-ast/src/lib.rs)), and the
   syntax lowering ([crates/osprey-syntax](../../crates/osprey-syntax)).
2. **Type the continuation.** In [crates/osprey-types](../../crates/osprey-types),
   bind the handler operation's parameters from the effect's `OpType`, and type
   `resume` as accepting the operation's result type and producing the handled
   computation's type. (Today handler params are fresh, unconnected vars —
   tighten that first.)
3. **Choose the runtime model.** Single-shot delimited continuations first:
   - **Option A — CPS transform** of effectful code paths (no stack copying;
     larger codegen change).
   - **Option B — stackful capture** (`ucontext`/saved stack segment) keyed off
     the handler stack frame; smaller codegen change, more runtime machinery.
   Recommend prototyping **B** for single-shot, since the handler stack and
   snapshot/restore plumbing already exist in `effects_runtime.c`.
4. **Implement `__osprey_handler_resume`** to reactivate the captured
   continuation and deliver the resume value back to the `perform` site.
5. **Codegen handler arms** to capture the continuation, bind `resume`, and emit
   resume calls.
6. **Defer multi-shot resume** to a follow-up; gate it behind a clear "not yet"
   diagnostic so single-shot ships first.

## Testing

- Extend
  [algebraic_effects_comprehensive.osp](../../compiler/examples/tested/effects/algebraic_effects_comprehensive.osp)
  with a state/generator-style handler that resumes (e.g. a counter or a
  `yield`-like producer) and asserts the performer continues; refresh
  `.expectedoutput`.
- `failscompilation` case: multi-shot resume rejected with a clear message (until
  implemented).

## Risks / considerations

- Highest-risk item here: continuations interact with memory management
  (currently none — [0018](../specs/0018-MemoryManagement.md)); a captured stack
  segment that is never freed compounds existing leaks. Note the dependency.
- Interaction with fibers: the handler snapshot/restore path must compose with
  captured continuations.
- Land single-shot first; keep the existing value-substitution behaviour working
  for handlers that never resume.

## TODO

- [ ] Add `resume <expr>` to grammar, AST, and syntax lowering.
- [ ] Bind handler-arm params from the effect `OpType`; type `resume`.
- [ ] Prototype single-shot delimited continuations (stackful capture on the
      existing handler stack).
- [ ] Implement `__osprey_handler_resume` in `effects_runtime.c`.
- [ ] Codegen handler arms to capture the continuation + emit resume.
- [ ] Reject multi-shot resume with a clear diagnostic (follow-up to implement).
- [ ] Extend the effects example with a resuming handler; refresh `.expectedoutput`.
- [ ] Update 0017 §Status once single-shot resume lands.
- [ ] `make ci` green.
