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
	libDir string // Optional lib directory override (for tests)
}

// NewJITExecutor creates a new JIT executor.
func NewJITExecutor() *JITExecutor {
	return &JITExecutor{}
}

// NewJITExecutorWithLibDir creates a new JIT executor instance with custom lib directory
func NewJITExecutorWithLibDir(libDir string) *JITExecutor {
	return &JITExecutor{libDir: libDir}
}

// CompileAndRunInMemory compiles LLVM IR and runs it without external dependencies.
func (j *JITExecutor) CompileAndRunInMemory(ir string) error {
	// For immediate solution: use embedded compilation approach

	return j.compileAndRunEmbedded(ir)
}

// CompileAndCaptureOutput compiles LLVM IR and captures the program's output
func (j *JITExecutor) CompileAndCaptureOutput(ir string) (string, error) {
	// Setup compilation environment
	tempDir, irFile, exeFile, err := j.setupCompilation(ir)
	if err != nil {
		return "", err
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Compile IR to object file
	objFile, err := j.compileToObject(irFile, tempDir)
	if err != nil {
		return "", err
	}

	// Setup linking arguments
	linkArgs, err := j.setupLinkArgs(exeFile, objFile)
	if err != nil {
		return "", err
	}

	// Link to executable
	if err := j.linkExecutable(linkArgs); err != nil {
		return "", err
	}

	// Execute and capture output
	return j.executeProgramWithCapture(exeFile)
}

// setupCompilation creates temp directory and writes IR file
func (j *JITExecutor) setupCompilation(ir string) (string, string, string, error) {
	// Create temporary directory for compilation
	tempDir, err := os.MkdirTemp("", "osprey_compile_*")
	if err != nil {
		return "", "", "", fmt.Errorf("INTERNAL_COMPILER_ERROR: failed to create temp directory: %w", err)
	}

	// Write IR to file
	irFile := filepath.Join(tempDir, "program.ll")
	if writeErr := os.WriteFile(irFile, []byte(ir), FilePermissionsLess); writeErr != nil {
		return "", "", "", fmt.Errorf("INTERNAL_COMPILER_ERROR: failed to write IR file: %w", writeErr)
	}

	// Determine executable file name
	exeFile := filepath.Join(tempDir, "program")
	if runtime.GOOS == "windows" {
		exeFile += ".exe"
	}

	return tempDir, irFile, exeFile, nil
}

// compileToObject compiles LLVM IR to object file
func (j *JITExecutor) compileToObject(irFile, tempDir string) (string, error) {
	// Find LLVM tools in common locations
	llcPath, err := j.findLLVMTool("llc")
	if err != nil {
		return "", fmt.Errorf("LLVM tools not found. Please install LLVM or use a different execution method: %w", err)
	}

	// Compile IR to object file
	objFile := filepath.Join(tempDir, "program.o")
	// #nosec G204 - llcPath is validated through findLLVMTool
	llcCmd := exec.Command(llcPath, "-filetype=obj", "-o", objFile, irFile)

	llcOutput, err := llcCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("INTERNAL_COMPILER_ERROR: failed to compile IR: %w\nOutput: %s", err, string(llcOutput))
	}

	return objFile, nil
}

// setupLinkArgs builds the linking arguments for the executable
func (j *JITExecutor) setupLinkArgs(exeFile, objFile string) ([]string, error) {
	var linkArgs []string
	linkArgs = append(linkArgs, "-o", exeFile, objFile)

	// FAIL HARD: All runtime libraries must be available for JIT execution
	for _, libName := range RuntimeLibraries {
		libPath, err := getLibraryPathWithDir(libName, j.libDir)
		if err != nil {
			return nil, WrapMissingRuntimeLibrary(libName)
		}

		if _, err := os.Stat(libPath); err != nil {
			return nil, WrapMissingRuntimeLibrary(libName)
		}

		linkArgs = append(linkArgs, libPath)
	}

	linkArgs = append(linkArgs, "-lpthread")

	// Add OpenSSL libraries
	linkArgs = j.addOpenSSLFlags(linkArgs)

	return linkArgs, nil
}

// addOpenSSLFlags adds OpenSSL linking flags
func (j *JITExecutor) addOpenSSLFlags(linkArgs []string) []string {
	// Use pkg-config to get proper OpenSSL flags when available
	cmd := exec.Command("pkg-config", "--libs", "openssl")
	if output, err := cmd.Output(); err == nil {
		// Parse pkg-config output and add flags
		flags := strings.Fields(strings.TrimSpace(string(output)))
		return append(linkArgs, flags...)
	}

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
			return append(linkArgs, "-L"+opensslLibPath, "-lssl", "-lcrypto")
		}
		// Final fallback
		return append(linkArgs, "-L/opt/homebrew/lib", "-lssl", "-lcrypto")
	}
	// Linux and other systems
	return append(linkArgs, "-lssl", "-lcrypto")
}

