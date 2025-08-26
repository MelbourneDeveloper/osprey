package descriptions

// KeywordDesc represents documentation for a language keyword.
type KeywordDesc struct {
	Keyword     string
	Description string
	Example     string
}

// GetKeywordDescriptions returns all keyword descriptions.
func GetKeywordDescriptions() map[string]*KeywordDesc {
	return map[string]*KeywordDesc{
		"fn": {
			Keyword: "fn",
			Description: "Function declaration keyword. Used to define functions with " +
				"parameters and return types.",
			Example: `fn add(a: Int, b: Int) -> Int {\n    a + b\n}\n\nfn greet(name: String) {\n    print("Hello, " + name)\n}`,
		},
		"let": {
			Keyword: "let",
			Description: "Variable declaration keyword. Used to bind values to identifiers. " +
				"Variables are immutable by default in Osprey.",
			Example: `let name = "Alice"\nlet age: Int = 25\nlet isActive = true`,
		},
		"type": {
			Keyword:     "type",
			Description: "Type declaration keyword. Used to define custom types and type aliases.",
			Example:     `type UserId = Int\ntype Status = Active | Inactive\ntype User = { name: String, age: Int }`,
		},
		"match": {
			Keyword: "match",
			Description: "Pattern matching expression. Used for destructuring values and " +
				"control flow based on patterns.",
			Example: `match value {\n    Some(x) -> x\n    None -> 0\n}\n\n` +
				`match status {\n    Active -> "User is active"\n    ` +
				`Inactive -> "User is inactive"\n}`,
		},
		"import": {
			Keyword: "import",
			Description: "Import declaration keyword. Used to bring modules and their " +
				"exports into the current scope.",
			Example: `import { function1, Type1 } from "module"\nimport * as Utils from "utils"`,
		},
		"true": {
			Keyword:     "true",
			Description: "Boolean literal representing the logical value true.",
			Example:     `let isReady = true\nif (isReady) { print("Ready!") }`,
		},
		"false": {
			Keyword:     "false",
			Description: "Boolean literal representing the logical value false.",
			Example:     `let isComplete = false\nif (!isComplete) { print("Not done yet") }`,
		},
	}
}

// GetKeywordDescription returns description for a single keyword.
func GetKeywordDescription(keyword string) *KeywordDesc {
	descriptions := GetKeywordDescriptions()
	if desc, exists := descriptions[keyword]; exists {
		return desc
	}

	return nil
}
