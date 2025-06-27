package codegen_test

import (
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// TestEffectDeclarationCompilation tests basic effect declaration compilation
func TestEffectDeclarationCompilation(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		checkIR     string // What to look for in the generated IR
	}{
		{
			name: "simple Logger effect",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				fn main() -> Unit = print("test")
			`,
			expectError: false,
			checkIR:     "main",
		},
		{
			name: "multiple operation effect",
			source: `
				effect Logger {
					log: fn(string) -> Unit
					error: fn(string) -> Unit
					debug: fn(string) -> Unit
				}
				fn main() -> Unit = print("test")
			`,
			expectError: false,
			checkIR:     "main",
		},
		{
			name: "State effect with different types",
			source: `
				effect State {
					get: fn() -> int
					set: fn(int) -> Unit
				}
				fn main() -> Unit = print("test")
			`,
			expectError: false,
			checkIR:     "main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llvmIR, err := codegen.CompileToLLVM(tt.source)

			if tt.expectError {
				if err == nil {
					t.Error("Expected compilation error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected compilation error: %v", err)
				return
			}

			if tt.checkIR != "" && !strings.Contains(llvmIR, tt.checkIR) {
				t.Errorf("Expected IR to contain '%s'", tt.checkIR)
			}
		})
	}
}

// TestPerformExpressionCompilation tests perform expression compilation
func TestPerformExpressionCompilation(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		errorMsg    string
	}{
		{
			name: "unhandled perform should fail",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				fn test() -> Unit !Logger = perform Logger.log("test")
				fn main() -> Unit = test()
			`,
			expectError: true,
			errorMsg:    "no handler found for effect",
		},
		{
			name: "multiple unhandled effects should fail",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				effect State {
					get: fn() -> int
				}
				fn test() -> int ![Logger, State] = {
					perform Logger.log("test")
					perform State.get()
				}
				fn main() -> Unit = print(toString(test()))
			`,
			expectError: true,
			errorMsg:    "no handler found for effect",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(tt.source)

			if tt.expectError {
				if err == nil {
					t.Error("Expected compilation error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected compilation error: %v", err)
				}
			}
		})
	}
}

// TestHandlerExpressionCompilation tests handler expression compilation
func TestHandlerExpressionCompilation(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		checkIR     []string // What to look for in the generated IR
	}{
		{
			name: "simple handler should compile",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				fn test() -> Unit !Logger = perform Logger.log("test")
				fn main() -> Unit = with handler Logger
					log(msg) => print(msg)
				{
					test()
				}
			`,
			expectError: false,
			checkIR:     []string{"__handler_Logger_log"},
		},
		{
			name: "multiple operation handler should compile",
			source: `
				effect Logger {
					log: fn(string) -> Unit
					error: fn(string) -> Unit
				}
				fn test() -> Unit !Logger = {
					perform Logger.log("info")
					perform Logger.error("error")
				}
				fn main() -> Unit = with handler Logger
					log(msg) => print("[LOG] " + msg)
					error(msg) => print("[ERROR] " + msg)
				{
					test()
				}
			`,
			expectError: false,
			checkIR:     []string{"__handler_Logger_log", "__handler_Logger_error"},
		},
		{
			name: "nested handlers should compile",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				effect State {
					get: fn() -> int
				}
				fn test() -> int ![Logger, State] = {
					perform Logger.log("getting state")
					perform State.get()
				}
				fn main() -> Unit = with handler Logger
					log(msg) => print(msg)
				{
					with handler State
						get() => 42
					{
						let result = test()
						print(toString(result))
					}
				}
			`,
			expectError: false,
			checkIR:     []string{"__handler_Logger_log", "__handler_State_get"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llvmIR, err := codegen.CompileToLLVM(tt.source)

			if tt.expectError {
				if err == nil {
					t.Error("Expected compilation error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected compilation error: %v", err)
				return
			}

			// Check that expected handler functions are generated
			for _, checkStr := range tt.checkIR {
				if !strings.Contains(llvmIR, checkStr) {
					t.Errorf("Expected IR to contain '%s'", checkStr)
				}
			}
		})
	}
}

// TestFunctionSignatureWithEffects tests function signature generation with effects
func TestFunctionSignatureWithEffects(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
		checkIR     []string
	}{
		{
			name: "function with single effect",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				fn test() -> Unit !Logger = perform Logger.log("test")
				fn main() -> Unit = with handler Logger
					log(msg) => print(msg)
				{
					test()
				}
			`,
			expectError: false,
			checkIR:     []string{"define", "test", "__evidence_Logger_log"},
		},
		{
			name: "function with multiple effects",
			source: `
				effect Logger {
					log: fn(string) -> Unit
				}
				effect State {
					get: fn() -> int
				}
				fn test() -> int ![Logger, State] = {
					perform Logger.log("test")
					perform State.get()
				}
				fn main() -> Unit = with handler Logger
					log(msg) => print(msg)
				{
					with handler State
						get() => 42
					{
						let result = test()
						print(toString(result))
					}
				}
			`,
			expectError: false,
			checkIR:     []string{"define", "test", "__evidence_Logger_log", "__evidence_State_get"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llvmIR, err := codegen.CompileToLLVM(tt.source)

			if tt.expectError {
				if err == nil {
					t.Error("Expected compilation error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected compilation error: %v", err)
				return
			}

			// Check that function signatures include evidence parameters
			for _, checkStr := range tt.checkIR {
				if !strings.Contains(llvmIR, checkStr) {
					t.Errorf("Expected IR to contain '%s'", checkStr)
				}
			}
		})
	}
}

// TestHandlerIsolationIR tests that different handlers generate separate functions
func TestHandlerIsolationIR(t *testing.T) {
	source := `
		effect Logger {
			log: fn(string) -> Unit
		}
		fn task(value: int) -> int !Logger = {
			perform Logger.log("processing: " + toString(value))
			value * value
		}
		fn main() -> Unit = {
			// Production handler
			let result1 = with handler Logger
				log(msg) => print("[PROD] " + msg)
			{
				task(value: 5)
			}
			
			// Silent handler
			let result2 = with handler Logger
				log(msg) => 0
			{
				task(value: 7)
			}
		}
	`

	llvmIR, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	// Should generate different handler functions
	expectedHandlers := []string{
		"__handler_Logger_log_0", // First handler
		"__handler_Logger_log_1", // Second handler
	}

	for _, handler := range expectedHandlers {
		if !strings.Contains(llvmIR, handler) {
			t.Errorf("Expected to find handler function '%s' in IR", handler)
		}
	}

	// This test documents the current behavior - helps isolate the isolation bug
	t.Logf("Handler isolation test generated IR with %d handler functions", len(expectedHandlers))
}

// TestEvidencePassingIR tests that evidence parameters are correctly passed
func TestEvidencePassingIR(t *testing.T) {
	source := `
		effect Logger {
			log: fn(string) -> Unit
		}
		fn inner() -> Unit !Logger = perform Logger.log("inner")
		fn outer() -> Unit !Logger = {
			perform Logger.log("outer start")
			inner()
			perform Logger.log("outer end")
		}
		fn main() -> Unit = with handler Logger
			log(msg) => print(msg)
		{
			outer()
		}
	`

	llvmIR, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	// Check that functions with effects receive evidence parameters
	expectedElements := []string{
		"define", "inner", // inner function definition
		"define", "outer", // outer function definition
		"__evidence_Logger_log", // evidence parameter
	}

	for _, element := range expectedElements {
		if !strings.Contains(llvmIR, element) {
			t.Errorf("Expected to find '%s' in IR", element)
		}
	}

	// This test documents evidence passing - helps isolate evidence bugs
	t.Logf("Evidence passing test compiled successfully")
}

// TestOperationSpecificRoutingIR tests that different operations get routed correctly
func TestOperationSpecificRoutingIR(t *testing.T) {
	source := `
		effect Logger {
			log: fn(string) -> Unit
			error: fn(string) -> Unit
			debug: fn(string) -> Unit
		}
		fn test() -> Unit !Logger = {
			perform Logger.log("log message")
			perform Logger.error("error message")
			perform Logger.debug("debug message")
		}
		fn main() -> Unit = with handler Logger
			log(msg) => print("[LOG] " + msg)
			error(msg) => print("[ERROR] " + msg)
			debug(msg) => print("[DEBUG] " + msg)
		{
			test()
		}
	`

	llvmIR, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	// Should generate operation-specific handler functions
	expectedHandlers := []string{
		"__handler_Logger_log_",
		"__handler_Logger_error_",
		"__handler_Logger_debug_",
	}

	for _, handlerPrefix := range expectedHandlers {
		if !strings.Contains(llvmIR, handlerPrefix) {
			t.Errorf("Expected to find handler function with prefix '%s' in IR", handlerPrefix)
		}
	}

	// This test documents operation-specific routing - helps isolate routing bugs
	t.Logf("Operation-specific routing test compiled successfully")
}
