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
# 1. Dynamically gather all Go packages in the module, excluding the generated
#    parser code (we do not want to track coverage for generated files).
# 2. Execute the full test-suite with race detection and atomic coverage.
# 3. Produce both textual and HTML coverage summaries.
# =============================================================================

# Emoji-rich status messaging keeps things fun but concise.
echo "🧪 Running comprehensive code-coverage analysis…"

# -----------------------------------------------------------------------------
# Clean up any previous artifacts
# -----------------------------------------------------------------------------
rm -f coverage.out coverage.html

# -----------------------------------------------------------------------------
# Build package list (exclude generated parser)
# -----------------------------------------------------------------------------
ALL_PKGS=$(go list ./...)
PKGS=$(echo "$ALL_PKGS" | grep -v "/parser$")

# Convert package list to comma-separated string for -coverpkg
COVERPKG=$(echo "$PKGS" | tr '\n' ',' | sed 's/,$//')

# -----------------------------------------------------------------------------
# Run tests with coverage across all selected packages
# -----------------------------------------------------------------------------
echo "📊 Running tests with coverage…"
go test -v -race -covermode=atomic -coverpkg="$COVERPKG" -coverprofile=coverage.out $PKGS

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