#!/bin/bash

# Osprey C Runtime Test Runner
# This script runs our Unity tests and outputs results in a format VS Code can parse

set -e

echo "🧪 Running Osprey C Runtime Tests..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to run a test and capture results
run_test() {
    local test_name="$1"
    local test_executable="$2"
    
    echo ""
    echo -e "${YELLOW}=== Running $test_name ===${NC}"
    
    if [ ! -f "$test_executable" ]; then
        echo -e "${RED}ERROR: $test_executable not found${NC}"
        return 1
    fi
    
    if ./"$test_executable"; then
        echo -e "${GREEN}✅ $test_name PASSED${NC}"
        return 0
    else
        echo -e "${RED}❌ $test_name FAILED${NC}"
        return 1
    fi
}

# Build tests first
echo "📦 Building Unity tests..."
echo "Building system runtime tests..."
clang -o test_system_runtime_unity test_system_runtime_unity.c unity.c system_runtime.c -pthread -std=c11
echo "Building fiber runtime tests..."
clang -o test_fiber_runtime_unity test_fiber_runtime_unity.c unity.c fiber_runtime.c system_runtime.c -pthread -std=c11

# Track test results
total_tests=0
failed_tests=0

# Run system runtime tests
total_tests=$((total_tests + 1))
if ! run_test "System Runtime Unity Tests" "test_system_runtime_unity"; then
    failed_tests=$((failed_tests + 1))
fi

# Run fiber runtime tests  
total_tests=$((total_tests + 1))
if ! run_test "Fiber Runtime Unity Tests" "test_fiber_runtime_unity"; then
    failed_tests=$((failed_tests + 1))
fi

# Cleanup
echo "🧹 Cleaning up test executables..."
rm -f test_system_runtime_unity test_fiber_runtime_unity

# Summary
echo ""
echo "========================================="
if [ $failed_tests -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! ($total_tests/$total_tests)${NC}"
    exit 0
else
    echo -e "${RED}💥 $failed_tests/$total_tests TESTS FAILED!${NC}"
    exit 1
fi 