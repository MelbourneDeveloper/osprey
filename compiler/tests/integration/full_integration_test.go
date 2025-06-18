package integration

// DO NOT EVER SKIP TESTS!!!!

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// TestMain runs before all tests in this package.
func TestMain(m *testing.M) {
	// Clean and rebuild everything before running any tests
	cleanAndRebuildAll()

	// Run all tests
	code := m.Run()

	// Exit with the test result code
	os.Exit(code)
}

// cleanAndRebuildAll cleans and rebuilds all dependencies.
func cleanAndRebuildAll() {
	// Get project root
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}
	projectRoot := filepath.Join(wd, "..", "..")

	// Clean everything including Rust
	cmd := exec.Command("make", "clean")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to clean: " + err.Error() + "\nOutput: " + string(output))
	}
	// Rebuild runtime libraries
	cmd = exec.Command("make", "fiber-runtime", "http-runtime")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build runtime libraries: " + err.Error() + "\nOutput: " + string(output))
	}

	// Build Rust interop library
	rustDir := filepath.Join(projectRoot, "examples", "rust_integration")
	if _, err := os.Stat(rustDir); err == nil {
		cmd = exec.Command("cargo", "build")
		cmd.Dir = rustDir
		if output, err := cmd.CombinedOutput(); err != nil {
			panic("Failed to build Rust interop: " + err.Error() + "\nOutput: " + string(output))
		}
	}

	// Build compiler (skip linting for tests)
	cmd = exec.Command("go", "build", "-o", "bin/osprey", "./cmd/osprey")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build compiler: " + err.Error() + "\nOutput: " + string(output))
	}
}

// ErrRustToolsNotFound indicates that Rust tools could not be located.
var ErrRustToolsNotFound = errors.New("rust tools not found in common locations")

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// checkLLVMTools verifies that required LLVM tools are available.
func checkLLVMTools(t *testing.T) {
	t.Helper()

	// Check for llc
	if _, err := exec.LookPath("llc"); err != nil {
		t.Fatalf("llc not found in PATH - required for integration tests. Install LLVM tools: brew install llvm")
	}

	// Check for clang
	if _, err := exec.LookPath("clang"); err != nil {
		t.Fatalf("clang not found in PATH - required for integration tests. Install clang: brew install llvm")
	}
}

// captureJITOutput captures stdout during JIT execution of source code.
func captureJITOutput(source string) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := codegen.CompileAndRunJIT(source)

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String(), err
}

// TestBasicCompilation tests that basic syntax compiles without errors.
func TestBasicCompilation(t *testing.T) {
	basicTests := map[string]string{
		"simple_let":             `let x = 42`,
		"simple_function":        `fn double(x) = x * 2`,
		"simple_print":           `print(42)`,
		"basic_match":            `let x = match 42 { 42 => 1 }`,
		"function_call":          `fn add(x, y) = x + y` + "\n" + `let result = add(x: 1, y: 2)`,
		"string_interpolation":   `let name = "Alice"` + "\n" + `print("Hello ${name}")`,
		"valid_type_declaration": `type Color = Red | Green | Blue`,
		"type_with_fields":       `type User = Admin { name: String, perms: Int } | Guest`,
	}

	for name, source := range basicTests {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Basic syntax %s failed to compile: %v", name, err)
			}
		})
	}
}

// TestErrorHandling tests that invalid syntax fails gracefully.
func TestErrorHandling(t *testing.T) {
	invalidTests := map[string]string{
		"undefined_variable": `print("Hello ${undefined_var}!")`,
		"undefined_function": `print("Result: ${unknownFunction()}")`,
		"missing_braces":     `fn test() = match x { 42 => 1`,
		"unclosed_string":    `let x = "hello`,
		"invalid_operator":   `let x = 1 @@ 2`,
		"missing_expression": `let x =`,
	}

	for name, source := range invalidTests {
		t.Run(name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err == nil {
				t.Errorf("Invalid syntax %s should have failed to compile", name)
			}
		})
	}
}

