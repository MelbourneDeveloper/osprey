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
		err := validateStatement(stmt)
		if err != nil {
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
	// With Hindley-Milner type inference, we trust the type system to handle inference
	// Only validate basic language constraints, not type inference capabilities

	// Check for built-in function redefinition
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
