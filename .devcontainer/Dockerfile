# Multi-stage build for dev container
# Base stage with common tools
FROM golang:1.24-bookworm AS base

# Install basic tools and build dependencies
RUN echo 'Acquire::AllowInsecureRepositories "true";' > /etc/apt/apt.conf.d/99allow-insecure \
    && echo 'Acquire::AllowDowngradeToInsecureRepositories "true";' >> /etc/apt/apt.conf.d/99allow-insecure \
    && echo 'APT::Get::AllowUnauthenticated "true";' >> /etc/apt/apt.conf.d/99allow-insecure \
    && rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/* \
    && apt-get update --allow-unauthenticated \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends --allow-unauthenticated -o APT::Keep-Downloaded-Packages=false \
        git \
        curl \
        wget \
        build-essential \
        ca-certificates \
        gnupg \
        lsb-release \
        apt-transport-https \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/* \
    && mkdir -p /etc/apt/keyrings

# Node.js stage - install Node.js and npm
FROM base AS node-stage

# Install Node.js 20.x with GPG key handling
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y nodejs

# Rust stage - install Rust toolchain
FROM base AS rust-stage

# Install Rust
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y \
    && . ~/.cargo/env \
    && rustup component add clippy rustfmt

# Java stage - install Java and ANTLR
FROM base AS java-stage

RUN apt-get update --allow-unauthenticated \
    && apt-get install -y --no-install-recommends --allow-unauthenticated \
        openjdk-17-jdk \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install ANTLR
RUN wget -q https://www.antlr.org/download/antlr-4.13.2-complete.jar -O /usr/local/lib/antlr.jar \
    && echo '#!/bin/sh\njava -Xmx500M -cp "/usr/local/lib/antlr.jar:$CLASSPATH" org.antlr.v4.Tool "$@"' > /usr/local/bin/antlr \
    && echo '#!/bin/sh\njava -Xmx500M -cp "/usr/local/lib/antlr.jar:$CLASSPATH" org.antlr.v4.gui.TestRig "$@"' > /usr/local/bin/grun \
    && chmod +x /usr/local/bin/antlr /usr/local/bin/grun

# Final development stage
FROM base AS development

# Install Node.js directly in development stage with GPG key handling
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install -y nodejs

# Copy Rust installation from rust-stage
COPY --from=rust-stage /root/.cargo /root/.cargo
COPY --from=rust-stage /root/.rustup /root/.rustup

# Copy Java and ANTLR from java-stage
COPY --from=java-stage /usr/lib/jvm /usr/lib/jvm
COPY --from=java-stage /usr/local/lib/antlr.jar /usr/local/lib/antlr.jar
COPY --from=java-stage /usr/local/bin/antlr /usr/local/bin/antlr
COPY --from=java-stage /usr/local/bin/grun /usr/local/bin/grun

# Install additional dev tools that weren't in base
RUN apt-get update --allow-unauthenticated \
    && apt-get install -y --no-install-recommends --allow-unauthenticated \
        sudo \
        vim \
        nano \
        tree \
        less \
        htop \
        procps \
        gdb \
        valgrind \
        llvm-14 \
        clang-14 \
        clangd-14 \
        llvm-14-dev \
        libclang-14-dev \
        llvm-14-runtime \
        libllvm14 \
        libclang-cpp14 \
        gcc \
        build-essential \
        libssl-dev \
        libssl3 \
        openssl \
        libcrypto++-dev \
        libcrypto++8 \
        pkg-config \
        jq \
        unzip \
        libcunit1-dev \
        check \
        cmake \
        ninja-build \
        gcovr \
        lcov \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && update-alternatives --install /usr/bin/clang clang /usr/bin/clang-14 60 \
    && update-alternatives --install /usr/bin/clang++ clang++ /usr/bin/clang++-14 60

# Install Docker CLI (client only - talks to host daemon via socket)
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce-cli \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | gpg --dearmor -o /etc/apt/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
    && apt-get update \
    && apt-get install gh -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set up environment variables
ENV PATH="/root/.cargo/bin:/usr/lib/llvm-14/bin:/usr/lib/jvm/java-17-openjdk-arm64/bin:${PATH}"
ENV CLASSPATH=".:/usr/local/lib/antlr.jar:$CLASSPATH"
ENV JAVA_HOME="/usr/lib/jvm/java-17-openjdk-arm64"

# Set working directory
WORKDIR /workspace

# Expose ports for development servers
EXPOSE 8080 3001

# Install Go tools for development
RUN go install golang.org/x/tools/gopls@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest \
    && go install golang.org/x/tools/cmd/goimports@latest

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6

RUN sudo npm install -g \
  @anthropic-ai/claude-code 

# Install additional C testing tools and frameworks
RUN git clone https://github.com/ThrowTheSwitch/Unity.git /tmp/unity \
    && cp /tmp/unity/src/unity.c /usr/local/src/ \
    && cp /tmp/unity/src/unity.h /usr/local/include/ \
    && cp /tmp/unity/src/unity_internals.h /usr/local/include/ \
    && rm -rf /tmp/unity

# Create test discovery helper script
RUN echo '#!/bin/bash\n\
echo "🔍 Discovering C Tests..."\n\
find /workspace -name "test_*_unity" -executable | while read test; do\n\
    echo "Found C test: $test"\n\
    if [ -x "$test" ]; then\n\
        echo "  ✅ Executable"\n\
    else\n\
        echo "  ❌ Not executable"\n\
    fi\n\
done\n\
echo "🧪 Running C Runtime Tests..."\n\
cd /workspace/compiler/runtime\n\
if [ -f "./run_tests.sh" ]; then\n\
    ./run_tests.sh\n\
else\n\
    echo "❌ run_tests.sh not found"\n\
fi' > /usr/local/bin/discover-c-tests \
    && chmod +x /usr/local/bin/discover-c-tests

# Install Fly CLI
RUN curl -L https://fly.io/install.sh | sh \
    && mv /root/.fly/bin/flyctl /usr/local/bin/fly \
    && chmod +x /usr/local/bin/fly

# Set up Git configuration for development
RUN git config --global init.defaultBranch main \
    && git config --global advice.detachedHead false

