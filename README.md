<p align="center">
  <img src="website/src/assets/images/logo.png" alt="Osprey logo" width="160" />
</p>

<h1 align="center">Osprey Programming Language</h1>

<p align="center">
  A modern functional programming language designed for elegance, safety, and
  performance.<br/>Written in Rust, outputs to LLVM.
</p>

⭐ **[Star us on GitHub](https://github.com/MelbourneDeveloper/osprey)** to support the project and allow us to submit to Homebrew! ⭐

## Installation

```bash
# macOS / Linux (Homebrew)
brew install nimblesite/tap/osprey

# Windows (Scoop)
scoop bucket add nimblesite https://github.com/Nimblesite/scoop-bucket
scoop install osprey
```

Osprey shells out to LLVM (`llc`) and a C compiler at compile time; the package
managers pull those in as dependencies (`llvm` for brew; `llvm` + `gcc` for scoop).

The [VS Code extension](https://marketplace.visualstudio.com/items?itemName=nimblesite.osprey)
(`nimblesite.osprey`) bundles a version-matched compiler and a Rust language
server (`osprey lsp`, built on [lspkit](https://github.com/Nimblesite/lspkit)) for
live diagnostics, hover, go-to-definition, and completion. The same server is
editor-agnostic — Neovim and Zed are on the roadmap. See
[Language Server & Editors](docs/specs/0020-LanguageServerAndEditors.md).

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

- `crates/` - Main Osprey compiler (Rust workspace: osprey-ast, osprey-syntax, osprey-types, osprey-codegen, osprey-runtime-sys, osprey-cli)
- `tree-sitter-osprey/` - Tree-sitter grammar (parser)
- `compiler/` - Pure-C runtime sources (`runtime/`) + example programs (`examples/`)
- `vscode-extension/` - VSCode language support
- `website/` - Documentation site
- `webcompiler/` - Browser-based compiler
- `homebrew-package/` - Homebrew tap
- `.devcontainer` - Configuration for the dev container

## Documentation

- [Language specification](docs/specs/)
- [API reference](website/src/docs/)
- [Contributing guide](CONTRIBUTING.md)
- [Release process](docs/RELEASING.md) — tag `v*` to release; CI runs only on PRs to `main`.

## Development

Built on proven tech: Rust for the compiler, tree-sitter for parsing, and LLVM for code generation.

**AI-Assisted Development**: Claude Sonnet 4 with Cursor makes implementing language features accessible. Check out [CONTRIBUTING.md](CONTRIBUTING.md) for the workflow.

**Use VS Code Dev Containers** - strongly recommended. Open in VS Code and hit "Reopen in Container".

```bash
make build         # C runtime archives + cargo build --release + extension
make test          # Run all tests + coverage thresholds
make lint          # cargo clippy + extension lint
make ci            # lint + test + build (full CI simulation)
make install       # Install compiler + runtime archives locally
```

The compiler binary lands at `target/release/osprey`.

## Status

🚧 **Alpha**: Core language features implemented. Algebraic effects system working with compile-time safety, but are missing some features. HTTP and advanced features in development.

See [docs/specs/](docs/specs/) for implementation status.

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
