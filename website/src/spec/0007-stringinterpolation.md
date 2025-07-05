---
layout: page
title: "String Interpolation"
description: "Osprey Language Specification: String Interpolation"
date: 2025-07-05
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0007-stringinterpolation/"
---

7. [String Interpolation](0007-StringInterpolation.md)
   - [Syntax](#71-syntax)
   - [Expression Support](#72-expression-support)
   - [Type Handling](#73-type-handling)
   - [Implementation](#74-implementation)

## 7. String Interpolation

✅ **IMPLEMENTED**: String interpolation is fully implemented and working with comprehensive test coverage.

### 7.1 Syntax

String interpolation uses `${}` syntax within double-quoted strings:

```osprey
let name = "Alice"
let age = 30
let message = "Hello ${name}, you are ${age} years old"
```

### 7.2 Expression Support

Any expression is valid inside interpolation:

```osprey
let x = 10
let y = 5
print("Sum: ${x + y}")
print("Product: ${x * y}")
print("Complex: ${(x + y) * 2 - 1}")
```

### 7.3 Type Handling

- **String variables**: Use `%s` format specifier
- **Integer expressions**: Use `%ld` format specifier  
- **Function calls**: Supported for single-parameter functions

### 7.4 Implementation

Interpolated strings are compiled to:
1. Allocate a buffer (`alloca [1024 x i8]`)
2. Use `sprintf` to format the string
3. Use `puts` to output the result