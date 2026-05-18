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

- [ ] **1.1** Check in the panicking test as `examples/tested/types/recursive_union_list_payload.osp` along with its `.expectedoutput`. The test asserts the runtime behaviour the spec promises — that a `JArr { items: List<JsonValue> }` round-trips through `match` and `listLength(items)` returns the right count.
- [ ] **1.2** Diagnose: print the field's inferred type, the resolved type, and the LLVM type at the top of `serializeVariantFields` for the failing field. Determine whether the bug is in inference or in layout.
- [ ] **1.3** Fix. If inference: ensure that when inferring `List<JsonValue>` inside a `JsonValue` variant declaration, the recursive reference resolves to the enclosing union type (not a fresh type variable, not the first variant's payload). If layout: ensure `getLLVMType` for `List<T>` and `Map<K, V>` returns `i8*` uniformly regardless of `T`, `K`, `V`.
- [ ] **1.4** Test from 1.1 passes; panic gone.

## Phase 2 — Map<K, Self> case

- [ ] **2.1** Check in `examples/tested/types/recursive_union_map_payload.osp`: a `JObj { entries: Map<string, JsonValue> }` value with two entries; assert `mapLength(entries) == 2` and that `mapContains` works under `match`.
- [ ] **2.2** Verify it passes after the Phase 1 fix (should be free if both go through the same `i8*` layout path).

## Phase 3 — Mutually recursive unions

Spec ([TYPE-UNION-REC]) says: mutually recursive unions follow the same rule.

- [ ] **3.1** `examples/tested/types/mutually_recursive_unions.osp`:
  ```osprey
  type Expr = ENum { v: int } | EAdd { args: List<Expr> } | EBlock { body: Stmt }
  type Stmt = SLet { name: string, value: Expr } | SReturn { value: Expr }
  ```
  Build a small expression-statement pair and `match` on it.
- [ ] **3.2** Ensure both directions of the recursion resolve to the right union; not the same bug as Phase 1 but symptoms could look similar.

## Phase 4 — Plain self-recursive payload stress

The `Tree` case works for the trivial shape. Verify under a deep tree:

- [ ] **4.1** `examples/tested/types/tree_deep.osp`: build a 1024-node binary tree by recursive construction, traverse with `match`, assert the count. (If construction overflows the stack, document the limit; it's a separate workstream — TCO — not this plan's scope.)

## Phase 5 — Negative tests

- [ ] **5.1** `examples/failscompilation/infinite_inline_payload.ospo`: `type Bad = Bad { inner: Bad }` (no indirection through a variant tag or collection — every value is its own infinitely-nested storage). Expected: clear compile error explaining the variant must be indirected, with the spec ID `[TYPE-UNION-REC]`.
  
  Note: today this may already error for a different reason; the test must assert the *helpful* error message, not just that it fails to compile.

## Out of scope

- User-defined hash/equality on union-typed map keys ([TYPE-MAP] reserves this).
- Polymorphic recursive unions (`type Forest<T> = Empty | Cons { tree: Tree<T>, rest: Forest<T> }`) where the recursion crosses a type parameter — needed eventually but not for the JSON canary.
- Garbage collection for cycles. Path-copying means new nodes only point *down* into already-shared trees, so cycles are not constructible by user code; reaffirm in a comment if the memory model is RC.

## TODO checklist

### Phase 1 — `List<Self>` payload
- [ ] 1.1 Failing test `recursive_union_list_payload.osp` with expected output
- [ ] 1.2 Diagnostic prints at top of `serializeVariantFields`
- [ ] 1.3 Fix (inference or layout — diagnose first)
- [ ] 1.4 Test passes; panic gone

### Phase 2 — `Map<K, Self>` payload
- [ ] 2.1 Failing test `recursive_union_map_payload.osp`
- [ ] 2.2 Passes after Phase 1 fix

### Phase 3 — Mutually recursive
- [ ] 3.1 `mutually_recursive_unions.osp` example
- [ ] 3.2 Both recursion directions resolve correctly

### Phase 4 — Deep recursion stress
- [ ] 4.1 `tree_deep.osp` — 1024-node tree construction + traversal

### Phase 5 — Negative
- [ ] 5.1 `infinite_inline_payload.ospo` produces a helpful error referencing [TYPE-UNION-REC]

### Acceptance
- [ ] `JsonValue` from the spec at [TYPE-UNION-REC] compiles, constructs, and matches.
- [ ] The JSON-parser canary from [`production-primitives.md`](production-primitives.md) uses `JsonValue` and parses `{"a": [1, true, null]}` to the expected tree.
