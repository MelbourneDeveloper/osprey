# Osprey Development Container

This directory contains configuration for a development container that provides all necessary dependencies for the Osprey compiler and VS Code extension development.
This directory contains configuration for a development container that provides all necessary dependencies for Osprey compiler and VS Code extension development.

## Features

- **Go 1.23.4** for compiler development
- **LLVM 14** for IR generation and compilation  
- **ANTLR 4.13.1** for parser generation
- **Node.js 20.17.0** for VS Code extension development
- **Rust** for the rust_integration examples
- **Java 17** for ANTLR runtime
- **All necessary VS Code extensions** pre-installed

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop)
- [Visual Studio Code](https://code.visualstudio.com/)
- [Dev Containers extension for VS Code](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)

### Opening the Project in the Dev Container

1. Open the `osprey.code-workspace` file in VS Code
2. When prompted to "Reopen in Container", click "Reopen in Container"
   - Alternatively, press F1, type "Dev Containers: Reopen in Container" and press Enter

The container will automatically run the post-create script which sets up both the compiler and VS Code extension.

## Available Scripts

### Individual Scripts
- `build-compiler.sh` - Builds the Osprey compiler
- `test-setup.sh` - Tests that all tools are working correctly

### Running Scripts
```bash
# Build the compiler
.devcontainer/build-compiler.sh

# Build VS Code extension (use the proper build script)
cd vscode-extension && ./build.sh
```

## Development Tasks

### Compiler Development

Navigate to the compiler directory and use make commands:
```bash
cd /compiler
make build          # Build the compiler
make test           # Run all tests
make test-llvm      # Run LLVM tests only
make lint           # Run linter
make clean          # Clean build artifacts
```

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
├── compiler/           # Osprey compiler (Go + ANTLR + LLVM)
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
The container uses Node.js 20.17.0 which is compatible with the latest npm. If you encounter version conflicts, the setup scripts handle the compatibility automatically.

## Manual Testing

After the container starts, you can verify everything works:

```bash
# Test individual components
go version              # Go
antlr -version         # ANTLR  
llc --version          # LLVM
node --version         # Node.js
npm --version          # npm
rustc --version        # Rust
```

## Notes

- The container uses a non-root user `vscode` to avoid permission issues
- All tools are pre-configured and ready to use
- The post-create script automatically sets up both projects
