---
layout: page
title: "forEach (Function)"
description: "Applies a function to each element in an iterator."
---

**Signature:** `forEach(iterator: List<t0>, function: (t0) -> Unit) -> Unit`

**Description:** Applies a function to each element in an iterator.

## Parameters

- **iterator** (List<t0>): The iterator to process
- **function** ((t0) -> Unit): The function to apply to each element

**Returns:** Unit

## Example

```osprey
forEach(range(1, 4), fn(x) { print(x * 2) })  // Prints: 2, 4, 6
```

```osprey-ml
forEach (range (1, 4), \x => print x * 2)  // Prints: 2, 4, 6
```
