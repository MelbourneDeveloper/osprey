#!/bin/bash

echo "🔧 CMake C Test Integration Verification"
echo "========================================="
echo "⚠️  This script must be run inside the VS Code Dev Container"
echo "   Use Command Palette → 'Dev Containers: Rebuild Container'"
echo ""

# Navigate to runtime directory
cd "$(dirname "$0")"

# Clean and create build directory
echo "📁 Setting up build directory..."
rm -rf build
mkdir -p build
cd build

# Configure with CMake
echo "⚙️  Configuring CMake..."
cmake .. -DCMAKE_BUILD_TYPE=Debug

if [ $? -ne 0 ]; then
    echo "❌ CMake configuration failed!"
    exit 1
fi

echo "✅ CMake configuration successful!"

# Build the tests
echo "🔨 Building C runtime tests..."
make -j$(nproc)

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful!"

# List available tests
echo "📋 Available CTest tests:"
ctest --show-only=json-v1 | jq -r '.tests[].name' 2>/dev/null || ctest -N

# Run tests with CTest
echo "🧪 Running CTest..."
ctest --verbose

if [ $? -eq 0 ]; then
    echo "✅ All C runtime tests passed!"
    echo ""
    echo "🎉 VS Code Integration Ready!"
    echo "   - Open this project in VS Code with Dev Container"
    echo "   - Go to Test Explorer (flask icon in sidebar)"
    echo "   - You should see 'SystemRuntimeTests' and 'FiberRuntimeTests'"
    echo "   - Click the play button to run individual tests"
    echo "   - Use debug button to debug tests with breakpoints"
else
    echo "❌ Some tests failed!"
    exit 1
fi 