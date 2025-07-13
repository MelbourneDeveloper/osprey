# Compiler Fixes Needed

## 1. Field Access in Pattern Matching (CRITICAL)
**Problem**: `field_access_comprehensive.osp` fails - pattern variables not bound
**Fix**: In `compiler/internal/codegen/llvm.go`, function `extractDiscriminatedUnionFields`:
- Change from `variables[field.Name] = fieldValue` 
- To `variables[pattern.Fields[i]] = fieldValue` (positional mapping)
- Test: `osprey examples/tested/basics/field_access_comprehensive.osp --run`

## 2. Test Infrastructure (CRITICAL)
**Problem**: `examples_test.go`, `fiber_examples_test.go` fail to compile
**Fix**: Missing utility functions `checkLLVMTools`, `captureJITOutput`, `runTestExamplesRecursive`, which you can get from the `main` branch.
- Or create shared testutils package

## 3. Compilation Failure Tests (11 remaining)
**Problem**: These should fail but pass:
- `http_get_wrong_args.ospo` - add HTTP argument validation
- `length_wrong_args.ospo` - add builtin function validation  
- `print_convert_error.ospo` - add type conversion checking
- `unsupported_unary_op.ospo` - add operator validation
- `named_args_violation.ospo` - add named argument validation
- `match_not_exhaustive.ospo` - add exhaustiveness checking
- `match_duplicate_arm.ospo` - add duplicate pattern detection
- `match_invalid_pattern.ospo` - add pattern syntax validation
- `match_mixed_patterns.ospo` - add pattern type consistency
- `match_wildcard_not_last.ospo` - add wildcard position validation
- `http_streaming_example.ospo` - fix variable scoping in HTTP callbacks

**Fix**: Add validation logic to respective codegen files to make these properly fail compilation

## 4. Integration Tests
**Problem**: `fiber_integration_test.go`, `full_integration_test.go` failing
**Fix**: Debug C runtime integration and memory management issues

## Priority Order
1. Field access (blocks examples)
2. Test infrastructure (blocks running tests)
3. Compilation failures (blocks validation)
4. Integration tests (final verification)

## Testing
After each fix: `make clean && make build && make test` 