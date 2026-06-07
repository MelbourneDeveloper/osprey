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
- **Phase 2 — port the compiler core (months, strangler-fig).** `osprey-syntax`→`-ast`→`-types`→`-codegen`
  →`-runtime-sys`→`-cli`, gated on golden + differential tests vs the Go oracle.
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

**Phase 0 COMPLETE; Phase 2 front-end landed (AST + CST→AST lowering + runtime-sys
FFI slice + CLI).** The Rust workspace builds, is `clippy`-clean, and all tests pass:
`tree-sitter-osprey` parses **45/45 valid examples with 0 ERROR nodes** (+6 corpus
tests), and `osprey-rs` lowers all 45 to a typed AST. **Not yet started: the hard
core** — `osprey-types` (HM inference, ~4K Rust LOC) and `osprey-codegen` (LLVM IR
text, ~10K LOC) — plus the Phase-1 LSP and the Phase-2 differential gate. Those are
the multi-week items this plan always scoped as "months"; the foundation they build
on is now in place and verified.

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
- [ ] 2.3 `crates/osprey-types`: port [`type_inference.go`](../../compiler/internal/codegen/type_inference.go) — **arena + index union-find**, match-exhaustiveness, effect rows. **(THE hard core — not started; ~4K Rust LOC.)**
- [ ] 2.4 `crates/osprey-codegen`: LLVM IR **text** emission — port `llvm.go` + `*_generation.go`; `builtin_registry.go` → const map; shell to `clang`. **(Not started; ~10K LOC, mostly transcription.)**
- [~] 2.5 [`crates/osprey-runtime-sys`](../../crates/osprey-runtime-sys/src/lib.rs): `cc`-built FFI to `compiler/runtime/*.c` (same hardening flags; **no C rewrite**). Self-contained FFI-pointer unit (`ffi_runtime.c`) linked + tested; pthread/OpenSSL units link the same way as their callers land.
- [~] 2.6 [`crates/osprey-cli`](../../crates/osprey-cli/src/main.rs): `osprey-rs` binary — `--ast` / `--check` / `--version` today; full clap surface (`--llvm --compile --run`, sandbox flags) grows with 2.3/2.4.
- [ ] 2.7 Differential test harness: Rust vs Go binary vs `.expectedoutput`, byte-for-byte, across all goldens. **(Blocked on 2.4.)**
- [ ] 2.8 **Gate:** 100% of `tested/` + `failscompilation/` pass; `make c-test` green

### Phase 3 — flip & retire (Gleam endgame)
- [ ] 3.1 `osprey-lsp` semantics go in-process (drop the Go sidecar)
- [ ] 3.2 Delete Go compiler, `go.mod`, ANTLR, generated `compiler/parser/` (**grammar count → 1**)
- [ ] 3.3 WASM target for front-end + checker → playground runs in-browser; retire [`webcompiler/`](../../webcompiler/)
- [ ] 3.4 Single binary: compiler + LSP + MCP + formatter; update [`RELEASING.md`](../RELEASING.md) + Shipwright manifest
- [ ] 3.5 Port specs verbatim once Rust front-end is golden-clean (separate task — **specs untouched until then**)

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
