---
name: spec-check
description: Audit spec/plan documents against the codebase. Ensures every spec section has implementing code, tests, and matching logic. Use when the user says "check specs", "spec audit", or "verify specs".
argument-hint: "[optional spec ID or filename filter]"
---
<!-- agent-pmo:74cf183 -->

# spec-check

> **Portable skill.** This skill adapts to the current repository. The agent MUST inspect the repo structure and use judgment to apply these instructions appropriately.

Audit spec/plan documents against the codebase. Ensures every spec section has implementing code, tests, and that the code logic matches the spec.

## Arguments

- `$ARGUMENTS` — optional spec name or ID to check (e.g., `AUTH-TOKEN-VERIFY` or `repo-standards`). If empty, check ALL specs. Spec IDs are descriptive slugs, NEVER numbered (see Step 1).

## Instructions

Follow these steps exactly. Be strict and pedantic. Stop on the first failure.

---

### Step 1: Validate spec ID structure

Before checking code/test references, verify that the specs themselves are well-formed.

1. Find all spec documents (see locations in Step 2).
2. Extract every section ID using the regex `\[([A-Z][A-Z0-9]*(-[A-Z0-9]+)+)\]`.
3. **Flag invalid IDs:**
   - Numbered IDs (`[SPEC-001]`, `[REQ-003]`, `[CI-004]`) — must be renamed to descriptive hierarchical slugs.
   - Single-word IDs (`[TIMEOUT]`) — must have a group prefix.
   - IDs with trailing numbers (`[FEAT-AUTH-01]`) — the number is meaningless, remove it.
4. **Check group clustering:** The first word of each ID is its group. All sections in the same group MUST appear together (adjacent) in the document. If they're scattered, flag it.
5. **Check for missing IDs:** Any heading that defines a requirement or behavior should have an ID. Flag headings in spec files that look like they define behavior but lack an ID.

If any ID violations are found, report them all and **STOP**:
```
SPEC ID VIOLATIONS:

- docs/specs/AUTH-SPEC.md line 12: [SPEC-001] → rename to descriptive ID (e.g., [AUTH-LOGIN])
- docs/specs/AUTH-SPEC.md line 30: [AUTH-TOKEN-VERIFY] and [AUTH-LOGIN] are not adjacent (scattered group)
- docs/specs/CI-SPEC.md line 5: "## Coverage thresholds" has no spec ID

Fix spec IDs first, then re-run spec-check.
```

If all IDs are valid, proceed to Step 2.

---

### Step 2: Find all spec/plan documents

Search for markdown files that contain spec sections with IDs. Look in these locations:

- `docs/*.md`
- `docs/**/*.md`
- `SPEC.md`
- `PLAN.md`
- `specs/*.md`
- `compiler/spec/*.md`

Use Glob to find candidate files, then use Grep to confirm they contain spec IDs.

**Spec ID patterns** — IDs appear in square brackets, typically at the start of a heading or section line. Match this regex pattern:

```
\[([A-Z][A-Z0-9]*(-[A-Z0-9]+)+)\]
```

Spec IDs are **hierarchical descriptive slugs, NEVER numbered.** The format is `[GROUP-TOPIC]` or `[GROUP-TOPIC-DETAIL]`. The first word is the **group** — all sections sharing the same group MUST appear together in the spec's table of contents. IDs are uppercase, hyphen-separated, unique across the repo, and MUST NOT contain sequential numbers.

For each file, extract every spec ID and its associated section title and the full section content (everything until the next heading of equal or higher level).

---

### Step 3: Filter specs

- If `$ARGUMENTS` is non-empty, filter the discovered specs:
  - If it matches a spec ID exactly (e.g., `AUTH-TOKEN-VERIFY`), check only that spec.
  - If it matches a partial name (e.g., `repo-standards`), check all specs in files whose path contains that string.
- If `$ARGUMENTS` is empty, process ALL discovered specs.

If filtering produces zero specs, report an error:
```
ERROR: No specs found matching "$ARGUMENTS". Discovered spec files: [list them]
```

---

### Step 4: Check each spec section

For EACH spec section that has an ID, perform checks A, B, and C below. **Stop on the first failure.**

#### Check A: Code references the spec ID

Search the entire codebase for the spec ID string, **excluding** these directories:
- `docs/`
- `node_modules/`
- `.git/`
- `*.md` files (markdown is docs, not code)

Use Grep with the literal spec ID (e.g., `[AUTH-TOKEN-VERIFY]`) to find references in code files.

**If NO code files reference the spec ID:**

```
SPEC VIOLATION: [AUTH-TOKEN-VERIFY] "Section Title" has no implementing code.

Every spec section must have at least one code file that references it via a comment
containing the spec ID (e.g., `// Implements [AUTH-TOKEN-VERIFY]`).

ACTION REQUIRED: Add a comment referencing [AUTH-TOKEN-VERIFY] in the file(s) that implement
this spec section, then re-run spec-check.
```

**STOP HERE. Do not continue to other checks.**

#### Check B: Tests reference the spec ID

Search test files for the spec ID. Test files are found in:
- `compiler/tests/`
- `compiler/**/*_test.go`
- `vscode-extension/**/*.test.ts`

**If NO test files reference the spec ID:**

```
SPEC VIOLATION: [AUTH-TOKEN-VERIFY] "Section Title" has no tests.

Every spec section must have corresponding tests that reference the spec ID.

ACTION REQUIRED: Add tests for [AUTH-TOKEN-VERIFY] with a comment or test name containing
the spec ID, then re-run spec-check.
```

**STOP HERE. Do not continue to other checks.**

#### Check C: Code logic matches the spec

1. **Read the spec section content carefully.** Understand exactly what behavior, logic, ordering, conditions, and constraints the spec describes.
2. **Read the implementing code.** Use the references found in Check A to locate the implementing files.
3. **Compare spec vs. code.** Be SENSITIVE and PEDANTIC. Check for ordering violations, missing conditions, wrong logic, missing steps, wrong defaults.
4. **If the code deviates from the spec**, report a detailed error with quotes from both spec and code. **STOP.**
5. **If the code matches the spec**, this check passes. Move to the next spec.

---

### Step 5: Report results

#### On failure (any check fails):

Output ONLY the first violation found. End with:
```
spec-check FAILED. Fix the violation above and re-run.
```

#### On success (all specs pass):

Output a summary table:

```
spec-check PASSED. All specs verified.

| Spec ID        | Title                    | Code References | Test References | Logic Match |
|----------------|--------------------------|-----------------|-----------------|-------------|
| [AUTH-TOKEN]   | Authentication flow      | src/auth.go     | tests/auth_test.go | PASS     |

Checked N spec sections across M files. All have implementing code, tests, and matching logic.
```

## Key principles

- **Fail fast.** Stop on the first violation. One fix at a time.
- **Be pedantic.** If the spec says it, the code must do it. No "close enough".
- **Quote everything.** Always quote the spec text and the code in error messages.
- **Be actionable.** Every error must tell the developer what file to change and what to do.
- **No numbered IDs.** Spec IDs are hierarchical descriptive slugs, NEVER sequential numbers.
