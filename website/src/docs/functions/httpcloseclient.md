---
layout: page
title: "httpCloseClient (Function)"
description: "Closes the HTTP client and cleans up resources."
---

**Signature:** `httpCloseClient(clientID: int) -> int`

**Description:** Closes the HTTP client and cleans up resources.

## Parameters

- **clientID** (int): Client identifier to close

**Returns:** int

## Example

```osprey
let result = httpCloseClient(clientId)
print("Client closed")
```
