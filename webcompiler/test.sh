#!/bin/bash

# Osprey Web Compiler API Test
# Tests the local container running on localhost:3001

echo "🧪 Testing Osprey Web Compiler API..."
echo "===================================="

# Test the local API
echo "Testing local API at http://localhost:3001/api/run"
RESPONSE=$(curl -s -X POST http://localhost:3001/api/run \
  -H 'Content-Type: application/json' \
  -d '{"code":"print(\"Hello World\")", "language":"osprey"}')

echo "Response: $RESPONSE"

# Verify the response contains expected output
if echo "$RESPONSE" | grep -q "Hello World"; then
    echo "✅ Test PASSED: API returned expected output"
    exit 0
else
    echo "❌ Test FAILED: API did not return expected output"
    echo "Expected: Hello World"
    echo "Got: $RESPONSE"
    exit 1
fi