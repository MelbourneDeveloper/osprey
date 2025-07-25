// Package main provides the command-line interface for the Osprey compiler.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/christianfindlay/osprey/internal/cli"
)

const (
	// MinArgs is the minimum number of arguments required
	MinArgs = 2
	// MinHoverArgs is the minimum number of arguments required for hover
	MinHoverArgs = 3
	// DocsFlag represents the --docs command line flag
	DocsFlag = "--docs"
	// DocsDirFlag represents the --docs-dir command line flag
	DocsDirFlag = "--docs-dir"
	// HoverFlag represents the --hover command line flag
	HoverFlag = "--hover"
)

// ErrUnknownOption is returned when an unknown command line option is encountered
var ErrUnknownOption = errors.New("unknown option")

// RunCLI is the main CLI function that handles all argument parsing and execution
func RunCLI(args []string) cli.CommandResult {
	// Handle insufficient arguments
	if len(args) < MinArgs {
		ShowHelp()
		return cli.CommandResult{Success: true, Output: ""}
	}

	// Handle help and version flags
	if result := handleBasicFlags(args); result != nil {
		return *result
	}

	// Handle special modes (docs, hover)
	if result := handleSpecialModes(args); result != nil {
		return *result
	}

	// Handle regular file-based operations
	return handleFileBasedOperations(args)
}

// handleBasicFlags processes help and version flags
func handleBasicFlags(args []string) *cli.CommandResult {
	switch args[1] {
	case "--help", "-h":
		ShowHelp()
		return &cli.CommandResult{Success: true, Output: ""}
	case "--version":
		fmt.Println("Osprey Compiler 1.0.0")
		return &cli.CommandResult{Success: true, Output: ""}
	}
	return nil
}

// handleSpecialModes processes docs and hover modes
func handleSpecialModes(args []string) *cli.CommandResult {
	switch args[1] {
	case DocsFlag:
		return handleDocsMode(args)
	case HoverFlag:
		return handleHoverMode(args)
	}
	return nil
}

// handleDocsMode processes the --docs flag
func handleDocsMode(args []string) *cli.CommandResult {
	// Find --docs-dir argument
	docsDir := ""
	for i := 1; i < len(args); i++ {
		if args[i] == DocsDirFlag && i+1 < len(args) {
			docsDir = args[i+1]
			break
		}
	}

	if docsDir == "" {
		return &cli.CommandResult{
			Success:  false,
			ErrorMsg: "--docs requires --docs-dir <directory> to specify output location",
		}
	}

	result := cli.RunCommand("", cli.OutputModeDocs, docsDir, false, cli.NewDefaultSecurityConfig())
	return &result
}

// handleHoverMode processes the --hover flag
func handleHoverMode(args []string) *cli.CommandResult {
	if len(args) < MinHoverArgs {
		fmt.Println("Error: --hover requires an element name")
		fmt.Println("Example: osprey --hover print")
		return &cli.CommandResult{Success: false, ErrorMsg: "Missing element name for --hover"}
	}
	result := cli.RunCommand(args[2], cli.OutputModeHover, "", false, cli.NewDefaultSecurityConfig())
	return &result
}

// handleFileBasedOperations processes regular file operations
func handleFileBasedOperations(args []string) cli.CommandResult {
	// Regular file-based operations need at least 2 args
	if len(args) < MinArgs {
		ShowHelp()
		return cli.CommandResult{Success: true, Output: ""}
	}

	filename := args[1]
	outputMode := cli.OutputModeLLVM // default to LLVM IR
	docsDir := ""

	// Create security config with defaults
	security := cli.NewDefaultSecurityConfig()

	// Parse arguments and handle them
	parsedMode, parsedDocsDir, quiet, parseErr := parseArgumentsForFile(args, security)
	if parseErr != nil {
		return cli.CommandResult{
			Success:  false,
			ErrorMsg: parseErr.Error(),
		}
	}

	if parsedMode != "" {
		outputMode = parsedMode
	}
	if parsedDocsDir != "" {
		docsDir = parsedDocsDir
	}

	// Execute the command with appropriate security settings
	return executeCommand(filename, outputMode, docsDir, quiet, security)
}

// parseArgumentsForFile parses command line arguments for file-based operations
func parseArgumentsForFile(args []string, security *cli.SecurityConfig) (string, string, bool, error) {
	outputMode := ""
	docsDir := ""
	quiet := false

	// Parse remaining arguments
	for i := MinArgs; i < len(args); i++ {
		arg := args[i]

		if mode := parseOutputMode(arg); mode != "" {
			outputMode = mode
		} else if arg == DocsDirFlag && i+1 < len(args) {
			docsDir = args[i+1]
			i++ // Skip next argument since we consumed it
		} else if arg == "--quiet" {
			quiet = true
		} else if !parseSecurityMode(arg, security) {
			return "", "", false, fmt.Errorf("%w: %s", ErrUnknownOption, arg)
		}
	}

	return outputMode, docsDir, quiet, nil
}

