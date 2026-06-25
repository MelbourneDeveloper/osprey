---
layout: page
title: "httpStopServer (Function)"
description: "Stops the HTTP server and closes all connections."
---

**Signature:** `httpStopServer(serverID: int) -> Unit`

**Description:** Stops the HTTP server and closes all connections.

## Parameters

- **serverID** (int): Server identifier to stop

**Returns:** Unit

## Example

```osprey
let result = httpStopServer(serverId)
print("Server stopped")
```
