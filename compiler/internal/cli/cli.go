// Package cli provides testable command-line interface functionality for the Osprey compiler.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/codegen"
	"github.com/christianfindlay/osprey/internal/language/descriptions"
	"github.com/christianfindlay/osprey/parser"
)

const (
	// OutputModeLLVM represents LLVM IR output mode.
	OutputModeLLVM = "llvm"
	// OutputModeAST represents AST output mode.
	OutputModeAST = "ast"
	// OutputModeCompile represents compilation output mode.
	OutputModeCompile = "compile"
	// OutputModeRun represents run output mode.
	OutputModeRun = "run"
	// OutputModeSymbols represents symbols output mode.
	OutputModeSymbols = "symbols"
	// OutputModeDocs represents documentation generation output mode.
	OutputModeDocs = "docs"
	// OutputModeHover represents hover documentation output mode.
	OutputModeHover = "hover"
	// File permissions
	dirPermissions  = 0o755
	filePermissions = 0o644
	dirPerms        = 0o750
	anyType         = "Any"
	typeString      = "string"
	typeInt         = "int"
	typeBool        = "bool"
)

// CommandResult holds the result of running a Osprey command.
type CommandResult struct {
	Output     string `json:"output"`
	Success    bool   `json:"success"`
	ErrorMsg   string `json:"errorMsg,omitempty"`
	OutputFile string `json:"outputFile,omitempty"`
}

// VexErrorListener handles parsing errors for the Osprey compiler.
type VexErrorListener struct {
	*antlr.DefaultErrorListener
	HasErrors bool
}

// SyntaxError handles syntax errors during parsing.
func (v *VexErrorListener) SyntaxError(
	_ antlr.Recognizer,
	_ interface{},
	line, column int,
	msg string,
	_ antlr.RecognitionException,
) {
	fmt.Printf("Syntax error at line %d:%d - %s\n", line, column, msg)
	v.HasErrors = true
}

// RunCommand executes a Osprey command with the given filename and mode.
// This function mimics exactly what the CLI does but returns testable results.
func RunCommand(filename, outputMode, docsDir string) CommandResult {
	// Handle docs mode (no file required)
	if outputMode == OutputModeDocs {
		return runGenerateDocs(docsDir)
	}

	// Handle hover mode (element name passed as filename)
	if outputMode == OutputModeHover {
		return runGetHoverDocumentation(filename)
	}

	// Read source file for all other modes
	content, err := os.ReadFile(filename) // #nosec G304 - filename is from CLI args
	if err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Error reading file: %v", err),
		}
	}

	source := string(content)

	switch outputMode {
	case OutputModeAST:
		return runShowAST(source, filename)
	case OutputModeLLVM:
		return runShowLLVM(source, filename)
	case OutputModeCompile:
		return runCompileToExecutable(source, filename)
	case OutputModeRun:
		return runRunProgram(source)
	case OutputModeSymbols:
		return runShowSymbols(source, filename)
	default:
		return CommandResult{
			Success:  false,
			ErrorMsg: "Unknown output mode: " + outputMode,
		}
	}
}

// RunCommandWithSecurity executes a Osprey command with the given filename, mode, and security configuration.
// This function is used for testing security restrictions.
func RunCommandWithSecurity(filename, outputMode string, security *SecurityConfig) CommandResult {
	// Handle docs mode (no file required)
	if outputMode == OutputModeDocs {
		return runGenerateDocs("") // Security doesn't affect docs generation
	}

	// Handle hover mode (element name passed as filename)
	if outputMode == OutputModeHover {
		return runGetHoverDocumentation(filename) // Security doesn't affect hover docs
	}

	// Read source file for all other modes
	content, err := os.ReadFile(filename) // #nosec G304 - filename is from CLI args
	if err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Error reading file: %v", err),
		}
	}

	source := string(content)

	switch outputMode {
	case OutputModeAST:
		return runShowASTWithSecurity(source, filename, security)
	case OutputModeLLVM:
		return runShowLLVMWithSecurity(source, filename, security)
	case OutputModeCompile:
		return runCompileToExecutableWithSecurity(source, filename, security)
	case OutputModeRun:
		return runRunProgramWithSecurity(source, security)
	case OutputModeSymbols:
		return runShowSymbolsWithSecurity(source, filename, security)
	default:
		return CommandResult{
			Success:  false,
			ErrorMsg: "Unknown output mode: " + outputMode,
		}
	}
}

