package codegen

import (
	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/llir/llvm/ir/value"
)

// BuiltInFunctionRegistry contains all built-in function definitions in one place
type BuiltInFunctionRegistry struct {
	functions map[string]*BuiltInFunction
}

// BuiltInFunction represents a complete built-in function definition
type BuiltInFunction struct {
	// Basic information
	Name        string // The Osprey name
	CName       string // The function at the C level (if relevant). TODO: NULLABLE
	Signature   string
	Description string

	// Type information - TODO: Rename to ParameterDefinitions!!
	ParameterTypes []BuiltInParameter
	ReturnType     Type

	// Generation information
	Category FunctionCategory
	// TODO: DELETE THIS. Get these from the BuiltInParameter list instead!!
	//ExpectedArgs int
	IsProtected  bool
	SecurityFlag SecurityPermission

	// Code generation
	Generator func(g *LLVMGenerator, callExpr *ast.CallExpression) (value.Value, error)

	// Documentation
	Example string
}

// BuiltInParameter represents a function parameter
type BuiltInParameter struct {
	Name        string
	Type        Type
	Description string
}

// FunctionCategory represents the category of a built-in function
type FunctionCategory int

const (
	// CategoryCore represents core language functions like print, input
	CategoryCore FunctionCategory = iota
	// CategorySystem represents system-level functions
	CategorySystem
	// CategoryString represents string manipulation functions
	CategoryString
	// CategoryFunctional represents functional programming functions
	CategoryFunctional
	// CategoryFile represents file I/O functions
	CategoryFile
	// CategoryHTTP represents HTTP client/server functions
	CategoryHTTP
	// CategoryWebSocket represents WebSocket functions
	CategoryWebSocket
	// CategoryProcess represents process management functions
	CategoryProcess
	// CategoryIterator represents iterator functions
	CategoryIterator
	// CategoryFiber represents fiber and concurrency functions
	CategoryFiber
)

// SecurityPermission represents security permissions required for a function
type SecurityPermission int

const (
	// PermissionNone indicates no special permissions required
	PermissionNone SecurityPermission = iota
	// PermissionHTTP indicates HTTP access permission required
	PermissionHTTP
	// PermissionWebSocket indicates WebSocket access permission required
	PermissionWebSocket
	// PermissionFileRead indicates file read permission required
	PermissionFileRead
	// PermissionFileWrite indicates file write permission required
	PermissionFileWrite
	// PermissionProcess indicates process spawning permission required
	PermissionProcess
	// PermissionFFI indicates foreign function interface permission required
	PermissionFFI
)

// NewBuiltInFunctionRegistry creates and initializes the built-in function registry
func NewBuiltInFunctionRegistry() *BuiltInFunctionRegistry {
	registry := &BuiltInFunctionRegistry{
		functions: make(map[string]*BuiltInFunction),
	}
	registry.initializeFunctions()
	return registry
}

// GetFunction retrieves a built-in function by name
func (r *BuiltInFunctionRegistry) GetFunction(name string) (*BuiltInFunction, bool) {
	fn, exists := r.functions[name]
	return fn, exists
}

// GetAllFunctions returns all built-in functions
func (r *BuiltInFunctionRegistry) GetAllFunctions() map[string]*BuiltInFunction {
	return r.functions
}

// GetFunctionsByCategory returns functions by category
func (r *BuiltInFunctionRegistry) GetFunctionsByCategory(category FunctionCategory) []*BuiltInFunction {
	var functions []*BuiltInFunction
	for _, fn := range r.functions {
		if fn.Category == category {
			functions = append(functions, fn)
		}
	}
	return functions
}

// IsProtectedFunction checks if a function name is protected (built-in)
func (r *BuiltInFunctionRegistry) IsProtectedFunction(name string) bool {
	fn, exists := r.GetFunction(name)
	return exists && fn.IsProtected
}

// RequiresPermission checks if a function requires specific security permissions
func (r *BuiltInFunctionRegistry) RequiresPermission(name string, permission SecurityPermission) bool {
	fn, exists := r.GetFunction(name)
	return exists && fn.SecurityFlag == permission
}

// ValidateArguments validates function arguments
func (r *BuiltInFunctionRegistry) ValidateArguments(name string, argCount int, position *ast.Position) error {
	fn, exists := r.GetFunction(name)
	if !exists {
		return nil // Not a built-in function
	}

	expectedArgs := len(fn.ParameterTypes)
	if argCount != expectedArgs {
		return WrapFunctionArgsWithPos(name, expectedArgs, argCount, position)
	}

	return nil
}

// GetFunctionNames returns all built-in function names
func (r *BuiltInFunctionRegistry) GetFunctionNames() []string {
	names := make([]string, 0, len(r.functions))
	for name := range r.functions {
		names = append(names, name)
	}
	return names
}

// initializeFunctions initializes all built-in function definitions
func (r *BuiltInFunctionRegistry) initializeFunctions() {
	// Core I/O functions
	r.registerCoreIOFunctions()

	// String functions
	r.registerStringFunctions()

	// Functional programming functions
	r.registerFunctionalFunctions()

	// File I/O functions
	r.registerFileIOFunctions()

	// Process functions
	r.registerProcessFunctions()

	// HTTP functions
	r.registerHTTPFunctions()

	// WebSocket functions
	r.registerWebSocketFunctions()

	// System functions
	r.registerSystemFunctions()

	// Fiber functions
	r.registerFiberFunctions()
}

