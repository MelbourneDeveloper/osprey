package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// ensureHTTPFunctionDeclaration ensures an HTTP function is declared using builtin registry
func (g *LLVMGenerator) ensureHTTPFunctionDeclaration(functionName string) *ir.Func {
	if fn, exists := g.functions[functionName]; exists {
		return fn
	}

	// Get function details from builtin registry
	builtinFunc, exists := GlobalBuiltInRegistry.GetFunction(functionName)
	if !exists {
		return nil
	}

	// Convert builtin parameters to LLVM parameters
	params := make([]*ir.Param, len(builtinFunc.ParameterTypes))
	for i, param := range builtinFunc.ParameterTypes {
		llvmType := g.getLLVMType(param.Type)
		params[i] = ir.NewParam(param.Name, llvmType)
	}

	// Convert return type to LLVM type
	returnType := g.getLLVMType(builtinFunc.ReturnType)

	// Create function with the correct C runtime name (which is now in the registry)
	fn := g.module.NewFunc(builtinFunc.CName, returnType, params...)
	g.functions[functionName] = fn

	return fn
}

// generateHTTPFunctionCall generates a call to an HTTP function with argument validation
func (g *LLVMGenerator) generateHTTPFunctionCall(functionName string, callExpr *ast.CallExpression,
	expectedArgs int) (value.Value, error) {
	// Handle named arguments if present
	if len(callExpr.NamedArguments) > 0 {
		return g.generateHTTPFunctionCallNamed(functionName, callExpr, expectedArgs)
	}

	// Validate argument count
	if len(callExpr.Arguments) != expectedArgs {
		return nil, WrapHTTPFunctionWrongArgs(functionName, expectedArgs, len(callExpr.Arguments))
	}

	// Generate argument values
	argValues := make([]value.Value, len(callExpr.Arguments))
	for i, arg := range callExpr.Arguments {
		val, err := g.generateExpression(arg)
		if err != nil {
			return nil, err
		}

		argValues[i] = val
	}

	// Ensure function is declared
	fn := g.ensureHTTPFunctionDeclaration(functionName)
	if fn == nil {
		return nil, WrapHTTPFunctionNotFound(functionName)
	}

	// Call the function
	return g.builder.NewCall(fn, argValues...), nil
}

// generateHTTPFunctionCallNamed generates a call to an HTTP function with named arguments
func (g *LLVMGenerator) generateHTTPFunctionCallNamed(functionName string, callExpr *ast.CallExpression,
	expectedArgs int) (value.Value, error) {
	// Get function details from builtin registry
	builtinFunc, exists := GlobalBuiltInRegistry.GetFunction(functionName)
	if !exists {
		return nil, WrapHTTPFunctionNotFound(functionName)
	}

	// Validate argument count
	if len(callExpr.NamedArguments) != expectedArgs {
		return nil, WrapHTTPFunctionWrongArgs(functionName, expectedArgs, len(callExpr.NamedArguments))
	}

	// Generate argument values in correct order
	argValues := make([]value.Value, len(builtinFunc.ParameterTypes))
	for i, param := range builtinFunc.ParameterTypes {
		found := false

		for _, namedArg := range callExpr.NamedArguments {
			if namedArg.Name == param.Name {
				val, err := g.generateExpression(namedArg.Value)
				if err != nil {
					return nil, err
				}

				argValues[i] = val
				found = true

				break
			}
		}

		if !found {
			return nil, WrapHTTPFunctionMissingNamedArg(functionName, param.Name)
		}
	}

	// Ensure function is declared
	fn := g.ensureHTTPFunctionDeclaration(functionName)
	if fn == nil {
		return nil, WrapHTTPFunctionNotFound(functionName)
	}

	// Call the function
	return g.builder.NewCall(fn, argValues...), nil
}

// HTTP Server Functions
func (g *LLVMGenerator) generateHTTPCreateServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(HTTPCreateServerOsprey, callExpr, TwoArgs)
}

func (g *LLVMGenerator) generateHTTPListenCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(HTTPListenOsprey, callExpr, TwoArgs)
}

