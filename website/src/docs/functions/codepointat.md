---
layout: page
title: "codePointAt (Function)"
description: "Returns the Unicode code point that begins at the given byte index. Fails on an invalid index or malformed UTF-8."
---

**Signature:** `codePointAt(text: string, index: int) -> Result<int, Error>`

**Description:** Returns the Unicode code point that begins at the given byte index. Fails on an invalid index or malformed UTF-8.

## Parameters

- **text** (string): The string to read from
- **index** (int): Byte offset where the code point starts

**Returns:** Result<int, Error>

## Example

```osprey
match codePointAt("héllo", 1) {
  Success { value } => print("U+${value}")
  Error { message } => print(message)
}
```
