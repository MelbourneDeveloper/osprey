---
layout: page
title: "websocketConnect (Function)"
description: "Connects to a WebSocket server at the given URL and returns a connection id."
---

**Signature:** `websocketConnect(url: string) -> int`

**Description:** Connects to a WebSocket server at the given URL and returns a connection id.

## Parameters

- **url** (string): WebSocket URL, e.g. "ws://localhost:8080/chat"

**Returns:** int

## Example

```osprey
let conn = websocketConnect("ws://localhost:8080/chat")
```
