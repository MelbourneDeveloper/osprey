package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/christianfindlay/osprey/internal/ast"
)

// ===========================================
// HTTP SERVER FUNCTION GENERATION
// ===========================================

// generateHttpCreateServerCall generates calls to http_create_server(port, address).
func (g *LLVMGenerator) generateHTTPCreateServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	// Handle named arguments by extracting them in the correct order
	if len(callExpr.NamedArguments) == TwoArgs {
		// Extract port and address from named arguments
		var portVal, addressVal value.Value
		var err error

		for _, namedArg := range callExpr.NamedArguments {
			switch namedArg.Name {
			case "port":
				portVal, err = g.generateExpression(namedArg.Value)
				if err != nil {
					return nil, err
				}
			case "address":
				addressVal, err = g.generateExpression(namedArg.Value)
				if err != nil {
					return nil, err
				}
			}
		}

		if portVal == nil || addressVal == nil {
			return nil, WrapHTTPCreateServerWrongArgs(len(callExpr.NamedArguments))
		}

		// Ensure http_create_server function is declared
		httpCreateServerFunc := g.ensureHTTPCreateServerDeclaration()

		// Call http_create_server(port, address)
		return g.builder.NewCall(httpCreateServerFunc, portVal, addressVal), nil
	}

	// Handle positional arguments
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapHTTPCreateServerWrongArgs(len(callExpr.Arguments))
	}

	// Get port argument (int)
	portVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get address argument (string)
	addressVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Ensure http_create_server function is declared
	httpCreateServerFunc := g.ensureHTTPCreateServerDeclaration()

	// Call http_create_server(port, address)

	return g.builder.NewCall(httpCreateServerFunc, portVal, addressVal), nil
}

// generateHttpListenCall generates calls to http_listen(server_id, handler_callback).
func (g *LLVMGenerator) generateHTTPListenCall(callExpr *ast.CallExpression) (value.Value, error) {
	// Handle named arguments by extracting them in the correct order
	if len(callExpr.NamedArguments) == TwoArgs {
		// Extract serverID and handler from named arguments
		var serverIDVal, handlerVal value.Value
		var err error

		for _, namedArg := range callExpr.NamedArguments {
			switch namedArg.Name {
			case "serverID":
				serverIDVal, err = g.generateExpression(namedArg.Value)
				if err != nil {
					return nil, err
				}
			case "handler":
				handlerVal, err = g.generateExpression(namedArg.Value)
				if err != nil {
					return nil, err
				}
			}
		}

		if serverIDVal == nil || handlerVal == nil {
			return nil, WrapHTTPListenWrongArgs(len(callExpr.NamedArguments))
		}

		// Ensure http_listen function is declared
		httpListenFunc := g.ensureHTTPListenDeclaration()

		// Call http_listen(server_id, handler)
		return g.builder.NewCall(httpListenFunc, serverIDVal, handlerVal), nil
	}

	// Handle positional arguments
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapHTTPListenWrongArgs(len(callExpr.Arguments))
	}

	// Get server ID argument
	serverIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get handler function argument (callback)
	handlerVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Ensure http_listen function is declared
	httpListenFunc := g.ensureHTTPListenDeclaration()

	// Call http_listen(server_id, handler)

	return g.builder.NewCall(httpListenFunc, serverIDVal, handlerVal), nil
}

// generateHttpStopServerCall generates calls to http_stop_server(server_id).
func (g *LLVMGenerator) generateHTTPStopServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	// Handle named arguments by extracting them in the correct order
	if len(callExpr.NamedArguments) == OneArg {
		// Extract serverID from named arguments
		var serverIDVal value.Value
		var err error

		for _, namedArg := range callExpr.NamedArguments {
			switch namedArg.Name {
			case "serverID":
				serverIDVal, err = g.generateExpression(namedArg.Value)
				if err != nil {
					return nil, err
				}
			default:
				return nil, WrapHTTPStopServerUnknownNamedArg(namedArg.Name)
			}
		}

		// Call the C function
		httpStopServerFunc := g.ensureHTTPStopServerDeclaration()

		return g.builder.NewCall(httpStopServerFunc, serverIDVal), nil
	}

	// Handle positional arguments (legacy support)
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapHTTPStopServerWrongArgCount(len(callExpr.Arguments))
	}

	serverID, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Call the C function
	httpStopServerFunc := g.ensureHTTPStopServerDeclaration()

	return g.builder.NewCall(httpStopServerFunc, serverID), nil
}

