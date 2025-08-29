---
layout: page
title: "websocketKeepAlive (Function)"
description: "⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. Keeps the WebSocket server running indefinitely until interrupted (blocking operation)."
---

**Signature:** `websocketKeepAlive() -> Unit`

**Description:** ⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. Keeps the WebSocket server running indefinitely until interrupted (blocking operation).

**Returns:** Unit

## Example

```osprey
websocketKeepAlive()  // Blocks until Ctrl+C
```
