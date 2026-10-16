# Contributing to Osprey

Want to help build a programming language? You're in the right place. Building a language is a several year long process. Right now, the aims are building community, getting the first example apps deployed and shaping the ergonomics of the language.

Even discussions are great at this point.

**Compiler changes**: always create an issue before starting work on the compiler. Explain the issue and propose a solution for the issue. Please get approval before moving into the implementation phase. All contributions are appreciated but compiler work requires strict review and planning.

Start with the [repository map](docs/contributing/repository-map.md) to find the component and build rules for your change.

## The Tech Stack

Osprey is built on battle-tested tools:
- **[Rust](https://www.rust-lang.org/)** - Compiler implementation (`crates/`)
- **[tree-sitter](https://tree-sitter.github.io/)** - Grammar parsing (`tree-sitter-osprey/`)
- **[LLVM](https://llvm.org/)** - Code generation (textual IR handed to clang)
- **C** - The runtime (fibers, HTTP/WebSocket, collections) in `compiler/runtime/`

## Getting Started

**Use VS Code Dev Containers** - strongly recommended. Open in VS Code and hit "Reopen in Container". Everything's already configured.

```bash
# Fork and clone
git clone https://github.com/Nimblesite/osprey.git
cd osprey

# Build and test (from the repo root)
make build
make test
```

## Repository coordination

Git worktrees are banned. There may be only **ONE active feature branch across ALL machines and clones combined**. Coordinate all work on the existing active branch; do not create a branch per machine, contributor, or task. A replacement branch may start only after the previous branch is merged or explicitly retired and all machines are coordinated. See [CLAUDE.md](CLAUDE.md) for the shared instructions.

## What to Work On

Osprey is in stabilization mode. Prioritize reliability of existing features, truthful diagnostics, and conformance to the documented contract. Active [plans](docs/plans/README.md) distinguish unfinished requirements from shipped behavior.

1. Harden existing language features and built-ins.
2. Expand acceptance and rejection coverage for edge cases.
3. Reconcile specs, code, tests, plans, and public status claims when behavior changes.
4. Improve existing examples and application boundaries.

## The AI Workflow

1. **Check the contract** — Read the relevant spec, implementation status, tests, and active plan.
2. **Understand the pattern** - Ask your AI: "How does pattern matching work in this codebase?"
3. **Implement incrementally** - Start with parsing, then AST, then codegen
4. **Test everything** - Add assertion suites to `tests/`
5. **Fix edge cases** - Let the AI help debug LLVM IR issues

Example prompts:
- "Add a new arithmetic operator to the tree-sitter grammar"
- "Implement string interpolation in the AST builder"
- "Generate LLVM IR for this pattern match"

## Code Guidelines

- Follow existing patterns - the clippy lint wall (`make lint`) enforces a lot
- Test new features thoroughly
- Keep changes focused
- Fix linter errors before submitting
- **Never hard-code a version.** All version fields stay at `0.0.0-dev` in source;
  releases stamp the real version from the git tag. See [docs/RELEASING.md](docs/RELEASING.md).

## CI & Releases

- **CI runs only on PRs to `main`** (`ci.yml` + `ci-windows.yml`). Merging to
  `main` triggers nothing.
- **Releases are tag-triggered**: push a tag `vX.Y.Z` and `release.yml` builds
  every platform, cuts the GitHub Release, updates the Homebrew tap + Scoop
  bucket, publishes the VS Code extension, and deploys the website. Full details
  in [docs/RELEASING.md](docs/RELEASING.md).

## Getting Help

- Open an issue for discussion
- Check existing issues and documentation
- Review similar implementations in the codebase
- The spec is the source of truth

## License

By contributing, you agree that your contributions will be licensed under the MIT License. 
