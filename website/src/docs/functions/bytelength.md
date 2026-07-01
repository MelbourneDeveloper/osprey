---
layout: page
title: "byteLength (Function)"
description: "Returns the number of bytes in the string's UTF-8 encoding."
---

**Signature:** `byteLength(text: string) -> int`

**Description:** Returns the number of bytes in the string's UTF-8 encoding.

## Parameters

- **text** (string): The string to measure

**Returns:** int

## Example

```osprey
let n = byteLength("héllo")  // 6
```

```osprey-ml
n = byteLength "héllo"  // 6
```
