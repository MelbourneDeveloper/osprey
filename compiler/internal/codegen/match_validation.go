package codegen

import (
	"fmt"
	"strconv"

	"github.com/christianfindlay/osprey/internal/ast"
)

const (
	// MinParametersForNamedArgs is the minimum number of parameters required to enforce named arguments
	MinParametersForNamedArgs = 2
)

// validateMatchExpressionWithType validates match expressions with discriminant type information
func (g *LLVMGenerator) validateMatchExpressionWithType(expr *ast.MatchExpression, discriminantType string) error {
	// Check if match expression has no arms
	if len(expr.Arms) == 0 {
		if expr.Position != nil {
			//nolint:err113 // Dynamic error needed for exact test format matching
			return fmt.Errorf("line %d:%d: match expression must have at least one arm",
				expr.Position.Line, expr.Position.Column)
		}

		return ErrMatchNoArms
	}

	for _, arm := range expr.Arms {
		err := g.validateMatchArmWithTypeAndPosition(arm, discriminantType, expr.Position)
		if err != nil {
			return err
		}
	}

	// Check for exhaustiveness (using the existing function)
	err := g.validateMatchExhaustiveness(expr)
	if err != nil {
		return err
	}

	return nil
}

// validateMatchExhaustiveness checks if all union variants are covered in match expressions
func (g *LLVMGenerator) validateMatchExhaustiveness(expr *ast.MatchExpression) error {
	// First, we need to infer the type of the expression being matched
	exprType, err := g.typeInferer.InferType(expr.Expression)
	if err != nil {
		return err
	}

	// Resolve the type to get the concrete type
	resolvedType := g.typeInferer.ResolveType(exprType)
	typeName := resolvedType.String()

	// Check if this is a union type
	typeDecl, exists := g.typeDeclarations[typeName]
	if !exists {
		// Not a union type, no exhaustiveness check needed
		return nil
	}

	// Only check exhaustiveness for union types with multiple variants
	if len(typeDecl.Variants) <= 1 {
		return nil
	}

	// Collect all patterns from match arms
	coveredVariants := make(map[string]bool)
	hasWildcard := false

	for _, arm := range expr.Arms {
		pattern := arm.Pattern.Constructor
		if pattern == "_" {
			hasWildcard = true
			break // Wildcard covers all remaining cases
		}

		if pattern != "" {
			coveredVariants[pattern] = true
		}
	}

	// If there's a wildcard, the match is exhaustive
	if hasWildcard {
		return nil
	}

	// Check if all variants are covered
	var missingVariants []string

	for _, variant := range typeDecl.Variants {
		if !coveredVariants[variant.Name] {
			missingVariants = append(missingVariants, variant.Name)
		}
	}

	if len(missingVariants) > 0 {
		return WrapMatchNotExhaustiveWithPos(missingVariants, expr.Position)
	}

	return nil
}

// validateNamedArguments validates that multi-parameter functions require named arguments
func (g *LLVMGenerator) validateNamedArguments(funcName string, callExpr *ast.CallExpression) error {
	// Check if this is a user-defined function with multiple parameters
	paramNames, exists := g.functionParameters[funcName]
	if !exists {
		return nil // Built-in or unknown function
	}

	// Only enforce named arguments for multi-parameter functions
	if len(paramNames) < MinParametersForNamedArgs {
		return nil
	}

	// If the function has 2 or more parameters and positional arguments are used, require named arguments
	if len(callExpr.Arguments) > 0 && len(callExpr.NamedArguments) == 0 {
		return WrapFunctionRequiresNamedArgsWithPos(funcName, len(paramNames), callExpr.Position)
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

func (g *LLVMGenerator) validateMatchArmWithTypeAndPosition(
	arm ast.MatchArm, discriminantType string, matchPos *ast.Position,
) error {
	return g.validateMatchPatternWithTypeAndPosition(arm.Pattern, discriminantType, matchPos)
}

func (g *LLVMGenerator) validateMatchPatternWithTypeAndPosition(
	pattern ast.Pattern, discriminantType string, matchPos *ast.Position,
) error {
	// Wildcard patterns and variable patterns are always valid
	if pattern.Constructor == "_" || pattern.Constructor == "" {
		return nil
	}

	// Literal patterns (integers, strings, booleans) are always valid for their type
	if isLiteralPattern(pattern.Constructor) {
		return nil
	}

	// Special constructors that are always allowed (Result types, etc.)
	if isSpecialConstructor(pattern.Constructor) {
		return nil
	}

	// Check if this pattern matches a variant of the discriminant's union type
	if typeDecl, exists := g.typeDeclarations[discriminantType]; exists {
		// Check if the pattern constructor is a valid variant of this type
		isValidVariant := false
		for _, variant := range typeDecl.Variants {
			if variant.Name == pattern.Constructor {
				isValidVariant = true
				break
			}
		}

		// If not a valid variant, return error
		if !isValidVariant {
			return WrapUnknownVariantWithPos(pattern.Constructor, discriminantType, matchPos)
		}
	}

	// Pattern validation is complete - type inference and variable binding
	// happen during the match expression type inference phase, not here
	return nil
}

// isLiteralPattern checks if a pattern constructor is a literal value
func isLiteralPattern(constructor string) bool {
	// Boolean literals
	if constructor == "true" || constructor == "false" {
		return true
	}

	// Integer literals (try parsing)
	_, err := strconv.ParseInt(constructor, 10, 64)
	if err == nil {
		return true
	}

	// String literals (quoted)
	if len(constructor) >= 2 && constructor[0] == '"' && constructor[len(constructor)-1] == '"' {
		return true
	}

	return false
}

// isSpecialConstructor checks if a constructor is a special built-in constructor
func isSpecialConstructor(constructor string) bool {
	// Result type constructors
	return constructor == SuccessPattern || constructor == ErrorPattern
}
