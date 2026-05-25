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

# Build the Osprey compiler (without linting on first run to ensure Rust utils are built)
echo "🔨 Building Osprey compiler..."
make clean
# Build Rust utilities first
make rust-interop || echo "⚠️ Rust interop build failed"
# Then do the full build
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

echo "🦀 Setting up Rust environment..."
# Rust is already installed in /opt/rust, just ensure it's in PATH
export PATH="/opt/rust/bin:$PATH"
export CARGO_HOME="/opt/rust"
export RUSTUP_HOME="/opt/rust"

# Verify Rust installation
echo "🔧 Verifying Rust installation..."
which rustc || echo "⚠️ rustc not found in PATH"
which cargo || echo "⚠️ cargo not found in PATH"

# Build Rust utilities library
echo "🦀 Building Rust utilities library..."
cd /workspace/compiler
if [ -d "examples/rust_integration" ]; then
  cd examples/rust_integration
  if cargo build --release; then
    mkdir -p /workspace/compiler/lib
    cp target/release/libosprey_math_utils.a /workspace/compiler/lib/librust_utils.a
    # Also copy to bin directory as fallback
    mkdir -p /workspace/compiler/bin
    cp target/release/libosprey_math_utils.a /workspace/compiler/bin/librust_utils.a
    echo "✅ Rust utilities library built and installed successfully"
  else
    echo "⚠️ Failed to build Rust utilities library"
  fi
  cd /workspace/compiler
fi

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