// parseOutputMode returns the output mode for a given argument
func parseOutputMode(arg string) string {
	modes := map[string]string{
		"--ast":     cli.OutputModeAST,
		"--llvm":    cli.OutputModeLLVM,
		"--compile": cli.OutputModeCompile,
		"--run":     cli.OutputModeRun,
		"--symbols": cli.OutputModeSymbols,
		DocsFlag:    cli.OutputModeDocs,
		HoverFlag:   cli.OutputModeHover,
	}
	return modes[arg]
}

// parseSecurityMode handles security-related arguments
func parseSecurityMode(arg string, security *cli.SecurityConfig) bool {
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

// executeCommand executes the command with the given security settings
func executeCommand(
	filename, outputMode, docsDir string,
	quiet bool,
	security *cli.SecurityConfig,
) cli.CommandResult {
	// Just call the merged function - it handles all security modes
	return cli.RunCommand(filename, outputMode, docsDir, quiet, security)
}

// ShowHelp displays the help message for the Osprey compiler
func ShowHelp() {
	fmt.Println("Osprey Compiler")
	fmt.Println()
	fmt.Println("Usage: osprey <source-file> [options]")
	fmt.Println("       osprey --docs --docs-dir <directory>")
	fmt.Println("       osprey --hover <element-name>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --ast      Show the Abstract Syntax Tree")
	fmt.Println("  --llvm     Show LLVM IR (default)")
	fmt.Println("  --compile  Compile to executable")
	fmt.Println("  --run      Compile and run immediately")
	fmt.Println("  --symbols  Output symbol information as JSON")
	fmt.Println("  --docs     Generate API reference documentation (no file required)")
	fmt.Println("  --docs-dir <directory> Output directory for documentation (REQUIRED with --docs)")
	fmt.Println("  --hover    Get hover documentation for language element")
	fmt.Println("  --quiet    Suppress compiler messages (errors still shown)")
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
func ParseArgs(args []string) (string, string, string, *cli.SecurityConfig) {
	if len(args) < MinArgs {
		ShowHelp()
		return "", "", "", nil
	}

	// Handle help flags
	if args[1] == "--help" || args[1] == "-h" {
		ShowHelp()
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
	if args[1] == DocsFlag {
		var docsDir string
		// Check for --docs-dir argument (REQUIRED)
		for i := MinArgs; i < len(args); i++ {
			if args[i] == DocsDirFlag && i+1 < len(args) {
				docsDir = args[i+1]
				break
			}
		}
		return "", cli.OutputModeDocs, docsDir
	}

	// Handle hover flag (element name required)
	if args[1] == HoverFlag {
		if len(args) < MinHoverArgs {
			fmt.Println("Error: --hover requires an element name")
			fmt.Println("Example: osprey --hover print")
			return "", "", ""
		}
		return args[2], cli.OutputModeHover, ""
	}

	return "", "", ""
}

// ParseFileBasedArgs parses arguments for file-based operations
func ParseFileBasedArgs(args []string) (string, string, *cli.SecurityConfig) {
	outputMode := cli.OutputModeLLVM // default to LLVM IR
	docsDir := ""

	// Create security config with defaults
	security := cli.NewDefaultSecurityConfig()

	// Parse remaining arguments
	for i := MinArgs; i < len(args); i++ {
		arg := args[i]

		if ParseOutputModeArg(arg) != "" {
			outputMode = ParseOutputModeArg(arg)
		} else if arg == DocsDirFlag && i+1 < len(args) {
			docsDir = args[i+1]
			i++ // Skip next argument since we consumed it
		} else if !ParseSecurityArg(arg, security) {
			fmt.Printf("Unknown option: %s\n", arg)
			return "", "", nil
		}
	}

	return outputMode, docsDir, security
}

// ParseOutputModeArg parses output mode arguments
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
	case DocsFlag:
		return cli.OutputModeDocs
	case HoverFlag:
		return cli.OutputModeHover
	default:
		return ""
	}
}

// ParseSecurityArg parses security-related arguments
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
	result := RunMain(os.Args)
	if !result.Success {
		fmt.Fprintf(os.Stderr, "%s\n", result.ErrorMsg)
		os.Exit(1)
	}
	if result.Output != "" {
		fmt.Print(result.Output)
	}
	if result.OutputFile != "" {
		fmt.Printf("Output written to: %s\n", result.OutputFile)
	}
}

// RunMain is the testable main function that takes args as parameter
func RunMain(args []string) cli.CommandResult {
	return RunCLI(args)
}
