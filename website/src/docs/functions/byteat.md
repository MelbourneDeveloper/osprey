---
layout: page
title: "byteAt (Function)"
description: "Returns the byte at the given index (0-255), or an error if the index is out of range."
---

**Signature:** `byteAt(text: string, index: int) -> Result<int, Error>`

**Description:** Returns the byte at the given index (0-255), or an error if the index is out of range.

## Parameters

- **text** (string): The string to read from
- **index** (int): Zero-based byte offset

**Returns:** Result<int, Error>

## Example

```osprey
match byteAt("hi", 0) {
  Success { value } => print("byte: ${value}")
  Error { message } => print(message)
}
```

```osprey-ml
match byteAt ("hi", 0)
    Success { value } => print "byte: ${value}"
    Error { message } => print message
```
