# Plan 0013 — ML Flavor Frontend

## Summary

Add the **ML flavor** — a layout-based, curry-by-default source surface — as a
second frontend **alongside** the existing Default (brace) flavor, not as a
replacement. Both frontends lower to the same `osprey_ast::Program`; everything
from type inference onward is shared and flavor-blind. The normative contract is
[spec 0023 — Language Flavors](../specs/0023-LanguageFlavors.md); the ML surface
is [spec 0024 — ML Flavor Syntax](../specs/0024-MLFlavorSyntax.md).

This plan supersedes the earlier "one canonical layout form, remove braces"
rollout drafts. Osprey keeps both surfaces permanently. The work is therefore
**additive**: a new parser, a new lowerer, a flavor selector, and one
shared-core feature the ML examples depend on — never a migration that rewrites
the Default flavor out of existence.

**Implementation decision — hand-written Rust layout frontend.** The ML
frontend is implemented as a **hand-written Rust layout lexer +
recursive-descent (Pratt / precedence-climbing) parser** in
`crates/osprey-syntax/src/ml/` (`token.rs`, `lexer.rs`, `cst.rs`, `parser.rs`,
`lower.rs`, `mod.rs`). The parser produces an ML **concrete syntax tree (CST)**;
a separate lowerer (`lower.rs`) converts it to canonical `osprey_ast::Program`
(clean **CST→AST separation**). The lexer derives
layout markers (`Indent`/`Dedent`/`Newline`) from the **offside rule**
(Landin 1966) via an explicit indentation stack, with bracket depth suppressing
layout inside parentheses. This **supersedes** the earlier plan of a
`tree-sitter-osprey-ml` grammar with an external C scanner. Rationale: the
offside rule is naturally expressed with an explicit indent stack in safe Rust;
it stays panic-free / `Result`-returning and unit-testable (project rules), with
no `unsafe` C and no codegen-tool build dependency. Per
[`[FLAVOR-BOUNDARY]`](../specs/0023-LanguageFlavors.md#the-one-law) the parser
**mechanism** is a below-the-AST, flavor-internal concern, so this swap does not
change the architecture (many CSTs, one AST). The tree-sitter + `scanner.c`
approach is retained as a documented **fallback (escape hatch)** in Phase 2. The
parsing techniques are cited in
[spec 0024 References](../specs/0024-MLFlavorSyntax.md#references).

**Current state.** Phase 1 (flavor frontend seam) and Phase 4 (flavor
selection) are **implemented and green**, and the diff harness discovers
`.ospml` additively. Phases 2–3 (ML lexer/parser/lowerer) are in active
development. Phase 0 (first-class handler values + effects) remains deferred, so
ML handler/effect syntax errors loudly until it lands.

## Why this is cheap (and where it is not)

The post-AST pipeline is already flavor-agnostic by construction:

- The type checker `check_program` / `infer_program`
  (`crates/osprey-types/src/check.rs:480`/`:493`) and code generator
  `compile_program` (`crates/osprey-codegen/src/lower.rs:20`) consume **only**
  `osprey_ast::Program` and the inferred type tables. Neither imports
  `osprey_syntax` or `tree_sitter`. No string `"flavor"` exists in the compiler.
- The Default lowerer (`crates/osprey-syntax/src/lower.rs`, `…/expr.rs`) already
  walks generic CST nodes by `kind()` and field name, so a second lowerer reuses
  the canonical AST vocabulary directly.
- **Currying needs no core change.** `Type::Fun` (`…/osprey-types/src/ty.rs:67`)
  is flat multi-arity; a curried function is nested `Fun` + nested one-param
  `Expr::Lambda` + nested one-arg `Expr::Call` — all implemented today
  (lambdas-as-values: [plan 0002](0002-codegen-generic-function-values.md)). The
  ML lowerer does the currying desugar; the checker and codegen are untouched.

The genuinely new work is two things: **(a)** a layout-sensitive parser — a
hand-written Rust layout lexer + recursive-descent (Pratt /
precedence-climbing) parser in `crates/osprey-syntax/src/ml/`, deriving layout
from the offside rule via an explicit indentation stack — and **(b)** one
shared-core feature — **first-class handler values + multi-install** — because
`Expr::Handler { effect, arms, body }` (`crates/osprey-ast/src/lib.rs:451`) fuses
construction and installation and cannot express `db = handler Db …; handle db
log do body`. That feature is flavor-neutral and lands first.

## Architecture (grounded)

| Stage | Today | After |
| --- | --- | --- |
| entry | `parse_program(src)` (`osprey-syntax/src/lib.rs:37`) | `parse_program_with_flavor(src, flavor)`; `parse_program` = Default wrapper |
| parse | tree-sitter brace grammar (`tree-sitter-osprey/`) | + hand-written Rust layout lexer + recursive-descent parser (`osprey-syntax/src/ml/`); tree-sitter + `scanner.c` retained as fallback |
| lower | `Lowerer` (`lower.rs`/`expr.rs`) → `Program` | + ML `lower.rs`: ML CST → the same `Program` (the parser builds the CST) |
| select | n/a | CLI flag > marker > extension > Default (`osprey-cli/src/main.rs:119`/`:200`) |
| check/codegen | `Program`-only, flavor-blind | **unchanged** |

## Phase 0 — Shared-core: first-class handler values

Flavor-neutral. Lands before the ML frontend because the ML (and new Default)
examples depend on it. See
[FLAVOR-HANDLER-VALUE](../specs/0023-LanguageFlavors.md#shared-core-additions).

TODO:

- [ ] Add `Expr::HandlerValue { effect, arms }` and
      `Expr::Install { handlers: Vec<Expr>, body }` to `osprey-ast`.
- [ ] Make the existing `Expr::Handler { effect, arms, body }` sugar for
      `Install { [HandlerValue { … }], body }` so all current Default programs
      keep compiling unchanged.
- [ ] Add a `Handler E` type to `osprey-types`; check arm/operation coverage.
- [ ] Type-check `Install` handler lists; detect duplicate installed handlers.
- [ ] Preserve handler-owned `mut` state on the handler value
      ([Algebraic Effects](../specs/0017-AlgebraicEffects.md) `[EFFECTS-HANDLER-STATE]`).
- [ ] Codegen: a runtime handler-value representation; lower `Install` of N
      values to nested handler installation; preserve behaviour across the C
      HTTP-callback and fiber boundaries; keep `resume` working.
- [ ] Default-flavor surface for the feature: `let h = handler E { … }` value
      form and multi-handler `handle h1 h2 in { body }`; grammar + lowerer.
- [ ] Tests: handler value bound/returned/passed; state isolation vs sharing;
      multi-install; existing effect examples still pass byte-for-byte.

## Phase 1 — Flavor frontend seam

**Implemented and green.** No behaviour change; Default stays the default.

TODO:

- [x] Add `enum Flavor { Default, Ml }` and `flavor: Flavor` on `Parsed`
      (`osprey-syntax/src/lib.rs:28`).
- [x] Add `parse_program_with_flavor(src, flavor) -> Parsed`; keep
      `parse_program` as the `Flavor::Default` wrapper.
- [x] Define the `FlavorFrontend` trait (`parse_tree` / `lower` /
      `collect_errors`); reorganise the current code as `default_frontend`.
- [x] Thread flavor through the interpolation re-entry (`expr.rs`
      `parse_fragment`, which recurses into `parse_program`).
- [x] Update callers (CLI, LSP, tests) to pass a flavor; all default to
      `Default`.

## Phase 2 — ML layout lexer + recursive-descent parser

Hand-written Rust frontend in `crates/osprey-syntax/src/ml/` (`token.rs`,
`lexer.rs`, `cst.rs`, `parser.rs`, `lower.rs`, `mod.rs`): the parser builds an ML
**concrete syntax tree (CST)** and a separate `lower.rs` converts the CST to
canonical `osprey_ast::Program` (clean **CST→AST separation**). The tree-sitter + `scanner.c` approach is the
documented **fallback (escape hatch)** below, not the primary path.

TODO:

- [ ] Layout lexer (`lexer.rs` + `token.rs`): derive `Indent` / `Dedent` /
      `Newline` from the **offside rule** via an explicit indentation stack;
      **bracket depth suppresses layout inside parentheses**; ignore blank and
      comment-only lines; preserve row/column on every token. Panic-free and
      `Result`-returning; unit-tested.
- [ ] ML CST types (`cst.rs`): surface nodes that preserve ML spelling —
      multi-parameter `funDef` heads, whitespace `application` as a callee +
      argument list (not yet nested), layout `block`/`match`/record, `effect`,
      `handler E` value, `handle … do`, `\… => …` lambdas. **Not** desugared.
- [ ] Recursive-descent parser (`parser.rs`): tokens → ML CST. Layout `block`,
      `funDef` heads, `:=` mutation, whitespace `application`, layout `match`.
- [ ] Lowerer (`lower.rs`): ML CST → canonical `osprey_ast::Program`. The
      **currying desugar lives here** (multi-param head → nested one-param
      `Lambda`; application list → nested one-arg `Call`), plus `${…}`
      interpolation. Clean CST→AST separation, symmetric with the Default flavor.
- [ ] Pratt / precedence-climbing expression layer: right-associative `->`,
      left-associative application; one binding-power table for the ML operators.
- [ ] Rust unit tests for indentation, match/handler arms, and edge cases (blank
      lines, comments, trailing newlines, tabs vs spaces, bracketed multi-line
      expressions where layout is suppressed).
- [ ] Module wiring: `mod ml` under `osprey-syntax`; no external build step, no
      `unsafe`, no codegen-tool dependency.

> **Escape hatch (documented fallback, not the primary path).** If the
> hand-written layout frontend becomes onerous or accrues parsing bugs we cannot
> tame, we fall back to a `tree-sitter-osprey-ml` grammar with an external
> `INDENT`/`DEDENT`/`NEWLINE` `scanner.c` (an indentation-stack scanner the brace
> grammar has never needed — `tree-sitter-osprey/` ships no `scanner.c` today),
> a tree-sitter grammar for the ML rules, a tree-sitter corpus test suite, and a
> separate `MlLowerer`. The boundary law
> ([`[FLAVOR-BOUNDARY]`](../specs/0023-LanguageFlavors.md#the-one-law)) makes the
> parser mechanism a flavor-internal swap that leaves the AST and everything
> above it untouched.

## Phase 3 — ML lowerer (CST → canonical AST)

Obeys the [lowering contract](../specs/0023-LanguageFlavors.md#the-lowering-contract).

TODO:

- [ ] `MlLowerer` producing `osprey_ast::Program`; preserve spans + doc comments;
      generated nodes carry the source span they desugar from.
- [ ] Bindings: `x = e` → `Let{mutable:false}`; `mut x = e` →
      `Let{mutable:true}`; `x := e` → `Assignment`.
- [ ] **Currying desugar** ([FLAVOR-CURRY](../specs/0023-LanguageFlavors.md#currying-canonicalisation)):
      `f x y = body` → one-param binding returning nested one-param `Lambda`;
      `f a b` → nested one-arg `Call`. Verify it equals the Default explicit-curry
      AST and differs from the Default multi-param AST.
- [ ] Effects: `op : P => R` → `EffectOperation { parameters:[P], return_type:R }`.
- [ ] Handlers: `handler E` → `HandlerValue`; `handle a b do body` → `Install`.
- [ ] Match: layout arms → `Match`/`MatchArm`; `Success value` →
      `Constructor { fields:["value"] }`.
- [ ] Records: layout block → `TypeConstructor`; layout update → `Update`.
- [ ] Diagnostics: same-scope `=` rebinding, write-to-immutable, unknown
      effect/operation — flavor-aware fix wording (`:=` vs `mut`/`=`).

## Phase 4 — Flavor selection wiring

**Implemented and green.**

TODO:

- [x] CLI `--flavor default|ml` on `Cli` (`osprey-cli/src/main.rs:34`), parsed in
      `parse_args` (`:119`); update `USAGE` (`:25`).
- [x] File marker `// osprey: flavor=ml` via the `directive` parser (`:521`),
      read in `run` (`:200`) before parsing.
- [x] Extension detection: `.ospml` ⇒ ML, `.osp` ⇒ Default (`Path::extension`).
- [x] Precedence flag > marker > extension > Default; **error** (hard, not a
      silent guess) when extension and marker disagree.
- [x] Diff harness (`crates/diff_examples.sh`) discovers `.ospml` **additively**
      and resolves flavor by extension; existing `.osp` discovery unchanged.
- [ ] Optional `osprey.toml` `flavor` key (deferred; not in the current
      precedence chain).
- [ ] LSP resolves the same precedence per document.

## Phase 5 — Tests, examples, equivalence

TODO:

- [ ] **LOADS of working `.ospml` tested examples** under `examples/tested/ml/`,
      each with a byte-for-byte `.expectedoutput`. Cover curried functions and
      partial application, `=>` effect operations, first-class handler values
      with owned `mut` state, `handle … do`, layout match, layout records,
      bindings/mutation, and string interpolation — concise files mixing many
      constructs. Discovered additively by `crates/diff_examples.sh`
      (`.ospml`, flavor by extension, `diff_examples.sh:23`).
- [ ] **No regressions: ALL existing `.osp` examples must continue to pass
      byte-for-byte.** `.ospml` discovery is purely additive; the Default harness
      output must not change for any current fixture.
- [ ] **Cross-flavor equivalence test**
      ([FLAVOR-TEST](../specs/0023-LanguageFlavors.md#cross-flavor-equivalence-tests)):
      parse a `.osp`/`.ospml` pair, strip spans + generated ids, assert canonical
      ASTs equal (equivalent bucket) or differ (non-equivalent bucket). Implement
      the comparison in Rust, not shell.
- [ ] ML must-reject cases under `examples/failscompilation/` (flavor resolved by
      extension/marker); keep the `FC_EXPECTED_ESCAPES` ratchet honest.
- [ ] Decide the ML extension story for negative cases (`.ospml` + marker vs a
      dedicated extension); document it in `examples/README.md`.
- [ ] WASM harness (`diff_wasm_examples.sh`): run any portable ML examples; keep
      the feature-gap SKIP classification.

## VS Code extension (VSIX) — hard requirement

**Explicit project-owner requirement. This is its own deliverable, not a Phase 6
footnote.** The published/built VSIX (`nimblesite.osprey`) must ship full ML
flavor support. Checklist:

- [ ] **Register the ML language.** Add `.ospml` to the extension's
      `contributes.languages` and register an `osprey-ml` language id (distinct
      from the existing `osprey`/`.osp` Default id).
- [ ] **ML TextMate / syntax grammar.** A dedicated ML grammar covering:
      keywords `mut`, `match`, `effect`, `handler`, `handle`, `do`, `perform`,
      `type` (and `true`/`false`); operators `:=`, `->`, `=>`, `\`; `//` line
      comments; strings with `${…}` interpolation; binding/function heads
      (`name = …`, `name param* = …`); and effect-operation lines
      `name : T => R`.
- [ ] **Layout-aware language configuration.** A separate
      `language-configuration` for `osprey-ml`: **no `{}` auto-pairing**;
      indentation `onEnter` rules so layout blocks indent correctly; keep `()`
      auto-pairing for grouping.
- [ ] **ML snippets.** An `osprey-ml` snippet set (binding, function, `effect`
      block, `handler` value, `handle … do`, `match`, layout record).
- [ ] **Commands include the ML flavor.** Ensure the run / compile / check
      commands and any "new file" / language-picker UI offer and handle the ML
      flavor (`.ospml`, `osprey-ml`), not just Default.
- [ ] **Ship it all in the VSIX.** The built/published VSIX bundles the ML
      grammar, language-configuration, snippets, language registration, and
      command wiring — verified in the packaged extension, not just the dev tree.

## Phase 6 — Tooling

TODO:

- [ ] VS Code ML editor support: see the dedicated
      [VS Code extension (VSIX)](#vs-code-extension-vsix--hard-requirement)
      checklist above (ML grammar, layout-aware config, snippets, command
      wiring); add folding and highlighting for `handler`, `handle`, `do`, `:=`,
      `=>`.
- [ ] LSP: hover/completion/signature help rendered in the **authoring** flavor;
      completion around effect operations and handler arms; curried-function
      signature help.
- [ ] Formatter: format within a flavor. Optional `osprey convert` to
      transliterate Default ⇄ ML (separate from the formatter).

## Phase 7 — Docs

TODO:

- [ ] Tag docs/website examples by flavor; mirror specs 0023/0024 to the website
      spec generator (`website/scripts/copy-spec.js`).
- [ ] Add flavor cross-reference notes to the existing language specs that gain a
      second spelling: `0003-Syntax`, `0005-FunctionCalls`, `0007-PatternMatching`,
      `0008-BlockExpressions`, `0017-AlgebraicEffects`.
- [ ] Update `examples/README.md` with the `.osp`/`.ospml` convention.

## Risks

- **ML lowerer must be its own exhaustive matcher.** The hand-written ML
  `lower.rs` converts the ML CST to canonical AST; it must produce only canonical
  nodes and never reuse the Default `Lowerer`'s `kind()` matching (whose wildcard
  arms on unknown kinds would silently corrupt the AST). (frontend-parse map)
- **Layout-lexer correctness.** Indentation tracking across tabs/spaces, blank
  lines, comments, trailing newlines, and bracket-suppressed layout is the
  hardest single piece; budget for it and cover the hand-written lexer with Rust
  unit tests. (frontend-parse map)
- **Currying conflation.** Default multi-param and ML curried functions must stay
  distinct in the AST; the golden non-equivalent bucket guards this. (types map)
- **Diagnostic hardcoding.** Existing fix messages assume Default spelling; ML
  needs its own fix wording behind the flavor-blind semantic code. (cli map)
- **Escape-hatch drift.** If the tree-sitter + `scanner.c` fallback is ever
  taken, it must remain a flavor-internal swap that produces the identical
  canonical AST; rely on the cross-flavor equivalence test to catch semantic
  drift. (frontend-parse map)

## Acceptance

- A `.ospml` program with curried functions, `=>` effect operations, first-class
  handlers, and `handle … do` compiles, runs, and matches its `.expectedoutput`
  byte-for-byte under `make test`.
- The equivalent-bucket golden tests prove Default explicit-curry ≡ ML curry and
  Default `handle … in` ≡ ML `handle … do` at the canonical AST.
- The non-equivalent-bucket golden tests prove Default multi-param ≢ ML curry.
- `grep` finds no flavor inspection in `osprey-types` or `osprey-codegen`.
- Every existing Default `.osp` example still passes unchanged.

## References

The parsing techniques behind the hand-written ML frontend — recursive-descent /
predictive parsing, the Pratt (precedence-climbing) expression layer, and the
offside-rule layout lexer — are cited with verified sources in
[spec 0024 References](../specs/0024-MLFlavorSyntax.md#references).
