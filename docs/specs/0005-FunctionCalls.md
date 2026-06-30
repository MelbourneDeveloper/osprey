# Function Calls

> **Flavor layer — surface (CST).**  The call spellings in this chapter — named-argument `f(x: a, y: b)`, positional `f(x)`, and `f()` — are the **Default** surface (`.osp`); every form lowers to a single canonical node, `Expr::Call { function, arguments, named_arguments }`. The **ML** flavor (`.ospml`) writes the same calls as whitespace application `f a b` ([FLAVOR-ML-CALL](0024-MLFlavorSyntax.md)), which lowers to the *same* `Expr::Call`. The arity, named-argument, and saturation rules below are **shared-core** — the type checker enforces them on the AST and is flavor-blind. The one honest surface difference is currying ([FLAVOR-CURRY](0023-LanguageFlavors.md#currying-canonicalisation)). See [Language Flavors](0023-LanguageFlavors.md) and [ML Flavor Syntax](0024-MLFlavorSyntax.md).

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
// (Default surface; lowers to Expr::Call. ML spells the saturated form `add 10 20`.)
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

## Rules

1. Zero parameters: empty parentheses, `f()`.
2. One parameter: positional or named.
3. Two or more parameters: every argument must be named. Mixing positional and named arguments is a compilation error.
4. **Built-in functions** ([Built-in Functions](0012-Built-InFunctions.md)) are exempt: they take positional arguments in subject-first order — `split("a,b,c", ",")`, `fold(xs, 0, add)` — so the pipe can supply the subject as the first argument: `xs |> fold(0, add)`.

Argument order at the call site is independent of declaration order; the compiler reorders by name.

These rules are enforced on the canonical `Expr::Call` after lowering, so they hold identically regardless of flavor: ML whitespace application `add 10 20` produces the same node and is checked the same way. The single surface-visible difference is currying — see [FLAVOR-CURRY](0023-LanguageFlavors.md#currying-canonicalisation) and [ML Flavor Syntax](0024-MLFlavorSyntax.md).

## Cross-references

- [Language Flavors](0023-LanguageFlavors.md)
- [ML Flavor Syntax](0024-MLFlavorSyntax.md)