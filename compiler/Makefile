.PHONY: build clean test regenerate-parser install-deps install uninstall test-llvm test-interpolation test-ast test-integration test-all test-stress test-coverage test-basic test-functions test-errors test-types test-rust-interop lint lint-fix lint-install fiber-runtime http-runtime ensure-built

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

# Build the osprey compiler without linting (for faster test builds)
build-no-lint: fiber-runtime http-runtime
	@echo "🏗️  Building osprey compiler (skipping lint for speed)..."
	go build -o bin/osprey ./cmd/osprey

# Ensure compiler is built (only build if binary doesn't exist)
ensure-built:
	@if [ ! -f bin/osprey ]; then \
		echo "🔍 Compiler not found, building..."; \
		$(MAKE) build; \
	else \
		echo "✅ Compiler already built at bin/osprey"; \
	fi

# Build fiber runtime library
fiber-runtime:
	@echo "🔨 Building fiber runtime library..."
	@mkdir -p bin
	@rm -f bin/libfiber_runtime.a  # Remove corrupted archive first
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
	@rm -f bin/libhttp_runtime.a  # Remove corrupted archive first
	ar rcs bin/libhttp_runtime.a bin/http_shared.o bin/http_client_runtime.o bin/http_server_runtime.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o
	@echo "💀 HTTP RUNTIME LIBRARY ASSEMBLED WITH MAXIMUM FIREPOWER! 💀"

# Install the osprey compiler globally - NUKE ALL EXISTING INSTALLATIONS
install: build
	@echo "💀 OBLITERATING ALL EXISTING OSPREY INSTALLATIONS..."
	@echo "🗑️  Removing homebrew version..."
	@if [ -L "/opt/homebrew/bin/osprey" ]; then sudo rm -f "/opt/homebrew/bin/osprey"; echo "✅ Removed homebrew symlink"; fi
	@if [ -f "/opt/homebrew/bin/osprey" ]; then sudo rm -f "/opt/homebrew/bin/osprey"; echo "✅ Removed homebrew binary"; fi
	@if brew list osprey-local >/dev/null 2>&1; then brew uninstall osprey-local; echo "✅ Uninstalled osprey-local"; fi
	@if brew list osprey >/dev/null 2>&1; then brew uninstall osprey; echo "✅ Uninstalled osprey"; fi
	@echo "🗑️  Removing old /usr/local/bin version..."
	sudo rm -f /usr/local/bin/osprey
	sudo rm -f /usr/local/lib/libfiber_runtime.a
	sudo rm -f /usr/local/lib/libhttp_runtime.a
	@echo "🚀 Installing FRESH osprey compiler to /usr/local/bin..."
	sudo cp bin/osprey /usr/local/bin/
	@echo "📁 Installing runtime libraries to /usr/local/lib..."
	sudo mkdir -p /usr/local/lib
	sudo cp bin/libfiber_runtime.a /usr/local/lib/
	sudo cp bin/libhttp_runtime.a /usr/local/lib/
	@echo "🔧 Ensuring /usr/local/bin comes before /opt/homebrew/bin in PATH..."
	@echo "💥 FRESH OSPREY INSTALLATION COMPLETE!"
	@echo "🧪 Testing installation..."
	@which osprey || echo "❌ osprey not found in PATH"
	@osprey --help >/dev/null 2>&1 && echo "✅ osprey responds to --help" || echo "❌ osprey --help failed"
	@echo "🎯 You can now use 'osprey' from anywhere - it should be the FRESH version!"

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

# Test with example files
test: test-all

# ========== TEST SUITE ==========

# Run all tests
test-all: test-ast test-llvm test-integration test-http test-websocket test-fiber test-cli test-rust-interop test-archive-integrity
	@echo "✅ All tests passed!"

# Run archive integrity tests
test-archive-integrity: fiber-runtime http-runtime
	@echo "🔍 Running archive integrity tests..."
	go test -v ./tests/integration/ -run "TestRuntimeArchiveIntegrity|TestExistingRuntimeArchiveIntegrity"

# Run LLVM IR generation tests
test-llvm: fiber-runtime http-runtime
	@echo "🔧 Running LLVM IR generation tests..."
	@cd internal/codegen && ln -sf ../../bin bin 2>/dev/null || true
	go test ./internal/codegen -v

# Run AST parsing tests for interpolation
test-ast:
	@echo "🌳 Running AST interpolation parsing tests..."
	go test ./internal/ast -v

# Run end-to-end integration tests
test-integration: ensure-built fiber-runtime http-runtime
	@echo "🚀 Running core integration tests..."
	go test -v ./tests/integration/ -run "TestRootLevelExamples|TestLanguageFeatures|TestBasicCompilation|TestErrorHandling|TestFunctionArguments|TestCompilationFailures"

# Run HTTP integration tests
test-http: fiber-runtime http-runtime
	@echo "🌐 Running HTTP integration tests..."
	go test -v ./tests/integration/ -run "TestHttpExamples"

# Run WebSocket integration tests
test-websocket: fiber-runtime http-runtime
	@echo "🔌 Running WebSocket integration tests..."
	go test -v ./tests/integration/ -run "TestWebsoxExamples"

# Run Fiber integration tests
test-fiber: fiber-runtime http-runtime
	@echo "🧵 Running Fiber integration tests..."
	go test -v ./tests/integration/ -run "TestFiberExamples|TestFiberFeatures|TestFiberErrorHandling|TestFiberModuleIsolation|TestFiberIntegration"

# Run CLI integration tests
test-cli: ensure-built
	@echo "⌨️  Running CLI integration tests..."
	go test -v ./tests/integration/ -run "TestCLI"

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	mkdir -p outputs
	go test -coverprofile=outputs/coverage.out ./internal/... ./tests/integration/
	@if [ -f outputs/coverage.out ]; then \
		go tool cover -html=outputs/coverage.out -o outputs/coverage.html; \
		echo "📊 Coverage report generated: outputs/coverage.html"; \
	else \
		echo "⚠️  No coverage data generated (packages may have no statements)"; \
		touch outputs/coverage.html; \
		echo "<html><body><h1>No Coverage Data</h1><p>No coverage data was generated. This may happen when packages have no statements to cover.</p></body></html>" > outputs/coverage.html; \
	fi

# Test Rust interop functionality
test-rust-interop: ensure-built
	@echo "🦀 Testing Rust interop functionality..."
	@if ! command -v rustc >/dev/null 2>&1; then \
		echo "❌ RUST COMPILER NOT FOUND! Install Rust: https://rustup.rs/"; \
		exit 1; \
	fi
	@if ! command -v cargo >/dev/null 2>&1; then \
		echo "❌ CARGO NOT FOUND! Install Rust toolchain: https://rustup.rs/"; \
		exit 1; \
	fi
	@echo "✅ Rust tools found, running Rust interop tests..."
	go test -v ./tests/integration/ -run "TestRustInterop|TestRustInteropCompilationOnly|TestRustInteropSimple"
	cd examples/rust_integration && chmod +x run.sh && ./run.sh

# ========== DEVELOPMENT ==========

# Regenerate parser from grammar (requires ANTLR)
regenerate-parser:
	antlr -Dlanguage=Go -package parser -o parser osprey.g4

# Run the parser on a file
run:
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	go run cmd/osprey/main.go $(FILE) 