# Plan: List Patterns in `match`

Spec: [`0004-TypeSystem.md` — Patterns (TYPE-LIST-PATTERNS)](../specs/0004-TypeSystem.md#patterns--type-list-patterns), [`0007-PatternMatching.md`](../specs/0007-PatternMatching.md).

Parent: [`production-primitives.md`](production-primitives.md).

## Problem

The spec at [TYPE-LIST-PATTERNS](../specs/0004-TypeSystem.md#patterns--type-list-patterns) advertises four list-pattern forms:

```osprey
fn classify(xs: List<int>) -> string = match xs {
    []                 => "empty"
    [single]           => "one"
    [first, second]    => "two"
    [head, ...tail]    => "many starting with ${head}"
}
```

None of them are implemented. The grammar accepts `[head, ...tail]` (per the deleted `collections.md` plan's Phase 5.1 entry) but there's no AST node, no builder support, and no match-codegen for any of the four forms.

Without list patterns, recursive-descent parsers — every JSON, query-string, header, CSV, or markdown reader — must use the existing `listLength` + `osprey_list_get` pattern. That's verbose enough (~3× more code, by inspection of typical Haskell vs. for-loop equivalents) that the "build it in Osprey" promise reads as bait-and-switch. Escalated from `deferred` to **critical-path** by [`production-primitives.md`](production-primitives.md).

## Scope (four pattern forms)

| Pattern | Semantics |
|---|---|
| `[]` | Matches iff `osprey_list_length(xs) == 0`. Binds nothing. |
| `[x, y, z]` | Matches iff `osprey_list_length(xs) == 3`. Binds positions to the named bindings. |
| `[head, ...tail]` | Matches iff `osprey_list_length(xs) >= 1`. Binds `head` to position 0; binds `tail` to a sub-list starting at position 1. |
| `[a, b, ...rest]` | Generalisation: matches iff length ≥ 2; binds two heads + a rest list. Any prefix length ≥ 0. |

Map patterns (`{ "key": binding }`) are a separate workstream — out of scope here. They were Phase 5.4 in the deleted collections plan and remain deferred.

## Phase 1 — Grammar + AST

- [ ] **1.1** Extend the `pattern` rule in [`tree-sitter-osprey/grammar.js`](../../tree-sitter-osprey/grammar.js) with a list-pattern alternative:
  ```
  pattern : '[' (pattern (',' pattern)*)? ( ',' '...' ID )? ']'
          | <existing alternatives>
          ;
  ```
  The `...` rest binder is mandatory at the tail position only; arbitrary middle ellipsis is **not** allowed (matches Haskell and Elm; Scala-style mid-list patterns are too expressive for this iteration).
- [ ] **1.2** Regenerate the parser (`tree-sitter generate` in `tree-sitter-osprey/`).
- [ ] **1.3** Add the AST shape in [`crates/osprey-ast`](../../crates/osprey-ast/src/):
  ```rust
  ListPattern {
      elements: Vec<Pattern>,   // bindings for the fixed-prefix positions
      rest: Option<String>,     // None if no `...rest`; otherwise the binder name
  }
  ```
  Both `[]` and `[x, y]` use this same node — `elements` empty / non-empty, `rest` `None`.
- [ ] **1.4** Update the match-arm builder in [`crates/osprey-syntax`](../../crates/osprey-syntax/src/) to construct `ListPattern` from the tree-sitter parse tree, walking children in source order so the rest binder lands correctly.

## Phase 2 — Match codegen

Extend the match lowering in [`crates/osprey-codegen/src/pattern.rs`](../../crates/osprey-codegen/src/pattern.rs) to cover all four forms. The generation pattern, for each pattern form, is a length-guard branch then per-position binding:

- [ ] **2.1** **Empty pattern `[]`**: emit `icmp eq i64 (osprey_list_length xs), 0` → conditional branch to the arm body if true, fallthrough to next arm if false.
- [ ] **2.2** **Fixed-length pattern `[x, y, z]`**: emit `icmp eq i64 (osprey_list_length xs), 3`. On true, for each binding, call `osprey_list_get(xs, i)` and bind the unwrapped value (length guard proves bounds, so we can use `osprey_list_get_unchecked` if one exists — otherwise call `osprey_list_get` and assume Success).
- [ ] **2.3** **Head/tail pattern `[head, ...tail]`**: emit `icmp uge i64 (osprey_list_length xs), 1`. On true, bind `head` via `osprey_list_get(xs, 0)`; bind `tail` via a new runtime call `osprey_list_drop(xs, 1)` (see Phase 3).
- [ ] **2.4** **Prefix + rest pattern `[a, b, ...rest]`**: length guard `>= 2`; positional binds for `a` and `b`; `rest = osprey_list_drop(xs, 2)`.
- [ ] **2.5** Integrate with the existing `match` codegen so list patterns coexist with union-variant patterns, literal patterns, and the catch-all `_`. Each arm's predicate is OR'd into the discriminator chain.

## Phase 3 — Runtime helper: `osprey_list_drop`

The list runtime ([`compiler/runtime/list_runtime.c`](../../compiler/runtime/list_runtime.c)) currently exposes `osprey_list_get`, `osprey_list_append`, `osprey_list_concat`, `osprey_list_length`, iter helpers, and a transient builder — but no `drop(n)`.

- [ ] **3.1** Add `void *osprey_list_drop(void *list, int64_t n)` to `list_runtime.c`: returns a new persistent list containing elements `[n, length)`. Persistent / structural-sharing — must NOT copy elements. Implementation: if the underlying structure is a bitmapped vector trie, the slice can usually share the suffix subtree and only re-create the path to the new root.
- [ ] **3.2** `osprey_list_drop(xs, n)` with `n >= length` returns the empty list (matches Haskell `drop`).
- [ ] **3.3** `osprey_list_drop(xs, n)` with `n < 0` MAY be treated as `n = 0` (no error path — the length-guard in codegen prevents that input anyway).
- [ ] **3.4** Declare the extern in [`crates/osprey-codegen/src/collections.rs`](../../crates/osprey-codegen/src/collections.rs) (with the other `osprey_list_*` decls).
- [ ] **3.5** C unit test in [`runtime/list_tests.c`](../../compiler/runtime/list_tests.c):
  - Empty → drop(0) returns empty; drop(5) returns empty.
  - 10-element list, drop(0) returns the same list; drop(10) returns empty; drop(3) returns the suffix `[3..10)`.
  - **Persistence**: drop(3) of a list does not mutate the original — original is still 10 elements after the call.

## Phase 4 — Tested examples

- [ ] **4.1** `examples/tested/basics/lists/list_patterns.osp` covering the four forms:
  ```osprey
  fn classify(xs) = match xs {
      []              => "empty"
      [single]        => "one"
      [first, second] => "two"
      [head, ...tail] => "many starting with ${head}"
  }
  print(classify(List()))                              // "empty"
  print(classify(List() |> listAppend(1)))             // "one"
  print(classify(List() |> listAppend(1) |> listAppend(2)))            // "two"
  print(classify(List() |> listAppend(1) |> listAppend(2) |> listAppend(3)))   // "many starting with 1"
  ```
  Pin the four output lines in `.expectedoutput`.
- [ ] **4.2** `examples/tested/basics/lists/list_pattern_recursion.osp` — recursive `sum` using head/tail pattern:
  ```osprey
  fn sumList(xs) = match xs {
      []              => 0
      [head, ...tail] => head + sumList(tail)
  }
  print(sumList(List() |> listAppend(1) |> listAppend(2) |> listAppend(3)))   // 6
  ```
  This is the test that proves the pattern is usable for recursive descent (and by extension, for the JSON parser canary).

## Phase 5 — Negative tests

- [ ] **5.1** `examples/failscompilation/list_pattern_middle_rest.ospo` — `[a, ...mid, b]` — expected: clear error that the rest binder must be at the tail.
- [ ] **5.2** `examples/failscompilation/list_pattern_double_rest.ospo` — `[...a, ...b]` — expected: at most one rest binder.

## Out of scope

- Map patterns (`{ "key": binding }`) — separate workstream; the spec at [TYPE-MAP-PATTERNS](../specs/0004-TypeSystem.md#patterns--type-map-patterns) describes them, but they aren't blocking the JSON parser canary.
- List comprehensions (`[x*x for x in xs]`) — separate workstream; spec at [TYPE-LIST-COMP](../specs/0004-TypeSystem.md#comprehensions--type-list-comp).
- Nested-pattern destructuring (`[Some { v: x }, ...]`) — separate workstream; the four forms above all bind via plain identifier patterns at the element positions in this iteration.

## TODO checklist

### Phase 1 — Grammar + AST
- [ ] 1.1 Extend the `pattern` rule in `tree-sitter-osprey/grammar.js`
- [ ] 1.2 `tree-sitter generate`
- [ ] 1.3 `ListPattern` AST node in `crates/osprey-ast`
- [ ] 1.4 Builder support in `crates/osprey-syntax`

### Phase 2 — Match codegen
- [ ] 2.1 Empty-list pattern
- [ ] 2.2 Fixed-length pattern
- [ ] 2.3 Head/tail pattern
- [ ] 2.4 Prefix + rest pattern
- [ ] 2.5 Integration with existing match codegen

### Phase 3 — `osprey_list_drop`
- [ ] 3.1 Implement in `list_runtime.c` (structural-sharing slice)
- [ ] 3.2 `n >= length` → empty list
- [ ] 3.3 Negative `n` → no-op (defensive)
- [ ] 3.4 Declare extern in `collections.rs`
- [ ] 3.5 C unit tests in `list_tests.c` covering all branches + persistence invariant

### Phase 4 — Tested examples
- [ ] 4.1 `list_patterns.osp` — the four forms
- [ ] 4.2 `list_pattern_recursion.osp` — recursive `sumList`

### Phase 5 — Negative
- [ ] 5.1 `list_pattern_middle_rest.ospo`
- [ ] 5.2 `list_pattern_double_rest.ospo`

### Acceptance
- [ ] The JSON-parser canary in [`production-primitives.md`](production-primitives.md) uses `[head, ...tail]` in at least one place.
- [ ] `sumList` recursion runs over a 1000-element list without stack overflow (or, if it does overflow, document the recursion-depth limit and open a TCO follow-up issue).
