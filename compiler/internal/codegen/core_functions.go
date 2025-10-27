package codegen

import (
	"fmt"
	"strings"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// validateBuiltInArgs validates argument count for built-in functions using the registry
func validateBuiltInArgs(funcName string, callExpr *ast.CallExpression) error {
	fn, exists := GlobalBuiltInRegistry.GetFunction(funcName)
	if !exists {
		return WrapUndefinedFunction(funcName)
	}

	if len(callExpr.Arguments) != len(fn.ParameterTypes) {
		return WrapFunctionArgsWithPos(funcName, len(fn.ParameterTypes), len(callExpr.Arguments), callExpr.Position)
	}

	return nil
}

// generateToStringCall handles toString function calls.
func (g *LLVMGenerator) generateToStringCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 1 {
		return nil, WrapToStringWrongArgs(len(callExpr.Arguments))
	}

	arg, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// TODO: check if result type. The value should have type info attached to it
	// If not, something has gone wrong
	if g.isResultType(arg) {
		if ptrType, ok := arg.Type().(*types.PointerType); ok {
			if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
				return g.convertResultToString(arg, structType)
			}
		}
		// Handle struct value directly (not pointer)
		if structType, ok := arg.Type().(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
			return g.convertResultToString(arg, structType)
		}
	}

	inferredType, err := g.typeInferer.InferType(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	//WRONG! TODO: We can't pass the type around as a string. We need a more robust solution.

	// CRITICAL: Resolve the type to get concrete type instead of type variables
	resolvedType := g.typeInferer.ResolveType(inferredType)
	argType := resolvedType.String()

	// If the resolved type is still a type variable (like t16), fall back to LLVM type
	if strings.HasPrefix(argType, "t") && len(argType) > 1 {
		// Check if it's a type variable (starts with 't' followed by digits)
		isTypeVar := true

		for _, c := range argType[1:] {
			if c < '0' || c > '9' {
				isTypeVar = false
				break
			}
		}

		if isTypeVar {
			// Fallback to determine type from LLVM value
			switch arg.Type() {
			case types.I64:
				argType = TypeInt
			case types.Double:
				argType = TypeFloat
			case types.I8Ptr:
				argType = TypeString
			case types.I1:
				argType = TypeBool
			default:
				// Default to int for unknown types
				argType = TypeInt
			}
		}
	}

	// Check if the resolved type IS boolean (not Result<bool, Error>!)
	isBooleanType := g.isSemanticBooleanType(resolvedType)
	if isBooleanType {
		argType = TypeBool
	}

	return g.convertValueToStringByType(argType, arg)
}

