# Function Calls

## Named Arguments Requirement

Functions with more than one parameter must be called with named arguments.

### Valid Function Calls

```osprey
// Zero parameters
fn getValue() = 42
let value = getValue()

// Single parameter - positional allowed
fn double(x) = x * 2
let result = double(5)

// Multiple parameters - named arguments required
fn add(x, y) = x + y
let sum = add(x: 10, y: 20)

// Order doesn't matter with named arguments
let sum2 = add(y: 20, x: 10)

// Works with type annotations
fn multiply(a: int, b: int) -> int = a * b
let product = multiply(a: 5, b: 3)
```

### Invalid Function Calls

```osprey
// ERROR: Multi-parameter function with positional arguments
fn add(x, y) = x + y
let sum = add(10, 20)  // Compilation error

// ERROR: Mixed positional and named arguments
let sum = add(10, y: 20)  // Compilation error

// ERROR: Missing parameter name
let result = multiply(5, b: 3)  // Compilation error
```

## Compilation Rules

1. **Zero parameters**: Called with empty parentheses `()`
2. **Single parameter**: May use positional or named argument
3. **Multiple parameters**: All arguments must be named
4. **Argument order**: Named arguments are reordered to match parameter declaration order during compilation

## Rationale

Named arguments improve readability and prevent argument order errors in multi-parameter functions:

```osprey
// Clear intent with named arguments
httpGet(clientID: client, path: "/users", headers: "")

// Unclear with positional arguments (forbidden)
httpGet(client, "/users", "")  // What does "" mean?
```