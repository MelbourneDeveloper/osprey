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
		// With Hindley-Milner type inference, validation is handled by the type system
		// AST validation no longer rejects missing type annotations
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
		if err != nil {
			t.Errorf("With HM type inference, AST validation should pass, got: %v", err)
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
		// With Hindley-Milner type inference, the type system handles all inference
		// AST validation no longer requires explicit type annotations
		fn := &ast.FunctionDeclaration{
			Name:       "mystery",
			Parameters: []ast.Parameter{{Name: "x", Type: &ast.TypeExpression{Name: "int"}}},
			ReturnType: nil,
			Body:       &ast.Identifier{Name: "unknown_var"},
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err != nil {
			t.Errorf("With HM type inference, AST validation should pass, got: %v", err)
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
		// With Hindley-Milner type inference, the type system handles all inference
		// AST validation no longer requires explicit type annotations
		fn := &ast.FunctionDeclaration{
			Name: "unknown_func",
			Parameters: []ast.Parameter{
				{Name: "param", Type: nil}, // No explicit type
			},
			ReturnType: &ast.TypeExpression{Name: "int"},
			Body:       &ast.IntegerLiteral{Value: 42}, // Parameter not used
		}

		err := ast.ValidateProgram(&ast.Program{Statements: []ast.Statement{fn}})
		if err != nil {
			t.Errorf("With HM type inference, AST validation should pass, got: %v", err)
		}
	})
}

func TestCanInferReturnType(t *testing.T) {
	// With Hindley-Milner type inference, all type inference is handled by the type system
	// AST validation always passes regardless of whether types can be inferred
	tests := []struct {
		name string
		body ast.Expression
	}{
		{
			name: "integer_literal",
			body: &ast.IntegerLiteral{Value: 42},
		},
		{
			name: "string_literal",
			body: &ast.StringLiteral{Value: "hello"},
		},
		{
			name: "boolean_literal",
			body: &ast.BooleanLiteral{Value: true},
		},
		{
			name: "arithmetic_expression",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
		},
		{
			name: "non_arithmetic_expression",
			body: &ast.BinaryExpression{
				Left:     &ast.IntegerLiteral{Value: 1},
				Operator: "==",
				Right:    &ast.IntegerLiteral{Value: 2},
			},
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
		},
		{
			name: "failed_result_expression",
			body: &ast.ResultExpression{
				Success: false,
				Value:   &ast.StringLiteral{Value: "error"},
			},
		},
		{
			name: "call_expression",
			body: &ast.CallExpression{Function: &ast.Identifier{Name: "unknown"}},
		},
		{
			name: "identifier",
			body: &ast.Identifier{Name: "param"},
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
			// With HM type inference, AST validation always passes
			if err != nil {
				t.Errorf("With HM type inference, AST validation should always pass, got: %v", err)
			}
		})
	}
}

