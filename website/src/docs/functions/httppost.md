---
layout: page
title: "httpPost (Function)"
description: "Makes an HTTP POST request with a request body."
---

**Signature:** `httpPost(clientID: int, path: string, body: string, headers: string) -> int`

**Description:** Makes an HTTP POST request with a request body.

## Parameters

- **clientID** (int): Client identifier from httpCreateClient
- **path** (string): Request path
- **body** (string): Request body data
- **headers** (string): Additional headers

**Returns:** int

## Example

```osprey
let status = httpPost(clientId, "/post", "{\"key\":\"value\"}", "Content-Type: application/json")
print("POST status: ${status}")
```
