# Osprey Development Container

This directory contains configuration for a development container that provides all necessary dependencies for the Osprey compiler and VS Code extension development.

## Features

- **Rust (stable)** for compiler development (`crates/` workspace)
- **LLVM 14 / clang** for compiling generated IR and the C runtime
- **Node.js 20** for the VS Code extension, webcompiler, and website
- **All necessary VS Code extensions** pre-installed (rust-analyzer, C/C++, Makefile Tools)

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers extension for VS Code](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Opening the Project in the Dev Container

1. Open the `osprey.code-workspace` file in VS Code
2. When prompted to "Reopen in Container", click "Reopen in Container"
   - Alternatively, press F1, type "Dev Containers: Reopen in Container" and press Enter

The container will automatically run the post-create script (`make setup` + an
initial `make build`).

## Available Scripts

- `build-compiler.sh` - Builds the Osprey compiler (`make build` at the repo root)

```bash
# Build the compiler
.devcontainer/build-compiler.sh

# Build VS Code extension (use the proper build script)
cd vscode-extension && ./build.sh
```

## Development Tasks

### Compiler Development

Use the root Makefile:
```bash
make build          # C runtime archives + cargo build --release + extension
make test           # All tests + coverage thresholds + differential harness
make lint           # cargo clippy + extension lint
make fmt            # cargo fmt + prettier
make ci             # lint + test + build (full CI simulation)
make clean          # Clean build artifacts
```

The compiler binary lands at `target/release/osprey`.

### VS Code Extension Development

Navigate to the extension directory:
```bash
cd /vscode-extension
npm install         # Install dependencies
npm run compile     # Compile the extension
npm run watch       # Watch for changes
npm run package     # Package the extension
```

To debug the extension:
1. Open the Run and Debug view in VS Code (Ctrl+Shift+D)
2. Select "Run Extension"
3. Press F5

## Project Structure

```
/osprey/
├── crates/             # Osprey compiler (Rust workspace)
├── tree-sitter-osprey/ # Tree-sitter grammar (parser)
├── compiler/           # C runtime sources (runtime/) + examples
├── vscode-extension/   # VS Code extension (TypeScript)
├── webcompiler/        # Web compiler (ignored in dev container)
└── .devcontainer/      # Dev container configuration
```

## Troubleshooting

### Container Build Issues
If you encounter issues building the container:

1. Try rebuilding: F1 → "Dev Containers: Rebuild Container"
2. Clear Docker cache: `docker system prune -a`
3. Check Docker logs for specific errors

### WSL Issues (Windows)
If you see WSL-related errors:

1. Restart Docker Desktop
2. Update WSL2: `wsl --update`
3. Restart VS Code

### Node.js/npm Issues
The container uses Node.js 20 which is compatible with the latest npm. If you encounter version conflicts, the setup scripts handle the compatibility automatically.

## Manual Testing

After the container starts, you can verify everything works:

```bash
# Test individual components
rustc --version        # Rust
cargo --version        # Cargo
llc --version          # LLVM
node --version         # Node.js
npm --version          # npm
```

## Notes

- The container uses a non-root user `vscode` to avoid permission issues
- All tools are pre-configured and ready to use
- The post-create script automatically sets up all sub-projects
