---
layout: page
title: "send (Function)"
description: "Sends a value to a channel. Returns 1 for success, 0 for failure."
---

**Signature:** `send(channel: Channel<t0>, value: t0) -> Unit`

**Description:** Sends a value to a channel. Returns 1 for success, 0 for failure.

## Parameters

- **channel** (Channel<t0>): The channel to send to
- **value** (t0): The value to send

**Returns:** Unit

## Example

```osprey
let success = send(ch, 42)
```
