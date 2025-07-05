package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// generatePerformExpression generates LLVM IR for perform expressions
func (g *LLVMGenerator) generatePerformExpression(perform *ast.PerformExpression) (value.Value, error) {
	// Use REAL algebraic effects system!
	return g.generateRealPerformExpression(perform)
}

// generateHandlerExpression generates LLVM IR for handler expressions
func (g *LLVMGenerator) generateHandlerExpression(handler *ast.HandlerExpression) (value.Value, error) {
	// Use REAL algebraic effects system with the EXISTING codegen instance!
	if g.effectCodegen == nil {
		g.InitializeEffects()
	}
	return g.effectCodegen.GenerateHandlerExpression(handler)
}

// generateEffectDeclaration generates LLVM IR for effect declarations
func (g *LLVMGenerator) generateEffectDeclaration(effect *ast.EffectDeclaration) error {
	// Register the effect with the effect system
	if err := g.RegisterEffectDeclaration(effect); err != nil {
		return err
	}

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
