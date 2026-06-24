---
layout: page
title: "httpResponseBody (Function)"
description: "Returns the body of a response handle as a string."
---

**Signature:** `httpResponseBody(responseID: int) -> Result<string, Error>`

**Description:** Returns the body of a response handle as a string.

## Parameters

- **responseID** (int): Handle returned by httpGetResponse

**Returns:** Result<string, Error>

## Example

```osprey
match httpResponseBody(response) {
  Success { value } => print(value)
  Error { message } => print(message)
}
```
