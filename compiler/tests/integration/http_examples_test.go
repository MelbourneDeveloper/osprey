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
		"http_advanced_example.osp": "=== Advanced HTTP Test ===\n" +
			"Creating HTTP server on port 8080...\n" +
			"Server created with ID: 1\n" +
			"Starting server listener...\n" +
			"Server listening on http://127.0.0.1:8080\n" +
			"=== Creating Multiple Clients ===\n" +
			"Creating client 1...\n" +
			"Client 1 created with ID: 2\n" +
			"Creating client 2...\n" +
			"Client 2 created with ID: 3\n" +
			"Creating client 3...\n" +
			"Client 3 created with ID: 4\n" +
			"=== Concurrent Requests ===\n" +
			"Client 1: GET /api/users\n" +
			"Client 1 GET result: -7\n" +
			"Client 2: POST /api/posts\n" +
			"Client 2 POST result: -7\n" +
			"Client 3: GET /api/health\n" +
			"Client 3 health check: -7\n" +
			"=== API Versioning ===\n" +
			"Client 1: GET /v1/users\n" +
			"v1 API result: -7\n" +
			"Client 2: GET /v2/users\n" +
			"v2 API result: -7\n" +
			"=== Content Types ===\n" +
			"Client 1: POST /api/upload (XML)\n" +
			"XML POST result: -7\n" +
			"Client 2: PUT /api/config (YAML)\n" +
			"YAML PUT result: -7\n" +
			"Client 3: POST /api/data (Form)\n" +
			"Form POST result: -7\n" +
			"=== Authentication ===\n" +
			"Client 1: POST /auth/login\n" +
			"Login result: -5\n" +
			"Client 2: GET /protected (with token)\n" +
			"Protected GET result: -5\n" +
			"Client 3: DELETE /auth/logout\n" +
			"Logout result: -5\n" +
			"=== Error Scenarios ===\n" +
			"Client 1: GET /nonexistent\n" +
			"404 test result: -5\n" +
			"Client 2: POST /api/invalid (bad JSON)\n" +
			"Bad JSON result: -5\n" +
			"Stopping server...\n" +
			"Server stopped with result: 0\n" +
			"=== Advanced HTTP Test Complete ===\n",
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
			"Starting server listener...\n" +
			"Server listening on http://127.0.0.1:8080\n" +
			"Creating HTTP client...\n" +
			"Client created with ID: 2\n" +
			"=== Testing HTTP Methods ===\n" +
			"GET /api/users\n" +
			"GET status: -5\n" +
			"POST /api/users\n" +
			"POST status: -5\n" +
			"PUT /api/users/1\n" +
			"PUT status: -5\n" +
			"DELETE /api/users/1\n" +
			"DELETE status: -5\n" +
			"=== Testing Additional Endpoints ===\n" +
			"GET /health\n" +
			"Health check status: -5\n" +
			"POST /api/auth/login\n" +
			"Login status: -5\n" +
			"Stopping server...\n" +
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
		"websocket_server_live.osp": "",
	}
}
