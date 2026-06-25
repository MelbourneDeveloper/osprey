---
layout: page
title: "fromCodePoint (Function)"
description: "Returns the single-character string for a Unicode code point, or an error if it is not a valid scalar value."
---

**Signature:** `fromCodePoint(codePoint: int) -> Result<string, Error>`

**Description:** Returns the single-character string for a Unicode code point, or an error if it is not a valid scalar value.

## Parameters

- **codePoint** (int): The Unicode scalar value to encode

**Returns:** Result<string, Error>

## Example

```osprey
match fromCodePoint(233) {
  Success { value } => print(value)  // é
  Error { message } => print(message)
}
```
