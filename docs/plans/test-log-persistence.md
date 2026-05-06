# Plan: Test Log Persistence

GitHub Issue: https://github.com/MelbourneDeveloper/osprey/issues/56

## Motivation

During a `/ci-prep` session, context was compacted mid-run. The next session had no test logs available and was forced to re-run the full test suite (several minutes) just to see what was failing. This is wasteful and slows down iterative fixing.

## Steps

### 1. Tee `make test` output to `compiler/test.log`

In `compiler/Makefile`, change the `test` target so stdout+stderr are written to `test.log` while still printing to the terminal:

```makefile
test:
	go test ./... 2>&1 | tee test.log; exit $${PIPESTATUS[0]}
```

### 2. Gitignore the log file

Add to `compiler/.gitignore` (or root `.gitignore`):

```
compiler/test.log
```

### 3. Upload log as CI artifact on failure

In `.github/workflows/ci.yml`, add an upload step after the test step:

```yaml
- name: Upload test log on failure
  if: failure()
  uses: actions/upload-artifact@v4
  with:
    name: test-log
    path: compiler/test.log
```

## Acceptance Criteria

- `make test` always writes output to `compiler/test.log`
- A fresh agent session can `Read compiler/test.log` without re-running tests
- CI uploads `test.log` as a downloadable artifact on failure
- `test.log` is not committed to git
