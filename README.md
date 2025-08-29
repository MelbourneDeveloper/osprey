# Osprey Programming Language

A modern functional programming language designed for elegance, safety, and performance. Written in Go, outputs to LLVM.

⭐ **[Star us on GitHub](https://github.com/MelbourneDeveloper/osprey)** to support the project and allow us to submit to Homebrew! ⭐

## Installation

```bash
# Add the tap
brew tap melbournedeveloper/osprey

# Install Osprey
brew install osprey
```

## Language Features

- **Functional-first**: Immutable data, pattern matching, pipe operators
- **Algebraic Effects**: First-class effects system with compile-time safety
- **Type-safe**: Algebraic data types with variant types
- **HTTP-native**: Built-in server/client with streaming support
- **Fiber concurrency**: Lightweight isolated execution contexts
- **Zero-cost abstractions**: Compiles to efficient LLVM IR

## Revolutionary Safety

🚀 **World's first language with 100% compile-time effect safety** - unhandled effects cause compilation errors, not runtime crashes!

## Syntax Example

```osprey
// 🔒 HANDLER ISOLATION SIMPLE TEST 🔒

effect Logger {
    log: fn(string) -> Unit
}

// Main function with different handlers
fn main() -> Unit = {
    print("🔒 Testing Handler Isolation")
    
    // Production handler
    let result1 = handle Logger
        log msg => print("[PROD] " + msg)
    in {
        perform Logger.log("Processing task: 5")
        10
    }
    
    // Debug handler
    let result2 = handle Logger
        log msg => print("[TEST] " + msg)
    in {
        perform Logger.log("Processing task: 12")
        24
    }
    
    // Silent handler
    let result3 = handle Logger
        log msg => 0
    in {
        perform Logger.log("Processing task: 0")
        0
    }
    
    print("📊 Results: Prod=" + toString(result1) + ", Test=" + toString(result2) + ", Silent=" + toString(result3))
} 
```

## Project Structure

- `compiler/` - Main Osprey compiler (Go + ANTLR)
- `vscode-extension/` - VSCode language support
- `website/` - Documentation site
- `webcompiler/` - Browser-based compiler
- `homebrew-package/` - Homebrew tap
- `.devcontainer` - Configuration for the dev container

## Documentation

- [Language specification](compiler/spec/)
- [API reference](website/src/docs/)
- [Contributing guide](CONTRIBUTING.md)

## Development

Built on proven tech: Go for the compiler, ANTLR for parsing, and LLVM for code generation.

**AI-Assisted Development**: Claude Sonnet 4 with Cursor makes implementing language features accessible. Check out [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow.

**Use VS Code Dev Containers** - strongly recommended. Open in VS Code and hit "Reopen in Container".

```bash
cd compiler
make build         # Build compiler
make test          # Run tests
make install       # Install locally
```

## Status

🚧 **Alpha**: Core language features implemented. Algebraic effects system working with compile-time safety, but are missing some features. HTTP and advanced features in development.

See [compiler/spec/](compiler/spec/) for implementation status.

## Recent Major Updates

- **Algebraic Effects System**: Complete implementation with compile-time safety guarantees
- **Effect Declarations**: `effect` keyword for defining effect operations
- **Perform Expressions**: `perform` keyword for effect operations
- **Handler Expressions**: `handle...in` syntax for effect handling
- **Compile-Time Verification**: Unhandled effects cause compilation errors (world-first!)

## License

MIT License - see [LICENSE](LICENSE)

---

⭐ **[Give us a star on GitHub](https://github.com/MelbourneDeveloper/osprey)** if you like what we're building! ⭐ 
