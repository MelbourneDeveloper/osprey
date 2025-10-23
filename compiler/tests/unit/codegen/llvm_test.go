package codegen_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/codegen"
)

// captureOutput captures stdout during execution of a function.
func captureOutput(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	_ = w.Close()

	os.Stdout = old

	var buf bytes.Buffer

	_, _ = io.Copy(&buf, r)

	return buf.String(), err
}

// testCompileAndRunWithOutput tests source code and captures its output.
func testCompileAndRunWithOutput(t *testing.T, source string, expectedOutput string, description string) {
	t.Helper()

	output, err := captureOutput(func() error {
		return codegen.CompileAndRunJIT(source)
	})
	if err != nil {
		// If JIT execution fails due to missing tools, FAIL LOUDLY
		if strings.Contains(err.Error(), "LLVM tools not found") ||
			strings.Contains(err.Error(), "no suitable compiler found") {
			t.Fatalf("CRITICAL MISSING DEPENDENCY for %s: %v", description, err)

			return
		}

		t.Fatalf("Failed to compile and run %s: %v", description, err)
	}

	output = strings.TrimSpace(output)
	expectedOutput = strings.TrimSpace(expectedOutput)

	if output != expectedOutput {
		t.Errorf("%s: expected output %q, got %q", description, expectedOutput, output)
	}
}

// testCompileOnly tests that source code compiles successfully without execution.
func testCompileOnly(t *testing.T, source string, description string) {
	t.Helper()

	_, err := codegen.CompileToLLVM(source)
	if err != nil {
		t.Fatalf("❌ COMPILATION FAILED for %s: %v", description, err)
	}
}

func TestBasicCompilation(t *testing.T) {
	testCompileOnly(t, "let x = 42", "basic variable declaration")
}

func TestFunctionDeclaration(t *testing.T) {
	testCompileOnly(t, "fn add(x, y) = x + y", "function declaration")
}

func TestArithmetic(t *testing.T) {
	tests := []struct {
		name   string
		source string
		output string
	}{
		{
			name:   "addition",
			source: "print(10 + 5)",
			output: "15",
		},
		{
			name:   "subtraction",
			source: "print(10 - 3)",
			output: "7",
		},
		{
			name:   "multiplication",
			source: "print(6 * 7)",
			output: "42",
		},
		{
			name:   "division",
			source: "print(24 / 4)",
			output: "6",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCompileAndRunWithOutput(t, tt.source, tt.output, tt.name)
		})
	}
}

func TestStringLiterals(t *testing.T) {
	testCompileAndRunWithOutput(t, `print("hello")`, "hello", "string literal printing")
}

func TestVariableAssignment(t *testing.T) {
	source := `let x = 10
let y = x
let z = x + y
print(toString(z))`

	testCompileAndRunWithOutput(t, source, "20", "variable assignment and arithmetic")
}

func TestFunctionWithMultipleParameters(t *testing.T) {
	source := `fn calculate(a: int, b: int, c: int) = a + b * c
print(toString(calculate(a: 2, b: 3, c: 4)))`

	testCompileAndRunWithOutput(t, source, "14", "function with multiple parameters")
}

func TestFunctionCall(t *testing.T) {
	source := `fn double(x) = x * 2
print(toString(double(21)))`

	testCompileAndRunWithOutput(t, source, "42", "function call")
}

func TestTypeDeclaration(t *testing.T) {
	testCompileOnly(t, "type Color = Red | Green | Blue", "type declaration")
}

func TestImportStatement(t *testing.T) {
	testCompileOnly(t, "import std", "import statement")
}

func TestComplexExpression(t *testing.T) {
	source := `let result = (10 + 5) * 2 - 3
print(toString(result))`
	testCompileAndRunWithOutput(t, source, "27", "complex arithmetic expression")
}

func TestNestedFunctionCalls(t *testing.T) {
	source := `fn add(x, y) = x + y
fn multiply(a, b) = a * b
print(toString(add(x: multiply(a: 2, b: 3), y: 4)))`

	testCompileAndRunWithOutput(t, source, "10", "nested function calls with named arguments")
}

func TestMatchExpression(t *testing.T) {
	source := `let x = match 42 { 42 => 1 }
print(toString(x))`
	testCompileAndRunWithOutput(t, source, "1", "match expression")
}

func TestMultipleLetDeclarations(t *testing.T) {
	source := `let a = 10
let b = 20
let c = a + b
print(toString(c))`

	testCompileAndRunWithOutput(t, source, "30", "multiple let declarations")
}

func TestZeroValues(t *testing.T) {
	source := `let zero = 0
print(toString(zero))`
	testCompileAndRunWithOutput(t, source, "0", "zero values")
}

func TestNegativeNumbers(t *testing.T) {
	source := `let negative = -42
print(toString(negative))`
	testCompileAndRunWithOutput(t, source, "-42", "negative numbers")
}

func TestEmptyStringLiteral(t *testing.T) {
	source := `print("")`
	testCompileAndRunWithOutput(t, source, "", "empty string literal")
}

func TestPrintInteger(t *testing.T) {
	testCompileAndRunWithOutput(t, `print(42)`, "42", "print integer")
}

func TestPrintString(t *testing.T) {
	testCompileAndRunWithOutput(t, `print("hello")`, "hello", "print string")
}

func TestUnaryMinus(t *testing.T) {
	source := `print(-42)`
	testCompileAndRunWithOutput(t, source, "-42", "unary minus")
}

