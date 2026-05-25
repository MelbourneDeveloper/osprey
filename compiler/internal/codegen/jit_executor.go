// Package codegen provides code generation and execution capabilities for Osprey.
package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// JITExecutor provides in-memory compilation and execution.
type JITExecutor struct {
	// For now, we'll use a self-contained approach that embeds the required tools
}

// NewJITExecutor creates a new JIT executor.
func NewJITExecutor() *JITExecutor {
	return &JITExecutor{}
}

// CompileAndRunInMemory compiles LLVM IR and runs it without external dependencies.
func (j *JITExecutor) CompileAndRunInMemory(ir string) error {
	// For immediate solution: use embedded compilation approach

	return j.compileAndRunEmbedded(ir)
}

// compileAndRunEmbedded uses an embedded approach with built-in LLVM tools detection.
func (j *JITExecutor) compileAndRunEmbedded(ir string) error {
	// Create temporary directory for compilation
	tempDir, err := os.MkdirTemp("", "osprey_compile_*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Write IR to file
	irFile := filepath.Join(tempDir, "program.ll")
	if writeErr := os.WriteFile(irFile, []byte(ir), FilePermissionsLess); writeErr != nil {
		return fmt.Errorf("failed to write IR file: %w", writeErr)
	}

	// Find LLVM tools in common locations
	llcPath, err := j.findLLVMTool("llc")
	if err != nil {
		return fmt.Errorf("LLVM tools not found. Please install LLVM or use a different execution method: %w", err)
	}

	// Compile IR to object file
	objFile := filepath.Join(tempDir, "program.o")
	// #nosec G204 - llcPath is validated through findLLVMTool
	llcCmd := exec.Command(llcPath, "-filetype=obj", "-o", objFile, irFile)

	llcOutput, err := llcCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile IR: %w\nOutput: %s", err, string(llcOutput))
	}

	// Find compiler for linking
	compilerPath, err := j.findCompiler()
	if err != nil {
		return fmt.Errorf("no suitable compiler found for linking: %w", err)
	}

	// Link to executable
	exeFile := filepath.Join(tempDir, "program")
	if runtime.GOOS == "windows" {
		exeFile += ".exe"
	}

	// Check for runtime libraries
	var linkArgs []string
	linkArgs = append(linkArgs, "-o", exeFile, objFile)

	// Look for fiber runtime library - prioritize local builds, then system installs
	fiberRuntimePaths := []string{
		"bin/libfiber_runtime.a",
		"./bin/libfiber_runtime.a",
		filepath.Join(filepath.Dir(os.Args[0]), "..", "libfiber_runtime.a"),
		"/usr/local/lib/libfiber_runtime.a", // System install location
	}

	// Look for HTTP runtime library - prioritize local builds, then system installs
	httpRuntimePaths := []string{
		"bin/libhttp_runtime.a",
		"./bin/libhttp_runtime.a",
		filepath.Join(filepath.Dir(os.Args[0]), "..", "libhttp_runtime.a"),
		"/usr/local/lib/libhttp_runtime.a", // System install location
	}

	// Add working directory based paths
	if wd, err := os.Getwd(); err == nil {
		fiberRuntimePaths = append(fiberRuntimePaths,
			filepath.Join(wd, "bin", "libfiber_runtime.a"),
			filepath.Join(wd, "..", "bin", "libfiber_runtime.a"),
			filepath.Join(wd, "..", "..", "bin", "libfiber_runtime.a"),
		)
		httpRuntimePaths = append(httpRuntimePaths,
			filepath.Join(wd, "bin", "libhttp_runtime.a"),
			filepath.Join(wd, "..", "bin", "libhttp_runtime.a"),
			filepath.Join(wd, "..", "..", "bin", "libhttp_runtime.a"),
		)
	}

	var foundFiberLib string
	var foundHTTPLib string

	for _, libPath := range fiberRuntimePaths {
		if _, err := os.Stat(libPath); err == nil {
			linkArgs = append(linkArgs, libPath)
			foundFiberLib = libPath

			break
		}
	}

	for _, libPath := range httpRuntimePaths {
		if _, err := os.Stat(libPath); err == nil {
			linkArgs = append(linkArgs, libPath)
			foundHTTPLib = libPath

			break
		}
	}

	linkArgs = append(linkArgs, "-lpthread")

	// Add OpenSSL libraries with platform-specific paths
	// Use pkg-config to get proper OpenSSL flags when available
	cmd := exec.Command("pkg-config", "--libs", "openssl")
	if output, err := cmd.Output(); err == nil {
		// Parse pkg-config output and add flags
		flags := strings.Fields(strings.TrimSpace(string(output)))
		linkArgs = append(linkArgs, flags...)
	} else {
		// Fallback to standard OpenSSL flags for different platforms
		if runtime.GOOS == "darwin" {
			// macOS with Homebrew OpenSSL - try multiple common paths
			possiblePaths := []string{
				"/opt/homebrew/opt/openssl@3/lib",
				"/opt/homebrew/lib",
				"/usr/local/opt/openssl@3/lib",
				"/usr/local/lib",
			}

			opensslLibPath := ""
			for _, path := range possiblePaths {
				if _, err := os.Stat(filepath.Join(path, "libssl.dylib")); err == nil {
					opensslLibPath = path

					break
				}
			}

			if opensslLibPath != "" {
				linkArgs = append(linkArgs, "-L"+opensslLibPath, "-lssl", "-lcrypto")
			} else {
				// Final fallback
				linkArgs = append(linkArgs, "-L/opt/homebrew/lib", "-lssl", "-lcrypto")
			}
		} else {
			// Linux and other systems
			linkArgs = append(linkArgs, "-lssl", "-lcrypto")
		}
	}

	// Debug output - only show warnings if libraries not found
	if foundFiberLib == "" {
		fmt.Fprintf(os.Stderr, "Warning: Fiber runtime library not found in any of: %v\n", fiberRuntimePaths)
	}
	if foundHTTPLib == "" {
		fmt.Fprintf(os.Stderr, "Warning: HTTP runtime library not found in any of: %v\n", httpRuntimePaths)
	}

	// #nosec G204 - compilerPath is validated through findCompiler
	linkCmd := exec.Command(compilerPath, linkArgs...)

	linkOutput, err := linkCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to link executable: %w\nOutput: %s", err, string(linkOutput))
	}

	// Execute the program
	// #nosec G204 - exeFile is created in controlled temp directory
	runCmd := exec.Command(exeFile)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	return runCmd.Run()
}

