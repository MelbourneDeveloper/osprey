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

	funcIdent, ok := funcArg.(*ast.Identifier)
	if !ok {
		return nil, ErrForEachNotFunction
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
	LoopCond *ir.Block
	LoopBody *ir.Block
	LoopEnd  *ir.Block
}

// createForEachLoopBlocks creates the basic blocks needed for a forEach loop.
func (g *LLVMGenerator) createForEachLoopBlocks(callExpr *ast.CallExpression) *ForEachLoopBlocks {
	blockSuffix := fmt.Sprintf("_%p", callExpr)

	return &ForEachLoopBlocks{
		LoopCond: g.function.NewBlock("loop_cond" + blockSuffix),
		LoopBody: g.function.NewBlock("loop_body" + blockSuffix),
		LoopEnd:  g.function.NewBlock("loop_end" + blockSuffix),
	}
}

// generateForEachLoop generates the actual loop logic for forEach.
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

	_, err := g.callFunctionWithValue(funcIdent, counterValue)
	if err != nil {
		return err
	}

	one := constant.NewInt(types.I64, 1)
	incrementedValue := g.builder.NewAdd(counterValue, one)
	g.builder.NewStore(incrementedValue, counterPtr)

	g.builder.NewBr(blocks.LoopCond)

	g.builder = blocks.LoopEnd

	return nil
}

// generateMapCall handles map function calls.
func (g *LLVMGenerator) generateMapCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapBuiltInFunctionWrongArgs(MapFunc, len(callExpr.Arguments))
	}

	rangeValue, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	if _, ok := callExpr.Arguments[1].(*ast.Identifier); !ok {
		return nil, ErrMapNotFunction
	}

	// TODO: Implement proper lazy map
	return rangeValue, nil
}

// generateFilterCall handles filter function calls.
func (g *LLVMGenerator) generateFilterCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapBuiltInFunctionWrongArgs(FilterFunc, len(callExpr.Arguments))
	}

	iterator, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	funcArg := callExpr.Arguments[1]
	if funcIdent, ok := funcArg.(*ast.Identifier); ok {
		return g.callFunctionWithValue(funcIdent, iterator)
	}

	return nil, ErrFilterNotFunction
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

	funcIdent, ok := funcArg.(*ast.Identifier)
	if !ok {
		return nil, ErrFoldNotFunction
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

	// Call the fold function with (accumulator, currentValue)
	newAccumulator, err := g.callFunctionWithTwoValues(funcIdent, currentAccumulator, counterValue)
	if err != nil {
		return nil, err
	}

	// Store the new accumulator value
	g.builder.NewStore(newAccumulator, accumulatorPtr)

	// Increment counter
	one := constant.NewInt(types.I64, 1)
	incrementedValue := g.builder.NewAdd(counterValue, one)
	g.builder.NewStore(incrementedValue, counterPtr)

	// Jump back to condition
	g.builder.NewBr(blocks.LoopCond)

	// Loop end block
	g.builder = blocks.LoopEnd
	finalResult := g.builder.NewLoad(types.I64, accumulatorPtr)

	return finalResult, nil
}
