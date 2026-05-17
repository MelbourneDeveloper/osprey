---
layout: page
title: "take (Function)"
description: "Returns at most the first n bytes of s. Clamps; never fails."
---

**Signature:** `take(s: string, n: int) -> string`

**Description:** Returns at most the first n bytes of s. Clamps; never fails.

## Parameters

- **s** (string): The source string
- **n** (int): How many bytes to take

**Returns:** string

## Example

```osprey
take("hello", 3)  // "hel"
```