func runShowAST(source, filename string) CommandResult {
	// Create input stream from source
	input := antlr.NewInputStream(source)

	// Create lexer
	lexer := parser.NewospreyLexer(input)
	errorListener := &VexErrorListener{}
	errorListener.HasErrors = false

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser
	p := parser.NewospreyParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)

	// Parse starting from program rule
	tree := p.Program()

	if errorListener.HasErrors {
		return CommandResult{
			Success:  false,
			ErrorMsg: "Found syntax errors",
		}
	}

	// Build AST from parse tree
	builder := ast.NewBuilder()
	program := builder.BuildProgram(tree)

	output := fmt.Sprintf("AST for %s:\nProgram with %d statements\n", filename, len(program.Statements))

	return CommandResult{
		Output:  output,
		Success: true,
	}
}

func runShowLLVM(source, filename string) CommandResult {
	ir, err := codegen.CompileToLLVM(source)
	if err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Error generating LLVM IR: %v", err),
		}
	}

	output := fmt.Sprintf("; LLVM IR for %s\n%s\n", filename, ir)

	return CommandResult{
		Output:  output,
		Success: true,
	}
}

func runCompileToExecutable(source, filename string) CommandResult {
	// Determine output filename (remove extension, put in outputs/ relative to source file)
	baseName := filepath.Base(filename)
	if ext := filepath.Ext(baseName); ext != "" {
		baseName = baseName[:len(baseName)-len(ext)]
	}

	// Create outputs directory relative to the source file
	sourceDir := filepath.Dir(filename)
	outputsDir := filepath.Join(sourceDir, "outputs")
	if err := os.MkdirAll(outputsDir, dirPerms); err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Failed to create outputs directory: %v", err),
		}
	}

	outputName := filepath.Join(outputsDir, baseName)

	if err := codegen.CompileToExecutable(source, outputName); err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Compilation failed: %v", err),
		}
	}

	output := fmt.Sprintf("Compiling %s to %s...\n", filename, outputName)

	return CommandResult{
		Output:     output,
		Success:    true,
		OutputFile: outputName,
	}
}

func runRunProgram(source string) CommandResult {
	if err := codegen.CompileAndRun(source); err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Execution failed: %v", err),
		}
	}

	return CommandResult{
		Output:  "Running program...\n",
		Success: true,
	}
}

func runShowSymbols(source, _ string) CommandResult {
	// Create input stream from source
	input := antlr.NewInputStream(source)

	// Create lexer
	lexer := parser.NewospreyLexer(input)
	errorListener := &VexErrorListener{}
	errorListener.HasErrors = false

	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errorListener)

	// Create token stream
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create parser
	p := parser.NewospreyParser(stream)
	p.RemoveErrorListeners()
	p.AddErrorListener(errorListener)

	// Parse starting from program rule
	tree := p.Program()

	if errorListener.HasErrors {
		return CommandResult{
			Success:  false,
			ErrorMsg: "Found syntax errors",
		}
	}

	// Build AST from parse tree
	builder := ast.NewBuilder()
	program := builder.BuildProgram(tree)

	// Extract symbols
	symbols := extractSymbols(program)

	// Convert to JSON
	jsonData, err := json.MarshalIndent(symbols, "", "  ")
	if err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Error generating JSON: %v", err),
		}
	}

	return CommandResult{
		Output:  string(jsonData),
		Success: true,
	}
}

// SymbolInfo represents symbol information for the language server.
type SymbolInfo struct {
	Name          string          `json:"name"`
	Kind          string          `json:"kind"`          // "function", "variable", "type"
	Type          string          `json:"type"`          // The type of the symbol
	Line          int             `json:"line"`          // 1-based line number
	Column        int             `json:"column"`        // 1-based column number
	Documentation string          `json:"documentation"` // Documentation from comments
	Signature     string          `json:"signature"`     // Function signature
	Parameters    []ParameterInfo `json:"parameters"`    // Function parameters
	ReturnType    string          `json:"returnType"`    // Function return type
}

