package codegen

import (
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

	return nil
}
