---
layout: page
title: "split (Function)"
description: "Splits s on separator. Error(InvalidArgument) on empty separator."
---

**Signature:** `split(s: string, separator: string) -> Result<List<string>, Error>`

**Description:** Splits s on separator. Error(InvalidArgument) on empty separator.

## Parameters

- **s** (string): The string to split
- **separator** (string): Non-empty separator

**Returns:** Result<List<string>, Error>

## Example

```osprey
split("a,b,c", ",")  // Success { value: ["a","b","c"] }
```

```osprey-ml
split("a,b,c", ",")  // Success
    value = ["a","b","c"]
```
