# Plan: Go → Rust Migration (strangler-fig, tooling-first)

Touches the whole compiler. No spec changes yet — specs in [`../specs/`](../specs/) stay the
behavioural contract and are ported verbatim once the Rust front-end is golden-clean. Each design
choice is cited in [Authorities](#authorities).

## Decision

**Commit to Rust as the destination. Migrate incrementally, tooling-first. No big-bang rewrite.**
Endgame is one native binary — compiler + LSP + MCP + formatter — Gleam-style, living inside the
existing Rust stack (lspkit / Deslop / Shipwright). The sequence ships editor value in weeks and keeps
a working compiler at every commit.

## Why

- The rest of the portfolio is Rust: Basilisk (~212K LOC), **lspkit** ("one engine, two surfaces"
  LSP+MCP framework), **Deslop** (tree-sitter dup-detection over C#/Rust/Python/Dart), **SharpLsp**
  (Rust host + language sidecar), **Shipwright** (release contract — already mandated by this repo's
  `CLAUDE.md`). Osprey-in-Go is the lone holdout, locked out of all of it.
- The stated goal — *share the LSP tooling* — needs a `tree-sitter-osprey` grammar, **not** a compiler
  rewrite. That grammar is also step 1 of the eventual rewrite, so it is never wasted work.
- A **batch compiler** is the one place Go's GC is a non-issue ([esbuild][es] proved speed is
  architecture, not language). A **long-running LSP** is the one place Rust wins decisively — no GC
  pauses, incremental ([rust-analyzer][ra] salsa model), drops onto lspkit/Deslop. The seam falls
  cleanly *between* compiler and tooling.
- Syntax is currently duplicated **3×**: ANTLR [`compiler/osprey.g4`](../../compiler/osprey.g4)
  (370 lines), TextMate [`osprey.tmGrammar.json`](../../vscode-extension/syntaxes/osprey.tmGrammar.json),
  Monaco Monarch in the website playground. tree-sitter collapses this to one.
- Current LSP is ~1,335 LOC of TS ([`server.ts`](../../vscode-extension/src/server.ts)) shelling out to
  the Go binary via temp files + `execFile`. Fragile; replaced by a Rust LSP on lspkit sharing all infra.
- **Survives any path unchanged:** the C runtime ([`compiler/runtime/*.c`](../../compiler/runtime/),
  ~7.8K LOC — fibers, HTTP/WS, list/HAMT-map, JSON, string, system) is already FFI; the golden tests
  ([`examples/tested/`](../../compiler/examples/tested/) 39 cases,
  [`examples/failscompilation/`](../../compiler/examples/failscompilation/) 64 cases, each with
  `.expectedoutput`) are a binary-behaviour contract.

## Locked design decisions (driven by "Rust builds are shit")

1. **LLVM IR as TEXT, not `inkwell`/`llvm-sys`.** Linking LLVM into `cargo build` is the worst thing
   you can do to build times. The Go side already emits IR as text via the pure-Go [`llir/llvm`][llir]
   (no LLVM linkage), so the codegen port is *mechanical* and keeps LLVM out of the Rust build. Final
   native compile shells out to `clang`/`llc` exactly as today.
2. **Cargo workspace of small crates from day one** ([Servo/Vector pattern][ws]) → incremental builds
   stay ~5–15s, not monolithic-crate Rust hell. No single crate over ~5K LOC (mirrors the repo's
   "files under 500 LOC" rule at crate granularity).
3. **`tree-sitter-osprey` is the single canonical grammar** — compiler parser + LSP + highlighting +
   Deslop. **Fallback:** if compiler-grade error messages suffer, hand-write a recursive-descent parser
   for the *compiler only* (tree-sitter stays for editors regardless). The `failscompilation` golden
   tests are the forcing function that surfaces this early — they fail loudly if error text drifts.
4. **Differential testing.** During the port the Go `osprey` binary is the oracle; the golden tests are
   the contract. The Rust compiler must match every `.expectedoutput` byte-for-byte before anything flips.
5. **No big-bang.** A working, shippable compiler exists at every commit. Phases are independently useful.

## Target workspace layout

```
osprey/
  tree-sitter-osprey/        # the one grammar (own dir; publishable to npm + crates)
  crates/
    osprey-syntax/           # tree-sitter bindings + CST→AST lowering
    osprey-ast/              # AST enums  (mirror compiler/internal/ast/ast.go)
    osprey-types/            # Hindley-Milner inference (port of type_inference.go)
    osprey-codegen/          # LLVM IR *text* emission (port of llvm.go + *_generation.go)
    osprey-runtime-sys/      # -sys crate: FFI to the unchanged C runtime archives
    osprey-cli/              # the `osprey` binary (clap) — same flag surface as today
    osprey-lsp/              # implements lspkit EngineApi
    osprey-mcp/              # via lspkit-mcp (free once EngineApi exists)
```

## Technical detail per layer

### Grammar (`tree-sitter-osprey`)
Translate the 12 parser + 28 lexer rules of `osprey.g4` to `grammar.js`. ANTLR precedence climbing maps to
tree-sitter `prec.left/right` + `prec` numbers. The ANTLR-generated Go parser
([`compiler/parser/`](../../compiler/parser/), ~36K LOC) is **discarded**, not ported — it's machine output.
Ship `highlights.scm`, `locals.scm`, `folds.scm` queries; these replace the TextMate + Monaco copies.

### Front-end (`osprey-syntax` + `osprey-ast`)
`osprey-ast` is plain Rust enums mirroring `ast.go` (~520 LOC of node types). `osprey-syntax` walks the
tree-sitter CST → typed AST, replacing the listener-based builders in
[`compiler/internal/ast/`](../../compiler/internal/ast/) (~2.9K LOC across `builder_*.go`). Go's
`BaseospreyListener` traversal → an explicit recursive descent over CST named nodes (cleaner, no visitor
plumbing). UFCS dispatch (`builder_calls.go`) and string-interpolation splitting (`builder_interpolation.go`)
port node-for-node.

### Type inference (`osprey-types`) — the real work
Port [`type_inference.go`](../../compiler/internal/codegen/type_inference.go) (3.3K LOC HM unification +
constraint solving) → ~4K Rust LOC. Where Go leaned on casual pointer mutation of `TypeVar`s, use
`Rc<RefCell<…>>` union-find or an arena + indices (preferred: arena + `u32` ids — borrow-checker-friendly,
cache-friendly, no `Rc` churn). Match-exhaustiveness (`match_validation.go`) and effect-row checking
(`effects_generation.go` type side) live here. This is the one layer that is genuine intellectual work; budget
accordingly.

### Codegen (`osprey-codegen`)
Port the IR builders to a thin text emitter (`Vec<String>` / `fmt::Write`), one module per current file:
`llvm.go` (2.8K), `expression_generation.go` (3.2K), `function_signatures.go` (1.6K), `core_functions.go`,
`collection_codegen.go`, `fiber_generation.go`, `iterator_generation.go`, `string_functions.go`,
`http_generation.go`, `system_generation.go`. `builtin_registry.go` (1.5K) → a `phf`/`HashMap` const map.
Because the Go side already produces `.ll` text, this is largely transcription, not redesign. Emit `.ll` →
shell out to `clang` (mirrors `compilation.go` / `jit_executor.go`).

### Runtime FFI (`osprey-runtime-sys`)
`build.rs` compiles the existing `compiler/runtime/*.c` with the same hardening flags
(`-D_FORTIFY_SOURCE=2 -fstack-protector-strong`, warnings-as-errors) via the `cc` crate, links the static
archives, exposes `extern "C"` decls. **No C is rewritten.** The C memory-safety tests (`make c-test`) keep
running unchanged.

### CLI (`osprey-cli`)
`clap` reproducing today's surface from [`cli.go`](../../compiler/internal/cli/cli.go): positional file,
`--ast --llvm --compile --run --symbols --docs --hover`, and the sandbox flags
`--sandbox --no-http --no-websocket --no-fs --no-ffi` ([spec 0016](../specs/0016-SecurityAndSandboxing.md)).
`--symbols` (JSON) and `--hover` (markdown) keep their exact output shape — the Phase-1 LSP depends on them.

### LSP/MCP (`osprey-lsp` + `osprey-mcp`)
Implement lspkit's `EngineApi` once; lspkit vends it as both LSP (`lspkit-server` + `lspkit-vfs` rope docs +
`lspkit-live` watcher) and MCP (`lspkit-mcp`). Phase 1: syntax/outline/highlight/folding in-process via
`lspkit-treesitter`; deep semantics by shelling to the Go binary (`--symbols`/`--hover`/stderr diagnostics) —
the **SharpLsp sidecar pattern already in production**. Phase 3: swap the sidecar for in-process
`osprey-types`/`osprey-codegen`.

