---
layout: page
title: "mapContains (Function)"
description: "True iff key is present in map."
---

**Signature:** `mapContains(map: Map<K, V>, key: K) -> bool`

**Description:** True iff key is present in map.

## Parameters

- **map** (Map): The map
- **key** (any): Key to find

**Returns:** bool

## Example

```osprey
mapContains({"a": 1}, "a")  // true
```
