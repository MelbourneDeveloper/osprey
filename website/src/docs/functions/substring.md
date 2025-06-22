---
layout: page
title: "substring (Function)"
description: "Extracts a substring from start to end index."
---

**Signature:** `substring(s: string, start: int, end: int) -> string`

**Description:** Extracts a substring from start to end index.

## Parameters

- **s** (string): The source string
- **start** (int): Starting index (inclusive)
- **end** (int): Ending index (exclusive)

**Returns:** string

## Example

```osprey
let sub = substring("hello", 1, 4)
print(sub)  // Prints: ell
```
