---
layout: page
title: "websocketConnect (Function)"
description: "Establishes a WebSocket connection with a message handler callback."
---

**Signature:** `websocketConnect(url: String, messageHandler: (String) -> Result<Success, String>) -> Result<WebSocketID, String>`

**Description:** Establishes a WebSocket connection with a message handler callback. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **url** (string): WebSocket URL (e.g., "ws://localhost:8080/chat")
- **messageHandler** ((string) -> Result<Success, String>): Callback function to handle incoming messages

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
