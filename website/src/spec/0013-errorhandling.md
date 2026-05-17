---
layout: page
title: "Error Handling"
description: "Osprey Language Specification: Error Handling"
date: 2026-05-17
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0013-errorhandling/"
---

# Error Handling

Osprey has no exceptions, panics, or null. Any function that can fail returns a `Result`.

## The Result Type

```osprey
type Result<T, E> = Success { value: T } | Error { message: E }
```

The compiler rejects any direct access to the contained value. Callers must pattern-match the `Result` (see [Pattern Matching](0007-PatternMatching.md)):

```osprey
let result = someFunctionThatCanFail()

match result {
    Success { value }   => print("Success: ${value}")
    Error   { message } => print("Error: ${message}")
}
```

## Arithmetic Returns Result

Every arithmetic operation returns `Result<T, MathError>` so overflow, underflow, and division by zero surface as values, not panics.

| Operator    | int, int                   | float, float               | int, float / float, int                   |
| ----------- | -------------------------- | -------------------------- | ----------------------------------------- |
| `+ - * %`   | `Result<int,   MathError>` | `Result<float, MathError>` | `Result<float, MathError>` (int promoted) |
| `/`         | `Result<float, MathError>` | `Result<float, MathError>` | `Result<float, MathError>`                |

`/` always yields `float`. There is no implicit `int`/`float` conversion outside this table; use `toFloat` and `toInt` for explicit conversion.

```osprey
let sum       = 1 + 3      // Result<int,   MathError>
let quotient  = 10 / 3     // Result<float, MathError>
let remainder = 10 % 3     // Result<int,   MathError>
let mixed     = 10 + 5.5   // Result<float, MathError>
let divZero   = 10 / 0     // Error(DivisionByZero)
```

#### Chaining Arithmetic

Each operation returns a `Result`, so chaining requires either nested matches or, in the future, Result-aware operators:

```osprey
let step1 = 10 + 5
match step1 {
    Success { val1 } => match val1 * 2 {
        Success { val2 }    => print("Final: ${val2}")
        Error   { message } => print("Multiplication error: ${message}")
    }
    Error { message } => print("Addition error: ${message}")
}
```

### toString Format

A `Result` formats as `Success(<value>)` or `Error(<message>)`:

```osprey
print(toString(15 / 3))   // "Success(5)"
print(toString(10 / 0))   // "Error(DivisionByZero)"
```