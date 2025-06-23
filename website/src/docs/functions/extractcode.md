---
layout: page
title: "extractCode (Function)"
description: "Extracts code from a JSON structure."
---

**Signature:** `extractCode(json: string) -> Result<string, string>`

**Description:** Extracts code from a JSON structure.

## Parameters

- **json** (string): JSON string containing code

**Returns:** Result<string, string>

## Example

```osprey
let code = extractCode("{\"code\": \"print(42)\"}")
print("Code extracted")
```
