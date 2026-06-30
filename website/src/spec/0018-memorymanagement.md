---
layout: page
title: "Memory Management"
description: "Osprey Language Specification: Memory Management"
date: 2026-06-30
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0018-memorymanagement/"
---

# Memory Management

Osprey programs do not manage memory. Reclamation is a property of the
*implementation*, never of the language: these rules define semantics under
which a conforming implementation may reclaim memory with reference counting,
a tracing collector, fully static frees, or any mix — with no observable
difference to any program. The developer's only obligation is the one every
garbage-collected language already imposes: don't keep references to values
you no longer need.

## Status

Partially implemented — the *boundary* exists and a first *reclaiming* backend
ships (tracing GC, opt-in via `--memory=gc`); the ARC default is next
([implementation plan 0011](https://github.com/Nimblesite/osprey/blob/main/docs/plans/0011-arc-gc-implementation.md)).

- **Swappable backend boundary [MEM-BACKENDS]: done.** All codegen heap
  allocation funnels through a single `@osp_alloc` hook (osprey-codegen
  `builder.rs::heap_alloc` / `OSP_ALLOC_DECL`); the emitted IR names no
  allocator, so a manager is chosen at link time. The default backend
  (`compiler/runtime/memory_runtime.c`) is a `malloc` passthrough that never
  frees during a run — sound because reclamation is unobservable [MEM-OPAQUE].
- **Static reclamation of non-escaping values: done, by the optimizer.** The
  `@osp_alloc` declaration carries allocator attributes, so at `-O2` LLVM proves
  provably-dead allocations (the common case — per-operation `Result` blocks,
  temporaries) non-escaping and removes them entirely. This realises the
  [MEM-OWNERSHIP] "free at last use, statically" ideal for everything whose
  lifetime LLVM can see.
- **Reclaiming *escaping* values: tracing GC ships, ARC pending.** Values that
  genuinely outlive their allocation site (e.g. nodes of a built-and-held tree)
  still leak under the *default* backend, but the opt-in tracing collector
  (`--memory=gc`, `compiler/runtime/memory_gc.c`) reclaims them — a conservative
  mark & sweep linked behind `@osp_alloc`, complete because the heap is acyclic
  [MEM-ACYCLIC]. On `binarytrees` it cuts peak RSS ~80× (905 MiB → ~11 MiB) with
  byte-identical output across every differential example (`make _conformance-gc`).
  The ARC default and a precise copying GC are the remaining work
  ([plan 0011](https://github.com/Nimblesite/osprey/blob/main/docs/plans/0011-arc-gc-implementation.md)); this spec is the
  contract they must satisfy.

## Collection Is Unobservable [MEM-OPAQUE]

No Osprey program can observe when, whether, or how memory is reclaimed.
Concretely:

- There are no finalizers or destructors, and there never will be — no code
  runs because a value died [MEM-OPAQUE-NO-FINALIZERS].
- There are no destruction-order or destruction-timing guarantees.
- No API exposes addresses, object identity beyond structural equality, or
  collector state.

A program whose output depends on reclamation behavior is not a valid Osprey
program; conforming implementations are free to differ on it. This rule is
what makes every backend below interchangeable.

## Resources Are Effects, Not Destructors [MEM-RESOURCES]

External resources (files, sockets, processes, handles) MUST be released by
*scoped* constructs — an effect handler that brackets acquire/release around
the code that uses the resource — never by tying release to a value's death.
This is forced by [MEM-OPAQUE]: value death has no observable timing, so it
can never be a release point.

## The Value Heap Is Acyclic [MEM-ACYCLIC]

Immutable values cannot reference values created after them, so reference
cycles cannot be constructed. Consequences:

- Reference counting is *complete* — no cycle collector is required, and a
  refcounting backend and a tracing backend are observationally identical.
- `mut` does not break this: reassignment rebinds the *name* to a new value;
  it never mutates a heap value in place (closure captures snapshot at
  creation per [TYPE-FN-CLOSURE]).

This is a constraint on language evolution: any future feature that allows a
heap value to be mutated to point at a younger value either preserves
acyclicity by construction or is rejected.

## Fibers Share Nothing [MEM-FIBER-ISOLATION]

Values cross fiber boundaries — `spawn` captures and channel `send` — by
move or by copy, never by sharing. No value is ever co-owned by two fibers.
Consequences:

- Each fiber's heap is single-threaded, so all reference counts are
  non-atomic.
- A fiber's values are reclaimable when the fiber completes, independent of
  other fibers.

## Ownership and the Shared Residue [MEM-OWNERSHIP]

Every heap value has an owner. The compiler infers ownership and statically
places the free wherever a value's last use is provable — the common case in
an immutable language.

The single construct that defeats static placement is **sharing**: two or
more live references to one value whose last use depends on runtime control
flow [MEM-OWNERSHIP-SHARED]. Canonical forms: structural sharing in
persistent data (`prepend(x, xs)` leaves `xs` and the result sharing a
spine), and aliased escaping closures. Shared values carry a non-atomic
reference count at runtime; everything else is freed statically. Sharing is
inferred — the developer never annotates it.

## Static Mode [MEM-STATIC-MODE]

Under `--static-memory`, compilation FAILS at every point where the
ownership analysis would insert a reference count, with a diagnostic naming
the shared value and the conflicting owners. A program accepted in static
mode contains **zero** runtime memory-management operations (no refcounts,
no collector) — Rust-class output without a borrow checker the developer
fights — and behaves byte-for-byte identically under the default mode.
Static-mode programs are a strict subset of Osprey, not a dialect.

### Barred Constructs [MEM-STATIC-MODE-BARRED]

Static mode bars exactly the constructs that create a shared residue:

1. **Live aliasing** — holding two or more references to one heap value
   past the point where a unique last owner is provable: `let g = f` where
   both escape, storing a value into a record or closure capture while the
   original binding stays live with divergent control flow.
2. **Built-in persistent collections** — `List` and `Map` (their spine/HAMT
   nodes share structure internally in the runtime); barred in static mode
   v1.

Everything else stays available: escaping closures with a unique owner,
records, unions, strings, `Result`, pattern matching, algebraic effects —
and fibers, because [MEM-FIBER-ISOLATION] moves or copies across the
boundary rather than sharing.

## Backend Conformance [MEM-BACKENDS]

Two backends ship out of the box, chosen at build time and invisible in
source code:

- **ARC (default)** — non-atomic reference counting on the shared residue,
  statically elided wherever ownership is provable. Complete without a
  cycle collector because the heap is acyclic [MEM-ACYCLIC].
- **Tracing GC** — the conformance oracle that keeps [MEM-OPAQUE] honest.

**Backend portability.** The two reclaiming backends need different things from
the host. The conservative tracing GC finds roots by scanning the native stack,
machine registers and data/BSS segments, so it runs on native targets only. ARC
is *precise* (the compiler inserts retain/release) and non-atomic, so it carries
to every target — including `wasm32`, where it is the *only* reclaiming option:
the conservative GC cannot scan a wasm stack, and the WebAssembly-GC proposal is
a separate, untargeted mechanism. See
[spec 0022](/spec/0022-webassemblytarget/) [WASM-TARGET-MEMORY].

A reclamation backend is conforming iff every differential-harness example
produces byte-identical output and reports zero leaked language values under
it.

### Custom Managers [MEM-BACKENDS-CUSTOM]

The backend boundary is a small C interface (alloc/retain/release/collect
hooks), and anyone may link their own manager against it — arenas, pools,
debugging allocators. Soundness of a custom manager is the supplier's
responsibility: the language's memory-safety guarantee covers only the
shipped backends, and a build linking a custom manager must say so visibly
(e.g. in `--version` output).