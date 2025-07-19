package codegen

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

func (g *LLVMGenerator) declareFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	if fnDecl.Name == ToStringFunc {
		return ErrToStringReserved
	}

	state := g.saveTypeInferenceState()

	paramTypes := g.inferParameterTypesForSignature(fnDecl)

	returnTypeVar := g.determineReturnTypeForSignature(fnDecl)

	err := g.unifyBodyWithReturnType(fnDecl, paramTypes, returnTypeVar, state)
	if err != nil {
		return err
	}

	finalFnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: g.typeInferer.prune(returnTypeVar),
	}

	g.restoreTypeInferenceState(state)
	g.typeInferer.env.Set(fnDecl.Name, finalFnType)

	err = g.createLLVMFunctionSignature(fnDecl, finalFnType)
	if err != nil {
		return err
	}

	return nil
}

// typeInferenceState holds the state that needs to be saved and restored
type typeInferenceState struct {
	env     *TypeEnv
	subst   Substitution
	effects []string
}

// saveTypeInferenceState saves the current type inference state
func (g *LLVMGenerator) saveTypeInferenceState() *typeInferenceState {
	oldSubst := make(Substitution)
	for k, v := range g.typeInferer.subst {
		oldSubst[k] = v
	}

	var oldEffects []string
	if g.effectCodegen != nil {
		oldEffects = g.effectCodegen.currentFunctionEffects
	}

	return &typeInferenceState{
		env:     g.typeInferer.env,
		subst:   oldSubst,
		effects: oldEffects,
	}
}

// restoreTypeInferenceState restores the saved type inference state
func (g *LLVMGenerator) restoreTypeInferenceState(state *typeInferenceState) {
	g.typeInferer.env = state.env
	g.typeInferer.subst = state.subst
	if g.effectCodegen != nil {
		g.effectCodegen.currentFunctionEffects = state.effects
	}
}

// inferParameterTypesForSignature infers the parameter types for function signature
func (g *LLVMGenerator) inferParameterTypesForSignature(fnDecl *ast.FunctionDeclaration) []Type {
	newEnv := g.typeInferer.env.Clone()
	g.typeInferer.env = newEnv

	paramTypes := make([]Type, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		if param.Type != nil {
			var paramType Type
			if param.Type.IsFunction {
				paramType = g.buildFunctionTypeFromAST(param.Type)
			} else {
				paramType = &ConcreteType{name: param.Type.Name}
			}
			g.typeInferer.env.Set(param.Name, paramType)
			paramTypes[i] = paramType
		} else {
			paramType := g.typeInferer.Fresh()
			g.typeInferer.env.Set(param.Name, paramType)
			paramTypes[i] = paramType
		}
	}

	return paramTypes
}

// determineReturnTypeForSignature determines the return type for function signature
func (g *LLVMGenerator) determineReturnTypeForSignature(fnDecl *ast.FunctionDeclaration) Type {
	if fnDecl.ReturnType != nil {
		return &ConcreteType{name: fnDecl.ReturnType.Name}
	}
	return g.typeInferer.Fresh()
}

// unifyBodyWithReturnType performs type inference on the function body and unifies with return type
func (g *LLVMGenerator) unifyBodyWithReturnType(
	fnDecl *ast.FunctionDeclaration,
	paramTypes []Type,
	returnTypeVar Type,
	state *typeInferenceState,
) error {
	fnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: returnTypeVar,
	}
	g.typeInferer.env.Set(fnDecl.Name, fnType)

	if g.effectCodegen != nil && len(fnDecl.Effects) > 0 {
		g.effectCodegen.currentFunctionEffects = fnDecl.Effects
	}

	bodyType, err := g.typeInferer.InferType(fnDecl.Body)
	if err != nil {
		g.restoreTypeInferenceState(state)
		return err
	}

	if !g.canImplicitlyConvert(bodyType, returnTypeVar, fnDecl) {
		if err := g.typeInferer.Unify(returnTypeVar, bodyType); err != nil {
			g.restoreTypeInferenceState(state)
			var positionInfo string
			if fnDecl.Position != nil {
				positionInfo = fmt.Sprintf(" at line %d, column %d", fnDecl.Position.Line, fnDecl.Position.Column)
			}
			return fmt.Errorf("return type mismatch in function '%s'%s: %w", fnDecl.Name, positionInfo, err)
		}
	}

	if g.effectCodegen != nil && len(fnDecl.Effects) > 0 {
		g.effectCodegen.currentFunctionEffects = state.effects
	}

	return nil
}

