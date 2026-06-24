---
layout: page
title: "websocketServerBroadcast (Function)"
description: "Broadcasts a message to all connected WebSocket clients. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*"
---

**Signature:** `websocketServerBroadcast(serverID: int, message: string) -> int`

**Description:** Broadcasts a message to all connected WebSocket clients. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **serverID** (int): Server identifier
- **message** (string): Message to broadcast to all clients

**Returns:** int

## Example

```osprey
let broadcastResult = websocketServerBroadcast(serverID: serverId, message: "Welcome to Osprey Chat!")
match broadcastResult {
    Success _ => print("Message broadcasted to all clients")
    Err message => print("Failed to broadcast: ${message}")
}
```
