---
layout: page
title: "websocketClose (Function)"
description: "Closes the WebSocket connection and cleans up resources. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*"
---

**Signature:** `websocketClose(wsID: int) -> Unit`

**Description:** Closes the WebSocket connection and cleans up resources. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **wsID** (int): WebSocket identifier to close

**Returns:** Unit

## Example

```osprey
let closeResult = websocketClose(wsID: wsId)
match closeResult {
    Success _ => print("Connection closed")
    Err message => print("Failed to close: ${message}")
}
```

```osprey-ml
closeResult = websocketClose wsId
match closeResult
    Success _ => print "Connection closed"
    Err message => print "Failed to close: ${message}"
```
