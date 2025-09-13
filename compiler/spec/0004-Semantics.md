# Semantics

- [Variable Binding](#variable-binding)
- [Function Semantics](#function-semantics)
- [Evaluation Order](#evaluation-order)

## Variable Binding

- `let` creates immutable bindings
- `mut` creates mutable bindings
- Variables must be initialized at declaration
- Shadowing is allowed in nested scopes

## Function Semantics

- Functions are first-class values
- All functions are pure (no side effects except I/O)
- Recursive functions are supported
- Tail recursion is optimized

## Evaluation Order

- Expressions are evaluated left-to-right
- Function arguments are evaluated before the function call
- Short-circuit evaluation for logical operators