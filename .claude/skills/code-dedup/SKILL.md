---
name: code-dedup
description: Searches for duplicate code, duplicate tests, and dead code, then safely merges or removes them. Use when the user says "deduplicate", "find duplicates", "remove dead code", "DRY up", or "code dedup". Requires test coverage — refuses to touch untested code.
---
<!-- agent-pmo:74cf183 -->

# Code Dedup

Carefully search for duplicate code, duplicate tests, and dead code across the repo. Merge duplicates and delete dead code — but only when test coverage proves the change is safe.

## Prerequisites — hard gate

Before touching ANY code, verify these conditions. If any fail, stop and report why.

1. Run `make test` — all tests must pass. If tests fail, stop. Do not dedup a broken codebase.
2. Run `make test` — tests are fail-fast AND enforce the coverage threshold from `coverage-thresholds.json`. If anything fails, stop and fix it before deduping.
3. Verify the project uses **static typing**:
   - Go (compiler/): typed by default — proceed
   - TypeScript (vscode-extension/, webcompiler/, website/): check `tsconfig.json` has `"strict": true` — proceed if yes
   - C (runtime/): typed by default — proceed

## Steps

Copy this checklist and track progress:

```
Dedup Progress:
- [ ] Step 1: Prerequisites passed (tests green, coverage met, typed)
- [ ] Step 2: Dead code scan complete
- [ ] Step 3: Duplicate code scan complete
- [ ] Step 4: Duplicate test scan complete
- [ ] Step 5: Changes applied
- [ ] Step 6: Verification passed (tests green, coverage stable)
```

### Step 1 — Inventory test coverage

Before deciding what to touch, understand what is tested.

1. Run `make test` to confirm green baseline. `make test` is fail-fast AND enforces the coverage threshold from `coverage-thresholds.json` (REPO-STANDARDS-SPEC [TEST-RULES], [COVERAGE-THRESHOLDS-JSON]). It exits non-zero on any test failure OR coverage shortfall.
2. Note the current coverage percentage — this is the floor. It must not drop.
3. Identify which files/modules have coverage and which do not. Only files WITH coverage are candidates for dedup.

### Step 2 — Scan for dead code

Search for code that is never called, never imported, never referenced.

1. Look for unused exports, unused functions, unused variables
2. Use language-appropriate tools:
   - Go: the compiler already catches unused imports/variables; `golangci-lint run` reports unused code
   - TypeScript: check for `noUnusedLocals`/`noUnusedParameters` in tsconfig, look for unexported functions with zero references
   - C: compiler warnings for unused functions/variables (already caught by `-Wall -Wextra`)
3. For each candidate: **grep the entire codebase** for references (including tests, scripts, configs). Only mark as dead if truly zero references.
4. List all dead code found with file paths and line numbers. Do NOT delete yet.

### Step 3 — Scan for duplicate code

Search for code blocks that do the same thing in multiple places.

1. Look for functions/methods with identical or near-identical logic
2. Look for copy-pasted blocks (same structure, maybe different variable names)
3. Look for multiple implementations of the same algorithm or pattern
4. Check across module boundaries — duplicates often hide in different packages
5. For each duplicate pair: note both locations, what they do, and how they differ (if at all)
6. List all duplicates found. Do NOT merge yet.

### Step 4 — Scan for duplicate tests

Search for tests that verify the same behavior.

1. Look for test functions with identical assertions against the same code paths
2. Look for test fixtures/helpers that are duplicated across test files
3. Look for integration tests that fully cover what a unit test also covers
4. List all duplicate tests found. Do NOT delete yet.

### Step 5 — Apply changes (one at a time)

For each change, follow this cycle: **change → test → verify coverage → continue or revert**.

#### 5a. Remove dead code
- Delete dead code identified in Step 2
- After each deletion: run `make test` (fail-fast + coverage + threshold all in one)
- If `make test` exits non-zero (test failure OR coverage drop): **revert immediately** and investigate

#### 5b. Merge duplicate code
- For each duplicate pair: extract the shared logic into a single function/module
- Update all call sites to use the shared version
- After each merge: run `make test`
- If tests fail: **revert immediately**

#### 5c. Remove duplicate tests
- Delete the redundant test (keep the more thorough one)
- After each deletion: run `make test`
- If coverage drops below threshold, `make test` exits non-zero — **revert immediately**

### Step 6 — Final verification

1. Run `make lint` — all linters must pass
2. Run `make test` — tests must pass AND coverage must remain ≥ the baseline from Step 1
3. Report: what was removed, what was merged, final coverage vs baseline

## Rules

- **No test coverage = do not touch.** If a file has no tests covering it, leave it alone entirely.
- **Coverage must not drop.** The coverage floor from Step 1 is sacred.
- **One change at a time.** Make one dedup change, run tests, verify coverage.
- **When in doubt, leave it.** If two code blocks look similar but you're not 100% sure they're functionally identical, leave both.
- **Preserve public API surface.** Do not change function signatures or exported names.
- **Three similar lines is fine.** Only dedup when the shared logic is substantial (>10 lines) or when there are 3+ copies.
