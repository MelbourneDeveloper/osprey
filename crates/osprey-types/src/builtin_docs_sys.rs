//! Built-in documentation data (system, network & runtime). Generated companion to
//! `builtins.rs`: every entry's prose pairs with the type scheme of the
//! same name. Edit prose here; edit types in `builtins.rs`. The parity
//! test in `builtin_docs.rs` guarantees the two stay in lockstep.
//!
//! Param order and count MUST match the builtin's real arity.

use crate::builtin_docs::{BuiltinDoc, ParamDoc};

/// `files` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static FILES: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "readFile",
        summary: "Reads the entire contents of a file as a string.",
        params: &[ParamDoc { name: "filename", description: "Path to the file to read" }],
        example: "let content = readFile(\"input.txt\")\nprint(\"File read\")",
    },
    BuiltinDoc {
        name: "writeFile",
        summary: "Writes content to a file. Creates the file if it doesn't exist. Returns number of bytes written.",
        params: &[ParamDoc { name: "filename", description: "Path to the file to write" }, ParamDoc { name: "content", description: "Content to write to the file" }],
        example: "let result = writeFile(\"output.txt\", \"Hello, World!\")\nprint(\"File written\")",
    },
    BuiltinDoc {
        name: "deleteFile",
        summary: "Deletes the file at the given path, returning Unit on success or an error.",
        params: &[ParamDoc { name: "path", description: "Filesystem path of the file to delete" }],
        example: "match deleteFile(\"temp.txt\") {\n  Success { value } => print(\"deleted\")\n  Error { message } => print(message)\n}",
    },
];