// linkExecutable links the object file into an executable
func (j *JITExecutor) linkExecutable(linkArgs []string) error {
	// Find compiler for linking
	compilerPath, err := j.findCompiler()
	if err != nil {
		return fmt.Errorf("no suitable compiler found for linking: %w", err)
	}

	// #nosec G204 - compilerPath is validated through findCompiler
	linkCmd := exec.Command(compilerPath, linkArgs...)

	linkOutput, err := linkCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("INTERNAL_COMPILER_ERROR: failed to link executable: %w\nOutput: %s", err, string(linkOutput))
	}

	return nil
}

// executeProgram runs the compiled executable
func (j *JITExecutor) executeProgram(exeFile string) error {
	// #nosec G204 - exeFile is created in controlled temp directory
	runCmd := exec.Command(exeFile)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	return runCmd.Run()
}

// executeProgramWithCapture runs the compiled executable and captures its output
func (j *JITExecutor) executeProgramWithCapture(exeFile string) (string, error) {
	// #nosec G204 - exeFile is created in controlled temp directory
	runCmd := exec.Command(exeFile)

	// CAPTURE STDOUT instead of outputting directly to terminal
	output, err := runCmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// compileAndRunEmbedded uses an embedded approach with built-in LLVM tools detection.
func (j *JITExecutor) compileAndRunEmbedded(ir string) error {
	// Setup compilation environment
	tempDir, irFile, exeFile, err := j.setupCompilation(ir)
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Compile IR to object file
	objFile, err := j.compileToObject(irFile, tempDir)
	if err != nil {
		return err
	}

	// Setup linking arguments
	linkArgs, err := j.setupLinkArgs(exeFile, objFile)
	if err != nil {
		return err
	}

	// Link to executable
	if err := j.linkExecutable(linkArgs); err != nil {
		return err
	}

	// Execute the program
	return j.executeProgram(exeFile)
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

// CompileAndRunJITWithLibDir is the main entry point for JIT compilation with custom lib directory.
func CompileAndRunJITWithLibDir(source, libDir string) error {
	return CompileAndRunJITWithSecurityAndLibDir(source, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	}, libDir)
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

// CompileAndRunJITWithSecurityAndLibDir is the main entry point for JIT compilation with security and lib directory.
func CompileAndRunJITWithSecurityAndLibDir(source string, security SecurityConfig, libDir string) error {
	// Generate LLVM IR with security configuration
	ir, err := CompileToLLVMWithSecurity(source, security)
	if err != nil {
		return fmt.Errorf("failed to generate LLVM IR: %w", err)
	}

	// Use JIT executor with custom lib directory
	executor := NewJITExecutorWithLibDir(libDir)
	return executor.CompileAndRunInMemory(ir)
}

// CompileAndCaptureJIT compiles and captures program output with default (permissive) security.
func CompileAndCaptureJIT(source string) (string, error) {
	return CompileAndCaptureJITWithSecurity(source, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
}

// CompileAndCaptureJITWithLibDir compiles and captures program output with custom lib directory.
func CompileAndCaptureJITWithLibDir(source, libDir string) (string, error) {
	return CompileAndCaptureJITWithSecurityAndLibDir(source, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	}, libDir)
}

// CompileAndCaptureJITWithSecurity compiles and captures program output with specified security configuration.
func CompileAndCaptureJITWithSecurity(source string, security SecurityConfig) (string, error) {
	// Generate LLVM IR with security configuration
	ir, err := CompileToLLVMWithSecurity(source, security)
	if err != nil {
		return "", fmt.Errorf("failed to generate LLVM IR: %w", err)
	}

	// Use JIT executor to compile and capture output
	executor := NewJITExecutor()
	return executor.CompileAndCaptureOutput(ir)
}

// CompileAndCaptureJITWithSecurityAndLibDir compiles and captures program output with security and lib directory.
func CompileAndCaptureJITWithSecurityAndLibDir(source string, security SecurityConfig, libDir string) (string, error) {
	// Generate LLVM IR with security configuration
	ir, err := CompileToLLVMWithSecurity(source, security)
	if err != nil {
		return "", fmt.Errorf("failed to generate LLVM IR: %w", err)
	}

	// Use JIT executor with custom lib directory to compile and capture output
	executor := NewJITExecutorWithLibDir(libDir)
	return executor.CompileAndCaptureOutput(ir)
}
