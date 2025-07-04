package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/parser"
)

// CompileToLLVM compiles source code directly to LLVM IR string with default (permissive) security.
// This is a convenience function that encapsulates the entire compilation pipeline.
func CompileToLLVM(source string) (string, error) {
	return CompileToLLVMWithSecurity(source, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
}

// CompileToLLVMWithSecurity compiles source code to LLVM IR string with specified security configuration.
func CompileToLLVMWithSecurity(source string, security SecurityConfig) (string, error) {
	// Parse the source
	input := antlr.NewInputStream(source)
	lexer := parser.NewospreyLexer(input)

	// Create error listener to collect parse errors
	errorListener := &ParseErrorListener{}
	errorListener.Errors = make([]string, 0)

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewospreyParser(stream)

	// Add error listener to parser as well
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)

	tree := p.Program()

	// Check for parse errors before proceeding
	if len(errorListener.Errors) > 0 {
		return "", WrapParseErrors(strings.Join(errorListener.Errors, "\n"))
	}

	// Check if parse tree is valid
	if tree == nil {
		return "", ErrParseTreeNil
	}

	// Build AST
	builder := ast.NewBuilder()
	program := builder.BuildProgram(tree)

	// Check if AST building succeeded
	if program == nil {
		return "", ErrASTBuildFailed
	}

	// Validate AST for type inference rules
	if err := ast.ValidateProgram(program); err != nil {
		return "", fmt.Errorf("validation error: %w", err)
	}

	// Generate LLVM IR with security configuration
	generator := NewLLVMGeneratorWithSecurity(security)

	_, err := generator.GenerateProgram(program)
	if err != nil {
		return "", err
	}

	return generator.GenerateIR(), nil
}

// ParseErrorListener collects parse errors instead of panicking.
type ParseErrorListener struct {
	antlr.DefaultErrorListener
	Errors []string
}

// SyntaxError handles syntax errors during parsing.
func (p *ParseErrorListener) SyntaxError(
	_ antlr.Recognizer,
	_ interface{},
	line, column int,
	msg string,
	_ antlr.RecognitionException,
) {
	errorMsg := fmt.Sprintf("line %d:%d %s", line, column, msg)
	p.Errors = append(p.Errors, errorMsg)
}

