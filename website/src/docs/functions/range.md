---
layout: page
title: "range (Function)"
description: "Creates an iterator that generates numbers from start to end (exclusive)."
---

**Signature:** `range(start: int, end: int) -> List<int>`

**Description:** Creates an iterator that generates numbers from start to end (exclusive).

## Parameters

- **start** (int): The starting number (inclusive)
- **end** (int): The ending number (exclusive)

**Returns:** List<int>

## Example

```osprey
forEach(range(0, 5), fn(x) { print(x) })  // Prints: 0, 1, 2, 3, 4
```