## Phases

- **Phase 0 — grammar (days).** `tree-sitter-osprey` + queries. Verify against every example.
- **Phase 1 — Rust tooling over the Go compiler (weeks, ships value).** `osprey-lsp`/`osprey-mcp` on lspkit,
  Go binary as semantic sidecar. Retire `server.ts` + TextMate + Monaco. Add Osprey to Deslop. Adopt
  Shipwright versioning. Compiler still Go; fast builds untouched.

- **Phase 2 — strangler-fig? HELL NO! ONE SHOT!!** `osprey-syntax`→`-ast`→`-types`→`-codegen`
  →`-runtime-sys`→`-cli`, gated on golden + differential tests vs the Go oracle. GO HARD!!!

- **Phase 3 — flip & retire (the Gleam endgame).** LSP goes in-process; delete Go compiler + `go.mod` +
  ANTLR (grammar count → 1). Add WASM target → playground runs in-browser, retiring `webcompiler/`. One
  binary: compiler + LSP + MCP + formatter.

## Verification

- **P0:** `tree-sitter test` + parse-all over `examples/*.osp` → zero ERROR/MISSING on valid files.
- **P1:** VSCode extension drives the Rust LSP (hover/completion/diagnostics/outline end-to-end); Deslop
  flags a planted dup across two `.osp` files; `osprey-lsp --version --json` satisfies Shipwright.
