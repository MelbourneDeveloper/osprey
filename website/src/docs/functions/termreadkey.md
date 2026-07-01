---
layout: page
title: "termReadKey (Function)"
description: "Reads a single keypress from the terminal and returns it as a string."
---

**Signature:** `termReadKey() -> Result<string, Error>`

**Description:** Reads a single keypress from the terminal and returns it as a string.

**Returns:** Result<string, Error>

## Example

```osprey
match termReadKey() {
  Success { value } => print("key: ${value}")
  Error { message } => print(message)
}
```

```osprey-ml
match termReadKey
    Success value => print "key: ${value}"
    Error message => print message
```
