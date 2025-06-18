#!/bin/bash

# Simple Osprey Web Compiler API Test
# Assumes the service is running on http://localhost:3001

echo "🧪 Testing Osprey Web Compiler API..."
echo "===================================="

curl -X POST https://osprey-web-compiler-gateway.mail-bff.workers.dev/api/run -H 'Content-Type: application/json' -d '{"code":"print(\"Hello World\")", "language":"osprey"}'