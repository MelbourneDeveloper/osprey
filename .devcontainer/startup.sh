#!/bin/bash

echo "Starting Osprey dev container..."

# Kill any processes using port 3001 on the host
echo "Checking for processes using port 3001..."
docker ps -q --filter 'publish=3001' | xargs -r docker stop || true

echo "Port 3001 cleared, starting container services..."

# Keep container running
sleep infinity
