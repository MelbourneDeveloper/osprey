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
| [contains](functions/contains/) | Checks if a string contains a substring. |
| [extractCode](functions/extractcode/) | Extracts code from a JSON structure. |
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
| [parseJSON](functions/parsejson/) | Parses a JSON string and returns the parsed result. |
| [print](functions/print/) | Prints a value to the console. Automatically converts the value to a string representation. |
| [range](functions/range/) | Creates an iterator that generates numbers from start to end (exclusive). |
| [readFile](functions/readfile/) | Reads the entire contents of a file as a string. |
| [sleep](functions/sleep/) | Pauses execution for the specified number of milliseconds. |
| [spawnProcess](functions/spawnprocess/) | Spawns an external process and returns the result. Currently supports simple command execution. |
| [substring](functions/substring/) | Extracts a substring from start to end index. |
| [toString](functions/tostring/) | Converts a value to its string representation. |
| [webSocketKeepAlive](functions/websocketkeepalive/) | Keeps the WebSocket server running indefinitely until interrupted. |
| [websocketClose](functions/websocketclose/) | Closes the WebSocket connection. |
| [websocketConnect](functions/websocketconnect/) | Establishes a WebSocket connection to the specified URL. |
| [websocketCreateServer](functions/websocketcreateserver/) | Creates a WebSocket server bound to the specified port, address, and path. |
| [websocketSend](functions/websocketsend/) | Sends a message through the WebSocket connection. |
| [websocketServerBroadcast](functions/websocketserverbroadcast/) | Broadcasts a message to all connected WebSocket clients. |
| [websocketServerListen](functions/websocketserverlisten/) | Starts the WebSocket server listening for connections. |
| [websocketStopServer](functions/websocketstopserver/) | Stops the WebSocket server and closes all connections. |
| [writeFile](functions/writefile/) | Writes content to a file. Creates the file if it doesn't exist. |

## Type Reference

| Type | Description |
|------|-------------|
| [Any](types/any/) | A type that can represent any value. Useful for generic programming but should be used carefully as it bypasses type checking. |
| [Bool](types/bool/) | A boolean type that can be either true or false. Used for logical operations and conditionals. |
| [HttpResponse](types/httpresponse/) | A built-in type representing an HTTP response with status code, headers, content type, body, and streaming capabilities. Used by HTTP server handlers to return structured responses to clients. |
| [Int](types/int/) | A 64-bit signed integer type. Can represent whole numbers from -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807. |
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

