// Package descriptions provides comprehensive documentation for built-in functions.
package descriptions

import (
	"github.com/christianfindlay/osprey/internal/codegen"
)

// BuiltinFunctionDesc represents documentation for a built-in function.
type BuiltinFunctionDesc struct {
	Name        string
	Signature   string
	Description string
	Parameters  []ParameterDesc
	ReturnType  string
	Example     string
}

// ParameterDesc represents documentation for a function parameter.
type ParameterDesc struct {
	Name        string
	Type        string
	Description string
}

// GetBuiltinFunctionDescriptions returns all built-in function descriptions from the unified registry.
func GetBuiltinFunctionDescriptions() map[string]*BuiltinFunctionDesc {
	// Convert from the unified registry format to the legacy format
	descriptions := make(map[string]*BuiltinFunctionDesc)
	
	for _, fn := range codegen.GlobalBuiltInRegistry.GetAllFunctions() {
		// Convert parameters
		params := make([]ParameterDesc, len(fn.ParameterTypes))
		for i, param := range fn.ParameterTypes {
			params[i] = ParameterDesc{
				Name:        param.Name,
				Type:        param.Type.String(),
				Description: param.Description,
			}
		}
		
		descriptions[fn.Name] = &BuiltinFunctionDesc{
			Name:        fn.Name,
			Signature:   fn.Signature,
			Description: fn.Description,
			Parameters:  params,
			ReturnType:  fn.ReturnType.String(),
			Example:     fn.Example,
		}
	}
	
	return descriptions
}

