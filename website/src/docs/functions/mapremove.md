---
layout: page
title: "mapRemove (Function)"
description: "Returns a new map without key. No-op if key is absent."
---

**Signature:** `mapRemove(map: Map<t0, t1>, key: t0) -> Map<t0, t1>`

**Description:** Returns a new map without key. No-op if key is absent.

## Parameters

- **map** (Map<t0, t1>): The map
- **key** (t0): Key

**Returns:** Map<t0, t1>

## Example

```osprey
mapRemove({"a": 1, "b": 2}, "a")  // {"b": 2}
```

```osprey-ml
mapRemove ({"a": 1, "b": 2}, "a")  // {"b": 2}
```
