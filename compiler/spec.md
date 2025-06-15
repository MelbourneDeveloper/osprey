**Version:** 0.1.0-alpha

**Date:** June 2025

**Author** Christian Findlay

<div class="table-of-contents">

1. [Introduction](#1-introduction)
   - [Completeness](#11-completeness)
   - [Principles](#12-principles)
2. [Lexical Structure](#2-lexical-structure)
   - [Identifiers](#21-identifiers)
   - [Keywords](#22-keywords)
   - [Literals](#23-literals)
   - [Operators](#24-operators)
   - [Delimiters](#25-delimiters)
3. [Syntax](#3-syntax)
   - [Program Structure](#31-program-structure)
   - [Import Statements](#32-import-statements)
   - [Let Declarations](#33-let-declarations)
   - [Function Declarations](#34-function-declarations)
   - [Extern Declarations](#35-extern-declarations)
   - [Type Declarations](#36-type-declarations)
   - [Record Types and Type Constructors](#37-record-types-and-type-constructors)
   - [Expressions](#38-expressions)
   - [Block Expressions](#39-block-expressions)
   - [Match Expressions](#310-match-expressions)
4. [Semantics](#4-semantics)
5. [Type System](#5-type-system)
   - [Built-in Types](#51-built-in-types)
   - [Built-in Error Types](#52-built-in-error-types)
   - [Type Inference Rules](#53-type-inference-rules)
   - [Type Safety and Explicit Typing](#54-type-safety-and-explicit-typing)
   - [Any Type Handling and Pattern Matching Requirement](#55-any-type-handling-and-pattern-matching-requirement)
   - [Type Compatibility](#56-type-compatibility)
6. [Function Calls](#6-function-calls)
7. [String Interpolation](#7-string-interpolation)
8. [Pattern Matching](#8-pattern-matching)
9. [Block Expressions](#9-block-expressions)
10. [Boolean Operations](#10-boolean-operations)
11. [Loop Constructs and Functional Iterators](#11-loop-constructs-and-functional-iterators)
12. [Lightweight Fibers and Concurrency](#12-lightweight-fibers-and-concurrency)
13. [Fiber-Isolated Module System](#13-fiber-isolated-module-system)
14. [Built-in Functions](#14-built-in-functions)
    - [HTTP Core Types](#141-http-core-types)
    - [HTTP Server Functions](#142-http-server-functions)
    - [HTTP Client Functions](#143-http-client-functions)
    - [WebSocket Support](#144-websocket-support-two-way-communication)
    - [Streaming Response Bodies](#145-streaming-response-bodies)
    - [Error Handling in HTTP](#146-error-handling-in-http)
    - [Fiber-Based Concurrency](#147-fiber-based-concurrency)
    - [Complete HTTP Server Example](#148-complete-http-server-example)
15. [Error Handling](#15-error-handling)
16. [Examples](#16-examples)
17. [Built-in Functions Reference](#17-built-in-functions-reference)
    - [Basic I/O Functions](#171-basic-io-functions)
    - [Functional Iterator Functions](#172-functional-iterator-functions)
    - [Pipe Operator](#173-pipe-operator)
    - [Functional Programming Patterns](#174-functional-programming-patterns)
    - [Fiber Types and Concurrency](#175-fiber-types-and-concurrency)
18. [Security and Sandboxing](#18-security-and-sandboxing)
    - [Security Flags](#181-security-flags)
    - [Security Policies](#182-security-policies)
    - [Blocked Functions by Category](#183-blocked-functions-by-category)
    - [Error Messages](#184-error-messages)
    - [Programming Best Practices](#185-programming-best-practices)
    - [Implementation Details](#186-implementation-details)

</div>

## 1. Introduction

Osprey is a modern functional programming oriented language designed for elegance, safety, and performance.. It emphasizes:

- **Named arguments** for multi-parameter functions to improve readability
- **Strong type inference** to reduce boilerplate while maintaining safety
- **String interpolation** for convenient text formatting
- **Pattern matching** for elegant conditional logic
- **Immutable-by-default** variables with explicit mutability
- **Fast HTTP servers and clients** with built-in streaming support
- **WebSocket support** for real-time two-way communication

## 1.1 Completeness

🚧 **IMPLEMENTATION STATUS**: The Osprey language and compiler are not complete. The documentation here does not represent the language at this present time. This specification represents the design aims for the language, a description of the syntax and a description of where the roadmap is taking the language. Developers should pay attention to the spec first are foremost as the single source of truth in regards to the syntax.

## 1.2 Principles

- Elegance (simplicity, ergonomics, efficiency), safety (fewer footguns, security at every level), performance (uses the most efficient approach and allows the use of Rust interop for extreme performance)
- No more than 1 way to do anything
- Make illegal states unrepresentable. There are no exceptions or panics. Anything than can result in an error state returns a result object
- Referential transparency
- Simplicity
- Interopability with Rust for high performance workloads
- Interopability with Haskell (future) for fundamental correctness
- Static/strong typing. Nothing should be "any" unless EXPLICITLY declared as any
- Minimal ceremony. No main function necessary for example.
- **Fast HTTP performance** as a core design principle
- **Streaming by default** for large responses to prevent memory issues

## 2. Lexical Structure

### 2.1 Identifiers

Identifiers must start with a letter or underscore, followed by letters, digits, or underscores:

```
ID := [a-zA-Z_][a-zA-Z0-9_]*
```

### 2.2 Keywords

Reserved keywords:
```
fn let mut type import match if then else case of extern
```

### 2.3 Literals

#### Integer Literals
```
INTEGER := [0-9]+
```

#### String Literals
```
STRING := '"' (CHAR | ESCAPE_SEQUENCE)* '"'
ESCAPE_SEQUENCE := '\n' | '\t' | '\r' | '\\' | '\"'
```

#### Interpolated String Literals
```
INTERPOLATED_STRING := '"' (CHAR | INTERPOLATION)* '"'
INTERPOLATION := '${' EXPRESSION '}'
```

### 2.4 Operators

#### Arithmetic Operators (All Safe by Default)
- `+` Addition: `(Int, Int) -> Result<Int, MathError>`
- `-` Subtraction: `(Int, Int) -> Result<Int, MathError>`
- `*` Multiplication: `(Int, Int) -> Result<Int, MathError>`
- `/` Division: `(Int, Int) -> Result<Int, MathError>`

**CRITICAL**: All arithmetic operators return Result types to handle overflow, underflow, and division by zero.

#### Assignment Operator
- `=` Assignment

#### Other Operators
- `=>` Lambda/Match arm arrow
- `|` Union type separator
- `::` Type annotation

### 2.5 Delimiters

- `(` `)` Parentheses
- `{` `}` Braces
- `[` `]` Brackets
- `,` Comma
- `:` Colon
- `;` Semicolon
- `.` Dot

## 3. Syntax

### 3.1 Program Structure

A Osprey program consists of a sequence of top-level statements and modules.

```
program := statement* EOF

statement := importStmt
          | letDecl
          | fnDecl
          | externDecl
          | typeDecl
          | moduleDecl
          | exprStmt
```

### 3.2 Import Statements

🚧 **PARTIAL IMPLEMENTATION**: Basic import parsing is implemented but module resolution is limited.

```
importStmt := IMPORT ID (DOT ID)*
```

**Examples:**
```osprey
import std
import std.io
import graphics.canvas
```

### 3.3 Let Declarations

```
letDecl := (LET | MUT) ID (COLON type)? EQ expr
```

**Examples:**
```osprey
let x = 42
let name = "Alice"
mut counter = 0
let result = calculateValue(input: data)
```

### 3.4 Function Declarations

```
fnDecl := docComment? FN ID LPAREN paramList? RPAREN (ARROW type)? (EQ expr | LBRACE blockBody RBRACE)

paramList := param (COMMA param)*
param := ID (COLON type)?
```

**Examples:**
```osprey
fn double(x) = x * 2
fn add(x, y) = x + y
fn greet(name) = "Hello " + name
fn getValue() = 42
```

### 3.5 Extern Declarations

Extern declarations allow Osprey programs to call functions implemented in other languages (such as Rust, C, or C++). These declarations specify the interface for external functions without providing an implementation.

```
externDecl := docComment? EXTERN FN ID LPAREN externParamList? RPAREN (ARROW type)?

externParamList := externParam (COMMA externParam)*
externParam := ID COLON type
```

**Key Features:**
- **Required type annotations**: All parameters must have explicit type annotations
- **Optional return type**: Return type can be specified with `-> type` syntax
- **No function body**: Extern declarations only specify the interface
- **Interoperability**: Enables calling functions from Rust, C, and other languages

**Examples:**
```osprey
// Basic extern function declarations
extern fn rust_add(a: int, b: int) -> int
extern fn rust_multiply(a: int, b: int) -> int
extern fn rust_factorial(n: int) -> int

// Using extern functions with named arguments
let sum = rust_add(a: 15, b: 25)
let product = rust_multiply(a: 6, b: 7)
let factorial = rust_factorial(5)

// Extern functions with different return types
extern fn rust_is_prime(n: int) -> bool
extern fn rust_format_number(n: int) -> string

let isPrime = rust_is_prime(17)
let formatted = rust_format_number(42)
```

**Type Mapping:**
- Osprey `int` ↔ Rust `i64` ↔ C `int64_t`
- Osprey `bool` ↔ Rust `bool` ↔ C `bool`
- Osprey `string` ↔ Rust `*const c_char` ↔ C `char*`

**Implementation Requirements:**
- External functions must use C ABI (`extern "C"` in Rust)
- Functions must be marked with `#[no_mangle]` in Rust
- Static libraries must be linked during compilation

### 3.6 Type Declarations

```
typeDecl := docComment? TYPE ID (LT typeParamList GT)? EQ (unionType | recordType)

unionType := variant (BAR variant)*
recordType := LBRACE fieldDeclarations RBRACE

variant := ID (LBRACE fieldDeclarations RBRACE)?

fieldDeclarations := fieldDeclaration (COMMA fieldDeclaration)*
fieldDeclaration := ID COLON type constraint?

constraint := WHERE functionCall
```

**Examples:**
```osprey
type Color = Red | Green | Blue

type Shape = Circle { radius: Int } 
           | Rectangle { width: Int, height: Int }

type Result = Success { value: String } 
            | Error { message: String }
```

### 3.7 Record Types and Type Constructors

❌ **NOT FULLY IMPLEMENTED**: Record types with constraints are a design goal but not yet implemented. Basic type declarations work but constraint validation is not implemented.

#### 3.7.1 Record Type Declaration

Record types define structured data with named fields:

```
record_type := 'type' ID '=' '{' field_declarations '}'

field_declarations := field_declaration (',' field_declaration)*
field_declaration := ID ':' type constraint?
constraint := 'where' function_call
```

**Examples:**
```osprey
// Simple record type
type Point = { x: Int, y: Int }

// Record with constraints
type Person = { 
    name: String where isValidName(name),
    age: Int where between(age, 0, 150),
    email: String where validateEmail(email),
    confirmEmail: String where equals(email, confirmEmail)  // Cross-field
}

// Mixed record and union types
type User = Guest { sessionId: String }
          | Member { 
              id: Int where id > 0,
              name: String where isValidName(name),
              joinDate: String
            }
```

#### 3.7.2 Type Construction

Type constructors create instances of record types using elegant field syntax. If there are any where constraints on any fields, the constructor returns the same Result type as arithmetic operators

```
type_constructor := type_name '{' field_assignments '}'
field_assignments := field_assignment (',' field_assignment)*
field_assignment := ID ':' expression
```

**Construction Examples:**
```osprey
// Simple construction
let point = Point { x: 10, y: 20 }

// Construction with constraints (returns Result)
let person = Person { 
    name: "Alice", 
    age: 25, 
    email: "alice@example.com" 
}

// Handle construction results
match person {
    Ok { value } => print("Created person: ${value.name}")
    Err { error } => print("Construction failed: ${error}")
}
```

#### 3.7.3 Construction Result Types

**CRITICAL**: Type constructors with constraints return `Result<T, ConstructionError>`:

- **Unconstrained types**: Direct construction returns the type
- **Constrained types**: Construction returns `Result<T, ConstructionError>`

**Rule**: If ANY field in a type has a WHERE constraint, the constructor ALWAYS returns a Result type.

```osprey
// No constraints = direct construction
type Point = { x: Int, y: Int }
let point = Point { x: 10, y: 20 }  // Returns Point

// With constraints = Result construction  
type Person = { 
    name: String where notEmpty(name),
    age: Int where isPositive(age) 
}
let person = Person { name: "Alice", age: 25 }  // Returns Result<Person, ConstructionError>

// Even ONE constraint makes it return Result
type User = {
    id: Int,                                    // No constraint
    name: String where notEmpty(name),          // Has constraint
    email: String                               // No constraint  
}
let user = User { id: 1, name: "Alice", email: "alice@test.com" }  // Returns Result<User, ConstructionError>

// Construction error handling
match Person { name: "", age: 25 } {
    Ok { value } => useValidPerson(value)
    Err { error } => match error {
        ConstraintViolation { field, message } => 
            print("Field ${field} failed: ${message}")
        MissingField { field } => 
            print("Required field missing: ${field}")
    }
}
```

#### 3.7.4 Non-Destructive Mutation (Structural Updates)

Records support elegant non-destructive updates that create modified copies:

```osprey
// Original record
let person = Person { name: "Alice", age: 25, email: "alice@example.com" }

// Non-destructive update (creates new instance)
let olderPerson = person { age: 26 }           // Only age changes
let renamedPerson = person { name: "Alicia" }  // Only name changes

// Multiple field updates
let updatedPerson = person { 
    age: 26, 
    email: "alicia@newdomain.com" 
}

// Original person unchanged - all updates create new instances
print(person.age)        // Still 25
print(olderPerson.age)   // Now 26
```

#### 3.7.5 Update Result Types

Updates that involve constrained fields also return Results:

```osprey
// Update with constraint validation
let result = person { age: 200 }  // Returns Result<Person, ConstraintViolation>

match result {
    Ok { value } => useUpdatedPerson(value)
    Err { error } => handleConstraintError(error)
}

// Valid update
let validUpdate = person { age: 30 }  // Returns Ok<Person>
```

#### 3.7.6 Field Access

Record fields are accessed using dot notation:

```osprey
let person = Person { name: "Alice", age: 25, email: "alice@example.com" }

print("Name: ${person.name}")     // "Alice"
print("Age: ${person.age}")       // 25
print("Email: ${person.email}")   // "alice@example.com"
```

#### 3.7.7 Pattern Matching on Records

Records can be destructured in pattern matching:

```osprey
match person {
    Person { name, age: 25, email } => 
        print("25-year-old ${name} with email ${email}")
    Person { name, age, email } => 
        print("${name} is ${age} years old")
}

// Partial destructuring
match person {
    Person { name: "Alice", ... } => print("It's Alice!")
    Person { age, ... } if age < 18 => print("Minor")
    _ => print("Other person")
}
```

#### 3.7.8 Constraint Functions

Constraints are function calls that return boolean values. The constraint system supports both compile-time and runtime evaluation:

- **Compile-time constraints**: When all arguments are constants/literals, functions execute at compile time
- **Runtime constraints**: When any argument is a runtime value, functions execute during construction

**Constraint Syntax:**
```
constraint := 'where' function_call
function_call := ID '(' argument_list ')'
```

**Constraint Categories:**
- **Field validation**: Direct field value checking
- **Cross-field validation**: Constraints involving multiple fields  
- **Complex validation**: Custom validation functions
- **Built-in constraints**: Standard validation functions

**Examples:**
```osprey
type Person = {
    name: String where isValidName(name),
    age: Int where between(age, 0, 150),
    email: String where validateEmail(email),
    confirmEmail: String where equals(email, confirmEmail)  // Cross-field
}

type Rectangle = {
    width: Int where positive(width),
    height: Int where positive(height),
    area: Int where equals(area, multiply(width, height))  // Computed constraint
}

type CreditCard = {
    number: String where luhnCheck(number),
    expiryMonth: Int where between(expiryMonth, 1, 12),
    expiryYear: Int where greaterThan(expiryYear, currentYear()),
    cvv: String where matchesPattern(cvv, "^[0-9]{3,4}$")
}
```

**Compile-Time vs Runtime Evaluation:**

```osprey
// All constraints evaluated at COMPILE TIME (constants/literals)
let person1 = Person { 
    name: "Alice",           // isValidName("Alice") → compile time
    age: 25,                 // between(25, 0, 150) → compile time  
    email: "alice@test.com"  // validateEmail("alice@test.com") → compile time
}

// Mixed compile-time and runtime evaluation
let inputName = input()
let person2 = Person {
    name: inputName,         // isValidName(inputName) → RUNTIME
    age: 30,                 // between(30, 0, 150) → compile time
    email: "bob@test.com"    // validateEmail("bob@test.com") → compile time
}

// All constraints evaluated at RUNTIME
let inputAge = input()
let inputEmail = input()
let person3 = Person {
    name: inputName,         // isValidName(inputName) → runtime
    age: inputAge,           // between(inputAge, 0, 150) → runtime
    email: inputEmail        // validateEmail(inputEmail) → runtime
}
```

**Custom Constraint Functions:**

```osprey
// Basic validation functions using match expressions
fn notEmpty(s: String) -> Bool = match s {
    "" => false
    _ => true
}

fn isPositive(n: Int) -> Bool = match n {
    0 => false  
    _ => true
}

// Complex validation with multiple rules
fn validateUsername(username: String) -> Bool = match username {
    "" => false           // Empty
    " " => false          // Whitespace only
    "admin" => false      // Reserved word
    "root" => false       // Reserved word
    "a" => false          // Too short
    _ => true             // Everything else valid
}

// Numeric range and reserved value validation
fn validatePort(port: Int) -> Bool = match port {
    0 => false           // Invalid port
    1 => false           // Reserved
    22 => false          // SSH reserved
    80 => true           // HTTP valid
    443 => true          // HTTPS valid
    65536 => false       // Too high
    _ => true            // Most ports valid
}

// Complex password validation
fn isValidPassword(password: String) -> Bool = 
    length(password) >= 8 && 
    hasUppercase(password) && 
    hasLowercase(password) && 
    hasDigits(password)

fn isBusinessHour(hour: Int) -> Bool = 
    between(hour, 9, 17)

fn isWeekend(dayOfWeek: String) -> Bool = 
    equals(dayOfWeek, "Saturday") || equals(dayOfWeek, "Sunday")

// Use in type definitions
type UserAccount = {
    username: String where validateUsername(username),
    password: String where isValidPassword(password),
    loginHour: Int where isBusinessHour(loginHour)
}

type NetworkConfig = {
    port: Int where validatePort(port),
    host: String where notEmpty(host)
}

type Appointment = {
    dayOfWeek: String where not(isWeekend(dayOfWeek)),
    hour: Int where isBusinessHour(hour),
    duration: Int where between(duration, 15, 240)  // 15 min to 4 hours
}
```

**Constraint Function Requirements:**
- Must return `Bool` type
- Can reference the field being constrained
- Can reference other fields in the same record
- Can call other functions (including built-ins)
- Must be pure functions (no side effects)

**Performance Characteristics:**
- **Compile-time evaluation**: Zero runtime cost for constant expressions
- **Runtime evaluation**: Minimal overhead, constraints evaluated once during construction
- **Caching**: Constraint results cached to avoid re-evaluation
- **Short-circuiting**: Constraint evaluation stops at first failure

#### 3.7.9 Built-in Constraint Functions

Standard constraint functions available for validation:

```osprey
// String validation functions
length(s: String) -> Int                    // String length
isEmpty(s: String) -> Bool                  // Empty string check
notEmpty(s: String) -> Bool                 // Non-empty string check
contains(s: String, sub: String) -> Bool    // Substring check
startsWith(s: String, prefix: String) -> Bool
endsWith(s: String, suffix: String) -> Bool
matchesPattern(s: String, regex: String) -> Bool  // Regex validation

// String content validation
isAlphabetic(s: String) -> Bool            // Only letters
isNumeric(s: String) -> Bool               // Only digits
isAlphanumeric(s: String) -> Bool          // Letters and digits
hasUppercase(s: String) -> Bool            // Contains uppercase
hasLowercase(s: String) -> Bool            // Contains lowercase  
hasDigits(s: String) -> Bool               // Contains digits

// Numeric validation functions
between(n: Int, min: Int, max: Int) -> Bool  // Range check: min <= n <= max
positive(n: Int) -> Bool                     // n > 0
negative(n: Int) -> Bool                     // n < 0
nonNegative(n: Int) -> Bool                  // n >= 0
nonPositive(n: Int) -> Bool                  // n <= 0
greaterThan(a: Int, b: Int) -> Bool          // a > b
lessThan(a: Int, b: Int) -> Bool             // a < b
equals(a: Any, b: Any) -> Bool               // a == b

// Logical functions
not(condition: Bool) -> Bool                 // Logical negation
and(a: Bool, b: Bool) -> Bool               // Logical AND
or(a: Bool, b: Bool) -> Bool                // Logical OR

// Mathematical functions (for constraints)
multiply(a: Int, b: Int) -> Int             // a * b
add(a: Int, b: Int) -> Int                  // a + b
subtract(a: Int, b: Int) -> Int             // a - b
divide(a: Int, b: Int) -> Int               // a / b (constraint context)

// Specialized validation functions
luhnCheck(cardNumber: String) -> Bool        // Credit card validation
isValidEmail(email: String) -> Bool          // Email format validation
isValidURL(url: String) -> Bool              // URL format validation
isValidIPAddress(ip: String) -> Bool         // IP address validation

// Temporal functions (for date/time constraints)
currentYear() -> Int                         // Current year (compile-time constant)
currentMonth() -> Int                        // Current month
currentDay() -> Int                          // Current day
```

**Function-Based Constraint Examples:**

```osprey
type StrongPassword = {
    value: String where and(
        between(length(value), 8, 128),
        and(hasUppercase(value), 
            and(hasLowercase(value), hasDigits(value))
        )
    )
}

type ValidatedEmail = {
    address: String where and(
        notEmpty(address),
        and(contains(address, "@"), isValidEmail(address))
    )
}

type DateOfBirth = {
    year: Int where between(year, 1900, currentYear()),
    month: Int where between(month, 1, 12),
    day: Int where between(day, 1, 31)
}

type BankAccount = {
    accountNumber: String where and(
        equals(length(accountNumber), 10),
        isNumeric(accountNumber)
    ),
    routingNumber: String where luhnCheck(routingNumber),
    balance: Int where nonNegative(balance)
}
```

#### 3.7.10 Error Types for Construction

```osprey
type ConstructionError = 
    ConstraintViolation { 
        field: String, 
        value: String, 
        constraint: String,
        message: String 
    }
  | MissingField { field: String }
  | TypeMismatch { 
        field: String, 
        expected: String, 
        actual: String 
    }
  | ConstraintFunctionError {
        field: String,
        function: String,
        error: String
    }
  | MultipleConstraintViolations {
        violations: String  // List of all constraint failures
    }
```

#### 3.7.11 Compilation Errors for Field Access

**CRITICAL**: Attempting to access fields directly on constrained type constructor results must produce specific compilation errors.

**Field Access on Result Types:**
When a type has WHERE constraints, its constructor returns `Result<T, ConstructionError>`. Attempting to access fields directly on this Result type should produce a clear compilation error:

```osprey
type User = { 
    name: String where notEmpty(name)
}

let user = User { name: "alice" }  // Returns Result<User, ConstructionError>

// COMPILATION ERROR: Cannot access field on Result type
print("${user.name}")  
// Should produce: "cannot access field 'name' on Result<User, ConstructionError> type - pattern matching required"

let name = user.name
// Should produce: "field access requires pattern matching on Result type"
```

**Required Error Messages:**
- **Field access on Result**: `"cannot access field 'FIELD' on Result<TYPE, ConstructionError> type - pattern matching required"`
- **Assignment from Result field**: `"field access requires pattern matching on Result type"`
- **Missing pattern matching**: `"constrained types return Result - use match expression to handle success/failure"`

**Correct Pattern:**
```osprey
match user {
    Ok { value } => print("Name: ${value.name}")
    Err { error } => print("Construction failed: ${error}")
}
```

**Current Implementation Issue:**
The current compiler incorrectly reports field access attempts as "undefined variable" errors instead of proper Result type access errors. This should be fixed to provide clear guidance on Result type handling.

### 3.8 Expressions

#### Binary Expressions
```
binary_expression := multiplicative_expression (('+' | '-') multiplicative_expression)*

multiplicative_expression := unary_expression (('*' | '/') unary_expression)*
```

#### Unary Expressions
```
unary_expression := ('+' | '-')? pipe_expression
```

#### Function Calls
```
call_expression := primary ('.' ID '(' argument_list? ')')* 
                | primary ('(' argument_list? ')')?

argument_list := named_argument_list 
              | positional_argument_list

named_argument_list := named_argument (',' named_argument)+
named_argument := ID ':' expression

positional_argument_list := expression (',' expression)*
```

#### Primary Expressions
```
primary_expression := literal
                   | identifier
                   | '(' expression ')'
                   | lambda_expression
                   | block_expression
```

### 3.9 Block Expressions

Block expressions allow grouping multiple statements together and returning a value from the final expression. They create a new scope for variable declarations and enable sequential execution with proper scoping rules.

```
block_expression := '{' statement* expression? '}'
```

**Examples:**
```osprey
// Simple block with local variables
let result = {
    let x = 10
    let y = 20
    x + y
}
print("Result: ${result}")  // prints "Result: 30"

// Nested blocks
let complex = {
    let outer = 100
    let inner_result = {
        let inner = 50
        outer + inner
    }
    inner_result * 2
}
print("Complex: ${complex}")  // prints "Complex: 300"

// Block with function calls
fn multiply(a: int, b: int) -> int = a * b
let calc = {
    let a = 5
    let b = 6
    multiply(a: a, b: b)
}
print("Calculation: ${calc}")  // prints "Calculation: 30"
```

#### 3.9.1 Block Scoping Rules

Block expressions create a new lexical scope:
- Variables declared inside a block are only visible within that block
- Variables from outer scopes can be accessed (lexical scoping)
- Variables declared in a block shadow outer variables with the same name
- Variables go out of scope when the block ends

**Scoping Examples:**
```osprey
let x = 100
let result = {
    let x = 50        // Shadows outer x
    let y = 25        // Only visible in this block
    x + y             // Uses inner x (50)
}
print("Result: ${result}")  // 75
print("Outer x: ${x}")      // 100 (unchanged)
// print("${y}")            // ERROR: y not in scope
```

#### 3.9.2 Block Return Values

Block expressions return the value of their final expression:
- If the block ends with an expression, that value is returned
- If the block has no final expression, it returns the unit type
- The block's type is determined by the type of the final expression

**Return Value Examples:**
```osprey
// Block returns integer
let number = {
    let a = 10
    let b = 20
    a + b           // Returns 30
}

// Block returns string
let message = {
    let name = "Alice"
    let age = 25
    "Hello ${name}, age ${age}"  // Returns string
}

// Block with statements only (returns unit)
let side_effect = {
    print("Doing work...")
    print("Work complete")
    // No final expression - returns unit
}
```

#### 3.9.3 Block Expressions in Match Arms

Block expressions are particularly useful in match expressions for complex logic:

```osprey
let result = match status {
    Success => {
        print("Operation succeeded")
        let timestamp = getCurrentTime()
        "Success at ${timestamp}"
    }
    Error => {
        print("Operation failed")
        let error_code = getErrorCode()
        "Error ${error_code}"
    }
    _ => "Unknown status"
}
```

#### 3.9.4 Function Bodies as Blocks

Functions can use block expressions as their body instead of single expressions:

```osprey
fn processData(input: string) -> string = {
    let cleaned = cleanInput(input)
    let validated = validateInput(cleaned)
    let processed = transformData(validated)
    formatOutput(processed)
}

// Equivalent to expression-bodied function:
fn processData(input: string) -> string = 
    formatOutput(transformData(validateInput(cleanInput(input))))
```

#### 3.9.5 Type Safety and Inference

Block expressions follow Osprey's type safety rules:
- The block's type is inferred from the final expression
- All statements in the block must be well-typed
- Variable declarations in blocks follow the same type inference rules
- Return type must be compatible with the expected type

**Type Inference Examples:**
```osprey
// Block type inferred as Int
let num: int = {
    let a = 10
    let b = 20
    a + b              // Type: int
}

// Block type inferred as String
let text: string = {
    let name = "Bob"
    "Hello ${name}"    // Type: string
}

// ERROR: Type mismatch
let wrong: int = {
    let x = 10
    "not a number"     // ERROR: Expected int, got string
}
```

#### 3.9.6 Performance Characteristics

Block expressions are zero-cost abstractions:
- **Compile-time scoping**: All variable scoping resolved at compile time
- **No runtime overhead**: Blocks compile to sequential instructions
- **Stack allocation**: Local variables allocated on the stack
- **Optimized away**: Simple blocks with no local variables are optimized away

#### 3.9.7 Best Practices

**Use block expressions when:**
- You need local variables for complex calculations
- Breaking down complex expressions into readable steps
- Implementing complex match arm logic
- Creating temporary scopes to avoid variable name conflicts

**Avoid block expressions when:**
- A simple expression would suffice
- The block only contains a single expression
- Creating unnecessary nesting levels

**Good Examples:**
```osprey
// Good: Complex calculation with intermediate steps
let result = {
    let base = getUserInput()
    let squared = base * base
    let doubled = squared * 2
    squared + doubled
}

// Good: Complex match logic
let response = match request.method {
    POST => {
        let body = parseBody(request.body)
        let validated = validateData(body)
        processCreation(validated)
    }
    _ => "Method not allowed"
}
```

**Bad Examples:**
```osprey
// Bad: Unnecessary block for simple expression
let bad = {
    42
}
// Better: let bad = 42

// Bad: Single operation doesn't need block
let also_bad = {
    x + y
}
// Better: let also_bad = x + y
```

### 3.10 Match Expressions

```
matchExpr := MATCH expr LBRACE matchArm+ RBRACE

matchArm := pattern LAMBDA expr

pattern := unaryExpr                                   // Support negative numbers: -1, +42, etc.
        | ID (LBRACE fieldPattern RBRACE)?          // Pattern destructuring: Ok { value }
        | ID (LPAREN pattern (COMMA pattern)* RPAREN)?  // Constructor patterns
        | ID (ID)?                                   // Variable capture
        | ID COLON type                              // Type annotation pattern: value: Int
        | ID COLON LBRACE fieldPattern RBRACE       // Named structural: person: { name, age }
        | LBRACE fieldPattern RBRACE                // Anonymous structural: { name, age }
        | UNDERSCORE                                 // Wildcard

fieldPattern := ID (COMMA ID)*
```

**Example:**
```osprey
let result = match status {
    Success => "OK"
    Error msg => "Failed: " + msg
    _ => "Unknown"
}
```

## 4. Semantics

### 4.1 Variable Binding

- `let` creates immutable bindings
- `mut` creates mutable bindings
- Variables must be initialized at declaration
- Shadowing is allowed in nested scopes

### 4.2 Function Semantics

- Functions are first-class values
- All functions are pure (no side effects except I/O)
- Recursive functions are supported
- Tail recursion is optimized

### 4.3 Evaluation Order

- Expressions are evaluated left-to-right
- Function arguments are evaluated before the function call
- Short-circuit evaluation for logical operators

## 5. Type System

### 5.1 Built-in Types

- `Int`: 64-bit signed integers
- `String`: UTF-8 encoded strings  
- `Bool`: Boolean values (`true`, `false`)
- `Unit`: Type for functions that don't return a meaningful value
- `Result<T, E>`: Built-in generic type for error handling

### 5.2 Built-in Error Types

- `MathError`: For arithmetic operations (DivisionByZero, Overflow, Underflow)
- `ParseError`: For string parsing operations
- `IndexError`: For array/string indexing operations

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

## 6. Function Calls

### 6.1 Named Arguments Requirement

**CRITICAL RULE**: Functions with more than one parameter **MUST** be called with named arguments.

#### Valid Function Calls

```osprey
// Single parameter - positional allowed
fn double(x) = x * 2
let result = double(5)

// Zero parameters - no arguments
fn getValue() = 42
let value = getValue()

// Multiple parameters - named arguments required
fn add(x, y) = x + y
let sum = add(x: 10, y: 20)

// Multiple parameters - order doesn't matter with named args
let sum2 = add(y: 20, x: 10)
```

#### Invalid Function Calls

```osprey
// ERROR: Multi-parameter function with positional arguments
fn add(x, y) = x + y
let sum = add(10, 20)  // ❌ Compilation error

// ERROR: Mixed positional and named arguments
let sum = add(10, y: 20)  // ❌ Compilation error
```

### 6.2 Function Call Compilation Rules

1. **Single Parameter Functions**: May use positional arguments
2. **Zero Parameter Functions**: Called with empty parentheses `()`
3. **Multi-Parameter Functions**: Must use named arguments for ALL parameters
4. **Argument Reordering**: Named arguments are reordered to match parameter declaration order

## 7. String Interpolation

✅ **IMPLEMENTED**: String interpolation is fully implemented and working with comprehensive test coverage.

### 7.1 Syntax

String interpolation uses `${}` syntax within double-quoted strings:

```osprey
let name = "Alice"
let age = 30
let message = "Hello ${name}, you are ${age} years old"
```

### 7.2 Expression Support

Any expression is valid inside interpolation:

```osprey
let x = 10
let y = 5
print("Sum: ${x + y}")
print("Product: ${x * y}")
print("Complex: ${(x + y) * 2 - 1}")
```

### 7.3 Type Handling

- **String variables**: Use `%s` format specifier
- **Integer expressions**: Use `%ld` format specifier  
- **Function calls**: Supported for single-parameter functions

### 7.4 Implementation

Interpolated strings are compiled to:
1. Allocate a buffer (`alloca [1024 x i8]`)
2. Use `sprintf` to format the string
3. Use `puts` to output the result

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

## 9. Block Expressions

## 10. Boolean Operations

## 11. Loop Constructs and Functional Iterators

🚧 **PARTIAL IMPLEMENTATION**: Basic iterator functions (`range`, `forEach`, `map`, `filter`, `fold`) are implemented and working. The pipe operator (`|>`) is implemented.

## 12. Lightweight Fibers and Concurrency

🚧 **IMPLEMENTATION STATUS**: Fiber syntax is partially implemented. Basic fiber operations (`spawn`, `await`, `yield`) are in the grammar but runtime support is limited.

### 12.1 Fiber Types and Concurrency

Osprey provides lightweight concurrency through fiber types. Unlike traditional function-based approaches, fibers are proper type instances constructed using Osprey's standard type construction syntax.

#### Core Fiber Types

**`Fiber<T>`** - A lightweight concurrent computation that produces a value of type T
**`Channel<T>`** - A communication channel for passing values of type T between fibers

#### Fiber Construction

Fibers are created using standard type construction syntax:

```osprey
// Create a fiber that computes a value
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(n: 1000) 
}

// Create a fiber with more complex computation
let worker = Fiber<String> { 
    computation: fn() => {
        processData()
        "completed"
    }
}

// Create a parameterized fiber
let calculator = Fiber<Int> { 
    computation: fn() => multiply(x: 10, y: 20) 
}
```

#### Spawn Syntax Sugar

For convenience, Osprey provides `spawn` as syntax sugar for creating and immediately starting a fiber:

```osprey
// Using spawn (syntax sugar)
let result = spawn 42

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => 42 }
let result = fiber

// More complex spawn
let computation = spawn (x * 2 + y)

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => x * 2 + y }
let computation = fiber
```

The `spawn` keyword immediately evaluates the expression in a new fiber context, making it convenient for quick concurrent computations without the full type construction syntax.

#### Channel Construction

Channels are created using type construction syntax:

```osprey
// Unbuffered (synchronous) channel
let sync_channel = Channel<Int> { capacity: 0 }

// Buffered (asynchronous) channel  
let async_channel = Channel<String> { capacity: 10 }

// Large buffer channel
let buffer_channel = Channel<Int> { capacity: 100 }
```

#### Fiber Operations

Once created, fibers and channels are manipulated using functional operations:

**`await(fiber: Fiber<T>) -> T`** - Wait for fiber completion and get result
**`send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`** - Send value to channel
**`recv(channel: Channel<T>) -> Result<T, ChannelError>`** - Receive value from channel
**`yield() -> Unit`** - Voluntarily yield control to scheduler

```osprey
// Create and await a fiber
let task = Fiber<Int> { computation: fn() => heavyComputation() }
let result = await(task)

// Channel communication
let ch = Channel<String> { capacity: 5 }
send(ch, "hello")
let message = recv(ch)

// Yielding control
yield()
```

#### Complete Fiber Example

```osprey
// Producer fiber
let producer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        send(ch, 1)
        send(ch, 2) 
        send(ch, 3)
    }
}

// Consumer fiber
let consumer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        let value1 = recv(ch)
        let value2 = recv(ch)
        let value3 = recv(ch)
        print("Received: ${value1}, ${value2}, ${value3}")
    }
}

// Start both fibers
await(producer)
await(consumer)
```

#### Select Expression for Channel Multiplexing

The `select` expression allows waiting on multiple channel operations:

```osprey
let ch1 = Channel<String> { capacity: 1 }
let ch2 = Channel<Int> { capacity: 1 }

let result = select {
    msg => recv(ch1) => process_string(msg)
    num => recv(ch2) => process_number(num)
    _ => timeout_handler()
}
```

#### Rust Interoperability

Osprey fibers are designed to interoperate with Rust's async/await system:

```osprey
// Osprey fiber that calls Rust async function
extern fn rust_async_task() -> Future<Int>

let osprey_task = Fiber<Int> {
    computation: fn() => await(rust_async_task())
}

let result = await(osprey_task)
```

#### Fiber-Isolated Modules

Each fiber gets its own isolated instance of modules, preventing data races:

```osprey
module Counter {
    let mut count = 0
    fn increment() -> Int = { count = count + 1; count }
    fn get() -> Int = count
}

// Each fiber has its own Counter instance
let fiber1 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}
let fiber2 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}

// These will both return 1, not 1 and 2
let result1 = await(fiber1)  // 1
let result2 = await(fiber2)  // 1 (separate instance)
```

## 13. Fiber-Isolated Module System

❌ **NOT IMPLEMENTED**: The fiber-isolated module system is a design goal but not yet implemented. Current module support is basic.

## 14. Built-in Functions

🚀 **IMPLEMENTATION STATUS**: HTTP functions are implemented and working. WebSocket functions are implemented but undergoing testing. Fiber operations are partially implemented.

Osprey provides first-class support for HTTP servers and clients, designed with performance, safety, and streaming as core principles. All HTTP functions follow Osprey's functional programming paradigms and comply with:

- **RFC 7230**: HTTP/1.1 Message Syntax and Routing ([https://tools.ietf.org/html/rfc7230](https://tools.ietf.org/html/rfc7230))
- **RFC 7231**: HTTP/1.1 Semantics and Content ([https://tools.ietf.org/html/rfc7231](https://tools.ietf.org/html/rfc7231))
- **RFC 7232**: HTTP/1.1 Conditional Requests ([https://tools.ietf.org/html/rfc7232](https://tools.ietf.org/html/rfc7232))
- **RFC 7233**: HTTP/1.1 Range Requests ([https://tools.ietf.org/html/rfc7233](https://tools.ietf.org/html/rfc7233))
- **RFC 7234**: HTTP/1.1 Caching ([https://tools.ietf.org/html/rfc7234](https://tools.ietf.org/html/rfc7234))
- **RFC 7235**: HTTP/1.1 Authentication ([https://tools.ietf.org/html/rfc7235](https://tools.ietf.org/html/rfc7235))

- **Result types** instead of exceptions for error handling
- **Immutable response objects** that cannot be modified after creation
- **Streaming by default** for large response bodies to prevent memory issues
- **Fiber-based concurrency** for handling thousands of concurrent connections

### 14.1 HTTP Core Types

#### HTTP Method Union Type
```osprey
type HttpMethod = GET | POST | PUT | DELETE | PATCH | HEAD | OPTIONS
```

#### HTTP Request Type (Immutable)
```osprey
type HttpRequest = {
    method: HttpMethod,
    path: String,
    headers: String,
    body: String,
    queryParams: String
}
```

#### HTTP Response Type (Immutable with Streaming)
```osprey
type HttpResponse = {
    status: Int,
    headers: String,
    contentType: String,
    contentLength: Int,
    streamFd: Int,        // File descriptor for streaming
    isComplete: Bool,     // Whether response is fully loaded
    partialBody: String,  // Current chunk of body data
    partialLength: Int    // Length of current chunk
}
```

### 14.2 HTTP Server Functions

#### `httpCreateServer(port: Int, address: String) -> Result<ServerID, String>`

Creates an HTTP server bound to the specified port and address.

**Parameters:**
- `port`: Port number (1-65535)
- `address`: IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")

**Returns:**
- `Success(serverID)`: Unique server identifier
- `Err(message)`: Error description (invalid port, bind failure, etc.)

**Example:**
```osprey
let serverResult = httpCreateServer(port: 8080, address: "127.0.0.1")
match serverResult {
    Success serverId => print("Server created with ID: ${serverId}")
    Err message => print("Failed to create server: ${message}")
}
```

#### `httpListen(serverID: Int, handler: fn(HttpRequest) -> Result<HttpResponse, String>) -> Result<Success, String>`

Starts the HTTP server listening for requests. Each request is handled in a separate fiber for maximum concurrency.

**Parameters:**
- `serverID`: Server identifier from `httpCreateServer`
- `handler`: Request handler function that processes incoming requests

**Returns:**
- `Success()`: Server started successfully
- `Err(message)`: Error description

**Example:**
```osprey
fn handleRequest(request: HttpRequest) -> Result<HttpResponse, String> = match request.method {
    GET => match request.path {
        "/health" => Success(HttpResponse {
            status: 200,
            contentType: "application/json",
            partialBody: "{\"status\": \"healthy\"}",
            isComplete: true,
            streamFd: -1
        })
        "/users" => Success(HttpResponse {
            status: 200,
            contentType: "application/json", 
            partialBody: "[{\"id\": 1, \"name\": \"Alice\"}]",
            isComplete: true,
            streamFd: -1
        })
        _ => Success(HttpResponse {
            status: 404,
            contentType: "text/plain",
            partialBody: "Not Found",
            isComplete: true,
            streamFd: -1
        })
    }
    POST => handlePostRequest(request)
    _ => Err("Method not supported")
}

let listenResult = httpListen(serverID: serverId, handler: handleRequest)
```

#### `httpStopServer(serverID: Int) -> Result<Success, String>`

Stops the HTTP server and cleans up resources.

**Parameters:**
- `serverID`: Server identifier to stop

**Returns:**
- `Success()`: Server stopped successfully  
- `Err(message)`: Error description

### 14.3 HTTP Client Functions

#### `httpCreateClient(baseUrl: String, timeout: Int) -> Result<ClientID, String>`

Creates an HTTP client for making requests.

**Parameters:**
- `baseUrl`: Base URL for requests (e.g., "http://api.example.com")
- `timeout`: Request timeout in milliseconds

**Returns:**
- `Success(clientID)`: Unique client identifier
- `Err(message)`: Error description

**Example:**
```osprey
let clientResult = httpCreateClient(baseUrl: "http://jsonplaceholder.typicode.com", timeout: 5000)
```

#### `httpGet(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP GET request.

**Parameters:**
- `clientID`: Client identifier from `httpCreateClient`
- `path`: Request path (e.g., "/users/1")
- `headers`: Additional headers (e.g., "Authorization: Bearer token\r\n")

**Returns:**
- `Success(statusCode)`: HTTP status code (200, 404, etc.)
- `Err(message)`: Error description

**Example:**
```osprey
let getResult = httpGet(clientID: clientId, path: "/users", headers: "")
match getResult {
    Success statusCode => print("Request completed with status: ${statusCode}")
    Err message => print("Request failed: ${message}")
}
```

#### `httpPost(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP POST request with a request body.

**Parameters:**
- `clientID`: Client identifier
- `path`: Request path
- `body`: Request body data
- `headers`: Additional headers

**Example:**
```osprey
let postData = "{\"name\": \"John\", \"email\": \"john@example.com\"}"
let headers = "Content-Type: application/json\r\n"
let postResult = httpPost(clientID: clientId, path: "/users", body: postData, headers: headers)
```

#### `httpPut(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP PUT request.

#### `httpDelete(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP DELETE request.

#### `httpRequest(clientID: Int, method: HttpMethod, path: String, headers: String, body: String) -> Result<StatusCode, String>`

Generic HTTP request function for any HTTP method.

#### `httpCloseClient(clientID: Int) -> Result<Success, String>`

Closes the HTTP client and cleans up resources.

### 14.4 WebSocket Support (Two-Way Communication)

🔒 **IMPLEMENTATION STATUS**: WebSocket functions are implemented with security features but are currently undergoing stability testing.

WebSockets provide real-time, bidirectional communication between client and server. Osprey implements WebSocket support with **MILITARY-GRADE SECURITY** following industry best practices for preventing attacks and ensuring bulletproof operation.

#### 14.4.1 WebSocket Security Implementation

Osprey's WebSocket implementation follows the **OWASP WebSocket Security Guidelines** and implements multiple layers of security protection:

**🛡️ TITANIUM-ARMORED Compilation Security:**
- `_FORTIFY_SOURCE=3`: Maximum buffer overflow protection
- `fstack-protector-all`: Complete stack canary protection  
- `fstack-clash-protection`: Stack clash attack prevention
- `fcf-protection=full`: Control Flow Integrity (CFI) protection
- `ftrapv`: Integer overflow trapping
- `fno-delete-null-pointer-checks`: Prevent null pointer optimizations
- `Wl,-z,relro,-z,now`: Full RELRO with immediate binding
- `Wl,-z,noexecstack`: Non-executable stack protection

**🔒 Cryptographic Security:**
- **OpenSSL SHA-1**: RFC 6455 compliant WebSocket handshake using industry-standard OpenSSL
- **Secure key validation**: 24-character base64 key format validation
- **Constant-time operations**: Memory clearing to prevent timing attacks
- **Error checking**: All OpenSSL operations validated for success

**⚔️ Input Validation Fortress:**
- **WebSocket key format validation**: Strict RFC 6455 compliance
- **Base64 character validation**: Only valid characters accepted
- **Buffer length validation**: Maximum 4096 character keys prevent DoS
- **Integer overflow protection**: All memory calculations checked
- **Memory boundary checking**: No buffer overruns possible

**🏰 Memory Security:**
- **Secure memory allocation**: `calloc()` with zero-initialization
- **Memory clearing**: All sensitive data zeroed before deallocation
- **Bounds checking**: All `snprintf()` operations validated for truncation
- **Safe string operations**: `memcpy()` instead of unsafe `strcpy()`/`strcat()`

#### 14.4.2 Security Standards Compliance

Osprey WebSocket implementation follows these security standards:

**RFC 6455 - WebSocket Protocol Security** ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455)):
- Proper Sec-WebSocket-Accept calculation using SHA-1 + base64
- Origin validation support for CSRF protection
- Secure WebSocket handshake implementation

**OWASP WebSocket Security Cheat Sheet** ([https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets)):
- Input validation on all WebSocket frames
- Authentication and authorization enforcement
- Rate limiting and DoS protection
- Secure error handling without information leakage

**NIST Cybersecurity Framework:**
- Defense in depth through multiple security layers
- Secure coding practices with compiler hardening
- Memory safety through bounds checking
- Cryptographic integrity using OpenSSL

**CWE (Common Weakness Enumeration) Mitigation:**
- CWE-120: Buffer overflow prevention through bounds checking
- CWE-190: Integer overflow protection with `ftrapv`
- CWE-200: Information exposure prevention through secure error handling
- CWE-416: Use-after-free prevention through memory clearing

#### 14.4.3 Security Architecture

```
┌─────────────────────────────────────────────────────────┐
│                TITANIUM SECURITY LAYERS                 │
├─────────────────────────────────────────────────────────┤
│ 🏰 Application Layer: Input Validation Fortress        │
│    • WebSocket key format validation                   │
│    • Base64 character validation                       │
│    • Buffer length enforcement                         │
│    • Memory boundary checking                          │
├─────────────────────────────────────────────────────────┤
│ 🔒 Cryptographic Layer: OpenSSL SHA-1                  │
│    • RFC 6455 compliant handshake                      │
│    • Secure hash computation                           │
│    • Constant-time operations                          │
│    • Error validated operations                        │
├─────────────────────────────────────────────────────────┤
│ ⚔️ Memory Layer: Bulletproof Memory Management         │
│    • Secure allocation with calloc()                   │
│    • Memory clearing before deallocation               │
│    • Safe string operations                            │
│    • Integer overflow protection                       │
├─────────────────────────────────────────────────────────┤
│ 🛡️ Compiler Layer: Military-Grade Hardening           │
│    • Stack protection (canaries + clash protection)    │
│    • Control Flow Integrity (CFI)                      │
│    • FORTIFY_SOURCE=3 buffer overflow protection       │
│    • RELRO + immediate binding                         │
│    • Non-executable stack                              │
└─────────────────────────────────────────────────────────┘
```

#### 14.4.4 Security Testing and Validation

Osprey WebSocket security is validated through:

**🧪 Automated Security Testing:**
- Buffer overflow attack simulation
- Malformed WebSocket key injection
- Integer overflow boundary testing
- Memory corruption detection

**🔍 Static Analysis:**
- Compiler warnings elevated to errors
- Memory safety analysis
- Control flow analysis
- Buffer bounds verification

**⚡ Dynamic Testing:**
- Address Sanitizer (ASan) testing
- Valgrind memory error detection
- Fuzzing with malformed inputs
- DoS resilience testing

#### 14.4.5 Security References and Standards

**Primary Security Standards:**
- **RFC 6455**: "The WebSocket Protocol" - Official WebSocket specification ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455))
- **OWASP WebSocket Security Cheat Sheet**: ([https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets))
- **NIST SP 800-53**: Security Controls for Federal Information Systems
- **ISO 27001**: Information Security Management Standards

**Compiler Security References:**
- **GCC Security Options**: ([https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html))
- **Red Hat Security Guide**: "Defensive Coding Practices"
- **Microsoft SDL**: Security Development Lifecycle practices
- **Google Safe Coding Practices**: Memory safety guidelines

**Cryptographic Standards:**
- **FIPS 180-4**: SHA-1 cryptographic hash standard
- **RFC 3174**: US Secure Hash Algorithm 1 (SHA1) ([https://tools.ietf.org/html/rfc3174](https://tools.ietf.org/html/rfc3174))
- **OpenSSL Security Advisories**: ([https://www.openssl.org/news/secadv.html](https://www.openssl.org/news/secadv.html))

**Memory Security Research:**
- **"Control Flow Integrity"** by Abadi et al. - CFI protection principles
- **"Stack Canaries"** - Buffer overflow detection mechanisms  
- **"RELRO"** - Read-only relocations for exploit mitigation
- **"FORTIFY_SOURCE"** - Compile-time and runtime buffer overflow detection

#### `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`

Establishes a WebSocket connection.

**Parameters:**
- `url`: WebSocket URL (e.g., "ws://localhost:8080/chat")
- `messageHandler`: Callback function to handle incoming messages

**Returns:**
- `Success(wsID)`: WebSocket connection identifier
- `Err(message)`: Connection error

**Example:**
```osprey
fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}

let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)
```

#### `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

Sends a message through the WebSocket connection.

**Parameters:**
- `wsID`: WebSocket identifier
- `message`: Message to send

**Example:**
```osprey
let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
```

#### `websocketClose(wsID: Int) -> Result<Success, String>`

Closes the WebSocket connection.

### 14.5 Streaming Response Bodies

Osprey automatically handles response streaming to prevent memory issues with large responses:

#### Complete Response
For small responses that fit in memory:
```osprey
HttpResponse {
    status: 200,
    contentType: "application/json",
    partialBody: "{\"data\": \"small response\"}",
    isComplete: true,
    streamFd: -1
}
```

#### Streamed Response
For large responses that should be streamed:
```osprey
HttpResponse {
    status: 200,
    contentType: "application/octet-stream",
    streamFd: fileDescriptor,  // File descriptor to stream from
    isComplete: false,
    contentLength: -1,         // -1 indicates chunked encoding
    partialBody: "",
    partialLength: 0
}
```

### 14.6 Error Handling in HTTP

All HTTP functions return Result types following Osprey's principle of making illegal states unrepresentable:

```osprey
// Server creation with error handling
let serverResult = httpCreateServer(port: 8080, address: "127.0.0.1")
match serverResult {
    Success serverId => {
        print("Server created successfully")
        let listenResult = httpListen(serverID: serverId, handler: myHandler)
        match listenResult {
            Success _ => print("Server is listening")
            Err error => print("Failed to start listening: ${error}")
        }
    }
    Err error => print("Failed to create server: ${error}")
}
```

### 14.7 Fiber-Based Concurrency

HTTP servers automatically spawn a new fiber for each incoming request, enabling thousands of concurrent connections:

```osprey
// Each request handler runs in its own fiber
fn handleRequest(request: HttpRequest) -> Result<HttpResponse, String> = {
    // This function runs in a separate fiber for each request
    // Multiple requests can be processed concurrently
    let result = processLongRunningTask(request.body)
    Success(HttpResponse {
        status: 200,
        contentType: "application/json",
        partialBody: result,
        isComplete: true,
        streamFd: -1
    })
}
```

### 14.8 Complete HTTP Server Example

```osprey
// Create and start an HTTP server
let serverResult = httpCreateServer(port: 8080, address: "0.0.0.0")
match serverResult {
    Success serverId => {
        fn apiHandler(request: HttpRequest) -> Result<HttpResponse, String> = match request.method {
            GET => match request.path {
                "/api/health" => Success(HttpResponse {
                    status: 200,
                    contentType: "application/json",
                    partialBody: "{\"status\": \"healthy\", \"timestamp\": \"${getCurrentTime()}\"}",
                    isComplete: true,
                    streamFd: -1
                })
                "/api/users" => Success(HttpResponse {
                    status: 200,
                    contentType: "application/json",
                    partialBody: getUsersJson(),
                    isComplete: true,
                    streamFd: -1
                })
                _ => Success(HttpResponse {
                    status: 404,
                    contentType: "application/json", 
                    partialBody: "{\"error\": \"Not Found\"}",
                    isComplete: true,
                    streamFd: -1
                })
            }
            POST => match request.path {
                "/api/users" => createUser(request.body)
                _ => Err("Unsupported POST endpoint")
            }
            _ => Err("Method not allowed")
        }
        
        let listenResult = httpListen(serverID: serverId, handler: apiHandler)
        match listenResult {
            Success _ => print("🚀 HTTP server listening on http://0.0.0.0:8080")
            Err error => print("❌ Failed to start server: ${error}")
        }
    }
    Err error => print("❌ Failed to create server: ${error}")
}
```

## 15. Error Handling

## 16. Examples

## 17. Built-in Functions Reference

### 17.1 Basic I/O Functions

#### `print(value: int | string | bool) -> int`
Prints the given value to standard output with automatic type conversion.

**Parameters:**
- `value: int | string | bool` - The value to print (int, bool, string, or expression)

**Returns:** `int` - Exit code from puts function

**Examples:**
```osprey
print("Hello World")
print(42)
print(true)
print(x + y)
```

#### `input() -> int`
Reads an integer from stdin. Blocks until user enters a number.

**Parameters:** None

**Returns:** `int` - The number entered by user

**Examples:**
```osprey
let x = input()
let age = input()
```

#### `toString(value: int | string | bool) -> string`
Converts any value to its string representation. Used internally by print and string interpolation.

**Parameters:**
- `value: int | string | bool` - The value to convert to string

**Returns:** `string` - String representation of the value

**Examples:**
```osprey
let str = toString(42)
let msg = toString(true)
```

### 17.2 Functional Iterator Functions

#### `range(start: int, end: int) -> Iterator<int>`
Creates an iterator that generates integers from start (inclusive) to end (exclusive). Used with functional iterator functions like forEach, map, filter, and fold.

**Parameters:**
- `start: int` - Starting value (inclusive)
- `end: int` - Ending value (exclusive)

**Returns:** `Iterator<int>` - Iterator struct containing start and end values

**Examples:**
```osprey
range(1, 5)    // generates 1, 2, 3, 4
range(0, 3)    // generates 0, 1, 2
range(10, 13)  // generates 10, 11, 12
range(1, 10) |> forEach(print)
```

#### `forEach(iterator: Iterator<T>, function: T -> U) -> T`
Applies a function to each element in an iterator for side effects. This is the primary way to iterate through ranges and apply operations to each element.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to traverse (usually from range())
- `function: T -> U` - Function to apply to each element

**Returns:** `T` - Final counter value after iteration

**Examples:**
```osprey
range(1, 5) |> forEach(print)          // prints 1, 2, 3, 4
forEach(range(0, 3), double)           // calls double(0), double(1), double(2)
range(1, 10) |> forEach(square)
forEach(range(-2, 3), print)          // prints -2, -1, 0, 1, 2
```

#### `map(iterator: Iterator<T>, function: T -> U) -> U`
Transforms each element in an iterator by applying a function. Returns the result of the transformation function applied to each element.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to transform (usually from range())
- `function: T -> U` - Transformation function to apply

**Returns:** `U` - Result of applying function to each element

**Examples:**
```osprey
range(1, 5) |> map(double)    // applies double to 1, 2, 3, 4
map(range(0, 3), square)      // applies square to 0, 1, 2
range(1, 6) |> map(addFive)
```

#### `filter(iterator: Iterator<T>, predicate: T -> bool) -> T`
Selects elements from an iterator based on a predicate function. Only elements where the predicate returns true are processed.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to filter (usually from range())
- `predicate: T -> bool` - Function that returns true for elements to keep

**Returns:** `T` - Filtered results

**Examples:**
```osprey
range(1, 10) |> filter(isEven)
filter(range(0, 20), isPositive)
range(-5, 6) |> filter(isGreaterThanZero)
```

#### `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`
Reduces an iterator to a single value by repeatedly applying a function. Also known as reduce or accumulate in other languages.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to reduce (usually from range())
- `initial: U` - Initial accumulator value
- `function: (U, T) -> U` - Function that combines accumulator with each element

**Returns:** `U` - Final accumulated value

**Examples:**
```osprey
range(1, 5) |> fold(0, add)          // sum: 0+1+2+3+4 = 10
fold(range(1, 6), 1, multiply)       // product: 1*1*2*3*4*5 = 120
range(0, 10) |> fold(0, max)
```

### 17.3 Pipe Operator

#### `|>` - Pipe Operator
The pipe operator takes the result of the left expression and passes it as the first argument to the right function. This enables elegant functional programming and method chaining.

**Syntax:** `expression |> function`

**Type:** `T |> (T -> U) -> U`

**Description:**
The pipe operator creates clean, readable function composition by allowing you to chain operations from left to right, making the data flow explicit and natural to read.

**Rules:**
- Left side can be any expression
- Right side must be a function or function call
- Creates clean, readable function composition
- Enables Haskell/Rust-style functional programming

**Examples:**
```osprey
// Basic piping
5 |> double |> print                 // Equivalent to: print(double(5))

// Iterator chaining
range(1, 10) |> forEach(print)
range(1, 5) |> map(square) |> fold(0, add)

// Complex chains
range(0, 20) |> filter(isEven) |> map(double) |> forEach(print)

// Multiple operations
let result = input() |> double |> square |> toString

// Nested operations
range(1, 10) 
  |> map(square) 
  |> filter(isEven) 
  |> fold(0, add) 
  |> print
```

### 17.4 Functional Programming Patterns

The combination of iterator functions and the pipe operator enables powerful functional programming patterns:

#### Chaining Pattern
```osprey
// Transform -> Filter -> Aggregate
range(1, 20)
  |> map(square)           // Square each number
  |> filter(isEven)        // Keep only even results
  |> fold(0, add)          // Sum them up
  |> print                 // Print the result
```

#### Side Effect Pattern
```osprey
// Process each element for side effects
range(1, 100)
  |> filter(isPrime)
  |> forEach(print)        // Print each prime
```

#### Data Transformation Pattern
```osprey
// Transform data through multiple stages
input()
  |> validateInput
  |> normalizeData
  |> processData
  |> formatOutput
  |> print
```

### 17.5 Fiber Types and Concurrency

Osprey provides lightweight concurrency through fiber types. Unlike traditional function-based approaches, fibers are proper type instances constructed using Osprey's standard type construction syntax.

#### Core Fiber Types

**`Fiber<T>`** - A lightweight concurrent computation that produces a value of type T
**`Channel<T>`** - A communication channel for passing values of type T between fibers

#### Fiber Construction

Fibers are created using standard type construction syntax:

```osprey
// Create a fiber that computes a value
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(n: 1000) 
}

// Create a fiber with more complex computation
let worker = Fiber<String> { 
    computation: fn() => {
        processData()
        "completed"
    }
}

// Create a parameterized fiber
let calculator = Fiber<Int> { 
    computation: fn() => multiply(x: 10, y: 20) 
}
```

#### Spawn Syntax Sugar

For convenience, Osprey provides `spawn` as syntax sugar for creating and immediately starting a fiber:

```osprey
// Using spawn (syntax sugar)
let result = spawn 42

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => 42 }
let result = fiber

// More complex spawn
let computation = spawn (x * 2 + y)

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => x * 2 + y }
let computation = fiber
```

The `spawn` keyword immediately evaluates the expression in a new fiber context, making it convenient for quick concurrent computations without the full type construction syntax.

#### Channel Construction

Channels are created using type construction syntax:

```osprey
// Unbuffered (synchronous) channel
let sync_channel = Channel<Int> { capacity: 0 }

// Buffered (asynchronous) channel  
let async_channel = Channel<String> { capacity: 10 }

// Large buffer channel
let buffer_channel = Channel<Int> { capacity: 100 }
```

#### Fiber Operations

Once created, fibers and channels are manipulated using functional operations:

**`await(fiber: Fiber<T>) -> T`** - Wait for fiber completion and get result
**`send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`** - Send value to channel
**`recv(channel: Channel<T>) -> Result<T, ChannelError>`** - Receive value from channel
**`yield() -> Unit`** - Voluntarily yield control to scheduler

```osprey
// Create and await a fiber
let task = Fiber<Int> { computation: fn() => heavyComputation() }
let result = await(task)

// Channel communication
let ch = Channel<String> { capacity: 5 }
send(ch, "hello")
let message = recv(ch)

// Yielding control
yield()
```

#### Complete Fiber Example

```osprey
// Producer fiber
let producer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        send(ch, 1)
        send(ch, 2) 
        send(ch, 3)
    }
}

// Consumer fiber
let consumer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        let value1 = recv(ch)
        let value2 = recv(ch)
        let value3 = recv(ch)
        print("Received: ${value1}, ${value2}, ${value3}")
    }
}

// Start both fibers
await(producer)
await(consumer)
```

#### Select Expression for Channel Multiplexing

The `select` expression allows waiting on multiple channel operations:

```osprey
let ch1 = Channel<String> { capacity: 1 }
let ch2 = Channel<Int> { capacity: 1 }

let result = select {
    msg => recv(ch1) => process_string(msg)
    num => recv(ch2) => process_number(num)
    _ => timeout_handler()
}
```

#### Rust Interoperability

Osprey fibers are designed to interoperate with Rust's async/await system:

```osprey
// Osprey fiber that calls Rust async function
extern fn rust_async_task() -> Future<Int>

let osprey_task = Fiber<Int> {
    computation: fn() => await(rust_async_task())
}

let result = await(osprey_task)
```

#### Fiber-Isolated Modules

Each fiber gets its own isolated instance of modules, preventing data races:

```osprey
module Counter {
    let mut count = 0
    fn increment() -> Int = { count = count + 1; count }
    fn get() -> Int = count
}

// Each fiber has its own Counter instance
let fiber1 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}
let fiber2 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}

// These will both return 1, not 1 and 2
let result1 = await(fiber1)  // 1
let result2 = await(fiber2)  // 1 (separate instance)
```

---

## 18. Security and Sandboxing

The Osprey compiler includes built-in security controls to restrict access to potentially dangerous functionality like network operations and file system access. This is essential for safe code execution in environments like web compilers where untrusted code may be executed.

### 18.1 Security Flags

#### `--sandbox`
Enables sandbox mode, which disables all potentially risky operations:
- HTTP/HTTPS operations (httpCreateServer, httpGet, httpPost, etc.)
- WebSocket operations (websocketConnect, websocketSend, etc.)
- File system access (when implemented)
- Foreign Function Interface (FFI)
- Process execution

**Example:**
```bash
osprey program.osp --sandbox --llvm
```

#### Granular Security Flags

For more granular control, you can disable specific categories of operations:

- `--no-http`: Disable HTTP client and server functions
- `--no-websocket`: Disable WebSocket client and server functions  
- `--no-fs`: Disable file system read/write operations
- `--no-ffi`: Disable foreign function interface

**Examples:**
```bash
# Disable only HTTP operations
osprey program.osp --no-http --compile

# Disable HTTP and WebSocket operations
osprey program.osp --no-http --no-websocket --run

# Disable file system access only
osprey program.osp --no-fs --llvm
```

### 18.2 Security Policies

#### Default Security (Permissive)
By default, all operations are allowed for backward compatibility and normal development use.

#### Sandbox Security (Restrictive)
When `--sandbox` is used, all potentially dangerous functions are unavailable. This is recommended for:
- Web-based code execution
- Untrusted code evaluation
- Educational environments
- Code review systems

### 18.3 Blocked Functions by Category

#### HTTP Functions
When HTTP access is disabled (`--no-http` or `--sandbox`), these functions are unavailable:
- `httpCreateServer` - Create HTTP server
- `httpListen` - Start HTTP server listening
- `httpStopServer` - Stop HTTP server
- `httpCreateClient` - Create HTTP client
- `httpGet` - HTTP GET request
- `httpPost` - HTTP POST request
- `httpPut` - HTTP PUT request
- `httpDelete` - HTTP DELETE request
- `httpRequest` - Generic HTTP request
- `httpCloseClient` - Close HTTP client

#### WebSocket Functions
When WebSocket access is disabled (`--no-websocket` or `--sandbox`), these functions are unavailable:
- `websocketConnect` - Connect to WebSocket server
- `websocketSend` - Send WebSocket message
- `websocketClose` - Close WebSocket connection
- `websocketCreateServer` - Create WebSocket server
- `websocketServerListen` - Start WebSocket server
- `websocketServerSend` - Send message to specific client
- `websocketServerBroadcast` - Broadcast message to all clients
- `websocketStopServer` - Stop WebSocket server

#### File System Functions (Future)
When file system access is disabled (`--no-fs` or `--sandbox`), these functions will be unavailable:
- `readFile` - Read file contents
- `writeFile` - Write file contents
- `deleteFile` - Delete file
- `createDirectory` - Create directory
- `listDirectory` - List directory contents

### 18.4 Function Availability

In different security modes, certain functions are simply not available in the language:

**Sandbox Mode**: Only safe functions like `print`, `toString`, `range`, etc. are available. Dangerous functions like `httpCreateServer` or `websocketConnect` result in "undefined function" compile errors.

**Partial Restrictions**: When specific categories are disabled (e.g., `--no-http`), those functions are unavailable while others remain accessible.

**Default Mode**: All functions are available.
- A human-readable explanation

### 18.5 Programming Best Practices

#### For Safe Code
Write code that doesn't use security-sensitive functions:
```osprey
// Safe operations - work in all security modes
let x = 42
let y = 24
let sum = x + y
print("Sum: ")
print(sum)
```

#### For Network Code
When writing network code, be aware that it may be restricted:
```osprey
// This will fail in sandbox mode or with --no-http
let serverID = httpCreateServer(port: 8080, address: "127.0.0.1")
```

### 18.6 Implementation Details

#### Security Configuration
Security settings are configured at compilation time and cannot be bypassed by the compiled program. The security checks happen during the LLVM IR generation phase, preventing security-sensitive functions from being included in the generated code.

#### Performance Impact
Security checks add minimal overhead during compilation and no runtime overhead, as restricted functions are simply not compiled into the final program.

#### Backward Compatibility
All existing code continues to work with default settings. Security restrictions are opt-in and don't affect normal development workflows.

#### Integration with Web Compiler
The security features are designed specifically for web compiler integration:

```javascript
// Example web compiler usage
const result = await compileOsprey(sourceCode, {
    mode: 'sandbox',  // Enable sandbox mode
    outputFormat: 'llvm'
});
```

#### Security Summary
When using security restrictions, the compiler will display a security summary:

```bash
# Sandbox mode
Security: SANDBOX MODE - All risky operations disabled

# Partial restrictions
Security: Allowed=[FileRead,FileWrite,FFI] Blocked=[HTTP,WebSocket]
```

---

## Summary

Osprey is a modern functional programming language with:

- **Type Safety**: No runtime panics, all errors handled explicitly via Result types
- **Named Arguments**: Multi-parameter functions require named arguments for clarity
- **Functional Programming**: Powerful iterator functions with pipe operator
- **Lightweight Fibers**: Zero-cost concurrency with Rust-like async/await
- **Fiber-Isolated Modules**: No global state, each fiber gets its own module instances
- **Rust Interoperability**: Seamless integration with Rust libraries
- **Memory Safety**: No shared mutable state between fibers

**Key Innovation**: The fiber-isolated module system eliminates data races by design while maintaining clean encapsulation through accessor patterns.

---

**End of Specification**

This specification defines the complete syntax and semantics of the Osprey programming language, including its revolutionary fiber-isolated module system and lightweight concurrency features. The accompanying `osprey.g4` grammar file provides the formal ANTLR4 grammar definition for parsing.