---
layout: page
title: "writeFile (Function)"
description: "Writes content to a file. Creates the file if it doesn't exist. Returns number of bytes written."
---

**Signature:** `writeFile(filename: string, content: string) -> Result<Unit, Error>`

**Description:** Writes content to a file. Creates the file if it doesn't exist. Returns number of bytes written.

## Parameters

- **filename** (string): Path to the file to write
- **content** (string): Content to write to the file

**Returns:** Result<Unit, Error>

## Example

```osprey
let result = writeFile("output.txt", "Hello, World!")
print("File written")
```
