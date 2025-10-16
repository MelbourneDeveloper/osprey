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
	module             *ir.Module
	builder            *ir.Block
	function           *ir.Func
	variables          map[string]value.Value
	mutableVariables   map[string]bool // Track which variables were declared as mutable
	functions          map[string]*ir.Func
	functionParameters map[string][]string // Track function parameter names for named argument reordering
	typeMap            map[string]types.Type
	// Union type tracking
	typeDeclarations map[string]*ast.TypeDeclaration // Track all type declarations
	unionVariants    map[string]int64                // Track union variants and their discriminant values
	// Monomorphization tracking
	monomorphizedInstances map[string]string                   // Track monomorphized function instances
	functionDeclarations   map[string]*ast.FunctionDeclaration // Track original function declarations
	// Fiber closure counter
	closureCounter int
	// Result pattern matching support
	currentResultValue value.Value // Current Result value being pattern matched
	// Security configuration
	security SecurityConfig
	// New fields for the new constructor
	stringConstants       map[string]value.Value
	currentFunction       *ir.Func
	currentFunctionParams map[string]value.Value
	// Real algebraic effects system
	effectCodegen *EffectCodegen
	// Context for type-aware literal generation
	expectedReturnType    types.Type
	expectedParameterType types.Type
	// Hindley-Milner type inference system
	typeInferer *TypeInferer
	// HINDLEY-MILNER FIX: Single source of truth for record field mappings
	// Maps record type name to field name -> LLVM index mapping
	recordFieldMappings map[string]map[string]int
	// Stream Fusion: Track pending transformations for map/filter
	pendingMapFunc    *ast.Identifier // Pending map transformation function
	pendingFilterFunc *ast.Identifier // Pending filter predicate function
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
		module:             module,
		variables:          make(map[string]value.Value),
		mutableVariables:   make(map[string]bool),
		functions:          make(map[string]*ir.Func),
		functionParameters: make(map[string][]string),
		typeMap:            typeMap,
		// Initialize union type tracking
		typeDeclarations: make(map[string]*ast.TypeDeclaration),
		unionVariants:    make(map[string]int64),
		// Set security configuration
		security:              security,
		stringConstants:       make(map[string]value.Value),
		currentResultValue:    nil,
		currentFunction:       nil,
		currentFunctionParams: make(map[string]value.Value),
		// Initialize Hindley-Milner type inference system
		typeInferer: NewTypeInferer(),
		// HINDLEY-MILNER FIX: Initialize record field mappings
		recordFieldMappings: make(map[string]map[string]int),
	}

	// Declare external functions for FFI
	generator.declareExternalFunctions()

	// Register built-in types
	generator.registerBuiltInTypes()

	// Set the generator reference in the type inferer for constraint checking
	generator.typeInferer.SetGenerator(generator)

	// Initialize fiber runtime declarations will happen on first use

	return generator
}

// GenerateIR returns the LLVM IR as a string.
func (g *LLVMGenerator) GenerateIR() string {
	return g.module.String()
}

// InitializeEffects initializes the effect system for the generator
func (g *LLVMGenerator) InitializeEffects() {
	g.effectCodegen = g.NewEffectCodegen()
}

// RegisterEffectDeclaration registers an effect declaration with the effect system
func (g *LLVMGenerator) RegisterEffectDeclaration(effect *ast.EffectDeclaration) error {
	if g.effectCodegen == nil {
		g.InitializeEffects()
	}

	return g.effectCodegen.RegisterEffect(effect)
}

// generateRealPerformExpression generates real algebraic effects perform expressions
func (g *LLVMGenerator) generateRealPerformExpression(perform *ast.PerformExpression) (value.Value, error) {
	if g.effectCodegen == nil {
		g.InitializeEffects()
	}

	return g.effectCodegen.GeneratePerformExpression(perform)
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

// registerBuiltInTypes registers built-in types in the type system.
func (g *LLVMGenerator) registerBuiltInTypes() {
	// Register HttpResponse as a built-in struct type (REMOVED REDUNDANT LENGTH FIELDS)
	httpResponseType := types.NewStruct(
		types.I64,   // status: Int
		types.I8Ptr, // headers: String
		types.I8Ptr, // contentType: String
		types.I64,   // streamFd: Int
		types.I1,    // isComplete: Bool
		types.I8Ptr, // partialBody: String (runtime calculates length automatically)
	)

	g.typeMap[TypeHTTPResponse] = httpResponseType

	// Register ProcessHandle as Int64 (process ID)
	g.typeMap["ProcessHandle"] = types.I64
}
