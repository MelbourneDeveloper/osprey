# Introduction

Osprey is a functional programming language designed for safety, performance, and expressiveness.

## Core Features

- Named arguments for multi-parameter functions
- Hindley-Milner type inference with strong static typing
- Pattern matching for all conditional logic
- Immutable-by-default with explicit mutability
- Algebraic effects with compile-time safety
- Result types for all error cases (no exceptions or panics)
- Built-in HTTP/WebSocket support with streaming
- Lightweight fiber-based concurrency

## Design Principles

- **Safety**: Make illegal states unrepresentable through static verification
- **Simplicity**: One idiomatic way to accomplish each task
- **Performance**: LLVM compilation with Rust interop for performance-critical code
- **Functional**: Referential transparency, immutable data structures, pure functions
- **Type Safety**: Strong static typing with Hindley-Milner inference; `any` type requires explicit declaration
- **No Exceptions**: All error cases return Result types, enforced at compile time
- **ML Heritage**: Syntax and semantics inspired by ML family languages

## Development Status

This specification is the authoritative source for Osprey syntax and behavior. The language and compiler are under active development; implementation status is noted where relevant.