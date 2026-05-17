---
layout: page
title: "fold (Function)"
description: "Reduces an iterator to a single value by repeatedly applying a function."
---

**Signature:** `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`

**Description:** Reduces an iterator to a single value by repeatedly applying a function.

## Parameters

- **iterator** (any): The iterator to reduce
- **initial** (any): The initial value for the accumulator
- **fn** (any): The reduction function that takes (accumulator, current) and returns new accumulator

**Returns:** any

## Example

```osprey
range(1, 5) |> fold(0, add)  // sum: 0+1+2+3+4 = 10
```
