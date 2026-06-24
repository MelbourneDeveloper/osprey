---
layout: page
title: "mapMerge (Function)"
description: "Right-biased union. Same as left + right."
---

**Signature:** `mapMerge(left: Map<t0, t1>, right: Map<t0, t1>) -> Map<t0, t1>`

**Description:** Right-biased union. Same as left + right.

## Parameters

- **left** (Map<t0, t1>): Left
- **right** (Map<t0, t1>): Right

**Returns:** Map<t0, t1>

## Example

```osprey
mapMerge({"a": 1}, {"b": 2})  // {"a": 1, "b": 2}
```
