**Version:** 0.2.0-alpha

**Date:** June 2025

**Author** Christian Findlay

<div class="table-of-contents">

- [18. Security and Sandboxing](#18-security-and-sandboxing)
  - [18.1 Security Flags](#181-security-flags)
    - [`--sandbox`](#--sandbox)
    - [Granular Security Flags](#granular-security-flags)
  - [18.2 Security Policies](#182-security-policies)
    - [Default Security (Permissive)](#default-security-permissive)
    - [Sandbox Security (Restrictive)](#sandbox-security-restrictive)
  - [18.3 Blocked Functions by Category](#183-blocked-functions-by-category)
    - [HTTP Functions](#http-functions)
    - [WebSocket Functions](#websocket-functions)
    - [File System Functions (Future)](#file-system-functions-future)
  - [18.4 Function Availability](#184-function-availability)
  - [18.5 Programming Best Practices](#185-programming-best-practices)
    - [For Safe Code](#for-safe-code)
    - [For Network Code](#for-network-code)
  - [18.6 Implementation Details](#186-implementation-details)
    - [Security Configuration](#security-configuration)
    - [Performance Impact](#performance-impact)
    - [Backward Compatibility](#backward-compatibility)
    - [Integration with Web Compiler](#integration-with-web-compiler)
    - [Security Summary](#security-summary)
- [Summary](#summary)
- [7. Server Applications and Long-Running Processes](#7-server-applications-and-long-running-processes)
  - [7.1 Functional Approaches to Server Persistence](#71-functional-approaches-to-server-persistence)
    - [7.1.1 Fiber-Based Server Persistence](#711-fiber-based-server-persistence)
    - [7.1.2 Recursive Function Patterns](#712-recursive-function-patterns)
    - [7.1.3 Event-Driven Architecture with Channels](#713-event-driven-architecture-with-channels)
    - [7.1.4 Functional Iterator-Based Processing](#714-functional-iterator-based-processing)
  - [7.2 Why No Imperative Loops?](#72-why-no-imperative-loops)
  - [7.3 Performance Considerations](#73-performance-considerations)

</div>


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