// TODO: This is wrong. We cannot convert fibers to string unless they return a String or
// there is a toString implementation
//
//nolint:gocognit // TODO: Refactor this function to reduce complexity
func (g *LLVMGenerator) convertValueToStringByType(
	//TODO: types must not be passed around as strings. This is wrong.
	theType string, arg value.Value) (value.Value, error) {
	// TODO: unhard code this!!! DO NOT IGNORE THIS! FIX IT!!
	// but the actual LLVM value is a plain int, treat as int
	if theType == "Result<int, Error>" && arg.Type() == types.I64 {
		return g.generateIntToString(arg)
	}

	switch theType {
	case TypeString:
		return arg, nil // Identity conversion
	case TypeInt, TypeAny: // Treat Any as Int at LLVM level for toString
		return g.generateIntToString(arg)
	case TypeBool:
		return g.generateBoolToString(arg)
	case TypeFloat:
		return g.generateFloatToString(arg)
	case TypeUnit:
		// Unit type should return "()"
		return g.createGlobalString("()"), nil
	default:
		// Check if it's a Fiber type - show the fiber ID, not await the result
		if theType == TypeFiber || strings.HasPrefix(theType, TypeFiber+"[") {
			// Fiber is just an integer ID, convert it to string
			return g.generateIntToString(arg)
		}

		// Check if it's a Channel type - show the channel ID, not a generic string
		if theType == TypeChannel {
			// Channel is just an integer ID, convert it to string
			return g.generateIntToString(arg)
		}

		// Check if it's a Result type (with either angle or square brackets)
		if strings.HasPrefix(theType, "Result<") || strings.HasPrefix(theType, "Result[") {
			// For Result types, check if it's a struct pointer
			if ptrType, ok := arg.Type().(*types.PointerType); ok {
				if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
					return g.convertResultToString(arg, structType)
				}
			}
			// Also handle struct value directly (not pointer)
			if structType, ok := arg.Type().(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
				return g.convertResultToString(arg, structType)
			}

			// AUTO-UNWRAP FIX: If inferred type is Result but actual LLVM value is NOT a Result struct,
			// the value has been auto-unwrapped per spec (0004-TypeSystem.md:115-160).
			// Convert the unwrapped value directly based on its LLVM type.
			// Example: fn double(x) = x * 2 returns unwrapped i64, not Result struct
			switch arg.Type() {
			case types.I64:
				return g.generateIntToString(arg)
			case types.Double:
				return g.generateFloatToString(arg)
			case types.I8Ptr:
				return arg, nil // Already a string
			case types.I1:
				return g.generateBoolToString(arg)
			}
		}

		// Fallback: Check LLVM type directly for Result-like structs (2-field struct with i64 and i8)
		// This handles cases where type inference has unresolved type variables
		if structType, ok := arg.Type().(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
			// Check if it looks like a Result struct: {i64, i8}
			if structType.Fields[1] == types.I8 {
				return g.convertResultToString(arg, structType)
			}
		}

		// For other complex types, return a generic representation
		return g.createGlobalString(fmt.Sprintf("<%s>", theType)), nil
	}
}

