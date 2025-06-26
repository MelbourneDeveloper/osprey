---
layout: page
title: "Error Handling"
description: "Osprey Language Specification: Error Handling"
date: 2025-06-26
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