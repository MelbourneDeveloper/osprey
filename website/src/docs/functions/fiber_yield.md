---
layout: page
title: "fiber_yield (Function)"
description: "Yields control to the fiber scheduler with an optional value."
---

**Signature:** `fiber_yield(value: int) -> int`

**Description:** Yields control to the fiber scheduler with an optional value.

## Parameters

- **value** (int): The value to yield

**Returns:** int

## Example

```osprey
let result = fiber_yield(42)
```

```osprey-ml
result = fiber_yield 42
```
