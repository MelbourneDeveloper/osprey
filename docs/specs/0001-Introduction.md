# Introduction

Osprey is a statically-typed functional language in the ML family. It compiles to native code via LLVM.

## Core Features

- Hindley-Milner type inference; explicit annotations are optional.
- Pattern matching as the only conditional construct (no `if`/`else`).
- Immutable bindings by default; `mut` opts in to mutability.
- Algebraic effects checked at compile time.
- `Result<T, E>` for all fallible operations; no exceptions, panics, or null.
- Named arguments required for functions of two or more parameters.
- Lightweight fibers and channel-based concurrency.
- Automatic memory management with no observable collector — ARC by default, tracing GC selectable, and a `--static-memory` mode with zero runtime memory operations.
- Built-in HTTP and WebSocket support.

## Status

This specification is the authoritative source for Osprey syntax and behaviour. The language and compiler are under active development; implementation status is called out per chapter where it diverges from the specification.