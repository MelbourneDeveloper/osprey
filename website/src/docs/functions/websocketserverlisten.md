---
layout: page
title: "websocketServerListen (Function)"
description: "Starts the WebSocket server listening for connections."
---

**Signature:** `websocketServerListen(serverID: int) -> int`

**Description:** Starts the WebSocket server listening for connections.

## Parameters

- **serverID** (int): Server identifier from websocketCreateServer

**Returns:** int

## Example

```osprey
let result = websocketServerListen(serverId)
print("WebSocket server listening")
```
