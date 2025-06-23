---
layout: page
title: "fold (Function)"
description: "Reduces an iterator to a single value using an accumulator function."
---

**Signature:** `fold(iterator: iterator, initial: any, fn: function) -> any`

**Description:** Reduces an iterator to a single value using an accumulator function.

## Parameters

- **iterator** (iterator): The iterator to reduce
- **initial** (any): The initial value for the accumulator
- **fn** (function): The reduction function that takes (accumulator, current) and returns new accumulator

**Returns:** any

## Example

```osprey
let sum = fold(range(1, 5), 0, fn(acc, x) { acc + x })
print(sum)  // Prints: 10
```
