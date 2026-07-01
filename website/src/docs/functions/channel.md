---
layout: page
title: "Channel (Function)"
description: "Creates a new channel with the specified capacity."
---

**Signature:** `Channel(capacity: int) -> Channel<t0>`

**Description:** Creates a new channel with the specified capacity.

## Parameters

- **capacity** (int): The capacity of the channel

**Returns:** Channel<t0>

## Example

```osprey
let ch = Channel(10)
```

```osprey-ml
ch = Channel 10
```
