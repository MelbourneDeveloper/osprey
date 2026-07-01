---
layout: page
title: "parseInt (Function)"
description: "Strict base-10 signed-int parser. No whitespace tolerance."
---

**Signature:** `parseInt(s: string) -> Result<int, Error>`

**Description:** Strict base-10 signed-int parser. No whitespace tolerance.

## Parameters

- **s** (string): The string to parse

**Returns:** Result<int, Error>

## Example

```osprey
parseInt("42")  // Success { value: 42 }
```

```osprey-ml
parseInt("42")  // Success
    value = 42
```