// ParameterInfo represents function parameter information.
type ParameterInfo struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	Documentation string `json:"documentation"`
}

// Removed old documentation types - now using descriptions package

func extractSymbols(program *ast.Program) []SymbolInfo {
	var symbols []SymbolInfo

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.FunctionDeclaration:
			symbols = append(symbols, extractFunctionSymbol(s))
		case *ast.LetDeclaration:
			symbols = append(symbols, extractVariableSymbol(s))
		case *ast.TypeDeclaration:
			symbols = append(symbols, extractTypeSymbol(s))
		}
	}

	return symbols
}

func extractFunctionSymbol(fn *ast.FunctionDeclaration) SymbolInfo {
	var params []ParameterInfo

	for _, param := range fn.Parameters {
		paramType := anyType
		if param.Type != nil {
			paramType = param.Type.Name
		}

		params = append(params, ParameterInfo{
			Name:          param.Name,
			Type:          paramType,
			Documentation: "", // Could extract from comments if available
		})
	}

	var returnType string
	if fn.ReturnType != nil {
		returnType = fn.ReturnType.Name
	} else {
		// Use compiler's type inference logic for functions without explicit return types
		returnType = inferFunctionReturnType(fn)
	}

	signature := fmt.Sprintf("%s(%s) -> %s", fn.Name, getParameterSignature(params), returnType)

	return SymbolInfo{
		Name:          fn.Name,
		Kind:          "function",
		Type:          fmt.Sprintf("Function(%s) -> %s", getParameterSignature(params), returnType),
		Line:          1, // Would need position info from parser
		Column:        1,
		Documentation: "", // Could extract from comments
		Signature:     signature,
		Parameters:    params,
		ReturnType:    returnType,
	}
}

func getParameterSignature(params []ParameterInfo) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string

	for _, param := range params {
		parts = append(parts, fmt.Sprintf("%s: %s", param.Name, param.Type))
	}

	return strings.Join(parts, ", ")
}

func extractVariableSymbol(letDecl *ast.LetDeclaration) SymbolInfo {
	varType := anyType
	if letDecl.Type != nil {
		varType = normalizeTypeName(letDecl.Type.Name)
	} else if letDecl.Value != nil {
		varType = normalizeTypeName(inferTypeFromExpression(letDecl.Value))
	}

	return SymbolInfo{
		Name:          letDecl.Name,
		Kind:          "variable",
		Type:          varType,
		Line:          1,
		Column:        1,
		Documentation: "",
		Signature:     fmt.Sprintf("%s: %s", letDecl.Name, varType),
		Parameters:    []ParameterInfo{},
		ReturnType:    "",
	}
}

func extractTypeSymbol(typeDecl *ast.TypeDeclaration) SymbolInfo {
	return SymbolInfo{
		Name:          typeDecl.Name,
		Kind:          "type",
		Type:          "Type",
		Line:          1,
		Column:        1,
		Documentation: "",
		Signature:     "type " + typeDecl.Name,
		Parameters:    []ParameterInfo{},
		ReturnType:    "",
	}
}

// inferFunctionReturnType infers the return type of a function using simplified logic.
func inferFunctionReturnType(fn *ast.FunctionDeclaration) string {
	if fn.Body == nil {
		return anyType
	}

	// Simple return type inference based on body expression type
	return inferTypeFromExpression(fn.Body)
}

func normalizeTypeName(typeName string) string {
	switch typeName {
	case "int", "Int":
		return "Int"
	case "string", "String":
		return "String"
	case "bool", "Bool":
		return "Bool"
	default:
		return typeName
	}
}