// TestFunctionArguments tests function argument requirements.
func TestFunctionArguments(t *testing.T) {
	// Valid cases
	validCases := map[string]string{
		"single_param": `fn double(x) = x * 2` + "\n" + `let result = double(5)`,
		"zero_param":   `fn getValue() = 42` + "\n" + `let result = getValue()`,
		"named_args":   `fn add(x: int, y: int) = x + y` + "\n" + `let result = add(x: 5, y: 10)`,
	}

	for name, source := range validCases {
		t.Run("valid_"+name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err != nil {
				t.Errorf("Valid case %s should compile: %v", name, err)
			}
		})
	}

	// Invalid cases (multi-param functions without named args)
	invalidCases := map[string]string{
		"two_param_positional":   `fn add(x: int, y: int) = x + y` + "\n" + `let result = add(5, 10)`,
		"three_param_positional": `fn combine(a: int, b: int, c: int) = a + b + c` + "\n" + `let result = combine(1, 2, 3)`,
	}

	for name, source := range invalidCases {
		t.Run("invalid_"+name, func(t *testing.T) {
			_, err := codegen.CompileToLLVM(source)
			if err == nil {
				t.Errorf("Invalid case %s should have failed", name)
			}
		})
	}
}

// findRustTools attempts to find Rust tools in common locations.
func findRustTools() (string, string, error) {
	// Common Rust installation paths
	commonPaths := []string{
		os.Getenv("HOME") + "/.cargo/bin",
		"/usr/local/bin",
		"/opt/homebrew/bin",
		"/usr/bin",
	}

	// Add current PATH
	currentPath := os.Getenv("PATH")

	var rustc, cargo string

	// Check each common path
	for _, path := range commonPaths {
		rustcPath := filepath.Join(path, "rustc")
		cargoPath := filepath.Join(path, "cargo")

		if _, err := os.Stat(rustcPath); err == nil {
			rustc = rustcPath
		}
		if _, err := os.Stat(cargoPath); err == nil {
			cargo = cargoPath
		}

		if rustc != "" && cargo != "" {
			// Update PATH to include this directory
			newPath := path + ":" + currentPath
			if err := os.Setenv("PATH", newPath); err != nil {
				return "", "", err
			}

			return rustc, cargo, nil
		}
	}

	// Try using exec.LookPath as fallback
	if rustcPath, err := exec.LookPath("rustc"); err == nil {
		if cargoPath, err := exec.LookPath("cargo"); err == nil {
			return rustcPath, cargoPath, nil
		}
	}

	return "", "", ErrRustToolsNotFound
}

// TestRustInterop tests the Rust-Osprey interop functionality.
func TestRustInterop(t *testing.T) {
	// Ensure compiler is built before running the test

	// Force the test to be visible in test explorers
	t.Log("🦀 Starting Rust interop test")

	// Find Rust tools in common locations
	rustc, cargo, err := findRustTools()
	if err != nil {
		t.Fatalf("❌ RUST TOOLS NOT FOUND - TEST FAILED. Install Rust: https://rustup.rs/ - Error: %v", err)
	}

	t.Logf("✅ Found Rust tools: rustc=%s, cargo=%s", rustc, cargo)

	// Check for clang
	if _, err := exec.LookPath("clang"); err != nil {
		t.Fatalf("❌ CLANG NOT FOUND - TEST FAILED. Install LLVM/Clang - Error: %v", err)
	}

	t.Log("✅ All required tools found")

	// Navigate to rust integration directory
	rustDir := "../../examples/rust_integration"

	// Check if directory exists
	if _, err := os.Stat(rustDir); os.IsNotExist(err) {
		t.Fatalf("❌ RUST INTEGRATION DIRECTORY NOT FOUND: %s", rustDir)
	}

	// Clean up any previous build artifacts first
	t.Log("🧹 Cleaning up previous Rust build artifacts...")
	cleanCmd := exec.Command("cargo", "clean")
	cleanCmd.Dir = rustDir
	if output, err := cleanCmd.CombinedOutput(); err != nil {
		t.Logf("⚠️ Warning: Failed to clean Rust artifacts: %v\nOutput: %s", err, output)
	}

	// Build the Rust library
	t.Log("🦀 Building Rust library...")
	buildCmd := exec.Command(cargo, "build", "--release")
	buildCmd.Dir = rustDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("❌ FAILED TO BUILD RUST LIBRARY: %v\nOutput: %s", err, output)
	}
	t.Log("✅ Rust library built successfully")

	// Verify the Rust library was created
	rustLibPath := filepath.Join(rustDir, "target/release/libosprey_math_utils.a")
	if !fileExists(rustLibPath) {
		t.Fatalf("❌ RUST LIBRARY NOT FOUND AT: %s", rustLibPath)
	}
	t.Log("✅ Rust library verified at:", rustLibPath)

	// Test the interop by running the demo script
	t.Log("🚀 Running Rust interop demo...")
	runCmd := exec.Command("./run.sh")
	runCmd.Dir = rustDir
	output, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("❌ FAILED TO RUN RUST INTEROP DEMO: %v\nOutput: %s", err, output)
	}

	// Verify the expected output contains Rust function results
	outputStr := string(output)
	expectedSubstrings := []string{
		"Rust add(15, 25) = 40",
		"Rust multiply(6, 7) = 42",
		"Rust factorial(5) = 120",
		"Rust fibonacci(10) = 55",
		"Rust is_prime(17) = 1",
		"✅ Rust-Osprey integration demo completed successfully!",
	}

	for _, expected := range expectedSubstrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("❌ EXPECTED OUTPUT MISSING: %q\nFull output:\n%s", expected, outputStr)
		}
	}

	t.Log("✅ Rust interop test completed successfully")
}

