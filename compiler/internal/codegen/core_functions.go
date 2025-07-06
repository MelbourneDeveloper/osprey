package codegen

import (
	"fmt"
	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// generateToStringCall handles toString function calls.
func (g *LLVMGenerator) generateToStringCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 1 {
		return nil, WrapToStringWrongArgs(len(callExpr.Arguments))
	}

	arg, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Check if this is a Result type (struct pointer)
	if ptrType, ok := arg.Type().(*types.PointerType); ok {
		if structType, ok := ptrType.ElemType.(*types.StructType); ok && len(structType.Fields) == 2 {
			// This is a Result type - extract the value and convert to string
			return g.convertResultToString(arg, structType)
		}
	}

	argType, err := g.getTypeOfExpression(callExpr.Arguments[0])
	if err != nil {
		// Fallback for literals without explicit type info
		switch arg.Type().(type) {
		case *types.IntType:
			argType = TypeInt
		case *types.PointerType:
			argType = TypeString
		default:
			return nil, err
		}
	}

	return g.convertValueToStringByType(argType, arg)
}

func (g *LLVMGenerator) convertValueToStringByType(theType string, arg value.Value) (value.Value, error) {
	switch theType {
	case TypeString:
		return arg, nil // Identity conversion
	case TypeInt, TypeAny: // Treat Any as Int at LLVM level for toString
		return g.generateIntToString(arg)
	case TypeBool:
		return g.generateBoolToString(arg)
	default:
		return nil, WrapNoToStringImpl(theType)
	}
}

// convertResultToString extracts the value from a Result type and converts it to string
func (g *LLVMGenerator) convertResultToString(
	resultPtr value.Value, structType *types.StructType,
) (value.Value, error) {
	// Check the discriminant first to see if it's Success or Error
	discriminantPtr := g.builder.NewGetElementPtr(structType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	discriminant := g.builder.NewLoad(types.I8, discriminantPtr)

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

	// Success case: extract and convert the value
	g.builder = successBlock
	valuePtr := g.builder.NewGetElementPtr(structType, resultPtr,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	resultValue := g.builder.NewLoad(structType.Fields[0], valuePtr)

	var successStr value.Value
	var err error

	// Convert based on the value type
	switch structType.Fields[0] {
	case types.I64:
		successStr, err = g.generateIntToString(resultValue)
	case types.I1:
		successStr, err = g.generateBoolToString(resultValue)
	case types.I8Ptr:
		successStr = resultValue // Already a string
	default:
		// For complex types (like ProcessHandle), convert to a generic string
		successStr = g.createGlobalString("Success")
	}

	if err != nil {
		return nil, err
	}

	successBlock.NewBr(endBlock)

	// Error case: return "Error"
	g.builder = errorBlock
	errorStr := g.createGlobalString("Error")
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

// generatePrintCall handles print function calls.
func (g *LLVMGenerator) generatePrintCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 1 {
		return nil, WrapPrintWrongArgsWithPos(len(callExpr.Arguments), callExpr.Position)
	}

	argExpr := callExpr.Arguments[0]
	arg, err := g.generateExpression(argExpr)
	if err != nil {
		return nil, err
	}

	// Try to get the semantic type of the expression
	argType, typeErr := g.getTypeOfExpression(argExpr)

	var stringArg value.Value
	switch arg.Type().(type) {
	case *types.PointerType: // Assuming i8* is string
		stringArg = arg
	case *types.IntType:
		// Check if this is semantically a boolean (either i1 or i64 with bool type)
		if arg.Type().(*types.IntType).BitSize == 1 || (typeErr == nil && argType == TypeBool) {
			stringArg, err = g.generateBoolToString(arg)
			if err != nil {
				return nil, err
			}
		} else { // int (i64)
			stringArg, err = g.generateIntToString(arg)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, ErrPrintCannotConvert
	}

	puts := g.functions["puts"]
	result := g.builder.NewCall(puts, stringArg)
	// puts returns i32, but our functions should return i64
	return g.builder.NewSExt(result, types.I64), nil
}

// generateInputCall handles input function calls.
func (g *LLVMGenerator) generateInputCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 0 {
		return nil, WrapInputWrongArgsWithPos(len(callExpr.Arguments), callExpr.Position)
	}
	// ... (rest of the function)
	return nil, ErrUnsupportedCall
}

func (g *LLVMGenerator) generateLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapLengthWrongArgsWithPos(len(callExpr.Arguments), callExpr.Position)
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

	// Call strlen(arg)
	length := g.builder.NewCall(strlenFunc, arg)

	// Create a Result<Int, NoError>
	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)

	// Store the length in the value field
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(length, valuePtr)

	// Store the discriminant (0 for Success)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr)

	return result, nil
}

func (g *LLVMGenerator) getResultType(valueType types.Type) *types.StructType {
	// A Result is a struct { value, discriminant }
	return types.NewStruct(valueType, types.I8)
}

// generateContainsCall handles contains(haystack: string, needle: string) -> bool function calls.
func (g *LLVMGenerator) generateContainsCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapContainsWrongArgsWithPos(len(callExpr.Arguments), callExpr.Position)
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
	if len(callExpr.Arguments) != ThreeArgs {
		return nil, WrapSubstringWrongArgsWithPos(len(callExpr.Arguments), callExpr.Position)
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
