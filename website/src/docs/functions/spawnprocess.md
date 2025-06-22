---
layout: page
title: "spawnProcess (Function)"
description: "Spawns an external process and returns the result. Currently supports simple command execution."
---

**Signature:** `spawnProcess(command: string) -> Result<ProcessResult, string>`

**Description:** Spawns an external process and returns the result. Currently supports simple command execution.

## Parameters

- **command** (string): The command to execute

**Returns:** Result<ProcessResult, string>

## Example

```osprey
let result = spawnProcess("echo hello")
print("Command executed")
```
