# Remaining Compiler Issues Fix Plan

## Failing Test Suites Analysis

### 1. `compilation_failure_test.go` (11 remaining failures)
**Status**: 74% complete (31/42 passing)
**Specific Failing Tests**:
- `http_get_wrong_args.ospo` (HTTP validation issue)
- `http_streaming_example.ospo` (HTTP streaming scope issue)
- `length_wrong_args.ospo` (built-in function validation)
- `match_duplicate_arm.ospo` (pattern matching validation)
- `match_invalid_pattern.ospo` (pattern matching validation)
- `match_mixed_patterns.ospo` (pattern matching validation)
- `match_not_exhaustive.ospo` (exhaustiveness checking)
- `match_wildcard_not_last.ospo` (pattern matching validation)
- `named_args_violation.ospo` (named arguments validation)
- `print_convert_error.ospo` (type conversion validation)
- `unsupported_unary_op.ospo` (operator validation)

### 2. `examples_test.go` (19.2s runtime - failing)
**Build Issues**:
- `undefined: checkLLVMTools` (missing utility function)
- `undefined: captureJITOutput` (missing utility function)
**Root Cause**: Test infrastructure functions are not properly shared between test files
**Solution**: Extract common test utilities into shared package

### 3. `fiber_examples_test.go` (1.8s runtime - failing)
**Build Issues**:
- `undefined: checkLLVMTools` (missing utility function)
- `undefined: runTestExamplesRecursive` (missing utility function)
**Root Cause**: Test infrastructure functions are not properly shared between test files
**Solution**: Extract common test utilities into shared package

### 4. `fiber_integration_test.go` (1.x runtime - failing)
**Likely Issues**:
- Fiber runtime integration problems
- C runtime binding issues
- Memory management in concurrent contexts

### 5. `full_integration_test.go` (1.4s runtime - failing)
**Likely Issues**:
- End-to-end compilation pipeline problems
- Multiple language features interacting incorrectly
- Integration with external systems

## Phase 3: Test Infrastructure (Priority: CRITICAL)

### 3.1 Fix Test Utility Functions
**Problem**: Multiple test files fail to compile due to missing utility functions
**Files Affected**:
- `examples_test.go` - missing `checkLLVMTools`, `captureJITOutput`
- `fiber_examples_test.go` - missing `checkLLVMTools`, `runTestExamplesRecursive`
- Other test files may have similar issues

**Solution**:
1. **Create shared test utilities package**: `tests/integration/testutils/`
2. **Extract common functions**:
   - `checkLLVMTools` (from `full_integration_test.go`)
   - `captureJITOutput` (from `full_integration_test.go`)
   - `runTestExamplesRecursive` (from `examples_test.go`)
3. **Update imports**: Modify all test files to import from testutils package
4. **Verify all tests compile**: Ensure no build failures before proceeding

**Priority**: This must be fixed first as it blocks running other test suites

## Phase 4: Advanced Language Features (Priority: HIGH)

### 4.1 Match Exhaustiveness Checking
**File**: `match_not_exhaustive.ospo` should fail but passes
**Implementation Plan**:
1. **Analyze**: `compiler/internal/codegen/match_validation.go`
2. **Add**: Exhaustiveness analysis for enum/variant types
3. **Implement**: Pattern coverage validation
4. **Test**: Ensure all match expressions cover all cases

**Key Files to Modify**:
- `internal/codegen/match_validation.go`
- `internal/codegen/type_inference.go`
- Add exhaustiveness checking logic

### 4.2 Named Arguments Validation
**File**: `named_args_violation.ospo` should fail but passes
**Implementation Plan**:
1. **Analyze**: Current function call validation
2. **Add**: Named argument position/type checking
3. **Implement**: Parameter name validation
4. **Test**: Ensure incorrect named args are rejected

**Key Files to Modify**:
- `internal/codegen/function_calls.go`
- `internal/codegen/type_inference.go`
- `internal/ast/builder_arguments.go`