// ===========================================
// HTTP CLIENT FUNCTION GENERATION
// ===========================================

// generateHttpCreateClientCall generates calls to http_create_client(base_url, timeout).
func (g *LLVMGenerator) generateHTTPCreateClientCall(callExpr *ast.CallExpression) (value.Value, error) {
	if err := validateBuiltInArgs(HTTPCreateClientFunc, callExpr); err != nil {
		return nil, err
	}

	// Get base URL argument (string)
	baseURLVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get timeout argument (int)
	timeoutVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Ensure http_create_client function is declared
	httpCreateClientFunc := g.ensureHTTPCreateClientDeclaration()

	// Call http_create_client(base_url, timeout)

	return g.builder.NewCall(httpCreateClientFunc, baseURLVal, timeoutVal), nil
}

// generateHttpGetCall generates calls to http_request(client_id, GET, path, headers, "").
func (g *LLVMGenerator) generateHTTPGetCall(callExpr *ast.CallExpression) (value.Value, error) {
	if err := validateBuiltInArgs(HTTPGetFunc, callExpr); err != nil {
		return nil, err
	}

	// Get client ID
	clientIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get path
	pathVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Get headers
	headersVal, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	// Create HTTP method constant (GET = 0)
	methodVal := constant.NewInt(types.I64, 0)

	// Create empty body
	emptyBodyStr := constant.NewCharArrayFromString("")
	emptyBodyGlobal := g.module.NewGlobalDef("", emptyBodyStr)
	emptyBodyPtr := g.builder.NewGetElementPtr(emptyBodyStr.Typ, emptyBodyGlobal,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))

	// Ensure http_request function is declared
	httpRequestFunc := g.ensureHTTPRequestDeclaration()

	// Call http_request(client_id, GET, path, headers, "")

	return g.builder.NewCall(httpRequestFunc, clientIDVal, methodVal, pathVal, headersVal, emptyBodyPtr), nil
}

// generateHttpPostCall generates calls to http_request(client_id, POST, path, headers, body).
func (g *LLVMGenerator) generateHTTPPostCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != FourArgs {
		return nil, WrapHTTPPostWrongArgs(len(callExpr.Arguments))
	}

	// Get client ID
	clientIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get path
	pathVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Get body
	bodyVal, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	// Get headers
	headersVal, err := g.generateExpression(callExpr.Arguments[3])
	if err != nil {
		return nil, err
	}

	// Create HTTP method constant (POST = 1)
	methodVal := constant.NewInt(types.I64, 1)

	// Ensure http_request function is declared
	httpRequestFunc := g.ensureHTTPRequestDeclaration()

	// Call http_request(client_id, POST, path, headers, body)

	return g.builder.NewCall(httpRequestFunc, clientIDVal, methodVal, pathVal, headersVal, bodyVal), nil
}

// generateHttpPutCall generates calls to http_request(client_id, PUT, path, headers, body).
func (g *LLVMGenerator) generateHTTPPutCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != FourArgs {
		return nil, WrapHTTPPutWrongArgs(len(callExpr.Arguments))
	}

	// Get client ID
	clientIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get path
	pathVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Get body
	bodyVal, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	// Get headers
	headersVal, err := g.generateExpression(callExpr.Arguments[3])
	if err != nil {
		return nil, err
	}

	// Create HTTP method constant (PUT = 2)
	methodVal := constant.NewInt(types.I64, HTTPMethodPut)

	// Ensure http_request function is declared
	httpRequestFunc := g.ensureHTTPRequestDeclaration()

	// Call http_request(client_id, PUT, path, headers, body)

	return g.builder.NewCall(httpRequestFunc, clientIDVal, methodVal, pathVal, headersVal, bodyVal), nil
}

