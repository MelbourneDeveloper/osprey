#!/bin/bash
set -e

echo "Running tests with coverage..."

# Clean and rebuild
make clean && make build

# Get packages for testing and coverage
TEST_PKGS=$(go list ./... | grep -v "/parser$")
# Only include source packages in coverpkg, not test packages
COVERPKG=$(go list ./... | grep -v "/parser$" | grep -v "/tests/" | tr '\n' ',' | sed 's/,$//')

# Run tests - fail fast (suppress verbose output for cleaner CI)
# First run with normal output to capture the result
if ! go test -covermode=atomic -coverpkg="$COVERPKG" -coverprofile=coverage.out $TEST_PKGS 2>&1 | tee test_output.tmp; then
    echo ""
    echo "❌ TESTS FAILED!"
    echo ""
    echo "=== FAILURE DETAILS ==="
    # Show the specific test failures with context
    grep -A 5 -B 2 "--- FAIL:" test_output.tmp || true
    # Also show any runtime errors/panics
    grep -E "(panic:|runtime error:|signal:|core dumped)" test_output.tmp || true
    # Show any assertion failures
    grep -E "(Expected:|Got:|Error:|Failed to|Output mismatch)" test_output.tmp | head -20 || true
    echo "==================="
    rm -f test_output.tmp
    exit 1
fi
rm -f test_output.tmp

echo "✅ All tests passed"

# Filter out AST interface methods from coverage report
grep -v "isStatement\|isExpression" coverage.out > coverage_filtered.out || cp coverage.out coverage_filtered.out

# Show coverage summary only
echo ""
echo "📊 Coverage Summary:"
go tool cover -func=coverage_filtered.out | grep -E "(total:|\.go:)"
TOTAL=$(go tool cover -func=coverage_filtered.out | awk '/^total:/ {print $3}')
echo ""
echo "🎯 Total Coverage: $TOTAL"

# Generate HTML
go tool cover -html=coverage_filtered.out -o coverage.html
echo "HTML report: coverage.html"