// createLLVMFunctionSignature creates the LLVM function signature
func (g *LLVMGenerator) createLLVMFunctionSignature(fnDecl *ast.FunctionDeclaration, finalFnType *FunctionType) error {
	llvmReturnType := g.getLLVMType(finalFnType.returnType)

	if fnDecl.Name == MainFunctionName {
		llvmReturnType = types.I32
	}

	params := make([]*ir.Param, len(finalFnType.paramTypes))
	for i, paramType := range finalFnType.paramTypes {
		params[i] = ir.NewParam(fnDecl.Parameters[i].Name, g.getLLVMType(paramType))
	}

	fn := g.module.NewFunc(fnDecl.Name, llvmReturnType, params...)
	g.functions[fnDecl.Name] = fn

	g.functionParameters[fnDecl.Name] = make([]string, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		g.functionParameters[fnDecl.Name][i] = param.Name
	}

	return nil
}

// generateFunctionDeclaration generates the LLVM function signature and body.
func (g *LLVMGenerator) generateFunctionDeclaration(fnDecl *ast.FunctionDeclaration) error {
	fn, exists := g.functions[fnDecl.Name]
	if !exists {
		return fmt.Errorf("%w: %s", ErrFunctionNotDeclared, fnDecl.Name)
	}

	paramTypes, fnEnv := g.setupFunctionEnvironment(fnDecl)

	returnTypeVar := g.determineReturnType(fnDecl)
	finalFnType, err := g.inferAndValidateTypes(fnDecl, paramTypes, returnTypeVar, fnEnv)
	if err != nil {
		return err
	}

	err = g.generateFunctionBody(fnDecl, fn, finalFnType)
	if err != nil {
		return err
	}

	return nil
}

// setupFunctionEnvironment creates the type environment and infers parameter types
func (g *LLVMGenerator) setupFunctionEnvironment(
	fnDecl *ast.FunctionDeclaration,
) ([]Type, *TypeEnv) {
	fnEnv := g.typeInferer.env.Clone()

	paramTypes := make([]Type, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		if param.Type != nil {
			var paramType Type
			if param.Type.IsFunction {
				paramType = g.buildFunctionTypeFromAST(param.Type)
			} else {
				paramType = &ConcreteType{name: param.Type.Name}
			}
			fnEnv.Set(param.Name, paramType)
			paramTypes[i] = paramType
		} else {
			paramType := g.typeInferer.Fresh()
			fnEnv.Set(param.Name, paramType)
			paramTypes[i] = paramType
		}
	}

	return paramTypes, fnEnv
}

// determineReturnType gets the return type from annotation or creates fresh type variable
func (g *LLVMGenerator) determineReturnType(fnDecl *ast.FunctionDeclaration) Type {
	if fnDecl.ReturnType != nil {
		return &ConcreteType{name: fnDecl.ReturnType.Name}
	}
	return g.typeInferer.Fresh()
}

// inferAndValidateTypes performs type inference and validation for the function
func (g *LLVMGenerator) inferAndValidateTypes(
	fnDecl *ast.FunctionDeclaration,
	paramTypes []Type,
	returnTypeVar Type,
	fnEnv *TypeEnv,
) (*FunctionType, error) {
	oldEnv := g.typeInferer.env

	fnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: returnTypeVar,
	}
	fnEnv.Set(fnDecl.Name, fnType)
	g.typeInferer.env = fnEnv

	inferredReturnType, err := g.typeInferer.InferType(fnDecl.Body)
	if err != nil {
		g.typeInferer.env = oldEnv
		return nil, fmt.Errorf("failed to infer return type: %w", err)
	}

	if !g.canImplicitlyConvert(inferredReturnType, returnTypeVar, fnDecl) {
		if err := g.typeInferer.Unify(returnTypeVar, inferredReturnType); err != nil {
			g.typeInferer.env = oldEnv
			var positionInfo string
			if fnDecl.Position != nil {
				positionInfo = fmt.Sprintf(" at line %d, column %d", fnDecl.Position.Line, fnDecl.Position.Column)
			}
			return nil, fmt.Errorf("return type mismatch in function '%s'%s: %w", fnDecl.Name, positionInfo, err)
		}
	}

	finalFnType := &FunctionType{
		paramTypes: paramTypes,
		returnType: g.typeInferer.prune(returnTypeVar),
	}

	g.typeInferer.env = oldEnv
	g.typeInferer.env.Set(fnDecl.Name, finalFnType)

	return finalFnType, nil
}

