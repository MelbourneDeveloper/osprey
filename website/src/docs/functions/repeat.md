---
layout: page
title: "repeat (Function)"
description: "Concatenates s with itself n times. Error(InvalidArgument) on negative n."
---

**Signature:** `repeat(s: string, n: int) -> Result<string, Error>`

**Description:** Concatenates s with itself n times. Error(InvalidArgument) on negative n.

## Parameters

- **s** (string): The string to repeat
- **n** (int): Repeat count, must be >= 0

**Returns:** Result<string, Error>

## Example

```osprey
repeat("ab", 3)  // Success { value: "ababab" }
```

```osprey-ml
repeat ("ab", 3)  // Success { value: "ababab" }
```