// registerCoreIOFunctions registers core I/O functions
func (r *BuiltInFunctionRegistry) registerCoreIOFunctions() {
	// print function
	r.functions[PrintFunc] = &BuiltInFunction{
		Name: PrintFunc,
		// 🪲
		//TODO: Incorrect. Print accepts a string as the only parameter - not any.
		Signature:   "print(value: any) -> Unit",
		Description: "Prints a value to the console. Automatically converts the value to a string representation.",
		ParameterTypes: []BuiltInParameter{
			//TODO: Incorrect. Print accepts a string as the only parameter - not any.
			{Name: "value", Type: &ConcreteType{name: TypeAny}, Description: "The value to print"},
		},
		ReturnType:   &ConcreteType{name: TypeUnit},
		Category:     CategoryCore,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generatePrintCall,
		Example: `print("Hello, World!")  // Prints: Hello, World!\n` +
			`print(42)             // Prints: 42\n` +
			`print(true)           // Prints: true`,
	}

	// toString function
	r.functions[ToStringFunc] = &BuiltInFunction{
		Name:        ToStringFunc,
		Signature:   "toString(value: any) -> string",
		Description: "Converts a value to its string representation.",
		ParameterTypes: []BuiltInParameter{
			{Name: "value", Type: &ConcreteType{name: TypeAny}, Description: "The value to convert to string"},
		},
		ReturnType:   &ConcreteType{name: TypeString},
		Category:     CategoryCore,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateToStringCall,
		Example:      `let str = toString(42)\nprint(str)  // Prints: 42`,
	}

	// input function
	r.functions[InputFunc] = &BuiltInFunction{
		Name:           InputFunc,
		Signature:      "input() -> int",
		Description:    "Reads an integer from the user's input.",
		ParameterTypes: []BuiltInParameter{},
		ReturnType:     &ConcreteType{name: "Result<int, Error>"},
		Category:       CategoryCore,
		IsProtected:    true,
		SecurityFlag:   PermissionNone,
		Generator:      (*LLVMGenerator).generateInputCall,
		Example:        `let userInput = input()\nprint(userInput)`,
	}
}

