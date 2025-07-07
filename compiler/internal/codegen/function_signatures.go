package codegen

import (
	"strings"

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
			// Check if it's a Result type (generic type)
			if strings.HasPrefix(fnDecl.ReturnType.Name, "Result<") {
				returnType = fnDecl.ReturnType.Name // Keep full Result<T, E> type
			} else if _, exists := g.typeDeclarations[fnDecl.ReturnType.Name]; exists {
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

	// Check if this is a Result type
	if isResultType(returnType) {
		g.functionReturnTypes[functionName] = returnType
		return g.getResultTypeForSignature(returnType)
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

		return types.I1
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

// isResultType checks if the return type is a Result type
func isResultType(returnType string) bool {
	return strings.HasPrefix(returnType, "Result<")
}

// getResultTypeForSignature returns the LLVM type for a Result type signature
func (g *LLVMGenerator) getResultTypeForSignature(resultType string) types.Type {
	// Parse the Result<T, E> type to get the inner type
	if strings.Contains(resultType, "Result<Unit,") {
		return g.getResultType(types.Void)
	}
	if strings.Contains(resultType, "Result<string,") {
		return g.getResultType(types.I8Ptr)
	}
	if strings.Contains(resultType, "Result<bool,") {
		return g.getResultType(types.I1)
	}
	// Default to Result<i64, Error>
	return g.getResultType(types.I64)
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

	// Check if this is a record type (single variant with fields)
	if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
		// Create a proper struct type for record types
		variant := typeDecl.Variants[0]
		fieldTypes := make([]types.Type, len(variant.Fields))

		for i, field := range variant.Fields {
			fieldTypes[i] = g.getFieldType(field.Type)
		}

		// Create the struct type
		structType := types.NewStruct(fieldTypes...)
		g.typeMap[typeDecl.Name] = structType

		// Store field names for field access
		fieldNames := make([]string, len(variant.Fields))
		for i, field := range variant.Fields {
			fieldNames[i] = field.Name
		}
		g.storeRecordFieldNames(typeDecl.Name, fieldNames)
	} else if len(typeDecl.Variants) > 1 {
		// This is a discriminated union with multiple variants
		g.declareDiscriminatedUnion(typeDecl)
	} else {
		// Simple enum (variants without fields)
		g.typeMap[typeDecl.Name] = types.I64

		// Register each variant as a constant with a discriminant value
		for i, variant := range typeDecl.Variants {
			discriminantValue := int64(i)
			g.unionVariants[variant.Name] = discriminantValue
		}
	}
}

// declareDiscriminatedUnion creates a tagged union structure for discriminated unions
func (g *LLVMGenerator) declareDiscriminatedUnion(typeDecl *ast.TypeDeclaration) {
	// For discriminated unions, we need to create:
	// struct { i8 tag, [largest_variant_size x i8] data }

	// Find the largest variant data size
	maxDataSize := int64(0)
	hasFieldVariants := false

	for _, variant := range typeDecl.Variants {
		if len(variant.Fields) > 0 {
			hasFieldVariants = true
			// Calculate size needed for this variant's fields
			variantSize := int64(0)
			for _, field := range variant.Fields {
				fieldType := g.getFieldType(field.Type)
				switch fieldType {
				case types.I64:
					variantSize += 8
				case types.I8Ptr:
					variantSize += 8 // pointer size
				case types.I1:
					variantSize += 1
				default:
					variantSize += 8 // default to 8 bytes
				}
			}
			if variantSize > maxDataSize {
				maxDataSize = variantSize
			}
		}
	}

	if !hasFieldVariants {
		// If no variants have fields, treat as simple enum
		g.typeMap[typeDecl.Name] = types.I64
		for i, variant := range typeDecl.Variants {
			g.unionVariants[variant.Name] = int64(i)
		}
		return
	}

	// Create tagged union: struct { i8 tag, [maxDataSize x i8] data }
	tagType := types.I8
	dataType := types.NewArray(uint64(maxDataSize), types.I8)
	unionType := types.NewStruct(tagType, dataType)

	g.typeMap[typeDecl.Name] = unionType

	// Store variant information for construction and pattern matching
	for i, variant := range typeDecl.Variants {
		g.unionVariants[variant.Name] = int64(i)
		// Store variant field information
		if len(variant.Fields) > 0 {
			fieldNames := make([]string, len(variant.Fields))
			for j, field := range variant.Fields {
				fieldNames[j] = field.Name
			}
			g.storeVariantFieldNames(typeDecl.Name, variant.Name, fieldNames)
		}
	}
}

// storeVariantFieldNames stores field names for discriminated union variants
func (g *LLVMGenerator) storeVariantFieldNames(typeName, variantName string, fieldNames []string) {
	key := typeName + "_" + variantName + "_fields"
	g.functionParameters[key] = fieldNames
}

// getVariantFieldNames retrieves field names for a specific variant
func (g *LLVMGenerator) getVariantFieldNames(typeName, variantName string) ([]string, bool) {
	key := typeName + "_" + variantName + "_fields"
	fieldNames, exists := g.functionParameters[key]
	return fieldNames, exists
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

// getFieldType converts an Osprey field type to LLVM type
func (g *LLVMGenerator) getFieldType(fieldType string) types.Type {
	switch fieldType {
	case "Int":
		return types.I64
	case "String":
		return types.I8Ptr
	case "Bool":
		return types.I1
	default:
		// Check if it's a user-defined type
		if llvmType, exists := g.typeMap[fieldType]; exists {
			return llvmType
		}
		// Default to i64 for unknown types
		return types.I64
	}
}

// storeRecordFieldNames stores field names for record types to enable field access
func (g *LLVMGenerator) storeRecordFieldNames(typeName string, fieldNames []string) {
	// For now, store in the existing functionParameters map with a special key
	// In a full implementation, you'd have a dedicated recordFields map
	g.functionParameters[typeName+"_fields"] = fieldNames
}

// getRecordFieldNames retrieves field names for a record type
func (g *LLVMGenerator) getRecordFieldNames(typeName string) ([]string, bool) {
	fieldNames, exists := g.functionParameters[typeName+"_fields"]
	return fieldNames, exists
}

// findFieldIndex finds the index of a field in a record type
func (g *LLVMGenerator) findFieldIndex(typeName, fieldName string) int {
	if fieldNames, exists := g.getRecordFieldNames(typeName); exists {
		for i, name := range fieldNames {
			if name == fieldName {
				return i
			}
		}
	}
	return -1
}
