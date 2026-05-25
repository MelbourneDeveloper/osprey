#!/bin/bash

# Helper script to switch between ARM64 and AMD64 development containers

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

show_usage() {
    echo "Usage: $0 [arm64|amd64|status]"
    echo ""
    echo "Commands:"
    echo "  arm64   - Switch to ARM64 development container (Apple Silicon)"
    echo "  amd64   - Switch to AMD64 development container (GitHub Codespaces/x64)"
    echo "  status  - Show current architecture configuration"
    echo ""
    echo "Examples:"
    echo "  $0 arm64    # Switch to ARM64 for local Apple Silicon development"
    echo "  $0 amd64    # Switch to AMD64 for GitHub Codespaces"
    echo "  $0 status   # Check current configuration"
}

check_status() {
    echo "🔍 Checking current devcontainer configuration..."
    
    if [ -f "$SCRIPT_DIR/devcontainer.json" ]; then
        if grep -q "arm64" "$SCRIPT_DIR/devcontainer.json"; then
            echo "✅ Current configuration: ARM64 (Apple Silicon)"
        elif grep -q "amd64" "$SCRIPT_DIR/devcontainer.json"; then
            echo "✅ Current configuration: AMD64 (GitHub Codespaces/x64)"
        else
            echo "❓ Current configuration: Unknown"
        fi
    else
        echo "❌ No devcontainer.json found in .devcontainer/"
    fi
    
    if [ -f "$ROOT_DIR/compiler/.devcontainer/devcontainer.json" ]; then
        echo "📁 Compiler devcontainer: AMD64 (GitHub Codespaces default)"
    fi
}

switch_to_arm64() {
    echo "🔄 Switching to ARM64 configuration..."
    
    if [ ! -f "$SCRIPT_DIR/devcontainer-arm64.json" ]; then
        echo "❌ ARM64 configuration file not found: $SCRIPT_DIR/devcontainer-arm64.json"
        exit 1
    fi
    
    cp "$SCRIPT_DIR/devcontainer-arm64.json" "$SCRIPT_DIR/devcontainer.json"
    echo "✅ Switched to ARM64 configuration"
    echo "💡 To apply changes:"
    echo "   1. Close VS Code"
    echo "   2. Reopen the project"
    echo "   3. Select 'Reopen in Container' when prompted"
    echo ""
    echo "   Or use Docker Compose directly:"
    echo "   cd .devcontainer && docker-compose --profile arm64 up -d"
}

switch_to_amd64() {
    echo "🔄 Switching to AMD64 configuration..."
    
    if [ ! -f "$SCRIPT_DIR/devcontainer-amd64.json" ]; then
        echo "❌ AMD64 configuration file not found: $SCRIPT_DIR/devcontainer-amd64.json"
        exit 1
    fi
    
    cp "$SCRIPT_DIR/devcontainer-amd64.json" "$SCRIPT_DIR/devcontainer.json"
    echo "✅ Switched to AMD64 configuration"
    echo "💡 To apply changes:"
    echo "   1. Close VS Code"
    echo "   2. Reopen the project"
    echo "   3. Select 'Reopen in Container' when prompted"
    echo ""
    echo "   Or use Docker Compose directly:"
    echo "   cd .devcontainer && docker-compose --profile amd64 up -d"
}

# Main script logic
case "${1:-}" in
    "arm64")
        switch_to_arm64
        ;;
    "amd64")
        switch_to_amd64
        ;;
    "status")
        check_status
        ;;
    "")
        echo "❌ No command specified"
        echo ""
        show_usage
        exit 1
        ;;
    *)
        echo "❌ Unknown command: $1"
        echo ""
        show_usage
        exit 1
        ;;
esac
