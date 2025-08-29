---
layout: page
title: "httpRequest (Function)"
description: "Makes a generic HTTP request with any method."
---

**Signature:** `httpRequest(clientID: int, method: int, path: string, headers: string, body: string) -> int`

**Description:** Makes a generic HTTP request with any method.

## Parameters

- **clientID** (int): Client identifier from httpCreateClient
- **method** (int): HTTP method (0=GET, 1=POST, 2=PUT, 3=DELETE)
- **path** (string): Request path
- **headers** (string): Additional headers
- **body** (string): Request body data

**Returns:** int

## Example

```osprey
let status = httpRequest(clientId, 0, "/custom", "", "")
print("Custom request status: ${status}")
```