/// `http` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static HTTP: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "httpCreateClient",
        summary: "Creates an HTTP client for making requests to a base URL.",
        params: &[ParamDoc { name: "base_url", description: "Base URL for requests (e.g., \"http://api.example.com\")" }, ParamDoc { name: "timeout", description: "Request timeout in milliseconds" }],
        example: "let clientId = httpCreateClient(\"http://httpbin.org\", 5000)\nprint(\"Client created\")",
    },
    BuiltinDoc {
        name: "httpCloseClient",
        summary: "Closes the HTTP client and cleans up resources.",
        params: &[ParamDoc { name: "clientID", description: "Client identifier to close" }],
        example: "let result = httpCloseClient(clientId)\nprint(\"Client closed\")",
    },
    BuiltinDoc {
        name: "httpGet",
        summary: "Makes an HTTP GET request to the specified path.",
        params: &[ParamDoc { name: "clientID", description: "Client identifier from httpCreateClient" }, ParamDoc { name: "path", description: "Request path (e.g., \"/api/users\")" }, ParamDoc { name: "headers", description: "Additional headers (e.g., \"Authorization: Bearer token\")" }],
        example: "let status = httpGet(clientId, \"/get\", \"\")\nprint(\"GET request status: ${status}\")",
    },
    BuiltinDoc {
        name: "httpGetResponse",
        summary: "Sends an HTTP GET request and returns a response handle for inspecting the status, headers, and body.",
        params: &[ParamDoc { name: "clientID", description: "Client identifier from httpCreateClient" }, ParamDoc { name: "path", description: "Request path, e.g. \"/api/users\"" }, ParamDoc { name: "headers", description: "Additional request headers, or \"\" for none" }],
        example: "match httpGetResponse(client, \"/users\", \"\") {\n  Success { value } => print(\"status: ${httpResponseStatus(value)}\")\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "httpResponseBody",
        summary: "Returns the body of a response handle as a string.",
        params: &[ParamDoc { name: "responseID", description: "Handle returned by httpGetResponse" }],
        example: "match httpResponseBody(response) {\n  Success { value } => print(value)\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "httpResponseFree",
        summary: "Releases a response handle obtained from httpGetResponse.",
        params: &[ParamDoc { name: "responseID", description: "Handle returned by httpGetResponse" }],
        example: "httpResponseFree(response)",
    },
    BuiltinDoc {
        name: "httpResponseStatus",
        summary: "Returns the HTTP status code of a response handle.",
        params: &[ParamDoc { name: "responseID", description: "Handle returned by httpGetResponse" }],
        example: "let code = httpResponseStatus(response)  // 200",
    },
    BuiltinDoc {
        name: "httpResponseHeader",
        summary: "Returns the value of the named header from a response handle.",
        params: &[ParamDoc { name: "responseID", description: "Handle returned by httpGetResponse" }, ParamDoc { name: "name", description: "Header name, e.g. \"Content-Type\"" }],
        example: "match httpResponseHeader(response, \"Content-Type\") {\n  Success { value } => print(value)\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "httpPost",
        summary: "Makes an HTTP POST request with a request body.",
        params: &[ParamDoc { name: "clientID", description: "Client identifier from httpCreateClient" }, ParamDoc { name: "path", description: "Request path" }, ParamDoc { name: "body", description: "Request body data" }, ParamDoc { name: "headers", description: "Additional headers" }],
        example: "let status = httpPost(clientId, \"/post\", \"{\\\"key\\\":\\\"value\\\"}\", \"Content-Type: application/json\")\nprint(\"POST status: ${status}\")",
    },
    BuiltinDoc {
        name: "httpPut",
        summary: "Makes an HTTP PUT request with a request body.",
        params: &[ParamDoc { name: "clientID", description: "Client identifier from httpCreateClient" }, ParamDoc { name: "path", description: "Request path" }, ParamDoc { name: "body", description: "Request body data" }, ParamDoc { name: "headers", description: "Additional headers" }],
        example: "let status = httpPut(clientId, \"/put\", \"{\\\"updated\\\":\\\"data\\\"}\", \"Content-Type: application/json\")\nprint(\"PUT status: ${status}\")",
    },
    BuiltinDoc {
        name: "httpDelete",
        summary: "Makes an HTTP DELETE request to the specified path.",
        params: &[ParamDoc { name: "clientID", description: "Client identifier from httpCreateClient" }, ParamDoc { name: "path", description: "Request path" }, ParamDoc { name: "headers", description: "Additional headers" }],
        example: "let status = httpDelete(clientId, \"/delete\", \"\")\nprint(\"DELETE status: ${status}\")",
    },
    BuiltinDoc {
        name: "httpCreateServer",
        summary: "Creates an HTTP server bound to the specified port and address.",
        params: &[ParamDoc { name: "port", description: "Port number to bind to (1-65535)" }, ParamDoc { name: "address", description: "IP address to bind to (e.g., \"127.0.0.1\", \"0.0.0.0\")" }],
        example: "let serverId = httpCreateServer(8080, \"127.0.0.1\")\nprint(\"Server created with ID: ${serverId}\")",
    },
    BuiltinDoc {
        name: "httpListen",
        summary: "Starts the HTTP server listening for requests with a handler function.",
        params: &[ParamDoc { name: "serverID", description: "Server identifier from httpCreateServer" }, ParamDoc { name: "handler", description: "Request handler function" }],
        example: "let result = httpListen(serverId, requestHandler)\nprint(\"Server listening\")",
    },
    BuiltinDoc {
        name: "httpStopServer",
        summary: "Stops the HTTP server and closes all connections.",
        params: &[ParamDoc { name: "serverID", description: "Server identifier to stop" }],
        example: "let result = httpStopServer(serverId)\nprint(\"Server stopped\")",
    },
];

/// `json` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static JSON: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "jsonParse",
        summary: "Parses a JSON string and returns an opaque document handle for querying, or an error on malformed input.",
        params: &[ParamDoc { name: "text", description: "The JSON text to parse" }],
        example: "match jsonParse(\"{\\\"name\\\": \\\"osprey\\\"}\") {\n  Success { value } => print(\"parsed\")\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "jsonGet",
        summary: "Returns the string value at the given path within a parsed JSON document.",
        params: &[ParamDoc { name: "document", description: "Handle returned by jsonParse" }, ParamDoc { name: "path", description: "Dotted path to the value, e.g. \"user.name\"" }],
        example: "match jsonGet(doc, \"name\") {\n  Success { value } => print(value)\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "jsonLength",
        summary: "Returns the number of elements in the JSON array at the given path.",
        params: &[ParamDoc { name: "document", description: "Handle returned by jsonParse" }, ParamDoc { name: "path", description: "Dotted path to the array" }],
        example: "let n = jsonLength(doc, \"items\")",
    },
    BuiltinDoc {
        name: "jsonFree",
        summary: "Releases a parsed JSON document handle obtained from jsonParse.",
        params: &[ParamDoc { name: "document", description: "Handle returned by jsonParse" }],
        example: "jsonFree(doc)",
    },
];

