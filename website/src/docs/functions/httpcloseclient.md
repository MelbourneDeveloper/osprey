---
layout: page
title: "httpCloseClient (Function)"
description: "Closes the HTTP client and cleans up resources."
---

**Signature:** `httpCloseClient(clientID: int) -> Unit`

**Description:** Closes the HTTP client and cleans up resources.

## Parameters

- **clientID** (int): Client identifier to close

**Returns:** Unit

## Example

```osprey
let result = httpCloseClient(clientId)
print("Client closed")
```

```osprey-ml
result = httpCloseClient clientId
print "Client closed"
```
