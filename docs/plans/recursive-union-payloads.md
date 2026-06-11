# Plan: Recursive Union Payloads (`List<Self>`, `Map<K, Self>`)

Spec: [`0004-TypeSystem.md` — Recursive Variants](../specs/0004-TypeSystem.md#recursive-variants--type-union-rec) ([TYPE-UNION-REC]).

Parent: [`production-primitives.md`](production-primitives.md).

## Problem (historical — fixed; see Status below)

This declaration used to panic the compiler in the variant-field serializer with
`store operands are not compatible: src=i8*; dst=i1*`:

```osprey
type JsonValue =
    JNull
    | JBool { v: bool }
    | JNum  { v: int }
    | JStr  { v: string }
    | JArr  { items:   List<JsonValue> }
    | JObj  { entries: Map<string, JsonValue> }

let arr = JArr { items: List() |> listAppend(JNum { v: 42 }) }
```

Plain self-recursive unions WITHOUT a collection wrapper work fine — `type Tree = Leaf | Node { value: int, left: Tree, right: Tree }` compiles and runs. So the bug is specifically the **`List<Self>` / `Map<K, Self>` payload** case.

The blast radius: every tree-shaped data structure depends on this — JSON, HTML, file tree, query AST, scene graph, expression tree, every interpreter's `Value` type. The current state means **no user library that builds a tree can be written in Osprey**.

## Why it was broken (diagnosed)

The panic message `src=i8*; dst=i1*` was the LLVM store-instruction operand-compatibility check: the
variant-field serializer stored a value typed `i8*` (a list/map handle pointer) into a slot typed `i1*`
(pointer-to-bool). Root cause (confirmed in Phase 1.2): variants were matched by **field-name set**, so
`JNum {v}` / `JStr {v}` / `JBool {v}` collapsed to whichever variant shared the field name — the recursive
payload was incidental.

In the current compiler the equivalents live in [`crates/osprey-codegen/src/aggregate.rs`](../../crates/osprey-codegen/src/aggregate.rs)
(variant construction — tagged heap blocks `{ i64 tag, fields… }`, selected by constructor *name*) and the
field-type mapping in [`crates/osprey-codegen/src/types.rs`](../../crates/osprey-codegen/src/types.rs)
(user-type and collection payloads are pointer-indirected per [TYPE-UNION-REC]). Constructor field types
come straight from inference ([`crates/osprey-types`](../../crates/osprey-types/src/)) via `ProgramTypes` —
there is no second name-based type computation to disagree with it.

## Approach

The spec ([TYPE-UNION-REC]) commits to **indirect storage** for recursive payloads. The C runtime already represents `List<T>` and `Map<K, V>` as opaque `i8*` handles ([`collection_runtime.h`](../../compiler/runtime/collection_runtime.h)), so the storage *width* of any payload field is already `i8*` (8 bytes on 64-bit). The question is solely whether the codegen's *type* for that slot matches what it's storing.

**The fix is therefore a type-layout fix, not a runtime-layout fix.** No C changes are required for the `List<Self>` / `Map<K, Self>` case — both already use opaque pointers.

For the *plain* self-recursive case (`Node { left: Tree }`), Phase 3 widens this to recursive non-collection payloads. That case may need a small indirection: today `Node`'s `left` field is presumably laid out inline, which is only possible because the test case used the simplest possible shape. Validate it under stress.

## Phase 1 — Reproduce, isolate, fix the `List<Self>` case

- [x] **1.1** Test checked in — now consolidated into [`examples/tested/basics/types/recursive_unions.osp`](../../compiler/examples/tested/basics/types/recursive_unions.osp) (covers all four payload shapes) with `.expectedoutput`.
- [x] **1.2** Diagnosed. The diagnostic in `serializeVariantFields` revealed the root cause was not the recursive payload — it was `findVariantByConstructorCall` matching variants by **field-name set** (so `JNum {v}`, `JStr {v}`, `JBool {v}` all collapsed to JBool because they share the field name `v`). The `List<Self>` payload was incidental; the same bug would have hit any union where two variants share a field name.
- [x] **1.3** Fixes landed (carried into the Rust port):
  - Variant lookup is name-based — the constructor name selects the variant, never a field-set match
    (`crates/osprey-codegen/src/aggregate.rs`).
  - `List`, `Map`, and user-defined-type payload fields are `i8*` pointer-indirected
    (`crates/osprey-codegen/src/types.rs`). Per spec [TYPE-UNION-REC].
  - Pointer-to-pointer casts on store, so a struct pointer goes into an `i8*` slot without silently
    producing NULL (`crates/osprey-codegen/src/cast.rs`).
  - Bare nullary variant identifiers like `JNull` lower to tagged union values when the variant belongs
    to a multi-variant union, not raw `i64` discriminants.
- [x] **1.4** Test passes; panic gone; the differential harness (`crates/diff_examples.sh`) is green.

## Phase 2 — Map<K, Self> case

- [x] **2.1** Test checked in (consolidated into `recursive_unions.osp`, the `JObj { entries: Map<string, JsonValue> }` shape).
- [x] **2.2** Passes; the Phase 1 fixes covered Map identically.

## Phase 3 — Mutually recursive unions

Spec ([TYPE-UNION-REC]) says: mutually recursive unions follow the same rule.

- [x] **3.1** Test checked in (consolidated into `recursive_unions.osp`, the `Expr <-> Stmt` shape):
  ```osprey
  type Expr = ENum { v: int } | EAdd { args: List<Expr> } | EBlk { body: Stmt }
  type Stmt = SLet { name: string, value: Expr } | SRet { value: Expr }
  ```
  Builds `Expr.EBlk { body: Stmt.SRet { value: Expr.EAdd { args: [ENum, ENum] } } }` and pattern-matches both directions.
- [x] **3.2** Both directions resolve correctly. The same `getFieldType`-via-`g.typeMap` indirection that fixed `List<Self>` covered cross-type recursion.

## Phase 4 — Plain self-recursive payload stress

- [x] **4.1** Test checked in (consolidated into `recursive_unions.osp`, the direct-`Self` `Tree` shape). A balanced tree hand-built, with functions that recursively walk left and right children; verifies value propagation across recursion levels. (Note: the original 1024-node iterative-construction variant was reduced because Osprey lacks the iteration primitive to build that programmatically without separate dependencies; the small tree exercises the same code paths.)

## Phase 5 — Negative tests

- [x] **5.1** ~~`examples/failscompilation/infinite_inline_payload.ospo`: `type Bad = Bad { inner: Bad }` should error.~~

  **Obsoleted by the Phase 1 fix.** The implementation now stores **every** user-defined type field as an `i8*` indirection (field-type mapping in [`crates/osprey-codegen/src/types.rs`](../../crates/osprey-codegen/src/types.rs)). There is no inline-storage failure mode to reject: a type like `type Bad = Bad { inner: Bad }` compiles cleanly because each `inner` field is just a pointer slot. You can never construct a valid non-trivial `Bad` value (no base case), but that is a separate concern that surfaces only when type-checking the constructor expression, not when declaring the type.

## Out of scope

- User-defined hash/equality on union-typed map keys ([TYPE-MAP] reserves this).
- Polymorphic recursive unions (`type Forest<T> = Empty | Cons { tree: Tree<T>, rest: Forest<T> }`) where the recursion crosses a type parameter — needed eventually but not for the JSON canary.
- Garbage collection for cycles. Path-copying means new nodes only point *down* into already-shared trees, so cycles are not constructible by user code; reaffirm in a comment if the memory model is RC.

## TODO checklist

### Phase 1 — `List<Self>` payload
- [x] 1.1 Failing test (now in `recursive_unions.osp`) with expected output
- [x] 1.2 Diagnostic confirmed root cause was variant-by-field-set matching, not the recursive payload
- [x] 1.3 Fix landed (4 sub-fixes; see Phase 1 above)
- [x] 1.4 Test passes; panic gone

### Phase 2 — `Map<K, Self>` payload
- [x] 2.1 Failing test (now in `recursive_unions.osp`)
- [x] 2.2 Passes (Phase 1 fix covered it)

### Phase 3 — Mutually recursive
- [x] 3.1 `Expr <-> Stmt` example (now in `recursive_unions.osp`)
- [x] 3.2 Both recursion directions resolve correctly

### Phase 4 — Recursive walk stress
- [x] 4.1 Balanced `Tree` walk (now in `recursive_unions.osp`) — recursive-walk functions verify value/shape propagation

### Phase 5 — Negative
- [x] 5.1 ~~`infinite_inline_payload.ospo`~~ — obsoleted by Phase 1 (universal indirection); see body of Phase 5 above

### Acceptance
- [x] `JsonValue` from the spec at [TYPE-UNION-REC] compiles, constructs, and matches — proved by `recursive_unions.osp`.
- [ ] The JSON-parser canary from [`production-primitives.md`](production-primitives.md) uses `JsonValue` and parses `{"a": [1, true, null]}` to the expected tree — pending the other production-primitives plans (string-cursor, closures, error-payloads, list-patterns) shipping.
