#!/bin/bash

set -euo pipefail

echo "🚀 Deploying Osprey Website"

# Check if we're in the website directory
if [ ! -f "package.json" ]; then
    echo "❌ Error: Must be run from the website directory"
    exit 1
fi

# Check if the compiler exists
COMPILER_PATH="../Osprey Compiler/bin/osprey"
if [ ! -f "$COMPILER_PATH" ]; then
    echo "❌ Error: Osprey compiler not found at $COMPILER_PATH"
    echo "Please build the compiler first from the 'Osprey Compiler' directory:"
    echo "  cd '../Osprey Compiler' && make build"
    exit 1
fi

echo "✅ Found Osprey compiler at $COMPILER_PATH"

# Install dependencies if needed
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# Generate documentation with the compiler
echo "📚 Generating documentation..."

# Create docs directory if it doesn't exist
mkdir -p src/docs/generated

# Generate spec documentation if the compiler supports it
if "$COMPILER_PATH" --help | grep -q "generate-docs\|docs\|spec"; then
    echo "🔧 Running compiler documentation generation..."
    "$COMPILER_PATH" generate-docs --output src/docs/generated/ || echo "⚠️ Docs generation not available yet"
fi

# Build the website
echo "🏗️ Building website..."
npm run build

echo "✅ Website built successfully!"
echo "📁 Output directory: _site/"

# If we're in GitHub Actions, the deployment will be handled by the workflow
if [ "${GITHUB_ACTIONS:-false}" = "true" ]; then
    echo "🔄 Running in GitHub Actions - deployment will be handled by workflow"
else
    echo "💡 To serve locally, run: npm run start"
    echo "💡 Website is ready for deployment from the _site/ directory"
fi 