func inferTypeFromExpression(expr ast.Expression) string {
	switch e := expr.(type) {
	case *ast.IntegerLiteral:
		return typeInt
	case *ast.StringLiteral:
		return typeString
	case *ast.BooleanLiteral:
		return typeBool
	case *ast.BinaryExpression:
		switch e.Operator {
		case "+", "-", "*", "/", "%":
			// For arithmetic operations, check operands
			leftType := inferTypeFromExpression(e.Left)
			rightType := inferTypeFromExpression(e.Right)
			if leftType == typeString || rightType == typeString {
				return typeString
			}
			return typeInt
		case "==", "!=", "<", "<=", ">", ">=":
			return typeBool
		default:
			return typeInt
		}
	case *ast.CallExpression:
		// Could analyze known function return types here
		return anyType
	case *ast.MatchExpression:
		// Could analyze match arms to determine return type
		return anyType
	case *ast.Identifier:
		// For identifiers, we can't determine type without context
		return anyType
	default:
		return anyType
	}
}

func runGenerateDocs(docsDir string) CommandResult {
	err := generateAPIDocumentationFiles(docsDir)
	if err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Failed to generate documentation: %v", err),
		}
	}

	return CommandResult{
		Success: true,
		Output:  "API reference documentation generated successfully!\n",
	}
}

func runGetHoverDocumentation(elementName string) CommandResult {
	hoverDoc := descriptions.GetHoverDocumentation(elementName)
	if hoverDoc == "" {
		return CommandResult{
			Success:  false,
			ErrorMsg: "No documentation found for element: " + elementName,
		}
	}

	return CommandResult{
		Success: true,
		Output:  hoverDoc,
	}
}

func generateAPIDocumentationFiles(docsDir string) error {
	if err := os.MkdirAll(docsDir, dirPermissions); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	if err := generateFunctionDocs(docsDir); err != nil {
		return fmt.Errorf("failed to generate function docs: %w", err)
	}

	if err := generateTypeDocs(docsDir); err != nil {
		return fmt.Errorf("failed to generate type docs: %w", err)
	}

	if err := generateOperatorDocs(docsDir); err != nil {
		return fmt.Errorf("failed to generate operator docs: %w", err)
	}

	if err := generateKeywordDocs(docsDir); err != nil {
		return fmt.Errorf("failed to generate keyword docs: %w", err)
	}

	if err := generateIndexFiles(docsDir); err != nil {
		return fmt.Errorf("failed to generate index files: %w", err)
	}

	return nil
}

func generateFunctionDocs(docsDir string) error {
	functionsDir := filepath.Join(docsDir, "functions")
	if err := os.MkdirAll(functionsDir, dirPermissions); err != nil {
		return err
	}

	functionDocs := descriptions.GetBuiltinFunctionDescriptions()

	// Sort function names alphabetically for deterministic output
	var functionNames []string
	for name := range functionDocs {
		functionNames = append(functionNames, name)
	}
	sort.Strings(functionNames)

	for _, name := range functionNames {
		doc := functionDocs[name]
		content := generateFunctionMarkdown(doc)
		filename := filepath.Join(functionsDir, strings.ToLower(name)+".md")
		if err := os.WriteFile(filename, []byte(content), filePermissions); err != nil {
			return fmt.Errorf("failed to write function doc %s: %w", name, err)
		}
	}

	return nil
}

func generateTypeDocs(docsDir string) error {
	typesDir := filepath.Join(docsDir, "types")
	if err := os.MkdirAll(typesDir, dirPermissions); err != nil {
		return err
	}

	typeDocs := descriptions.GetBuiltinTypeDescriptions()

	// Sort type names alphabetically for deterministic output
	var typeNames []string
	for name := range typeDocs {
		typeNames = append(typeNames, name)
	}
	sort.Strings(typeNames)

	for _, name := range typeNames {
		doc := typeDocs[name]
		content := generateTypeMarkdown(doc)
		filename := filepath.Join(typesDir, strings.ToLower(name)+".md")
		if err := os.WriteFile(filename, []byte(content), filePermissions); err != nil {
			return fmt.Errorf("failed to write type doc %s: %w", name, err)
		}
	}

	return nil
}

