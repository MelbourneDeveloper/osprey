#!/bin/bash

echo "🚀 Starting Osprey Web Compiler (Development Mode)..."
echo "===================================================="
echo ""
echo "This will start the Osprey web compiler in development mode:"
echo "- API/WebSocket server on port 3001"
echo "- LSP bridge for language features"
echo "- Compile/Run endpoints with sandbox security"
echo "- Sandbox mode: HTTP, WebSocket, file system, and FFI disabled"
echo ""
echo "Access the service at: http://localhost:3001"
echo ""
echo "===================================================="

set -e

# Navigate to webcompiler directory
cd "$(dirname "$0")"

# Find the workspace root
ROOT_DIR="$(cd .. && pwd)"

# Check if dependencies are installed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm ci
fi

# Build the Osprey compiler if not already built
if [ ! -f "$ROOT_DIR/bin/osprey" ]; then
    echo "🔨 Building Osprey compiler..."
    cd "$ROOT_DIR"
    make fiber-runtime
    go build -o bin/osprey ./cmd/osprey
    cd "$ROOT_DIR/webcompiler"
fi

# Build the LSP server if not already built
if [ ! -f "$ROOT_DIR/vscode-extension/server/out/src/server.js" ]; then
    echo "🔧 Building VSCode extension LSP server..."
    cd "$ROOT_DIR/vscode-extension/server"
    if [ ! -d "node_modules" ]; then
        echo "📦 Installing VSCode extension dependencies..."
        npm install
    fi
    echo "🔨 Compiling TypeScript LSP server..."
    npm run compile
    cd "$ROOT_DIR/webcompiler"
fi

# Start the web compiler
echo "🌐 Starting web compiler server..."
npm run server 