# Plan 0014 - Modules, Namespaces, and Multi-File Apps

## Summary

Implement [spec 0025 - Modules and Namespaces](../specs/0025-ModulesAndNamespaces.md):
.NET-style open namespaces for path-independent names, ML-style modules and
signatures for abstraction, and explicit state modules for centralising mutable
state.

The existing compiler already has `Stmt::Import`, `Stmt::Module`, Default-flavor
syntax for `import`/`module`, child-scope module checking, and project-adjacent
LSP work. The missing pieces are the project namespace graph, import resolution,
exports, signatures, qualified names, state-ownership checks, and codegen/LSP
support for fully qualified symbols.

## Current State

- `osprey_ast::Stmt::Import { module: Vec<String> }` and
  `Stmt::Module { name, body }` exist.
- The Default lowerer parses `import_statement` and `module_declaration`.
- The type checker recurses into module bodies with a child scope, but module
  declarations do not export/import usable symbols.
- The LSP outline recurses into modules but flattens names.
- `docs/specs/0011-LightweightFibersAndConcurrency.md` has an older
  fiber-isolated module sketch. Spec 0025 supersedes it.
- Cross-file LSP is already planned in
  [plan 0009](0009-lsp-context-and-cross-file.md), but it needs the module graph
  from this plan.

## Non-Goals

- No package manager yet.
- No recursive modules in the first implementation.
- No higher-order parameterised modules before basic signatures work.
- No wildcard imports in library code by default.
- No implicit path-to-namespace mapping.

## Phase 0 - Spec And Parser Contract

TODO:

- [x] Add spec 0025.
- [x] Update `docs/specs/0003-Syntax.md` so `import` and `module` point to spec
      0025 for semantics.
- [x] Update `docs/specs/0011-LightweightFibersAndConcurrency.md` to mark the
      old fiber-isolated module paragraph as superseded by spec 0025.
- [ ] Decide exact surface grammar for ML-flavor `namespace`, `module`,
      `signature`, `export`, and `state module`.
- [ ] Reserve new keywords in both flavors: `namespace`, `signature`, `export`,
      `opaque`, `state`, `as`.

## Phase 1 - AST And Project Model

TODO:

- [ ] Add `QualifiedName(Vec<String>)` to `osprey-ast`.
- [ ] Replace string-only `Stmt::Module { name }` with qualified module names and
      visibility/export metadata.
- [ ] Add `Stmt::Namespace { name, body }`.
- [ ] Add `Stmt::Signature { name, items }`.
- [ ] Add import forms: namespace/module import, member import list, alias, and
      wildcard.
- [ ] Add `Visibility::{Private, Exported}` on module items.
- [ ] Add opaque/manifest type export metadata.
- [ ] Add `ModuleKind::{Plain, State}` and `state_boundary` metadata.
- [ ] Preserve source spans on every new declaration for diagnostics and LSP.

## Phase 2 - Frontend Lowering

TODO:

- [ ] Default flavor: parse block-scoped and file-scoped `namespace`.
- [ ] Default flavor: parse `module A.B { ... }`, `state module`, `signature`,
      `export`, `opaque type`, import aliases, member lists, and wildcards.
- [ ] Default flavor: parse qualified names in expressions and types.
- [ ] ML flavor: add the same constructs in layout form and lower to the same
      canonical AST.
- [ ] Interpolation re-entry must parse qualified names using the current file's
      flavor.
- [ ] Add parser tests for path-independent namespace declarations, duplicate
      namespace blocks, exports, signatures, aliases, and qualified calls.

## Phase 3 - Project Loader And Namespace Graph

TODO:

- [ ] Add `Project` / `ProjectGraph` in a compiler-facing crate, shared by CLI
      and LSP.
- [ ] Read `osprey.toml` with `source_roots`, `root_namespace`, and module
      policy. Single-file mode remains unchanged.
- [ ] Scan `.osp` and `.ospml` files under source roots.
- [ ] Resolve flavor per file using the existing precedence rules.
- [ ] Parse every file independently, then merge namespace declarations by
      qualified name.
- [ ] Build an import table that maps imports to namespaces/modules, never file
      paths.
- [ ] Detect duplicate exported declarations in one namespace/module.
- [ ] Detect ambiguous imports and emit actionable diagnostics.
- [ ] Enforce one project entry point: designated entry file or `fn main()`.
- [ ] Reject executable top-level statements in non-entry project files.

