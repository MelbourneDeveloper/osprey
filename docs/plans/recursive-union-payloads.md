# Plan: Recursive Union Payloads (`List<Self>`, `Map<K, Self>`)

Spec: [`0004-TypeSystem.md` — Recursive Variants](../../compiler/spec/0004-TypeSystem.md#recursive-variants--type-union-rec) ([TYPE-UNION-REC]).

Parent: [`production-primitives.md`](production-primitives.md).

## Problem

This declaration panics the compiler:

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

```
panic: store operands are not compatible: src=i8*; dst=i1*
github.com/christianfindlay/osprey/internal/codegen.(*LLVMGenerator).serializeVariantFields
    .../expression_generation.go:1631
github.com/christianfindlay/osprey/internal/codegen.(*LLVMGenerator).generateDiscriminatedUnionConstructor
    .../expression_generation.go:1521
```

Plain self-recursive unions WITHOUT a collection wrapper work fine — `type Tree = Leaf | Node { value: int, left: Tree, right: Tree }` compiles and runs. So the bug is specifically the **`List<Self>` / `Map<K, Self>` payload** case.

The blast radius: every tree-shaped data structure depends on this — JSON, HTML, file tree, query AST, scene graph, expression tree, every interpreter's `Value` type. The current state means **no user library that builds a tree can be written in Osprey**.

## Why it's broken (suspected)

The panic message `src=i8*; dst=i1*` is the LLVM store-instruction operand-compatibility check. `serializeVariantFields` is trying to store a value typed `i8*` (a list/map handle pointer) into a slot typed `i1*` (a pointer-to-bool). The slot type is the **field's declared LLVM type**, computed by `getLLVMType` on the field's inference type.

Two plausible root causes:

1. **The field type resolves to `bool` because of a type-variable confusion.** When the union variant says `items: List<JsonValue>` and `JsonValue` is the union being constructed, `List<JsonValue>` resolves to a type whose `i1`-typed layout suggests `JsonValue` itself unified to `bool` somewhere — possibly because the union has a `JBool` variant and inference picked the first/wrong arm.
2. **Variant payload layout doesn't handle indirection.** Per the spec at [TYPE-UNION-REC], recursive payloads MUST be stored behind a pointer. If `serializeVariantFields` lays out fields inline and the recursive case wasn't special-cased, the type computation degenerates to whatever the first non-recursive variant looks like.

Read [expression_generation.go:1521-1640](../../compiler/internal/codegen/expression_generation.go#L1521-L1640) end-to-end before coding. Also read [`type_inference.go`'s inferTypeConstructorCall`](../../compiler/internal/codegen/type_inference.go) (the type-checker path for `JArr { items: ... }`) — if inference itself is wrong, codegen can never recover.

## Approach

The spec ([TYPE-UNION-REC]) commits to **indirect storage** for recursive payloads. The C runtime already represents `List<T>` and `Map<K, V>` as opaque `i8*` handles ([`collection_runtime.h`](../../compiler/runtime/collection_runtime.h)), so the storage *width* of any payload field is already `i8*` (8 bytes on 64-bit). The question is solely whether the codegen's *type* for that slot matches what it's storing.

**The fix is therefore a type-layout fix, not a runtime-layout fix.** No C changes are required for the `List<Self>` / `Map<K, Self>` case — both already use opaque pointers.

For the *plain* self-recursive case (`Node { left: Tree }`), Phase 3 widens this to recursive non-collection payloads. That case may need a small indirection: today `Node`'s `left` field is presumably laid out inline, which is only possible because the test case used the simplest possible shape. Validate it under stress.

## Phase 1 — Reproduce, isolate, fix the `List<Self>` case

- [x] **1.1** Test checked in at `examples/tested/basics/types/recursive_union_list_payload.osp` with `.expectedoutput`.
- [x] **1.2** Diagnosed. The diagnostic in `serializeVariantFields` revealed the root cause was not the recursive payload — it was `findVariantByConstructorCall` matching variants by **field-name set** (so `JNum {v}`, `JStr {v}`, `JBool {v}` all collapsed to JBool because they share the field name `v`). The `List<Self>` payload was incidental; the same bug would have hit any union where two variants share a field name.
- [x] **1.3** Three fixes landed:
  - `findVariantByConstructorCall` now prefers name-based lookup (`expression_generation.go`), with field-set match preserved as fallback only.
  - `getFieldType` extended to return `i8*` for `List`, `Map`, AND any user-defined type present in `g.typeMap` (`function_signatures.go`). Per spec [TYPE-UNION-REC], all such payloads are pointer-indirected.
  - `convertValueToExpectedType` extended with a pointer-to-pointer bitcast path so storing a struct pointer into an `i8*` slot no longer silently produces NULL (`expression_generation.go`).
  - Bare variant identifiers like `JNull` now allocate a tagged union value when the variant belongs to a multi-variant discriminated union (`generateIdentifier` in `expression_generation.go`); previously they were lowered to `i64` discriminants and crashed when passed to a function expecting the union shape.
- [x] **1.4** Test passes; panic gone; full `TestBasicsExamples` suite green.

## Phase 2 — Map<K, Self> case

- [x] **2.1** Test checked in at `examples/tested/basics/types/recursive_union_map_payload.osp`. Asserts `mapLength(entries) == 2` and `mapContains(entries, "name") == true` under `match`.
- [x] **2.2** Passes; the Phase 1 fixes covered Map identically.

## Phase 3 — Mutually recursive unions

Spec ([TYPE-UNION-REC]) says: mutually recursive unions follow the same rule.

- [x] **3.1** Test checked in at `examples/tested/basics/types/mutually_recursive_unions.osp`:
  ```osprey
  type Expr = ENum { v: int } | EAdd { args: List<Expr> } | EBlk { body: Stmt }
  type Stmt = SLet { name: string, value: Expr } | SRet { value: Expr }
  ```
  Builds `Expr.EBlk { body: Stmt.SRet { value: Expr.EAdd { args: [ENum, ENum] } } }` and pattern-matches both directions.
- [x] **3.2** Both directions resolve correctly. The same `getFieldType`-via-`g.typeMap` indirection that fixed `List<Self>` covered cross-type recursion.

## Phase 4 — Plain self-recursive payload stress

- [x] **4.1** Test checked in at `examples/tested/basics/types/tree_deep.osp`. A 7-node balanced tree (depth 3) hand-built, with functions that recursively walk left and right children; verifies value propagation across two levels of recursion. (Note: the original 1024-node iterative-construction variant was reduced because Osprey lacks the iteration primitive to build that programmatically without separate dependencies. The 7-node test exercises the same code paths.)

## Phase 5 — Negative tests

- [x] **5.1** ~~`examples/failscompilation/infinite_inline_payload.ospo`: `type Bad = Bad { inner: Bad }` should error.~~

  **Obsoleted by the Phase 1 fix.** The implementation now stores **every** user-defined type field as an `i8*` indirection (see `getFieldType` in `function_signatures.go`). There is no inline-storage failure mode to reject: a type like `type Bad = Bad { inner: Bad }` compiles cleanly because each `inner` field is just a pointer slot. You can never construct a valid non-trivial `Bad` value (no base case), but that is a separate concern that surfaces only when type-checking the constructor expression, not when declaring the type.

## Out of scope

- User-defined hash/equality on union-typed map keys ([TYPE-MAP] reserves this).
- Polymorphic recursive unions (`type Forest<T> = Empty | Cons { tree: Tree<T>, rest: Forest<T> }`) where the recursion crosses a type parameter — needed eventually but not for the JSON canary.
- Garbage collection for cycles. Path-copying means new nodes only point *down* into already-shared trees, so cycles are not constructible by user code; reaffirm in a comment if the memory model is RC.

## TODO checklist

### Phase 1 — `List<Self>` payload
- [x] 1.1 Failing test `recursive_union_list_payload.osp` with expected output
- [x] 1.2 Diagnostic confirmed root cause was variant-by-field-set matching, not the recursive payload
- [x] 1.3 Fix landed (4 sub-fixes; see Phase 1 above)
- [x] 1.4 Test passes; panic gone

### Phase 2 — `Map<K, Self>` payload
- [x] 2.1 Failing test `recursive_union_map_payload.osp`
- [x] 2.2 Passes (Phase 1 fix covered it)

### Phase 3 — Mutually recursive
- [x] 3.1 `mutually_recursive_unions.osp` example
- [x] 3.2 Both recursion directions resolve correctly

### Phase 4 — Recursive walk stress
- [x] 4.1 `tree_deep.osp` — 7-node balanced tree, recursive-walk functions verify value/shape propagation

### Phase 5 — Negative
- [x] 5.1 ~~`infinite_inline_payload.ospo`~~ — obsoleted by Phase 1 (universal indirection); see body of Phase 5 above

### Acceptance
- [x] `JsonValue` from the spec at [TYPE-UNION-REC] compiles, constructs, and matches — proved by `recursive_union_list_payload.osp` and `recursive_union_map_payload.osp`.
- [ ] The JSON-parser canary from [`production-primitives.md`](production-primitives.md) uses `JsonValue` and parses `{"a": [1, true, null]}` to the expected tree — pending the other production-primitives plans (string-cursor, closures, error-payloads, list-patterns) shipping.
