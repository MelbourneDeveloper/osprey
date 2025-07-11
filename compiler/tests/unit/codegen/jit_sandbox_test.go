package codegen

import (
	"testing"

	"github.com/christianfindlay/osprey/internal/cli"
)

func TestSandboxModeDisablesAllSecurityFlags(t *testing.T) {
	// Test that --sandbox flag sets all security flags to false
	args := []string{"osprey", "test.osp", "--sandbox"}

	filename, outputMode, docsDir, quiet, security := cli.ParseArgs(args)
	_ = filename
	_ = outputMode
	_ = docsDir
	_ = quiet

	if security == nil {
		t.Fatal("Security config should not be nil")
	}

	// Verify ALL security flags are disabled in sandbox mode
	if security.AllowHTTP {
		t.Error("AllowHTTP should be false in sandbox mode")
	}
	if security.AllowWebSocket {
		t.Error("AllowWebSocket should be false in sandbox mode")
	}
	if security.AllowFileRead {
		t.Error("AllowFileRead should be false in sandbox mode")
	}
	if security.AllowFileWrite {
		t.Error("AllowFileWrite should be false in sandbox mode")
	}
	if security.AllowFFI {
		t.Error("AllowFFI should be false in sandbox mode")
	}
	if security.AllowProcessExecution {
		t.Error("AllowProcessExecution should be false in sandbox mode")
	}
}

func TestDefaultModeEnablesAllSecurityFlags(t *testing.T) {
	// Test that without --sandbox flag, all security flags are enabled by default
	args := []string{"osprey", "test.osp"}

	filename, outputMode, docsDir, quiet, security := cli.ParseArgs(args)
	_ = filename
	_ = outputMode
	_ = docsDir
	_ = quiet

	if security == nil {
		t.Fatal("Security config should not be nil")
	}

	// Verify ALL security flags are enabled by default
	if !security.AllowHTTP {
		t.Error("AllowHTTP should be true by default")
	}
	if !security.AllowWebSocket {
		t.Error("AllowWebSocket should be true by default")
	}
	if !security.AllowFileRead {
		t.Error("AllowFileRead should be true by default")
	}
	if !security.AllowFileWrite {
		t.Error("AllowFileWrite should be true by default")
	}
	if !security.AllowFFI {
		t.Error("AllowFFI should be true by default")
	}
	if !security.AllowProcessExecution {
		t.Error("AllowProcessExecution should be true by default")
	}
}

func TestSandboxModeWithOtherFlags(t *testing.T) {
	// Test that --sandbox overrides other security flags
	args := []string{"osprey", "test.osp", "--sandbox", "--run"}

	filename, outputMode, docsDir, quiet, security := cli.ParseArgs(args)
	_ = filename
	_ = outputMode
	_ = docsDir
	_ = quiet

	if security == nil {
		t.Fatal("Security config should not be nil")
	}

	// Even with other flags, sandbox mode should disable everything
	if security.AllowHTTP || security.AllowWebSocket || security.AllowFileRead ||
		security.AllowFileWrite || security.AllowFFI || security.AllowProcessExecution {
		t.Error("Sandbox mode should disable all security flags regardless of other flags")
	}
}
