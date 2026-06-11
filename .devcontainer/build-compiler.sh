#!/usr/bin/env bash
# This script builds the Osprey compiler (C runtime archives + Rust workspace)

set -e

cd /workspace

echo "🔧 Building Osprey (C runtime + Rust workspace + VSCode extension)..."
make build

echo "✅ Osprey compiler built successfully!"
echo ""
echo "The compiler is available at: ./target/release/osprey"
echo "To install it globally, run: make install"
echo "To run tests, run: make test"