// GetBuiltinFunctionDescriptionsLegacy returns all built-in function descriptions using the old hardcoded method.
// This is kept for reference but should not be used in new code.
//
//nolint:maintidx // Large function with comprehensive function documentation
func GetBuiltinFunctionDescriptionsLegacy() map[string]*BuiltinFunctionDesc {
	return map[string]*BuiltinFunctionDesc{
		"print": {
			Name:      "print",
			Signature: "print(value: any) -> int",
			Description: "Prints a value to the console. " +
				"Automatically converts the value to a string representation.",
			Parameters: []ParameterDesc{
				{
					Name:        "value",
					Type:        "any",
					Description: "The value to print",
				},
			},
			ReturnType: "int",
			Example: `print("Hello, World!")  // Prints: Hello, World!\n` +
				`print(42)             // Prints: 42\n` +
				`print(true)           // Prints: true`,
		},
		"input": {
			Name:        "input",
			Signature:   "input() -> int",
			Description: "Reads an integer from the user's input.",
			Parameters:  []ParameterDesc{},
			ReturnType:  "int",
			Example:     `let userInput = input()\nprint(userInput)`,
		},
		"toString": {
			Name:        "toString",
			Signature:   "toString(value: any) -> string",
			Description: "Converts a value to its string representation.",
			Parameters: []ParameterDesc{
				{
					Name:        "value",
					Type:        "any",
					Description: "The value to convert to string",
				},
			},
			ReturnType: "string",
			Example:    `let str = toString(42)\nprint(str)  // Prints: 42`,
		},
		"range": {
			Name:      "range",
			Signature: "range(start: int, end: int) -> iterator",
			Description: "Creates an iterator that generates numbers from start to end " +
				"(exclusive).",
			Parameters: []ParameterDesc{
				{
					Name:        "start",
					Type:        "int",
					Description: "The starting number (inclusive)",
				},
				{
					Name:        "end",
					Type:        "int",
					Description: "The ending number (exclusive)",
				},
			},
			ReturnType: "iterator",
			Example:    `forEach(range(0, 5), fn(x) { print(x) })  // Prints: 0, 1, 2, 3, 4`,
		},
		"forEach": {
			Name:        "forEach",
			Signature:   "forEach(iterator: iterator, fn: function) -> int",
			Description: "Applies a function to each element in an iterator.",
			Parameters: []ParameterDesc{
				{
					Name:        "iterator",
					Type:        "iterator",
					Description: "The iterator to process",
				},
				{
					Name:        "fn",
					Type:        "function",
					Description: "The function to apply to each element",
				},
			},
			ReturnType: "int",
			Example:    `forEach(range(1, 4), fn(x) { print(x * 2) })  // Prints: 2, 4, 6`,
		},
		"map": {
			Name:      "map",
			Signature: "map(iterator: iterator, fn: function) -> iterator",
			Description: "Transforms each element in an iterator using a function, " +
				"returning a new iterator.",
			Parameters: []ParameterDesc{
				{
					Name:        "iterator",
					Type:        "iterator",
					Description: "The iterator to transform",
				},
				{
					Name:        "fn",
					Type:        "function",
					Description: "The transformation function",
				},
			},
			ReturnType: "iterator",
			Example:    `let doubled = map(range(1, 4), fn(x) { x * 2 })\nforEach(doubled, print)  // Prints: 2, 4, 6`,
		},
		"filter": {
			Name:        "filter",
			Signature:   "filter(iterator: iterator, predicate: function) -> iterator",
			Description: "Filters elements in an iterator based on a predicate function.",
			Parameters: []ParameterDesc{
				{
					Name:        "iterator",
					Type:        "iterator",
					Description: "The iterator to filter",
				},
				{
					Name: "predicate",
					Type: "function",
					Description: "The predicate function that returns true for " +
						"elements to keep",
				},
			},
			ReturnType: "iterator",
			Example:    `let evens = filter(range(1, 6), fn(x) { x % 2 == 0 })\nforEach(evens, print)  // Prints: 2, 4`,
		},
		"fold": {
			Name:        "fold",
			Signature:   "fold(iterator: iterator, initial: any, fn: function) -> any",
			Description: "Reduces an iterator to a single value using an accumulator function.",
			Parameters: []ParameterDesc{
				{
					Name:        "iterator",
					Type:        "iterator",
					Description: "The iterator to reduce",
				},
				{
					Name:        "initial",
					Type:        "any",
					Description: "The initial value for the accumulator",
				},
				{
					Name: "fn",
					Type: "function",
					Description: "The reduction function that takes (accumulator, current) " +
						"and returns new accumulator",
				},
			},
			ReturnType: "any",
			Example:    `let sum = fold(range(1, 5), 0, fn(acc, x) { acc + x })\nprint(sum)  // Prints: 10`,
		},
		"length": {
			Name:        "length",
			Signature:   "length(s: string) -> int",
			Description: "Returns the length of a string.",
			Parameters: []ParameterDesc{
				{
					Name:        "s",
					Type:        "string",
					Description: "The string to measure",
				},
			},
			ReturnType: "int",
			Example:    `let len = length("hello")\nprint(len)  // Prints: 5`,
		},
		"contains": {
			Name:        "contains",
			Signature:   "contains(haystack: string, needle: string) -> bool",
			Description: "Checks if a string contains a substring.",
			Parameters: []ParameterDesc{
				{
					Name:        "haystack",
					Type:        "string",
					Description: "The string to search in",
				},
				{
					Name:        "needle",
					Type:        "string",
					Description: "The substring to search for",
				},
			},
			ReturnType: "bool",
			Example:    `let found = contains("hello world", "world")\nprint(found)  // Prints: true`,
		},
		"substring": {
			Name:        "substring",
			Signature:   "substring(s: string, start: int, end: int) -> Result<string, Error>",
			Description: "Extracts a substring from start to end index, or returns an error if indices are invalid.",
			Parameters: []ParameterDesc{
				{
					Name:        "s",
					Type:        "string",
					Description: "The source string",
				},
				{
					Name:        "start",
					Type:        "int",
					Description: "Starting index (inclusive)",
				},
				{
					Name:        "end",
					Type:        "int",
					Description: "Ending index (exclusive)",
				},
			},
			ReturnType: "Result<string, Error>",
			Example:    `let sub = substring("hello", 1, 4)\nprint(sub)  // Prints: Result containing "ell"`,
		},

		// === CORE SYSTEM FUNCTIONS ===
		"spawnProcess": {
			Name:      "spawnProcess",
			Signature: "spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessHandle, string>",
			Description: "Spawns an external async process with MANDATORY callback for stdout/stderr capture. " +
				"The callback function receives (processID: int, eventType: int, data: string) and " +
				"is called for stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. " +
				"CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!",
			Parameters: []ParameterDesc{
				{
					Name:        "command",
					Type:        "string",
					Description: "The command to execute",
				},
				{
					Name:        "callback",
					Type:        "fn(int, int, string) -> Unit",
					Description: "MANDATORY callback function for process events (processID, eventType, data)",
				},
			},
			ReturnType: "Result<ProcessHandle, string>",
			Example: `fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {` + "\n" +
				`    match eventType {` + "\n" +
				`        1 => print("STDOUT: ${data}")` + "\n" +
				`        2 => print("STDERR: ${data}")` + "\n" +
				`        3 => print("EXIT: ${data}")` + "\n" +
				`        _ => print("Unknown event")` + "\n" +
				`    }` + "\n" +
				`}` + "\n" +
				`let result = spawnProcess("echo hello", processEventHandler)` + "\n" +
				`match result {` + "\n" +
				`    Success { value } => {` + "\n" +
				`        let exitCode = awaitProcess(value)` + "\n" +
				`        cleanupProcess(value)` + "\n" +
				`    }` + "\n" +
				`    Error { message } => print("Failed")` + "\n" +
				`}`,
		},
		"awaitProcess": {
			Name:        "awaitProcess",
			Signature:   "awaitProcess(handle: ProcessHandle) -> int",
			Description: "Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes.",
			Parameters: []ParameterDesc{
				{
					Name:        "handle",
					Type:        "ProcessHandle",
					Description: "Process handle from spawnProcess",
				},
			},
			ReturnType: "int",
			Example:    `let exitCode = awaitProcess(processHandle)\nprint("Process exited with code: ${toString(exitCode)}")`,
		},
		"cleanupProcess": {
			Name:        "cleanupProcess",
			Signature:   "cleanupProcess(handle: ProcessHandle) -> void",
			Description: "Cleans up resources associated with a completed process. Should be called after awaitProcess.",
			Parameters: []ParameterDesc{
				{
					Name:        "handle",
					Type:        "ProcessHandle",
					Description: "Process handle from spawnProcess",
				},
			},
			ReturnType: "void",
			Example:    `cleanupProcess(processHandle)  // Free process resources`,
		},
		"sleep": {
			Name:        "sleep",
			Signature:   "sleep(milliseconds: int) -> int",
			Description: "Pauses execution for the specified number of milliseconds.",
			Parameters: []ParameterDesc{
				{
					Name:        "milliseconds",
					Type:        "int",
					Description: "Number of milliseconds to sleep",
				},
			},
			ReturnType: "int",
			Example:    `sleep(1000)  // Sleep for 1 second\nprint("Awake!")`,
		},
		"writeFile": {
			Name:        "writeFile",
			Signature:   "writeFile(filename: string, content: string) -> Result<Unit, Error>",
			Description: "Writes content to a file. Creates the file if it doesn't exist.",
			Parameters: []ParameterDesc{
				{
					Name:        "filename",
					Type:        "string",
					Description: "Path to the file to write",
				},
				{
					Name:        "content",
					Type:        "string",
					Description: "Content to write to the file",
				},
			},
			ReturnType: "Result<Unit, Error>",
			Example:    `let result = writeFile("output.txt", "Hello, World!")\nprint("File written")`,
		},
		"readFile": {
			Name:        "readFile",
			Signature:   "readFile(filename: string) -> Result<string, Error>",
			Description: "Reads the entire contents of a file as a string.",
			Parameters: []ParameterDesc{
				{
					Name:        "filename",
					Type:        "string",
					Description: "Path to the file to read",
				},
			},
			ReturnType: "Result<string, Error>",
			Example:    `let content = readFile("input.txt")\nprint("File read")`,
		},

		// === HTTP SERVER FUNCTIONS ===
		"httpCreateServer": {
			Name:        "httpCreateServer",
			Signature:   "httpCreateServer(port: int, address: string) -> int",
			Description: "Creates an HTTP server bound to the specified port and address.",
			Parameters: []ParameterDesc{
				{
					Name:        "port",
					Type:        "int",
					Description: "Port number to bind to (1-65535)",
				},
				{
					Name:        "address",
					Type:        "string",
					Description: "IP address to bind to (e.g., \"127.0.0.1\", \"0.0.0.0\")",
				},
			},
			ReturnType: "int",
			Example:    `let serverId = httpCreateServer(8080, "127.0.0.1")\nprint("Server created with ID: ${serverId}")`,
		},
		"httpListen": {
			Name:        "httpListen",
			Signature:   "httpListen(serverID: int, handler: function) -> int",
			Description: "Starts the HTTP server listening for requests with a handler function.",
			Parameters: []ParameterDesc{
				{
					Name:        "serverID",
					Type:        "int",
					Description: "Server identifier from httpCreateServer",
				},
				{
					Name:        "handler",
					Type:        "function",
					Description: "Request handler function",
				},
			},
			ReturnType: "int",
			Example:    `let result = httpListen(serverId, requestHandler)\nprint("Server listening")`,
		},
		"httpStopServer": {
			Name:        "httpStopServer",
			Signature:   "httpStopServer(serverID: int) -> int",
			Description: "Stops the HTTP server and closes all connections.",
			Parameters: []ParameterDesc{
				{
					Name:        "serverID",
					Type:        "int",
					Description: "Server identifier to stop",
				},
			},
			ReturnType: "int",
			Example:    `let result = httpStopServer(serverId)\nprint("Server stopped")`,
		},

		// === HTTP CLIENT FUNCTIONS ===
		"httpCreateClient": {
			Name:        "httpCreateClient",
			Signature:   "httpCreateClient(baseUrl: string, timeout: int) -> int",
			Description: "Creates an HTTP client for making requests to a base URL.",
			Parameters: []ParameterDesc{
				{
					Name:        "baseUrl",
					Type:        "string",
					Description: "Base URL for requests (e.g., \"http://api.example.com\")",
				},
				{
					Name:        "timeout",
					Type:        "int",
					Description: "Request timeout in milliseconds",
				},
			},
			ReturnType: "int",
			Example:    `let clientId = httpCreateClient("http://httpbin.org", 5000)\nprint("Client created")`,
		},
		"httpGet": {
			Name:        "httpGet",
			Signature:   "httpGet(clientID: int, path: string, headers: string) -> int",
			Description: "Makes an HTTP GET request to the specified path.",
			Parameters: []ParameterDesc{
				{
					Name:        "clientID",
					Type:        "int",
					Description: "Client identifier from httpCreateClient",
				},
				{
					Name:        "path",
					Type:        "string",
					Description: "Request path (e.g., \"/api/users\")",
				},
				{
					Name:        "headers",
					Type:        "string",
					Description: "Additional headers (e.g., \"Authorization: Bearer token\")",
				},
			},
			ReturnType: "int",
			Example:    `let status = httpGet(clientId, "/get", "")\nprint("GET request status: ${status}")`,
		},
		"httpPost": {
			Name:        "httpPost",
			Signature:   "httpPost(clientID: int, path: string, body: string, headers: string) -> int",
			Description: "Makes an HTTP POST request with a request body.",
			Parameters: []ParameterDesc{
				{
					Name:        "clientID",
					Type:        "int",
					Description: "Client identifier from httpCreateClient",
				},
				{
					Name:        "path",
					Type:        "string",
					Description: "Request path",
				},
				{
					Name:        "body",
					Type:        "string",
					Description: "Request body data",
				},
				{
					Name:        "headers",
					Type:        "string",
					Description: "Additional headers",
				},
			},
			ReturnType: "int",
			Example: `let status = httpPost(clientId, "/post", "{\"key\":\"value\"}", ` +
				`"Content-Type: application/json")\nprint("POST status: ${status}")`,
		},
		"httpPut": {
			Name:        "httpPut",
			Signature:   "httpPut(clientID: int, path: string, body: string, headers: string) -> int",
			Description: "Makes an HTTP PUT request with a request body.",
			Parameters: []ParameterDesc{
				{
					Name:        "clientID",
					Type:        "int",
					Description: "Client identifier from httpCreateClient",
				},
				{
					Name:        "path",
					Type:        "string",
					Description: "Request path",
				},
				{
					Name:        "body",
					Type:        "string",
					Description: "Request body data",
				},
				{
					Name:        "headers",
					Type:        "string",
					Description: "Additional headers",
				},
			},
			ReturnType: "int",
			Example: `let status = httpPut(clientId, "/put", "{\"updated\":\"data\"}", ` +
				`"Content-Type: application/json")\nprint("PUT status: ${status}")`,
		},
		"httpDelete": {
			Name:        "httpDelete",
			Signature:   "httpDelete(clientID: int, path: string, headers: string) -> int",
			Description: "Makes an HTTP DELETE request to the specified path.",
			Parameters: []ParameterDesc{
				{
					Name:        "clientID",
					Type:        "int",
					Description: "Client identifier from httpCreateClient",
				},
				{
					Name:        "path",
					Type:        "string",
					Description: "Request path",
				},
				{
					Name:        "headers",
					Type:        "string",
					Description: "Additional headers",
				},
			},
			ReturnType: "int",
			Example:    `let status = httpDelete(clientId, "/delete", "")\nprint("DELETE status: ${status}")`,
		},
		"httpRequest": {
			Name:        "httpRequest",
			Signature:   "httpRequest(clientID: int, method: int, path: string, headers: string, body: string) -> int",
			Description: "Makes a generic HTTP request with any method.",
			Parameters: []ParameterDesc{
				{
					Name:        "clientID",
					Type:        "int",
					Description: "Client identifier from httpCreateClient",
				},
				{
					Name:        "method",
					Type:        "int",
					Description: "HTTP method (0=GET, 1=POST, 2=PUT, 3=DELETE)",
				},
				{
					Name:        "path",
					Type:        "string",
					Description: "Request path",
				},
				{
					Name:        "headers",
					Type:        "string",
					Description: "Additional headers",
				},
				{
					Name:        "body",
					Type:        "string",
					Description: "Request body data",
				},
			},
			ReturnType: "int",
			Example:    `let status = httpRequest(clientId, 0, "/custom", "", "")\nprint("Custom request status: ${status}")`,
		},
		"httpCloseClient": {
			Name:        "httpCloseClient",
			Signature:   "httpCloseClient(clientID: int) -> int",
			Description: "Closes the HTTP client and cleans up resources.",
			Parameters: []ParameterDesc{
				{
					Name:        "clientID",
					Type:        "int",
					Description: "Client identifier to close",
				},
			},
			ReturnType: "int",
			Example:    `let result = httpCloseClient(clientId)\nprint("Client closed")`,
		},

		// === WEBSOCKET CLIENT FUNCTIONS ===
		"websocketConnect": {
			Name: "websocketConnect",
			Signature: "websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> " +
				"Result<WebSocketID, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
				"Result<WebSocketID, String> and takes string handler instead of function pointer. " +
				"Establishes a WebSocket connection with a message handler callback.",
			Parameters: []ParameterDesc{
				{
					Name:        "url",
					Type:        "String",
					Description: "WebSocket URL (e.g., \"ws://localhost:8080/chat\")",
				},
				{
					Name:        "messageHandler",
					Type:        "fn(String) -> Result<Success, String>",
					Description: "Callback function to handle incoming messages",
				},
			},
			ReturnType: "Result<WebSocketID, String>",
			Example: `fn handleMessage(message: String) -> Result<Success, String> = {` + "\n" +
				`    print("Received: ${message}")` + "\n" +
				`    Success()` + "\n" +
				`}` + "\n" +
				`let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)` + "\n" +
				`match wsResult {` + "\n" +
				`    Success wsId => print("Connected with ID: ${wsId}")` + "\n" +
				`    Err message => print("Failed to connect: ${message}")` + "\n" +
				`}`,
		},
		"websocketSend": {
			Name:      "websocketSend",
			Signature: "websocketSend(wsID: Int, message: String) -> Result<Success, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
				"Result<Success, String>. Sends a message through the WebSocket connection.",
			Parameters: []ParameterDesc{
				{
					Name:        "wsID",
					Type:        "Int",
					Description: "WebSocket identifier from websocketConnect",
				},
				{
					Name:        "message",
					Type:        "String",
					Description: "Message to send",
				},
			},
			ReturnType: "Result<Success, String>",
			Example: `let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")` + "\n" +
				`match sendResult {` + "\n" +
				`    Success _ => print("Message sent successfully")` + "\n" +
				`    Err message => print("Failed to send: ${message}")` + "\n" +
				`}`,
		},
		"websocketClose": {
			Name:      "websocketClose",
			Signature: "websocketClose(wsID: Int) -> Result<Success, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
				"Result<Success, String>. Closes the WebSocket connection and cleans up resources.",
			Parameters: []ParameterDesc{
				{
					Name:        "wsID",
					Type:        "Int",
					Description: "WebSocket identifier to close",
				},
			},
			ReturnType: "Result<Success, String>",
			Example: `let closeResult = websocketClose(wsID: wsId)` + "\n" +
				`match closeResult {` + "\n" +
				`    Success _ => print("Connection closed")` + "\n" +
				`    Err message => print("Failed to close: ${message}")` + "\n" +
				`}`,
		},

		// === WEBSOCKET SERVER FUNCTIONS ===
		"websocketCreateServer": {
			Name:      "websocketCreateServer",
			Signature: "websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
				"Result<ServerID, String> and has critical runtime issues with port binding failures. " +
				"Creates a WebSocket server bound to the specified port, address, and path.",
			Parameters: []ParameterDesc{
				{
					Name:        "port",
					Type:        "Int",
					Description: "Port number to bind to (1-65535)",
				},
				{
					Name:        "address",
					Type:        "String",
					Description: "IP address to bind to (e.g., \"127.0.0.1\", \"0.0.0.0\")",
				},
				{
					Name:        "path",
					Type:        "String",
					Description: "WebSocket endpoint path (e.g., \"/chat\", \"/live\")",
				},
			},
			ReturnType: "Result<ServerID, String>",
			Example: `let serverResult = websocketCreateServer(port: 8080, address: "127.0.0.1", path: "/chat")` + "\n" +
				`match serverResult {` + "\n" +
				`    Success serverId => print("WebSocket server created with ID: ${serverId}")` + "\n" +
				`    Err message => print("Failed to create server: ${message}")` + "\n" +
				`}`,
		},
		"websocketServerListen": {
			Name:      "websocketServerListen",
			Signature: "websocketServerListen(serverID: Int) -> Result<Success, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
				"Result<Success, String> and currently returns -4 (bind failed) due to port binding issues. " +
				"Starts the WebSocket server listening for connections.",
			Parameters: []ParameterDesc{
				{
					Name:        "serverID",
					Type:        "Int",
					Description: "Server identifier from websocketCreateServer",
				},
			},
			ReturnType: "Result<Success, String>",
			Example: `let listenResult = websocketServerListen(serverID: serverId)` + "\n" +
				`match listenResult {` + "\n" +
				`    Success _ => print("Server listening on ws://127.0.0.1:8080/chat")` + "\n" +
				`    Err message => print("Failed to start listening: ${message}")` + "\n" +
				`}`,
		},
		"websocketServerBroadcast": {
			Name:      "websocketServerBroadcast",
			Signature: "websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t (number of clients " +
				"sent to) instead of Result<Success, String>. Broadcasts a message to all connected " +
				"WebSocket clients.",
			Parameters: []ParameterDesc{
				{
					Name:        "serverID",
					Type:        "Int",
					Description: "Server identifier",
				},
				{
					Name:        "message",
					Type:        "String",
					Description: "Message to broadcast to all clients",
				},
			},
			ReturnType: "Result<Success, String>",
			Example: `let broadcastResult = websocketServerBroadcast(serverID: serverId, ` +
				`message: "Welcome to Osprey Chat!")` + "\n" +
				`match broadcastResult {` + "\n" +
				`    Success _ => print("Message broadcasted to all clients")` + "\n" +
				`    Err message => print("Failed to broadcast: ${message}")` + "\n" +
				`}`,
		},
		"websocketStopServer": {
			Name:      "websocketStopServer",
			Signature: "websocketStopServer(serverID: Int) -> Result<Success, String>",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
				"Result<Success, String>. Stops the WebSocket server and closes all connections.",
			Parameters: []ParameterDesc{
				{
					Name:        "serverID",
					Type:        "Int",
					Description: "Server identifier to stop",
				},
			},
			ReturnType: "Result<Success, String>",
			Example: `let stopResult = websocketStopServer(serverID: serverId)` + "\n" +
				`match stopResult {` + "\n" +
				`    Success _ => print("Server stopped successfully")` + "\n" +
				`    Err message => print("Failed to stop server: ${message}")` + "\n" +
				`}`,
		},
		"webSocketKeepAlive": {
			Name:      "webSocketKeepAlive",
			Signature: "webSocketKeepAlive() -> Unit",
			Description: "⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. " +
				"Keeps the WebSocket server running indefinitely until interrupted (blocking operation).",
			Parameters: []ParameterDesc{},
			ReturnType: codegen.TypeUnit,
			Example:    `webSocketKeepAlive()  // Blocks until Ctrl+C`,
		},
	}
}

