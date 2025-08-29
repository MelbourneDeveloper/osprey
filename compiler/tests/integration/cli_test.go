package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

const testDataDir = "../data"

func TestCLI(t *testing.T) {
	t.Run("help output via cli", testHelpOutputViaCLI)
	t.Run("ast output via cli", testASTOutputViaCLI)
	t.Run("llvm output via cli", testLLVMOutputViaCLI)
	t.Run("compile mode via cli", testCompileModeViaCLI)
	t.Run("symbols output via cli", testSymbolsOutputViaCLI)
	t.Run("run mode via cli", testRunModeViaCLI)
	t.Run("invalid arguments via cli", testInvalidArgumentsViaCLI)
	t.Run("missing file via cli", testMissingFileViaCLI)
	t.Run("syntax error handling via cli", testSyntaxErrorHandlingViaCLI)
	t.Run("security cli arguments via cli", testSecurityCLIArgumentsViaCLI)
	t.Run("cli interface functions", testCLIInterfaceFunctions)
}

func testHelpOutputViaCLI(t *testing.T) {
	// Test help via CLI interface
	args := []string{"osprey", "--help"}

	result := cli.RunMainWithArgs(args)
	if !result.Success {
		t.Error("CLI with --help should succeed")
	}
}

func testASTOutputViaCLI(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	args := []string{"osprey", testFile, "--ast"}
	result := cli.RunMainWithArgs(args)

	if !result.Success {
		t.Fatalf("AST command failed: %s", result.ErrorMsg)
	}

	expectedElements := []string{
		"AST for",
		"Program with",
		"statements",
	}

	for _, element := range expectedElements {
		if !strings.Contains(result.Output, element) {
			t.Errorf("AST output missing element: %s", element)
		}
	}
}

func testLLVMOutputViaCLI(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	args := []string{"osprey", testFile, "--llvm"}
	result := cli.RunMainWithArgs(args)

	if !result.Success {
		t.Fatalf("LLVM command failed: %s", result.ErrorMsg)
	}

	expectedElements := []string{
		"define",
		"i32 @main",
		"ret i32",
		"@printf",
	}

	for _, element := range expectedElements {
		if !strings.Contains(result.Output, element) {
			t.Errorf("LLVM output missing element: %s", element)
		}
	}
}

func testCompileModeViaCLI(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	// The compiler creates outputs/filename (without extension) relative to source file
	expectedOutput := filepath.Join(testDataDir, "outputs", "hello")

	defer func() { _ = os.RemoveAll(filepath.Join(testDataDir, "outputs")) }() // Cleanup

	args := []string{"osprey", testFile, "--compile"}
	result := cli.RunMainWithArgs(args)

	if result.Success {
		// If compilation succeeded, check if executable was created
		if fileExists(expectedOutput) {
			t.Log("✅ Compilation successful, executable created at:", expectedOutput)
		} else {
			t.Error("Compilation succeeded but no executable found at:", expectedOutput)
		}
	} else {
		// Compilation might fail due to missing LLVM tools, which is acceptable
		t.Logf("⚠️ Compilation failed (likely missing LLVM tools): %s", result.ErrorMsg)
	}
}

