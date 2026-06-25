---
layout: page
title: "recv (Function)"
description: "Receives a value from a channel."
---

**Signature:** `recv(channel: Channel<t0>) -> t0`

**Description:** Receives a value from a channel.

## Parameters

- **channel** (Channel<t0>): The channel to receive from

**Returns:** t0

## Example

```osprey
let value = recv(ch)
```
