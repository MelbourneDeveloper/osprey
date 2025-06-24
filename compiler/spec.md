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