---
layout: page
title: "ProcessHandle (Type)"
description: "A handle to a spawned async process. Contains the process ID and allows waiting for completion and cleanup. Process output is delivered via callbacks registered with the runtime."
---

**Description:** A handle to a spawned async process. Contains the process ID and allows waiting for completion and cleanup. Process output is delivered via callbacks registered with the runtime.

## Example

```osprey
let result = spawnProcess("echo hello")
match result {
    Success { value } => {
        let exitCode = awaitProcess(value)
        cleanupProcess(value)
    }
    Error { message } => print("Process failed")
}
```
