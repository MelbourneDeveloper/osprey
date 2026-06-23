# Plan 0009 — LSP Context-Awareness & Cross-File Resolution

**Subsystem:** `crates/osprey-lsp` (Rust language server)
**Status:** Partially implemented
**Spec:** [0020-LanguageServerAndEditors.md](../specs/0020-LanguageServerAndEditors.md) (`[LSP-CAPABILITIES]`)

## Summary

All advertised capabilities are wired and diagnostics, go-to-definition, find
references, and document symbols are solid. Three capabilities — **hover**,
**completion**, and **signature help** — are functional but shallow: they resolve
only at identifier positions, ignore the surrounding syntactic context, and never
look beyond the current file. Finishing this plan makes them context-aware and
workspace-aware.

## What works today

- Diagnostics (syntax + type), go-to-definition, find-references, document symbols
  — [crates/osprey-lsp/src/features.rs](../../crates/osprey-lsp/src/features.rs),
  [crates/osprey-lsp/src/diagnostics.rs](../../crates/osprey-lsp/src/diagnostics.rs),
  [crates/osprey-lsp/src/analysis.rs](../../crates/osprey-lsp/src/analysis.rs).
- Hover for declared symbols + builtins; completion of keywords + document
  symbols; signature help with active-parameter tracking inside an open call.
- VSCode client delegates cleanly over stdio —
  [vscode-extension/client/src/extension.ts](../../vscode-extension/client/src/extension.ts).

## Gaps (all three are "works narrowly")

1. **Hover** — only at identifier characters; nothing on type annotations,
   keywords, operators, or expressions (no inferred-type hover). Single-file only.
   [features.rs](../../crates/osprey-lsp/src/features.rs) `hover`.
2. **Completion** — ignores cursor position/context; returns all keywords + all
   symbols everywhere (offers `fn` snippet inside a type annotation, etc.); no
   member completion after `.`; single-file only.
   [features.rs](../../crates/osprey-lsp/src/features.rs) `completion`.
3. **Signature help** — only when the cursor is inside an already-open call; not on
   the function name itself; single-file only.
   [features.rs](../../crates/osprey-lsp/src/features.rs) `signature_help`.

Root cause for the single-file limitation: every feature takes only `text: &str`;
there is no workspace/file-graph context — [crates/osprey-lsp/src/model.rs](../../crates/osprey-lsp/src/model.rs)
`Query`, [crates/osprey-lsp/src/engine.rs](../../crates/osprey-lsp/src/engine.rs)
`answer`.

## Implementation plan

1. **Thread position + a lightweight syntactic context** into completion (and
   hover). Classify the cursor: top-level, inside a type annotation, inside a
   pattern, after `.`, inside a call. Filter keyword/symbol suggestions by what is
   valid there.
2. **Inferred-type hover.** Reuse `osprey_types` inference results to show the type
   of the expression under the cursor (not just declared names).
3. **Signature help on the function name**, not only inside the parens; broaden
   `enclosing_call` to start from the callee identifier.
4. **Introduce a workspace index.** Add a multi-file model to `Query`/`engine` so
   definitions, references, hover, completion, and signature help resolve symbols
   declared in other files. Build it incrementally on didOpen/didChange.
5. **Member completion after `.`** using record/type member info from the type
   table.

## Testing

- Extend the LSP unit tests in
  [crates/osprey-lsp/src/features.rs](../../crates/osprey-lsp/src/features.rs):
  context-filtered completion (no `fn` snippet inside a type annotation), hover on
  an expression returning an inferred type, signature help triggered on the
  function name, and a two-file cross-reference resolution.

## Risks / considerations

- The workspace index is the largest piece; land the single-file context-awareness
  (1–3) first as independent wins, then the index (4–5).
- Keep the `initialize` capabilities honest — only advertise what is implemented
  ([wire.rs](../../crates/osprey-lsp/src/wire.rs)).
- Out of scope (currently neither spec'd nor implemented): rename, formatting,
  semantic tokens, folding, code actions, inlay hints — track separately if
  desired.

## TODO

- [ ] Thread cursor position + syntactic-context classification into completion.
- [ ] Context-filter completion suggestions (type annotation / pattern / after-dot
      / call).
- [ ] Inferred-type hover via `osprey_types` results (beyond declared names).
- [ ] Trigger signature help on the function name, not only inside parens.
- [ ] Add a multi-file workspace index to `Query`/`engine`; make
      definition/references/hover/completion/sig-help cross-file.
- [ ] Member completion after `.` from type-table member info.
- [ ] Extend `features.rs` tests (context completion, inferred hover, name-position
      sig help, two-file resolution).
- [ ] `make ci` green.
