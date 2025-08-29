#!/bin/bash

# Osprey Web Compiler API Test
# Tests the local container running on localhost:3001

echo "🧪 Testing Osprey Web Compiler API..."
echo "===================================="

# Define paths to the test files
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OSP_FILE="$SCRIPT_DIR/../compiler/examples/tested/basics/osprey_mega_showcase.osp"
EXPECTED_OUTPUT_FILE="$SCRIPT_DIR/../compiler/examples/tested/basics/osprey_mega_showcase.osp.expectedoutput"

# Check if files exist
if [ ! -f "$OSP_FILE" ]; then
    echo "❌ Error: Osprey file not found at $OSP_FILE"
    exit 1
fi

if [ ! -f "$EXPECTED_OUTPUT_FILE" ]; then
    echo "❌ Error: Expected output file not found at $EXPECTED_OUTPUT_FILE"
    exit 1
fi

# Read the Osprey code and expected output
OSP_CODE=$(cat "$OSP_FILE")
EXPECTED_OUTPUT=$(cat "$EXPECTED_OUTPUT_FILE")

echo "📄 Loaded Osprey code from: $OSP_FILE"
echo "📄 Loaded expected output from: $EXPECTED_OUTPUT_FILE"

# Test the local API
echo "Testing local API at http://localhost:3001/api/run"
RESPONSE=$(curl -s -X POST http://localhost:3001/api/run \
  -H 'Content-Type: application/json' \
  -d "{\"code\":$(echo "$OSP_CODE" | jq -Rs .)}")

echo "Response received from API"

# Extract the program output from the JSON response
PROGRAM_OUTPUT=$(echo "$RESPONSE" | jq -r '.programOutput // empty')

if [ $? -ne 0 ]; then
    echo "❌ Test FAILED: Failed to parse JSON response"
    echo "Response: $RESPONSE"
    exit 1
fi

# Verify the response contains expected structure
if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo "✅ API returned success: true"
else
    echo "❌ Test FAILED: API did not return success: true"
    echo "Response: $RESPONSE"
    exit 1
fi

# Compare the program output with expected output
if [ "$PROGRAM_OUTPUT" = "$EXPECTED_OUTPUT" ]; then
    echo "✅ Test PASSED: Program output matches expected output exactly"
    exit 0
else
    echo "❌ Test FAILED: Program output does not match expected output"
    echo ""
    echo "Expected output:"
    echo "=================="
    echo "$EXPECTED_OUTPUT"
    echo ""
    echo "Actual output:"
    echo "=============="
    echo "$PROGRAM_OUTPUT"
    echo ""
    echo "JSON Response:"
    echo "=============="
    echo "$RESPONSE"
    exit 1
fi