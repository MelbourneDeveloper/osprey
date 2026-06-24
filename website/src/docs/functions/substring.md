---
layout: page
title: "substring (Function)"
description: "Extracts s[start, end). Returns Error(IndexOutOfRange) if start<0, end>len, or start>end."
---

**Signature:** `substring(s: string, start: int, end: int) -> Result<string, Error>`

**Description:** Extracts s[start, end). Returns Error(IndexOutOfRange) if start<0, end>len, or start>end.

## Parameters

- **s** (string): The source string
- **start** (int): Starting index (inclusive)
- **end** (int): Ending index (exclusive)

**Returns:** Result<string, Error>

## Example

```osprey
substring("hello", 1, 4)  // Success { value: "ell" }
```
