---
layout: page
title: "writeFile (Function)"
description: "Writes content to a file. Creates the file if it doesn't exist."
---

**Signature:** `writeFile(filename: string, content: string) -> Result<Success, string>`

**Description:** Writes content to a file. Creates the file if it doesn't exist.

## Parameters

- **filename** (string): Path to the file to write
- **content** (string): Content to write to the file

**Returns:** Result<Success, string>

## Example

```osprey
let result = writeFile("output.txt", "Hello, World!")
print("File written")
```
