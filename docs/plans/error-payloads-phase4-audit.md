# Phase 4 audit — auto-unwrap drops the err_msg payload

Companion to [`error-payloads.md`](error-payloads.md). Author: OspreyServant (assist).
Scope: read-only audit of items **4.1** and **4.2** as of `2abe50c`.

## Finding

`unwrapIfResult` at [`expression_generation.go:1020`](../../compiler/internal/codegen/expression_generation.go#L1020) extracts slot 0 (value) of a 3-field Result struct unconditionally:

```go
return g.builder.NewExtractValue(val, 0)
```

Slot 1 (discriminant) and slot 2 (`err_msg: i8*`) are dropped. The comment on the function ("assumes the Result is Success — errors will propagate at runtime") is misleading because nothing in the generated IR actually checks the discriminant on this path; the unwrap is unconditional and the `err_msg` slot is never read.

Call sites that lose the payload as a result:

| File:line | Context |
|---|---|
| [`expression_generation.go:1044`](../../compiler/internal/codegen/expression_generation.go#L1044) | `generateComparisonOperationWithPos` — both operands unwrapped before `<`, `==`, etc. |
| [`expression_generation.go:1045`](../../compiler/internal/codegen/expression_generation.go#L1045) | (same call, right operand) |
| [`function_signatures.go:944`](../../compiler/internal/codegen/function_signatures.go#L944) | `maybeUnwrapResult` — last-resort unwrap when neither declared nor inferred return type is Result |
| [`effects_generation.go:573`](../../compiler/internal/codegen/effects_generation.go#L573), [`:710`](../../compiler/internal/codegen/effects_generation.go#L710) | Effect-arg coercion |
| [`core_functions.go:42`](../../compiler/internal/codegen/core_functions.go#L42) | Result-arg coercion in core helpers |

`maybeUnwrapResult` at [`function_signatures.go:909`](../../compiler/internal/codegen/function_signatures.go#L909) decides *whether* to unwrap based on declared/inferred return type. When neither is Result but the body is, it calls `unwrapIfResult` — losing the payload. The decision logic itself is fine; only the act of unwrapping is lossy.

## What spec [ERR-PAYLOAD] requires

> "the original payload must survive that flattening"

For `let x = parseInt(parseInt("notanumber") + "5")`:

- Inner `parseInt` produces `Error { discriminant=1, value=0, err_msg="parseInt: input is not a valid integer" }`
- The `+ "5"` operator must NOT silently discard slot 2 — either the error propagates as-is (short-circuit) or, if the surrounding expression is itself in Result context, the outer Result must carry the inner's err_msg.

## Recommended fix design

Two paths, in increasing scope:

### Option A — error short-circuit at every unwrap site (smaller, runtime-cost on hot path)

Replace `unwrapIfResult`'s unconditional `NewExtractValue(val, 0)` with a conditional that branches on slot 1:

```
if discriminant != 0:
    propagate the entire Result (return early from the enclosing fn,
    or store-and-skip in expression context)
else:
    extract slot 0
```

This requires `unwrapIfResult` to know what to do on the error branch. In a function with a Result return type it can `ret` the same Result. In a sub-expression it needs a sentinel value of the success-slot type plus a flag to surface up — basically the same problem the spec is trying to fix, recursively. Verdict: messy.

### Option B — `?`-style propagation in `maybeUnwrapResult`, leave arithmetic loud (recommended)

Distinguish *function-return unwrap* from *arithmetic-arg unwrap*:

1. **Function-return unwrap** (`maybeUnwrapResult`, line 944): when the declared return type is `T` but the body is `Result<T, E>`, the compiler should reject the program at type-check time rather than silently extract slot 0. The `unwrapIfResult` fallback exists because the type system can't always prove the body is Result; once inference catches all such cases this branch becomes dead code.
2. **Arithmetic-arg unwrap** (`generateComparisonOperationWithPos`, lines 1044-1045): if either operand is a Result-typed expression with discriminant=1, the *enclosing expression* must produce that Error. Implementing this means tagging the comparison's result as `Result<bool, E>` and threading the payload through. For chained comparisons (`a < b < c`) the payload travels all the way out.

Option B keeps the runtime check in one place — at the boundary where a Result-producing expression feeds a non-Result context. Inside arithmetic the value is always Result-typed and the err_msg slot is carried by struct identity.

## Test cases to land alongside the fix

These should go into a new file `compiler/examples/tested/basics/types/nested_error_payload.osp` (currently not present — leaving the slot for the fixer):

```osprey
// Inner parseInt fails; the outer Result must carry the inner's specific message.
match parseInt(parseInt("notanumber") + "5") {
    Success { value }   => print("unreachable: ${value}")
    Error   { message } => print(message)
    // Expected: "parseInt: input is not a valid integer"
}

// Chained: division by zero inside a comparison.
let cmp = (10 / 0) < 5
match cmp {
    Success { value }   => print("unreachable: ${value}")
    Error   { message } => print(message)
    // Expected: "math: division by zero"
}

// Mixed: substring failure inside arithmetic via length.
match length(substring("ab", 5, 10)) + 1 {
    Success { value }   => print("unreachable: ${value}")
    Error   { message } => print(message)
    // Expected: "substring: index out of range"
}
```

## Out of scope for this audit

- Implementing the fix — Master is currently in `expression_generation.go` (Phase 5); a second writer there would conflict.
- `effects_generation.go` payload propagation through effect handlers — depends on Option A/B decision.
- The `core_functions.go:42` site — needs a separate look once 4.2 lands.
