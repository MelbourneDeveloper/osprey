---
layout: page
title: "httpCreateServer (Function)"
description: "Creates an HTTP server bound to the specified port and address."
---

**Signature:** `httpCreateServer(port: int, address: string) -> int`

**Description:** Creates an HTTP server bound to the specified port and address.

## Parameters

- **port** (int): Port number to bind to (1-65535)
- **address** (string): IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")

**Returns:** int

## Example

```osprey
let serverId = httpCreateServer(8080, "127.0.0.1")
print("Server created with ID: ${serverId}")
```
