---
layout: page
title: "deleteFile (Function)"
description: "Deletes the file at the given path, returning Unit on success or an error."
---

**Signature:** `deleteFile(path: string) -> Result<Unit, Error>`

**Description:** Deletes the file at the given path, returning Unit on success or an error.

## Parameters

- **path** (string): Filesystem path of the file to delete

**Returns:** Result<Unit, Error>

## Example

```osprey
match deleteFile("temp.txt") {
  Success { value } => print("deleted")
  Error { message } => print(message)
}
```
