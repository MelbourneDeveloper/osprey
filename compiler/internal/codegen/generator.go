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
	// Result pattern matching support
	currentResultValue value.Value // Current Result value being pattern matched
	// Security configuration
	security SecurityConfig
	// New fields for the new constructor
	stringConstants       map[string]value.Value
	currentFunction       *ir.Func
	currentFunctionParams map[string]value.Value
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
		security:              security,
		stringConstants:       make(map[string]value.Value),
		currentResultValue:    nil,
		currentFunction:       nil,
		currentFunctionParams: make(map[string]value.Value),
	}

	// Declare external functions for FFI
	generator.declareExternalFunctions()

	// Register built-in function return types
	generator.registerBuiltInFunctionReturnTypes()

	// Register built-in types
	generator.registerBuiltInTypes()

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

	// Declare strlen: i64 @strlen(i8* %str)
	strlen := g.module.NewFunc("strlen", types.I64, ir.NewParam("str", types.I8Ptr))
	g.functions["strlen"] = strlen

	// Declare strstr: i8* @strstr(i8* %haystack, i8* %needle)
	strstr := g.module.NewFunc("strstr", types.I8Ptr,
		ir.NewParam("haystack", types.I8Ptr),
		ir.NewParam("needle", types.I8Ptr),
	)
	g.functions["strstr"] = strstr

	// Declare malloc: i8* @malloc(i64 %size)
	malloc := g.module.NewFunc("malloc", types.I8Ptr, ir.NewParam("size", types.I64))
	g.functions["malloc"] = malloc

	// Declare memcpy: i8* @memcpy(i8* %dest, i8* %src, i64 %n)
	memcpy := g.module.NewFunc("memcpy", types.I8Ptr,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr),
		ir.NewParam("n", types.I64),
	)
	g.functions["memcpy"] = memcpy
}

// registerBuiltInFunctionReturnTypes registers return types for built-in functions.
func (g *LLVMGenerator) registerBuiltInFunctionReturnTypes() {
	// TODO: Most of these are WRONG!
	// Anything that COULD fail MUST return a RESULT
	// Especially IO functions like readFile, writeFile, etc.
	// DON'T IGNORE THIS!!! FIX IT!!

	// Core functions
	g.functionReturnTypes["toString"] = TypeString
	g.functionReturnTypes["print"] = TypeInt // Returns exit code
	g.functionReturnTypes["input"] = TypeInt // Returns input as integer
	g.functionReturnTypes["range"] = TypeInt // Returns range object (simplified as int)
	// STRING FUNCTIONS RETURN RESULT TYPES - THEY CAN FAIL!
	g.functionReturnTypes["length"] = TypeResult + "<Int, string>"       // Returns Result<Int, string>
	g.functionReturnTypes["contains"] = TypeResult + "<Bool, string>"    // Returns Result<Bool, string>
	g.functionReturnTypes["substring"] = TypeResult + "<String, string>" // Returns Result<String, string>

	// Process and file functions - MUST return Result types per spec
	g.functionReturnTypes["spawnProcess"] = TypeResult + "<ProcessResult, string>" // Returns Result<ProcessResult, string>
	g.functionReturnTypes["writeFile"] = TypeResult + "<Success, string>"          // Returns Result<Success, string>
	g.functionReturnTypes["readFile"] = TypeResult + "<string, string>"            // Returns Result<string, string>
	g.functionReturnTypes["parseJSON"] = TypeResult + "<string, string>"           // Returns Result<string, string>
	g.functionReturnTypes["extractCode"] = TypeResult + "<string, string>"         // Returns Result<string, string>

	// HTTP functions
	g.functionReturnTypes["httpCreateServer"] = TypeInt // Returns server ID
	g.functionReturnTypes["httpListen"] = TypeInt       // Returns status code
	g.functionReturnTypes["httpStopServer"] = TypeInt   // Returns status code
	g.functionReturnTypes["httpCreateClient"] = TypeInt // Returns client ID
	g.functionReturnTypes["httpGet"] = TypeInt          // Returns status code
	g.functionReturnTypes["httpPost"] = TypeInt         // Returns status code
	g.functionReturnTypes["httpPut"] = TypeInt          // Returns status code
	g.functionReturnTypes["httpDelete"] = TypeInt       // Returns status code
	g.functionReturnTypes["httpRequest"] = TypeInt      // Returns status code
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

// registerBuiltInTypes registers built-in types in the type system.
func (g *LLVMGenerator) registerBuiltInTypes() {
	// Register HttpResponse as a built-in struct type
	httpResponseType := types.NewStruct(
		types.I64,   // status: Int
		types.I8Ptr, // headers: String
		types.I8Ptr, // contentType: String
		types.I64,   // contentLength: Int
		types.I64,   // streamFd: Int
		types.I1,    // isComplete: Bool
		types.I8Ptr, // partialBody: String
		types.I64,   // partialLength: Int
	)

	g.typeMap[TypeHTTPResponse] = httpResponseType
}
