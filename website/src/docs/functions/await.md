---
layout: page
title: "await (Function)"
description: "Waits for a fiber to finish and returns its result, suspending the current fiber until then."
---

**Signature:** `await(fiber: Fiber<t0>) -> t0`

**Description:** Waits for a fiber to finish and returns its result, suspending the current fiber until then.

## Parameters

- **fiber** (Fiber<t0>): The fiber to await

**Returns:** t0

## Example

```osprey
let result = await(worker)
```
