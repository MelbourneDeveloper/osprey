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