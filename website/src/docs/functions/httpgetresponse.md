---
layout: page
title: "httpGetResponse (Function)"
description: "Sends an HTTP GET request and returns a response handle for inspecting the status, headers, and body."
---

**Signature:** `httpGetResponse(clientID: int, path: string, headers: string) -> Result<int, Error>`

**Description:** Sends an HTTP GET request and returns a response handle for inspecting the status, headers, and body.

## Parameters

- **clientID** (int): Client identifier from httpCreateClient
- **path** (string): Request path, e.g. "/api/users"
- **headers** (string): Additional request headers, or "" for none

**Returns:** Result<int, Error>

## Example

```osprey
match httpGetResponse(client, "/users", "") {
  Success { value } => print("status: ${httpResponseStatus(value)}")
  Error { message } => print(message)
}
```

```osprey-ml
match httpGetResponse (client, "/users", "")
    Success { value } => print "status: ${httpResponseStatus(value)}"
    Error { message } => print message
```
