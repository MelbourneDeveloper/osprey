---
layout: page
title: "fold (Function)"
description: "Reduces an iterator to a single value by repeatedly applying a function."
---

**Signature:** `fold(iterator: List<t0>, initial: t1, fn: (t1, t0) -> t1) -> t1`

**Description:** Reduces an iterator to a single value by repeatedly applying a function.

## Parameters

- **iterator** (List<t0>): The iterator to reduce
- **initial** (t1): The initial value for the accumulator
- **fn** ((t1, t0) -> t1): The reduction function that takes (accumulator, current) and returns new accumulator

**Returns:** t1

## Example

```osprey
range(1, 5) |> fold(0, add)  // sum: 0+1+2+3+4 = 10
```

```osprey-ml
range (1, 5) |> fold (0, add)  // sum: 0+1+2+3+4 = 10
```
