---
layout: page
title: "contains (Function)"
description: "True if needle appears anywhere in s. Empty needle returns true."
---

**Signature:** `contains(s: string, needle: string) -> bool`

**Description:** True if needle appears anywhere in s. Empty needle returns true.

## Parameters

- **s** (string): The string to search in
- **needle** (string): The substring to search for

**Returns:** bool

## Example

```osprey
let found = contains("hello world", "world")  // true
```
