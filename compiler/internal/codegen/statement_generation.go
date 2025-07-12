package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) generateStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ImportStatement:
		// Imports are handled at compile time
		return nil

	case *ast.LetDeclaration:
		_, err := g.generateLetDeclaration(s)
		return err

	case *ast.AssignmentStatement:
		return g.generateAssignmentStatement(s)

	case *ast.FunctionDeclaration:
		// Use the unified function declaration generator from function_signatures.go
		return g.declareFunctionSignature(s)

	case *ast.ExternDeclaration:
		return g.generateExternDeclaration(s)

	case *ast.TypeDeclaration:
		// Type declarations are handled in first pass
		return nil

	case *ast.EffectDeclaration:
		return g.generateEffectDeclaration(s)

	case *ast.ExpressionStatement:
		_, err := g.generateExpression(s.Expression)
		return err

	default:
		return WrapUnsupportedStatement(stmt)
	}
}

func (g *LLVMGenerator) generateExternDeclaration(externDecl *ast.ExternDeclaration) error {
	// Convert extern parameters to LLVM parameters
	params := make([]*ir.Param, len(externDecl.Parameters))
	paramNames := make([]string, len(externDecl.Parameters))

	for i, param := range externDecl.Parameters {
		llvmType := g.typeExpressionToLLVMType(&param.Type)
		params[i] = ir.NewParam(param.Name, llvmType)
		paramNames[i] = param.Name
	}

	// Determine return type
	var returnType types.Type = types.I64 // Default to int
	if externDecl.ReturnType != nil {
		returnType = g.typeExpressionToLLVMType(externDecl.ReturnType)
	}

	// Declare the external function
	externFunc := g.module.NewFunc(externDecl.Name, returnType, params...)
	g.functions[externDecl.Name] = externFunc
	// Built-in functions are handled by Hindley-Milner type inference
	g.functionParameters[externDecl.Name] = paramNames

	return nil
}

// typeExpressionToLLVMType converts an Osprey TypeExpression to an LLVM type.
func (g *LLVMGenerator) typeExpressionToLLVMType(typeExpr *ast.TypeExpression) types.Type {
	// Handle function types
	if typeExpr.IsFunction {
		// Build parameter types
		paramTypes := make([]types.Type, len(typeExpr.ParameterTypes))
		for i, paramType := range typeExpr.ParameterTypes {
			paramTypes[i] = g.typeExpressionToLLVMType(&paramType)
		}

		// Build return type
		var returnType types.Type = types.I64 // Default to int
		if typeExpr.ReturnType != nil {
			returnType = g.typeExpressionToLLVMType(typeExpr.ReturnType)
		}

		// Create function signature
		funcSig := types.NewFunc(returnType, paramTypes...)

		// Return pointer to function (function pointer type)
		return types.NewPointer(funcSig)
	}

	switch typeExpr.Name {
	case "Int":
		return types.I64
	case "String":
		return types.I8Ptr
	case TypeUnit:
		return types.Void
	case TypeHTTPResponse:
		// Return pointer to HttpResponse struct
		return types.NewPointer(g.typeMap[TypeHTTPResponse])
	default:
		// Check if it's a user-defined type
		if llvmType, exists := g.typeMap[typeExpr.Name]; exists {
			return llvmType
		}
		// Default to i64 for unknown types
		return types.I64
	}
}

func (g *LLVMGenerator) generateLetDeclaration(letDecl *ast.LetDeclaration) (value.Value, error) {
	value, err := g.generateExpression(letDecl.Value)
	if err != nil {
		return nil, err
	}

	// Store the value in our variable map
	g.variables[letDecl.Name] = value

	// CRITICAL FIX: Use explicit type annotation if present, otherwise infer from value
	var varType Type
	if letDecl.Type != nil {
		// Use the explicit type annotation
		varType = &ConcreteType{name: letDecl.Type.Name}
	} else {
		// Use unified type inference system for untyped declarations
		inferredType, err := g.typeInferer.InferType(letDecl.Value)
		if err != nil {
			return nil, err
		}
		varType = inferredType
	}

	// Store the type in the Hindley-Milner environment
	g.typeInferer.env.Set(letDecl.Name, varType)

	// Track if this variable is mutable
	g.mutableVariables[letDecl.Name] = letDecl.Mutable

	return value, nil
}

// generateAssignmentStatement generates LLVM IR for mutable variable assignments.
func (g *LLVMGenerator) generateAssignmentStatement(assignStmt *ast.AssignmentStatement) error {
	// Check if the variable exists in the Hindley-Milner type environment (single source of truth)
	if _, exists := g.typeInferer.env.Get(assignStmt.Name); !exists {
		return WrapUndefinedVariableWithPos(assignStmt.Name, assignStmt.Position)
	}

	// Check if the variable is mutable
	if mutable, exists := g.mutableVariables[assignStmt.Name]; !exists || !mutable {
		return WrapImmutableAssignmentErrorWithPos(assignStmt.Name, assignStmt.Position)
	}

	// Generate the new value
	newValue, err := g.generateExpression(assignStmt.Value)
	if err != nil {
		return err
	}

	// Use unified type inference system
	inferredType, err := g.typeInferer.InferType(assignStmt.Value)
	if err != nil {
		return err
	}

	// Verify type compatibility using unification
	existingType := g.typeInferer.env.vars[assignStmt.Name]
	if err := g.typeInferer.Unify(existingType, inferredType); err != nil {
		return fmt.Errorf("type mismatch in assignment: %w", err)
	}

	// Update the variable
	g.variables[assignStmt.Name] = newValue
	// Update type in Hindley-Milner environment
	g.typeInferer.env.Set(assignStmt.Name, inferredType)

	return nil
}
