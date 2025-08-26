# Hindley-Milner Type System Fix Plan

## Current Issues

The Osprey compiler's Hindley-Milner type inference implementation has fundamental inconsistencies causing test failures:

1. **Type Unification Failures**: `parameter 0 unification failed: int vs string`
2. **Record Field Access Issues**: Field ordering dependencies instead of name-only access
3. **Pattern Matching Errors**: String patterns generating `br i1 false` instead of proper comparisons
4. **Validation Function Type Mismatches**: Record types inconsistent between construction and validation

## Root Cause Analysis

Based on Hindley-Milner theory from the authoritative sources, the issues stem from:

### 1. Incomplete Type Environment Management
**Theory**: HM requires a consistent type environment (Γ) throughout inference
**Problem**: Our type environment doesn't properly track record field types by name

### 2. Improper Substitution Application
**Theory**: Substitutions must be applied consistently across all type expressions
**Problem**: Record field access applies substitutions inconsistently between construction and usage

### 3. Incorrect Unification Algorithm
**Theory**: Unification must handle structural types (records) properly with occurs check
**Problem**: Record type unification relies on field ordering instead of structural equivalence

## Proper Hindley-Milner Implementation Plan

### Phase 1: Fix Type Environment (Γ) Management

#### 1.1 Record Type Environment Consistency
```go
// Ensure record types are stored by field names, not positions
type RecordTypeEnvironment struct {
    name       string
    fieldTypes map[string]Type  // ONLY access by name
    // Remove fieldOrder - ordering must not matter
}
```

#### 1.2 Type Environment Operations
- **ADD**: `env.addRecord(name, fieldMap)` - stores field types by name only
- **LOOKUP**: `env.getFieldType(recordName, fieldName)` - retrieves by name
- **UNIFY**: `env.unifyRecords(r1, r2)` - structural equivalence by field names

### Phase 2: Correct Unification Algorithm

#### 2.1 Structural Record Unification
Based on HM theory, record unification must check:
```
unify(Record(fields1), Record(fields2)) =
  if |fields1| ≠ |fields2| then FAIL
  else for each field f:
    if f ∈ fields1 ∧ f ∈ fields2 then
      unify(type(fields1[f]), type(fields2[f]))
    else FAIL
```

#### 2.2 Function Type Unification Fix
The current error suggests function parameter types aren't being unified correctly:
```go
// Fix: Ensure parameter substitutions are applied before unification
func (ti *TypeInferer) unifyFunctionCall(fnType *FunctionType, argTypes []Type) error {
    // Apply current substitution to function type
    fnType = ti.applySubstitution(fnType)
    
    // Unify each parameter with argument
    for i, argType := range argTypes {
        if err := ti.Unify(fnType.paramTypes[i], argType); err != nil {
            return fmt.Errorf("parameter %d unification failed: %w", i, err)
        }
    }
    return nil
}
```

### Phase 3: Fix Pattern Matching Type Inference

#### 3.1 String Pattern Matching
Current issue: Empty string patterns return `constant.NewBool(false)` instead of string comparison.

**Fix**: Ensure pattern matching preserves discriminant type information:
```go
func (g *LLVMGenerator) createStringPatternCondition(pattern string, discriminant value.Value) value.Value {
    // For empty string patterns like ""
    if pattern == `""` {
        emptyStr := constant.NewCharArrayFromString("\0")
        emptyGlobal := g.module.NewGlobalDef("", emptyStr)
        emptyPtr := g.builder.NewGetElementPtr(emptyStr.Typ, emptyGlobal, ...)
        
        // Generate actual strcmp call, not constant false
        strcmp := g.functions["strcmp"]
        result := g.builder.NewCall(strcmp, discriminant, emptyPtr)
        zero := constant.NewInt(types.I32, 0)
        return g.builder.NewICmp(enum.IPredEQ, result, zero)
    }
    // ... handle other string patterns
}
```

