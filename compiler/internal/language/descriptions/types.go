package descriptions

// BuiltinTypeDesc represents documentation for a built-in type.
type BuiltinTypeDesc struct {
	Name        string
	Description string
	Example     string
}

// GetBuiltinTypeDescriptions returns all built-in type descriptions.
func GetBuiltinTypeDescriptions() map[string]*BuiltinTypeDesc {
	return map[string]*BuiltinTypeDesc{
		"Int": {
			Name: "Int",
			Description: "A 64-bit signed integer type. Can represent whole numbers from " +
				"-9,223,372,036,854,775,808 to 9,223,372,036,854,775,807.",
			Example: `let number: Int = 42\nlet negative: Int = -100`,
		},
		"String": {
			Name: "String",
			Description: "A sequence of characters representing text. Supports string " +
				"interpolation and escape sequences.",
			Example: `let greeting: String = "Hello, World!"\nlet name = "Alice"\nlet message = "Hello, ${name}!"`,
		},
		"Bool": {
			Name: "Bool",
			Description: "A boolean type that can be either true or false. Used for " +
				"logical operations and conditionals.",
			Example: `let isValid: Bool = true\nlet isComplete: Bool = false`,
		},
		"Any": {
			Name: "Any",
			Description: "A type that can represent any value. Useful for generic " +
				"programming but should be used carefully as it bypasses type checking.",
			Example: `let value: Any = 42\nlet text: Any = "Hello"`,
		},
		"HttpResponse": {
			Name: "HttpResponse",
			Description: "A built-in type representing an HTTP response with status code, " +
				"headers, content type, body, and streaming capabilities. Used by HTTP server " +
				"handlers to return structured responses to clients.",
			Example: `HttpResponse {` + "\n" +
				`    status: 200,` + "\n" +
				`    headers: "Content-Type: application/json",` + "\n" +
				`    contentType: "application/json",` + "\n" +
				`    streamFd: -1,` + "\n" +
				`    isComplete: true,` + "\n" +
				`    partialBody: "{\"message\": \"Hello\"}"` + "\n" +
				`}`,
		},
		"ProcessHandle": {
			Name: "ProcessHandle",
			Description: "A handle to a spawned async process. Contains the process ID and allows " +
				"waiting for completion and cleanup. Process output is delivered via callbacks " +
				"registered with the runtime.",
			Example: `let result = spawnProcess("echo hello")` + "\n" +
				`match result {` + "\n" +
				`    Success { value } => {` + "\n" +
				`        let exitCode = awaitProcess(value)` + "\n" +
				`        cleanupProcess(value)` + "\n" +
				`    }` + "\n" +
				`    Error { message } => print("Process failed")` + "\n" +
				`}`,
		},
	}
}

// GetBuiltinTypeDescription returns description for a single built-in type.
func GetBuiltinTypeDescription(name string) *BuiltinTypeDesc {
	descriptions := GetBuiltinTypeDescriptions()
	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return nil
}