func generateOperatorDocs(docsDir string) error {
	operatorsDir := filepath.Join(docsDir, "operators")
	if err := os.MkdirAll(operatorsDir, dirPermissions); err != nil {
		return err
	}

	operatorDocs := descriptions.GetOperatorDescriptions()

	// Sort operator symbols alphabetically for deterministic output
	var operatorSymbols []string
	for symbol := range operatorDocs {
		operatorSymbols = append(operatorSymbols, symbol)
	}
	sort.Strings(operatorSymbols)

	for _, symbol := range operatorSymbols {
		doc := operatorDocs[symbol]
		content := generateOperatorMarkdown(doc)
		safeFilename := getOperatorFilename(symbol)
		filename := filepath.Join(operatorsDir, safeFilename+".md")

		if err := os.WriteFile(filename, []byte(content), filePermissions); err != nil {
			return fmt.Errorf("failed to write operator doc %s: %w", symbol, err)
		}
	}

	return nil
}

func generateKeywordDocs(docsDir string) error {
	keywordsDir := filepath.Join(docsDir, "keywords")
	if err := os.MkdirAll(keywordsDir, dirPermissions); err != nil {
		return err
	}

	keywordDocs := descriptions.GetKeywordDescriptions()

	// Sort keyword names alphabetically for deterministic output
	var keywordNames []string
	for name := range keywordDocs {
		keywordNames = append(keywordNames, name)
	}
	sort.Strings(keywordNames)

	for _, name := range keywordNames {
		doc := keywordDocs[name]
		content := generateKeywordMarkdown(doc)
		filename := filepath.Join(keywordsDir, strings.ToLower(name)+".md")
		if err := os.WriteFile(filename, []byte(content), filePermissions); err != nil {
			return fmt.Errorf("failed to write keyword doc %s: %w", name, err)
		}
	}

	return nil
}

func generateIndexFiles(docsDir string) error {
	if err := generateMainIndex(docsDir); err != nil {
		return err
	}

	if err := generateFunctionIndex(docsDir); err != nil {
		return err
	}

	if err := generateTypeIndex(docsDir); err != nil {
		return err
	}

	if err := generateOperatorIndex(docsDir); err != nil {
		return err
	}

	if err := generateKeywordIndex(docsDir); err != nil {
		return err
	}

	return nil
}

func generateMainIndex(docsDir string) error {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString("title: \"API Reference - Osprey Programming Language\"\n")
	content.WriteString("description: \"Complete reference documentation for the Osprey programming language\"\n")
	content.WriteString("---\n\n")

	content.WriteString("## Quick Navigation\n\n")
	content.WriteString("- [Functions](functions/) - Built-in functions for I/O, iteration, and data transformation\n")
	content.WriteString("- [Types](types/) - Built-in data types (Int, String, Bool, Any)\n")
	content.WriteString("- [Operators](operators/) - Arithmetic, comparison, and logical operators\n")
	content.WriteString("- [Keywords](keywords/) - Language keywords (fn, let, type, match, import)\n\n")

	// Generate Function Reference table dynamically with sorted order
	content.WriteString("## Function Reference\n\n")
	content.WriteString("| Function | Description |\n")
	content.WriteString("|----------|-------------|\n")

	functionDocs := descriptions.GetBuiltinFunctionDescriptions()
	var functionNames []string
	for name := range functionDocs {
		functionNames = append(functionNames, name)
	}
	sort.Strings(functionNames)

	for _, name := range functionNames {
		doc := functionDocs[name]
		linkName := strings.ToLower(name)
		content.WriteString(fmt.Sprintf("| [%s](functions/%s/) | %s |\n", name, linkName, doc.Description))
	}
	content.WriteString("\n")

	// Generate Type Reference table dynamically with sorted order
	content.WriteString("## Type Reference\n\n")
	content.WriteString("| Type | Description |\n")
	content.WriteString("|------|-------------|\n")

	typeDocs := descriptions.GetBuiltinTypeDescriptions()
	var typeNames []string
	for name := range typeDocs {
		typeNames = append(typeNames, name)
	}
	sort.Strings(typeNames)

	for _, name := range typeNames {
		doc := typeDocs[name]
		linkName := strings.ToLower(name)
		content.WriteString(fmt.Sprintf("| [%s](types/%s/) | %s |\n", name, linkName, doc.Description))
	}
	content.WriteString("\n")

	// Generate Operator Reference table dynamically with sorted order
	content.WriteString("## Operator Reference\n\n")
	content.WriteString("| Operator | Name | Description |\n")
	content.WriteString("|----------|------|-------------|\n")

	operatorDocs := descriptions.GetOperatorDescriptions()
	var operatorSymbols []string
	for symbol := range operatorDocs {
		operatorSymbols = append(operatorSymbols, symbol)
	}
	sort.Strings(operatorSymbols)

	for _, symbol := range operatorSymbols {
		doc := operatorDocs[symbol]
		filename := getOperatorFilename(symbol)
		content.WriteString(fmt.Sprintf("| [%s](operators/%s/) | %s | %s |\n", symbol, filename, doc.Name, doc.Description))
	}
	content.WriteString("\n")

	// Generate Keyword Reference table dynamically with sorted order
	content.WriteString("## Keyword Reference\n\n")
	content.WriteString("| Keyword | Description |\n")
	content.WriteString("|---------|-------------|\n")

	keywordDocs := descriptions.GetKeywordDescriptions()
	var keywordNames []string
	for name := range keywordDocs {
		keywordNames = append(keywordNames, name)
	}
	sort.Strings(keywordNames)

	for _, name := range keywordNames {
		doc := keywordDocs[name]
		linkName := strings.ToLower(name)
		content.WriteString(fmt.Sprintf("| [%s](keywords/%s/) | %s |\n", name, linkName, doc.Description))
	}
	content.WriteString("\n")

	filename := filepath.Join(docsDir, "index.md")
	return os.WriteFile(filename, []byte(content.String()), filePermissions)
}