// registerStringFunctions registers string manipulation functions
func (r *BuiltInFunctionRegistry) registerStringFunctions() {
	// length function
	r.functions[LengthFunc] = &BuiltInFunction{
		Name:        LengthFunc,
		Signature:   "length(text: string) -> int",
		Description: "Returns the length of a string.",
		ParameterTypes: []BuiltInParameter{
			{Name: "text", Type: &ConcreteType{name: TypeString}, Description: "The string to measure"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryString,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateLengthCall,
		Example:      `let len = length("hello")\nprint(len)  // Prints: 5`,
	}

	// contains function
	r.functions[ContainsFunc] = &BuiltInFunction{
		Name:        ContainsFunc,
		Signature:   "contains(haystack: string, needle: string) -> bool",
		Description: "Checks if a string contains a substring.",
		ParameterTypes: []BuiltInParameter{
			{Name: "haystack", Type: &ConcreteType{name: TypeString}, Description: "The string to search in"},
			{Name: "needle", Type: &ConcreteType{name: TypeString}, Description: "The substring to search for"},
		},
		ReturnType:   &ConcreteType{name: "Result<bool, Error>"},
		Category:     CategoryString,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateContainsCall,
		Example:      `let found = contains("hello world", "world")\nprint(found)  // Prints: true`,
	}

	// substring function
	r.functions[SubstringFunc] = &BuiltInFunction{
		Name:        SubstringFunc,
		Signature:   "substring(s: string, start: int, end: int) -> Result<string, Error>",
		Description: "Extracts a substring from start to end index, or returns an error if indices are invalid.",
		ParameterTypes: []BuiltInParameter{
			{Name: "s", Type: &ConcreteType{name: TypeString}, Description: "The source string"},
			{Name: "start", Type: &ConcreteType{name: TypeInt}, Description: "Starting index (inclusive)"},
			{Name: "end", Type: &ConcreteType{name: TypeInt}, Description: "Ending index (exclusive)"},
		},
		ReturnType:   &ConcreteType{name: "Result<string, Error>"},
		Category:     CategoryString,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateSubstringCall,
		Example:      `let sub = substring("hello", 1, 4)\nprint(sub)  // Prints: Result containing "ell"`,
	}
}

// registerFunctionalFunctions registers functional programming functions
func (r *BuiltInFunctionRegistry) registerFunctionalFunctions() {
	// range function
	r.functions[RangeFunc] = &BuiltInFunction{
		Name:        RangeFunc,
		Signature:   "range(start: int, end: int) -> iterator",
		Description: "Creates an iterator that generates numbers from start to end (exclusive).",
		ParameterTypes: []BuiltInParameter{
			{Name: "start", Type: &ConcreteType{name: TypeInt}, Description: "The starting number (inclusive)"},
			{Name: "end", Type: &ConcreteType{name: TypeInt}, Description: "The ending number (exclusive)"},
		},
		ReturnType:   &ConcreteType{name: "Iterator<int>"},
		Category:     CategoryFunctional,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateRangeCall,
		Example:      `forEach(range(0, 5), fn(x) { print(x) })  // Prints: 0, 1, 2, 3, 4`,
	}

	// forEach function
	r.functions[ForEachFunc] = &BuiltInFunction{
		Name:        ForEachFunc,
		Signature:   "forEach(iterator: iterator, function: function) -> int",
		Description: "Applies a function to each element in an iterator.",
		ParameterTypes: []BuiltInParameter{
			{Name: "iterator", Type: &ConcreteType{name: "Iterator<T>"}, Description: "The iterator to process"},
			{Name: "function", Type: &ConcreteType{name: "T -> Unit"}, Description: "The function to apply to each element"},
		},
		ReturnType:   &ConcreteType{name: TypeUnit},
		Category:     CategoryFunctional,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateForEachCall,
		Example:      `forEach(range(1, 4), fn(x) { print(x * 2) })  // Prints: 2, 4, 6`,
	}

	// map function
	r.functions[MapFunc] = &BuiltInFunction{
		Name:        MapFunc,
		Signature:   "map(iterator: iterator, fn: function) -> iterator",
		Description: "Transforms each element in an iterator using a function, returning a new iterator.",
		ParameterTypes: []BuiltInParameter{
			{Name: "iterator", Type: &ConcreteType{name: "Iterator<T>"}, Description: "The iterator to transform"},
			{Name: "fn", Type: &ConcreteType{name: "T -> U"}, Description: "The transformation function"},
		},
		ReturnType:   &ConcreteType{name: "Iterator<U>"},
		Category:     CategoryFunctional,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateMapCall,
		Example:      `let doubled = map(range(1, 4), fn(x) { x * 2 })\nforEach(doubled, print)  // Prints: 2, 4, 6`,
	}

	// filter function
	r.functions[FilterFunc] = &BuiltInFunction{
		Name:        FilterFunc,
		Signature:   "filter(iterator: iterator, predicate: function) -> iterator",
		Description: "Filters elements in an iterator based on a predicate function.",
		ParameterTypes: []BuiltInParameter{
			{Name: "iterator", Type: &ConcreteType{name: "Iterator<T>"}, Description: "The iterator to filter"},
			{Name: "predicate", Type: &ConcreteType{name: "T -> bool"},
				Description: "The predicate function that returns true for elements to keep"},
		},
		ReturnType:   &ConcreteType{name: "Iterator<T>"},
		Category:     CategoryFunctional,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateFilterCall,
		Example:      `let evens = filter(range(1, 6), fn(x) { x % 2 == 0 })\nforEach(evens, print)  // Prints: 2, 4`,
	}

	// fold function
	r.functions[FoldFunc] = &BuiltInFunction{
		Name:        FoldFunc,
		Signature:   "fold(iterator: iterator, initial: any, fn: function) -> any",
		Description: "Reduces an iterator to a single value using an accumulator function.",
		ParameterTypes: []BuiltInParameter{
			{Name: "iterator", Type: &ConcreteType{name: "Iterator<T>"}, Description: "The iterator to reduce"},
			{Name: "initial", Type: &ConcreteType{name: "U"}, Description: "The initial value for the accumulator"},
			{Name: "fn", Type: &ConcreteType{name: "(U, T) -> U"},
				Description: "The reduction function that takes (accumulator, current) and returns new accumulator"},
		},
		ReturnType:   &ConcreteType{name: "U"},
		Category:     CategoryFunctional,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateFoldCall,
		Example:      `let sum = fold(range(1, 5), 0, fn(acc, x) { acc + x })\nprint(sum)  // Prints: 10`,
	}
}

// registerFileIOFunctions registers file I/O functions
func (r *BuiltInFunctionRegistry) registerFileIOFunctions() {
	// readFile function
	r.functions[ReadFileFunc] = &BuiltInFunction{
		Name:        ReadFileFunc,
		Signature:   "readFile(filename: string) -> Result<string, Error>",
		Description: "Reads the entire contents of a file as a string.",
		ParameterTypes: []BuiltInParameter{
			{Name: "filename", Type: &ConcreteType{name: TypeString}, Description: "Path to the file to read"},
		},
		ReturnType:   &ConcreteType{name: "Result<string, Error>"},
		Category:     CategoryFile,
		IsProtected:  true,
		SecurityFlag: PermissionFileRead,
		Generator:    (*LLVMGenerator).generateReadFileCall,
		Example:      `let content = readFile("input.txt")\nprint("File read")`,
	}

	// writeFile function
	r.functions[WriteFileFunc] = &BuiltInFunction{
		Name:        WriteFileFunc,
		Signature:   "writeFile(filename: string, content: string) -> Result<Unit, Error>",
		Description: "Writes content to a file. Creates the file if it doesn't exist.",
		ParameterTypes: []BuiltInParameter{
			{Name: "filename", Type: &ConcreteType{name: TypeString}, Description: "Path to the file to write"},
			{Name: "content", Type: &ConcreteType{name: TypeString}, Description: "Content to write to the file"},
		},
		ReturnType:   &ConcreteType{name: "Result<Unit, Error>"},
		Category:     CategoryFile,
		IsProtected:  true,
		SecurityFlag: PermissionFileWrite,
		Generator:    (*LLVMGenerator).generateWriteFileCall,
		Example:      `let result = writeFile("output.txt", "Hello, World!")\nprint("File written")`,
	}
}

// registerProcessFunctions registers process management functions
func (r *BuiltInFunctionRegistry) registerProcessFunctions() {
	// spawnProcess function
	r.functions[SpawnProcessFunc] = &BuiltInFunction{
		Name:      SpawnProcessFunc,
		Signature: "spawnProcess(command: string, callback: fn(int, int, string) -> Unit) -> Result<ProcessHandle, string>",
		Description: "Spawns an external async process with MANDATORY callback for stdout/stderr capture. " +
			"The callback function receives (processID: int, eventType: int, data: string) and is called for " +
			"stdout (1), stderr (2), and exit (3) events. Returns a handle for the running process. " +
			"CALLBACK IS REQUIRED - NO FUNCTION OVERLOADING!",
		ParameterTypes: []BuiltInParameter{
			{Name: "command", Type: &ConcreteType{name: TypeString}, Description: "The command to execute"},
			{Name: "callback", Type: &ConcreteType{name: "fn(int, int, string) -> Unit"},
				Description: "MANDATORY callback function for process events (processID, eventType, data)"},
		},
		ReturnType:   &ConcreteType{name: "Result<ProcessHandle, string>"},
		Category:     CategoryProcess,
		IsProtected:  true,
		SecurityFlag: PermissionProcess,
		Generator:    (*LLVMGenerator).generateSpawnProcessCall,
		Example: `fn processEventHandler(processID: int, eventType: int, data: string) -> Unit = {
    match eventType {
        1 => print("STDOUT: ${data}")
        2 => print("STDERR: ${data}")
        3 => print("EXIT: ${data}")
        _ => print("Unknown event")
    }
}
let result = spawnProcess("echo hello", processEventHandler)
match result {
    Success { value } => {
        let exitCode = awaitProcess(value)
        cleanupProcess(value)
    }
    Error { message } => print("Failed")
}`,
	}

	// awaitProcess function
	r.functions[AwaitProcessFunc] = &BuiltInFunction{
		Name:        AwaitProcessFunc,
		Signature:   "awaitProcess(handle: ProcessHandle) -> int",
		Description: "Waits for a spawned process to complete and returns its exit code. Blocks until the process finishes.",
		ParameterTypes: []BuiltInParameter{
			{Name: "handle", Type: &ConcreteType{name: "ProcessHandle"}, Description: "Process handle from spawnProcess"},
		},
		ReturnType:   &ConcreteType{name: "Result<int, Error>"},
		Category:     CategoryProcess,
		IsProtected:  true,
		SecurityFlag: PermissionProcess,
		Generator:    (*LLVMGenerator).generateAwaitProcessCall,
		Example:      `let exitCode = awaitProcess(processHandle)\nprint("Process exited with code: ${toString(exitCode)}")`,
	}

	// cleanupProcess function
	r.functions[CleanupProcessFunc] = &BuiltInFunction{
		Name:        CleanupProcessFunc,
		Signature:   "cleanupProcess(handle: ProcessHandle) -> void",
		Description: "Cleans up resources associated with a completed process. Should be called after awaitProcess.",
		ParameterTypes: []BuiltInParameter{
			{Name: "handle", Type: &ConcreteType{name: "ProcessHandle"}, Description: "Process handle from spawnProcess"},
		},
		ReturnType:   &ConcreteType{name: "Result<Unit, Error>"},
		Category:     CategoryProcess,
		IsProtected:  true,
		SecurityFlag: PermissionProcess,
		Generator:    (*LLVMGenerator).generateCleanupProcessCall,
		Example:      `cleanupProcess(processHandle)  // Free process resources`,
	}

	// sleep function
	r.functions[SleepFunc] = &BuiltInFunction{
		Name:        SleepFunc,
		Signature:   "sleep(milliseconds: int) -> int",
		Description: "Pauses execution for the specified number of milliseconds.",
		ParameterTypes: []BuiltInParameter{
			{Name: "milliseconds", Type: &ConcreteType{name: TypeInt}, Description: "Number of milliseconds to sleep"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategorySystem,
		IsProtected:  true,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateSleepCall,
		Example:      `sleep(1000)  // Sleep for 1 second\nprint("Awake!")`,
	}
}

// registerHTTPFunctions registers HTTP-related functions
func (r *BuiltInFunctionRegistry) registerHTTPFunctions() {
	// httpCreateServer function
	r.functions[HTTPCreateServerOsprey] = &BuiltInFunction{
		Name:        HTTPCreateServerOsprey,
		CName:       HTTPCreateServerFunc,
		Signature:   "httpCreateServer(port: int, address: string) -> int",
		Description: "Creates an HTTP server bound to the specified port and address.",
		ParameterTypes: []BuiltInParameter{
			{Name: "port", Type: &ConcreteType{name: TypeInt}, Description: "Port number to bind to (1-65535)"},
			{Name: "address", Type: &ConcreteType{name: TypeString},
				Description: "IP address to bind to (e.g., \"127.0.0.1\", \"0.0.0.0\")"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPCreateServerCall,
		Example:      `let serverId = httpCreateServer(8080, "127.0.0.1")\nprint("Server created with ID: ${serverId}")`,
	}

	// httpListen function
	r.functions[HTTPListenOsprey] = &BuiltInFunction{
		Name:        HTTPListenOsprey,
		CName:       HTTPListenFunc,
		Signature:   "httpListen(serverID: int, handler: (string, string, string, string) -> HttpResponse) -> int",
		Description: "Starts the HTTP server listening for requests with a handler function.",
		ParameterTypes: []BuiltInParameter{
			{Name: "serverID", Type: &ConcreteType{name: TypeInt}, Description: "Server identifier from httpCreateServer"},
			{Name: "handler", Type: NewFunctionType(
				[]Type{
					&ConcreteType{name: TypeString},
					&ConcreteType{name: TypeString},
					&ConcreteType{name: TypeString},
					&ConcreteType{name: TypeString},
				},
				&ConcreteType{name: TypeHTTPResponse},
			), Description: "Request handler function"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPListenCall,
		Example:      `let result = httpListen(serverId, requestHandler)\nprint("Server listening")`,
	}

	// httpStopServer function
	r.functions[HTTPStopServerOsprey] = &BuiltInFunction{
		Name:        HTTPStopServerOsprey,
		CName:       HTTPStopServerFunc,
		Signature:   "httpStopServer(serverID: int) -> int",
		Description: "Stops the HTTP server and closes all connections.",
		ParameterTypes: []BuiltInParameter{
			{Name: "serverID", Type: &ConcreteType{name: TypeInt}, Description: "Server identifier to stop"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPStopServerCall,
		Example:      `let result = httpStopServer(serverId)\nprint("Server stopped")`,
	}

	// httpCreateClient function
	r.functions[HTTPCreateClientOsprey] = &BuiltInFunction{
		Name:        HTTPCreateClientOsprey,
		CName:       HTTPCreateClientFunc,
		Signature:   "httpCreateClient(base_url: string, timeout: int) -> int",
		Description: "Creates an HTTP client for making requests to a base URL.",
		ParameterTypes: []BuiltInParameter{
			{Name: "base_url", Type: &ConcreteType{name: TypeString},
				Description: "Base URL for requests (e.g., \"http://api.example.com\")"},
			{Name: "timeout", Type: &ConcreteType{name: TypeInt}, Description: "Request timeout in milliseconds"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPCreateClientCall,
		Example:      `let clientId = httpCreateClient("http://httpbin.org", 5000)\nprint("Client created")`,
	}

	// httpGet function
	r.functions[HTTPGetOsprey] = &BuiltInFunction{
		Name:        HTTPGetOsprey,
		CName:       HTTPGetFunc,
		Signature:   "httpGet(clientID: int, path: string, headers: string) -> int",
		Description: "Makes an HTTP GET request to the specified path.",
		ParameterTypes: []BuiltInParameter{
			{Name: "clientID", Type: &ConcreteType{name: TypeInt}, Description: "Client identifier from httpCreateClient"},
			{Name: "path", Type: &ConcreteType{name: TypeString}, Description: "Request path (e.g., \"/api/users\")"},
			{Name: "headers", Type: &ConcreteType{name: TypeString},
				Description: "Additional headers (e.g., \"Authorization: Bearer token\")"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPGetCall,
		Example:      `let status = httpGet(clientId, "/get", "")\nprint("GET request status: ${status}")`,
	}

	// httpPost function
	r.functions[HTTPPostOsprey] = &BuiltInFunction{
		Name:        HTTPPostOsprey,
		CName:       HTTPPostFunc,
		Signature:   "httpPost(clientID: int, path: string, body: string, headers: string) -> int",
		Description: "Makes an HTTP POST request with a request body.",
		ParameterTypes: []BuiltInParameter{
			{Name: "clientID", Type: &ConcreteType{name: TypeInt}, Description: "Client identifier from httpCreateClient"},
			{Name: "path", Type: &ConcreteType{name: TypeString}, Description: "Request path"},
			{Name: "body", Type: &ConcreteType{name: TypeString}, Description: "Request body data"},
			{Name: "headers", Type: &ConcreteType{name: TypeString}, Description: "Additional headers"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPPostCall,
		Example: `let status = httpPost(clientId, "/post", "{\"key\":\"value\"}", "Content-Type: application/json")\n` +
			`print("POST status: ${status}")`,
	}

	// httpPut function
	r.functions[HTTPPutOsprey] = &BuiltInFunction{
		Name:        HTTPPutOsprey,
		CName:       HTTPPutFunc,
		Signature:   "httpPut(clientID: int, path: string, body: string, headers: string) -> int",
		Description: "Makes an HTTP PUT request with a request body.",
		ParameterTypes: []BuiltInParameter{
			{Name: "clientID", Type: &ConcreteType{name: TypeInt}, Description: "Client identifier from httpCreateClient"},
			{Name: "path", Type: &ConcreteType{name: TypeString}, Description: "Request path"},
			{Name: "body", Type: &ConcreteType{name: TypeString}, Description: "Request body data"},
			{Name: "headers", Type: &ConcreteType{name: TypeString}, Description: "Additional headers"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPPutCall,
		Example: `let status = httpPut(clientId, "/put", "{\"updated\":\"data\"}", "Content-Type: application/json")\n` +
			`print("PUT status: ${status}")`,
	}

	// httpDelete function
	r.functions[HTTPDeleteOsprey] = &BuiltInFunction{
		Name:        HTTPDeleteOsprey,
		CName:       HTTPDeleteFunc,
		Signature:   "httpDelete(clientID: int, path: string, headers: string) -> int",
		Description: "Makes an HTTP DELETE request to the specified path.",
		ParameterTypes: []BuiltInParameter{
			{Name: "clientID", Type: &ConcreteType{name: TypeInt}, Description: "Client identifier from httpCreateClient"},
			{Name: "path", Type: &ConcreteType{name: TypeString}, Description: "Request path"},
			{Name: "headers", Type: &ConcreteType{name: TypeString}, Description: "Additional headers"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPDeleteCall,
		Example:      `let status = httpDelete(clientId, "/delete", "")\nprint("DELETE status: ${status}")`,
	}

	// httpRequest function
	r.functions[HTTPRequestOsprey] = &BuiltInFunction{
		Name:        HTTPRequestOsprey,
		CName:       HTTPRequestFunc,
		Signature:   "httpRequest(clientID: int, method: int, path: string, headers: string, body: string) -> int",
		Description: "Makes a generic HTTP request with any method.",
		ParameterTypes: []BuiltInParameter{
			{Name: "clientID", Type: &ConcreteType{name: TypeInt}, Description: "Client identifier from httpCreateClient"},
			{Name: "method", Type: &ConcreteType{name: TypeInt}, Description: "HTTP method (0=GET, 1=POST, 2=PUT, 3=DELETE)"},
			{Name: "path", Type: &ConcreteType{name: TypeString}, Description: "Request path"},
			{Name: "headers", Type: &ConcreteType{name: TypeString}, Description: "Additional headers"},
			{Name: "body", Type: &ConcreteType{name: TypeString}, Description: "Request body data"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPRequestCall,
		Example:      `let status = httpRequest(clientId, 0, "/custom", "", "")\nprint("Custom request status: ${status}")`,
	}

	// httpCloseClient function
	r.functions[HTTPCloseClientOsprey] = &BuiltInFunction{
		Name:        HTTPCloseClientOsprey,
		CName:       HTTPCloseClientFunc,
		Signature:   "httpCloseClient(clientID: int) -> int",
		Description: "Closes the HTTP client and cleans up resources.",
		ParameterTypes: []BuiltInParameter{
			{Name: "clientID", Type: &ConcreteType{name: TypeInt}, Description: "Client identifier to close"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryHTTP,
		IsProtected:  true,
		SecurityFlag: PermissionHTTP,
		Generator:    (*LLVMGenerator).generateHTTPCloseClientCall,
		Example:      `let result = httpCloseClient(clientId)\nprint("Client closed")`,
	}
}

// registerWebSocketFunctions registers WebSocket-related functions
func (r *BuiltInFunctionRegistry) registerWebSocketFunctions() {
	// websocketConnect function
	r.functions[WebSocketConnectOsprey] = &BuiltInFunction{
		Name:  WebSocketConnectOsprey,
		CName: WebSocketConnectFunc,
		Signature: "websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> " +
			"Result<WebSocketID, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
			"Result<WebSocketID, String> and takes string handler instead of function pointer. " +
			"Establishes a WebSocket connection with a message handler callback.",
		ParameterTypes: []BuiltInParameter{
			{Name: "url", Type: &ConcreteType{name: TypeString},
				Description: "WebSocket URL (e.g., \"ws://localhost:8080/chat\")"},
			{Name: "messageHandler", Type: &ConcreteType{name: "fn(String) -> Result<Success, String>"},
				Description: "Callback function to handle incoming messages"},
		},
		ReturnType:   &ConcreteType{name: "Result<WebSocketID, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketConnectCall,
		Example: `fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}
let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)
match wsResult {
    Success wsId => print("Connected with ID: ${wsId}")
    Err message => print("Failed to connect: ${message}")
}`,
	}

	// websocketSend function
	r.functions[WebSocketSendOsprey] = &BuiltInFunction{
		Name:      WebSocketSendOsprey,
		CName:     WebSocketSendFunc,
		Signature: "websocketSend(wsID: Int, message: String) -> Result<Success, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
			"Result<Success, String>. Sends a message through the WebSocket connection.",
		ParameterTypes: []BuiltInParameter{
			{Name: "wsID", Type: &ConcreteType{name: TypeInt}, Description: "WebSocket identifier from websocketConnect"},
			{Name: "message", Type: &ConcreteType{name: TypeString}, Description: "Message to send"},
		},
		ReturnType:   &ConcreteType{name: "Result<Success, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketSendCall,
		Example: `let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
match sendResult {
    Success _ => print("Message sent successfully")
    Err message => print("Failed to send: ${message}")
}`,
	}

	// websocketClose function
	r.functions[WebSocketCloseOsprey] = &BuiltInFunction{
		Name:      WebSocketCloseOsprey,
		CName:     WebSocketCloseFunc,
		Signature: "websocketClose(wsID: Int) -> Result<Success, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
			"Result<Success, String>. Closes the WebSocket connection and cleans up resources.",
		ParameterTypes: []BuiltInParameter{
			{Name: "wsID", Type: &ConcreteType{name: TypeInt}, Description: "WebSocket identifier to close"},
		},
		ReturnType:   &ConcreteType{name: "Result<Success, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketCloseCall,
		Example: `let closeResult = websocketClose(wsID: wsId)
match closeResult {
    Success _ => print("Connection closed")
    Err message => print("Failed to close: ${message}")
}`,
	}

	// websocketCreateServer function
	r.functions[WebSocketCreateServerOsprey] = &BuiltInFunction{
		Name:      WebSocketCreateServerOsprey,
		CName:     WebSocketCreateServerFunc,
		Signature: "websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
			"Result<ServerID, String> and has critical runtime issues with port binding failures. " +
			"Creates a WebSocket server bound to the specified port, address, and path.",
		ParameterTypes: []BuiltInParameter{
			{Name: "port", Type: &ConcreteType{name: TypeInt}, Description: "Port number to bind to (1-65535)"},
			{Name: "address", Type: &ConcreteType{name: TypeString},
				Description: "IP address to bind to (e.g., \"127.0.0.1\", \"0.0.0.0\")"},
			{Name: "path", Type: &ConcreteType{name: TypeString},
				Description: "WebSocket endpoint path (e.g., \"/chat\", \"/live\")"},
		},
		ReturnType:   &ConcreteType{name: "Result<ServerID, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketCreateServerCall,
		Example: `let serverResult = websocketCreateServer(port: 8080, address: "127.0.0.1", path: "/chat")
match serverResult {
    Success serverId => print("WebSocket server created with ID: ${serverId}")
    Err message => print("Failed to create server: ${message}")
}`,
	}

	// websocketServerListen function
	r.functions[WebSocketServerListenOsprey] = &BuiltInFunction{
		Name:      WebSocketServerListenOsprey,
		CName:     WebSocketServerListenFunc,
		Signature: "websocketServerListen(serverID: Int) -> Result<Success, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
			"Result<Success, String> and currently returns -4 (bind failed) due to port binding issues. " +
			"Starts the WebSocket server listening for connections.",
		ParameterTypes: []BuiltInParameter{
			{Name: "serverID", Type: &ConcreteType{name: TypeInt}, Description: "Server identifier from websocketCreateServer"},
		},
		ReturnType:   &ConcreteType{name: "Result<Success, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketServerListenCall,
		Example: `let listenResult = websocketServerListen(serverID: serverId)
match listenResult {
    Success _ => print("Server listening on ws://127.0.0.1:8080/chat")
    Err message => print("Failed to start listening: ${message}")
}`,
	}

	// websocketServerBroadcast function
	r.functions[WebSocketServerBroadcastOsprey] = &BuiltInFunction{
		Name:      WebSocketServerBroadcastOsprey,
		CName:     WebSocketServerBroadcastFunc,
		Signature: "websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t (number of " +
			"clients sent to) instead of Result<Success, String>. Broadcasts a message to all " +
			"connected WebSocket clients.",
		ParameterTypes: []BuiltInParameter{
			{Name: "serverID", Type: &ConcreteType{name: TypeInt}, Description: "Server identifier"},
			{Name: "message", Type: &ConcreteType{name: TypeString}, Description: "Message to broadcast to all clients"},
		},
		ReturnType:   &ConcreteType{name: "Result<Success, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketServerBroadcastCall,
		Example: `let broadcastResult = websocketServerBroadcast(serverID: serverId, message: "Welcome to Osprey Chat!")
match broadcastResult {
    Success _ => print("Message broadcasted to all clients")
    Err message => print("Failed to broadcast: ${message}")
}`,
	}

	// websocketStopServer function
	r.functions[WebSocketStopServerOsprey] = &BuiltInFunction{
		Name:      WebSocketStopServerOsprey,
		CName:     WebSocketStopServerFunc,
		Signature: "websocketStopServer(serverID: Int) -> Result<Success, String>",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns raw int64_t instead of " +
			"Result<Success, String>. Stops the WebSocket server and closes all connections.",
		ParameterTypes: []BuiltInParameter{
			{Name: "serverID", Type: &ConcreteType{name: TypeInt}, Description: "Server identifier to stop"},
		},
		ReturnType:   &ConcreteType{name: "Result<Success, String>"},
		Category:     CategoryWebSocket,
		IsProtected:  true,
		SecurityFlag: PermissionWebSocket,
		Generator:    (*LLVMGenerator).generateWebSocketStopServerCall,
		Example: `let stopResult = websocketStopServer(serverID: serverId)
match stopResult {
    Success _ => print("Server stopped successfully")
    Err message => print("Failed to stop server: ${message}")
}`,
	}
}

// registerSystemFunctions registers system-related functions
func (r *BuiltInFunctionRegistry) registerSystemFunctions() {
	// webSocketKeepAlive function
	r.functions[WebSocketKeepAliveOsprey] = &BuiltInFunction{
		Name:      WebSocketKeepAliveOsprey,
		CName:     WebSocketKeepAliveFunc,
		Signature: "websocket_keep_alive() -> Unit",
		Description: "⚠️ SPEC VIOLATION: Current implementation returns int instead of Unit. " +
			"Keeps the WebSocket server running indefinitely until interrupted (blocking operation).",
		ParameterTypes: []BuiltInParameter{},
		ReturnType:     &ConcreteType{name: TypeUnit},
		Category:       CategoryWebSocket,
		IsProtected:    true,
		SecurityFlag:   PermissionWebSocket,
		Generator:      (*LLVMGenerator).generateWebSocketKeepAliveCall,
		Example:        `webSocketKeepAlive()  // Blocks until Ctrl+C`,
	}

}

// registerFiberFunctions registers fiber and concurrency functions
func (r *BuiltInFunctionRegistry) registerFiberFunctions() {
	// fiber_spawn function
	r.functions["fiber_spawn"] = &BuiltInFunction{
		Name:        "fiber_spawn",
		CName:       "fiber_spawn",
		Signature:   "fiber_spawn(fn: () -> any) -> Fiber",
		Description: "Spawns a new fiber to execute the given function concurrently.",
		ParameterTypes: []BuiltInParameter{
			{Name: "fn", Type: &ConcreteType{name: "function"}, Description: "The function to execute in the fiber"},
		},
		ReturnType:   &ConcreteType{name: "Fiber"},
		Category:     CategoryFiber,
		IsProtected:  false,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateSpawnCall,
		Example:      `let fiber = fiber_spawn(() -> print("Hello from fiber"))`,
	}

	// fiber_yield function
	r.functions["fiber_yield"] = &BuiltInFunction{
		Name:        "fiber_yield",
		CName:       "fiber_yield",
		Signature:   "fiber_yield(value: any) -> any",
		Description: "Yields control to the fiber scheduler with an optional value.",
		ParameterTypes: []BuiltInParameter{
			{Name: "value", Type: &ConcreteType{name: "any"}, Description: "The value to yield"},
		},
		ReturnType:   &ConcreteType{name: "any"},
		Category:     CategoryFiber,
		IsProtected:  false,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateYieldCall,
		Example:      `let result = fiber_yield(42)`,
	}

	// fiber_await function
	r.functions["fiber_await"] = &BuiltInFunction{
		Name:        "fiber_await",
		CName:       "fiber_await",
		Signature:   "fiber_await(fiber: Fiber) -> any",
		Description: "Waits for a fiber to complete and returns its result.",
		ParameterTypes: []BuiltInParameter{
			{Name: "fiber", Type: &ConcreteType{name: "Fiber"}, Description: "The fiber to await"},
		},
		ReturnType:   &ConcreteType{name: "any"},
		Category:     CategoryFiber,
		IsProtected:  false,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateAwaitCall,
		Example:      `let result = fiber_await(fiberHandle)`,
	}

	// Channel function
	r.functions["Channel"] = &BuiltInFunction{
		Name:        "Channel",
		CName:       "channel_create",
		Signature:   "Channel(capacity: int) -> Channel",
		Description: "Creates a new channel with the specified capacity.",
		ParameterTypes: []BuiltInParameter{
			{Name: "capacity", Type: &ConcreteType{name: TypeInt}, Description: "The capacity of the channel"},
		},
		ReturnType:   &ConcreteType{name: "Channel"},
		Category:     CategoryFiber,
		IsProtected:  false,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateChannelCreateCall,
		Example:      `let ch = Channel(10)`,
	}

	// send function
	r.functions["send"] = &BuiltInFunction{
		Name:        "send",
		CName:       "channel_send",
		Signature:   "send(channel: Channel, value: any) -> int",
		Description: "Sends a value to a channel. Returns 1 for success, 0 for failure.",
		ParameterTypes: []BuiltInParameter{
			{Name: "channel", Type: &ConcreteType{name: "Channel"}, Description: "The channel to send to"},
			{Name: "value", Type: &ConcreteType{name: "any"}, Description: "The value to send"},
		},
		ReturnType:   &ConcreteType{name: TypeInt},
		Category:     CategoryFiber,
		IsProtected:  false,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateChannelSendCall,
		Example:      `let success = send(ch, 42)`,
	}

	// recv function
	r.functions["recv"] = &BuiltInFunction{
		Name:        "recv",
		CName:       "channel_recv",
		Signature:   "recv(channel: Channel) -> any",
		Description: "Receives a value from a channel.",
		ParameterTypes: []BuiltInParameter{
			{Name: "channel", Type: &ConcreteType{name: "Channel"}, Description: "The channel to receive from"},
		},
		ReturnType:   &ConcreteType{name: "any"},
		Category:     CategoryFiber,
		IsProtected:  false,
		SecurityFlag: PermissionNone,
		Generator:    (*LLVMGenerator).generateChannelRecvCall,
		Example:      `let value = recv(ch)`,
	}
}

// GlobalBuiltInRegistry is the global instance of the built-in function registry
//
//nolint:gochecknoglobals // This is a necessary global registry for built-in functions
var GlobalBuiltInRegistry *BuiltInFunctionRegistry

// init initializes the global registry
func init() {
	GlobalBuiltInRegistry = NewBuiltInFunctionRegistry()
}
