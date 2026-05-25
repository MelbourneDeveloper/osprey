// Package main provides the command-line interface for the Osprey compiler.
package main

import (
	"fmt"
	"os"

	"github.com/christianfindlay/osprey/internal/cli"
)

// RunCLI is the main CLI function that handles all argument parsing and execution
func RunCLI(args []string) cli.CommandResult {
	// Handle insufficient arguments
	if len(args) < 2 {
		ShowHelp()
		return cli.CommandResult{Success: true, Output: ""}
	}

	// Handle help flags
	if args[1] == "--help" || args[1] == "-h" {
		ShowHelp()
		return cli.CommandResult{Success: true, Output: ""}
	}

	// Handle version flag
	if args[1] == "--version" {
		fmt.Println("Osprey Compiler 1.0.0")
		return cli.CommandResult{Success: true, Output: ""}
	}

	// Handle docs flag (no file required)
	if args[1] == "--docs" {
		docsDir := "../website/src/docs" // default directory
		// Check for --docs-dir argument
		for i := 2; i < len(args); i++ {
			if args[i] == "--docs-dir" && i+1 < len(args) {
				docsDir = args[i+1]
				break
			}
		}
		return cli.RunCommand("", cli.OutputModeDocs, docsDir)
	}

	// Handle hover flag (element name required)
	if args[1] == "--hover" {
		if len(args) < 3 {
			fmt.Println("Error: --hover requires an element name")
			fmt.Println("Example: osprey --hover print")
			return cli.CommandResult{Success: false, ErrorMsg: "Missing element name for --hover"}
		}
		return cli.RunCommand(args[2], cli.OutputModeHover, "")
	}

	// Regular file-based operations need at least 2 args
	if len(args) < 2 {
		ShowHelp()
		return cli.CommandResult{Success: true, Output: ""}
	}

	filename := args[1]
	outputMode := cli.OutputModeLLVM // default to LLVM IR
	docsDir := ""

	// Create security config with defaults
	security := cli.NewDefaultSecurityConfig()

	// Parse remaining arguments
	for i := 2; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "--ast":
			outputMode = cli.OutputModeAST
		case "--llvm":
			outputMode = cli.OutputModeLLVM
		case "--compile":
			outputMode = cli.OutputModeCompile
		case "--run":
			outputMode = cli.OutputModeRun
		case "--symbols":
			outputMode = cli.OutputModeSymbols
		case "--docs":
			outputMode = cli.OutputModeDocs
		case "--hover":
			outputMode = cli.OutputModeHover
		case "--docs-dir":
			if i+1 < len(args) {
				docsDir = args[i+1]
				i++ // Skip next argument since we consumed it
			}
		case "--sandbox":
			security.ApplySandboxMode()
		case "--no-http":
			security.AllowHTTP = false
		case "--no-websocket":
			security.AllowWebSocket = false
		case "--no-fs":
			security.AllowFileRead = false
			security.AllowFileWrite = false
		case "--no-ffi":
			security.AllowFFI = false
		default:
			return cli.CommandResult{
				Success:  false,
				ErrorMsg: fmt.Sprintf("Unknown option: %s", arg),
			}
		}
	}

	// Execute the command with appropriate security settings
	var result cli.CommandResult

	// Use security-aware functions if security settings are non-default
	if security.SandboxMode || !security.AllowHTTP || !security.AllowWebSocket ||
		!security.AllowFileRead || !security.AllowFileWrite || !security.AllowFFI {

		// Use security-aware command execution
		result = cli.RunCommandWithSecurity(filename, outputMode, security)
	} else {
		// Use regular command execution for default/permissive mode
		result = cli.RunCommand(filename, outputMode, docsDir)
	}

	return result
}

const (
	minArgs      = 2
	minHoverArgs = 3
)

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
	fmt.Println("  --version  Show version information")
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

func ParseArgs(args []string) (string, string, string, *cli.SecurityConfig) {
	if len(args) < minArgs {
		ShowHelp()
		return "", "", "", nil
	}

	// Handle help flags
	if args[1] == "--help" || args[1] == "-h" {
		ShowHelp()
		return "", "", "", nil
	}

	// Handle version flag
	if os.Args[1] == "--version" {
		fmt.Println("Osprey Compiler v1.0.0")
		os.Exit(0)
		return "", "", "", nil
	}

	// Handle special modes (docs, hover)
	if filename, outputMode, docsDir := HandleSpecialModes(args); filename != "" || outputMode != "" {
		return filename, outputMode, docsDir, nil
	}

	// Regular file-based operations need at least 2 args
	if len(args) < minArgs {
		ShowHelp()
		return "", "", "", nil
	}

	filename := args[1]
	outputMode, docsDir, security := ParseFileBasedArgs(args)

	return filename, outputMode, docsDir, security
}

func HandleSpecialModes(args []string) (string, string, string) {
	// Handle docs flag (no file required)
	if args[1] == "--docs" {
		docsDir := "../website/src/docs" // default directory
		// Check for --docs-dir argument
		for i := 2; i < len(args); i++ {
			if args[i] == "--docs-dir" && i+1 < len(args) {
				docsDir = args[i+1]
				break
			}
		}
		return "", cli.OutputModeDocs, docsDir
	}

	// Handle hover flag (element name required)
	if args[1] == "--hover" {
		if len(args) < minHoverArgs {
			fmt.Println("Error: --hover requires an element name")
			fmt.Println("Example: osprey --hover print")
			return "", "", ""
		}
		return args[2], cli.OutputModeHover, ""
	}

	return "", "", ""
}

func ParseFileBasedArgs(args []string) (string, string, *cli.SecurityConfig) {
	outputMode := cli.OutputModeLLVM // default to LLVM IR
	docsDir := ""

	// Create security config with defaults
	security := cli.NewDefaultSecurityConfig()

	// Parse remaining arguments
	for i := 2; i < len(args); i++ {
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

func ParseOutputModeArg(arg string) string {
	switch arg {
	case "--ast":
		return cli.OutputModeAST
	case "--llvm":
		return cli.OutputModeLLVM
	case "--compile":
		return cli.OutputModeCompile
	case "--run":
		return cli.OutputModeRun
	case "--symbols":
		return cli.OutputModeSymbols
	case "--docs":
		return cli.OutputModeDocs
	case "--hover":
		return cli.OutputModeHover
	default:
		return ""
	}
}

func ParseSecurityArg(arg string, security *cli.SecurityConfig) bool {
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

func main() {
	result := RunCLI(os.Args)

	if !result.Success {
		fmt.Println(result.ErrorMsg)
		os.Exit(1)
	}

	fmt.Print(result.Output)
	if result.OutputFile != "" {
		fmt.Printf("Successfully compiled to %s\n", result.OutputFile)
	}
}

// Testable main function that takes args as parameter
func RunMain(args []string) cli.CommandResult {
	filename, outputMode, docsDir, security := ParseArgs(args)
	if filename == "" && outputMode == "" {
		return cli.CommandResult{Success: true, Output: ""}
	}

	var result cli.CommandResult

	// Use security-aware functions if security settings are non-default
	if security != nil && (security.SandboxMode || !security.AllowHTTP || !security.AllowWebSocket ||
		!security.AllowFileRead || !security.AllowFileWrite || !security.AllowFFI) {

		// Use security-aware command execution
		result = cli.RunCommandWithSecurity(filename, outputMode, security)
	} else {
		// Use regular command execution for default/permissive mode
		result = cli.RunCommand(filename, outputMode, docsDir)
	}

	return result
}
