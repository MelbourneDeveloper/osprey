## 8. Pattern Matching

### 8.1 Basic Patterns

```osprey
let result = match value {
    0 => "zero"
    1 => "one"
    n => "other: " + toString(n)
}
```

### 8.2 Union Type Patterns

```osprey
type Option = Some { value: Int } | None

let message = match option {
    Some x => "Value: " + toString(x.value)
    None => "No value"
}
```

### 8.3 Wildcard Patterns

The underscore `_` matches any value:

```osprey
let category = match score {
    100 => "perfect"
    90 => "excellent" 
    _ => "good"
}
```

### 8.4 Type Annotation Patterns

Type annotation patterns use the `:` operator to match values of specific types. This is **REQUIRED** for `any` types.

```
type_pattern := ID ':' type
structural_pattern := ID ':' '{' field_list '}'
anonymous_structural_pattern := '{' field_list '}'
constructor_pattern := ID ('(' pattern (',' pattern)* ')')?
variable_pattern := ID
wildcard_pattern := '_'
```

**Examples:**
```osprey
// Required for any types
match anyValue {
    num: Int => num + 1
    text: String => length(text)
    flag: Bool => if flag then 1 else 0
    _ => 0
}

// Structural matching - matches any type with these fields
match anyValue {
    { name, age } => print("${name}: ${age}")           // Anonymous structural
    p: { name, age } => print("Person ${p.name}: ${p.age}")  // Named structural
    u: User { id } => print("User ${id}")               // Traditional typed
    _ => print("Unknown")
}

// Advanced structural patterns
match anyValue {
    { x, y } => print("Point: (${x}, ${y})")           // Any type with x, y fields
    p: { name } => print("Named thing: ${p.name}")     // Any type with name field
    { id, email, active: Bool } => print("Active user: ${id}")  // Mixed field patterns
    _ => print("No match")
}

// Type patterns with field destructuring
match result {
    success: Success { value, timestamp } => processSuccess(value, timestamp)
    error: Error { code, message } => handleError(code, message)
    _ => defaultHandler()
}
```

### Pattern Matching Features

#### **1. Type Annotation Patterns**
```osprey
match anyValue {
    i: Int => i * 2                    // Bind as 'i' if Int
    s: String => s + "!"               // Bind as 's' if String
    user: User => user.name            // Bind as 'user' if User type
}
```

#### **2. Anonymous Structural Matching**
Match on structure without requiring specific type names:
```osprey
match anyValue {
    { name, age } => print("${name} is ${age}")        // ANY type with name, age
    { x, y, z } => print("3D point: ${x},${y},${z}")   // ANY type with x, y, z
    { id } => print("Has ID: ${id}")                    // ANY type with id field
}
```

#### **3. Named Structural Matching**
Bind the whole object AND destructure fields:
```osprey
match anyValue {
    person: { name, age } => {
        print("Person: ${person}")      // Access whole object
        print("Name: ${name}")          // Access destructured field
        print("Age: ${age}")            // Access destructured field
    }
    point: { x, y } => calculateDistance(point, origin)
}
```

#### **4. Mixed Type and Structural Patterns**
```osprey
match anyValue {
    user: User { id, name } => print("User ${id}: ${name}")     // Explicit type
    { email, active } => print("Has email: ${email}")           // Structural only
    data: { values: Array<Int> } => processArray(data.values)   // Nested types
    _ => print("Unknown structure")
}
```

#### **5. Partial Field Matching**
```osprey
match anyValue {
    { name, ... } => print("Has name: ${name}")        // Match any object with 'name'
    user: User { id, ... } => print("User ID: ${id}")  // User with at least 'id' field
    { x, y, ... } => print("At least 2D: ${x}, ${y}")  // Match with extra fields
}
```

### 8.5 Match Expression Type Safety Rules

**CRITICAL**: Osprey enforces strict type safety and exhaustiveness checking for match expressions.

#### 8.5.1 Type Compatibility Requirement

Match expressions must have **type-compatible** patterns. The expression being matched and all pattern arms must be of compatible types.

