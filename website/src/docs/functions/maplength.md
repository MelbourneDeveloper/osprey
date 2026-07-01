---
layout: page
title: "mapLength (Function)"
description: "Returns the number of entries in a map. O(1)."
---

**Signature:** `mapLength(map: Map<t0, t1>) -> int`

**Description:** Returns the number of entries in a map. O(1).

## Parameters

- **map** (Map<t0, t1>): The map

**Returns:** int

## Example

```osprey
mapLength({"a": 1, "b": 2})  // 2
```

```osprey-ml
mapLength ["a" => 1, "b" => 2]  // 2
```
