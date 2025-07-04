.PHONY: build clean test regenerate-parser install-deps install uninstall test-llvm test-interpolation test-ast test-integration test-all test-stress test-basic test-functions test-errors test-types lint lint-fix lint-install fiber-runtime http-runtime c-lint c-test test-runtime

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
build: lint c-lint fiber-runtime http-runtime
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
	@echo "🔧 Building http_shared.o..."
	@gcc -c -fPIC -O2 \
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
	ar rcs bin/libhttp_runtime.a bin/http_shared.o bin/http_client_runtime.o bin/http_server_runtime.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o bin/system_runtime.o bin/fiber_runtime.o
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
	gcc -o runtime/fiber_runtime_test runtime/fiber_runtime_tests.c runtime/fiber_runtime.c runtime/system_runtime.c -pthread -std=gnu11
	@echo "🚀 Running fiber runtime tests..."
	./runtime/fiber_runtime_test
	@echo ""
	@echo "📦 Creating HTTP test stub for missing bridge functions..."
	@echo '#include <stddef.h>' > runtime/http_bridge_stub.c
	@echo 'int osprey_handle_http_request(int server_id, char *method, char *full_url, char *raw_headers, char *body, size_t body_length, int *response_status, char **response_headers, char **response_body, size_t *response_body_length) { return -1; }' >> runtime/http_bridge_stub.c
	@echo 'char* osprey_handle_websocket_message(int server_id, int connection_id, char *message, size_t message_length) { return NULL; }' >> runtime/http_bridge_stub.c
	@echo 'int osprey_handle_websocket_connection(int server_id, int connection_id, char *client_ip) { return -1; }' >> runtime/http_bridge_stub.c
	@echo 'int osprey_handle_websocket_disconnect(int server_id, int connection_id) { return -1; }' >> runtime/http_bridge_stub.c
	@echo 'void osprey_free_websocket_response(char *response) { }' >> runtime/http_bridge_stub.c
	@echo "📦 Compiling HTTP runtime tests with stub functions..."
	gcc -o runtime/http_runtime_test runtime/http_runtime_tests.c \
		runtime/http_shared.c \
		runtime/http_client_runtime.c \
		runtime/http_server_runtime.c \
		runtime/websocket_client_runtime.c \
		runtime/websocket_server_runtime.c \
		runtime/fiber_runtime.c \
		runtime/system_runtime.c \
		runtime/http_bridge_stub.c \
		-pthread `pkg-config --cflags --libs openssl 2>/dev/null || echo "-lssl -lcrypto"`
	@echo "🚀 Running HTTP runtime tests..."
	@echo "⚠️  NOTE: HTTP server tests focus on core runtime functionality"
	./runtime/http_runtime_test
	@echo ""
	@echo "🧹 Cleaning up test executables..."
	rm -f runtime/fiber_runtime_test runtime/http_runtime_test runtime/http_bridge_stub.c
	@echo "✅ Runtime tests passed - no segfaults detected!"

# Run all tests without coverage (fast feedback)
test: build test-runtime c-test
	@go test ./... -p 1 -v

# MAXIMUM STRICTNESS C LINTING - ERRORS FOR EVERYTHING!!!
c-lint:
	@echo "💀 DEPLOYING MAXIMUM C STRICTNESS - NO MERCY!!! 💀"
	@echo "🔥 Running C linter with MILITARY-GRADE strictness..."
	@cd runtime && clang \
		-Wall -Wextra -Werror -Wpedantic \
		-Wformat=2 -Wconversion -Wsign-conversion \
		-Wnull-dereference -Wdouble-promotion \
		-Wshadow -Wunreachable-code -Wcast-qual \
		-Wwrite-strings \
		-Wstrict-prototypes -Wundef -Wredundant-decls \
		-Winit-self -Wswitch-default -Wswitch-enum \
		-Wfloat-equal -Wbad-function-cast \
		-Wstrict-overflow=5 -Wstack-protector \
		-Winline -Wdisabled-optimization \
		-std=c11 -fsyntax-only \
		system_runtime.c fiber_runtime.c test_system_runtime.c fiber_runtime_tests.c
	@echo "✅ C CODE IS BULLETPROOF! NO WARNINGS SURVIVED!"

# C TESTS - Vanilla C Tests with standard assert()
c-test:
	@echo "🧪 RUNNING C RUNTIME TESTS..."
	@cd runtime && echo "📦 Compiling system runtime tests..."
	@cd runtime && clang -o test_system_runtime test_system_runtime.c system_runtime.c -pthread -std=c11 -D_GNU_SOURCE -g
	@cd runtime && echo "🚀 Running system runtime tests..."
	@cd runtime && ./test_system_runtime
	@cd runtime && echo "📦 Compiling fiber runtime tests..."
	@cd runtime && clang -o test_fiber_runtime fiber_runtime_tests.c fiber_runtime.c system_runtime.c -pthread -std=c11 -D_GNU_SOURCE -g
	@cd runtime && echo "🚀 Running fiber runtime tests..."
	@cd runtime && ./test_fiber_runtime
	@cd runtime && echo "🧹 Cleaning up test executables..."
	@cd runtime && rm -f test_system_runtime test_fiber_runtime
	@echo "✅ ALL C RUNTIME TESTS PASSED!"

# Regenerate parser from grammar (requires ANTLR)
regenerate-parser:
	antlr -Dlanguage=Go -package parser -o parser osprey.g4

# Run the parser on a file
run:
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	go run cmd/osprey/main.go $(FILE)