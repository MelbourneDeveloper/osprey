---
layout: page
title: "fiber_yield (Function)"
description: "Yields control to the fiber scheduler with an optional value."
---

**Signature:** `fiber_yield(value: any) -> any`

**Description:** Yields control to the fiber scheduler with an optional value.

## Parameters

- **value** (any): The value to yield

**Returns:** any

## Example

```osprey
let result = fiber_yield(42)
```
