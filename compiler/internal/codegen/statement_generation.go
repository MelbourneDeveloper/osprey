package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
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

	case *ast.PluginFunctionDeclaration:
		err := g.generatePluginFunctionDeclaration(s)

		return err

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

	// Add extern function to type environment for type inference
	// Convert extern parameters to type inference types
	paramTypes := make([]Type, len(externDecl.Parameters))
	for i, param := range externDecl.Parameters {
		paramTypes[i] = g.typeExpressionToInferenceType(&param.Type)
	}

	// Determine inference return type
	var inferenceReturnType Type = &ConcreteType{name: TypeInt} // Default to int
	if externDecl.ReturnType != nil {
		inferenceReturnType = g.typeExpressionToInferenceType(externDecl.ReturnType)
	}

	// Add extern function to type environment
	functionType := &FunctionType{
		paramTypes: paramTypes,
		returnType: inferenceReturnType,
	}
	g.typeInferer.env.Set(externDecl.Name, functionType)

	return nil
}

func (g *LLVMGenerator) generatePluginFunctionDeclaration(pluginDecl *ast.PluginFunctionDeclaration) error {
	// Generate LLVM function that represents a plugin function call

	// Build parameter list for LLVM function signature
	var params []*ir.Param
	var paramNames []string

	for _, param := range pluginDecl.Parameters {
		paramType := g.typeExpressionToLLVMType(param.Type)
		params = append(params, ir.NewParam(param.Name, paramType))
		paramNames = append(paramNames, param.Name)
	}

	// Plugin functions return strings for now (would be Result<T, Error> in full implementation)
	returnType := types.I8Ptr // String type

	// Create the plugin function
	pluginFunc := g.module.NewFunc(pluginDecl.FunctionName, returnType, params...)
	g.functions[pluginDecl.FunctionName] = pluginFunc
	g.functionReturnTypes[pluginDecl.FunctionName] = TypeString
	g.functionParameters[pluginDecl.FunctionName] = paramNames

	// Save current builder context (might be nil during first pass)
	savedBuilder := g.builder

	// Create function body
	entry := pluginFunc.NewBlock("entry")
	g.builder = entry

	// Generate plugin call representation
	pluginResult := g.generatePluginCallResult(pluginDecl)

	// Return the result
	entry.NewRet(pluginResult)

	// Restore original builder context (only if it wasn't nil)
	if savedBuilder != nil {
		g.builder = savedBuilder
	} else {
		g.builder = nil
	}

	return nil
}

