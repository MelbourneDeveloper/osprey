<!-- agent-pmo:74cf183 -->
# Engineering Specs

This directory holds **engineering specs** for repo behaviour, processes, and
cross-cutting concerns — anything an engineer needs as the source of truth for
*how this repository operates*.

For the **Osprey language specification** (the language this compiler
implements), see [`../../compiler/spec/`](../../compiler/spec/).

## Conventions

- Spec IDs are hierarchical descriptive slugs in the format `[GROUP-TOPIC]` or
  `[GROUP-TOPIC-DETAIL]`. **NEVER** numbered IDs (`[SPEC-001]`). See the
  project's [`CLAUDE.md`](../../CLAUDE.md) for the full convention.
- Code implementing a spec section MUST reference its ID in a comment, e.g.
  `// Implements [PARSER-EFFECTS-HANDLE]`.
- Sibling [`../plans/`](../plans/) holds *plans* (how-to docs with TODO
  checklists) rather than specs (what/why).
