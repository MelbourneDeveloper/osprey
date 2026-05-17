---
layout: page
title: "map (Function)"
description: "Transforms each element in an iterator using a function, returning a new iterator."
---

**Signature:** `map(iterator: iterator, fn: function) -> iterator`

**Description:** Transforms each element in an iterator using a function, returning a new iterator.

## Parameters

- **iterator** (Iterator<T>): The iterator to transform
- **fn** (T -> U): The transformation function

**Returns:** Iterator<U>

## Example

```osprey
let doubled = map(range(1, 4), fn(x) { x * 2 })
forEach(doubled, print)  // Prints: 2, 4, 6
```