// convertResultToString extracts the value from a Result type and converts it to string
func (g *LLVMGenerator) convertResultToString(
	result value.Value, structType *types.StructType,
) (value.Value, error) {
	var discriminant value.Value
	var resultValue value.Value

	// Handle both pointer and value cases
	if _, ok := result.Type().(*types.PointerType); ok {
		// Pointer case: use getelementptr and load
		discriminantPtr := g.builder.NewGetElementPtr(structType, result,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
		discriminant = g.builder.NewLoad(types.I8, discriminantPtr)
	} else {
		// Struct value case: use extractvalue
		discriminant = g.builder.NewExtractValue(result, 1)
		resultValue = g.builder.NewExtractValue(result, 0)
	}

	// Check if discriminant == 0 (Success)
	zero := constant.NewInt(types.I8, 0)
	isSuccess := g.builder.NewICmp(enum.IPredEQ, discriminant, zero)

	// Create blocks with unique names to avoid conflicts
	blockID := len(g.function.Blocks) // Use block count as unique ID
	successBlockName := fmt.Sprintf("result_toString_success_%d", blockID)
	errorBlockName := fmt.Sprintf("result_toString_error_%d", blockID)
	endBlockName := fmt.Sprintf("result_toString_end_%d", blockID)

	successBlock := g.function.NewBlock(successBlockName)
	errorBlock := g.function.NewBlock(errorBlockName)
	endBlock := g.function.NewBlock(endBlockName)

	g.builder.NewCondBr(isSuccess, successBlock, errorBlock)

	// Declare sprintf and malloc for formatting (used by both success and error cases)
	sprintf := g.ensureSprintfDeclaration()
	malloc := g.ensureMallocDeclaration()
	bufferSize := constant.NewInt(types.I64, BufferSize64Bytes)

	// Success case: extract and convert the value
	g.builder = successBlock

	// Get the value - handle both pointer and struct cases
	if _, ok := result.Type().(*types.PointerType); ok {
		// Pointer case: use getelementptr and load
		valuePtr := g.builder.NewGetElementPtr(structType, result,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		resultValue = g.builder.NewLoad(structType.Fields[0], valuePtr)
	}
	// Struct value case: resultValue was already extracted above with NewExtractValue

	var (
		successStr value.Value
		err        error
	)

	// Convert the success value to string and wrap in Success(...)
	// CRITICAL: Result types MUST always format as "Success(value)" or "Error(message)"
	var innerValueStr value.Value
	switch structType.Fields[0] {
	case types.I64:
		innerValueStr, err = g.generateIntToString(resultValue)
	case types.Double:
		// Float value - convert to string
		innerValueStr, err = g.generateFloatToString(resultValue)
	case types.I1:
		innerValueStr, err = g.generateBoolToString(resultValue)
	case types.I8Ptr:
		innerValueStr = resultValue // Already a string
	default:
		// For complex types (like ProcessHandle), convert to a generic string
		innerValueStr = g.createGlobalString("complex_value")
	}

	if err != nil {
		return nil, err
	}

	// Wrap the value in Success(...) format
	successFormatStr := g.createGlobalString("Success(%s)")
	successBuffer := g.builder.NewCall(malloc, bufferSize)
	g.builder.NewCall(sprintf, successBuffer, successFormatStr, innerValueStr)
	successStr = successBuffer

	successBlock.NewBr(endBlock)

	// Error case: format as "Error(message)"
	g.builder = errorBlock

	// Extract the error message from the Result struct
	var errorMsg value.Value
	if _, ok := result.Type().(*types.PointerType); ok {
		// Pointer case: use getelementptr and load
		errorMsgPtr := g.builder.NewGetElementPtr(structType, result,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		errorMsg = g.builder.NewLoad(structType.Fields[0], errorMsgPtr)
	} else {
		// Struct value case: extract the value field
		errorMsg = g.builder.NewExtractValue(result, 0)
	}

	// Format as "Error(message)" - handle different error message types
	var errorStr value.Value
	if structType.Fields[0] == types.I8Ptr {
		// String error message - format as Error(message)
		errorFormatStr := g.createGlobalString("Error(%s)")
		errorBuffer := g.builder.NewCall(malloc, bufferSize)
		g.builder.NewCall(sprintf, errorBuffer, errorFormatStr, errorMsg)
		errorStr = errorBuffer
	} else {
		// Non-string error - just use "Error" for now
		// TODO: Handle other error types properly
		errorStr = g.createGlobalString("Error")
	}

	errorBlock.NewBr(endBlock)

	// End block: PHI node to select result
	g.builder = endBlock
	phi := endBlock.NewPhi(
		ir.NewIncoming(successStr, successBlock),
		ir.NewIncoming(errorStr, errorBlock),
	)

	return phi, nil
}

// createGlobalString creates a global string constant and returns a pointer to it
func (g *LLVMGenerator) createGlobalString(str string) value.Value {
	strConstant := constant.NewCharArrayFromString(str + "\x00")
	global := g.module.NewGlobalDef("", strConstant)

	return g.builder.NewGetElementPtr(strConstant.Typ, global,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
}

// isSemanticBooleanType checks if the inferred type is semantically a boolean
func (g *LLVMGenerator) isSemanticBooleanType(inferredType Type) bool {
	if concrete, ok := inferredType.(*ConcreteType); ok {
		return concrete.name == TypeBool
	}

	return false
}

// generatePrintCall handles print function calls.
func (g *LLVMGenerator) generatePrintCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(PrintFunc, callExpr)
	if err != nil {
		return nil, err
	}

	argExpr := callExpr.Arguments[0]
	arg, err := g.generateExpression(argExpr)
	if err != nil {
		return nil, err
	}

	inferredType, err := g.typeInferer.InferType(argExpr)
	if err != nil {
		return nil, err
	}

	stringArg, err := g.convertValueToStringForPrint(arg, inferredType)
	if err != nil {
		return nil, err
	}

	puts := g.functions["puts"]
	g.builder.NewCall(puts, stringArg)
	return nil, nil
}

// convertValueToStringForPrint converts any value to a string for printing.
func (g *LLVMGenerator) convertValueToStringForPrint(arg value.Value, inferredType Type) (value.Value, error) {
	if g.isResultType(arg) {
		return g.convertResultValueToString(arg)
	}
	return g.convertPrimitiveToString(arg, inferredType)
}

// convertResultValueToString handles Result type conversion to string.
func (g *LLVMGenerator) convertResultValueToString(arg value.Value) (value.Value, error) {
	if structType, ok := arg.Type().(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
		return g.convertResultToString(arg, structType)
	}
	if ptrType, ok := arg.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == ResultFieldCount {
			return g.convertResultToString(arg, structType)
		}
	}
	return nil, ErrPrintCannotConvert
}

// convertPrimitiveToString handles primitive type conversion to string.
func (g *LLVMGenerator) convertPrimitiveToString(arg value.Value, inferredType Type) (value.Value, error) {
	switch arg.Type().(type) {
	case *types.PointerType:
		return arg, nil
	case *types.IntType:
		return g.convertIntTypeToString(arg, inferredType)
	case *types.FloatType:
		return g.generateFloatToString(arg)
	default:
		return nil, ErrPrintCannotConvert
	}
}

// convertIntTypeToString handles int type conversion, distinguishing between bool and int.
func (g *LLVMGenerator) convertIntTypeToString(arg value.Value, inferredType Type) (value.Value, error) {
	intType := arg.Type().(*types.IntType)
	if intType.BitSize == 1 || g.isSemanticBooleanType(inferredType) {
		return g.generateBoolToString(arg)
	}
	return g.generateIntToString(arg)
}

// generateInputCall handles input function calls.
func (g *LLVMGenerator) generateInputCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(InputFunc, callExpr)
	if err != nil {
		return nil, err
	}

	// Declare fgets function if not already declared
	fgetsFunc, ok := g.functions["fgets"]
	if !ok {
		fgetsFunc = g.module.NewFunc("fgets", types.I8Ptr,
			ir.NewParam("str", types.I8Ptr),
			ir.NewParam("size", types.I32),
			ir.NewParam("stream", types.I8Ptr))
		g.functions["fgets"] = fgetsFunc
	}

	// Use stdin directly - it's an external global in libc
	stdinGlobal := g.module.NewGlobalDef("stdin", constant.NewNull(types.I8Ptr))

	// Allocate buffer for input string (256 chars should be enough)
	const inputBufferSize = 256

	bufferSize := constant.NewInt(types.I32, inputBufferSize)
	inputBuffer := g.builder.NewAlloca(types.NewArray(inputBufferSize, types.I8))

	// Cast array to i8* for fgets
	bufferPtr := g.builder.NewGetElementPtr(types.NewArray(inputBufferSize, types.I8), inputBuffer,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))

	// Call fgets to read the string
	fgetsResult := g.builder.NewCall(fgetsFunc, bufferPtr, bufferSize, stdinGlobal)

	// Create a Result<String, Error>
	resultType := g.getResultType(types.I8Ptr)
	result := g.builder.NewAlloca(resultType)

	// Check if fgets succeeded (returns non-null pointer)
	nullPtr := constant.NewNull(types.I8Ptr)
	fgetsSucceeded := g.builder.NewICmp(enum.IPredNE, fgetsResult, nullPtr)

	// Create blocks for success and error cases
	successBlock := g.function.NewBlock("input_success")
	errorBlock := g.function.NewBlock("input_error")
	endBlock := g.function.NewBlock("input_end")

	g.builder.NewCondBr(fgetsSucceeded, successBlock, errorBlock)

	// Success case: store the input string
	g.builder = successBlock

	// Remove trailing newline if present using strlen and string manipulation
	strlenFunc, ok := g.functions["strlen"]
	if !ok {
		strlenFunc = g.module.NewFunc("strlen", types.I64, ir.NewParam("str", types.I8Ptr))
		g.functions["strlen"] = strlenFunc
	}

	// Get string length
	strLength := successBlock.NewCall(strlenFunc, bufferPtr)

	// Check if last character is newline (ASCII 10) and remove it
	one := constant.NewInt(types.I64, 1)
	lastCharIdx := successBlock.NewSub(strLength, one)
	lastCharPtr := successBlock.NewGetElementPtr(types.I8, bufferPtr, lastCharIdx)
	lastChar := successBlock.NewLoad(types.I8, lastCharPtr)

	const asciiNewline = 10

	newlineChar := constant.NewInt(types.I8, asciiNewline) // ASCII newline
	isNewline := successBlock.NewICmp(enum.IPredEQ, lastChar, newlineChar)

	// Replace newline with null terminator if it exists
	nullChar := constant.NewInt(types.I8, 0)
	// Conditionally replace the newline character with null terminator
	charToStore := successBlock.NewSelect(isNewline, nullChar, lastChar)
	successBlock.NewStore(charToStore, lastCharPtr)

	valuePtr := successBlock.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	successBlock.NewStore(bufferPtr, valuePtr)
	discriminantPtr := successBlock.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	successBlock.NewStore(constant.NewInt(types.I8, 0), discriminantPtr) // 0 for Success
	successBlock.NewBr(endBlock)

	// Error case: store error discriminant
	g.builder = errorBlock
	errorDiscriminantPtr := errorBlock.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	errorBlock.NewStore(constant.NewInt(types.I8, 1), errorDiscriminantPtr) // 1 for Error
	// Set value to null for error case
	errorValuePtr := errorBlock.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	errorBlock.NewStore(nullPtr, errorValuePtr)
	errorBlock.NewBr(endBlock)

	// Continue with end block
	g.builder = endBlock

	return result, nil
}

