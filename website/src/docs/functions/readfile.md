---
layout: page
title: "readFile (Function)"
description: "Reads the entire contents of a file as a string."
---

**Signature:** `readFile(filename: string) -> Result<string, string>`

**Description:** Reads the entire contents of a file as a string.

## Parameters

- **filename** (string): Path to the file to read

**Returns:** Result<string, string>

## Example

```osprey
let content = readFile("input.txt")
print("File read")
```
