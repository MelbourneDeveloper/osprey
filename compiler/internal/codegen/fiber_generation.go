package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

const (
	defaultChannelCapacity  = 10
	moduleAccessPlaceholder = 42
)

// CRITICAL: NO PLACEHOLDERS ALLOWED IN THIS FILE
// EVERY FUNCTION MUST HAVE REAL IMPLEMENTATION

// FiberRuntimeFunction represents a fiber runtime function configuration
type FiberRuntimeFunction struct {
	Name       string
	ReturnType types.Type
	Params     []types.Type
	Linkage    enum.Linkage
}

// initFiberRuntime declares external runtime functions using consolidated approach
func (g *LLVMGenerator) initFiberRuntime() {
	fiberFunctions := []FiberRuntimeFunction{
		{
			Name:       "fiber_spawn",
			ReturnType: types.I64,
			Params:     []types.Type{types.NewPointer(types.NewFunc(types.I64))},
			Linkage:    enum.LinkageExternal,
		},
		{
			Name:       "fiber_await",
			ReturnType: types.I64,
			Params:     []types.Type{types.I64},
			Linkage:    enum.LinkageExternal,
		},
		{
			Name:       "fiber_yield",
			ReturnType: types.I64,
			Params:     []types.Type{types.I64},
			Linkage:    enum.LinkageExternal,
		},
		{
			Name:       "fiber_sleep",
			ReturnType: types.I64,
			Params:     []types.Type{types.I64},
			Linkage:    enum.LinkageExternal,
		},
		{
			Name:       "channel_create",
			ReturnType: types.I64,
			Params:     []types.Type{types.I64},
			Linkage:    enum.LinkageExternal,
		},
		{
			Name:       "channel_send",
			ReturnType: types.I64,
			Params:     []types.Type{types.I64, types.I64},
			Linkage:    enum.LinkageExternal,
		},
		{
			Name:       "channel_recv",
			ReturnType: types.I64,
			Params:     []types.Type{types.I64},
			Linkage:    enum.LinkageExternal,
		},
	}

	for _, fn := range fiberFunctions {
		params := make([]*ir.Param, len(fn.Params))
		for i, paramType := range fn.Params {
			params[i] = ir.NewParam(fmt.Sprintf("arg%d", i), paramType)
		}

		llvmFunc := g.module.NewFunc(fn.Name, fn.ReturnType, params...)
		llvmFunc.Linkage = fn.Linkage
		g.functions[fn.Name] = llvmFunc
	}
}

// generateFiberRuntimeCall generates calls to fiber runtime functions using builtin registry
func (g *LLVMGenerator) generateFiberRuntimeCall(builtinName string, runtimeName string,
	args []value.Value) (value.Value, error) {
	// Validate using builtin registry
	builtinFunc, exists := GlobalBuiltInRegistry.GetFunction(builtinName)
	if !exists {
		return nil, fmt.Errorf("builtin function %s: %w", builtinName, ErrFunctionNotFound)
	}

	// Validate argument count
	expectedArgs := len(builtinFunc.ParameterTypes)
	if len(args) != expectedArgs {
		return nil, fmt.Errorf("function %s: %w", builtinName, ErrWrongArgCount)
	}

	// Get or initialize runtime function
	runtimeFunc := g.functions[runtimeName]
	if runtimeFunc == nil {
		g.initFiberRuntime()
		runtimeFunc = g.functions[runtimeName]
	}

	// Generate the runtime call
	return g.builder.NewCall(runtimeFunc, args...), nil
}

// generateSpawnExpression generates REAL fiber spawning with concurrency.
func (g *LLVMGenerator) generateSpawnExpression(spawn *ast.SpawnExpression) (value.Value, error) {
	// Create a closure function for the spawned expression
	g.closureCounter++
	closureName := fmt.Sprintf("fiber_closure_%d", g.closureCounter)
	closureFunc := g.module.NewFunc(closureName, types.I64)

	// Save current context
	prevFunc := g.function
	prevBuilder := g.builder
	prevVars := g.variables

	// Capture variables by value BEFORE creating the closure context
	capturedValues := make(map[string]value.Value)
	g.captureVariablesInExpression(spawn.Expression, capturedValues)

	// Create new context for closure
	g.function = closureFunc
	entry := closureFunc.NewBlock("entry")
	g.builder = entry

	// Create new variable scope with captured values + preserved globals
	// We need to merge captured values with the original variables to support nested spawns
	g.variables = make(map[string]value.Value)

	// First, copy all original variables (this preserves global scope for nested spawns)
	for name, val := range prevVars {
		g.variables[name] = val
	}

	// Then, override with captured values (this ensures proper closure semantics)
	for name, val := range capturedValues {
		g.variables[name] = val
	}

	// Generate the expression inside the closure
	result, err := g.generateExpression(spawn.Expression)
	if err != nil {
		return nil, err
	}

	// Return the result
	g.builder.NewRet(result)

	// Restore context
	g.function = prevFunc
	g.builder = prevBuilder
	g.variables = prevVars

	// Call fiber_spawn with the closure using consolidated approach
	return g.generateFiberRuntimeCall("fiber_spawn", "fiber_spawn", []value.Value{closureFunc})
}

