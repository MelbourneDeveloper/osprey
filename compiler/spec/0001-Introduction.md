# Introduction

- [Completeness](#completeness)
- [Core Principles](#core-principles)

Osprey is a modern functional programming language designed for elegance, safety, and performance. It emphasizes:

- **Named arguments** for multi-parameter functions to improve readability
- **Strong type inference** to reduce boilerplate while maintaining safety
- **String interpolation** for convenient text formatting
- **Pattern matching** for elegant conditional logic
- **Immutable-by-default** variables with explicit mutability
- **Fast HTTP servers and clients** with built-in streaming support
- **WebSocket support** for real-time two-way communication

## Completeness

**Note**: The Osprey language and compiler are under active development. This specification represents design goals and planned features. The spec is the authoritative source for syntax and behavior.

## Core Principles

- Elegance (simplicity, ergonomics, efficiency), safety (fewer footguns, security at every level), performance (uses the most efficient approach and allows the use of Rust interop for extreme performance)
- No more than 1 way to do anything
- ML style syntax by default
- Make illegal states unrepresentable. There are no exceptions or panics. Anything than can result in an error state returns a result object
- Referential transparency
- Simplicity
- Interopability with Rust for high performance workloads
- Interopability with Haskell (future) for fundamental correctness
- Static/strong typing. Nothing should be "any" unless EXPLICITLY declared as any
- Minimal ceremony. No main function necessary for example.
- **Fast HTTP performance** as a core design principle
- **Streaming by default** for large responses to prevent memory issues