func (g *LLVMGenerator) generateLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(LengthFunc, callExpr)
	if err != nil {
		return nil, err
	}

	arg, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Declare or get the strlen function
	strlenFunc, ok := g.functions["strlen"]
	if !ok {
		strlenFunc = g.module.NewFunc("strlen", types.I64, ir.NewParam("str", types.I8Ptr))
		g.functions["strlen"] = strlenFunc
	}

	// Call strlen(arg) and return the result directly (no Result wrapper)
	length := g.builder.NewCall(strlenFunc, arg)

	return length, nil
}

func (g *LLVMGenerator) getResultType(valueType types.Type) *types.StructType {
	// A Result is a struct { value, discriminant }
	return types.NewStruct(valueType, types.I8)
}

// generateContainsCall handles contains(haystack: string, needle: string) -> bool function calls.
func (g *LLVMGenerator) generateContainsCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(ContainsFunc, callExpr)
	if err != nil {
		return nil, err
	}

	haystack, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	needle, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Declare or get the strstr function
	strstrFunc, ok := g.functions["strstr"]
	if !ok {
		strstrFunc = g.module.NewFunc("strstr", types.I8Ptr,
			ir.NewParam("haystack", types.I8Ptr),
			ir.NewParam("needle", types.I8Ptr))
		g.functions["strstr"] = strstrFunc
	}

	// Call strstr(haystack, needle)
	resultPtr := g.builder.NewCall(strstrFunc, haystack, needle)

	// Check if result is not null (convert to bool)
	nullPtr := constant.NewNull(types.I8Ptr)
	isNotNull := g.builder.NewICmp(enum.IPredNE, resultPtr, nullPtr)

	// Create a Result<Bool, NoError>
	resultType := g.getResultType(types.I1)
	result := g.builder.NewAlloca(resultType)

	// Store the boolean in the value field
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(isNotNull, valuePtr)

	// Store the discriminant (0 for Success)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	return result, nil
}