// TestRustInteropCompilationOnly tests that Rust interop code compiles without execution.
func TestRustInteropCompilationOnly(t *testing.T) {
	checkLLVMTools(t)

	rustInteropSource := `
extern fn rust_add(a: Int, b: Int) -> Int
extern fn rust_multiply(a: Int, b: Int) -> Int

let result1 = rust_add(a: 10, b: 20)
let result2 = rust_multiply(a: 5, b: 6)
printf("Sum: ", result1)
printf("Product: ", result2)
`

	// Test that the code compiles to LLVM IR without errors
	_, err := codegen.CompileToLLVM(rustInteropSource)
	if err != nil {
		t.Fatalf("Failed to compile Rust interop code: %v", err)
	}

	t.Logf("✅ Rust interop compilation test passed")
}

// TestRustInteropSimple is a simplified test that always runs in test explorers.
func TestRustInteropSimple(t *testing.T) {
	t.Log("🦀 Testing Rust interop compilation (simple)")

	// This test just verifies that Rust interop syntax compiles correctly
	rustInteropSource := `
extern fn rust_add(a: int, b: int) -> int
extern fn rust_multiply(a: int, b: int) -> int

let result1 = rust_add(a: 10, b: 20)
let result2 = rust_multiply(a: 5, b: 6)
print("Sum: ${result1}")
print("Product: ${result2}")
`

	// Test that the code compiles to LLVM IR without errors
	llvmIR, err := codegen.CompileToLLVM(rustInteropSource)
	if err != nil {
		t.Fatalf("Failed to compile Rust interop code: %v", err)
	}

	// Verify that external function declarations are in the LLVM IR
	expectedDeclarations := []string{
		"declare i64 @rust_add(i64 %a, i64 %b)",
		"declare i64 @rust_multiply(i64 %a, i64 %b)",
	}

	for _, expected := range expectedDeclarations {
		if !strings.Contains(llvmIR, expected) {
			t.Errorf("LLVM IR should contain declaration: %s", expected)
		}
	}

	// Verify function calls are generated
	expectedCalls := []string{
		"call i64 @rust_add(i64 10, i64 20)",
		"call i64 @rust_multiply(i64 5, i64 6)",
	}

	for _, expected := range expectedCalls {
		if !strings.Contains(llvmIR, expected) {
			t.Errorf("LLVM IR should contain function call: %s", expected)
		}
	}

	t.Log("✅ Rust interop compilation test passed")
}

// TestSystemLibraryInstallation tests that runtime libraries can be found from system locations
func TestSystemLibraryInstallation(t *testing.T) {
	checkLLVMTools(t)

	// Test that system-installed libraries work from any directory
	tempDir := t.TempDir()

	// Change to temp directory to simulate VSCode working from file directory
	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Test with a simple example that requires runtime linking
	source := `print("System libraries work!")`

	_, err = codegen.CompileToLLVM(source)
	if err != nil {
		t.Fatalf("Failed to compile from temp directory: %v", err)
	}

	// Test execution which requires runtime library linking
	output, err := captureJITOutput(source)
	if err != nil {
		t.Fatalf("Failed to execute from temp directory: %v", err)
	}

	expected := "System libraries work!\n"
	if output != expected {
		t.Errorf("Output mismatch: expected %q, got %q", expected, output)
	}

	t.Logf("✅ System library installation test passed from directory: %s", tempDir)
}