func TestParenthesizedExpressions(t *testing.T) {
	source := `let result = (10 + 5) * 2
print(toString(result))`
	testCompileAndRunWithOutput(t, source, "30", "parenthesized expressions")
}

func TestOperatorPrecedence(t *testing.T) {
	source := `let result = 2 + 3 * 4
print(toString(result))`
	testCompileAndRunWithOutput(t, source, "14", "operator precedence")
}

func TestComplexArithmetic(t *testing.T) {

	//FIX THIS!!!
	//t.Skip("TODO: Nested arithmetic with division not yet supported - division returns Result, other ops expect Int")
	source := `let result = ((5 + 3) * 2) - (4 / 2)
print(toString(result))`
	testCompileAndRunWithOutput(t, source, "14", "complex arithmetic")
}

func TestFunctionWithNoParameters(t *testing.T) {
	source := `fn getValue() = 42
print(toString(getValue()))`

	testCompileAndRunWithOutput(t, source, "42", "function with no parameters")
}

func TestFunctionCallNoArgs(t *testing.T) {
	source := `fn getValue() = 42
let x = getValue()
print(toString(x))`

	testCompileAndRunWithOutput(t, source, "42", "function call with no arguments")
}

func TestChainedFunctionCalls(t *testing.T) {
	source := `fn double(x) = x * 2
fn quad(x: int) -> int = double(double(x))
print(toString(quad(5)))`

	testCompileAndRunWithOutput(t, source, "20", "chained function calls")
}

func TestStringEscaping(t *testing.T) {
	// Test basic string without problematic escapes for now
	source := `print("Hello World!")`
	testCompileAndRunWithOutput(t, source, "Hello World!", "string with spaces")
}

func TestLargeNumbers(t *testing.T) {
	source := `let big = 1234567890
print(toString(big))`
	testCompileAndRunWithOutput(t, source, "1234567890", "large numbers")
}

func TestZeroAndOne(t *testing.T) {
	source := `let zero = 0
let one = 1
print(toString(zero))
print(toString(one))`

	// We'll need to handle multiple outputs - for now test compilation
	testCompileOnly(t, source, "zero and one values")
}

func TestMatchExpressionSyntax(t *testing.T) {
	source := `let x = match 42 {
    42 => "found"
    _ => "not found"
}
print(toString(x))`

	testCompileOnly(t, source, "complex match expression")
}

func TestTypeDeclarationSyntax(t *testing.T) {
	testCompileOnly(t, `type Color = Red | Green | Blue`, "type declaration syntax")
}

func TestImportSyntax(t *testing.T) {
	testCompileOnly(t, `import std.io`, "import with module path")
}

func TestMutableVariables(t *testing.T) {
	testCompileOnly(t, `mut counter = 0`, "mutable variables")
}

func TestAllOperators(t *testing.T) {
	tests := []struct {
		name   string
		source string
		output string
	}{
		{
			name:   "addition_result",
			source: "print(toString(5 + 3))",
			output: "8",
		},
		{
			name:   "subtraction_result",
			source: "print(toString(10 - 4))",
			output: "6",
		},
		{
			name:   "multiplication_result",
			source: "print(toString(3 * 4))",
			output: "12",
		},
		{
			name:   "division_result",
			source: "print(toString(15 / 3))",
			output: "Success(5)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCompileAndRunWithOutput(t, tt.source, tt.output, tt.name)
		})
	}
}

func TestCompleteProgram(t *testing.T) {
	source := `fn factorial(n: int) -> int = match n {
    0 => 1
    1 => 1
    _ => n * factorial(n - 1)
}

let result = factorial(5)
print(toString(result))`

	testCompileAndRunWithOutput(t, source, "120", "complete program with recursion")
}

// ========== STRING INTERPOLATION TESTS ==========
// Note: Most string interpolation tests are covered comprehensively in integration tests.
// These are kept for unit testing specific LLVM codegen aspects.

func TestBasicStringInterpolation(t *testing.T) {
	source := `
		let name = "Alice"
		print("Hello ${name}!")
	`
	testCompileAndRunWithOutput(t, source, "Hello Alice!", "basic string interpolation")
}

// Test that we can still access the LLVM IR generation for debugging.
func TestLLVMIRGeneration(t *testing.T) {
	ir, err := codegen.CompileToLLVM("let x = 42")
	if err != nil {
		t.Fatalf("Error generating LLVM IR: %v", err)
	}

	if ir == "" {
		t.Fatalf("Generated IR is empty")
	}

	// Basic sanity check - should contain a main function
	if !strings.Contains(ir, "define") || !strings.Contains(ir, "main") {
		t.Errorf("Generated IR doesn't appear to contain a main function:\n%s", ir)
	}
}

// Test JIT compilation detection.
func TestJITToolDetection(t *testing.T) {
	executor := codegen.NewJITExecutor()

	// Try to find LLVM tools - if not found, fail the test
	err := executor.CompileAndRunInMemory("define i32 @main() { ret i32 0 }")
	if err != nil {
		if strings.Contains(err.Error(), "LLVM tools not found") ||
			strings.Contains(err.Error(), "no suitable compiler found") {
			t.Fatalf("❌ LLVM TOOLS NOT AVAILABLE - TEST FAILED: %v", err)
		} else {
			t.Fatalf("Unexpected error during JIT compilation: %v", err)
		}
	}
}
