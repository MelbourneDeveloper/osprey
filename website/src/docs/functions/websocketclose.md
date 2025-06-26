---
layout: page
title: "websocketClose (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Closes the WebSocket connection and cleans up resources."
---

**Signature:** `websocketClose(wsID: Int) -> Result<Success, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Closes the WebSocket connection and cleans up resources.

## Parameters

- **wsID** (Int): WebSocket identifier to close

**Returns:** Result<Success, String>

## Example

```osprey
let closeResult = websocketClose(wsID: wsId)
match closeResult {
    Success _ => print("Connection closed")
    Err message => print("Failed to close: ${message}")
}
```
