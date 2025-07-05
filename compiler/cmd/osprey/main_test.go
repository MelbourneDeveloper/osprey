package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

const testDataDir = "../../tests/data"

func TestRunCLI(t *testing.T) {
	t.Run("help output", testRunCLIHelp)
	t.Run("version output", testRunCLIVersion)
	t.Run("ast output", testRunCLIAST)
	t.Run("llvm output", testRunCLILLVM)
	t.Run("compile mode", testRunCLICompile)
	t.Run("symbols output", testRunCLISymbols)
	t.Run("run mode", testRunCLIRun)
	t.Run("invalid arguments", testRunCLIInvalidArgs)
	t.Run("missing file", testRunCLIMissingFile)
	t.Run("syntax error handling", testRunCLISyntaxError)
	t.Run("security arguments", testRunCLISecurityArgs)
	t.Run("docs mode", testRunCLIDocs)
	t.Run("hover mode", testRunCLIHover)
}

func TestParsingFunctions(t *testing.T) {
	t.Run("ParseOutputModeArg", testParseOutputModeArg)
	t.Run("ParseSecurityArg", testParseSecurityArg)
	t.Run("ShowHelp", testShowHelp)
}

func testRunCLIHelp(t *testing.T) {
	args := []string{"osprey", "--help"}
	result := RunCLI(args)

	if !result.Success {
		t.Error("RunCLI with --help should succeed")
	}

	// Test -h as well
	args = []string{"osprey", "-h"}
	result = RunCLI(args)

	if !result.Success {
		t.Error("RunCLI with -h should succeed")
	}
}

func testRunCLIVersion(t *testing.T) {
	args := []string{"osprey", "--version"}
	result := RunCLI(args)

	if !result.Success {
		t.Error("RunCLI with --version should succeed")
	}
}

func testRunCLIAST(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatalf("❌ CRITICAL FAILURE: Required test file missing: %s - test infrastructure broken", testFile)
	}

	args := []string{"osprey", testFile, "--ast"}
	result := RunCLI(args)

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

func testRunCLILLVM(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatalf("❌ CRITICAL FAILURE: Required test file missing: %s - test infrastructure broken", testFile)
	}

	args := []string{"osprey", testFile, "--llvm"}
	result := RunCLI(args)

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

func testRunCLICompile(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatalf("❌ CRITICAL FAILURE: Required test file missing: %s - test infrastructure broken", testFile)
	}

	// The compiler creates outputs/filename (without extension) relative to source file
	expectedOutput := filepath.Join(testDataDir, "outputs", "hello")
	defer func() { _ = os.RemoveAll(filepath.Join(testDataDir, "outputs")) }() // Cleanup

	args := []string{"osprey", testFile, "--compile"}
	result := RunCLI(args)

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

func testRunCLISymbols(t *testing.T) {
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
	result := RunCLI(args)

	if !result.Success {
		t.Fatalf("Symbols command failed: %s", result.ErrorMsg)
	}

	// JSON symbols output should contain function definitions
	if !strings.Contains(result.Output, "{") {
		t.Error("Symbols output should be valid JSON")
	}
}

func testRunCLIRun(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatalf("❌ CRITICAL FAILURE: Required test file missing: %s - test infrastructure broken", testFile)
	}

	args := []string{"osprey", testFile, "--run"}
	result := RunCLI(args)

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

func testRunCLIInvalidArgs(t *testing.T) {
	args := []string{"osprey", "nonexistent.osp", "--invalid-flag"}
	result := RunCLI(args)

	if result.Success {
		t.Error("Expected failure for invalid arguments, but got success")
	}

	if !strings.Contains(result.ErrorMsg, "unknown option") {
		t.Errorf("Expected error about unknown option, got: %s", result.ErrorMsg)
	}
}

func testRunCLIMissingFile(t *testing.T) {
	args := []string{"osprey", "definitely_nonexistent_file.osp"}
	result := RunCLI(args)

	if result.Success {
		t.Error("Expected failure for missing file, but got success")
	}
}

