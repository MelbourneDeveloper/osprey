# Fiber Non-Determinism Bug

## Problem Description

The Osprey compiler's fiber runtime exhibits non-deterministic execution order when multiple fibers are spawned and awaited sequentially. This causes test failures because the output order varies between runs.

## Root Cause

The issue is in `/workspace/compiler/runtime/fiber_runtime.c`:

1. **Default Mode**: `deterministic_mode` is `false` by default (line 40)
2. **pthread Threads**: Fibers use `pthread_create()` for concurrent execution (line 112)
3. **OS Scheduling**: Thread execution order depends on OS scheduler, which is non-deterministic
4. **Race Conditions**: Multiple fibers can complete in any order, even when awaited sequentially

## Evidence

### Failing Tests
- `osprey_mega_showcase.osp` - fails ~22% of the time (11/50 runs)
- `fiber_exact_replica_test.osp` - fails when cache is cleared
- `fiber_nondeterminism_fail_test.osp` - fails when cache is cleared

### Example Failure
```
Expected: task-alpha-1, task-alpha-2, task-alpha-3
Got:      task-alpha-2, task-alpha-1, task-alpha-3
```

## Technical Details

### Fiber Spawn Process
```c
// fiber_runtime.c:77
int64_t fiber_spawn(int64_t (*fn)(void)) {
    // ...
    if (!deterministic_mode) {
        // Creates pthread thread - NON-DETERMINISTIC
        pthread_create(&fiber->thread, NULL, fiber_thread_func, fiber);
    } else {
        // Queues for sequential execution - DETERMINISTIC
        execution_queue[queue_size++] = id;
    }
}
```

### Fiber Await Process
```c  
// fiber_runtime.c:134
int64_t fiber_await(int64_t fiber_id) {
    if (is_deterministic) {
        // Execute fibers in queue order - DETERMINISTIC
        for (int64_t i = 0; i < queue_size; i++) {
            // Execute in spawn order
        }
    } else {
        // Wait for thread completion - NON-DETERMINISTIC
        pthread_cond_wait(&fiber->cond, &fiber->mutex);
    }
}
```

## Solutions

### Option 1: Enable Deterministic Mode by Default
```c
// Change line 40 in fiber_runtime.c
static bool deterministic_mode = true; // Change from false to true
```

**Pros:**
- Fixes all test failures immediately
- Predictable behavior for users
- No API changes needed

**Cons:**
- Loses true parallelism (fibers execute sequentially)
- May impact performance for concurrent workloads

### Option 2: Add Compiler Flag for Deterministic Mode
Add a command-line flag to control fiber execution mode:

```go
// In CLI code
if *deterministicFlag {
    C.fiber_set_deterministic_mode(C.bool(true))
}
```

**Pros:**
- Best of both worlds - deterministic when needed, parallel when desired
- Backward compatible
- Users can choose behavior

**Cons:**
- Requires CLI changes
- Tests need to specify the flag

### Option 3: Fix pthread Ordering
Implement proper ordering in pthread mode using synchronization primitives:

```c
// Add ordering constraints in fiber_thread_func
// Use barriers, semaphores, or condition variables
// to ensure spawned fibers execute in spawn order
```

**Pros:**
- True parallelism with deterministic ordering
- Best performance

**Cons:**
- Complex implementation
- May introduce deadlocks
- Significant runtime changes needed

### Option 4: Test-Specific Solution
Modify test framework to handle non-deterministic output:

```go
// In examples_test.go
func normalizeOutput(output string) string {
    // Sort fiber execution logs by task ID
    // Return normalized output for comparison
}
```

**Pros:**
- No runtime changes
- Tests become more robust

**Cons:**
- Doesn't fix underlying issue
- Users still experience non-determinism

## Recommended Solution

**Implement Option 2: Add Compiler Flag**

1. Add `--deterministic-fibers` flag to CLI
2. Set `fiber_set_deterministic_mode(true)` when flag is present
3. Update test runner to use the flag for all tests
4. Document the behavior for users

This provides:
- ✅ Fixed tests (deterministic mode for testing)
- ✅ Performance (parallel mode for production)
- ✅ User choice (explicit control over behavior)
- ✅ Backward compatibility

## Implementation Steps

1. **CLI Changes**: Add deterministic fiber flag
2. **Runtime Integration**: Call `fiber_set_deterministic_mode()` based on flag
3. **Test Updates**: Use deterministic flag in test runner
4. **Documentation**: Update docs explaining the two modes

## Files to Modify

- `internal/cli/cli.go` - Add flag parsing
- `internal/codegen/compilation.go` - Call runtime function
- `tests/integration/examples_test.go` - Use flag in tests
- `Makefile` - Update test commands
- `CLAUDE.md` - Document the flag

## Current Workaround

For immediate testing, manually enable deterministic mode:
```bash
# Before running tests, patch the runtime:
sed -i 's/static bool deterministic_mode = false;/static bool deterministic_mode = true;/' runtime/fiber_runtime.c
make build
make test
```

This bug affects the reliability of the test suite and creates unpredictable behavior for users working with fiber concurrency.