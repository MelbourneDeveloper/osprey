# Hindley-Milner Type System Fixes

## Problem Summary

The Hindley-Milner type inference system was implemented with incorrect built-in function type definitions that don't match the actual LLVM implementations. This causes widespread type mismatches and compilation failures.

Essential, any function that CAN fail, must return a Result type. If it can't fail, for example length on string, it can return a raw type.

## Root Cause

The type definitions in `initializeBuiltInFunctions()` (in `type_inference.go`) declare what the type system THINKS each function returns, but the actual LLVM generation code (in various files) generates different return types.

## Built-in Function Return Type Analysis

After examining the actual LLVM implementations, here's what each function ACTUALLY returns:

### Functions that return plain values (no Result wrapper):
- **print(any)** → `i64` (puts return code extended to i64)
- **toString(any)** → `string` (i8* - direct conversion)
- **length(string)** → `int` never fails so doesn't return a result

### Functions that return Result types:
- **input()** → `Result<Int, Error>` (currently unimplemented - returns ErrUnsupportedCall)
- **contains(string, string)** → `Result<Bool, NoError>` (never fails but uses Result wrapper)
- **substring(string, int, int)** → `Result<String, NoError>` (can fail with bounds)
- **readFile(string)** → `Result<String, Error>` (can fail)
- **writeFile(string, string)** → `Result<Unit, Error>` (can fail)
- **spawnProcess(string, handler)** → `Result<Int, Error>` (can fail)
- **awaitProcess(ProcessHandle)** → `Result<Int, Error>` (can fail)
- **cleanupProcess(ProcessHandle)** → `Result<Unit, Error>` (can fail)
- HTTP functions → `Result<T, Error>` (can fail)

## Key Issues to Fix

### 1. Fix Type Definitions in `type_inference.go`

Update `initializeBuiltInFunctions()` to match actual implementations:

```go
// CORRECT definitions based on actual LLVM implementations:

// print(value: any) -> Int (NOT Result!)
ti.env.Set("print", &FunctionType{
    paramTypes: []Type{anyType},
    returnType: intType,  // Returns i64 directly
})

// toString(value: any) -> String (NOT Result!)
ti.env.Set("toString", &FunctionType{
    paramTypes: []Type{anyType},
    returnType: stringType,  // Returns string directly
})

// input() -> Result<Int, Error> (needs implementation)
ti.env.Set("input", &FunctionType{
    paramTypes: []Type{},
    returnType: &ConcreteType{name: "Result<Int, Error>"},
})


// contains, substring, file I/O, etc. already correct as Result types
```

### 2. Implement Missing Functions

**input()** function in `core_functions.go` currently returns `ErrUnsupportedCall`. Need to implement:

```go
func (g *LLVMGenerator) generateInputCall(callExpr *ast.CallExpression) (value.Value, error) {
    if len(callExpr.Arguments) != 0 {
        return nil, WrapInputWrongArgsWithPos(len(callExpr.Arguments), callExpr.Position)
    }
    
    // Declare scanf function
    scanf := g.ensureScanfDeclaration()
    
    // Create buffer for input
    bufferType := types.NewArray(256, types.I8)
    buffer := g.builder.NewAlloca(bufferType)
    
    // Call scanf to read integer
    formatStr := g.createGlobalString("%ld")
    result := g.builder.NewCall(scanf, formatStr, buffer)
    
    // Convert to Result<Int, Error>
    // Implementation details...
    
    return resultValue, nil
}
```

### 3. Fix Binary Operations for Result Types

The arithmetic operators expect `Int` but get `Result<Int, Error>` from functions like `length()`. Need to either:

A. **Add automatic Result unwrapping** in binary operations
B. **Change functions like length() to return plain Int** (breaking change)
C. **Require explicit pattern matching** on Result types

**Recommended: Option A** - Add automatic unwrapping for "safe" Result types:

```go
func (g *LLVMGenerator) generateBinaryExpression(binExpr *ast.BinaryExpression) (value.Value, error) {
    left, err := g.generateExpression(binExpr.Left)
    if err != nil {
        return nil, err
    }
    
    right, err := g.generateExpression(binExpr.Right)
    if err != nil {
        return nil, err
    }
    
    // Auto-unwrap Result types for arithmetic operations
    left = g.autoUnwrapResultIfNeeded(left, binExpr.Left)
    right = g.autoUnwrapResultIfNeeded(right, binExpr.Right)
    
    return g.generateBinaryOperationWithPos(binExpr.Operator, left, right, binExpr.Position)
}
```

### 4. Remove Debug Output

Remove all `DEBUG:` prints contaminating test output:

```go
// REMOVE THESE LINES:
fmt.Printf("DEBUG: Block has %d statements, expression: %t\n", len(e.Statements), e.Expression != nil)
fmt.Printf("DEBUG: Last statement type: %T\n", e.Statements[len(e.Statements)-1])
```

### 5. Fix LLVM Type Mismatches

Several functions have return type mismatches:

- Boolean functions returning `i64` instead of `i1`
- Functions expecting `i64` return type but generating structs
- Fix in `generateBooleanLiteral()` and comparison operators

### 6. Fix Struct Field Assignment Panic

The panic `store operands are not compatible: src=i8*; dst=i64*` indicates type mismatch in struct field assignment. Check field type mappings in `generateUnconstrainedRecordConstructor()`.

### 7. Update Error Message Formats

Several tests expect specific error message formats. Update error wrappers to match expected formats.

## Systematic Fix Strategy

1. **First**: Fix the type definitions in `initializeBuiltInFunctions()`
2. **Second**: Implement missing `input()` function
3. **Third**: Add Result type auto-unwrapping for arithmetic
4. **Fourth**: Remove debug prints and fix LLVM type issues
5. **Fifth**: Run tests and fix remaining format issues

## Files to Modify

- `compiler/internal/codegen/type_inference.go` - Fix type definitions
- `compiler/internal/codegen/core_functions.go` - Implement input(), fix signatures
- `compiler/internal/codegen/expression_generation.go` - Add Result unwrapping, remove debug
- `compiler/internal/codegen/errors.go` - Fix error message formats
- `compiler/internal/codegen/llvm.go` - Fix LLVM type issues

## Testing Strategy

After each fix:
```bash
cd compiler && make clean && make build && make test
```

Focus on these failing tests first:
- `TestBasicsExamples/blocks_block_statements_basic` (debug output)
- `TestBasicsExamples/comprehensive` (Result type arithmetic)
- `TestBasicsExamples/operators_equality_test` (LLVM type mismatch)

## Long-term Considerations

- **Design Decision**: Should "safe" functions like `length()` return Result types?
- **Type System**: Consider explicit Result handling vs automatic unwrapping
- **Error Handling**: Standardize which functions can fail vs which cannot
- **Performance**: Result types add overhead for functions that never fail

The current system is inconsistent and over-engineered. A simpler approach would be:
- Functions that CAN'T fail return plain types
- Functions that CAN fail return Result types
- Require explicit handling of Result types (no auto-unwrapping) 