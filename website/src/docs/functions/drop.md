---
layout: page
title: "drop (Function)"
description: "Returns s without its first n bytes. Clamps; never fails."
---

**Signature:** `drop(s: string, n: int) -> string`

**Description:** Returns s without its first n bytes. Clamps; never fails.

## Parameters

- **s** (string): The source string
- **n** (int): How many bytes to drop

**Returns:** string

## Example

```osprey
drop("hello", 3)  // "lo"
```

```osprey-ml
drop ("hello", 3)  // "lo"
```
