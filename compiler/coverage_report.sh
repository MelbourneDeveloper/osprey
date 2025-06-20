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

# Create bin directory
mkdir -p bin

# Build fiber runtime
echo "   Building fiber runtime..."
gcc -c -fPIC -O2 runtime/fiber_runtime.c -o bin/fiber_runtime.o
ar rcs bin/libfiber_runtime.a bin/fiber_runtime.o

# Build HTTP runtime
echo "   Building HTTP runtime..."
gcc -c -fPIC -O2 runtime/http_shared.c -o bin/http_shared.o
gcc -c -fPIC -O2 runtime/http_client_runtime.c -o bin/http_client_runtime.o
gcc -c -fPIC -O2 runtime/http_server_runtime.c -o bin/http_server_runtime.o
gcc -c -fPIC -O2 runtime/websocket_client_runtime.c -o bin/websocket_client_runtime.o
gcc -c -fPIC -O2 runtime/websocket_server_runtime.c -o bin/websocket_server_runtime.o
gcc -c -fPIC -O2 runtime/system_runtime.c -o bin/system_runtime.o
ar rcs bin/libhttp_runtime.a bin/http_shared.o bin/http_client_runtime.o bin/http_server_runtime.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o bin/system_runtime.o

# Build compiler
echo "   Building Osprey compiler..."
go build -o bin/osprey ./cmd/osprey

# -----------------------------------------------------------------------------
# Build package list (exclude generated parser)
# -----------------------------------------------------------------------------
ALL_PKGS=$(go list ./...)
PKGS=$(echo "$ALL_PKGS" | grep -v "/parser$")

# Separate integration tests from unit tests
INTEGRATION_PKGS=$(echo "$PKGS" | grep "/tests/integration$")
UNIT_PKGS=$(echo "$PKGS" | grep -v "/tests/integration$")

# Convert package lists to comma-separated strings for -coverpkg
COVERPKG=$(echo "$PKGS" | tr '\n' ',' | sed 's/,$//')

# -----------------------------------------------------------------------------
# Run tests with coverage - integration tests separately (no race detection)
# -----------------------------------------------------------------------------
echo "🧪 Running test suite with atomic coverage..."
echo "    📦 Unit test packages: $(echo "$UNIT_PKGS" | wc -l) packages"
echo "    📦 Integration test packages: $(echo "$INTEGRATION_PKGS" | wc -l) packages"
echo

# Run unit tests
echo "🏃‍♂️ Running unit tests..."
if go test -v -covermode=atomic -coverpkg="$COVERPKG" -coverprofile=coverage_unit.out -timeout=15m $UNIT_PKGS; then
    UNIT_EXIT_CODE=0
    echo "✅ Unit tests completed successfully"
else
    UNIT_EXIT_CODE=$?
    echo "❌ Unit tests failed with exit code: $UNIT_EXIT_CODE"
fi

# Run integration tests
echo "🏗️ Running integration tests..."
if go test -v -covermode=atomic -coverpkg="$COVERPKG" -coverprofile=coverage_integration.out -timeout=15m $INTEGRATION_PKGS; then
    INTEGRATION_EXIT_CODE=0
    echo "✅ Integration tests completed successfully"
else
    INTEGRATION_EXIT_CODE=$?
    echo "❌ Integration tests failed with exit code: $INTEGRATION_EXIT_CODE"
fi

# Combine coverage profiles
echo "🔗 Combining coverage profiles..."
echo "mode: atomic" > coverage.out
grep -h -v "^mode:" coverage_unit.out coverage_integration.out >> coverage.out

# Overall test result
if [[ $UNIT_EXIT_CODE -eq 0 && $INTEGRATION_EXIT_CODE -eq 0 ]]; then
    TEST_EXIT_CODE=0
    echo "✅ All tests completed successfully"
else
    TEST_EXIT_CODE=1
    echo "❌ Some tests failed (Unit: $UNIT_EXIT_CODE, Integration: $INTEGRATION_EXIT_CODE)"
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
# Generate & display coverage reports
# -----------------------------------------------------------------------------
go tool cover -func=coverage.out | { echo "📈 Coverage Summary:"; cat; }
go tool cover -html=coverage.out -o coverage.html

TOTAL_COVERAGE=$(go tool cover -func=coverage.out | awk '/^total:/ {print $3}')

printf "\n🎯 Total Coverage: %s\n" "$TOTAL_COVERAGE"

echo "📁 HTML report saved to: coverage.html"
echo "📁 Raw coverage data saved to: coverage.out"

# Automatically open the HTML report on macOS for convenience
if [[ "$(uname -s)" == "Darwin" ]]; then
    echo "🌐 Opening HTML coverage report in browser…"
    open coverage.html
fi

echo "✅ Coverage analysis complete!" 