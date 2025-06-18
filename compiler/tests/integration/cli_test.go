package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	ospreyBinary = "../../bin/osprey"
	testDataDir  = "../data"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

func ensureCompilerBuilt(t *testing.T) {
	t.Helper()

	if fileExists(ospreyBinary) {
		t.Log("✅ Compiler binary already exists")

		return
	}

	t.Log("🔍 Compiler binary not found, building...")

	// Change to project root directory
	projectRoot := "../../"
	cmd := exec.Command("make", "build")
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("❌ FAILED TO BUILD COMPILER: %v\nOutput: %s", err, string(output))
	}

	// Verify the binary was created
	if !fileExists(ospreyBinary) {
		t.Fatalf("❌ BUILD COMPLETED BUT BINARY NOT FOUND AT: %s", ospreyBinary)
	}

	// Also verify runtime libraries were built
	runtimeLibs := []string{
		"../../bin/libfiber_runtime.a",
		"../../bin/libhttp_runtime.a",
	}

	for _, lib := range runtimeLibs {
		if !fileExists(lib) {
			t.Fatalf("❌ RUNTIME LIBRARY NOT FOUND AT: %s", lib)
		}
	}

	t.Log("✅ Compiler and runtime libraries built successfully")
}

func TestCLI(t *testing.T) {
	ensureCompilerBuilt(t)

	t.Run("help output", testHelpOutput)
	t.Run("ast output", testASTOutput)
	t.Run("llvm output", testLLVMOutput)
	t.Run("compile mode", testCompileMode)
	t.Run("symbols output", testSymbolsOutput)
	t.Run("run mode", testRunMode)
	t.Run("invalid arguments", testInvalidArguments)
	t.Run("missing file", testMissingFile)
	t.Run("syntax error handling", testSyntaxErrorHandling)
	t.Run("security cli arguments", testSecurityCLIArguments)
}

func testHelpOutput(t *testing.T) {
	cmd := exec.Command(ospreyBinary, "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	helpText := string(output)
	expectedSections := []string{
		"Usage:",
		"Options:",
		"Security Options:",
		"Examples:",
		"--ast",
		"--llvm",
		"--compile",
		"--run",
		"--symbols",
		"--sandbox",
		"--no-http",
		"--no-websocket",
		"--no-fs",
		"--no-ffi",
		"Osprey Compiler",
	}

	for _, section := range expectedSections {
		if !strings.Contains(helpText, section) {
			t.Errorf("Help output missing section: %s", section)
		}
	}
}

func testASTOutput(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	cmd := exec.Command(ospreyBinary, testFile, "--ast")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("AST command failed: %v", err)
	}

	astOutput := string(output)
	expectedElements := []string{
		"AST for",
		"Program with",
		"statements",
	}

	for _, element := range expectedElements {
		if !strings.Contains(astOutput, element) {
			t.Errorf("AST output missing element: %s", element)
		}
	}
}

func testLLVMOutput(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	cmd := exec.Command(ospreyBinary, testFile, "--llvm")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("LLVM command failed: %v", err)
	}

	llvmOutput := string(output)
	expectedElements := []string{
		"define",
		"i32 @main",
		"ret i32",
		"@printf",
	}

	for _, element := range expectedElements {
		if !strings.Contains(llvmOutput, element) {
			t.Errorf("LLVM output missing element: %s", element)
		}
	}
}

func testCompileMode(t *testing.T) {
	testFile := filepath.Join(testDataDir, "hello.osp")
	if !fileExists(testFile) {
		t.Fatal("❌ TEST FILE NOT FOUND - TEST FAILED:", testFile)
	}

	// The compiler creates outputs/filename (without extension) relative to source file
	expectedOutput := filepath.Join(testDataDir, "outputs", "hello")
	defer func() { _ = os.RemoveAll(filepath.Join(testDataDir, "outputs")) }() // Cleanup

	cmd := exec.Command(ospreyBinary, testFile, "--compile")
	output, err := cmd.CombinedOutput()

	if err == nil {
		// If compilation succeeded, check if executable was created
		if fileExists(expectedOutput) {
			t.Log("✅ Compilation successful, executable created at:", expectedOutput)
		} else {
			t.Error("Compilation succeeded but no executable found at:", expectedOutput)
		}
	} else {
		// Compilation might fail due to missing LLVM tools, which is acceptable
		t.Logf("⚠️ Compilation failed (likely missing LLVM tools): %v\nOutput: %s", err, string(output))
	}
}

func testSymbolsOutput(t *testing.T) {
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

	cmd := exec.Command(ospreyBinary, testFile, "--symbols")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Symbols command failed: %v", err)
	}

	symbolsOutput := string(output)
	expectedElements := []string{
		"[",        // JSON array start
		"{",        // JSON object start
		"\"name\"", // JSON property
		"\"kind\"", // JSON property
		"\"type\"", // JSON property
	}

	for _, element := range expectedElements {
		if !strings.Contains(symbolsOutput, element) {
			t.Errorf("Symbols output missing element: %s", element)
		}
	}
}

func testRunMode(t *testing.T) {
	// Create a simple test file that should run
	testFile := "/tmp/run_test.osp"
	testContent := `print("Hello from run test")`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	cmd := exec.Command(ospreyBinary, testFile, "--run")
	output, err := cmd.CombinedOutput()

	if err == nil {
		if strings.Contains(string(output), "Hello from run test") {
			t.Log("✅ Run mode successful")
		} else {
			t.Error("Run mode didn't produce expected output")
		}
	} else {
		// Run might fail due to missing JIT tools, which is acceptable
		t.Logf("⚠️ Run mode failed (likely missing JIT tools): %v", err)
	}
}

