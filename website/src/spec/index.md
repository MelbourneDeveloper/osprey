---
layout: page
title: "Osprey Language Specification"
description: "Complete language specification and syntax reference for the Osprey programming language"
date: 2025-07-11
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/"
---

# Osprey Language Specification

**Version:** 0.2.0-alpha  
**Date:** 2025-07-11  
**Author:** Christian Findlay

## Table of Contents

1. [Introduction](/spec/0001-introduction/)
2. [Lexical Structure](/spec/0002-lexicalstructure/)
3. [Syntax](/spec/0003-syntax/)
4. [Semantics](/spec/0004-semantics/)
5. [Type System](/spec/0005-typesystem/)
6. [Function Calls](/spec/0006-functioncalls/)
7. [String Interpolation](/spec/0007-stringinterpolation/)
8. [Pattern Matching](/spec/0008-patternmatching/)
9. [Block Expressions](/spec/0009-blockexpressions/)
10. [Boolean Operations](/spec/0010-booleanoperations/)
11. [Loop Constructs and Functional Iterators](/spec/0011-loopconstructsandfunctionaliterators/)
12. [Lightweight Fibers and Concurrency](/spec/0012-lightweightfibersandconcurrency/)
13. [Built-in Functions](/spec/0013-built-infunctions/)
14. [Error Handling](/spec/0014-errorhandling/)
15. [HTTP](/spec/0015-http/)
16. [WebSocket Functions](/spec/0016-websockets/)
17. [Security and Sandboxing](/spec/0017-securityandsandboxing/)
18. [Algebraic Effects](/spec/0018-algebraiceffects/)

## About This Specification

This specification defines the complete syntax and semantics of the Osprey programming language. Each section is available as a separate page for easy navigation and reference.

The Osprey language is designed for elegance, safety, and performance, emphasizing:

- **Named arguments** for multi-parameter functions to improve readability
- **Strong type inference** to reduce boilerplate while maintaining safety
- **String interpolation** for convenient text formatting
- **Pattern matching** for elegant conditional logic
- **Immutable-by-default** variables with explicit mutability
- **Fast HTTP servers and clients** with built-in streaming support
- **WebSocket support** for real-time two-way communication

## Implementation Status

🚧 **NOTE**: The Osprey language and compiler are actively under development. This specification represents the design goals and planned features. Please refer to individual sections for current implementation status.