func (g *LLVMGenerator) generateHTTPStopServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(HTTPStopServerOsprey, callExpr, OneArg)
}

// HTTP Client Functions
func (g *LLVMGenerator) generateHTTPCreateClientCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(HTTPCreateClientOsprey, callExpr, TwoArgs)
}

func (g *LLVMGenerator) generateHTTPGetCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPRequestWithMethod(callExpr, ThreeArgs, HTTPMethodGet)
}

func (g *LLVMGenerator) generateHTTPPostCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPRequestWithMethod(callExpr, FourArgs, HTTPMethodPost)
}

func (g *LLVMGenerator) generateHTTPPutCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPRequestWithMethod(callExpr, FourArgs, HTTPMethodPut)
}

func (g *LLVMGenerator) generateHTTPDeleteCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPRequestWithMethod(callExpr, ThreeArgs, HTTPMethodDelete)
}

func (g *LLVMGenerator) generateHTTPRequestCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(HTTPRequestOsprey, callExpr, FiveArgs)
}

func (g *LLVMGenerator) generateHTTPCloseClientCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(HTTPCloseClientOsprey, callExpr, OneArg)
}

// generateHTTPRequestWithMethod generates HTTP requests with specific methods
func (g *LLVMGenerator) generateHTTPRequestWithMethod(callExpr *ast.CallExpression, expectedArgs int,
	method int64) (value.Value, error) {
	if len(callExpr.Arguments) != expectedArgs {
		return nil, WrapHTTPFunctionWrongArgs(HTTPRequestOsprey, expectedArgs, len(callExpr.Arguments))
	}

	// Generate client ID
	clientID, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Generate path
	path, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Create method constant
	methodVal := constant.NewInt(types.I64, method)

	// Generate headers and body based on method
	var headers, body value.Value
	if method == HTTPMethodGet || method == HTTPMethodDelete {
		// GET/DELETE: headers from arg 2, empty body
		headers, err = g.generateExpression(callExpr.Arguments[2])
		if err != nil {
			return nil, err
		}

		body = g.createEmptyStringConstant()
	} else {
		// POST/PUT: body from arg 2, headers from arg 3
		body, err = g.generateExpression(callExpr.Arguments[2])
		if err != nil {
			return nil, err
		}

		headers, err = g.generateExpression(callExpr.Arguments[3])
		if err != nil {
			return nil, err
		}
	}

	// Call http_request
	fn := g.ensureHTTPFunctionDeclaration(HTTPRequestOsprey)

	return g.builder.NewCall(fn, clientID, methodVal, path, headers, body), nil
}

// createEmptyStringConstant creates a constant empty string
func (g *LLVMGenerator) createEmptyStringConstant() value.Value {
	emptyStr := constant.NewCharArrayFromString("")
	emptyGlobal := g.module.NewGlobalDef("", emptyStr)

	return g.builder.NewGetElementPtr(emptyStr.Typ, emptyGlobal,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
}

// WebSocket Functions
func (g *LLVMGenerator) generateWebSocketConnectCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketConnectOsprey, callExpr, TwoArgs)
}

func (g *LLVMGenerator) generateWebSocketSendCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketSendOsprey, callExpr, TwoArgs)
}

func (g *LLVMGenerator) generateWebSocketCloseCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketCloseOsprey, callExpr, OneArg)
}

func (g *LLVMGenerator) generateWebSocketCreateServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketCreateServerOsprey, callExpr, ThreeArgs)
}

func (g *LLVMGenerator) generateWebSocketServerListenCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketServerListenOsprey, callExpr, OneArg)
}

func (g *LLVMGenerator) generateWebSocketServerBroadcastCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketServerBroadcastOsprey, callExpr, TwoArgs)
}

func (g *LLVMGenerator) generateWebSocketStopServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketStopServerOsprey, callExpr, OneArg)
}

func (g *LLVMGenerator) generateWebSocketKeepAliveCall(callExpr *ast.CallExpression) (value.Value, error) {
	return g.generateHTTPFunctionCall(WebSocketKeepAliveOsprey, callExpr, 0)
}
