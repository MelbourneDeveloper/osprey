package codegen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// Type mangling constants for generating unique function names
const (
	TypeManglingInt64   = "i64"
	TypeManglingStr     = "str"
	TypeManglingBool    = "bool"
	TypeManglingI32     = "i32"
	TypeManglingVoid    = "void"
	TypeManglingFn      = "fn"
	TypeManglingUnknown = "unknown"
)

func (g *LLVMGenerator) declareFunctionSignature(fnDecl *ast.FunctionDeclaration) error {
	if fnDecl.Name == ToStringFunc {
		return ErrToStringReserved
	}

	// Store function declaration for monomorphization
	if g.functionDeclarations == nil {
		g.functionDeclarations = make(map[string]*ast.FunctionDeclaration)
	}

	g.functionDeclarations[fnDecl.Name] = fnDecl

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

	// HINDLEY-MILNER: Generalize the function type to create a type scheme
	// This allows the function to be polymorphic across multiple call sites
	scheme := g.typeInferer.Generalize(finalFnType)
	g.typeInferer.env.Set(fnDecl.Name, scheme)

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
				paramType = g.typeExpressionToInferenceType(param.Type)
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
		// Use proper type conversion instead of hardcoded ConcreteType strings
		return g.typeExpressionToInferenceType(fnDecl.ReturnType)
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

	// CRUCIAL: After type inference, update parameter types which may now be constrained
	for i, param := range fnDecl.Parameters {
		if param.Type == nil { // Only for parameters without explicit types
			if paramType, exists := g.typeInferer.env.Get(param.Name); exists {
				prunedType := g.typeInferer.prune(paramType)
				g.typeInferer.env.Set(param.Name, prunedType)
				// Update the parameter type in the function's paramTypes
				if i < len(paramTypes) {
					paramTypes[i] = prunedType
				}
			}
		}
	}

	if !g.canImplicitlyConvert(bodyType, returnTypeVar, fnDecl) {
		err := g.typeInferer.Unify(returnTypeVar, bodyType)
		if err != nil {
			g.restoreTypeInferenceState(state)

			var positionInfo string
			if fnDecl.Position != nil {
				positionInfo = fmt.Sprintf(" at line %d, column %d", fnDecl.Position.Line, fnDecl.Position.Column)
			}

			return fmt.Errorf("return type mismatch in function '%s'%s: body type=%s (%T), return type=%s (%T), error: %w",
				fnDecl.Name, positionInfo, bodyType.String(), bodyType, returnTypeVar.String(), returnTypeVar, err)
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

	// HINDLEY-MILNER FIX: During initial declaration, use original name for all functions
	// Monomorphization will happen later during function calls
	mangledName := fnDecl.Name

	params := make([]*ir.Param, len(finalFnType.paramTypes))
	for i, paramType := range finalFnType.paramTypes {
		params[i] = ir.NewParam(fnDecl.Parameters[i].Name, g.getLLVMType(paramType))
	}

	fn := g.module.NewFunc(mangledName, llvmReturnType, params...)
	g.functions[mangledName] = fn

	// HINDLEY-MILNER FIX: During initial declaration, always store under original name
	// This allows generateFunctionDeclaration to find the function later
	// Monomorphized instances will be created on-demand with different names
	g.functions[fnDecl.Name] = fn

	g.functionParameters[mangledName] = make([]string, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		g.functionParameters[mangledName][i] = param.Name
	}

	return nil
}

// getMonomorphizedName creates a unique mangled name for monomorphized function instances
func (g *LLVMGenerator) getMonomorphizedName(baseName string, fnType *FunctionType) string {
	// For non-polymorphic or first instance, use original name
	if !g.isPolymorphicFunction(baseName, fnType) {
		return baseName
	}

	// Create a type signature for mangling
	typeSignature := g.createTypeSignature(fnType)
	mangledName := fmt.Sprintf("%s_%s", baseName, typeSignature)

	// Track this monomorphization
	if g.monomorphizedInstances == nil {
		g.monomorphizedInstances = make(map[string]string)
	}

	g.monomorphizedInstances[typeSignature] = mangledName

	return mangledName
}

// isPolymorphicFunction checks if a function needs monomorphization
func (g *LLVMGenerator) isPolymorphicFunction(name string, fnType *FunctionType) bool {
	// HINDLEY-MILNER FIX: Check if function has type variables or concrete types that differ from base
	// Always monomorphize if there are already monomorphized instances
	for existingName := range g.functions {
		if strings.HasPrefix(existingName, name+"_") {
			return true
		}
	}

	// Check if this function has a type scheme in the environment (polymorphic)
	if funcTypeFromEnv, exists := g.typeInferer.env.Get(name); exists {
		if _, isScheme := funcTypeFromEnv.(*TypeScheme); isScheme {
			return true
		}
	}

	// Check if we already have an instance of this function with a different signature
	if existingFn, exists := g.functions[name]; exists {
		// Compare signatures - if they differ, this is a polymorphic instantiation
		existingSignature := g.getFunctionTypeSignature(existingFn)
		newSignature := g.createTypeSignature(fnType)

		return existingSignature != newSignature
	}

	return false
}

// createTypeSignature creates a string representation of a function type for mangling
func (g *LLVMGenerator) createTypeSignature(fnType *FunctionType) string {
	var parts []string

	// Add parameter types
	for _, paramType := range fnType.paramTypes {
		parts = append(parts, g.getTypeString(paramType))
	}

	// Add return type
	parts = append(parts, g.getTypeString(fnType.returnType))

	return strings.Join(parts, "_")
}

// getTypeString gets a short string representation of a type for mangling
func (g *LLVMGenerator) getTypeString(t Type) string {
	switch typ := t.(type) {
	case *PrimitiveType:
		switch typ.name {
		case TypeInt:
			return TypeManglingInt64
		case TypeString:
			return TypeManglingStr
		case TypeBool:
			return TypeManglingBool
		default:
			return typ.name
		}
	case *ConcreteType:
		switch typ.name {
		case TypeInt:
			return TypeManglingInt64
		case TypeString:
			return TypeManglingStr
		case TypeBool:
			return TypeManglingBool
		default:
			return typ.name
		}
	case *FunctionType:
		return TypeManglingFn
	default:
		return TypeManglingUnknown
	}
}

// getFunctionTypeSignature extracts type signature from an existing LLVM function
func (g *LLVMGenerator) getFunctionTypeSignature(fn *ir.Func) string {
	var parts []string

	// Add parameter types
	for _, param := range fn.Params {
		parts = append(parts, g.getLLVMTypeString(param.Type()))
	}

	// Add return type
	parts = append(parts, g.getLLVMTypeString(fn.Sig.RetType))

	return strings.Join(parts, "_")
}

// getLLVMTypeString gets a short string representation of an LLVM type
func (g *LLVMGenerator) getLLVMTypeString(t types.Type) string {
	switch t {
	case types.I64:
		return TypeManglingInt64
	case types.I8Ptr:
		return TypeManglingStr
	case types.I1:
		return TypeManglingBool
	case types.I32:
		return TypeManglingI32
	case types.Void:
		return TypeManglingVoid
	default:
		return TypeManglingUnknown
	}
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
				// Check if this is a known record type first
				if typeDecl, exists := g.typeDeclarations[param.Type.Name]; exists {
					// Check if this is a record type (single variant with fields)
					if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
						variant := typeDecl.Variants[0]
						// Create a proper RecordType for the type inference environment
						recordFieldTypes := make(map[string]Type)

						fieldOrder := make([]string, 0, len(variant.Fields))

						for _, field := range variant.Fields {
							recordFieldTypes[field.Name] = g.getInferenceFieldType(field.Type)
							fieldOrder = append(fieldOrder, field.Name)
						}

						paramType = NewOrderedRecordType(typeDecl.Name, recordFieldTypes, fieldOrder)
					} else {
						paramType = g.typeExpressionToInferenceType(param.Type)
					}
				} else {
					paramType = g.typeExpressionToInferenceType(param.Type)
				}
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
		// Use proper type conversion to handle generic types like Result<int, MathError>
		return g.typeExpressionToInferenceType(fnDecl.ReturnType)
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

	// CRUCIAL: After type inference, update parameter types which may now be constrained
	for i, param := range fnDecl.Parameters {
		if param.Type == nil { // Only for parameters without explicit types
			if paramType, exists := g.typeInferer.env.Get(param.Name); exists {
				prunedType := g.typeInferer.prune(paramType)
				g.typeInferer.env.Set(param.Name, prunedType)
				// Update the parameter type in the function's paramTypes
				if i < len(paramTypes) {
					paramTypes[i] = prunedType
				}
			}
		}
	}

	if !g.canImplicitlyConvert(inferredReturnType, returnTypeVar, fnDecl) {
		err := g.typeInferer.Unify(returnTypeVar, inferredReturnType)
		if err != nil {
			g.typeInferer.env = oldEnv

			var positionInfo string
			if fnDecl.Position != nil {
				positionInfo = fmt.Sprintf(" at line %d, column %d", fnDecl.Position.Line, fnDecl.Position.Column)
			}

			return nil, fmt.Errorf("return type mismatch in function '%s'%s: body type=%s (%T), return type=%s (%T), error: %w",
				fnDecl.Name, positionInfo, inferredReturnType.String(), inferredReturnType,
				returnTypeVar.String(), returnTypeVar, err)
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

	// Set up parameters in type inference environment
	g.setupParametersInEnvironment(fnDecl, fn)

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

// setupParametersInEnvironment sets up parameters in the type inference environment
func (g *LLVMGenerator) setupParametersInEnvironment(fnDecl *ast.FunctionDeclaration, fn *ir.Func) {
	params := fn.Params
	for i, param := range fnDecl.Parameters {
		g.variables[param.Name] = params[i]
		// Also add to type inference environment for runtime generation
		if param.Type != nil {
			paramType := g.getParameterType(param.Type)
			g.typeInferer.env.Set(param.Name, paramType)
		} else {
			// HINDLEY-MILNER: For parameters without explicit types,
			// get the type from the function signature that was already inferred
			g.setInferredParameterType(fnDecl.Name, param.Name, i)
		}
	}
}

// getParameterType gets the type for a parameter with explicit type annotation
func (g *LLVMGenerator) getParameterType(paramType *ast.TypeExpression) Type {
	if paramType.IsFunction {
		return g.buildFunctionTypeFromAST(paramType)
	}
	// Check if this is a known record type first
	if typeDecl, exists := g.typeDeclarations[paramType.Name]; exists {
		// Check if this is a record type (single variant with fields)
		if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
			variant := typeDecl.Variants[0]
			// Create a proper RecordType for the type inference environment
			recordFieldTypes := make(map[string]Type)

			fieldOrder := make([]string, 0, len(variant.Fields))

			for _, field := range variant.Fields {
				recordFieldTypes[field.Name] = g.getInferenceFieldType(field.Type)
				fieldOrder = append(fieldOrder, field.Name)
			}

			return NewOrderedRecordType(typeDecl.Name, recordFieldTypes, fieldOrder)
		}
	}

	return g.typeExpressionToInferenceType(paramType)
}

// setInferredParameterType sets the type for a parameter without explicit type annotation
func (g *LLVMGenerator) setInferredParameterType(fnName, paramName string, paramIndex int) {
	if fnType, exists := g.typeInferer.env.Get(fnName); exists {
		// Check if it's a type scheme and instantiate it
		if scheme, ok := fnType.(*TypeScheme); ok {
			instantiated := g.typeInferer.Instantiate(scheme)
			if ft, ok := instantiated.(*FunctionType); ok && paramIndex < len(ft.paramTypes) {
				g.typeInferer.env.Set(paramName, ft.paramTypes[paramIndex])
			}
		} else if ft, ok := fnType.(*FunctionType); ok && paramIndex < len(ft.paramTypes) {
			g.typeInferer.env.Set(paramName, ft.paramTypes[paramIndex])
		}
	}
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
		// Unwrap Result types if function return type is not a Result
		finalReturnValue = g.maybeUnwrapResult(finalReturnValue, fnDecl)
		g.builder.NewRet(finalReturnValue)
	}
}

// maybeWrapInResult wraps a plain value in a Result structure if the function declares a Result return type
func (g *LLVMGenerator) maybeWrapInResult(bodyValue value.Value, fnDecl *ast.FunctionDeclaration) value.Value {
	// Check if function declares a Result return type
	if fnDecl.ReturnType != nil && fnDecl.ReturnType.Name == TypeResult && len(fnDecl.ReturnType.GenericParams) >= 2 {
		// Check if bodyValue is already a Result struct type
		if _, isStruct := bodyValue.Type().(*types.StructType); isStruct {
			// Already a Result struct, no wrapping needed
			return bodyValue
		}

		successType := fnDecl.ReturnType.GenericParams[0].Name
		errorType := fnDecl.ReturnType.GenericParams[1].Name

		// Check if the body value is a plain int and function expects Result<int, MathError>
		if successType == "int" && errorType == TypeMathError && bodyValue.Type() == types.I64 {
			return g.wrapInMathResult(bodyValue)
		}
		// Check if the body value is a plain bool and function expects Result<bool, MathError>
		if successType == "bool" && errorType == TypeMathError &&
			(bodyValue.Type() == types.I1 || bodyValue.Type() == types.I64) {
			return g.wrapInBoolResult(bodyValue)
		}
		// Add other Result type mappings as needed
	}

	// No wrapping needed, return original value
	return bodyValue
}

// wrapInMathResult wraps a plain int value in a Result<int, MathError> structure
func (g *LLVMGenerator) wrapInMathResult(intValue value.Value) value.Value {
	// Create Result<int, MathError> structure by value
	resultType := g.getResultType(types.I64)

	// Use InsertValue to build the struct value directly
	undefStruct := constant.NewUndef(resultType)
	resultWithValue := g.builder.NewInsertValue(undefStruct, intValue, 0)
	resultComplete := g.builder.NewInsertValue(resultWithValue, constant.NewInt(types.I8, 0), 1)

	return resultComplete
}

// wrapInBoolResult wraps a plain bool value in a Result<bool, MathError> structure
func (g *LLVMGenerator) wrapInBoolResult(boolValue value.Value) value.Value {
	// Create Result<bool, MathError> structure (using i1 as the value type for booleans)
	resultType := g.getResultType(types.I1)

	// Ensure we store the correct type - if value is i64, convert to i1
	var valueToStore value.Value
	if boolValue.Type() == types.I64 {
		// Truncate i64 to i1 for boolean values
		valueToStore = g.builder.NewTrunc(boolValue, types.I1)
	} else {
		valueToStore = boolValue
	}

	// Use InsertValue to build the struct value directly
	undefStruct := constant.NewUndef(resultType)
	resultWithValue := g.builder.NewInsertValue(undefStruct, valueToStore, 0)
	resultComplete := g.builder.NewInsertValue(resultWithValue, constant.NewInt(types.I8, 0), 1)

	return resultComplete
}

// maybeUnwrapResult unwraps a Result value if the function return type is not a Result
func (g *LLVMGenerator) maybeUnwrapResult(bodyValue value.Value, fnDecl *ast.FunctionDeclaration) value.Value {
	// FIRST: Check explicit declaration (highest priority)
	if fnDecl.ReturnType != nil && fnDecl.ReturnType.Name == TypeResult {
		// Function explicitly declares Result return type - don't unwrap
		return bodyValue
	}

	// SECOND: Check INFERRED return type from type environment
	// This handles inferred types for functions without explicit annotations
	if fnType, exists := g.typeInferer.env.Get(fnDecl.Name); exists {
		if funcType, ok := fnType.(*FunctionType); ok {
			// Prune to resolve any type variables
			returnType := g.typeInferer.prune(funcType.returnType)

			// Check if the return type is a Result type
			if concrete, ok := returnType.(*ConcreteType); ok {
				typeName := concrete.String()
				// If return type is Result (or Result<...>), don't unwrap
				if typeName == TypeResult ||
					(strings.HasPrefix(typeName, "Result<") && strings.HasSuffix(typeName, ">")) {
					return bodyValue
				}
			}
		}
	}

	// AUTO-PROPAGATION: If body returns Result but function declares non-Result return,
	// unwrap the Result value. This allows: fn add(a,b)->int = a+b (body is Result<int>)
	if g.isResultType(bodyValue) {
		return g.unwrapIfResult(bodyValue)
	}

	return bodyValue
}

// canImplicitlyConvert checks if we can implicitly convert from one type to another
func (g *LLVMGenerator) canImplicitlyConvert(fromType, toType Type, _ *ast.FunctionDeclaration) bool {
	// Handle the case where Result types are still ConcreteType due to type inference issues
	if fromConcrete, ok := fromType.(*ConcreteType); ok {
		if toConcrete, ok := toType.(*ConcreteType); ok {
			// Check for primitive to Result conversion
			if (fromConcrete.name == TypeInt || fromConcrete.name == TypeBool) &&
				toConcrete.name == TypeResult {
				return true
			}

			// Check for Result<T> to T conversion (auto-unwrapping)
			// e.g., Result<int, MathError> can be converted to int
			if strings.HasPrefix(fromConcrete.name, "Result<") &&
				strings.HasSuffix(fromConcrete.name, ">") {
				// Extract the success type from Result<T, E>
				successType := g.typeInferer.extractResultSuccessType(fromConcrete.name)
				// Check if success type matches the target type
				return successType.String() == toConcrete.name
			}
		}
	}

	return g.canConvertToResult(fromType, toType)
}

// canConvertToResult checks if we can convert a primitive type to a Result type
func (g *LLVMGenerator) canConvertToResult(fromType, toType Type) bool {
	toGeneric, ok := toType.(*GenericType)
	if !ok || toGeneric.name != TypeResult || len(toGeneric.typeArgs) < 2 {
		return false
	}

	fromConcrete, ok := fromType.(*ConcreteType)
	if !ok {
		return false
	}

	successType, ok := toGeneric.typeArgs[0].(*ConcreteType)
	if !ok {
		return false
	}

	_, ok = toGeneric.typeArgs[1].(*ConcreteType)
	if !ok {
		return false
	}

	// Check if from type matches success type - allow any error type for Result conversion
	return fromConcrete.name == successType.name &&
		(successType.name == TypeInt || successType.name == TypeBool)
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
	case TypeFloat:
		return types.Double
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
	case TypeFloat:
		return types.Double
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
		// Parser truncates generic parameters, so "Result" without parameters
		// defaults to Result<string, Error> for now
		if ct.name == TypeResult {
			return g.getResultType(types.I8Ptr)
		}

		// Handle Result types like "Result<int, MathError>"
		if strings.HasPrefix(ct.name, "Result<") {
			// Result types are represented as structs with { value, discriminant }
			if ct.name == "Result<int, MathError>" {
				return g.getResultType(types.I64)
			}

			if ct.name == "Result<float, MathError>" {
				return g.getResultType(types.Double)
			}

			if ct.name == "Result<bool, MathError>" {
				return g.getResultType(types.I64)
			}

			if ct.name == "Result<string, Error>" {
				return g.getResultType(types.I8Ptr)
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
	case TypeResult:
		if len(gt.typeArgs) >= TwoArgs {
			// Result<T, E> - use value semantics for Result struct
			innerType := g.getLLVMType(gt.typeArgs[0])
			return g.getResultType(innerType)
		}

		return types.I64 // fallback
	case TypeList:
		if len(gt.typeArgs) >= 1 {
			// List<T> - represent as pointer to array struct { i64 length, i8* data }
			arrayStructType := types.NewStruct(types.I64, types.I8Ptr)
			return types.NewPointer(arrayStructType)
		}

		// Fallback - empty list type
		arrayStructType := types.NewStruct(types.I64, types.I8Ptr)
		return types.NewPointer(arrayStructType)
	case TypeMap:
		if len(gt.typeArgs) >= TwoTypeArgs {
			// Map<K, V> - represent as pointer to map struct { i64 length, i8* data }
			// For now, use the same structure as List - will be enhanced with C runtime
			mapStructType := types.NewStruct(types.I64, types.I8Ptr)
			return types.NewPointer(mapStructType)
		}

		// Fallback - empty map type
		mapStructType := types.NewStruct(types.I64, types.I8Ptr)
		return types.NewPointer(mapStructType)
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

// HINDLEY-MILNER FIX: Consistent field mapping system
// getOrCreateRecordFieldMapping creates a consistent field-to-index mapping for a record type
func (g *LLVMGenerator) getOrCreateRecordFieldMapping(recordTypeName string, fields map[string]Type) map[string]int {
	if mapping, exists := g.recordFieldMappings[recordTypeName]; exists {
		return mapping
	}

	// Create deterministic field ordering using sorted field names
	// This ensures consistent LLVM struct layout regardless of field declaration order
	fieldNames := make([]string, 0, len(fields))
	for fieldName := range fields {
		fieldNames = append(fieldNames, fieldName)
	}

	// Sort field names for deterministic ordering
	sort.Strings(fieldNames)

	// Create field-to-index mapping
	mapping := make(map[string]int)
	for i, fieldName := range fieldNames {
		mapping[fieldName] = i
	}

	// Cache the mapping
	g.recordFieldMappings[recordTypeName] = mapping

	return mapping
}

// createRecordTypeKey creates a unique key for record types that includes field types
func (g *LLVMGenerator) createRecordTypeKey(rt *RecordType) string {
	var keyParts []string

	keyParts = append(keyParts, rt.name)

	// Sort field names for consistent key generation
	fieldNames := make([]string, 0, len(rt.fields))
	for fieldName := range rt.fields {
		fieldNames = append(fieldNames, fieldName)
	}

	sort.Strings(fieldNames)

	// Add field types to the key
	for _, fieldName := range fieldNames {
		fieldType := rt.fields[fieldName]
		keyParts = append(keyParts, fmt.Sprintf("%s:%s", fieldName, fieldType.String()))
	}

	return strings.Join(keyParts, "_")
}

// getLLVMRecordType converts record types to LLVM struct types
func (g *LLVMGenerator) getLLVMRecordType(rt *RecordType) types.Type {
	// Create a unique key that includes both the record name and field types
	// This ensures that polymorphic records like Point<int,int> and Point<string,string>
	// get different LLVM struct types
	typeKey := g.createRecordTypeKey(rt)

	// Check if we already have this specific record type in the type map
	if llvmType, exists := g.typeMap[typeKey]; exists {
		// Return struct type by value, not pointer
		return llvmType
	}

	// HINDLEY-MILNER FIX: Use consistent field mapping
	fieldMapping := g.getOrCreateRecordFieldMapping(rt.name, rt.fields)

	// Create field types in the mapped order
	fieldTypes := make([]types.Type, len(rt.fields))

	for fieldName, fieldType := range rt.fields {
		index := fieldMapping[fieldName]
		fieldTypes[index] = g.getLLVMType(fieldType)
	}

	structType := types.NewStruct(fieldTypes...)
	g.typeMap[typeKey] = structType

	return structType
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

		// Sort field names to match ObjectLiteral and field access ordering
		fieldNames := make([]string, len(variant.Fields))
		fieldTypeMap := make(map[string]string)

		for i, field := range variant.Fields {
			fieldNames[i] = field.Name
			fieldTypeMap[field.Name] = field.Type
		}

		sort.Strings(fieldNames)

		// Create field types in sorted order
		fieldTypes := make([]types.Type, len(fieldNames))
		for i, fieldName := range fieldNames {
			fieldTypes[i] = g.getFieldType(fieldTypeMap[fieldName])
		}

		// Create the struct type
		structType := types.NewStruct(fieldTypes...)
		g.typeMap[typeDecl.Name] = structType

		// Create a proper RecordType for the type inference environment
		// Infer field types for the record
		recordFieldTypes := make(map[string]Type)

		fieldOrder := make([]string, 0, len(variant.Fields))

		for _, field := range variant.Fields {
			recordFieldTypes[field.Name] = g.getInferenceFieldType(field.Type)
			fieldOrder = append(fieldOrder, field.Name)
		}

		recordType := NewOrderedRecordType(typeDecl.Name, recordFieldTypes, fieldOrder)
		g.typeInferer.env.Set(typeDecl.Name, recordType)

		// Also register the variant name in the type environment
		// This allows type constructors like Warrior { health: 100, strength: 20 } to work
		g.typeInferer.env.Set(variant.Name, recordType)

		// Register the variant name in unionVariants for pattern matching
		// This allows pattern matching to work correctly for single-variant types with fields
		g.unionVariants[variant.Name] = 0 // Use 0 as discriminant for single-variant types

		// Store field names for field access (already sorted above)
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
	//ALWAYS LOWERCASE!
	case TypeString: // "string"
		return types.I8Ptr
		//ALWAYS LOWERCASE!
	case TypeInt: // "int"
		return types.I64
		//ALWAYS LOWERCASE!
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

// getInferenceFieldType converts a field type string to a Type for type inference
func (g *LLVMGenerator) getInferenceFieldType(fieldType string) Type {
	switch fieldType {
	case TypeInt, "Int":
		return &ConcreteType{name: TypeInt}
	case TypeString, "String":
		return &ConcreteType{name: TypeString}
	case TypeBool, "Bool":
		return &ConcreteType{name: TypeBool}
	default:
		// For user-defined types, check if they exist in the type map
		if _, exists := g.typeMap[fieldType]; exists {
			return &ConcreteType{name: fieldType}
		}
		// Default to creating a concrete type
		return &ConcreteType{name: fieldType}
	}
}