// generateHttpDeleteCall generates calls to http_request(client_id, DELETE, path, headers, "").
func (g *LLVMGenerator) generateHTTPDeleteCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != ThreeArgs {
		return nil, WrapHTTPDeleteWrongArgs(len(callExpr.Arguments))
	}

	// Get client ID
	clientIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get path
	pathVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Get headers
	headersVal, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	// Create HTTP method constant (DELETE = 3)
	methodVal := constant.NewInt(types.I64, HTTPMethodDelete)

	// Create empty body
	emptyBodyStr := constant.NewCharArrayFromString("")
	emptyBodyGlobal := g.module.NewGlobalDef("", emptyBodyStr)
	emptyBodyPtr := g.builder.NewGetElementPtr(emptyBodyStr.Typ, emptyBodyGlobal,
		constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))

	// Ensure http_request function is declared
	httpRequestFunc := g.ensureHTTPRequestDeclaration()

	// Call http_request(client_id, DELETE, path, headers, "")

	return g.builder.NewCall(httpRequestFunc, clientIDVal, methodVal, pathVal, headersVal, emptyBodyPtr), nil
}

// generateHttpRequestCall generates calls to http_request(client_id, method, path, headers, body).
func (g *LLVMGenerator) generateHTTPRequestCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != FiveArgs {
		return nil, WrapHTTPRequestWrongArgs(len(callExpr.Arguments))
	}

	// Get all arguments
	clientIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	methodVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	pathVal, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	headersVal, err := g.generateExpression(callExpr.Arguments[3])
	if err != nil {
		return nil, err
	}

	bodyVal, err := g.generateExpression(callExpr.Arguments[4])
	if err != nil {
		return nil, err
	}

	// Ensure http_request function is declared
	httpRequestFunc := g.ensureHTTPRequestDeclaration()

	// Call http_request(client_id, method, path, headers, body)

	return g.builder.NewCall(httpRequestFunc, clientIDVal, methodVal, pathVal, headersVal, bodyVal), nil
}

// generateHttpCloseClientCall generates calls to http_close_client(client_id).
func (g *LLVMGenerator) generateHTTPCloseClientCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapHTTPCloseClientWrongArgs(len(callExpr.Arguments))
	}

	// Get client ID argument
	clientIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Ensure http_close_client function is declared
	httpCloseClientFunc := g.ensureHTTPCloseClientDeclaration()

	// Call http_close_client(client_id)

	return g.builder.NewCall(httpCloseClientFunc, clientIDVal), nil
}

// ===========================================
// WEBSOCKET FUNCTION GENERATION
// ===========================================

// generateWebSocketConnectCall generates calls to websocket_connect(url, message_handler).
func (g *LLVMGenerator) generateWebSocketConnectCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapWebSocketConnectWrongArgs(len(callExpr.Arguments))
	}

	// Get URL argument
	urlVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get message handler argument (callback function)
	handlerVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_connect function is declared
	websocketConnectFunc := g.ensureWebSocketConnectDeclaration()

	// Call websocket_connect(url, message_handler)

	return g.builder.NewCall(websocketConnectFunc, urlVal, handlerVal), nil
}

// generateWebSocketSendCall generates calls to websocket_send(ws_id, message).
func (g *LLVMGenerator) generateWebSocketSendCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapWebSocketSendWrongArgs(len(callExpr.Arguments))
	}

	// Get WebSocket ID argument
	wsIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get message argument
	messageVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_send function is declared
	websocketSendFunc := g.ensureWebSocketSendDeclaration()

	// Call websocket_send(ws_id, message)

	return g.builder.NewCall(websocketSendFunc, wsIDVal, messageVal), nil
}

