---
layout: page
title: "jsonFree (Function)"
description: "Releases a parsed JSON document handle obtained from jsonParse."
---

**Signature:** `jsonFree(document: int) -> Unit`

**Description:** Releases a parsed JSON document handle obtained from jsonParse.

## Parameters

- **document** (int): Handle returned by jsonParse

**Returns:** Unit

## Example

```osprey
jsonFree(doc)
```
