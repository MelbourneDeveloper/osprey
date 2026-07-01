---
layout: page
title: "lines (Function)"
description: "Splits on '\\n'. A trailing newline does not produce an empty entry."
---

**Signature:** `lines(s: string) -> List<string>`

**Description:** Splits on '\n'. A trailing newline does not produce an empty entry.

## Parameters

- **s** (string): The string to split

**Returns:** List<string>

## Example

```osprey
lines("a\
b\
c")  // ["a","b","c"]
```