// GetBuiltinFunctionDescription returns description for a single built-in function.
func GetBuiltinFunctionDescription(name string) *BuiltinFunctionDesc {
	descriptions := GetBuiltinFunctionDescriptions()
	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return nil
}

// ValidateAllBuiltinFunctionsDocumented checks that all built-in functions are documented.
// This function should be called during build/test to ensure documentation completeness.
func ValidateAllBuiltinFunctionsDocumented() []string {
	// Get the authoritative list of built-in functions from the compiler's constants
	builtinFunctions := GetCompilerBuiltinFunctionNames()

	descriptions := GetBuiltinFunctionDescriptions()
	var missing []string

	for _, funcName := range builtinFunctions {
		if _, exists := descriptions[funcName]; !exists {
			missing = append(missing, funcName)
		}
	}

	return missing
}

// GetCompilerBuiltinFunctionNames returns all built-in function names from the unified registry.
// This is the authoritative source - it reads directly from the unified built-in function registry.
func GetCompilerBuiltinFunctionNames() []string {
	return codegen.GlobalBuiltInRegistry.GetFunctionNames()
}

// GetAllBuiltinFunctionNames returns all documented built-in function names.
// This can be used to cross-check against the actual compiler implementation.
func GetAllBuiltinFunctionNames() []string {
	descriptions := GetBuiltinFunctionDescriptions()
	var names []string
	for name := range descriptions {
		names = append(names, name)
	}
	return names
}
