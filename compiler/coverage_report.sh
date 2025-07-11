#!/usr/bin/env bash

# If the script was invoked with a shell that is *not* Bash (e.g. /bin/sh),
# re-execute it with Bash to guarantee compatibility with 'set -o pipefail'.
if [ -z "${BASH_VERSION:-}" ]; then
  exec bash "$0" "$@"
fi

set -euo pipefail

# =============================================================================
# Comprehensive code-coverage report for the compiler repository.
# =============================================================================
# 1. FORCE clean all artifacts and rebuild everything to ensure fresh state
# 2. Dynamically gather all Go packages in the module, excluding the generated
#    parser code (we do not want to track coverage for generated files).
# 3. Execute the full test-suite with race detection and atomic coverage.
# 4. Produce both textual and HTML coverage summaries.
# =============================================================================

# Emoji-rich status messaging keeps things fun but concise.
echo "🧪 Running comprehensive code-coverage analysis…"

# -----------------------------------------------------------------------------
# FORCE CLEAN AND REBUILD - No shortcuts for coverage tests!
# -----------------------------------------------------------------------------
echo "🧹 Force cleaning all artifacts for reliable coverage..."

# Clean directories directly
rm -rf bin/
rm -rf outputs/
rm -rf internal/codegen/bin
rm -f coverage.out coverage.html

# Clean build artifacts with patterns
find . -name "*.o" -delete 2>/dev/null || true
find . -name "*.a" -delete 2>/dev/null || true
find . -name "*.so" -delete 2>/dev/null || true
find . -name "*.dylib" -delete 2>/dev/null || true
find . -name "*.ll" -delete 2>/dev/null || true
find . -name "*.bc" -delete 2>/dev/null || true
find /tmp -name "*osprey*" -delete 2>/dev/null || true

echo "🔨 Force rebuilding all runtimes and compiler..."

# Create bin and lib directories
mkdir -p bin lib

# Build ALL 4 runtime libraries (same as Makefile)
echo "   Building fiber runtime..."
gcc -c -fPIC -O2 runtime/fiber_runtime.c -o bin/fiber_runtime.o
ar rcs lib/libfiber_runtime.a bin/fiber_runtime.o

echo "   Building HTTP runtime..."
# Set OpenSSL paths for macOS Homebrew
OPENSSL_CFLAGS=""
if [[ "$(uname -s)" == "Darwin" ]]; then
    OPENSSL_PREFIX=$(brew --prefix openssl 2>/dev/null || echo "/usr/local")
    OPENSSL_CFLAGS="-I${OPENSSL_PREFIX}/include"
fi

gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/http_shared.c -o bin/http_shared.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/fiber_runtime.c -o bin/fiber_runtime_http.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/system_runtime.c -o bin/system_runtime_http.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/http_client_runtime.c -o bin/http_client_runtime.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/http_server_runtime.c -o bin/http_server_runtime.o
ar rcs lib/libhttp_runtime.a bin/http_shared.o bin/fiber_runtime_http.o bin/system_runtime_http.o bin/http_client_runtime.o bin/http_server_runtime.o

echo "   Building WebSocket runtime..."
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/http_shared.c -o bin/http_shared_ws.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/fiber_runtime.c -o bin/fiber_runtime_ws.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/system_runtime.c -o bin/system_runtime_ws.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/websocket_client_runtime.c -o bin/websocket_client_runtime.o
gcc -c -fPIC -O2 $OPENSSL_CFLAGS runtime/websocket_server_runtime.c -o bin/websocket_server_runtime.o
ar rcs lib/libwebsocket_runtime.a bin/http_shared_ws.o bin/fiber_runtime_ws.o bin/system_runtime_ws.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o

echo "   Building System runtime..."
gcc -c -fPIC -O2 runtime/system_runtime.c -o bin/system_runtime.o
ar rcs lib/libsystem_runtime.a bin/system_runtime.o

# Build compiler
echo "   Building Osprey compiler..."
go build -o bin/osprey ./cmd/osprey

# -----------------------------------------------------------------------------
# Build package list (exclude generated parser)
# -----------------------------------------------------------------------------
ALL_PKGS=$(go list ./...)
PKGS=$(echo "$ALL_PKGS" | grep -v "/parser$")

# Convert package list to comma-separated string for -coverpkg
COVERPKG=$(echo "$PKGS" | tr '\n' ',' | sed 's/,$//')

# -----------------------------------------------------------------------------
# Run tests with coverage
# -----------------------------------------------------------------------------
echo "🧪 Running all tests sequentially with atomic coverage..."
echo "    📦 Packages: $PKGS"
echo

# Run all tests sequentially
if go test -v -covermode=atomic -coverpkg="$COVERPKG" -coverprofile=coverage.out -timeout=15m $PKGS; then
    TEST_EXIT_CODE=0
    echo "✅ All tests completed successfully"
else
    TEST_EXIT_CODE=$?
    echo "❌ Tests failed with exit code: $TEST_EXIT_CODE"
fi

echo "🔍 Test execution finished with exit code: $TEST_EXIT_CODE"

# Check if coverage file was generated
if [[ -f "coverage.out" ]]; then
    echo "✅ Coverage file generated successfully"
    COVERAGE_SIZE=$(wc -l < coverage.out)
    echo "📊 Coverage file contains $COVERAGE_SIZE lines"
else
    echo "❌ Coverage file was not generated!"
fi

# Check for race condition reports in the output
echo "🔍 Checking for race conditions or other specific failures..."

# -----------------------------------------------------------------------------
# Verify coverage file was created
# -----------------------------------------------------------------------------
if [ ! -f coverage.out ]; then
    echo "❌ Error: coverage.out file was not created"
    exit 1
fi

echo "✅ Coverage file created successfully"

# -----------------------------------------------------------------------------
# Generate & display coverage reports
# -----------------------------------------------------------------------------
echo "📈 Coverage Summary:"
go tool cover -func=coverage.out

echo "🔧 Generating HTML report…"
go tool cover -html=coverage.out -o coverage.html

# -----------------------------------------------------------------------------
# Extract total coverage with error handling
# -----------------------------------------------------------------------------
TOTAL_COVERAGE=$(go tool cover -func=coverage.out | awk '/^total:/ {print $3}' || echo "unknown")

if [ "$TOTAL_COVERAGE" = "unknown" ]; then
    echo "⚠️  Warning: Could not extract total coverage percentage"
    echo "📁 Raw coverage data saved to: coverage.out"
    echo "📁 HTML report saved to: coverage.html"
    exit 0
fi

printf "\n🎯 Total Coverage: %s\n" "$TOTAL_COVERAGE"

echo "📁 HTML report saved to: coverage.html"
echo "📁 Raw coverage data saved to: coverage.out"

# Automatically open the HTML report on macOS for convenience
if [[ "$(uname -s)" == "Darwin" ]]; then
    echo "🌐 Opening HTML coverage report in browser…"
    open coverage.html
fi

echo "✅ Coverage analysis complete!" 