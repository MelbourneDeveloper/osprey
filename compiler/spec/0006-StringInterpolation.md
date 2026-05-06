# String Interpolation

String interpolation provides convenient inline expression evaluation within string literals.

## Syntax

String interpolation uses `${}` syntax:

```osprey
let name = "Alice"
let age = 30
let message = "Hello ${name}, you are ${age} years old"
```

## Expression Support

Any expression can be interpolated:

```osprey
let x = 10
let y = 5
print("Sum: ${x + y}")
print("Product: ${x * y}")
print("Complex: ${(x + y) * 2 - 1}")

// Function calls
fn double(n) = n * 2
print("Doubled: ${double(5)}")

// Field access
type Person = { name: string, age: int }
let person = Person { name: "Bob", age: 25 }
print("Person: ${person.name}, age ${person.age}")
```

## Type Handling

Interpolated expressions are automatically converted to strings:

- **Primitive types**: int, float, bool converted directly
- **String types**: Inserted as-is
- **Result types**: Must be pattern-matched before interpolation
- **Complex types**: Use `toString()` for explicit conversion

```osprey
let num = 42
let flag = true
print("Number: ${num}, Flag: ${flag}")

// Result types require unwrapping
let result = 10 + 5
match result {
    Success { value } => print("Result: ${value}")
    Error { message } => print("Error: ${message}")
}
```

## Escaping

Use backslash to escape special characters:

```osprey
let literal = "Dollar sign: \${not interpolated}"
let newline = "Line 1\nLine 2"
let quote = "He said \"Hello\""
let backslash = "Path: C:\\Users\\Name"
```

Supported escape sequences:
- `\n` - Newline
- `\t` - Tab
- `\r` - Carriage return
- `\\` - Backslash
- `\"` - Double quote
- `\${` - Literal `${` (prevents interpolation)

## Implementation

Interpolated strings compile to efficient buffer operations:
1. Allocate a buffer (`alloca [1024 x i8]`)
2. Use `sprintf` with appropriate format specifiers (`%s`, `%ld`, etc.)
3. Return the formatted string for use in expressions