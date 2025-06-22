---
layout: page
title: "websocketStopServer (Function)"
description: "Stops the WebSocket server and closes all connections."
---

**Signature:** `websocketStopServer(serverID: int) -> int`

**Description:** Stops the WebSocket server and closes all connections.

## Parameters

- **serverID** (int): Server identifier to stop

**Returns:** int

## Example

```osprey
let result = websocketStopServer(serverId)
print("WebSocket server stopped")
```
