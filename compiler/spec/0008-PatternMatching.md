8. [Pattern Matching](0008-PatternMatching.md)
   - [Basic Patterns](#81-basic-patterns)
   - [Union Type Patterns](#82-union-type-patterns)
   - [Wildcard Patterns](#83-wildcard-patterns)
   - [Type Annotation Patterns](#84-type-annotation-patterns)
       - [Type Annotation Patterns](#1-type-annotation-patterns)
       - [Anonymous Structural Matching](#2-anonymous-structural-matching)
       - [Named Structural Matching](#3-named-structural-matching)
       - [Mixed Type and Structural Patterns](#4-mixed-type-and-structural-patterns)
   - [Match Expression Type Safety Rules](#85-match-expression-type-safety-rules)

## 8. Pattern Matching

### 8.1 Basic Patterns

```osprey
let result = match value {
    0 => "zero"
    1 => "one"
    n => "other: " + toString(n)
}
```

## 8.2 Union Type Patterns

```osprey
type Option = Some { value: Int } | None

let message = match option {
    Some x => "Value: " + toString(x.value)
    None => "No value"
}
```

## 8.3 Wildcard Patterns

The underscore `_` matches any value:

```osprey
let category = match score {
    100 => "perfect"
    90 => "excellent" 
    _ => "good"
}
```

## 8.4 Type Annotation Patterns

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

## Pattern Matching Features

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

## 8.5 Result Type Pattern Matching (Arithmetic Expressions)

**🔥 CRITICAL**: All arithmetic expressions return `Result<T, MathError>`. You **MUST** handle them with pattern matching.

### Simple Arithmetic Result Handling
```osprey
let calculation = 1 + 3 + (300 / 5)  // Result<int, MathError>

match calculation {
    Success { value } => print("Result: ${value}")
    Error { message } => print("Math error: ${message}")
}
```

### Compound Expression Examples (CRYSTAL CLEAR)
```osprey
// Each of these returns a SINGLE Result for the ENTIRE expression
let simple = 10 + 5                    // Result<int, MathError>
let complex = 1 + 2 * 3 - 4 / 2        // Result<int, MathError>  
let nested = ((a + b) * c) / (d - e)   // Result<int, MathError>

// Handle ALL of them the SAME WAY
match simple {
    Success { value } => print("10 + 5 = ${value}")
    Error { message } => print("Failed: ${message}")
}

match complex {
    Success { value } => print("Complex calc = ${value}")
    Error { message } => print("Overflow/error: ${message}")
}

match nested {
    Success { value } => print("Nested result = ${value}")
    Error { message } => print("Division by zero or overflow: ${message}")
}
```

### Function Return Results
```osprey
fn calculate(x: int, y: int) -> Result<int, MathError> = x + y * 2 - 5

let result = calculate(10, 3)  // Result<int, MathError>
match result {
    Success { value } => print("Function result: ${value}")
    Error { message } => print("Function failed: ${message}")
}
```

### Advanced Result Chains
```osprey
// Multiple Results in sequence
let step1 = 100 + 50        // Result<int, MathError>
let step2 = 200 * 3         // Result<int, MathError>

// Handle each step
match step1 {
    Success { value1 } => {
        match step2 {
            Success { value2 } => {
                let final = value1 + value2  // This is also Result<int, MathError>!
                match final {
                    Success { total } => print("Final: ${total}")
                    Error { message } => print("Final calc failed: ${message}")
                }
            }
            Error { message } => print("Step 2 failed: ${message}")
        }
    }
    Error { message } => print("Step 1 failed: ${message}")
}
```

## 8.6 Match Expression Type Safety Rules

```

## 8.7 Ternary Match Expression (Syntactic Sugar)

To reduce verbosity for common two-armed match scenarios, Osprey provides a concise ternary match expression. This is **purely syntactic sugar** and desugars to a standard `match` expression internally.

**Syntax:**
`<expression> { <pattern> } ? <then_expr> : <else_expr>`

This is exactly equivalent to:
```osprey
match <expression> {
    { <pattern> } => <then_expr>
    _ => <else_expr>
}
```

### Breakdown

- **`<expression>`**: The value to be matched.
- **`{ <pattern> }`**: A structural pattern to match against the expression. This can be used for destructuring.
- **`? <then_expr>`**: The expression to evaluate if the pattern matches.
- **`: <else_expr>`**: The expression to evaluate if the pattern does not match.

### Examples

**Example 1: Handling a `Result` Type**

Given the type: `type Result = Success { value: String } | Error`

The ternary match provides a clean way to extract a value or provide a default.

```osprey
let myResult = Success { value: "It worked" }

// Extracts "It worked" using the ternary match
let message = myResult { value } ? value : "An unknown error occurred"
// message is now "It worked"
```

**Example 2: Handling Booleans**

The ternary match can also work with boolean values by matching on the implicit structure of a `true` value.

```osprey
let is_active = true
let status_text = is_active ? "Active" : "Inactive"
// status_text is now "Active"
```