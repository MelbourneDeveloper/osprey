package cli

import (
	"fmt"
	"os"
)

const (
	// MinArgs is the minimum number of arguments required for most commands
	MinArgs = 2
	// MinHoverArgs is the minimum number of arguments required for hover command
	MinHoverArgs = 3
)

// ShowHelp displays the help message for the Osprey compiler
func ShowHelp() {
	fmt.Println("Osprey Compiler")
	fmt.Println()
	fmt.Println("Usage: osprey <source-file> [options]")
	fmt.Println("       osprey --docs [--docs-dir <directory>]")
	fmt.Println("       osprey --hover <element-name>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --ast      Show the Abstract Syntax Tree")
	fmt.Println("  --llvm     Show LLVM IR (default)")
	fmt.Println("  --compile  Compile to executable")
	fmt.Println("  --run      Compile and run immediately")
	fmt.Println("  --symbols  Output symbol information as JSON")
	fmt.Println("  --docs     Generate API reference documentation (no file required)")
	fmt.Println("  --docs-dir <directory> Output directory for documentation (used with --docs)")
	fmt.Println("  --hover    Get hover documentation for language element")
	fmt.Println("  --help, -h Show this help message")
	fmt.Println()
	fmt.Println("Security Options:")
	fmt.Println("  --sandbox      Enable sandbox mode (disable all risky operations)")
	fmt.Println("  --no-http      Disable HTTP functions")
	fmt.Println("  --no-websocket Disable WebSocket functions")
	fmt.Println("  --no-fs        Disable file system access")
	fmt.Println("  --no-ffi       Disable foreign function interface")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  osprey program.osp --run         # Compile and run")
	fmt.Println("  osprey program.osp --compile     # Compile to executable")
	fmt.Println("  osprey program.osp --llvm        # Show LLVM IR")
	fmt.Println("  osprey program.osp --ast         # Show AST")
	fmt.Println("  osprey --docs --docs-dir ./docs  # Generate docs to ./docs")
	fmt.Println("  osprey --hover print             # Get hover docs for print function")
	fmt.Println("  osprey program.osp --sandbox     # Compile with all security restrictions")
	fmt.Println("  osprey program.osp --no-http     # Compile without HTTP functions")
}

// ParseArgs parses command line arguments and returns parsed values
func ParseArgs(args []string) (string, string, string, *SecurityConfig) {
	if len(args) < MinArgs {
		ShowHelp()
		return "", "", "", nil
	}

	// Handle help flags
	if args[1] == "--help" || args[1] == "-h" {
		ShowHelp()
		return "", "", "", nil
	}

	// Handle version flag
	if args[1] == "--version" {
		fmt.Println("Osprey Compiler 1.0.0")
		return "", "", "", nil
	}

	// Handle special modes (docs, hover)
	if filename, outputMode, docsDir := HandleSpecialModes(args); filename != "" || outputMode != "" {
		return filename, outputMode, docsDir, nil
	}

	// Regular file-based operations need at least 2 args
	if len(args) < MinArgs {
		ShowHelp()
		return "", "", "", nil
	}

	filename := args[1]
	outputMode, docsDir, security := ParseFileBasedArgs(args)

	return filename, outputMode, docsDir, security
}

// HandleSpecialModes handles special command modes like docs and hover
func HandleSpecialModes(args []string) (string, string, string) {
	// Handle docs flag (no file required)
	if args[1] == "--docs" {
		var docsDir string
		// Check for --docs-dir argument (REQUIRED)
		for i := MinArgs; i < len(args); i++ {
			if args[i] == "--docs-dir" && i+1 < len(args) {
				docsDir = args[i+1]
				break
			}
		}
		return "", OutputModeDocs, docsDir
	}

	// Handle hover flag (element name required)
	if args[1] == "--hover" {
		if len(args) < MinHoverArgs {
			fmt.Println("Error: --hover requires an element name")
			fmt.Println("Example: osprey --hover print")
			return "", "", ""
		}
		return args[MinArgs], OutputModeHover, ""
	}

	return "", "", ""
}

// ParseFileBasedArgs parses arguments for file-based operations
func ParseFileBasedArgs(args []string) (string, string, *SecurityConfig) {
	outputMode := OutputModeLLVM // default to LLVM IR
	docsDir := ""

	// Create security config with defaults
	security := NewDefaultSecurityConfig()

	// Parse remaining arguments
	for i := MinArgs; i < len(args); i++ {
		arg := args[i]

		if newOutputMode := ParseOutputModeArg(arg); newOutputMode != "" {
			outputMode = newOutputMode
		} else if arg == "--docs-dir" && i+1 < len(args) {
			docsDir = args[i+1]
			i++ // Skip next argument since we consumed it
		} else if !ParseSecurityArg(arg, security) {
			fmt.Printf("Unknown option: %s\n", arg)
			return "", "", nil
		}
	}

	return outputMode, docsDir, security
}

// ParseOutputModeArg parses output mode arguments and returns the corresponding mode
func ParseOutputModeArg(arg string) string {
	switch arg {
	case "--ast":
		return OutputModeAST
	case "--llvm":
		return OutputModeLLVM
	case "--compile":
		return OutputModeCompile
	case "--run":
		return OutputModeRun
	case "--symbols":
		return OutputModeSymbols
	case "--docs":
		return OutputModeDocs
	case "--hover":
		return OutputModeHover
	default:
		return ""
	}
}

// ParseSecurityArg parses security-related arguments and applies them to the security config
func ParseSecurityArg(arg string, security *SecurityConfig) bool {
	switch arg {
	case "--sandbox":
		security.ApplySandboxMode()
		return true
	case "--no-http":
		security.AllowHTTP = false
		return true
	case "--no-websocket":
		security.AllowWebSocket = false
		return true
	case "--no-fs":
		security.AllowFileRead = false
		security.AllowFileWrite = false
		return true
	case "--no-ffi":
		security.AllowFFI = false
		return true
	default:
		return false
	}
}

// RunMainWithArgs is the testable main function that takes args as parameter
func RunMainWithArgs(args []string) CommandResult {
	filename, outputMode, docsDir, security := ParseArgs(args)
	if filename == "" && outputMode == "" {
		return CommandResult{Success: true, Output: ""}
	}

	var result CommandResult

	// Use security-aware functions if security settings are non-default
	if security != nil && (security.SandboxMode || !security.AllowHTTP || !security.AllowWebSocket ||
		!security.AllowFileRead || !security.AllowFileWrite || !security.AllowFFI) {

		// Use security-aware command execution
		result = RunCommandWithSecurity(filename, outputMode, security)
	} else {
		// Use regular command execution for default/permissive mode
		result = RunCommand(filename, outputMode, docsDir)
	}

	return result
}

// RunMainFromOS runs the main function using os.Args
func RunMainFromOS() CommandResult {
	return RunMainWithArgs(os.Args)
}
