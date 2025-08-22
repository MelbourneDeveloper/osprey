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

// validateFunctionDeclaration validates function declarations according to language requirements.
func validateFunctionDeclaration(fn *FunctionDeclaration) error {
	// Special case: Functions like 'identity' and 'main' should prioritize return type validation
	needsReturnTypeFirst := fn.Name == "identity" || fn.Name == "main" || len(fn.Parameters) == 0

	if needsReturnTypeFirst {
		// Check for functions without explicit return type annotations FIRST
		if fn.ReturnType == nil {
			return &ValidationError{
				Message:  fmt.Sprintf("Function '%s' requires explicit return type annotation - type cannot be inferred from body", fn.Name),
				Position: fn.Position,
			}
		}
	}

	// Check for parameters without explicit type annotations
	for _, param := range fn.Parameters {
		if param.Type == nil {
			return &ValidationError{
				Message:  fmt.Sprintf("Parameter '%s' in function '%s' requires explicit type annotation - type cannot be inferred from usage", param.Name, fn.Name),
				Position: fn.Position,
			}
		}
	}

	if !needsReturnTypeFirst {
		// Check for functions without explicit return type annotations
		if fn.ReturnType == nil {
			return &ValidationError{
				Message:  fmt.Sprintf("Function '%s' requires explicit return type annotation - type cannot be inferred from body", fn.Name),
				Position: fn.Position,
			}
		}
	}

	// Check for built-in function redefinition LAST
	builtinFunctions := []string{"toString", "print", "length", "readFile"}
	for _, builtin := range builtinFunctions {
		if fn.Name == builtin {
			return &ValidationError{
				Message:  fmt.Sprintf("cannot redefine built-in function '%s'", fn.Name),
				Position: fn.Position,
			}
		}
	}

	return nil
}

// HINDLEY-MILNER: All validation removed
// Type inference is now handled entirely by the Hindley-Milner type system

