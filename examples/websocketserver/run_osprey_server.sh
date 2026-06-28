#!/bin/bash

echo "Starting Osprey WebSocket Server..."

# Kill any existing process on port 54321
echo "Killing any existing process on port 54321..."
lsof -ti:54321 | xargs kill -9 2>/dev/null || true

# Always build the compiler first
echo "Building Osprey compiler..."
cd ../.. && make build && cd examples/websocketserver

# Run the server and show output immediately
echo "Running Osprey WebSocket Server..."
echo "=================================="
../../bin/osprey osprey_websocket_server.osp --run 