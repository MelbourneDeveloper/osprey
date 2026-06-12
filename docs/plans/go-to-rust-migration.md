# Plan: Single-Binary Toolchain (LSP + MCP + WASM)

The Go → Rust compiler migration this plan originally tracked is **complete**:
the compiler is the Rust workspace (`crates/`, binary `target/release/osprey`),
the Go compiler / `go.mod` / ANTLR / generated parser are deleted, CI runs the
Rust gate only, and the differential harness holds **41/41** goldens byte-exact
with the `failscompilation/` must-reject ratchet (`FC_EXPECTED_ESCAPES`). The
full migration record lives in git history of this file.

What remains is the **Gleam-style endgame**: one native binary — compiler +
LSP + MCP + formatter — plus an in-browser playground via WASM.

## Remaining design intent

- `tree-sitter-osprey` stays the single canonical grammar (compiler parser,
  LSP, highlighting, Deslop). TextMate + Monaco copies still exist and must be
  retired.
- The LSP implements lspkit's `EngineApi` once; lspkit vends it as both LSP
  and MCP. Semantics run **in-process** over `osprey-syntax`/`osprey-types`
  (the original Go-sidecar phase is obsolete — there is no Go binary).
- LLVM stays out of the Rust build (IR as text, shell to `clang`), so the
  WASM target covers front-end + checker; native codegen remains host-only.

## TODO

- [ ] `crates/osprey-lsp`: implement lspkit `EngineApi` — tree-sitter
  syntax/outline/highlight/folding + in-process `osprey-types` diagnostics
- [ ] `crates/osprey-mcp` via `lspkit-mcp` (vends the same `EngineApi`)
- [ ] Repoint the VSCode extension at the Rust LSP (the old `server.ts` is
  already deleted)
- [ ] Replace TextMate + Monaco with tree-sitter highlighting (grammar
  count → 1)
- [ ] Add `deslop-core/src/lang/osprey.rs` (Deslop repo) following the
  C#/Rust/Python/Dart pattern
- [ ] WASM target for front-end + checker → playground runs in-browser;
  retire [`webcompiler/`](../../webcompiler/)
- [ ] Single binary: compiler + LSP + MCP + formatter; update
  [`RELEASING.md`](../RELEASING.md) + Shipwright manifest