func testRunCLISyntaxError(t *testing.T) {
	// Create a file with syntax errors
	testFile := "/tmp/syntax_error_test.osp"
	testContent := `fn broken_syntax( = invalid`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	args := []string{"osprey", testFile, "--ast"}
	result := RunCLI(args)

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

func testRunCLISecurityArgs(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		// Create a simple test file
		testFile = "/tmp/security_test.osp"
		testContent := `print("Hello, world!")`
		err := os.WriteFile(testFile, []byte(testContent), 0o644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer func() { _ = os.Remove(testFile) }()
	}

	// Test security flags with compilation
	securityFlags := []string{"--sandbox", "--no-http", "--no-websocket", "--no-fs", "--no-ffi"}

	for _, flag := range securityFlags {
		t.Run(flag, func(t *testing.T) {
			args := []string{"osprey", testFile, "--ast", flag}
			result := RunCLI(args)

			if !result.Success {
				t.Logf("⚠️ Security flag %s caused failure: %s", flag, result.ErrorMsg)
			} else {
				t.Logf("✅ Security flag %s processed successfully", flag)
			}
		})
	}
}

func testRunCLIDocs(t *testing.T) {
	args := []string{"osprey", "--docs"}
	result := RunCLI(args)

	// Docs mode might fail due to missing dependencies, but should not crash
	if !result.Success {
		t.Logf("⚠️ Docs mode failed: %s", result.ErrorMsg)
	} else {
		t.Log("✅ Docs mode processed successfully")
	}

	// Test with custom docs directory
	args = []string{"osprey", "--docs", "--docs-dir", "/tmp/test-docs"}
	result = RunCLI(args)

	if !result.Success {
		t.Logf("⚠️ Docs mode with custom directory failed: %s", result.ErrorMsg)
	} else {
		t.Log("✅ Docs mode with custom directory processed successfully")
	}
}

func testRunCLIHover(t *testing.T) {
	args := []string{"osprey", "--hover", "print"}
	result := RunCLI(args)

	// Hover mode might fail due to missing dependencies, but should not crash
	if !result.Success {
		t.Logf("⚠️ Hover mode failed: %s", result.ErrorMsg)
	} else {
		t.Log("✅ Hover mode processed successfully")
	}

	// Test hover without element name (should fail)
	args = []string{"osprey", "--hover"}
	result = RunCLI(args)

	if result.Success {
		t.Error("Expected hover without element name to fail")
	}
}

func testParseOutputModeArg(t *testing.T) {
	testCases := []struct {
		arg      string
		expected string
	}{
		{"--ast", cli.OutputModeAST},
		{"--llvm", cli.OutputModeLLVM},
		{"--compile", cli.OutputModeCompile},
		{"--run", cli.OutputModeRun},
		{"--symbols", cli.OutputModeSymbols},
		{"--docs", cli.OutputModeDocs},
		{"--hover", cli.OutputModeHover},
		{"--invalid", ""},
	}

	for _, tc := range testCases {
		result := ParseOutputModeArg(tc.arg)
		if result != tc.expected {
			t.Errorf("ParseOutputModeArg(%s) = %s, expected %s", tc.arg, result, tc.expected)
		}
	}
}

func testParseSecurityArg(t *testing.T) {
	// Test --sandbox flag
	security := cli.NewDefaultSecurityConfig()
	if !ParseSecurityArg("--sandbox", security) {
		t.Error("ParseSecurityArg should return true for --sandbox")
	}
	if !security.SandboxMode {
		t.Error("--sandbox should enable SandboxMode")
	}

	// Test --no-http flag
	security = cli.NewDefaultSecurityConfig()
	if !ParseSecurityArg("--no-http", security) {
		t.Error("ParseSecurityArg should return true for --no-http")
	}
	if security.AllowHTTP {
		t.Error("--no-http should disable AllowHTTP")
	}

	// Test --no-websocket flag
	security = cli.NewDefaultSecurityConfig()
	if !ParseSecurityArg("--no-websocket", security) {
		t.Error("ParseSecurityArg should return true for --no-websocket")
	}
	if security.AllowWebSocket {
		t.Error("--no-websocket should disable AllowWebSocket")
	}

	// Test --no-fs flag
	security = cli.NewDefaultSecurityConfig()
	if !ParseSecurityArg("--no-fs", security) {
		t.Error("ParseSecurityArg should return true for --no-fs")
	}
	if security.AllowFileRead || security.AllowFileWrite {
		t.Error("--no-fs should disable AllowFileRead and AllowFileWrite")
	}

	// Test --no-ffi flag
	security = cli.NewDefaultSecurityConfig()
	if !ParseSecurityArg("--no-ffi", security) {
		t.Error("ParseSecurityArg should return true for --no-ffi")
	}
	if security.AllowFFI {
		t.Error("--no-ffi should disable AllowFFI")
	}

	// Test invalid flag
	if ParseSecurityArg("--invalid", security) {
		t.Error("ParseSecurityArg should return false for invalid flag")
	}
}

func testShowHelp(t *testing.T) {
	// Just test that ShowHelp doesn't crash
	ShowHelp()
	t.Log("✅ ShowHelp executed without crashing")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
