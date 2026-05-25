#!/bin/bash

echo "🚀 Starting Osprey ARM64 Dev Container..."

# Start the dev container
echo "🐳 Starting dev container..."
docker compose -f .devcontainer/docker-compose.yml --profile arm64 up -d

if [ $? -eq 0 ]; then
    echo "✅ Dev container started successfully!"
    
    # Get the dynamically assigned ports
    CONTAINER_ID=$(docker ps -q --filter "name=devcontainer-osprey-dev-arm64")
    if [ ! -z "$CONTAINER_ID" ]; then
        PORT_3001=$(docker port $CONTAINER_ID 3001 | cut -d: -f2)
        PORT_8080=$(docker port $CONTAINER_ID 8080 | cut -d: -f2)
        
        if [ ! -z "$PORT_3001" ]; then
            echo "🌐 Web Compiler available at: http://localhost:$PORT_3001"
        fi
        if [ ! -z "$PORT_8080" ]; then
            echo "🔧 Development Server available at: http://localhost:$PORT_8080"
        fi
    fi
else
    echo "❌ Failed to start dev container"
    exit 1
fi
