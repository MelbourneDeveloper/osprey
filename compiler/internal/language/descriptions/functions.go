package descriptions

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

// GetBuiltinFunctionDescriptions returns all built-in function descriptions.
func GetBuiltinFunctionDescriptions() map[string]*BuiltinFunctionDesc {
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
			Signature:   "substring(s: string, start: int, end: int) -> string",
			Description: "Extracts a substring from start to end index.",
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
			ReturnType: "string",
			Example:    `let sub = substring("hello", 1, 4)\nprint(sub)  // Prints: ell`,
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
