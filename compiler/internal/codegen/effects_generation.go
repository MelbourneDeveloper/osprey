package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// generatePerformExpression generates LLVM IR for perform expressions
func (g *LLVMGenerator) generatePerformExpression(perform *ast.PerformExpression) (value.Value, error) {
	// For now, implement as a direct function call to demonstrate functionality
	// In a full implementation, this would involve effect handlers and continuations

	// Generate function name from effect and operation
	functionName := fmt.Sprintf("__effect_%s_%s", perform.EffectName, perform.OperationName)

	// Generate arguments
	var args []value.Value
	for _, argExpr := range perform.Arguments {
		argVal, err := g.generateExpression(argExpr)
		if err != nil {
			return nil, err
		}
		args = append(args, argVal)
	}

	// Check if function exists in module
	var targetFunc *ir.Func
	for _, fn := range g.module.Funcs {
		if fn.Name() == functionName {
			targetFunc = fn
			break
		}
	}

	// If function doesn't exist, create a stub that prints the operation
	if targetFunc == nil {
		targetFunc = g.createEffectStub(functionName, len(args))
	}

	// Generate the call
	return g.builder.NewCall(targetFunc, args...), nil
}

// generateHandlerExpression generates LLVM IR for handler expressions
func (g *LLVMGenerator) generateHandlerExpression(handler *ast.HandlerExpression) (value.Value, error) {
	// For this basic implementation, we'll execute the body directly
	// A full implementation would set up effect handlers

	// Store current effect handlers (for nested handlers)
	// TODO: Implement proper handler stack

	// Generate the body expression
	return g.generateExpression(handler.Body)
}

// createEffectStub creates a stub function for effect operations
func (g *LLVMGenerator) createEffectStub(name string, argCount int) *ir.Func {
	// DEBUG: Print what we're creating
	fmt.Printf("DEBUG: Creating effect stub for %s with %d args\n", name, argCount)

	// Create parameter types (string for most effect operations)
	var paramTypes []types.Type
	for i := 0; i < argCount; i++ {
		paramTypes = append(paramTypes, types.I8Ptr) // String type for most operations
	}

	// Create function type (void return for now)
	funcType := types.NewFunc(types.Void, paramTypes...)

	// Create the function
	stubFunc := g.module.NewFunc(name, funcType)

	// Create entry block
	entry := stubFunc.NewBlock("entry")

	// For logging effects, call printf to demonstrate they're working
	printfFunc := g.functions["printf"]
	if printfFunc != nil {
		if argCount > 0 && len(stubFunc.Params) > 0 {
			// Create format string for the effect operation with argument
			formatStr := fmt.Sprintf("[EFFECT %s] %%s\\n", name)
			formatConst := constant.NewCharArrayFromString(formatStr + "\x00")
			formatGlobal := g.module.NewGlobalDef("", formatConst)
			formatPtr := constant.NewGetElementPtr(types.NewArray(uint64(len(formatStr)+1), types.I8), formatGlobal, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))

			// Call printf with the format and the first argument
			entry.NewCall(printfFunc, formatPtr, stubFunc.Params[0])
		} else {
			// Create format string for effect operation without arguments
			formatStr := fmt.Sprintf("[EFFECT %s] called\\n", name)
			formatConst := constant.NewCharArrayFromString(formatStr + "\x00")
			formatGlobal := g.module.NewGlobalDef("", formatConst)
			formatPtr := constant.NewGetElementPtr(types.NewArray(uint64(len(formatStr)+1), types.I8), formatGlobal, constant.NewInt(types.I64, 0), constant.NewInt(types.I64, 0))

			// Call printf with just the format string
			entry.NewCall(printfFunc, formatPtr)
		}
	}

	// Add return
	entry.NewRet(nil)

	return stubFunc
}

// generateEffectDeclaration generates LLVM IR for effect declarations
func (g *LLVMGenerator) generateEffectDeclaration(effect *ast.EffectDeclaration) error {
	// Register effect operations as function declarations
	for _, operation := range effect.Operations {
		functionName := fmt.Sprintf("__effect_%s_%s", effect.Name, operation.Name)

		// For now, create a simple function signature
		// TODO: Parse the actual operation type signature
		funcType := types.NewFunc(types.Void, types.I64) // Simple stub

		// Check if function already exists
		exists := false
		for _, fn := range g.module.Funcs {
			if fn.Name() == functionName {
				exists = true
				break
			}
		}

		if !exists {
			// Create the function declaration
			g.module.NewFunc(functionName, funcType)
		}
	}

	return nil
}
