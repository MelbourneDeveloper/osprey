package codegen

import (
	"strings"

	"github.com/christianfindlay/osprey/internal/ast"
)

// analyzeReturnType analyzes an expression to determine its return type.
func (g *LLVMGenerator) analyzeReturnType(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.StringLiteral, *ast.IntegerLiteral, *ast.BooleanLiteral:
		return g.analyzeLiteralType(expr)
	case *ast.Identifier:
		return g.analyzeIdentifierType(e)
	case *ast.BinaryExpression:
		return g.analyzeBinaryExpressionType(e)
	case *ast.CallExpression:
		return g.analyzeCallExpressionType(e)
	case *ast.MatchExpression:
		return g.analyzeMatchExpressionType(e)
	case *ast.ResultExpression:
		return g.analyzeResultExpressionType(e)
	default:
		// CRITICAL: Unknown expression types should return 'any' to trigger validation error
		return TypeAny
	}
}

// analyzeLiteralType analyzes literal expressions to determine their type.
func (g *LLVMGenerator) analyzeLiteralType(expr ast.Expression) string {
	switch expr.(type) {
	case *ast.StringLiteral:

		return TypeString
	case *ast.IntegerLiteral:

		return TypeInt
	case *ast.BooleanLiteral:

		return TypeBool
	default:

		return TypeInt
	}
}

// analyzeBinaryExpressionType analyzes binary expressions to determine their return type.
func (g *LLVMGenerator) analyzeBinaryExpressionType(expr *ast.BinaryExpression) string {
	switch expr.Operator {
	case "+", "-", "*", "/", "%":
		// For arithmetic operations, check if we can determine concrete types
		leftType := g.getOperandType(expr.Left)
		rightType := g.getOperandType(expr.Right)

		// If both operands have concrete types (not 'any'), we can infer int
		if leftType != TypeAny && rightType != TypeAny {
			return TypeInt
		}

		// If we can't determine operand types, this would be 'any'
		return TypeAny
	case "==", "!=", "<", "<=", ">", ">=":
		return TypeBool
	default:
		return TypeInt
	}
}

// getOperandType determines the type of an operand in an expression.
func (g *LLVMGenerator) getOperandType(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return TypeInt
	case *ast.StringLiteral:
		return TypeString
	case *ast.BooleanLiteral:
		return TypeBool
	case *ast.Identifier:
		// Check if the identifier is a parameter with explicit type
		if g.currentFunctionParameterTypes != nil {
			if paramType, exists := g.currentFunctionParameterTypes[e.Name]; exists && paramType != TypeAny {
				return paramType
			}
		}
		// Check if it's a variable with known type
		if varType, exists := g.variableTypes[e.Name]; exists && varType != TypeAny {
			return varType
		}
		// Unknown identifiers would be 'any'
		return TypeAny
	case *ast.BinaryExpression:
		// Handle nested binary expressions recursively
		return g.analyzeBinaryExpressionType(e)
	case *ast.CallExpression:
		// Handle function calls
		return g.analyzeCallExpressionType(e)
	case *ast.ResultExpression:
		// Handle result expressions by analyzing the wrapped value
		return g.getOperandType(e.Value)
	default:
		// Other expressions would need to be analyzed recursively, but for now return 'any'
		return TypeAny
	}
}

// analyzeIdentifierType analyzes an identifier to determine its type based on naming patterns and context.
func (g *LLVMGenerator) analyzeIdentifierType(ident *ast.Identifier) string {
	if g.isStringTypeByName(ident.Name) {
		return TypeString
	}

	if varType, exists := g.variableTypes[ident.Name]; exists {
		return varType
	}

	// Check if it's a union type variant
	if _, exists := g.unionVariants[ident.Name]; exists {
		// Find which union type this variant belongs to
		for typeName, typeDecl := range g.typeDeclarations {
			for _, variant := range typeDecl.Variants {
				if variant.Name == ident.Name {
					return typeName
				}
			}
		}
	}

	if g.isIntTypeByName(ident.Name) {
		return TypeInt
	}

	// Check if this identifier is a function parameter with explicit type
	if g.currentFunctionParameterTypes != nil {
		if paramType, exists := g.currentFunctionParameterTypes[ident.Name]; exists {
			return paramType
		}
	}

	// CRITICAL: If we can't determine the type, return 'any' to trigger validation error
	// This catches cases like 'fn identity(x) = x' where x could be anything
	return TypeAny
}

