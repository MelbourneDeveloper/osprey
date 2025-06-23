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

	// 2. Re-build native runtime libraries (build sequentially to avoid race conditions)
	// First ensure the bin directory exists
	binDir := filepath.Join(projectRoot, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		panic("Failed to create bin directory: " + err.Error())
	}

	// Build fiber runtime first
	cmd = exec.Command("make", "fiber-runtime")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build fiber runtime: " + err.Error() + "\nOutput: " + string(output))
	}

	// Build HTTP runtime second
	cmd = exec.Command("make", "http-runtime")
	cmd.Dir = projectRoot
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build HTTP runtime: " + err.Error() + "\nOutput: " + string(output))
	}

	// 3. Re-build Rust interop library (when present)
	rustDir := filepath.Join(projectRoot, "examples", "rust_integration")
	if _, err := os.Stat(rustDir); err == nil {
		// Only run cargo clean if there's a target directory with content to clean
		targetDir := filepath.Join(rustDir, "target")
		if _, targetErr := os.Stat(targetDir); targetErr == nil {
			// Check if the target directory has any content
			entries, err := os.ReadDir(targetDir)
			if err != nil {
				panic("Failed to read target directory: " + err.Error())
			}

			// Only clean if there are files/directories to clean
			if len(entries) > 0 {
				cleanCmd := exec.Command("cargo", "clean")
				cleanCmd.Dir = rustDir
				if output, err := cleanCmd.CombinedOutput(); err != nil {
					panic("Failed to clean Rust artifacts: " + err.Error() + "\nOutput: " + string(output))
				}
			}
		}

		tempTargetDir := filepath.Join(os.TempDir(), "osprey_rust_target_"+strconv.Itoa(os.Getpid()))
		cmd = exec.Command("cargo", "build", "--target-dir", tempTargetDir, "-j", "1")
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
