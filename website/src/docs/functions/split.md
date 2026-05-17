---
layout: page
title: "split (Function)"
description: "Splits s on separator. Error(InvalidArgument) on empty separator."
---

**Signature:** `split(s: string, separator: string) -> Result<List<string>, StringError>`

**Description:** Splits s on separator. Error(InvalidArgument) on empty separator.

## Parameters

- **s** (string): The string to split
- **separator** (string): Non-empty separator

**Returns:** Result<List<string>, StringError>

## Example

```osprey
split("a,b,c", ",")  // Success { value: ["a","b","c"] }
```
