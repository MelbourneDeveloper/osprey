---
layout: page
title: "Function Calls"
description: "Osprey Language Specification: Function Calls"
date: 2025-07-06
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0006-functioncalls/"
---

6. [Function Calls](0006-FunctionCalls.md)
   - [Named Arguments Requirement](#61-named-arguments-requirement)
       - [Valid Function Calls](#valid-function-calls)
       - [Invalid Function Calls](#invalid-function-calls)
   - [Function Call Compilation Rules](#62-function-call-compilation-rules)

## 6. Function Calls

### 6.1 Named Arguments Requirement

**CRITICAL RULE**: Functions with more than one parameter **MUST** be called with named arguments.

#### Valid Function Calls

```osprey
// Single parameter - positional allowed
fn double(x) = x * 2
let result = double(5)

// Zero parameters - no arguments
fn getValue() = 42
let value = getValue()

// Multiple parameters - named arguments required
fn add(x, y) = x + y
let sum = add(x: 10, y: 20)

// Multiple parameters - order doesn't matter with named args
let sum2 = add(y: 20, x: 10)
```

#### Invalid Function Calls

```osprey
// ERROR: Multi-parameter function with positional arguments
fn add(x, y) = x + y
let sum = add(10, 20)  // ❌ Compilation error

// ERROR: Mixed positional and named arguments
let sum = add(10, y: 20)  // ❌ Compilation error
```

### 6.2 Function Call Compilation Rules

1. **Single Parameter Functions**: May use positional arguments
2. **Zero Parameter Functions**: Called with empty parentheses `()`
3. **Multi-Parameter Functions**: Must use named arguments for ALL parameters
4. **Argument Reordering**: Named arguments are reordered to match parameter declaration order