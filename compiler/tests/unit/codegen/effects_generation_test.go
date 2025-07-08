package codegen

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
)

func TestGenerateEffectDeclaration_NoMalformedLLVM(t *testing.T) {
	// Test that generateEffectDeclaration doesn't generate malformed LLVM IR
	generator := codegen.NewLLVMGenerator()

	// Create a simple effect declaration like: effect Logger { log: fn(string) -> Unit }
	effectDecl := &ast.EffectDeclaration{
		Name: "Logger",
		Operations: []ast.EffectOperation{
			{
				Name: "log",
				Parameters: []ast.Parameter{
					{
						Name: "message",
						Type: &ast.TypeExpression{Name: codegen.TypeString},
					},
				},
				ReturnType: "Unit",
			},
		},
	}

	// Generate the effect declaration
	err := generator.RegisterEffectDeclaration(effectDecl)
	if err != nil {
		t.Fatalf("Failed to register effect declaration: %v", err)
	}

	// Get the generated LLVM IR
	llvmIR := generator.GenerateIR()

	// CRITICAL: Verify no malformed function declarations are generated
	// The old bug would create: declare void (i64) @__effect_Logger_log()
	// which is malformed LLVM IR

	// Should NOT contain malformed stub declarations
	if strings.Contains(llvmIR, "declare void (i64) @__effect_Logger_log()") {
		t.Errorf("Generated malformed LLVM IR: found malformed stub declaration")
	}

	// Should NOT contain any __effect_ stub functions at all
	if strings.Contains(llvmIR, "__effect_") {
		t.Errorf("Generated stub effect functions - effects should not generate stubs")
	}

	// Should NOT contain malformed parameter syntax like "(i64) @func()"
	if strings.Contains(llvmIR, ") @") && strings.Contains(llvmIR, "()") {
		t.Errorf("Generated malformed LLVM function syntax with parameters in wrong place")
	}

	// Verify basic LLVM IR structure is valid (has standard declarations)
	if !strings.Contains(llvmIR, "@printf") {
		t.Errorf("Generated LLVM IR missing standard function declarations")
	}
}

func TestGenerateEffectDeclaration_EffectRegistration(t *testing.T) {
	// Test that effect declarations are properly registered in the effect system
	generator := codegen.NewLLVMGenerator()

	// Create a complex effect declaration with multiple operations
	effect := &ast.EffectDeclaration{
		Name: "FileSystem",
		Operations: []ast.EffectOperation{
			{
				Name: "read",
				Parameters: []ast.Parameter{
					{
						Name: "path",
						Type: &ast.TypeExpression{Name: codegen.TypeString},
					},
				},
				ReturnType: codegen.TypeString,
			},
			{
				Name: "write",
				Parameters: []ast.Parameter{
					{
						Name: "path",
						Type: &ast.TypeExpression{Name: codegen.TypeString},
					},
					{
						Name: "content",
						Type: &ast.TypeExpression{Name: codegen.TypeString},
					},
				},
				ReturnType: "Unit",
			},
		},
	}

	// Register the effect
	err := generator.RegisterEffectDeclaration(effect)
	if err != nil {
		t.Fatalf("Failed to register effect declaration: %v", err)
	}

	// Verify the effect registry was populated
	// This is a bit tricky to test directly since the registry is internal
	// But we can verify by trying to register the same effect again
	err = generator.RegisterEffectDeclaration(effect)
	if err != nil {
		t.Errorf("Failed to register effect declaration twice: %v", err)
	}
}

