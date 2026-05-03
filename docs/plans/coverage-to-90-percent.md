# Plan: Ratchet Code Coverage to 90%

## Current State

- **Total coverage**: 78.1% (after excluding generated parser).
- **Threshold**: 78% (current achievable; ratchets up only).
- **Aim**: 90%.
- **Excluded from measurement**: `compiler/parser/osprey_*.go` (ANTLR-generated; see `Makefile` `_test` target).
- **Coverage configuration**: `go test -coverpkg=./...` is required so integration tests (which call `codegen.CompileAndRunJIT` directly) instrument all packages, not just the package under test.

## Coverage Gaps by File (avg per file, < 90%)

Source: `make _test` followed by `go tool cover -func=coverage.out`.

| Coverage | File | Functions | Notes |
|---:|---|---:|---|
| 0.0% | `internal/ast/ast.go` | 40 | Marker methods for type assertions (`isStatement()`, `isExpression()`). Dead-but-required — see Untestable section below. |
| 0.0% | `internal/ast/expressions.go` | 2 | Marker methods (`PerformExpression.isExpression`, `HandlerExpression.isExpression`). Untestable. |
| 0.0% | `internal/codegen/websocket_bridge.go` | 4 | CGo functions called FROM C runtime, never from Go. See Untestable section. |
| 49.5% | `internal/logging/logging.go` | 6 | Newly-added logger; no unit tests yet. Easy win. |
| 53.8% | `internal/codegen/http_generation.go` | 23 | HTTP server/client codegen paths not exercised by examples. |
| 60.6% | `internal/codegen/statement_generation.go` | 6 | Some statement variants not in any example. |
| 63.2% | `internal/codegen/errors.go` | 55 | Many error wrappers never hit by current examples. |
| 64.3% | `internal/codegen/fiber_generation.go` | 20 | Channel select / select-with-default paths. |
| 64.8% | `internal/codegen/system_generation.go` | 13 | Process spawning edge cases. |
| 66.2% | `internal/codegen/core_functions.go` | 20 | Some core builtins not exercised. |
| 67.5% | `internal/cli/main_interface.go` | 8 | CLI flag combinations not all tested. |
| 69.3% | `internal/codegen/logging.go` | 8 | Codegen-side logging hooks. |
| 71.7% | `internal/codegen/compilation.go` | 14 | Compilation pipeline error branches. |
| 73.3% | `internal/codegen/expression_generation.go` | 73 | Many rarely-used expression forms. |
| 73.7% | `internal/codegen/effects_generation.go` | 24 | Effect resumption / shallow handler paths. |
| 74.1% | `internal/codegen/type_inference.go` | 144 | Many narrow inference paths. |
| 77.1% | `internal/cli/cli.go` | 37 | Subcommand error branches. |
| 77.1% | `internal/language/descriptions/functions.go` | 5 | A couple of public getters not unit-tested. |
| 79.1% | `internal/codegen/function_signatures.go` | 52 | Mostly OK; a few mangling branches missed. |
| 79.2% | `internal/ast/builder_fibers.go` | 6 | Channel/spawn AST construction. |
| 80.5% | `internal/codegen/jit_executor.go` | 19 | Linker/clang invocation error paths. |
| 81.3% | `internal/codegen/llvm.go` | 81 | Many helpers; long tail of small gaps. |
| 81.5% | `internal/codegen/iterator_generation.go` | 10 | Stream-fusion edge cases. |
| 83.9% | `internal/ast/builder_match.go` | 11 | Some pattern shapes missed. |
| 84.4% | `internal/ast/builder_literals.go` | 15 | Float/string literal edge cases. |
| 85.0% | `internal/codegen/builtin_registry.go` | 20 | Registry lookup misses. |
| 85.4% | `internal/ast/builder_arguments.go` | 2 | Named-argument validation edge cases. |
| 86.5% | `internal/ast/builder_calls.go` | 11 | A few call shapes. |
| 88.2% | `cmd/osprey/main.go` | 14 | Top-level flag handling. |
| 89.2% | `internal/ast/builder_core.go` | 7 | Almost-complete. |
| 89.9% | `internal/ast/builder_statements.go` | 12 | Almost-complete. |

## Untestable Code (Document & Exempt)

The following files contain code that **cannot be meaningfully exercised by `go test`**. Plan: add to a generated-or-bridge exclusion list in the `_test` Makefile target so they don't drag the average down.

### `internal/ast/ast.go` and `internal/ast/expressions.go` — Interface Marker Methods

These methods exist purely to satisfy interface type assertions:

