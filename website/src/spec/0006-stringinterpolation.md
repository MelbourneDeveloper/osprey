---
layout: page
title: "String Interpolation"
description: "Osprey Language Specification: String Interpolation"
date: 2026-06-30
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0006-stringinterpolation/"
---

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
- **Result types**: interpolation auto-unwraps — the success payload is rendered (context 5 of [Result Auto-Unwrapping](/spec/0004-typesystem/#result-auto-unwrapping)); an `Error` renders as `Error(<message>)`, preserving the payload per [ERR-PAYLOAD](/spec/0013-errorhandling/#error-payload-propagation--err-payload). To render the wrapper of a success, use `toString`.
- **Complex types**: Use `toString()` for explicit conversion

```osprey
let num = 42
let flag = true
print("Number: ${num}, Flag: ${flag}")

let result = 10 + 5
print("Result: ${result}")        // "Result: 15"  (auto-unwrapped)
print(toString(result))           // "Success(15)" (wrapper kept)
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