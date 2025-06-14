---
layout: page
title: "forEach (Function)"
description: "Applies a function to each element in an iterator."
---

**Signature:** `forEach(iterator: iterator, fn: function) -> int`

**Description:** Applies a function to each element in an iterator.

## Parameters

- **iterator** (iterator): The iterator to process
- **fn** (function): The function to apply to each element

**Returns:** int

## Example

```osprey
forEach(range(1, 4), fn(x) { print(x * 2) })  // Prints: 2, 4, 6
```