func TestCanInferParameterType(t *testing.T) {
	// With Hindley-Milner type inference, all type inference is handled by the type system
	// AST validation always passes regardless of whether types can be inferred
	tests := []struct {
		name       string
		param      string
		body       ast.Expression
		returnType *ast.TypeExpression
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
		},
		{
			name:       "parameter_directly_returned_with_type",
			param:      "value",
			body:       &ast.Identifier{Name: "value"},
			returnType: &ast.TypeExpression{Name: "int"},
		},
		{
			name:       "parameter_directly_returned_without_type",
			param:      "value",
			body:       &ast.Identifier{Name: "value"},
			returnType: nil,
		},
		{
			name:       "parameter_not_used",
			param:      "unused",
			body:       &ast.IntegerLiteral{Value: 42},
			returnType: &ast.TypeExpression{Name: "int"},
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
			// With HM type inference, AST validation always passes
			if err != nil {
				t.Errorf("With HM type inference, AST validation should always pass, got: %v", err)
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

func TestASTStatementInterfaceMethods(t *testing.T) {
	// Test that all statement types implement the isStatement() method
	statements := []ast.Statement{
		&ast.ImportStatement{Module: []string{"test"}},
		&ast.LetDeclaration{Name: "x", Value: &ast.IntegerLiteral{Value: 42}},
		&ast.AssignmentStatement{Name: "x", Value: &ast.IntegerLiteral{Value: 10}},
		&ast.FunctionDeclaration{Name: "test", Body: &ast.IntegerLiteral{Value: 1}},
		&ast.ExternDeclaration{Name: "extern_func"},
		&ast.ExpressionStatement{Expression: &ast.IntegerLiteral{Value: 5}},
		&ast.TypeDeclaration{Name: "MyType"},
		&ast.EffectDeclaration{Name: "MyEffect"},
	}
	
	for i, stmt := range statements {
		// Verify they can be used as statements in a program
		program := &ast.Program{Statements: []ast.Statement{stmt}}
		_ = program // Use the program
		
		t.Logf("Statement %d implements Statement interface", i)
	}
}

func TestASTExpressionInterfaceMethods(t *testing.T) {
	// Test that expression types implement the isExpression() method
	expressions := []ast.Expression{
		&ast.IntegerLiteral{Value: 42},
		&ast.StringLiteral{Value: "test"},
		&ast.BooleanLiteral{Value: true},
		&ast.InterpolatedStringLiteral{Parts: []ast.InterpolatedPart{}},
		&ast.Identifier{Name: "x"},
		&ast.BinaryExpression{Left: &ast.IntegerLiteral{Value: 1}, Operator: "+", Right: &ast.IntegerLiteral{Value: 2}},
		&ast.UnaryExpression{Operator: "-", Operand: &ast.IntegerLiteral{Value: 5}},
		&ast.CallExpression{Function: &ast.Identifier{Name: "func"}},
		&ast.FunctionCallExpression{Function: "test", Arguments: []ast.Expression{}},
		&ast.ResultExpression{Success: true, Value: &ast.IntegerLiteral{Value: 1}},
		&ast.FieldAccessExpression{Object: &ast.Identifier{Name: "obj"}, FieldName: "field"},
		&ast.MatchExpression{Expression: &ast.Identifier{Name: "x"}, Arms: []ast.MatchArm{}},
		&ast.LambdaExpression{Parameters: []ast.Parameter{}},
		&ast.YieldExpression{Value: &ast.IntegerLiteral{Value: 1}},
		&ast.AwaitExpression{Expression: &ast.Identifier{Name: "fiber"}},
		&ast.SpawnExpression{Expression: &ast.Identifier{Name: "task"}},
		&ast.MethodCallExpression{Object: &ast.Identifier{Name: "obj"}, MethodName: "method"},
		&ast.ModuleAccessExpression{ModuleName: "module", MemberName: "member"},
		&ast.ChannelExpression{ElementType: ast.TypeExpression{Name: "int"}},
		&ast.ChannelSendExpression{Channel: &ast.Identifier{Name: "ch"}, Value: &ast.IntegerLiteral{Value: 1}},
		&ast.ChannelRecvExpression{Channel: &ast.Identifier{Name: "ch"}},
		&ast.SelectExpression{Arms: []ast.SelectArm{}},
		&ast.ChannelCreateExpression{Capacity: &ast.IntegerLiteral{Value: 10}},
		&ast.TypeConstructorExpression{TypeName: "MyType", Fields: map[string]ast.Expression{}},
		&ast.BlockExpression{Statements: []ast.Statement{}},
		&ast.ListLiteral{Elements: []ast.Expression{}},
		&ast.ObjectLiteral{Fields: map[string]ast.Expression{}},
		&ast.MapLiteral{Entries: []ast.MapEntry{}},
		&ast.ListAccessExpression{List: &ast.Identifier{Name: "list"}, Index: &ast.IntegerLiteral{Value: 0}},
		&ast.UpdateExpression{Target: &ast.Identifier{Name: "obj"}, Fields: map[string]ast.Expression{}},
	}
	
	for i, expr := range expressions {
		// Verify they can be used as expressions
		letDecl := &ast.LetDeclaration{Name: "temp", Value: expr}
		program := &ast.Program{Statements: []ast.Statement{letDecl}}
		_ = program // Use the program
		
		t.Logf("Expression %d implements Expression interface", i)
	}
}