### 4.3 Pattern Matching Validation
**Files**: Multiple pattern matching tests failing
**Implementation Plan**:
1. **Fix**: `match_duplicate_arm.ospo` - detect duplicate patterns
2. **Fix**: `match_invalid_pattern.ospo` - validate pattern syntax
3. **Fix**: `match_mixed_patterns.ospo` - ensure consistent pattern types
4. **Fix**: `match_wildcard_not_last.ospo` - enforce wildcard position rules

**Key Files to Modify**:
- `internal/codegen/match_validation.go`
- Add comprehensive pattern validation logic

## Phase 5: HTTP and Networking (Priority: HIGH)

### 5.1 HTTP Streaming Scope Issue
**File**: `http_streaming_example.ospo` shows wrong error
**Implementation Plan**:
1. **Debug**: Scope resolution in HTTP contexts
2. **Fix**: Variable binding in streaming callbacks
3. **Test**: HTTP streaming examples compile correctly

**Key Files to Investigate**:
- `internal/codegen/http_generation.go`
- `internal/codegen/scope_management.go`
- `examples/tested/http/` directory

### 5.2 HTTP Argument Validation
**File**: `http_get_wrong_args.ospo` should fail but passes
**Implementation Plan**:
1. **Analyze**: HTTP function argument validation
2. **Fix**: Argument count and type checking
3. **Test**: Ensure HTTP functions reject invalid arguments

## Phase 6: Examples and Basic Features (Priority: MEDIUM)

### 6.1 Fix `examples_test.go` Failures
**Strategy**:
1. **First**: Fix test infrastructure issues (Phase 3)
2. **Run**: Individual example tests to identify specific failures
3. **Categorize**: Issues by language feature (types, functions, etc.)
4. **Fix**: Systematic repair of basic language features
5. **Verify**: All examples compile and run correctly

### 6.2 Built-in Function Validation
**Files**: `length_wrong_args.ospo`, `print_convert_error.ospo`
**Implementation Plan**:
1. **Fix**: Built-in function argument validation
2. **Add**: Type conversion error checking
3. **Test**: Ensure built-in functions reject invalid usage

### 6.3 Operator Validation
**File**: `unsupported_unary_op.ospo` should fail but passes
**Implementation Plan**:
1. **Analyze**: Unary operator validation
2. **Add**: Unsupported operator detection
3. **Test**: Ensure invalid operators are rejected

## Phase 7: Fiber and Concurrency (Priority: MEDIUM)

### 7.1 Fix `fiber_examples_test.go` Failures
**Strategy**:
1. **First**: Fix test infrastructure issues (Phase 3)
2. **Analyze**: Fiber system implementation
3. **Debug**: Algebraic effects compilation
4. **Fix**: Async/await pattern generation
5. **Test**: Concurrent program execution

**Key Areas**:
- Fiber creation and management
- Effect handlers implementation
- Async operation compilation
- Memory safety in concurrent contexts

### 7.2 Fix `fiber_integration_test.go` Failures
**Strategy**:
1. **Debug**: C runtime integration
2. **Fix**: Memory management issues
3. **Test**: Fiber runtime stability
4. **Verify**: Cross-language bindings

## Phase 8: Full Integration (Priority: LOW)

### 8.1 Fix `full_integration_test.go` Failures
**Strategy**:
1. **Identify**: End-to-end pipeline issues
2. **Fix**: Integration between all systems
3. **Test**: Complete compilation workflow
4. **Verify**: Production-ready compiler

## Implementation Strategy

### Week 1: Phase 3 (Test Infrastructure) + Phase 4 (Advanced Language Features)
- **Days 1-2**: Fix test infrastructure (Phase 3)
- **Days 3-4**: Match exhaustiveness and named arguments (Phase 4)
- **Day 5**: Pattern matching validation (Phase 4)

