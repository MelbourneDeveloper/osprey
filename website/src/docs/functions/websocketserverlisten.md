---
layout: page
title: "websocketServerListen (Function)"
description: "Starts the WebSocket server listening for connections. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*"
---

**Signature:** `websocketServerListen(serverID: int) -> int`

**Description:** Starts the WebSocket server listening for connections. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **serverID** (int): Server identifier from websocketCreateServer

**Returns:** int

## Example

```osprey
let listenResult = websocketServerListen(serverID: serverId)
match listenResult {
    Success _ => print("Server listening on ws://127.0.0.1:8080/chat")
    Err message => print("Failed to start listening: ${message}")
}
```

```osprey-ml
listenResult = websocketServerListen (serverID: serverId)
match listenResult
    Success _ => print "Server listening on ws://127.0.0.1:8080/chat"
    Err message => print "Failed to start listening: ${message}"
```
