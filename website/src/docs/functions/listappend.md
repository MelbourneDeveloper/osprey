---
layout: page
title: "listAppend (Function)"
description: "Returns a new list with value at the end. O(log32 n) amortised."
---

**Signature:** `listAppend(list: List<t0>, value: t0) -> List<t0>`

**Description:** Returns a new list with value at the end. O(log32 n) amortised.

## Parameters

- **list** (List<t0>): The list
- **value** (t0): Value to append

**Returns:** List<t0>

## Example

```osprey
listAppend([1, 2], 3)  // [1, 2, 3]
```
