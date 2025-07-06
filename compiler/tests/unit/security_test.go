package cli_test

import (
	"os"
	"strings"
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

func TestSandboxModeBlocksHTTPFunctions(t *testing.T) {
	// Test program that uses HTTP functions
	source := `let serverID = httpCreateServer(port: 8080, address: "127.0.0.1")
print(serverID)`

	// Create test file
	testFile := createTestFile(t, "test_http.osp", source)
	defer func() { _ = os.Remove(testFile) }()

	// Create sandbox security configuration
	security := cli.NewSandboxSecurityConfig()

	// Try to compile with sandbox mode
	result := cli.RunCommandWithSecurity(testFile, cli.OutputModeLLVM, false, security)

	// Should fail because httpCreateServer function doesn't exist in sandbox mode
	if result.Success {
		t.Errorf("Expected compilation to fail in sandbox mode, but it succeeded")
	}

	// Error message should mention unsupported call expression (function doesn't exist)
	if !strings.Contains(result.ErrorMsg, "unsupported call expression") {
		t.Errorf("Expected unsupported call expression error for httpCreateServer, got: %s", result.ErrorMsg)
	}
}

func TestSandboxModeBlocksWebSocketFunctions(t *testing.T) {
	// Test program that uses WebSocket functions
	source := `let wsID = websocketConnect("ws://localhost:8080", "handler")
print(wsID)`

	// Create test file
	testFile := createTestFile(t, "test_ws.osp", source)
	defer func() { _ = os.Remove(testFile) }()

	// Create sandbox security configuration
	security := cli.NewSandboxSecurityConfig()

	// Try to compile with sandbox mode
	result := cli.RunCommandWithSecurity(testFile, cli.OutputModeLLVM, false, security)

	// Should fail because websocketConnect function doesn't exist in sandbox mode
	if result.Success {
		t.Errorf("Expected compilation to fail in sandbox mode, but it succeeded")
	}

	// Error message should mention unsupported call expression (function doesn't exist)
	if !strings.Contains(result.ErrorMsg, "unsupported call expression") {
		t.Errorf("Expected unsupported call expression error for websocketConnect, got: %s", result.ErrorMsg)
	}
}

func TestDefaultSecurityAllowsAllFunctions(t *testing.T) {
	// Test program that uses safe functions
	source := `let x = 42
print(x)`

	// Create test file
	testFile := createTestFile(t, "test_safe.osp", source)
	defer func() { _ = os.Remove(testFile) }()

	// Create default (permissive) security configuration
	security := cli.NewDefaultSecurityConfig()

	// Should succeed with default security
	result := cli.RunCommandWithSecurity(testFile, cli.OutputModeLLVM, false, security)

	if !result.Success {
		t.Errorf("Expected compilation to succeed with default security, but it failed: %s", result.ErrorMsg)
	}

	// Result should contain LLVM IR
	if !strings.Contains(result.Output, "; LLVM IR for") {
		t.Errorf("Expected LLVM IR output, got: %s", result.Output)
	}
}

func TestSecuritySummaryGeneration(t *testing.T) {
	// Test sandbox mode summary
	sandbox := cli.NewSandboxSecurityConfig()
	summary := sandbox.GetSecuritySummary()

	if !strings.Contains(summary, "SANDBOX MODE") {
		t.Errorf("Expected sandbox mode in summary, got: %s", summary)
	}

	// Test partial restrictions summary
	partial := cli.NewDefaultSecurityConfig()
	partial.AllowHTTP = false
	partial.AllowWebSocket = false

	summary = partial.GetSecuritySummary()

	if !strings.Contains(summary, "Unavailable=[HTTP,WebSocket]") {
		t.Errorf("Expected unavailable HTTP and WebSocket in summary, got: %s", summary)
	}

	if !strings.Contains(summary, "Available=[FileRead,FileWrite,FFI]") {
		t.Errorf("Expected available FileRead, FileWrite, FFI in summary, got: %s", summary)
	}
}

func TestBlockedFunctionsList(t *testing.T) {
	// Test sandbox mode blocked functions
	sandbox := cli.NewSandboxSecurityConfig()
	blocked := sandbox.GetBlockedFunctions()

	// Should block HTTP functions
	httpFunctions := []string{"httpCreateServer", "httpListen", "httpGet", "httpPost"}
	for _, fn := range httpFunctions {
		found := false
		for _, blocked := range blocked {
			if blocked == fn {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected function %s to be unavailable in sandbox mode", fn)
		}
	}

	// Should block WebSocket functions
	wsFunctions := []string{"websocketConnect", "websocketSend", "websocketClose"}
	for _, fn := range wsFunctions {
		found := false
		for _, blocked := range blocked {
			if blocked == fn {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected function %s to be unavailable in sandbox mode", fn)
		}
	}
}