- **P2:** differential harness — for every golden case, Rust `--compile` + run output matches the Go binary
  **and** the `.expectedoutput`; `make c-test` stays green.
- **P3:** playground compiles Osprey in-browser via WASM; single binary exposes all subcommands; Go removed
  from the build graph.

## Risks

- **tree-sitter error quality** for compiler diagnostics → caught by `failscompilation` goldens; hand-written
  parser fallback for the compiler only.
- **HM inference port** is the genuine hard work (~4K Rust LOC); arena-of-indices design decided up front to
  dodge `Rc<RefCell>` borrow pain.
- **Now-or-never** applies *only* to the compiler core (P2), which grows with LOC. P0–P1 are correct to start
  immediately regardless of when P2 lands, and they make P2 cheaper.

---

## Status (current)

**Phase 0 COMPLETE. Phase 2 front-end + type checker + codegen all landed and
running end-to-end.** The differential gate is now **41 / 41** golden cases
byte-exact (whole-string `TrimSpace`) against the `.expectedoutput` oracle, via
[`crates/diff_examples.sh`](../../crates/diff_examples.sh) — **every**
`examples/tested` case passes, with **zero NOEXP**: the six cases the Go test
verified through in-code maps (`getExpectedOutputs`/`getHTTPExpectedOutputs` —
`function_composition_test`, `adventure_game`, `string_edge_cases`,
`http_client_example`, `http_create_client`, and the platform-branching
`callback_stdout_demo`) now carry on-disk `.expectedoutput` files (a single
source of truth both the Go file-first path and the Rust harness read;
`callback_stdout_demo` uses per-`uname` files since its subprocess error text +
exit code differ by OS). Closing the last of these required two codegen
additions: **lambda-lifting** (an inline lambda passed to a function-typed
parameter is lifted to a top-level `@__lambda_*` and passed as a code pointer —
[`genfn::lift_lambda`](../../crates/osprey-codegen/src/genfn.rs)) and **uniform
`i64` boxing of generic union-variant payloads** (`Full { value: T }` stores its
payload in one `i64` slot, matching the Go backend, instead of an `i8*` that
mis-dispatched `toString` into a segfault). The TUI (`tui_repo_table`) and the
live HTTP server (`http_server_example`: server + GET/POST/PUT/DELETE + 404 +
graceful shutdown) were both verified to match the Go binary byte-for-byte. The
whole Rust workspace is green under **maximum strictness**: `cargo build
--release`, `cargo clippy --workspace --all-targets` (clippy::all + pedantic +
restriction: no `unwrap`/`expect`/`panic`/`indexing`/`as`), `cargo fmt --check`,
and `cargo test --workspace` (all unit + corpus tests) all pass. **CI now runs
the Rust gate**: a `rust` job in [`ci.yml`](../../.github/workflows/ci.yml)
(alongside the still-required Go `ci` job, since the Go binary is still the
released artifact) builds the C runtime archives and runs fmt + clippy
(`-D warnings`) + workspace tests + the 41/41 differential harness on every PR
to `main`.

