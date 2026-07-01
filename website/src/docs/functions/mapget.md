---
layout: page
title: "mapGet (Function)"
description: "Returns the value associated with the key, or an error if the key is absent."
---

**Signature:** `mapGet(map: Map<t0, t1>, key: t0) -> Result<t1, Error>`

**Description:** Returns the value associated with the key, or an error if the key is absent.

## Parameters

- **map** (Map<t0, t1>): The map to look up in
- **key** (t0): The key to find

**Returns:** Result<t1, Error>

## Example

```osprey
match mapGet(scores, "alice") {
  Success { value } => print(value)
  Error { message } => print(message)
}
```

```osprey-ml
match mapGet (scores, "alice")
    Success { value } => print value
    Error { message } => print message
```
