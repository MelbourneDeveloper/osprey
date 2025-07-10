#!/usr/bin/env bash
# This script builds the Osprey compiler

set -e

cd /workspace/compiler

echo "📦 Installing Go dependencies..."
go mod tidy

echo "🔧 Building runtime libraries..."
make fiber-runtime http-runtime

echo "🔧 Creating symlinks for tests..."
cd internal/codegen && ln -sf ../../bin bin && cd ../..

echo "🔧 Building Osprey compiler..."
make build

echo "✅ Osprey compiler built successfully!"
echo ""
echo "The compiler is available at: ./bin/osprey"
echo "To install it globally, run: make install"
echo "To run tests, run: make test"
