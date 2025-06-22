---
layout: page
title: "parseJSON (Function)"
description: "Parses a JSON string and returns the parsed result."
---

**Signature:** `parseJSON(json: string) -> Result<string, string>`

**Description:** Parses a JSON string and returns the parsed result.

## Parameters

- **json** (string): JSON string to parse

**Returns:** Result<string, string>

## Example

```osprey
let parsed = parseJSON("{\"key\": \"value\"}")
print("JSON parsed")
```
