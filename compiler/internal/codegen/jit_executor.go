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

// setupCompilation creates temp directory and writes IR file
func (j *JITExecutor) setupCompilation(ir string) (string, string, string, error) {
	// Create temporary directory for compilation
	tempDir, err := os.MkdirTemp("", "osprey_compile_*")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Write IR to file
	irFile := filepath.Join(tempDir, "program.ll")
	if writeErr := os.WriteFile(irFile, []byte(ir), FilePermissionsLess); writeErr != nil {
		return "", "", "", fmt.Errorf("failed to write IR file: %w", writeErr)
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
		return "", fmt.Errorf("failed to compile IR: %w\nOutput: %s", err, string(llcOutput))
	}

	return objFile, nil
}

// setupLinkArgs builds the linking arguments for the executable
func (j *JITExecutor) setupLinkArgs(exeFile, objFile string) []string {
	var linkArgs []string
	linkArgs = append(linkArgs, "-o", exeFile, objFile)

	// Find and add runtime libraries (order matters: dependents before dependencies)
	linkArgs = j.findAndAddRuntimeLibrary("http_runtime", linkArgs)
	linkArgs = j.findAndAddRuntimeLibrary("fiber_runtime", linkArgs)

	linkArgs = append(linkArgs, "-lpthread")

	// Add OpenSSL libraries
	linkArgs = j.addOpenSSLFlags(linkArgs)

	return linkArgs
}

// findAndAddRuntimeLibrary finds a runtime library and adds it to link args
func (j *JITExecutor) findAndAddRuntimeLibrary(libName string, linkArgs []string) []string {
	paths := j.buildRuntimeLibraryPaths(libName)

	var foundLib string
	for _, libPath := range paths {
		if _, err := os.Stat(libPath); err == nil {
			linkArgs = append(linkArgs, libPath)
			foundLib = libPath
			break
		}
	}

	// Debug output - only show warnings if libraries not found
	if foundLib == "" {
		fmt.Fprintf(os.Stderr, "Warning: %s runtime library not found in any of: %v\n", libName, paths)
	}

	return linkArgs
}

// buildRuntimeLibraryPaths builds search paths for a specific runtime library
func (j *JITExecutor) buildRuntimeLibraryPaths(libName string) []string {
	paths := []string{
		fmt.Sprintf("bin/lib%s.a", libName),
		fmt.Sprintf("./bin/lib%s.a", libName),
		fmt.Sprintf("../../bin/lib%s.a", libName),    // For tests running from tests/integration
		fmt.Sprintf("../../../bin/lib%s.a", libName), // For deeper test directories
		filepath.Join(filepath.Dir(os.Args[0]), "..", fmt.Sprintf("lib%s.a", libName)),
		fmt.Sprintf("/usr/local/lib/lib%s.a", libName), // System install location
	}

	// Add working directory based paths
	if wd, err := os.Getwd(); err == nil {
		paths = append(paths,
			filepath.Join(wd, "bin", fmt.Sprintf("lib%s.a", libName)),
			filepath.Join(wd, "..", "bin", fmt.Sprintf("lib%s.a", libName)),
			filepath.Join(wd, "..", "..", "bin", fmt.Sprintf("lib%s.a", libName)),
			filepath.Join(wd, "..", "..", "..", "bin", fmt.Sprintf("lib%s.a", libName)), // For test directories
		)
	}

	return paths
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
		return fmt.Errorf("failed to link executable: %w\nOutput: %s", err, string(linkOutput))
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
	linkArgs := j.setupLinkArgs(exeFile, objFile)

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
