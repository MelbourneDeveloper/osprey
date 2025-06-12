---
title: Feature Status
description: Current implementation status of Osprey language features
date: 2025-01-17
tags: ["status", "features", "roadmap"]
author: "Christian Findlay"
---

Current version: **0.1.0-alpha** (not yet released)

## ✅ Complete Features

### Core Language
- **Variables & Constants**: `let` declarations, immutable by default
- **Data Types**: `int`, `string`, `bool`, basic type inference
- **Functions**: Function declarations, expression bodies, named arguments (2+ params)
- **String Interpolation**: `${}` syntax with expressions
- **Pattern Matching**: Basic match expressions, wildcards, type annotation patterns
- **Block Expressions**: Local scoping, multi-statement blocks
- **Arithmetic Operations**: Safe arithmetic returning `Result` types
- **Boolean Operations**: Logical operators and boolean expressions

### Advanced Features
- **Functional Programming**: 
  - Iterator functions (`range`, `forEach`, `map`, `filter`, `fold`)
  - Pipe operator (`|>`)
  - Function composition and chaining
- **Any Type Handling**: Explicit `any` types with pattern matching requirement
- **Result Types**: Error handling without exceptions
- **Type Safety**: No implicit conversions, compile-time type checking

### Concurrency & HTTP
- **Fiber-based Concurrency**: 
  - Lightweight fiber spawning with `spawn`
  - Fiber isolation and communication
  - `await`, `yield` operations
- **HTTP Server**: 
  - Server creation (`httpCreateServer`)
  - Request handling with all methods (GET, POST, PUT, DELETE)
  - Concurrent request processing
- **HTTP Client**:
  - Client creation (`httpCreateClient`)
  - All HTTP methods with custom headers
  - Request/response handling
- **WebSocket Support**:
  - WebSocket server and client connections
  - Real-time bidirectional communication
  - Military-grade security implementation

### Built-in Functions
- **I/O**: `print()`, `input()`, `toString()`
- **Functional Iterators**: Complete pipe operator support
- **Safe Math**: All arithmetic operations return `Result` types

## 🚧 Roadmap Features

### Type System Extensions
- **Record Types with Constraints**: `where` clause validation (partially implemented)
- **Union Types**: Complex algebraic data types with destructuring
- **Generic Types**: Type parameters and polymorphism
- **Module System**: Fiber-isolated modules with proper imports

### Advanced Language Features
- **Extern Declarations**: Full Rust/C interoperability (syntax ready)
- **Advanced Pattern Matching**: Constructor patterns, guards, exhaustiveness checking
- **Select Expressions**: Channel multiplexing for concurrent operations
- **Streaming Responses**: Large HTTP response streaming

### Tooling & Ecosystem
- **Package Manager**: Dependency management system
- **Standard Library**: Comprehensive built-in functions
- **REPL**: Interactive development environment
- **Language Server**: Full IDE support beyond VS Code extension

### Performance & Optimization
- **Compile-time Optimization**: Dead code elimination, inlining
- **Memory Management**: Advanced memory safety features
- **Rust Integration**: Seamless Rust library integration

---

**Note**: Features marked as complete have working examples in the [`examples/tested/`](https://github.com/MelbourneDeveloper/osprey/tree/main/compiler/examples/tested) directory and pass integration tests. 