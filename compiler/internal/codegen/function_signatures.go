package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) declareFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	if fnDecl.Name == ToStringFunc {
		return ErrToStringReserved
	}

	// Check for protected built-in functions
	if err := CheckProtectedFunction(fnDecl); err != nil {
		return err
	}

	return g.createAndStoreFunctionSignature(fnDecl)
}

// createAndStoreFunctionSignature creates and stores the function signature.
func (g *LLVMGenerator) createAndStoreFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	// Use explicit return type if provided, otherwise analyze the body
	var returnType string
	if fnDecl.ReturnType != nil {
		switch fnDecl.ReturnType.Name {
		case "string", StringTypeName:
			returnType = TypeString
		case TypeInt, IntTypeName:
			returnType = TypeInt
		case "bool", "Bool":
			returnType = TypeBool
		case "any":
			returnType = TypeAny
		case TypeUnit:
			returnType = TypeUnit
		case TypeHTTPResponse:
			returnType = TypeHTTPResponse
		default:
			// Check if it's a user-defined union type
			if _, exists := g.typeDeclarations[fnDecl.ReturnType.Name]; exists {
				returnType = fnDecl.ReturnType.Name
			} else {
				returnType = TypeInt // Default fallback
			}
		}
	} else {
		// Set up parameter types for analysis
		g.setParameterTypesForAnalysis(fnDecl)

		returnType = g.analyzeReturnType(fnDecl.Body)

		// CRITICAL RULE: Functions cannot implicitly return 'any' type
		if returnType == TypeAny {
			// Clear parameter types before returning error
			g.clearParameterTypesForAnalysis()

			return WrapImplicitAnyReturn(fnDecl.Name)
		}

		// Clear parameter types after successful analysis
		g.clearParameterTypesForAnalysis()
	}

	params := g.createFunctionParameters(fnDecl)
	llvmReturnType := g.getLLVMReturnType(returnType, fnDecl.Name)

	fn := g.module.NewFunc(fnDecl.Name, llvmReturnType, params...)
	g.functions[fnDecl.Name] = fn
	g.functionReturnTypes[fnDecl.Name] = returnType
	g.storeFunctionParameterNames(fnDecl)

	return nil
}

// createFunctionParameters creates LLVM parameters for a function declaration.
func (g *LLVMGenerator) createFunctionParameters(fnDecl *ast.FunctionDeclaration) []*ir.Param {
	params := make([]*ir.Param, len(fnDecl.Parameters))

	for i, param := range fnDecl.Parameters {
		var paramType string

		// Use explicit parameter type if provided
		if param.Type != nil {
			// Check if this is a function type
			if param.Type.IsFunction {
				paramType = TypeFunction
			} else {
				switch param.Type.Name {
				case TypeString, StringTypeName:
					paramType = TypeString
				case TypeInt, IntTypeName:
					paramType = TypeInt
				case TypeBool, BoolTypeName:
					paramType = TypeBool
				default:
					// Check if it's a user-defined union type
					if _, exists := g.typeDeclarations[param.Type.Name]; exists {
						paramType = TypeInt // Union types are represented as integers
					} else {
						paramType = TypeInt // Default fallback
					}
				}
			}
		} else {
			// Try to infer type from usage in function body
			inferredType := g.analyzeParameterUsage(param.Name, fnDecl.Body)
			if inferredType != "" && inferredType != TypeAny {
				paramType = inferredType
			} else {
				// If we can't infer, default to int
				paramType = TypeInt
			}
		}

		llvmParamType := g.getLLVMParameterType(paramType, param.Type)
		params[i] = ir.NewParam(param.Name, llvmParamType)
	}

	return params
}

// getLLVMParameterType converts a parameter type string to LLVM type.
func (g *LLVMGenerator) getLLVMParameterType(paramType string, paramTypeExpr *ast.TypeExpression) types.Type {
	if paramType == TypeString {
		return types.I8Ptr
	}

	if paramType == TypeFunction && paramTypeExpr != nil {
		// Use the type expression to LLVM type conversion for function types
		return g.typeExpressionToLLVMType(paramTypeExpr)
	}

	return types.I64
}

