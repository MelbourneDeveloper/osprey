---
layout: page
title: "listAppend (Function)"
description: "Returns a new list with value at the end. O(log32 n) amortised."
---

**Signature:** `listAppend(list: List<T>, value: T) -> List<T>`

**Description:** Returns a new list with value at the end. O(log32 n) amortised.

## Parameters

- **list** (List): The list
- **value** (any): Value to append

**Returns:** List

## Example

```osprey
listAppend([1, 2], 3)  // [1, 2, 3]
```
