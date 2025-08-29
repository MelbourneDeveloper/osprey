---
layout: page
title: "Built-in Functions"
description: "Complete reference for all built-in functions in Osprey"
---

All built-in functions available in Osprey.

## [awaitProcess](awaitprocess/)

**Signature:** `awaitProcess(handle: ProcessHandle) -> int`

Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes.

## [cleanupProcess](cleanupprocess/)

**Signature:** `cleanupProcess(handle: ProcessHandle) -> void`

Cleans up resources associated with a completed process. Should be called after awaitProcess.

## [contains](contains/)

**Signature:** `contains(haystack: string, needle: string) -> bool`

Checks if a string contains a substring.

## [filter](filter/)

**Signature:** `filter(iterator: iterator, predicate: function) -> iterator`

Filters elements in an iterator based on a predicate function.

## [fold](fold/)

**Signature:** `fold(iterator: iterator, initial: any, fn: function) -> any`

Reduces an iterator to a single value using an accumulator function.

## [forEach](foreach/)

**Signature:** `forEach(iterator: iterator, fn: function) -> int`

Applies a function to each element in an iterator.

## [httpCloseClient](httpcloseclient/)

**Signature:** `httpCloseClient(clientID: int) -> int`

Closes the HTTP client and cleans up resources.

## [httpCreateClient](httpcreateclient/)

**Signature:** `httpCreateClient(baseUrl: string, timeout: int) -> int`

Creates an HTTP client for making requests to a base URL.

## [httpCreateServer](httpcreateserver/)

**Signature:** `httpCreateServer(port: int, address: string) -> int`

Creates an HTTP server bound to the specified port and address.

## [httpDelete](httpdelete/)

**Signature:** `httpDelete(clientID: int, path: string, headers: string) -> int`

Makes an HTTP DELETE request to the specified path.

## [httpGet](httpget/)

**Signature:** `httpGet(clientID: int, path: string, headers: string) -> int`

Makes an HTTP GET request to the specified path.

## [httpListen](httplisten/)

**Signature:** `httpListen(serverID: int, handler: function) -> int`

Starts the HTTP server listening for requests with a handler function.

## [httpPost](httppost/)

**Signature:** `httpPost(clientID: int, path: string, body: string, headers: string) -> int`

Makes an HTTP POST request with a request body.

## [httpPut](httpput/)

**Signature:** `httpPut(clientID: int, path: string, body: string, headers: string) -> int`

Makes an HTTP PUT request with a request body.

## [httpRequest](httprequest/)

**Signature:** `httpRequest(clientID: int, method: int, path: string, headers: string, body: string) -> int`

Makes a generic HTTP request with any method.

## [httpStopServer](httpstopserver/)

**Signature:** `httpStopServer(serverID: int) -> int`

Stops the HTTP server and closes all connections.

## [input](input/)

**Signature:** `input() -> int`

Reads an integer from the user's input.

## [length](length/)

**Signature:** `length(s: string) -> int`

Returns the length of a string.

## [map](map/)

**Signature:** `map(iterator: iterator, fn: function) -> iterator`

Transforms each element in an iterator using a function, returning a new iterator.

## [print](print/)

**Signature:** `print(value: any) -> int`

Prints a value to the console. Automatically converts the value to a string representation.

## [range](range/)

**Signature:** `range(start: int, end: int) -> iterator`

Creates an iterator that generates numbers from start to end (exclusive).

## [readFile](readfile/)

**Signature:** `readFile(filename: string) -> Result<string, string>`

Reads the entire contents of a file as a string.

## [sleep](sleep/)

**Signature:** `sleep(milliseconds: int) -> int`

Pauses execution for the specified number of milliseconds.

## [spawnProcess](spawnprocess/)

**Signature:** `spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessHandle, string>`

Spawns an external async process with MANDATORY callback for stdout/stderr capture. The callback function receives (processID: int, eventType: int, data: string) and is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!

## [substring](substring/)

**Signature:** `substring(s: string, start: int, end: int) -> string`

Extracts a substring from start to end index.

## [toString](tostring/)

**Signature:** `toString(value: any) -> string`

Converts a value to its string representation.

## [websocketKeepAlive](websocketKeepAlive/)

**Signature:** `websocketKeepAlive() -> Unit`

⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. Keeps the WebSocket server running indefinitely until interrupted (blocking operation).

## [websocketClose](websocketclose/)

**Signature:** `websocketClose(wsID: Int) -> Result<Success, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Closes the WebSocket connection and cleans up resources.

## [websocketConnect](websocketconnect/)

**Signature:** `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<WebSocketID, String> and takes string handler instead of function pointer. Establishes a WebSocket connection with a message handler callback.

## [websocketCreateServer](websocketcreateserver/)

**Signature:** `websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<ServerID, String> and has critical runtime issues with port binding failures. Creates a WebSocket server bound to the specified port, address, and path.

## [websocketSend](websocketsend/)

**Signature:** `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Sends a message through the WebSocket connection.

## [websocketServerBroadcast](websocketserverbroadcast/)

**Signature:** `websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t (number of clients sent to) instead of Result<Success, String>. Broadcasts a message to all connected WebSocket clients.

## [websocketServerListen](websocketserverlisten/)

**Signature:** `websocketServerListen(serverID: Int) -> Result<Success, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String> and currently returns -4 (bind failed) due to port binding issues. Starts the WebSocket server listening for connections.

## [websocketStopServer](websocketstopserver/)

**Signature:** `websocketStopServer(serverID: Int) -> Result<Success, String>`

⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Stops the WebSocket server and closes all connections.

## [writeFile](writefile/)

**Signature:** `writeFile(filename: string, content: string) -> Result<Success, string>`

Writes content to a file. Creates the file if it doesn't exist.

