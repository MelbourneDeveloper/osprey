#!/bin/bash
set -e

echo "🔧 Running post-creation setup for Osprey dev environment..."

# Install Claude Code and MCP servers
echo "📦 Installing Claude Code and MCP servers..."
sudo npm install -g \
      @modelcontextprotocol/server-filesystem \
      @modelcontextprotocol/server-memory \
      @modelcontextprotocol/server-everything \
      @anthropic-ai/claude-code \
      mcp-smart-crawler \
      @playwright/mcp

echo "📦 Installing GitHub CLI..."
curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg
sudo chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null
sudo apt-get update
sudo apt-get install gh -y

# Set up Go workspace and dependencies  
echo "📦 Setting up Go workspace..."
cd /workspace/compiler

# Fix Go module cache permissions
echo "🔧 Fixing Go module cache permissions..."
sudo chown -R vscode:vscode /go/pkg/mod || true

# Set proper GOPROXY and GOSUMDB for better module resolution
export GOPROXY=https://proxy.golang.org,direct
export GOSUMDB=sum.golang.org

go mod download
go mod tidy

# Build the Osprey compiler
echo "🔨 Building Osprey compiler..."
make clean
make build || echo "⚠️ Initial build failed - this is expected on first setup"

# Install compiler globally
echo "📦 Installing Osprey compiler globally..."
make install || echo "⚠️ Installation failed - may need manual intervention"

# Set up VSCode extension dependencies
if [ -d "/workspace/vscode-extension" ]; then
  echo "📦 Setting up VSCode extension dependencies..."
  cd /workspace/vscode-extension
  npm install
fi

# Set up website dependencies
if [ -d "/workspace/website" ]; then
  echo "📦 Setting up website dependencies..."
  cd /workspace/website
  npm install
fi

# Set up webcompiler dependencies
if [ -d "/workspace/webcompiler" ]; then
  echo "📦 Setting up webcompiler dependencies..."
  cd /workspace/webcompiler
  npm install
fi

# Return to workspace root
cd /workspace

echo "🦀 Setting up Rust for vscode user..."
# Ensure Rust is properly set up for vscode user
if [ ! -f ~/.cargo/env ]; then
  echo "🔧 Rust not found for vscode user, installing..."
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
fi

# Source Rust environment and set default toolchain
source ~/.cargo/env
rustup default stable
rustup component add clippy rustfmt

echo "🎯 Verifying installation..."
go version
node --version
npm --version
rustc --version || echo "⚠️ Rust not properly installed"
cargo --version || echo "⚠️ Cargo not properly installed"
claude --version || echo "⚠️ Claude Code not installed"

echo "🎉 Post-creation setup complete!"
echo ""
echo "📝 Available commands:"
echo "  make build      - Build the Osprey compiler"
echo "  make test       - Run all tests"
echo "  make lint       - Run Go linter"
echo "  make install    - Install compiler globally"
echo "  claude-code     - Run Claude Code"
echo ""
echo "🚀 Ready to develop Osprey!"