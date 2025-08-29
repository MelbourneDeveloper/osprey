---
layout: page
title: "httpDelete (Function)"
description: "Makes an HTTP DELETE request to the specified path."
---

**Signature:** `httpDelete(clientID: int, path: string, headers: string) -> int`

**Description:** Makes an HTTP DELETE request to the specified path.

## Parameters

- **clientID** (int): Client identifier from httpCreateClient
- **path** (string): Request path
- **headers** (string): Additional headers

**Returns:** int

## Example

```osprey
let status = httpDelete(clientId, "/delete", "")
print("DELETE status: ${status}")
```
