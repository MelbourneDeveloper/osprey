---
layout: page
title: "websocketCreateServer (Function)"
description: "Creates a WebSocket server bound to the specified port, address, and path. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*"
---

**Signature:** `websocketCreateServer(port: int, address: string, path: string) -> int`

**Description:** Creates a WebSocket server bound to the specified port, address, and path. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **port** (int): Port number to bind to (1-65535)
- **address** (string): IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")
- **path** (string): WebSocket endpoint path (e.g., "/chat", "/live")

**Returns:** int

## Example

```osprey
let serverResult = websocketCreateServer(port: 8080, address: "127.0.0.1", path: "/chat")
match serverResult {
    Success serverId => print("WebSocket server created with ID: ${serverId}")
    Err message => print("Failed to create server: ${message}")
}
```

```osprey-ml
match serverResult
    Success serverId => print "WebSocket server created with ID: ${serverId}"
    Err message => print "Failed to create server: ${message}"
```