**Hardening round (post-41/41).** Five structural upgrades landed on top of the
behavioural parity above:

1. **The type checker now gates the pipeline.** `--run`/`--llvm` run
   `check_program` before codegen — an ill-typed program never compiles. Two
   checker false positives this exposed were fixed (a user union re-using
   `Success`/`Error` no longer hijacks patterns over a real `Result`; bool-literal
   arms over a `Result` are the ternary/Elvis truthiness test, yielding the
   unwrapped payload and covering both variants).
2. **The must-reject suite is wired and ratcheted.**
   [`crates/diff_examples.sh`](../../crates/diff_examples.sh) now also runs every
   `examples/failscompilation/*.ospo` and fails CI if more than
   `FC_EXPECTED_ESCAPES` ill-formed programs are accepted. Rejections went
   33 → **50 / 62** (checker wiring + ported validations: constructor field-set
   checking, unknown match variant, built-in redefinition, immutable-assignment).
   The 12 remaining escapes (effects safety ×5, `any`-position rules ×4,
   named-args arity, print-on-record, `comprehensive_bugs_demo`) are the ratchet
   — port a validation, decrease the constant.
3. **Fibers are REAL** — the same C runtime the language always used: `spawn e`
   lifts `e` into a no-arg `i64 ()` thunk for `fiber_spawn` (a pthread; effect
   handlers snapshot across), with closed-over locals spilled through per-spawn
   module globals; `await`/`fiberDone` are `fiber_await`/`fiber_done`; channels
   are `channel_create/send/recv`. `sleep` maps to `fiber_sleep(ms)` (an unmapped
   fall-through had linked libc `sleep(seconds)`).
4. **One type system.** `ProgramTypes` now carries *resolved* `Type`s for
   constructor fields (a generic field resolves to a type variable — the boxed
   `i64` rule falls out) and exported effect-operation signatures; the backend's
   written-string re-parsers (`ltype_of_name`, `parse_op_sig`) are deleted.
   Inference is the single source of type truth; codegen only maps `Type` →
   LLVM.
5. **The 411-line effects TUI (`examples/tui/api_browser.osp`) runs end-to-end**
   under `osprey-rs` — live GitHub API, animated fiber-overlap loader, raw-mode
   keys — matching the Go binary's screens byte-for-byte. Closing it needed the
   eight `[BUILTIN-TERM]` builtins in codegen and `Fiber`/`Channel` lowering as
   `i64` runtime ids.

- **`osprey-types` — the HM core — DONE and verified.** A complete Hindley-Milner
  engine: enum `Type` language, index-addressed union-find substitution with
  path-compressed `prune`/occurs-check, unification (with the Osprey rules: `any`
  wildcard, bare-collection generics, structural records, Result auto-unwrap **and**
  auto-wrap), let-polymorphism (generalize/instantiate, incl. top-level functions),
  a builtin registry, full expression/statement inference (arithmetic Result
  semantics, records, generic unions, lambdas, UFCS, effects), pattern inference and
  match exhaustiveness.
- **`osprey-codegen` — LLVM-IR-text backend — working END-TO-END across the
  language.** Emits textual LLVM IR (no inkwell), handed to `clang`. Now covers,
  on top of the compute core (literals, arithmetic-as-`Result<…,MathError>`,
  comparison, `print`/`toString`, sprintf interpolation, `let`/blocks, recursive
  functions + named-arg ordering, synthesized `main`): **records & union variants,
  object literals, `match` (literal / Result / union / Elvis), pattern field-bind
  by name, negative-literal patterns, lists + maps + 2D indexing, iterators
  (stream-fused range + eager list ops), fibers (eager-eval), algebraic effects
  (`perform`/`handle` via a dynamic handler stack), division-by-zero → `Error`,
  generic field-access by-name fallback.** Idiomatic Rust, **no panics,
  `Result<T,E>` throughout**; genuinely-unsupported nodes (a bare lambda used as a
  runtime value) **fail loudly**, never emit a placeholder or invalid IR.

