---
layout: page
title: "websocketKeepAlive (Function)"
description: "Keeps the WebSocket server running indefinitely until interrupted (blocking operation). *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*"
---

**Signature:** `websocketKeepAlive() -> Unit`

**Description:** Keeps the WebSocket server running indefinitely until interrupted (blocking operation). *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*

**Returns:** Unit

## Example

```osprey
websocketKeepAlive()  // Blocks until Ctrl+C
```
