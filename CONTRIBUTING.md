# Contributing to Osprey

Want to help build a programming language? You're in the right place. Building a language is a several year long process. Right now, the aims are building community, getting the first example apps deployed and shaping the ergonomics of the language.

Even discussions are great at this point.

## The Tech Stack

Osprey is built on battle-tested tools:
- **[Go](https://golang.org/)** - Compiler implementation 
- **[ANTLR](https://www.antlr.org/)** - Grammar parsing
- **[LLVM](https://llvm.org/)** - Code generation

## AI Assisted Development

You don't need a CS degree to implement language features. I've been using Claude Sonnet 4 with Cursor, and it's the first combo that actually guided me through the process of building a compiler. Other AI agents will work too, but this setup lets anyone contribute to language design.

The AI can help you:
- Parse/create ANTLR grammars and understand AST patterns
- Implement new operators and language constructs
- Debug LLVM IR generation
- Write comprehensive tests

## Getting Started

**Use VS Code Dev Containers** - strongly recommended. Open in VS Code and hit "Reopen in Container". Everything's already configured.

```bash
# Fork and clone
git clone https://github.com/MelbourneDeveloper/osprey.git
cd osprey/compiler

# Build and test
make install-deps
make build
make test
```

## What to Work On

1. **Language features** - Check [spec.md](compiler/spec.md) for "NOT IMPLEMENTED" 
2. **New operators** - Add arithmetic, comparison, or logical operators
3. **Pattern matching** - Extend match expressions
4. **Standard library** - Add built-in functions
5. **Make examples** - The HTTP server works. Try building an API
6. **Test compilation errors** - Make sure the compiler is forcing you to do things the right way

## The AI Workflow

1. **Understand the pattern** - Ask your AI: "How does pattern matching work in this codebase?"
2. **Implement incrementally** - Start with parsing, then AST, then codegen
3. **Test everything** - Add examples to `compiler/examples/tested/`
4. **Fix edge cases** - Let the AI help debug LLVM IR issues

Example prompts:
- "Add a new arithmetic operator to the ANTLR grammar"
- "Implement string interpolation in the AST builder"
- "Generate LLVM IR for this pattern match"

## Code Guidelines

- Follow existing patterns - Go lints enforce a lot
- Test new features thoroughly
- Keep changes focused
- Fix linter errors before submitting

## Getting Help

- Open an issue for discussion
- Check existing issues and documentation
- Review similar implementations in the codebase
- The spec is the source of truth

## License

By contributing, you agree that your contributions will be licensed under the MIT License. 