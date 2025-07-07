3. [Syntax](0003-Syntax.md)
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

constraint := WHERE function_name
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
record_type := 'type' ID '=' '{' field_declarations '}' validation?

field_declarations := field_declaration (',' field_declaration)*
field_declaration := ID ':' type
validation := 'where' function_name
```

**Examples:**
```osprey
// Simple record type (no validation)
type Point = { x: Int, y: Int }

// Record with type-level validation
type Person = { 
    name: String,
    age: Int,
    email: String,
    confirmEmail: String
} where validatePerson

// Validation function defines all constraints
fn validatePerson(person: Person) -> Result<Person, String> = {
    if person.name == "" then Error("Name cannot be empty")
    else if person.age < 0 then Error("Age must be positive")
    else if person.age > 150 then Error("Age must be realistic")
    else if not(isValidEmail(person.email)) then Error("Invalid email format")
    else if person.email != person.confirmEmail then Error("Email confirmation mismatch")
    else Success(person)
}

// Mixed record and union types
type User = Guest { sessionId: String }
          | Member { 
              id: Int,
              name: String,
              joinDate: String
            } where validateMember
```

#### 3.7.2 Type Construction

Type constructors create instances of record types using elegant field syntax. If a type has a `where` validation function, the constructor returns a `Result` type.

```
type_constructor := type_name '{' field_assignments '}'
field_assignments := field_assignment (',' field_assignment)*
field_assignment := ID ':' expression
```

**Construction Examples:**
```osprey
// Simple construction (no validation)
type Point = { x: Int, y: Int }
let point = Point { x: 10, y: 20 }  // Returns Point directly

// Type with validation function
type Person = { 
    name: String, 
    age: Int, 
    email: String,
    confirmEmail: String
} where validatePerson

fn validatePerson(person: Person) -> Result<Person, String> = match person.name {
    "" => Error("Name cannot be empty")
    _ => match person.age {
        0 => Error("Age must be positive")
        _ => match person.email {
            email if email == person.confirmEmail => Success(person)
            _ => Error("Email confirmation mismatch")
        }
    }
}

// Construction with validation (returns Result)
let person = Person { 
    name: "Alice", 
    age: 25, 
    email: "alice@example.com",
    confirmEmail: "alice@example.com"
}  // Returns Result<Person, String>

// Handle construction results
match person {
    Success { value } => print("Created person: ${value.name}")
    Error { message } => print("Construction failed: ${message}")
}
```

#### 3.7.3 Construction Result Types

**CRITICAL**: Type constructors with validation functions return `Result<T, String>`:

- **Unvalidated types**: Direct construction returns the type
- **Validated types**: Construction returns `Result<T, String>`

**Rule**: If a type has a WHERE validation function, the constructor ALWAYS returns a Result type.

```osprey
// No validation = direct construction
type Point = { x: Int, y: Int }
let point = Point { x: 10, y: 20 }  // Returns Point

// With validation = Result construction  
type Person = { 
    name: String,
    age: Int
} where validatePerson

fn validatePerson(person: Person) -> Result<Person, String> = match person.name {
    "" => Error("Name cannot be empty")
    _ => match person.age {
        0 => Error("Age must be positive")
        _ => Success(person)
    }
}

let person = Person { name: "Alice", age: 25 }  // Returns Person directly (no validation)

