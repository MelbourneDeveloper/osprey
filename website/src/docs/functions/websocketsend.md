---
layout: page
title: "websocketSend (Function)"
description: "Sends a message through the WebSocket connection. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*"
---

**Signature:** `websocketSend(wsID: int, message: string) -> int`

**Description:** Sends a message through the WebSocket connection. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

## Parameters

- **wsID** (int): WebSocket identifier from websocketConnect
- **message** (string): Message to send

**Returns:** int

## Example

```osprey
let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
match sendResult {
    Success _ => print("Message sent successfully")
    Err message => print("Failed to send: ${message}")
}
```

```osprey-ml
sendResult = websocketSend wsId "Hello, WebSocket!"
match sendResult
    Success _ => print "Message sent successfully"
    Err message => print "Failed to send: ${message}"
```
