# Simple, minimal Dockerfile for Osprey development (Rust compiler + C runtime)
FROM debian:bookworm

# Install basic system packages
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    wget \
    git \
    build-essential \
    cmake \
    clang-14 \
    llvm-14 \
    llvm-14-dev \
    libssl-dev \
    pkg-config \
    vim \
    nano \
    tree \
    htop \
    procps \
    jq \
    unzip \
    zip \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs

# Install Rust (the compiler toolchain for crates/)
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y \
    && . ~/.cargo/env \
    && rustup component add clippy rustfmt llvm-tools-preview

# Create vscode user
RUN apt-get update && apt-get install -y sudo \
    && groupadd --gid 1000 vscode \
    && useradd --uid 1000 --gid vscode --shell /bin/bash --create-home vscode \
    && mkdir -p /etc/sudoers.d \
    && echo 'vscode ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers.d/50-vscode \
    && mkdir -p /home/vscode/.cargo /home/vscode/.rustup \
    && cp -r /root/.cargo/* /home/vscode/.cargo/ || true \
    && cp -r /root/.rustup/* /home/vscode/.rustup/ || true \
    && chown -R vscode:vscode /home/vscode \
    && echo 'source ~/.cargo/env' >> /home/vscode/.bashrc

# Set environment variables
ENV PATH="/home/vscode/.cargo/bin:/usr/lib/llvm-14/bin:${PATH}"

# Set working directory
WORKDIR /workspace

# Switch to vscode user
USER vscode

# Expose ports for development
EXPOSE 8080 3001
