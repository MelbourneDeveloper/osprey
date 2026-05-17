---
layout: page
title: "repeat (Function)"
description: "Concatenates s with itself n times. Error(InvalidArgument) on negative n."
---

**Signature:** `repeat(s: string, n: int) -> Result<string, StringError>`

**Description:** Concatenates s with itself n times. Error(InvalidArgument) on negative n.

## Parameters

- **s** (string): The string to repeat
- **n** (int): Repeat count, must be >= 0

**Returns:** Result<string, StringError>

## Example

```osprey
repeat("ab", 3)  // Success { value: "ababab" }
```