// For types with validation, construction error handling required
match Person { name: "", age: 25 } {
    Success { value } => useValidPerson(value)
    Error { message } => print("Validation failed: ${message}")
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
// NOTE: The old field-level constraint syntax has been replaced with type-level validation

type Person = {
    name: String,
    age: Int,
    email: String,
    confirmEmail: String
} where validatePerson

fn validatePerson(person: Person) -> Result<Person, String> = {
    // Type-level validation can check all fields and cross-field relationships
    if person.name == "" then Error("Name cannot be empty")
    else if person.age < 0 || person.age > 150 then Error("Age must be between 0 and 150")
    else if person.email == "" then Error("Email cannot be empty")
    else if person.email != person.confirmEmail then Error("Email confirmation must match")
    else Success(person)
}

// For examples without validation functions, these would be simple types:
type Rectangle = {
    width: Int,
    height: Int,
    area: Int
}

type CreditCard = {
    number: String,
    expiryMonth: Int,
    expiryYear: Int,
    cvv: String
} where validateCreditCard

fn validateCreditCard(card: CreditCard) -> Result<CreditCard, String> = {
    // Type-level validation for credit card
    if card.number == "" then Error("Card number required")
    else if card.expiryMonth < 1 || card.expiryMonth > 12 then Error("Invalid expiry month")
    else Success(card)
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
// Simple types without validation (return the type directly)
type UserAccount = {
    username: String,
    password: String,
    loginHour: Int
}

type NetworkConfig = {
    port: Int,
    host: String
}

type Appointment = {
    dayOfWeek: String,
    hour: Int,
    duration: Int
}
```

**Type-Level Validation Function Requirements:**
- Must return `Result<T, String>` type - Note: that string is a placeholder and arbitrary type arguments will be available here
- Receive the entire constructed object as parameter
- Can validate individual fields and cross-field relationships  
- Can call other functions (including built-ins)
- Must be pure functions (no side effects)

**Performance Characteristics:**
- **Single validation pass**: All validation happens once during construction
- **Type safety**: Returns Result types that must be pattern matched
- **Detailed errors**: String error messages provide clear validation feedback
- **Cross-field validation**: Can validate relationships between multiple fields

#### 3.7.9 Type-Level Validation Examples

Type-level validation functions receive the entire constructed object and can implement complex validation logic:

```osprey
// Validation function with cross-field checking
type StrongPassword = {
    value: String,
    confirmValue: String
} where validateStrongPassword

fn validateStrongPassword(pwd: StrongPassword) -> Result<StrongPassword, String> = match pwd.value {
    "" => Error("Password cannot be empty")
    _ => match length(pwd.value) {
        len if len < 8 => Error("Password must be at least 8 characters")
        _ => match pwd.value {
            val if val != pwd.confirmValue => Error("Password confirmation must match")
            _ => Success(pwd)
        }
    }
}

// Validation function for business logic  
type UserAccount = {
    username: String,
    email: String,
    age: Int
} where validateUserAccount

fn validateUserAccount(user: UserAccount) -> Result<UserAccount, String> = match user.username {
    "" => Error("Username cannot be empty")
    _ => match user.age {
        age if age < 13 => Error("User must be at least 13 years old")
        age if age > 120 => Error("Age must be realistic")
        _ => Success(user)
    }
}

// Simple types without validation (return the type directly)
type Point = {
    x: Int,
    y: Int
}

type Rectangle = {
    width: Int,
    height: Int,
    area: Int
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
    name: String
} where validateUser

fn validateUser(user: User) -> Result<User, String> = match user.name {
    "" => Error("Name cannot be empty")
    _ => Success(user)
}

let user = User { name: "alice" }  // Returns Result<User, String>

// COMPILATION ERROR: Cannot access field on Result type
print("${user.name}")  
// Should produce: "cannot access field 'name' on Result<User, String> type - pattern matching required"

let name = user.name
// Should produce: "field access requires pattern matching on Result type"
```

**Required Error Messages:**
- **Field access on Result**: `"cannot access field 'FIELD' on Result<TYPE, String> type - pattern matching required"`
- **Assignment from Result field**: `"field access requires pattern matching on Result type"`
- **Missing pattern matching**: `"validated types return Result - use match expression to handle success/failure"`

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


#### Binary Expressions
```
binary_expression := multiplicative_expression (('+' | '-') multiplicative_expression)*
multiplicative_expression := unary_expression (('*' | '/') unary_expression)*
```

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

#### List Access (Safe)
```
list_access := expression '[' INT ']'  // Returns Result<T, IndexError>
```

🚨 **CRITICAL SAFETY GUARANTEE**: List access **ALWAYS** returns `Result<T, IndexError>` - **NO PANICS, NO NULLS, NO EXCEPTIONS**

**MANDATORY PATTERN MATCHING REQUIRED:**
```osprey
let numbers = [1, 2, 3, 4]

// ✅ CORRECT: Pattern matching required
let firstResult = numbers[0]  // Returns Result<Int, IndexError>
match firstResult {
    Success { value } => print("First: ${value}")
    Error { message } => print("Index out of bounds: ${message}")
}

// ✅ CORRECT: Inline pattern matching
let second = match numbers[1] {
    Success { value } => value
    Error { _ } => -1  // Default value for out-of-bounds
}

// ✅ CORRECT: Bounds-safe iteration
let commands = ["echo hello", "echo world"]
match commands[0] {
    Success { value } => {
        print("Executing: ${value}")
        spawnProcess(value)
    }
    Error { message } => print("No command at index 0: ${message}")
}
```

**FUNDAMENTAL SAFETY PRINCIPLE**: Array access can fail (index out of bounds), therefore it MUST return Result types to enforce explicit error handling and prevent runtime crashes.

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

#### Field Access Rules and Restrictions

**✅ ALLOWED - Field Access on Record Types:**
```osprey
type User = { id: Int, name: String, email: String }
let user = User { id: 1, name: "Alice", email: "alice@example.com" }

let userId = user.id          // Valid: direct field access
let userName = user.name      // Valid: direct field access
let userEmail = user.email    // Valid: direct field access
```

**❌ FORBIDDEN - Field Access on `any` Types:**
```osprey
fn processAnyValue(value: any) -> String = {
    // ERROR: Cannot access fields on 'any' type
    let result = value.name   // Compilation error
    return result
}

// CORRECT: Use pattern matching for 'any' types
fn processAnyValue(value: any) -> String = match value {
    person: { name } => person.name        // Extract field via pattern matching
    user: User { name } => name           // Type-specific pattern matching
    _ => "unknown"
}
```

**❌ FORBIDDEN - Field Access on Result Types:**
```osprey
type Person = { 
    name: String
} where validatePerson

fn validatePerson(person: Person) -> Result<Person, String> = match person.name {
    "" => Error("Name cannot be empty")
    _ => Success(person)
}

let personResult = Person { name: "Alice" }  // Returns Result<Person, String>

// ERROR: Cannot access field on Result type
let name = personResult.name   // Compilation error

// CORRECT: Use pattern matching on Result types
match personResult {
    Ok { value } => print("Name: ${value.name}")    // Access field after unwrapping
    Err { error } => print("Construction failed: ${error}")
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