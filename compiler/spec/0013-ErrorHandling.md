14. [Error Handling](0015-ErrorHandling.md)
    - [The Result Type](#the-result-type)

## Error Handling
### The Result Type

**CRITICAL**: All functions that can fail **MUST** return a `Result` type. There are no exceptions, panics, or nulls. This is a core design principle of the language to ensure safety and eliminate entire classes of runtime errors.

The `Result` type is a generic union type with two variants:

- `Success { value: T }`: Represents a successful result, containing the value of type `T`.
- `Error { message: E }`: Represents an error, containing an error message or object of type `E`.

**Example:**
```osprey
type Result<T, E> = Success { value: T } | Error { message: E }
```

The compiler **MUST** enforce that `Result` types are always handled with a `match` expression, preventing direct access to the underlying value and ensuring that all possible outcomes are considered.

**Pattern Matching Results:**
```osprey
let result = someFunctionThatCanFail()

match result {
    Success { value } => print("Success: ${value}")
    Error { message } => print("Error: ${message}")
}
```

This approach guarantees that error handling is explicit, robust, and checked at compile time.

### Arithmetic Operations and Result Types

All arithmetic operations return `Result<T, MathError>` to handle overflow, underflow, and division by zero:

**Operator Signatures:**
- `+` Addition: `(int, int) -> Result<int, MathError>`
- `-` Subtraction: `(int, int) -> Result<int, MathError>`
- `*` Multiplication: `(int, int) -> Result<int, MathError>`
- `/` Division: `(int, int) -> Result<float, MathError>` — always returns float
- `%` Modulo: `(int, int) -> Result<int, MathError>`

**Type Promotion:**
- `int ⊕ int` → `Result<int, MathError>` (where ⊕ is +, -, *, %)
- `float ⊕ float` → `Result<float, MathError>`
- `int ⊕ float` → `Result<float, MathError>` (int promoted to float)
- `int / int` → `Result<float, MathError>` (division always returns float)

#### Arithmetic Examples

```osprey
// ALL arithmetic returns Result types
let sum = 1 + 3                    // Result<int, MathError> - could overflow
let diff = 10 - 5                  // Result<int, MathError> - could underflow
let product = 2 * 4                // Result<int, MathError> - could overflow
let quotient = 10 / 2              // Result<float, MathError> - could divide by zero
let remainder = 10 % 3             // Result<int, MathError> - could divide by zero

// Float arithmetic also returns Result
let fsum = 3.14 + 2.86             // Result<float, MathError>
let fproduct = 2.5 * 4.0           // Result<float, MathError>

// Type promotion
let mixed = 10 + 5.5               // Result<float, MathError> (int promoted)
```

#### Working with Arithmetic Results

```osprey
// Must pattern match to extract values
let calculation = 10 + 5           // Result<int, MathError>

match calculation {
    Success { value } => print("Result: ${value}")
    Error { message } => print("Error: ${message}")
}

// Chaining arithmetic requires nested matches or Result operators
let step1 = 10 + 5                 // Result<int, MathError>
match step1 {
    Success { val1 } => {
        let step2 = val1 * 2       // Result<int, MathError>
        match step2 {
            Success { val2 } => print("Final: ${val2}")
            Error { message } => print("Multiplication error: ${message}")
        }
    }
    Error { message } => print("Addition error: ${message}")
}

// Printing Result types directly
print(10 + 5)                      // Outputs: Success(15)
print(10 / 0)                      // Outputs: Error(DivisionByZero)
```

### Result Type toString Format

Result types convert to strings in the format `Success(value)` or `Error(message)`:

```osprey
let divisionResult = 15 / 3              // Result<float, MathError>
print(toString(divisionResult))          // "Success(5)"

let divisionByZero = 10 / 0              // Result<float, MathError>
print(toString(divisionByZero))          // "Error(DivisionByZero)"

let calculation = 10 + 5                 // Result<int, MathError>
print(toString(calculation))             // "Success(15)"
```

This format is consistent across all Result types.