func (g *LLVMGenerator) generateSubstringCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(SubstringFunc, callExpr)
	if err != nil {
		return nil, err
	}

	str, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	start, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	end, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	// Calculate length: end - start
	length := g.builder.NewSub(end, start)

	// Allocate memory for the substring (length + 1 for null terminator)
	lengthPlusOne := g.builder.NewAdd(length, constant.NewInt(types.I64, 1))

	mallocFunc, ok := g.functions["malloc"]
	if !ok {
		mallocFunc = g.module.NewFunc("malloc", types.I8Ptr, ir.NewParam("size", types.I64))
		g.functions["malloc"] = mallocFunc
	}

	newStr := g.builder.NewCall(mallocFunc, lengthPlusOne)

	// Calculate source pointer: str + start
	srcPtr := g.builder.NewGetElementPtr(types.I8, str, start)

	// Copy the substring using memcpy
	memcpyFunc, ok := g.functions["memcpy"]
	if !ok {
		memcpyFunc = g.module.NewFunc("memcpy", types.I8Ptr,
			ir.NewParam("dest", types.I8Ptr),
			ir.NewParam("src", types.I8Ptr),
			ir.NewParam("n", types.I64))
		g.functions["memcpy"] = memcpyFunc
	}

	g.builder.NewCall(memcpyFunc, newStr, srcPtr, length)

	// Null-terminate the new string
	endPtr := g.builder.NewGetElementPtr(types.I8, newStr, length)
	g.builder.NewStore(constant.NewInt(types.I8, 0), endPtr)

	// Create a Result<String, NoError>
	resultType := g.getResultType(types.I8Ptr)
	result := g.builder.NewAlloca(resultType)

	// Store the new string in the value field
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(newStr, valuePtr)

	// Store the discriminant (0 for Success)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	return result, nil
}

