---
layout: page
title: "httpStopServer (Function)"
description: "Stops the HTTP server and closes all connections."
---

**Signature:** `httpStopServer(serverID: int) -> int`

**Description:** Stops the HTTP server and closes all connections.

## Parameters

- **serverID** (int): Server identifier to stop

**Returns:** int

## Example

```osprey
let result = httpStopServer(serverId)
print("Server stopped")
```
