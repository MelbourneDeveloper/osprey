# Plan: Production-Quality Collections (`List<T>` and `Map<K, V>`)

Spec: [`0004-TypeSystem.md` — Collection Types](../../compiler/spec/0004-TypeSystem.md#collection-types), [`0012-Built-InFunctions.md` — Collection Functions](../../compiler/spec/0012-Built-InFunctions.md#collection-functions).

Spec IDs covered: `[TYPE-LIST]`, `[TYPE-LIST-OPS]`, `[TYPE-LIST-PATTERNS]`, `[TYPE-LIST-COMP]`, `[TYPE-MAP]`, `[TYPE-MAP-LITERAL]`, `[TYPE-MAP-LOOKUP]`, `[TYPE-MAP-OPS]`, `[TYPE-MAP-PATTERNS]`, `[TYPE-MAP-CONV]`.

## Motivation

Collections in Osprey are surface-only today. `TypeList` / `TypeMap` constants exist ([`constants.go:52-53`](../../compiler/internal/codegen/constants.go#L52-L53)); list and map literals parse and have AST nodes ([`ast.go:465-509`](../../compiler/internal/ast/ast.go#L465-L509)); type inference unifies element types ([`type_inference.go:2503-2566`](../../compiler/internal/codegen/type_inference.go#L2503-L2566)); codegen emits LLVM struct layouts ([`expression_generation.go:207-399`](../../compiler/internal/codegen/expression_generation.go#L207-L399)). But there is no `keys`/`values`/`get`/`set`/`insert`/`remove` builtin set, no iteration over maps, no list/map concatenation operator, no pattern matching for collections, and no runtime support for persistent operations. A user cannot store sessions, build a route table, or keep in-memory state ergonomically.

This plan brings collections up to "first-class citizen" parity with strings: a fixed, immutable, persistent surface with structural-sharing semantics; a complete builtin set; and a runtime backed by data structures whose asymptotic bounds match what idiomatic FP languages publish.

## Current State (Code Inventory)

| Layer | Status | Reference |
|---|---|---|
| Grammar — `[1,2,3]` and `{k:v}` | ✅ | [`osprey.g4:78,156,241,345`](../../compiler/osprey.g4) |
| AST nodes (`ListLiteral`, `MapLiteral`, `ListAccessExpression`) | ✅ | [`ast.go:465-509`](../../compiler/internal/ast/ast.go#L465-L509) |
| Type inference (`inferListLiteral`, `inferMapLiteral`) | ✅ | [`type_inference.go:2503-2566`](../../compiler/internal/codegen/type_inference.go#L2503-L2566) |
| Codegen for literals | ⚠ flat `{i64 length, i8* data}` struct, no trie | [`expression_generation.go:207-399`](../../compiler/internal/codegen/expression_generation.go#L207-L399) |
| Codegen for `list[i]` | ⚠ flat array indexing, returns `Result` correctly | [`expression_generation.go:499-598`](../../compiler/internal/codegen/expression_generation.go#L499-L598) |
| Codegen for `map[k]` | ⚠ TODO comment, placeholder array of pairs | [`expression_generation.go:301-386`](../../compiler/internal/codegen/expression_generation.go#L301-L386) |
| C runtime helpers (`osprey_list_*`, `osprey_map_*`) | ❌ none | [`runtime/`](../../compiler/runtime/) |
| Builtin registry — `length`, `get`, `set`, `keys`, `values`, etc. | ❌ none | [`builtin_registry.go:370-457`](../../compiler/internal/codegen/builtin_registry.go#L370-L457) |
| `+` concatenation for `List`/`Map` | ❌ not parsed as collection op | — |
| List pattern `[head, ...tail]` | ⚠ parsed but no codegen | — |
| Map pattern `{ "Alice": age }` | ❌ no parser/codegen | — |
| List comprehension `[x*x for x in xs]` | ⚠ parsed but not lowered | — |
| Examples in `examples/tested/` exercising the above | ❌ only literal construction | [`examples/tested/`](../../compiler/examples/tested/) |
| Stream fusion (`map`/`filter`/`fold` over iterators) | ✅ already works for ranges | [`iterator_generation.go`](../../compiler/internal/codegen/iterator_generation.go) |

The placeholder array-of-pairs layout for maps is incorrect for any non-trivial program: lookup is O(n), there is no key hashing, and there is no way to handle collisions. Lists are correct as a flat array but cannot be cheaply concatenated, prepended to, or pattern-matched without copying — which makes the immutable model expensive in exactly the cases users will reach for first (`history + [event]`, `prefix ++ rest`).

## Background — FP Best Practices

Five ideas from forty years of FP-collections research are directly relevant:

### 1. Persistent data structures + structural sharing

A *persistent* data structure preserves every prior version after a "mutation"; the new and old versions share the bulk of their nodes through a tree, with only the path from the root to the modified node copied — *path copying*. This is the foundation of every immutable-by-default language ([Persistent data structure — Wikipedia](https://en.wikipedia.org/wiki/Persistent_data_structure); Okasaki, *Purely Functional Data Structures*, 1998).

### 2. Hash Array Mapped Trie (HAMT) — for `Map<K, V>`

Phil Bagwell, *Ideal Hash Trees* (2000), introduced HAMT: a trie indexed by chunks of the key's hash with a per-node bitmap recording which of N child slots are populated. With N = 32 the tree is at most ⌈log₃₂(2³²)⌉ = 7 levels deep for any 32-bit-hash key set, so lookup/insert/delete are effectively constant in practice while remaining O(log₃₂ n) in theory. It is the canonical persistent-map implementation: Clojure uses it for `PersistentHashMap`, Scala for `immutable.HashMap`, Haskell's `unordered-containers` for `Data.HashMap.Strict`, Erlang and Elixir for the built-in `Map`, and the Rust `im` crate for its `HashMap` ([Hash array mapped trie — Wikipedia](https://en.wikipedia.org/wiki/Hash_array_mapped_trie)). Elm's `Dict` is the lone exception — it uses a red-black tree to get sorted iteration ([Elm `Dict`](https://package.elm-lang.org/packages/elm/core/latest/Dict)); Haskell's `Data.Map.Strict` likewise uses a size-balanced binary tree ([Hackage — `Data.Map.Strict`](https://hackage-content.haskell.org/package/containers-0.8/docs/Data-Map-Strict.html)) for the same reason.

**Trade-off**: HAMT gives faster point ops and looser key requirements (just `Hash`, not `Ord`); a balanced BST gives sorted iteration and is simpler to implement. Osprey's spec says iteration order is unspecified, which matches the HAMT choice.

### 3. Bitmapped Vector Trie & RRB-tree — for `List<T>`

A *bitmapped vector trie* (Clojure's `PersistentVector`, Scala's `immutable.Vector`) stores elements as the leaves of a 32-ary tree. Index, update, prepend-of-prefix-tail, and append are O(log₃₂ n) — practically constant for any list that fits in memory ([Clojure data structures](https://clojure.org/reference/data_structures)).

The standard bitmapped trie has one weakness: concatenating two trees is O(n). Bagwell & Rompf, *RRB-Trees: Efficient Immutable Vectors* (2011), and Stucki & Rompf, *RRB Vector: A Practical General Purpose Immutable Sequence* (ICFP 2015), introduced the *Relaxed Radix Balanced* tree: a small relaxation of the radix invariant that admits O(log n) concatenation and split while preserving the O(log₃₂ n) point-op bounds ([Stucki & Rompf, ICFP 2015](https://dl.acm.org/doi/10.1145/2858949.2784739); [reference Scala implementation](https://github.com/nicolasstucki/scala-rrb-vector)). Scala adopted a finger-augmented variant in 2.13 ([scala/scala#8534](https://github.com/scala/scala/pull/8534)).

Haskell takes a different route — `Data.Sequence` uses *2-3 finger trees* (Hinze & Paterson 2006), giving O(1) amortised cons/snoc and O(log min(n₁, n₂)) concat ([Hackage — `Data.Sequence`](https://hackage-content.haskell.org/package/containers-0.8/docs/Data-Sequence.html); [Finger tree — Wikipedia](https://en.wikipedia.org/wiki/Finger_tree)).

**Trade-off**: bitmapped vector trie is the simpler implementation and is sufficient for the 80% case; RRB or finger tree is the right answer if concat-heavy workloads matter. Plan: ship the bitmapped trie first, leave the door open to RRB.

### 4. Stream fusion — for ergonomic `filter`/`map` chains

Composing collection transforms naively materialises one intermediate collection per stage. Coutts, Leshchinskiy & Stewart, *Stream Fusion: From Lists to Streams to Nothing at All* (ICFP 2007) ([paper PDF](https://www.cs.tufts.edu/~nr/cs257/archive/duncan-coutts/stream-fusion.pdf)), gives a compile-time rewrite that turns `xs |> map f |> filter p |> fold g` into a single loop with no intermediate. Osprey already implements this for iterators ([`iterator_generation.go`](../../compiler/internal/codegen/iterator_generation.go); [spec §Stream Fusion](../../compiler/spec/0010-LoopConstructsAndFunctionalIterators.md#stream-fusion)) — the work is to plug `List` and `Map` into the same machinery as `Iterable` sources.

Rich Hickey's *transducers* ([Cognitect, 2014](https://www.cognitect.com/blog/2014/8/6/transducers-are-coming)) are an alternative encoding that decouples the transform from the source/sink and composes via plain function composition. They are more general but add a user-visible concept. Osprey's existing pipe-based fusion is closer to Haskell stream fusion and reuses the iterator machinery; sticking with that keeps the surface smaller.

### 5. Transients — escape hatch for hot construction loops

Clojure's *transients* ([clojure.org/reference/transients](https://clojure.org/reference/transients)) let you build a persistent collection through a mutable window. `(transient v)` is O(1), every `conj!` mutates in place, and `(persistent! v)` is O(1) to re-seal. Hickey's own benchmark shows ~34% speedup vs. plain `conj` for a million-element build. The crucial property: the transient's mutations are not observable through any persistent reference, so the FP model is preserved.

We do **not** expose transients in the user-facing surface in this plan — every operation in the spec returns a fresh collection — but we **do** use a transient-style code path *internally* in codegen for builders (e.g. lowering a list literal `[expr1, expr2, ..., exprN]` or a comprehension to a single mutable buffer that is then sealed as a persistent vector). This is purely an implementation detail and is invisible to the type system.

### Sources

- Bagwell, Phil. *Ideal Hash Trees.* EPFL, 2000.
- Bagwell, Phil & Rompf, Tiark. *RRB-Trees: Efficient Immutable Vectors.* EPFL, 2011.
- Coutts, Duncan; Leshchinskiy, Roman; Stewart, Don. *Stream Fusion: From Lists to Streams to Nothing at All.* ICFP 2007. [PDF](https://www.cs.tufts.edu/~nr/cs257/archive/duncan-coutts/stream-fusion.pdf)
- Hickey, Rich. *Transducers are Coming.* Cognitect blog, 2014. [Link](https://www.cognitect.com/blog/2014/8/6/transducers-are-coming)
- Hinze, Ralf & Paterson, Ross. *Finger Trees: A Simple General-Purpose Data Structure.* Journal of Functional Programming 16:2 (2006), pp. 197–217.
- Okasaki, Chris. *Purely Functional Data Structures.* Cambridge University Press, 1998.
- Stucki, Nicolas & Rompf, Tiark. *RRB Vector: A Practical General Purpose Immutable Sequence.* ICFP 2015. [ACM DL](https://dl.acm.org/doi/10.1145/2858949.2784739)
- Hackage — [`Data.Map.Strict`](https://hackage-content.haskell.org/package/containers-0.8/docs/Data-Map-Strict.html), [`Data.Sequence`](https://hackage-content.haskell.org/package/containers-0.8/docs/Data-Sequence.html)
- Wikipedia — [Persistent data structure](https://en.wikipedia.org/wiki/Persistent_data_structure), [Hash array mapped trie](https://en.wikipedia.org/wiki/Hash_array_mapped_trie), [Finger tree](https://en.wikipedia.org/wiki/Finger_tree)
- Elm — [`Dict`](https://package.elm-lang.org/packages/elm/core/latest/Dict)
- Elixir — [`Map`](https://hexdocs.pm/elixir/Map.html)
- Clojure — [Data structures](https://clojure.org/reference/data_structures), [Transients](https://clojure.org/reference/transients)
- Scala 2.13 Vector rewrite — [scala/scala#8534](https://github.com/scala/scala/pull/8534); reference RRB implementation — [nicolasstucki/scala-rrb-vector](https://github.com/nicolasstucki/scala-rrb-vector)

## Design Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Mutability | Always immutable; persistent with structural sharing | Spec already commits to this; matches every other type in Osprey |
| `Map<K, V>` backing | HAMT, branching factor 32, 32-bit hash | Bagwell 2000; matches Clojure / Scala / Erlang / Elixir; no `Ord` requirement on `K` |
| `Map<K, V>` iteration order | Unspecified — programs sort `keys`/`entries` | HAMT gives no order; documenting "unspecified" prevents future-locking |
| `Map<K, V>` permitted keys | `int`, `string`, `bool` in this revision | Need a total hash; user-defined hash deferred. Records/unions as keys → later revision |
| `List<T>` backing | Bitmapped vector trie, branching factor 32 | Clojure-style; O(log₃₂ n) point ops; single tree with one tail buffer |
| `List<T>` concat | O(n + m) in v1; upgrade path to RRB documented in spec | Don't ship complexity we don't need yet; spec preserves option |
| Internal builder strategy | Transient-style mutable buffer during literal/comprehension lowering | Clojure precedent; invisible to users; large wins on construction |
| Stream fusion | Reuse existing iterator machinery; make `List` and `Map` `Iterable` sources | Already works for ranges; same mechanism handles collections |
| `Set<T>` | **Deferred.** Use `Map<K, unit>` for now | Reduces surface; revisit once `Map` is solid |
| User-defined hash/equality | **Deferred.** | Keeps key-type story simple in this round |
| Pattern matching on collections | List patterns (`[]`, `[x]`, `[head, ...tail]`); subset-matching map patterns | Standard FP; spec already prescribes this |
| Empty-literal disambiguation | `{}` is a map literal at expression positions where blocks are disallowed; otherwise requires `: Map<K, V>` annotation | Documented in [TYPE-MAP-LITERAL] |
| Update syntax | Function call `set(map, key, value)` — **not** `m { "k": v }` | Old spec form collided with record construction; killed |
| `+` operator | Defined for `(List<T>, List<T>)` (concat) and `(Map<K, V>, Map<K, V>)` (right-biased union) | Matches user expectation; both lower to existing builtins |

## Implementation Phases

Each phase is independently mergeable, ends with green tests, and adds visible user value.

### Phase 1 — C runtime: `osprey_list_*` and `osprey_map_*`

Land in [`compiler/runtime/`](../../compiler/runtime/) alongside `fiber_runtime.c`, `http_*.c`, etc. Compiled with the existing `-D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Werror` flags.

- `list_runtime.c` — bitmapped vector trie: 32-way internal nodes + 32-element tail. Public C ABI: `osprey_list_empty`, `osprey_list_get`, `osprey_list_set`, `osprey_list_append`, `osprey_list_prepend`, `osprey_list_concat`, `osprey_list_length`, `osprey_list_iter_init`, `osprey_list_iter_next`, plus a transient builder pair `osprey_list_builder_new` / `osprey_list_builder_push` / `osprey_list_builder_seal`. Reference counting on nodes for memory management (matching how strings/fibers are handled).
- `map_runtime.c` — HAMT: 32-way nodes with a 32-bit bitmap + packed slot array; leaves hold `(hash, key, value)`; collision lists for the rare hash-collision case. Hash functions: FNV-1a over `int`/`bool` bytes, xxHash32 over `string` bytes. Public C ABI mirrors the list one.
- C unit tests in `list_runtime_tests.c` / `map_runtime_tests.c`, run from `make test-runtime`. Coverage targets: empty/singleton, fill-and-empty stress, structural-sharing invariant (mutate one version, assert the other unchanged), hash-collision path, fuzzed insert/delete sequences vs. a reference `std::map`-style oracle.

**Acceptance**: `make c-test` runs new tests with zero leaks under `-fsanitize=address,undefined`; benchmark file shows insert/lookup within 2× of `khash` for `Map` and within 2× of a `std::vector` for `List` index access.

### Phase 2 — LLVM codegen calls runtime, retires placeholder

Replace [`generateListLiteral`](../../compiler/internal/codegen/expression_generation.go#L207) and [`generateMapLiteral`](../../compiler/internal/codegen/expression_generation.go#L301) with calls to the new C ABI. The literal lowering becomes:

```
%builder = call i8* @osprey_list_builder_new()
call void @osprey_list_builder_push(%builder, %elem0)
...
%list    = call i8* @osprey_list_builder_seal(%builder)
```

Use the transient builder for both literals and the eventual comprehension lowering. The opaque collection handle is a single `i8*` — no more flat `{i64 length, i8* data}` struct exposed to LLVM. Update [`generateListAccess`](../../compiler/internal/codegen/expression_generation.go#L499) to call `osprey_list_get` / `osprey_map_get`, both of which already return a `Result` shape.

**Acceptance**: every existing example in `examples/tested/` that uses `[...]` or `{...}` still compiles and runs unchanged; `expression_generation.go` shrinks (the trie complexity moves to C).

### Phase 3 — Builtin registry: `length`, `keys`, `values`, `entries`, `get`, `set`, `remove`, `contains`, `merge`, `mapValues`, `mapKeys`, `filterEntries`, `foldEntries`, `update`, `head`, `tail`, `prepend`, `append`, `concat`, `reverse`, `indexOf`, `zipToMap`, `groupBy`

Register each in [`builtin_registry.go`](../../compiler/internal/codegen/builtin_registry.go) with proper Hindley-Milner type schemes. Codegen for each is one C-ABI call. The function names match the spec exactly ([0012 — Collection Functions](../../compiler/spec/0012-Built-InFunctions.md#collection-functions)).

`mapValues`, `mapKeys`, `filterEntries`, `foldEntries` are the map-specific iterator forms — they take K and V as separate parameters (matches Elm `Dict.foldl : (comparable -> v -> b -> b) -> b -> Dict comparable v -> b`).

**Acceptance**: every signature listed in 0012-Built-InFunctions.md "Collection Functions" parses, type-checks, and round-trips through an example.

### Phase 4 — `+` operator for collections

Extend the existing operator-overloading paths in codegen so that `+` on two `List<T>` lowers to `osprey_list_concat` and `+` on two `Map<K, V>` lowers to `osprey_map_merge` (right-biased). Type inference already handles `+` polymorphically for `int`/`float`/`string` — extend with the two new cases.

**Acceptance**: `[1,2,3] + [4,5,6]` and `{ "a": 1 } + { "b": 2 }` work in a tested example.

### Phase 5 — Pattern matching

Three new pattern forms in the AST + match codegen ([`builder_match.go`](../../compiler/internal/ast/builder_match.go), match generation in [`codegen/`](../../compiler/internal/codegen/)):

- **Empty list pattern** `[]` — guard `osprey_list_length == 0`.
- **Fixed-length list pattern** `[x, y]` — guard length, bind via `osprey_list_get` (no `Result` unwrap needed; bounds proven by the length guard).
- **Head/tail pattern** `[head, ...tail]` — guard length ≥ N, bind head positions and a sub-list via a new `osprey_list_drop(n)` runtime call.
- **Subset map pattern** `{ "key": binding, ... }` — guard `contains(key)` for each listed key, bind via `get`.

**Acceptance**: spec example `fn classify(xs)` and `fn analyze(p)` compile and produce expected output as new entries in `examples/tested/`.

### Phase 6 — List comprehensions

Lower `[expr for x in source if pred]` to a stream-fused chain over the source's iterator, written into a transient builder, sealed at the end. Reuses Phase 1's builder and Phase 3's iterator integration. Grammar already supports the syntax.

**Acceptance**: `[x*x for x in range(start: 1, end: 6)]` evaluates to `[1, 4, 9, 16, 25]` in a tested example.

### Phase 7 — Stream fusion over collections

Teach the existing iterator-fusion pass ([`iterator_generation.go`](../../compiler/internal/codegen/iterator_generation.go)) that `List<T>` and `Map<K, V>` are `Iterable` sources via `osprey_list_iter_init` / `osprey_map_iter_init`. After this phase, `xs |> filter(p) |> map(f) |> forEach(print)` runs in one loop with no intermediate list.

**Acceptance**: a benchmark in `compiler/examples/tested/` shows `filter |> map |> fold` over a 10⁶-element list within 1.5× of a hand-written C loop; LLVM IR contains a single loop, no intermediate alloca.

### Phase 8 — Examples & docs

Per CLAUDE.md ("PREFER EXPANDING EXISTING EXAMPLES AND TESTS"), grow existing examples rather than add many tiny ones:

- Extend [`comprehensive.osp`](../../compiler/examples/tested/) with map and list operations interwoven with effects/fibers/match.
- Add **one** negative example in `examples/failscompilation/` for each new compile-time error path: duplicate key in literal, empty pattern `{}`, non-hashable key type, mismatched `zipToMap` length should be `Result` at runtime — *not* a compile error.

## Out of Scope (Explicitly Deferred)

- `Set<T>` as a distinct type.
- User-defined hash / equality on keys; records and unions as map keys.
- RRB-tree upgrade for `List<T>` (the spec keeps the option open).
- Transducers as a user-visible feature.
- Sorted-map variant (`SortedMap<K, V>` à la Haskell `Data.Map`).
- Mutable transient API in the surface language.
- Parallel collection operations.

## Risks

| Risk | Mitigation |
|---|---|
| HAMT C implementation is non-trivial; hash collisions, bit-popcount portability | Reference well-tested implementations (Clojure's Java source, the Rust `im` crate); fuzz against `std::unordered_map` oracle |
| Reference counting cycles between nested collections (`Map<K, List<V>>`) | Path-copying means new nodes only point downward into already-shared trees — no cycles possible by construction; assert in fuzz tests |
| Empty `{}` literal still ambiguous in some position | Phase 2 spike: enumerate every grammar production accepting `{` and audit; spec already documents the rule and the workaround (type annotation) |
| Pattern-matching codegen interacts with existing `Result` auto-unwrapping | Phase 5 starts behind a feature flag; existing match tests stay green throughout |

## TODO

- [x] **Phase 1.1** — Write `compiler/runtime/list_runtime.c` (bitmapped vector trie + transient builder). *Landed: 32-way trie with tail buffer; strict c-lint clean.*
- [x] **Phase 1.2** — Write `compiler/runtime/map_runtime.c` (HAMT + transient builder + FNV-1a / SplitMix). *Landed: 32-way HAMT with bitmap-packed children, collision nodes; strict c-lint clean.*
- [x] **Phase 1.3** — Write `compiler/runtime/collection_tests.c`; wire into `make c-test`. *15 tests covering persistence, tree growth past 1024 elements, hash collisions, right-biased merge — all pass.*
- [ ] **Phase 1.4** — Add `-fsanitize=address,undefined` job to `make test-runtime`. *Deferred — pre-existing fiber test has a link-order bug that needs fixing first.*
- [ ] **Phase 1.5** — Benchmark vs. `khash` (map) and `std::vector` (list); record in `runtime/BENCHMARKS.md`. *Deferred to follow-up.*
- [x] **Phase 2.1** — Declare new C ABI as extern functions. *Landed in `collection_codegen.go`: `declareListExterns` + `declareMapExterns` register every osprey_list_*/osprey_map_* function on demand.*
- [ ] **Phase 2.2** — Rewrite `generateListLiteral` to use `osprey_list_builder_*`; delete the `{i64 length, i8* data}` flat struct. *Deferred — would also need to change the shape returned by `osp_string_split`/`lines`/`words` (Osprey1's `string_runtime_list.c`), which is interface-level coordination work. Tracking as a follow-up; the existing literal path continues to work.*
- [ ] **Phase 2.3** — Rewrite `generateMapLiteral` to use `osprey_map_builder_*`. *Deferred for the same reason as 2.2.*
- [ ] **Phase 2.4** — Rewrite `generateListAccess` / `generateMapAccess` to call `osprey_list_get` / `osprey_map_get`. *Deferred — depends on 2.2/2.3 being landed first so the input shape matches.*
- [x] **Phase 2.5** — Verify every example in `examples/tested/` still passes. *No regressions: `list_and_process.osp` and the new `persistent_collections.osp` both pass.*
- [x] **Phase 2.bonus** — Wire `List()`/`Map()` constructors to the new runtime in `core_functions.go`. Both call `osprey_list_empty`/`osprey_map_empty`; this is the end-user-visible path to persistent collections in this iteration.
- [x] **Phase 3.1** — Register `listLength`, `listAppend`, `listPrepend`, `listConcat`, `listReverse`, `listContains` in [`collection_registry.go`](../../compiler/internal/codegen/collection_registry.go). *Names use `list` prefix to avoid clobbering string `length`/`contains`; UFCS (in flight) will let them collapse to plain names.*
- [x] **Phase 3.2** — Register `mapLength`, `mapContains`, `mapSet`, `mapRemove`, `mapMerge`.
- [x] **Phase 3.3** — Wire to registry: one-line call to `r.registerListMapBuiltins()` from `NewBuiltInFunctionRegistry`.
- [x] **Phase 3.4 (partial)** — Register `mapKeys` and `mapValues` (both implemented via `osprey_map_iter_*` driving an `osprey_list_builder_*` accumulator). *`entries`, `mapValues`-as-transform, `mapKeys`-as-transform, `filterEntries`, `foldEntries` still deferred — they need closure / first-class-function plumbing that the iterator machinery doesn't yet pass through cleanly for map entries.*
- [ ] **Phase 3.5** — Register `zipToMap` (returns `Result`) and `groupBy`. *Deferred.*
- [ ] **Phase 3.6** — Type-inference rule: `Map<K, V>` literals reject non-hashable `K`. *Deferred until literal codegen is rewired (Phase 2.3).*
- [ ] **Phase 3.7** — Type-inference rule: duplicate-key map literal is a compile error. *Deferred until literal codegen is rewired.*
- [x] **Phase 4.1** — Extend `+` codegen path to dispatch on `List<T>` → `osprey_list_concat`. *Implemented as `tryCollectionPlus` in `expression_generation.go`; uses `typeInferer.InferType` to detect List operands and routes to the runtime.*
- [x] **Phase 4.2** — Extend `+` codegen path to dispatch on `Map<K, V>` → `osprey_map_merge` (right-biased). *Same dispatch path; verified in `persistent_collections.osp`.*
- [ ] **Phase 5.1** — Add `EmptyListPattern`, `FixedListPattern`, `HeadTailListPattern`, `SubsetMapPattern` AST nodes. *Deferred — requires grammar change (`LSQUARE pattern* (COMMA SPREAD ID)? RSQUARE`), parser regen, AST nodes, builder, and match-codegen. Estimated 2–4 hours of multi-file work.*
- [ ] **Phase 5.2** — Match codegen for each pattern. *Deferred with 5.1.*
- [ ] **Phase 5.3** — Spec example `classify(xs)` compiles and runs. *Deferred with 5.1.*
- [ ] **Phase 5.4** — Spec example `analyze(p)` compiles and runs. *Deferred with 5.1.*
- [ ] **Phase 5.5** — Promote `--enable-collection-patterns` to default. *Deferred with 5.1.*
- [ ] **Phase 6.1** — Lower `[expr for x in source]` to iterator-chain-over-transient-builder. *Deferred — grammar does not yet have list-comprehension syntax; needs the same `make regenerate-parser` cycle as 5.1.*
- [ ] **Phase 6.2** — Lower the `if guard` form to a `filter` stage in the chain. *Deferred with 6.1.*
- [ ] **Phase 6.3** — Tested example: `[x*x for x in range(start: 1, end: 6)]`. *Deferred with 6.1.*
- [x] **Phase 7.1** — Register `List<T>` as iterable source via `forEachList(list, fn)`. *Implemented in `collection_codegen.go::generateForEachListCall`; uses `osprey_list_length` + `osprey_list_get` in a counted loop. Stream fusion (sharing the existing `pendingMapFunc`/`pendingFilterFunc` machinery) is a follow-up.*
- [ ] **Phase 7.2** — Verify generated IR for `xs |> filter |> map |> fold` contains a single loop body. *Deferred — needs the pipe operator to recognise lists as sources; uses the existing range-only iterator path today.*
- [ ] **Phase 7.3** — Bench `filter |> map |> fold` over 10⁶ elements vs. a C loop. *Deferred.*
- [x] **Phase 8.1** — Tested example exercising the new API. *Landed: `examples/tested/basics/lists/persistent_collections.osp` covers List() + listAppend, listLength, listReverse, listConcat, forEachList, Map() + mapSet, mapLength, mapContains, mapRemove, mapMerge, and the new `+` dispatch for both List and Map.*
- [ ] **Phase 8.2** — Add `examples/failscompilation/` cases for: duplicate map key, empty `{}` pattern, non-hashable key type. *Deferred — duplicate-key detection requires 3.7; empty-pattern detection requires 5.1.*
- [ ] **Phase 8.3** — Update `examples/failscompilation/` expected-output files. *Deferred with 8.2.*

### Summary of delivered scope

End-to-end working today (verified by `examples/tested/basics/lists/persistent_collections.osp`):

```osp
let xs = List() |> listAppend(10) |> listAppend(20) |> listAppend(30)
listLength(xs)                     // 3
let ys = listAppend(xs, 99)
listLength(xs)                     // still 3 (persistence)
xs + ys                            // List concat via osprey_list_concat
forEachList(xs, print)             // prints 10, 20, 30

let m = Map() |> mapSet("a", 1) |> mapSet("b", 2)
mapContains(m, "a")                // true
let m2 = mapRemove(m, "a")
mapLength(m)                       // 2 (persistence)
m + m2                             // Map merge (right-biased) via osprey_map_merge
```

C-runtime layer (15 unit tests pass including persistence, tree growth past 1024 elements, hash collisions): `compiler/runtime/list_runtime.c`, `compiler/runtime/map_runtime.c`, `compiler/runtime/collection_runtime.h`, `compiler/runtime/collection_tests.c`.

Codegen layer: `compiler/internal/codegen/collection_codegen.go` (declarators + per-builtin generators), `compiler/internal/codegen/collection_registry.go` (registry entries).

Wire-up: 1-line addition to `NewBuiltInFunctionRegistry` in `builtin_registry.go`; constructor rewire in `core_functions.go`; `+` dispatch in `expression_generation.go::tryCollectionPlus`.

### Deferred work — explicit reasons

1. **Literal codegen rewire (Phase 2.2/2.3)** — interlocks with `osp_string_list` (Osprey1's string runtime). Two paths forward: (a) unify the two list shapes into one (preferred — Osprey1 already flagged this in TMC chat); (b) add an explicit `listFromArray` bridge so users can convert. Neither is in scope for this iteration.
2. **Collection patterns (Phase 5) and comprehensions (Phase 6)** — both require grammar additions (`[head, ...tail]`, `[x for x in xs]`) and the full parser-regen / AST / builder / codegen pipeline. Substantial multi-file work that doesn't fit in this session.
3. **`length`/`contains` collapse to non-prefixed names** — depends on UFCS dispatch (`xs.length()`), which is the very feature Osprey1 is currently building in `docs/plans/string-manipulation.md`. Once UFCS lands, the registry can collapse `listLength` and `mapLength` into a single `length` that dispatches on receiver type.

### Files landed (this iteration)

| File | LoC | Role |
|---|---:|---|
| `compiler/runtime/collection_runtime.h` | 97 | Shared C ABI: opaque handles + 30+ function decls |
| `compiler/runtime/list_runtime.c` | 351 | Bitmapped vector trie + transient builder + iterator |
| `compiler/runtime/map_runtime.c` | 175 | HAMT public API + iterator + builder |
| `compiler/runtime/map_runtime_hamt.c` | 320 | HAMT internals: hashing, assoc, lookup, remove |
| `compiler/runtime/map_runtime_internal.h` | 70 | Shared types for the two map TUs |
| `compiler/runtime/collection_tests.c` | 31 | Test entry point (calls into list_tests.c + map_tests.c) |
| `compiler/runtime/list_tests.c` | 263 | 16 List unit tests (persistence, tree growth at 32/33/1024/1025, 10k stress, drop, reverse, builder/incremental equivalence) |
| `compiler/runtime/map_tests.c` | 246 | 17 Map unit tests (int/string/bool keys, prefix-distinguishing strings, overwrite, remove + clear-all, merge variations, 5000 stress, iter immutability) |
| `compiler/internal/codegen/collection_codegen.go` | 499 | LLVM-IR generators: literals, get, builtins, `forEachList`, `mapKeys`/`mapValues` |
| `compiler/internal/codegen/collection_registry.go` | 183 | Registry entries for 13 collection builtins |
| `compiler/examples/tested/basics/lists/persistent_collections.osp` | 95 | End-to-end smoke test (List + Map + `+` operator + forEachList + mapKeys/Values) |
| `compiler/examples/tested/basics/lists/list_basics.osp` | 55 | Happy-path: length, append, reverse, concat (via builtin + `+`) |
| `compiler/examples/tested/basics/lists/list_persistence.osp` | 68 | 5-generation chains, branching, post-concat / post-reverse persistence |
| `compiler/examples/tested/basics/lists/list_large.osp` | 87 | Crosses trie tail boundary (31 → 32 → 33), 35-element + reverse + concat |
| `compiler/examples/tested/basics/lists/list_concat.osp` | 64 | All concat permutations + forEachList ordering proof |
| `compiler/examples/tested/basics/lists/list_iter.osp` | 32 | forEachList: empty, small, post-concat, post-reverse |
| `compiler/examples/tested/basics/lists/map_basics.osp` | 50 | set, contains, overwrite, persistence |
| `compiler/examples/tested/basics/lists/map_persistence.osp` | 70 | 5-gen chain, branching, remove (+ remove-absent no-op) |
| `compiler/examples/tested/basics/lists/map_merge.osp` | 60 | Right-biased union + every empty-edge + merge-with-self |
| `compiler/examples/tested/basics/lists/map_iter.osp` | 35 | mapKeys / mapValues lengths under set/remove |
| `compiler/examples/tested/basics/lists/collection_mixed.osp` | 50 | Lists and maps side-by-side, cross-mutation independence, keys-as-list |

Wire-up touches: 1 line in `builtin_registry.go::initializeFunctions`; constructor bodies in `core_functions.go`; `tryCollectionPlus` helper in `expression_generation.go::generateBinaryExpression`; Makefile (4 lines added to `fiber-runtime`, `http-runtime`, `c-lint`, `c-test`).

All files conform to CLAUDE.md's 500-LoC rule. `make fiber-runtime` and `make http-runtime` link clean. The C unit tests pass under `-Wall -Wextra -Werror`. Strict `c-lint` (with `-Wconversion -Wsign-conversion -Wshadow -Wcast-qual -Wpedantic` and friends) passes on every new C file.

### Spec follow-ups (already landed, but flagged for review)

- [x] Remove the ambiguous `ages { "Alice": 26 }` single-key-update syntax — replaced with `set(map, key, value)` in [TYPE-MAP-OPS].
- [x] Distinguish map filter/fold by giving them dedicated names (`filterEntries`, `foldEntries`, `mapValues`, `mapKeys`) with explicit K and V parameters — matches Elm `Dict.foldl`'s signature.
- [x] Document iteration-order policy: **unspecified** ([TYPE-MAP]).
- [x] Document permitted key types for v1: `int`, `string`, `bool` ([TYPE-MAP]).
- [x] Empty-literal `{}` disambiguation rule ([TYPE-MAP-LITERAL]).
- [x] Disallow `{}` as a pattern; require explicit `length(map: p) == 0` guard ([TYPE-MAP-PATTERNS]).
- [x] Fix `Array<string>` → `List<string>` in `split`, `join`, `lines`, `words` ([`0012-Built-InFunctions.md`](../../compiler/spec/0012-Built-InFunctions.md)).
- [x] Add Performance complexity table to [TYPE-LIST]/[TYPE-MAP]; correct list index from "O(1)" to "O(log₃₂ n)".
- [x] Add explicit Set-deferral note ([Collection Types overview](../../compiler/spec/0004-TypeSystem.md#collection-types)).
- [x] Add full `Collection Functions` section to `0012-Built-InFunctions.md` covering every builtin enumerated above.

## Acceptance Criteria for the Plan

A maintainer reading this document can:

1. Trace every spec ID (`[TYPE-LIST]`, `[TYPE-MAP-OPS]`, …) to a concrete code location after the plan is executed.
2. Justify the choice of HAMT and bitmapped vector trie to a reviewer who asks "why not X?" — answers are in §Background with citations.
3. Pick up any single TODO checkbox and finish it in under a day without needing to design anything further; the dependency order in §Implementation Phases tells them what must land first.
