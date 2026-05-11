package codegen

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// pluginParamPlaceholderRegex matches `$paramName` placeholders inside a plugin language body
// (e.g. SQL: `WHERE id = $userId`). Letters/digits/underscore are accepted to match the
// SQL plugin's parameter validation regex.
var pluginParamPlaceholderRegex = regexp.MustCompile(`\$(\w+)`)

// Sentinel errors for plugin codegen. We wrap these with %w so callers can errors.Is
// against them while still getting a context-rich message.
var (
	ErrPluginSystemUninit       = errors.New("plugin system not initialised")
	ErrPluginValidation         = errors.New("plugin failed to validate function")
	ErrPluginPlaceholderUnbound = errors.New("plugin body references parameter not declared on function")
	ErrPluginParamType          = errors.New("unsupported plugin parameter type")
)

const (
	// sprintfInt64Spec is the printf-family format specifier for an i64 value.
	// Plugin functions promote bool/i32/etc. up to i64 before printing.
	sprintfInt64Spec = "%lld"
	// int64BitSize is the bit width of LLVM i64 — used to decide when an integer
	// parameter needs sign-extending before being passed to sprintf.
	int64BitSize = 64
)

func (g *LLVMGenerator) generateStatement(stmt ast.Statement) error {
	switch s := stmt.(type) {
	case *ast.ImportStatement:
		// Imports are handled at compile time
		return nil

	case *ast.LetDeclaration:
		_, err := g.generateLetDeclaration(s)
		return err

	case *ast.AssignmentStatement:
		return g.generateAssignmentStatement(s)

	case *ast.FunctionDeclaration:
		// Use the unified function declaration generator from function_signatures.go
		return g.declareFunctionSignature(s)

	case *ast.ExternDeclaration:
		return g.generateExternDeclaration(s)

	case *ast.PluginFunctionDeclaration:
		err := g.generatePluginFunctionDeclaration(s)

		return err

	case *ast.TypeDeclaration:
		// Type declarations are handled in first pass
		return nil

	case *ast.EffectDeclaration:
		return g.generateEffectDeclaration(s)

	case *ast.ExpressionStatement:
		_, err := g.generateExpression(s.Expression)
		return err

	default:
		return WrapUnsupportedStatement(stmt)
	}
}

func (g *LLVMGenerator) generateExternDeclaration(externDecl *ast.ExternDeclaration) error {
	// Convert extern parameters to LLVM parameters
	params := make([]*ir.Param, len(externDecl.Parameters))
	paramNames := make([]string, len(externDecl.Parameters))

	for i, param := range externDecl.Parameters {
		llvmType := g.typeExpressionToLLVMType(&param.Type)
		params[i] = ir.NewParam(param.Name, llvmType)
		paramNames[i] = param.Name
	}

	// Determine return type
	var returnType types.Type = types.I64 // Default to int
	if externDecl.ReturnType != nil {
		returnType = g.typeExpressionToLLVMType(externDecl.ReturnType)
	}

	// Declare the external function
	externFunc := g.module.NewFunc(externDecl.Name, returnType, params...)
	g.functions[externDecl.Name] = externFunc
	// Built-in functions are handled by Hindley-Milner type inference
	g.functionParameters[externDecl.Name] = paramNames

	// Add extern function to type environment for type inference
	// Convert extern parameters to type inference types
	paramTypes := make([]Type, len(externDecl.Parameters))
	for i, param := range externDecl.Parameters {
		paramTypes[i] = g.typeExpressionToInferenceType(&param.Type)
	}

	// Determine inference return type
	var inferenceReturnType Type = &ConcreteType{name: TypeInt} // Default to int
	if externDecl.ReturnType != nil {
		inferenceReturnType = g.typeExpressionToInferenceType(externDecl.ReturnType)
	}

	// Add extern function to type environment
	functionType := &FunctionType{
		paramTypes: paramTypes,
		returnType: inferenceReturnType,
	}
	g.typeInferer.env.Set(externDecl.Name, functionType)

	return nil
}

