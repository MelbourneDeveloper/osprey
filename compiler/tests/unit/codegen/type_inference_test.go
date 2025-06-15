package codegen_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
)

const (
	stringType    = "string"
	intType       = "int"
	testFuncName  = "testFunc"
	testValue     = "test"
	testAge       = 25
	testScore     = 42
	maxFuncLength = 60
)

// TestFunctionReturnTypeInference tests the analyzeReturnType method.
func TestFunctionReturnTypeInference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		expr     ast.Expression
		expected string
	}{
		{
			name:     "string literal",
			expr:     &ast.StringLiteral{Value: "hello"},
			expected: stringType,
		},
		{
			name:     "integer literal",
			expr:     &ast.IntegerLiteral{Value: testScore},
			expected: intType,
		},
		{
			name:     "identifier with name keyword",
			expr:     &ast.Identifier{Name: "userName"},
			expected: stringType,
		},
		{
			name:     "identifier with text keyword",
			expr:     &ast.Identifier{Name: "inputText"},
			expected: stringType,
		},
		{
			name:     "identifier with str keyword",
			expr:     &ast.Identifier{Name: "myStr"},
			expected: stringType,
		},
		{
			name:     "regular identifier",
			expr:     &ast.Identifier{Name: "value"},
			expected: stringType, // Default to string for identifiers
		},
		{
			name: "binary expression",
			expr: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 5},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Value: 3},
			},
			expected: intType,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			generator := codegen.NewLLVMGenerator()

			result := generator.AnalyzeReturnType(testCase.expr)
			if result != testCase.expected {
				t.Errorf("Expected %s, got %s", testCase.expected, result)
			}
		})
	}
}

// TestParameterTypeInference tests parameter type inference during function signature analysis.
func TestParameterTypeInference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		paramName    string
		returnType   string
		numParams    int
		expectedType string
	}{
		{
			name:         "regular parameter",
			paramName:    "value",
			returnType:   intType,
			numParams:    2,
			expectedType: intType,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runParameterTypeTest(t, testCase)
		})
	}
}

func runParameterTypeTest(t *testing.T, testCase struct {
	name         string
	paramName    string
	returnType   string
	numParams    int
	expectedType string
},
) {
	t.Helper()

	generator := codegen.NewLLVMGenerator()

	// Create a mock function declaration with body matching expected return type
	var body ast.Expression
	if testCase.returnType == stringType {
		body = &ast.StringLiteral{Value: testValue}
	} else {
		body = &ast.IntegerLiteral{Value: testScore}
	}

	// Create parameters based on numParams
	params := make([]ast.Parameter, testCase.numParams)
	params[0] = ast.Parameter{Name: testCase.paramName, Type: nil}
	for i := 1; i < testCase.numParams; i++ {
		params[i] = ast.Parameter{Name: fmt.Sprintf("param%d", i), Type: nil}
	}

	fnDecl := &ast.FunctionDeclaration{
		Name:       testFuncName,
		Parameters: params,
		Body:       body,
		ReturnType: nil,
	}

	// Test the parameter type inference logic
	err := generator.DeclareFunctionSignature(fnDecl)
	if err != nil {
		t.Fatalf("Failed to declare function signature: %v", err)
	}

	validateParameterType(t, generator, testCase.expectedType)
}

func validateParameterType(t *testing.T, generator *codegen.LLVMGenerator, expectedType string) {
	t.Helper()

	// Check if the function was created with the right parameter type
	function := generator.GetFunction(testFuncName)
	if function == nil {
		t.Fatal("Function not found after declaration")
		return // This helps staticcheck understand that function is not nil after this point
	}

	if len(function.Params) == 0 {
		t.Fatal("Function has no parameters")
		return
	}

	paramType := function.Params[0].Type()

	expectedLLVMType := "i64" // Default to int
	if expectedType == stringType {
		expectedLLVMType = "i8"
	}

	if !strings.Contains(paramType.String(), expectedLLVMType) {
		t.Errorf("Expected parameter type to contain %s, got %s", expectedLLVMType, paramType.String())
	}
}

// TestComplexTypeInference tests complex type inference scenarios.
func TestComplexTypeInference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		functions    []ast.FunctionDeclaration
		testFunction string
		expectedType string
	}{
		{
			name:         "function calling another function",
			functions:    createStringFunctionChain(),
			testFunction: "processString",
			expectedType: stringType,
		},
		{
			name:         "arithmetic function",
			functions:    createArithmeticFunction(),
			testFunction: "calculate",
			expectedType: intType,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runComplexTypeTest(t, testCase)
		})
	}
}

func createStringFunctionChain() []ast.FunctionDeclaration {
	return []ast.FunctionDeclaration{
		{
			Name:       "getString",
			Parameters: []ast.Parameter{},
			Body:       &ast.StringLiteral{Value: "hello"},
			ReturnType: nil,
		},
		{
			Name:       "processString",
			Parameters: []ast.Parameter{{Name: "text", Type: nil}},
			Body: &ast.CallExpression{
				Function:       &ast.Identifier{Name: "getString"},
				Arguments:      []ast.Expression{},
				NamedArguments: []ast.NamedArgument{},
			},
			ReturnType: nil,
		},
	}
}

