---
layout: page
title: "httpPut (Function)"
description: "Makes an HTTP PUT request with a request body."
---

**Signature:** `httpPut(clientID: int, path: string, body: string, headers: string) -> int`

**Description:** Makes an HTTP PUT request with a request body.

## Parameters

- **clientID** (int): Client identifier from httpCreateClient
- **path** (string): Request path
- **body** (string): Request body data
- **headers** (string): Additional headers

**Returns:** int

## Example

```osprey
let status = httpPut(clientId, "/put", "{\"updated\":\"data\"}", "Content-Type: application/json")
print("PUT status: ${status}")
```
