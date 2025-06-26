13. [Built-in Functions](0014-Built-InFunctions.md)
    - [HTTP Core Types](#141-http-core-types)
        - [HTTP Method Union Type](#http-method-union-type)
        - [HTTP Request Type (Immutable)](#http-request-type-immutable)
        - [HTTP Response Type (Immutable with Streaming)](#http-response-type-immutable-with-streaming)
    - [HTTP Server Functions](#142-http-server-functions)
        - [`httpCreateServer(port: Int, address: String) -> Result<ServerID, String>`](#httpcreateserverport-int-address-string---resultserverid-string)
        - [`httpListen(serverID: Int, handler: fn(String, String, String, String) -> String) -> Result<Success, String>`](#httplistenserverid-int-handler-fnstring-string-string-string---string---resultsuccess-string)
        - [`httpStopServer(serverID: Int) -> Result<Success, String>`](#httpstopserverserverid-int---resultsuccess-string)
        - [HTTP Request Handling Bridge](#1421-http-request-handling-bridge)
            - [Request Handling Architecture](#request-handling-architecture)
            - [Bridge Function Specification](#bridge-function-specification)
            - [Raw Function Pointer Callbacks](#raw-function-pointer-callbacks)
            - [Legacy Bridge Function (Deprecated)](#legacy-bridge-function-deprecated)
            - [New Raw Callback Architecture Flow](#new-raw-callback-architecture-flow)
            - [Implementation Requirements](#implementation-requirements)
            - [Example Implementation](#example-implementation)
    - [HTTP Client Functions](#143-http-client-functions)
        - [`httpCreateClient(baseUrl: String, timeout: Int) -> Result<ClientID, String>`](#httpcreeateclientbaseurl-string-timeout-int---resultclientid-string)
        - [`httpGet(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`](#httpgetclientid-int-path-string-headers-string---resultstatuscode-string)
        - [`httpPost(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`](#httppostclientid-int-path-string-body-string-headers-string---resultstatuscode-string)
        - [`httpPut(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`](#httpputclientid-int-path-string-body-string-headers-string---resultstatuscode-string)
        - [`httpDelete(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`](#httpdeleteclientid-int-path-string-headers-string---resultstatuscode-string)
        - [`httpRequest(clientID: Int, method: HttpMethod, path: String, headers: String, body: String) -> Result<StatusCode, String>`](#httprequestclientid-int-method-httpmethod-path-string-headers-string-body-string---resultstatuscode-string)
        - [`httpCloseClient(clientID: Int) -> Result<Success, String>`](#httpcloseclientclientid-int---resultsuccess-string)
    - [WebSocket Support (Two-Way Communication)](#144-websocket-support-two-way-communication)
        - [WebSocket Security Implementation](#1441-websocket-security-implementation)
        - [Security Standards Compliance](#1442-security-standards-compliance)
        - [Security Architecture](#1443-security-architecture)
        - [Security Testing and Validation](#1444-security-testing-and-validation)
        - [Security References and Standards](#1445-security-references-and-standards)
        - [`websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`](#websocketconnecturl-string-messagehandler-fnstring---resultsuccess-string---resultwebsocketid-string)
        - [`websocketSend(wsID: Int, message: String) -> Result<Success, String>`](#websocketsendwsid-int-message-string---resultsuccess-string)
        - [`websocketClose(wsID: Int) -> Result<Success, String>`](#websocketclosewsid-int---resultsuccess-string)
        - [WebSocket Server Functions](#1441-websocket-server-functions)
            - [`websocketCreateServer(port: Int, address: String, path: String) -> Int`](#websocketcreateserverport-int-address-string-path-string---int)
            - [`websocketServerListen(serverID: Int) -> Int`](#websocketserverlistenserverid-int---int)
            - [`websocketServerBroadcast(serverID: Int, message: String) -> Int`](#websocketserverbroadcastserverid-int-message-string---int)
            - [`websocketStopServer(serverID: Int) -> Int`](#websocketstopserverserverid-int---int)
            - [`websocketKeepAlive() -> Void`](#websocketkeepalive---void)
    - [Streaming Response Bodies](#145-streaming-response-bodies)
        - [Complete Response](#complete-response)
        - [Streamed Response](#streamed-response)
    - [Error Handling in HTTP](#146-error-handling-in-http)
    - [Fiber-Based Concurrency](#147-fiber-based-concurrency)
    - [Complete HTTP Server Example](#148-complete-http-server-example)


# 13. Built-in Functions

🚀 **IMPLEMENTATION STATUS**: HTTP and basic I/O functions are implemented and working. WebSocket functions are implemented but undergoing testing. Fiber operations are partially implemented.

Osprey provides built-in functions for I/O, networking, concurrency, and functional programming. All functions follow Osprey's functional programming paradigms with Result types for error handling.

## 13.1 Basic I/O Functions

### `print(value: int | string | bool) -> int`
Prints values to standard output with automatic type conversion.

```osprey
print("Hello World")
print(42)
print(true)
```

### `input() -> int`
Reads an integer from stdin.

```osprey
let x = input()
```

### `toString(value: int | string | bool) -> string`
Converts any value to its string representation.

### String Functions

#### `length(s: string) -> Result<int, StringError>`
Returns string length. Requires pattern matching for safety.

```osprey
match length("hello") {
    Success { value } => print("Length: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `contains(haystack: string, needle: string) -> Result<bool, StringError>`
Checks if a string contains a substring.

```osprey
match contains("hello", "ell") {
    Success { value } => print("Found: ${value}")
    Error { message } => print("Error: ${message}")
}
```

#### `substring(s: string, start: int, end: int) -> Result<string, StringError>`
Extracts a substring from start to end.

## 13.2 File System Functions

### `writeFile(path: string, content: string) -> Result<Success, string>`
Writes content to a file.

### `readFile(path: string) -> Result<string, string>`
Reads file content as string.

### `deleteFile(path: string) -> Result<Success, string>`
Deletes a file.

### `createDirectory(path: string) -> Result<Success, string>`
Creates a directory.

### `fileExists(path: string) -> bool`
Checks if file exists.

## 13.3 Process Operations

### `spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessResult, string>`
Spawns external process with asynchronous stdout/stderr collection via callbacks.

```osprey
fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {
    match eventType {
        1 => print("[STDOUT] ${data}")
        2 => print("[STDERR] ${data}")
        3 => print("[EXIT] Code: ${data}")
        _ => print("[UNKNOWN] ${data}")
    }
}

let result = spawnProcess("echo 'Hello'", processEventHandler)
```

### `awaitProcess(processId: int) -> int`
Waits for process completion and returns exit code.

### `cleanupProcess(processId: int) -> void`
Cleans up process resources after completion.

## 13.4 Functional Programming

### Iterator Functions

#### `range(start: int, end: int) -> Iterator<int>`
Creates an iterator from start (inclusive) to end (exclusive).

```osprey
range(1, 5)    // generates 1, 2, 3, 4
```

#### `forEach(iterator: Iterator<T>, function: T -> U) -> T`
Applies a function to each element for side effects.

```osprey
range(1, 5) |> forEach(print)          // prints 1, 2, 3, 4
```

#### `map(iterator: Iterator<T>, function: T -> U) -> U`
Transforms each element by applying a function.

```osprey
range(1, 5) |> map(double)    // applies double to 1, 2, 3, 4
```

#### `filter(iterator: Iterator<T>, predicate: T -> bool) -> T`
Selects elements based on a predicate function.

```osprey
range(1, 10) |> filter(isEven)
```

#### `fold(iterator: Iterator<T>, initial: U, function: (U, T) -> U) -> U`
Reduces an iterator to a single value.

```osprey
range(1, 5) |> fold(0, add)          // sum: 0+1+2+3+4 = 10
```

### Pipe Operator `|>`

The pipe operator passes the left expression as the first argument to the right function.

```osprey
5 |> double |> print                 // Equivalent to: print(double(5))
range(1, 10) |> map(square) |> filter(isEven) |> forEach(print)
```

## 13.5 HTTP Functions

### Core Types

```osprey
type HttpMethod = GET | POST | PUT | DELETE | PATCH | HEAD | OPTIONS

type HttpRequest = {
    method: HttpMethod,
    path: String,
    headers: String,
    body: String,
    queryParams: String
}

type HttpResponse = {
    status: Int,
    headers: String,
    contentType: String,
    body: String
}
```

### HTTP Server Functions

#### `httpCreateServer(port: Int, address: String) -> Result<ServerID, String>`
Creates an HTTP server bound to the specified port and address.

```osprey
let serverResult = httpCreateServer(port: 8080, address: "127.0.0.1")
```

#### `httpListen(serverID: Int, handler: fn(String, String, String, String) -> String) -> Result<Success, String>`
Starts the HTTP server with a raw request handler. The handler receives method, path, headers, and body and returns the response body.

```osprey
fn handleRequest(method: String, path: String, headers: String, body: String) -> String = 
    match method {
        "GET" => match path {
            "/health" => "{\"status\": \"healthy\"}"
            "/api/users" => "[{\"id\": 1, \"name\": \"Alice\"}]"
            _ => "Not Found"
        }
        "POST" => "{\"message\": \"Created\"}"
        _ => "Method not allowed"
    }

let listenResult = httpListen(serverId, handleRequest)
```

#### `httpStopServer(serverID: Int) -> Result<Success, String>`
Stops the HTTP server and cleans up resources.

### HTTP Client Functions

#### `httpCreateClient(baseUrl: String, timeout: Int) -> Result<ClientID, String>`
Creates an HTTP client for making requests.

#### `httpGet(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`
Makes an HTTP GET request.

#### `httpPost(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`
Makes an HTTP POST request with a request body.

#### `httpPut(clientID: Int, path: String, body: String, headers: String) -> Result<StatusCode, String>`
Makes an HTTP PUT request.

#### `httpDelete(clientID: Int, path: String, headers: String) -> Result<StatusCode, String>`
Makes an HTTP DELETE request.

#### `httpRequest(clientID: Int, method: HttpMethod, path: String, headers: String, body: String) -> Result<StatusCode, String>`
Generic HTTP request function for any HTTP method.

#### `httpCloseClient(clientID: Int) -> Result<Success, String>`
Closes the HTTP client and cleans up resources.

## 13.6 WebSocket Functions

WebSocket functions provide real-time, bidirectional communication. Implementation includes security features and follows RFC 6455 standards.

### Client Functions

#### `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`
Establishes a WebSocket connection.

```osprey
fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}

let wsResult = websocketConnect("ws://localhost:8080/chat", handleMessage)
```

#### `websocketSend(wsID: Int, message: String) -> Result<Success, String>`
Sends a message through the WebSocket connection.

#### `websocketClose(wsID: Int) -> Result<Success, String>`
Closes the WebSocket connection.

### Server Functions

#### `websocketCreateServer(port: Int, address: String, path: String) -> Int`
Creates a WebSocket server.

🚧 **IMPLEMENTATION STATUS**: Currently has port binding issues in test environments.

#### `websocketServerListen(serverID: Int) -> Int`
Starts the WebSocket server listening for connections.

#### `websocketServerBroadcast(serverID: Int, message: String) -> Int`
Broadcasts a message to all connected clients.

#### `websocketStopServer(serverID: Int) -> Int`
Stops the WebSocket server.

#### `websocketKeepAlive() -> Void`
Maintains WebSocket connections.

## 13.7 Fiber and Concurrency Functions

Osprey provides lightweight concurrency through fibers.

### Fiber Types

```osprey
// Create a fiber
let task = Fiber<Int> { 
    computation: fn() => calculatePrimes(1000) 
}

// Spawn syntax sugar
let result = spawn 42

// Channels for communication
let ch = Channel<String> { capacity: 10 }
```

### Fiber Operations

#### `await(fiber: Fiber<T>) -> T`
Wait for fiber completion and get result.

#### `send(channel: Channel<T>, value: T) -> Result<Unit, ChannelError>`
Send value to channel.

#### `recv(channel: Channel<T>) -> Result<T, ChannelError>`
Receive value from channel.

#### `yield() -> Unit`
Voluntarily yield control to scheduler.

### Example Usage

```osprey
// Producer-consumer pattern
let ch = Channel<Int> { capacity: 3 }

let producer = spawn {
    send(ch, 1)
    send(ch, 2)
    send(ch, 3)
}

let consumer = spawn {
    let value1 = recv(ch)
    let value2 = recv(ch)
    let value3 = recv(ch)
    print("Received values")
}

await(producer)
await(consumer)
```

## 13.8 Complete Example

A simple HTTP server with functional programming:

```osprey
fn main() -> Int = {
    let serverResult = httpCreateServer(8080, "127.0.0.1")
    match serverResult {
        Success serverId => {
            print("Server created, starting listener...")
            let listenResult = httpListen(serverId, handleRequest)
            match listenResult {
                Success _ => print("Server stopped")
                Err message => print("Server error: ${message}")
            }
        }
        Err message => print("Failed to create server: ${message}")
    }
    0
}

fn handleRequest(method: String, path: String, headers: String, body: String) -> String = 
    match method {
        "GET" => match path {
            "/health" => "{\"status\": \"healthy\"}"
            "/users" => generateUserList()
            _ => "Not Found"
        }
        "POST" => match path {
            "/users" => createUser(body)
            _ => "Endpoint not found"
        }
        _ => "Method not allowed"
    }

fn generateUserList() -> String = 
    range(1, 6)
    |> map(createUser)
    |> fold("[]", appendJson)

fn createUser(id: Int) -> String = 
    "{\"id\": ${toString(id)}, \"name\": \"User${toString(id)}\"}"