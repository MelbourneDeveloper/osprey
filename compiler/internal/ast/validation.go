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

// validateFunctionDeclaration validates function declarations according to Hindley-Milner type inference.
// With Hindley-Milner, we allow all type inference to be handled by the type system itself.
func validateFunctionDeclaration(fn *FunctionDeclaration) error {
	// HINDLEY-MILNER: Allow complete type inference
	// No validation needed - let the type inference system handle everything
	_ = fn // Use fn to avoid unused parameter warning
	return nil
}

// HINDLEY-MILNER: All validation removed
// Type inference is now handled entirely by the Hindley-Milner type system

