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

// resolveCallbackIdent normalises an iterator callback argument into a
// named-function identifier the loop generators can call back into.
//
//   - *ast.Identifier  → used as-is (the existing path).
//   - *ast.LambdaExpression → emitted as a fresh top-level function and
//     synthesised as an Identifier pointing at the generated name, so
//     `forEach(xs, fn(x) => print(x))` works alongside the named-fn form.
//   - anything else → returns notFnErr (the original "callback is not a
//     function" rejection).
func (g *LLVMGenerator) resolveCallbackIdent(
	funcArg ast.Expression,
	notFnErr error,
) (*ast.Identifier, error) {
	switch v := funcArg.(type) {
	case *ast.Identifier:
		return v, nil
	case *ast.LambdaExpression:
		// Run the inferer first so the lambda body's identifiers (including
		// parameters) resolve once the codegen path below kicks in.
		_, err := g.typeInferer.InferType(v)
		if err != nil {
			return nil, err
		}
		lambdaValue, err := g.generateLambdaExpression(v)
		if err != nil {
			return nil, err
		}
		lambdaFunc, ok := lambdaValue.(*ir.Func)
		if !ok {
			return nil, notFnErr
		}
		// Register under the generated mangled name so the existing
		// callFunctionWithValue / monomorphisation paths can resolve it.
		if g.functions == nil {
			g.functions = make(map[string]*ir.Func)
		}
		g.functions[lambdaFunc.Name()] = lambdaFunc
		return &ast.Identifier{Name: lambdaFunc.Name()}, nil
	}
	return nil, notFnErr
}

// generateRangeCall handles range function calls - creates an iterator from start to end.
func (g *LLVMGenerator) generateRangeCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(RangeFunc, callExpr)
	if err != nil {
		return nil, err
	}

	start, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	end, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Create a struct to hold range data: {start, end}
	rangeStructType := types.NewStruct(types.I64, types.I64)
	rangeValue := g.builder.NewAlloca(rangeStructType)

	// Store start value at index 0
	startPtr := g.builder.NewGetElementPtr(rangeStructType, rangeValue,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	g.builder.NewStore(start, startPtr)

	// Store end value at index 1
	endPtr := g.builder.NewGetElementPtr(rangeStructType, rangeValue,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	g.builder.NewStore(end, endPtr)

	return rangeValue, nil
}

// generateForEachCall handles forEach function calls - applies a function to each element.
func (g *LLVMGenerator) generateForEachCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(ForEachFunc, callExpr)
	if err != nil {
		return nil, err
	}

	// Get the range struct from first argument
	rangeValue, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get the function to apply from second argument
	funcArg := callExpr.Arguments[1]

	funcIdent, err := g.resolveCallbackIdent(funcArg, ErrForEachNotFunction)
	if err != nil {
		return nil, err
	}

	// Extract range bounds
	start, end := g.extractRangeBounds(rangeValue)

	// Create loop blocks
	blocks := g.createForEachLoopBlocks(callExpr)

	// Generate loop logic
	err = g.generateForEachLoop(start, end, funcIdent, blocks)
	if err != nil {
		return nil, err
	}

	// Return the original range struct for potential pipe chaining
	return rangeValue, nil
}

