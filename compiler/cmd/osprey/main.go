// Package main provides the command-line interface for the Osprey compiler.
package main

import (
	"fmt"
	"os"

	"github.com/christianfindlay/osprey/internal/cli"
)

const (
	minArgs      = 2
	minHoverArgs = 3
)

func showHelp() {
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

func parseArgs() (string, string, string, *cli.SecurityConfig) {
	if len(os.Args) < minArgs {
		showHelp()
		os.Exit(1)
	}

	// Handle help flags
	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		showHelp()
		return "", "", "", nil
	}

	// Handle special modes (docs, hover)
	if filename, outputMode, docsDir := handleSpecialModes(); filename != "" || outputMode != "" {
		return filename, outputMode, docsDir, nil
	}

	// Regular file-based operations need at least 2 args
	if len(os.Args) < minArgs {
		showHelp()
		os.Exit(1)
	}

	filename := os.Args[1]
	outputMode, docsDir, security := parseFileBasedArgs()

	return filename, outputMode, docsDir, security
}

func handleSpecialModes() (string, string, string) {
	// Handle docs flag (no file required)
	if os.Args[1] == "--docs" {
		docsDir := "../website/src/docs" // default directory
		// Check for --docs-dir argument
		for i := 2; i < len(os.Args); i++ {
			if os.Args[i] == "--docs-dir" && i+1 < len(os.Args) {
				docsDir = os.Args[i+1]
				break
			}
		}
		return "", cli.OutputModeDocs, docsDir
	}

	// Handle hover flag (element name required)
	if os.Args[1] == "--hover" {
		if len(os.Args) < minHoverArgs {
			fmt.Println("Error: --hover requires an element name")
			fmt.Println("Example: osprey --hover print")
			os.Exit(1)
		}
		return os.Args[2], cli.OutputModeHover, ""
	}

	return "", "", ""
}

func parseFileBasedArgs() (string, string, *cli.SecurityConfig) {
	outputMode := cli.OutputModeLLVM // default to LLVM IR
	docsDir := ""

	// Create security config with defaults
	security := cli.NewDefaultSecurityConfig()

	// Parse remaining arguments
	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]

		if newOutputMode := parseOutputModeArg(arg); newOutputMode != "" {
			outputMode = newOutputMode
		} else if arg == "--docs-dir" && i+1 < len(os.Args) {
			docsDir = os.Args[i+1]
			i++ // Skip next argument since we consumed it
		} else if !parseSecurityArg(arg, security) {
			fmt.Printf("Unknown option: %s\n", arg)
			os.Exit(1)
		}
	}

	return outputMode, docsDir, security
}

func parseOutputModeArg(arg string) string {
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

func parseSecurityArg(arg string, security *cli.SecurityConfig) bool {
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
	filename, outputMode, docsDir, security := parseArgs()
	if filename == "" && outputMode == "" {
		return
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

	if !result.Success {
		fmt.Println(result.ErrorMsg)
		os.Exit(1)
	}

	fmt.Print(result.Output)
	if result.OutputFile != "" {
		fmt.Printf("Successfully compiled to %s\n", result.OutputFile)
	}
}
