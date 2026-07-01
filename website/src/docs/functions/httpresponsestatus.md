---
layout: page
title: "httpResponseStatus (Function)"
description: "Returns the HTTP status code of a response handle."
---

**Signature:** `httpResponseStatus(responseID: int) -> int`

**Description:** Returns the HTTP status code of a response handle.

## Parameters

- **responseID** (int): Handle returned by httpGetResponse

**Returns:** int

## Example

```osprey
let code = httpResponseStatus(response)  // 200
```

```osprey-ml
code = httpResponseStatus response  // 200
```
