---
layout: page
title: "jsonParse (Function)"
description: "Parses a JSON string and returns an opaque document handle for querying, or an error on malformed input."
---

**Signature:** `jsonParse(text: string) -> Result<int, Error>`

**Description:** Parses a JSON string and returns an opaque document handle for querying, or an error on malformed input.

## Parameters

- **text** (string): The JSON text to parse

**Returns:** Result<int, Error>

## Example

```osprey
match jsonParse("{\"name\": \"osprey\"}") {
  Success { value } => print("parsed")
  Error { message } => print(message)
}
```