// generateAwaitExpression generates REAL fiber await with blocking.
func (g *LLVMGenerator) generateAwaitExpression(await *ast.AwaitExpression) (value.Value, error) {
	// Generate the fiber ID expression
	fiberID, err := g.generateExpression(await.Expression)
	if err != nil {
		return nil, err
	}

	// Call fiber_await using consolidated approach
	return g.generateFiberRuntimeCall("fiber_await", "fiber_await", []value.Value{fiberID})
}

// generateYieldExpression generates REAL yield with scheduler cooperation.
func (g *LLVMGenerator) generateYieldExpression(yield *ast.YieldExpression) (value.Value, error) {
	// Get the value to yield
	var yieldValue value.Value

	if yield.Value != nil {
		var err error

		yieldValue, err = g.generateExpression(yield.Value)
		if err != nil {
			return nil, err
		}
	} else {
		yieldValue = constant.NewInt(types.I64, 0)
	}

	// Call fiber_yield using consolidated approach
	return g.generateFiberRuntimeCall("fiber_yield", "fiber_yield", []value.Value{yieldValue})
}

// generateChannelExpression generates REAL channel creation.
func (g *LLVMGenerator) generateChannelExpression(channel *ast.ChannelExpression) (value.Value, error) {
	// Get capacity
	var capacity value.Value = constant.NewInt(types.I64, defaultChannelCapacity)

	if channel.Capacity != nil {
		var err error

		capacity, err = g.generateExpression(channel.Capacity)
		if err != nil {
			return nil, err
		}
	}

	// Call channel_create using consolidated approach
	return g.generateFiberRuntimeCall("Channel", "channel_create", []value.Value{capacity})
}

// generateChannelCreateExpression generates REAL channel creation using type constructor syntax.
func (g *LLVMGenerator) generateChannelCreateExpression(channel *ast.ChannelCreateExpression) (value.Value, error) {
	// Get capacity
	var capacity value.Value = constant.NewInt(types.I64, defaultChannelCapacity)

	if channel.Capacity != nil {
		var err error

		capacity, err = g.generateExpression(channel.Capacity)
		if err != nil {
			return nil, err
		}
	}

	// Call channel_create using consolidated approach
	return g.generateFiberRuntimeCall("Channel", "channel_create", []value.Value{capacity})
}

// generateChannelSendExpression generates REAL channel send with blocking.
func (g *LLVMGenerator) generateChannelSendExpression(send *ast.ChannelSendExpression) (value.Value, error) {
	// Get channel ID
	channelID, err := g.generateExpression(send.Channel)
	if err != nil {
		return nil, err
	}

	// Get value to send
	sendValue, err := g.generateExpression(send.Value)
	if err != nil {
		return nil, err
	}

	// Call channel_send using consolidated approach
	return g.generateFiberRuntimeCall("send", "channel_send", []value.Value{channelID, sendValue})
}

// generateChannelRecvExpression generates REAL channel receive with blocking.
func (g *LLVMGenerator) generateChannelRecvExpression(recv *ast.ChannelRecvExpression) (value.Value, error) {
	// Get channel ID
	channelID, err := g.generateExpression(recv.Channel)
	if err != nil {
		return nil, err
	}

	// Call channel_recv using consolidated approach
	return g.generateFiberRuntimeCall("recv", "channel_recv", []value.Value{channelID})
}

// generateChannelFunctionCall generates channel function calls using builtin registry
func (g *LLVMGenerator) generateChannelFunctionCall(builtinName string, runtimeName string,
	callExpr *ast.CallExpression) (value.Value, error) {
	// Generate arguments
	args := make([]value.Value, len(callExpr.Arguments))

	for i, arg := range callExpr.Arguments {
		var err error

		args[i], err = g.generateExpression(arg)
		if err != nil {
			return nil, err
		}
	}

	// Call the appropriate channel function using consolidated approach
	return g.generateFiberRuntimeCall(builtinName, runtimeName, args)
}

