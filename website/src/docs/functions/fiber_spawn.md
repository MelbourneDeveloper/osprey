---
layout: page
title: "fiber_spawn (Function)"
description: "Spawns a new fiber to execute the given function concurrently."
---

**Signature:** `fiber_spawn(fn: () -> any) -> Fiber`

**Description:** Spawns a new fiber to execute the given function concurrently.

## Parameters

- **fn** (() -> any): The function to execute in the fiber

**Returns:** Fiber

## Example

```osprey
let fiber = fiber_spawn(() -> print("Hello from fiber"))
```
