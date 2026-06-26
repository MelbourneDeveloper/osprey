// Default Osprey memory backend — the swappable allocation boundary.
//
// Implements [MEM-BACKENDS] / [MEM-BACKENDS-CUSTOM] (docs/specs/0018). Compiler
// codegen emits calls to `osp_alloc` and never names `malloc`, so the memory
// manager is chosen at link time, never baked into the IR. This default backend
// is a `malloc` passthrough with no reclamation yet — matching the current
// "allocate, never free during a run" semantics, which is sound because
// reclamation is unobservable [MEM-OPAQUE]. A custom manager (ARC, tracing GC,
// arena, pool) replaces this object by linking its own `osp_alloc` — plus the
// future `osp_retain` / `osp_release` / `osp_collect` hooks — against the same
// symbols.
//
// The IR-level allocator attributes on `@osp_alloc` (see osprey-codegen
// builder.rs OSP_ALLOC_DECL) let LLVM remove provably non-escaping allocations
// at -O2, so most allocations never reach this function at all.

#include <stdint.h>
#include <stdlib.h>

void *osp_alloc(int64_t size) { return malloc((size_t)size); }
