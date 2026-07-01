---
layout: page
title: "padEnd (Function)"
description: "Pads s on the right with copies of fill to reach targetLength bytes."
---

**Signature:** `padEnd(s: string, targetLength: int, fill: string) -> Result<string, Error>`

**Description:** Pads s on the right with copies of fill to reach targetLength bytes.

## Parameters

- **s** (string): The string to pad
- **targetLength** (int): Desired total length
- **fill** (string): Padding string (non-empty)

**Returns:** Result<string, Error>

## Example

```osprey
padEnd("7", 3, ".")  // Success { value: "7.." }
```

```osprey-ml
padEnd ("7", 3, ".")  // Success { value: "7.." }
```
