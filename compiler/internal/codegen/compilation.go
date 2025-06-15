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

	objFile, err := compileIRToObject(ir, outputPath)
	if err != nil {
		return err
	}
	defer func() { _ = os.Remove(objFile) }()

	return linkObjectToExecutable(objFile, outputPath)
}

func compileIRToObject(ir, outputPath string) (string, error) {
	irFile := outputPath + ".ll"
	if err := os.WriteFile(irFile, []byte(ir), FilePermissions); err != nil {
		return "", WrapWriteIRFile(err)
	}
	defer func() { _ = os.Remove(irFile) }()

	objFile := outputPath + ".o"
	llcCmd := exec.Command("llc", "-filetype=obj", "-o", objFile, irFile) // #nosec G204 - args are controlled

	llcOutput, err := llcCmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to compile IR to object file: %w\nllc output: %s", err, string(llcOutput))
	}

	return objFile, nil
}

func linkObjectToExecutable(objFile, outputPath string) error {
	linkArgs := buildLinkArguments(objFile, outputPath)
	clangCommands := buildClangCommands(linkArgs, outputPath, objFile)

	var lastErr error
	for _, cmd := range clangCommands {
		fmt.Fprintf(os.Stderr, "DEBUG: Trying link command: %v\n", cmd)
		linkCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204 - predefined safe commands

		linkOutput, err := linkCmd.CombinedOutput()
		if err == nil {
			return nil
		}

		lastErr = fmt.Errorf("failed to link executable with %s: %w\nOutput: %s", cmd[0], err, string(linkOutput))
	}

	return fmt.Errorf("failed to link executable with any available compiler: %w", lastErr)
}

func buildLinkArguments(objFile, outputPath string) []string {
	linkArgs := []string{"-o", outputPath, objFile}

	// Add runtime libraries if they exist
	fiberRuntimeLib := "bin/libfiber_runtime.a"
	httpRuntimeLib := "bin/libhttp_runtime.a"

	if _, err := os.Stat(fiberRuntimeLib); err == nil {
		linkArgs = append(linkArgs, fiberRuntimeLib)
	}
	if _, err := os.Stat(httpRuntimeLib); err == nil {
		linkArgs = append(linkArgs, httpRuntimeLib)
	}

	linkArgs = append(linkArgs, "-lpthread")
	linkArgs = append(linkArgs, getOpenSSLFlags()...)

	return linkArgs
}

func getOpenSSLFlags() []string {
	cmd := exec.Command("pkg-config", "--libs", "openssl")
	if output, err := cmd.Output(); err == nil {
		flags := strings.Fields(strings.TrimSpace(string(output)))
		fmt.Fprintf(os.Stderr, "DEBUG: pkg-config openssl flags: %v\n", flags)
		return flags
	}

	fmt.Fprintf(os.Stderr, "DEBUG: pkg-config failed, using fallback\n")
	if runtime.GOOS == "darwin" {
		return []string{"-L/opt/homebrew/lib", "-lssl", "-lcrypto"}
	}
	return []string{"-lssl", "-lcrypto"}
}

func buildClangCommands(linkArgs []string, outputPath, objFile string) [][]string {
	fiberExists := fileExists("bin/libfiber_runtime.a")
	httpExists := fileExists("bin/libhttp_runtime.a")

	if fiberExists || httpExists {
		return [][]string{
			append([]string{"clang"}, linkArgs...),
			append([]string{"/usr/bin/clang"}, linkArgs...),
			append([]string{"/opt/homebrew/bin/clang"}, linkArgs...),
			append([]string{"gcc"}, linkArgs...),
		}
	}

	return [][]string{
		{"clang", "-o", outputPath, objFile},
		{"/usr/bin/clang", "-o", outputPath, objFile},
		{"/opt/homebrew/bin/clang", "-o", outputPath, objFile},
		{"gcc", "-o", outputPath, objFile},
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
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
