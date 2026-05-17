---
layout: page
title: "send (Function)"
description: "Sends a value to a channel. Returns 1 for success, 0 for failure."
---

**Signature:** `send(channel: Channel, value: any) -> int`

**Description:** Sends a value to a channel. Returns 1 for success, 0 for failure.

## Parameters

- **channel** (Channel): The channel to send to
- **value** (any): The value to send

**Returns:** int

## Example

```osprey
let success = send(ch, 42)
```
