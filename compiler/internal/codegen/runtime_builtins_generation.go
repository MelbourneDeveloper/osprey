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

// Declared Osprey return types understood by the generic C-runtime bridge.
const (
	runtimeRetResultString = "Result<string, Error>"
	runtimeRetResultInt    = "Result<int, Error>"
	runtimeRetRawString    = "string"
)

// generateRuntimeBuiltinCall is the shared generator for the HTTP-response,
// JSON and terminal builtins. It looks the function up in the registry, declares
// the backing C function with its raw return type (i8* for string-producing
// calls, i64 otherwise), emits the call, and wraps the result according to the
// registry's declared return type. Implements [HTTP-RESPONSE-HANDLE],
// [BUILTIN-JSON] and [BUILTIN-TERM].
func (g *LLVMGenerator) generateRuntimeBuiltinCall(callExpr *ast.CallExpression) (value.Value, error) {
	ident, ok := callExpr.Function.(*ast.Identifier)
	if !ok {
		return nil, WrapUndefinedFunction("<runtime builtin>")
	}

	name := ident.Name
	builtinFunc, exists := GlobalBuiltInRegistry.GetFunction(name)
	if !exists {
		return nil, WrapUndefinedFunction(name)
	}

	expectedArgs := len(builtinFunc.ParameterTypes)
	if len(callExpr.Arguments) != expectedArgs {
		return nil, WrapWrongArgCount(name, expectedArgs, len(callExpr.Arguments))
	}

	args := make([]value.Value, expectedArgs)
	for i, arg := range callExpr.Arguments {
		val, err := g.generateExpression(arg)
		if err != nil {
			return nil, err
		}
		args[i] = val
	}

	retName := ""
	if ct, isConcrete := builtinFunc.ReturnType.(*ConcreteType); isConcrete {
		retName = ct.name
	}

	rawRet := types.Type(types.I64)
	if retName == runtimeRetResultString || retName == runtimeRetRawString {
		rawRet = types.I8Ptr
	}

	fn := g.ensureRawCFunction(builtinFunc.CName, rawRet, builtinFunc.ParameterTypes)
	call := g.builder.NewCall(fn, args...)

	switch retName {
	case runtimeRetResultString:
		return g.wrapStringResult(call, name), nil
	case runtimeRetResultInt:
		return g.wrapIntResult(call), nil
	default:
		// Raw int / raw string returns are passed straight through.
		return call, nil
	}
}

// ensureRawCFunction declares (once) a C runtime function with the given raw
// return type and parameter types taken from the registry entry.
func (g *LLVMGenerator) ensureRawCFunction(cName string, retType types.Type,
	params []BuiltInParameter) *ir.Func {
	if fn, ok := g.functions[cName]; ok {
		return fn
	}

	llvmParams := make([]*ir.Param, len(params))
	for i, param := range params {
		llvmParams[i] = ir.NewParam(param.Name, g.getLLVMType(param.Type))
	}

	fn := g.module.NewFunc(cName, retType, llvmParams...)
	g.functions[cName] = fn

	return fn
}

// wrapStringResult turns an i8* C return into Result<string, Error>: NULL maps
// to Error, anything else to Success. Mirrors generateReadFileCall.
func (g *LLVMGenerator) wrapStringResult(callValue value.Value, label string) value.Value {
	resultType := g.getResultType(types.I8Ptr)
	result := g.builder.NewAlloca(resultType)

	nullPtr := constant.NewNull(types.I8Ptr)
	isError := g.builder.NewICmp(enum.IPredEQ, callValue, nullPtr)

	blockID := len(g.function.Blocks)
	successBlock := g.function.NewBlock(fmt.Sprintf("rt_str_succ_%d", blockID))
	errorBlock := g.function.NewBlock(fmt.Sprintf("rt_str_err_%d", blockID))
	continueBlock := g.function.NewBlock(fmt.Sprintf("rt_str_cont_%d", blockID))

	g.builder.NewCondBr(isError, errorBlock, successBlock)

	g.builder = successBlock
	g.storeResultFields(resultType, result, callValue, 0)
	g.builder.NewBr(continueBlock)

	g.builder = errorBlock
	errMsg := g.createGlobalString(label + " failed")
	g.storeResultFields(resultType, result, errMsg, 1)
	g.builder.NewBr(continueBlock)

	g.builder = continueBlock

	return result
}

// wrapIntResult turns an i64 C return into Result<int, Error>: a negative value
// maps to Error (the code is preserved), zero or positive to Success. Mirrors
// generateWriteFileCall.
func (g *LLVMGenerator) wrapIntResult(callValue value.Value) value.Value {
	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)

	zero := constant.NewInt(types.I64, 0)
	isError := g.builder.NewICmp(enum.IPredSLT, callValue, zero)

	blockID := len(g.function.Blocks)
	successBlock := g.function.NewBlock(fmt.Sprintf("rt_int_succ_%d", blockID))
	errorBlock := g.function.NewBlock(fmt.Sprintf("rt_int_err_%d", blockID))
	continueBlock := g.function.NewBlock(fmt.Sprintf("rt_int_cont_%d", blockID))

	g.builder.NewCondBr(isError, errorBlock, successBlock)

	g.builder = successBlock
	g.storeResultFields(resultType, result, callValue, 0)
	g.builder.NewBr(continueBlock)

	g.builder = errorBlock
	g.storeResultFields(resultType, result, callValue, 1)
	g.builder.NewBr(continueBlock)

	g.builder = continueBlock

	return result
}

// storeResultFields writes the value field and the discriminant (0=Success,
// 1=Error) of a Result struct allocation.
func (g *LLVMGenerator) storeResultFields(resultType *types.StructType,
	result value.Value, fieldValue value.Value, discriminant int64) {
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(fieldValue, valuePtr)

	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, discriminant), discriminantPtr)
}
