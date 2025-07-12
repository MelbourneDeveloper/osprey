package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) declareFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	if fnDecl.Name == ToStringFunc {
		return ErrToStringReserved
	}

	// Save old environment AND substitution map
	oldEnv := g.typeInferer.env
	oldSubst := make(Substitution)
	for k, v := range g.typeInferer.subst {
		oldSubst[k] = v
	}

	// Clone the existing environment to preserve access to previously declared functions
	newEnv := g.typeInferer.env.Clone()
	g.typeInferer.env = newEnv

	// Infer parameter types
	paramTypes := make([]Type, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		// CRITICAL FIX: Use explicit type annotation if present
		if param.Type != nil {
			// Use the explicit type annotation
			paramType := &ConcreteType{name: param.Type.Name}
			g.typeInferer.env.Set(param.Name, paramType)
			paramTypes[i] = paramType
		} else {
			// Create fresh type variable for parameter without explicit type
			paramType := g.typeInferer.Fresh()
			g.typeInferer.env.Set(param.Name, paramType)
			paramTypes[i] = paramType
		}
	}

	// For recursion: Add the function to the environment before processing the body
	// We'll use a fresh type variable for the return type initially
	returnTypeVar := g.typeInferer.Fresh()
	fnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: returnTypeVar,
	}
	g.typeInferer.env.Set(fnDecl.Name, fnType)

	// Infer body type
	bodyType, err := g.typeInferer.InferType(fnDecl.Body)
	if err != nil {
		// Restore environment AND substitution map on error
		g.typeInferer.env = oldEnv
		g.typeInferer.subst = oldSubst
		return err
	}

	// Unify the return type variable with the actual body type
	if err := g.typeInferer.Unify(returnTypeVar, bodyType); err != nil {
		// Restore environment AND substitution map on error
		g.typeInferer.env = oldEnv
		g.typeInferer.subst = oldSubst
		return fmt.Errorf("return type mismatch: %w", err)
	}

	// Create the final function type with the unified return type
	finalFnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: g.typeInferer.prune(returnTypeVar), // Use the unified type
	}

	// Restore original environment and substitution map
	g.typeInferer.env = oldEnv
	g.typeInferer.subst = oldSubst

	// Add function to the clean environment
	g.typeInferer.env.Set(fnDecl.Name, finalFnType)

	// Generate LLVM function signature ONLY (no body)
	llvmReturnType := g.getLLVMType(finalFnType.returnType.String())
	params := make([]*ir.Param, len(paramTypes))
	for i, paramType := range paramTypes {
		params[i] = ir.NewParam(fnDecl.Parameters[i].Name, g.getLLVMType(paramType.String()))
	}

	// Create LLVM function declaration ONLY
	fn := g.module.NewFunc(fnDecl.Name, llvmReturnType, params...)
	g.functions[fnDecl.Name] = fn

	// Store parameter names for named argument support
	g.functionParameters[fnDecl.Name] = make([]string, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		g.functionParameters[fnDecl.Name][i] = param.Name
	}

	return nil
}

// generateFunctionDeclaration generates the LLVM function signature and body.
func (g *LLVMGenerator) generateFunctionDeclaration(fnDecl *ast.FunctionDeclaration) error {
	// Get the already-declared function
	fn, exists := g.functions[fnDecl.Name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrFunctionNotDeclared, fnDecl.Name)
	}

	// Create function environment for parameter type inference
	oldEnv := g.typeInferer.env
	fnEnv := g.typeInferer.env.Clone()

	// Infer parameter types (same as in declareFunctionSignature)
	paramTypes := make([]Type, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		// CRITICAL FIX: Use explicit type annotation if present
		if param.Type != nil {
			// Use the explicit type annotation
			paramType := &ConcreteType{name: param.Type.Name}
			fnEnv.Set(param.Name, paramType)
			paramTypes[i] = paramType
		} else {
			// Create fresh type variable for parameter without explicit type
			paramType := g.typeInferer.Fresh()
			fnEnv.Set(param.Name, paramType)
			paramTypes[i] = paramType
		}
	}

	// For recursion: Add the function to the environment before processing the body
	returnTypeVar := g.typeInferer.Fresh()
	fnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: returnTypeVar,
	}
	fnEnv.Set(fnDecl.Name, fnType)

	// Set the function environment with parameters for body inference
	g.typeInferer.env = fnEnv

	// Infer the return type in the context with parameters and the function itself
	inferredReturnType, err := g.typeInferer.InferType(fnDecl.Body)
	if err != nil {
		g.typeInferer.env = oldEnv
		return fmt.Errorf("failed to infer return type: %w", err)
	}

	// Unify the return type variable with the actual body type
	if err := g.typeInferer.Unify(returnTypeVar, inferredReturnType); err != nil {
		g.typeInferer.env = oldEnv
		return fmt.Errorf("return type mismatch: %w", err)
	}

	// Create final function type
	finalFnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: g.typeInferer.prune(returnTypeVar),
	}

	// Generate function body (keep the function environment with parameters and function itself)
	entry := fn.NewBlock("")
	g.builder = entry
	g.function = fn

	// Set expected return type context for literals and expressions
	oldExpectedReturnType := g.expectedReturnType
	g.expectedReturnType = fn.Sig.RetType

	// Set up parameter values in runtime environment
	// (type environment already has them from fnEnv)
	params := fn.Params
	for i, param := range fnDecl.Parameters {
		g.variables[param.Name] = params[i]
	}

	bodyValue, err := g.generateExpression(fnDecl.Body)
	if err != nil {
		// Clean up before returning error
		for _, param := range fnDecl.Parameters {
			delete(g.variables, param.Name)
		}
		g.expectedReturnType = oldExpectedReturnType
		g.typeInferer.env = oldEnv
		return err
	}

	// Clean up parameter variables from runtime environment
	for _, param := range fnDecl.Parameters {
		delete(g.variables, param.Name)
	}

	// Restore original environment and add function to it
	g.typeInferer.env = oldEnv
	g.typeInferer.env.Set(fnDecl.Name, finalFnType)

	// Restore expected return type context
	g.expectedReturnType = oldExpectedReturnType

	// Add return instruction
	g.builder.NewRet(bodyValue)

	return nil
}

// getLLVMType converts our type system types to LLVM types
func (g *LLVMGenerator) getLLVMType(typeName string) types.Type {
	switch typeName {
	case TypeInt: // "int"
		return types.I64
	case TypeString: // "string"
		return types.I8Ptr
	case TypeBool: // "bool"
		return types.I1
	default:
		// For now, default to i64 for unknown types
		return types.I64
	}
}

// declareType processes type declarations and sets up LLVM structures
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
			// Add variant to the type environment as a concrete type
			variantType := &ConcreteType{name: typeDecl.Name}
			g.typeInferer.env.Set(variant.Name, variantType)
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
					variantSize++
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
			// Add variant to the type environment as a concrete type
			variantType := &ConcreteType{name: typeDecl.Name}
			g.typeInferer.env.Set(variant.Name, variantType)
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
		// Add variant to the type environment as a concrete type
		variantType := &ConcreteType{name: typeDecl.Name}
		g.typeInferer.env.Set(variant.Name, variantType)
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

// getFieldType converts an Osprey field type to LLVM type
func (g *LLVMGenerator) getFieldType(fieldType string) types.Type {
	switch fieldType {
	case TypeString: // "string"
		return types.I8Ptr
	case TypeInt: // "int"
		return types.I64
	case TypeBool: // "bool"
		return types.I1
	default:
		return types.I64 // default to i64
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
