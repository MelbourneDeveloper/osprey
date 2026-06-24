---
layout: page
title: "jsonGet (Function)"
description: "Returns the string value at the given path within a parsed JSON document."
---

**Signature:** `jsonGet(document: int, path: string) -> Result<string, Error>`

**Description:** Returns the string value at the given path within a parsed JSON document.

## Parameters

- **document** (int): Handle returned by jsonParse
- **path** (string): Dotted path to the value, e.g. "user.name"

**Returns:** Result<string, Error>

## Example

```osprey
match jsonGet(doc, "name") {
  Success { value } => print(value)
  Error { message } => print(message)
}
```
