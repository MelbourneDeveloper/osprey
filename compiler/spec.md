**Version:** 0.2.0-alpha

**Date:** June 2025

**Author** Christian Findlay

<div class="table-of-contents">

1. [Introduction](#1-introduction)
   - [Completeness](#11-completeness)
   - [Principles](#12-principles)
2. [Lexical Structure](#2-lexical-structure)
   - [Identifiers](#21-identifiers)
   - [Keywords](#22-keywords)
   - [Literals](#23-literals)
   - [Operators](#24-operators)
   - [Delimiters](#25-delimiters)
3. [Syntax](#3-syntax)
   - [Program Structure](#31-program-structure)
   - [Import Statements](#32-import-statements)
   - [Let Declarations](#33-let-declarations)
   - [Function Declarations](#34-function-declarations)
   - [Extern Declarations](#35-extern-declarations)
   - [Type Declarations](#36-type-declarations)
   - [Record Types and Type Constructors](#37-record-types-and-type-constructors)
   - [Expressions](#38-expressions)
   - [Block Expressions](#39-block-expressions)
   - [Match Expressions](#310-match-expressions)
4. [Semantics](#4-semantics)
5. [Type System](#5-type-system)
   - [Built-in Types](#51-built-in-types)
   - [Built-in Error Types](#52-built-in-error-types)
   - [Type Inference Rules](#53-type-inference-rules)
   - [Type Safety and Explicit Typing](#54-type-safety-and-explicit-typing)
   - [Any Type Handling and Pattern Matching Requirement](#55-any-type-handling-and-pattern-matching-requirement)
   - [Type Compatibility](#56-type-compatibility)
6. [Function Calls](#6-function-calls)
7. [String Interpolation](#7-string-interpolation)
8. [Pattern Matching](#8-pattern-matching)
9. [Block Expressions](#9-block-expressions)
10. [Boolean Operations](#10-boolean-operations)
11. [Loop Constructs and Functional Iterators](#11-loop-constructs-and-functional-iterators)
12. [Lightweight Fibers and Concurrency](#12-lightweight-fibers-and-concurrency)
13. [Fiber-Isolated Module System](#13-fiber-isolated-module-system)
14. [Built-in Functions](#14-built-in-functions)
    - [HTTP Core Types](#141-http-core-types)
    - [HTTP Server Functions](#142-http-server-functions)
    - [HTTP Client Functions](#143-http-client-functions)
    - [WebSocket Support](#144-websocket-support-two-way-communication)
    - [Streaming Response Bodies](#145-streaming-response-bodies)
    - [Error Handling in HTTP](#146-error-handling-in-http)
    - [Fiber-Based Concurrency](#147-fiber-based-concurrency)
    - [Complete HTTP Server Example](#148-complete-http-server-example)
15. [Error Handling](#15-error-handling)
    - [The Result Type](#151-the-result-type)
16. [Examples](#16-examples)
17. [Built-in Functions Reference](#17-built-in-functions-reference)
    - [Basic I/O Functions](#171-basic-io-functions)
    - [Functional Iterator Functions](#172-functional-iterator-functions)
    - [Pipe Operator](#173-pipe-operator)
    - [Functional Programming Patterns](#174-functional-programming-patterns)
    - [Fiber Types and Concurrency](#175-fiber-types-and-concurrency)
18. [Security and Sandboxing](#18-security-and-sandboxing)
    - [Security Flags](#181-security-flags)
    - [Security Policies](#182-security-policies)
    - [Blocked Functions by Category](#183-blocked-functions-by-category)
    - [Error Messages](#184-error-messages)
    - [Programming Best Practices](#185-programming-best-practices)
    - [Implementation Details](#186-implementation-details)

</div>













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
- `-6`: Thread creation failed

#### `websocketServerBroadcast(serverID: Int, message: String) -> Int`
Broadcasts a message to all connected WebSocket clients.

#### `websocketStopServer(serverID: Int) -> Int`
Stops the WebSocket server and closes all connections.

#### `websocketKeepAlive() -> Void`
Keeps the WebSocket server running indefinitely until interrupted.

### 14.5 Streaming Response Bodies

Osprey automatically handles response streaming to prevent memory issues with large responses:

#### Complete Response
For small responses that fit in memory:
```osprey
HttpResponse {
    status: 200,
    contentType: "application/json",
    partialBody: "{\"data\": \"small response\"}",
    isComplete: true,
    streamFd: -1
}
```

#### Streamed Response
For large responses that should be streamed:
```osprey
HttpResponse {
    status: 200,
    contentType: "application/octet-stream",
    streamFd: fileDescriptor,  // File descriptor to stream from
    isComplete: false,
    contentLength: -1,         // -1 indicates chunked encoding
    partialBody: "",
    partialLength: 0
}
```

### 14.6 Error Handling in HTTP

All HTTP functions return Result types following Osprey's principle of making illegal states unrepresentable:

```osprey
// Server creation with error handling
let serverResult = httpCreateServer(port: 8080, address: "127.0.0.1")
match serverResult {
    Success serverId => {
        print("Server created successfully")
        let listenResult = httpListen(serverID: serverId, handler: myHandler)
        match listenResult {
            Success _ => print("Server is listening")
            Err error => print("Failed to start listening: ${error}")
        }
    }
    Err error => print("Failed to create server: ${error}")
}
```

### 14.7 Fiber-Based Concurrency

HTTP servers automatically spawn a new fiber for each incoming request, enabling thousands of concurrent connections:

```osprey
// Each request handler runs in its own fiber
fn handleRequest(request: HttpRequest) -> Result<HttpResponse, String> = {
    // This function runs in a separate fiber for each request
    // Multiple requests can be processed concurrently
    let result = processLongRunningTask(request.body)
    Success(HttpResponse {
        status: 200,
        contentType: "application/json",
        partialBody: result,
        isComplete: true,
        streamFd: -1
    })
}
```

### 14.8 Complete HTTP Server Example

```osprey
// Create and start an HTTP server
let serverResult = httpCreateServer(port: 8080, address: "0.0.0.0")
match serverResult {
    Success serverId => {
        fn apiHandler(request: HttpRequest) -> Result<HttpResponse, String> = match request.method {
            GET => match request.path {
                "/api/health" => Success(HttpResponse {
                    status: 200,
                    contentType: "application/json",
                    partialBody: "{\"status\": \"healthy\", \"timestamp\": \"${getCurrentTime()}\"}",
                    isComplete: true,
                    streamFd: -1
                })
                "/api/users" => Success(HttpResponse {
                    status: 200,
                    contentType: "application/json",
                    partialBody: getUsersJson(),
                    isComplete: true,
                    streamFd: -1
                })
                _ => Success(HttpResponse {
                    status: 404,
                    contentType: "application/json", 
                    partialBody: "{\"error\": \"Not Found\"}",
                    isComplete: true,
                    streamFd: -1
                })
            }
            POST => match request.path {
                "/api/users" => createUser(request.body)
                _ => Err("Unsupported POST endpoint")
            }
            _ => Err("Method not allowed")
        }
        
        let listenResult = httpListen(serverID: serverId, handler: apiHandler)
        match listenResult {
            Success _ => print("🚀 HTTP server listening on http://0.0.0.0:8080")
            Err error => print("❌ Failed to start server: ${error}")
        }
    }
    Err error => print("❌ Failed to create server: ${error}")
}
```

## 15. Error Handling

### 15.1 The Result Type

**CRITICAL**: All functions that can fail **MUST** return a `Result` type. There are no exceptions, panics, or nulls. This is a core design principle of the language to ensure safety and eliminate entire classes of runtime errors.

The `Result` type is a generic union type with two variants:

- `Success { value: T }`: Represents a successful result, containing the value of type `T`.
- `Error { message: E }`: Represents an error, containing an error message or object of type `E`.

**Example:**
```osprey
type Result<T, E> = Success { value: T } | Error { message: E }
```

The compiler **MUST** enforce that `Result` types are always handled with a `match` expression, preventing direct access to the underlying value and ensuring that all possible outcomes are considered.

```osprey
let result = someFunctionThatCanFail()

match result {
    Success { value } => print("Success: ${value}")
    Error { message } => print("Error: ${message}")
}
```

This approach guarantees that error handling is explicit, robust, and checked at compile time.

## 16. Examples

The `examples/` directory contains a variety of sample programs demonstrating Osprey's features. These examples are tested as part of the standard build process to ensure they remain up-to-date and functional.

## 17. Built-in Functions Reference

### 17.1 Basic I/O Functions

#### `print(value: int | string | bool) -> int`
Prints the given value to standard output with automatic type conversion.

**Parameters:**
- `value: int | string | bool` - The value to print (int, bool, string, or expression)

**Returns:** `int` - Exit code from puts function

**Examples:**
```osprey
print("Hello World")
print(42)
print(true)
print(x + y)
```

#### `input() -> int`
Reads an integer from stdin. Blocks until user enters a number.

**Parameters:** None

**Returns:** `int` - The number entered by user

**Examples:**
```osprey
let x = input()
let age = input()
```

#### `toString(value: int | string | bool) -> string`
Converts any value to its string representation.

#### `length(s: string) -> Result<int, StringError>`
🚨 **CRITICAL**: Returns the length of a string wrapped in a Result type for safety.

**MANDATORY PATTERN MATCHING:**
```osprey
match length("hello") {
    Success { value } => print("Length: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `contains(haystack: string, needle: string) -> Result<bool, StringError>`
🚨 **CRITICAL**: Checks if a string contains a substring, returns Result for safety.

**MANDATORY PATTERN MATCHING:**
```osprey
match contains("hello", "ell") {
    Success { value } => print("Found: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `substring(s: string, start: int, end: int) -> Result<string, StringError>`
🚨 **CRITICAL**: Extracts a substring from start to end, returns Result for bounds safety.

**MANDATORY PATTERN MATCHING:**
```osprey
match substring("hello", 1, 3) {
    Success { value } => print("Substring: ${value}")
    Error { message } => print("Error: ${message}")
}
```

**FUNDAMENTAL PRINCIPLE**: All string operations that could conceptually fail MUST return Result types. This enforces explicit error handling and prevents runtime panics.

### 17.2 File System Functions

#### `writeFile(path: string, content: string) -> Result<Success, string>`
Writes content to a file.

#### `readFile(path: string) -> Result<string, string>`
Reads file content as string.

#### `deleteFile(path: string) -> Result<Success, string>`
Deletes a file.

#### `createDirectory(path: string) -> Result<Success, string>`
Creates a directory.

#### `fileExists(path: string) -> bool`
Checks if file exists.

### 17.3 Process Operations

#### `spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessResult, string>`
Spawns external process with asynchronous stdout/stderr collection via callbacks.

**Key Features:**
- **Event-driven output collection** - Uses callbacks from the runtime for things like new stdout
- **Non-blocking execution** - Runs in background threads via fiber runtime
- **Real-time stdout/stderr** - Process output is sent to Osprey via callbacks as it's generated
- **Thread-safe** - Safe for concurrent use with multiple processes

**Architecture:**
The process runtime follows the same callback pattern as the HTTP server:
1. Osprey calls `spawnProcess(command)` 
2. C runtime spawns process and monitoring thread
3. C runtime calls back to Osprey with stdout/stderr events as they occur
4. Process completion triggers exit event callback

**ProcessResult Type:**
```osprey
type ProcessResult = {
    processId: int
}
```

**Event Types (handled by C runtime callbacks):**
- `PROCESS_STDOUT_DATA` - New stdout data available
- `PROCESS_STDERR_DATA` - New stderr data available  
- `PROCESS_EXIT` - Process completed with exit code

**Stdout Callback Handling:**

The process runtime uses callback functions to deliver stdout/stderr data as it's generated. The C runtime calls back into Osprey functions for real-time process events.

🚨 **CALLBACK IS MANDATORY!** The callback parameter is **REQUIRED** for `spawnProcess` - it cannot be omitted.

```osprey
// Define the callback function that C runtime will call for ALL process events
fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {
    match eventType {
        1 => print("[STDOUT] Process ${toString(processID)}: ${data}")
        2 => print("[STDERR] Process ${toString(processID)}: ${data}")
        3 => print("[EXIT] Process ${toString(processID)} exited with code: ${data}")
        _ => print("[UNKNOWN] Process ${toString(processID)} event ${toString(eventType)}: ${data}")
    }
}

// Process with callback-based stdout collection (CALLBACK IS MANDATORY!)
let result = spawnProcess("echo 'Hello from callback!'", processEventHandler)
match result {
    Success { value } => {
        print("Process spawned with ID: ${toString(value)}")
        let exitCode = awaitProcess(value)
        print("Process finished with exit code: ${toString(exitCode)}")
        cleanupProcess(value)
        print("Process cleaned up")
    }
    Error { message } => print("Failed to spawn process")
}
```

**Advanced Usage:**
```osprey
// Wait for process completion and get exit code
let result = spawnProcess("gcc myprogram.c -o myprogram")
match result {
    Success { value } => {
        let exitCode = awaitProcess(value.processId)
        cleanupProcess(value.processId)
        print("Compilation finished with exit code: ${toString(exitCode)}")
    }
    Error { message } => print("Compilation failed: ${message}")
}
```

#### `awaitProcess(processId: int) -> int`
Waits for process completion and returns exit code.

**Parameters:**
- `processId: int` - Process ID returned by spawnProcess

**Returns:** `int` - Process exit code (0 = success, non-zero = error)

#### `cleanupProcess(processId: int) -> void`
Cleans up process resources after completion.

**Parameters:**  
- `processId: int` - Process ID to clean up

**Note:** Always call this after `awaitProcess` to prevent memory leaks.

### 17.2 Functional Iterator Functions

#### `range(start: int, end: int) -> Iterator<int>`
Creates an iterator that generates integers from start (inclusive) to end (exclusive). Used with functional iterator functions like forEach, map, filter, and fold.

**Parameters:**
- `start: int` - Starting value (inclusive)
- `end: int` - Ending value (exclusive)

**Returns:** `Iterator<int>` - Iterator struct containing start and end values

**Examples:**
```osprey
range(1, 5)    // generates 1, 2, 3, 4
range(0, 3)    // generates 0, 1, 2
range(10, 13)  // generates 10, 11, 12
range(1, 10) |> forEach(print)
```

#### `forEach(iterator: Iterator<T>, function: T -> U) -> T`
Applies a function to each element in an iterator for side effects. This is the primary way to iterate through ranges and apply operations to each element.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to traverse (usually from range())
- `function: T -> U` - Function to apply to each element

**Returns:** `T` - Final counter value after iteration

**Examples:**
```osprey
range(1, 5) |> forEach(print)          // prints 1, 2, 3, 4
forEach(range(0, 3), double)           // calls double(0), double(1), double(2)
range(1, 10) |> forEach(square)
forEach(range(-2, 3), print)          // prints -2, -1, 0, 1, 2
```

#### `map(iterator: Iterator<T>, function: T -> U) -> U`
Transforms each element in an iterator by applying a function. Returns the result of the transformation function applied to each element.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to transform (usually from range())
- `function: T -> U` - Transformation function to apply

**Returns:** `U` - Result of applying function to each element

**Examples:**
```osprey
range(1, 5) |> map(double)    // applies double to 1, 2, 3, 4
map(range(0, 3), square)      // applies square to 0, 1, 2
range(1, 6) |> map(addFive)
```

#### `filter(iterator: Iterator<T>, predicate: T -> bool) -> T`
Selects elements from an iterator based on a predicate function. Only elements where the predicate returns true are processed.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to filter (usually from range())
- `predicate: T -> bool` - Function that returns true for elements to keep

**Returns:** `T` - Filtered results

**Examples:**
```osprey
range(1, 10) |> filter(isEven)
filter(range(0, 20), isPositive)
range(-5, 6) |> filter(isGreaterThanZero)
```

#### `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`
Reduces an iterator to a single value by repeatedly applying a function. Also known as reduce or accumulate in other languages.

**Parameters:**
- `iterator: Iterator<T>` - Iterator to reduce (usually from range())
- `initial: U` - Initial accumulator value
- `function: (U, T) -> U` - Function that combines accumulator with each element

**Returns:** `U` - Final accumulated value

**Examples:**
```osprey
range(1, 5) |> fold(0, add)          // sum: 0+1+2+3+4 = 10
fold(range(1, 6), 1, multiply)       // product: 1*1*2*3*4*5 = 120
range(0, 10) |> fold(0, max)
```

### 17.3 Pipe Operator

#### `|>` - Pipe Operator
The pipe operator takes the result of the left expression and passes it as the first argument to the right function. This enables elegant functional programming and method chaining.

**Syntax:** `expression |> function`

**Type:** `T |> (T -> U) -> U`

**Description:**
The pipe operator creates clean, readable function composition by allowing you to chain operations from left to right, making the data flow explicit and natural to read.

**Rules:**
- Left side can be any expression
- Right side must be a function or function call
- Creates clean, readable function composition
- Enables Haskell/Rust-style functional programming

**Examples:**
```osprey
// Basic piping
5 |> double |> print                 // Equivalent to: print(double(5))

// Iterator chaining
range(1, 10) |> forEach(print)
range(1, 5) |> map(square) |> fold(0, add)

// Complex chains
range(0, 20) |> filter(isEven) |> map(double) |> forEach(print)

// Multiple operations
let result = input() |> double |> square |> toString

// Nested operations
range(1, 10) 
  |> map(square) 
  |> filter(isEven) 
  |> fold(0, add) 
  |> print
```

### 17.4 Functional Programming Patterns

The combination of iterator functions and the pipe operator enables powerful functional programming patterns:

#### Chaining Pattern
```osprey
// Transform -> Filter -> Aggregate
range(1, 20)
  |> map(square)           // Square each number
  |> filter(isEven)        // Keep only even results
  |> fold(0, add)          // Sum them up
  |> print                 // Print the result
```

#### Side Effect Pattern
```osprey
// Process each element for side effects
range(1, 100)
  |> filter(isPrime)
  |> forEach(print)        // Print each prime
```

#### Data Transformation Pattern
```osprey
// Transform data through multiple stages
input()
  |> validateInput
  |> normalizeData
  |> processData
  |> formatOutput
  |> print
```

### 17.5 Fiber Types and Concurrency

Osprey provides lightweight concurrency through fiber types. Unlike traditional function-based approaches, fibers are proper type instances constructed using Osprey's standard type construction syntax.

#### Core Fiber Types

**`Fiber<T>`** - A lightweight concurrent computation that produces a value of type T
**`Channel<T>`** - A communication channel for passing values of type T between fibers

#### Fiber Construction

Fibers are created using standard type construction syntax:

```osprey
// Create a fiber that computes a value
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(n: 1000) 
}

// Create a fiber with more complex computation
let worker = Fiber<String> { 
    computation: fn() => {
        processData()
        "completed"
    }
}

// Create a parameterized fiber
let calculator = Fiber<Int> { 
    computation: fn() => multiply(x: 10, y: 20) 
}
```

#### Spawn Syntax Sugar

For convenience, Osprey provides `spawn` as syntax sugar for creating and immediately starting a fiber:

```osprey
// Using spawn (syntax sugar)
let result = spawn 42

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => 42 }
let result = fiber

// More complex spawn
let computation = spawn (x * 2 + y)

// Equivalent to:
let fiber = Fiber<Int> { computation: fn() => x * 2 + y }
let computation = fiber
```

The `spawn` keyword immediately evaluates the expression in a new fiber context, making it convenient for quick concurrent computations without the full type construction syntax.

#### Channel Construction

Channels are created using type construction syntax:

```osprey
// Unbuffered (synchronous) channel
let sync_channel = Channel<Int> { capacity: 0 }

// Buffered (asynchronous) channel  
let async_channel = Channel<String> { capacity: 10 }

// Large buffer channel
let buffer_channel = Channel<Int> { capacity: 100 }
```

#### Fiber Operations

Once created, fibers and channels are manipulated using functional operations:

**`await(fiber: Fiber<T>) -> T`** - Wait for fiber completion and get result
**`send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`** - Send value to channel
**`recv(channel: Channel<T>) -> Result<T, ChannelError>`** - Receive value from channel
**`yield() -> Unit`** - Voluntarily yield control to scheduler

```osprey
// Create and await a fiber
let task = Fiber<Int> { computation: fn() => heavyComputation() }
let result = await(task)

// Channel communication
let ch = Channel<String> { capacity: 5 }
send(ch, "hello")
let message = recv(ch)

// Yielding control
yield()
```

#### Complete Fiber Example

```osprey
// Producer fiber
let producer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        send(ch, 1)
        send(ch, 2) 
        send(ch, 3)
    }
}

// Consumer fiber
let consumer = Fiber<Unit> {
    computation: fn() => {
        let ch = Channel<Int> { capacity: 3 }
        let value1 = recv(ch)
        let value2 = recv(ch)
        let value3 = recv(ch)
        print("Received: ${value1}, ${value2}, ${value3}")
    }
}

// Start both fibers
await(producer)
await(consumer)
```

#### Select Expression for Channel Multiplexing

The `select` expression allows waiting on multiple channel operations:

```osprey
let ch1 = Channel<String> { capacity: 1 }
let ch2 = Channel<Int> { capacity: 1 }

let result = select {
    msg => recv(ch1) => process_string(msg)
    num => recv(ch2) => process_number(num)
    _ => timeout_handler()
}
```

#### Rust Interoperability

Osprey fibers are designed to interoperate with Rust's async/await system:

```osprey
// Osprey fiber that calls Rust async function
extern fn rust_async_task() -> Future<Int>

let osprey_task = Fiber<Int> {
    computation: fn() => await(rust_async_task())
}

let result = await(osprey_task)
```

#### Fiber-Isolated Modules

Each fiber gets its own isolated instance of modules, preventing data races:

```osprey
module Counter {
    let mut count = 0
    fn increment() -> Int = { count = count + 1; count }
    fn get() -> Int = count
}

// Each fiber has its own Counter instance
let fiber1 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}
let fiber2 = Fiber<Int> { 
    computation: fn() => Counter.increment() 
}

// These will both return 1, not 1 and 2
let result1 = await(fiber1)  // 1
let result2 = await(fiber2)  // 1 (separate instance)
```

---

## 18. Security and Sandboxing

The Osprey compiler includes built-in security controls to restrict access to potentially dangerous functionality like network operations and file system access. This is essential for safe code execution in environments like web compilers where untrusted code may be executed.

### 18.1 Security Flags

#### `--sandbox`
Enables sandbox mode, which disables all potentially risky operations:
- HTTP/HTTPS operations (httpCreateServer, httpGet, httpPost, etc.)
- WebSocket operations (websocketConnect, websocketSend, etc.)
- File system access (when implemented)
- Foreign Function Interface (FFI)
- Process execution

**Example:**
```bash
osprey program.osp --sandbox --llvm
```

#### Granular Security Flags

For more granular control, you can disable specific categories of operations:

- `--no-http`: Disable HTTP client and server functions
- `--no-websocket`: Disable WebSocket client and server functions  
- `--no-fs`: Disable file system read/write operations
- `--no-ffi`: Disable foreign function interface

**Examples:**
```bash
# Disable only HTTP operations
osprey program.osp --no-http --compile

# Disable HTTP and WebSocket operations
osprey program.osp --no-http --no-websocket --run

# Disable file system access only
osprey program.osp --no-fs --llvm
```

### 18.2 Security Policies

#### Default Security (Permissive)
By default, all operations are allowed for backward compatibility and normal development use.

#### Sandbox Security (Restrictive)
When `--sandbox` is used, all potentially dangerous functions are unavailable. This is recommended for:
- Web-based code execution
- Untrusted code evaluation
- Educational environments
- Code review systems

### 18.3 Blocked Functions by Category

#### HTTP Functions
When HTTP access is disabled (`--no-http` or `--sandbox`), these functions are unavailable:
- `httpCreateServer` - Create HTTP server
- `httpListen` - Start HTTP server listening
- `httpStopServer` - Stop HTTP server
- `httpCreateClient` - Create HTTP client
- `httpGet` - HTTP GET request
- `httpPost` - HTTP POST request
- `httpPut` - HTTP PUT request
- `httpDelete` - HTTP DELETE request
- `httpRequest` - Generic HTTP request
- `httpCloseClient` - Close HTTP client

#### WebSocket Functions
When WebSocket access is disabled (`--no-websocket` or `--sandbox`), these functions are unavailable:
- `websocketConnect` - Connect to WebSocket server
- `websocketSend` - Send WebSocket message
- `websocketClose` - Close WebSocket connection
- `websocketCreateServer` - Create WebSocket server
- `websocketServerListen` - Start WebSocket server
- `websocketServerSend` - Send message to specific client
- `websocketServerBroadcast` - Broadcast message to all clients
- `websocketStopServer` - Stop WebSocket server

#### File System Functions (Future)
When file system access is disabled (`--no-fs` or `--sandbox`), these functions will be unavailable:
- `readFile` - Read file contents
- `writeFile` - Write file contents
- `deleteFile` - Delete file
- `createDirectory` - Create directory
- `listDirectory` - List directory contents

### 18.4 Function Availability

In different security modes, certain functions are simply not available in the language:

**Sandbox Mode**: Only safe functions like `print`, `toString`, `range`, etc. are available. Dangerous functions like `httpCreateServer` or `websocketConnect` result in "undefined function" compile errors.

**Partial Restrictions**: When specific categories are disabled (e.g., `--no-http`), those functions are unavailable while others remain accessible.

**Default Mode**: All functions are available.
- A human-readable explanation

### 18.5 Programming Best Practices

#### For Safe Code
Write code that doesn't use security-sensitive functions:
```osprey
// Safe operations - work in all security modes
let x = 42
let y = 24
let sum = x + y
print("Sum: ")
print(sum)
```

#### For Network Code
When writing network code, be aware that it may be restricted:
```osprey
// This will fail in sandbox mode or with --no-http
let serverID = httpCreateServer(port: 8080, address: "127.0.0.1")
```

### 18.6 Implementation Details

#### Security Configuration
Security settings are configured at compilation time and cannot be bypassed by the compiled program. The security checks happen during the LLVM IR generation phase, preventing security-sensitive functions from being included in the generated code.

#### Performance Impact
Security checks add minimal overhead during compilation and no runtime overhead, as restricted functions are simply not compiled into the final program.

#### Backward Compatibility
All existing code continues to work with default settings. Security restrictions are opt-in and don't affect normal development workflows.

#### Integration with Web Compiler
The security features are designed specifically for web compiler integration:

```javascript
// Example web compiler usage
const result = await compileOsprey(sourceCode, {
    mode: 'sandbox',  // Enable sandbox mode
    outputFormat: 'llvm'
});
```

#### Security Summary
When using security restrictions, the compiler will display a security summary:

```bash
# Sandbox mode
Security: SANDBOX MODE - All risky operations disabled

# Partial restrictions
Security: Allowed=[FileRead,FileWrite,FFI] Blocked=[HTTP,WebSocket]
```

---

## Summary

Osprey is a modern functional programming language with:

- **Type Safety**: No runtime panics, all errors handled explicitly via Result types
- **Named Arguments**: Multi-parameter functions require named arguments for clarity
- **Functional Programming**: Powerful iterator functions with pipe operator
- **Lightweight Fibers**: Zero-cost concurrency with Rust-like async/await
- **Fiber-Isolated Modules**: No global state, each fiber gets its own module instances
- **Rust Interoperability**: Seamless integration with Rust libraries
- **Memory Safety**: No shared mutable state between fibers

**Key Innovation**: The fiber-isolated module system eliminates data races by design while maintaining clean encapsulation through accessor patterns.

---

**End of Specification**

This specification defines the complete syntax and semantics of the Osprey programming language, including its revolutionary fiber-isolated module system and lightweight concurrency features. The accompanying `osprey.g4` grammar file provides the formal ANTLR4 grammar definition for parsing.

## 7. Server Applications and Long-Running Processes

### 7.1 Functional Approaches to Server Persistence

**Osprey is a functional language and does NOT support imperative loop constructs.** Server applications that need to stay alive should use functional patterns instead:

#### 7.1.1 Fiber-Based Server Persistence

Use fibers to handle concurrent requests and keep the server process alive:

```osprey
// HTTP server with fiber-based request handling
fn handleRequest(requestId: Int) -> Int = {
    // Process the request
    let response = processData(requestId)
    response
}

fn serverMain() -> Unit = {
    let server = httpCreateServer(port: 8080, address: "0.0.0.0")
    
    // Spawn fibers to handle requests concurrently
    let requestHandler = spawn {
        // Use functional iteration to process incoming requests
        range(1, 1000000) |> forEach(handleRequest)
    }
    
    // Keep server alive by awaiting the handler fiber
    await(requestHandler)
}
```

#### 7.1.2 Recursive Function Patterns

Use tail-recursive functions for continuous processing:

```osprey
fn serverLoop(state: ServerState) -> Unit = match getNextRequest(state) {
    Some { request } => {
        let newState = processRequest(request, state)
        serverLoop(newState)  // Tail recursion keeps server alive
    }
    None => serverLoop(state)  // Continue waiting for requests
}

fn main() -> Unit = {
    let initialState = initializeServer()
    serverLoop(initialState)  // Functional "loop" via recursion
}
```

#### 7.1.3 Event-Driven Architecture with Channels

Use channels for event-driven server architectures:

```osprey
fn eventProcessor(eventChannel: Channel<Event>) -> Unit = {
    let event = recv(eventChannel)
    match event {
        Success { value } => {
            processEvent(value)
            eventProcessor(eventChannel)  // Continue processing
        }
        Error { _ } => eventProcessor(eventChannel)  // Retry on error
    }
}

fn serverWithEvents() -> Unit = {
    let eventChan = Channel<Event> { capacity: 100 }
    
    // Spawn event processor fiber
    let processor = spawn eventProcessor(eventChan)
    
    // Spawn request handlers that send events
    let handler1 = spawn handleHTTPRequests(eventChan)
    let handler2 = spawn handleWebSocketRequests(eventChan)
    
    // Wait for all handlers
    await(processor)
    await(handler1)
    await(handler2)
}
```

#### 7.1.4 Functional Iterator-Based Processing

Use functional iterators for continuous data processing:

```osprey
// Stream processing with functional iterators
fn processIncomingData() -> Unit = {
    // Process data in batches using functional approach
    range(1, Int.MAX_VALUE) 
    |> map(getBatch)
    |> filter(isValidBatch)
    |> forEach(processBatch)
}

fn webSocketServer() -> Unit = {
    let server = websocketCreateServer(port: 8080, address: "0.0.0.0", path: "/ws")
    
    // Use functional processing instead of loops
    let dataProcessor = spawn processIncomingData()
    let connectionHandler = spawn manageConnections(server)
    
    await(dataProcessor)
    await(connectionHandler)
}
```

### 7.2 Why No Imperative Loops?

**Functional Superiority:**
1. **Composability** - Functional iterators can be chained with `|>`
2. **Safety** - No mutable state, no infinite loop bugs
3. **Concurrency** - Fibers provide better parallelism than loops
4. **Testability** - Pure functions are easier to test than stateful loops

**Anti-Pattern:**
```osprey
// ❌ WRONG - Imperative loops (NOT SUPPORTED)
loop {
    let request = getRequest()
    processRequest(request)
}
```

**Functional Pattern:**
```osprey
// ✅ CORRECT - Functional approach
fn serverHandler() -> Unit = {
    requestStream() 
    |> map(processRequest)
    |> forEach(sendResponse)
}
```

### 7.3 Performance Considerations

Functional approaches in Osprey are optimized for:
- **Tail call optimization** prevents stack overflow in recursive functions
- **Fiber scheduling** provides efficient concurrency without OS threads
- **Channel buffering** enables high-throughput event processing
- **Iterator fusion** optimizes chained functional operations

This functional approach provides better maintainability, testability, and performance than traditional imperative loops.