// isStringTypeByName checks if an identifier name suggests a string type.
func (g *LLVMGenerator) isStringTypeByName(name string) bool {
	lowerName := strings.ToLower(name)
	stringIndicators := []string{
		"name", "text", "str", "title", "message", "user", "data",
		"input", "result", "output", "value", "letter", "word",
		"char", "content", "label",
	}

	for _, indicator := range stringIndicators {
		if strings.Contains(lowerName, indicator) {
			return true
		}
	}

	return false
}

// isIntTypeByName checks if an identifier name suggests an integer type.
func (g *LLVMGenerator) isIntTypeByName(name string) bool {
	lowerName := strings.ToLower(name)
	intIndicators := []string{"age", "score", "count", "num", "val"}

	for _, indicator := range intIndicators {
		if strings.Contains(lowerName, indicator) {
			return true
		}
	}

	return false
}

// analyzeCallExpressionType analyzes a call expression to determine its return type.
func (g *LLVMGenerator) analyzeCallExpressionType(callExpr *ast.CallExpression) string {
	if ident, ok := callExpr.Function.(*ast.Identifier); ok {
		if returnType, exists := g.functionReturnTypes[ident.Name]; exists {
			return returnType
		}
	}

	// CRITICAL: Unknown function calls should return 'any' to trigger validation error
	// This catches cases like 'fn process() = someUnknownFunction()'
	return TypeAny
}

// analyzeMatchExpressionType analyzes a match expression to determine its return type.
func (g *LLVMGenerator) analyzeMatchExpressionType(matchExpr *ast.MatchExpression) string {
	for _, arm := range matchExpr.Arms {
		if g.analyzeReturnType(arm.Expression) == TypeString {
			return TypeString
		}
	}

	return TypeInt
}

// analyzeResultExpressionType analyzes a result expression to determine its return type.
func (g *LLVMGenerator) analyzeResultExpressionType(resultExpr *ast.ResultExpression) string {
	// Analyze the inner value of the result expression
	return g.analyzeReturnType(resultExpr.Value)
}

// analyzeParameterUsage analyzes how a parameter is used in a function body to infer its type.
func (g *LLVMGenerator) analyzeParameterUsage(paramName string, body ast.Expression) string {
	switch expr := body.(type) {
	case *ast.BinaryExpression:
		return g.analyzeParameterInBinaryExpr(paramName, expr)
	case *ast.Identifier:
		return g.analyzeParameterInIdentifier(paramName, expr)
	case *ast.CallExpression:
		return g.analyzeParameterInCallExpr(paramName, expr)
	case *ast.MatchExpression:
		return g.analyzeParameterInMatchExpr(paramName, expr)
	case *ast.ResultExpression:
		// Analyze the inner value of the result expression
		return g.analyzeParameterUsage(paramName, expr.Value)
	default:
		return ""
	}
}

// analyzeParameterInBinaryExpr checks if parameter is used in arithmetic operations.
func (g *LLVMGenerator) analyzeParameterInBinaryExpr(paramName string, expr *ast.BinaryExpression) string {
	if g.isParameterInArithmetic(paramName, expr) {
		return TypeInt
	}

	return ""
}

// analyzeParameterInIdentifier checks if the parameter is directly returned.
func (g *LLVMGenerator) analyzeParameterInIdentifier(paramName string, expr *ast.Identifier) string {
	if expr.Name == paramName {
		return ""
	}

	return ""
}

// analyzeParameterInCallExpr checks if parameter is passed to another function.
func (g *LLVMGenerator) analyzeParameterInCallExpr(paramName string, expr *ast.CallExpression) string {
	for _, arg := range expr.Arguments {
		if ident, ok := arg.(*ast.Identifier); ok && ident.Name == paramName {
			return ""
		}
	}

	return ""
}

