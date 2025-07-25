---
layout: page
title: "Error Handling"
description: "Osprey Language Specification: Error Handling"
date: 2025-07-06
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0014-errorhandling/"
---

14. [Error Handling](0015-ErrorHandling.md)
    - [The Result Type](#151-the-result-type)

## 14. Error Handling
### 14.1 The Result Type

**CRITICAL**: All functions that can fail **MUST** return a `Result` type. There are no exceptions, panics, or nulls. This is a core design principle of the language to ensure safety and eliminate entire classes of runtime errors.

The `Result` type is a generic union type with two variants:

- `Success { value: T }`: Represents a successful result, containing the value of type `T`.
- `Error { message: E }`: Represents an error, containing an error message or object of type `E`.

**Example:**
```osprey
type Result<T, E> = Success { value: T } | Error { message: E }
```

The compiler **MUST** enforce that `Result` types are always handled with a `match` expression, preventing direct access to the underlying value and ensuring that all possible outcomes are considered.

```osprey
let result = someFunctionThatCanFail()

match result {
    Success { value } => print("Success: ${value}")
    Error { message } => print("Error: ${message}")
}
```

This approach guarantees that error handling is explicit, robust, and checked at compile time.

### 14.2 Compound Expression Result Propagation

**🚨 CRITICAL DESIGN PRINCIPLE 🚨**: When multiple arithmetic operations are combined in a single expression, the **entire expression** returns a single `Result` type, not each individual operation.

## **THE GOLDEN RULE:**
```osprey
let a = 1 + 3                    // ✅ a: Result<int, MathError>
let a = 1 + 3 + (300 / 5)        // ✅ a: Result<int, MathError> (WHOLE EXPRESSION)
// Individual operations inside are AUTOMATICALLY PROPAGATED
```

#### Individual Operations (Each Returns Result)
```osprey
let a = 1 + 3          // a: Result<int, MathError>
let b = 300 / 5        // b: Result<int, MathError>
let c = 2 * 4          // c: Result<int, MathError>
```

#### Compound Expressions (Single Result, Auto-Propagation)
```osprey
let result = 1 + 3 + (300 / 5)         // result: Result<int, MathError>
let complex = (a * b) + (c / d) - 10   // complex: Result<int, MathError>
let nested = ((x + y) * z) / (a - b)   // nested: Result<int, MathError>
let mega = 1 + 2 * 3 - 4 / 2 + 5       // mega: Result<int, MathError>
```

## **WHAT THIS MEANS:**
- ❌ **NO**: You don't handle Results for each `+`, `-`, `*`, `/` inside an expression
- ✅ **YES**: You handle the Result **ONCE** for the entire compound expression
- ⚡ **AUTO**: If any operation fails → whole expression fails
- ⚡ **AUTO**: If all operations succeed → expression returns Success with final value

#### Error Propagation Rules
- **Any operation fails** → **Entire expression fails**
- **All operations succeed** → **Expression returns Success with final value**
- **Individual operations** inside compound expressions **don't need explicit Result handling**

#### Example Behavior
```osprey
// If any operation overflows, the whole expression fails
let calculation = 1000000 * 1000000 + 50 / 2  // Result<int, MathError>

match calculation {
    Success { value } => print("Final result: ${value}")
    Error { message } => print("Calculation failed: ${message}")
}
```

#### Rationale
This design provides:
- **Ergonomic code**: No need to unwrap Results for every operation
- **Safety**: All potential arithmetic errors are still caught
- **Clarity**: Single error handling point for compound expressions
- **Performance**: Runtime can optimize arithmetic chains

**KEY INSIGHT**: You handle the Result **once** for the entire expression, not for each individual operation within it.