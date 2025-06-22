---
layout: page
title: "Built-in Functions"
description: "Complete reference for all built-in functions in Osprey"
---

All built-in functions available in Osprey.

## [contains](contains/)

**Signature:** `contains(haystack: string, needle: string) -> bool`

Checks if a string contains a substring.

## [extractCode](extractcode/)

**Signature:** `extractCode(json: string) -> Result<string, string>`

Extracts code from a JSON structure.

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

## [parseJSON](parsejson/)

**Signature:** `parseJSON(json: string) -> Result<string, string>`

Parses a JSON string and returns the parsed result.

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

**Signature:** `spawnProcess(command: string) -> Result<ProcessResult, string>`

Spawns an external process and returns the result. Currently supports simple command execution.

## [substring](substring/)

**Signature:** `substring(s: string, start: int, end: int) -> string`

Extracts a substring from start to end index.

## [toString](tostring/)

**Signature:** `toString(value: any) -> string`

Converts a value to its string representation.

## [webSocketKeepAlive](websocketkeepalive/)

**Signature:** `webSocketKeepAlive() -> int`

Keeps the WebSocket server running indefinitely until interrupted.

## [websocketClose](websocketclose/)

**Signature:** `websocketClose(wsID: int) -> int`

Closes the WebSocket connection.

## [websocketConnect](websocketconnect/)

**Signature:** `websocketConnect(url: string, messageHandler: string) -> int`

Establishes a WebSocket connection to the specified URL.

## [websocketCreateServer](websocketcreateserver/)

**Signature:** `websocketCreateServer(port: int, address: string, path: string) -> int`

Creates a WebSocket server bound to the specified port, address, and path.

## [websocketSend](websocketsend/)

**Signature:** `websocketSend(wsID: int, message: string) -> int`

Sends a message through the WebSocket connection.

## [websocketServerBroadcast](websocketserverbroadcast/)

**Signature:** `websocketServerBroadcast(serverID: int, message: string) -> int`

Broadcasts a message to all connected WebSocket clients.

## [websocketServerListen](websocketserverlisten/)

**Signature:** `websocketServerListen(serverID: int) -> int`

Starts the WebSocket server listening for connections.

## [websocketStopServer](websocketstopserver/)

**Signature:** `websocketStopServer(serverID: int) -> int`

Stops the WebSocket server and closes all connections.

## [writeFile](writefile/)

**Signature:** `writeFile(filename: string, content: string) -> Result<Success, string>`

Writes content to a file. Creates the file if it doesn't exist.

