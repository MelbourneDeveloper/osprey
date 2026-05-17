---
layout: page
title: "parseFloat (Function)"
description: "Strict base-10 floating-point parser. No whitespace tolerance."
---

**Signature:** `parseFloat(s: string) -> Result<float, StringError>`

**Description:** Strict base-10 floating-point parser. No whitespace tolerance.

## Parameters

- **s** (string): The string to parse

**Returns:** Result<float, StringError>

## Example

```osprey
parseFloat("3.14")  // Success { value: 3.14 }
```
