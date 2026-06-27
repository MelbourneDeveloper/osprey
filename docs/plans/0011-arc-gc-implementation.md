# Plan 0011 — Swappable reclaiming memory backends (ARC + tracing GC)

Realises the [MEM-BACKENDS] contract of
[spec 0018](../specs/0018-MemoryManagement.md): two robust, swappable memory
managers behind the existing `@osp_alloc` link-time boundary, plus a static
`--static-memory` mode (the "borrow-checker" subset). Reclamation stays
unobservable [MEM-OPAQUE], so every backend is observationally identical and
selected at link time, never in source.

## The governing facts (and the papers that justify them)

Three properties of the Osprey value heap collapse the usual GC design space:

1. **The heap is acyclic** [MEM-ACYCLIC]. Immutable values cannot reference
   values created after them, so cycles are unconstructable. Bacon, Cheng &
   Rajan, *A Unified Theory of Garbage Collection* (OOPSLA 2004) prove tracing
   and reference counting are duals computing the least / greatest fix-point of
   the same reference-count equation, and their difference is **exactly the
   cyclic garbage**. Acyclic ⇒ the fix-points coincide ⇒ **naive reference
   counting is complete** — no cycle collector, no trial deletion, no backup
   trace. This is the licence for ARC as the primary backend.
2. **Fibers share nothing** [MEM-FIBER-ISOLATION]. Each fiber's heap is
   single-threaded, so reference counts are **non-atomic** and a fiber's heap is
   collectable independently when it completes.
3. **Reclamation is unobservable** [MEM-OPAQUE]. No finalizers, no timing. Any
   two conforming backends produce byte-identical output — the conformance
   oracle (below).

## Backends

### `[GC-TRACE-CONSERVATIVE]` — tracing GC, **shipped first** (this plan, phase 1)

A conservative, non-moving **mark & sweep** over the managed heap reachable from
the C stack, machine registers, and the program's data/BSS segments —
Boehm & Weiser, *Garbage Collection in an Uncooperative Environment* (SP&E
1988), specialised to the acyclic heap so a single mark pass is complete
(Bacon/Cheng/Rajan 2004). A machine word is treated as a root iff it equals the
**base address of a known managed allocation**; false positives (integers that
look like pointers) only *retain* an object — they never corrupt it, because the
collector never moves. This needs **zero codegen changes**: it slots in behind
`@osp_alloc` purely at link time, which is why it ships first and why it is the
safe way to validate the whole boundary end-to-end.

- **Soundness scope v1:** collection runs only while the process is effectively
  single-threaded (the main thread is the sole allocator); the first allocation
  from any other thread permanently disables collection (a fiber's isolated heap
  is future work — precise per-fiber GC, phase 3). Every allocation and the
  whole collection run hold one mutex, so disabling is race-free.
- **Managed heap:** codegen allocations (`@osp_alloc`) plus the value-container
  runtime units (`list_runtime`, `map_runtime`, `map_runtime_hamt`) whose nodes
  store boxed Osprey values; recompiled in the GC archive with `malloc`/`calloc`/
  `realloc`/`free` redirected to the collector (`osp_gc_shim.h`). Fiber / HTTP /
  effect runtime keep libc `malloc` (never collected — status quo, sound).

### `[GC-ARC-PERCEUS]` — reference counting, **default backend** (phase 2)

Precise reference counting following **Perceus** (Reinking, Xie, de Moura,
Leijen, *Perceus: Garbage Free Reference Counting with Reuse*, PLDI 2021):

- **Borrow inference** `[GC-ARC-BORROW]` — owned vs borrowed parameters via the
  `collectO` fix-point of Ullrich & de Moura, *Counting Immutable Beans* (IFL
  2019). Inspectors compile reference-count-free.
- **dup/drop insertion** with borrowing to delay `dup` to the actual last use
  (Perceus λ¹ rules); **drop specialization** (the `is-unique` test) on the hot
  path.
- **Reuse analysis** `[GC-ARC-REUSE]` — `drop-reuse` tokens turn a unique
  matched cell into an in-place write (FBIP), so a functional `map`/`tree-map`
  runs as in-place mutation when uniquely owned.
