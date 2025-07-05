// Package ast provides validation rules for the Osprey AST.
package ast

import (
	"fmt"
)

// ValidationError represents a validation error with optional position information.
type ValidationError struct {
	Message  string
	Position *Position
}

// Error returns the formatted error message with position information.
func (e *ValidationError) Error() string {
	if e.Position != nil {
		return fmt.Sprintf("line %d:%d: validation error: %s", e.Position.Line, e.Position.Column, e.Message)
	}
	return e.Message
}

// ValidateProgram validates the entire program AST and returns any validation errors.
func ValidateProgram(program *Program) error {
	for _, stmt := range program.Statements {
		if err := validateStatement(stmt); err != nil {
			return err
		}
	}

	return nil
}

// validateStatement validates a single statement.
func validateStatement(stmt Statement) error {
	switch s := stmt.(type) {
	case *FunctionDeclaration:
		return validateFunctionDeclaration(s)
	default:
		return nil
	}
}

// validateFunctionDeclaration validates function declarations according to type inference rules.
func validateFunctionDeclaration(fn *FunctionDeclaration) error {
	// Rule: Functions without explicit return types must have bodies that allow type inference
	if fn.ReturnType == nil {
		if !canInferReturnType(fn.Body) {
			return &ValidationError{
				Message: fmt.Sprintf("Function '%s' requires explicit return type annotation - "+
					"type cannot be inferred from body", fn.Name),
				Position: fn.Position,
			}
		}

		// CRITICAL RULE: Functions cannot implicitly return 'any' type
		if wouldInferAsAny(fn.Body, fn) {
			return &ValidationError{
				Message: fmt.Sprintf("Function '%s' cannot implicitly return 'any' type - "+
					"if 'any' return type is intended, declare it explicitly with '-> any'", fn.Name),
				Position: fn.Position,
			}
		}
	}

	// Rule: Parameters without explicit types must be used in a way that allows type inference
	for _, param := range fn.Parameters {
		if param.Type == nil {
			if !canInferParameterType(param.Name, fn.Body, fn.ReturnType) {
				return &ValidationError{
					Message: fmt.Sprintf("Parameter '%s' in function '%s' requires explicit type annotation - "+
						"type cannot be inferred from usage", param.Name, fn.Name),
					Position: fn.Position,
				}
			}
		}
	}

	return nil
}

// canInferReturnType checks if the return type can be inferred from the function body.
func canInferReturnType(body Expression) bool {
	switch e := body.(type) {
	case *IntegerLiteral, *StringLiteral, *BooleanLiteral:
		return true
	case *BinaryExpression:
		return canInferFromBinaryExpression(e)
	case *ResultExpression:
		return canInferFromResultExpression(e)
	case *CallExpression:
		return canInferFromCallExpression(e)
	case *Identifier:
		return canInferFromIdentifier(e)
	default:
		return false
	}
}

// canInferFromBinaryExpression checks if type can be inferred from binary expressions.
func canInferFromBinaryExpression(e *BinaryExpression) bool {
	return isArithmeticOperator(e.Operator) || isComparisonOperator(e.Operator)
}

// canInferFromResultExpression checks if type can be inferred from result expressions.
func canInferFromResultExpression(e *ResultExpression) bool {
	if !e.Success {
		return false
	}

	binExpr, ok := e.Value.(*BinaryExpression)
	if !ok {
		return false
	}

	return isArithmeticOperator(binExpr.Operator) || isComparisonOperator(binExpr.Operator)
}

// canInferFromCallExpression checks if type can be inferred from call expressions.
func canInferFromCallExpression(_ *CallExpression) bool {
	// For now, we can't reliably infer from function calls without more context
	return false
}

// canInferFromIdentifier checks if type can be inferred from identifiers.
func canInferFromIdentifier(_ *Identifier) bool {
	// For now, we can't reliably infer from identifiers without more context
	return false
}

// wouldInferAsAny checks if a function body would infer as 'any' type.
func wouldInferAsAny(body Expression, fn *FunctionDeclaration) bool {
	switch e := body.(type) {
	case *Identifier:
		// Direct return of a parameter - check if parameter has explicit type
		for _, param := range fn.Parameters {
			if param.Name == e.Name && param.Type != nil {
				// Parameter has explicit type, so return type can be inferred
				return false
			}
		}
		// Parameter without explicit type would be 'any'
		return true
	case *BinaryExpression:
		// Arithmetic expressions with explicitly typed parameters can be inferred
		if isArithmeticOperator(e.Operator) {
			return !allOperandsHaveExplicitTypes(e, fn)
		}

		return false
	case *CallExpression:
		// Function calls with unknown return types would be 'any'
		return true
	case *ResultExpression:
		// For result expressions, check the wrapped value
		if e.Success {
			return wouldInferAsAny(e.Value, fn)
		}

		return true
	default:
		return false
	}
}

