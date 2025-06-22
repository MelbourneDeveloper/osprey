package ast

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
)

func TestValidateProgram(t *testing.T) {
	t.Run("empty_program", func(t *testing.T) {
		program := &ast.Program{
			Statements: []ast.Statement{},
		}

		err := ast.ValidateProgram(program)
		if err != nil {
			t.Errorf("Empty program should validate successfully, got: %v", err)
		}
	})

	t.Run("program_with_valid_function", func(t *testing.T) {
		program := &ast.Program{
			Statements: []ast.Statement{
				&ast.FunctionDeclaration{
					Name: "add",
					Parameters: []ast.Parameter{
						{Name: "x", Type: &ast.TypeExpression{Name: "int"}},
						{Name: "y", Type: &ast.TypeExpression{Name: "int"}},
					},
					ReturnType: &ast.TypeExpression{Name: "int"},
					Body: &ast.BinaryExpression{
						Left:     &ast.Identifier{Name: "x"},
						Operator: "+",
						Right:    &ast.Identifier{Name: "y"},
					},
				},
			},
		}

		err := ast.ValidateProgram(program)
		if err != nil {
			t.Errorf("Valid function should validate successfully, got: %v", err)
		}
	})

	t.Run("program_with_invalid_function", func(t *testing.T) {
		program := &ast.Program{
			Statements: []ast.Statement{
				&ast.FunctionDeclaration{
					Name:       "mystery",
					Parameters: []ast.Parameter{{Name: "x", Type: nil}},
					ReturnType: nil,
					Body:       &ast.Identifier{Name: "unknown"},
				},
			},
		}

		err := ast.ValidateProgram(program)
		if err == nil {
			t.Error("Invalid function should fail validation")
		}
	})
}