// CompileToExecutable compiles source code to an executable binary with default (permissive) security.
func CompileToExecutable(source, outputPath string) error {
	return CompileToExecutableWithSecurity(source, outputPath, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
}

// buildLibraryPaths builds the search paths for runtime libraries
func buildLibraryPaths(libName string) []string {
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

// findAndAddLibrary finds a library in the given paths and adds it to linkArgs
func findAndAddLibrary(libName string, linkArgs []string) []string {
	paths := buildLibraryPaths(libName)
	for _, libPath := range paths {
		if _, err := os.Stat(libPath); err == nil {
			return append(linkArgs, libPath)
		}
	}
	return linkArgs
}

// addOpenSSLFlags adds OpenSSL linking flags using pkg-config or platform-specific fallbacks
func addOpenSSLFlags(linkArgs []string) []string {
	// Use pkg-config to get proper OpenSSL flags when available
	cmd := exec.Command("pkg-config", "--libs", "openssl")
	if output, err := cmd.Output(); err == nil {
		// Parse pkg-config output and add flags
		flags := strings.Fields(strings.TrimSpace(string(output)))
		fmt.Fprintf(os.Stderr, "DEBUG: pkg-config openssl flags: %v\n", flags)
		return append(linkArgs, flags...)
	}

	// Fallback to standard OpenSSL flags for different platforms
	if runtime.GOOS == "darwin" {
		// macOS with Homebrew OpenSSL
		return append(linkArgs, "-L/opt/homebrew/lib", "-lssl", "-lcrypto")
	}
	// Linux and other systems
	return append(linkArgs, "-lssl", "-lcrypto")
}

// checkLibraryAvailability checks if any runtime libraries are available
func checkLibraryAvailability() (bool, bool) {
	fiberExists := false
	httpExists := false

	// Check if any fiber runtime library was found
	for _, libPath := range buildLibraryPaths("fiber_runtime") {
		if _, err := os.Stat(libPath); err == nil {
			fiberExists = true
			break
		}
	}

	// Check if any HTTP runtime library was found
	for _, libPath := range buildLibraryPaths("http_runtime") {
		if _, err := os.Stat(libPath); err == nil {
			httpExists = true
			break
		}
	}

	return fiberExists, httpExists
}

// tryLinkWithCompilers attempts to link the executable using multiple compiler options
func tryLinkWithCompilers(outputPath, objFile string, linkArgs []string, fiberExists, httpExists bool) error {
	var clangCommands [][]string

	if fiberExists || httpExists {
		clangCommands = [][]string{
			append([]string{"clang"}, linkArgs...),                   // System clang
			append([]string{"/usr/bin/clang"}, linkArgs...),          // System path clang
			append([]string{"/opt/homebrew/bin/clang"}, linkArgs...), // Homebrew clang
			append([]string{"gcc"}, linkArgs...),                     // Fallback to gcc
		}
	} else {
		clangCommands = [][]string{
			{"clang", "-o", outputPath, objFile},                   // System clang
			{"/usr/bin/clang", "-o", outputPath, objFile},          // System path clang
			{"/opt/homebrew/bin/clang", "-o", outputPath, objFile}, // Homebrew clang
			{"gcc", "-o", outputPath, objFile},                     // Fallback to gcc
		}
	}

	var lastErr error
	for _, cmd := range clangCommands {
		fmt.Fprintf(os.Stderr, "DEBUG: Trying link command: %v\n", cmd)
		linkCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204 - predefined safe commands

		linkOutput, err := linkCmd.CombinedOutput()
		if err == nil {
			return nil // Success!
		}

		lastErr = fmt.Errorf("failed to link executable with %s: %w\nOutput: %s", cmd[0], err, string(linkOutput))
	}

	return fmt.Errorf("failed to link executable with any available compiler: %w", lastErr)
}

// CompileToExecutableWithSecurity compiles source code to an executable binary with specified security configuration.
func CompileToExecutableWithSecurity(source, outputPath string, security SecurityConfig) error {
	// Ensure the output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, DirPermissions); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate LLVM IR with security configuration
	ir, err := CompileToLLVMWithSecurity(source, security)
	if err != nil {
		return fmt.Errorf("failed to generate LLVM IR: %w", err)
	}

	// Write IR to temporary file
	irFile := outputPath + ".ll"
	if err := os.WriteFile(irFile, []byte(ir), FilePermissions); err != nil {
		return WrapWriteIRFile(err)
	}
	defer func() { _ = os.Remove(irFile) }() // Clean up temp file

	// Compile IR to object file using llc
	objFile := outputPath + ".o"
	llcCmd := exec.Command("llc", "-filetype=obj", "-o", objFile, irFile) // #nosec G204 - args are controlled

	llcOutput, err := llcCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile IR to object file: %w\nllc output: %s", err, string(llcOutput))
	}

	defer func() { _ = os.Remove(objFile) }() // Clean up temp file

	// Build link arguments with runtime libraries
	var linkArgs []string
	linkArgs = append(linkArgs, "-o", outputPath, objFile)

	// Find and add runtime libraries (order matters: dependents before dependencies)
	linkArgs = findAndAddLibrary("http_runtime", linkArgs)
	linkArgs = findAndAddLibrary("fiber_runtime", linkArgs)

	linkArgs = append(linkArgs, "-lpthread")

	// Add OpenSSL libraries with platform-specific paths
	linkArgs = addOpenSSLFlags(linkArgs)

	// Check library availability and try linking
	fiberExists, httpExists := checkLibraryAvailability()
	return tryLinkWithCompilers(outputPath, objFile, linkArgs, fiberExists, httpExists)
}

// CompileAndRun compiles and runs source code using smart JIT execution with default (permissive) security.
func CompileAndRun(source string) error {
	return CompileAndRunWithSecurity(source, SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
}

// CompileAndRunWithSecurity compiles and runs source code with specified security configuration.
func CompileAndRunWithSecurity(source string, security SecurityConfig) error {
	// Try JIT execution first (smart tool detection)
	return CompileAndRunJITWithSecurity(source, security)
}
