#!/bin/bash
set -e

echo "Running tests with coverage..."

# Clean and rebuild
make clean && make build

# Get packages for testing and coverage
TEST_PKGS=$(go list ./... | grep -v "/parser$")
# Only include source packages in coverpkg, not test packages
COVERPKG=$(go list ./... | grep -v "/parser$" | grep -v "/tests/" | tr '\n' ',' | sed 's/,$//')

# Run tests - fail fast
if ! go test -v -covermode=atomic -coverpkg="$COVERPKG" -coverprofile=coverage.out $TEST_PKGS; then
    echo "❌ TESTS FAILED!"
    echo "Last failing test output:"
    go test -v $TEST_PKGS 2>&1 | grep -E "(FAIL|--- FAIL:)" | tail -5
    exit 1
fi

echo "✅ All tests passed"

# Show coverage
go tool cover -func=coverage.out
TOTAL=$(go tool cover -func=coverage.out | awk '/^total:/ {print $3}')
echo "Total Coverage: $TOTAL"

# Generate HTML
go tool cover -html=coverage.out -o coverage.html
echo "HTML report: coverage.html"