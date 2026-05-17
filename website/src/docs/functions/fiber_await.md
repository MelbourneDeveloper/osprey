---
layout: page
title: "fiber_await (Function)"
description: "Waits for a fiber to complete and returns its result."
---

**Signature:** `fiber_await(fiber: Fiber) -> any`

**Description:** Waits for a fiber to complete and returns its result.

## Parameters

- **fiber** (Fiber): The fiber to await

**Returns:** any

## Example

```osprey
let result = fiber_await(fiberHandle)
```
