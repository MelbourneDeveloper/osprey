---
layout: page
title: "Type System"
description: "Osprey Language Specification: Type System"
date: 2025-07-06
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0005-typesystem/"
---

5. [Type System](0005-TypeSystem.md)
   - [Built-in Types](#51-built-in-types)
       - [Function Types](#function-types)
   - [Built-in Error Types](#52-built-in-error-types)
   - [Type Inference Rules](#53-type-inference-rules)
       - [Function Return Types](#function-return-types)
       - [Parameter Types](#parameter-types)
       - [Type Inference Examples](#type-inference-examples)
       - [Rationale](#rationale)
       - [Function Return Type "any" Restriction](#function-return-type-any-restriction)
       - [Common Validation Fixes](#common-validation-fixes)
   - [Type Safety and Explicit Typing](#54-type-safety-and-explicit-typing)
       - [Mandatory Type Safety](#mandatory-type-safety)
   - [Any Type Handling and Pattern Matching Requirement](#55-any-type-handling-and-pattern-matching-requirement)
       - [Forbidden Operations on `any` Types](#forbidden-operations-on-any-types)
       - [Legal Operations on `any` Types](#legal-operations-on-any-types)
       - [Pattern Matching Requirement](#pattern-matching-requirement)
       - [Direct Access Compilation Errors](#direct-access-compilation-errors)
       - [Function Return Type Handling](#function-return-type-handling)
       - [Type Annotation Pattern Syntax](#type-annotation-pattern-syntax)
       - [Compilation Error Messages](#compilation-error-messages)
       - [Exhaustiveness Checking for Any Types](#exhaustiveness-checking-for-any-types)
       - [Default Wildcard Behavior for Any Types](#default-wildcard-behavior-for-any-types)
       - [Type Constraint Checking](#type-constraint-checking)
       - [Context-Aware Type Validation](#context-aware-type-validation)
       - [Compilation Errors for Impossible Types](#compilation-errors-for-impossible-types)
       - [Performance and Safety Characteristics](#performance-and-safety-characteristics)
       - [Type Annotation Requirements](#type-annotation-requirements)
       - [Compilation Errors for Type Ambiguity](#compilation-errors-for-type-ambiguity)
       - [Error Handling Requirements](#error-handling-requirements)
   - [Type Compatibility](#56-type-compatibility)

## 5. Type System

Osprey's type system is one of the most important aspects of the language. It puts type safety and expressiveness as the top priorities. It is inspired by popular functional programming languages like Haskell and ML. The type system aims towards making illegal states unrepresentable. Some features like constructor where constraints are critical features that mean that instances cannot exist where they fail the criteria in construction.

### 5.1 Built-in Types

- `Int`: 64-bit signed integers
- `String`: UTF-8 encoded strings  
- `Bool`: Boolean values (`true`, `false`)
- `Unit`: Type for functions that don't return a meaningful value
- `Result<T, E>`: Built-in generic type for error handling
- `List<T, N>`: Immutable fixed-size lists with N elements of type T
- `Function Types`: First-class function types with syntax `(T1, T2, ...) -> R`

#### Function Types

Function types represent functions as first-class values, enabling higher-order functions and function composition.

**Syntax:**
```
FunctionType := '(' (Type (',' Type)*)? ')' '->' Type
```

**Examples:**
```osprey
(Int) -> Int              // Function taking an Int, returning an Int
(Int, String) -> Bool     // Function taking Int and String, returning Bool
() -> String              // Function with no parameters, returning String
(String) -> (Int) -> Bool // Higher-order function returning another function
```

**Function Type Declarations:**
```osprey
// Function parameter with explicit function type
fn applyFunction(value: Int, transform: (Int) -> Int) -> Int = 
    transform(value)

// Variable with function type
let doubler: (Int) -> Int = fn(x: Int) -> Int = x * 2

// Higher-order function that returns a function
fn createAdder(n: Int) -> (Int) -> Int = 
    fn(x: Int) -> Int = x + n
```

**Function Composition Examples:**
```osprey
// Define some simple functions
fn double(x: Int) -> Int = x * 2
fn square(x: Int) -> Int = x * x
fn addFive(x: Int) -> Int = x + 5

// Higher-order function with strong typing
fn applyTwice(value: Int, func: (Int) -> Int) -> Int = 
    func(func(value))

// Usage with type safety
let result1 = applyTwice(5, double)  // 20
let result2 = applyTwice(3, square)  // 81
let result3 = applyTwice(10, addFive) // 20

// Composition of functions
fn compose(f: (Int) -> Int, g: (Int) -> Int) -> (Int) -> Int =
    fn(x: Int) -> Int = f(g(x))

let doubleSquare = compose(double, square)
let result4 = doubleSquare(3) // double(square(3)) = double(9) = 18
```

**Type Safety Benefits:**
- **Compile-time validation**: Function signatures are checked at compile time
- **No runtime type errors**: Mismatched function types caught early
- **Clear documentation**: Function types serve as documentation
- **Enables optimization**: Compiler can optimize based on known function signatures

### 5.2 Built-in Error Types

- `MathError`: For arithmetic operations (DivisionByZero, Overflow, Underflow)
- `ParseError`: For string parsing operations  
- `IndexError`: For list/string indexing operations (OutOfBounds)
- `Success`: Successful result wrapper

### 5.3 Type Inference Rules

**Core Principle**: The `any` type is invalid unless explicitly declared. All types must be either explicitly annotated or inferrable from context.

#### Function Return Types

Return type annotations may be omitted **only** when the return type can be definitively inferred from the function body:

##### Allowed (Return Type Inferred)
- **Literal expressions**: `fn getNumber() = 42` → infers `int`
- **String literals**: `fn getText() = "hello"` → infers `string`  
- **Boolean literals**: `fn getFlag() = true` → infers `bool`
- **Arithmetic expressions**: `fn calculate() = 1 + 2` → infers `int`

##### Disallowed (Requires Explicit Return Type)
- **Direct parameter return**: `fn identity(x) = x` → **ERROR** (would be `any`)
- **Function calls**: `fn process() = someFunction()` → **ERROR** (unknown return type)
- **Complex expressions**: Without clear type resolution

#### Parameter Types

Parameter type annotations may be omitted **only** when the parameter type can be definitively inferred from usage within the function body:

##### Allowed (Parameter Type Inferred)
- **Arithmetic usage**: `fn addOne(x) = x + 1` → `x` infers as `int`
- **With explicit return type**: `fn identity(x) -> int = x` → `x` infers as `int` from return type
- **Direct return with explicit type**: `fn process(data) -> string = data` → `data` infers as `string`

##### Disallowed (Requires Explicit Parameter Type)
- **Direct return without return type**: `fn identity(x) = x` → **ERROR**
- **String parameter return**: `fn greet(name) = name` → **ERROR** (use `name: string` or `-> string`)
- **Mixed parameter scenarios**: `fn formatScore(name, score) = name` → **ERROR**
- **Ambiguous usage**: Where type cannot be determined from context

#### Type Inference Examples

**Valid Code:**
```osprey
// Literals allow return type inference
fn getAge() = 25
fn getName() = "Alice"
fn isActive() = true

// Arithmetic allows both return and parameter type inference  
fn increment(x) = x + 1
fn add(a, b) = a + b

// Explicit types always allowed
fn identity(x) -> int = x
fn process(data: string) -> string = data
```

**Invalid Code:**
```osprey
// ERROR: Cannot infer return type from parameter
fn identity(x) = x

// ERROR: String parameter without type annotation
fn greet(name) = name

// ERROR: Mixed parameters without explicit types
fn formatScore(name, score) = name

// ERROR: Cannot infer parameter type from function call
fn process(data) = someFunction(data)

// ERROR: Ambiguous type inference
fn conditional(flag, a, b) = if flag then a else b
```

#### Rationale

This design ensures:
1. **Type Safety**: No implicit `any` types that could lead to runtime errors
2. **Readability**: Clear type contracts without excessive annotation
3. **Maintainability**: Predictable type behavior for code evolution
4. **Performance**: Compile-time type checking without runtime overhead

**Summary Rule**: "Type annotations may be omitted only when the type can be unambiguously determined from constants, literals, well-defined operations, or explicit return types that constrain parameter types."

#### Function Return Type "any" Restriction

**CRITICAL RULE**: Functions CANNOT return `any` type unless the return type is EXPLICITLY declared as `any`.

**✅ ALLOWED - Explicit any return type:**
```osprey
fn parseValue(input: string) -> any = processInput(input)
fn getDynamicValue() -> any = readFromConfig()
```

**❌ FORBIDDEN - Implicit any return type:**
```osprey
fn identity(x) = x                    // ERROR: Would infer as 'any'
fn callUnknown() = someFunction()     // ERROR: Would infer as 'any'
fn processData(data) = data           // ERROR: Would infer as 'any'
```

**Rationale**: This prevents accidental `any` type propagation that could lead to runtime type errors and maintains Osprey's strong type safety guarantees.

**Built-in Functions**: No built-in functions return `any` type. All built-in functions have concrete, well-defined return types.

#### Common Validation Fixes

When the compiler reports type inference errors, use these patterns:

```osprey
// ❌ ERROR: Function 'greet' requires explicit return type annotation
fn greet(name) = name

// ✅ FIX: Add explicit parameter type
fn greet(name: string) = name

// ✅ FIX: Add explicit return type  
fn greet(name) -> string = name

// ❌ ERROR: Parameter 'name' requires explicit type annotation
fn formatScore(name, score) = name

// ✅ FIX: Add explicit parameter types
fn formatScore(name: string, score: int) = name

// ✅ FIX: Add explicit return type to enable inference
fn formatScore(name, score) -> string = name
```

### 5.4 Type Safety and Explicit Typing

**CRITICAL RULE**: Osprey is fully type-safe with no exceptions.

#### Mandatory Type Safety
- **No implicit type conversions** - all type mismatches are compile-time errors
- **No runtime type errors** - all type issues caught at compile time
- **No panics or exceptions** - all error conditions must be handled explicitly

### 5.5 Any Type Handling and Pattern Matching Requirement

🔄 **IMPLEMENTATION STATUS**: `any` type validation is partially implemented. Basic validation for function arguments is working, but complete pattern matching enforcement is in progress.

Osprey provides the `any` type for maximum flexibility, but enforces strict access rules to maintain type safety. Direct access to `any` types is forbidden in most contexts - all `any` values must be accessed through pattern matching to extract their actual types.

#### Forbidden Operations on `any` Types

The following operations on `any` types will result in compilation errors:

1. **Direct variable access** - `any` variables cannot be used directly
2. **Function arguments** - `any` values cannot be passed to functions expecting concrete types  
3. **Field access** - Properties cannot be accessed directly on `any` types
4. **Implicit conversions** - `any` cannot be implicitly converted to other types

#### Legal Operations on `any` Types

**Arithmetic operations** with `any` types are explicitly allowed and return `Result` types:

```osprey
let x: any = 42
let result = x + 5  // Returns Result<Int, ArithmeticError>

let y: any = "hello" 
let sum = y + 10    // Returns Result<Int, TypeError>
```

These operations are safe because they return `Result` types that encapsulate potential errors, maintaining type safety while allowing flexible computation.

#### Pattern Matching Requirement

**Pattern Matching on `any` Types:**
```osprey
// Pattern matching on any type
match anyValue {
    value: Int => handleInteger(value)
    value: String => handleString(value)
    value: Bool => handleBoolean(value)
    _ => handleUnknownType()
}
```

#### Direct Access Compilation Errors

**❌ FORBIDDEN - Direct Access:**
```osprey
fn processAny(value: any) -> int = value + 1
// ERROR: cannot use 'any' type directly in arithmetic operation

fn getLength(value: any) -> int = value.length
// ERROR: cannot access field on 'any' type without pattern matching

let result: int = someAnyFunction()
// ERROR: cannot assign 'any' to 'int' without pattern matching

fn callFunction(value: any) = someFunction(value)
// ERROR: cannot pass 'any' type to function expecting specific type

let converted = toString(value)  // where value: any
// ERROR: cannot implicitly convert 'any' to expected parameter type
```

**✅ REQUIRED - Pattern Matching:**
```osprey
fn processAny(value: any) -> int = match value {
    num: Int => num + 1
    str: String => length(str)
    _ => 0
}

fn getLength(value: any) -> int = match value {
    str: String => length(str)
    arr: Array<T> => arrayLength(arr)
    _ => 0
}

let result = match someAnyFunction() {
    value: Int => value
    _ => 0
}
```

#### Function Return Type Handling

Functions returning `any` types require immediate pattern matching:

```osprey
// Function that returns any
extern fn parseValue(input: string) -> any

// ERROR: Direct usage
let number = parseValue("42") + 1

// CORRECT: Pattern matching
let number = match parseValue("42") {
    value: Int => value + 1
    _ => 0
}
```

#### Type Annotation Pattern Syntax

The `:` operator is used for type annotation in patterns:
- `value: Int` - Matches if value is an Int, binds to `value`
- `text: String` - Matches if value is a String, binds to `text`
- `flag: Bool` - Matches if value is a Bool, binds to `flag`
- `{ name, age }` - Structural match on any type with name and age fields
- `person: { name, age }` - Named structural match, binds whole object and fields
- `_` - Wildcard matches any remaining type

#### Compilation Error Messages

The compiler **MUST** emit these specific errors for `any` type violations:

```osprey
// Direct arithmetic operation
"cannot use 'any' type directly in arithmetic operation - pattern matching required"

// Direct field access
"cannot access field on 'any' type without pattern matching"

// Direct assignment to typed variable
"cannot assign 'any' to 'TYPE' without pattern matching"

// Direct function argument
"cannot pass 'any' type to function expecting 'TYPE' - pattern matching required"

// Implicit conversion attempt
"cannot implicitly convert 'any' to 'TYPE' - use pattern matching to extract specific type"

// Variable access on any
"cannot access variable of type 'any' directly - pattern matching required"

// Missing pattern match arm
"pattern matching on 'any' type must handle all possible types or include wildcard"

// Impossible type patterns
"pattern 'TYPE' is not a possible type for expression of documented types [TYPE1, TYPE2, ...]"

// Unreachable patterns
"unreachable pattern: 'TYPE' cannot occur based on context analysis"
```

#### Exhaustiveness Checking for Any Types

Pattern matching on `any` types **MUST** be exhaustive:
- Handle all expected types, OR
- Include a wildcard pattern (`_`) to handle unexpected types

```osprey
// Non-exhaustive (ERROR)
match anyValue {
    value: Int => processInt(value)
    value: String => processString(value)
    // ERROR: missing wildcard or Bool case
}

// Exhaustive (CORRECT)
match anyValue {
    value: Int => processInt(value)
    value: String => processString(value)
    _ => handleOther()
}
```

#### Default Wildcard Behavior for Any Types

The wildcard pattern (`_`) in `any` type matching preserves the `any` type:

```osprey
// Wildcard returns any type
let result = match someAnyValue {
    value: Int => processInt(value)    // Returns specific type
    value: String => processString(value)  // Returns specific type
    _ => someAnyValue  // Returns any type (unchanged)
}
// result type: any (due to wildcard arm)

// To avoid any type in result, handle all expected cases explicitly
let result = match someAnyValue {
    value: Int => processInt(value)
    value: String => processString(value)
    _ => defaultInt()  // Convert to specific type
}
// result type: Int (all arms return Int)
```

#### Type Constraint Checking

The compiler **MUST** validate that pattern types are actually possible for the value being matched:

**✅ VALID - Realistic Type Patterns:**
```osprey
// Function known to return Int or String
extern fn parseIntOrString(input: string) -> any

match parseIntOrString("42") {
    value: Int => value + 1
    value: String => length(value)
    _ => 0  // Valid: handles any unexpected types
}
```

**❌ INVALID - Impossible Type Patterns:**
```osprey
// Function documented to only return Int or String
extern fn parseIntOrString(input: string) -> any

match parseIntOrString("42") {
    value: Int => value + 1
    value: String => length(value)
    value: Bool => if value then 1 else 0  // ERROR: Bool not possible
    _ => 0
}
// ERROR: pattern 'Bool' is not a possible type for function 'parseIntOrString'
```

#### Context-Aware Type Validation

When the compiler has information about possible types (from documentation, extern declarations, or analysis), it **MUST** enforce realistic pattern matching:

```osprey
// Extern function with documented return types
extern fn getUserInput() -> any  // Documentation: returns Int | String | Bool only

// VALID: Only realistic types
match getUserInput() {
    value: Int => processInt(value)
    value: String => processString(value) 
    value: Bool => processBool(value)
    _ => handleUnexpected()  // Still allowed for safety
}

// INVALID: Unrealistic types
match getUserInput() {
    value: Int => processInt(value)
    value: Array<String> => processArray(value)  // ERROR: Array not documented
    _ => handleOther()
}
// ERROR: pattern 'Array<String>' is not a documented return type for 'getUserInput'
```

#### Compilation Errors for Impossible Types

```osprey
"pattern 'TYPE' is not a possible type for expression of documented types [TYPE1, TYPE2, ...]"
"unreachable pattern: 'TYPE' cannot occur based on context analysis"
"pattern matching includes impossible type 'TYPE' - check function documentation"
```

#### Performance and Safety Characteristics

- **Compile-time type checking**: Pattern matching enables compile-time verification
- **Zero runtime cost**: Type patterns compiled to efficient type tags
- **Memory safety**: No type confusion or invalid casts possible
- **Explicit control**: Developers must explicitly handle all type cases

#### Type Annotation Requirements
When the compiler cannot infer types, explicit type annotations are **REQUIRED**:

```osprey
// Type annotations required when inference is ambiguous
fn complexOperation(data: String, count: Int) = processData(data, count)

// Generic functions require type parameters
fn parseValue<T>(input: String) -> Result<T, ParseError> = ...

// Union types with fields require explicit typing
type Result<T, E> = Ok { value: T } | Err { error: E }
```

#### Compilation Errors for Type Ambiguity
The compiler **MUST** emit errors when:
1. Function parameter types cannot be inferred from usage
2. Return types are ambiguous
3. Variable types cannot be determined from initializers
4. Generic type parameters are not specified

#### Error Handling Requirements
- **No exceptions or panics** - all failing operations return Result types
- **Explicit error handling** - all Result types must be pattern matched
- **Safe arithmetic** - operations like division must return Result<T, Error>

```osprey
// REQUIRED: Safe division that cannot panic
fn safeDivide(a: Int, b: Int) -> Result<Int, MathError> = match b {
  0 => Err { error: DivisionByZero }
  _ => Ok { value: a / b }
}

// REQUIRED: All results must be handled
let result = safeDivide(a: 10, b: 2)
match result {
  Ok { value } => print("Result: ${value}")
  Err { error } => handleError(error)
}
```

### 5.6 Type Compatibility

- Pattern matching for type discrimination
- Union types for representing alternatives
- Result types for error handling instead of exceptions