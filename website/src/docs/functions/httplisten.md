---
layout: page
title: "httpListen (Function)"
description: "Starts the HTTP server listening for requests with a handler function."
---

**Signature:** `httpListen(serverID: int, handler: (string, string, string, string) -> HttpResponse) -> int`

**Description:** Starts the HTTP server listening for requests with a handler function.

## Parameters

- **serverID** (int): Server identifier from httpCreateServer
- **handler** ((string, string, string, string) -> HttpResponse): Request handler function

**Returns:** int

## Example

```osprey
let result = httpListen(serverId, requestHandler)
print("Server listening")
```