// canInferParameterType checks if a parameter's type can be inferred from its usage.
func canInferParameterType(paramName string, body Expression, returnType *TypeExpression) bool {
	usage := analyzeParameterUsage(paramName, body)

	// If we found concrete usage patterns, we can infer
	if usage.usedInArithmetic || usage.usedAsLiteral {
		return true
	}

	// If parameter is directly returned and return type is specified, we can infer
	if usage.directlyReturned && returnType != nil {
		return true
	}

	return false
}

// ParameterUsage tracks how a parameter is used in a function body.
type ParameterUsage struct {
	usedInArithmetic bool
	usedAsLiteral    bool
	directlyReturned bool
}

// analyzeParameterUsage analyzes how a parameter is used in the function body.
func analyzeParameterUsage(paramName string, body Expression) ParameterUsage {
	usage := ParameterUsage{}

	switch e := body.(type) {
	case *Identifier:
		if e.Name == paramName {
			usage.directlyReturned = true
		}
	case *BinaryExpression:
		if containsParameterInArithmetic(paramName, e) {
			usage.usedInArithmetic = true
		}
	case *ResultExpression:
		// Check if parameter is used in arithmetic within ResultExpression
		if e.Success {
			if binExpr, ok := e.Value.(*BinaryExpression); ok {
				if containsParameterInArithmetic(paramName, binExpr) {
					usage.usedInArithmetic = true
				}
			}
		}
	}

	return usage
}

// containsParameterInArithmetic checks if a parameter is used in arithmetic operations.
func containsParameterInArithmetic(paramName string, expr *BinaryExpression) bool {
	if !isArithmeticOperator(expr.Operator) && !isComparisonOperator(expr.Operator) {
		return false
	}

	return containsParameter(paramName, expr.Left) || containsParameter(paramName, expr.Right)
}

// containsParameter recursively checks if an expression contains a parameter reference.
func containsParameter(paramName string, expr Expression) bool {
	switch e := expr.(type) {
	case *Identifier:
		return e.Name == paramName
	case *BinaryExpression:
		return containsParameter(paramName, e.Left) || containsParameter(paramName, e.Right)
	case *CallExpression:
		for _, arg := range e.Arguments {
			if containsParameter(paramName, arg) {
				return true
			}
		}
	}

	return false
}

// isArithmeticOperator checks if an operator is arithmetic.
func isArithmeticOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/" || op == "%"
}

// isComparisonOperator checks if an operator is a comparison operator.
func isComparisonOperator(op string) bool {
	return op == "==" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">="
}

// allOperandsHaveExplicitTypes checks if all operands in a binary expression have explicit types.
func allOperandsHaveExplicitTypes(expr *BinaryExpression, fn *FunctionDeclaration) bool {
	return operandHasExplicitType(expr.Left, fn) && operandHasExplicitType(expr.Right, fn)
}

// operandHasExplicitType checks if an operand has an explicit type.
func operandHasExplicitType(operand Expression, fn *FunctionDeclaration) bool {
	switch e := operand.(type) {
	case *IntegerLiteral, *StringLiteral, *BooleanLiteral:
		// Literals have explicit types
		return true
	case *Identifier:
		// Check if this identifier is a parameter with explicit type
		for _, param := range fn.Parameters {
			if param.Name == e.Name {
				// Parameter has explicit type annotation
				if param.Type != nil {
					return true
				}
				// Parameter can be inferred from usage in arithmetic
				if canInferParameterType(param.Name, fn.Body, fn.ReturnType) {
					return true
				}
			}
		}

		return false
	case *BinaryExpression:
		// Nested binary expressions - check recursively if it's an arithmetic operation
		if isArithmeticOperator(e.Operator) {
			return allOperandsHaveExplicitTypes(e, fn)
		}

		return false
	case *ResultExpression:
		// For result expressions, check the wrapped value
		if e.Success {
			return operandHasExplicitType(e.Value, fn)
		}

		return false
	default:
		return false
	}
}
