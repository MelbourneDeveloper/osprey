package descriptions

// OperatorDesc represents documentation for an operator.
type OperatorDesc struct {
	Symbol      string
	Name        string
	Description string
	Example     string
}

// GetOperatorDescriptions returns all operator descriptions.
func GetOperatorDescriptions() map[string]*OperatorDesc {
	return map[string]*OperatorDesc{
		"+": {
			Symbol:      "+",
			Name:        "Addition",
			Description: "Adds two numbers together.",
			Example:     `let result = 5 + 3  // result = 8`,
		},
		"-": {
			Symbol:      "-",
			Name:        "Subtraction",
			Description: "Subtracts the second number from the first.",
			Example:     `let result = 10 - 4  // result = 6`,
		},
		"*": {
			Symbol:      "*",
			Name:        "Multiplication",
			Description: "Multiplies two numbers.",
			Example:     `let result = 6 * 7  // result = 42`,
		},
		"/": {
			Symbol:      "/",
			Name:        "Division",
			Description: "Divides the first number by the second.",
			Example:     `let result = 15 / 3  // result = 5`,
		},
		"%": {
			Symbol:      "%",
			Name:        "Modulo",
			Description: "Returns the remainder of dividing the first number by the second.",
			Example:     `let result = 17 % 5  // result = 2`,
		},
		"==": {
			Symbol:      "==",
			Name:        "Equality",
			Description: "Compares two values for equality.",
			Example:     `let isEqual = 5 == 5  // isEqual = true`,
		},
		"!=": {
			Symbol:      "!=",
			Name:        "Inequality",
			Description: "Compares two values for inequality.",
			Example:     `let isNotEqual = 5 != 3  // isNotEqual = true`,
		},
		"<": {
			Symbol:      "<",
			Name:        "Less Than",
			Description: "Checks if the first value is less than the second.",
			Example:     `let isLess = 3 < 5  // isLess = true`,
		},
		"<=": {
			Symbol:      "<=",
			Name:        "Less Than or Equal",
			Description: "Checks if the first value is less than or equal to the second.",
			Example:     `let isLessOrEqual = 5 <= 5  // isLessOrEqual = true`,
		},
		">": {
			Symbol:      ">",
			Name:        "Greater Than",
			Description: "Checks if the first value is greater than the second.",
			Example:     `let isGreater = 7 > 3  // isGreater = true`,
		},
		">=": {
			Symbol:      ">=",
			Name:        "Greater Than or Equal",
			Description: "Checks if the first value is greater than or equal to the second.",
			Example:     `let isGreaterOrEqual = 5 >= 5  // isGreaterOrEqual = true`,
		},
		"|>": {
			Symbol: "|>",
			Name:   "Pipe Operator",
			Description: "Takes the result of the left expression and passes it as the " +
				"first argument to the right function. Enables functional programming " +
				"and method chaining.",
			Example: `5 |> double |> print  // (5) -> double -> print\nrange(1, 10) |> forEach(print)`,
		},
	}
}

// GetOperatorDescription returns description for a single operator.
func GetOperatorDescription(symbol string) *OperatorDesc {
	descriptions := GetOperatorDescriptions()
	if desc, exists := descriptions[symbol]; exists {
		return desc
	}

	return nil
}
