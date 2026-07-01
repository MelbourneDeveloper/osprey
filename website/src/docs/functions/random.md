---
layout: page
title: "random (Function)"
description: "A cryptographically-secure uniform random non-negative integer (0 .. 2^63-1), drawn fresh from the OS entropy source. Unseeded and unpredictable."
---

**Signature:** `random() -> int`

**Description:** A cryptographically-secure uniform random non-negative integer (0 .. 2^63-1), drawn fresh from the OS entropy source. Unseeded and unpredictable.

**Returns:** int

## Example

```osprey
let big = random()  // e.g. 7240982340198
```

```osprey-ml
big = random  // e.g. 7240982340198
```