The hardest intellectual core (HM inference) and the full pipeline (parse → check
→ LLVM IR → clang → run) are proven. **All five remaining subsystems have landed
and the differential gate is at 35 / 35** — file I/O, processes (first-class
function pointers + callback ABI), HTTP server/client + JSON + the `HttpResponse`
C-struct layout, SQLite via generic FFI (extern-fn signature publishing + Result
phi propagation), and generic monomorphisation (call-site inlining + concrete
per-instance record layouts + indirect higher-order calls + uniform `i64`
boxing). See [Remaining examples](#remaining-examplestested--the-last-11) for the
per-subsystem record.

## TODO

### Phase 0 — grammar (prereq)  ✅ DONE
- [x] 0.1 [`tree-sitter-osprey/grammar.js`](../../tree-sitter-osprey/grammar.js) from [`osprey.g4`](../../compiler/osprey.g4) (precedence via `prec.*`; left-recursive postfix call chain; `//` vs `/` and `///` vs `//` lexer precedence fixed). Rust bindings + `cc` build.
- [x] 0.2 [`highlights.scm`](../../tree-sitter-osprey/queries/highlights.scm) / [`locals.scm`](../../tree-sitter-osprey/queries/locals.scm) / [`folds.scm`](../../tree-sitter-osprey/queries/folds.scm) queries
- [x] 0.3 Parse-all harness ([`test/parse-all.js`](../../tree-sitter-osprey/test/parse-all.js)) → **45/45 valid examples, 0 ERROR nodes** (`failscompilation/` are `.ospo` error cases)
- [x] 0.4 [`test/corpus/osprey.txt`](../../tree-sitter-osprey/test/corpus/osprey.txt) — 6 corpus tests (fn, union, effect, lambda, match-destructure, pipe/UFCS), all green

### Phase 1 — Rust tooling over the Go compiler (ships value)
- [ ] 1.1 `crates/osprey-lsp`: implement lspkit `EngineApi`; in-process tree-sitter for syntax/outline/highlight/folding
- [ ] 1.2 Semantic sidecar: shell to Go binary — `--symbols` (JSON), `--hover`, stderr diagnostics (SharpLsp pattern)
- [ ] 1.3 `crates/osprey-mcp` via `lspkit-mcp` (vends the same `EngineApi`)
- [ ] 1.4 Repoint VSCode extension at the Rust LSP; **delete [`server.ts`](../../vscode-extension/src/server.ts)**
- [ ] 1.5 Replace TextMate + Monaco with tree-sitter highlighting (2 of 3 grammar copies gone)
- [ ] 1.6 Add `deslop-core/src/lang/osprey.rs` (Deslop repo) following the C#/Rust/Python/Dart pattern
- [ ] 1.7 Adopt Shipwright contract: `--version` / `--version --json` on the Rust binaries

### Phase 2 — port the compiler core (strangler-fig, Go = oracle)
- [x] 2.1 [`crates/osprey-ast`](../../crates/osprey-ast/src/lib.rs): `Stmt`/`Expr` enums mirroring [`ast.go`](../../compiler/internal/ast/ast.go) (exhaustively matchable for the checker/codegen ports)
- [x] 2.2 [`crates/osprey-syntax`](../../crates/osprey-syntax/src/lib.rs): CST→AST lowering (replaces `internal/ast/builder_*.go`) — all core constructs incl. UFCS method calls, named args, string interpolation; lowers **45/45 examples** clean; 7 unit tests
- [x] 2.3 [`crates/osprey-types`](../../crates/osprey-types/src/lib.rs): port of [`type_inference.go`](../../compiler/internal/codegen/type_inference.go) — index union-find (`ctx.rs`), unification with the Osprey rules (`unify.rs`), let-polymorphism (`env.rs`), builtin registry (`builtins.rs`), expr/stmt inference (`expr.rs`/`check.rs`), pattern inference + **match-exhaustiveness** (`pattern.rs`). **`--check` passes 44/45 examples; 26 tests.** (Effect *rows* tracked structurally, not yet row-unified — a follow-up.)
- [x] 2.4 [`crates/osprey-codegen`](../../crates/osprey-codegen/src/lib.rs): LLVM IR **text** emission (port of `llvm.go`/`*_generation.go`); shells to `clang`. **Works end-to-end across the language (`--run`):** type-driven signatures (string/record/float params + returns), arithmetic-as-`Result`, `print`/interpolation, `match` (literal/Result/union/Elvis), records & unions, object literals, lists/maps + 2D indexing, iterators, fibers, algebraic effects (`perform`/`handle`), division-by-zero → `Error`, the runtime-builtin table ([`extern_call.rs`](../../crates/osprey-codegen/src/extern_call.rs): file/process/HTTP/JSON), first-class function pointers + indirect higher-order calls + generic-function inlining + **lambda-lifting of value-position lambdas** ([`genfn.rs`](../../crates/osprey-codegen/src/genfn.rs)), uniform `i64` boxing of generic union-variant payloads, and the `HttpResponse` C-struct constructor. Drives **41/41** goldens byte-exact (2.9).
- [~] 2.5 [`crates/osprey-runtime-sys`](../../crates/osprey-runtime-sys/src/lib.rs): `cc`-built FFI to `compiler/runtime/*.c` (same hardening flags; **no C rewrite**). Self-contained FFI-pointer unit (`ffi_runtime.c`) linked + tested; pthread/OpenSSL units link the same way as their callers land.
- [~] 2.6 [`crates/osprey-cli`](../../crates/osprey-cli/src/main.rs): `osprey-rs` binary — `--ast` / `--check` / `--llvm` / `--run` / `--version` today; remaining clap surface (`--compile`, sandbox flags) grows with 2.4.
- [x] 2.7 Differential test harness: [`crates/diff_examples.sh`](../../crates/diff_examples.sh) — Rust `--run` vs `.expectedoutput`, whole-string `TrimSpace`, across all goldens; portable repo-root (runs in CI), per-`uname` fallback for platform-specific outputs. **41/41, 0 NOEXP.**
- [x] 2.8 **Gate:** 100% of `tested/` **passes (41/41, zero NOEXP)** — including the six cases the Go test held in in-code maps, now materialised as on-disk `.expectedoutput`. **CI runs the Rust gate** (a `rust` job in [`ci.yml`](../../.github/workflows/ci.yml): C-runtime build + fmt + clippy `-D warnings` + `cargo test` + the differential harness, on every PR to `main`). `failscompilation/` is now wired into the harness as a ratcheted must-reject suite (50/62, target 62/62); `make c-test` remains out of scope this round.
- [x] 2.9 **All 11 remaining goldens landed** (see [Remaining examples](#remaining-examplestested--the-last-11)):
  - [x] R2 **files — runtime symbols.** Name-map `readFile`→`read_file` / `writeFile`→`write_file` with the right `Result<…>` wrapping, via the [`extern_call.rs`](../../crates/osprey-codegen/src/extern_call.rs) builtin table. `file_io_json_workflow` passes.
  - [x] R2 **db×2 — sqlite FFI.** Root cause was the type checker never publishing extern signatures (`collect_extern` skipped `fn_sigs`), so every `Ptr` collapsed to `i64`; publishing the signature + propagating `result_inner` through `match` phis (`finish_phi`, which fixed a Result double-wrap) makes `sqlite_basics` and `database_effect` pass.
  - [x] R1 **types×3 — generic monomorphisation.** Done **without** per-expression type tables: a generic function is specialised by inlining its body at each call site ([`genfn.rs`](../../crates/osprey-codegen/src/genfn.rs)), a generic *record* is built with the concrete field types present at construction (per-instance layout like an object literal), function-valued params lower to indirect calls, and `any`/generic values box uniformly through `i64` (`ptrtoint`/`inttoptr`). `any_type_comprehensive`, `pure_hindley_milner_test`, `type_equality_comprehensive` pass.
  - [x] R3 **processes — first-class function pointers.** A bare top-level function name lowers to `bitcast <sig>* @name to i8*` (`expr::fn_pointer`); `spawnProcess`/`awaitProcess`/`cleanupProcess` map to `spawn_process_with_handler`/`fiber_await_process`/`fiber_cleanup_process`. `async_process_management` passes.
  - [x] R4 **http×3 — HttpResponse layout + HTTP runtime.** `HttpResponse` reconciled to the C struct (`status, headers, contentType, streamFd, isComplete, partialBody`) and emitted **tag-free, `bool` as `i8`** to match `runtime/http_shared.h`; the `http*`/`json*` builtins go through the `extern_call` table; `httpListen` receives the handler via the R3 code pointer. Also fixed string-escape lowering (`\r`, `\e`/ESC). `http_response_handle`, `http_server_example`, `tui_repo_table` pass.
  - [x] R5 **comprehensive_math — product decision (resolved).** The `.expectedoutput` `complex = 2` is a confirmed Go type-inference bug in 3+-operand `Result`-arithmetic chains (NOT a precedence bug — plain arithmetic is correct). Per the owner's call, the oracle was corrected to the arithmetically-correct `complex = 18`; the Go quirk was **not** replicated.
  - [x] **Regression gate (held throughout):** `cargo build --release`, `cargo clippy --workspace --all-targets`, `cargo fmt --check`, `cargo test --workspace`, and `crates/diff_examples.sh` all green; zero regressions at every step.

### Phase 3 — flip & retire (Gleam endgame)
- [ ] 3.1 `osprey-lsp` semantics go in-process (drop the Go sidecar)
- [ ] 3.2 Delete Go compiler, `go.mod`, ANTLR, generated `compiler/parser/` (**grammar count → 1**)
- [ ] 3.3 WASM target for front-end + checker → playground runs in-browser; retire [`webcompiler/`](../../webcompiler/)
- [ ] 3.4 Single binary: compiler + LSP + MCP + formatter; update [`RELEASING.md`](../RELEASING.md) + Shipwright manifest
- [ ] 3.5 Port specs verbatim once Rust front-end is golden-clean (separate task — **specs untouched until then**)

## `examples/tested` — the last 11 (all resolved)

Status: **35 / 35** byte-exact — every case below now passes. The analysis that
guided the work is kept as a record; each fell into one of **4 subsystems** (+ the
math product-decision). See the per-item record in the Phase-2 TODO (2.9) for how
each was closed.

### R1. Generic monomorphization — unlocks 3
**Examples:** `types/any_type_comprehensive`, `types/pure_hindley_milner_test`,
`types/type_equality_comprehensive`. **Symptom:** codegen aborts —
`invalid program: expected an integer, found a string/handle`.

**Root cause.** A generic record field (`type Generic<T> = { data: T }`) has its
written type `T`, which `types::ltype_of_name` maps to `LType::Ptr`. So
`makeGenericInt(42)` stores `inttoptr 42` and `makeGenericString("x")` stores a
string pointer — both `i8*`. At a use site (`${gen1.data}`, `gen1.data * 2`) the
backend can't tell an int payload from a string payload. There is **no
per-expression type info**: `osprey_types::ProgramTypes` publishes only
`functions`, `ctors`, `unions` — frozen tables — so the backend can't recover that
*this* `gen1.data` is `int` while *that* `gen3.data` is `string`.

**Fix (recommended = A).** **(A)** Add `expr_types: HashMap<ExprId, Type>` to
`ProgramTypes`; record each node's resolved type in `infer_expr` after the final
substitution; consume it in `gen_field_access`/`gen_object`/`gen_index` to pick
the concrete `LType` instead of `T → Ptr`. Needs stable AST ids (`id: u32` on
`Expr`, or a `Position`-keyed side table). Only approach that fixes the general
case. **(B)** Monomorphize by call site (read a function's concrete resolved
return type, attach a per-variable layout override) — cheaper, partial. **(C)**
Uniform `{ i8 tag, i64 payload }` boxing + runtime dispatch — avoids inference
plumbing, changes the ABI. Effort: ~1–2 days. Watch record/generic regressions.

### R2. Runtime FFI symbols — unlocks `files` (1); same shape unlocks `db` (2)
**Examples:** `files/file_io_json_workflow`; `db/sqlite_basics`,
`db/database_effect`. **Symptom:** link error
`Undefined symbols … "_readFile", "_writeFile"` (files); empty output / `create err`
(sqlite — runtime no-op).

**Root cause.** The program calls builtins (`readFile`/`writeFile`/JSON; sqlite
ops); the backend emits `call @readFile(...)` (unknown callees are auto-declared
as externs, so codegen succeeds) but the linked static runtimes
(`compiler/bin/libfiber_runtime.a`, `libhttp_runtime.a`) export no such symbol.
Purely **missing runtime symbols** at link time.

**Fix.** Grep `compiler/runtime/` for the real names (likely `osp_read_file`-style);
then either (a) a backend name-map (`runtime.rs`/`call_with_values`), (b) thin C
shims rebuilt into the static libs, or (c) implement in `osprey-runtime-sys`.
Ensure `osprey-cli`'s `link_args` includes the exporting archive. Effort: ~½ day
once symbols are located. Check the `.expectedoutput` embeds no absolute
paths/timestamps.

### R3. First-class function pointers + process runtime — unlocks `processes` (1)
**Example:** `processes/async_process_management`. **Symptom:**
`unknown name processEventHandler`.

**Root cause.** `spawnProcess("echo …", processEventHandler)` passes a *named
function* as a value; `gen_expr(Identifier("processEventHandler"))` finds no
binding/ctor → `unknown name`. The backend supports only inline/let-bound lambdas
and direct calls — no first-class function pointer. (Also needs a deterministic
`spawnProcess` runtime + callback ABI.)

**Fix.** In `gen_expr`'s `Identifier` arm, if the name is a known top-level
function, emit `bitcast <sig>* @name to i8*` (a code-pointer `Value`); add an
indirect-call path that bitcasts an `i8*` callee back to its fn-ptr type before
`call` — **mirror `effects::gen_perform`**, which already does this for handler
pointers. Keep `gen_user_call` direct calls unchanged. Effort: ~1–2 days incl.
runtime. **Prerequisite for R4.**

### R4. HTTP server/client runtime + `HttpResponse` builtin — unlocks 3
**Examples:** `http/http_response_handle`, `http/http_server_example`,
`http/tui_repo_table`. **Symptom:**
`invalid program: missing field 'body' for 'HttpResponse'`.

**Root cause.** `HttpResponse` is a *builtin* type; the registered ctor layout
expects `body` but the literal supplies `partialBody` — the two definitions are
out of sync. Beyond that, these need a **live HTTP server + client** producing
byte-exact output (`status=200`, `body=hello body`, header lookups, double-free →
`Error`).

**Fix.** Reconcile the `HttpResponse` field set with the examples
(`status, headers, contentType, streamFd, isComplete, partialBody`); link
`libhttp_runtime.a` and wire the `http*` builtins to its symbols; drive the
lifecycle deterministically. `httpListen(serverId, handleRequest)` passes a named
handler — **depends on R3 (do R3 first).** Highest risk (sockets, dispatch, port
binding, ordering). Effort: ~2–3 days.

### R5. Go parser precedence bug — `comprehensive_math` (1), do last / maybe never
**Symptom:** one line — `complex = 18` (ours, arithmetically correct) vs
`complex = 2` (expected). `fn complex(a,b) = match (a*2)+(b*3)-1 {…}`, `a=5,b=3` →
`10+9-1 = 18`; the `.expectedoutput` `2` is a **precedence bug in the Go ANTLR
parser**; our tree-sitter grammar parses it correctly.

**Fix — product decision, do not act unilaterally.** **(A)** Replicate the bug —
**strongly discouraged**: corrupts correct arithmetic precedence everywhere, would
break passing examples. **(B)** Regenerate the oracle (`complex = 18`) if the
contract is "match the *new* compiler," treating the Go file as the stale
artifact. Raise with the team: are `.expectedoutput` files frozen bug-for-bug Go
output, or the spec of correct behaviour?

**Dependency order:** `R2 files` → `R1 types×3` → `R3 fn-ptrs → R4 http×3`; `R2 db`
alongside files; `R5` last (decision).

## Authorities

[Gleam — Erlang→Rust rewrite for static types & a unified single-binary toolchain][gleam] ·
[esbuild — Go over Rust for a batch tool][es] · [rust-analyzer — incremental LSP model][ra] ·
[llir/llvm — pure-Go IR text emission][llir] · [Cargo workspaces for compile time][ws] ·
lspkit / Deslop / SharpLsp / Shipwright (sibling repos under `~/Documents/Code`).

[gleam]: https://gleam.run/frequently-asked-questions/
[es]: https://news.ycombinator.com/item?id=30079403
[ra]: https://rust-analyzer.github.io/
[llir]: https://github.com/llir/llvm
[ws]: https://corrode.dev/blog/tips-for-faster-rust-compile-times/
