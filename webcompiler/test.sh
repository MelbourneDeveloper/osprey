#!/bin/bash

# Osprey Web Compiler API Test
# Tests the local container running on localhost:3001

echo "🧪 Testing Osprey Web Compiler API..."
echo "===================================="

# Test the local API
echo "Testing local API at http://localhost:3001/api/run"
RESPONSE=$(curl -s -X POST http://localhost:3001/api/run \
  -H 'Content-Type: application/json' \
  -d '{"code":"print(\"Testing API Response\")"}')

echo "Response: $RESPONSE"

# Verify the response contains expected structure and output
if echo "$RESPONSE" | grep -q "Testing API Response" && \
   echo "$RESPONSE" | grep -q "\"programOutput\":" && \
   echo "$RESPONSE" | grep -q "\"success\":true"; then
    echo "✅ Test PASSED: API returned expected response format with correct output"
    exit 0
else
    echo "❌ Test FAILED: API did not return expected response format"
    echo "Expected: JSON with success:true, programOutput field, and 'Testing API Response' in programOutput"
    echo "Got: $RESPONSE"
    exit 1
fi