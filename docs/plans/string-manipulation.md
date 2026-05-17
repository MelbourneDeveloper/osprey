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
- [ ] A1.1 Unwrap `length` return type in registry + generator
- [ ] A1.2 Unwrap `contains` return type in registry + generator
- [ ] A1.3 Migrate all `match length(...)` / `match contains(...)` callsites
- [ ] A2.1 `isEmpty`
- [ ] A2.2 `startsWith`
- [ ] A2.3 `endsWith`
- [ ] A2.4 `toUpperCase`
- [ ] A2.5 `toLowerCase`
- [ ] A2.6 `trim`
- [ ] A2.7 `trimStart`
- [ ] A2.8 `trimEnd`
- [ ] A2.9 `reverse`
- [ ] A2.10 `take`
- [ ] A2.11 `drop`
- [ ] A2.12 `lines`
- [ ] A2.13 `words`
- [ ] A3.1 `indexOf`
- [ ] A3.2 `split`
- [ ] A3.3 `replace`
- [ ] A3.4 `repeat`
- [ ] A3.5 `padStart`
- [ ] A3.6 `padEnd`
- [ ] A3.7 `parseFloat`
- [ ] A3.8 Strict `parseInt` (fix silent-zero bug)
- [ ] A4 Verify `string + string -> string` is infallible
- [ ] A5.1 `string_pipeline.osp` query-string demo
- [ ] A5.2 Port `http_server_example.osp` to `startsWith`

### Workstream B — UFCS
- [ ] B1.1 Implement codegen rewrite `x.f(a, b)` → `f(x, a, b)`
- [ ] B1.2 Field-access precedence rule (field wins over UFCS)
- [ ] B1.3 Grammar test: `x.field` is field access, never a call
- [ ] B2 Type inference parity for `MethodCallExpression`
- [ ] B3 Helpful error when UFCS dispatch fails
- [ ] B4.1 `ufcs_string.osp` (pipe == UFCS == direct)
- [ ] B4.2 `failscompilation/ufcs_field_collision.osp`

### Out of scope (separate plans needed)
- [ ] `Char` type + higher-order `String.map` / `filter` / `foldl` / `any` / `all`
- [ ] Grapheme-cluster aware operations
- [ ] Regex
- [ ] Formatting (`printf`-style or `String.format`)
- [ ] `Maybe`/`Option` (would let `indexOf` return `Maybe<int>` instead of `Result`)