func TestGenerateEffectDeclaration_ErrorHandling(t *testing.T) {
	// Test error handling for malformed effect declarations
	generator := codegen.NewLLVMGenerator()

	// Create an effect declaration with missing parameter type
	effect := &ast.EffectDeclaration{
		Name: "BadEffect",
		Operations: []ast.EffectOperation{
			{
				Name: "badOp",
				Parameters: []ast.Parameter{
					{
						Name: "param",
						Type: nil, // Missing type - should cause error
					},
				},
				ReturnType: "Unit",
			},
		},
	}

	// Should fail with proper error message
	err := generator.RegisterEffectDeclaration(effect)
	if err == nil {
		t.Errorf("Expected error for missing parameter type, got none")
	}

	if !strings.Contains(err.Error(), "INTERNAL COMPILER ERROR") {
		t.Errorf("Expected internal compiler error message, got: %v", err)
	}
}

func TestGenerateEffectDeclaration_MultipleEffects(t *testing.T) {
	// Test that multiple effect declarations work correctly
	generator := codegen.NewLLVMGenerator()

	// Create multiple effect declarations
	logger := &ast.EffectDeclaration{
		Name: "Logger",
		Operations: []ast.EffectOperation{
			{
				Name: "log",
				Parameters: []ast.Parameter{
					{
						Name: "message",
						Type: &ast.TypeExpression{Name: codegen.TypeString},
					},
				},
				ReturnType: "Unit",
			},
		},
	}

	state := &ast.EffectDeclaration{
		Name: "State",
		Operations: []ast.EffectOperation{
			{
				Name:       "get",
				Parameters: []ast.Parameter{},
				ReturnType: codegen.TypeInt,
			},
			{
				Name: "set",
				Parameters: []ast.Parameter{
					{
						Name: "value",
						Type: &ast.TypeExpression{Name: codegen.TypeInt},
					},
				},
				ReturnType: "Unit",
			},
		},
	}

	// Register both effects
	err := generator.RegisterEffectDeclaration(logger)
	if err != nil {
		t.Fatalf("Failed to register Logger effect: %v", err)
	}

	err = generator.RegisterEffectDeclaration(state)
	if err != nil {
		t.Fatalf("Failed to register State effect: %v", err)
	}

	// Verify no malformed LLVM IR generated
	llvmIR := generator.GenerateIR()

	if strings.Contains(llvmIR, "__effect_") {
		t.Errorf("Generated stub effect functions for multiple effects")
	}

	// Verify basic LLVM structure is maintained
	if !strings.Contains(llvmIR, "@printf") {
		t.Errorf("Generated LLVM IR missing standard function declarations with multiple effects")
	}
}

func TestGenerateEffectDeclaration_NoStubGeneration(t *testing.T) {
	// CRITICAL: Test that no stub functions are generated
	// This prevents the regression of malformed LLVM IR like:
	// declare void (i64) @__effect_Logger_log()
	generator := codegen.NewLLVMGenerator()

	effect := &ast.EffectDeclaration{
		Name: "TestEffect",
		Operations: []ast.EffectOperation{
			{
				Name: "operation",
				Parameters: []ast.Parameter{
					{
						Name: "param1",
						Type: &ast.TypeExpression{Name: codegen.TypeString},
					},
					{
						Name: "param2",
						Type: &ast.TypeExpression{Name: codegen.TypeInt},
					},
				},
				ReturnType: codegen.TypeBool,
			},
		},
	}

	// Get initial LLVM IR
	initialIR := generator.GenerateIR()

	// Register effect
	err := generator.RegisterEffectDeclaration(effect)
	if err != nil {
		t.Fatalf("Failed to register effect declaration: %v", err)
	}

	// Get final LLVM IR
	finalIR := generator.GenerateIR()

	// CRITICAL: The LLVM IR should be identical before and after effect registration
	// because effects should NOT generate any LLVM IR - they're handled by the real
	// algebraic effects system

	if finalIR != initialIR {
		t.Errorf("Effect declaration generated LLVM IR when it should not have")
		t.Errorf("Initial IR length: %d, Final IR length: %d", len(initialIR), len(finalIR))

		// Find the difference
		if len(finalIR) > len(initialIR) {
			diff := finalIR[len(initialIR):]
			t.Errorf("Extra LLVM IR generated: %s", diff)
		}
	}
}
