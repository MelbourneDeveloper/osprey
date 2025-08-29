# Simple, minimal Dockerfile for Osprey development
FROM golang:1.23-bookworm

# Install basic system packages
RUN apt-get update && apt-get install -y \
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
    openjdk-17-jdk \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install Node.js
RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y nodejs

# Continue with the same image for Rust installation

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y \
    && . ~/.cargo/env \
    && rustup component add clippy rustfmt

# Install ANTLR
RUN wget -q https://www.antlr.org/download/antlr-4.13.2-complete.jar -O /usr/local/lib/antlr.jar \
    && echo '#!/bin/sh\njava -Xmx500M -cp "/usr/local/lib/antlr.jar:$CLASSPATH" org.antlr.v4.Tool "$@"' > /usr/local/bin/antlr \
    && echo '#!/bin/sh\njava -Xmx500M -cp "/usr/local/lib/antlr.jar:$CLASSPATH" org.antlr.v4.gui.TestRig "$@"' > /usr/local/bin/grun \
    && chmod +x /usr/local/bin/antlr /usr/local/bin/grun

# Install Go tools
RUN go install golang.org/x/tools/gopls@v0.16.2 \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install golang.org/x/tools/cmd/goimports@latest

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0

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
ENV JAVA_HOME="/usr/lib/jvm/java-17-openjdk-arm64"
ENV CLASSPATH=".:/usr/local/lib/antlr.jar"
ENV PATH="/home/vscode/.cargo/bin:/usr/lib/llvm-14/bin:${PATH}"

# Set working directory
WORKDIR /workspace

# Switch to vscode user
USER vscode

# Expose ports for development
EXPOSE 8080 3001