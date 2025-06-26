---
layout: page
title: "cleanupProcess (Function)"
description: "Cleans up resources associated with a completed process. Should be called after awaitProcess."
---

**Signature:** `cleanupProcess(handle: ProcessHandle) -> void`

**Description:** Cleans up resources associated with a completed process. Should be called after awaitProcess.

## Parameters

- **handle** (ProcessHandle): Process handle from spawnProcess

**Returns:** void

## Example

```osprey
cleanupProcess(processHandle)  // Free process resources
```
