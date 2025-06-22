package integration

// DO NOT EVER SKIP TESTS!!!!

import (
	"testing"
)

// TestHTTPExamples tests HTTP and WebSocket examples in the examples/tested/http directory.
func TestHTTPExamples(t *testing.T) {
	checkLLVMTools(t)

	examplesDir := "../../examples/tested/http"
	runTestExamplesRecursive(t, examplesDir, getHTTPExpectedOutputs())
}

func getHTTPExpectedOutputs() map[string]string {
	return map[string]string{
		"http_create_client.osp": "Creating client\n1\n",
		"http_client_example.osp": "=== HTTP Client Test ===\n" +
			"Testing HTTP client creation...\n" +
			"SUCCESS: Client created, making HTTP request...\n" +
			"HTTP GET request completed\n" +
			"Cleaning up...\n" +
			"HTTP test complete\n",
		"http_server_example.osp": "=== Comprehensive HTTP API Test ===\n" +
			"Creating HTTP server on port 8080...\n" +
			"Server created with ID: 1\n" +
			"Starting server listener with callback handler...\n" +
			"HTTP server listening on 127.0.0.1:8080\n" +
			"Server listening on http://127.0.0.1:8080\n" +
			"Creating HTTP client...\n" +
			"Client created with ID: 2\n" +
			"=== Testing HTTP Methods ===\n" +
			"GET /api/users\n" +
			"🌐 HTTP Request: GET /api/users\n" +
			"GET /api/users result: 200\n" +
			"GET /health\n" +
			"🌐 HTTP Request: GET /health\n" +
			"GET /health result: 200\n" +
			"POST /api/users\n" +
			"🌐 HTTP Request: POST /api/users\n" +
			"POST /api/users result: 201\n" +
			"POST /api/auth/login\n" +
			"🌐 HTTP Request: POST /api/auth/login\n" +
			"POST /api/auth/login result: 200\n" +
			"PUT /api/users/1\n" +
			"🌐 HTTP Request: PUT /api/users/1\n" +
			"PUT /api/users/1 result: 200\n" +
			"DELETE /api/users/1\n" +
			"🌐 HTTP Request: DELETE /api/users/1\n" +
			"DELETE /api/users/1 result: 200\n" +
			"=== Testing Error Cases ===\n" +
			"GET /nonexistent\n" +
			"🌐 HTTP Request: GET /nonexistent\n" +
			"GET /nonexistent result: 404\n" +
			"=== Server Keep-Alive Test ===\n" +
			"🌐 Server now accepting external connections...\n" +
			"💡 Test with: curl http://127.0.0.1:8080/api/users\n" +
			"💡 Test with: curl http://127.0.0.1:8080/health\n" +
			"⏰ Server will stay alive for 3 seconds...\n" +
			"🛑 Initiating graceful server shutdown...\n" +
			"Server stopped with result: 0\n" +
			"=== HTTP API Test Complete ===\n",
		"websocket_example.osp": "=== WebSocket Communication Test ===\n" +
			"Connecting to WebSocket server...\n" +
			"WebSocket connected with ID: 1\n" +
			"=== WebSocket Messages ===\n" +
			"Sending: Hello WebSocket!\n" +
			"Send result: 0\n" +
			"Sending: Real-time message 1\n" +
			"Send result: 0\n" +
			"Sending: Real-time message 2\n" +
			"Send result: 0\n" +
			"Sending JSON data...\n" +
			"JSON send result: 0\n" +
			"Sending status update...\n" +
			"Status send result: 0\n" +
			"=== WebSocket Chat Simulation ===\n" +
			"Sending chat message 1...\n" +
			"Chat 1 result: 0\n" +
			"Sending chat message 2...\n" +
			"Chat 2 result: 0\n" +
			"Sending typing indicator...\n" +
			"Typing result: 0\n" +
			"Closing WebSocket connection...\n" +
			"WebSocket closed with result: 0\n" +
			"=== WebSocket Test Complete ===\n",
		"websocket_local_test.osp": "=== WebSocket Local Server Test ===\n" +
			"Creating local WebSocket server...\n" +
			"Server listening with result: 0\n" +
			"Connecting to local WebSocket server...\n" +
			"WebSocket connected with ID: 1\n" +
			"=== Basic Text Messages ===\n" +
			"Sending: Hello from Osprey!\n" +
			"Send result: 0\n" +
			"Sending: Osprey WebSocket client test\n" +
			"Send result: 0\n" +
			"=== JSON Messages ===\n" +
			"Sending JSON chat message...\n" +
			"JSON send result: 0\n" +
			"Sending status update...\n" +
			"Status send result: 0\n" +
			"=== Rapid Message Test ===\n" +
			"Sending message 1...\n" +
			"Rapid 1 result: 0\n" +
			"Sending message 2...\n" +
			"Rapid 2 result: 0\n" +
			"Sending message 3...\n" +
			"Rapid 3 result: 0\n" +
			"=== Closing Connection ===\n" +
			"Closing WebSocket connection...\n" +
			"WebSocket closed with result: 0\n" +
			"Stopping WebSocket server...\n" +
			"Server stopped with result: 0\n" +
			"=== WebSocket Local Test Complete ===\n",
		"websocket_server_example.osp": "=== Osprey WebSocket Server ===\n" +
			"Creating WebSocket server on port 8080...\n" +
			"WebSocket server created with ID: 1\n" +
			"Starting WebSocket server...\n" +
			"WebSocket server listening with result: 0\n" +
			"🚀 WebSocket server running on ws://127.0.0.1:8080/chat\n" +
			"💡 You can now connect browsers to this server!\n" +
			"\n" +
			"🧪 Test with browser:\n" +
			"1. Open examples/websocketserver/websocket_test.html\n" +
			"2. Browser will connect to ws://localhost:8080/chat\n" +
			"3. Send messages and see real-time echo responses\n" +
			"\n" +
			"🔄 Broadcasting welcome message to all connections...\n" +
			"Broadcast result: 0\n" +
			"\n" +
			"🛑 Stopping WebSocket server...\n" +
			"WebSocket server stopped with result: 0\n" +
			"=== WebSocket Server Example Complete ===\n",
		"websocket_server_live.osp": "🚀 Starting Osprey WebSocket Server...\n" +
			"📁 Will serve WebSocket at: ws://127.0.0.1:8080/chat\n" +
			"🔌 Ready to replace Node.js websocket_server.js!\n" +
			"\n" +
			"Creating WebSocket server on port 8080...\n" +
			"✅ WebSocket server created with ID: 1\n" +
			"Starting WebSocket server...\n" +
			"✅ WebSocket server listening with result: 0\n" +
			"\n" +
			"🚀 Osprey WebSocket Server is LIVE!\n" +
			"🔌 WebSocket endpoint: ws://127.0.0.1:8080/chat\n" +
			"🌐 Connect your browsers, test apps, and clients!\n" +
			"\n" +
			"🔄 Broadcasting welcome message to all connections...\n" +
			"✅ Welcome broadcast sent with result: 0\n" +
			"\n" +
			"💡 Server Instructions:\n" +
			"1. Open examples/websocketserver/websocket_test.html in browser\n" +
			"2. Browser will connect to ws://localhost:8080/chat\n" +
			"3. Send messages and see real-time responses\n" +
			"4. Press Ctrl+C to stop the server\n" +
			"\n" +
			"🎉 Osprey WebSocket Server has replaced Node.js!\n" +
			"📦 You can now delete websocket_server.js\n" +
			"\n" +
			"⚠️  Current Status: Demo mode (exits immediately)\n" +
			"🚧 TODO: Add infinite loop to keep server alive\n" +
			"🚧 TODO: Add signal handling for graceful shutdown\n" +
			"\n" +
			"🛑 Stopping WebSocket server (demo complete)...\n" +
			"✅ WebSocket server stopped with result: 0\n" +
			"=== Osprey WebSocket Server Demo Complete ===\n",
	}
}
