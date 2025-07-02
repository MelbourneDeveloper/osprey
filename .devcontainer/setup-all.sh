#!/usr/bin/env bash
# This script runs all setup scripts for the Osprey development environment

set -e

echo "🚀 Setting up complete Osprey development environment..."
echo ""

# Make all scripts executable
chmod +x /workspaces/osprey/.devcontainer/*.sh

# Build the compiler
echo "📦 Building compiler..."
/workspaces/osprey/.devcontainer/build-compiler.sh
echo ""

# Setup VS Code extension
echo "📦 Setting up VS Code extension..."
/workspaces/osprey/.devcontainer/setup-vscode-extension.sh
echo ""

# Test everything
echo "🧪 Testing setup..."
/workspaces/osprey/.devcontainer/test-setup.sh

echo ""
echo "✅ Complete development environment setup finished!"
echo ""
echo "🎯 You can now:"
echo "- Build the compiler: cd compiler && make build"
echo "- Run compiler tests: cd compiler && make test"
echo "- Develop VS Code extension: Press F5 in VS Code"
echo "- Package extension: cd vscode-extension && npm run package" 