---
layout: page
title: "Block Expressions"
description: "Osprey Language Specification: Block Expressions"
date: 2026-06-18
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0008-blockexpressions/"
---

# Block Expressions

A block expression groups statements and returns the value of its final expression. Each block introduces a new lexical scope.

```ebnf
blockExpression ::= "{" statement* expression? "}"
```

**Examples:**
```osprey
// Simple block with local variables
let result = {
    let x = 10
    let y = 20
    x + y
}
print("Result: ${result}")  // prints "Result: 30"

// Nested blocks
let complex = {
    let outer = 100
    let inner_result = {
        let inner = 50
        outer + inner
    }
    inner_result * 2
}
print("Complex: ${complex}")  // prints "Complex: 300"

// Block with function calls
fn multiply(a: int, b: int) -> int = a * b
let calc = {
    let a = 5
    let b = 6
    multiply(a: a, b: b)
}
print("Calculation: ${calc}")  // prints "Calculation: 30"
```

## Block Scoping Rules

Block expressions create a new lexical scope:
- Variables declared inside a block are only visible within that block
- Variables from outer scopes can be accessed (lexical scoping)
- Variables declared in a block shadow outer variables with the same name
- Variables go out of scope when the block ends

**Scoping Examples:**
```osprey
let x = 100
let result = {
    let x = 50        // Shadows outer x
    let y = 25        // Only visible in this block
    x + y             // Uses inner x (50)
}
print("Result: ${result}")  // 75
print("Outer x: ${x}")      // 100 (unchanged)
// print("${y}")            // ERROR: y not in scope
```

## Block Return Values

A block ending with an expression returns that expression's value and adopts its type. A block ending with a statement returns `unit`.