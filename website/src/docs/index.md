---
layout: page
title: "API Reference - Osprey Programming Language"
description: "Complete reference documentation for the Osprey programming language"
---

## Quick Navigation

- [Functions](functions/) - Built-in functions for I/O, iteration, and data transformation
- [Types](types/) - Built-in data types (Int, String, Bool, Any)
- [Operators](operators/) - Arithmetic, comparison, and logical operators
- [Keywords](keywords/) - Language keywords (fn, let, type, match, import)

## Function Reference

| Function | Description |
|----------|-------------|
| [awaitProcess](functions/awaitprocess/) | Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes. |
| [cleanupProcess](functions/cleanupprocess/) | Cleans up resources associated with a completed process. Should be called after awaitProcess. |
| [contains](functions/contains/) | Checks if a string contains a substring. |
| [filter](functions/filter/) | Filters elements in an iterator based on a predicate function. |
| [fold](functions/fold/) | Reduces an iterator to a single value using an accumulator function. |
| [forEach](functions/foreach/) | Applies a function to each element in an iterator. |
| [httpCloseClient](functions/httpcloseclient/) | Closes the HTTP client and cleans up resources. |
| [httpCreateClient](functions/httpcreateclient/) | Creates an HTTP client for making requests to a base URL. |
| [httpCreateServer](functions/httpcreateserver/) | Creates an HTTP server bound to the specified port and address. |
| [httpDelete](functions/httpdelete/) | Makes an HTTP DELETE request to the specified path. |
| [httpGet](functions/httpget/) | Makes an HTTP GET request to the specified path. |
| [httpListen](functions/httplisten/) | Starts the HTTP server listening for requests with a handler function. |
| [httpPost](functions/httppost/) | Makes an HTTP POST request with a request body. |
| [httpPut](functions/httpput/) | Makes an HTTP PUT request with a request body. |
| [httpRequest](functions/httprequest/) | Makes a generic HTTP request with any method. |
| [httpStopServer](functions/httpstopserver/) | Stops the HTTP server and closes all connections. |
| [input](functions/input/) | Reads an integer from the user's input. |
| [length](functions/length/) | Returns the length of a string. |
| [map](functions/map/) | Transforms each element in an iterator using a function, returning a new iterator. |
| [print](functions/print/) | Prints a value to the console. Automatically converts the value to a string representation. |
| [range](functions/range/) | Creates an iterator that generates numbers from start to end (exclusive). |
| [readFile](functions/readfile/) | Reads the entire contents of a file as a string. |
| [sleep](functions/sleep/) | Pauses execution for the specified number of milliseconds. |
| [spawnProcess](functions/spawnprocess/) | Spawns an external async process with MANDATORY callback for stdout/stderr capture. The callback function receives (processID: int, eventType: int, data: string) and is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING! |
| [substring](functions/substring/) | Extracts a substring from start to end index. |
| [toString](functions/tostring/) | Converts a value to its string representation. |
| [webSocketKeepAlive](functions/websocketkeepalive/) | ⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. Keeps the WebSocket server running indefinitely until interrupted (blocking operation). |
| [websocketClose](functions/websocketclose/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Closes the WebSocket connection and cleans up resources. |
| [websocketConnect](functions/websocketconnect/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<WebSocketID, String> and takes string handler instead of function pointer. Establishes a WebSocket connection with a message handler callback. |
| [websocketCreateServer](functions/websocketcreateserver/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<ServerID, String> and has critical runtime issues with port binding failures. Creates a WebSocket server bound to the specified port, address, and path. |
| [websocketSend](functions/websocketsend/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Sends a message through the WebSocket connection. |
| [websocketServerBroadcast](functions/websocketserverbroadcast/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t (number of clients sent to) instead of Result<Success, String>. Broadcasts a message to all connected WebSocket clients. |
| [websocketServerListen](functions/websocketserverlisten/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String> and currently returns -4 (bind failed) due to port binding issues. Starts the WebSocket server listening for connections. |
| [websocketStopServer](functions/websocketstopserver/) | ⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of Result<Success, String>. Stops the WebSocket server and closes all connections. |
| [writeFile](functions/writefile/) | Writes content to a file. Creates the file if it doesn't exist. |

## Type Reference

| Type | Description |
|------|-------------|
| [Any](types/any/) | A type that can represent any value. Useful for generic programming but should be used carefully as it bypasses type checking. |
| [Bool](types/bool/) | A boolean type that can be either true or false. Used for logical operations and conditionals. |
| [HttpResponse](types/httpresponse/) | A built-in type representing an HTTP response with status code, headers, content type, body, and streaming capabilities. Used by HTTP server handlers to return structured responses to clients. |
| [Int](types/int/) | A 64-bit signed integer type. Can represent whole numbers from -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807. |
| [ProcessHandle](types/processhandle/) | A handle to a spawned async process. Contains the process ID and allows waiting for completion and cleanup. Process output is delivered via callbacks registered with the runtime. |
| [String](types/string/) | A sequence of characters representing text. Supports string interpolation and escape sequences. |

## Operator Reference

| Operator | Name | Description |
|----------|------|-------------|
| [!=](operators/not-equal/) | Inequality | Compares two values for inequality. |
| [%](operators/modulo/) | Modulo | Returns the remainder of dividing the first number by the second. |
| [*](operators/multiply/) | Multiplication | Multiplies two numbers. |
| [+](operators/plus/) | Addition | Adds two numbers together. |
| [-](operators/minus/) | Subtraction | Subtracts the second number from the first. |
| [/](operators/divide/) | Division | Divides the first number by the second. |
| [<](operators/less-than/) | Less Than | Checks if the first value is less than the second. |
| [<=](operators/less-equal/) | Less Than or Equal | Checks if the first value is less than or equal to the second. |
| [==](operators/equal/) | Equality | Compares two values for equality. |
| [>](operators/greater-than/) | Greater Than | Checks if the first value is greater than the second. |
| [>=](operators/greater-equal/) | Greater Than or Equal | Checks if the first value is greater than or equal to the second. |
| [|>](operators/pipe-operator/) | Pipe Operator | Takes the result of the left expression and passes it as the first argument to the right function. Enables functional programming and method chaining. |

## Keyword Reference

| Keyword | Description |
|---------|-------------|
| [false](keywords/false/) | Boolean literal representing the logical value false. |
| [fn](keywords/fn/) | Function declaration keyword. Used to define functions with parameters and return types. |
| [import](keywords/import/) | Import declaration keyword. Used to bring modules and their exports into the current scope. |
| [let](keywords/let/) | Variable declaration keyword. Used to bind values to identifiers. Variables are immutable by default in Osprey. |
| [match](keywords/match/) | Pattern matching expression. Used for destructuring values and control flow based on patterns. |
| [true](keywords/true/) | Boolean literal representing the logical value true. |
| [type](keywords/type/) | Type declaration keyword. Used to define custom types and type aliases. |

