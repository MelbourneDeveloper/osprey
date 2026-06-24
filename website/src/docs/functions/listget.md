---
layout: page
title: "listGet (Function)"
description: "Returns the element at the given index, or an error if the index is out of range."
---

**Signature:** `listGet(list: List<t0>, index: int) -> Result<t0, Error>`

**Description:** Returns the element at the given index, or an error if the index is out of range.

## Parameters

- **list** (List<t0>): The list to read from
- **index** (int): Zero-based element index

**Returns:** Result<t0, Error>

## Example

```osprey
match listGet(myList, 0) {
  Success { value } => print(value)
  Error { message } => print(message)
}
```
