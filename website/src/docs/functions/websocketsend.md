---
layout: page
title: "websocketSend (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Sends a message through the WebSocket connection."
---

**Signature:** `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Sends a message through the WebSocket connection.

## Parameters

- **wsID** (Int): WebSocket identifier from websocketConnect
- **message** (String): Message to send

**Returns:** Result<Success, String>

## Example

```osprey
let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
match sendResult {
    Success _ => print("Message sent successfully")
    Err message => print("Failed to send: ${message}")
}
```
