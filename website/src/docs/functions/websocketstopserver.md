---
layout: page
title: "websocketStopServer (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Stops the WebSocket server and closes all connections."
---

**Signature:** `websocketStopServer(serverID: Int) -> Result<Success, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Stops the WebSocket server and closes all connections.

## Parameters

- **serverID** (Int): Server identifier to stop

**Returns:** Result<Success, String>

## Example

```osprey
let stopResult = websocketStopServer(serverID: serverId)
match stopResult {
    Success _ => print("Server stopped successfully")
    Err message => print("Failed to stop server: ${message}")
}
```
