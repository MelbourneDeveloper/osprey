---
layout: page
title: "indexOf (Function)"
description: "Returns byte-index of first occurrence of needle, or Error(NotFound)."
---

**Signature:** `indexOf(s: string, needle: string) -> Result<int, StringError>`

**Description:** Returns byte-index of first occurrence of needle, or Error(NotFound).

## Parameters

- **s** (string): The string to search in
- **needle** (string): The substring to locate

**Returns:** Result<int, StringError>

## Example

```osprey
match indexOf("foo=bar", "=") { Success { value } => print(value) ... }
```
