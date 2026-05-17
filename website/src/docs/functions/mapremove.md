---
layout: page
title: "mapRemove (Function)"
description: "Returns a new map without key. No-op if key is absent."
---

**Signature:** `mapRemove(map: Map<K, V>, key: K) -> Map<K, V>`

**Description:** Returns a new map without key. No-op if key is absent.

## Parameters

- **map** (Map): The map
- **key** (any): Key

**Returns:** Map

## Example

```osprey
mapRemove({"a": 1, "b": 2}, "a")  // {"b": 2}
```
