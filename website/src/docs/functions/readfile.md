---
layout: page
title: "readFile (Function)"
description: "Reads the entire contents of a file as a string."
---

**Signature:** `readFile(filename: string) -> Result<string, Error>`

**Description:** Reads the entire contents of a file as a string.

## Parameters

- **filename** (string): Path to the file to read

**Returns:** Result<string, Error>

## Example

```osprey
let content = readFile("input.txt")
print("File read")
```

```osprey-ml
content = readFile "input.txt"
print "File read"
```
