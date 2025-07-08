package codegen

import (
	"fmt"

	"github.com/christianfindlay/osprey/internal/ast"
)

// validateMatchExpression validates match expressions for exhaustiveness and unknown variants.
func (g *LLVMGenerator) validateMatchExpression(expr *ast.MatchExpression) error {
	for _, arm := range expr.Arms {
		if err := g.validateMatchArm(arm); err != nil {
			return err
		}
	}
	return nil
}

// reorderNamedArguments reorders function arguments based on parameter names.
func (g *LLVMGenerator) reorderNamedArguments(fnName string, args []ast.NamedArgument) ([]ast.Expression, error) {
	paramNames, exists := g.functionParameters[fnName]
	if !exists {
		// Convert NamedArguments to Expressions
		exprs := make([]ast.Expression, len(args))
		for i, arg := range args {
			exprs[i] = arg.Value
		}
		return exprs, nil // No parameter info, keep original order
	}

	// Create mapping of parameter names to positions
	paramPositions := make(map[string]int)
	for i, name := range paramNames {
		paramPositions[name] = i
	}

	// Create new argument slice in correct order
	reorderedArgs := make([]ast.Expression, len(args))

	// Handle named arguments
	for _, arg := range args {
		pos, exists := paramPositions[arg.Name]
		if !exists {
			return nil, fmt.Errorf("%w: %s", ErrUnknownParameterName, arg.Name)
		}
		reorderedArgs[pos] = arg.Value
	}

	return reorderedArgs, nil
}

func (g *LLVMGenerator) validateMatchPattern(pattern ast.Pattern) error {
	// Infer pattern type
	_, err := g.typeInferer.InferPattern(pattern)
	if err != nil {
		return err
	}
	return nil
}

func (g *LLVMGenerator) validateMatchArm(arm ast.MatchArm) error {
	return g.validateMatchPattern(arm.Pattern)
}