// generateWebSocketCloseCall generates calls to websocket_close(ws_id).
func (g *LLVMGenerator) generateWebSocketCloseCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapWebSocketCloseWrongArgs(len(callExpr.Arguments))
	}

	// Get WebSocket ID argument
	wsIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_close function is declared
	websocketCloseFunc := g.ensureWebSocketCloseDeclaration()

	// Call websocket_close(ws_id)

	return g.builder.NewCall(websocketCloseFunc, wsIDVal), nil
}

// ===========================================
// HTTP FUNCTION DECLARATIONS
// ===========================================

// ensureHttpCreateServerDeclaration ensures http_create_server is declared.
func (g *LLVMGenerator) ensureHTTPCreateServerDeclaration() *ir.Func {
	if fn, exists := g.functions["http_create_server"]; exists {
		return fn
	}

	// int64_t http_create_server(int64_t port, char* address)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("http_create_server", types.I64,
		ir.NewParam("port", types.I64),
		ir.NewParam("address", types.I8Ptr))

	g.functions["http_create_server"] = fn

	return fn
}

// ensureHttpListenDeclaration ensures http_listen is declared.
func (g *LLVMGenerator) ensureHTTPListenDeclaration() *ir.Func {
	if fn, exists := g.functions["http_listen"]; exists {
		return fn
	}

	// int64_t http_listen(int64_t server_id, HttpResponse* (*handler)(char*, char*, char*, char*))
	// Function that takes server ID and a callback function pointer
	handlerType := types.NewPointer(types.NewFunc(
		types.NewPointer(g.typeMap[TypeHTTPResponse]), // Return type: HttpResponse*
		types.I8Ptr, // method parameter
		types.I8Ptr, // path parameter
		types.I8Ptr, // headers parameter
		types.I8Ptr, // body parameter
	))

	fn := g.module.NewFunc("http_listen", types.I64,
		ir.NewParam("server_id", types.I64),
		ir.NewParam("handler", handlerType))

	g.functions["http_listen"] = fn

	return fn
}

// ensureHttpStopServerDeclaration ensures http_stop_server is declared.
func (g *LLVMGenerator) ensureHTTPStopServerDeclaration() *ir.Func {
	if fn, exists := g.functions["http_stop_server"]; exists {
		return fn
	}

	// int64_t http_stop_server(int64_t server_id)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("http_stop_server", types.I64,
		ir.NewParam("server_id", types.I64))

	g.functions["http_stop_server"] = fn

	return fn
}

// ensureHttpCreateClientDeclaration ensures http_create_client is declared.
func (g *LLVMGenerator) ensureHTTPCreateClientDeclaration() *ir.Func {
	if fn, exists := g.functions["http_create_client"]; exists {
		return fn
	}

	// int64_t http_create_client(char* base_url, int64_t timeout)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("http_create_client", types.I64,
		ir.NewParam("base_url", types.I8Ptr),
		ir.NewParam("timeout", types.I64))

	g.functions["http_create_client"] = fn

	return fn
}

// ensureHTTPRequestDeclaration ensures http_request is declared.
func (g *LLVMGenerator) ensureHTTPRequestDeclaration() *ir.Func {
	if fn, exists := g.functions["http_request"]; exists {
		return fn
	}

	// int64_t http_request(int64_t client_id, int64_t method, char* path, char* headers, char* body)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("http_request", types.I64,
		ir.NewParam("client_id", types.I64),
		ir.NewParam("method", types.I64),
		ir.NewParam("path", types.I8Ptr),
		ir.NewParam("headers", types.I8Ptr),
		ir.NewParam("body", types.I8Ptr))

	g.functions["http_request"] = fn

	return fn
}

