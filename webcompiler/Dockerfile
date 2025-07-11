# Build stage for Osprey compiler
FROM golang:1.23-bookworm AS osprey-builder

# Install build dependencies including gcc for the fiber runtime and OpenSSL for HTTP runtime
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

# Copy Osprey source code from the compiler directory
COPY compiler /workspace
WORKDIR /workspace

# Build ALL 4 runtime libraries (fiber, http, websocket, system)
RUN make fiber-runtime http-runtime websocket-runtime system-runtime

# VALIDATE that all 4 libraries were built correctly in the builder stage
RUN echo "🔍 Validating built libraries in builder stage..." && \
    ls -la lib/ && \
    test -f lib/libfiber_runtime.a || (echo "❌ Failed to build libfiber_runtime.a" && exit 1) && \
    test -f lib/libhttp_runtime.a || (echo "❌ Failed to build libhttp_runtime.a" && exit 1) && \
    test -f lib/libwebsocket_runtime.a || (echo "❌ Failed to build libwebsocket_runtime.a" && exit 1) && \
    test -f lib/libsystem_runtime.a || (echo "❌ Failed to build libsystem_runtime.a" && exit 1) && \
    echo "✅ All 4 runtime libraries built successfully in builder stage!"

# Build the Osprey compiler directly without linting
RUN go build -o bin/osprey ./cmd/osprey

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

# Create app directory and bin/lib subdirectories
WORKDIR /app
RUN mkdir -p bin lib temp && chmod 755 temp

# Copy built Osprey compiler from the build stage
COPY --from=osprey-builder /workspace/bin/osprey /usr/local/bin/osprey

# Copy ALL 4 runtime libraries to /usr/local/lib/ where osprey expects them (FHS standard)
COPY --from=osprey-builder /workspace/lib/libfiber_runtime.a /usr/local/lib/libfiber_runtime.a
COPY --from=osprey-builder /workspace/lib/libhttp_runtime.a /usr/local/lib/libhttp_runtime.a
COPY --from=osprey-builder /workspace/lib/libwebsocket_runtime.a /usr/local/lib/libwebsocket_runtime.a
COPY --from=osprey-builder /workspace/lib/libsystem_runtime.a /usr/local/lib/libsystem_runtime.a

# VALIDATE that all runtime libraries are actually present
RUN echo "🔍 Validating runtime libraries..." && \
    ls -la /usr/local/lib/lib*runtime.a && \
    test -f /usr/local/lib/libfiber_runtime.a || (echo "❌ Missing libfiber_runtime.a" && exit 1) && \
    test -f /usr/local/lib/libhttp_runtime.a || (echo "❌ Missing libhttp_runtime.a" && exit 1) && \
    test -f /usr/local/lib/libwebsocket_runtime.a || (echo "❌ Missing libwebsocket_runtime.a" && exit 1) && \
    test -f /usr/local/lib/libsystem_runtime.a || (echo "❌ Missing libsystem_runtime.a" && exit 1) && \
    echo "✅ All 4 runtime libraries validated successfully!"

# Copy built LSP server
COPY --from=lsp-builder /workspace/vscode-extension/server /app/server

# Copy web compiler source
COPY webcompiler/package*.json ./
RUN npm ci --only=production && npm cache clean --force

COPY webcompiler/src ./src

# Change ownership to non-root user (including temp directory)
RUN chown -R osprey:nodejs /app && chmod -R 755 /app/temp

# Make sure the osprey binary is executable
RUN chmod +x /usr/local/bin/osprey

# Set environment variables
ENV DOCKER_ENV=true
ENV NODE_ENV=production

# Switch to non-root user
USER osprey

# Expose port
EXPOSE 3001

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3001/api || exit 1

# Start the server
CMD ["npm", "run", "server"] 