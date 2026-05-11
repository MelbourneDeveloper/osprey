# Type System

- [Hindley-Milner Inference](#hindley-milner-inference)
- [Built-in Types](#built-in-types)
- [Result Auto-Unwrapping](#result-auto-unwrapping)
- [Function Types](#function-types)
- [Record Types](#record-types)
- [Validated Records (`where`)](#validated-records-where)
- [Collection Types](#collection-types)
- [Built-in Error Types](#built-in-error-types)
- [The `any` Type](#the-any-type)
- [Type Annotation Requirements](#type-annotation-requirements)

## Hindley-Milner Inference

Osprey uses Hindley-Milner inference. Every well-typed expression has a unique most general type, inference always terminates, and a successful type-check guarantees no runtime type errors.

Type annotations are optional everywhere they can be inferred:

```osprey
fn identity(x)         = x                       // <T>(T) -> T
fn add(a, b)           = a + b                   // (int, int) -> Result<int, MathError>
fn greet(name)         = "Hello, " + name        // (string) -> string
fn makeUser(n, a)      = User { name: n, age: a }  // (string, int) -> User
fn getName(u)          = u.name                  // (User) -> string
fn twice(f, x)         = f(f(x))                 // <T>((T) -> T, T) -> T
fn compose(f, g)       = fn(x) => f(g(x))        // <A,B,C>((B)->C,(A)->B) -> (A)->C
```

A polymorphic function is monomorphised independently at each call site:

```osprey
let i = identity(42)        // identity<int>
let s = identity("hello")   // identity<string>
```

### Record Type Unification

Two record types unify iff they have the same set of field names and corresponding field types unify. Field order is irrelevant in both declaration and construction.

```
unify(R1, R2) :=
    if names(R1) ≠ names(R2) then FAIL
    else for each f ∈ names(R1): unify(typeOf(R1, f), typeOf(R2, f))
```

### Polymorphic Variables vs `any`

Inference produces polymorphic variables (`<T>`, `<A>`, …), not `any`. The `any` type is opt-in; see [The `any` Type](#the-any-type).

## Built-in Types

All primitive types are lowercase.

| Type             | Description                                                        |
| ---------------- | ------------------------------------------------------------------ |
| `int`            | 64-bit signed integer (LLVM `i64`)                                 |
| `float`          | 64-bit IEEE 754 (LLVM `double`)                                    |
| `string`         | UTF-8 encoded                                                      |
| `bool`           | `true` \| `false`                                                  |
| `unit`           | The single value `()`; the return type of a function with no result|
| `any`            | Erased value; access requires pattern matching                     |
| `Result<T, E>`   | Error-handling sum type (see [Error Handling](0013-ErrorHandling.md)) |
| `List<T>`        | Immutable sequential collection                                    |
| `Map<K, V>`      | Immutable key/value collection                                     |

`int` and `float` do not implicitly convert. Use `toFloat(int)` and `toInt(float)`.

Arithmetic returns `Result<T, MathError>`. The full operator-by-operand table and chaining rules are in [Error Handling](0013-ErrorHandling.md#arithmetic-returns-result).

## Result Auto-Unwrapping

A bare arithmetic expression has type `Result<T, MathError>`. The compiler auto-unwraps the inner `Result` in three contexts so authors do not write nested `match` chains:

1. **Nested arithmetic.** `(10 + 5) * 2` has type `Result<int, MathError>`, not `Result<Result<int, MathError>, MathError>`. If any sub-expression errors, the chain errors.
2. **User function arguments.** Passing a `Result`-typed expression to a function that expects the underlying type unwraps it. `double(add(a: 5, b: 3))` is well-typed when `add` returns `Result<int, MathError>` and `double` expects `int`.
3. **Fiber operations.** `spawn`, `await`, `send`, and `recv` unwrap `Result` arguments before storing them.

Auto-unwrap does **not** apply to:

- `toString`. `toString(add(a: 5, b: 3))` produces `"Success(8)"`, never `"8"`. Use `toString` to inspect the `Result` itself.
- A function's declared return type. `fn add(x, y) = x + y` returns `Result<int, MathError>`; `fn compute() -> int = 5` returns `int`.

The top-level value of an arithmetic expression must still be either matched or stored as a `Result`.

## Function Types

```ebnf
functionType ::= "(" (type ("," type)*)? ")" "->" type
```

```osprey
(int) -> int
(int, string) -> bool
() -> string
(string) -> (int) -> bool          // higher-order
```

```osprey
fn applyFunction(value: int, transform: (int) -> int) -> int = transform(value)

let doubler: (int) -> int = fn(x: int) => x * 2

fn createAdder(n: int) -> (int) -> int = fn(x: int) => x + n
```

Multi-argument call syntax (named arguments are required for two or more parameters) is in [Function Calls](0005-FunctionCalls.md).

## Record Types

```ebnf
recordType ::= "type" ID "=" "{" field ("," field)* "}" constraint?
field      ::= ID ":" type
constraint ::= "where" function_name
```

```osprey
type Point   = { x: int, y: int }
type Person  = { name: string, age: int, active: bool }
```

### Construction

```osprey
let point  = Point  { x: 10, y: 20 }
let person = Person { name: "Alice", age: 30, active: true }

// Field order at construction is irrelevant
let person2 = Person { active: true, name: "Bob", age: 22 }
```

All fields are required. Missing or unknown fields, or type mismatches, are compilation errors.

### Field Access

Direct field access is permitted only on a record value. Field access on `any`, `Result`, or any union type requires `match` to narrow the value first.

```osprey
let n = person.name        // ok

// any: pattern-match
fn nameOf(v: any) -> string = match v {
    p: { name } => p.name
    _           => "unknown"
}

// Result: match before access
match personResult {
    Success { value }   => print(value.name)
    Error   { message } => print(message)
}

// Union: discriminate first
let area = match shape {
    Circle    { radius }         => 3.14 * radius * radius
    Rectangle { width, height }  => width * height
}
```

The compiler implementation must look up fields by name; positional access is forbidden in code generation.

### Immutability and Non-Destructive Update

Records cannot be modified. To produce a record that differs in some fields from an existing one, use the update form:

```osprey
let p2 = point  { x: 15 }                // y carried over
let p3 = person { age: 26, active: false }
```

### Nested Records

```osprey
type Address = { street: string, city: string, zipCode: string }
type Company = { name: string, address: Address }

let company = Company {
    name:    "Tech Corp",
    address: Address { street: "456 Tech Ave", city: "Sydney", zipCode: "2000" }
}

let companyCity = company.address.city
```

## Validated Records (`where`)

A `where` clause attaches a validation function to a record type. The constructor of a validated type returns `Result<T, string>` instead of `T`, and the validation function runs at construction time.

```osprey
type Product = {
    name:  string,
    price: int
} where validateProduct

fn validateProduct(p: Product) -> Result<Product, string> = match p.name {
    ""  => Error   { message: "name cannot be empty" }
    _   => match p.price {
        0 => Error   { message: "price must be positive" }
        _ => Success { value:   p }
    }
}

// Construction returns Result<Product, string>
let r = Product { name: "Widget", price: 100 }
match r {
    Success { value }   => print("ok: ${value.name}")
    Error   { message } => print("validation failed: ${message}")
}
```

Field access on a validated value is only legal after matching on the `Result`.

## Collection Types

### `List<T>`

Immutable, homogeneous, with structural sharing. Index access is bounds-checked and returns `Result<T, IndexError>`.

```osprey
let numbers = [1, 2, 3, 4, 5]            // List<int>
let names   = ["Alice", "Bob"]           // List<string>

// Empty literal cannot infer its element type unless the context provides it
let empty: List<int> = []                // ok
let total = sumOfInts([])                // ok if sumOfInts: (List<int>) -> int

match numbers[0] {
    Success { value }   => print(value)
    Error   { message } => print(message)
}
```

Operations (see [Iterators and Iteration](0010-LoopConstructsAndFunctionalIterators.md) for full signatures):

```osprey
let doubled  = numbers |> map(x => x * 2)
let evens    = numbers |> filter(x => x % 2 == 0)
let total    = numbers |> fold(initial: 0, function: (acc, x) => acc + x)
let combined = numbers + [6, 7, 8]                       // concatenation produces a new list
numbers |> forEach(x => print(toString(x)))
```

Patterns:

```osprey
fn classify(xs: List<int>) -> string = match xs {
    []                 => "empty"
    [single]           => "one"
    [first, second]    => "two"
    [head, ...tail]    => "many starting with ${head}"
}
```

List comprehensions and head/tail destructuring:

```osprey
let squares  = [x * x for x in range(start: 1, end: 5)]   // [1, 4, 9, 16, 25]
let filtered = [x for x in numbers if x > 3]
let [head, ...tail] = numbers
```

### `Map<K, V>`

Immutable, persistent. Lookup, insert, and remove are O(log n) and return new maps.

```osprey
let ages = {
    "Alice":   25,
    "Bob":     30,
    "Charlie": 35
}                                                 // Map<string, int>

let scores: Map<string, int> = {}                 // empty needs annotation

match ages["Alice"] {
    Success { value }   => print(toString(value))
    Error   { message } => print(message)
}

let hasAlice = contains(map: ages, key: "Alice")
let n        = length(ages)
```

Operations:

```osprey
let bumped     = ages |> mapValues(age => age + 1)
let upper      = ages |> mapKeys(name => toUpperCase(name))
let thirties   = ages |> filter((name, age) => age >= 30)
let totalAge   = ages |> fold(initial: 0, function: (acc, name, age) => acc + age)
let withDave   = ages + { "Dave": 28 }
let updated    = ages { "Alice": 26 }                       // single-key update
let withoutBob = removeKey(map: ages, key: "Bob")
```

Patterns:

```osprey
fn analyze(p: Map<string, int>) -> string = match p {
    {}                                       => "none"
    { "Alice": age }                         => "only Alice (${age})"
    { "Alice": _, "Bob": _ }                 => "Alice and Bob"
    p when length(p) > 5                     => "large"
    _                                        => "other"
}
```

Conversions:

```osprey
let names    = keys(ages)                                       // List<string>
let agesList = values(ages)                                     // List<int>
let pairs    = entries(ages)                                    // List<(string, int)>
let m        = zipToMap(keys: names, values: agesList)
let byGrade  = groupBy(items: students, function: s => s.grade) // Map<string, List<Student>>
```

### Performance

`List<T>`: O(1) indexed access, O(n) concatenation and bulk transforms. `Map<K, V>`: O(log n) lookup/insert/update, O(n) iteration. Both use structural sharing; small collections may be stack-allocated.

## Built-in Error Types

| Type          | Used by                                              |
| ------------- | ---------------------------------------------------- |
| `MathError`   | Arithmetic (`DivisionByZero`, `Overflow`, `Underflow`)|
| `ParseError`  | String parsing                                       |
| `IndexError`  | List and string indexing (`OutOfBounds`)             |
| `StringError` | String operations that can fail (`length`, `substring`, `contains`) |
| `ChannelError`| Channel send/recv                                    |

`Success` and `Error` are the constructors of `Result<T, E>` (see [Error Handling](0013-ErrorHandling.md)).

## The `any` Type

`any` is an erased type. It exists so a function may receive a value whose type is not known statically — for example, parsed JSON or a foreign-function return value. Direct use of an `any` value is forbidden; the value must be narrowed by `match` first.

### Forbidden Operations

```osprey
fn processAny(v: any) -> int = v + 1                       // ❌ direct arithmetic
fn getLength(v: any) -> int = v.length                     // ❌ direct field access
let n: int = someAnyFunction()                             // ❌ implicit conversion
fn callIt(v: any) = someFunction(v)                        // ❌ pass to typed parameter
let s = toString(v)              // where v: any            // ❌ implicit conversion
```

Each of these produces a compilation error.

### Required Form

```osprey
fn process(v: any) -> int = match v {
    n: int    => n + 1
    s: string => length(s)
    _         => 0
}
```

Pattern syntax (`name: Type`, `name: { fields }`, `{ fields }`, `_`) is defined in [Pattern Matching](0007-PatternMatching.md).

### Exhaustiveness

A `match` on `any` must either cover every type the value may take or include a wildcard `_`:

```osprey
match v {
    n: int    => processInt(n)
    s: string => processString(s)
    _         => handleOther()
}
```

A wildcard arm that returns the matched value preserves the `any` type; to escape `any`, every arm (including the wildcard) must return the same concrete type.

### Compiler Errors

The compiler emits the following message strings on `any`-type misuse:

```
cannot use 'any' type directly in arithmetic operation - pattern matching required
cannot access field on 'any' type without pattern matching
cannot assign 'any' to 'TYPE' without pattern matching
cannot pass 'any' type to function expecting 'TYPE' - pattern matching required
cannot implicitly convert 'any' to 'TYPE' - use pattern matching to extract specific type
cannot access variable of type 'any' directly - pattern matching required
pattern matching on 'any' type must handle all possible types or include wildcard
```

### Documented Type Sets

When the compiler has information about which types an `any`-returning function may produce — for example, from an `extern` declaration or annotation — it rejects patterns that match impossible types and unreachable patterns:

```
pattern 'TYPE' is not a possible type for expression of documented types [TYPE1, TYPE2, ...]
unreachable pattern: 'TYPE' cannot occur based on context analysis
pattern matching includes impossible type 'TYPE' - check function documentation
```

## Type Annotation Requirements

Annotations are required when the compiler cannot infer a type:

- An empty literal (`[]`, `{}`) with no contextual type.
- A function whose return type is ambiguous (for example, a value returned from an `extern`).
- A polymorphic function whose type variables are not constrained by any argument.

```osprey
let xs: List<int> = []
fn parseValue<T>(input: string) -> Result<T, ParseError> = ...
```

The compiler emits an error in each case where inference is ambiguous; explicit annotations resolve them.
