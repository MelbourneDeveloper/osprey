# Plan: Thread Real Error Payloads Through `Result<T, E>`

Spec: [`0013-ErrorHandling.md` — Error Payload Propagation](../../compiler/spec/0013-ErrorHandling.md#error-payload-propagation--err-payload) ([ERR-PAYLOAD]).

Parent: [`production-primitives.md`](production-primitives.md).

## Problem

Every `Error { message }` branch in user code is bound to the same hardcoded global string `"Error occurred"`, regardless of which builtin failed or what reason the runtime had.

Source: [`internal/codegen/llvm.go:2305`](../../compiler/internal/codegen/llvm.go#L2305):

```go
errorStr := g.module.NewGlobalDef(
    "error_msg"+blockSuffix,
    constant.NewCharArrayFromString("Error occurred\\x00"),
)
errorPtr := g.builder.NewGetElementPtr(...)
g.variables[fieldName] = errorPtr
```

This is acknowledged in [`docs/plans/string-manipulation.md:241`](string-manipulation.md#L241):

> Fallible builtins set discriminant=1 and null value; the match-expression `Error { message }` branch always binds the same static `"Error occurred\x00"` global string regardless of which builtin failed.

A parser that cannot say *what went wrong* and *where* is unshippable. Same for every HTTP handler, file I/O failure, parse failure, division by zero. The spec's [ERR-PAYLOAD] section, added alongside this plan, makes the current implementation explicitly non-conforming.

## What needs to change

`Result<T, E>` at runtime is a two-field layout: an `i64` discriminant (0 = Success, 1 = Error) and a payload. Today the payload slot for the error case is unused — every fallible runtime function returns a Result whose error slot is null/garbage, and codegen "recovers" by binding the global stand-in.

The fix has three pieces:

1. **Define the runtime payload contract.** For `E = string`, the payload slot holds a `char*` to the error message (heap-allocated by the runtime, or a static string constant — either is fine as long as it's read-only and null-terminated). For `E` = a discriminated union like `StringError`, the payload slot holds a pointer to the union value (same layout as any other discriminated union).
2. **Make every runtime function that produces an Error actually populate the slot.** Today many runtime functions return what is effectively `Error { message: nullptr }`. Audit and fix.
3. **Make codegen read the slot instead of generating a fresh global.** Replace [llvm.go:2298-2308](../../compiler/internal/codegen/llvm.go#L2298-L2308) with a load from the matched expression's Result struct.

## Phase 1 — Runtime contract for `Result<T, string>`

**Implementation chose codegen-side messages over a runtime contract.** The C
runtime never produced `Result` structs directly (it returned raw `char*` /
status `int64_t`); codegen wraps every return. So instead of a new
`runtime/result_runtime.c`, the wrapping sites in
[`internal/codegen/`](../../compiler/internal/codegen/) intern function-specific
static messages and store them in the new err_msg slot. This keeps the C ABI
unchanged and avoids a parallel struct layout in `string_runtime.h`.

- [x] **1.1** Layout decision shipped: extended `getResultType` in
  [`core_functions.go:528`](../../compiler/internal/codegen/core_functions.go#L528)
  from `{value, discriminant}` to `{value, discriminant, err_msg: i8*}`.
  `ResultFieldCount = 3` in [`constants.go`](../../compiler/internal/codegen/constants.go).
  `UnionFieldCount = 2` introduced alongside to keep tagged-union dispatch
  separate from Result detection (a previous shared `len==2` check silently
  routed unions through `createSimpleEnumCondition` mid-implementation —
  fixed once observed; covered by `recursive_union_list_payload.osp` and
  friends).
- [x] **1.2** Helpers shipped in
  [`internal/codegen/result_helpers.go`](../../compiler/internal/codegen/result_helpers.go):
  `makeSuccessValue`, `makeErrorValueWithMessage`, `storeResultFields`,
  `loadResultErrorMessage`, `internErrorMessage`. The intern table reuses
  the previously-unused `stringConstants` map on `LLVMGenerator` so repeated
  messages share one global (`@osp_err_msg_N`).
- [x] **1.3** Builtins audited and updated (per-callsite, in codegen — not in C):
  - `parseInt` → `"parseInt: input is not a valid integer"`
  - `parseFloat` → `"parseFloat: input is not a valid number"`
  - `substring` → `"substring: index out of range"`
  - `split` → `"split: separator must not be empty"`
  - `replace` → `"replace: needle must not be empty"`
  - `repeat` → `"repeat: count must be non-negative"`
  - `padStart` / `padEnd` → `"<name>: fill must not be empty"`
  - `indexOf` → `"indexOf: needle not found"`
  - `input` → `"input: failed to read line"`
  - `readFile` / `writeFile` → `"<name>: failed to <op> file"`
  - `spawnProcess` → `"spawnProcess: failed to spawn process"`
  - list `[idx]` out-of-bounds → `"list: index out of range"`
  - map `[key]` missing / empty → `"map: key not found"` / `"map: lookup on empty map"`
  - integer / float division by zero → `"math: division by zero"`
- [x] **1.4** All shipped messages are .rodata static strings. The intern
  helper marks each global as `Immutable`. No `strdup`/heap path is needed
  for any current builtin.

## Phase 2 — Codegen reads the slot

- [x] **2.1** Failing test landed at
  [`examples/tested/basics/types/error_payload.osp`](../../compiler/examples/tested/basics/types/error_payload.osp)
  exercising `parseInt("oops")` and `substring("abc", 5, 10)`. The
  `.expectedoutput` pins the exact text per the messages in 1.3.
- [x] **2.2** Hardcoded global removed from
  [`llvm.go`](../../compiler/internal/codegen/llvm.go) `generateErrorBlock`.
  Replaced with `g.variables[fieldName] = g.loadResultErrorMessage(g.currentResultValue)`.
- [x] **2.3** Loaded pointer typed as `i8*` (matches the new slot 2 layout).
- [x] **2.4** Test passes; same fix incidentally made `feature_omnibus2`,
  `file_io_json_workflow`, `pattern_matching_result_tests`, and
  `string_utils_combined` report real messages — their `.expectedoutput`
  files were updated to assert the new text.

### Implementation note — nested Result matches

`g.currentResultValue` is a single field on the generator. When a Result match's
success arm contains another Result match, the inner match overwrites
`currentResultValue` and the *outer* error block was then loading from a
non-dominating sibling block (`llc` rejected with "Instruction does not
dominate all uses"). Fixed in `generateResultMatchExpression` by
save/restore around success-arm generation; covered by
`script_style_working.osp`'s nested `factorial`.

## Phase 3 — Discriminated-union payload types (e.g. `StringError`)

Spec defines `StringError` as a discriminated union ([0012-Built-InFunctions.md:67-73](../../compiler/spec/0012-Built-InFunctions.md#L67-L73)). Currently `Error { message }` always binds a `string`, never a `StringError` variant. Once the spec is taken literally, fallible string functions return `Result<T, StringError>` and the match arm binds an actual `StringError` value that the caller pattern-matches:

```osprey
match split(s, "") {
    Success { value } => ...
    Error   { message: InvalidArgument { message: m } } => print("bad arg: ${m}")
    Error   { message: IndexOutOfRange { index, length } } => print("oor: ${index} of ${length}")
    Error   { message: NotFound } => print("not found")
    Error   { message: ParseFailed { input } } => print("parse: ${input}")
}
```

- [ ] **3.1** Decide: ship Phase 3 in this plan, or defer? **Recommend defer** — Phase 1+2 unblocks JSON-parser canary with `Result<T, string>`. Phase 3 needs union-type runtime layout work; the recursive-union prerequisite shipped in PR #67 (see [TYPE-UNION-REC]), so Phase 3 is no longer blocked.
- [ ] **3.2** When deferred work resumes: update runtime helpers to construct `StringError` values; update codegen so the Result payload slot holds a union pointer rather than a string pointer; update the registry signatures for fallible string functions from `Result<T, string>` to `Result<T, StringError>`; migrate the examples.

## Phase 4 — Auto-unwrap preserves the payload

Spec auto-unwrap rules ([0004-TypeSystem.md:71-77](../../compiler/spec/0004-TypeSystem.md#result-auto-unwrapping)) flatten nested Results. The spec ([ERR-PAYLOAD]) requires the original payload to survive that flattening.

- [ ] **4.1** Test: `let x = parseInt(parseInt("notanumber") + "5")` — the outer Result must carry the inner parseInt's error message (or a wrapping message that mentions it), not a generic "Error occurred".
- [ ] **4.2** Audit `maybeUnwrapResult` in [`expression_generation.go`](../../compiler/internal/codegen/expression_generation.go) — when it sees a nested Error, it must propagate the payload, not drop it.

## Phase 5 — Codegen-side `Error { message: "literal" }` construction

User code constructs Error values directly: `Error { message: "name cannot be empty" }`. That path must populate the same slot Phase 1 defined.

- [ ] **5.1** Verify by reading `generateDiscriminatedUnionConstructor` in [`expression_generation.go`](../../compiler/internal/codegen/expression_generation.go) that user-constructed Error values store their message in the slot Phase 1.1 documented. If not, fix.
- [ ] **5.2** Test: a user function `fn fail() -> Result<int, string> = Error { message: "boom" }` — `match fail() { Error { message } => print(message) }` must print `"boom"`.

## Phase 6 — Coverage and tests

- [ ] **6.1** Add a `tests/integration/error_payload_test.go` (or extend an existing one) that asserts every fallible builtin's specific error message text. One Go test case per error path. If any new builtin lands without an asserted message, the test must fail.
- [ ] **6.2** Update `examples/tested/basics/strings/string_edge_cases.osp` to assert the specific messages now that they are real, not `"Error occurred"`.

## Out of scope

- Structured error context (file/line/column tracking) — useful but not in scope.
- Stack traces — Osprey has no exceptions; not applicable.
- Localised error messages — defer.
- Error chaining (`caused_by`) — defer.

## TODO checklist

### Phase 1 — Runtime contract for `Result<T, string>`
- [x] 1.1 Layout extended in codegen (`getResultType`), not in C runtime header
- [x] 1.2 Helpers in `internal/codegen/result_helpers.go` (codegen-side, not C)
- [x] 1.3 Per-callsite messages wired in codegen for all listed builtins
- [x] 1.4 All messages are .rodata static strings; intern table dedupes

### Phase 2 — Codegen reads the slot
- [x] 2.1 Failing test landed at `examples/tested/basics/types/error_payload.osp`
- [x] 2.2 Hardcoded global removed
- [x] 2.3 GEP+load of the new err_msg slot (index 2) wired
- [x] 2.4 Test passes; four other expected-outputs updated to assert real messages

### Phase 3 — `Result<T, StringError>` (prerequisite shipped: union payloads work per [TYPE-UNION-REC])
- [ ] 3.1 Decision: defer or include in this iteration
- [ ] 3.2 Migration plan from `string` → `StringError` once union payloads work

### Phase 4 — Auto-unwrap preserves payload
- [ ] 4.1 Nested-error propagation test
- [ ] 4.2 Audit `maybeUnwrapResult`

### Phase 5 — User-constructed Error values
- [ ] 5.1 Audit `generateDiscriminatedUnionConstructor`
- [ ] 5.2 User-defined fail() test

### Phase 6 — Coverage
- [ ] 6.1 `tests/integration/error_payload_test.go` per-builtin message assertions
- [ ] 6.2 Update `string_edge_cases.osp` to specific messages

### Acceptance
- [x] Every existing test still passes (`go test ./... -count=1` clean across `compiler/{cmd,internal,tests/...}`; preexisting C-runtime fiber-tests linker failure on `main` is unrelated).
- [x] No more `"Error occurred"` global emitted by `osprey` for any input (the literal global allocation was deleted from `generateErrorBlock`; see grep on `internal/codegen/`).
- [ ] The JSON-parser canary in [`production-primitives.md`](production-primitives.md) reports `line N column M: expected ':'` (or equivalent) — pending. Phase 1+2 unblock real messages from builtins; the JSON parser still needs `closures.md` + `string-cursor.md` + `list-patterns.md` to be written in pure Osprey.

## Now next

With Phase 1+2 shipped, **Phase 5 (user-constructed `Error { message: "literal" }`)** is the next contained slice. A latent bug surfaced during this work: a function returning `Result<int, string>` whose Error arm uses the `Error { message: "boom" }` constructor produces struct type `{i8*, i8, i8*}` from the Error branch but `{i64, i8, i8*}` from a Success branch — these don't unify at LLVM level, so `llc` rejects the function. That's the third sub-test that was scoped out of [`error_payload.osp`](../../compiler/examples/tested/basics/types/error_payload.osp). Fixing it requires the Error constructor to pick the value-slot type from the *inferred Result success type at the construction site*, not from the message expression's type. After that, the user-defined-Error test can be re-enabled and Phase 5 ticked off.
