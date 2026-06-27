# Plan 0008 — Effect `resume` / Continuations

**Subsystem:** `crates/osprey-syntax`, `crates/osprey-ast`, `crates/osprey-types`,
`crates/osprey-codegen`, `compiler/runtime`
**Status:** Single-shot deep `resume` landing — thread-as-continuation (Option B)
**Spec:** [0017-AlgebraicEffects.md](../specs/0017-AlgebraicEffects.md)

## Summary

Effects are real and compile-time safe: declarations, `perform`, `handle … in`,
effect annotations, and unhandled-effect rejection all work. A handler arm's
value now becomes the `perform`'s result and the performer continues past the
`perform` site — the common **single-shot tail-resume** — and handlers may own
mutable state (`[EFFECTS-HANDLER-STATE]`, see below), so the `State` effect is
fully usable today. What remains is an *explicit* `resume` expression: capturing
the continuation as a value so an arm can run code *after* resuming, resume in a
non-tail position, or resume more than once (multi-shot). That is the remaining
capstone.

## Update — handler-owned state landed

`[EFFECTS-HANDLER-STATE]` Handler arms can read and write a `mut` captured from
the enclosing scope; such a `mut` is promoted to a shared heap cell and the
handler stack carries a per-region environment (`__osprey_handler_push` gained an
`env` pointer; `__osprey_handler_lookup_env` resolves it). This delivers the
canonical State-effect pattern (handler owns the cell, effectful code stays
pure) without general continuations, because value-substitution *is* tail-resume.
Implemented in [crates/osprey-codegen/src/effects.rs](../../crates/osprey-codegen/src/effects.rs)
(`capture_list`/`build_env`/`reload_env`), the cell read/write in
[expr.rs](../../crates/osprey-codegen/src/expr.rs) and
[lower.rs](../../crates/osprey-codegen/src/lower.rs), and
[compiler/runtime/effects_runtime.c](../../compiler/runtime/effects_runtime.c).
Reference app: `examples/tested/effects/http_state_levels.osp`. The
remaining `resume` work below is unchanged.

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
  [examples/tested/effects/algebraic_effects_comprehensive.osp](../../examples/tested/effects/algebraic_effects_comprehensive.osp).

## Where it stops

A handler arm is lowered as an ordinary function returning a value; there is no
captured continuation, so `resume v` (continue the performer with `v`) cannot be
expressed. Multi-shot resume (resuming the same continuation more than once) is
likewise impossible.

## Chosen design — thread-as-continuation (Option B)

The runtime is already thread-based (fibers are pthreads,
[fiber_runtime.c](../../compiler/runtime/fiber_runtime.c)) and already
snapshots/restores the handler stack across threads
([effects_runtime.c](../../compiler/runtime/effects_runtime.c)). There is no
`ucontext`/`setjmp` in the tree, so a suspended **thread** is the continuation —
no stack-segment copying, no CPS pass. This also makes single-shot fall out for
free: a live pthread stack cannot be cloned, so multi-shot is naturally excluded
(and rejected with a diagnostic).

### Static gate

A `handle E arm… in body` is a **resuming region** iff any arm body contains a
`resume`. Resuming regions emit the coroutine path below; every other region
keeps the existing zero-overhead function-call path (handler = function, `ret` =
tail-resume). Detected by an AST walk over arm bodies in codegen.

### Runtime ABI (`compiler/runtime/effects_runtime.c`)

A per-region `Coro` control block carries: the captured user `env`, a turn flag +
mutex/cond pair, an operation-id + argument buffer (body→host), a `resume_value`
(host→body), the body's final result, and `done`/`abort` flags.

- `Coro *__osprey_coro_new(void *env)` — allocate.
- `void __osprey_coro_start(Coro*, i64 (*body)(void*), HandlerSnapshot*)` — spawn
  the body thread with the inherited handler stack, then block the host until the
  body first suspends or completes.
- `i64 __osprey_coro_suspend(Coro*, i64 op_id, i64 *args)` — **body side**, called
  by each arm's suspend-trampoline at a `perform`: publishes `(op_id,args)`, hands
  control to the host, blocks until resumed, returns `resume_value`. If the host
  set `abort`, it `pthread_exit`s (single-shot teardown).
- `i64 __osprey_coro_resume(Coro*, i64 v)` — **host side**, lowering of `resume`:
  delivers `v`, runs the body until its next suspend or completion, blocks the
  host meanwhile; returns `done ? body_result : <sentinel: body performed again>`.
- accessors `__osprey_coro_done/op/arg/result` and `__osprey_coro_abort` + free.

### Dispatch model (codegen-emitted, host thread)

The host side is ordinary host-thread call recursion; the body thread runs
straight-line and only suspends:

```
region(env):
  coro = coro_new(env); push suspend-trampolines(env=coro); snapshot = handler_snapshot()
  coro_start(coro, body_thunk, snapshot)
  return drive(coro)

drive(coro):                 // also re-entered after each resume that performed again
  if coro_done(coro): return coro_result(coro)
  return dispatch_arm(coro, coro_op(coro))     // an arm finishing IS the region answer

dispatch_arm(coro, op):      // switch op → __arm_E_op(env, args…); arm body may call resume
resume(v)  ⇒  r = coro_resume(coro, v); if !coro_done(coro) then drive(coro) else r
```

Tail-resume arms bottom out when the body completes (answer = `body_result`),
matching today's semantics; non-tail arms run their post-`resume` code as the
recursion unwinds (LIFO), the behaviour value-substitution can't express. An arm
that never resumes returns directly → host sets `abort`, joins the body, frees the
`Coro`.

### Phases

1. **Surface syntax + AST.** `resume "(" expr? ")"` in the grammar, AST
   `Expr::Resume(Option<Box<Expr>>)`, syntax lowering.
2. **Types.** Bind handler-arm params from the effect's `OpType`; type `resume`'s
   argument against the operation result; reject `resume` outside a handler arm.
3. **Runtime.** The `Coro` ABI above in `effects_runtime.c`.
4. **Codegen.** Static gate; suspend-trampolines; `body_thunk`; the host
   `drive`/`dispatch_arm`; lower `resume`.
5. **Multi-shot** rejected with a clear diagnostic.

## Testing

- Extend
  [algebraic_effects_comprehensive.osp](../../examples/tested/effects/algebraic_effects_comprehensive.osp)
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
