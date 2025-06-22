---
layout: page
title: "websocketClose (Function)"
description: "Closes the WebSocket connection."
---

**Signature:** `websocketClose(wsID: int) -> int`

**Description:** Closes the WebSocket connection.

## Parameters

- **wsID** (int): WebSocket identifier to close

**Returns:** int

## Example

```osprey
let result = websocketClose(wsId)
print("WebSocket closed")
```