// ensureHTTPCloseClientDeclaration ensures http_close_client is declared.
func (g *LLVMGenerator) ensureHTTPCloseClientDeclaration() *ir.Func {
	if fn, exists := g.functions["http_close_client"]; exists {
		return fn
	}

	// int64_t http_close_client(int64_t client_id)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("http_close_client", types.I64,
		ir.NewParam("client_id", types.I64))

	g.functions["http_close_client"] = fn

	return fn
}

// ensureWebSocketConnectDeclaration ensures websocket_connect is declared.
func (g *LLVMGenerator) ensureWebSocketConnectDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_connect"]; exists {
		return fn
	}

	// int64_t websocket_connect(char* url, void* message_handler)
	// Simple return type like fiber functions
	handlerType := types.I8Ptr // Simplified to void pointer
	fn := g.module.NewFunc("websocket_connect", types.I64,
		ir.NewParam("url", types.I8Ptr),
		ir.NewParam("handler", handlerType))

	g.functions["websocket_connect"] = fn

	return fn
}

// ensureWebSocketSendDeclaration ensures websocket_send is declared.
func (g *LLVMGenerator) ensureWebSocketSendDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_send"]; exists {
		return fn
	}

	// int64_t websocket_send(int64_t ws_id, char* message)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("websocket_send", types.I64,
		ir.NewParam("ws_id", types.I64),
		ir.NewParam("message", types.I8Ptr))

	g.functions["websocket_send"] = fn

	return fn
}

// ensureWebSocketCloseDeclaration ensures websocket_close is declared.
func (g *LLVMGenerator) ensureWebSocketCloseDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_close"]; exists {
		return fn
	}

	// int64_t websocket_close(int64_t ws_id)
	// Simple return type like fiber functions
	fn := g.module.NewFunc("websocket_close", types.I64,
		ir.NewParam("ws_id", types.I64))

	g.functions["websocket_close"] = fn

	return fn
}

// ===========================================
// WEBSOCKET SERVER FUNCTIONS
// ===========================================

// generateWebSocketCreateServerCall generates calls to websocket_create_server(port, address, path).
func (g *LLVMGenerator) generateWebSocketCreateServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != ThreeArgs {
		return nil, WrapWebSocketCreateServerWrongArgs(len(callExpr.Arguments))
	}

	// Get port argument
	portVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get address argument
	addressVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Get path argument
	pathVal, err := g.generateExpression(callExpr.Arguments[2])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_create_server function is declared
	websocketCreateServerFunc := g.ensureWebSocketCreateServerDeclaration()

	// Call websocket_create_server(port, address, path)

	return g.builder.NewCall(websocketCreateServerFunc, portVal, addressVal, pathVal), nil
}

// generateWebSocketServerListenCall generates calls to websocket_server_listen(server_id).
func (g *LLVMGenerator) generateWebSocketServerListenCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapWebSocketServerListenWrongArgs(len(callExpr.Arguments))
	}

	// Get server ID argument
	serverIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_server_listen function is declared
	websocketServerListenFunc := g.ensureWebSocketServerListenDeclaration()

	// Call websocket_server_listen(server_id)

	return g.builder.NewCall(websocketServerListenFunc, serverIDVal), nil
}

// generateWebSocketServerBroadcastCall generates calls to websocket_server_broadcast(server_id, message).
func (g *LLVMGenerator) generateWebSocketServerBroadcastCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != TwoArgs {
		return nil, WrapWebSocketServerBroadcastWrongArgs(len(callExpr.Arguments))
	}

	// Get server ID argument
	serverIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Get message argument
	messageVal, err := g.generateExpression(callExpr.Arguments[1])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_server_broadcast function is declared
	websocketServerBroadcastFunc := g.ensureWebSocketServerBroadcastDeclaration()

	// Call websocket_server_broadcast(server_id, message)

	return g.builder.NewCall(websocketServerBroadcastFunc, serverIDVal, messageVal), nil
}