// analyzeParameterInMatchExpr analyzes parameter usage in match expressions.
func (g *LLVMGenerator) analyzeParameterInMatchExpr(paramName string, expr *ast.MatchExpression) string {
	if matchDiscriminant, ok := expr.Expression.(*ast.Identifier); ok && matchDiscriminant.Name == paramName {
		return g.analyzeMatchPatterns(expr.Arms)
	}

	if nestedType := g.analyzeParameterUsage(paramName, expr.Expression); nestedType != "" {
		return nestedType
	}

	for _, arm := range expr.Arms {
		if armType := g.analyzeParameterUsage(paramName, arm.Expression); armType != "" {
			return armType
		}
	}

	return ""
}

// analyzeMatchPatterns analyzes match patterns to determine the return type.
func (g *LLVMGenerator) analyzeMatchPatterns(arms []ast.MatchArm) string {
	for _, arm := range arms {
		if g.analyzeReturnType(arm.Expression) == TypeString {
			return TypeString
		}
	}

	return TypeInt
}

// isParameterInArithmetic checks if a parameter is used in arithmetic operations.
func (g *LLVMGenerator) isParameterInArithmetic(paramName string, expr *ast.BinaryExpression) bool {
	// Check if this binary expression uses arithmetic operators
	if expr.Operator == "+" || expr.Operator == "-" || expr.Operator == "*" || expr.Operator == "/" {
		// Check if the parameter appears in either side

		return g.containsParameter(paramName, expr.Left) || g.containsParameter(paramName, expr.Right)
	}

	return false
}

// containsParameter recursively checks if an expression contains a reference to a parameter.
func (g *LLVMGenerator) containsParameter(paramName string, expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.Identifier:

		return e.Name == paramName
	case *ast.BinaryExpression:

		return g.containsParameter(paramName, e.Left) || g.containsParameter(paramName, e.Right)
	case *ast.CallExpression:
		for _, arg := range e.Arguments {
			if g.containsParameter(paramName, arg) {
				return true
			}
		}
	}

	return false
}

// isAnyType checks if an expression evaluates to the 'any' type.
func (g *LLVMGenerator) isAnyType(expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.Identifier:
		// Check if the variable is of type 'any'
		if varType, exists := g.variableTypes[e.Name]; exists {
			return varType == TypeAny
		}

		return false
	case *ast.CallExpression:
		// Check if the function returns 'any'
		if ident, ok := e.Function.(*ast.Identifier); ok {
			if returnType, exists := g.functionReturnTypes[ident.Name]; exists {
				return returnType == TypeAny
			}
		}

		return false
	default:

		return false
	}
}

// validateNotAnyType validates that an expression is not of type 'any' for operations that require concrete types.
func (g *LLVMGenerator) validateNotAnyType(expr ast.Expression, operation string) error {
	if g.isAnyType(expr) {
		switch operation {
		case AnyOpArithmetic:
			return WrapAnyDirectArithmetic(operation)
		case AnyOpFieldAccess:
			return WrapAnyDirectFieldAccess("unknown")
		case AnyOpFunctionArgument:
			return WrapAnyDirectFunctionArg("unknown", "unknown")
		default:
			return WrapAnyDirectVariableAccess("unknown")
		}
	}

	return nil
}

// validateFunctionArgument validates function arguments with strong typing rules.
// Allows function composition (function types) but prevents 'any' from being passed to non-function parameters.
func (g *LLVMGenerator) validateFunctionArgument(expr ast.Expression, functionName, paramName string,
	pos *ast.Position) error {
	// Only validate if the expression is of type 'any'
	if !g.isAnyType(expr) {
		return nil
	}

	// Check if this is a function reference (for function composition)
	// If the 'any' variable actually contains a function reference, allow it
	if ident, ok := expr.(*ast.Identifier); ok {
		if varType, exists := g.variableTypes[ident.Name]; exists && varType == TypeFunction {
			return nil // Allow function composition - this is a function reference
		}
	}

	// For all other cases, 'any' type is not allowed (strong typing)
	return WrapAnyDirectFunctionArgWithPos(functionName, paramName, pos)
}
