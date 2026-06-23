# Plan 0006 — Cooperative `yield()`

**Subsystem:** `compiler/runtime` (C fiber runtime), with codegen already in place
**Status:** Partially implemented (front end done; runtime is a no-op stub)
**Spec:** [0011-LightweightFibersAndConcurrency.md](../specs/0011-LightweightFibersAndConcurrency.md)

## Summary

`yield` parses, type-checks, and lowers — but at runtime it does nothing. The
codegen emits `yield e` as the identity of `e`, and the C runtime's `fiber_yield`
returns its argument without switching context. So a fiber cannot cooperatively
hand the CPU back to the scheduler.

## What works today

- Parser accepts `yield` / `yield(value)`.
- Codegen lowers `yield` as identity — `gen_yield` in
  [crates/osprey-codegen/src/fiber.rs](../../crates/osprey-codegen/src/fiber.rs)
  (`Some(inner) => gen_expr(cg, inner)`).
- Spawn / await / channels / `fiberDone` are fully working —
  [compiler/runtime/fiber_runtime.c](../../compiler/runtime/fiber_runtime.c),
  [compiler/examples/tested/fiber/fiber_showcase.osp](../../compiler/examples/tested/fiber/fiber_showcase.osp).

## Where it bails

```c
// compiler/runtime/fiber_runtime.c:229
// TODO: Implement proper fiber yielding with context switching
int64_t fiber_yield(int64_t value) {
    return value;        // no-op pass-through
}
```

The spec intends `yield` to donate the CPU to the scheduler and resume the fiber
when it is next scheduled (cooperative multitasking) —
[0011 §yield](../specs/0011-LightweightFibersAndConcurrency.md).

## Implementation plan

The current runtime backs fibers with **pthreads**. Two routes:

- **A (recommended, smaller):** Implement `yield` as a cooperative scheduling
  point on top of the existing threaded model — a `sched_yield()`-style hand-off
  plus a fiber run-queue and a condition variable, so the yielding fiber blocks
  until the scheduler re-selects it. Deterministic mode (already present —
  `fiber_set_deterministic_mode`) must remain reproducible.
- **B (larger):** Introduce stackful context switching (`ucontext`/custom asm) and
  a real M:N scheduler. Higher fidelity, much larger surface; defer unless A
  proves insufficient.

Steps for A:

1. **Add a run-queue + scheduler state** to the fiber runtime (ready/blocked
   lists guarded by the existing mutex/condvar).
2. **Implement `fiber_yield`** to: mark the current fiber ready, signal the
   scheduler, and wait until re-selected before returning `value`.
3. **Integrate with `spawn`/`await`/channels** so a fiber blocked on a channel
   yields rather than busy-waits (improves the existing blocking paths too).
4. **Preserve deterministic mode** ordering for the differential harness.
5. **Carry the yielded value** through if/when `yield` gains generator-style
   semantics; for now `yield` returns control and forwards `value` unchanged.

## Testing

- Extend [fiber_showcase.osp](../../compiler/examples/tested/fiber/fiber_showcase.osp)
  with two fibers that interleave via `yield`, asserting a deterministic
  interleaving order under deterministic mode; refresh `.expectedoutput`.

## Risks / considerations

- Concurrency correctness: the run-queue and condvar handshake must avoid lost
  wakeups and deadlocks.
- Determinism: the differential harness compares output byte-for-byte, so the
  scheduler must be deterministic under the test mode.
- Keep `fiber_yield`'s signature/ABI stable so codegen needs no change.

## TODO

- [ ] Add a fiber run-queue + scheduler state to the C runtime.
- [ ] Implement `fiber_yield` as a real cooperative hand-off (replace the stub at
      `fiber_runtime.c:229`).
- [ ] Make channel/await blocking yield to the scheduler instead of busy-wait.
- [ ] Preserve deterministic-mode ordering.
- [ ] Extend `fiber_showcase.osp` with an interleaving test; refresh `.expectedoutput`.
- [ ] `make ci` green (watch for flakiness under concurrency).
