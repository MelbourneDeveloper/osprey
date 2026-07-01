---
layout: page
title: "httpResponseFree (Function)"
description: "Releases a response handle obtained from httpGetResponse."
---

**Signature:** `httpResponseFree(responseID: int) -> Unit`

**Description:** Releases a response handle obtained from httpGetResponse.

## Parameters

- **responseID** (int): Handle returned by httpGetResponse

**Returns:** Unit

## Example

```osprey
httpResponseFree(response)
```

```osprey-ml
httpResponseFree response
```
