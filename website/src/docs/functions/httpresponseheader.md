---
layout: page
title: "httpResponseHeader (Function)"
description: "Returns the value of the named header from a response handle."
---

**Signature:** `httpResponseHeader(responseID: int, name: string) -> Result<string, Error>`

**Description:** Returns the value of the named header from a response handle.

## Parameters

- **responseID** (int): Handle returned by httpGetResponse
- **name** (string): Header name, e.g. "Content-Type"

**Returns:** Result<string, Error>

## Example

```osprey
match httpResponseHeader(response, "Content-Type") {
  Success { value } => print(value)
  Error { message } => print(message)
}
```
