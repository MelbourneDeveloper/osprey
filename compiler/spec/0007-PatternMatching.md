# Pattern Matching

`match` is the only branching construct in Osprey. Record patterns are matched structurally by field name, not by field order. See [Type System](0004-TypeSystem.md) for type unification rules.

### Basic Patterns

```osprey
let result = match value {
    0 => "zero"
    1 => "one"
    n => "other: " + toString(n)
}
```

## Union Type Patterns

```osprey
type Option = Some { value: Int } | None

let message = match option {
    Some x => "Value: " + toString(x.value)
    None => "No value"
}
```

## Wildcard Patterns

The underscore `_` matches any value:

```osprey
let category = match score {
    100 => "perfect"
    90 => "excellent"
    _ => "good"
}
```

## Type Annotation Patterns

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

## Result Patterns

`Result<T, E>` is matched the same way as any other union. See [Error Handling](0013-ErrorHandling.md) for the type and arithmetic semantics.

```osprey
let calculation = 1 + 3 + (300 / 5)  // Result<int, MathError>

match calculation {
    Success { value }   => print("Result: ${value}")
    Error   { message } => print("Math error: ${message}")
}
```

Compound arithmetic expressions yield a single `Result`, not nested `Result`s; the compiler unwraps intermediate values inside the chain. Only the final value needs to be matched.

## Ternary Match (Syntactic Sugar)

A two-arm match has a shorthand. Two equivalent forms exist:

```ebnf
ternary ::= expr "{" pattern "}" "?" expr ":" expr   (* structural form *)
          | expr "?:" expr                            (* Result default form *)
```

Structural form — pick out a field, fall back if the pattern fails:

```osprey
let calculation = 10 + 5
let value = calculation { value } ? value : -1   // 15
```

Desugars to:

```osprey
match calculation {
    { value } => value
    _         => -1
}
```

Result-default form — extract `Success { value }` or use the default on `Error`:

```osprey
let safeValue = divide(a: 10, b: 2) ?: -1   // 5
let errorVal  = divide(a: 10, b: 0) ?: -1   // -1
```

A boolean expression with `?:` works because `true`/`false` desugar to the same match:

```osprey
let status = isActive ? "Active" : "Inactive"
```