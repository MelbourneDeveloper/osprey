---
layout: page
title: "cleanupProcess (Function)"
description: "Cleans up resources associated with a completed process. Should be called after awaitProcess."
---

**Signature:** `cleanupProcess(handle: int) -> Unit`

**Description:** Cleans up resources associated with a completed process. Should be called after awaitProcess.

## Parameters

- **handle** (int): Process ID from spawnProcess

**Returns:** Unit

## Example

```osprey
cleanupProcess(processHandle)  // Free process resources
```
