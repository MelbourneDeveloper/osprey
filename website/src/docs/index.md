---
layout: page
title: Osprey documentation
description: Start with algebraic effects and reusable handlers, then build native, web and mobile applications with Osprey.
---

Osprey uses typed algebraic effects to separate application requests from their implementations. Start with a handler, then see the same approach in the banking and mobile applications.

## Start here

- **[Algebraic effects](/docs/effects/)** — declare an operation, call a reusable handler, replace it in a test, and choose an arithmetic policy.
- **[One model, two backends](/docs/one-model-two-backends/)** — write GPT-2 once against a static `Tensor` effect, then run it over host lists or on the GPU surface by choosing a handler.
- **[Install Osprey](/docs/installation/)** — set up the compiler and LLVM/clang toolchain on macOS, Linux, or Windows.
- **[Build My First App](/docs/my-first-app/)** — create a native CLI that models JSON as a recursive algebraic data type, reads a JSON file, and safely writes the next state.
- **[Try the Playground](/playground/)** — run small programs in the browser before setting up a local toolchain.

## Build applications

- **[Build an iOS and Android app](/docs/mobile-apps/)** — handle application requests as native commands for SQLite, GitHub and reactive controls.
- **[Build a web app](/docs/web-apps/)** — follow Talon Bank from storage and audit handlers to shared web, Android and iOS clients.
- **[Explore the WebAssembly studio](/wasm/)** — compare the Default and ML flavors in a browser-hosted example.
- **[Browse working examples](https://github.com/Nimblesite/osprey/tree/main/examples)** — study native, HTTP, WebSocket, terminal, graphics, and WebAssembly programs.
- **[Document your modules](/docs/documentation/)** — generate HTML API references, add Markdown guides and custom CSS, and test documentation examples.

## Language

- **[Read the language specification](/spec/)** — the intended language contract; check feature status for current implementation limits.
- **[Use dense computation buffers](/spec/0034-gpucomputation/)** — typed buffers and pure kernels, run on Metal or OpenCL when `OSPREY_GPU` selects a device and on the CPU otherwise.
- **[Check feature status](/status/)** — see what is stable, partial, experimental, or planned.
- **[Browse keywords](/docs/keywords/)** — declarations, bindings, matching, imports, and literals.

A flavor is a surface syntax: its concrete syntax tree lowers into one canonical AST, and everything after that is shared. Default files use `.osp`, braces, `fn`, and parenthesized calls. ML files use `.ospml`, indentation, currying by default, and whitespace application. Python files use `.py` with the `# osprey: flavor=python` marker, typed Python syntax, and decorators for effects.

## Reference

- **[Functions](/docs/functions/)** — file I/O, strings, collections, processes, networking, JSON document queries, and concurrency.
- **[Types](/docs/types/)** — built-in value and runtime types.
- **[Operators](/docs/operators/)** — arithmetic, comparison, and pipeline operators.
- **[Keywords](/docs/keywords/)** — language syntax by keyword.
