---
layout: page
title: "map (Function)"
description: "Transforms each element in an iterator using a function, returning a new iterator."
---

**Signature:** `map(iterator: List<t0>, fn: (t0) -> t1) -> List<t1>`

**Description:** Transforms each element in an iterator using a function, returning a new iterator.

## Parameters

- **iterator** (List<t0>): The iterator to transform
- **fn** ((t0) -> t1): The transformation function

**Returns:** List<t1>

## Example

```osprey
let doubled = map(range(1, 4), fn(x) { x * 2 })
forEach(doubled, print)  // Prints: 2, 4, 6
```

```osprey-ml
doubled = map(range (1, 4), fn(x) { x * 2 })
forEach (doubled, print)  // Prints: 2, 4, 6
```