```go
func (i *ImportStatement) isStatement() {}
func (pe *PerformExpression) isExpression() {}
```

They have empty bodies. Calling them in a test wouldn't prove anything (they do nothing). Go's coverage tool counts them as 0% because they are never called via the interface — interface dispatch is invisible to the coverage instrumenter on these no-op methods.

**Action**: refactor these into a single embedded struct/marker so there's only one of each (DRY) OR exclude `ast.go` and `expressions.go` from coverage measurement after confirming they only contain marker methods.

### `internal/codegen/websocket_bridge.go` — CGo Callbacks Invoked From C

These are CGo-exported functions called by the C WebSocket runtime, never from Go code:

```go
//export osprey_handle_websocket_connection
func osprey_handle_websocket_connection(serverID C.int, ...) C.int { ... }
```

Go's coverage tool only instruments calls made from Go. C-to-Go calls bypass instrumentation entirely.

**Action**: write Go-level integration tests that spin up an Osprey WebSocket server, connect to it from a Go client, and assert behavior — that exercises the C runtime which calls back into these functions. Their internal logic IS observable via the assertions on the connection lifecycle, even though the coverage tool won't credit them. Alternatively, exclude this file from coverage (it's a thin bridge, ~50 LOC).

## Ratcheting Plan

Each step below should land as a single PR raising `default_threshold` in `coverage-thresholds.json`. The rule "monotonically increasing — only ratchet UP" means we never go below the previous threshold.

### Stage 1 — 78% → 82% (quick wins)

- Add unit tests for `internal/logging/logging.go` (49% → 90%+). 6 functions, all pure.
- Add unit tests for `internal/language/descriptions/functions.go` exported getters (77% → 100%).
- Cover the missing branches in `internal/cli/main_interface.go` flag combos.

### Stage 2 — 82% → 86% (codegen breadth)

- Add example programs covering missing forms in `expression_generation.go`, `type_inference.go`, `effects_generation.go`. Each new `.osp` example with a `.expectedoutput` lifts coverage on the codegen path it exercises.
- Add unit tests for error wrappers in `errors.go` so the error formatting paths are checked (not just the happy path).

### Stage 3 — 86% → 90% (HTTP and fibers)

- Add HTTP server/client examples to `examples/tested/effects/` to lift `http_generation.go` from 54%.
- Add channel-select examples covering `fiber_generation.go` select paths.
- Add process-spawn variants to lift `system_generation.go`.

### Stage 4 — exclude untestable code, lock at 90%

- Refactor or exclude AST marker methods.
- Exclude or integration-test `websocket_bridge.go`.
- Update `_test` Makefile target's grep filter to drop these files alongside the generated parser.
- Re-measure; threshold should sit above 90% organically.

## TODOs

- [ ] Stage 1a: add unit tests for `internal/logging/logging.go` (target: lift file coverage from 49% to 90%+).
- [ ] Stage 1b: add unit tests for `internal/language/descriptions/functions.go` exported getters.
- [ ] Stage 1c: add unit tests for `internal/cli/main_interface.go` flag-combination branches.
- [ ] Ratchet `default_threshold` from 78 to 82 once Stage 1 lands.
- [ ] Stage 2a: write `.osp` examples covering missing branches in `expression_generation.go`.
- [ ] Stage 2b: write `.osp` examples covering missing branches in `type_inference.go`.
- [ ] Stage 2c: write `.osp` examples covering shallow-handler / resumption paths in `effects_generation.go`.
- [ ] Stage 2d: add unit tests for `errors.go` wrapper formatting.
- [ ] Ratchet `default_threshold` from 82 to 86 once Stage 2 lands.
- [ ] Stage 3a: add HTTP server/client `.osp` examples to lift `http_generation.go` above 80%.
- [ ] Stage 3b: add channel-select `.osp` examples to lift `fiber_generation.go` above 80%.
- [ ] Stage 3c: add process-spawn variant `.osp` examples for `system_generation.go`.
- [ ] Ratchet `default_threshold` from 86 to 90 once Stage 3 lands.
- [ ] Stage 4a: refactor AST marker methods (`isStatement`/`isExpression`) into an embedded marker, OR exclude `internal/ast/ast.go` and `internal/ast/expressions.go` marker-only sections from coverage.
- [ ] Stage 4b: write Go-level WebSocket integration test that exercises `websocket_bridge.go` end-to-end.
- [ ] Stage 4c: extend `_test` Makefile filter to drop confirmed-untestable bridge files (alongside the generated parser exclusion).
- [ ] Ratchet `default_threshold` from 90 to whatever organic value Stage 4 produces; keep it there.
