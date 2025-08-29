# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Rules

- **DO NOT COMMIT/PUSH CODE** unless explicitly requested
- **NEVER DUPLICATE CODE** - Edit in place, **SEARCH** for code **BEFORE** creating new functions/constants
- **NO PLACEHOLDERS** - Fix existing placeholders or fail with error
- **NEVER IGNORE TESTS** - Don't reduce assertions to make tests pass, fail loudly
- **KEEP ALL FILES UNDER 500 LOC** - Break large files into focused modules  
- **FOLLOW STATIC ANALYSIS** - Pay attention to linters and fix issues
- **MOVE FILES, DON'T COPY** - Use CLI commands to move files
- **NO SWEARING IN CODE** - Keep code professional
- **USE CONSTANTS** - Name values meaningfully instead of using literals

## Essential Commands

**Build & Development:**
- `make build` - Build the osprey compiler (includes linting and C runtime builds)
- `make clean` - Clean all build artifacts
- `make install` - Install compiler globally to /usr/local/bin

**Testing:**
- `make test` - Run all tests with runtime verification
- `make test-runtime` - Run C runtime tests to verify segfault fixes
- `make c-test` - Run C runtime tests with strict error checking
- `./coverage_report.sh` - Generate comprehensive coverage report with HTML output

**Linting:**
- `make lint` - Run Go linter (golangci-lint)
- `make lint-fix` - Run linter with auto-fix
- `make c-lint` - Run C linter with strict warnings-as-errors

**Development:**
- `make run FILE=<path>` - Run compiler on specific file
- `make regenerate-parser` - Regenerate ANTLR parser from grammar

## Architecture Overview

**Core Components:**
- `cmd/osprey/` - Main CLI entry point
- `internal/ast/` - Abstract Syntax Tree definitions and builders
- `internal/codegen/` - LLVM code generation with Hindley-Milner type inference
- `internal/cli/` - Command line interface and security controls
- `parser/` - Generated ANTLR parser (do not edit manually)
- `runtime/` - C runtime libraries for fiber, HTTP, WebSocket, and system operations

**Key Files:**
- `osprey.g4` - ANTLR grammar definition
- `Makefile` - Build system with strict C compiler flags
- `go.mod` - Go module using ANTLR v4, LLVM IR, and testing frameworks

**Type System:**
- Implements Hindley-Milner type inference (`internal/codegen/type_inference.go`)
- Supports algebraic effects, fibers, and pattern matching
- Strong type safety with compile-time verification

**Runtime Architecture:**
- Multi-language: Go compiler generates LLVM IR, C runtime for performance-critical operations
- Fiber-based concurrency with isolation
- HTTP/WebSocket support with OpenSSL integration
- All C code compiled with warnings-as-errors for security

**Testing Structure:**
- `tests/integration/` - Full compiler integration tests
- `tests/unit/` - Unit tests for individual components
- `examples/tested/` - Working code examples with expected outputs
- `examples/failscompilation/` - Error case examples with expected error messages

**Security:**
- Configurable security policies in codegen
- Sandboxing capabilities for HTTP, file access, and process execution
- C runtime built with security hardening flags

## Development Guidelines

**Code Generation:**
- All LLVM IR generation happens in `internal/codegen/`
- Type inference must complete before code generation
- Effect handlers are generated for algebraic effects system

**Testing:**
- Use `make test` for comprehensive testing including runtime verification
- C runtime tests focus on preventing segfaults and memory safety
- Integration tests verify end-to-end compilation and execution

**C Runtime:**
- Located in `runtime/` directory
- Compiled with maximum warning levels and security hardening
- Test files use standard C assertions
- HTTP runtime requires OpenSSL 3.5.0+

**Parser Changes:**
- Modify `osprey.g4` grammar file
- Run `make regenerate-parser` to update generated code
- Requires ANTLR installation

This compiler implements a functional programming language with algebraic effects, fiber-based concurrency, and strong type safety through Hindley-Milner inference.