---
layout: page
title: "mapMerge (Function)"
description: "Right-biased union. Same as left + right."
---

**Signature:** `mapMerge(left: Map<K, V>, right: Map<K, V>) -> Map<K, V>`

**Description:** Right-biased union. Same as left + right.

## Parameters

- **left** (Map): Left
- **right** (Map): Right

**Returns:** Map

## Example

```osprey
mapMerge({"a": 1}, {"b": 2})  // {"a": 1, "b": 2}
```