func (g *LLVMGenerator) generatePluginCallResult(pluginDecl *ast.PluginFunctionDeclaration) value.Value {
	// Generate a result string that represents the plugin function call
	// In a full implementation, this would:
	// 1. Call the plugin system runtime
	// 2. Serialize parameters to JSON
	// 3. Execute the plugin executable
	// 4. Handle the JSON response
	// 5. Return Result<T, PluginError>

	// For now, create a descriptive string
	pluginInfo := fmt.Sprintf("[Plugin %s:%s]", pluginDecl.PluginName, pluginDecl.FunctionName)

	// Create a global string constant using the same approach as string literals
	str := constant.NewCharArrayFromString(pluginInfo + StringTerminator)
	global := g.module.NewGlobalDef("", str)

	// Return pointer to the string
	return g.builder.NewGetElementPtr(str.Typ, global,
		constant.NewInt(types.I32, ArrayIndexZero),
		constant.NewInt(types.I32, ArrayIndexZero))
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
	case TypeInt:
		return types.I64
	case TypeString:
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

// typeExpressionToInferenceType converts an Osprey TypeExpression to a type inference Type.
func (g *LLVMGenerator) typeExpressionToInferenceType(typeExpr *ast.TypeExpression) Type {
	// Handle function types
	if typeExpr.IsFunction {
		// Build parameter types
		paramTypes := make([]Type, len(typeExpr.ParameterTypes))
		for i, paramType := range typeExpr.ParameterTypes {
			paramTypes[i] = g.typeExpressionToInferenceType(&paramType)
		}

		// Build return type
		var returnType Type = &ConcreteType{name: TypeInt} // Default to int
		if typeExpr.ReturnType != nil {
			returnType = g.typeExpressionToInferenceType(typeExpr.ReturnType)
		}

		// Create function type
		return &FunctionType{
			paramTypes: paramTypes,
			returnType: returnType,
		}
	}

	// Handle generic types like Result<bool, MathError>
	if len(typeExpr.GenericParams) > 0 {
		// Convert generic parameters to Type arguments
		typeArgs := make([]Type, len(typeExpr.GenericParams))
		for i, genericParam := range typeExpr.GenericParams {
			typeArgs[i] = g.typeExpressionToInferenceType(&genericParam)
		}

		genericType := NewGenericType(typeExpr.Name, typeArgs)

		return genericType
	}

	switch typeExpr.Name {
	case TypeInt:
		return &ConcreteType{name: TypeInt}
	case "string":
		return &ConcreteType{name: TypeString}
	case TypeBool:
		return &ConcreteType{name: TypeBool}
	case TypeUnit:
		return &ConcreteType{name: TypeUnit}
	case TypeHTTPResponse:
		return &ConcreteType{name: TypeHTTPResponse}
	case TypeFiber:
		return &ConcreteType{name: TypeFiber}
	case TypeChannel:
		return &ConcreteType{name: TypeChannel}
	default:
		// Check if this is a user-defined type (record or union)
		if typeDecl, exists := g.typeDeclarations[typeExpr.Name]; exists {
			// If it's a single-variant record type, return RecordType
			if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
				fields := make(map[string]Type)

				variant := &typeDecl.Variants[0]
				for _, field := range variant.Fields {
					// For now, use the field's declared type or default to int
					// This should ideally be more sophisticated type resolution
					var fieldType Type

					switch field.Type {
					case TypeInt:
						fieldType = &ConcreteType{name: TypeInt}
					case TypeString:
						fieldType = &ConcreteType{name: TypeString}
					case TypeBool:
						fieldType = &ConcreteType{name: TypeBool}
					default:
						fieldType = &ConcreteType{name: field.Type}
					}

					fields[field.Name] = fieldType
				}

				return NewRecordType(typeExpr.Name, fields)
			}
		}

		// For unknown types without generic parameters, return as concrete type
		return &ConcreteType{name: typeExpr.Name}
	}
}

func (g *LLVMGenerator) generateLetDeclaration(letDecl *ast.LetDeclaration) (value.Value, error) {
	value, err := g.generateExpression(letDecl.Value)
	if err != nil {
		return nil, err
	}

	// Store the value in our variable map
	g.variables[letDecl.Name] = value

	// Use explicit type annotation if present, otherwise infer from value
	var varType Type

	if letDecl.Type != nil {
		// Use the explicit type annotation, but validate it matches the value
		annotatedType := g.typeExpressionToInferenceType(letDecl.Type)

		// Infer the actual type of the value
		valueType, err := g.typeInferer.InferType(letDecl.Value)
		if err != nil {
			return nil, err
		}

		// Check if the value type is compatible with the annotation
		err = g.typeInferer.Unify(annotatedType, valueType)
		if err != nil {
			return nil, WrapTypeMismatchWithPos(
				valueType.String(), letDecl.Name, annotatedType.String(), letDecl.Position)
		}

		varType = annotatedType
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
	err = g.typeInferer.Unify(existingType, inferredType)
	if err != nil {
		return fmt.Errorf("type mismatch in assignment: %w", err)
	}

	// Update the variable
	g.variables[assignStmt.Name] = newValue
	// Update type in Hindley-Milner environment
	g.typeInferer.env.Set(assignStmt.Name, inferredType)

	return nil
}
