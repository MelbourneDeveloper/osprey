---
layout: page
title: "contains (Function)"
description: "Checks if a string contains a substring."
---

**Signature:** `contains(haystack: string, needle: string) -> bool`

**Description:** Checks if a string contains a substring.

## Parameters

- **haystack** (string): The string to search in
- **needle** (string): The substring to search for

**Returns:** bool

## Example

```osprey
let found = contains("hello world", "world")
print(found)  // Prints: true
```