func generateFunctionIndex(docsDir string) error {
	var content strings.Builder
	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString("title: \"Built-in Functions\"\n")
	content.WriteString("description: \"Complete reference for all built-in functions in Osprey\"\n")
	content.WriteString("---\n\n")
	content.WriteString("All built-in functions available in Osprey.\n\n")

	functionDocs := descriptions.GetBuiltinFunctionDescriptions()

	// Sort function names alphabetically for deterministic output
	var functionNames []string
	for name := range functionDocs {
		functionNames = append(functionNames, name)
	}
	sort.Strings(functionNames)

	for _, name := range functionNames {
		doc := functionDocs[name]
		content.WriteString(fmt.Sprintf("## [%s](%s/)\n\n", doc.Name, strings.ToLower(doc.Name)))
		content.WriteString(fmt.Sprintf("**Signature:** `%s`\n\n", doc.Signature))
		content.WriteString(doc.Description + "\n\n")
	}

	filename := filepath.Join(docsDir, "functions", "index.md")
	return os.WriteFile(filename, []byte(content.String()), filePermissions)
}

func generateTypeIndex(docsDir string) error {
	var content strings.Builder
	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString("title: \"Built-in Types\"\n")
	content.WriteString("description: \"Complete reference for all built-in types in Osprey\"\n")
	content.WriteString("---\n\n")
	content.WriteString("All built-in types available in Osprey.\n\n")

	typeDocs := descriptions.GetBuiltinTypeDescriptions()

	// Sort type names alphabetically for deterministic output
	var typeNames []string
	for name := range typeDocs {
		typeNames = append(typeNames, name)
	}
	sort.Strings(typeNames)

	for _, name := range typeNames {
		doc := typeDocs[name]
		content.WriteString(fmt.Sprintf("## [%s](%s/)\n\n", doc.Name, strings.ToLower(doc.Name)))
		content.WriteString(doc.Description + "\n\n")
	}

	filename := filepath.Join(docsDir, "types", "index.md")
	return os.WriteFile(filename, []byte(content.String()), filePermissions)
}

func generateOperatorIndex(docsDir string) error {
	var content strings.Builder
	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString("title: \"Operators\"\n")
	content.WriteString("description: \"Complete reference for all operators in Osprey\"\n")
	content.WriteString("---\n\n")
	content.WriteString("All operators available in Osprey.\n\n")

	operatorDocs := descriptions.GetOperatorDescriptions()

	// Sort operator symbols alphabetically for deterministic output
	var operatorSymbols []string
	for symbol := range operatorDocs {
		operatorSymbols = append(operatorSymbols, symbol)
	}
	sort.Strings(operatorSymbols)

	for _, symbol := range operatorSymbols {
		doc := operatorDocs[symbol]
		filename := getOperatorFilename(symbol)
		content.WriteString(fmt.Sprintf("## [%s](%s/) - %s\n\n", doc.Symbol, filename, doc.Name))
		content.WriteString(doc.Description + "\n\n")
	}

	filename := filepath.Join(docsDir, "operators", "index.md")
	return os.WriteFile(filename, []byte(content.String()), filePermissions)
}

