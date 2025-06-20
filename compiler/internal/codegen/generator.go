// Package codegen provides code generation and execution capabilities for Osprey.
package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// LLVMGenerator generates LLVM IR from AST.
type LLVMGenerator struct {
	module              *ir.Module
	builder             *ir.Block
	function            *ir.Func
	variables           map[string]value.Value
	variableTypes       map[string]string // Track variable types: "string" or "int"
	functions           map[string]*ir.Func
	functionReturnTypes map[string]string   // Track function return types: "string" or "int"
	functionParameters  map[string][]string // Track function parameter names for named argument reordering
	typeMap             map[string]types.Type
	// Union type tracking
	typeDeclarations map[string]*ast.TypeDeclaration // Track all type declarations
	unionVariants    map[string]int64                // Track union variants and their discriminant values
	// Fiber closure counter
	closureCounter int
	// Temporary parameter types for return type analysis
	currentFunctionParameterTypes map[string]string
	// Security configuration
	security SecurityConfig
}

// SecurityConfig defines security policies for the code generator.
// This is a copy of the SecurityConfig from cli package to avoid circular dependencies.
type SecurityConfig struct {
	AllowHTTP             bool
	AllowWebSocket        bool
	AllowFileRead         bool
	AllowFileWrite        bool
	AllowFFI              bool
	AllowProcessExecution bool
	SandboxMode           bool
}

// NewLLVMGenerator creates a new LLVM IR generator with default (permissive) security.
func NewLLVMGenerator() *LLVMGenerator {
	return NewLLVMGeneratorWithSecurity(SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	})
}

// NewLLVMGeneratorWithSecurity creates a new LLVM IR generator with specified security configuration.
func NewLLVMGeneratorWithSecurity(security SecurityConfig) *LLVMGenerator {
	module := ir.NewModule()

	// Define built-in types
	typeMap := map[string]types.Type{
		"Int":    types.I64,
		"String": types.I8Ptr, // Simplified string representation
	}

	generator := &LLVMGenerator{
		module:              module,
		variables:           make(map[string]value.Value),
		variableTypes:       make(map[string]string),
		functions:           make(map[string]*ir.Func),
		functionReturnTypes: make(map[string]string),
		functionParameters:  make(map[string][]string),
		typeMap:             typeMap,
		// Initialize union type tracking
		typeDeclarations: make(map[string]*ast.TypeDeclaration),
		unionVariants:    make(map[string]int64),
		// Set security configuration
		security: security,
	}

	// Declare external functions for FFI
	generator.declareExternalFunctions()

	// Register built-in function return types
	generator.registerBuiltInFunctionReturnTypes()

	// Initialize fiber runtime declarations will happen on first use

	return generator
}

// GenerateIR returns the LLVM IR as a string.
func (g *LLVMGenerator) GenerateIR() string {
	return g.module.String()
}

// declareExternalFunctions declares external C library functions.
func (g *LLVMGenerator) declareExternalFunctions() {
	// Declare printf: i32 @printf(i8*, ...)
	printf := g.module.NewFunc("printf", types.I32, ir.NewParam("format", types.I8Ptr))
	printf.Sig.Variadic = true
	g.functions["printf"] = printf

	// Declare puts: i32 @puts(i8* %str)
	puts := g.module.NewFunc("puts", types.I32, ir.NewParam("str", types.I8Ptr))
	g.functions["puts"] = puts

	// Declare scanf: i32 @scanf(i8* %format, ...)
	scanf := g.module.NewFunc("scanf", types.I32, ir.NewParam("format", types.I8Ptr))
	scanf.Sig.Variadic = true
	g.functions["scanf"] = scanf

	// Declare strcmp: i32 @strcmp(i8* %str1, i8* %str2)
	strcmp := g.module.NewFunc("strcmp", types.I32,
		ir.NewParam("str1", types.I8Ptr),
		ir.NewParam("str2", types.I8Ptr))
	g.functions["strcmp"] = strcmp
}

// registerBuiltInFunctionReturnTypes registers return types for built-in functions.
func (g *LLVMGenerator) registerBuiltInFunctionReturnTypes() {
	
	//TODO: Most of these are WRONG!
	//Anything that COULD fail MUST return a RESULT
	//Especially IO functions like readFile, writeFile, etc.
	
	// Core functions
	g.functionReturnTypes["toString"] = TypeString
	g.functionReturnTypes["print"] = TypeInt     // Returns exit code
	g.functionReturnTypes["input"] = TypeInt     // Returns input as integer
	g.functionReturnTypes["range"] = TypeInt     // Returns range object (simplified as int)
	g.functionReturnTypes["length"] = TypeInt    // Returns string/array length
	g.functionReturnTypes["contains"] = TypeBool // Returns boolean
	g.functionReturnTypes["substring"] = TypeString

	// Process and file functions
	g.functionReturnTypes["spawnProcess"] = TypeInt   // Returns exit code
	g.functionReturnTypes["writeFile"] = TypeInt      // Returns bytes written or error code
	g.functionReturnTypes["readFile"] = TypeString    // Returns file content
	g.functionReturnTypes["parseJSON"] = TypeString   // Returns parsed JSON as string
	g.functionReturnTypes["extractCode"] = TypeString // Returns extracted code

	// HTTP functions
	g.functionReturnTypes["httpCreateServer"] = TypeInt // Returns server ID
	g.functionReturnTypes["httpListen"] = TypeInt       // Returns status code
	g.functionReturnTypes["httpStopServer"] = TypeInt   // Returns status code
	g.functionReturnTypes["httpCreateClient"] = TypeInt // Returns client ID
	g.functionReturnTypes["httpGet"] = TypeString       // Returns response body
	g.functionReturnTypes["httpPost"] = TypeString      // Returns response body
	g.functionReturnTypes["httpPut"] = TypeString       // Returns response body
	g.functionReturnTypes["httpDelete"] = TypeString    // Returns response body
	g.functionReturnTypes["httpRequest"] = TypeString   // Returns response body
	g.functionReturnTypes["httpCloseClient"] = TypeInt  // Returns status code

	// WebSocket functions
	g.functionReturnTypes["webSocketConnect"] = TypeInt         // Returns connection ID
	g.functionReturnTypes["webSocketSend"] = TypeInt            // Returns status code
	g.functionReturnTypes["webSocketClose"] = TypeInt           // Returns status code
	g.functionReturnTypes["webSocketCreateServer"] = TypeInt    // Returns server ID
	g.functionReturnTypes["webSocketServerListen"] = TypeInt    // Returns status code
	g.functionReturnTypes["webSocketServerBroadcast"] = TypeInt // Returns status code
	g.functionReturnTypes["webSocketStopServer"] = TypeInt      // Returns status code
	g.functionReturnTypes["webSocketKeepAlive"] = TypeInt       // Returns status code

	// Functional programming functions return various types
	g.functionReturnTypes["forEach"] = TypeInt // Returns status/count
	g.functionReturnTypes["map"] = TypeInt     // Returns transformed array (simplified as int)
	g.functionReturnTypes["filter"] = TypeInt  // Returns filtered array (simplified as int)
	g.functionReturnTypes["fold"] = TypeInt    // Returns accumulated value (could be any type, simplified as int)
}