func TestValidateFunctionDeclaration(t *testing.T) {
	t.Run("function_with_explicit_types", func(t *testing.T) {
		fn := &ast.FunctionDeclaration{
			Name: "multiply",
			Parameters: []ast.Parameter{
				{Name: "a", Type: &ast.TypeExpression{Name: "int"}},
				{Name: "b", Type: &ast.TypeExpression{Name: "int"}},
			},
			ReturnType: &ast.TypeExpression{Name: "int"},
			Body: &ast.BinaryExpression{
				Left:     &ast.Identifier{Name: "a"},
				Operator: "*",
				Right:    &ast.Identifier{Name: "b"},
			},
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err != nil {
			t.Errorf("Function with explicit types should validate, got: %v", err)
		}
	})

	t.Run("function_without_return_type_inferrable", func(t *testing.T) {
		fn := &ast.FunctionDeclaration{
			Name: "double",
			Parameters: []ast.Parameter{
				{Name: "x", Type: &ast.TypeExpression{Name: "int"}},
			},
			ReturnType: nil, // No explicit return type
			Body: &ast.BinaryExpression{
				Left:     &ast.Identifier{Name: "x"},
				Operator: "*",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err != nil {
			t.Errorf("Function with inferrable return type should validate, got: %v", err)
		}
	})

	t.Run("function_without_return_type_not_inferrable", func(t *testing.T) {
		fn := &ast.FunctionDeclaration{
			Name:       "mystery",
			Parameters: []ast.Parameter{{Name: "x", Type: &ast.TypeExpression{Name: "int"}}},
			ReturnType: nil,
			Body:       &ast.Identifier{Name: "unknown_var"},
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err == nil {
			t.Error("Function with non-inferrable return type should fail validation")
		}

		if !containsText(err.Error(), "requires explicit return type annotation") {
			t.Errorf("Error should mention return type annotation, got: %v", err)
		}
	})

	t.Run("parameter_without_type_inferrable", func(t *testing.T) {
		fn := &ast.FunctionDeclaration{
			Name: "square",
			Parameters: []ast.Parameter{
				{Name: "num", Type: nil}, // No explicit type
			},
			ReturnType: &ast.TypeExpression{Name: "int"},
			Body: &ast.BinaryExpression{
				Left:     &ast.Identifier{Name: "num"},
				Operator: "*",
				Right:    &ast.Identifier{Name: "num"},
			},
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err != nil {
			t.Errorf("Function with inferrable parameter type should validate, got: %v", err)
		}
	})

	t.Run("parameter_without_type_not_inferrable", func(t *testing.T) {
		fn := &ast.FunctionDeclaration{
			Name: "unknown_func",
			Parameters: []ast.Parameter{
				{Name: "param", Type: nil}, // No explicit type
			},
			ReturnType: &ast.TypeExpression{Name: "int"},
			Body:       &ast.IntegerLiteral{Value: 42}, // Parameter not used
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err == nil {
			t.Error("Function with non-inferrable parameter type should fail validation")
		}

		if !containsText(err.Error(), "requires explicit type annotation") {
			t.Errorf("Error should mention type annotation, got: %v", err)
		}
	})
}

func TestCanInferReturnType(t *testing.T) {
	tests := []struct {
		name     string
		body     ast.Expression
		expected bool
	}{
		{
			name:     "integer_literal",
			body:     &ast.IntegerLiteral{Value: 42},
			expected: true,
		},
		{
			name:     "string_literal",
			body:     &ast.StringLiteral{Value: "hello"},
			expected: true,
		},
		{
			name:     "boolean_literal",
			body:     &ast.BooleanLiteral{Value: true},
			expected: true,
		},
		{
			name: "arithmetic_expression",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: true,
		},
		{
			name: "non_arithmetic_expression",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: "==",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			expected: true,
		},
		{
			name: "successful_result_expression",
			body: &ast.ResultExpression{
				Success: true,
				Value: &ast.BinaryExpression{
					Left:     &ast.IntegerLiteral{Value: 5},
					Operator: "*",
					Right:    &ast.IntegerLiteral{Value: 3},
				},
			},
			expected: true,
		},
		{
			name: "failed_result_expression",
			body: &ast.ResultExpression{
				Success: false,
				Value:   &ast.StringLiteral{Value: "error"},
			},
			expected: false,
		},
		{
			name:     "call_expression",
			body:     &ast.CallExpression{Function: &ast.Identifier{Name: "unknown"}},
			expected: false,
		},
		{
			name:     "identifier",
			body:     &ast.Identifier{Name: "param"},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a dummy function to test return type inference
			fn := &ast.FunctionDeclaration{
				Name:       "test",
				Parameters: []ast.Parameter{},
				ReturnType: nil,
				Body:       test.body,
			}

			err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
			canInfer := err == nil

			if canInfer != test.expected {
				t.Errorf("Expected can infer return type = %v, got %v (error: %v)",
					test.expected, canInfer, err)
			}
		})
	}
}

func TestCanInferParameterType(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		body       ast.Expression
		returnType *ast.TypeExpression
		expected   bool
	}{
		{
			name:  "parameter_used_in_arithmetic",
			param: "x",
			body: &ast.BinaryExpression{
				Left:     &ast.Identifier{Name: "x"},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Value: 1},
			},
			returnType: nil,
			expected:   true,
		},
		{
			name:       "parameter_directly_returned_with_type",
			param:      "value",
			body:       &ast.Identifier{Name: "value"},
			returnType: &ast.TypeExpression{Name: "int"},
			expected:   true,
		},
		{
			name:       "parameter_directly_returned_without_type",
			param:      "value",
			body:       &ast.Identifier{Name: "value"},
			returnType: nil,
			expected:   false,
		},
		{
			name:       "parameter_not_used",
			param:      "unused",
			body:       &ast.IntegerLiteral{Value: 42},
			returnType: &ast.TypeExpression{Name: "int"},
			expected:   false,
		},
		{
			name:  "parameter_in_result_expression",
			param: "num",
			body: &ast.ResultExpression{
				Success: true,
				Value: &ast.BinaryExpression{
					Left:     &ast.Identifier{Name: "num"},
					Operator: "*",
					Right:    &ast.IntegerLiteral{Value: 2},
				},
			},
			returnType: nil,
			expected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn := &ast.FunctionDeclaration{
				Name: "test",
				Parameters: []ast.Parameter{
					{Name: test.param, Type: nil}, // No explicit type
				},
				ReturnType: test.returnType,
				Body:       test.body,
			}

			err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
			canInfer := err == nil

			if canInfer != test.expected {
				t.Errorf("Expected can infer parameter type = %v, got %v (error: %v)",
					test.expected, canInfer, err)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := &ast.ValidationError{Message: "test error"}

	if err.Error() != "test error" {
		t.Errorf("ValidationError.Error() should return message, got: %v", err.Error())
	}
}

func TestArithmeticOperators(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		body     ast.Expression
		valid    bool
	}{
		{
			name:     "addition",
			operator: "+",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			valid: true,
		},
		{
			name:     "subtraction",
			operator: "-",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 5},
				Operator: "-",
				Right:    &ast.IntegerLiteral{Value: 3},
			},
			valid: true,
		},
		{
			name:     "multiplication",
			operator: "*",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 4},
				Operator: "*",
				Right:    &ast.IntegerLiteral{Value: 3},
			},
			valid: true,
		},
		{
			name:     "division",
			operator: "/",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 8},
				Operator: "/",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
			valid: true,
		},
		{
			name:     "comparison",
			operator: "==",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: "==",
				Right:    &ast.IntegerLiteral{Value: 1},
			},
			valid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fn := &ast.FunctionDeclaration{
				Name:       "test",
				Parameters: []ast.Parameter{},
				ReturnType: nil, // Should be inferrable for arithmetic, not for comparison
				Body:       test.body,
			}

			err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
			isValid := err == nil

			if isValid != test.valid {
				t.Errorf("Expected operator %s validity = %v, got %v (error: %v)",
					test.operator, test.valid, isValid, err)
			}
		})
	}
}

func containsText(text, substring string) bool {
	return len(text) > 0 && len(substring) > 0 &&
		len(text) >= len(substring) &&
		indexOf(text, substring) >= 0
}

func indexOf(text, substring string) int {
	for i := 0; i <= len(text)-len(substring); i++ {
		if text[i:i+len(substring)] == substring {
			return i
		}
	}

	return -1
}
