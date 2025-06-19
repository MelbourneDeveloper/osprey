// Package util provides shared helper utilities for test suites.
// In particular it contains CleanAndRebuildAll which performs a full, deterministic
// rebuild of the compiler and its native/Rust runtime dependencies before running tests.
package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// CleanAndRebuildAll cleans the repository, rebuilds runtime libraries, the Rust interop library and the compiler.
// It panics on any failure, ensuring test suites fail fast.
func CleanAndRebuildAll(projectRoot string) {
	// 1. Clean artifacts (including Rust targets and /tmp outputs)
	cmd := exec.Command("make", "clean")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to clean: " + err.Error() + "\nOutput: " + string(output))
	}

	// 2. Re-build native runtime libraries
	cmd = exec.Command("make", "fiber-runtime", "http-runtime")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build runtime libraries: " + err.Error() + "\nOutput: " + string(output))
	}

	// 3. Re-build Rust interop library (when present)
	rustDir := filepath.Join(projectRoot, "examples", "rust_integration")
	if _, err := os.Stat(rustDir); err == nil {
		cleanCmd := exec.Command("cargo", "clean")
		cleanCmd.Dir = rustDir
		if output, err := cleanCmd.CombinedOutput(); err != nil {
			panic("Failed to clean Rust artifacts: " + err.Error() + "\nOutput: " + string(output))
		}

		targetDir := filepath.Join(os.TempDir(), "osprey_rust_target_"+strconv.Itoa(os.Getpid()))
		cmd = exec.Command("cargo", "build", "--target-dir", targetDir, "-j", "1")
		cmd.Dir = rustDir
		if output, err := cmd.CombinedOutput(); err != nil {
			panic("Failed to build Rust interop: " + err.Error() + "\nOutput: " + string(output))
		}
	}

	// 4. Re-build the compiler binary (skipping lint for speed)
	cmd = exec.Command("go", "build", "-o", "bin/osprey", "./cmd/osprey")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build compiler: " + err.Error() + "\nOutput: " + string(output))
	}
}
