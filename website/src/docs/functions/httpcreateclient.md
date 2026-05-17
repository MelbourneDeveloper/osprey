---
layout: page
title: "httpCreateClient (Function)"
description: "Creates an HTTP client for making requests to a base URL."
---

**Signature:** `httpCreateClient(base_url: string, timeout: int) -> int`

**Description:** Creates an HTTP client for making requests to a base URL.

## Parameters

- **base_url** (string): Base URL for requests (e.g., "http://api.example.com")
- **timeout** (int): Request timeout in milliseconds

**Returns:** int

## Example

```osprey
let clientId = httpCreateClient("http://httpbin.org", 5000)
print("Client created")
```