func generateKeywordIndex(docsDir string) error {
	var content strings.Builder
	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString("title: \"Language Keywords\"\n")
	content.WriteString("description: \"Complete reference for all language keywords in Osprey\"\n")
	content.WriteString("---\n\n")
	content.WriteString("All language keywords available in Osprey.\n\n")

	keywordDocs := descriptions.GetKeywordDescriptions()

	// Sort keyword names alphabetically for deterministic output
	var keywordNames []string
	for name := range keywordDocs {
		keywordNames = append(keywordNames, name)
	}
	sort.Strings(keywordNames)

	for _, name := range keywordNames {
		doc := keywordDocs[name]
		content.WriteString(fmt.Sprintf("## [%s](%s/)\n\n", doc.Keyword, strings.ToLower(doc.Keyword)))
		content.WriteString(doc.Description + "\n\n")
	}

	filename := filepath.Join(docsDir, "keywords", "index.md")
	return os.WriteFile(filename, []byte(content.String()), filePermissions)
}

func generateFunctionMarkdown(doc *descriptions.BuiltinFunctionDesc) string {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString(fmt.Sprintf("title: \"%s (Function)\"\n", doc.Name))
	content.WriteString(fmt.Sprintf("description: \"%s\"\n", doc.Description))
	content.WriteString("---\n\n")
	content.WriteString(fmt.Sprintf("**Signature:** `%s`\n\n", doc.Signature))
	content.WriteString(fmt.Sprintf("**Description:** %s\n\n", doc.Description))

	if len(doc.Parameters) > 0 {
		content.WriteString("## Parameters\n\n")
		for _, param := range doc.Parameters {
			content.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", param.Name, param.Type, param.Description))
		}
		content.WriteString("\n")
	}

	content.WriteString(fmt.Sprintf("**Returns:** %s\n\n", doc.ReturnType))

	if doc.Example != "" {
		content.WriteString("## Example\n\n")
		content.WriteString("```osprey\n")
		content.WriteString(strings.ReplaceAll(doc.Example, "\\n", "\n"))
		content.WriteString("\n```\n")
	}

	return content.String()
}

func generateTypeMarkdown(doc *descriptions.BuiltinTypeDesc) string {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString(fmt.Sprintf("title: \"%s (Type)\"\n", doc.Name))
	content.WriteString(fmt.Sprintf("description: \"%s\"\n", doc.Description))
	content.WriteString("---\n\n")
	content.WriteString(fmt.Sprintf("**Description:** %s\n\n", doc.Description))

	if doc.Example != "" {
		content.WriteString("## Example\n\n")
		content.WriteString("```osprey\n")
		content.WriteString(strings.ReplaceAll(doc.Example, "\\n", "\n"))
		content.WriteString("\n```\n")
	}

	return content.String()
}

func generateOperatorMarkdown(doc *descriptions.OperatorDesc) string {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString(fmt.Sprintf("title: \"%s (%s Operator)\"\n", doc.Symbol, doc.Name))
	content.WriteString(fmt.Sprintf("description: \"%s\"\n", doc.Description))
	content.WriteString("---\n\n")
	content.WriteString(fmt.Sprintf("**Description:** %s\n\n", doc.Description))

	if doc.Example != "" {
		content.WriteString("## Example\n\n")
		content.WriteString("```osprey\n")
		content.WriteString(strings.ReplaceAll(doc.Example, "\\n", "\n"))
		content.WriteString("\n```\n")
	}

	return content.String()
}

