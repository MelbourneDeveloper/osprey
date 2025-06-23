---
layout: page
title: "websocketSend (Function)"
description: "Sends a message through the WebSocket connection."
---

**Signature:** `websocketSend(wsID: int, message: string) -> int`

**Description:** Sends a message through the WebSocket connection.

## Parameters

- **wsID** (int): WebSocket identifier from websocketConnect
- **message** (string): Message to send

**Returns:** int

## Example

```osprey
let result = websocketSend(wsId, "Hello, WebSocket!")
print("Message sent")
```
