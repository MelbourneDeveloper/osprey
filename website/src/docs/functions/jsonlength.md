---
layout: page
title: "jsonLength (Function)"
description: "Returns the number of elements in the JSON array at the given path."
---

**Signature:** `jsonLength(document: int, path: string) -> int`

**Description:** Returns the number of elements in the JSON array at the given path.

## Parameters

- **document** (int): Handle returned by jsonParse
- **path** (string): Dotted path to the array

**Returns:** int

## Example

```osprey
let n = jsonLength(doc, "items")
```