func generateKeywordMarkdown(doc *descriptions.KeywordDesc) string {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString("layout: page\n")
	content.WriteString(fmt.Sprintf("title: \"%s (Keyword)\"\n", doc.Keyword))
	content.WriteString(fmt.Sprintf("description: \"%s\"\n", doc.Description))
	content.WriteString("---\n\n")
	content.WriteString(fmt.Sprintf("**Description:** %s\n\n", doc.Description))

	if doc.Example != "" {
		content.WriteString("## Example\n\n")
		content.WriteString("```osprey\n")
		content.WriteString(strings.ReplaceAll(doc.Example, "\\n", "\n"))
		content.WriteString("\n```\n")
	}

	return content.String()
}

func getOperatorFilename(symbol string) string {
	switch symbol {
	case "+":
		return "plus"
	case "-":
		return "minus"
	case "*":
		return "multiply"
	case "/":
		return "divide"
	case "%":
		return "modulo"
	case "==":
		return "equal"
	case "!=":
		return "not-equal"
	case "<":
		return "less-than"
	case "<=":
		return "less-equal"
	case ">":
		return "greater-than"
	case ">=":
		return "greater-equal"
	case "|":
		return "pipe"
	case "|>":
		return "pipe-operator"
	default:
		return strings.ReplaceAll(symbol, ".", "_")
	}
}

// Security-aware versions of the functions
func runShowASTWithSecurity(source, filename string, _ *SecurityConfig) CommandResult {
	// AST generation doesn't need security restrictions
	return runShowAST(source, filename)
}

func runShowLLVMWithSecurity(source, filename string, security *SecurityConfig) CommandResult {
	// Convert CLI security config to codegen security config
	codegenSecurity := convertToCodegenSecurity(security)

	ir, err := codegen.CompileToLLVMWithSecurity(source, codegenSecurity)
	if err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Error generating LLVM IR: %v", err),
		}
	}

	output := fmt.Sprintf("; LLVM IR for %s\n%s\n", filename, ir)

	return CommandResult{
		Output:  output,
		Success: true,
	}
}

func runCompileToExecutableWithSecurity(source, filename string, security *SecurityConfig) CommandResult {
	// Determine output filename (remove extension, put in outputs/ relative to source file)
	baseName := filepath.Base(filename)
	if ext := filepath.Ext(baseName); ext != "" {
		baseName = baseName[:len(baseName)-len(ext)]
	}

	// Create outputs directory relative to the source file
	sourceDir := filepath.Dir(filename)
	outputsDir := filepath.Join(sourceDir, "outputs")
	if err := os.MkdirAll(outputsDir, dirPerms); err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Failed to create outputs directory: %v", err),
		}
	}

	outputName := filepath.Join(outputsDir, baseName)

	// Convert CLI security config to codegen security config
	codegenSecurity := convertToCodegenSecurity(security)

	if err := codegen.CompileToExecutableWithSecurity(source, outputName, codegenSecurity); err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Compilation failed: %v", err),
		}
	}

	output := fmt.Sprintf("Compiling %s to %s...\n", filename, outputName)

	return CommandResult{
		Output:     output,
		Success:    true,
		OutputFile: outputName,
	}
}

func runRunProgramWithSecurity(source string, security *SecurityConfig) CommandResult {
	// Convert CLI security config to codegen security config
	codegenSecurity := convertToCodegenSecurity(security)

	if err := codegen.CompileAndRunWithSecurity(source, codegenSecurity); err != nil {
		return CommandResult{
			Success:  false,
			ErrorMsg: fmt.Sprintf("Execution failed: %v", err),
		}
	}

	return CommandResult{
		Output:  "Running program...\n",
		Success: true,
	}
}

func runShowSymbolsWithSecurity(source, filename string, _ *SecurityConfig) CommandResult {
	// Symbol extraction doesn't need security restrictions
	return runShowSymbols(source, filename)
}

// convertToCodegenSecurity converts CLI SecurityConfig to codegen SecurityConfig.
func convertToCodegenSecurity(security *SecurityConfig) codegen.SecurityConfig {
	return codegen.SecurityConfig{
		AllowHTTP:             security.AllowHTTP,
		AllowWebSocket:        security.AllowWebSocket,
		AllowFileRead:         security.AllowFileRead,
		AllowFileWrite:        security.AllowFileWrite,
		AllowFFI:              security.AllowFFI,
		AllowProcessExecution: security.AllowProcessExecution,
		SandboxMode:           security.SandboxMode,
	}
}
