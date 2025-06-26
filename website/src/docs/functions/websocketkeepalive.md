---
layout: page
title: "webSocketKeepAlive (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. Keeps the WebSocket server running indefinitely until interrupted (blocking operation)."
---

**Signature:** `webSocketKeepAlive() -> Unit`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. Keeps the WebSocket server running indefinitely until interrupted (blocking operation).

**Returns:** Unit

## Example

```osprey
webSocketKeepAlive()  // Blocks until Ctrl+C
```
