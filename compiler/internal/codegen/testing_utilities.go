package codegen

import (
	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// Public methods for testing

// AnalyzeReturnType exposes the analyzeReturnType method for testing.
func (g *LLVMGenerator) AnalyzeReturnType(expr ast.Expression) string {
	return g.analyzeReturnType(expr)
}

// DeclareFunctionSignature exposes the declareFunctionSignature method for testing.
func (g *LLVMGenerator) DeclareFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	return g.declareFunctionSignature(fnDecl)
}

// GetFunction returns a function by name for testing.
func (g *LLVMGenerator) GetFunction(name string) *ir.Func {
	return g.functions[name]
}

// GetFunctionReturnType returns the return type of a function for testing.
func (g *LLVMGenerator) GetFunctionReturnType(name string) string {
	if returnType, exists := g.functionReturnTypes[name]; exists {
		return returnType
	}

	return "int" // Default
}

// SetVariableType sets a variable type for testing.
func (g *LLVMGenerator) SetVariableType(name, varType string) {
	g.variableTypes[name] = varType
}

// SetParameterTypesForAnalysis exposes the setParameterTypesForAnalysis method for testing.
func (g *LLVMGenerator) SetParameterTypesForAnalysis(fnDecl *ast.FunctionDeclaration) {
	g.setParameterTypesForAnalysis(fnDecl)
}

// ClearParameterTypesForAnalysis exposes the clearParameterTypesForAnalysis method for testing.
func (g *LLVMGenerator) ClearParameterTypesForAnalysis() {
	g.clearParameterTypesForAnalysis()
}

// CreateAndStoreFunctionSignature exposes the createAndStoreFunctionSignature method for testing.
func (g *LLVMGenerator) CreateAndStoreFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	return g.createAndStoreFunctionSignature(fnDecl)
}

// reorderNamedArguments reorders named arguments to match the function's parameter order
func (g *LLVMGenerator) reorderNamedArguments(
	functionName string,
	namedArgs []ast.NamedArgument,
) ([]value.Value, error) {
	paramNames, exists := g.functionParameters[functionName]
	if !exists {
		return nil, WrapFunctionNotFound(functionName)
	}

	if len(namedArgs) != len(paramNames) {
		return nil, WrapWrongArgCount(functionName, len(paramNames), len(namedArgs))
	}

	// Create a map for quick lookup of arguments by name
	argMap := make(map[string]ast.NamedArgument)
	for _, arg := range namedArgs {
		argMap[arg.Name] = arg
	}

	// Reorder arguments to match parameter order
	orderedArgs := make([]value.Value, len(paramNames))

	for i, paramName := range paramNames {
		namedArg, exists := argMap[paramName]
		if !exists {
			return nil, WrapMissingArgument(paramName, functionName)
		}

		// STRONG TYPING: Validate that 'any' type cannot be passed to non-function parameters
		// This preserves type safety while allowing function composition
		// Extract position from identifier arguments only
		var pos *ast.Position
		if ident, ok := namedArg.Value.(*ast.Identifier); ok {
			pos = ident.Position
		}

		if err := g.validateFunctionArgument(namedArg.Value, functionName, paramName, pos); err != nil {
			return nil, err
		}

		val, err := g.generateExpression(namedArg.Value)
		if err != nil {
			return nil, err
		}

		orderedArgs[i] = val
	}

	return orderedArgs, nil
}
