package codegen

import (
	"sort"

	"github.com/christianfindlay/osprey/internal/ast"
)

// validateMatchExpression validates match expressions for exhaustiveness and unknown variants.
func (g *LLVMGenerator) validateMatchExpression(matchExpr *ast.MatchExpression) error {
	// First, validate that all match arms have consistent types
	if err := g.validateMatchArmTypes(matchExpr); err != nil {
		return err
	}

	// Get the discriminant type if it's a known union type
	var unionType *ast.TypeDeclaration
	if ident, ok := matchExpr.Expression.(*ast.Identifier); ok {
		if varType, exists := g.variableTypes[ident.Name]; exists {
			if typeDecl, exists := g.typeDeclarations[varType]; exists {
				unionType = typeDecl
			}
		}
	}

	if unionType == nil {
		// If we can't determine the union type, skip validation for now
		return nil
	}

	// Validate that all patterns are known variants
	if err := g.validatePatternVariants(matchExpr.Arms, unionType, matchExpr.Position); err != nil {
		return err
	}

	// Validate exhaustiveness
	if err := g.validateExhaustiveness(matchExpr.Arms, unionType, matchExpr.Position); err != nil {
		return err
	}

	return nil
}

// validateMatchArmTypes validates that all match arms have consistent return types.
func (g *LLVMGenerator) validateMatchArmTypes(matchExpr *ast.MatchExpression) error {
	if len(matchExpr.Arms) <= 1 {
		return nil // Single arm or empty match can't have type mismatch
	}

	// Get the type of the first arm
	firstArmType := g.analyzeReturnType(matchExpr.Arms[0].Expression)

	// Don't validate if the first arm is 'any' - this might be a complex expression
	if firstArmType == TypeAny {
		return nil
	}

	// Check all subsequent arms for type consistency
	for i := 1; i < len(matchExpr.Arms); i++ {
		armType := g.analyzeReturnType(matchExpr.Arms[i].Expression)

		// Skip validation if we get 'any' type - might be complex expressions or type constructors
		if armType == TypeAny {
			continue
		}

		// Only validate if we have concrete, simple types and they differ
		if armType != firstArmType && firstArmType != TypeInt && firstArmType != TypeString && firstArmType != TypeBool {
			continue // Skip validation for complex types
		}

		// Only enforce strict validation for simple literal type mismatches (like int vs string)
		if armType != firstArmType &&
			((firstArmType == TypeInt && armType == TypeString) ||
				(firstArmType == TypeString && armType == TypeInt) ||
				(firstArmType == TypeBool && armType != TypeBool)) {
			return WrapMatchArmTypeMismatchWithPos(i, armType, firstArmType, matchExpr.Position)
		}
	}

	return nil
}

// validatePatternVariants ensures all patterns in the match arms are valid variants.
func (g *LLVMGenerator) validatePatternVariants(
	arms []ast.MatchArm, unionType *ast.TypeDeclaration, pos *ast.Position,
) error {
	for _, arm := range arms {
		if arm.Pattern.Constructor == "_" || arm.Pattern.Constructor == UnknownPattern {
			continue // Wildcard patterns are always valid
		}

		// Check if the pattern is a valid variant
		found := false
		for _, variant := range unionType.Variants {
			if arm.Pattern.Constructor == variant.Name {
				found = true

				break
			}
		}

		if !found {
			return WrapMatchUnknownVariantWithPos(arm.Pattern.Constructor, unionType.Name, pos)
		}
	}

	return nil
}

// validateExhaustiveness ensures all variants of the union type are covered.
func (g *LLVMGenerator) validateExhaustiveness(
	arms []ast.MatchArm, unionType *ast.TypeDeclaration, pos *ast.Position,
) error {
	// Collect all covered patterns
	coveredPatterns := make(map[string]bool)
	hasWildcard := false

	for _, arm := range arms {
		if arm.Pattern.Constructor == "_" || arm.Pattern.Constructor == UnknownPattern {
			hasWildcard = true
		} else {
			coveredPatterns[arm.Pattern.Constructor] = true
		}
	}

	// If there's a wildcard, match is exhaustive
	if hasWildcard {
		return nil
	}

	// Check if all variants are covered
	var missingPatterns []string
	for _, variant := range unionType.Variants {
		if !coveredPatterns[variant.Name] {
			missingPatterns = append(missingPatterns, variant.Name)
		}
	}

	if len(missingPatterns) > 0 {
		sort.Strings(missingPatterns)

		return WrapMatchNotExhaustiveWithPos(missingPatterns, pos)
	}

	return nil
}