/// `concurrency` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static CONCURRENCY: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "await",
        summary: "Waits for a fiber to finish and returns its result, suspending the current fiber until then.",
        params: &[ParamDoc { name: "fiber", description: "The fiber to await" }],
        example: "let result = await(worker)",
    },
    BuiltinDoc {
        name: "fiberDone",
        summary: "Returns 1 if the given fiber has finished, 0 otherwise.",
        params: &[ParamDoc { name: "fiber", description: "The fiber to test" }],
        example: "let finished = fiberDone(worker)  // 0 or 1",
    },
    BuiltinDoc {
        name: "yield",
        summary: "Yields control from the current fiber, letting other ready fibers run.",
        params: &[],
        example: "yield()",
    },
    BuiltinDoc {
        name: "fiber_yield",
        summary: "Yields control to the fiber scheduler with an optional value.",
        params: &[ParamDoc { name: "value", description: "The value to yield" }],
        example: "let result = fiber_yield(42)",
    },
    BuiltinDoc {
        name: "Channel",
        summary: "Creates a new channel with the specified capacity.",
        params: &[ParamDoc { name: "capacity", description: "The capacity of the channel" }],
        example: "let ch = Channel(10)",
    },
    BuiltinDoc {
        name: "send",
        summary: "Sends a value to a channel. Returns 1 for success, 0 for failure.",
        params: &[ParamDoc { name: "channel", description: "The channel to send to" }, ParamDoc { name: "value", description: "The value to send" }],
        example: "let success = send(ch, 42)",
    },
    BuiltinDoc {
        name: "recv",
        summary: "Receives a value from a channel.",
        params: &[ParamDoc { name: "channel", description: "The channel to receive from" }],
        example: "let value = recv(ch)",
    },
];

/// `websocket` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static WEBSOCKET: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "websocketCreateServer",
        summary: "Creates a WebSocket server bound to the specified port, address, and path. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*",
        params: &[ParamDoc { name: "port", description: "Port number to bind to (1-65535)" }, ParamDoc { name: "address", description: "IP address to bind to (e.g., \"127.0.0.1\", \"0.0.0.0\")" }, ParamDoc { name: "path", description: "WebSocket endpoint path (e.g., \"/chat\", \"/live\")" }],
        example: "let serverResult = websocketCreateServer(port: 8080, address: \"127.0.0.1\", path: \"/chat\")\nmatch serverResult {\n    Success serverId => print(\"WebSocket server created with ID: ${serverId}\")\n    Err message => print(\"Failed to create server: ${message}\")\n}",
    },
    BuiltinDoc {
        name: "websocketServerListen",
        summary: "Starts the WebSocket server listening for connections. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*",
        params: &[ParamDoc { name: "serverID", description: "Server identifier from websocketCreateServer" }],
        example: "let listenResult = websocketServerListen(serverID: serverId)\nmatch listenResult {\n    Success _ => print(\"Server listening on ws://127.0.0.1:8080/chat\")\n    Err message => print(\"Failed to start listening: ${message}\")\n}",
    },
    BuiltinDoc {
        name: "websocketServerBroadcast",
        summary: "Broadcasts a message to all connected WebSocket clients. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*",
        params: &[ParamDoc { name: "serverID", description: "Server identifier" }, ParamDoc { name: "message", description: "Message to broadcast to all clients" }],
        example: "let broadcastResult = websocketServerBroadcast(serverID: serverId, message: \"Welcome to Osprey Chat!\")\nmatch broadcastResult {\n    Success _ => print(\"Message broadcasted to all clients\")\n    Err message => print(\"Failed to broadcast: ${message}\")\n}",
    },
    BuiltinDoc {
        name: "websocketKeepAlive",
        summary: "Keeps the WebSocket server running indefinitely until interrupted (blocking operation). *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*",
        params: &[],
        example: "websocketKeepAlive()  // Blocks until Ctrl+C",
    },
    BuiltinDoc {
        name: "websocketConnect",
        summary: "Connects to a WebSocket server at the given URL and returns a connection id.",
        params: &[ParamDoc { name: "url", description: "WebSocket URL, e.g. \"ws://localhost:8080/chat\"" }],
        example: "let conn = websocketConnect(\"ws://localhost:8080/chat\")",
    },
    BuiltinDoc {
        name: "websocketSend",
        summary: "Sends a message through the WebSocket connection. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*",
        params: &[ParamDoc { name: "wsID", description: "WebSocket identifier from websocketConnect" }, ParamDoc { name: "message", description: "Message to send" }],
        example: "let sendResult = websocketSend(wsID: wsId, message: \"Hello, WebSocket!\")\nmatch sendResult {\n    Success _ => print(\"Message sent successfully\")\n    Err message => print(\"Failed to send: ${message}\")\n}",
    },
    BuiltinDoc {
        name: "websocketClose",
        summary: "Closes the WebSocket connection and cleans up resources. *(Implementation note: currently returns an integer status code; the `Result`-typed API shown in the signature is planned.)*",
        params: &[ParamDoc { name: "wsID", description: "WebSocket identifier to close" }],
        example: "let closeResult = websocketClose(wsID: wsId)\nmatch closeResult {\n    Success _ => print(\"Connection closed\")\n    Err message => print(\"Failed to close: ${message}\")\n}",
    },
];

