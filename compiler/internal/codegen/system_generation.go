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

// There is an insane amount of duplication in this file.
// TODO: Merge a bunch of this code!
// TODO: Code REUSE!!!!

// -----------------------------------------------------------------------------
// Built-in helper dispatchers used from iterator_generation.go
// -----------------------------------------------------------------------------

// callFunctionWithValue calls a built-in or user-defined function with a single
// LLVM IR value argument. It is used by iterator helpers such as forEach,
// filter and map.
func (g *LLVMGenerator) callFunctionWithValue(
	funcIdent *ast.Identifier,
	val value.Value,
) (value.Value, error) {
	switch funcIdent.Name {
	// Built-ins that accept a single argument.
	case PrintFunc:
		return g.callBuiltInPrint(val)
	case ToStringFunc:
		return g.callBuiltInToString(val)
	case InputFunc:
		return nil, ErrInputNoArgs
	case "testAny": // Helper used by validation tests.
		return g.generateTestAnyCall()
	}

	// User-defined function.
	if fn, ok := g.functions[funcIdent.Name]; ok {
		return g.builder.NewCall(fn, val), nil
	}

	return nil, WrapFunctionNotFound(funcIdent.Name)
}

// callFunctionWithTwoValues is similar to callFunctionWithValue but passes two
// arguments to the callee.
func (g *LLVMGenerator) callFunctionWithTwoValues(
	funcIdent *ast.Identifier,
	val1, val2 value.Value,
) (value.Value, error) {
	// None of the current built-ins legitimately take two positional
	// arguments.  Treat any attempt as an error so that misuse is surfaced
	// during compilation rather than silently ignored.
	switch funcIdent.Name {
	case PrintFunc, ToStringFunc, InputFunc:
		return nil, WrapBuiltInTwoArgs(funcIdent.Name)
	}

	if fn, ok := g.functions[funcIdent.Name]; ok {
		return g.builder.NewCall(fn, val1, val2), nil
	}

	return nil, WrapFunctionNotFound(funcIdent.Name)
}

// callBuiltInPrint prints a single LLVM IR value using the C library puts
// function. It first converts any integer or boolean values to strings.
func (g *LLVMGenerator) callBuiltInPrint(val value.Value) (value.Value, error) {
	var strArg value.Value
	var err error

	switch val.Type().(type) {
	case *types.PointerType: // i8* – already a C-string.
		strArg = val
	case *types.IntType:
		// i1 == bool, otherwise treat as int64.
		if val.Type().(*types.IntType).BitSize == 1 {
			strArg, err = g.generateBoolToString(val)
		} else {
			strArg, err = g.generateIntToString(val)
		}
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrPrintCannotConvert
	}

	puts := g.functions["puts"]
	res := g.builder.NewCall(puts, strArg)
	return g.builder.NewSExt(res, types.I64), nil // normalise to i64
}

// callBuiltInToString converts the provided value to its string representation
// and returns it.
func (g *LLVMGenerator) callBuiltInToString(val value.Value) (value.Value, error) {
	switch val.Type().(type) {
	case *types.PointerType:
		return val, nil // Already a string.
	case *types.IntType:
		if val.Type().(*types.IntType).BitSize == 1 {
			return g.generateBoolToString(val)
		}
		return g.generateIntToString(val)
	default:
		return nil, ErrNoToStringImpl
	}
}

// generateTestAnyCall returns a constant placeholder value used in unit tests
// for validating that the compiler correctly handles the `any` type.
func (g *LLVMGenerator) generateTestAnyCall() (value.Value, error) {
	return constant.NewInt(types.I64, DefaultPlaceholder), nil
}

// -----------------------------------------------------------------------------
// Process / File / JSON built-ins
// -----------------------------------------------------------------------------

