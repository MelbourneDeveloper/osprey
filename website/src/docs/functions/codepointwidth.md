---
layout: page
title: "codePointWidth (Function)"
description: "Returns how many bytes the given Unicode code point occupies in UTF-8 (1-4)."
---

**Signature:** `codePointWidth(codePoint: int) -> Result<int, Error>`

**Description:** Returns how many bytes the given Unicode code point occupies in UTF-8 (1-4).

## Parameters

- **codePoint** (int): The Unicode scalar value

**Returns:** Result<int, Error>

## Example

```osprey
match codePointWidth(233) {
  Success { value } => print("${value} bytes")
  Error { message } => print(message)
}
```
