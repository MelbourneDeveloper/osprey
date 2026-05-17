---
layout: page
title: "recv (Function)"
description: "Receives a value from a channel."
---

**Signature:** `recv(channel: Channel) -> any`

**Description:** Receives a value from a channel.

## Parameters

- **channel** (Channel): The channel to receive from

**Returns:** any

## Example

```osprey
let value = recv(ch)
```
