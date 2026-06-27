# Osprey WebSocket Implementation

WebSocket implementation for Osprey with browser testing tools.

## Quick Start

```bash
# Install dependencies
npm install

# Start WebSocket test server
node websocket_server.js

# Open browser test interface
# Visit http://localhost:8080

# Test Osprey client
./bin/osprey examples/tested/websocket_example.osp --run
```

## Files

- `websocket_server.js` - Node.js WebSocket server
- `websocket_test.html` - Browser test interface
- `examples/tested/websocket_*.osp` - Osprey WebSocket examples
- `runtime/http_runtime.c` - C implementation

## WebSocket Functions

```osprey
websocketConnect(url: String, messageHandler: String) -> Int
websocketSend(wsId: Int, message: String) -> Int
websocketClose(wsId: Int) -> Int
```

## Return Codes

- Positive: Success (WebSocket ID for connect, 0 for operations)
- -1: Invalid parameters
- -2: Connection not established
- -3 to -9: Various connection/handshake errors

## Status

✅ **Working**: Browser client, Node.js server, basic Osprey connectivity
⚠️ **Limited**: URL parsing, error handling, no SSL support 