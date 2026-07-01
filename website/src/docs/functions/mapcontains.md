---
layout: page
title: "mapContains (Function)"
description: "True iff key is present in map."
---

**Signature:** `mapContains(map: Map<t0, t1>, key: t0) -> bool`

**Description:** True iff key is present in map.

## Parameters

- **map** (Map<t0, t1>): The map
- **key** (t0): Key to find

**Returns:** bool

## Example

```osprey
mapContains({"a": 1}, "a")  // true
```

```osprey-ml
mapContains (["a" => 1], "a")  // true
```
