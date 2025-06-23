#!/bin/bash

echo "🚀 Starting Osprey Web Compiler (Docker)..."
echo "================================================"
echo ""
echo "This will start the Dockerized Osprey web compiler:"
echo "- API/WebSocket server on port 3001"
echo "- LSP bridge for language features"
echo "- Compile/Run endpoints with sandbox security"
echo "- Sandbox mode: HTTP, WebSocket, file system, and FFI disabled"
echo ""
echo "Access the service at: http://localhost:3001"
echo ""
echo "================================================"

# Navigate to webcompiler directory
cd "$(dirname "$0")"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is required but not installed!"
    echo "Run: curl -fsSL https://get.docker.com | sh"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "❌ Docker Compose is required but not installed!"
    echo "Run: sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-\$(uname -s)-\$(uname -m)\" -o /usr/local/bin/docker-compose"
    echo "     sudo chmod +x /usr/local/bin/docker-compose"
    exit 1
fi

# Stop any existing containers
echo "🛑 Stopping any existing containers..."
if command -v docker-compose &> /dev/null; then
    docker-compose down 2>/dev/null || true
else
    docker compose down 2>/dev/null || true
fi

# Build and start the container
echo ""
echo "🔨 Building and starting Docker container..."
echo "   This may take a few minutes on first run..."
if command -v docker-compose &> /dev/null; then
    docker-compose up --build
else
    docker compose up --build
fi

# Note: docker-compose up runs in foreground, so this only executes on Ctrl+C
echo ""
echo "🛑 Shutting down containers..."
if command -v docker-compose &> /dev/null; then
    docker-compose down
else
    docker compose down
fi 