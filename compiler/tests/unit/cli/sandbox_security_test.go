package cli

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

// TestSandboxWithRunMode verifies --sandbox properly disables all security features in run mode
func TestSandboxWithRunMode(t *testing.T) {
	args := []string{"osprey", "test.osp", "--sandbox", "--run"}
	verifyIsSandboxed(t, args, "run")
}

// TestSandboxWithCompileMode verifies --sandbox properly disables all security features in compile mode
func TestSandboxWithCompileMode(t *testing.T) {
	args := []string{"osprey", "test.osp", "--sandbox", "--compile"}
	verifyIsSandboxed(t, args, "compile")
}

// TestSandboxWithAstMode verifies --sandbox properly disables all security features in ast mode
func TestSandboxWithAstMode(t *testing.T) {
	args := []string{"osprey", "test.osp", "--sandbox", "--ast"}
	verifyIsSandboxed(t, args, "ast")
}

// TestSandboxWithLlvmMode verifies --sandbox properly disables all security features in llvm mode
func TestSandboxWithLlvmMode(t *testing.T) {
	args := []string{"osprey", "test.osp", "--sandbox", "--llvm"}
	verifyIsSandboxed(t, args, "llvm")
}

// verifyIsSandboxed is a helper function that verifies all security features are disabled
func verifyIsSandboxed(t *testing.T, args []string, expectedMode string) {
	filename, outputMode, docsDir, quiet, security := cli.ParseArgs(args)

	// Verify basic parsing worked
	if filename != "test.osp" {
		t.Errorf("Expected filename 'test.osp', got '%s'", filename)
	}
	if outputMode != expectedMode {
		t.Errorf("Expected mode '%s', got '%s'", expectedMode, outputMode)
	}
	if docsDir != "" {
		t.Errorf("Expected empty docsDir, got '%s'", docsDir)
	}
	if quiet {
		t.Errorf("Expected quiet=false, got true")
	}
	if security == nil {
		t.Fatal("Expected security config, got nil")
	}

	// THE CRITICAL TEST: Verify ALL security flags are disabled
	verifyAllSecurityFlagsDisabled(t, security)
}

// verifyAllSecurityFlagsDisabled verifies that all security flags are properly disabled
func verifyAllSecurityFlagsDisabled(t *testing.T, security *cli.SecurityConfig) {
	if security.AllowHTTP {
		t.Error("Expected AllowHTTP=false, got true - SECURITY VIOLATION!")
	}
	if security.AllowWebSocket {
		t.Error("Expected AllowWebSocket=false, got true - SECURITY VIOLATION!")
	}
	if security.AllowFileRead {
		t.Error("Expected AllowFileRead=false, got true - SECURITY VIOLATION!")
	}
	if security.AllowFileWrite {
		t.Error("Expected AllowFileWrite=false, got true - SECURITY VIOLATION!")
	}
	if security.AllowFFI {
		t.Error("Expected AllowFFI=false, got true - SECURITY VIOLATION!")
	}
	if security.AllowProcessExecution {
		t.Error("Expected AllowProcessExecution=false, got true - SECURITY VIOLATION!")
	}
}

// TestDefaultConfigIsPermissive verifies that without --sandbox, all features are enabled
func TestDefaultConfigIsPermissive(t *testing.T) {
	args := []string{"osprey", "test.osp", "--run"}
	filename, outputMode, docsDir, quiet, security := cli.ParseArgs(args)

	// Verify basic parsing worked
	if filename != "test.osp" {
		t.Errorf("Expected filename 'test.osp', got '%s'", filename)
	}
	if outputMode != "run" {
		t.Errorf("Expected mode 'run', got '%s'", outputMode)
	}
	if docsDir != "" {
		t.Errorf("Expected empty docsDir, got '%s'", docsDir)
	}
	if quiet {
		t.Errorf("Expected quiet=false, got true")
	}
	if security == nil {
		t.Fatal("Expected security config, got nil")
	}

	// Verify ALL security flags are enabled by default
	if !security.AllowHTTP {
		t.Error("Expected AllowHTTP=true, got false")
	}
	if !security.AllowWebSocket {
		t.Error("Expected AllowWebSocket=true, got false")
	}
	if !security.AllowFileRead {
		t.Error("Expected AllowFileRead=true, got false")
	}
	if !security.AllowFileWrite {
		t.Error("Expected AllowFileWrite=true, got false")
	}
	if !security.AllowFFI {
		t.Error("Expected AllowFFI=true, got false")
	}
	if !security.AllowProcessExecution {
		t.Error("Expected AllowProcessExecution=true, got false")
	}
}

// TestSandboxConfigCreation verifies NewSandboxSecurityConfig creates proper sandbox config
func TestSandboxConfigCreation(t *testing.T) {
	config := cli.NewSandboxSecurityConfig()

	if config == nil {
		t.Fatal("Expected security config, got nil")
	}

	// Verify ALL security flags are disabled in sandbox config
	verifyAllSecurityFlagsDisabled(t, config)
}

// TestApplySandboxMode verifies ApplySandboxMode disables all security features
func TestApplySandboxMode(t *testing.T) {
	// Start with permissive config
	config := cli.NewDefaultSecurityConfig()

	// Verify it starts permissive
	if !config.AllowHTTP {
		t.Error("Expected AllowHTTP=true initially, got false")
	}
	if !config.AllowWebSocket {
		t.Error("Expected AllowWebSocket=true initially, got false")
	}

	// Apply sandbox mode
	config.ApplySandboxMode()

	// Verify ALL security flags are now disabled
	verifyAllSecurityFlagsDisabled(t, config)
}

// TestSandboxSecuritySummary verifies the security summary indicates sandbox mode
func TestSandboxSecuritySummary(t *testing.T) {
	config := cli.NewSandboxSecurityConfig()
	summary := config.GetSecuritySummary()

	expectedSummary := "Security: SANDBOX MODE - Only safe functions available"
	if summary != expectedSummary {
		t.Errorf("Expected summary '%s', got '%s'", expectedSummary, summary)
	}
}

// TestSandboxBlockedFunctions verifies sandbox mode blocks dangerous functions
func TestSandboxBlockedFunctions(t *testing.T) {
	config := cli.NewSandboxSecurityConfig()
	blocked := config.GetBlockedFunctions()

	// These functions MUST be blocked in sandbox mode
	expectedBlocked := []string{
		// HTTP functions
		"httpCreateServer", "httpListen", "httpStopServer",
		"httpCreateClient", "httpGet", "httpPost", "httpPut",
		"httpDelete", "httpRequest", "httpCloseClient",
		// WebSocket functions
		"websocketConnect", "websocketSend", "websocketClose",
		"websocketCreateServer", "websocketServerListen",
		"websocketServerSend", "websocketServerBroadcast",
		"websocketStopServer",
		// File I/O functions
		"readFile", "openFile", "writeFile", "createFile", "deleteFile",
		// FFI functions
		"extern",
	}

	// Verify all expected functions are blocked
	for _, expectedFunc := range expectedBlocked {
		found := false
		for _, blockedFunc := range blocked {
			if blockedFunc == expectedFunc {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected function '%s' to be blocked in sandbox mode, but it wasn't - SECURITY VIOLATION!", expectedFunc)
		}
	}

	// Verify we have all the expected blocked functions
	if len(blocked) != len(expectedBlocked) {
		t.Errorf("Expected %d blocked functions, got %d. Blocked: %v", len(expectedBlocked), len(blocked), blocked)
	}
}
