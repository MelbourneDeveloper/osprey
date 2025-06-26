---
layout: page
title: "Introduction"
description: "Osprey Language Specification: Introduction"
date: 2025-06-26
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0001-introduction/"
---

1. [Introduction](0001-Introduction.md)
   - [Completeness](#11-completeness)
   - [Principles](#12-principles)

## 1. Introduction

Osprey is a modern functional programming oriented language designed for elegance, safety, and performance.. It emphasizes:

- **Named arguments** for multi-parameter functions to improve readability
- **Strong type inference** to reduce boilerplate while maintaining safety
- **String interpolation** for convenient text formatting
- **Pattern matching** for elegant conditional logic
- **Immutable-by-default** variables with explicit mutability
- **Fast HTTP servers and clients** with built-in streaming support
- **WebSocket support** for real-time two-way communication

## 1.1 Completeness

🚧 **IMPLEMENTATION STATUS**: The Osprey language and compiler are not complete. The documentation here does not represent the language at this present time. This specification represents the design aims for the language, a description of the syntax and a description of where the roadmap is taking the language. Developers should pay attention to the spec first are foremost as the single source of truth in regards to the syntax.

## 1.2 Principles

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