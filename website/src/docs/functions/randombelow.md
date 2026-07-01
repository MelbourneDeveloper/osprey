---
layout: page
title: "randomBelow (Function)"
description: "A cryptographically-secure uniform random integer in [0, n), unbiased by rejection sampling. Returns Result<int, MathError> — Error when n <= 0."
---

**Signature:** `randomBelow(n: int) -> Result<int, Error>`

**Description:** A cryptographically-secure uniform random integer in [0, n), unbiased by rejection sampling. Returns Result<int, MathError> — Error when n <= 0.

## Parameters

- **n** (int): Exclusive upper bound; must be positive

**Returns:** Result<int, Error>

## Example

```osprey
let d = randomBelow(6) ?: 0  // a fair die face 0..5
```

```osprey-ml
d = randomBelow 6 ?: 0  // a fair die face 0..5
```