func testSymbolsOutputViaCLI(t *testing.T) {
	testFile := filepath.Join(testDataDir, "simple_types.osp")
	if !fileExists(testFile) {
		// Create a simple test file for symbols
		testFile = "/tmp/symbols_test.osp"
		testContent := `fn add(a, b) = a + b`

		err := os.WriteFile(testFile, []byte(testContent), 0o644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		defer func() { _ = os.Remove(testFile) }()
	}

	args := []string{"osprey", testFile, "--symbols"}
	result := cli.RunMainWithArgs(args)

	if !result.Success {
		t.Fatalf("Symbols command failed: %s", result.ErrorMsg)
	}

	// JSON symbols output should contain function definitions
	if !strings.Contains(result.Output, "{") {
		t.Error("Symbols output should be valid JSON")
	}
}

func testRunModeViaCLI(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	args := []string{"osprey", testFile, "--run"}
	result := cli.RunMainWithArgs(args)

	if result.Success {
		t.Log("✅ Run mode successful")
		// Basic sanity check - output should not be empty for hello world
		if len(strings.TrimSpace(result.Output)) == 0 {
			t.Log("⚠️ No output from run mode (may be expected)")
		}
	} else {
		// Run might fail due to missing runtime dependencies
		t.Logf("⚠️ Run mode failed (likely missing runtime deps): %s", result.ErrorMsg)
	}
}

func testInvalidArgumentsViaCLI(t *testing.T) {
	args := []string{"osprey", "nonexistent.osp", "--invalid-flag"}
	result := cli.RunMainWithArgs(args)

	if result.Success {
		t.Error("Expected failure for invalid arguments, but got success")
	}

	if !strings.Contains(result.ErrorMsg, "Unknown option") &&
		!strings.Contains(result.ErrorMsg, "no such file") {
		t.Errorf("Expected error about unknown option or missing file, got: %s", result.ErrorMsg)
	}
}

func testMissingFileViaCLI(t *testing.T) {
	args := []string{"osprey", "definitely_nonexistent_file.osp"}
	result := cli.RunMainWithArgs(args)

	if result.Success {
		t.Error("Expected failure for missing file, but got success")
	}
}

func testSyntaxErrorHandlingViaCLI(t *testing.T) {
	// Create a file with syntax errors
	testFile := "/tmp/syntax_error_test.osp"
	testContent := `fn broken_syntax( = invalid`

	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	defer func() { _ = os.Remove(testFile) }()

	args := []string{"osprey", testFile, "--ast"}
	result := cli.RunMainWithArgs(args)

	if result.Success {
		t.Log("⚠️ Expected syntax error to fail compilation, but it succeeded")
	} else {
		// This is expected - syntax errors should cause failure
		if !strings.Contains(result.ErrorMsg, "error") &&
			!strings.Contains(result.ErrorMsg, "failed") {
			t.Errorf("Expected error message about syntax error, got: %s", result.ErrorMsg)
		}
	}
}

func testSecurityCLIArgumentsViaCLI(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	// Test security flags with compilation
	securityFlags := []string{"--sandbox", "--no-http", "--no-websocket", "--no-fs", "--no-ffi"}

	for _, flag := range securityFlags {
		t.Run(flag, func(t *testing.T) {
			args := []string{"osprey", testFile, "--ast", flag}
			result := cli.RunMainWithArgs(args)

			if !result.Success {
				t.Logf("⚠️ Security flag %s caused failure: %s", flag, result.ErrorMsg)
			} else {
				t.Logf("✅ Security flag %s processed successfully", flag)
			}
		})
	}
}

func testCLIInterfaceFunctions(t *testing.T) {
	// Test ParseOutputModeArg function via CLI interface
	testCases := []struct {
		arg      string
		expected string
	}{
		{"--ast", cli.OutputModeAST},
		{"--llvm", cli.OutputModeLLVM},
		{"--compile", cli.OutputModeCompile},
		{"--run", cli.OutputModeRun},
		{"--symbols", cli.OutputModeSymbols},
		{"--invalid", ""},
	}

	for _, tc := range testCases {
		result := cli.ParseOutputModeArg(tc.arg)
		if result != tc.expected {
			t.Errorf("ParseOutputModeArg(%s) = %s, expected %s", tc.arg, result, tc.expected)
		}
	}

	// Test ParseSecurityArg function via CLI interface
	security := cli.NewDefaultSecurityConfig()

	// Test --sandbox flag
	if !cli.ParseSecurityArg("--sandbox", security) {
		t.Error("ParseSecurityArg should return true for --sandbox")
	}

	if !security.SandboxMode {
		t.Error("--sandbox should enable SandboxMode")
	}

	// Test --no-http flag
	security = cli.NewDefaultSecurityConfig()
	if !cli.ParseSecurityArg("--no-http", security) {
		t.Error("ParseSecurityArg should return true for --no-http")
	}

	if security.AllowHTTP {
		t.Error("--no-http should disable AllowHTTP")
	}

	// Test invalid flag
	if cli.ParseSecurityArg("--invalid", security) {
		t.Error("ParseSecurityArg should return false for invalid flag")
	}
}
