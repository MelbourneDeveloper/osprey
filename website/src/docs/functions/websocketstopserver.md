---
layout: page
title: "websocketStopServer (Function)"
description: "Stops the WebSocket server and closes all connections."
---

**Signature:** `websocketStopServer(serverID: Int) -> Result<Success, String>`

**Description:** Stops the WebSocket server and closes all connections. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **serverID** (int): Server identifier to stop

**Returns:** Result<Success, String>

## Example

```osprey
let stopResult = websocketStopServer(serverID: serverId)
match stopResult {
    Success _ => print("Server stopped successfully")
    Err message => print("Failed to stop server: ${message}")
}
```
