---
layout: page
title: "filter (Function)"
description: "Filters elements in an iterator based on a predicate function."
---

**Signature:** `filter(iterator: List<t0>, predicate: (t0) -> bool) -> List<t0>`

**Description:** Filters elements in an iterator based on a predicate function.

## Parameters

- **iterator** (List<t0>): The iterator to filter
- **predicate** ((t0) -> bool): The predicate function that returns true for elements to keep

**Returns:** List<t0>

## Example

```osprey
let evens = filter(range(1, 6), fn(x) { x % 2 == 0 })
forEach(evens, print)  // Prints: 2, 4
```