### Week 2: Phase 5 (HTTP and Networking)
- **Days 1-3**: HTTP streaming scope issue
- **Days 4-5**: HTTP argument validation

### Week 3: Phase 6 (Examples and Basic Features)
- **Days 1-3**: Fix examples_test.go failures
- **Days 4-5**: Built-in function and operator validation

### Week 4: Phase 7 (Fiber and Concurrency)
- **Days 1-3**: Fiber examples fixes
- **Days 4-5**: Fiber integration fixes

### Week 5: Phase 8 (Full Integration)
- **Days 1-3**: Full integration testing
- **Days 4-5**: Final verification and cleanup

## Success Metrics

### Phase 3 Complete (Test Infrastructure)
- [ ] All test files compile successfully
- [ ] Common test utilities properly shared
- [ ] No more undefined function errors

### Phase 4 Complete (Advanced Language Features)
- [ ] `match_not_exhaustive.ospo` fails compilation (correctly)
- [ ] `named_args_violation.ospo` fails compilation (correctly)
- [ ] All pattern matching validation tests pass
- [ ] 95%+ compilation failure tests passing (40+/42)

### Phase 5 Complete (HTTP and Networking)
- [ ] `http_streaming_example.ospo` fails with correct error
- [ ] `http_get_wrong_args.ospo` fails compilation (correctly)
- [ ] All HTTP tests pass
- [ ] Network-related examples work

### Phase 6 Complete (Examples and Basic Features)
- [ ] `examples_test.go` passes completely
- [ ] `length_wrong_args.ospo` fails compilation (correctly)
- [ ] `print_convert_error.ospo` fails compilation (correctly)
- [ ] `unsupported_unary_op.ospo` fails compilation (correctly)
- [ ] 100% compilation failure tests passing (42/42)

### Phase 7 Complete (Fiber and Concurrency)
- [ ] `fiber_examples_test.go` passes
- [ ] `fiber_integration_test.go` passes
- [ ] Concurrency system stable

### Phase 8 Complete (Full Integration)
- [ ] `full_integration_test.go` passes
- [ ] All test suites green
- [ ] Production-ready compiler

## Risk Assessment

### High Risk
- **Fiber system complexity**: May require significant runtime changes
- **HTTP integration**: Network programming edge cases
- **Memory management**: Concurrent systems are complex

### Medium Risk
- **Match exhaustiveness**: Complex pattern analysis required
- **Named arguments**: Function call system changes
- **Example compilation**: Wide variety of potential issues

### Low Risk
- **Error message consistency**: Already proven approach
- **Documentation**: Straightforward cleanup tasks

## Next Steps

1. **Immediate**: Focus on Phase 3 (test infrastructure) - fix shared utility functions
2. **Then**: Phase 4 (advanced language features) - match exhaustiveness and named args
3. **Run diagnostics**: After Phase 3, run `go test -v ./tests/integration/` to get detailed failure info
4. **Create targeted fixes**: One issue at a time with proper testing
5. **Maintain momentum**: Keep build passing throughout the process

## Summary

This plan addresses the remaining compiler issues in a structured, priority-based approach:

1. **Phase 3 (CRITICAL)**: Fix test infrastructure to enable running all test suites
2. **Phase 4 (HIGH)**: Fix advanced language features (11 specific failing tests)
3. **Phase 5 (HIGH)**: Fix HTTP and networking issues
4. **Phase 6 (MEDIUM)**: Fix examples and basic language features
5. **Phase 7 (MEDIUM)**: Fix fiber and concurrency system
6. **Phase 8 (LOW)**: Fix full integration pipeline

The current 74% success rate (31/42) on compilation failure tests provides a solid foundation. With systematic execution of this plan, we can achieve 100% test suite success and a production-ready compiler.

## Notes
- Follow project rules: test after every change, keep files under 500 LOC
- Maintain backward compatibility
- Document all changes
- Use constants instead of literals
- Follow Go best practices 