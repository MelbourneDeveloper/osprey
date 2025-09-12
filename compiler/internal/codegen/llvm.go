package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// Constants for type sizes and arrays (shared with expression_generation.go)
const (
	LargeArraySizeForCasting = 1000 // Large array size for type casting operations
	TruePattern              = "true"
	FalsePattern             = "false"
)

func (g *LLVMGenerator) generateCallExpression(callExpr *ast.CallExpression) (value.Value, error) {
	// Extract function name and handle built-in functions
	funcName, isBuiltin, builtinResult, err := g.extractFunctionAndHandleBuiltins(callExpr)
	if err != nil {
		return nil, err
	}

	if isBuiltin {
		return builtinResult, nil
	}

	// Validate function call
	err = g.validateFunctionCall(funcName, callExpr)
	if err != nil {
		return nil, err
	}

	// Resolve function value
	funcValue, err := g.resolveFunctionValue(funcName, callExpr)
	if err != nil {
		return nil, err
	}

	// Validate that the resolved value is actually callable
	err = g.validateCallableType(funcValue, funcName, callExpr)
	if err != nil {
		return nil, err
	}

	// Generate arguments
	args, err := g.generateCallArguments(funcName, callExpr, funcValue)
	if err != nil {
		return nil, err
	}

	// Create the function call
	return g.builder.NewCall(funcValue, args...), nil
}

// extractFunctionAndHandleBuiltins extracts function name and handles built-in functions
func (g *LLVMGenerator) extractFunctionAndHandleBuiltins(
	callExpr *ast.CallExpression,
) (string, bool, value.Value, error) {
	ident, ok := callExpr.Function.(*ast.Identifier)
	if !ok {
		return "", false, nil, nil
	}

	funcName := ident.Name

	// Check if this is a built-in function
	result, err := g.handleBuiltInFunction(funcName, callExpr)
	if err != nil {
		return funcName, false, nil, err
	}

	if _, exists := GlobalBuiltInRegistry.GetFunction(funcName); exists {
		// This is a built-in function that was handled successfully
		return funcName, true, result, nil
	}

	return funcName, false, nil, nil
}

// validateFunctionCall validates the function call arguments
func (g *LLVMGenerator) validateFunctionCall(funcName string, callExpr *ast.CallExpression) error {
	if funcName == "" {
		return nil
	}

	// Type inference already validates any type mismatches

	// Validate named arguments requirement for multi-parameter functions
	return g.validateNamedArguments(funcName, callExpr)
}

// resolveFunctionValue resolves the function value for user-defined functions
func (g *LLVMGenerator) resolveFunctionValue(funcName string, callExpr *ast.CallExpression) (value.Value, error) {
	if funcName == "" {
		// Not an identifier (e.g., lambda or function parameter), generate normally
		return g.generateExpression(callExpr.Function)
	}

	// Check if this function name is a parameter (not a declared function)
	if val, isParam := g.variables[funcName]; isParam {
		// This is a function parameter being called
		return val, nil
	}

	// Infer the argument types to determine which monomorphized instance to use
	argTypes, err := g.inferCallArgumentTypes(callExpr)
	if err != nil {
		return nil, err
	}

	// Get the correct monomorphized function
	return g.resolveMonomorphizedFunction(funcName, argTypes)
}

// validateCallableType validates that a value is callable (function or function pointer)
func (g *LLVMGenerator) validateCallableType(
	funcValue value.Value,
	_ string,
	callExpr *ast.CallExpression,
) error {
	// Check if the value is a function type or function pointer
	valueType := funcValue.Type()

	// Handle pointer to function
	if ptrType, ok := valueType.(*types.PointerType); ok {
		if _, isFuncType := ptrType.ElemType.(*types.FuncType); isFuncType {
			return nil // Valid function pointer
		}
	}

	// Handle direct function type
	if _, isFuncType := valueType.(*types.FuncType); isFuncType {
		return nil // Valid function
	}

	// Not a callable type - return appropriate error
	if callExpr.Function != nil {
		if ident, ok := callExpr.Function.(*ast.Identifier); ok && ident.Position != nil {
			return fmt.Errorf("line %d:%d: %w",
				ident.Position.Line, ident.Position.Column, ErrNotAFunction)
		}
	}

	return fmt.Errorf("%w", ErrNotAFunction)
}

// generateCallArguments generates the arguments for the function call
func (g *LLVMGenerator) generateCallArguments(
	funcName string,
	callExpr *ast.CallExpression,
	funcValue value.Value,
) ([]value.Value, error) {
	if len(callExpr.NamedArguments) > 0 {
		return g.generateNamedArguments(funcName, callExpr, funcValue)
	}

	return g.generateRegularArguments(callExpr, funcValue)
}