## Phase 4 - Name Resolution And Type Checking

TODO:

- [ ] Add a resolver pass before type checking: local scope, module private
      scope, imported aliases, imported members, namespace qualified lookup,
      builtins.
- [ ] Store resolved symbol IDs on declarations/references or in a side table.
- [ ] Make type inference consume resolved symbols instead of raw strings.
- [ ] Enforce module privacy: private names are visible only inside their module.
- [ ] Check explicit exports and signature ascriptions.
- [ ] Implement opaque exported types: representation available inside the owning
      module, abstract outside.
- [ ] Check effect declarations and operations through signatures.
- [ ] Allow separate type checking of importers against signatures.
- [ ] Add tests for cross-file values, functions, types, effects, and flavor
      pairs.

## Phase 5 - State Ownership Rules

TODO:

- [ ] Reject namespace-level `mut`.
- [ ] Reject exported `mut` cells.
- [ ] Reject state-cell escape through exported pointers/references once pointer
      escape analysis exists; until then, reject direct export of `Ptr` derived
      from a state cell.
- [ ] Require every `state module` to expose a declared access path: handler,
      effect, or function API.
- [ ] Enforce one unannotated `state module` per namespace.
- [ ] Add `@state_boundary("reason")` parsing and diagnostics.
- [ ] Prefer handler-owned state for state modules that expose algebraic effects.
- [ ] Add LSP warnings and docs metadata that list all project state boundaries.
- [ ] Add compile-fail tests for scattered state, wildcard import of state
      modules, state cycles, and exported mutable cells.

## Phase 6 - Codegen And Runtime

TODO:

- [ ] Mangle fully qualified names deterministically.
- [ ] Update function, extern, effect-operation, handler, type-constructor, and
      generated lambda names to use qualified symbol IDs.
- [ ] Ensure imports do not emit runtime initialization.
- [ ] Lower pure namespace/module constants without hidden order dependence.
- [ ] Lower state module instances through explicit handler/instance
      construction.
- [ ] Reject cyclic state initialization before codegen.
- [ ] Preserve source-level qualified names in debug info and stack traces.
- [ ] Add IR equivalence tests for cross-flavor modules with identical canonical
      project graphs.

## Phase 7 - CLI, LSP, Formatter, Docs

TODO:

- [ ] CLI: `osprey build` / project mode reads `osprey.toml`; existing single-file
      commands keep working.
- [ ] CLI: diagnostics show qualified names and import candidates.
- [ ] LSP: maintain an incremental project graph across open files and source
      roots.
- [ ] LSP: go-to-definition, references, hover, completion, and document symbols
      understand namespaces/modules/imports.
- [ ] LSP: show state-boundary warnings and quick fixes for aliases/qualified
      names.
- [ ] Formatter: preserve file-scoped namespace and format module/signature
      blocks in both flavors.
- [ ] Docs generator: create namespace/module reference pages from exported
      signatures.

## Phase 8 - Tests And Examples

TODO:

- [ ] Add `examples/tested/modules/` with multi-file projects.
- [ ] Include path-independent namespaces where file paths intentionally do not
      match namespace names.
- [ ] Include Default imports ML and ML imports Default.
- [ ] Include explicit import lists, aliases, ambiguous import failures, and
      wildcard policy failures.
- [ ] Include signatures with opaque and manifest types.
- [ ] Include a state module that exposes an effect handler and a pure fake for
      tests.
- [ ] Add compile-fail examples for namespace-level `mut`, exported `mut`,
      duplicate exports, private-name leakage, and state cycles.
- [ ] Add LSP integration tests for cross-file completion/hover/definition.
- [ ] `make ci` green.

## Rollout Order

1. AST + parser support with no project mode, covered by unit tests.
2. Project graph and resolver in check-only mode.
3. Type checker on qualified/resolved names.
4. Exports/signatures/opaque types.
5. State rules.
6. Codegen for multi-file project mode.
7. LSP and formatter polish.
8. Parameterised modules after the basic module system is stable.

## Risks

- Resolver churn will touch type checking and codegen. Keep a raw-name fallback
  only temporarily and remove it before project mode is declared complete.
- Opaque types need careful interaction with existing union/record constructors.
- State modules overlap with existing handler-owned state. Treat state modules as
  a disciplined way to define handlers and access paths, not as process-global
  mutable singletons.
- The old top-level-script model is useful for examples. Preserve it in
  single-file mode while project mode enforces one entry point.
