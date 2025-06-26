---
layout: page
title: "websocketServerBroadcast (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t (number of clients sent to) instead of Result<Success, String>. Broadcasts a message to all connected WebSocket clients."
---

**Signature:** `websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t (number of clients sent to) instead of Result<Success, String>. Broadcasts a message to all connected WebSocket clients.

## Parameters

- **serverID** (Int): Server identifier
- **message** (String): Message to broadcast to all clients

**Returns:** Result<Success, String>

## Example

```osprey
let broadcastResult = websocketServerBroadcast(serverID: serverId, message: "Welcome to Osprey Chat!")
match broadcastResult {
    Success _ => print("Message broadcasted to all clients")
    Err message => print("Failed to broadcast: ${message}")
}
```