func createArithmeticFunction() []ast.FunctionDeclaration {
	return []ast.FunctionDeclaration{
		{
			Name:       "calculate",
			Parameters: []ast.Parameter{{Name: "x", Type: nil}, {Name: "y", Type: nil}},
			Body: &ast.BinaryExpression{
				Left:     &ast.Identifier{Name: "x"},
				Operator: "+",
				Right:    &ast.Identifier{Name: "y"},
			},
			ReturnType: nil,
		},
	}
}

func runComplexTypeTest(t *testing.T, testCase struct {
	name         string
	functions    []ast.FunctionDeclaration
	testFunction string
	expectedType string
},
) {
	t.Helper()

	generator := codegen.NewLLVMGenerator()

	// Declare all functions
	for _, function := range testCase.functions {
		err := generator.DeclareFunctionSignature(&function)
		if err != nil {
			t.Fatalf("Failed to declare function %s: %v", function.Name, err)
		}
	}

	// Check the return type of the test function
	returnType := generator.GetFunctionReturnType(testCase.testFunction)
	if returnType != testCase.expectedType {
		t.Errorf("Expected %s return type for function %s, got %s",
			testCase.expectedType,
			testCase.testFunction,
			returnType)
	}
}

// TestInterpolatedStringTypeInference tests type inference in interpolated strings.
func TestInterpolatedStringTypeInference(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		functions   []ast.FunctionDeclaration
		variables   map[string]string
		expressions []ast.Expression
		expectedFmt []string // Expected format specifiers
	}{
		{
			name:        "mixed types in interpolation",
			functions:   createMixedTypeFunctions(),
			variables:   map[string]string{"score": intType},
			expressions: createMixedTypeExpressions(),
			expectedFmt: []string{"%s", "%ld", "%ld"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // capture range variable
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			runInterpolatedStringTest(t, testCase)
		})
	}
}

func createMixedTypeFunctions() []ast.FunctionDeclaration {
	return []ast.FunctionDeclaration{
		{
			Name:       "getName",
			Parameters: []ast.Parameter{},
			Body:       &ast.StringLiteral{Value: "Alice"},
			ReturnType: nil,
		},
		{
			Name:       "getAge",
			Parameters: []ast.Parameter{},
			Body:       &ast.IntegerLiteral{Value: testAge},
			ReturnType: nil,
		},
	}
}

func createMixedTypeExpressions() []ast.Expression {
	return []ast.Expression{
		&ast.CallExpression{
			Function:       &ast.Identifier{Name: "getName"},
			Arguments:      []ast.Expression{},
			NamedArguments: []ast.NamedArgument{},
		},
		&ast.CallExpression{
			Function:       &ast.Identifier{Name: "getAge"},
			Arguments:      []ast.Expression{},
			NamedArguments: []ast.NamedArgument{},
		},
		&ast.Identifier{Name: "score"},
	}
}

func runInterpolatedStringTest(t *testing.T, testCase struct {
	name        string
	functions   []ast.FunctionDeclaration
	variables   map[string]string
	expressions []ast.Expression
	expectedFmt []string
},
) {
	t.Helper()

	generator := codegen.NewLLVMGenerator()

	// Declare functions
	for _, function := range testCase.functions {
		err := generator.DeclareFunctionSignature(&function)
		if err != nil {
			t.Fatalf("Failed to declare function %s: %v", function.Name, err)
		}
	}

	// Set up variables
	for name, varType := range testCase.variables {
		generator.SetVariableType(name, varType)
	}

	// Test each expression's format specifier
	expressions := testCase.expressions
	expectedFmt := testCase.expectedFmt
	validateExpressionTypes(t, generator, expressions, expectedFmt)
}

func validateExpressionTypes(
	t *testing.T,
	generator *codegen.LLVMGenerator,
	expressions []ast.Expression,
	expectedFmt []string,
) {
	t.Helper()

	for index, expr := range expressions {
		expectedFormat := expectedFmt[index]
		validateSingleExpression(t, generator, expr, expectedFormat, index)
	}
}

func validateSingleExpression(
	t *testing.T,
	generator *codegen.LLVMGenerator,
	expr ast.Expression,
	expectedFormat string,
	index int,
) {
	t.Helper()

	callExpr, isCallExpr := expr.(*ast.CallExpression)
	if !isCallExpr {
		return
	}

	ident, isIdent := callExpr.Function.(*ast.Identifier)
	if !isIdent {
		return
	}

	returnType := generator.GetFunctionReturnType(ident.Name)
	expectedReturnType := determineExpectedReturnType(expectedFormat)

	if returnType != expectedReturnType {
		t.Errorf("Expression %d: expected %s return type, got %s",
			index, expectedReturnType, returnType)
	}
}

func determineExpectedReturnType(expectedFormat string) string {
	if expectedFormat == "%ld" {
		return intType
	}

	return stringType
}
