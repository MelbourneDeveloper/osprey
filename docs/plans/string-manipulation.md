# Plan: String Manipulation + Uniform Function Call Syntax (UFCS)

Two related workstreams. Together they unblock the #1 gap to a usable HTTP server in Osprey: parsing query strings, paths, headers, and JSON bodies in idiomatic FP style.

Spec reference: [compiler/spec/0012-Built-InFunctions.md](../../compiler/spec/0012-Built-InFunctions.md), section "String Functions".

## Motivation

Today's builtin string surface is `length`, `contains`, `substring`, `parseInt`, `join` ([compiler/internal/codegen/constants.go:78-82](../../compiler/internal/codegen/constants.go#L78-L82)). Missing: `split`, `trim`, `toUpperCase`, `toLowerCase`, `startsWith`, `endsWith`, `replace`, `indexOf`, `repeat`, `reverse`, `take`, `drop`, `pad*`, `lines`, `words`, `parseFloat`, `isEmpty`. You cannot:

- Parse a query string (`?name=alice&age=30`) — no `split`, no `indexOf`.
- Normalise a header value — no `trim`, no `toLowerCase`.
- Test a path prefix — no `startsWith`.
- Build a response by concatenating parts — no infallible `+` on strings (current spec wraps everything in `Result`).

Separately, **method-call syntax** (`"a".getStringLength()`) parses today but its codegen is stubbed at [compiler/internal/codegen/expression_generation.go:1409-1414](../../compiler/internal/codegen/expression_generation.go#L1409-L1414). Wiring it gives users discoverability (autocomplete after `.`) without changing the language; it desugars to a plain call.

## Design summary

See the spec for the full function list and rationale. Headline decisions:

- **Total functions don't wrap in `Result`.** `length`, `contains`, `startsWith`, `toUpperCase`, `trim`, etc. return plain values. Only partial operations (`substring` with bad indices, `parseInt` on garbage, `split` with empty separator) return `Result<T, StringError>`. This matches Elm's `String` module and Haskell's `Data.Text`.
- **Subject-first argument order** so pipe and UFCS both work: `split(s, ",")`, `trim(s)`, `startsWith(s, "GET ")`.
- **Pipe (`|>`) is the documented primary form.** UFCS is sugar.
- **No `Char` type yet** — higher-order char operations (`String.map`, `filter`, `foldl`) are deferred.

## Workstream A — String functions

### A1. Fix existing wrongly-`Result`-wrapped functions

`length` and `contains` currently return `Result<T, StringError>` for no reason — they can't fail on well-formed UTF-8. Unwrap them.

- [ ] Update [compiler/internal/codegen/builtin_registry.go:252-282](../../compiler/internal/codegen/builtin_registry.go#L252-L282): change `length`'s `ReturnType` from `Result<int, StringError>` to plain `int`; change `contains`'s `ReturnType` from `Result<bool, StringError>` to plain `bool`. Update `Signature` strings to match.
- [ ] Update generators `generateLengthCall` and `generateContainsCall` in [compiler/internal/codegen/core_functions.go](../../compiler/internal/codegen/core_functions.go) to return raw scalars, not constructed `Success {}` records.
- [ ] Find every example/test using `match length(...)` / `match contains(...)` and inline the value. Search: `rg "match\s+(length|contains)\(" compiler/`.

### A2. Add infallible builtins

Each item below requires: (a) constant in `constants.go`, (b) registry entry in `builtin_registry.go`, (c) generator in a new file (`string_functions.go`) since `core_functions.go` is approaching the 500-LOC limit per CLAUDE.md, (d) C runtime helper in `compiler/runtime/string_runtime.c` (new file), (e) at least one example in `compiler/examples/tested/basics/` exercising it, (f) integration test asserting expected output.

- [ ] `isEmpty(s: string) -> bool`
- [ ] `startsWith(s: string, prefix: string) -> bool`
- [ ] `endsWith(s: string, suffix: string) -> bool`
- [ ] `toUpperCase(s: string) -> string`  *(Unicode simple case mapping; document the SS expansion in tests)*
- [ ] `toLowerCase(s: string) -> string`
- [ ] `trim(s: string) -> string`  *(Unicode `White_Space` per spec)*
- [ ] `trimStart(s: string) -> string`
- [ ] `trimEnd(s: string) -> string`
- [ ] `reverse(s: string) -> string`  *(codepoint reversal — note grapheme-cluster reversal is future work)*
- [ ] `take(s: string, n: int) -> string`  *(clamp on overflow/negative)*
- [ ] `drop(s: string, n: int) -> string`
- [ ] `lines(s: string) -> Array<string>`
- [ ] `words(s: string) -> Array<string>`

### A3. Add fallible builtins

Same scaffolding as A2 plus error-path tests.

- [ ] `indexOf(s: string, needle: string) -> Result<int, StringError>` — `Error(NotFound)` when absent.
- [ ] `split(s: string, separator: string) -> Result<Array<string>, StringError>` — `Error(InvalidArgument)` on empty separator (matches Haskell `splitOn`).
- [ ] `replace(s: string, needle: string, replacement: string) -> Result<string, StringError>` — `Error(InvalidArgument)` on empty needle.
- [ ] `repeat(s: string, n: int) -> Result<string, StringError>` — `Error(InvalidArgument)` on negative `n`.
- [ ] `padStart(s: string, targetLength: int, fill: string) -> Result<string, StringError>` — `Error(InvalidArgument)` on empty `fill`.
- [ ] `padEnd(s: string, targetLength: int, fill: string) -> Result<string, StringError>`
- [ ] `parseFloat(s: string) -> Result<float, StringError>`
- [ ] Tighten `parseInt` in `compiler/runtime/system_runtime.c` to reject non-numeric input strictly (current TODO at line ~470 about `atoll` returning 0 silently). Return `Error(ParseFailed)` on any invalid character; no leading/trailing whitespace tolerance.

### A4. Promote `+` on strings to infallible

String concatenation cannot fail. Make `string + string -> string` (not `Result<string, _>`).

- [ ] Audit binary-op codegen for the `string + string` case in `compiler/internal/codegen/`. Confirm the result type is plain `string`, not `Result`.
- [ ] Add a one-line example to [compiler/spec/0012-Built-InFunctions.md](../../compiler/spec/0012-Built-InFunctions.md) under "Concatenation Operator" (already in spec text — just verify behavior matches).

### A5. End-to-end demo

- [ ] Add `compiler/examples/tested/basics/string_pipeline.osp` that parses a fake query string (`"name=alice&age=30&role=admin"`) into key-value pairs using `split`, `indexOf`, `substring`, `trim` — entirely via `|>` pipes. Assert output in `.osp.expectedoutput`.
- [ ] Once the demo works, port [compiler/examples/tested/http/http_server_example.osp](../../compiler/examples/tested/http/http_server_example.osp) to use real path-prefix matching via `startsWith` instead of exact-string `match`.

## Workstream B — UFCS codegen

Goal: make `"a".getStringLength()` compile to `getStringLength("a")`.

Grammar and AST are already in place ([compiler/osprey.g4:154-155](../../compiler/osprey.g4#L154-L155), [compiler/internal/ast/ast.go:275](../../compiler/internal/ast/ast.go#L275)). Only codegen + type inference need work.

### B1. Codegen rewrite

- [ ] In [compiler/internal/codegen/expression_generation.go:1409-1414](../../compiler/internal/codegen/expression_generation.go#L1409-L1414), replace the stub:
  1. Build a `*ast.CallExpression` with `Function = Identifier{Name: methodCall.MethodName}` and `Arguments = [methodCall.Receiver, ...methodCall.Args]`.
  2. Delegate to `generateCallExpression`.
- [ ] Add field-access disambiguation: before rewriting, check whether `methodCall.Receiver`'s inferred type has a field named `methodCall.MethodName` of function type. If so, call the field instead of doing UFCS. This is the Scala precedence rule.
- [ ] Verify the no-parens form `x.field` still routes to field access, not to UFCS. (Grammar already enforces this — UFCS requires parens — but add a test.)

### B2. Type inference

- [ ] In [compiler/internal/codegen/type_inference.go:687](../../compiler/internal/codegen/type_inference.go#L687) (`case *ast.MethodCallExpression`), perform the same call-rewrite virtually for type checking. The inferred type of `x.f(a)` must equal the inferred type of `f(x, a)`.

### B3. Error messages

- [ ] When UFCS dispatch fails (no function with that name, or arity mismatch after prepending the receiver), the error should mention both the method-call form and the equivalent call form, e.g. `cannot call "a".getStringLength() — no function 'getStringLength(string) -> _' in scope`.

### B4. Tests

- [ ] `compiler/examples/tested/basics/ufcs_string.osp` — calls every Workstream A function via both pipe and UFCS, asserts identical output.
- [ ] `compiler/examples/failscompilation/ufcs_field_collision.osp` — record with a `string`-typed field named `trim`; `record.trim()` must give a clear error (calling a string is not allowed) rather than silently UFCS-ing.

## Sequencing

1. **A1 first** (fix the wrongly-`Result`-wrapped existing functions) — this is a small breaking change to the user-facing surface and should land before adding more.
2. **A2 + A3 in parallel** — independent builtins, each PR adds one or two.
3. **B1 + B2 + B3** as a single PR — small and atomic.
4. **A4** anytime; verify-only if codegen is already correct.
5. **A5** last — proves the whole thing hangs together.

## TODO checklist

### Workstream A — String functions
- [x] A1.1 Unwrap `length` return type in registry + generator
- [x] A1.2 Unwrap `contains` return type in registry + generator
- [x] A1.3 Migrate all `match length(...)` / `match contains(...)` callsites
- [x] A2.1 `isEmpty`
- [x] A2.2 `startsWith`
- [x] A2.3 `endsWith`
- [x] A2.4 `toUpperCase`
- [x] A2.5 `toLowerCase`
- [x] A2.6 `trim`
- [x] A2.7 `trimStart`
- [x] A2.8 `trimEnd`
- [x] A2.9 `reverse`
- [x] A2.10 `take`
- [x] A2.11 `drop`
- [x] A2.12 `lines`
- [x] A2.13 `words`
- [x] A3.1 `indexOf`
- [x] A3.2 `split`
- [x] A3.3 `replace`
- [x] A3.4 `repeat`
- [x] A3.5 `padStart`
- [x] A3.6 `padEnd`
- [x] A3.7 `parseFloat`
- [x] A3.8 Strict `parseInt` (fix silent-zero bug)
- [x] A4 Verify `string + string -> string` is infallible
- [x] A5.1 `string_pipeline.osp` query-string demo
- [x] A5.2 Port `http_server_example.osp` to `startsWith` — done as a new
  pure-logic sibling example [`route_match.osp`](../../compiler/examples/tested/basics/strings/route_match.osp)
  rather than touching the live HTTP server example (which is a network test
  with a brittle expected output).

### Workstream B — UFCS
- [x] B1.1 Implement codegen rewrite `x.f(a, b)` → `f(x, a, b)`
- [ ] B1.2 Field-access precedence rule (field wins over UFCS) — **deferred.**
  The negative test `ufcs_field_collision.ospo` proves we don't silently
  succeed (compilation fails with a type-mismatch error), which is the
  important property. Promoting the error message to mention field-vs-UFCS
  precedence requires plumbing record-type info through the UFCS rewrite
  step. Tracked as a follow-up.
- [x] B1.3 Grammar test: `x.field` is field access, never a call
- [x] B2 Type inference parity for `MethodCallExpression`
- [x] B3 Helpful error when UFCS dispatch fails (error now reads
  `UFCS call \`_.foo(...)\` rewrites to \`foo(_, ...)\`: <inner>`)
- [x] B4.1 `ufcs_string.osp` (pipe == UFCS == direct)
- [x] B4.2 `failscompilation/ufcs_field_collision.ospo`

### Test coverage

- [x] **C unit tests with hard assertions** — [`runtime/string_runtime_tests.c`](../../compiler/runtime/string_runtime_tests.c)
  exercises every `osp_string_*` helper. Each function covers the happy
  path, every documented error case, NULL inputs, and boundary values
  (empty strings, n=0, n<0, n>len, prefix/suffix longer than s, INT64_MIN,
  INT64_MAX, overflow). Compiled with `-Werror -Wall -Wextra -ftrapv
  -fsanitize=signed-integer-overflow` so any signed-overflow bug aborts
  the test before the linker. Wired into `compiler/Makefile`'s `c-lint`
  (under the full pedantic flag stack) and `c-test`. 19 test groups, all
  green.
- [x] **E2E edge-case proof** — [`examples/tested/basics/strings/string_edge_cases.osp`](../../compiler/examples/tested/basics/strings/string_edge_cases.osp)
  runs every error path and boundary through the actual compiled
  pipeline (Osprey source → LLVM IR → linked binary → stdout). The
  integration runner compares the full stdout to the expected output
  verbatim — any deviation in any builtin causes a mismatch failure.
  Catches `INT64_MIN` parsing as well as the empty-needle/negative-n/
  empty-fill/out-of-range/inverted-range/overflow rejections.
- [x] **E2E happy-path proof** — [`string_pipeline.osp`](../../compiler/examples/tested/basics/strings/string_pipeline.osp),
  [`ufcs_string.osp`](../../compiler/examples/tested/basics/strings/ufcs_string.osp),
  [`route_match.osp`](../../compiler/examples/tested/basics/strings/route_match.osp),
  [`string_utils_combined.osp`](../../compiler/examples/tested/basics/strings/string_utils_combined.osp),
  and [`result_type_workflow.osp`](../../compiler/examples/tested/basics/types/result_type_workflow.osp)
  all exact-match-asserted in `tests/integration/examples_test.go`.
- [x] **Negative compilation test** —
  [`examples/failscompilation/ufcs_field_collision.ospo`](../../compiler/examples/failscompilation/ufcs_field_collision.ospo)
  has a `.expectedoutput` file enforcing the exact error message.

### Drive-by fixes uncovered during implementation
These weren't on the original plan but were blockers found while
implementing the above. Documenting here for traceability.

- [x] **AST builder: chained method args mis-aligned.**
  `s.trim().take(3)` was building `MethodCallExpression` with `Arguments=[]`
  because `ctx.ArgList(i)` returns the *i-th present* ArgList, not the
  ArgList for chain position `i` (empty arg lists are skipped by ANTLR).
  Fixed in [`builder_calls.go:argListForChainElement`](../../compiler/internal/ast/builder_calls.go) by walking children in source order and pairing
  each LPAREN with its inner ArgList.

- [x] **`isModuleName` stole every uppercase-receiver UFCS call.**
  The old heuristic "uppercase identifier == module" routed
  `UpperCaseVar.method()` to a placeholder that always returned 42
  (`fiber_generation.go:moduleAccessPlaceholder = 42`). Replaced with a
  pre-pass in [`builder_core.go:BuildProgram`](../../compiler/internal/ast/builder_core.go) that records every actual
  `module Name { ... }` declaration; `isModuleName` now consults that set.
  This keeps the fiber-module tests passing while letting UFCS work on
  uppercase variable names.

- [x] **`substring` finally rejects out-of-range indices.**
  The existing `generateSubstringCall` declared a `Result<string, Error>`
  return type but never produced an Error — bad indices silently returned
  garbage. Rewritten to delegate bounds-checking to the C helper
  `osp_string_substring`, which returns NULL → wrapped as Error.

- [x] **`parseInt` of `INT64_MIN` no longer traps under `-ftrapv`.**
  Initial `osp_parse_int_strict` accumulated in `int64_t` and would
  overflow on `"-9223372036854775808"` (magnitude is INT64_MAX+1) before
  the negation could rescue it — the trap fired inside the linked
  runtime even though the standalone C test (without `-ftrapv`) passed.
  Rewrote the accumulator in `uint64_t`, special-cased INT64_MIN, and
  added `-ftrapv -fsanitize=signed-integer-overflow` to the C test
  compile so this category of bug surfaces during `make c-test` instead
  of during a runtime crash.

### Spec divergences from this implementation

The spec describes the target behaviour; the v1 implementation cuts the
following corners and tracks them as follow-ups:

- **Unicode codepoints vs bytes.** Spec says `length` etc. count
  codepoints. v1 counts bytes (matches the existing `strlen`-based code).
  Same for `take`, `drop`, `substring`, `indexOf`, `reverse`. UTF-8-aware
  rewrites are deferred to a follow-up.
- **Unicode simple case mapping.** Spec mentions German `ß` → `SS`.
  v1 uses ASCII-only `tolower`/`toupper`. Documented in registry
  description strings.
- **Error message payload.** Fallible builtins set discriminant=1 and
  null value; the match-expression `Error { message }` branch always
  binds the same static `"Error occurred\x00"` global string regardless
  of which builtin failed. Pre-existing limitation of the Result codegen
  (`llvm.go:generateErrorBlock`); not specific to strings. Replacing the
  static message with a real `StringError` discriminated union is its own
  workstream — left out of scope here.

### Out of scope (separate plans needed)
- [ ] `Char` type + higher-order `String.map` / `filter` / `foldl` / `any` / `all`
- [ ] Grapheme-cluster aware operations
- [ ] UTF-8 codepoint counting for `length` / `take` / `drop` / `substring` / `indexOf` / `reverse`
- [ ] Unicode simple case mapping for `toUpperCase` / `toLowerCase`
- [ ] Per-error-kind `StringError` payload routed into the match `message` binding
- [ ] Regex
- [ ] Formatting (`printf`-style or `String.format`)
- [ ] `Maybe`/`Option` (would let `indexOf` return `Maybe<int>` instead of `Result`)
