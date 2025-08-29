---
layout: page
title: "websocketConnect (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<WebSocketID, String> and takes string handler instead of function pointer. Establishes a WebSocket connection with a message handler callback."
---

**Signature:** `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<WebSocketID, String> and takes string handler instead of function pointer. Establishes a WebSocket connection with a message handler callback.

## Parameters

- **url** (String): WebSocket URL (e.g., "ws://localhost:8080/chat")
- **messageHandler** (fn(String) -> Result<Success, String>): Callback function to handle incoming messages

**Returns:** Result<WebSocketID, String>

## Example

```osprey
fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}
let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)
match wsResult {
    Success wsId => print("Connected with ID: ${wsId}")
    Err message => print("Failed to connect: ${message}")
}
```
