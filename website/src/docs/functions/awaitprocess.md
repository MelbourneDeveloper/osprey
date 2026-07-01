---
layout: page
title: "awaitProcess (Function)"
description: "Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes."
---

**Signature:** `awaitProcess(handle: int) -> int`

**Description:** Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes.

## Parameters

- **handle** (int): Process ID from spawnProcess

**Returns:** int

## Example

```osprey
let exitCode = awaitProcess(processHandle)
print("Process exited with code: ${toString(exitCode)}")
```

```osprey-ml
exitCode = awaitProcess processHandle
print "Process exited with code: ${toString exitCode}"
```