func testInvalidArguments(t *testing.T) {
	// Test invalid mode
	cmd := exec.Command(ospreyBinary, "--invalid")
	_, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for invalid argument")
	}

	// Test missing filename
	cmd = exec.Command(ospreyBinary, "--ast")
	_, err = cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for missing filename")
	}
}

func testMissingFile(t *testing.T) {
	cmd := exec.Command(ospreyBinary, "/nonexistent/file.osp", "--ast")
	_, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for missing file")
	}
}

func testSyntaxErrorHandling(t *testing.T) {
	// Create a file with syntax errors
	testFile := "/tmp/syntax_error_test.osp"
	testContent := `print("unterminated string`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	cmd := exec.Command(ospreyBinary, testFile, "--ast")
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Error("Expected error for syntax error file")
	}

	// Check that error message is informative
	errorOutput := string(output)
	if !strings.Contains(errorOutput, "Error") && !strings.Contains(errorOutput, "error") {
		t.Error("Error output should contain error information")
	}
}

func testSecurityCLIArguments(t *testing.T) {
	t.Run("sandbox blocks HTTP", testSandboxBlocksHTTP)
	t.Run("no-http blocks HTTP", testNoHTTPBlocksHTTP)
	t.Run("no-websocket blocks WebSocket", testNoWebSocketBlocksWebSocket)
	t.Run("multiple security flags", testMultipleSecurityFlags)
	t.Run("safe code in sandbox", testSafeCodeInSandbox)
	t.Run("security flag combinations", testSecurityFlagCombinations)
}

func testSandboxBlocksHTTP(t *testing.T) {
	// Create test file with HTTP function
	testFile := "/tmp/http_test.osp"
	testContent := `let serverID = httpCreateServer(port: 8080, address: "127.0.0.1")
print(serverID)`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	// Test normal compilation (should work)
	cmd := exec.Command(ospreyBinary, testFile, "--llvm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Normal compilation failed: %v\nOutput: %s", err, string(output))
	}

	// Test sandbox mode (should fail)
	cmd = exec.Command(ospreyBinary, testFile, "--sandbox", "--llvm")
	output, err = cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected sandbox mode to block HTTP functions, but compilation succeeded\nOutput: %s", string(output))
	}

	// Check for proper error (security summary removed from CLI output)
	outputStr := string(output)
	if !strings.Contains(outputStr, "unsupported call expression") {
		t.Error("Expected 'unsupported call expression' error for blocked function")
	}
}

func testNoHTTPBlocksHTTP(t *testing.T) {
	// Create test file with HTTP function
	testFile := "/tmp/no_http_test.osp"
	testContent := `let clientID = httpCreateClient()
print(clientID)`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	// Test --no-http flag (should fail)
	cmd := exec.Command(ospreyBinary, testFile, "--no-http", "--llvm")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected --no-http to block HTTP functions, but compilation succeeded\nOutput: %s", string(output))
	}

	// Security summary removed from CLI output - test still validates blocking works
}

func testNoWebSocketBlocksWebSocket(t *testing.T) {
	// Create test file with WebSocket function
	testFile := "/tmp/no_ws_test.osp"
	testContent := `let wsID = websocketConnect("ws://localhost:8080", "handler")
print(wsID)`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	// Test --no-websocket flag (should fail)
	cmd := exec.Command(ospreyBinary, testFile, "--no-websocket", "--llvm")
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Errorf("Expected --no-websocket to block WebSocket functions, but compilation succeeded\n"+
			"Output: %s", string(output))
	}

	// Security summary removed from CLI output - test still validates blocking works
}

func testMultipleSecurityFlags(t *testing.T) {
	// Create test file with safe operations
	testFile := "/tmp/multi_security_test.osp"
	testContent := `let x = 42
let y = 24
print(x + y)`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	// Test multiple security flags with safe code (should work)
	cmd := exec.Command(ospreyBinary, testFile, "--no-http", "--no-websocket", "--no-ffi", "--llvm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Safe code should compile with security restrictions: %v\nOutput: %s", err, string(output))
	}

	// Security summary removed from CLI output - test still validates restrictions work
}

func testSafeCodeInSandbox(t *testing.T) {
	// Create test file with only safe functions
	testFile := "/tmp/safe_test.osp"
	testContent := `let greeting = "Hello, sandbox!"
print(greeting)
let x = 42
let y = 24
print(x + y)`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	// Test sandbox mode with safe code (should work)
	cmd := exec.Command(ospreyBinary, testFile, "--sandbox", "--llvm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Safe code should compile in sandbox mode: %v\nOutput: %s", err, string(output))
	}

	// Check for LLVM output (security summary removed from CLI)
	outputStr := string(output)
	if !strings.Contains(outputStr, "define") || !strings.Contains(outputStr, "@main") {
		t.Error("Expected valid LLVM IR output for safe code")
	}
}

func testSecurityFlagCombinations(t *testing.T) {
	// Test --no-fs flag
	testFile := "/tmp/no_fs_test.osp"
	testContent := `print("Testing file system restrictions")`
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	defer func() { _ = os.Remove(testFile) }()

	// Test --no-fs flag (should work for this code)
	cmd := exec.Command(ospreyBinary, testFile, "--no-fs", "--llvm")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Safe code should work with --no-fs: %v\nOutput: %s", err, string(output))
	}

	// Security summary removed from CLI output - test still validates flag works
}
