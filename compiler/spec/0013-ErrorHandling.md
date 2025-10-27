14. [Error Handling](0015-ErrorHandling.md)
    - [The Result Type](#the-result-type)

## Error Handling
### The Result Type

**CRITICAL**: All functions that can fail **MUST** return a `Result` type. There are no exceptions, panics, or nulls. This is a core design principle of the language to ensure safety and eliminate entire classes of runtime errors.

The `Result` type is a generic union type with two variants:

- `Success { value: T }`: Represents a successful result. Success **IS A STRUCT** with a `value` field of type `T`.
- `Error E`: Represents an error. Error **CONTAINS A DIRECT VALUE** of type `E` (not a struct).

**Example:**
```osprey
type Result<T, E> = Success { value: T } | Error E
```

The compiler **MUST** enforce that `Result` types are always handled with a `match` expression, preventing direct access to the underlying value and ensuring that all possible outcomes are considered.

**Pattern Matching Results:**
```osprey
let result = someFunctionThatCanFail()

match result {
    Success value => print("Success: ${value}")
    Error message => print("Error: ${message}")
}
```

This approach guarantees that error handling is explicit, robust, and checked at compile time.

### Arithmetic Safety and Result Types

**🚨 CRITICAL DESIGN PRINCIPLE 🚨**: ALL arithmetic operations (`+`, `-`, `*`, `/`, `%`) return `Result<T, MathError>` because they can **ALL** fail (overflow, underflow, division by zero).

## **THE GOLDEN RULE:**

**ALL Arithmetic Operations Return Result Types:**
- `+` Addition: `(int, int) -> Result<int, MathError>` - Can overflow
- `-` Subtraction: `(int, int) -> Result<int, MathError>` - Can underflow
- `*` Multiplication: `(int, int) -> Result<int, MathError>` - Can overflow
- `/` Division: `(int, int) -> Result<float, MathError>` - Can divide by zero, ALWAYS returns float
- `%` Modulo: `(int, int) -> Result<int, MathError>` - Can divide by zero

**Why ALL operations return Result:**
- **Addition/Subtraction**: Can overflow/underflow (e.g., MAX_INT + 1)
- **Multiplication**: Can overflow (e.g., 1000000 * 1000000)
- **Division**: Can divide by zero, ALWAYS returns float for mathematical correctness
- **Modulo**: Can divide by zero

**Type Promotion Rules:**
- `int ⊕ int` → `Result<int, MathError>` (where ⊕ is +, -, *, %)
- `float ⊕ float` → `Result<float, MathError>`
- `int ⊕ float` → `Result<float, MathError>` (int promoted to float)
- `int / int` → `Result<float, MathError>` (division ALWAYS returns float!)

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

When converting a `Result` type to a string using `toString()`, the format **MUST ALWAYS** be:
- **`Success(value)`**: For successful results
- **`Error(message)`**: For error results

**Examples:**
```osprey
let divisionResult = 15 / 3              // Result<float, MathError>
print(toString(divisionResult))          // Outputs: "Success(5)"

let divisionByZero = 10 / 0              // Result<float, MathError>
print(toString(divisionByZero))          // Outputs: "Error(DivisionByZero)"

let calculation = 10 + 5                 // Result<int, MathError>
print(toString(calculation))             // Outputs: "Success(15)"
```

**ABSOLUTE RULES:**
- ✅ **ALWAYS** wrap values in `Success(...)` or `Error(...)`
- ❌ **NEVER** output raw values without the wrapper
- ❌ **NEVER** use different formats for different Result types

This ensures consistent, predictable string representations of Result types across the entire language.