// generateSpawnProcessCall emits a call to the simple process spawning function
// that uses a default callback for stdout/stderr events.
func (g *LLVMGenerator) generateSpawnProcessCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapFunctionArgsWithPos(
			SpawnProcessFunc, SpawnProcessExpectedArgs,
			len(callExpr.Arguments), callExpr.Position)
	}

	cmd, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	callback, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Use the callback-based process spawning function (MANDATORY callback)
	var fn *ir.Func
	if existing, ok := g.functions["spawn_process_with_handler"]; ok {
		fn = existing
	} else {
		// int64_t spawn_process_with_handler(char *command, ProcessEventHandler handler)
		// ProcessEventHandler: void (*)(int64_t process_id, int64_t event_type, char *data)
		handlerType := types.NewPointer(types.NewFunc(
			types.Void,  // Return type: void (matches C runtime)
			types.I64,   // process_id parameter
			types.I64,   // event_type parameter
			types.I8Ptr, // data parameter
		))

		fn = g.module.NewFunc("spawn_process_with_handler", types.I64,
			ir.NewParam("command", types.I8Ptr),
			ir.NewParam("handler", handlerType))
		g.functions["spawn_process_with_handler"] = fn
	}

	// Call the callback-based process spawning function
	processID := g.builder.NewCall(fn, cmd, callback)

	// Create a Result<Int, String> - process ID on success, error message on failure
	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)

	// Initialize the result struct with default values
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))

	// Check if the process ID is negative (error)
	zero := constant.NewInt(types.I64, 0)
	isError := g.builder.NewICmp(enum.IPredSLT, processID, zero)

	// Create blocks with unique names to avoid conflicts
	blockID := len(g.function.Blocks) // Use block count as unique ID
	successBlockName := fmt.Sprintf("spawn_success_%d", blockID)
	errorBlockName := fmt.Sprintf("spawn_error_%d", blockID)
	continueBlockName := fmt.Sprintf("spawn_continue_%d", blockID)

	successBlock := g.function.NewBlock(successBlockName)
	errorBlock := g.function.NewBlock(errorBlockName)
	continueBlock := g.function.NewBlock(continueBlockName)

	g.builder.NewCondBr(isError, errorBlock, successBlock)

	// Success case: store the process ID
	g.builder = successBlock
	g.builder.NewStore(processID, valuePtr)
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr) // 0 = Success
	g.builder.NewBr(continueBlock)

	// Error case: store error indicator (we'll use -1 as the "error value")
	g.builder = errorBlock
	g.builder.NewStore(constant.NewInt(types.I64, -1), valuePtr)
	g.builder.NewStore(constant.NewInt(types.I8, 1), discriminantPtr) // 1 = Error
	g.builder.NewBr(continueBlock)

	// Continue execution
	g.builder = continueBlock

	return result, nil
}

// generateSleepCall emits a call to the fiber_sleep function to sleep for the specified milliseconds.
func (g *LLVMGenerator) generateSleepCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapSleepWrongArgs(len(callExpr.Arguments))
	}

	milliseconds, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get the fiber_sleep function
	sleepFunc := g.functions["fiber_sleep"]
	if sleepFunc == nil {
		g.initFiberRuntime()
		sleepFunc = g.functions["fiber_sleep"]
	}

	return g.builder.NewCall(sleepFunc, milliseconds), nil
}

// generateWriteFileCall writes data to a file via an external helper returning
// a Result<Success, string> type.
func (g *LLVMGenerator) generateWriteFileCall(callExpr *ast.CallExpression) (value.Value, error) {
	const expectedArgs = TwoArgs
	if len(callExpr.Arguments) != expectedArgs {
		return nil, WrapFunctionArgsWithPos(WriteFileFunc, WriteFileExpectedArgs, len(callExpr.Arguments), callExpr.Position)
	}

	filename, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}
	content, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	var fn *ir.Func
	if existing, ok := g.functions["write_file"]; ok {
		fn = existing
	} else {
		fn = g.module.NewFunc("write_file", types.I64,
			ir.NewParam("filename", types.I8Ptr),
			ir.NewParam("content", types.I8Ptr))
		g.functions["write_file"] = fn
	}

	// Call the C function
	writeResult := g.builder.NewCall(fn, filename, content)

	// CRITICAL FIX: Create the Result allocation in the current block before branching
	resultType := g.getResultType(types.I64)
	result := g.builder.NewAlloca(resultType)

	// Check if the result is negative (error)
	zero := constant.NewInt(types.I64, 0)
	isError := g.builder.NewICmp(enum.IPredSLT, writeResult, zero)

	// Create blocks with unique names to avoid conflicts
	blockID := len(g.function.Blocks) // Use block count as unique ID
	successBlockName := fmt.Sprintf("write_success_%d", blockID)
	errorBlockName := fmt.Sprintf("write_error_%d", blockID)
	continueBlockName := fmt.Sprintf("write_continue_%d", blockID)

	successBlock := g.function.NewBlock(successBlockName)
	errorBlock := g.function.NewBlock(errorBlockName)
	continueBlock := g.function.NewBlock(continueBlockName)

	g.builder.NewCondBr(isError, errorBlock, successBlock)

	// Success case: store the bytes written
	g.builder = successBlock
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(writeResult, valuePtr)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr) // 0 = Success
	g.builder.NewBr(continueBlock)

	// Error case: store error value
	g.builder = errorBlock
	errorValuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(constant.NewInt(types.I64, -1), errorValuePtr)
	errorDiscriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 1), errorDiscriminantPtr) // 1 = Error
	g.builder.NewBr(continueBlock)

	// Continue execution
	g.builder = continueBlock

	return result, nil
}

