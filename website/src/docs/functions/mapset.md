---
layout: page
title: "mapSet (Function)"
description: "Returns a new map with key bound to value (replaces prior binding)."
---

**Signature:** `mapSet(map: Map<K, V>, key: K, value: V) -> Map<K, V>`

**Description:** Returns a new map with key bound to value (replaces prior binding).

## Parameters

- **map** (Map): The map
- **key** (any): Key
- **value** (any): Value

**Returns:** Map

## Example

```osprey
mapSet({"a": 1}, "b", 2)  // {"a": 1, "b": 2}
```
