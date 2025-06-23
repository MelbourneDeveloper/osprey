---
layout: page
title: "websocketServerBroadcast (Function)"
description: "Broadcasts a message to all connected WebSocket clients."
---

**Signature:** `websocketServerBroadcast(serverID: int, message: string) -> int`

**Description:** Broadcasts a message to all connected WebSocket clients.

## Parameters

- **serverID** (int): Server identifier
- **message** (string): Message to broadcast

**Returns:** int

## Example

```osprey
let result = websocketServerBroadcast(serverId, "Hello everyone!")
print("Message broadcasted")
```
