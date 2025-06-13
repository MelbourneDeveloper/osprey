# Osprey Programming Language

A modern functional programming oriented language designed for elegance, safety, and performance. It is written in Go and outputs to LLVM.

## Installation

```bash
# Add the tap
brew tap melbournedeveloper/osprey

# Install Osprey
brew install osprey
```

## Local Development Quick Start

```bash
git clone https://github.com/christianfindlay/vexels.git
cd vexels/compiler
make build
./bin/osprey examples/simple.osp
```

## Language Features

- **Functional-first**: Immutable data, pattern matching, pipe operators
- **Type-safe**: Algebraic data types with variant types
- **HTTP-native**: Built-in server/client with streaming support
- **Fiber concurrency**: Lightweight isolated execution contexts
- **Zero-cost abstractions**: Compiles to efficient LLVM IR

## Syntax Examples

```osprey
// Variables and functions
let x = 42
fn add(x, y) = x + y

// Types and pattern matching
type Result = Ok { value: Int } | Error { message: String }
let status = match result {
    Ok -> "success"
    Error -> "failed"
}

// HTTP server
httpServer(8080) |> onRequest((req) => 
    response(200, "Hello World")
)
```

## Project Structure

- [`compiler/`](compiler/) - Main Osprey compiler (Go + ANTLR)
- [`vscode-extension/`](vscode-extension/) - VSCode language support
- [`website/`](website/) - Documentation site
- [`webcompiler/`](webcompiler/) - Browser-based compiler

## Documentation

- [Language specification](compiler/spec.md)
- [API reference](website/src/docs/)
- [Contributing guide](CONTRIBUTING.md)

## Development

```bash
cd compiler
make install-deps  # Install Go dependencies
make build         # Build compiler
make test          # Run tests
make regenerate-parser  # Regenerate from grammar
```

## Status

🚧 **Alpha**: Core language features implemented. HTTP and fiber systems in development.

See [compiler/spec.md](compiler/spec.md) for implementation status and roadmap.

## License

MIT License - see [LICENSE](LICENSE) 