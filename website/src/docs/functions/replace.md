---
layout: page
title: "replace (Function)"
description: "Replaces every occurrence of needle. Error(InvalidArgument) on empty needle."
---

**Signature:** `replace(s: string, needle: string, replacement: string) -> Result<string, Error>`

**Description:** Replaces every occurrence of needle. Error(InvalidArgument) on empty needle.

## Parameters

- **s** (string): The source string
- **needle** (string): The substring to find
- **replacement** (string): The replacement string

**Returns:** Result<string, Error>

## Example

```osprey
replace("a-b-c", "-", "_")  // Success { value: "a_b_c" }
```

```osprey-ml
replace ("a-b-c", "-", "_")  // Success { value: "a_b_c" }
```