// generateNamedArguments handles named arguments
func (g *LLVMGenerator) generateNamedArguments(
	funcName string,
	callExpr *ast.CallExpression,
	funcValue value.Value,
) ([]value.Value, error) {
	reorderedExprs, err := g.reorderNamedArguments(funcName, callExpr.NamedArguments)
	if err != nil {
		return nil, err
	}

	args := make([]value.Value, len(reorderedExprs))

	for i, expr := range reorderedExprs {
		arg, err := g.generateSingleArgument(expr, i, callExpr, funcValue)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

// generateRegularArguments handles regular (positional) arguments
func (g *LLVMGenerator) generateRegularArguments(
	callExpr *ast.CallExpression,
	funcValue value.Value,
) ([]value.Value, error) {
	args := make([]value.Value, len(callExpr.Arguments))

	for i, argExpr := range callExpr.Arguments {
		arg, err := g.generateSingleArgument(argExpr, i, callExpr, funcValue)
		if err != nil {
			return nil, err
		}

		args[i] = arg
	}

	return args, nil
}

// generateSingleArgument generates a single argument with proper type context
func (g *LLVMGenerator) generateSingleArgument(
	expr ast.Expression,
	index int,
	callExpr *ast.CallExpression,
	funcValue value.Value,
) (value.Value, error) {
	// Set expected parameter type context for boolean literals
	g.setParameterTypeContext(funcValue, index)
	defer g.clearParameterTypeContext()

	return g.generateTypedArgumentExpression(expr, index, callExpr, funcValue)
}

// generateArgumentExpression generates an argument value, handling polymorphic functions specially
func (g *LLVMGenerator) generateArgumentExpression(
	expr ast.Expression,
	argIndex int,
	callExpr *ast.CallExpression,
) (value.Value, error) {
	// Handle polymorphic function identifiers
	if ident, ok := expr.(*ast.Identifier); ok {
		if result, handled, err := g.handlePolymorphicFunctionArgument(ident, argIndex, callExpr); handled {
			return result, err
		}
	}

	// For non-polymorphic function arguments, generate normally
	return g.generateExpression(expr)
}

// handlePolymorphicFunctionArgument handles polymorphic function arguments
func (g *LLVMGenerator) handlePolymorphicFunctionArgument(
	ident *ast.Identifier,
	argIndex int,
	callExpr *ast.CallExpression,
) (value.Value, bool, error) {
	// Check if this is a polymorphic function
	funcTypeFromEnv, exists := g.typeInferer.env.Get(ident.Name)
	if !exists {
		return nil, false, nil
	}

	_, isPolymorphic := funcTypeFromEnv.(*TypeScheme)
	if !isPolymorphic {
		return nil, false, nil
	}

	// This is a polymorphic function being passed as an argument
	// We need to infer the complete call expression to determine concrete types
	_, err := g.typeInferer.InferType(callExpr)
	if err != nil {
		return nil, true, fmt.Errorf("failed to infer call type: %w", err)
	}

	// Get the argument types from the inferred call
	argTypes, err := g.inferCallArgumentTypes(callExpr)
	if err != nil {
		return nil, true, fmt.Errorf("failed to infer argument types: %w", err)
	}

	// Try to resolve the function type for this argument
	result, resolved, err := g.resolvePolymorphicArgument(ident.Name, argIndex, argTypes)
	if resolved || err != nil {
		return result, true, err
	}

	// No alternative resolution needed - Osprey requires named args for multi-param functions
	return nil, false, nil
}

// resolvePolymorphicArgument resolves a polymorphic function argument by direct type matching
func (g *LLVMGenerator) resolvePolymorphicArgument(
	funcName string,
	argIndex int,
	argTypes []Type,
) (value.Value, bool, error) {
	// For the function argument, determine its concrete function type based on context
	if argIndex >= len(argTypes) {
		return nil, false, nil
	}

	fnType, ok := argTypes[argIndex].(*FunctionType)
	if !ok {
		return nil, false, nil
	}

	// Get the monomorphized function instance
	mangledName := g.getMonomorphizedName(funcName, fnType)

	// Check if we already have this monomorphized instance
	if fn, exists := g.functions[mangledName]; exists {
		return fn, true, nil
	}

	// Generate it on-demand
	result, err := g.generateMonomorphizedInstance(funcName, fnType)

	return result, true, err
}

// generateTypedArgumentExpression generates an argument with type awareness
func (g *LLVMGenerator) generateTypedArgumentExpression(
	expr ast.Expression,
	argIndex int,
	callExpr *ast.CallExpression,
	funcValue value.Value,
) (value.Value, error) {
	// Handle boolean literals directly if we know the expected type
	if boolLit, ok := expr.(*ast.BooleanLiteral); ok && funcValue != nil {
		// Try to get the function signature directly from the IR function
		if irFunc, ok := funcValue.(*ir.Func); ok {
			if argIndex < len(irFunc.Sig.Params) {
				expectedType := irFunc.Sig.Params[argIndex]

				// Generate boolean with the correct type from the start
				if expectedType == types.I1 {
					if boolLit.Value {
						return constant.NewBool(true), nil
					}

					return constant.NewBool(false), nil
				}
			}
		}
	}

	// For non-boolean literals or when type is not i1, use the polymorphic function handling
	val, err := g.generateArgumentExpression(expr, argIndex, callExpr)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// inferCallArgumentTypes infers the types of arguments in a function call
func (g *LLVMGenerator) inferCallArgumentTypes(callExpr *ast.CallExpression) ([]Type, error) {
	var argTypes []Type

	if len(callExpr.NamedArguments) > 0 {
		// Process named arguments
		for _, namedArg := range callExpr.NamedArguments {
			argType, err := g.typeInferer.InferType(namedArg.Value)
			if err != nil {
				return nil, err
			}

			argTypes = append(argTypes, argType)
		}
	} else {
		// Process regular arguments
		for _, arg := range callExpr.Arguments {
			argType, err := g.typeInferer.InferType(arg)
			if err != nil {
				return nil, err
			}

			argTypes = append(argTypes, argType)
		}
	}

	return argTypes, nil
}

// resolveMonomorphizedFunction finds or creates the correct monomorphized instance
func (g *LLVMGenerator) resolveMonomorphizedFunction(funcName string, argTypes []Type) (value.Value, error) {
	// Check if we have the function type in the environment
	funcTypeFromEnv, exists := g.typeInferer.env.Get(funcName)
	if !exists {
		// Fallback to existing function lookup
		if fn, exists := g.functions[funcName]; exists {
			return fn, nil
		}

		return nil, fmt.Errorf("%w: %s", ErrFunctionNotDeclared, funcName)
	}

	// Handle TypeScheme (polymorphic function)
	if scheme, ok := funcTypeFromEnv.(*TypeScheme); ok {
		// Instantiate the scheme to get a concrete function type
		instantiated := g.typeInferer.Instantiate(scheme)
		if fnType, ok := instantiated.(*FunctionType); ok {
			// Unify the argument types to determine the concrete instantiation
			expectedFnType := &FunctionType{
				paramTypes: argTypes,
				returnType: g.typeInferer.Fresh(),
			}

			err := g.typeInferer.Unify(fnType, expectedFnType)
			if err != nil {
				return nil, fmt.Errorf("function call type mismatch: %w", err)
			}

			// Use the declared return type from the function declaration if available
			fnDecl, exists := g.functionDeclarations[funcName]

			var concreteReturnType Type
			if exists && fnDecl.ReturnType != nil {
				// Use the declared return type for functions with explicit return type annotations
				concreteReturnType = g.typeExpressionToInferenceType(fnDecl.ReturnType)
			} else {
				// Re-infer the return type with concrete parameter types for accurate monomorphization
				var err error

				concreteReturnType, err = g.reInferReturnType(funcName, argTypes)
				if err != nil {
					// Fallback to original return type if re-inference fails
					concreteReturnType = g.typeInferer.prune(expectedFnType.returnType)
				}
			}

			// Create the concrete function type after unification
			concreteFnType := &FunctionType{
				paramTypes: argTypes,
				returnType: concreteReturnType,
			}

			// Get the monomorphized name
			mangledName := g.getMonomorphizedName(funcName, concreteFnType)

			// Check if we already have this monomorphized instance
			if fn, exists := g.functions[mangledName]; exists {
				return fn, nil
			}

			// If not, we need to generate it on-demand
			return g.generateMonomorphizedInstance(funcName, concreteFnType)
		}
	}

	// For non-polymorphic functions, just return the existing function
	if fn, exists := g.functions[funcName]; exists {
		return fn, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrFunctionNotDeclared, funcName)
}

// generateMonomorphizedInstance generates a new monomorphized instance of a polymorphic function
func (g *LLVMGenerator) generateMonomorphizedInstance(
	baseFuncName string,
	concreteFnType *FunctionType,
) (value.Value, error) {
	// Get the original function declaration
	fnDecl, exists := g.functionDeclarations[baseFuncName]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrFunctionNotDeclared, baseFuncName)
	}

	// Create the monomorphized name
	mangledName := g.getMonomorphizedName(baseFuncName, concreteFnType)

	// Check if it was already generated (race condition protection)
	if fn, exists := g.functions[mangledName]; exists {
		return fn, nil
	}

	// Create the LLVM function signature for this monomorphized instance
	llvmReturnType := g.getLLVMType(concreteFnType.returnType)

	params := make([]*ir.Param, len(concreteFnType.paramTypes))
	for i, paramType := range concreteFnType.paramTypes {
		params[i] = ir.NewParam(fnDecl.Parameters[i].Name, g.getLLVMType(paramType))
	}

	// Create the LLVM function
	fn := g.module.NewFunc(mangledName, llvmReturnType, params...)
	g.functions[mangledName] = fn

	// Store parameter names for this instance
	g.functionParameters[mangledName] = make([]string, len(fnDecl.Parameters))
	for i, param := range fnDecl.Parameters {
		g.functionParameters[mangledName][i] = param.Name
	}

	// Generate the function body
	err := g.generateMonomorphizedFunctionBody(fnDecl, fn, concreteFnType)
	if err != nil {
		// Clean up on error
		delete(g.functions, mangledName)
		delete(g.functionParameters, mangledName)

		return nil, err
	}

	return fn, nil
}

// generateMonomorphizedFunctionBody generates the body for a monomorphized function instance
func (g *LLVMGenerator) generateMonomorphizedFunctionBody(
	fnDecl *ast.FunctionDeclaration,
	fn *ir.Func,
	concreteFnType *FunctionType,
) error {
	// Save current state
	oldFunction := g.function
	oldBuilder := g.builder

	oldVariables := make(map[string]value.Value)
	for k, v := range g.variables {
		oldVariables[k] = v
	}

	oldExpectedReturnType := g.expectedReturnType

	// Set up new function context
	entry := fn.NewBlock("")
	g.builder = entry
	g.function = fn
	g.expectedReturnType = fn.Sig.RetType

	// Set up parameters with concrete types in the type environment
	savedTypeEnv := g.typeInferer.env.Clone()

	for i, param := range fnDecl.Parameters {
		// Add parameter value to variables
		g.variables[param.Name] = fn.Params[i]
		// Set concrete type in type environment for this monomorphization
		if i < len(concreteFnType.paramTypes) {
			g.typeInferer.env.Set(param.Name, concreteFnType.paramTypes[i])
		}
	}

	// Set the function type in the environment with the concrete type
	g.typeInferer.env.Set(fnDecl.Name, concreteFnType)

	// Set expected return type for proper expression generation (especially for boolean operators)
	g.expectedReturnType = g.getLLVMType(concreteFnType.returnType)

	// Generate the function body
	bodyValue, err := g.generateExpression(fnDecl.Body)
	if err != nil {
		// Restore state on error
		g.function = oldFunction
		g.builder = oldBuilder
		g.variables = oldVariables
		g.expectedReturnType = oldExpectedReturnType
		g.typeInferer.env = savedTypeEnv

		return err
	}

	// Generate return instruction
	if fn.Sig.RetType == types.Void {
		g.builder.NewRet(nil)
	} else if fnDecl.Name == MainFunctionName {
		g.builder.NewRet(constant.NewInt(types.I32, 0))
	} else {
		finalReturnValue := g.maybeWrapInResult(bodyValue, fnDecl)
		g.builder.NewRet(finalReturnValue)
	}

	// Restore original state
	g.function = oldFunction
	g.builder = oldBuilder
	g.variables = oldVariables
	g.expectedReturnType = oldExpectedReturnType
	g.typeInferer.env = savedTypeEnv

	return nil
}

// handleBuiltInFunction handles all built-in function calls using the unified registry.
func (g *LLVMGenerator) handleBuiltInFunction(name string, callExpr *ast.CallExpression) (value.Value, error) {
	// Check if this is a built-in function using the registry
	fn, exists := GlobalBuiltInRegistry.GetFunction(name)
	if !exists {
		return nil, nil // Not a built-in function
	}

	// Check security permissions
	if !g.checkSecurityPermission(fn.SecurityFlag) {
		return nil, WrapUnsupportedCallExpressionSecurity(name)
	}

	// Call the function's generator
	return fn.Generator(g, callExpr)
}

// checkSecurityPermission checks if the generator's security config allows the given permission
func (g *LLVMGenerator) checkSecurityPermission(permission SecurityPermission) bool {
	switch permission {
	case PermissionNone:
		return true
	case PermissionHTTP:
		return g.security.AllowHTTP
	case PermissionWebSocket:
		return g.security.AllowWebSocket
	case PermissionFileRead:
		return g.security.AllowFileRead
	case PermissionFileWrite:
		return g.security.AllowFileWrite
	case PermissionProcess:
		return g.security.AllowProcessExecution
	case PermissionFFI:
		return g.security.AllowFFI
	default:
		return false
	}
}

// generateInterpolatedString generates LLVM IR for interpolated strings by concatenating parts.
func (g *LLVMGenerator) generateInterpolatedString(interpStr *ast.InterpolatedStringLiteral) (value.Value, error) {
	// For now, we'll use a simple approach: build the string by calling printf multiple times
	// into a buffer. A more sophisticated implementation would use string concatenation.
	// If there's only one part and it's text, treat it as a regular string
	if len(interpStr.Parts) == 1 && !interpStr.Parts[0].IsExpression {
		str := constant.NewCharArrayFromString(interpStr.Parts[0].Text + StringTerminator)
		global := g.module.NewGlobalDef("", str)

		return g.builder.NewGetElementPtr(str.Typ, global,
			constant.NewInt(types.I32, ArrayIndexZero),
			constant.NewInt(types.I32, ArrayIndexZero)), nil
	}

	// For multiple parts or expressions, we'll use sprintf to build the string
	// First, we need sprintf declared
	sprintf := g.ensureSprintfDeclaration()

	// Build format string and collect arguments
	var formatParts []string

	var args []value.Value

	for _, part := range interpStr.Parts {
		if part.IsExpression {
			// Auto-call toString() on all expressions in string interpolation
			toStringCall := &ast.CallExpression{
				Function:  &ast.Identifier{Name: "toString"},
				Arguments: []ast.Expression{part.Expression},
			}

			// Generate the toString call which will return a string
			val, err := g.generateExpression(toStringCall)
			if err != nil {
				return nil, err
			}

			args = append(args, val)

			// All expressions become %s since toString() always returns string
			formatParts = append(formatParts, "%s")
		} else {
			// Escape % characters in literal text by replacing % with %%
			escapedText := strings.ReplaceAll(part.Text, "%", "%%")
			formatParts = append(formatParts, escapedText)
		}
	}

	// Create the format string
	formatString := strings.Join(formatParts, "") + StringTerminator
	formatStr := constant.NewCharArrayFromString(formatString)
	formatGlobal := g.module.NewGlobalDef("", formatStr)
	formatPtr := g.builder.NewGetElementPtr(formatStr.Typ, formatGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Allocate buffer for result string (simplified - use fixed size)
	bufferType := types.NewArray(BufferSize1KB, types.I8) // 1KB buffer
	buffer := g.builder.NewAlloca(bufferType)
	bufferPtr := g.builder.NewGetElementPtr(bufferType, buffer,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Call sprintf(buffer, format, args...)
	sprintfArgs := make([]value.Value, 0, len(args)+TwoArgs)
	sprintfArgs = append(sprintfArgs, bufferPtr, formatPtr)
	sprintfArgs = append(sprintfArgs, args...)

	g.builder.NewCall(sprintf, sprintfArgs...)

	return bufferPtr, nil
}

// ensureSprintfDeclaration declares sprintf if not already declared.
func (g *LLVMGenerator) ensureSprintfDeclaration() *ir.Func {
	if sprintf, exists := g.functions["sprintf"]; exists {
		return sprintf
	}

	// Declare sprintf: i32 @sprintf(i8* %str, i8* %format, ...)
	sprintf := g.module.NewFunc("sprintf", types.I32,
		ir.NewParam("str", types.I8Ptr),
		ir.NewParam("format", types.I8Ptr))
	sprintf.Sig.Variadic = true
	g.functions["sprintf"] = sprintf

	return sprintf
}

func (g *LLVMGenerator) generateIntToString(arg value.Value) (value.Value, error) {
	// Ensure sprintf and malloc are declared
	sprintf := g.ensureSprintfDeclaration()
	malloc := g.ensureMallocDeclaration()

	// Create format string for integer conversion
	formatStr := constant.NewCharArrayFromString("%ld\x00")
	formatGlobal := g.module.NewGlobalDef("", formatStr)
	formatPtr := g.builder.NewGetElementPtr(formatStr.Typ, formatGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Allocate buffer for result string using malloc (64 bytes should be enough for any 64-bit integer)
	bufferSize := constant.NewInt(types.I64, BufferSize64Bytes)
	bufferPtr := g.builder.NewCall(malloc, bufferSize)

	// Call sprintf(buffer, "%ld", arg)
	g.builder.NewCall(sprintf, bufferPtr, formatPtr, arg)

	return bufferPtr, nil
}

func (g *LLVMGenerator) generateBoolToString(arg value.Value) (value.Value, error) {
	// Create string constants for "true" and "false"
	trueStr := constant.NewCharArrayFromString("true\x00")
	trueGlobal := g.module.NewGlobalDef("", trueStr)
	truePtr := g.builder.NewGetElementPtr(trueStr.Typ, trueGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	falseStr := constant.NewCharArrayFromString("false\x00")
	falseGlobal := g.module.NewGlobalDef("", falseStr)
	falsePtr := g.builder.NewGetElementPtr(falseStr.Typ, falseGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Create condition based on argument type
	var condition value.Value
	if arg.Type() == types.I1 {
		// For i1 (boolean) types, use directly
		condition = arg
	} else {
		// For i64 types, compare against zero
		zero := constant.NewInt(types.I64, 0)
		condition = g.builder.NewICmp(enum.IPredNE, arg, zero)
	}

	// Use select instruction: select condition, trueValue, falseValue
	return g.builder.NewSelect(condition, truePtr, falsePtr), nil
}

func (g *LLVMGenerator) generateMatchExpression(matchExpr *ast.MatchExpression) (value.Value, error) {
	discriminant, err := g.generateExpression(matchExpr.Expression)
	if err != nil {
		return nil, err
	}

	// Validate match expression for exhaustiveness and unknown variants
	// Now that we have the discriminant, we can infer its type for proper validation
	discriminantType := g.inferDiscriminantTypeName(discriminant, matchExpr.Expression)
	err = g.validateMatchExpressionWithType(matchExpr, discriminantType)
	if err != nil {
		return nil, err
	}

	return g.generateMatchExpressionWithDiscriminant(matchExpr, discriminant)
}

// inferDiscriminantTypeName infers the type name of the discriminant for match validation
func (g *LLVMGenerator) inferDiscriminantTypeName(_ value.Value, expr ast.Expression) string {
	// Try to get the type from the type environment if the discriminant expression is an identifier
	if ident, ok := expr.(*ast.Identifier); ok {
		if varType, exists := g.typeInferer.env.Get(ident.Name); exists {
			resolvedType := g.typeInferer.ResolveType(varType)
			if concreteType, ok := resolvedType.(*ConcreteType); ok {
				return concreteType.name
			}

			if recordType, ok := resolvedType.(*RecordType); ok {
				// For record types, we need to find which type declaration this belongs to
				for _, typeDecl := range g.typeDeclarations {
					if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
						// Check if this record type matches the variant's fields
						variant := typeDecl.Variants[0]
						if len(recordType.fields) == len(variant.Fields) {
							return typeDecl.Name
						}
					}
				}
			}
		}
	}

	// Alternative approach: directly call type inference on the expression
	inferredType, err := g.typeInferer.InferType(expr)
	if err == nil {
		resolvedType := g.typeInferer.ResolveType(inferredType)
		if concreteType, ok := resolvedType.(*ConcreteType); ok {
			return concreteType.name
		}
	}

	// If that fails, try to find the type by looking at the type declarations
	// This is needed for custom ADTs where the variable has the type name, not constructor name
	return "UnknownType"
}

// generateMatchExpressionWithDiscriminant generates match expression with pre-computed discriminant.
func (g *LLVMGenerator) generateMatchExpressionWithDiscriminant(
	matchExpr *ast.MatchExpression,
	discriminant value.Value,
) (value.Value, error) {
	if g.hasResultPatterns(matchExpr.Arms) {
		// Check if the discriminant is already a Result type
		if g.isResultType(discriminant) {
			return g.generateResultMatchExpression(matchExpr, discriminant)
		}

		// If not a Result type, wrap it in a Success Result automatically
		wrappedDiscriminant := g.wrapInSuccessResult(discriminant)

		return g.generateResultMatchExpression(matchExpr, wrappedDiscriminant)
	}

	if len(matchExpr.Arms) == 0 {
		return constant.NewInt(types.I64, ArrayIndexZero), nil
	}

	return g.generateStandardMatchExpression(matchExpr, discriminant)
}

// hasResultPatterns checks if the match expression has Success/Err patterns.
func (g *LLVMGenerator) hasResultPatterns(arms []ast.MatchArm) bool {
	for _, arm := range arms {
		if arm.Pattern.Constructor == SuccessPattern || arm.Pattern.Constructor == ErrorPattern {
			return true
		}
	}

	return false
}

// wrapInSuccessResult wraps a value in a Success Result automatically
func (g *LLVMGenerator) wrapInSuccessResult(discriminant value.Value) value.Value {
	// Create a Result struct with the discriminant value as the success value
	// Result struct: [value, discriminant] where discriminant=0 for success
	resultType := types.NewStruct(discriminant.Type(), types.I8)

	// Allocate memory for the result
	resultPtr := g.builder.NewAlloca(resultType)

	// Store the value in the first field
	valuePtr := g.builder.NewGetElementPtr(resultType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(discriminant, valuePtr)

	// Store 0 (success) in the discriminant field
	discriminantPtr := g.builder.NewGetElementPtr(resultType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	return resultPtr
}

// isResultType checks if a value is a Result type (struct with two fields or pointer to such struct)
func (g *LLVMGenerator) isResultType(val value.Value) bool {
	// Check for pointer to struct (legacy pointer semantics)
	if ptrType, ok := val.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok {
			return len(structType.Fields) == ResultFieldCount
		}
	}
	// Check for struct value directly (value semantics)
	if structType, ok := val.Type().(*types.StructType); ok {
		return len(structType.Fields) == ResultFieldCount
	}

	return false
}

// generateStandardMatchExpression generates a standard (non-result) match expression.
func (g *LLVMGenerator) generateStandardMatchExpression(
	matchExpr *ast.MatchExpression,
	discriminant value.Value,
) (value.Value, error) {
	// Create unique block names for this match expression
	blockSuffix := fmt.Sprintf("_%p", matchExpr)
	endBlock := g.function.NewBlock("match_end" + blockSuffix)
	armBlocks := g.createMatchArmBlocks(matchExpr.Arms, blockSuffix)

	// Generate conditions first (creates branching logic)
	g.generateMatchConditions(matchExpr.Arms, armBlocks, discriminant, blockSuffix)

	// Then generate arm values (adds terminators to arm blocks)
	armValues, predecessorBlocks, err := g.generateMatchArmValues(matchExpr.Arms, armBlocks, endBlock, discriminant)
	if err != nil {
		return nil, err
	}

	return g.createMatchResult(armValues, predecessorBlocks, endBlock)
}

// createMatchArmBlocks creates LLVM blocks for each match arm.
func (g *LLVMGenerator) createMatchArmBlocks(arms []ast.MatchArm, blockSuffix string) []*ir.Block {
	var armBlocks []*ir.Block
	for i := range arms {
		armBlocks = append(armBlocks, g.function.NewBlock(fmt.Sprintf("match_arm_%d%s", i, blockSuffix)))
	}

	return armBlocks
}

// generateMatchArmValues generates values for each match arm.
func (g *LLVMGenerator) generateMatchArmValues(
	arms []ast.MatchArm,
	armBlocks []*ir.Block,
	endBlock *ir.Block,
	discriminant value.Value,
) ([]value.Value, []*ir.Block, error) {
	var (
		armValues         []value.Value
		predecessorBlocks []*ir.Block
	)

	for i, arm := range arms {
		// Set builder to the arm block at the start of each iteration
		g.builder = armBlocks[i]

		// Handle variable binding in patterns
		armValue, err := g.processMatchArm(arm, discriminant)
		if err != nil {
			return nil, nil, err
		}

		armValues = append(armValues, armValue)

		// Handle branch termination
		currentBuilderBlock := g.addBranchTermination(endBlock)
		predecessorBlocks = append(predecessorBlocks, currentBuilderBlock)
	}

	return armValues, predecessorBlocks, nil
}

// processMatchArm handles the processing of a single match arm
func (g *LLVMGenerator) processMatchArm(arm ast.MatchArm, discriminant value.Value) (value.Value, error) {
	if arm.Pattern.Variable != "" || len(arm.Pattern.Fields) > 0 {
		return g.processMatchArmWithBinding(arm, discriminant)
	}

	return g.processMatchArmWithoutBinding(arm)
}

// processMatchArmWithBinding handles match arms that have variable binding
func (g *LLVMGenerator) processMatchArmWithBinding(arm ast.MatchArm, discriminant value.Value) (value.Value, error) {
	// Save the current variable scope
	oldVariables := g.saveVariableScope()

	// Bind the pattern variable to the discriminant value
	if arm.Pattern.Variable != "" {
		g.variables[arm.Pattern.Variable] = discriminant
	}

	// Handle field extraction for structural matching and discriminated unions
	if len(arm.Pattern.Fields) > 0 {
		g.extractPatternFields(arm.Pattern, discriminant)
	}

	// Generate the arm expression
	armValue, err := g.generateExpression(arm.Expression)
	if err != nil {
		g.restoreVariableScope(oldVariables)
		return nil, err
	}

	// Restore the previous variable scope
	g.restoreVariableScope(oldVariables)

	return g.normalizeArmValue(armValue), nil
}

// processMatchArmWithoutBinding handles match arms without variable binding
func (g *LLVMGenerator) processMatchArmWithoutBinding(arm ast.MatchArm) (value.Value, error) {
	armValue, err := g.generateExpression(arm.Expression)
	if err != nil {
		return nil, err
	}

	return g.normalizeArmValue(armValue), nil
}

// saveVariableScope saves the current variable scope
func (g *LLVMGenerator) saveVariableScope() map[string]value.Value {
	oldVariables := make(map[string]value.Value)
	for k, v := range g.variables {
		oldVariables[k] = v
	}

	return oldVariables
}

// restoreVariableScope restores a saved variable scope
func (g *LLVMGenerator) restoreVariableScope(oldVariables map[string]value.Value) {
	g.variables = oldVariables
}

// extractPatternFields handles field extraction for patterns
func (g *LLVMGenerator) extractPatternFields(pattern ast.Pattern, discriminant value.Value) {
	if pattern.Constructor == "*" {
		// For structural matching, extract fields from the object
		g.extractStructuralFields(pattern.Fields, discriminant)
		return
	}

	// Try to extract fields for record types (single-variant)
	if g.extractRecordTypeFields(pattern, discriminant) {
		return
	}

	// Try to extract fields for discriminated unions (multi-variant)
	if g.extractDiscriminatedUnionTypeFields(pattern, discriminant) {
		return
	}

	// If no type found, bind fields to zero values
	g.bindFieldsToZeroValues(pattern.Fields)
}

// extractRecordTypeFields extracts fields for single-variant types
func (g *LLVMGenerator) extractRecordTypeFields(pattern ast.Pattern, discriminant value.Value) bool {
	for _, typeDecl := range g.typeDeclarations {
		if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
			if typeDecl.Variants[0].Name == pattern.Constructor {
				g.extractRecordFields(pattern, discriminant, typeDecl.Variants[0])
				return true
			}
		}
	}

	return false
}

// extractDiscriminatedUnionTypeFields extracts fields for multi-variant discriminated unions
func (g *LLVMGenerator) extractDiscriminatedUnionTypeFields(pattern ast.Pattern, discriminant value.Value) bool {
	if _, exists := g.unionVariants[pattern.Constructor]; !exists {
		return false
	}

	// Check if this is actually a multi-variant union
	if !g.isMultiVariantType(pattern.Constructor) {
		return false
	}

	// Handle discriminated union field extraction
	err := g.extractDiscriminatedUnionFields(discriminant, pattern, g.variables)
	if err != nil {
		// If field extraction fails, bind fields to zero values
		g.bindFieldsToZeroValues(pattern.Fields)
	}

	return true
}

// isMultiVariantType checks if a pattern constructor belongs to a multi-variant type
func (g *LLVMGenerator) isMultiVariantType(constructorName string) bool {
	for _, typeDecl := range g.typeDeclarations {
		if len(typeDecl.Variants) > 1 {
			for _, variant := range typeDecl.Variants {
				if variant.Name == constructorName {
					return true
				}
			}
		}
	}

	return false
}

// extractStructuralFields extracts fields for structural matching
func (g *LLVMGenerator) extractStructuralFields(fields []string, discriminant value.Value) {
	for _, fieldName := range fields {
		// Extract the actual field value from the object
		fieldValue, err := g.extractFieldFromObject(discriminant, fieldName)
		if err != nil {
			// If field extraction fails, bind to null/zero value
			fieldValue = constant.NewNull(types.I8Ptr)
		}

		g.variables[fieldName] = fieldValue
	}
}

// bindFieldsToZeroValues binds pattern fields to zero values when extraction fails
func (g *LLVMGenerator) bindFieldsToZeroValues(fields []string) {
	for _, fieldName := range fields {
		g.variables[fieldName] = constant.NewInt(types.I64, 0)
	}
}

// extractRecordFields extracts fields from a record type (single-variant struct)
func (g *LLVMGenerator) extractRecordFields(pattern ast.Pattern, discriminant value.Value, variant ast.TypeVariant) {
	// For record types, the discriminant is a pointer to the struct
	discriminantType := discriminant.Type()

	var (
		structType *types.StructType
		isPointer  bool
	)

	if ptrType, ok := discriminantType.(*types.PointerType); ok {
		if st, ok := ptrType.ElemType.(*types.StructType); ok {
			structType = st
			isPointer = true
		}
	} else if st, ok := discriminantType.(*types.StructType); ok {
		structType = st
		isPointer = false
	}

	if structType == nil {
		// If not a struct, bind fields to zero values
		g.bindFieldsToZeroValues(pattern.Fields)
		return
	}

	// Build field map from variant to match the struct generation logic
	fieldMap := g.buildFieldMapFromVariant(&variant)

	// Get the same field mapping that was used during struct creation
	// This ensures we use the same sorted field order
	fieldMapping := g.getOrCreateRecordFieldMapping(pattern.Constructor, fieldMap)

	// Extract fields from the struct using the correct mapping
	for _, patternFieldName := range pattern.Fields {
		// Get the actual struct index for this field name
		fieldIndex, exists := fieldMapping[patternFieldName]
		if !exists {
			// Field not found in type, bind to zero
			g.variables[patternFieldName] = constant.NewInt(types.I64, 0)
			continue
		}

		// Find the field info from the variant
		var field ast.TypeField

		for _, variantField := range variant.Fields {
			if variantField.Name == patternFieldName {
				field = variantField
				break
			}
		}

		if fieldIndex >= 0 && fieldIndex < len(structType.Fields) {
			// Get pointer to the field
			var fieldPtr value.Value
			if isPointer {
				// Discriminant is already a pointer to the struct
				fieldPtr = g.builder.NewGetElementPtr(
					structType,
					discriminant,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, int64(fieldIndex)),
				)
			} else {
				// Discriminant is a struct value, need to get its address first
				structAddr := g.builder.NewAlloca(structType)
				g.builder.NewStore(discriminant, structAddr)
				fieldPtr = g.builder.NewGetElementPtr(
					structType,
					structAddr,
					constant.NewInt(types.I32, 0),
					constant.NewInt(types.I32, int64(fieldIndex)),
				)
			}

			// Load the field value
			fieldType := structType.Fields[fieldIndex]
			fieldValue := g.builder.NewLoad(fieldType, fieldPtr)

			// Bind the pattern variable to the field value
			g.variables[patternFieldName] = fieldValue

			// Also register the variable in the Hindley-Milner type environment
			// For Result types, try to infer the semantic type of the value field
			var semanticType string
			if patternFieldName == "value" && pattern.Constructor == SuccessPattern {
				// This is likely a value extracted from a Result<T, E>
				// Try to infer the original type T from the context
				semanticType = g.inferResultValueType(discriminant, field.Type)
			} else {
				semanticType = field.Type
			}

			concreteType := &ConcreteType{name: semanticType}
			g.typeInferer.env.Set(patternFieldName, concreteType)
		}
	}
}

// inferResultValueType tries to infer the semantic type of a value extracted from a Result type
func (g *LLVMGenerator) inferResultValueType(_ value.Value, fieldType string) string {
	// Try to track back to the source of this Result value to determine the generic type parameter
	// This is a complex problem that requires proper generic type tracking
	// For now, we'll use heuristics based on the field type
	// Check if the field type suggests this could be a boolean value
	// Simple heuristic: if the field type suggests it could be boolean (i64 that might be boolean)
	// and this is a Success pattern match, assume it's a boolean
	// This handles the common case of Result<bool, Error> from comparison functions
	if fieldType == TypeInt || fieldType == "" {
		// This is likely a boolean value stored as i64 in the Result
		return TypeBool
	}

	return fieldType // Otherwise use the original field type
}

// normalizeArmValue handles Unit expressions in match arms
func (g *LLVMGenerator) normalizeArmValue(armValue value.Value) value.Value {
	// Handle Unit expressions in match arms
	// If the expression returns Unit, use void type
	if armValue == nil || armValue.Type() == types.Void {
		armValue = constant.NewUndef(types.Void)
	}

	return armValue
}

// addBranchTermination ensures proper branch termination for match arms
func (g *LLVMGenerator) addBranchTermination(endBlock *ir.Block) *ir.Block {
	// TODO: FIX THIS! DON'T IGNORE IT!!
	// After generating the expression (which might be a nested match),
	// the builder might be pointing to a different block. We need to ensure the
	// branch comes from the current builder block (where the expression ended),
	// but ONLY if that block doesn't already have a terminator.
	currentBuilderBlock := g.builder

	// Check if the current block already has a terminator instruction
	if currentBuilderBlock.Term == nil {
		currentBuilderBlock.NewBr(endBlock)
	}

	return currentBuilderBlock
}

// generateMatchConditions generates the conditional branches for pattern matching.
func (g *LLVMGenerator) generateMatchConditions(
	arms []ast.MatchArm,
	armBlocks []*ir.Block,
	discriminant value.Value,
	blockSuffix string,
) {
	currentBlock := g.builder

	for i, arm := range arms {
		// Ensure the builder is set to the current block
		g.builder = currentBlock
		condition := g.createPatternCondition(arm.Pattern, discriminant, currentBlock)

		if i == len(arms)-1 {
			currentBlock.NewBr(armBlocks[i])
		} else {
			nextCheckBlock := g.function.NewBlock(fmt.Sprintf("match_check_%d%s", i+1, blockSuffix))
			currentBlock.NewCondBr(condition, armBlocks[i], nextCheckBlock)
			currentBlock = nextCheckBlock
		}
	}
}

// createPatternCondition creates a condition for pattern matching.
func (g *LLVMGenerator) createPatternCondition(
	pattern ast.Pattern,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	// Handle wildcard and unknown patterns
	if pattern.Constructor == "_" || pattern.Constructor == UnknownPattern {
		return constant.NewBool(true)
	}

	// Handle structural matching
	if pattern.Constructor == "*" {
		return constant.NewBool(true)
	}

	// Handle variable binding patterns
	if pattern.Constructor == "" && pattern.Variable != "" {
		return constant.NewBool(true)
	}

	// Handle boolean patterns
	if pattern.Constructor == TruePattern || pattern.Constructor == FalsePattern {
		return g.createBooleanPatternCondition(pattern.Constructor, discriminant, currentBlock)
	}

	// Handle single-variant types FIRST (like PersonData { name, age })
	// Check if this constructor belongs to a single-variant type
	for _, typeDecl := range g.typeDeclarations {
		if len(typeDecl.Variants) == 1 && typeDecl.Variants[0].Name == pattern.Constructor {
			// For single-variant types, the constructor always matches
			// (there's only one possible variant)
			return constant.NewBool(true)
		}
	}

	// Handle union type variants (multi-variant types only)
	if discriminantValue, exists := g.unionVariants[pattern.Constructor]; exists {
		return g.createUnionPatternCondition(discriminantValue, discriminant, currentBlock)
	}

	// Handle numeric literals
	patternValue, err := strconv.ParseInt(pattern.Constructor, 10, 64)
	if err == nil {
		return g.createNumericPatternCondition(patternValue, discriminant, currentBlock)
	}

	// Default to string pattern condition
	return g.createStringPatternCondition(pattern.Constructor, discriminant, currentBlock)
}

// createBooleanPatternCondition handles boolean pattern matching
func (g *LLVMGenerator) createBooleanPatternCondition(
	constructor string,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	if constructor == TruePattern {
		return g.createTruePatternCondition(discriminant, currentBlock)
	}

	return g.createFalsePatternCondition(discriminant, currentBlock)
}

// createTruePatternCondition creates condition for "true" pattern
func (g *LLVMGenerator) createTruePatternCondition(discriminant value.Value, currentBlock *ir.Block) value.Value {
	// Handle struct types (boolean represented as struct)
	if _, isStruct := discriminant.Type().(*types.StructType); isStruct {
		return constant.NewBool(true)
	}

	// Handle integer types (check if non-zero)
	if intType, isInt := discriminant.Type().(*types.IntType); isInt {
		zero := constant.NewInt(intType, 0)
		return currentBlock.NewICmp(enum.IPredNE, discriminant, zero)
	}

	// Handle boolean types (compare directly)
	if discriminant.Type() == types.I1 {
		return discriminant
	}

	// Default fallback
	return constant.NewBool(true)
}

// createFalsePatternCondition creates condition for "false" pattern
func (g *LLVMGenerator) createFalsePatternCondition(discriminant value.Value, currentBlock *ir.Block) value.Value {
	// Handle struct types (boolean represented as struct)
	if _, isStruct := discriminant.Type().(*types.StructType); isStruct {
		return constant.NewBool(false)
	}

	// Handle pointer types (like strings)
	if _, isPtr := discriminant.Type().(*types.PointerType); isPtr {
		return constant.NewBool(false)
	}

	// Handle integer types (check if zero)
	if intType, isInt := discriminant.Type().(*types.IntType); isInt {
		zero := constant.NewInt(intType, 0)
		return currentBlock.NewICmp(enum.IPredEQ, discriminant, zero)
	}

	// Handle boolean types (negate)
	if discriminant.Type() == types.I1 {
		return currentBlock.NewXor(discriminant, constant.NewBool(true))
	}

	// Default fallback
	return constant.NewBool(false)
}

// createUnionPatternCondition handles union type variant pattern matching
func (g *LLVMGenerator) createUnionPatternCondition(
	discriminantValue int64,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	discriminantType := discriminant.Type()

	// Handle tagged union (pointer to struct with tag + data)
	if ptrType, ok := discriminantType.(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
			return g.createTaggedUnionCondition(
				discriminantValue, discriminant, structType, currentBlock,
			)
		}
	}

	// Fallback: simple enum discriminant (i64)
	return g.createSimpleEnumCondition(
		discriminantValue, discriminant, currentBlock,
	)
}

// createTaggedUnionCondition creates condition for tagged union
func (g *LLVMGenerator) createTaggedUnionCondition(
	discriminantValue int64,
	discriminant value.Value,
	structType *types.StructType,
	currentBlock *ir.Block,
) value.Value {
	// Extract tag field (index 0)
	tagPtr := currentBlock.NewGetElementPtr(structType, discriminant,
		constant.NewInt(types.I32, 0), // struct index
		constant.NewInt(types.I32, 0)) // tag field index
	tagValue := currentBlock.NewLoad(types.I8, tagPtr)

	// Convert discriminant value to i8 for comparison
	patternConst := constant.NewInt(types.I8, discriminantValue)

	return currentBlock.NewICmp(enum.IPredEQ, tagValue, patternConst)
}

// createSimpleEnumCondition creates condition for simple enum
func (g *LLVMGenerator) createSimpleEnumCondition(
	discriminantValue int64,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	// Only compare with integers if discriminant is actually an integer type
	if _, isInt := discriminant.Type().(*types.IntType); !isInt &&
		discriminant.Type() != types.I64 {
		return constant.NewBool(false)
	}

	patternConst := constant.NewInt(types.I64, discriminantValue)

	return currentBlock.NewICmp(enum.IPredEQ, discriminant, patternConst)
}

// createNumericPatternCondition handles numeric literal pattern matching
func (g *LLVMGenerator) createNumericPatternCondition(
	patternValue int64,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	// Only compare with integers if discriminant is actually an integer type
	if _, isInt := discriminant.Type().(*types.IntType); !isInt &&
		discriminant.Type() != types.I64 {
		return constant.NewBool(false)
	}

	patternConst := constant.NewInt(types.I64, patternValue)

	return currentBlock.NewICmp(enum.IPredEQ, discriminant, patternConst)
}

// extractFieldFromObject extracts a field value from an object literal
func (g *LLVMGenerator) extractFieldFromObject(objectValue value.Value, fieldName string) (value.Value, error) {
	// Check if this is a pointer to a struct
	objectType := objectValue.Type()

	var (
		structType *types.StructType
		isPointer  bool
	)

	if ptrType, ok := objectType.(*types.PointerType); ok {
		if st, ok := ptrType.ElemType.(*types.StructType); ok {
			structType = st
			isPointer = true
		}
	} else if st, ok := objectType.(*types.StructType); ok {
		structType = st
		isPointer = false
	}

	if structType == nil {
		// If not a struct, we can't extract fields
		return nil, fmt.Errorf("%w: %s", ErrCannotExtractField, fieldName)
	}

	// For now, assume the first field contains the value we want
	// In a real implementation, we'd need to map field names to indices
	fieldIndex := 0

	// Get pointer to the field
	var fieldPtr value.Value
	if isPointer {
		// Object is already a pointer to the struct
		fieldPtr = g.builder.NewGetElementPtr(
			structType,
			objectValue,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
	} else {
		// Object is a struct value, need to get its address first
		structAddr := g.builder.NewAlloca(structType)
		g.builder.NewStore(objectValue, structAddr)
		fieldPtr = g.builder.NewGetElementPtr(
			structType,
			structAddr,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(fieldIndex)),
		)
	}

	// Load the field value
	fieldType := structType.Fields[fieldIndex]
	fieldValue := g.builder.NewLoad(fieldType, fieldPtr)

	return fieldValue, nil
}

// extractDiscriminatedUnionFields extracts field values from a discriminated union variant
func (g *LLVMGenerator) extractDiscriminatedUnionFields(
	discriminant value.Value,
	pattern ast.Pattern,
	variables map[string]value.Value,
) error {
	// Check if this is a tagged union (pointer to struct with tag + data)
	discriminantType := discriminant.Type()

	ptrType, ok := discriminantType.(*types.PointerType)
	if !ok {
		return ErrDiscriminantNotPointer
	}

	structType, ok := ptrType.ElemType.(*types.StructType)
	if !ok || len(structType.Fields) != 2 {
		return ErrDiscriminantNotTaggedUnion
	}

	// Get pointer to the data area (second field in the tagged union)
	dataPtr := g.builder.NewGetElementPtr(structType, discriminant,
		constant.NewInt(types.I32, 0), // struct index
		constant.NewInt(types.I32, 1)) // data field index

	// Find the union type and variant information
	var variant *ast.TypeVariant

	// Find which union type this discriminant belongs to
	for _, typeDecl := range g.typeDeclarations {
		if len(typeDecl.Variants) > 1 {
			// Check if this pattern constructor matches any variant
			for _, v := range typeDecl.Variants {
				if v.Name == pattern.Constructor {
					variant = &v
					break
				}
			}

			if variant != nil {
				break
			}
		}
	}

	if variant == nil {
		return fmt.Errorf("%w: %s", ErrVariantNotFound, pattern.Constructor)
	}

	// Extract each field from the data area
	offset := int64(0)

	// Pattern fields are provided as a mapping from variant field names to pattern variable names
	// For example: pattern.Fields = ["d", "dur"] corresponds to variant fields ["damage", "durability"]
	// We need to map them correctly based on position/order

	for i, field := range variant.Fields {
		// Check if this field position has a corresponding pattern field
		var patternFieldName string

		fieldRequested := false

		if i < len(pattern.Fields) {
			patternFieldName = pattern.Fields[i]
			fieldRequested = true
		}

		if fieldRequested {
			// Get the field type and calculate its size
			fieldType := g.getFieldType(field.Type)

			// Get the actual data field type from the struct (second field)
			dataFieldType := structType.Fields[1] // Data field is at index 1

			// Cast data array to appropriate pointer type for this field
			fieldPtr := g.builder.NewBitCast(
				g.builder.NewGetElementPtr(
					dataFieldType, // Use the actual data field type
					dataPtr,
					constant.NewInt(types.I32, 0),      // array index
					constant.NewInt(types.I32, offset), // byte offset
				),
				types.NewPointer(fieldType),
			)

			// Load the field value
			fieldValue := g.builder.NewLoad(fieldType, fieldPtr)

			// Bind using the pattern variable name, not the variant field name
			variables[patternFieldName] = fieldValue

			// Also register the variable in the Hindley-Milner type environment
			// Infer the type from the field definition
			concreteType := &ConcreteType{name: field.Type}
			g.typeInferer.env.Set(patternFieldName, concreteType)
		}

		// Move to next field offset regardless of whether it was requested
		fieldSize := g.getTypeSize(g.getFieldType(field.Type))
		offset += fieldSize
	}

	return nil
}

// createStringPatternCondition creates a condition for string pattern matching.
func (g *LLVMGenerator) createStringPatternCondition(
	constructor string,
	discriminant value.Value,
	currentBlock *ir.Block,
) value.Value {
	if !strings.HasPrefix(constructor, "\"") || !strings.HasSuffix(constructor, "\"") {
		return constant.NewBool(false)
	}

	patternStr := constructor[1 : len(constructor)-1]
	patternConstStr := constant.NewCharArrayFromString(patternStr + StringTerminator)
	patternGlobal := g.module.NewGlobalDef("", patternConstStr)
	patternPtr := currentBlock.NewGetElementPtr(patternConstStr.Typ, patternGlobal,
		constant.NewInt(types.I32, ArrayIndexZero), constant.NewInt(types.I32, ArrayIndexZero))

	// Ensure discriminant is i8* for strcmp
	stringPtr := discriminant

	if discriminant.Type().String() != "i8*" {
		// If discriminant is not i8*, try to convert it
		if discriminant.Type().String() == "i64" {
			// This might be a string passed as i64 - need to cast or handle differently
			// For now, return false to avoid the strcmp issue
			return constant.NewBool(false)
		}
		// If it's another pointer type, try bitcast
		if _, isPtr := discriminant.Type().(*types.PointerType); isPtr {
			stringPtr = currentBlock.NewBitCast(discriminant, types.I8Ptr)
		} else {
			// Can't convert, pattern doesn't match
			return constant.NewBool(false)
		}
	}

	strcmp := g.functions["strcmp"]
	compareResult := currentBlock.NewCall(strcmp, stringPtr, patternPtr)
	zero := constant.NewInt(types.I32, ArrayIndexZero)

	return currentBlock.NewICmp(enum.IPredEQ, compareResult, zero)
}

// createMatchResult creates the final result value using PHI nodes if needed.
func (g *LLVMGenerator) createMatchResult(
	armValues []value.Value,
	predecessorBlocks []*ir.Block,
	endBlock *ir.Block,
) (value.Value, error) {
	g.builder = endBlock

	if len(armValues) == 1 {
		// For single arm, we still need to set the builder but don't need PHI
		return armValues[0], nil
	}

	// Check if all arm values are void type
	// PHI nodes cannot be created with void values
	allVoid := true

	for _, val := range armValues {
		if val.Type() != types.Void {
			allVoid = false
			break
		}
	}

	// If all arms return void, return void directly without PHI
	if allVoid {
		return constant.NewUndef(types.Void), nil
	}

	// Check if we need type coercion
	coercedValues, err := g.coerceArmValuesToCommonType(armValues)
	if err != nil {
		return nil, err
	}

	// Only include predecessors that actually have terminators
	var validIncomings []*ir.Incoming

	for i, val := range coercedValues {
		// Skip void values in PHI nodes
		if val.Type() != types.Void {
			// Check if the predecessor block has a terminator (meaning it actually branches)
			if i < len(predecessorBlocks) && predecessorBlocks[i].Term != nil {
				validIncomings = append(validIncomings, ir.NewIncoming(val, predecessorBlocks[i]))
			}
		}
	}

	// If no valid (non-void) values for PHI, return void
	if len(validIncomings) == 0 {
		return constant.NewUndef(types.Void), nil
	}

	// If only one valid incoming, don't create PHI
	if len(validIncomings) == 1 {
		return validIncomings[0].X, nil
	}

	// BUGFIX: Don't create PHI with void values - check if all values are void
	allVoid = true
	for _, incoming := range validIncomings {
		if !isVoidType(incoming.X.Type()) {
			allVoid = false
			break
		}
	}
	
	if allVoid && len(validIncomings) > 0 {
		// All arms return void - return nil to represent void, don't create PHI
		return nil, nil
	}

	phi := endBlock.NewPhi(validIncomings...)

	// The end block now has a PHI node and the builder is set to this block.
	// The calling function (like generateStandardMatchExpression) should handle
	// adding any necessary terminator when the match is used in a larger context.

	return phi, nil
}

// coerceArmValuesToCommonType ensures all arm values have compatible types.
func (g *LLVMGenerator) coerceArmValuesToCommonType(armValues []value.Value) ([]value.Value, error) {
	expectedType := armValues[0].Type()

	if !g.needsTypeCoercion(armValues, expectedType) {
		return armValues, nil
	}

	return g.performTypeCoercion(armValues, expectedType)
}

// needsTypeCoercion checks if type coercion is needed for arm values.
func (g *LLVMGenerator) needsTypeCoercion(armValues []value.Value, expectedType types.Type) bool {
	for i := 1; i < len(armValues); i++ {
		if armValues[i].Type() != expectedType {
			return true
		}
	}

	return false
}

// performTypeCoercion converts all values to the expected type.
func (g *LLVMGenerator) performTypeCoercion(armValues []value.Value, expectedType types.Type) ([]value.Value, error) {
	coercedValues := make([]value.Value, len(armValues))

	for i, val := range armValues {
		if val.Type() == expectedType {
			coercedValues[i] = val
		} else {
			coercedVal, err := g.coerceValueToType(val, expectedType)
			if err != nil {
				return nil, err
			}

			coercedValues[i] = coercedVal
		}
	}

	return coercedValues, nil
}

// coerceValueToType converts a single value to the target type.
func (g *LLVMGenerator) coerceValueToType(val value.Value, targetType types.Type) (value.Value, error) {
	switch targetType {
	case types.I8Ptr:
		return g.convertToString(val)
	case types.I64:
		return constant.NewInt(types.I64, 0), nil
	default:
		return val, nil
	}
}

// convertToString converts a value to string type.
func (g *LLVMGenerator) convertToString(val value.Value) (value.Value, error) {
	if val.Type() == types.I64 {
		return g.generateIntToString(val)
	}

	return val, nil
}

func (g *LLVMGenerator) generateResultMatchExpression(
	matchExpr *ast.MatchExpression,
	discriminant value.Value,
) (value.Value, error) {
	blocks := g.createResultMatchBlocks(matchExpr)

	// Store the Result value for pattern binding
	g.currentResultValue = discriminant

	g.generateResultMatchCondition(discriminant, blocks)

	// Track which blocks actually branch to the end
	var (
		actualSuccessBlock *ir.Block
		actualErrorBlock   *ir.Block
	)

	// Generate success block and track the actual ending block
	g.builder = blocks.Success

	successValue, err := g.generateSuccessBlock(matchExpr, blocks)
	if err != nil {
		return nil, err
	}
	// The builder is now pointing to the block that will branch to the end
	actualSuccessBlock = g.builder

	// Generate error block and track the actual ending block
	g.builder = blocks.Error

	errorValue, err := g.generateErrorBlock(matchExpr, blocks)
	if err != nil {
		return nil, err
	}
	// The builder is now pointing to the block that will branch to the end
	actualErrorBlock = g.builder

	// Create PHI with the actual predecessor blocks
	return g.createResultMatchPhiWithActualBlocks(successValue, errorValue, actualSuccessBlock, actualErrorBlock, blocks)
}

// ResultMatchBlocks holds the blocks for result match expressions.
type ResultMatchBlocks struct {
	Success *ir.Block
	Error   *ir.Block
	End     *ir.Block
}

// createResultMatchBlocks creates blocks for result match expressions.
func (g *LLVMGenerator) createResultMatchBlocks(matchExpr *ast.MatchExpression) *ResultMatchBlocks {
	blockSuffix := fmt.Sprintf("_%p", matchExpr)

	return &ResultMatchBlocks{
		Success: g.function.NewBlock("success_case" + blockSuffix),
		Error:   g.function.NewBlock("error_case" + blockSuffix),
		End:     g.function.NewBlock("match_end" + blockSuffix),
	}
}

// generateResultMatchCondition generates the condition for result matching.
func (g *LLVMGenerator) generateResultMatchCondition(discriminant value.Value, blocks *ResultMatchBlocks) {
	// Handle both struct values and pointers to structs
	if ptrType, ok := discriminant.Type().(*types.PointerType); ok {
		// Extract the discriminant field from the Result struct
		// Result struct: [value, discriminant] where discriminant is at index 1
		resultType := ptrType.ElemType
		discriminantPtr := g.builder.NewGetElementPtr(resultType, discriminant,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
		discriminantValue := g.builder.NewLoad(discriminantPtr.Type().(*types.PointerType).ElemType, discriminantPtr)

		// 0 = Success, 1 = Error
		zero := constant.NewInt(types.I8, 0)
		isSuccess := g.builder.NewICmp(enum.IPredEQ, discriminantValue, zero)
		g.builder.NewCondBr(isSuccess, blocks.Success, blocks.Error)
	} else if _, ok := discriminant.Type().(*types.StructType); ok {
		// Handle struct value directly (for value semantics)
		// Extract discriminant field (index 1) from struct value
		discriminantValue := g.builder.NewExtractValue(discriminant, 1)

		// 0 = Success, 1 = Error
		zero := constant.NewInt(types.I8, 0)
		isSuccess := g.builder.NewICmp(enum.IPredEQ, discriminantValue, zero)
		g.builder.NewCondBr(isSuccess, blocks.Success, blocks.Error)
	} else {
		// Fallback: treat as integer discriminant
		zero := constant.NewInt(types.I64, 0)
		isSuccess := g.builder.NewICmp(enum.IPredSGE, discriminant, zero)
		g.builder.NewCondBr(isSuccess, blocks.Success, blocks.Error)
	}
}

// generateSuccessBlock generates the success block for result matching.
func (g *LLVMGenerator) generateSuccessBlock(
	matchExpr *ast.MatchExpression,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.Success

	// Find the success arm and bind pattern variables
	successArm := g.findSuccessArm(matchExpr)
	if successArm != nil && len(successArm.Pattern.Fields) > 0 {
		// Bind the Result value to the pattern variable
		fieldName := successArm.Pattern.Fields[0] // First field is the value

		// Bind the extracted value type to the pattern variable
		g.bindPatternVariableType(fieldName, matchExpr.Expression)

		// Get the Result value from the matched expression
		// The Result struct has: [value, discriminant]
		// We need to extract the value field (index 0)
		if g.currentResultValue != nil {
			if ptrType, ok := g.currentResultValue.Type().(*types.PointerType); ok {
				resultType := ptrType.ElemType
				valuePtr := g.builder.NewGetElementPtr(resultType, g.currentResultValue,
					constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
				extractedValue := g.builder.NewLoad(valuePtr.Type().(*types.PointerType).ElemType, valuePtr)
				g.variables[fieldName] = extractedValue
			} else if _, ok := g.currentResultValue.Type().(*types.StructType); ok {
				// Handle struct value directly (for value semantics)
				// Extract value field (index 0) from struct value
				extractedValue := g.builder.NewExtractValue(g.currentResultValue, 0)
				g.variables[fieldName] = extractedValue
				// Note: Type information is handled elsewhere for struct values
			} else {
				// Fallback: use the discriminant value directly
				g.variables[fieldName] = g.currentResultValue
			}
		}
	}

	successExpr := g.findSuccessValue(matchExpr)

	var successValue value.Value

	if successExpr != nil {
		// Generate the expression (which might be a nested match)
		val, err := g.generateExpression(successExpr)
		if err != nil {
			return nil, err
		}

		successValue = val

		// After generating a nested expression, the builder might have changed
		// We need to ensure the branch to the end block comes from the correct block
		// But only add the branch if the current block doesn't already have a terminator
		if g.builder.Term == nil {
			g.builder.NewBr(blocks.End)
		}

		// For PHI node creation, we need to track which block actually branches to the end
		// This might be different from the original success block if we had nested expressions
		return successValue, nil
	}

	// Fallback: use the bound variable from pattern matching
	if successArm := g.findSuccessArm(matchExpr); successArm != nil && len(successArm.Pattern.Fields) > 0 {
		fieldName := successArm.Pattern.Fields[0]
		if extractedValue, exists := g.variables[fieldName]; exists {
			successValue = extractedValue
		} else {
			successValue = constant.NewInt(types.I64, ArrayIndexZero)
		}
	} else {
		successValue = constant.NewInt(types.I64, ArrayIndexZero)
	}

	// Only add branch if the current block doesn't already have a terminator
	if g.builder.Term == nil {
		g.builder.NewBr(blocks.End)
	}

	return successValue, nil
}

// generateErrorBlock generates the error block for result matching.
func (g *LLVMGenerator) generateErrorBlock(
	matchExpr *ast.MatchExpression,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.Error

	// Find the Error arm and bind pattern variables
	errorArm := g.findErrorArm(matchExpr)
	if errorArm != nil && len(errorArm.Pattern.Fields) > 0 {
		// Bind the Result error message to the pattern variable
		fieldName := errorArm.Pattern.Fields[0] // First field is the message
		// Create a unique global string for the error message
		// Include function context to ensure uniqueness across monomorphized instances
		funcContext := ""
		if g.function != nil {
			funcContext = g.function.Name()
		}
		blockSuffix := fmt.Sprintf("_%s_%p", funcContext, matchExpr)
		errorStr := g.module.NewGlobalDef("error_msg"+blockSuffix, constant.NewCharArrayFromString("Error occurred\\x00"))
		errorPtr := g.builder.NewGetElementPtr(errorStr.ContentType, errorStr,
			constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
		g.variables[fieldName] = errorPtr
	}

	errorExpr := g.findErrorValue(matchExpr)

	var errorValue value.Value

	if errorExpr != nil {
		// Generate the expression (which might be a nested match)
		val, err := g.generateExpression(errorExpr)
		if err != nil {
			return nil, err
		}

		errorValue = val

		// After generating a nested expression, the builder might have changed
		// We need to ensure the branch to the end block comes from the correct block
		// But only add the branch if the current block doesn't already have a terminator
		if g.builder.Term == nil {
			g.builder.NewBr(blocks.End)
		}

		return errorValue, nil
	}

	// Fallback: use the bound variable from pattern matching
	if errorArm := g.findErrorArm(matchExpr); errorArm != nil && len(errorArm.Pattern.Fields) > 0 {
		fieldName := errorArm.Pattern.Fields[0]
		if extractedError, exists := g.variables[fieldName]; exists {
			errorValue = extractedError
		} else {
			errorValue = constant.NewInt(types.I64, ArrayIndexZero)
		}
	} else {
		errorValue = constant.NewInt(types.I64, ArrayIndexZero)
	}

	// Only add branch if the current block doesn't already have a terminator
	if g.builder.Term == nil {
		g.builder.NewBr(blocks.End)
	}

	return errorValue, nil
}

// bindPatternVariableType binds the correct type for a pattern variable extracted from a Result
func (g *LLVMGenerator) bindPatternVariableType(fieldName string, matchedExpr ast.Expression) {
	// Infer the type of the matched expression to get the Result type
	matchedExprType, err := g.typeInferer.InferType(matchedExpr)
	if err == nil {
		resolvedType := g.typeInferer.ResolveType(matchedExprType)
		if genericType, ok := resolvedType.(*GenericType); ok {
			if genericType.name == TypeResult && len(genericType.typeArgs) >= 1 {
				// Extract the success type (first type argument of Result<T, E>)
				successType := genericType.typeArgs[0]
				g.typeInferer.env.Set(fieldName, successType)
			}
		}
	}
}

// findSuccessArm finds the success match arm.
func (g *LLVMGenerator) findSuccessArm(matchExpr *ast.MatchExpression) *ast.MatchArm {
	for i, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == SuccessPattern {
			return &matchExpr.Arms[i]
		}
	}

	return nil
}

// findErrorArm finds the error match arm.
func (g *LLVMGenerator) findErrorArm(matchExpr *ast.MatchExpression) *ast.MatchArm {
	for i, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == ErrorPattern {
			return &matchExpr.Arms[i]
		}
	}

	return nil
}

// findSuccessValue finds the success expression in match arms.
func (g *LLVMGenerator) findSuccessValue(matchExpr *ast.MatchExpression) ast.Expression {
	for _, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == SuccessPattern {
			return arm.Expression
		}
	}

	return nil
}

// findErrorValue finds the error expression in match arms.
func (g *LLVMGenerator) findErrorValue(matchExpr *ast.MatchExpression) ast.Expression {
	for _, arm := range matchExpr.Arms {
		if arm.Pattern.Constructor == ErrorPattern {
			return arm.Expression
		}
	}

	return nil
}

// createResultMatchPhiWithActualBlocks creates the PHI node for result matching with actual predecessor blocks.
func (g *LLVMGenerator) createResultMatchPhiWithActualBlocks(
	successValue, errorValue value.Value,
	actualSuccessBlock, actualErrorBlock *ir.Block,
	blocks *ResultMatchBlocks,
) (value.Value, error) {
	g.builder = blocks.End

	// Use the actual blocks that branch to the end
	var validPredecessors []*ir.Incoming

	// Check if the actual success block has a terminator and branches to end
	if actualSuccessBlock != nil && actualSuccessBlock.Term != nil && successValue != nil {
		// Don't add void values to PHI predecessors
		if !isVoidType(successValue.Type()) {
			validPredecessors = append(validPredecessors, ir.NewIncoming(successValue, actualSuccessBlock))
		}
	}

	// Check if the actual error block has a terminator and branches to end
	if actualErrorBlock != nil && actualErrorBlock.Term != nil && errorValue != nil {
		// Don't add void values to PHI predecessors
		if !isVoidType(errorValue.Type()) {
			validPredecessors = append(validPredecessors, ir.NewIncoming(errorValue, actualErrorBlock))
		}
	}

	// If we don't have valid predecessors, return a default value
	if len(validPredecessors) == 0 {
		// Return the success value as a fallback
		return successValue, nil
	}

	// If we only have one valid predecessor, don't create a PHI
	if len(validPredecessors) == 1 {
		return validPredecessors[0].X, nil
	}

	// BUGFIX: Check if both values are void (Unit) - can't create PHI with void values
	if successValue != nil && errorValue != nil {
		// Check if both are void types (nil represents void/Unit)
		successIsVoid := (successValue == nil) || isVoidType(successValue.Type())
		errorIsVoid := (errorValue == nil) || isVoidType(errorValue.Type())
		
		if successIsVoid && errorIsVoid {
			// Both arms return Unit - return nil to represent void, don't create PHI
			return nil, nil
		}
	}

	// Create PHI node with valid predecessors
	phi := blocks.End.NewPhi(validPredecessors...)

	return phi, nil
}

// isVoidType checks if a type represents void/Unit
func isVoidType(t types.Type) bool {
	// In LLVM, void is represented as nil or specific void type
	if t == nil {
		return true
	}
	// Check if it's LLVM void type
	if _, ok := t.(*types.VoidType); ok {
		return true
	}
	return false
}

// reInferReturnType re-infers the return type of a function with concrete parameter types
func (g *LLVMGenerator) reInferReturnType(funcName string, argTypes []Type) (Type, error) {
	// Get the original function declaration
	fnDecl, exists := g.functionDeclarations[funcName]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrFunctionNotDeclared, funcName)
	}

	// Save current type environment
	savedEnv := g.typeInferer.env.Clone()

	// Set up a new environment with concrete parameter types
	for i, param := range fnDecl.Parameters {
		if i < len(argTypes) {
			g.typeInferer.env.Set(param.Name, argTypes[i])
		}
	}

	// Re-infer the return type from the function body
	bodyType, err := g.typeInferer.InferType(fnDecl.Body)
	if err != nil {
		// Restore environment on error
		g.typeInferer.env = savedEnv
		return nil, err
	}

	// Restore environment
	g.typeInferer.env = savedEnv

	return g.typeInferer.prune(bodyType), nil
}

// setParameterTypeContext sets the expected parameter type context for argument generation
func (g *LLVMGenerator) setParameterTypeContext(funcValue value.Value, paramIndex int) {
	if funcValue != nil {
		if funcType, ok := funcValue.Type().(*types.FuncType); ok && paramIndex < len(funcType.Params) {
			g.expectedParameterType = funcType.Params[paramIndex]
		}
	}
}

// clearParameterTypeContext clears the parameter type context
func (g *LLVMGenerator) clearParameterTypeContext() {
	g.expectedParameterType = nil
}
