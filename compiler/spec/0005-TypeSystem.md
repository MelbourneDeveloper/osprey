5. [Type System](0005-TypeSystem.md)
   - [Hindley-Milner Type Inference Foundation](#50-hindley-milner-type-inference-foundation)
   - [Built-in Types](#51-built-in-types)
       - [Function Types](#function-types)
       - [Record Types](#record-types)
   - [Built-in Error Types](#52-built-in-error-types)
   - [Hindley-Milner Type Inference](#53-hindley-milner-type-inference)
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

### 5.0 Hindley-Milner Type Inference Foundation

**🔥 CORE SPECIFICATION**: Osprey implements complete **Hindley-Milner type inference** as its foundational type system. This is a **MANDATORY REQUIREMENT** for compiler implementation.

**Academic Foundation & Implementation Requirements:**
- **Hindley, R. (1969)**: "The Principal Type-Scheme of an Object in Combinatory Logic" - Communications of the ACM 12(12):719-721
- **Milner, R. (1978)**: "A Theory of Type Polymorphism in Programming" - Journal of Computer and System Sciences 17:348-375  
- **Damas, L. & Milner, R. (1982)**: "Principal type-schemes for functional programs" - POPL '82

**🔥 CRITICAL IMPLEMENTATION MANDATES:**

1. **COMPLETE TYPE INFERENCE**: Variables and functions MAY be declared without type annotations when types can be inferred through Hindley-Milner unification
2. **PRINCIPAL TYPES**: Every well-typed expression MUST have a unique most general (polymorphic) type
3. **SOUNDNESS GUARANTEE**: If type checker accepts a program, NO runtime type errors can occur
4. **COMPLETENESS GUARANTEE**: If a program has a valid typing, the type system MUST find it
5. **DECIDABILITY GUARANTEE**: Type inference MUST always terminate with definitive results

**🔥 STRUCTURAL TYPE REQUIREMENTS:**
- **Record Type Unification**: MUST use structural equivalence based on **FIELD NAMES ONLY**
- **Field Access**: MUST be **STRICTLY BY NAME** - never by position or ordering
- **Type Environment (Γ)**: MUST maintain consistent field name-to-type mappings
- **Substitution Application**: MUST apply substitutions consistently across all type expressions

**Hindley-Milner Algorithm Implementation Steps (MANDATORY):**
1. **Type Variable Generation**: Assign fresh type variables (α, β, γ) to untyped expressions
2. **Constraint Collection**: Gather type equality constraints from expression structure  
3. **Unification**: Solve constraints using Robinson's unification algorithm with occurs check
4. **Generalization**: Generalize types to introduce polymorphism at let-bindings
5. **Instantiation**: Create fresh instances of polymorphic types at usage sites

**Implementation References (REQUIRED READING):**
- **Robinson, J.A. (1965)**: "A Machine-Oriented Logic Based on the Resolution Principle" - Unification algorithm
- **Cardelli, L. (1987)**: "Basic Polymorphic Typechecking" - Implementation techniques  
- **Jones, M.P. (1995)**: "Functional Programming with Overloading and Higher-Order Polymorphism" - Advanced HM features

**🔥 COMPILER CORRECTNESS REQUIREMENT**: The implementation MUST pass all Hindley-Milner theoretical guarantees. Failure to implement proper HM inference is a **CRITICAL COMPILER BUG**.

---

Osprey's type system puts type safety and expressiveness as the top priorities. It is built upon the solid theoretical foundation of Hindley-Milner type inference, inspired by ML and Haskell. The type system aims towards making illegal states unrepresentable through complete static verification.

### 5.1 Built-in Types

**IMPORTANT**: All primitive types use lowercase names - `int`, `string`, `bool`. Capitalized forms (`Int`, `String`, `Bool`) are invalid.

- `int`: 64-bit signed integers
- `string`: UTF-8 encoded strings  
- `bool`: Boolean values (`true`, `false`)
- `unit`: Type for functions that don't return a meaningful value
- `Result<T, E>`: Built-in generic type for error handling
- `List<T, N>`: Immutable fixed-size lists with N elements of type T
- `Function Types`: First-class function types with syntax `(T1, T2, ...) -> R`
- `Record Types`: Immutable structured data types with named fields

#### Function Types

Function types represent functions as first-class values, enabling higher-order functions and function composition.

**Syntax:**
```
FunctionType := '(' (Type (',' Type)*)? ')' '->' Type
```

**Examples:**
```osprey
(int) -> int              // Function taking an int, returning an int
(int, string) -> bool     // Function taking int and string, returning bool
() -> string              // Function with no parameters, returning string
(string) -> (int) -> bool // Higher-order function returning another function
```

**Function Type Declarations:**
```osprey
// Function parameter with explicit function type
fn applyFunction(value: int, transform: (int) -> int) -> int = 
    transform(value)

// Variable with function type
let doubler: (int) -> int = fn(x: int) -> int = x * 2

// Higher-order function that returns a function
fn createAdder(n: int) -> (int) -> int = 
    fn(x: int) -> int = x + n
```

**Function Composition Examples:**
```osprey
// Define some simple functions
fn double(x: int) -> int = x * 2
fn square(x: int) -> int = x * x
fn addFive(x: int) -> int = x + 5

// Higher-order function with strong typing
fn applyTwice(value: int, func: (int) -> int) -> int = 
    func(func(value))

// Usage with type safety
let result1 = applyTwice(5, double)  // 20
let result2 = applyTwice(3, square)  // 81
let result3 = applyTwice(10, addFive) // 20

// Composition of functions
fn compose(f: (int) -> int, g: (int) -> int) -> (int) -> int =
    fn(x: int) -> int = f(g(x))

let doubleSquare = compose(double, square)
let result4 = doubleSquare(3) // double(square(3)) = double(9) = 18
```

**Type Safety Benefits:**
- **Compile-time validation**: Function signatures are checked at compile time
- **No runtime type errors**: Mismatched function types caught early
- **Clear documentation**: Function types serve as documentation
- **Enables optimization**: Compiler can optimize based on known function signatures

#### Record Types

Record types are immutable structured data types with named fields, providing structural equality semantics and compile-time field access validation.

**Syntax:**
```
RecordType := 'type' Identifier '=' '{' FieldList '}'
FieldList  := Field (',' Field)*
Field      := Identifier ':' Type
```

**Declaration Examples:**
```osprey
type Point = { x: int, y: int }
type Person = { name: string, age: int, active: bool }
type Address = { street: string, city: string, zipCode: string }
```

**Construction Syntax:**
```osprey
// Simple record construction
let point = Point { x: 10, y: 20 }
let person = Person { name: "Alice", age: 30, active: true }

// Field order is flexible during construction
let address = Address { 
    city: "Melbourne",
    street: "123 Main St", 
    zipCode: "3000"
}
```

**Field Access:**
```osprey
// Direct field access with dot notation
let x = point.x           // 10
let name = person.name    // "Alice"
let isActive = person.active  // true

// Field access in expressions
let distance = sqrt(point.x * point.x + point.y * point.y)
let greeting = "Hello, ${person.name}!"
```

**Key Properties:**
- **Immutability**: Records cannot be modified after creation
- **Structural Equality**: Two records are equal if all their fields are equal
- **Compile-time Field Validation**: Field access is validated at compile time
- **Type Safety**: Field types are enforced during construction and access
- **No Null Fields**: All fields must be provided during construction

**Pattern Matching with Records:**
Records support pattern matching for destructuring and value extraction (see [Pattern Matching](0008-PatternMatching.md) for complete details).

**Nested Records:**
```osprey
type Company = { name: string, address: Address }

let company = Company {
    name: "Tech Corp",
    address: Address {
        street: "456 Tech Ave",
        city: "Sydney", 
        zipCode: "2000"
    }
}

// Nested field access
let companyCity = company.address.city  // "Sydney"
```

**Record Updates (Non-destructive):**
Records support elegant non-destructive updates that create modified copies:

```osprey
// Original record
let person = Person { name: "Alice", age: 25, active: true }

// Non-destructive update (creates new instance)
let olderPerson = person { age: 26 }           // Only age changes
let renamedPerson = person { name: "Alicia" }  // Only name changes

// Multiple field updates
let updatedPerson = person { 
    age: 26, 
    active: false 
}

// Original person unchanged - all updates create new instances
print(person.age)        // Still 25
print(olderPerson.age)   // Now 26
```

**Record Type Constraints:**
Records can have validation constraints using `where` clauses:

```osprey
type ValidatedPerson = {
    name: string,
    age: int,
    email: string
} where validatePersonData

fn validatePersonData(person: ValidatedPerson) -> Result<ValidatedPerson, string> = 
    match person.age {
        age when age < 0 => Error("Age cannot be negative")
        age when age > 150 => Error("Age must be realistic")
        _ => match person.name {
            "" => Error("Name cannot be empty")
            _ => Success(person)
        }
    }

// Constrained record construction returns Result type
let result = ValidatedPerson { name: "Bob", age: 25, email: "bob@example.com" }
// result type: Result<ValidatedPerson, string>
// Must be handled with pattern matching (see Pattern Matching chapter)
```

**Field Access Rules:**

**🔥 CRITICAL SPECIFICATION: FIELD ACCESS IS STRICTLY BY NAME ONLY**

**ABSOLUTE REQUIREMENT**: Record field access is **EXCLUSIVELY BY NAME**. Field ordering, positioning, or indexing is **COMPLETELY FORBIDDEN** and must **NEVER** be relied upon by the compiler implementation.

**✅ ALLOWED - Field Access on Record Types (BY NAME ONLY):**
```osprey
type User = { id: int, name: string, email: string }
let user = User { id: 1, name: "Alice", email: "alice@example.com" }

let userId = user.id          // ✅ VALID: direct field access BY NAME
let userName = user.name      // ✅ VALID: direct field access BY NAME  
let userEmail = user.email    // ✅ VALID: direct field access BY NAME

// Field order during construction is IRRELEVANT
let user2 = User { 
    email: "bob@example.com",  // Different order - PERFECTLY VALID
    name: "Bob",               // Field position does NOT matter
    id: 2                      // Only field NAMES matter
}
let bobName = user2.name      // ✅ VALID: name-based access works regardless of declaration order
```

**❌ ABSOLUTELY FORBIDDEN - Positional or Indexed Access:**
```osprey
// NEVER ALLOWED - These are COMPILATION ERRORS
let value1 = user[0]          // ❌ FORBIDDEN: No indexed access  
let value2 = user.fields[1]   // ❌ FORBIDDEN: No positional access
let value3 = getFieldAt(user, 0)  // ❌ FORBIDDEN: No position-based access

// COMPILER IMPLEMENTATION MUST NEVER:
// - Rely on field declaration order for LLVM struct generation
// - Use field positioning for type unification
// - Access fields by index in any internal operation
// - Generate code that depends on field ordering
```

**🔥 COMPILER IMPLEMENTATION REQUIREMENT:**
The Osprey compiler **MUST** implement field access using **FIELD NAME LOOKUP ONLY**:
- ✅ Field-to-LLVM-index mapping by name
- ✅ Type unification based on field name matching  
- ✅ Pattern matching using field names
- ❌ **NEVER** field ordering dependencies
- ❌ **NEVER** positional field access in codegen
- ❌ **NEVER** field index assumptions

**❌ FORBIDDEN - Field Access on `any` Types:**
```osprey
fn processAnyValue(value: any) -> string = {
    // ERROR: Cannot access fields on 'any' type
    let result = value.name   // Compilation error
    return result
}

// CORRECT: Use pattern matching for 'any' types
fn processAnyValue(value: any) -> string = match value {
    person: { name } => person.name        // Extract field via pattern matching
    user: User { name } => name           // Type-specific pattern matching
    _ => "unknown"
}
```

**❌ FORBIDDEN - Field Access on Result Types:**
```osprey
type Person = { 
    name: string
} where validatePerson

fn validatePerson(person: Person) -> Result<Person, string> = match person.name {
    "" => Error("Name cannot be empty")
    _ => Success(person)
}

let personResult = Person { name: "Alice" }  // Returns Result<Person, string>

// ERROR: Cannot access fields on Result type
let name = personResult.name    // Compilation error

// CORRECT: Pattern match on Result first
let name = match personResult {
    Success { value } => value.name
    Error { message } => "error"
}
```

**❌ FORBIDDEN - Field Access on Union Types:**
```osprey
type Shape = Circle { radius: int } 
           | Rectangle { width: int, height: int }

let shape = Circle { radius: 5 }

// ERROR: Cannot access fields on union type
let radius = shape.radius     // Compilation error

// CORRECT: Pattern match on union type first
let area = match shape {
    Circle { radius } => 3.14 * radius * radius
    Rectangle { width, height } => width * height
}
```

**Compilation Errors:**
```osprey
// ERROR: Unknown field
let invalid = Point { x: 10, z: 30 }  // 'z' not defined in Point

// ERROR: Missing required field  
let incomplete = Person { name: "Alice" }  // Missing 'age' and 'active'

// ERROR: Type mismatch
let wrongType = Point { x: "ten", y: 20 }  // 'x' expects int, got string

// ERROR: Field access on non-record type
let num = 42
let invalid = num.x  // Cannot access field on non-record type

// ERROR: Cannot assign to fields (immutable)
let counter = Counter { value: 0 }
counter.value = 5    // Compilation error

// CORRECT: Create new instance with updated value
let newCounter = counter { value: 5 }
```

### 5.2 Built-in Error Types

- `MathError`: For arithmetic operations (DivisionByZero, Overflow, Underflow)
- `ParseError`: For string parsing operations  
- `IndexError`: For list/string indexing operations (OutOfBounds)
- `Success`: Successful result wrapper

### 5.3 Hindley-Milner Type Inference

**Core Implementation**: Osprey implements complete Hindley-Milner type inference (Hindley 1969, Milner 1978) enabling polymorphic type inference without explicit type annotations.

**Academic Foundation**:
- **Hindley, R. (1969)**: "The Principal Type-Scheme of an Object in Combinatory Logic" - Communications of the ACM 12(12):719-721
- **Milner, R. (1978)**: "A Theory of Type Polymorphism in Programming" - Journal of Computer and System Sciences 17:348-375
- **Damas, L. & Milner, R. (1982)**: "Principal type-schemes for functional programs" - POPL '82

**Implementation Principle**: Variables may be declared without type annotations when their types can be inferred through unification and constraint solving. The system performs automatic generalization and instantiation of polymorphic types.

**Hindley-Milner Algorithm Steps**:
1. **Type Variable Generation**: Assign fresh type variables to untyped expressions
2. **Constraint Collection**: Gather type equality constraints from expression structure
3. **Unification**: Solve constraints using Robinson's unification algorithm
4. **Generalization**: Generalize types to introduce polymorphism at let-bindings
5. **Instantiation**: Create fresh instances of polymorphic types at usage sites

#### Hindley-Milner Function Inference

**Complete Type Inference**: Both parameter and return types can be omitted when inferrable through Hindley-Milner constraint solving:

##### Hindley-Milner Inference Examples

**✅ POLYMORPHIC INFERENCE (No Type Annotations Required):**
```osprey
// Identity function - fully polymorphic
fn identity(x) = x              // Infers: <T>(T) -> T

// Arithmetic functions
fn add(a, b) = a + b           // Infers: (int, int) -> Result<int, MathError>
fn increment(x) = x + 1        // Infers: (int) -> Result<int, MathError>

// String operations  
fn concat(s1, s2) = s1 + s2    // Infers: (string, string) -> string

// Boolean operations
fn negate(x) = !x              // Infers: (bool) -> bool

// Field access polymorphism
fn getX(p) = p.x               // Infers: <T>(Point<T, _>) -> T
fn getValue(c) = c.value       // Infers: <T>(Container<T, _>) -> T

// Higher-order functions
fn apply(f, x) = f(x)          // Infers: <A, B>((A) -> B, A) -> B
fn compose(f, g) = fn(x) = f(g(x))  // Infers: <A, B, C>((B) -> C, (A) -> B) -> (A) -> C
```

**✅ MONOMORPHIC USAGE (Types Specialized at Call Sites):**
```osprey
// Same identity function used with different types
let intResult = identity(42)        // identity<int>
let stringResult = identity("test") // identity<string>
let boolResult = identity(true)     // identity<bool>

// Field accessors specialized by usage
let intPoint = Point { x: 10, y: 20 }
let stringPoint = Point { x: "a", y: "b" }

let intX = getX(intPoint)           // getX<int>
let stringX = getX(stringPoint)     // getX<string>
```

#### Constraint-Based Type Inference

**Unification Rules**: Osprey's Hindley-Milner implementation uses constraint unification to solve type equations:

**✅ CONSTRAINT SOLVING EXAMPLES:**
```osprey
// Function with multiple constraints
fn processData(item, transform) = transform(item.value)
// Constraints:
// - item must have field 'value' of type α
// - transform must be function (α) -> β  
// - return type is β
// Solution: <α, β>(Container<α>, (α) -> β) -> β

// Recursive constraint solving
fn chain(f, g, x) = f(g(x))
// Constraints:
// - g must be function (α) -> β
// - f must be function (β) -> γ
// - x has type α
// - return type is γ
// Solution: <α, β, γ>((β) -> γ, (α) -> β, α) -> γ
```

**✅ POLYMORPHIC GENERALIZATION:**
```osprey
// Generic data constructors
fn makePair(x, y) = Pair { first: x, second: y }
// Infers: <A, B>(A, B) -> Pair<A, B>

fn makePoint(a, b) = Point { x: a, y: b }
// Infers: <T>(T, T) -> Point<T, T>

// Generic extractors
fn getFirst(p) = p.first
// Infers: <A, B>(Pair<A, B>) -> A

fn getSecond(p) = p.second  
// Infers: <A, B>(Pair<A, B>) -> B
```

#### Hindley-Milner Implementation Examples

**✅ COMPLETE TYPE INFERENCE (All Types Derived):**
```osprey
// Polymorphic identity - no annotations needed
fn identity(x) = x              // <T>(T) -> T

// Arithmetic with constraint propagation
fn add(a, b) = a + b           // (int, int) -> Result<int, MathError>
fn multiply(x, y) = x * y      // (int, int) -> Result<int, MathError>

// String operations
fn greet(name) = "Hello, " + name  // (string) -> string

// Record construction and access
fn makeUser(n, a) = User { name: n, age: a }  // (string, int) -> User
fn getName(u) = u.name         // (User) -> string

// Higher-order functions
fn twice(f, x) = f(f(x))       // <T>((T) -> T, T) -> T
fn map(f, list) = [f(x) for x in list]  // <A, B>((A) -> B, List<A>) -> List<B>
```

**✅ MONOMORPHIZATION AT USAGE SITES:**
```osprey
// Same polymorphic functions, different instantiations
let id1 = identity(42)          // identity<int>
let id2 = identity("test")      // identity<string> 
let id3 = identity(true)        // identity<bool>

// Function used in different contexts
let intTwice = twice(increment, 5)      // twice<int>
let stringTwice = twice(greet, "World")  // twice<string>
```

**❌ INFERENCE LIMITATIONS (Explicit Types Required):**
```osprey
// Ambiguous conditional requires annotation
fn conditional(flag, a, b) -> T = if flag then a else b  // T must be specified

// External function calls may need hints
fn process(data) -> R = externalFunction(data)  // R may need annotation

// Complex constraints may require explicit polymorphic declaration
fn complex<T>(x: T, pred: (T) -> bool) -> Option<T> = 
    if pred(x) then Some(x) else None
```

#### Hindley-Milner Benefits

**Academic Guarantees (Milner 1978, Damas & Milner 1982)**:
1. **Principal Types**: Every well-typed expression has a unique most general type
2. **Completeness**: If a program has a type, Hindley-Milner will find it
3. **Soundness**: All inferred types are correct - no runtime type errors
4. **Decidability**: Type inference always terminates with definitive result

**Practical Benefits**:
- **Zero Annotation Burden**: Write polymorphic functions without type signatures
- **Maximum Reusability**: Functions automatically work with all compatible types
- **Compile-time Safety**: All type errors caught before execution
- **Performance**: Monomorphization enables optimal code generation

**Implementation References**:
- **Robinson, J.A. (1965)**: "A Machine-Oriented Logic Based on the Resolution Principle" - Unification algorithm
- **Cardelli, L. (1987)**: "Basic Polymorphic Typechecking" - Implementation techniques
- **Jones, M.P. (1995)**: "Functional Programming with Overloading and Higher-Order Polymorphism" - Advanced HM features

**Hindley-Milner Principle**: "Type annotations are optional for all expressions where types can be inferred through constraint unification. The system automatically finds the most general (polymorphic) type for each expression."

#### 🔥 CRITICAL: Record Type Structural Equivalence

**MANDATORY REQUIREMENT**: Osprey's Hindley-Milner implementation MUST treat record types using **structural equivalence based EXCLUSIVELY on field names**.

**✅ CORRECT Structural Unification:**
```osprey
// These record types are structurally equivalent (same field names and types)
type PersonA = { name: string, age: int }
type PersonB = { age: int, name: string }  // Different field ORDER - still equivalent

// Hindley-Milner MUST unify these as the same structural type
fn processA(p: PersonA) = p.name
fn processB(p: PersonB) = p.name

// These functions MUST be considered type-compatible
let result1 = processA(PersonB { age: 25, name: "Alice" })  // ✅ MUST work
let result2 = processB(PersonA { name: "Bob", age: 30 })    // ✅ MUST work
```

**✅ POLYMORPHIC Field Access Inference:**
```osprey
// Generic field accessor - inferred type based on field NAME only
fn getName(record) = record.name           // Infers: ∀α. {name: string, ...α} -> string
fn getAge(record) = record.age            // Infers: ∀α. {age: int, ...α} -> int

// Works with ANY record type that has the named field
let name1 = getName(Person { name: "Alice", age: 25 })      // ✅ Valid
let name2 = getName(User { name: "Bob", id: 1, email: "bob@example.com" })  // ✅ Valid
let age1 = getAge(Person { name: "Alice", age: 25 })       // ✅ Valid
```

**❌ FORBIDDEN Implementation Approaches:**
```
// NEVER ALLOWED in compiler implementation:
struct_field_0 = llvm_get_field_by_index(record, 0)  // ❌ Positional access
field_type = type_signature.params[field_position]   // ❌ Position-based type lookup
unify_by_field_order(record1, record2)              // ❌ Order-dependent unification
```

**🔥 UNIFICATION ALGORITHM REQUIREMENT:**
```
unify(RecordType1, RecordType2) := 
    if field_names(RecordType1) ≠ field_names(RecordType2) then FAIL
    else ∀ field_name ∈ field_names(RecordType1):
        unify(field_type(RecordType1, field_name), field_type(RecordType2, field_name))
        
// Field ordering is IRRELEVANT - only field names and their types matter
```

#### Polymorphic Type Variables vs Any Type

**CRITICAL DISTINCTION**: Hindley-Milner infers polymorphic type variables (α, β, γ), NOT the `any` type.

**✅ HINDLEY-MILNER POLYMORPHISM:**
```osprey
fn identity(x) = x                    // Infers: <T>(T) -> T (polymorphic)
fn getFirst(p) = p.first             // Infers: <A, B>(Pair<A, B>) -> A
fn apply(f, x) = f(x)                // Infers: <A, B>((A) -> B, A) -> B
```

**❌ ANY TYPE (Requires Explicit Declaration):**
```osprey
fn parseValue(input: string) -> any = processInput(input)  // Explicit any
fn getDynamicValue() -> any = readFromConfig()             // Explicit any
```

**Type Variable Instantiation**: Polymorphic type variables are instantiated to concrete types at usage sites:
```osprey
let intId = identity(42)        // T := int
let stringId = identity("test") // T := string
let boolId = identity(true)     // T := bool
```

**Safety Guarantee**: Polymorphic types are statically safe - all type checking occurs at compile time with no runtime type uncertainty.

#### Hindley-Milner Constraint Resolution

**Automatic Type Resolution**: The compiler uses constraint solving to resolve all type variables:

```osprey
// ✅ FULLY INFERRED: No annotations needed
fn greet(name) = "Hello, " + name    // String constraint from concatenation
fn formatScore(name, score) = "${name}: ${score}"  // String interpolation context
fn calculate(x, y) = x * y + 1       // Arithmetic constraints

// ✅ POLYMORPHIC RESOLUTION: Generic across multiple types
fn wrap(value) = Container { data: value }  // <T>(T) -> Container<T>
fn unwrap(container) = container.data       // <T>(Container<T>) -> T

// ✅ HIGHER-ORDER INFERENCE: Function parameters inferred
fn applyToList(func, items) = [func(x) for x in items]
// Infers: <A, B>((A) -> B, List<A>) -> List<B>
```

**Manual Annotation (When Desired)**: Explicit types can be added for documentation:
```osprey
fn greet(name: string) -> string = "Hello, " + name  // Explicit for clarity
fn identity<T>(x: T) -> T = x                        // Explicit polymorphism
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

**Pattern Matching Requirement:**
All `any` values must be accessed through pattern matching to extract their actual types (see [Pattern Matching](0008-PatternMatching.md) for complete syntax and examples).

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
Use pattern matching to safely extract values from `any` types (complete examples in [Pattern Matching](0008-PatternMatching.md)).

#### Function Return Type Handling

Functions returning `any` types require immediate pattern matching (see [Pattern Matching](0008-PatternMatching.md)).

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
match anyValue 
    value: Int => processInt(value)
    value: String => processString(value)
    // ERROR: missing wildcard or Bool case

// Exhaustive (CORRECT)
match anyValue 
    value: Int => processInt(value)
    value: String => processString(value)
    _ => handleOther()
```

#### Default Wildcard Behavior for Any Types

The wildcard pattern (`_`) in `any` type matching preserves the `any` type:

```osprey
// Wildcard returns any type
let result = match someAnyValue 
    value: Int => processInt(value)    // Returns specific type
    value: String => processString(value)  // Returns specific type
    _ => someAnyValue  // Returns any type (unchanged)
// result type: any (due to wildcard arm)

// To avoid any type in result, handle all expected cases explicitly
let result = match someAnyValue 
    value: Int => processInt(value)
    value: String => processString(value)
    _ => defaultInt()  // Convert to specific type
// result type: Int (all arms return Int)
```

#### Type Constraint Checking

The compiler **MUST** validate that pattern types are actually possible for the value being matched:

**✅ VALID - Realistic Type Patterns:**
```osprey
// Function known to return Int or String
extern fn parseIntOrString(input: string) -> any

match parseIntOrString("42") 
    value: Int => value + 1
    value: String => length(value)
    _ => 0  // Valid: handles any unexpected types
```

**❌ INVALID - Impossible Type Patterns:**
```osprey
// Function documented to only return Int or String
extern fn parseIntOrString(input: string) -> any

match parseIntOrString("42") 
    value: Int => value + 1
    value: String => length(value)
    value: Bool => if value then 1 else 0  // ERROR: Bool not possible
    _ => 0
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

#### Type-Level Validation

**CRITICAL**: When a type has a WHERE validation function, the type constructor **ALWAYS** returns `Result<T, String>` instead of the type directly.

**Validation Syntax:**
```
type_validation := 'where' function_name
```

**Examples:**
```osprey
// Unconstrained type - direct construction
type Point = { x: Int, y: Int }
let point = Point { x: 10, y: 20 }  // Returns Point directly

// Type with validation function - returns Result type
type Product = { 
    name: String,
    price: Int
} where validateProduct

fn validateProduct(product: Product) -> Result<Product, String> = match product.name {
    "" => Error("Product name cannot be empty")
    _ => match product.price {
        0 => Error("Price must be positive")
        _ => Success(product)
    }
}

let product = Product { name: "Widget", price: 100 }  // Returns Result<Product, String>

// Pattern matching required for validated types
match product {
    Success { value } => {
        print("Product: ${value.name}")
        print("Price: ${value.price}")
    }
    Error { message } => {
        print("Validation failed: ${message}")
    }
}

// Invalid validation fails at construction
let invalid = Product { name: "", price: -10 }  // Returns Error variant
match invalid {
    Success { value } => print("This won't execute")
    Error { message } => print("Expected error: ${message}")
}
```

**Field Access Rules:**
- **Unvalidated types**: Direct field access allowed (`point.x`)
- **Validated types**: Field access **ONLY** after pattern matching on Result

**Compilation Errors:**
```osprey
// ERROR: Cannot access field on Result type
let name = product.name  // Compilation error - pattern matching required

// CORRECT: Field access after unwrapping
match product {
    Success { value } => print("Name: ${value.name}")  // Valid field access
    Error { message } => print("Error: ${message}")
}
```

### 5.6 Type Compatibility

- Pattern matching for type discrimination
- Union types for representing alternatives
- Result types for error handling instead of exceptions