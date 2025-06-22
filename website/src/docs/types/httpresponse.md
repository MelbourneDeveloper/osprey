---
layout: page
title: "HttpResponse (Type)"
description: "A built-in type representing an HTTP response with status code, headers, content type, body, and streaming capabilities. Used by HTTP server handlers to return structured responses to clients."
---

**Description:** A built-in type representing an HTTP response with status code, headers, content type, body, and streaming capabilities. Used by HTTP server handlers to return structured responses to clients.

## Example

```osprey
HttpResponse {
    status: 200,
    headers: "Content-Type: application/json",
    contentType: "application/json",
    contentLength: 25,
    streamFd: -1,
    isComplete: true,
    partialBody: "{\"message\": \"Hello\"}",
    partialLength: 25
}
```
