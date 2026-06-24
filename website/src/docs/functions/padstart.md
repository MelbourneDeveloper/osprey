---
layout: page
title: "padStart (Function)"
description: "Pads s on the left with copies of fill to reach targetLength bytes."
---

**Signature:** `padStart(s: string, targetLength: int, fill: string) -> Result<string, Error>`

**Description:** Pads s on the left with copies of fill to reach targetLength bytes.

## Parameters

- **s** (string): The string to pad
- **targetLength** (int): Desired total length
- **fill** (string): Padding string (non-empty)

**Returns:** Result<string, Error>

## Example

```osprey
padStart("7", 3, "0")  // Success { value: "007" }
```