// generatePluginFunctionDeclaration compiles a `fn <plugin> <name>(...) = <body>` declaration.
//
// What it does at COMPILE TIME:
//  1. Invokes the named plugin as a subprocess (via PluginSystem), passing the function name,
//     parameters, and language body as JSON over stdin.
//  2. If the plugin reports an error (e.g. unknown SQL operation, parameter referenced in the
//     body but not declared on the function), compilation FAILS — no placeholder, no silent pass.
//
// What it emits as RUNTIME code:
//
//	An LLVM function returning i8* (string). The body uses sprintf to splice the function's
//	parameter values into the validated language body at the `$paramName` placeholders.
//	For SQL, that means `fn sql q(id: Int) = SELECT * FROM users WHERE id = $id` called with
//	id=42 returns the string "SELECT * FROM users WHERE id = 42" — ready to hand to a SQL driver.
//
// This is the bridge between compile-time plugin validation and runtime use.
func (g *LLVMGenerator) generatePluginFunctionDeclaration(pluginDecl *ast.PluginFunctionDeclaration) error {
	if g.pluginSystem == nil {
		return fmt.Errorf("%w: function %q", ErrPluginSystemUninit, pluginDecl.FunctionName)
	}

	// Invoke the plugin at compile time. Failure here MUST fail compilation —
	// silently emitting a placeholder would mask real validation errors (per CLAUDE.md).
	_, pluginErr := g.pluginSystem.ProcessPluginFunction(pluginDecl, "", 0)
	if pluginErr != nil {
		return fmt.Errorf("plugin %q failed to validate function %q: %w",
			pluginDecl.PluginName, pluginDecl.FunctionName, pluginErr)
	}

	params := make([]*ir.Param, 0, len(pluginDecl.Parameters))
	paramNames := make([]string, 0, len(pluginDecl.Parameters))
	paramInferTypes := make([]Type, 0, len(pluginDecl.Parameters))
	for _, param := range pluginDecl.Parameters {
		params = append(params, ir.NewParam(param.Name, g.typeExpressionToLLVMType(param.Type)))
		paramNames = append(paramNames, param.Name)
		paramInferTypes = append(paramInferTypes, g.typeExpressionToInferenceType(param.Type))
	}

	pluginFunc := g.module.NewFunc(pluginDecl.FunctionName, types.I8Ptr, params...)
	g.functions[pluginDecl.FunctionName] = pluginFunc
	g.functionParameters[pluginDecl.FunctionName] = paramNames

	// Register the plugin function in the type environment so call sites can resolve it.
	// Plugin functions always return String (the language body with placeholders filled in).
	g.typeInferer.env.Set(
		pluginDecl.FunctionName,
		NewFunctionType(paramInferTypes, NewPrimitiveType(TypeString)),
	)

	savedBuilder := g.builder
	entry := pluginFunc.NewBlock("entry")
	g.builder = entry

	result, err := g.emitPluginLanguageBody(pluginDecl, pluginFunc)
	if err != nil {
		g.builder = savedBuilder
		return err
	}
	entry.NewRet(result)
	g.builder = savedBuilder
	return nil
}

