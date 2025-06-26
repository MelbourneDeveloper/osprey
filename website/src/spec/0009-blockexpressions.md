---
layout: page
title: "9. Block Expressions"
description: "Osprey Language Specification: 9. Block Expressions"
date: 2025-06-25
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0009-blockexpressions/"
---

## 9.1 Block Scoping Rules

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

## 9.2 Block Return Values

Block expressions return the value of their final expression:
- If the block ends with an expression, that value is returned
- If the block has no final expression, it returns the unit type
- The block's type is determined by the type of the final expression

## 9.3 Performance Characteristics

Block expressions are zero-cost abstractions:
- **Compile-time scoping**: All variable scoping resolved at compile time
- **No runtime overhead**: Blocks compile to sequential instructions
- **Stack allocation**: Local variables allocated on the stack
- **Optimized away**: Simple blocks with no local variables are optimized away