package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"

	"github.com/christianfindlay/osprey/internal/ast"
)

// GenerateProgram generates LLVM IR for a complete program.
func (g *LLVMGenerator) GenerateProgram(program *ast.Program) (*ir.Module, error) {
	// Skip plugin processing for now
	// processedProgram, err := g.processPlugins(program)
	// if err != nil {
	//     return nil, err
	// }

	// First pass: collect ALL function declarations and types (including main)
	mainFunc, topLevelStatements, err := g.collectDeclarations(program)
	if err != nil {
		return nil, err
	}

	// CRITICAL: Resolve all type variables before code generation
	g.typeInferer.ResolveAllEnvironmentTypes()

	// Create main function
	err = g.createMainFunction(mainFunc, topLevelStatements)
	if err != nil {
		return nil, err
	}

	// Second pass: generate code for user-defined functions (not main)
	err = g.generateUserFunctions(program)
	if err != nil {
		return nil, err
	}

	return g.module, nil
}

// collectDeclarations collects function declarations and top-level statements.
func (g *LLVMGenerator) collectDeclarations(program *ast.Program) (*ast.FunctionDeclaration, []ast.Statement, error) {
	var (
		mainFunc           *ast.FunctionDeclaration
		topLevelStatements []ast.Statement
	)

	// FIRST PASS: Declare ALL types and externs first

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.ExternDeclaration:
			// Process extern declarations in the first pass
			err := g.generateExternDeclaration(s)
			if err != nil {
				return nil, nil, err
			}
		case *ast.TypeDeclaration:
			g.declareType(s)
		case *ast.EffectDeclaration:
			// Process effect declarations in the first pass!
			// This populates the effect registry before any handlers are generated
			err := g.generateEffectDeclaration(s)
			if err != nil {
				return nil, nil, err
			}
		}
	}

	// SECOND PASS: Declare function signatures after types are available
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.FunctionDeclaration:
			if s.Name == MainFunctionName {
				mainFunc = s
				// Also declare main function signature so other code can reference it
				err := g.declareFunctionSignature(s)
				if err != nil {
					return nil, nil, err
				}
			} else {
				err := g.declareFunctionSignature(s)
				if err != nil {
					return nil, nil, err
				}
			}
		default:
			// Only add non-type, non-function, non-extern, non-effect statements
			if g.isTopLevelStatement(s) {
				topLevelStatements = append(topLevelStatements, stmt)
			}
		}
	}

	return mainFunc, topLevelStatements, nil
}

// isTopLevelStatement checks if a statement should be added to top-level statements
func (g *LLVMGenerator) isTopLevelStatement(stmt ast.Statement) bool {
	if _, isType := stmt.(*ast.TypeDeclaration); isType {
		return false
	}

	if _, isExtern := stmt.(*ast.ExternDeclaration); isExtern {
		return false
	}

	if _, isEffect := stmt.(*ast.EffectDeclaration); isEffect {
		return false
	}

	if _, isFunc := stmt.(*ast.FunctionDeclaration); isFunc {
		return false
	}

	return true
}

// createMainFunction creates the main function based on user definition or top-level statements.
func (g *LLVMGenerator) createMainFunction(
	mainFunc *ast.FunctionDeclaration,
	topLevelStatements []ast.Statement,
) error {
	// If there's a user-defined main function, generate it
	if mainFunc != nil {
		return g.generateFunctionDeclaration(mainFunc)
	}

	// Create main function for top-level statements
	main := g.module.NewFunc(MainFunctionName, types.I32)
	g.function = main
	g.builder = main.NewBlock("")

	// Process top-level statements in the main function
	for _, stmt := range topLevelStatements {
		err := g.generateStatement(stmt)
		if err != nil {
			return err
		}
	}

	// Return 0 from main
	g.builder.NewRet(constant.NewInt(types.I32, 0))

	return nil
}

// NOTE: Plugin system temporarily removed to fix tests
// TODO: Re-implement plugin system properly

// generateUserFunctions generates code for user-defined functions (excluding main).
func (g *LLVMGenerator) generateUserFunctions(program *ast.Program) error {
	for _, stmt := range program.Statements {
		if fnDecl, ok := stmt.(*ast.FunctionDeclaration); ok && fnDecl.Name != MainFunctionName {
			err := g.generateFunctionDeclaration(fnDecl)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