// emitPluginLanguageBody generates LLVM IR that builds the plugin's language body string
// at runtime, substituting `$paramName` placeholders with the current parameter values.
//
// When the body has no placeholders we emit a single global string constant. When it does,
// we build a sprintf format string by replacing each `$name` with the format specifier that
// matches that parameter's LLVM type, then call sprintf into a stack buffer.
func (g *LLVMGenerator) emitPluginLanguageBody(
	pluginDecl *ast.PluginFunctionDeclaration, pluginFunc *ir.Func,
) (value.Value, error) {
	matches := pluginParamPlaceholderRegex.FindAllStringSubmatchIndex(pluginDecl.PluginContent, -1)
	if len(matches) == 0 {
		return g.emitStaticString(pluginDecl.PluginContent), nil
	}

	paramIndex := make(map[string]int, len(pluginDecl.Parameters))
	for i, p := range pluginDecl.Parameters {
		paramIndex[p.Name] = i
	}

	var formatParts []string
	args := make([]value.Value, 0, len(matches))
	cursor := 0
	for _, m := range matches {
		start, end, nameStart, nameEnd := m[0], m[1], m[2], m[3]
		formatParts = append(formatParts, escapeSprintfLiteral(pluginDecl.PluginContent[cursor:start]))

		paramName := pluginDecl.PluginContent[nameStart:nameEnd]
		idx, ok := paramIndex[paramName]
		if !ok {
			// Plugin should have caught this, but be defensive — fail compilation rather than emit bad IR.
			return nil, fmt.Errorf("%w: plugin %q, function %q, placeholder $%s",
				ErrPluginPlaceholderUnbound, pluginDecl.PluginName, pluginDecl.FunctionName, paramName)
		}

		paramVal := pluginFunc.Params[idx]
		spec, arg, specErr := g.sprintfSpecForParam(paramVal)
		if specErr != nil {
			return nil, fmt.Errorf("plugin %q, function %q, parameter %q: %w",
				pluginDecl.PluginName, pluginDecl.FunctionName, paramName, specErr)
		}
		formatParts = append(formatParts, spec)
		args = append(args, arg)
		cursor = end
	}
	formatParts = append(formatParts, escapeSprintfLiteral(pluginDecl.PluginContent[cursor:]))

	return g.callSprintfInto(strings.Join(formatParts, ""), args), nil
}

// emitStaticString creates a global string constant and returns an i8* pointer into it.
func (g *LLVMGenerator) emitStaticString(s string) value.Value {
	str := constant.NewCharArrayFromString(s + StringTerminator)
	global := g.module.NewGlobalDef("", str)
	return g.builder.NewGetElementPtr(str.Typ, global,
		constant.NewInt(types.I32, ArrayIndexZero),
		constant.NewInt(types.I32, ArrayIndexZero))
}

// sprintfSpecForParam returns the printf-family format specifier for an LLVM parameter
// value and the (possibly converted) argument to feed sprintf alongside it.
func (g *LLVMGenerator) sprintfSpecForParam(param *ir.Param) (string, value.Value, error) {
	switch t := param.Type().(type) {
	case *types.IntType:
		if t.BitSize == 1 {
			// LLVM bool — promote to i64 and let sprintf print 0/1.
			return sprintfInt64Spec, g.builder.NewZExt(param, types.I64), nil
		}
		if t.BitSize < int64BitSize {
			return sprintfInt64Spec, g.builder.NewSExt(param, types.I64), nil
		}
		return sprintfInt64Spec, param, nil
	case *types.FloatType:
		return "%g", param, nil
	case *types.PointerType:
		// String (i8*) or other pointer — render as a string.
		return "%s", param, nil
	default:
		return "", nil, fmt.Errorf("%w: %v", ErrPluginParamType, t)
	}
}

// callSprintfInto allocates a stack buffer, calls sprintf with the supplied format string
// and arguments, and returns the i8* pointing at the populated buffer.
func (g *LLVMGenerator) callSprintfInto(formatString string, args []value.Value) value.Value {
	sprintf := g.ensureSprintfDeclaration()

	formatConst := constant.NewCharArrayFromString(formatString + StringTerminator)
	formatGlobal := g.module.NewGlobalDef("", formatConst)
	formatPtr := g.builder.NewGetElementPtr(formatConst.Typ, formatGlobal,
		constant.NewInt(types.I32, ArrayIndexZero),
		constant.NewInt(types.I32, ArrayIndexZero))

	bufferType := types.NewArray(BufferSize1KB, types.I8)
	buffer := g.builder.NewAlloca(bufferType)
	bufferPtr := g.builder.NewGetElementPtr(bufferType, buffer,
		constant.NewInt(types.I32, ArrayIndexZero),
		constant.NewInt(types.I32, ArrayIndexZero))

	sprintfArgs := make([]value.Value, 0, len(args)+TwoArgs)
	sprintfArgs = append(sprintfArgs, bufferPtr, formatPtr)
	sprintfArgs = append(sprintfArgs, args...)
	g.builder.NewCall(sprintf, sprintfArgs...)

	return bufferPtr
}

