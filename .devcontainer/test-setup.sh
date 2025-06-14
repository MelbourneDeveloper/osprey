#!/usr/bin/env bash
# This script tests that all development tools are working correctly

set -e

echo "🧪 Testing Osprey development environment setup..."
echo ""

# Test Go
echo "🔍 Testing Go..."
go version
cd /workspaces/osprey/compiler
go mod tidy
echo "✅ Go is working!"
echo ""

# Test ANTLR
echo "🔍 Testing ANTLR..."
antlr -version
echo "✅ ANTLR is working!"
echo ""

# Test LLVM
echo "🔍 Testing LLVM..."
llc --version | head -1
echo "✅ LLVM is working!"
echo ""

# Test Node.js and npm
echo "🔍 Testing Node.js and npm..."
node --version
npm --version
echo "✅ Node.js and npm are working!"
echo ""

# Test Rust
echo "🔍 Testing Rust..."
rustc --version
echo "✅ Rust is working!"
echo ""

# Test compiler build
echo "🔍 Testing compiler build..."
cd /workspaces/osprey/compiler
make build
echo "✅ Compiler builds successfully!"
echo ""

# Test VS Code extension setup
echo "🔍 Testing VS Code extension setup..."
cd /workspaces/osprey/vscode-extension
npm install --silent
npm run compile
echo "✅ VS Code extension compiles successfully!"
echo ""

echo "🎉 All tests passed! Your development environment is ready!"
echo ""
echo "📋 Summary:"
echo "- ✅ Go $(go version | cut -d' ' -f3)"
echo "- ✅ ANTLR $(antlr -version | head -1)"
echo "- ✅ LLVM $(llc --version | head -1 | cut -d' ' -f3)"
echo "- ✅ Node.js $(node --version)"
echo "- ✅ npm $(npm --version)"
echo "- ✅ Rust $(rustc --version | cut -d' ' -f2)"
echo "- ✅ Osprey compiler builds"
echo "- ✅ VS Code extension compiles" 