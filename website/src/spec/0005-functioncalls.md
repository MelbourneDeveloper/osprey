---
layout: page
title: "Function Calls"
description: "Osprey Language Specification: Function Calls"
date: 2026-06-18
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0005-functioncalls/"
---

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

## Rules

1. Zero parameters: empty parentheses, `f()`.
2. One parameter: positional or named.
3. Two or more parameters: every argument must be named. Mixing positional and named arguments is a compilation error.
4. **Built-in functions** ([Built-in Functions](/spec/0012-built-infunctions/)) are exempt: they take positional arguments in subject-first order — `split("a,b,c", ",")`, `fold(xs, 0, add)` — so the pipe can supply the subject as the first argument: `xs |> fold(0, add)`.

Argument order at the call site is independent of declaration order; the compiler reorders by name.