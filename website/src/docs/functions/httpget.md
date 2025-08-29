---
layout: page
title: "httpGet (Function)"
description: "Makes an HTTP GET request to the specified path."
---

**Signature:** `httpGet(clientID: int, path: string, headers: string) -> int`

**Description:** Makes an HTTP GET request to the specified path.

## Parameters

- **clientID** (int): Client identifier from httpCreateClient
- **path** (string): Request path (e.g., "/api/users")
- **headers** (string): Additional headers (e.g., "Authorization: Bearer token")

**Returns:** int

## Example

```osprey
let status = httpGet(clientId, "/get", "")
print("GET request status: ${status}")
```
