---
layout: page
title: "filter (Function)"
description: "Filters elements in an iterator based on a predicate function."
---

**Signature:** `filter(iterator: iterator, predicate: function) -> iterator`

**Description:** Filters elements in an iterator based on a predicate function.

## Parameters

- **iterator** (iterator): The iterator to filter
- **predicate** (function): The predicate function that returns true for elements to keep

**Returns:** iterator

## Example

```osprey
let evens = filter(range(1, 6), fn(x) { x % 2 == 0 })
forEach(evens, print)  // Prints: 2, 4
```