/// `terminal` built-in documentation. Prose only — types come from the
/// authoritative scheme in `builtins.rs`, joined by name.
pub(crate) static TERMINAL: &[BuiltinDoc] = &[
    BuiltinDoc {
        name: "termReadKey",
        summary: "Reads a single keypress from the terminal and returns it as a string.",
        params: &[],
        example: "match termReadKey() {\n  Success { value } => print(\"key: ${value}\")\n  Error { message } => print(message)\n}",
    },
    BuiltinDoc {
        name: "termRawMode",
        summary: "Enables (1) or disables (0) raw terminal input mode, so keypresses arrive unbuffered.",
        params: &[ParamDoc { name: "enabled", description: "1 to enable raw mode, 0 to restore cooked mode" }],
        example: "termRawMode(1)",
    },
    BuiltinDoc {
        name: "termCols",
        summary: "Returns the terminal width in columns.",
        params: &[],
        example: "let width = termCols()",
    },
    BuiltinDoc {
        name: "termRows",
        summary: "Returns the terminal height in rows.",
        params: &[],
        example: "let height = termRows()",
    },
    BuiltinDoc {
        name: "termClear",
        summary: "Clears the terminal screen.",
        params: &[],
        example: "termClear()",
    },
    BuiltinDoc {
        name: "termMoveCursor",
        summary: "Moves the terminal cursor to the given row and column.",
        params: &[ParamDoc { name: "row", description: "Target row (1-based)" }, ParamDoc { name: "col", description: "Target column (1-based)" }],
        example: "termMoveCursor(1, 1)",
    },
    BuiltinDoc {
        name: "termHideCursor",
        summary: "Hides the terminal cursor.",
        params: &[],
        example: "termHideCursor()",
    },
    BuiltinDoc {
        name: "termShowCursor",
        summary: "Shows the terminal cursor.",
        params: &[],
        example: "termShowCursor()",
    },
    BuiltinDoc {
        name: "spawnProcess",
        summary: "Spawns an external async process with MANDATORY callback for stdout/stderr capture. The callback function receives (processID: int, eventType: int, data: string) and is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!",
        params: &[ParamDoc { name: "command", description: "The command to execute" }, ParamDoc { name: "callback", description: "MANDATORY callback function for process events (processID, eventType, data)" }],
        example: "fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {\n    match eventType {\n        1 => print(\"STDOUT: ${data}\")\n        2 => print(\"STDERR: ${data}\")\n        3 => print(\"EXIT: ${data}\")\n        _ => print(\"Unknown event\")\n    }\n}\nlet result = spawnProcess(\"echo hello\", processEventHandler)\nmatch result {\n    Success { value } => {\n        let exitCode = awaitProcess(value)\n        cleanupProcess(value)\n    }\n    Error { message } => print(\"Failed\")\n}",
    },
    BuiltinDoc {
        name: "awaitProcess",
        summary: "Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes.",
        params: &[ParamDoc { name: "handle", description: "Process ID from spawnProcess" }],
        example: "let exitCode = awaitProcess(processHandle)\nprint(\"Process exited with code: ${toString(exitCode)}\")",
    },
    BuiltinDoc {
        name: "cleanupProcess",
        summary: "Cleans up resources associated with a completed process. Should be called after awaitProcess.",
        params: &[ParamDoc { name: "handle", description: "Process ID from spawnProcess" }],
        example: "cleanupProcess(processHandle)  // Free process resources",
    },
];
