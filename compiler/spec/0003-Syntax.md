# Syntax

- [Program Structure](#program-structure)
- [Import Statements](#import-statements)
- [Let Declarations](#let-declarations)
- [Function Declarations](#function-declarations)
- [Extern Declarations](#extern-declarations)
- [Type Declarations](#type-declarations)
- [Record Types and Type Constructors](#record-types-and-type-constructors)
- [Expressions](#expressions)
- [Block Expressions](#block-expressions)
- [Match Expressions](#match-expressions)

## Program Structure

An Osprey program consists of a sequence of top-level statements and modules.

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

## Import Statements

Basic import parsing is implemented but module resolution is limited.

```
importStmt := IMPORT ID (DOT ID)*
```

**Examples:**
```osprey
import std
import std.io
import graphics.canvas
```

### Let Declarations

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

### Function Declarations

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

### Extern Declarations

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

### Type Declarations

```
typeDecl := docComment? TYPE ID (LT typeParamList GT)? EQ (unionType | recordType)

unionType := variant (BAR variant)*
recordType := LBRACE fieldDeclarations RBRACE

variant := ID (LBRACE fieldDeclarations RBRACE)?

fieldDeclarations := fieldDeclaration (COMMA fieldDeclaration)*
fieldDeclaration := ID COLON type constraint?

constraint := WHERE function_name
```

**Examples:**
```osprey
type Color = Red | Green | Blue

type Shape = Circle { radius: int } 
           | Rectangle { width: int, height: int }

type Result = Success { value: string } 
            | Error { message: string }
```

### Record Types and Type Constructors

Record types define structured data with named fields using the following syntax:

```
record_type := 'type' ID '=' '{' field_declarations '}' validation?

field_declarations := field_declaration (',' field_declaration)*
field_declaration := ID ':' type
validation := 'where' function_name

type_constructor := type_name '{' field_assignments '}'
field_assignments := field_assignment (',' field_assignment)*
field_assignment := ID ':' expression
```

**Examples:**
```osprey
type Point = { x: int, y: int }
type Person = { name: string, age: int } where validatePerson

let point = Point { x: 10, y: 20 }
let person = Person { name: "Alice", age: 25 }
```

For complete record type semantics, construction rules, field access, constraints, and validation behavior, see [Type System - Record Types](0004-TypeSystem.md#record-types).

### Expressions

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


#### Binary Expressions
```
binary_expression := logical_or_expression
logical_or_expression := logical_and_expression ('||' logical_and_expression)*
logical_and_expression := comparison_expression ('&&' comparison_expression)*
comparison_expression := additive_expression (('==' | '!=' | '<' | '>' | '<=' | '>=') additive_expression)*
additive_expression := multiplicative_expression (('+' | '-') multiplicative_expression)*
multiplicative_expression := unary_expression (('*' | '/') unary_expression)*
```

**Operator Precedence (highest to lowest):**
1. Unary operators (`!`, `-`)
2. Multiplicative (`*`, `/`)  
3. Additive (`+`, `-`)
4. Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`)
5. Logical AND (`&&`)
6. Logical OR (`||`)

#### Boolean Pattern Matching
Use pattern matching for conditional logic:

**Examples:**
```osprey
let result = match x > 0 {
    true => "positive"
    false => "zero or negative"
}

let max = match a > b {
    true => a
    false => b
}
```

#### List Access

```
list_access := expression '[' INT ']'  // Returns Result<T, IndexError>
```

List access always returns `Result<T, IndexError>` for bounds safety and must be handled with pattern matching:

```osprey
let numbers = [1, 2, 3, 4]

let firstResult = numbers[0]  // Result<int, IndexError>
match firstResult {
    Success { value } => print("First: ${value}")
    Error { message } => print("Index error: ${message}")
}

// Inline pattern matching
let second = match numbers[1] {
    Success { value } => value
    Error { _ } => -1
}
```

#### Field Access

Field access uses dot notation to access fields of record types:

```
field_access := expression '.' identifier
```

**Examples:**
```osprey
type Person = { name: String, age: Int }
let person = Person { name: "Alice", age: 25 }

// Field access on record types
let name = person.name        // "Alice"
let age = person.age          // 25

// Field access in expressions
print("Name: ${person.name}")
print("Age: ${person.age}")

// Field access in function calls
sendEmail(to: person.name, subject: "Hello")
```

#### Field Access Rules

Field access is allowed on record types:

```osprey
type User = { id: int, name: string }
let user = User { id: 1, name: "Alice" }
let userId = user.id          // Valid
let userName = user.name      // Valid
```

Field access requires pattern matching for:
- **`any` types**: Extract fields through structural patterns
- **Result types**: Unwrap Result before accessing fields
- **Union types**: Match variant before accessing fields

```osprey
// any type - use pattern matching
fn processAny(value: any) -> string = match value {
    person: { name } => person.name
    _ => "unknown"
}

// Result type - unwrap first
match personResult {
    Success { value } => print("Name: ${value.name}")
    Error { message } => print("Error: ${message}")
}
```

#### Structural Field Access

Osprey supports structural field access through pattern matching, allowing access to fields regardless of the specific type:

```osprey
// Any type with 'name' and 'age' fields can be matched
fn processEntity(entity: any) -> String = match entity {
    { name, age } => "Entity ${name} is ${age}"           // Structural matching
    person: { name } => "Named entity: ${person.name}"    // Named structural matching  
    _ => "Unknown entity"
}

// Works with any type that has these fields
type Person = { name: String, age: Int }
type Employee = { name: String, age: Int, department: String }

let person = Person { name: "Alice", age: 25 }
let employee = Employee { name: "Bob", age: 30, department: "Engineering" }

// Both work with the same function
print(processEntity(person))     // "Entity Alice is 25"
print(processEntity(employee))   // "Entity Bob is 30"
```

#### Immutability and Field Access

All record types are immutable by default. Field access returns the current value and cannot be used for assignment:

```osprey
type Point = { x: Int, y: Int }
let point = Point { x: 10, y: 20 }

let x = point.x              // ✅ Valid: read field value
let y = point.y              // ✅ Valid: read field value

// ERROR: Cannot assign to fields (immutable by default)
point.x = 15                 // ❌ Compilation error

// CORRECT: Create new instance with updated values
let newPoint = point { x: 15 }   // ✅ Non-destructive update syntax
print("New point: ${newPoint.x}, ${newPoint.y}")  // 15, 20
```

#### Primary Expressions
```
primary_expression := literal | list_literal | identifier | '(' expression ')' 
                   | list_access | field_access | lambda_expression | block_expression | match_expression
```

### Block Expressions

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

#### Block Scoping Rules

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

#### Block Return Values

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

#### Block Expressions in Match Arms

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

#### Function Bodies as Blocks

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

#### Type Safety and Inference

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

Block expressions are zero-cost abstractions: scoping is resolved at compile time, and simple blocks are optimized away. See [Block Expressions](0008-BlockExpressions.md) for complete details on semantics and usage patterns.

### Match Expressions

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

## Language Semantics

### Variable Binding

- `let` creates immutable bindings
- `mut` creates mutable bindings
- Variables must be initialized at declaration
- Shadowing is allowed in nested scopes

### Function Semantics

- Functions are first-class values
- All functions are pure (no side effects except I/O)
- Recursive functions are supported
- Tail recursion is optimized

### Evaluation Order

- Expressions are evaluated left-to-right
- Function arguments are evaluated before the function call
- Short-circuit evaluation for logical operators