// generateReadFileCall reads the entire contents of the specified file and
// returns a Result<string, string> type.
func (g *LLVMGenerator) generateReadFileCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapFunctionArgsWithPos(ReadFileFunc, ReadFileExpectedArgs, len(callExpr.Arguments), callExpr.Position)
	}

	filename, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	var fn *ir.Func
	if existing, ok := g.functions["read_file"]; ok {
		fn = existing
	} else {
		fn = g.module.NewFunc("read_file", types.I8Ptr, ir.NewParam("filename", types.I8Ptr))
		g.functions["read_file"] = fn
	}

	// Call the C function
	readResult := g.builder.NewCall(fn, filename)

	// Create a Result<string, string>
	resultType := g.getResultType(types.I8Ptr)
	result := g.builder.NewAlloca(resultType)

	// Check if the result is NULL (error)
	nullPtr := constant.NewNull(types.I8Ptr)
	isError := g.builder.NewICmp(enum.IPredEQ, readResult, nullPtr)

	// Create blocks with unique names to avoid conflicts (same fix as writeFile)
	blockID := len(g.function.Blocks) // Use block count as unique ID
	successBlockName := fmt.Sprintf("read_success_%d", blockID)
	errorBlockName := fmt.Sprintf("read_error_%d", blockID)
	continueBlockName := fmt.Sprintf("read_continue_%d", blockID)

	successBlock := g.function.NewBlock(successBlockName)
	errorBlock := g.function.NewBlock(errorBlockName)
	continueBlock := g.function.NewBlock(continueBlockName)

	g.builder.NewCondBr(isError, errorBlock, successBlock)

	// Success case: store the file content
	g.builder = successBlock
	valuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(readResult, valuePtr)
	discriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 0), discriminantPtr) // 0 = Success
	g.builder.NewBr(continueBlock)

	// Error case: store error placeholder
	g.builder = errorBlock
	// Create unique global name to avoid redefinition
	globalName := fmt.Sprintf("read_error_msg_%p", callExpr)
	errorStr := g.module.NewGlobalDef(globalName, constant.NewCharArrayFromString("File read error\x00"))
	errorPtr := g.builder.NewGetElementPtr(errorStr.ContentType, errorStr,
		constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))
	errorValuePtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(errorPtr, errorValuePtr)
	errorDiscriminantPtr := g.builder.NewGetElementPtr(resultType, result,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(constant.NewInt(types.I8, 1), errorDiscriminantPtr) // 1 = Error
	g.builder.NewBr(continueBlock)

	// Continue execution
	g.builder = continueBlock

	return result, nil
}

// Process management functions for async process handling

// generateAwaitProcessCall waits for process completion
func (g *LLVMGenerator) generateAwaitProcessCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapAwaitProcessWrongArgs(len(callExpr.Arguments))
	}

	processID, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	var fn *ir.Func
	if existing, ok := g.functions["fiber_await_process"]; ok {
		fn = existing
	} else {
		fn = g.module.NewFunc("fiber_await_process", types.I64, ir.NewParam("process_id", types.I64))
		g.functions["fiber_await_process"] = fn
	}

	return g.builder.NewCall(fn, processID), nil
}

// generateCleanupProcessCall cleans up process resources
func (g *LLVMGenerator) generateCleanupProcessCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapCleanupProcessWrongArgs(len(callExpr.Arguments))
	}

	processID, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	var fn *ir.Func
	if existing, ok := g.functions["fiber_cleanup_process"]; ok {
		fn = existing
	} else {
		fn = g.module.NewFunc("fiber_cleanup_process", types.Void, ir.NewParam("process_id", types.I64))
		g.functions["fiber_cleanup_process"] = fn
	}

	g.builder.NewCall(fn, processID)
	return constant.NewInt(types.I64, 0), nil // Return success
}
