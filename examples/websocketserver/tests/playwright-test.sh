#!/bin/bash

echo "Running Playwright WebSocket tests..."

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "Installing dependencies..."
    npm install
fi

# Run playwright tests
npx playwright test 