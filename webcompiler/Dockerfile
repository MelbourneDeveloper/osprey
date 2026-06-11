# Build stage for the Osprey compiler (Rust: osprey-rs) + the pure-C runtime
FROM rust:1-bookworm AS osprey-builder

# Build dependencies: gcc/make for the C runtime archives, OpenSSL for the HTTP runtime
RUN apt-get update && apt-get upgrade -y \
    && apt-get install -y --no-install-recommends \
    git \
    make \
    gcc \
    build-essential \
    libssl-dev \
    openssl \
    pkg-config \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy the Rust workspace + the C runtime sources (root Makefile builds both)
COPY Cargo.toml Cargo.lock Makefile rustfmt.toml clippy.toml /workspace/
COPY crates /workspace/crates
COPY tree-sitter-osprey /workspace/tree-sitter-osprey
COPY compiler/runtime /workspace/compiler/runtime
WORKDIR /workspace

# Build the C runtime archives (compiler/lib/lib*_runtime.a)
RUN make _runtime

# Build the Rust compiler binary (release)
RUN cargo build --release -p osprey-cli

# Build stage for LSP server
FROM node:20-bookworm AS lsp-builder

# Copy the VSCode extension for LSP server build
COPY vscode-extension /workspace/vscode-extension
WORKDIR /workspace

# Install dependencies and build LSP server
RUN cd vscode-extension/server && npm ci && npm run compile

# Runtime stage
FROM node:20-bookworm

# Install runtime dependencies including LLVM, compiler toolchain, and utilities
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    llvm-14 \
    llvm-14-dev \
    clang-14 \
    gcc \
    build-essential \
    libssl-dev \
    libssl3 \
    openssl \
    pkg-config \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && update-alternatives --install /usr/bin/clang clang /usr/bin/clang-14 60

# Add LLVM14 tools to PATH so Osprey compiler can find them
ENV PATH="/usr/lib/llvm-14/bin:${PATH}"

# Create non-root user for security
RUN groupadd -g 1001 nodejs && \
    useradd -r -u 1001 -g nodejs osprey

# Create app directory
WORKDIR /app

# Copy built Osprey compiler from the build stage
COPY --from=osprey-builder /workspace/bin/osprey /usr/local/bin/osprey

# Copy built runtime libraries from the lib directory (not bin)
COPY --from=osprey-builder /workspace/lib/libfiber_runtime.a /usr/local/lib/libfiber_runtime.a
COPY --from=osprey-builder /workspace/lib/libhttp_runtime.a /usr/local/lib/libhttp_runtime.a
COPY --from=osprey-builder /workspace/lib/libwebsocket_runtime.a /usr/local/lib/libwebsocket_runtime.a
COPY --from=osprey-builder /workspace/lib/libsystem_runtime.a /usr/local/lib/libsystem_runtime.a

# Copy built LSP server
COPY --from=lsp-builder /workspace/vscode-extension/server /app/server

# Copy web compiler source
COPY webcompiler/package*.json ./
RUN npm ci --only=production && npm cache clean --force

COPY webcompiler/src ./src

# Create temp directory for compiler operations
RUN mkdir -p /tmp/osprey-temp && chown osprey:nodejs /tmp/osprey-temp

# Change ownership to non-root user
RUN chown -R osprey:nodejs /app

# Make sure the osprey binary is executable
RUN chmod +x /usr/local/bin/osprey

# Set environment variables
ENV DOCKER_ENV=true
ENV NODE_ENV=production
ENV PORT=3001

# Switch to non-root user
USER osprey

# Expose port
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3001/api || exit 1

# Start the server
CMD ["npm", "run", "server"] 