// generateParseIntCall handles parseInt(s: string) -> Result<int, string> function calls.
func (g *LLVMGenerator) generateParseIntCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(ParseIntFunc, callExpr)
	if err != nil {
		return nil, err
	}

	str, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Declare or get the atoll function (ASCII to long long)
	atollFunc, ok := g.functions["atoll"]
	if !ok {
		atollFunc = g.module.NewFunc("atoll", types.I64, ir.NewParam("str", types.I8Ptr))
		g.functions["atoll"] = atollFunc
	}

	// Call atoll(str) to convert string to integer
	parsedValue := g.builder.NewCall(atollFunc, str)

	// TODO: Add proper error checking - atoll returns 0 for invalid strings
	// For now, assume parsing always succeeds

	// Create a Result<int, string>
	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)

	// Store the parsed integer in the value field
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(parsedValue, valuePtr)

	// Store the discriminant (0 for Success)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	return result, nil
}

// generateJoinCall handles join(list: List<string>, separator: string) -> string function calls.
func (g *LLVMGenerator) generateJoinCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(JoinFunc, callExpr)
	if err != nil {
		return nil, err
	}

	// For now, return a placeholder implementation
	// TODO: Implement proper list handling once List<T> type is fully supported

	separator, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Return the separator as a placeholder result
	// In a real implementation, we would iterate through the list and join with separator
	return separator, nil
}

// generateListConstructorCall handles List() constructor calls.
func (g *LLVMGenerator) generateListConstructorCall(_ *ast.CallExpression) (value.Value, error) {
	// For now, return a simple placeholder (null pointer)
	// TODO: Implement proper dynamic list structure
	return constant.NewNull(types.I8Ptr), nil
}

// generateMapConstructorCall handles Map() constructor calls.
func (g *LLVMGenerator) generateMapConstructorCall(_ *ast.CallExpression) (value.Value, error) {
	// For now, return a simple placeholder (null pointer)
	// TODO: Implement proper dynamic map structure
	return constant.NewNull(types.I8Ptr), nil
}
