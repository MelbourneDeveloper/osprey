# Plan to Fix Discriminated Unions and Where Constraint Validation

## **CRITICAL ISSUES TO FIX**

### 1. **Discriminated Union Support (Currently Broken)**
**Problem**: `declareType()` treats all unions as simple i64 enums, completely ignoring fields
**Example**: `type Character = Warrior { health: Int, strength: Int } | Mage { mana: Int, spells: Int }`

**What I Started**:
- ✅ Updated `declareType()` to detect discriminated unions vs simple enums
- ✅ Added `declareDiscriminatedUnion()` that creates tagged union structures
- ✅ Added helper functions `storeVariantFieldNames()` and `getVariantFieldNames()`

**Still Need to Fix**:
- **Type Constructor Generation**: Update `generateTypeConstructorExpression()` to handle discriminated union variant construction
- **Pattern Matching**: Update match expression generation to extract fields from tagged unions
- **Field Access**: Handle field access on discriminated union variants (should fail at compile time)
- **Variant Construction**: Generate proper tagged union construction code

### 2. **Where Constraint Validation (Parser Works, Codegen Doesn't)**
**Problem**: Parser correctly parses `where validateAge` but codegen never calls validation functions
**Example**: `type Person { name: String, age: Int } where validateAge`

**What Exists**:
- ✅ Parser correctly stores `ValidationFunc` in `TypeDeclaration`
- ✅ `generateRecordTypeConstructor()` has skeleton code to call validation

**Still Need to Fix**:
- **Validation Function Calls**: Actually implement calling the validation function
- **Result Type Generation**: Generate proper `Result<T, String>` return types
- **Error Handling**: Handle validation failures properly
- **Integration Testing**: Create working examples that demonstrate validation

### 3. **Type Constructor Expression Generation**
**File**: `compiler/internal/codegen/expression_generation.go:760-820`

**Current Issues**:
```go
// For union types without fields, return discriminant value
return constant.NewInt(types.I64, 1), nil  // ❌ WRONG - doesn't handle discriminated unions
```

**Need to Add**:
- **Discriminated Union Constructor**: Handle `Warrior { health: 100, strength: 50 }`
- **Tagged Union Assembly**: Create proper tagged union with discriminant + field data
- **Field Validation**: Ensure all required fields are provided
- **Memory Layout**: Proper serialization of variant fields into union data area

## **DETAILED IMPLEMENTATION PLAN**

### **STEP 1: Fix Type Constructor Generation**
```go
// In generateTypeConstructorExpression()
if len(typeDecl.Variants) > 1 {
    // This is a discriminated union - find which variant we're constructing
    return g.generateDiscriminatedUnionConstructor(typeConstructor, typeDecl)
}
```

**New Function Needed**: `generateDiscriminatedUnionConstructor()`
- Find which variant matches the constructor call
- Create tagged union structure
- Set discriminant (tag) field
- Serialize variant fields into data area
- Handle both constrained and unconstrained variants

### **STEP 2: Fix Where Constraint Validation**
```go
// In generateRecordTypeConstructor() - currently has skeleton code
if typeDecl.ValidationFunc != nil {
    // ❌ Current code calls validation but doesn't handle Result properly
    validationFunc, exists := g.functions[*typeDecl.ValidationFunc]
    if !exists {
        return nil, WrapUndefinedFunction(*typeDecl.ValidationFunc)
    }
    result := g.builder.NewCall(validationFunc, structValue)
    return result, nil  // ❌ This is wrong - need proper Result<T,E> handling
}
```

**Fix Required**:
- Generate proper `Result<T, String>` types
- Handle validation success/failure cases
- Return appropriate LLVM structures

### **STEP 3: Update Pattern Matching for Discriminated Unions And Structural Typing**
**File**: Need to update match generation code

**Current Issue**: Pattern matching assumes simple discriminants, not tagged unions

**Fix Required**:
- Extract discriminant from tagged union
- Extract and bind variant fields for patterns like `Warrior { health, strength } => ...`
- Handle nested field destructuring

- We could throw anything at a pattern match and we must be able to match on the fields

### **STEP 4: Add Variant Construction Functions**
**Need**: Helper functions to construct discriminated union variants

```go
func (g *LLVMGenerator) generateDiscriminatedUnionConstructor(
    typeConstructor *ast.TypeConstructorExpression, 
    typeDecl *ast.TypeDeclaration
) (value.Value, error)

func (g *LLVMGenerator) findVariantByConstructorCall(
    typeConstructor *ast.TypeConstructorExpression,
    typeDecl *ast.TypeDeclaration
) (*ast.TypeVariant, int, error)

func (g *LLVMGenerator) serializeVariantFields(
    variant *ast.TypeVariant,
    fieldValues map[string]Expression,
    unionDataPtr value.Value
) error
```

### **STEP 5: Fix Field Access Validation**
**File**: `compiler/internal/codegen/expression_generation.go:633`

**Current Issue**: Field access on discriminated unions should be compilation error
**Fix**: Detect discriminated union field access and generate proper error

### **STEP 6: Integration and Testing**
- Update failing test examples to use correct syntax
- Create new test cases for discriminated unions with fields
- Create test cases for where constraint validation
- Ensure `make clean && make build && make test` passes

## **FILES TO MODIFY**

1. **`compiler/internal/codegen/function_signatures.go`** - ✅ DONE (discriminated union type declaration)
2. **`compiler/internal/codegen/expression_generation.go`** - ❌ CRITICAL (type constructor generation)
3. **`compiler/internal/codegen/match_generation.go`** - ❌ TODO (pattern matching for discriminated unions)
4. **`compiler/internal/codegen/core_functions.go`** - ❌ TODO (Result type handling)
5. **Test files** - ❌ TODO (update examples with correct syntax)

## **TESTING STRATEGY**

### **Test Cases Needed**:
1. **Simple discriminated union**: `type Color = Red | Green | Blue`
2. **Discriminated union with fields**: `type Character = Warrior { health: Int } | Mage { mana: Int }`
3. **Where constraint validation**: `type Person { age: Int } where validateAge`
4. **Pattern matching on discriminated unions**: `match character | Warrior { health } => ...`
5. **Failed validation**: Test validation function returning error

### **Commands to Run After Each Change**:
```bash
make clean
make build
make test
make install  # Update VSCode integration
```

## **PRIORITY ORDER**
1. **URGENT**: Fix `generateTypeConstructorExpression()` for discriminated unions
2. **URGENT**: Implement where constraint validation properly
3. **HIGH**: Update pattern matching for discriminated unions
4. **MEDIUM**: Fix field access validation
5. **LOW**: Clean up test cases

**REMEMBER**: The user is frustrated with surface-level fixes. Focus on **ACTUAL FUNCTIONALITY** not just documentation updates!