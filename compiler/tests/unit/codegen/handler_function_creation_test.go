package codegen

import (
	"fmt"
	"testing"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// TestHandlerFunctionCreation isolates the LLVM function creation issue
func TestHandlerFunctionCreation(t *testing.T) {
	t.Run("CorrectLLVMFunctionCreation", func(t *testing.T) {
		// Create a minimal LLVM module to test function creation
		module := ir.NewModule()

		// CORRECT approach: Use ir.NewParam for parameters
		param1 := ir.NewParam("msg", types.I8Ptr)
		testFunc := module.NewFunc("test_handler", types.Void, param1)

		fmt.Printf("Created Function: %s\n", testFunc)
		fmt.Printf("Function Params: %d\n", len(testFunc.Params))

		// Emit LLVM IR
		irOutput := module.String()
		fmt.Printf("Generated LLVM IR:\n%s\n", irOutput)

		// CRITICAL TEST: Function should have 1 parameter
		if len(testFunc.Params) != 1 {
			t.Errorf("Expected 1 parameter, got %d", len(testFunc.Params))
		}

		// Test parameter name
		if testFunc.Params[0].Name() != "msg" {
			t.Errorf("Expected parameter name 'msg', got '%s'", testFunc.Params[0].Name())
		}
	})

	t.Run("MultipleParameterFunctionCreation", func(t *testing.T) {
		// Test with multiple parameters using correct approach
		module := ir.NewModule()

		param1 := ir.NewParam("msg", types.I8Ptr)
		param2 := ir.NewParam("count", types.I64)
		param3 := ir.NewParam("data", types.I8Ptr)

		multiFunc := module.NewFunc("multi_param_handler", types.Void, param1, param2, param3)

		fmt.Printf("Created Multi-param Function: %s\n", multiFunc)
		fmt.Printf("Multi-param Function Params: %d\n", len(multiFunc.Params))

		// Emit LLVM IR
		irOutput := module.String()
		fmt.Printf("Generated Multi-param LLVM IR:\n%s\n", irOutput)

		// CRITICAL TEST: Function should have 3 parameters
		if len(multiFunc.Params) != 3 {
			t.Errorf("Expected 3 parameters, got %d", len(multiFunc.Params))
		}

		// Test parameter names
		expectedNames := []string{"msg", "count", "data"}
		for i, expectedName := range expectedNames {
			if i < len(multiFunc.Params) {
				actualName := multiFunc.Params[i].Name()
				if actualName != expectedName {
					t.Errorf("Expected parameter %d name '%s', got '%s'", i, expectedName, actualName)
				}
			}
		}
	})
}
