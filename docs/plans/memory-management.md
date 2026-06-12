# Plan: Pluggable Memory Management + Static Mode

Implements [`0018-MemoryManagement.md`](../specs/0018-MemoryManagement.md).

## Why

Today the compiler leaks everything: codegen emits `malloc` for closure
cells, records/unions, Result blocks, and list literals
([`closure.rs`](../../crates/osprey-codegen/src/closure.rs),
[`aggregate.rs`](../../crates/osprey-codegen/src/aggregate.rs),
[`result.rs`](../../crates/osprey-codegen/src/result.rs)), and the C runtime
allocates strings, HAMT map nodes, and list nodes
([`string_runtime.c`](../../compiler/runtime/string_runtime.c),
[`map_runtime_hamt.c`](../../compiler/runtime/map_runtime_hamt.c),
[`list_runtime.c`](../../compiler/runtime/list_runtime.c)) — with no frees
anywhere on the language-value path. Fine for proving out semantics; fatal
for production.

The spec deliberately makes reclamation invisible ([MEM-OPAQUE]), so we can
fix this with **one analysis** and get three products from it:

1. **Default mode** — ownership inference statically frees every value whose
   last use is provable; the *shared residue* (structural sharing, aliased
   escaping closures — [MEM-OWNERSHIP-SHARED]) carries a non-atomic refcount.
   Perceus/Lobster-class: most RC traffic elided at compile time.
2. **`--static-memory` mode** — the SAME analysis, but where it would insert
   a refcount it errors instead ([MEM-STATIC-MODE]). Accepted programs have
   zero runtime memory management, Rust-class, without a borrow checker the
   dev fights — the subset is enforced, not annotated.
3. **Backends** — ARC is the default; a tracing GC ships out of the box as
   the conformance oracle ([MEM-BACKENDS]); custom managers link against the
   same C interface at their supplier's risk ([MEM-BACKENDS-CUSTOM]).
   Goldens must match byte-for-byte and leak-free under every shipped
   backend.

Design pillars already in our favor: immutability keeps the heap acyclic
([MEM-ACYCLIC]) so RC is complete without a cycle collector, and fiber
isolation ([MEM-FIBER-ISOLATION]) makes every refcount non-atomic.

## Phases

### Phase 0 — Guardrails (small, do first)

Leak instrumentation in the differential harness: run each golden under
macOS `leaks` (CI: LeakSanitizer), record a per-example baseline count.
Not gating yet — the baseline is the ratchet, same discipline as
`FC_EXPECTED_ESCAPES`. Audit specs/examples for anything that could observe
reclamation timing (expected: nothing — the language has no finalizer
surface).

### Phase 1 — Refcount backbone (default mode becomes leak-free)

- One heap header for all language values (closure cells, strings, list
  nodes, HAMT nodes, records/unions, Result blocks): non-atomic count +
  payload-kind tag so `release` can walk children. C runtime gains
  `osprey_alloc`/`osprey_retain`/`osprey_release` — this IS the backend
  interface custom managers link against ([MEM-BACKENDS-CUSTOM]); existing
  per-type allocators route through it.
- Codegen emits retain/release naively but correctly: retain on bind/escape,
  release at scope exit and last use. Correctness first; elision is Phase 2.
- `spawn` capture and channel `send` deep-copy (or move when the analysis
  already proves uniqueness) per [MEM-FIBER-ISOLATION] — today's shared
  closure-cell capture across fibers is replaced, killing the cross-fiber
  co-ownership case entirely.
- Flip the Phase 0 ratchet to a gate: PASS requires zero leaked language
  values on every golden.

### Phase 2 — Ownership inference (Perceus-style elision)

- Infer uniqueness/borrowing over the lowered program; delete
  retain/release pairs where ownership transfers linearly (the common case).
- Reuse analysis (functional-but-in-place): a uniquely-owned value freed and
  reallocated at the same shape mutates in place.
- Benchmark gate: an RC-ops counter in debug builds; goldens record elision
  rates so regressions in the analysis are visible, not vibes.

### Phase 3 — `--static-memory` mode

- Same analysis; at each would-be refcount site emit a compile error naming
  the shared value and the conflicting owners (e.g. "`xs` is shared between
  `ys` (line 4) and the closure at line 9; no static last owner").
- Harness: a `static-ok/` subset of goldens must compile in static mode and
  produce byte-identical output to default mode; `failscompilation/`-style
  cases pin the rejections.
- Emitted IR for static-mode programs must contain zero
  `osprey_retain`/`osprey_release` calls — assert by grep in the harness.

### Phase 4 — Second backend + arenas

- Tracing GC backend behind the Phase 1 interface; build flag selects
  backend; full golden run under both ([MEM-BACKENDS]).
- Custom-manager link path: document the C interface, stamp custom builds
  in `--version` ([MEM-BACKENDS-CUSTOM]).
- Effect-scoped arenas: a handler that allocates a region and frees it
  wholesale at scope exit (per-request arenas for the HTTP framework).
  Policy via effects, never per-allocation dispatch.

## Open questions

- **Runtime-internal allocations** (HTTP buffers, terminal state) are owned
  and freed by the C runtime itself — separate leak-audit workstream, not
  governed by [MEM-OWNERSHIP].
- **Copy vs move on `send`**: moving requires the analysis to prove the
  sender holds the unique reference; copy is always safe. Start with copy,
  let Phase 2 upgrade provable cases to moves.
- **Collections in static mode**: barred in v1 per
  [MEM-STATIC-MODE-BARRED]. Revisit (runtime-internal counts? unique-owned
  array type?) if static-mode users need them.

## TODO

- [ ] Phase 0: leak baseline in differential harness (`leaks`/LSan ratchet)
- [ ] Phase 0: spec/example audit for reclamation-timing observability
- [ ] Phase 1: heap header + `osprey_retain`/`osprey_release` in C runtime
- [ ] Phase 1: naive retain/release emission in codegen
- [ ] Phase 1: deep-copy (move-ready) `spawn`/`send` per [MEM-FIBER-ISOLATION]
- [ ] Phase 1: flip leak ratchet to hard gate (zero leaks on all goldens)
- [ ] Phase 2: ownership inference + RC elision + reuse analysis
- [ ] Phase 2: RC-ops counter + elision-rate tracking in harness
- [ ] Phase 3: `--static-memory` flag, sharing diagnostics
- [ ] Phase 3: `static-ok/` goldens + static-mode rejection pins + zero-RC grep
- [ ] Phase 4: tracing-GC backend + dual-backend golden run
- [ ] Phase 4: custom-manager interface docs + `--version` stamp
- [ ] Phase 4: effect-scoped arena handler