// generateWebSocketStopServerCall generates calls to websocket_stop_server(server_id).
func (g *LLVMGenerator) generateWebSocketStopServerCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != OneArg {
		return nil, WrapWebSocketStopServerWrongArgs(len(callExpr.Arguments))
	}

	// Get server ID argument
	serverIDVal, err := g.generateExpression(callExpr.Arguments[0])
	if err != nil {
		return nil, err
	}

	// Ensure websocket_stop_server function is declared
	websocketStopServerFunc := g.ensureWebSocketStopServerDeclaration()

	// Call websocket_stop_server(server_id)

	return g.builder.NewCall(websocketStopServerFunc, serverIDVal), nil
}

// generateWebSocketKeepAliveCall handles WebSocket server keep alive (persistent mode).
func (g *LLVMGenerator) generateWebSocketKeepAliveCall(callExpr *ast.CallExpression) (value.Value, error) {
	if len(callExpr.Arguments) != 0 {
		return nil, ErrWebSocketKeepAliveWrongArgs
	}

	// Ensure websocket_keep_alive function is declared
	websocketKeepAliveFunc := g.ensureWebSocketKeepAliveDeclaration()

	// Call websocket_keep_alive() - void function, but we return i64 for consistency
	g.builder.NewCall(websocketKeepAliveFunc)

	// Return 0 to indicate success (even though function is void)

	return constant.NewInt(types.I64, 0), nil
}

// ===========================================
// WEBSOCKET SERVER FUNCTION DECLARATIONS
// ===========================================

// ensureWebSocketCreateServerDeclaration ensures websocket_create_server is declared.
func (g *LLVMGenerator) ensureWebSocketCreateServerDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_create_server"]; exists {
		return fn
	}

	// int64_t websocket_create_server(int64_t port, char* address, char* path)
	fn := g.module.NewFunc("websocket_create_server", types.I64,
		ir.NewParam("port", types.I64),
		ir.NewParam("address", types.I8Ptr),
		ir.NewParam("path", types.I8Ptr))

	g.functions["websocket_create_server"] = fn

	return fn
}

// ensureWebSocketServerListenDeclaration ensures websocket_server_listen is declared.
func (g *LLVMGenerator) ensureWebSocketServerListenDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_server_listen"]; exists {
		return fn
	}

	// int64_t websocket_server_listen(int64_t server_id)
	fn := g.module.NewFunc("websocket_server_listen", types.I64,
		ir.NewParam("server_id", types.I64))

	g.functions["websocket_server_listen"] = fn

	return fn
}

// ensureWebSocketServerBroadcastDeclaration ensures websocket_server_broadcast is declared.
func (g *LLVMGenerator) ensureWebSocketServerBroadcastDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_server_broadcast"]; exists {
		return fn
	}

	// int64_t websocket_server_broadcast(int64_t server_id, char* message)
	fn := g.module.NewFunc("websocket_server_broadcast", types.I64,
		ir.NewParam("server_id", types.I64),
		ir.NewParam("message", types.I8Ptr))

	g.functions["websocket_server_broadcast"] = fn

	return fn
}

// ensureWebSocketStopServerDeclaration ensures websocket_stop_server is declared.
func (g *LLVMGenerator) ensureWebSocketStopServerDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_stop_server"]; exists {
		return fn
	}

	// int64_t websocket_stop_server(int64_t server_id)
	fn := g.module.NewFunc("websocket_stop_server", types.I64,
		ir.NewParam("server_id", types.I64))

	g.functions["websocket_stop_server"] = fn

	return fn
}

// ensureWebSocketKeepAliveDeclaration ensures websocket_keep_alive is declared.
func (g *LLVMGenerator) ensureWebSocketKeepAliveDeclaration() *ir.Func {
	if fn, exists := g.functions["websocket_keep_alive"]; exists {
		return fn
	}

	// void websocket_keep_alive()
	fn := g.module.NewFunc("websocket_keep_alive", types.Void)

	g.functions["websocket_keep_alive"] = fn

	return fn
}
