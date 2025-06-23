package cli_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

func TestRunCommand_AST(t *testing.T) {
	// Create test file
	testFile := createTestFile(t, "test_ast.osp", "fn add(a, b) = a + b")
	defer func() { _ = os.Remove(testFile) }()

	result := cli.RunCommand(testFile, cli.OutputModeAST, "")

	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.ErrorMsg)
	}

	if !strings.Contains(result.Output, "AST for") {
		t.Errorf("Expected AST output, got: %s", result.Output)
	}

	if !strings.Contains(result.Output, "Program with") {
		t.Errorf("Expected program statement count, got: %s", result.Output)
	}
}

func TestRunCommand_LLVM(t *testing.T) {
	// Create test file
	testFile := createTestFile(t, "test_llvm.osp", "fn add(a, b) = a + b")
	defer func() { _ = os.Remove(testFile) }()

	result := cli.RunCommand(testFile, cli.OutputModeLLVM, "")

	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.ErrorMsg)
	}

	if !strings.Contains(result.Output, "; LLVM IR for") {
		t.Errorf("Expected LLVM IR output, got: %s", result.Output)
	}
}

func TestRunCommand_Symbols(t *testing.T) {
	// Create test file with function
	testFile := createTestFile(t, "test_symbols.osp", "fn add(a, b) = a + b")
	defer func() { _ = os.Remove(testFile) }()

	result := cli.RunCommand(testFile, cli.OutputModeSymbols, "")

	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.ErrorMsg)
	}

	// Validate JSON output
	var symbols []interface{}
	err := json.Unmarshal([]byte(result.Output), &symbols)
	if err != nil {
		t.Fatalf("Expected valid JSON output, got error: %v", err)
	}
}

func TestRunCommand_Compile(t *testing.T) {
	// Create test file
	testFile := createTestFile(t, "test_compile.osp", "fn add(a, b) = a + b")
	defer func() { _ = os.Remove(testFile) }()
	defer cleanupOutputs(t, testFile)

	result := cli.RunCommand(testFile, cli.OutputModeCompile, "")

	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.ErrorMsg)
	}

	if !strings.Contains(result.Output, "Compiling") {
		t.Errorf("Expected compilation output, got: %s", result.Output)
	}

	if result.OutputFile == "" {
		t.Error("Expected output file to be set")
	}
}

func TestRunCommand_Run(t *testing.T) {
	// Create test file
	testFile := createTestFile(t, "test_run.osp", "fn add(a, b) = a + b")
	defer func() { _ = os.Remove(testFile) }()

	result := cli.RunCommand(testFile, cli.OutputModeRun, "")

	if !result.Success {
		t.Fatalf("Expected success, got error: %s", result.ErrorMsg)
	}

	// The output should be empty for successful runs (just the program output)
	if result.Output != "" {
		t.Errorf("Expected empty output for successful run, got: %s", result.Output)
	}
}

func TestRunCommand_InvalidMode(t *testing.T) {
	// Create test file
	testFile := createTestFile(t, "test_invalid.osp", "fn add(a, b) = a + b")
	defer func() { _ = os.Remove(testFile) }()

	result := cli.RunCommand(testFile, "invalid", "")

	if result.Success {
		t.Fatal("Expected failure for invalid mode")
	}

	if !strings.Contains(result.ErrorMsg, "Unknown output mode") {
		t.Errorf("Expected unknown mode error, got: %s", result.ErrorMsg)
	}
}

func TestRunCommand_FileNotFound(t *testing.T) {
	result := cli.RunCommand("nonexistent.osp", cli.OutputModeAST, "")

	if result.Success {
		t.Fatal("Expected failure for nonexistent file")
	}

	if !strings.Contains(result.ErrorMsg, "Error reading file") {
		t.Errorf("Expected file error, got: %s", result.ErrorMsg)
	}
}

func TestRunCommand_SyntaxError(t *testing.T) {
	// Create test file with syntax error
	testFile := createTestFile(t, "test_syntax_error.osp", "fn invalid syntax {{{}}")
	defer func() { _ = os.Remove(testFile) }()

	result := cli.RunCommand(testFile, cli.OutputModeAST, "")

	if result.Success {
		t.Fatal("Expected failure for syntax error")
	}

	if !strings.Contains(result.ErrorMsg, "Found syntax errors") {
		t.Errorf("Expected syntax error, got: %s", result.ErrorMsg)
	}
}

func TestRunCommand_AllModes(t *testing.T) {
	modes := []string{
		cli.OutputModeAST,
		cli.OutputModeLLVM,
		cli.OutputModeCompile,
		cli.OutputModeRun,
		cli.OutputModeSymbols,
	}

	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			testFile := createTestFile(t, "test_mode_"+mode+".osp", "fn add(a, b) = a + b")
			defer func() { _ = os.Remove(testFile) }()
			defer cleanupOutputs(t, testFile)

			result := cli.RunCommand(testFile, mode, "")

			if !result.Success {
				t.Fatalf("Mode %s failed: %s", mode, result.ErrorMsg)
			}

			// Run mode returns empty output (program output only), other modes should have output
			if mode == cli.OutputModeRun {
				if result.Output != "" {
					t.Errorf("Mode %s should produce empty output, got: %s", mode, result.Output)
				}
			} else {
				if result.Output == "" {
					t.Errorf("Mode %s produced no output", mode)
				}
			}
		})
	}
}

// Helper functions

func createTestFile(t *testing.T, filename, content string) string {
	t.Helper()

	testFile := filepath.Join(t.TempDir(), filename)
	err := os.WriteFile(testFile, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	return testFile
}

func cleanupOutputs(t *testing.T, testFile string) {
	t.Helper()

	dir := filepath.Dir(testFile)
	outputsDir := filepath.Join(dir, "outputs")
	_ = os.RemoveAll(outputsDir)
}