#### 3.2 Pattern Type Consistency
Ensure pattern matching preserves HM type information:
- Pattern discriminant type must match expression type
- Pattern arms must have consistent result types
- Match exhaustiveness based on type structure, not ordering

### Phase 4: Record Field Access Specification

#### 4.1 Strict Name-Only Access
**SPECIFICATION UPDATE REQUIRED**:
```
RULE: Record field access is STRICTLY by name only.
NEVER by position, index, or ordering.

type Person = { name: string, age: int }
let p = Person { name: "Alice", age: 25 }

VALID:   p.name    // Access by name
VALID:   p.age     // Access by name  
INVALID: p[0]      // No positional access
INVALID: p.field[i] // No indexed access
```

#### 4.2 LLVM Code Generation Fix
Ensure struct field access uses consistent name-to-index mapping:
```go
type RecordFieldMapper struct {
    recordName string
    nameToIndex map[string]int
}

func (g *LLVMGenerator) getFieldIndex(recordType string, fieldName string) int {
    mapper := g.recordMappers[recordType]
    return mapper.nameToIndex[fieldName]
}
```

### Phase 5: Validation Function Type Consistency

#### 5.1 Record Constructor vs Validation Type Alignment
Current issue: Validation functions expect different struct layout than constructors create.

**Fix**: Ensure both use identical field name-to-LLVM-field mapping:
```go
// Single source of truth for record type mapping
func (g *LLVMGenerator) getRecordStructType(recordName string) *types.StructType {
    if cached, exists := g.structTypeCache[recordName]; exists {
        return cached
    }
    
    // Build struct type with consistent field ordering
    typeDecl := g.typeDeclarations[recordName]
    fieldTypes := make([]types.Type, len(typeDecl.Variants[0].Fields))
    fieldMap := make(map[string]int)
    
    for i, field := range typeDecl.Variants[0].Fields {
        fieldTypes[i] = g.getLLVMType(g.getInferenceFieldType(field.Type))
        fieldMap[field.Name] = i
    }
    
    structType := types.NewStruct(fieldTypes...)
    g.structTypeCache[recordName] = structType
    g.fieldMappings[recordName] = fieldMap
    
    return structType
}
```

## Implementation Steps

### Step 1: Type Environment Audit (Priority: Critical)
1. Remove all field ordering dependencies from `RecordType`
2. Implement name-only field access throughout codebase
3. Add validation to prevent positional access

### Step 2: Unification Algorithm Rewrite (Priority: Critical)
1. Fix `unifyFunctionTypes` to handle parameter type mismatches
2. Implement proper structural record unification
3. Add occurs check for recursive types

### Step 3: Pattern Matching Overhaul (Priority: High)
1. Fix string pattern conditions to generate proper comparisons
2. Ensure pattern discriminant types are preserved
3. Remove type coercion in pattern matching that breaks HM invariants

### Step 4: LLVM Generation Consistency (Priority: High)
1. Single source of truth for struct field layouts
2. Consistent field name-to-index mapping
3. Remove duplicate struct type generation

### Step 5: Integration Testing (Priority: Medium)
1. Add HM-specific unit tests for edge cases
2. Verify constraint validation works with proper types
3. Test polymorphic function instantiation

## Success Criteria

1. ✅ `pure_hindley_milner_test.osp` passes - no "int vs string" errors
2. ✅ `constraint_validation_test.osp` passes - proper LLVM IR generation
3. ✅ `working_constraint_test.osp` passes - validation functions work correctly
4. ✅ All field access uses names only - no positional dependencies
5. ✅ Type inference is sound and complete per HM theory

## Theoretical Foundation

This plan follows Hindley-Milner theory principles:

1. **Soundness**: If the type checker accepts a program, it won't have runtime type errors
2. **Completeness**: If a program has a valid typing, the type checker will find it  
3. **Principal Types**: Every well-typed expression has a most general type
4. **Decidability**: Type inference terminates and gives correct results

The fixes ensure our implementation adheres to these theoretical guarantees while maintaining practical usability.