// generateFunctionBody generates the LLVM instructions for the function body
func (g *LLVMGenerator) generateFunctionBody(
	fnDecl *ast.FunctionDeclaration,
	fn *ir.Func,
	_ *FunctionType,
) error {
	entry := fn.NewBlock("")
	g.builder = entry
	g.function = fn

	oldExpectedReturnType := g.expectedReturnType
	g.expectedReturnType = fn.Sig.RetType

	params := fn.Params
	for i, param := range fnDecl.Parameters {
		g.variables[param.Name] = params[i]
		// Also add to type inference environment for runtime generation
		if param.Type != nil {
			var paramType Type
			if param.Type.IsFunction {
				paramType = g.buildFunctionTypeFromAST(param.Type)
			} else {
				paramType = &ConcreteType{name: param.Type.Name}
			}
			g.typeInferer.env.Set(param.Name, paramType)
		}
	}

	var oldEffects []string
	if g.effectCodegen != nil && len(fnDecl.Effects) > 0 {
		oldEffects = g.effectCodegen.currentFunctionEffects
		g.effectCodegen.currentFunctionEffects = fnDecl.Effects
	}

	bodyValue, err := g.generateExpression(fnDecl.Body)
	if err != nil {
		g.cleanupAfterError(fnDecl, oldExpectedReturnType, oldEffects)
		return err
	}

	g.cleanupAfterSuccess(fnDecl, oldExpectedReturnType, oldEffects)
	g.generateReturnInstruction(fn, fnDecl, bodyValue)

	return nil
}

// cleanupAfterError cleans up state when an error occurs during function generation
func (g *LLVMGenerator) cleanupAfterError(
	fnDecl *ast.FunctionDeclaration,
	oldExpectedReturnType types.Type,
	oldEffects []string,
) {
	// Don't delete parameters from g.variables - they need to remain for function calls
	g.expectedReturnType = oldExpectedReturnType
	if g.effectCodegen != nil && len(fnDecl.Effects) > 0 {
		g.effectCodegen.currentFunctionEffects = oldEffects
	}
}

// cleanupAfterSuccess cleans up state after successful function generation
func (g *LLVMGenerator) cleanupAfterSuccess(
	fnDecl *ast.FunctionDeclaration,
	oldExpectedReturnType types.Type,
	oldEffects []string,
) {
	// Don't delete parameters from g.variables - they need to remain for function calls
	g.expectedReturnType = oldExpectedReturnType
	if g.effectCodegen != nil && len(fnDecl.Effects) > 0 {
		g.effectCodegen.currentFunctionEffects = oldEffects
	}
}

// generateReturnInstruction generates the appropriate return instruction
func (g *LLVMGenerator) generateReturnInstruction(
	fn *ir.Func,
	fnDecl *ast.FunctionDeclaration,
	bodyValue value.Value,
) {
	if fn.Sig.RetType == types.Void {
		g.builder.NewRet(nil)
	} else if fnDecl.Name == MainFunctionName {
		g.builder.NewRet(constant.NewInt(types.I32, 0))
	} else {
		finalReturnValue := g.maybeWrapInResult(bodyValue, fnDecl)
		g.builder.NewRet(finalReturnValue)
	}
}

// maybeWrapInResult wraps a plain value in a Result structure if the function declares a Result return type
func (g *LLVMGenerator) maybeWrapInResult(bodyValue value.Value, fnDecl *ast.FunctionDeclaration) value.Value {
	// Check if function declares a Result return type
	if fnDecl.ReturnType != nil && strings.HasPrefix(fnDecl.ReturnType.Name, "Result<") {
		// Check if the body value is a plain int and function expects Result<int, MathError>
		if fnDecl.ReturnType.Name == "Result<int, MathError>" && bodyValue.Type() == types.I64 {
			return g.wrapInMathResult(bodyValue)
		}
		// Add other Result type mappings as needed
	}

	// No wrapping needed, return original value
	return bodyValue
}

// wrapInMathResult wraps a plain int value in a Result<int, MathError> structure
func (g *LLVMGenerator) wrapInMathResult(intValue value.Value) value.Value {
	// Create Result<int, MathError> structure
	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)

	// Store the int value in the success field
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(intValue, valuePtr)

	// Store success discriminant (0 = Success)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	// Return pointer to the Result structure
	return result
}

// canImplicitlyConvert checks if we can implicitly convert from one type to another
func (g *LLVMGenerator) canImplicitlyConvert(fromType, toType Type, fnDecl *ast.FunctionDeclaration) bool {
	// Check if we're trying to convert primitive types to Result types
	if fnDecl.ReturnType != nil && fnDecl.ReturnType.Name == TypeResult {
		if len(fnDecl.ReturnType.GenericParams) >= TwoArgs {
			expectedInnerType := fnDecl.ReturnType.GenericParams[0].Name
			expectedErrorType := fnDecl.ReturnType.GenericParams[1].Name

			// Check if it's Result<int, MathError> or Result<bool, MathError>
			if (expectedInnerType == "int" || expectedInnerType == "bool") &&
				expectedErrorType == "MathError" {

				if fromConcrete, ok := fromType.(*ConcreteType); ok {
					if fromConcrete.name == expectedInnerType {
						if toConcrete, ok := toType.(*ConcreteType); ok {
							return toConcrete.name == "Result"
						}
					}
				}
			}
		}
	}

	// Add other implicit conversions as needed
	return false
}

