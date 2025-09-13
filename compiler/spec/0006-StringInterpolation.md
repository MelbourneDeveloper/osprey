7. [String Interpolation](0007-StringInterpolation.md)
   - [Syntax](#syntax)
   - [Expression Support](#expression-support)
   - [Type Handling](#type-handling)
   - [Implementation](#implementation)

# String Interpolation

✅ **IMPLEMENTED**: String interpolation is fully implemented and working with comprehensive test coverage.

## Syntax

String interpolation uses `${}` syntax within double-quoted strings:

```osprey
let name = "Alice"
let age = 30
let message = "Hello ${name}, you are ${age} years old"
```

## Expression Support

Any expression is valid inside interpolation:

```osprey
let x = 10
let y = 5
print("Sum: ${x + y}")
print("Product: ${x * y}")
print("Complex: ${(x + y) * 2 - 1}")
```

## Type Handling

- **String variables**: Use `%s` format specifier
- **Integer expressions**: Use `%ld` format specifier  
- **Function calls**: Supported for single-parameter functions

## Implementation

Interpolated strings are compiled to:
1. Allocate a buffer (`alloca [1024 x i8]`)
2. Use `sprintf` to format the string
3. Use `puts` to output the result