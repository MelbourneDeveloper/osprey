---
layout: page
title: "spawnProcess (Function)"
description: "Spawns an external async process with MANDATORY callback for stdout/stderr capture. The callback function receives (processID: int, eventType: int, data: string) and is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!"
---

**Signature:** `spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessHandle, string>`

**Description:** Spawns an external async process with MANDATORY callback for stdout/stderr capture. The callback function receives (processID: int, eventType: int, data: string) and is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!

## Parameters

- **command** (string): The command to execute
- **callback** (fn(int, int, string) -> Unit): MANDATORY callback function for process events (processID, eventType, data)

**Returns:** Result<ProcessHandle, string>

## Example

```osprey
fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {
    match eventType {
        1 => print("STDOUT: ${data}")
        2 => print("STDERR: ${data}")
        3 => print("EXIT: ${data}")
        _ => print("Unknown event")
    }
}
let result = spawnProcess("echo hello", processEventHandler)
match result {
    Success { value } => {
        let exitCode = awaitProcess(value)
        cleanupProcess(value)
    }
    Error { message } => print("Failed")
}
```