// escapeSprintfLiteral escapes `%` so a literal segment of the plugin body cannot be
// interpreted as a sprintf directive.
func escapeSprintfLiteral(s string) string {
	return strings.ReplaceAll(s, "%", "%%")
}

// typeExpressionToLLVMType converts an Osprey TypeExpression to an LLVM type.
func (g *LLVMGenerator) typeExpressionToLLVMType(typeExpr *ast.TypeExpression) types.Type {
	// Handle function types
	if typeExpr.IsFunction {
		// Build parameter types
		paramTypes := make([]types.Type, len(typeExpr.ParameterTypes))
		for i, paramType := range typeExpr.ParameterTypes {
			paramTypes[i] = g.typeExpressionToLLVMType(&paramType)
		}

		// Build return type
		var returnType types.Type = types.I64 // Default to int
		if typeExpr.ReturnType != nil {
			returnType = g.typeExpressionToLLVMType(typeExpr.ReturnType)
		}

		// Create function signature
		funcSig := types.NewFunc(returnType, paramTypes...)

		// Return pointer to function (function pointer type)
		return types.NewPointer(funcSig)
	}

	switch typeExpr.Name {
	case TypeInt:
		return types.I64
	case TypeString:
		return types.I8Ptr
	case TypeUnit:
		return types.Void
	case TypeHTTPResponse:
		// Return pointer to HttpResponse struct
		return types.NewPointer(g.typeMap[TypeHTTPResponse])
	default:
		// Check if it's a user-defined type
		if llvmType, exists := g.typeMap[typeExpr.Name]; exists {
			return llvmType
		}
		// Default to i64 for unknown types
		return types.I64
	}
}

// typeExpressionToInferenceType converts an Osprey TypeExpression to a type inference Type.
func (g *LLVMGenerator) typeExpressionToInferenceType(typeExpr *ast.TypeExpression) Type {
	// Handle function types
	if typeExpr.IsFunction {
		// Build parameter types
		paramTypes := make([]Type, len(typeExpr.ParameterTypes))
		for i, paramType := range typeExpr.ParameterTypes {
			paramTypes[i] = g.typeExpressionToInferenceType(&paramType)
		}

		// Build return type
		var returnType Type = &ConcreteType{name: TypeInt} // Default to int
		if typeExpr.ReturnType != nil {
			returnType = g.typeExpressionToInferenceType(typeExpr.ReturnType)
		}

		// Create function type
		return &FunctionType{
			paramTypes: paramTypes,
			returnType: returnType,
		}
	}

	// Handle generic types like Result<bool, MathError>
	if len(typeExpr.GenericParams) > 0 {
		// Convert generic parameters to Type arguments
		typeArgs := make([]Type, len(typeExpr.GenericParams))
		for i, genericParam := range typeExpr.GenericParams {
			typeArgs[i] = g.typeExpressionToInferenceType(&genericParam)
		}

		genericType := NewGenericType(typeExpr.Name, typeArgs)

		return genericType
	}

	switch typeExpr.Name {
	case TypeInt:
		return &ConcreteType{name: TypeInt}
	case "string":
		return &ConcreteType{name: TypeString}
	case TypeBool:
		return &ConcreteType{name: TypeBool}
	case TypeUnit:
		return &ConcreteType{name: TypeUnit}
	case TypeHTTPResponse:
		return &ConcreteType{name: TypeHTTPResponse}
	case TypeFiber:
		return &ConcreteType{name: TypeFiber}
	case TypeChannel:
		return &ConcreteType{name: TypeChannel}
	default:
		// Check if this is a user-defined type (record or union)
		if typeDecl, exists := g.typeDeclarations[typeExpr.Name]; exists {
			// If it's a single-variant record type, return RecordType
			if len(typeDecl.Variants) == 1 && len(typeDecl.Variants[0].Fields) > 0 {
				fields := make(map[string]Type)

				variant := &typeDecl.Variants[0]
				for _, field := range variant.Fields {
					// For now, use the field's declared type or default to int
					// This should ideally be more sophisticated type resolution
					var fieldType Type

					switch field.Type {
					case TypeInt:
						fieldType = &ConcreteType{name: TypeInt}
					case TypeString:
						fieldType = &ConcreteType{name: TypeString}
					case TypeBool:
						fieldType = &ConcreteType{name: TypeBool}
					default:
						fieldType = &ConcreteType{name: field.Type}
					}

					fields[field.Name] = fieldType
				}

				return NewRecordType(typeExpr.Name, fields)
			}
		}

		// For unknown types without generic parameters, return as concrete type
		return &ConcreteType{name: typeExpr.Name}
	}
}