// getLLVMReturnType converts a return type string to LLVM type and stores it.
func (g *LLVMGenerator) getLLVMReturnType(returnType, functionName string) types.Type {
	// Special case: main function must return i32 for C runtime compatibility
	if functionName == MainFunctionName {
		g.functionReturnTypes[functionName] = TypeInt

		return types.I32
	}

	if returnType == TypeString {
		g.functionReturnTypes[functionName] = TypeString

		return types.I8Ptr
	}

	if returnType == TypeUnit {
		g.functionReturnTypes[functionName] = TypeUnit

		return types.Void
	}

	if returnType == TypeBool {
		g.functionReturnTypes[functionName] = TypeBool

		return types.I64 // Use i64 for consistency, but could use i1
	}

	if returnType == TypeAny {
		g.functionReturnTypes[functionName] = TypeAny

		return types.I64 // any types are represented as i64 at LLVM level
	}

	if returnType == TypeHTTPResponse {
		g.functionReturnTypes[functionName] = TypeHTTPResponse

		// Return pointer to HttpResponse struct
		return types.NewPointer(g.typeMap[TypeHTTPResponse])
	}

	g.functionReturnTypes[functionName] = TypeInt

	return types.I64
}

// storeFunctionParameterNames stores parameter names for named argument support.
func (g *LLVMGenerator) storeFunctionParameterNames(fnDecl *ast.FunctionDeclaration) {
	paramNames := make([]string, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		paramNames[i] = param.Name
	}
	g.functionParameters[fnDecl.Name] = paramNames
}

func (g *LLVMGenerator) declareType(typeDecl *ast.TypeDeclaration) {
	// Store the complete type declaration
	g.typeDeclarations[typeDecl.Name] = typeDecl

	// For now, just map the type name to i64 (simplified)
	// In a full implementation, you'd create proper struct types for variants
	g.typeMap[typeDecl.Name] = types.I64

	// Register each variant as a constant with a discriminant value
	for i, variant := range typeDecl.Variants {
		discriminantValue := int64(i)
		g.unionVariants[variant.Name] = discriminantValue
	}
}

// CheckProtectedFunction checks if a function name is protected (built-in).
func CheckProtectedFunction(fnDecl *ast.FunctionDeclaration) error {
	switch fnDecl.Name {
	case PrintFunc, InputFunc, RangeFunc, ForEachFunc, MapFunc, FilterFunc, FoldFunc,
		LengthFunc, ContainsFunc, SubstringFunc:
		return WrapBuiltInRedefine(fnDecl.Name)
	default:
		return nil
	}
}

// setParameterTypesForAnalysis temporarily stores parameter types for return type analysis.
func (g *LLVMGenerator) setParameterTypesForAnalysis(fnDecl *ast.FunctionDeclaration) {
	if g.currentFunctionParameterTypes == nil {
		g.currentFunctionParameterTypes = make(map[string]string)
	}

	for _, param := range fnDecl.Parameters {
		if param.Type != nil {
			// Explicit type annotation - match all possible type name variants
			switch param.Type.Name {
			case "int", IntTypeName:
				g.currentFunctionParameterTypes[param.Name] = TypeInt
			case "string", StringTypeName:
				g.currentFunctionParameterTypes[param.Name] = TypeString
			case TypeBool, BoolTypeName:
				g.currentFunctionParameterTypes[param.Name] = TypeBool
			case "any":
				g.currentFunctionParameterTypes[param.Name] = TypeAny
			default:
				// Check if it's a user-defined union type
				if _, exists := g.typeDeclarations[param.Type.Name]; exists {
					g.currentFunctionParameterTypes[param.Name] = TypeInt // Union types are represented as integers
				} else {
					g.currentFunctionParameterTypes[param.Name] = TypeInt // Default fallback
				}
			}
		} else {
			// Try to infer type from usage in function body
			inferredType := g.analyzeParameterUsage(param.Name, fnDecl.Body)
			if inferredType != "" && inferredType != TypeAny {
				g.currentFunctionParameterTypes[param.Name] = inferredType
			} else {
				// If we can't infer, mark as 'any' for now
				g.currentFunctionParameterTypes[param.Name] = TypeAny
			}
		}
	}
}

// clearParameterTypesForAnalysis clears temporary parameter types after analysis.
func (g *LLVMGenerator) clearParameterTypesForAnalysis() {
	g.currentFunctionParameterTypes = nil
}
