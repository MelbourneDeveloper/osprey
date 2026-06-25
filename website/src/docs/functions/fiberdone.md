---
layout: page
title: "fiberDone (Function)"
description: "Returns 1 if the given fiber has finished, 0 otherwise."
---

**Signature:** `fiberDone(fiber: any) -> int`

**Description:** Returns 1 if the given fiber has finished, 0 otherwise.

## Parameters

- **fiber** (any): The fiber to test

**Returns:** int

## Example

```osprey
let finished = fiberDone(worker)  // 0 or 1
```
