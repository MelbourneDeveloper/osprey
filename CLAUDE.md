# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Core Development Principles

- **NEVER DUPLICATE CODE** - Edit in place, never create new versions
- **NO PLACEHOLDERS** - Fix existing placeholders or fail with error
- **SEARCH BEFORE ADDING** - Check for existing code before creating new functions/constants
- **NEVER IGNORE FAILING TESTS** - Don't reduce assertions to make tests pass, fail loudly
- **KEEP ALL FILES UNDER 500 LOC** - Break large files into focused modules 
- **RUN LINTER REGULARLY** - lints are strict. Obey them!!
- **FP STYLE CODE** - pure functions over OOP style
- **NEVER COMMIT/PUSH** unless explicitly requested
- **FOLLOW STATIC ANALYSIS** - Pay attention to linters and fix issues
- **USE CONSTANTS** - Name values meaningfully instead of using literals

## Commands

**Primary Development Commands:**
```bash
cd compiler
make build         # Build compiler with linting and C runtime
make test          # Run all tests including C runtime verification
make lint          # Run Go linter (golangci-lint)  
make install       # Install compiler globally to /usr/local/bin
make clean         # Clean all build artifacts
```

**Testing Commands:**
```bash
make test-runtime  # Run C runtime tests for memory safety
make c-test        # Run C runtime tests with strict error checking
make c-lint        # Run C linter with warnings-as-errors
./coverage_report.sh # Generate HTML coverage report
```

**Development Commands:**
```bash
make run FILE=<path>              # Run compiler on specific file
make regenerate-parser            # Regenerate ANTLR parser from grammar
make lint-fix                     # Run linter with auto-fix
```

**VSCode Extension:**
```bash
cd vscode-extension
npm install && npm run compile    # Build VSCode extension
npm test                         # Run extension tests
```

**Website Development:**
```bash
cd website
npm install && npm run dev       # Start local development server
npm run build                    # Build static site
```

**WebCompiler (Browser-based):**
```bash
cd webcompiler  
npm install && npm start         # Start web-based compiler service
```

## High-Level Architecture

**Repository Structure:**
- `compiler/` - Core Osprey compiler (Go + ANTLR → LLVM)
- `vscode-extension/` - VSCode language support with TypeScript
- `website/` - Documentation site using 11ty static site generator
- `webcompiler/` - Node.js web service for browser compilation
- `homebrew-package/` - Homebrew tap for macOS installation

**Compiler Architecture (Go-based):**
- **Parser**: ANTLR4 grammar (`osprey.g4`) generates Go parser
- **AST**: Abstract Syntax Tree builders in `internal/ast/`
- **Type System**: Hindley-Milner type inference in `internal/codegen/type_inference.go`
- **Code Generation**: LLVM IR generation in `internal/codegen/`
- **Runtime**: C libraries for fiber concurrency, HTTP/WebSocket, system operations

**Language Features:**
- **Algebraic Effects**: First-class effects system with compile-time safety
- **Fiber Concurrency**: Lightweight isolated execution contexts
- **Pattern Matching**: Union types with exhaustiveness checking
- **Functional Programming**: Immutable data, pipe operators, iterators
- **HTTP/WebSocket**: Built-in networking with streaming support
- **Type Safety**: Strong static typing with inference

**Multi-Language Runtime:**
- **Go**: Compiler frontend, parsing, type checking
- **C**: Performance-critical runtime (fibers, HTTP, WebSockets)
- **LLVM IR**: Compilation target for optimized execution
- **Rust Integration**: Optional interop library for math utilities

**Key Technical Patterns:**
- Effects are declared with `effect` keyword and handled with `handle...in` expressions
- Unhandled effects cause compilation errors (world-first compile-time effect safety)
- Pattern matching is mandatory for `any` types and union types
- All HTTP/WebSocket operations return `Result<T, String>` for error handling
- Fiber isolation prevents shared memory bugs through message passing

**Testing Strategy:**
- `tests/integration/` - Full compiler integration tests
- `tests/unit/` - Component-specific unit tests  
- `examples/tested/` - Working code examples with expected outputs
- `examples/failscompilation/` - Error cases with expected error messages
- `runtime/` - C tests for memory safety and segfault prevention

**Security Architecture:**
- Configurable sandboxing for file access, HTTP, and process execution
- C runtime compiled with security hardening flags (`-D_FORTIFY_SOURCE=2`, `-fstack-protector-strong`)
- All warnings treated as errors in C compilation
- Effect system provides capability-based security

**Development Workflow:**
1. **Grammar Changes**: Edit `osprey.g4`, run `make regenerate-parser`
2. **Language Features**: Implement in AST builders, then codegen
3. **Testing**: Add examples to `examples/tested/` and error cases to `examples/failscompilation/`
4. **Type System**: Extend `type_inference.go` for new type rules
5. **Runtime**: Add C functions in `runtime/` for system operations

**AI-Assisted Development Notes:**
- This compiler is built using various AI agents and models
- AI can help with ANTLR grammars, LLVM IR generation, and type inference
- The codebase follows clear patterns that AI can recognize and extend
- Use VS Code Dev Container for consistent development environment

This is a functional programming language compiler with algebraic effects, fiber-based concurrency, and strong compile-time safety guarantees.