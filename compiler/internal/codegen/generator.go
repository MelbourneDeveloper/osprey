// Package codegen provides code generation and execution capabilities for Osprey.
package codegen

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
	"github.com/christianfindlay/osprey/internal/logging"
	"github.com/christianfindlay/osprey/internal/plugins"
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
	// Structured diagnostics
	logger *slog.Logger
	// HINDLEY-MILNER FIX: Single source of truth for record field mappings
	// Maps record type name to field name -> LLVM index mapping
	recordFieldMappings map[string]map[string]int
	// Stream Fusion: Track pending transformations for map/filter
	pendingMapFunc    *ast.Identifier // Pending map transformation function
	pendingFilterFunc *ast.Identifier // Pending filter predicate function
	// Language plugin system — invoked during codegen for `fn <plugin> <name>(...) = <body>` declarations
	pluginSystem *plugins.PluginSystem
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
		logger:      logging.Logger("codegen"),
		// HINDLEY-MILNER FIX: Initialize record field mappings
		recordFieldMappings: make(map[string]map[string]int),
	}

	// Initialize the language plugin system, rooted at the compiler directory
	// where the plugins/ directory lives. Resolution order:
	//   1) OSPREY_COMPILER_DIR env var (explicit override)
	//   2) parent of the running osprey binary's directory (binary lives in <compilerDir>/bin/)
	//   3) current working directory (covers `go test ./...` from compiler/)
	generator.pluginSystem = plugins.NewPluginSystem(resolveCompilerDir())

	// Declare external functions for FFI
	generator.declareExternalFunctions()

	// Register built-in types
	generator.registerBuiltInTypes()

	// Set the generator reference in the type inferer for constraint checking
	generator.typeInferer.SetGenerator(generator)

	// Initialize fiber runtime declarations will happen on first use

	return generator
}

// resolveCompilerDir locates the compiler directory containing plugins/.
// Resolution order, falling through to the working directory as a last resort:
//  1. OSPREY_COMPILER_DIR (explicit override — used by tests and packaged installs)
//  2. <compilerDir>/bin/osprey — derived from os.Executable() (covers `./bin/osprey`)
//  3. Walk up from the working directory looking for a `plugins/` dir
//     (covers `go test ./...` from any subdir under compiler/)
//  4. Working directory
//
// Step 3 matters because the Go test runner sets cwd to the test package's directory,
// not the compiler root — without the walk, `go test ./tests/integration/...` would
// look for plugins under tests/integration/plugins.
func resolveCompilerDir() string {
	dir := os.Getenv("OSPREY_COMPILER_DIR")
	if dir != "" {
		return dir
	}

	exe, exeErr := os.Executable()
	if exeErr == nil {
		candidate := filepath.Dir(filepath.Dir(exe))
		if candidate != "" && candidate != "." && hasPluginsDir(candidate) {
			return candidate
		}
	}

	wd, wdErr := os.Getwd()
	if wdErr != nil {
		return "."
	}

	if ancestor := walkUpForPluginsDir(wd); ancestor != "" {
		return ancestor
	}
	return wd
}

func hasPluginsDir(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, "plugins"))
	return err == nil
}

// walkUpForPluginsDir ascends from start looking for a directory that contains a
// plugins/ subdirectory. Returns "" if none is found before reaching the filesystem root.
func walkUpForPluginsDir(start string) string {
	current := start
	for {
		if hasPluginsDir(current) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
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

	// Declare effect runtime functions for dynamic handler resolution
	// i32 @__osprey_handler_push(i8* %effect_name, i8* %operation_name, i8* %handler_func_ptr)
	handlerPush := g.module.NewFunc("__osprey_handler_push", types.I32,
		ir.NewParam("effect_name", types.I8Ptr),
		ir.NewParam("operation_name", types.I8Ptr),
		ir.NewParam("handler_func_ptr", types.I8Ptr),
	)
	g.functions["__osprey_handler_push"] = handlerPush

	// i32 @__osprey_handler_pop()
	handlerPop := g.module.NewFunc("__osprey_handler_pop", types.I32)
	g.functions["__osprey_handler_pop"] = handlerPop

	// i8* @__osprey_handler_lookup(i8* %effect_name, i8* %operation_name)
	handlerLookup := g.module.NewFunc("__osprey_handler_lookup", types.I8Ptr,
		ir.NewParam("effect_name", types.I8Ptr),
		ir.NewParam("operation_name", types.I8Ptr),
	)
	g.functions["__osprey_handler_lookup"] = handlerLookup
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
