# Plan: Thread Real Error Payloads Through `Result<T, E>` — SHIPPED

Spec: [`0013-ErrorHandling.md` — Error Payload Propagation](../../compiler/spec/0013-ErrorHandling.md#error-payload-propagation--err-payload) ([ERR-PAYLOAD]).

Parent: [`production-primitives.md`](production-primitives.md).

> **Status:** all six phases shipped. Every fallible builtin now carries a
> function-specific message through the `Result<T, string>` payload, the
> hardcoded `"Error occurred"` global is gone, arithmetic auto-unwrap
> preserves the inner Error, and user-constructed `Error { message: "…" }`
> survives an explicit `Result<T, string>` return-type annotation. Phase 3
> (typed `StringError` payload) is the only deferred slice; rationale and
> migration are documented inline below.

## Problem (resolved)

Every `Error { message }` branch in user code used to bind the same
hardcoded global string `"Error occurred"`, regardless of which builtin
failed or what reason the runtime had. That made it impossible to ship a
parser, an HTTP handler, file I/O, parse failures, or division-by-zero
errors with any diagnostic value, and explicitly violated the spec
([ERR-PAYLOAD]). All of that is fixed.

## What the implementation actually changed

The C runtime never produced `Result` structs directly — it returned raw
`char *` / status `int64_t`. So instead of routing through a new
`runtime/result_runtime.c`, every wrapping site in
[`internal/codegen/`](../../compiler/internal/codegen/) now interns a
function-specific static message and stores it in the new `err_msg` slot
(index 2 of the Result struct). The C ABI is unchanged.

Key files:

- [`internal/codegen/result_helpers.go`](../../compiler/internal/codegen/result_helpers.go)
  — construction (`makeSuccessValue`, `makeErrorValueWithMessage`,
  `storeResultFields`), reader (`loadResultErrorMessage`,
  `extractResultDiscriminant`, `extractErrMsgSlot`), intern table
  (`internErrorMessage`), and arithmetic propagation
  (`withResultErrorPropagation`).
- [`internal/codegen/llvm.go`](../../compiler/internal/codegen/llvm.go)
  `generateErrorBlock` — reads `g.currentResultValue`'s `err_msg` slot
  instead of synthesising a placeholder global.
- [`internal/codegen/expression_generation.go`](../../compiler/internal/codegen/expression_generation.go)
  `generateErrorConstructor` — picks the value-slot type from
  `g.expectedReturnType` so `Error { message: "boom" }` produces a struct
  that unifies with the surrounding function's success branch at the
  trailing PHI / return site.
- [`internal/codegen/fiber_generation.go`](../../compiler/internal/codegen/fiber_generation.go)
  `generateLambdaExpression` — saves and restores `g.function` (and
  routes the terminating `Ret` through `g.builder`) so arithmetic
  propagation blocks land inside the lambda function, not the outer
  caller.

## Phase 1 — Runtime contract for `Result<T, string>` — SHIPPED

Implementation chose codegen-side messages over a runtime contract.
`getResultType` in
[`core_functions.go`](../../compiler/internal/codegen/core_functions.go)
returns `{value, discriminant: i8, err_msg: i8*}`.
`ResultFieldCount = 3` and `UnionFieldCount = 2` in
[`constants.go`](../../compiler/internal/codegen/constants.go) keep
tagged-union dispatch separate from Result detection — a previous shared
`len==2` check silently routed unions through `createSimpleEnumCondition`
mid-implementation; that's fixed and covered by
`recursive_union_list_payload.osp`.

Per-callsite messages (interned via `stringConstants`, all .rodata
`Immutable` globals):

| Builtin                                | Message                                       |
| -------------------------------------- | --------------------------------------------- |
| `parseInt`                             | `parseInt: input is not a valid integer`      |
| `parseFloat`                           | `parseFloat: input is not a valid number`     |
| `substring`                            | `substring: index out of range`               |
| `split`                                | `split: separator must not be empty`          |
| `replace`                              | `replace: needle must not be empty`           |
| `repeat`                               | `repeat: count must be non-negative`          |
| `padStart` / `padEnd`                  | `<name>: fill must not be empty`              |
| `indexOf`                              | `indexOf: needle not found`                   |
| `input`                                | `input: failed to read line`                  |
| `readFile` / `writeFile`               | `<name>: failed to <op> file`                 |
| `spawnProcess`                         | `spawnProcess: failed to spawn process`       |
| list `[idx]` out-of-bounds             | `list: index out of range`                    |
| map `[key]` missing / empty            | `map: key not found` / `map: lookup on empty map` |
| integer / float division-by-zero       | `math: division by zero`                      |

## Phase 2 — Codegen reads the slot — SHIPPED

`generateErrorBlock` reads `g.loadResultErrorMessage(g.currentResultValue)`.
The hardcoded `"Error occurred\\x00"` global is deleted. Loaded pointer is
typed `i8*` to match the new slot-2 layout. Covered by
[`examples/tested/basics/types/error_payload.osp`](../../compiler/examples/tested/basics/types/error_payload.osp)
plus four other expected-output files whose assertions were updated to the
real per-builtin text.

### Implementation note — nested Result matches

`g.currentResultValue` is a single field on the generator; a Result match's
success arm containing another Result match was overwriting it before the
*outer* error block ran, loading from a non-dominating sibling block (`llc`
rejected with "Instruction does not dominate all uses"). Fixed in
`generateResultMatchExpression` by save/restore around success-arm
generation; covered by `script_style_working.osp`'s nested `factorial`.

## Phase 3 — Discriminated-union payload types — DEFERRED (decision recorded)

Spec defines `StringError` as a discriminated union
([0012-Built-InFunctions.md:67-73](../../compiler/spec/0012-Built-InFunctions.md#L67-L73)).
Properly honouring that would let users pattern-match through the variant:

```osprey
match split(s, "") {
    Success { value } => ...
    Error   { message: InvalidArgument { message: m } } => print("bad arg: ${m}")
    Error   { message: IndexOutOfRange { index, length } } => print("oor: ${index} of ${length}")
}
```

**Decision (3.1): defer.** Phase 1+2+4+5+6 already deliver real
per-builtin messages end-to-end, which is what unblocks every downstream
production primitive (parsers, HTTP handlers, file I/O, the JSON-parser
canary). Phase 3 buys structured destructuring at the cost of a parallel
runtime layout and migrating every builtin's wrap site to construct a
`StringError` union value instead of an interned `char *`. The
recursive-union prerequisite shipped in PR #67 ([TYPE-UNION-REC]), so the
work is no longer *blocked* — it's just no longer on the critical path.

**Decision (3.2) — migration plan when Phase 3 resumes:**

1. **Type system flip.** Reverse the temporary `StringError → string`
   substitution in
   [`builtin_registry.go`](../../compiler/internal/codegen/builtin_registry.go).
   Update `getResultType` so the err_msg slot is `union *` instead of
   `i8*`. The slot index stays at 2; only the slot type changes. Add a
   `StringError` type declaration to the prelude (or as a built-in
   `*GenericType`) with variants `InvalidArgument { message: string }`,
   `IndexOutOfRange { index: int, length: int }`, `NotFound`,
   `ParseFailed { input: string }`.
2. **Runtime helpers.** Add a `makeStringError<Variant>(…)` helper per
   variant in `result_helpers.go`, mirroring the existing
   `internErrorMessage` shape. Each helper allocates / interns the
   union value and returns a pointer of the right shape for the err_msg
   slot.
3. **Builtin migration.** Replace every `g.internErrorMessage("<text>")`
   site listed in Phase 1's table with the matching
   `makeStringError<Variant>(…)` call. The substring text becomes the
   payload of `InvalidArgument` / `IndexOutOfRange` / `ParseFailed` etc.
4. **Codegen reader update.** `generateErrorBlock` already binds
   `g.variables[fieldName]` from `loadResultErrorMessage`; change that
   binding's type from `i8*` to the union pointer so subsequent
   pattern-matching arms type-check.
5. **Migrate the examples.** Update
   `examples/tested/basics/types/error_payload.osp` and
   `string_edge_cases.osp` to add at least one nested-variant pattern
   per fallible builtin, proving end-to-end variant binding.
6. **Drop the shim comment.** The "Phase-3 will re-tighten" note next to
   `generateSubstringCall` in
   [`core_functions.go`](../../compiler/internal/codegen/core_functions.go)
   goes away once that re-tightening lands.

Until Phase 3 ships, fallible builtin error types remain `string`. That
is fully spec-conformant under [ERR-PAYLOAD] (which mandates the message
must reach the caller, not the typed union shape).

## Phase 4 — Auto-unwrap preserves the payload — SHIPPED

Spec auto-unwrap rules
([0004-TypeSystem.md:71-77](../../compiler/spec/0004-TypeSystem.md#result-auto-unwrapping))
flatten nested Results; the spec ([ERR-PAYLOAD] last paragraph) requires
the original payload to survive. Phase 4 wires that propagation through
arithmetic:

- **4.2** `generateArithmeticOperationWithPos` now calls
  `withResultErrorPropagation` (in `result_helpers.go`). For each
  binary arithmetic op:
  - if neither operand is a Result, fast-path through the existing
    compute-and-wrap;
  - else emit an `arith_prop_N` block that builds a Result of the
    operator's natural value-slot type carrying the err_msg of the
    leftmost Error operand, plus an `arith_ok_N` block that extracts the
    value slots and runs the existing arithmetic, joined by a PHI in
    `arith_end_N`.
  - value-slot type for the propagated Result is picked by
    `arithmeticResultElemType` so it unifies with the success branch's
    Result at the PHI (`/` always returns `Result<float>`; any
    float-touching operand promotes; otherwise `Result<int>`).
- `maybeUnwrapResult` in
  [`function_signatures.go`](../../compiler/internal/codegen/function_signatures.go)
  was missing the `*GenericType` Result case (only handled the legacy
  `*ConcreteType` form), so functions whose body inferred to a
  `GenericType` Result were getting wrongly unwrapped at return. Fixed
  alongside Phase 4.
- `generateLambdaExpression` in
  [`fiber_generation.go`](../../compiler/internal/codegen/fiber_generation.go)
  was hard-coding `entryBlock.NewRet(bodyValue)` and not saving
  `g.function`. With propagation creating new blocks during body gen,
  those blocks landed in the *outer* function and the lambda's entry
  block ended with a branch to a label that didn't exist. Fixed by
  swapping `g.function`/`g.builder` around body gen and terminating
  through `g.builder.NewRet` instead.

- **4.1** Tests:
  - `examples/tested/basics/types/error_payload.osp` includes
    `match parseInt("nope") + 5 { … }` and a `(parseInt("nope") * 3) - 1`
    chain, plus a `parseInt("12") + 8` positive case.
  - `tests/integration/error_payload_test.go` has
    `arith_propagates_parseInt_message` and
    `arith_chain_propagates_message`.

## Phase 5 — User-constructed `Error { message: "literal" }` — SHIPPED

`generateErrorConstructor` in
[`expression_generation.go`](../../compiler/internal/codegen/expression_generation.go)
now picks the value slot type from `g.expectedReturnType` (via
`inferErrorValueSlotType`) and emits a typed zero placeholder
(`zeroValueForType`). The message goes into the err_msg slot. With this
fix:

```osprey
fn fail() -> Result<int, string> = Error { message: "boom" }
match fail() { Error { message } => print(message) }   // prints "boom"
```

- **5.1** `generateDiscriminatedUnionConstructor` was audited; only the
  user-facing `Error { … }` shorthand needed fixing because non-Result
  unions already pick their value-slot type from the union's tag layout.
- **5.2** `examples/tested/basics/types/error_payload.osp` covers
  `failInt() -> Result<int, string>`, `failFloat() -> Result<float, string>`,
  `okInt() -> Result<int, string>` (Success arm), and a `rethrowInt(s)`
  function that round-trips a builtin's err_msg through a user
  `Error { message: message }` re-throw.
- `tests/integration/error_payload_test.go` has
  `user_constructed_error` (declared-return) and
  `user_rethrow_forwards_builtin_message` (inferred-return).

The latent shim — a registry-wide `StringError → string` substitution —
is the single concession to keep Phase 5 / Phase 4 round-trippable until
Phase 3 lands. Documented inline above; reverse-migration steps in 3.2.

## Phase 6 — Coverage and tests — SHIPPED

- **6.1** New file
  [`tests/integration/error_payload_test.go`](../../compiler/tests/integration/error_payload_test.go)
  asserts the exact string emitted by every fallible builtin's Error path,
  plus the arithmetic propagation and user-constructed cases. A separate
  `TestErrorPayloadNoGenericFallback` re-runs five representative programs
  and fails if `"Error occurred"` shows up anywhere — the regression
  tripwire for the deleted hardcoded global.
- **6.2** [`examples/tested/basics/strings/string_edge_cases.osp`](../../compiler/examples/tested/basics/strings/string_edge_cases.osp)
  now binds and prints `message` for every Error path (substring,
  indexOf, replace, repeat, padStart, padEnd, parseFloat, parseInt). The
  matching `examples_test.go` entry asserts the full per-line output.

## Out of scope

- Structured error context (file/line/column tracking) — useful but not in scope.
- Stack traces — Osprey has no exceptions; not applicable.
- Localised error messages — defer.
- Error chaining (`caused_by`) — defer.

## TODO checklist — all closed

### Phase 1 — Runtime contract for `Result<T, string>`
- [x] 1.1 Layout extended in codegen (`getResultType`), not in C runtime header
- [x] 1.2 Helpers in `internal/codegen/result_helpers.go` (codegen-side, not C)
- [x] 1.3 Per-callsite messages wired in codegen for every listed builtin
- [x] 1.4 All messages are .rodata `Immutable` globals; intern table dedupes

### Phase 2 — Codegen reads the slot
- [x] 2.1 Failing test landed at `examples/tested/basics/types/error_payload.osp`
- [x] 2.2 Hardcoded global removed from `generateErrorBlock`
- [x] 2.3 GEP+load of the err_msg slot (index 2) wired
- [x] 2.4 Test passes; four other expected-output files updated to assert real messages

### Phase 3 — `Result<T, StringError>` — deferred with documented migration
- [x] 3.1 Decision recorded above: defer; rationale documented
- [x] 3.2 Six-step migration plan recorded above

### Phase 4 — Auto-unwrap preserves payload
- [x] 4.1 Nested-error propagation cases in both the `.osp` example and the Go integration test
- [x] 4.2 `withResultErrorPropagation` wired into `generateArithmeticOperationWithPos`; `maybeUnwrapResult` GenericType bug fixed; `generateLambdaExpression` block/function save-restore fixed

### Phase 5 — User-constructed Error values
- [x] 5.1 `generateErrorConstructor` audited and rewritten to honour `expectedReturnType`
- [x] 5.2 `fail()` / `rethrowInt()` test cases in the `.osp` example and Go integration test

### Phase 6 — Coverage
- [x] 6.1 `tests/integration/error_payload_test.go` per-builtin message assertions + no-fallback tripwire
- [x] 6.2 `string_edge_cases.osp` binds and prints every per-builtin message; expected output updated

### Acceptance
- [x] Every existing test still passes (`go test ./...` clean across `compiler/{cmd,internal,tests/...}`).
- [x] No more `"Error occurred"` global emitted anywhere by `osprey` for any input (`TestErrorPayloadNoGenericFallback` makes this a permanent invariant).
- [ ] The JSON-parser canary in [`production-primitives.md`](production-primitives.md) reports `line N column M: expected ':'` (or equivalent) on malformed input — still pending, but unblocked: depends on the remaining production-primitive plans (`closures.md`, `string-cursor.md`, `list-patterns.md`).

## Next, after Phase 3 resumes

Per `production-primitives.md` sequencing, the next contained slice is
**`closures.md`** — lambdas-with-capture, which is the prerequisite for
every higher-order combinator the JSON parser will compose. Phase 3 of
this plan can be picked up at any time without blocking that work.
