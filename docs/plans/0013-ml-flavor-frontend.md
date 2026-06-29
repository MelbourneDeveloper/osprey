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

The genuinely new work is two things: **(a)** a layout-sensitive parser (an
external scanner that the brace grammar has never needed), and **(b)** one
shared-core feature — **first-class handler values + multi-install** — because
`Expr::Handler { effect, arms, body }` (`crates/osprey-ast/src/lib.rs:451`) fuses
construction and installation and cannot express `db = handler Db …; handle db
log do body`. That feature is flavor-neutral and lands first.

## Architecture (grounded)

| Stage | Today | After |
| --- | --- | --- |
| entry | `parse_program(src)` (`osprey-syntax/src/lib.rs:37`) | `parse_program_with_flavor(src, flavor)`; `parse_program` = Default wrapper |
| parse | tree-sitter brace grammar (`tree-sitter-osprey/`) | + `tree-sitter-osprey-ml` with external scanner |
| lower | `Lowerer` (`lower.rs`/`expr.rs`) → `Program` | + `MlLowerer` → same `Program` |
| select | n/a | CLI flag > marker > extension > config > Default (`osprey-cli/src/main.rs:119`/`:200`) |
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

No behaviour change; Default stays the default.

TODO:

- [ ] Add `enum Flavor { Default, Ml }` and `flavor: Flavor` on `Parsed`
      (`osprey-syntax/src/lib.rs:28`).
- [ ] Add `parse_program_with_flavor(src, flavor) -> Parsed`; keep
      `parse_program` as the `Flavor::Default` wrapper.
- [ ] Define the `FlavorFrontend` trait (`parse_tree` / `lower` /
      `collect_errors`); reorganise the current code as `default_frontend`.
- [ ] Thread flavor through the interpolation re-entry (`expr.rs`
      `parse_fragment`, which recurses into `parse_program`).
- [ ] Update callers (CLI, LSP, tests) to pass a flavor; all default to
      `Default`.

## Phase 2 — ML grammar + layout scanner

TODO:

- [ ] Create `tree-sitter-osprey-ml` with an external `scanner.c` emitting
      `INDENT` / `DEDENT` / `NEWLINE` over an indentation stack; ignore blank and
      comment-only lines; preserve row/column on every token.
- [ ] Grammar rules: layout `block`, `funDef` binding heads, `:=` mutation,
      whitespace `application`, layout `match`, layout record, `effect` with
      `op : P => R`, `handler E` value, `handle … do`, `\… => …` lambdas.
- [ ] Right-associative `->`, left-associative application; ML precedence table.
- [ ] tree-sitter corpus tests for indentation, match/handler arms, and edge
      cases (blank lines, comments, trailing newlines, tabs vs spaces).
- [ ] Build wiring: independent `build.rs` / Cargo crate; regenerate `parser.c`.

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

TODO:

- [ ] CLI `--flavor default|ml` on `Cli` (`osprey-cli/src/main.rs:34`), parsed in
      `parse_args` (`:119`); update `USAGE` (`:25`).
- [ ] File marker `// osprey: flavor=ml` via the `directive` parser (`:521`),
      read in `run` (`:200`) before parsing.
- [ ] Extension detection: `.ospml` ⇒ ML, `.osp` ⇒ Default (`Path::extension`).
- [ ] Optional `osprey.toml` `flavor` key.
- [ ] Precedence flag > marker > extension > config > Default; **error** when
      extension and marker disagree.
- [ ] LSP resolves the same precedence per document.

## Phase 5 — Tests, examples, equivalence

TODO:

- [ ] Cross-flavor golden harness
      ([FLAVOR-TEST](../specs/0023-LanguageFlavors.md#cross-flavor-equivalence-tests)):
      parse a `.osp`/`.ospml` pair, strip spans + generated ids, assert canonical
      ASTs equal (equivalent bucket) or differ (non-equivalent bucket). Implement
      the comparison in Rust, not shell.
- [ ] Curated ML tested examples under `examples/tested/` (`.ospml`), each with a
      byte-for-byte `.expectedoutput`; extend `crates/diff_examples.sh` discovery
      to `.ospml` and resolve flavor by extension (`diff_examples.sh:23`).
- [ ] ML must-reject cases under `examples/failscompilation/` (flavor resolved by
      extension/marker); keep the `FC_EXPECTED_ESCAPES` ratchet honest.
- [ ] Decide the ML extension story for negative cases (`.ospml` + marker vs a
      dedicated extension); document it in `examples/README.md`.
- [ ] WASM harness (`diff_wasm_examples.sh`): run any portable ML examples; keep
      the feature-gap SKIP classification.

## Phase 6 — Tooling

TODO:

- [ ] VS Code: ML TextMate/highlight grammar; indentation language-config; folding;
      highlight `handler`, `handle`, `do`, `:=`, `=>`.
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

- **CST node-kind divergence.** The Default `Lowerer`'s exhaustive `kind()`
  matches fall through to wildcard arms on unknown kinds, silently corrupting the
  AST. The ML lowerer must be its own exhaustive matcher, not a patched Default
  one. (frontend-parse map)
- **External scanner correctness.** Indentation tracking across tabs/spaces, blank
  lines, comments, and trailing newlines is the hardest single piece; budget for
  it and cover it with corpus tests. (tree-sitter map)
- **Currying conflation.** Default multi-param and ML curried functions must stay
  distinct in the AST; the golden non-equivalent bucket guards this. (types map)
- **Diagnostic hardcoding.** Existing fix messages assume Default spelling; ML
  needs its own fix wording behind the flavor-blind semantic code. (cli map)
- **Two grammars, two pipelines.** Every shared grammar fix risks needing to be
  applied twice; keep the grammars genuinely independent and rely on the AST
  golden tests to catch semantic drift. (tree-sitter map)

## Acceptance

- A `.ospml` program with curried functions, `=>` effect operations, first-class
  handlers, and `handle … do` compiles, runs, and matches its `.expectedoutput`
  byte-for-byte under `make test`.
- The equivalent-bucket golden tests prove Default explicit-curry ≡ ML curry and
  Default `handle … in` ≡ ML `handle … do` at the canonical AST.
- The non-equivalent-bucket golden tests prove Default multi-param ≢ ML curry.
- `grep` finds no flavor inspection in `osprey-types` or `osprey-codegen`.
- Every existing Default `.osp` example still passes unchanged.
