package codegen

import (
	"errors"

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

// generatePrintCall handles print function calls.
func (g *LLVMGenerator) generatePrintCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 1 {
		return nil, WrapPrintWrongArgs(len(callExpr.Arguments))
	}

	argExpr := callExpr.Arguments[0]
	arg, err := g.generateExpression(argExpr)
	if err != nil {
		return nil, err
	}

	var stringArg value.Value
	switch arg.Type().(type) {
	case *types.PointerType: // Assuming i8* is string
		stringArg = arg
	case *types.IntType:
		if arg.Type().(*types.IntType).BitSize == 1 { // bool (i1)
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
		return nil, errors.New("unsupported type for print")
	}

	puts := g.functions["puts"]
	result := g.builder.NewCall(puts, stringArg)
	// puts returns i32, but our functions should return i64
	return g.builder.NewSExt(result, types.I64), nil
}

// generateInputCall handles input function calls.
func (g *LLVMGenerator) generateInputCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 0 {
		return nil, WrapInputWrongArgs(len(callExpr.Arguments))
	}
	// ... (rest of the function)
	return nil, errors.New("unsupported call expression")
}

func (g *LLVMGenerator) generateLengthCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapLengthWrongArgs(len(callExpr.Arguments))
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
		return nil, WrapContainsWrongArgs(len(callExpr.Arguments))
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
		return nil, WrapSubstringWrongArgs(len(callExpr.Arguments))
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
