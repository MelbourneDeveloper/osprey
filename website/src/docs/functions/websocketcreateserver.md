---
layout: page
title: "websocketCreateServer (Function)"
description: "Creates a WebSocket server bound to the specified port, address, and path."
---

**Signature:** `websocketCreateServer(port: int, address: string, path: string) -> int`

**Description:** Creates a WebSocket server bound to the specified port, address, and path.

## Parameters

- **port** (int): Port number to bind to (1-65535)
- **address** (string): IP address to bind to (e.g., "127.0.0.1")
- **path** (string): WebSocket endpoint path (e.g., "/chat")

**Returns:** int

## Example

```osprey
let serverId = websocketCreateServer(8080, "127.0.0.1", "/chat")
print("WebSocket server created")
```