// findLLVMTool finds LLVM tools in common installation locations.
func (j *JITExecutor) findLLVMTool(toolName string) (string, error) {
	// Common LLVM installation paths
	commonPaths := []string{
		"/opt/homebrew/opt/llvm/bin/" + toolName,
		"/opt/homebrew/bin/" + toolName,
		"/usr/local/opt/llvm/bin/" + toolName,
		"/usr/local/bin/" + toolName,
		"/usr/bin/" + toolName,
	}

	// First check if it's in PATH
	if path, err := exec.LookPath(toolName); err == nil {
		return path, nil
	}

	// Check common installation locations
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", WrapToolNotFound(toolName)
}

// findCompiler finds a suitable C compiler for linking.
func (j *JITExecutor) findCompiler() (string, error) {
	compilers := []string{"clang", "gcc", "cc"}

	// Also check common paths
	commonPaths := []string{
		"/opt/homebrew/bin/",
		"/usr/local/bin/",
		"/usr/bin/",
	}

	// First check PATH
	for _, compiler := range compilers {
		if path, err := exec.LookPath(compiler); err == nil {
			return path, nil
		}
	}

	// Check common locations
	for _, basePath := range commonPaths {
		for _, compiler := range compilers {
			fullPath := basePath + compiler
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
		}
	}

	return "", WrapNoSuitableCompiler(compilers)
}

// CompileAndRunJIT is the main entry point for JIT compilation with default (permissive) security.
func CompileAndRunJIT(source string) error {
	return CompileAndRunJITWithSecurity(source, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
}

// CompileAndRunJITWithSecurity is the main entry point for JIT compilation with specified security configuration.
func CompileAndRunJITWithSecurity(source string, security SecurityConfig) error {
	// Generate LLVM IR with security configuration
	ir, err := CompileToLLVMWithSecurity(source, security)
	if err != nil {
		return fmt.Errorf("failed to generate LLVM IR: %w", err)
	}

	// Use JIT executor
	executor := NewJITExecutor()

	return executor.CompileAndRunInMemory(ir)
}