- **Object header** (Koka model, one 8-byte word): pointer fields laid out
  first, a `scan_fsize` count for the generic drop/trace fallback, a 16-bit
  tag, and a **signed** non-atomic refcount (`0` = unique ⇒ cheapest free test,
  `<0` = cross-fiber / persistent ⇒ the only atomic path).

This requires the codegen work the conservative GC avoids: per-allocation type
info in the header, and a dup/drop insertion pass between type-checking and
codegen (which is today a direct AST→LLVM-text lowering with no SSA IR — the
pass introduces the structured form Perceus needs).

### `[GC-TRACE-CHENEY]` — precise copying GC (phase 3, conformance oracle)

Cheney semi-space copying (Cheney, *A Nonrecursive List Compacting Algorithm*,
CACM 1970) with **precise roots** via an LLVM shadow stack
(`llvm.gcroot`/`"shadow-stack"`) made per-fiber, reusing the phase-2 header type
info for tracing. Bump allocation, free compaction, GC cost ∝ live data. Immix
(Blackburn & McKinley, PLDI 2008) is the later upgrade. Primary role: the
oracle that keeps [MEM-OPAQUE] honest — must be byte-identical to ARC.

## The C ABI (uniform across backends)

```c
void* osp_alloc(int64_t size);          // the existing hook (all backends)
void  osp_retain(void* o);              // dup  — no-op under tracing
void  osp_release(void* o);             // drop — no-op under tracing
void  osp_collect(void);                // full GC — no-op under ARC (acyclic ⇒ complete)
```

`osp_retain`/`osp_release` are no-ops in the tracing backends; `osp_collect` is
a no-op under ARC. That asymmetry is exactly what makes the backends drop-in
swappable while observationally identical.

## Backend selection

Link-time, never in the IR (the IR names only `@osp_alloc`). `--memory=gc`
(default `--memory=default`, future `--memory=arc`) selects the runtime archive
(`libfiber_runtime_<backend>.a` / `libhttp_runtime_<backend>.a`) in the CLI's
`link_args`. The Makefile builds one archive set per backend; the default set is
untouched, so the default build/test path carries zero risk.

## Conformance `[MEM-BACKENDS]`

A backend is conforming iff every differential-harness example produces
byte-identical output and leaks zero language values under it. `make
conformance-gc` runs `crates/diff_examples.sh` with the backend selected; the
benchmark suite adds an `Osprey (GC)` column so `binarytrees` (905 MiB → a few
MiB) is visible next to the default.

## `[MEM-STATIC-MODE]` — the static "borrow-checker" subset (phase 4)

`--static-memory` fails compilation at every point the ownership analysis would
insert a reference count, naming the shared value and the conflicting owners —
Rust-class output with no runtime memory management, a strict subset of Osprey
that behaves byte-for-byte identically under the default mode. Built on the
phase-2 borrow/ownership analysis (a program is static-mode-clean iff that
analysis inserts no `dup`/`drop` on a shared residue).

## Phasing

1. **Conservative tracing GC** + link-time selection + benchmark column +
   conformance target. *(this change)*
2. Header type-info in codegen; Perceus borrow inference + dup/drop insertion +
   reuse ⇒ the ARC default backend.
3. Precise Cheney copying GC (per-fiber shadow-stack roots) as the oracle.
4. `--static-memory`.

## References

- Bacon, Cheng, Rajan. *A Unified Theory of Garbage Collection.* OOPSLA 2004.
- Reinking, Xie, de Moura, Leijen. *Perceus: Garbage Free Reference Counting
  with Reuse.* PLDI 2021.
- Ullrich, de Moura. *Counting Immutable Beans.* IFL 2019.
- Cheney. *A Nonrecursive List Compacting Algorithm.* CACM 13(11), 1970.
- Blackburn, McKinley. *Immix: A Mark-Region Garbage Collector.* PLDI 2008.
- Boehm, Weiser. *Garbage Collection in an Uncooperative Environment.* SP&E 1988.
