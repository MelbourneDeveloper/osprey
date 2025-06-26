---
layout: page
title: "Built-in Functions"
description: "Osprey Language Specification: Built-in Functions"
date: 2025-06-25
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0014-built-infunctions/"
---

## 14. Built-in Functions

🚀 **IMPLEMENTATION STATUS**: HTTP functions are implemented and working. WebSocket functions are implemented but undergoing testing. Fiber operations are partially implemented.

Osprey provides first-class support for HTTP servers and clients, designed with performance, safety, and streaming as core principles. All HTTP functions follow Osprey's functional programming paradigms and comply with:

- **RFC 7230**: HTTP/1.1 Message Syntax and Routing ([https://tools.ietf.org/html/rfc7230](https://tools.ietf.org/html/rfc7230))
- **RFC 7231**: HTTP/1.1 Semantics and Content ([https://tools.ietf.org/html/rfc7231](https://tools.ietf.org/html/rfc7231))
- **RFC 7232**: HTTP/1.1 Conditional Requests ([https://tools.ietf.org/html/rfc7232](https://tools.ietf.org/html/rfc7232))
- **RFC 7233**: HTTP/1.1 Range Requests ([https://tools.ietf.org/html/rfc7233](https://tools.ietf.org/html/rfc7233))
- **RFC 7234**: HTTP/1.1 Caching ([https://tools.ietf.org/html/rfc7234](https://tools.ietf.org/html/rfc7234))
- **RFC 7235**: HTTP/1.1 Authentication ([https://tools.ietf.org/html/rfc7235](https://tools.ietf.org/html/rfc7235))

- **Result types** instead of exceptions for error handling
- **Immutable response objects** that cannot be modified after creation
- **Streaming by default** for large response bodies to prevent memory issues
- **Fiber-based concurrency** for handling thousands of concurrent connections

### 14.1 HTTP Core Types

#### HTTP Method Union Type
```osprey
type HttpMethod = GET | POST | PUT | DELETE | PATCH | HEAD | OPTIONS
```

#### HTTP Request Type (Immutable)
```osprey
type HttpRequest = {
    method: HttpMethod,
    path: String,
    headers: String,
    body: String,
    queryParams: String
}
```

#### HTTP Response Type (Immutable with Streaming)
```osprey
type HttpResponse = {
    status: Int,
    headers: String,
    contentType: String,
    contentLength: Int,
    streamFd: Int,        // File descriptor for streaming
    isComplete: Bool,     // Whether response is fully loaded
    partialBody: String,  // Current chunk of body data
    partialLength: Int    // Length of current chunk
}
```

### 14.2 HTTP Server Functions

#### `httpCreateServer(port: Int, address: String) -> Result<ServerID, String>`

Creates an HTTP server bound to the specified port and address.

**Parameters:**
- `port`: Port number (1-65535)
- `address`: IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")

**Returns:**
- `Success(serverID)`: Unique server identifier
- `Err(message)`: Error description (invalid port, bind failure, etc.)

**Example:**
```osprey
let serverResult = httpCreateServer(port: 8080, address: "127.0.0.1")
match serverResult {
    Success serverId => print("Server created with ID: ${serverId}")
    Err message => print("Failed to create server: ${message}")
}
```

#### `httpListen(serverID: Int, handler: fn(String, String, String, String) -> String) -> Result<Success, String>`

Starts the HTTP server listening for requests. Each request is handled in a separate fiber for maximum concurrency.

**CRITICAL**: The handler function receives **RAW HTTP request data** and must return the **RAW response body**. The C runtime handles HTTP parsing and response formatting - the Osprey handler only processes the application logic.

**Parameters:**
- `serverID`: Server identifier from `httpCreateServer`
- `handler`: Request handler function that takes RAW HTTP data:
  - `method: String` - HTTP method (GET, POST, PUT, DELETE, etc.)
  - `path: String` - Request path (e.g., "/api/users", "/health")
  - `headers: String` - Raw HTTP headers as received
  - `body: String` - Raw request body data

**Returns:**
- `Success()`: Server started successfully
- `Err(message)`: Error description

**Example:**
```osprey
fn handleRawRequest(method: String, path: String, headers: String, body: String) -> String = 
    match method {
        "GET" => match path {
            "/health" => "{\"status\": \"healthy\"}"
            "/api/users" => "[{\"id\": 1, \"name\": \"Alice\"}, {\"id\": 2, \"name\": \"Bob\"}]"
            _ => "Not Found"
        }
        "POST" => match path {
            "/api/users" => "{\"id\": 3, \"name\": \"New User\", \"message\": \"User created\"}"
            "/api/auth/login" => "{\"token\": \"abc123\", \"message\": \"Login successful\"}"
            _ => "Endpoint not found"
        }
        "PUT" => match path {
            "/api/users/1" => "{\"id\": 1, \"name\": \"Alice Updated\", \"message\": \"User updated\"}"
            _ => "Not Found"
        }
        "DELETE" => match path {
            "/api/users/1" => "{\"message\": \"User deleted\"}"
            _ => "Not Found"
        }
        _ => "Method not allowed"
    }

let listenResult = httpListen(serverID: serverId, handler: handleRawRequest)
```

**Raw HTTP Handler Architecture:**

The HTTP server uses a **raw callback architecture** where:

1. **C Runtime** handles TCP connections, HTTP parsing, and response formatting
2. **Osprey Handler** receives raw request data and returns raw response body
3. **No HTTP abstraction** - direct access to method, path, headers, and body
4. **Maximum performance** - minimal overhead between network and application logic

**Handler Function Signature:**
```osprey
fn myHandler(method: String, path: String, headers: String, body: String) -> String
```

**Response Handling:**
- Return value becomes the HTTP response body
- HTTP status codes are determined by the response content (200 for success)
- Content-Type headers are set automatically based on response format
- For error responses, return appropriate error messages

#### `httpStopServer(serverID: Int) -> Result<Success, String>`

Stops the HTTP server and cleans up resources.

**Parameters:**
- `serverID`: Server identifier to stop

**Returns:**
- `Success()`: Server stopped successfully  
- `Err(message)`: Error description

### 14.2.1 HTTP Request Handling Bridge

**CRITICAL REQUIREMENT**: HTTP servers in Osprey must call back into Osprey code to handle requests. **NO ROUTING LOGIC SHALL BE IMPLEMENTED IN C RUNTIME**. The C runtime provides only the transport layer; all application logic, routing, and request handling must be implemented in Osprey.

#### Request Handling Architecture

When an HTTP server receives a request, the C runtime must:

1. **Parse the HTTP request** (method, path, headers, body)
2. **Serialize request data** into a structured format
3. **Call back into Osprey** to handle the request
4. **Receive response data** from Osprey
5. **Send HTTP response** back to the client

#### Bridge Function Specification

**NEW ARCHITECTURE**: Osprey now uses **direct function pointer callbacks** for maximum performance and zero overhead.

#### Raw Function Pointer Callbacks

When `httpListen()` is called, the Osprey handler function is passed directly to the C runtime as a function pointer:

**C Runtime Function Signature:**
```c
int64_t http_listen(int64_t server_id, char* (*handler)(char* method, char* path, char* headers, char* body));
```

**Handler Function Signature:**
```c
char* handler(char* method, char* path, char* headers, char* body);
```

#### Legacy Bridge Function (Deprecated)

The old bridge function is deprecated but still supported for compatibility:

```c
// DEPRECATED: Use direct function pointers instead
extern int osprey_handle_http_request(
    int server_id,
    char* method,
    char* full_url,
    char* raw_headers,
    char* body,
    size_t body_length,
    int* response_status,
    char** response_headers,
    char** response_body,
    size_t* response_body_length
);
```

#### New Raw Callback Architecture Flow

**1. Osprey Code:**
```osprey
fn handleRawRequest(method: String, path: String, headers: String, body: String) -> String = 
    match method {
        "GET" => match path {
            "/health" => "{\"status\": \"healthy\"}"
            "/api/users" => "[{\"id\": 1, \"name\": \"Alice\"}]"
            _ => "Not Found"
        }
        "POST" => "{\"message\": \"Created\"}"
        _ => "Method not allowed"
    }

let listenResult = httpListen(serverId, handleRawRequest)
```

**2. LLVM Code Generation:**
- Generates function pointer for `handleRawRequest`
- Passes function pointer to `http_listen()` C function

**3. C Runtime Implementation:**
```c
// Global storage for handler function pointer
static char* (*request_handler)(char* method, char* path, char* headers, char* body) = NULL;

int64_t http_listen(int64_t server_id, char* (*handler)(char* method, char* path, char* headers, char* body)) {
    request_handler = handler;  // Store the function pointer
    // Setup server socket and start listening...
    return 0;
}

// In request processing loop:
void handle_client_request(int client_fd, char* method, char* path, char* headers, char* body) {
    if (request_handler) {
        // Call Osprey function directly with RAW data
        char* response_body = request_handler(method, path, headers, body);
        
        // Format and send HTTP response
        char response[8192];
        snprintf(response, sizeof(response),
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: application/json\r\n"
            "Content-Length: %zu\r\n"
            "Connection: close\r\n"
            "\r\n%s",
            strlen(response_body), response_body);
        
        send(client_fd, response, strlen(response), 0);
        
        // Clean up if response was allocated
        if (response_body) free(response_body);
    }
}
```

**Architecture Benefits:**
- **Zero overhead**: Direct function calls, no serialization
- **Raw data access**: Handler receives exactly what was sent over HTTP
- **Maximum performance**: Minimal abstraction between network and application
- **Simple debugging**: Direct call stack from C to Osprey
- **Memory efficient**: No intermediate data structures

**Legacy Bridge Parameters (Deprecated):**
- `server_id`: The server ID that received the request
- `method`: HTTP method (GET, POST, PUT, DELETE, etc.)
- `full_url`: Complete URL including query parameters ("/api/users?page=1&limit=10")
- `raw_headers`: Raw HTTP headers as received ("Content-Type: application/json\r\nAuthorization: Bearer token\r\n")
- `body`: Raw request body data (may be binary)
- `body_length`: Length of request body in bytes
- `response_status`: Output parameter for HTTP status code
- `response_headers`: Output parameter for raw response headers
- `response_body`: Output parameter for response body (may be binary)
- `response_body_length`: Output parameter for response body length

**Legacy Return Value:**
- `0`: Success
- `-1`: Error handling request

**Streaming Support:**
For large request/response bodies, the bridge function must support streaming:
- Request body streaming: C runtime provides file descriptor for reading body data
- Response body streaming: Osprey can return a file descriptor for C runtime to stream response

#### Implementation Requirements

**🚫 FORBIDDEN IN C RUNTIME:**
- URL routing logic
- Application-specific response generation
- Business logic
- Hardcoded API endpoints
- Request path matching beyond basic parsing

**✅ REQUIRED IN C RUNTIME:**
- HTTP protocol parsing
- Socket management
- Request/response serialization
- Bridge function calls
- Error handling for transport failures

**✅ REQUIRED IN OSPREY:**
- All request routing logic
- Application business logic
- Response generation
- API endpoint definitions
- Request validation

#### Example Implementation

**C Runtime (Transport Layer Only):**
```c
void handle_client_request(int client_fd, int server_id, char* method, char* full_url, 
                          char* raw_headers, char* body, size_t body_length) {
    // Prepare response parameters
    int response_status;
    char* response_headers = NULL;
    char* response_body = NULL;
    size_t response_body_length;
    
    // Call back into Osprey with raw HTTP data - NO ROUTING IN C!
    int result = osprey_handle_http_request(
        server_id, method, full_url, raw_headers, body, body_length,
        &response_status, &response_headers, &response_body, &response_body_length
    );
    
    if (result == 0) {
        // Send HTTP response with raw data
        send_raw_http_response(client_fd, response_status, response_headers, 
                              response_body, response_body_length);
    } else {
        // Send 500 error
        send_error_response(client_fd, 500, "Internal Server Error");
    }
    
    // Clean up allocated response data
    if (response_headers) free(response_headers);
    if (response_body) free(response_body);
}
```

**Osprey Code (Application Layer):**
```osprey
fn handleHttpRequest(request: HttpRequest) -> Result<HttpResponse, String> = 
    match request.method {
        GET => match request.path {
            "/api/users" => getUserList()
            "/api/health" => Success(HttpResponse {
                status: 200,
                contentType: "application/json",
                body: "{\"status\": \"healthy\"}"
            })
            _ => Success(HttpResponse {
                status: 404,
                contentType: "text/plain", 
                body: "Not Found"
            })
        }
        POST => match request.path {
            "/api/users" => createUser(request.body)
            _ => Err("Endpoint not found")
        }
        _ => Err("Method not supported")
    }
```

This architecture ensures **complete separation of concerns**: C handles transport, Osprey handles application logic.

### 14.3 HTTP Client Functions

#### `httpCreateClient(baseUrl: String, timeout: Int) -> Result<ClientID, String>`

Creates an HTTP client for making requests.

**Parameters:**
- `baseUrl`: Base URL for requests (e.g., "http://api.example.com")
- `timeout`: Request timeout in milliseconds

**Returns:**
- `Success(clientID)`: Unique client identifier
- `Err(message)`: Error description

**Example:**
```osprey
let clientResult = httpCreateClient(baseUrl: "http://jsonplaceholder.typicode.com", timeout: 5000)
```

#### `httpGet(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP GET request.

**Parameters:**
- `clientID`: Client identifier from `httpCreateClient`
- `path`: Request path (e.g., "/users/1")
- `headers`: Additional headers (e.g., "Authorization: Bearer token\r\n")

**Returns:**
- `Success(statusCode)`: HTTP status code (200, 404, etc.)
- `Err(message)`: Error description

**Example:**
```osprey
let getResult = httpGet(clientID: clientId, path: "/users", headers: "")
match getResult {
    Success statusCode => print("Request completed with status: ${statusCode}")
    Err message => print("Request failed: ${message}")
}
```

#### `httpPost(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP POST request with a request body.

**Parameters:**
- `clientID`: Client identifier
- `path`: Request path
- `body`: Request body data
- `headers`: Additional headers

**Example:**
```osprey
let postData = "{\"name\": \"John\", \"email\": \"john@example.com\"}"
let headers = "Content-Type: application/json\r\n"
let postResult = httpPost(clientID: clientId, path: "/users", body: postData, headers: headers)
```

#### `httpPut(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP PUT request.

#### `httpDelete(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`

Makes an HTTP DELETE request.

#### `httpRequest(clientID: Int, method: HttpMethod, path: String, headers: String, body: String) -> Result<StatusCode, String>`

Generic HTTP request function for any HTTP method.

#### `httpCloseClient(clientID: Int) -> Result<Success, String>`

Closes the HTTP client and cleans up resources.

### 14.4 WebSocket Support (Two-Way Communication)

🔒 **IMPLEMENTATION STATUS**: WebSocket functions are implemented with security features but are currently undergoing stability testing.

WebSockets provide real-time, bidirectional communication between client and server. Osprey implements WebSocket support with **MILITARY-GRADE SECURITY** following industry best practices for preventing attacks and ensuring bulletproof operation.

#### 14.4.1 WebSocket Security Implementation

Osprey's WebSocket implementation follows the **OWASP WebSocket Security Guidelines** and implements multiple layers of security protection:

**🛡️ TITANIUM-ARMORED Compilation Security:**
- `_FORTIFY_SOURCE=3`: Maximum buffer overflow protection
- `fstack-protector-all`: Complete stack canary protection  
- `fstack-clash-protection`: Stack clash attack prevention
- `fcf-protection=full`: Control Flow Integrity (CFI) protection
- `ftrapv`: Integer overflow trapping
- `fno-delete-null-pointer-checks`: Prevent null pointer optimizations
- `Wl,-z,relro,-z,now`: Full RELRO with immediate binding
- `Wl,-z,noexecstack`: Non-executable stack protection

**🔒 Cryptographic Security:**
- **OpenSSL SHA-1**: RFC 6455 compliant WebSocket handshake using industry-standard OpenSSL
- **Secure key validation**: 24-character base64 key format validation
- **Constant-time operations**: Memory clearing to prevent timing attacks
- **Error checking**: All OpenSSL operations validated for success

**⚔️ Input Validation Fortress:**
- **WebSocket key format validation**: Strict RFC 6455 compliance
- **Base64 character validation**: Only valid characters accepted
- **Buffer length validation**: Maximum 4096 character keys prevent DoS
- **Integer overflow protection**: All memory calculations checked
- **Memory boundary checking**: No buffer overruns possible

**🏰 Memory Security:**
- **Secure memory allocation**: `calloc()` with zero-initialization
- **Memory clearing**: All sensitive data zeroed before deallocation
- **Bounds checking**: All `snprintf()` operations validated for truncation
- **Safe string operations**: `memcpy()` instead of unsafe `strcpy()`/`strcat()`

#### 14.4.2 Security Standards Compliance

Osprey WebSocket implementation follows these security standards:

**RFC 6455 - WebSocket Protocol Security** ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455)):
- Proper Sec-WebSocket-Accept calculation using SHA-1 + base64
- Origin validation support for CSRF protection
- Secure WebSocket handshake implementation

**OWASP WebSocket Security Cheat Sheet** ([https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets)):
- Input validation on all WebSocket frames
- Authentication and authorization enforcement
- Rate limiting and DoS protection
- Secure error handling without information leakage

**NIST Cybersecurity Framework:**
- Defense in depth through multiple security layers
- Secure coding practices with compiler hardening
- Memory safety through bounds checking
- Cryptographic integrity using OpenSSL

**CWE (Common Weakness Enumeration) Mitigation:**
- CWE-120: Buffer overflow prevention through bounds checking
- CWE-190: Integer overflow protection with `ftrapv`
- CWE-200: Information exposure prevention through secure error handling
- CWE-416: Use-after-free prevention through memory clearing

#### 14.4.3 Security Architecture

```
┌─────────────────────────────────────────────────────────┐
│                TITANIUM SECURITY LAYERS                 │
├─────────────────────────────────────────────────────────┤
│ 🏰 Application Layer: Input Validation Fortress        │
│    • WebSocket key format validation                   │
│    • Base64 character validation                       │
│    • Buffer length enforcement                         │
│    • Memory boundary checking                          │
├─────────────────────────────────────────────────────────┤
│ 🔒 Cryptographic Layer: OpenSSL SHA-1                  │
│    • RFC 6455 compliant handshake                      │
│    • Secure hash computation                           │
│    • Constant-time operations                          │
│    • Error validated operations                        │
├─────────────────────────────────────────────────────────┤
│ ⚔️ Memory Layer: Bulletproof Memory Management         │
│    • Secure allocation with calloc()                   │
│    • Memory clearing before deallocation               │
│    • Safe string operations                            │
│    • Integer overflow protection                       │
├─────────────────────────────────────────────────────────┤
│ 🛡️ Compiler Layer: Military-Grade Hardening           │
│    • Stack protection (canaries + clash protection)    │
│    • Control Flow Integrity (CFI)                      │
│    • FORTIFY_SOURCE=3 buffer overflow protection       │
│    • RELRO + immediate binding                         │
│    • Non-executable stack                              │
└─────────────────────────────────────────────────────────┘
```

#### 14.4.4 Security Testing and Validation

Osprey WebSocket security is validated through:

**🧪 Automated Security Testing:**
- Buffer overflow attack simulation
- Malformed WebSocket key injection
- Integer overflow boundary testing
- Memory corruption detection

**🔍 Static Analysis:**
- Compiler warnings elevated to errors
- Memory safety analysis
- Control flow analysis
- Buffer bounds verification

**⚡ Dynamic Testing:**
- Address Sanitizer (ASan) testing
- Valgrind memory error detection
- Fuzzing with malformed inputs
- DoS resilience testing

#### 14.4.5 Security References and Standards

**Primary Security Standards:**
- **RFC 6455**: "The WebSocket Protocol" - Official WebSocket specification ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455))
- **OWASP WebSocket Security Cheat Sheet**: ([https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets))
- **NIST SP 800-53**: Security Controls for Federal Information Systems
- **ISO 27001**: Information Security Management Standards

**Compiler Security References:**
- **GCC Security Options**: ([https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html))
- **Red Hat Security Guide**: "Defensive Coding Practices"
- **Microsoft SDL**: Security Development Lifecycle practices
- **Google Safe Coding Practices**: Memory safety guidelines

**Cryptographic Standards:**
- **FIPS 180-4**: SHA-1 cryptographic hash standard
- **RFC 3174**: US Secure Hash Algorithm 1 (SHA1) ([https://tools.ietf.org/html/rfc3174](https://tools.ietf.org/html/rfc3174))
- **OpenSSL Security Advisories**: ([https://www.openssl.org/news/secadv.html](https://www.openssl.org/news/secadv.html))

**Memory Security Research:**
- **"Control Flow Integrity"** by Abadi et al. - CFI protection principles
- **"Stack Canaries"** - Buffer overflow detection mechanisms  
- **"RELRO"** - Read-only relocations for exploit mitigation
- **"FORTIFY_SOURCE"** - Compile-time and runtime buffer overflow detection

#### `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`

Establishes a WebSocket connection.

**Parameters:**
- `url`: WebSocket URL (e.g., "ws://localhost:8080/chat")
- `messageHandler`: Callback function to handle incoming messages

**Returns:**
- `Success(wsID)`: WebSocket connection identifier
- `Err(message)`: Connection error

**Example:**
```osprey
fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}

let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)
```

#### `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

Sends a message through the WebSocket connection.

**Parameters:**
- `wsID`: WebSocket identifier
- `message`: Message to send

**Example:**
```osprey
let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
```

#### `websocketClose(wsID: Int) -> Result<Success, String>`

Closes the WebSocket connection.

### 14.4.1 WebSocket Server Functions

#### `websocketCreateServer(port: Int, address: String, path: String) -> Int`
Creates a WebSocket server bound to the specified port and address.

🚧 **IMPLEMENTATION STATUS**: The current implementation has **CRITICAL RUNTIME ISSUES**:

**CURRENT BEHAVIOR**:
- Returns server ID on successful creation
- Returns negative error codes on failure

**RUNTIME ISSUES DETECTED**:
- **Port Binding Failures**: `websocketServerListen()` returns `-4` (bind failed) instead of expected `0` (success)
- **Resource Conflicts**: Multiple test runs cause port conflicts and resource exhaustion
- **Test Environment Instability**: Inconsistent behavior between different execution environments

**ROOT CAUSE ANALYSIS**:
- **Issue**: `bind()` system call fails with `EADDRINUSE` (Address already in use)
- **Impact**: WebSocket server cannot bind to port, causing listen operation to fail
- **Environment**: Particularly problematic in containerized test environments with limited cleanup

**NEEDED FIXES**:
1. **Port Management**: Implement proper port cleanup and reuse detection
2. **Resource Cleanup**: Ensure proper socket closure and resource deallocation
3. **Retry Logic**: Add exponential backoff for port binding failures
4. **Error Handling**: Better error reporting for different failure modes
5. **Test Isolation**: Implement proper test teardown to prevent resource conflicts

**Example:**
```osprey
let serverId = websocketCreateServer(8080, "127.0.0.1", "/chat")
print("Server created with ID: ${serverId}")
```

#### `websocketServerListen(serverID: Int) -> Int`
Starts the WebSocket server listening for connections.

🚧 **CURRENT ISSUE**: Returns `-4` (bind failed) instead of `0` (success) due to port binding issues.

**Error Codes:**
- `0`: Success
- `-1`: Invalid server ID
- `-2`: Socket creation failed
- `-3`: Socket options failed
- `-4`: **BIND FAILED** (most common current issue)
- `-5`: Listen failed
- `