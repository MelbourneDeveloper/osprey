# Osprey Dev Container

This directory contains the development container configuration for the Osprey programming language project.

## Fixed Issues

✅ **Port Conflict Resolution**: The dev container now uses dynamic port allocation to avoid conflicts with other services running on ports 3001 and 8080.

## Quick Start

### ARM64 (Apple Silicon Macs)
```bash
# Use the convenient startup script
./.devcontainer/start-dev-container.sh

# Or manually with Docker Compose
docker compose -f .devcontainer/docker-compose.yml --profile arm64 up -d
```

### AMD64 (GitHub Codespaces / x64 systems)
```bash
docker compose -f .devcontainer/docker-compose.yml --profile amd64 up -d
```

## Features

- **Multi-architecture support**: ARM64 and AMD64 configurations
- **Dynamic port allocation**: No more port conflicts!
- **Complete development environment**: Go, Rust, Node.js, Java, ANTLR, C/C++ tools
- **VS Code integration**: Pre-configured extensions and settings
- **Persistent volumes**: Go cache, Cargo cache, and VS Code extensions

## Port Information

The dev container exposes two ports with dynamic allocation:
- **Port 3001**: Web Compiler service
- **Port 8080**: Development server

When you start the container, the actual port numbers will be displayed in the output.

## Troubleshooting

If you encounter any issues:

1. **Container won't start**: Try rebuilding the image:
   ```bash
   docker compose -f .devcontainer/docker-compose.yml --profile arm64 build --no-cache
   ```

2. **Port conflicts**: The new dynamic port allocation should prevent this, but if you still have issues, stop conflicting containers:
   ```bash
   docker ps | grep :3001
   docker stop <container-id>
   ```

3. **Clean restart**: Stop and remove everything:
   ```bash
   docker compose -f .devcontainer/docker-compose.yml --profile arm64 down
   docker system prune -f
   ```

## Files

- `devcontainer.json` - ARM64 dev container configuration
- `devcontainer-amd64.json` - AMD64 dev container configuration  
- `docker-compose.yml` - Multi-architecture container definitions
- `Dockerfile` - Container image definition
- `start-dev-container.sh` - Convenient startup script for ARM64
- `startup.sh` - Legacy startup script (deprecated)
