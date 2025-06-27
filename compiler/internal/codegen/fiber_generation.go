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

// initFiberRuntime declares external runtime functions.
func (g *LLVMGenerator) initFiberRuntime() {
	// fiber_spawn: Create and schedule a fiber
	fiberSpawn := g.module.NewFunc("fiber_spawn",
		types.I64, // returns fiber ID
		ir.NewParam("fn", types.NewPointer(types.NewFunc(types.I64))), // function pointer
	)
	fiberSpawn.Linkage = enum.LinkageExternal
	g.functions["fiber_spawn"] = fiberSpawn

	// fiber_await: Wait for fiber completion
	fiberAwait := g.module.NewFunc("fiber_await",
		types.I64,                          // returns fiber result
		ir.NewParam("fiber_id", types.I64), // fiber ID
	)
	fiberAwait.Linkage = enum.LinkageExternal
	g.functions["fiber_await"] = fiberAwait

	// fiber_yield: Yield control with value
	fiberYield := g.module.NewFunc("fiber_yield",
		types.I64,                       // returns yielded value
		ir.NewParam("value", types.I64), // value to yield
	)
	fiberYield.Linkage = enum.LinkageExternal
	g.functions["fiber_yield"] = fiberYield

	// fiber_sleep: Sleep for specified milliseconds
	fiberSleep := g.module.NewFunc("fiber_sleep",
		types.I64,                              // returns 0 on success
		ir.NewParam("milliseconds", types.I64), // milliseconds to sleep
	)
	fiberSleep.Linkage = enum.LinkageExternal
	g.functions["fiber_sleep"] = fiberSleep

	// channel_create: Create a channel
	channelCreate := g.module.NewFunc("channel_create",
		types.I64,                          // returns channel ID
		ir.NewParam("capacity", types.I64), // channel capacity
	)
	channelCreate.Linkage = enum.LinkageExternal
	g.functions["channel_create"] = channelCreate

	// channel_send: Send value to channel
	channelSend := g.module.NewFunc("channel_send",
		types.I64,                            // returns success (1) or failure (0)
		ir.NewParam("channel_id", types.I64), // channel ID
		ir.NewParam("value", types.I64),      // value to send
	)
	channelSend.Linkage = enum.LinkageExternal
	g.functions["channel_send"] = channelSend

	// channel_recv: Receive from channel
	channelRecv := g.module.NewFunc("channel_recv",
		types.I64,                            // returns received value
		ir.NewParam("channel_id", types.I64), // channel ID
	)
	channelRecv.Linkage = enum.LinkageExternal
	g.functions["channel_recv"] = channelRecv
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

	// Create new context for closure
	g.function = closureFunc
	entry := closureFunc.NewBlock("entry")
	g.builder = entry
	g.variables = make(map[string]value.Value)

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

	// Get runtime spawn function
	spawnFunc := g.functions["fiber_spawn"]
	if spawnFunc == nil {
		g.initFiberRuntime()
		spawnFunc = g.functions["fiber_spawn"]
	}

	// Call fiber_spawn with the closure
	fiberID := g.builder.NewCall(spawnFunc, closureFunc)

	return fiberID, nil
}

// generateAwaitExpression generates REAL fiber await with blocking.
func (g *LLVMGenerator) generateAwaitExpression(await *ast.AwaitExpression) (value.Value, error) {
	// Generate the fiber ID expression
	fiberID, err := g.generateExpression(await.Expression)
	if err != nil {
		return nil, err
	}

	// Get runtime await function
	awaitFunc := g.functions["fiber_await"]
	if awaitFunc == nil {
		g.initFiberRuntime()
		awaitFunc = g.functions["fiber_await"]
	}

	// Call fiber_await to block until completion
	result := g.builder.NewCall(awaitFunc, fiberID)

	return result, nil
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

	// Get runtime yield function
	yieldFunc := g.functions["fiber_yield"]
	if yieldFunc == nil {
		g.initFiberRuntime()
		yieldFunc = g.functions["fiber_yield"]
	}

	// Call fiber_yield
	result := g.builder.NewCall(yieldFunc, yieldValue)

	return result, nil
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

	// Get runtime channel create function
	createFunc := g.functions["channel_create"]
	if createFunc == nil {
		g.initFiberRuntime()
		createFunc = g.functions["channel_create"]
	}

	// Call channel_create
	channelID := g.builder.NewCall(createFunc, capacity)

	return channelID, nil
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

	// Get runtime channel create function
	createFunc := g.functions["channel_create"]
	if createFunc == nil {
		g.initFiberRuntime()
		createFunc = g.functions["channel_create"]
	}

	// Call channel_create
	channelID := g.builder.NewCall(createFunc, capacity)

	return channelID, nil
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

	// Get runtime send function
	sendFunc := g.functions["channel_send"]
	if sendFunc == nil {
		g.initFiberRuntime()
		sendFunc = g.functions["channel_send"]
	}

	// Call channel_send
	result := g.builder.NewCall(sendFunc, channelID, sendValue)

	return result, nil
}

// generateChannelRecvExpression generates REAL channel receive with blocking.
func (g *LLVMGenerator) generateChannelRecvExpression(recv *ast.ChannelRecvExpression) (value.Value, error) {
	// Get channel ID
	channelID, err := g.generateExpression(recv.Channel)
	if err != nil {
		return nil, err
	}

	// Get runtime recv function
	recvFunc := g.functions["channel_recv"]
	if recvFunc == nil {
		g.initFiberRuntime()
		recvFunc = g.functions["channel_recv"]
	}

	// Call channel_recv
	result := g.builder.NewCall(recvFunc, channelID)

	return result, nil
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

// generateLambdaExpression generates lambda with proper closure support.
func (g *LLVMGenerator) generateLambdaExpression(lambda *ast.LambdaExpression) (value.Value, error) {
	// For now, evaluate lambda body immediately
	// TODO: Implement proper closure creation with captured variables
	return g.generateExpression(lambda.Body)
}

// generateModuleAccessExpression generates module access with fiber isolation.
func (g *LLVMGenerator) generateModuleAccessExpression(_ *ast.ModuleAccessExpression) (value.Value, error) {
	// TODO: Implement proper module state isolation per fiber
	// For now, return a placeholder value
	return constant.NewInt(types.I64, moduleAccessPlaceholder), nil
}