// getLLVMType converts our type system types to LLVM types
func (g *LLVMGenerator) getLLVMType(ospreyType Type) types.Type {
	switch t := ospreyType.(type) {
	case *PrimitiveType:
		return g.getLLVMPrimitiveType(t)
	case *ConcreteType:
		return g.getLLVMConcreteType(t)
	case *GenericType:
		return g.getLLVMGenericType(t)
	case *FunctionType:
		return g.getLLVMFunctionType(t)
	case *RecordType:
		return g.getLLVMRecordType(t)
	case *UnionType:
		return g.getLLVMUnionType(t)
	case *TypeVar:
		// Type variables should be resolved before code generation
		// For now, default to i64 as a fallback
		return types.I64
	default:
		// Unknown type, fallback to i64
		return types.I64
	}
}

// getLLVMPrimitiveType converts primitive types to LLVM types
func (g *LLVMGenerator) getLLVMPrimitiveType(pt *PrimitiveType) types.Type {
	switch pt.name {
	case TypeInt:
		return types.I64
	case TypeString:
		return types.I8Ptr
	case TypeBool:
		return types.I1
	case TypeUnit:
		return types.Void
	default:
		return types.I64
	}
}

// getLLVMConcreteType converts concrete types to LLVM types
func (g *LLVMGenerator) getLLVMConcreteType(ct *ConcreteType) types.Type {
	switch ct.name {
	case TypeInt:
		return types.I64
	case TypeString:
		return types.I8Ptr
	case TypeBool:
		return types.I1
	case TypeUnit:
		return types.Void
	case TypeHTTPResponse:
		// Return pointer to HttpResponse struct
		if httpResponseType, exists := g.typeMap[TypeHTTPResponse]; exists {
			return types.NewPointer(httpResponseType)
		}
		return types.I64 // fallback
	case "Fiber":
		return types.I64 // Fiber handle is represented as i64
	case "Channel":
		return types.I64 // Channel handle is represented as i64
	case "ProcessHandle":
		return types.I64 // Process handle is represented as i64
	case "any":
		return types.I64 // any type is represented as i64
	default:
		// Handle Result types like "Result<int, MathError>"
		if strings.HasPrefix(ct.name, "Result<") {
			// Result types are represented as structs with { value, discriminant }
			if ct.name == "Result<int, MathError>" {
				return types.NewPointer(g.getResultType(types.I64))
			}
			// Add other Result type mappings as needed
		}

		// Handle function types like "fn(int, int, string) -> Unit" or "(int) -> int"
		if strings.HasPrefix(ct.name, "fn(") && strings.Contains(ct.name, ") -> ") {
			// For now, represent function types as function pointers (i8*)
			return types.I8Ptr
		}
		if strings.HasPrefix(ct.name, "(") && strings.Contains(ct.name, ") -> ") {
			// For now, represent function types as function pointers (i8*)
			return types.I8Ptr
		}

		// Check if it's a user-defined type
		if userType, exists := g.typeMap[ct.name]; exists {
			// For struct types, return pointer to the struct
			if _, ok := userType.(*types.StructType); ok {
				return types.NewPointer(userType)
			}
			return userType
		}

		// For now, default to i64 for unknown types
		return types.I64
	}
}

// getLLVMGenericType converts generic types to LLVM types
func (g *LLVMGenerator) getLLVMGenericType(gt *GenericType) types.Type {
	switch gt.name {
	case "Result":
		if len(gt.typeArgs) >= TwoArgs {
			// Result<T, E> - get the inner type for the value
			innerType := g.getLLVMType(gt.typeArgs[0])
			return types.NewPointer(g.getResultType(innerType))
		}
		return types.I64 // fallback
	case "List":
		if len(gt.typeArgs) >= 1 {
			// List<T> - for now, represent as pointer to dynamic array
			return types.I8Ptr // TODO: implement proper list types
		}
		return types.I8Ptr // fallback
	default:
		// Unknown generic type, fallback to i64
		return types.I64
	}
}