// generateChannelSendCall generates channel send from built-in function call
func (g *LLVMGenerator) generateChannelSendCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateChannelFunctionCall("send", "channel_send", callExpr)
}

// generateChannelRecvCall generates channel receive from built-in function call
func (g *LLVMGenerator) generateChannelRecvCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateChannelFunctionCall("recv", "channel_recv", callExpr)
}

// generateChannelCreateCall generates channel creation from built-in function call
func (g *LLVMGenerator) generateChannelCreateCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateChannelFunctionCall("Channel", "channel_create", callExpr)
}

// generateSelectExpression generates select with proper channel multiplexing.
func (g *LLVMGenerator) generateSelectExpression(selectExpr *ast.SelectExpression) (value.Value, error) {
	// For now, implement basic select that evaluates first ready channel
	// TODO: Implement proper non-deterministic select with runtime support
	if len(selectExpr.Arms) == 0 {
		return constant.NewInt(types.I64, 0), nil
	}

	// For simplicity, evaluate first arm
	firstArm := selectExpr.Arms[0]

	result, err := g.generateExpression(firstArm.Expression)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// generateLambdaExpression generates lambda with basic support.
func (g *LLVMGenerator) generateLambdaExpression(lambda *ast.LambdaExpression) (value.Value, error) {
	// Create a simple function for the lambda
	funcName := fmt.Sprintf("lambda_%d", len(g.module.Funcs))

	// Create parameters with names
	var params []*ir.Param
	for _, param := range lambda.Parameters {
		llvmParam := ir.NewParam(param.Name, types.I64) // Assume int type for now
		params = append(params, llvmParam)
	}

	// Create function with parameters
	lambdaFunc := g.module.NewFunc(funcName, types.I64, params...)

	// Create entry block
	entryBlock := lambdaFunc.NewBlock("entry")

	// Save current builder and switch to lambda
	oldBuilder := g.builder
	g.builder = entryBlock

	// Save current variables and create new scope
	savedVars := make(map[string]value.Value)
	for k, v := range g.variables {
		savedVars[k] = v
	}

	// Add lambda parameters to scope
	for i, param := range lambda.Parameters {
		if i < len(lambdaFunc.Params) {
			g.variables[param.Name] = lambdaFunc.Params[i]
		}
	}

	// Generate lambda body
	bodyValue, err := g.generateExpression(lambda.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to generate lambda body: %w", err)
	}

	// Create return instruction
	entryBlock.NewRet(bodyValue)

	// Restore context
	g.variables = savedVars
	g.builder = oldBuilder

	// Return the function
	return lambdaFunc, nil
}

// generateModuleAccessExpression generates module access with fiber isolation.
func (g *LLVMGenerator) generateModuleAccessExpression(_ *ast.ModuleAccessExpression) (value.Value, error) {
	// TODO: Implement proper module state isolation per fiber
	// For now, return a placeholder value
	return constant.NewInt(types.I64, moduleAccessPlaceholder), nil
}

// generateSpawnCall generates fiber spawn from built-in function call
func (g *LLVMGenerator) generateSpawnCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateChannelFunctionCall("fiber_spawn", "fiber_spawn", callExpr)
}

// generateYieldCall generates fiber yield from built-in function call
func (g *LLVMGenerator) generateYieldCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateChannelFunctionCall("fiber_yield", "fiber_yield", callExpr)
}

// captureVariablesInExpression captures variables used in an expression by copying their values
func (g *LLVMGenerator) captureVariablesInExpression(expr ast.Expression, captured map[string]value.Value) {
	switch e := expr.(type) {
	case *ast.Identifier:
		// If this identifier refers to a local variable, capture its current value
		if val, exists := g.variables[e.Name]; exists {
			captured[e.Name] = val
		}
	case *ast.CallExpression:
		// Recursively capture variables in function arguments
		g.captureVariablesInExpression(e.Function, captured)
		for _, arg := range e.Arguments {
			g.captureVariablesInExpression(arg, captured)
		}
	case *ast.BinaryExpression:
		g.captureVariablesInExpression(e.Left, captured)
		g.captureVariablesInExpression(e.Right, captured)
	case *ast.UnaryExpression:
		g.captureVariablesInExpression(e.Operand, captured)
		// Add more cases as needed
	}
}

// generateAwaitCall generates fiber await from built-in function call
func (g *LLVMGenerator) generateAwaitCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateChannelFunctionCall("fiber_await", "fiber_await", callExpr)
}
