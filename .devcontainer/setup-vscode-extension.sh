#!/usr/bin/env bash
# This script sets up the vscode extension development

set -e

cd /workspaces/osprey/vscode-extension

# Install dependencies for client
echo "📦 Installing VS Code extension dependencies..."
npm install

# Install dependencies for server
echo "📦 Installing language server dependencies..."
cd server
npm install
cd ..

# Compile the extension
echo "🔧 Compiling the extension..."
npm run compile

echo "✅ VS Code extension setup complete!"
echo ""
echo "To test the extension, press F5 in VS Code with 'Extension Development Host'"
echo "To package the extension, run 'npm run package' in the vscode-extension directory"
