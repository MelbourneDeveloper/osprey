package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"

	"github.com/christianfindlay/osprey/internal/ast"
)

// GenerateProgram generates LLVM IR for a complete program.
func (g *LLVMGenerator) GenerateProgram(program *ast.Program) (*ir.Module, error) {
	// First pass: collect ALL function declarations and types (including main)
	mainFunc, topLevelStatements, err := g.collectDeclarations(program)
	if err != nil {
		return nil, err
	}

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
	var mainFunc *ast.FunctionDeclaration
	var topLevelStatements []ast.Statement

	// FIRST: Declare ALL function signatures (including main if it's a function)
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.FunctionDeclaration:
			if s.Name == MainFunctionName {
				mainFunc = s
				// Also declare main function signature so other code can reference it
				if err := g.declareFunctionSignature(s); err != nil {
					return nil, nil, err
				}
			} else {
				if err := g.declareFunctionSignature(s); err != nil {
					return nil, nil, err
				}
			}
		case *ast.ExternDeclaration:
			// Process extern declarations in the first pass
			if err := g.generateExternDeclaration(s); err != nil {
				return nil, nil, err
			}
		case *ast.TypeDeclaration:
			g.declareType(s)
		case *ast.EffectDeclaration:
			// CRITICAL FIX: Process effect declarations in the first pass!
			// This populates the effect registry before any handlers are generated
			if err := g.generateEffectDeclaration(s); err != nil {
				return nil, nil, err
			}
		default:
			topLevelStatements = append(topLevelStatements, stmt)
		}
	}

	return mainFunc, topLevelStatements, nil
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
		if err := g.generateStatement(stmt); err != nil {
			return err
		}
	}

	// Return 0 from main
	g.builder.NewRet(constant.NewInt(types.I32, 0))

	return nil
}

// generateUserFunctions generates code for user-defined functions (excluding main).
func (g *LLVMGenerator) generateUserFunctions(program *ast.Program) error {
	for _, stmt := range program.Statements {
		if fnDecl, ok := stmt.(*ast.FunctionDeclaration); ok && fnDecl.Name != MainFunctionName {
			if err := g.generateFunctionDeclaration(fnDecl); err != nil {
				return err
			}
		}
	}

	return nil
}
