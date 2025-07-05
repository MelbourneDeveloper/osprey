---
layout: page
title: "Semantics"
description: "Osprey Language Specification: Semantics"
date: 2025-07-05
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0004-semantics/"
---

4. [Semantics](0004-Semantics.md)
   - [Variable Binding](#41-variable-binding)
   - [Function Semantics](#42-function-semantics)
   - [Evaluation Order](#43-evaluation-order)

## 4. Semantics

TODO: this section needs serious work. Add more detail.

### 4.1 Variable Binding

- `let` creates immutable bindings
- `mut` creates mutable bindings
- Variables must be initialized at declaration
- Shadowing is allowed in nested scopes

### 4.2 Function Semantics

- Functions are first-class values
- All functions are pure (no side effects except I/O)
- Recursive functions are supported
- Tail recursion is optimized

### 4.3 Evaluation Order

- Expressions are evaluated left-to-right
- Function arguments are evaluated before the function call
- Short-circuit evaluation for logical operators