func (g *LLVMGenerator) generateLetDeclaration(letDecl *ast.LetDeclaration) (value.Value, error) {
	value, err := g.generateExpression(letDecl.Value)
	if err != nil {
		return nil, err
	}

	// Store the value in our variable map
	g.variables[letDecl.Name] = value

	// Use explicit type annotation if present, otherwise infer from value
	var varType Type

	if letDecl.Type != nil {
		// Use the explicit type annotation, but validate it matches the value
		annotatedType := g.typeExpressionToInferenceType(letDecl.Type)

		// Infer the actual type of the value
		valueType, err := g.typeInferer.InferType(letDecl.Value)
		if err != nil {
			return nil, err
		}

		// Check if the value type is compatible with the annotation
		err = g.typeInferer.Unify(annotatedType, valueType)
		if err != nil {
			return nil, WrapTypeMismatchWithPos(
				valueType.String(), letDecl.Name, annotatedType.String(), letDecl.Position)
		}

		varType = annotatedType
	} else {
		// Use unified type inference system for untyped declarations
		inferredType, err := g.typeInferer.InferType(letDecl.Value)
		if err != nil {
			return nil, err
		}

		varType = inferredType
	}

	// Store the type in the Hindley-Milner environment
	g.typeInferer.env.Set(letDecl.Name, varType)

	// Track if this variable is mutable
	g.mutableVariables[letDecl.Name] = letDecl.Mutable

	return value, nil
}

// generateAssignmentStatement generates LLVM IR for mutable variable assignments.
func (g *LLVMGenerator) generateAssignmentStatement(assignStmt *ast.AssignmentStatement) error {
	// Check if the variable exists in the Hindley-Milner type environment (single source of truth)
	if _, exists := g.typeInferer.env.Get(assignStmt.Name); !exists {
		return WrapUndefinedVariableWithPos(assignStmt.Name, assignStmt.Position)
	}

	// Check if the variable is mutable
	if mutable, exists := g.mutableVariables[assignStmt.Name]; !exists || !mutable {
		return WrapImmutableAssignmentErrorWithPos(assignStmt.Name, assignStmt.Position)
	}

	// Generate the new value
	newValue, err := g.generateExpression(assignStmt.Value)
	if err != nil {
		return err
	}

	// Use unified type inference system
	inferredType, err := g.typeInferer.InferType(assignStmt.Value)
	if err != nil {
		return err
	}

	// Verify type compatibility using unification
	existingType := g.typeInferer.env.vars[assignStmt.Name]
	err = g.typeInferer.Unify(existingType, inferredType)
	if err != nil {
		return fmt.Errorf("type mismatch in assignment: %w", err)
	}

	// Update the variable
	g.variables[assignStmt.Name] = newValue
	// Update type in Hindley-Milner environment
	g.typeInferer.env.Set(assignStmt.Name, inferredType)

	return nil
}
