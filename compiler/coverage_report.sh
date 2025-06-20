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