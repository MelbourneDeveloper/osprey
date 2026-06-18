---
layout: page
title: "websocketClose (Function)"
description: "Closes the WebSocket connection and cleans up resources."
---

**Signature:** `websocketClose(wsID: Int) -> Result<Success, String>`

**Description:** Closes the WebSocket connection and cleans up resources. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **wsID** (int): WebSocket identifier to close

**Returns:** Result<Success, String>

## Example

```osprey
let closeResult = websocketClose(wsID: wsId)
match closeResult {
    Success _ => print("Connection closed")
    Err message => print("Failed to close: ${message}")
}
```