// extractRangeBounds extracts start and end values from a range struct.
func (g *LLVMGenerator) extractRangeBounds(rangeValue value.Value) (value.Value, value.Value) {
	// Define the range struct type
	rangeStructType := types.NewStruct(types.I64, types.I64)

	var startPtr, endPtr value.Value

	if rangeValue.Type().String() == rangeStructType.String()+"*" {
		startPtr = g.builder.NewGetElementPtr(rangeStructType, rangeValue,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		endPtr = g.builder.NewGetElementPtr(rangeStructType, rangeValue,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	} else {
		tempRange := g.builder.NewAlloca(rangeStructType)
		g.builder.NewStore(rangeValue, tempRange)

		startPtr = g.builder.NewGetElementPtr(rangeStructType, tempRange,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
		endPtr = g.builder.NewGetElementPtr(rangeStructType, tempRange,
			constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	}

	start := g.builder.NewLoad(types.I64, startPtr)
	end := g.builder.NewLoad(types.I64, endPtr)

	return start, end
}

// ForEachLoopBlocks holds the basic blocks for a forEach loop.
type ForEachLoopBlocks struct {
	LoopCond      *ir.Block
	LoopBody      *ir.Block
	LoopIncrement *ir.Block
	LoopEnd       *ir.Block
}

// createForEachLoopBlocks creates the basic blocks needed for a forEach loop.
func (g *LLVMGenerator) createForEachLoopBlocks(callExpr *ast.CallExpression) *ForEachLoopBlocks {
	blockSuffix := fmt.Sprintf("_%p", callExpr)

	return &ForEachLoopBlocks{
		LoopCond:      g.function.NewBlock("loop_cond" + blockSuffix),
		LoopBody:      g.function.NewBlock("loop_body" + blockSuffix),
		LoopIncrement: g.function.NewBlock("loop_increment" + blockSuffix),
		LoopEnd:       g.function.NewBlock("loop_end" + blockSuffix),
	}
}

// generateForEachLoop generates the actual loop logic for forEach with stream fusion.
// Applies pending map/filter transformations inline for zero-cost abstractions.
func (g *LLVMGenerator) generateForEachLoop(
	start, end value.Value,
	funcIdent *ast.Identifier,
	blocks *ForEachLoopBlocks,
) error {
	counterPtr := g.builder.NewAlloca(types.I64)
	g.builder.NewStore(start, counterPtr)

	g.builder.NewBr(blocks.LoopCond)

	g.builder = blocks.LoopCond
	currentCounter := g.builder.NewLoad(types.I64, counterPtr)
	condition := g.builder.NewICmp(enum.IPredSLT, currentCounter, end)
	g.builder.NewCondBr(condition, blocks.LoopBody, blocks.LoopEnd)

	g.builder = blocks.LoopBody
	counterValue := g.builder.NewLoad(types.I64, counterPtr)

	// Replay pipeline ops (map / filter) in source order so that
	// `xs |> map(f) |> filter(p)` filters on f(x), not on x.
	processedValue, err := g.applyPendingIterOps(counterValue, blocks.LoopIncrement, fmt.Sprintf("_%p", blocks))
	if err != nil {
		return err
	}

	_, err = g.callFunctionWithValue(funcIdent, processedValue)
	if err != nil {
		return err
	}
	g.builder.NewBr(blocks.LoopIncrement)

	// Increment counter in the common increment block
	g.builder = blocks.LoopIncrement
	one := constant.NewInt(types.I64, 1)
	incrementedValue := g.builder.NewAdd(counterValue, one)
	g.builder.NewStore(incrementedValue, counterPtr)

	g.builder.NewBr(blocks.LoopCond)

	g.builder = blocks.LoopEnd

	// STREAM FUSION: Clear pending transformations after consuming
	g.pendingMapFunc = nil
	g.pendingFilterFunc = nil
	g.pendingIterOps = nil

	return nil
}

// applyPendingIterOps replays the recorded map/filter pipeline stages in
// source order on `v`. Filter stages branch to skipBlock when the
// predicate returns zero (so the consumer is bypassed for that
// element). Returns the final transformed value, with the builder
// positioned in the block that should hold the consumer call.
func (g *LLVMGenerator) applyPendingIterOps(
	v value.Value, skipBlock *ir.Block, blockSuffix string,
) (value.Value, error) {
	current := v
	for i, op := range g.pendingIterOps {
		switch op.kind {
		case iterOpMap:
			mapped, err := g.callFunctionWithValue(op.fn, current)
			if err != nil {
				return nil, err
			}
			current = mapped
		case iterOpFilter:
			predicate, err := g.callFunctionWithValue(op.fn, current)
			if err != nil {
				return nil, err
			}
			passBlock := g.function.NewBlock(fmt.Sprintf("filter_pass_%d%s", i, blockSuffix))
			// Compare to a zero of the predicate's own int type so an
			// i1 (bool) result doesn't crash IR-emit against an i64 zero.
			intT, ok := predicate.Type().(*types.IntType)
			if !ok {
				return nil, fmt.Errorf("%w: predicate must return int/bool, got %s",
					ErrFilterNotFunction, predicate.Type())
			}
			zero := constant.NewInt(intT, 0)
			isNonZero := g.builder.NewICmp(enum.IPredNE, predicate, zero)
			g.builder.NewCondBr(isNonZero, passBlock, skipBlock)
			g.builder = passBlock
		}
	}
	return current, nil
}

// generateMapCall handles map function calls using stream fusion.
// Stores the transformation function and returns the range unchanged.
// The transformation will be fused into forEach/fold when they consume the iterator.
func (g *LLVMGenerator) generateMapCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(MapFunc, callExpr)
	if err != nil {
		return nil, err
	}

	// Get the range struct from first argument (iterator)
	rangeValue, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get the transformation function
	funcArg := callExpr.Arguments[1]

	funcIdent, err := g.resolveCallbackIdent(funcArg, ErrMapNotFunction)
	if err != nil {
		return nil, err
	}

	// STREAM FUSION: Store the map function for later fusion with forEach/fold
	g.pendingMapFunc = funcIdent
	g.pendingIterOps = append(g.pendingIterOps, pendingIterOp{kind: iterOpMap, fn: funcIdent})

	return rangeValue, nil
}

// generateFilterCall handles filter function calls using stream fusion.
// Stores the predicate function and returns the range unchanged.
// The filter will be fused into forEach/fold when they consume the iterator.
func (g *LLVMGenerator) generateFilterCall(callExpr *ast.CallExpression) (value.Value, error) {
	err := validateBuiltInArgs(FilterFunc, callExpr)
	if err != nil {
		return nil, err
	}

	// Get the range struct from first argument (iterator)
	rangeValue, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get the predicate function
	funcArg := callExpr.Arguments[1]

	funcIdent, err := g.resolveCallbackIdent(funcArg, ErrFilterNotFunction)
	if err != nil {
		return nil, err
	}

	// STREAM FUSION: Store the filter predicate for later fusion with forEach/fold
	g.pendingFilterFunc = funcIdent
	g.pendingIterOps = append(g.pendingIterOps, pendingIterOp{kind: iterOpFilter, fn: funcIdent})

	return rangeValue, nil
}

// generateFoldCall handles fold function calls.
func (g *LLVMGenerator) generateFoldCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != ThreeArgs {
		return nil, WrapBuiltInFunctionWrongArgs(FoldFunc, len(callExpr.Arguments))
	}

	// Get the range struct from first argument (iterator)
	rangeValue, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get the initial accumulator value
	initial, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Get the fold function
	funcArg := callExpr.Arguments[2]

	funcIdent, err := g.resolveCallbackIdent(funcArg, ErrFoldNotFunction)
	if err != nil {
		return nil, err
	}

	// Extract range bounds
	start, end := g.extractRangeBounds(rangeValue)

	// Create fold loop blocks
	blocks := g.createFoldLoopBlocks(callExpr)

	// Generate fold loop logic
	return g.generateFoldLoop(start, end, initial, funcIdent, blocks)
}

// FoldLoopBlocks holds the basic blocks for a fold loop.
type FoldLoopBlocks struct {
	LoopCond *ir.Block
	LoopBody *ir.Block
	LoopEnd  *ir.Block
}

// createFoldLoopBlocks creates the basic blocks needed for a fold loop.
func (g *LLVMGenerator) createFoldLoopBlocks(callExpr *ast.CallExpression) *FoldLoopBlocks {
	blockSuffix := fmt.Sprintf("_%p", callExpr)

	return &FoldLoopBlocks{
		LoopCond: g.function.NewBlock("fold_cond" + blockSuffix),
		LoopBody: g.function.NewBlock("fold_body" + blockSuffix),
		LoopEnd:  g.function.NewBlock("fold_end" + blockSuffix),
	}
}

// generateFoldLoop generates the actual loop logic for fold.
func (g *LLVMGenerator) generateFoldLoop(
	start, end, initial value.Value,
	funcIdent *ast.Identifier,
	blocks *FoldLoopBlocks,
) (value.Value, error) {
	// Create counter and accumulator variables
	counterPtr := g.builder.NewAlloca(types.I64)
	accumulatorPtr := g.builder.NewAlloca(types.I64)

	// Initialize counter and accumulator
	g.builder.NewStore(start, counterPtr)
	g.builder.NewStore(initial, accumulatorPtr)

	// Jump to condition check
	g.builder.NewBr(blocks.LoopCond)

	// Loop condition block
	g.builder = blocks.LoopCond
	currentCounter := g.builder.NewLoad(types.I64, counterPtr)
	condition := g.builder.NewICmp(enum.IPredSLT, currentCounter, end)
	g.builder.NewCondBr(condition, blocks.LoopBody, blocks.LoopEnd)

	// Loop body block
	g.builder = blocks.LoopBody
	counterValue := g.builder.NewLoad(types.I64, counterPtr)
	currentAccumulator := g.builder.NewLoad(types.I64, accumulatorPtr)

	// Build a skip block that bumps the counter past a filter-rejected
	// element and loops back without touching the accumulator. Created
	// up-front so applyPendingIterOps can branch to it.
	blockSuffix := fmt.Sprintf("_%p", blocks)
	skipBlock := g.function.NewBlock("fold_filter_skip" + blockSuffix)
	{
		savedBuilder := g.builder
		g.builder = skipBlock
		one := constant.NewInt(types.I64, 1)
		next := g.builder.NewAdd(counterValue, one)
		g.builder.NewStore(next, counterPtr)
		g.builder.NewBr(blocks.LoopCond)
		g.builder = savedBuilder
	}

	// Replay pipeline ops (map / filter) in source order so that
	// `xs |> map(f) |> filter(p)` filters on f(x), not on x.
	processedValue, err := g.applyPendingIterOps(counterValue, skipBlock, blockSuffix)
	if err != nil {
		return nil, err
	}

	// Call the fold function with (accumulator, currentValue)
	newAccumulator, err := g.callFunctionWithTwoValues(funcIdent, currentAccumulator, processedValue)
	if err != nil {
		return nil, err
	}

	// Check if newAccumulator is a Result type and extract value if needed
	valueToStore := newAccumulator
	if structType, ok := newAccumulator.Type().(*types.StructType); ok {
		// Check if this looks like a Result struct: {i64, i8}
		const resultTypeFieldCount = 2
		if len(structType.Fields) == resultTypeFieldCount {
			if _, ok := structType.Fields[0].(*types.IntType); ok {
				const errorFlagBitSize = 8
				if intType, ok := structType.Fields[1].(*types.IntType); ok && intType.BitSize == errorFlagBitSize {
					// This is a Result type - extract the value field (index 0)
					valueToStore = g.builder.NewExtractValue(newAccumulator, 0)
				}
			}
		}
	}

	// Store the new accumulator value
	g.builder.NewStore(valueToStore, accumulatorPtr)

	// Increment counter
	one := constant.NewInt(types.I64, 1)
	incrementedValue := g.builder.NewAdd(counterValue, one)
	g.builder.NewStore(incrementedValue, counterPtr)

	// Jump back to condition
	g.builder.NewBr(blocks.LoopCond)

	// Loop end block
	g.builder = blocks.LoopEnd
	finalResult := g.builder.NewLoad(types.I64, accumulatorPtr)

	// STREAM FUSION: clear pending transformations now that fold consumed them.
	g.pendingMapFunc = nil
	g.pendingFilterFunc = nil
	g.pendingIterOps = nil

	return finalResult, nil
}
