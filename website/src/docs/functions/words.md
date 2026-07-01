---
layout: page
title: "words (Function)"
description: "Splits on runs of whitespace; empty results dropped."
---

**Signature:** `words(s: string) -> List<string>`

**Description:** Splits on runs of whitespace; empty results dropped.

## Parameters

- **s** (string): The string to split

**Returns:** List<string>

## Example

```osprey
words("a  b\\tc")  // ["a","b","c"]
```

```osprey-ml
words "a  b\\tc"  // ["a","b","c"]
```
