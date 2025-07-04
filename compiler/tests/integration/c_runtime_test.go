package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// Test wrapper for C Runtime System Tests
func TestSystemRuntimeUnity(t *testing.T) {
	// Get the runtime directory
	runtimeDir := filepath.Join("..", "..", "runtime")

	// Build the system runtime test executable
	cmd := exec.Command("clang",
		"-o", "test_system_runtime_unity",
		"test_system_runtime_unity.c",
		"unity.c",
		"system_runtime.c",
		"-pthread",
		"-std=c11")
	cmd.Dir = runtimeDir

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build system runtime tests: %v\nOutput: %s", err, string(output))
	}

	// Run the test executable
	cmd = exec.Command("./test_system_runtime_unity")
	cmd.Dir = runtimeDir

	output, err := cmd.CombinedOutput()

	// Clean up the executable
	if removeErr := os.Remove(filepath.Join(runtimeDir, "test_system_runtime_unity")); removeErr != nil {
		t.Logf("Warning: Failed to clean up test executable: %v", removeErr)
	}

	if err != nil {
		t.Fatalf("System runtime tests failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("🎉 System runtime tests passed (20 tests):\n%s", string(output))
}

// Test wrapper for C Runtime Fiber Tests
func TestFiberRuntimeUnity(t *testing.T) {
	// Get the runtime directory
	runtimeDir := filepath.Join("..", "..", "runtime")

	// Build the fiber runtime test executable
	cmd := exec.Command("clang",
		"-o", "test_fiber_runtime_unity",
		"test_fiber_runtime_unity.c",
		"unity.c",
		"fiber_runtime.c",
		"system_runtime.c",
		"-pthread",
		"-std=c11")
	cmd.Dir = runtimeDir

	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build fiber runtime tests: %v\nOutput: %s", err, string(output))
	}

	// Run the test executable
	cmd = exec.Command("./test_fiber_runtime_unity")
	cmd.Dir = runtimeDir

	output, err := cmd.CombinedOutput()

	// Clean up the executable
	if removeErr := os.Remove(filepath.Join(runtimeDir, "test_fiber_runtime_unity")); removeErr != nil {
		t.Logf("Warning: Failed to clean up test executable: %v", removeErr)
	}

	if err != nil {
		t.Fatalf("Fiber runtime tests failed: %v\nOutput: %s", err, string(output))
	}

	t.Logf("🎉 Fiber runtime tests passed (60 tests):\n%s", string(output))
}
