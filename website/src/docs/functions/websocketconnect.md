---
layout: page
title: "websocketConnect (Function)"
description: "Establishes a WebSocket connection to the specified URL."
---

**Signature:** `websocketConnect(url: string, messageHandler: string) -> int`

**Description:** Establishes a WebSocket connection to the specified URL.

## Parameters

- **url** (string): WebSocket URL (e.g., "ws://localhost:8080/chat")
- **messageHandler** (string): Message handler identifier

**Returns:** int

## Example

```osprey
let wsId = websocketConnect("ws://localhost:8080/chat", "handler")
print("Connected with ID: ${wsId}")
```
