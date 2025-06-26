---
layout: page
title: "websocketCreateServer (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<ServerID, String> and has critical runtime issues with port binding failures. Creates a WebSocket server bound to the specified port, address, and path."
---

**Signature:** `websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<ServerID, String> and has critical runtime issues with port binding failures. Creates a WebSocket server bound to the specified port, address, and path.

## Parameters

- **port** (Int): Port number to bind to (1-65535)
- **address** (String): IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")
- **path** (String): WebSocket endpoint path (e.g., "/chat", "/live")

**Returns:** Result<ServerID, String>

## Example

```osprey
let serverResult = websocketCreateServer(port: 8080, address: "127.0.0.1", path: "/chat")
match serverResult {
    Success serverId => print("WebSocket server created with ID: ${serverId}")
    Err message => print("Failed to create server: ${message}")
}
```
