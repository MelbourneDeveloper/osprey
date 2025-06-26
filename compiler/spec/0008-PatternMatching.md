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

## 8.5 Match Expression Type Safety Rules