// getLLVMFunctionType converts function types to LLVM function pointer types
func (g *LLVMGenerator) getLLVMFunctionType(ft *FunctionType) types.Type {
	// Convert parameter types
	paramTypes := make([]types.Type, len(ft.paramTypes))
	for i, paramType := range ft.paramTypes {
		paramTypes[i] = g.getLLVMType(paramType)
	}

	// Convert return type
	returnType := g.getLLVMType(ft.returnType)

	// Create function signature and return pointer to it
	funcSignature := types.NewFunc(returnType, paramTypes...)
	return types.NewPointer(funcSignature)
}

// getLLVMRecordType converts record types to LLVM struct types
func (g *LLVMGenerator) getLLVMRecordType(rt *RecordType) types.Type {
	// Check if we already have this record type in the type map
	if llvmType, exists := g.typeMap[rt.name]; exists {
		// For struct types, return pointer to the struct
		if _, ok := llvmType.(*types.StructType); ok {
			return types.NewPointer(llvmType)
		}
		return llvmType
	}

	// Create struct type from record fields
	fieldTypes := make([]types.Type, 0, len(rt.fields))
	for _, fieldType := range rt.fields {
		fieldTypes = append(fieldTypes, g.getLLVMType(fieldType))
	}

	structType := types.NewStruct(fieldTypes...)
	g.typeMap[rt.name] = structType

	return types.NewPointer(structType)
}

// getLLVMUnionType converts union types to LLVM tagged union types
func (g *LLVMGenerator) getLLVMUnionType(ut *UnionType) types.Type {
	// Check if we already have this union type in the type map
	if llvmType, exists := g.typeMap[ut.name]; exists {
		return types.NewPointer(llvmType)
	}

	// For now, represent union types as tagged unions with discriminant
	// This is a simplified implementation - full implementation would need
	// proper analysis of variant sizes
	tagType := types.I8
	dataType := types.I64 // simplified data representation
	unionType := types.NewStruct(tagType, dataType)

	g.typeMap[ut.name] = unionType

	return types.NewPointer(unionType)
}

// buildFunctionTypeFromAST converts an AST TypeExpression with IsFunction=true to a FunctionType
func (g *LLVMGenerator) buildFunctionTypeFromAST(typeExpr *ast.TypeExpression) Type {
	if !typeExpr.IsFunction {
		return &ConcreteType{name: typeExpr.Name}
	}

	// Convert parameter types
	var paramTypes []Type
	for _, paramType := range typeExpr.ParameterTypes {
		if paramType.IsFunction {
			// Recursive function types
			paramTypes = append(paramTypes, g.buildFunctionTypeFromAST(&paramType))
		} else {
			paramTypes = append(paramTypes, &ConcreteType{name: paramType.Name})
		}
	}

	// Convert return type
	var returnType Type
	if typeExpr.ReturnType != nil {
		if typeExpr.ReturnType.IsFunction {
			returnType = g.buildFunctionTypeFromAST(typeExpr.ReturnType)
		} else {
			returnType = &ConcreteType{name: typeExpr.ReturnType.Name}
		}
	} else {
		returnType = &ConcreteType{name: "Unit"} // default return type
	}

	return &FunctionType{
		paramTypes: paramTypes,
		returnType: returnType,
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

		// Register the record type in the type inference environment
		recordType := &ConcreteType{name: typeDecl.Name}
		g.typeInferer.env.Set(typeDecl.Name, recordType)

		// Also register the variant name in the type environment
		// This allows type constructors like Warrior { health: 100, strength: 20 } to work
		g.typeInferer.env.Set(variant.Name, recordType)

		// Register the variant name in unionVariants for pattern matching
		// This allows pattern matching to work correctly for single-variant types with fields
		g.unionVariants[variant.Name] = 0 // Use 0 as discriminant for single-variant types

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

	// Register the union type in the type inference environment
	unionTypeInference := &ConcreteType{name: typeDecl.Name}
	g.typeInferer.env.Set(typeDecl.Name, unionTypeInference)

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

// CheckProtectedFunction checks if a function name is protected (built-in) using the unified registry.
func CheckProtectedFunction(fnDecl *ast.FunctionDeclaration) error {
	if GlobalBuiltInRegistry.IsProtectedFunction(fnDecl.Name) {
		return WrapBuiltInRedefine(fnDecl.Name)
	}
	return nil
}

// getFieldType converts an Osprey field type to LLVM type
func (g *LLVMGenerator) getFieldType(fieldType string) types.Type {
	switch fieldType {
	case TypeString, "String": // "string" or "String"
		return types.I8Ptr
	case TypeInt, "Int": // "int" or "Int"
		return types.I64
	case TypeBool, "Bool": // "bool" or "Bool"
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
