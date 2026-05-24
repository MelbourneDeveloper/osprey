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

Most fallible builtins today return `Result<T, string>` (file I/O, parseInt, parseFloat, http*). They are the easiest target because the payload is a single `char*`.

- [ ] **1.1** Document the layout in [`compiler/runtime/string_runtime.h`](../../compiler/runtime/string_runtime.h):
  ```c
  // Layout of Result<T, char*> when returned from C runtime to Osprey:
  //   { i64 discriminant, T value, char* error_message }
  // discriminant == 0: value valid, error_message ignored.
  // discriminant == 1: value zeroed, error_message points to a null-terminated
  //   string owned by the runtime (static or heap; runtime guarantees liveness
  //   for the lifetime of the Result).
  ```
- [ ] **1.2** Add `osp_result_make_error(const char *msg)` and `osp_result_make_ok(<value>)` helpers in `runtime/result_runtime.c` (new file, ≤200 LOC) so runtime functions don't open-code the struct layout.
- [ ] **1.3** Audit every `osp_*` function that returns a Result and verify it populates the message slot on the error path. Specific files:
  - [`runtime/string_runtime.c`](../../compiler/runtime/string_runtime.c) — `osp_string_substring`, `osp_parse_int_strict`, `osp_parse_float_strict`, `osp_string_split`, `osp_string_replace`, `osp_string_repeat`, `osp_string_pad_*`, `osp_string_index_of`. Each must pass a meaningful message (e.g., `"substring: start index out of range"`, `"parseInt: non-numeric input"`).
  - [`runtime/system_runtime.c`](../../compiler/runtime/system_runtime.c) — `osp_read_file`, `osp_write_file`, etc.
  - [`runtime/http_client_runtime.c`](../../compiler/runtime/http_client_runtime.c) — HTTP error paths.
- [ ] **1.4** Strings used as error messages MUST be either static-string constants in the .rodata segment (lifetime = process) or `strdup`'d heap allocations attached to the Result (lifetime managed by whoever frees the Result). Pick one per function; document inline.

## Phase 2 — Codegen reads the slot

- [ ] **2.1** Delete the global-creation lines at [`llvm.go:2298-2308`](../../compiler/internal/codegen/llvm.go#L2298-L2308). Replace with: load the message pointer from the Result struct of the matched expression (the second/third field per Phase 1.1's layout).
- [ ] **2.2** The matched expression's LLVM value is already live — it's the discriminant test that drives the match. Reuse that value (or its alloca) to GEP into the message slot. Look at how `generateSuccessBlock` (nearby in `llvm.go`) extracts the success value — the error path mirrors it.
- [ ] **2.3** Type the loaded pointer as `i8*` and bind it to `g.variables[fieldName]`. Downstream code that uses `message` as a string already knows how to read `i8*`.
- [ ] **2.4** Failing test FIRST: `examples/tested/types/error_payload.osp` calls `split("abc", "")` and asserts the printed message is `"split: separator must not be empty"` (or whatever specific string Phase 1.3 chose). The `.expectedoutput` pins the exact text.

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
- [ ] 1.1 Document layout in `string_runtime.h`
- [ ] 1.2 `runtime/result_runtime.c` with `osp_result_make_error` / `osp_result_make_ok`
- [ ] 1.3 Audit and fix every `osp_*` Result-returning function to populate message slot
- [ ] 1.4 Decide static-string vs strdup per function; document inline

### Phase 2 — Codegen reads the slot
- [ ] 2.1 Failing test `examples/tested/types/error_payload.osp` with exact expected output
- [ ] 2.2 Delete hardcoded global at `llvm.go:2298-2308`
- [ ] 2.3 Replace with GEP+load of the message slot
- [ ] 2.4 Test passes

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
- [ ] Every existing test still passes.
- [ ] No more `"Error occurred"` global emitted anywhere by `osprey` for any input.
- [ ] The JSON-parser canary in [`production-primitives.md`](production-primitives.md) reports `line N column M: expected ':'` (or equivalent) on malformed input — not `"Error occurred"`.
