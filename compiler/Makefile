.PHONY: build clean test regenerate-parser install-deps install uninstall test-llvm test-interpolation test-ast test-integration test-all test-stress test-basic test-functions test-errors test-types lint lint-fix lint-install fiber-runtime http-runtime

# Install golangci-lint
lint-install:
	@echo "📦 Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.1.6; \
	else \
		echo "✅ golangci-lint already installed"; \
	fi

# Run golangci-lint
lint: lint-install
	@echo "🔍 Running golangci-lint..."
	golangci-lint run --config .golangci.yml

# Run golangci-lint with auto-fix
lint-fix: lint-install
	@echo "🔧 Running golangci-lint with auto-fix..."
	golangci-lint run --config .golangci.yml --fix

# Build the osprey compiler with linting
build: lint fiber-runtime http-runtime
	@echo "🏗️  Building osprey compiler..."
	go build -o bin/osprey ./cmd/osprey

# Don't ever user this unless it's a life or death situation
build-no-lint: fiber-runtime http-runtime
	@echo "🏗️  Building osprey compiler (skipping lint for speed)..."
	go build -o bin/osprey ./cmd/osprey

# Build fiber runtime library
fiber-runtime:
	@echo "🔨 Building fiber runtime library..."
	@mkdir -p bin
	gcc -c -fPIC -O2 runtime/fiber_runtime.c -o bin/fiber_runtime.o
	ar rcs bin/libfiber_runtime.a bin/fiber_runtime.o
	@echo "✅ Fiber runtime library built: bin/libfiber_runtime.a"

# Build HTTP runtime library with TITANIUM-ARMORED MILITARY-GRADE SECURITY 🛡️💀
http-runtime:
	@echo "🔨 Building HTTP runtime with OpenSSL 3.5.0+ support..."
	@echo "⚠️  DEPLOYING MILITARY-GRADE PROTECTIONS..."
	@mkdir -p bin
	gcc -c -fPIC -O2 \
		-D_FORTIFY_SOURCE=2 \
		-fstack-protector-strong \
		-Wformat \
		-Werror=format-security \
		-Werror=implicit-function-declaration \
		-Werror=incompatible-pointer-types \
		-Werror=int-conversion \
		-Wall -Wextra \
		-Warray-bounds \
		-ftrapv \
		-fno-delete-null-pointer-checks \
		-fno-strict-overflow \
		-fno-strict-aliasing \
		-fPIE \
		-DWITH_OPENSSL \
		-DOPENSSL_SUPPRESS_DEPRECATED \
		-DOPENSSL_API_COMPAT=30000 \
		-Wno-deprecated-declarations \
		`pkg-config --cflags openssl 2>/dev/null || echo ""` \
		-o bin/http_shared.o runtime/http_shared.c
	gcc -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Wextra -ftrapv -fPIE -DOPENSSL_SUPPRESS_DEPRECATED -DOPENSSL_API_COMPAT=30000 -Wno-deprecated-declarations `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/http_client_runtime.c -o bin/http_client_runtime.o
	gcc -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Wextra -ftrapv -fPIE -DOPENSSL_SUPPRESS_DEPRECATED -DOPENSSL_API_COMPAT=30000 -Wno-deprecated-declarations `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/http_server_runtime.c -o bin/http_server_runtime.o
	gcc -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Wextra -ftrapv -fPIE -DOPENSSL_SUPPRESS_DEPRECATED -DOPENSSL_API_COMPAT=30000 -Wno-deprecated-declarations `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/websocket_client_runtime.c -o bin/websocket_client_runtime.o
	gcc -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Wextra -ftrapv -fPIE -DOPENSSL_SUPPRESS_DEPRECATED -DOPENSSL_API_COMPAT=30000 -Wno-deprecated-declarations `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/websocket_server_runtime.c -o bin/websocket_server_runtime.o
	gcc -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Wextra -ftrapv -fPIE runtime/system_runtime.c -o bin/system_runtime.o
	gcc -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Wall -Wextra -ftrapv -fPIE runtime/bridge_stubs.c -o bin/bridge_stubs.o
	ar rcs bin/libhttp_runtime.a bin/http_shared.o bin/http_client_runtime.o bin/http_server_runtime.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o bin/system_runtime.o bin/fiber_runtime.o bin/bridge_stubs.o
	@echo "💀 HTTP RUNTIME LIBRARY ASSEMBLED WITH MAXIMUM FIREPOWER! 💀"

# Install the osprey compiler globally
install: build
	@echo "🚀 Installing osprey compiler to /usr/local/bin..."
	sudo cp bin/osprey /usr/local/bin/
	@echo "📁 Creating /usr/local/lib if it doesn't exist..."
	sudo mkdir -p /usr/local/lib
	sudo cp bin/libfiber_runtime.a /usr/local/lib/
	sudo cp bin/libhttp_runtime.a /usr/local/lib/
	@echo "✅ osprey compiler installed! You can now use 'osprey' from anywhere."
	@echo "Test with: osprey --help"

# Uninstall the osprey compiler
uninstall:
	@echo "🗑️  Removing osprey compiler from /usr/local/bin..."
	sudo rm -f /usr/local/bin/osprey
	sudo rm -f /usr/local/lib/libfiber_runtime.a
	sudo rm -f /usr/local/lib/libhttp_runtime.a
	@echo "✅ osprey compiler uninstalled."

# Install dependencies
install-deps:
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf outputs/
	rm -rf internal/codegen/bin
	find /tmp -name "*osprey*" -delete 2>/dev/null || true
	find /tmp -name "*TestHTTP*" -delete 2>/dev/null || true
	find /tmp -name "*TestManual*" -delete 2>/dev/null || true

# ========== TEST SUITE ==========

# Run runtime tests (C tests that verify no segfaults)
test-runtime: build
	@echo "🧪 Running C Runtime Tests to verify segfault fixes..."
	@echo "📦 Compiling fiber runtime tests..."
	gcc -o runtime/fiber_runtime_test runtime/fiber_runtime_tests.c runtime/fiber_runtime.c -pthread -std=gnu11
	@echo "🚀 Running fiber runtime tests..."
	./runtime/fiber_runtime_test
	@echo ""
	@echo "📦 Compiling HTTP runtime tests..."
	gcc -o runtime/http_runtime_test runtime/http_runtime_tests.c \
		runtime/http_shared.c \
		runtime/http_client_runtime.c \
		runtime/http_server_runtime.c \
		runtime/websocket_client_runtime.c \
		runtime/websocket_server_runtime.c \
		runtime/fiber_runtime.c \
		runtime/bridge_stubs.c \
		-pthread `pkg-config --cflags --libs openssl 2>/dev/null || echo "-lssl -lcrypto"`
	@echo "🚀 Running HTTP runtime tests..."
	@echo "⚠️  NOTE: HTTP server tests use stub bridge functions - web API functionality not tested"
	./runtime/http_runtime_test
	@echo ""
	@echo "🧹 Cleaning up test executables..."
	rm -f runtime/fiber_runtime_test runtime/http_runtime_test
	@echo "✅ Runtime tests passed - no segfaults detected!"

# Run all tests without coverage (fast feedback)
test: build test-runtime
	@go test ./... -p 1 -v

# Regenerate parser from grammar (requires ANTLR)
regenerate-parser:
	antlr -Dlanguage=Go -package parser -o parser osprey.g4

# Run the parser on a file
run:
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	go run cmd/osprey/main.go $(FILE)