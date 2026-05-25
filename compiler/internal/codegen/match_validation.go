package codegen

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/christianfindlay/osprey/internal/ast"
)

const (
	// MinParametersForNamedArgs is the minimum number of parameters required to enforce named arguments
	MinParametersForNamedArgs = 2
)

// errUnsupportedPatternSyntax surfaces when the AST builder fell off
// its pattern-shape cases and emitted Constructor="unknown".
var errUnsupportedPatternSyntax = errors.New("unsupported pattern syntax — only literals, _, ID, ID name, " +
	"and `Variant { field, field }` destructuring are supported")

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

	// Spec 0009-BooleanOperations.md: "match on a boolean ... forces both
	// arms to be considered." Without this check the lowered IR
	// unconditionally branches to the last arm, so `b=false; match b { true => … }`
	// silently runs the true arm. Skip when the arms use Result-style
	// patterns (Success/Error) — the codegen auto-wraps the bool into a
	// Result there and the user's intent is to match Result, not bool.
	if typeName == TypeBool && !hasResultArms(expr) {
		return ensureBoolMatchExhaustive(expr)
	}

	// Same rationale for unbounded primitive types: int / float / string
	// matches over literals fall through to the last arm at runtime if
	// no wildcard / variable-bind catch-all is present, so
	// `match x { 1 => "one" }` silently returns "one" for every input.
	// Require an explicit catch-all to make the partial intent loud.
	if typeName == TypeInt || typeName == TypeFloat || typeName == TypeString {
		if !hasResultArms(expr) {
			return ensurePrimitiveMatchHasCatchAll(expr, typeName)
		}
	}

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

// hasResultArms reports whether any arm of the match uses a Success or
// Error pattern. Used to skip bool-exhaustiveness on matches whose intent
// is Result destructuring (the codegen auto-wraps the discriminant).
func hasResultArms(expr *ast.MatchExpression) bool {
	for _, arm := range expr.Arms {
		if arm.Pattern.Constructor == SuccessPattern || arm.Pattern.Constructor == ErrorPattern {
			return true
		}
	}
	return false
}

// ensurePrimitiveMatchHasCatchAll errors when a match on an unbounded
// primitive type (int / float / string) has no `_` or variable-bind
// catch-all arm. Without one the lowered IR unconditionally branches to
// the last arm, so `match x { 1 => "one" }` silently returns "one" for
// every input — the same hazard ensureBoolMatchExhaustive fixed for bool.
//
// Variable-bind arms like `n => print("Many!")` parse as
// `Constructor: "n", Variable: ""` (the AST builder can't tell at
// parse time whether `n` names a union variant or binds a fresh
// variable), so we treat any non-literal single identifier as a
// catch-all for primitive matches — int/float/string have no named
// variants.
func ensurePrimitiveMatchHasCatchAll(expr *ast.MatchExpression, typeName string) error {
	for _, arm := range expr.Arms {
		if arm.Pattern.Constructor == "_" {
			return nil
		}
		if arm.Pattern.Constructor == "" && arm.Pattern.Variable != "" {
			return nil
		}
		// Non-literal single identifier on a primitive match = variable bind.
		if arm.Pattern.Constructor != "" && !isLiteralPattern(arm.Pattern.Constructor) {
			return nil
		}
	}
	missing := fmt.Sprintf("_ (catch-all for %s)", typeName)
	return WrapMatchNotExhaustiveWithPos([]string{missing}, expr.Position)
}

// ensureBoolMatchExhaustive errors when a boolean match misses either
// `true` or `false` and has no wildcard / variable catch-all. Per spec
// 0009-BooleanOperations.md both arms must be present.
func ensureBoolMatchExhaustive(expr *ast.MatchExpression) error {
	seenTrue, seenFalse := false, false
	for _, arm := range expr.Arms {
		switch arm.Pattern.Constructor {
		case "_":
			return nil
		case TruePattern:
			seenTrue = true
		case FalsePattern:
			seenFalse = true
		case "":
			if arm.Pattern.Variable != "" {
				return nil // variable-bind catch-all
			}
		}
	}
	var missing []string
	if !seenTrue {
		missing = append(missing, TruePattern)
	}
	if !seenFalse {
		missing = append(missing, FalsePattern)
	}
	if len(missing) == 0 {
		return nil
	}
	return WrapMatchNotExhaustiveWithPos(missing, expr.Position)
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
		return WrapFunctionRequiresNamedArgsWithPos(funcName, paramNames, callExpr.Position)
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

	// Handle named arguments. Duplicate names (e.g. `add(a: 1, a: 2)`)
	// used to silently overwrite the earlier slot and leave another slot
	// nil, so codegen later barfed with "unsupported expression: <nil>".
	// Flag the duplicate up-front with a useful message.
	seen := make(map[string]bool, len(args))

	for _, arg := range args {
		if seen[arg.Name] {
			return nil, fmt.Errorf("%w: duplicate named argument '%s' for function '%s'",
				ErrUnknownParameterName, arg.Name, fnName)
		}

		seen[arg.Name] = true

		pos, exists := paramPositions[arg.Name]
		if !exists {
			return nil, fmt.Errorf("%w '%s' for function '%s'; valid parameters: %s",
				ErrUnknownParameterName, arg.Name, fnName, strings.Join(paramNames, ", "))
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

	// The AST builder emits Constructor="unknown" when it falls off the
	// end of its pattern-shape cases (e.g. for unsupported syntax like
	// `Some { value: v }` field-rename destructure). Without this branch
	// the user hits a baffling "variant 'unknown' is not defined" — make
	// the parser-level limitation visible instead.
	if pattern.Constructor == "unknown" {
		if matchPos != nil {
			return fmt.Errorf( //nolint:err113
				"line %d:%d: unsupported pattern syntax — only literals, _, ID, ID name, "+
					"and `Variant { field, field }` destructuring are supported",
				matchPos.Line, matchPos.Column)
		}
		return errUnsupportedPatternSyntax
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
	_, intErr := strconv.ParseInt(constructor, 10, 64)
	if intErr == nil {
		return true
	}

	// Float literals (e.g. `3.14`) — added so `match f { 3.14 => … }`
	// passes the validator alongside int literals.
	_, floatErr := strconv.ParseFloat(constructor, 64)
	if floatErr == nil {
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
