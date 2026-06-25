---
layout: page
title: "mapSet (Function)"
description: "Returns a new map with key bound to value (replaces prior binding)."
---

**Signature:** `mapSet(map: Map<t0, t1>, key: t0, value: t1) -> Map<t0, t1>`

**Description:** Returns a new map with key bound to value (replaces prior binding).

## Parameters

- **map** (Map<t0, t1>): The map
- **key** (t0): Key
- **value** (t1): Value

**Returns:** Map<t0, t1>

## Example

```osprey
mapSet({"a": 1}, "b", 2)  // {"a": 1, "b": 2}
```
