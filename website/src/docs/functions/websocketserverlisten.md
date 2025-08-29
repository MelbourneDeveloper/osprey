---
layout: page
title: "websocketServerListen (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String> and currently returns -4 (bind failed) due to port binding issues. Starts the WebSocket server listening for connections."
---

**Signature:** `websocketServerListen(serverID: Int) -> Result<Success, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String> and currently returns -4 (bind failed) due to port binding issues. Starts the WebSocket server listening for connections.

## Parameters

- **serverID** (Int): Server identifier from websocketCreateServer

**Returns:** Result<Success, String>

## Example

```osprey
let listenResult = websocketServerListen(serverID: serverId)
match listenResult {
    Success _ => print("Server listening on ws://127.0.0.1:8080/chat")
    Err message => print("Failed to start listening: ${message}")
}
```