**✅ Valid - Compatible Types:**
```osprey
// Matching int against int patterns
let x = 42
let result = match x {
    0 => "zero"
    1 => "one"  
    _ => "other"
}

// Matching union type against its variants
type Color = Red | Green | Blue
let color = Red
let description = match color {
    Red => "red color"
    Green => "green color"
    Blue => "blue color"
}
```

**❌ Invalid - Type Mismatch:**
```osprey
// COMPILER ERROR: Type mismatch
let x = 42  // int type
type Option = Some { value: String } | None

let result = match x {  // ERROR: cannot match int against Option patterns
    Some => "some"      // Some is Option variant, not int
    None => "none"      // None is Option variant, not int
}
// Error: match expression type mismatch: cannot match expression of type 'int' against pattern of type 'Option'
```

#### 8.5.2 Exhaustiveness Checking

All match expressions **MUST** be exhaustive - every possible value must be handled.

**✅ Valid - Exhaustive:**
```osprey
type Status = Success | Error | Pending

let result = match status {
    Success => "completed"
    Error => "failed"  
    Pending => "waiting"  // All variants covered
}

// Or with wildcard
let result = match status {
    Success => "completed"
    _ => "not completed"  // Covers Error and Pending
}
```

**❌ Invalid - Non-Exhaustive:**
```osprey
type Color = Red | Green | Blue

let description = match color {
    Red => "red color"
    Green => "green color"
    // Missing Blue case!
}
// Error: match expression is not exhaustive: missing patterns: [Blue]
```

#### 8.5.3 Pattern Validity Rules

1. **Literal Patterns**: Must match the expression type
2. **Constructor Patterns**: Must be valid variants of the union type
3. **Variable Patterns**: Capture the matched value
4. **Wildcard Pattern**: Must be the last arm if present

**❌ Invalid Examples:**
```osprey
// Unknown variant error
type Color = Red | Green | Blue
let result = match color {
    Red => "red"
    Green => "green"
    Blue => "blue"
    Purple => "purple"  // ERROR: Purple not a variant of Color
}
// Error: unknown variant 'Purple' is not defined in type 'Color'

// Wildcard not last
let result = match color {
    _ => "any color"    // ERROR: wildcard must be last
    Red => "red"
}
// Error: wildcard pattern must be the last arm

// Duplicate patterns
let result = match color {
    Red => "red"
    Green => "green"  
    Red => "also red"   // ERROR: duplicate pattern
}
// Error: duplicate match arm: pattern 'Red' appears multiple times
```

#### 8.5.4 Compilation Error Messages

The compiler provides specific error messages for match violations:

```osprey
// Type mismatch errors
"match expression type mismatch: cannot match expression of type 'T1' against pattern of type 'T2'"

// Exhaustiveness errors  
"match expression is not exhaustive: missing patterns: [Pattern1, Pattern2]"

// Unknown variant errors
"unknown variant 'VariantName' is not defined in type 'TypeName'"

// Pattern ordering errors
"wildcard pattern must be the last arm"

// Duplicate pattern errors
"duplicate match arm: pattern 'PatternName' appears multiple times"
```

#### 8.5.5 Implementation Status

🔄 **PATTERN MATCHING IMPLEMENTATION STATUS**:

**Currently Implemented:**
- ✅ Basic pattern matching with literals and identifiers
- ✅ Variable capture patterns
- ✅ Wildcard patterns (`_`)
- ✅ Type annotation patterns (`value: Int`)
- ✅ Named structural patterns (`person: { name, age }`)
- ✅ Anonymous structural patterns (`{ name, age }`)

**🚧 PARTIAL IMPLEMENTATION:**
- 🔄 Exhaustiveness checking for union types (in progress)
- 🔄 Unknown variant detection (error messages implemented)

**❌ NOT YET IMPLEMENTED:**
- ❌ Type compatibility checking between expression and patterns
- ❌ Constructor pattern validation
- ❌ Duplicate pattern detection
- ❌ Wildcard position validation

**Testing**: Examples in `examples/failscompilation/*.ospo` test these error conditions. Some tests are currently skipped as the features are not yet implemented.