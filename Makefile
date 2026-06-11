# =============================================================================
# Osprey — single root Makefile (Rust toolchain: osprey-rs).
#
# The compiler is the Rust workspace in crates/. `osprey-rs --run` emits LLVM IR
# and hands it to clang together with the prebuilt C runtime archives
# (compiler/bin/lib*_runtime.a) — those archives are pure C (no Go), built by
# the private `_runtime` helper below.
#
# Public targets are the handful documented with `## name:` lines. Everything
# else (the `_`-prefixed rules) is an internal helper — not meant to be run by
# hand. Cross-platform via GNU Make.
# =============================================================================

.PHONY: build tui run test lint fmt clean ci setup install vsix
.PHONY: BUILD TUI RUN TEST LINT FMT CLEAN CI SETUP INSTALL VSIX

# Case-insensitive convenience: `make TUI` == `make tui`, etc.
BUILD: build
TUI: tui
RUN: run
TEST: test
LINT: lint
FMT: fmt
CLEAN: clean
CI: ci
SETUP: setup
INSTALL: install
VSIX: vsix

# ---------------------------------------------------------------------------
# OS detection
# ---------------------------------------------------------------------------
ifeq ($(OS),Windows_NT)
  SHELL := powershell.exe
  .SHELLFLAGS := -NoProfile -Command
  RM = Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
else
  SHELL := /bin/bash
  RM = rm -rf
endif

# ---------------------------------------------------------------------------
# Toolchain / paths
# ---------------------------------------------------------------------------
# BIN: the built CLI. RTB: archive output dir (osprey-rs searches compiler/bin).
CC  ?= cc
AR  ?= ar
BIN := target/release/osprey-rs
RTB := compiler/bin
TUI_DEMO := compiler/examples/tui/api_browser.osp

# VSIX (VSCode extension) — macOS only. Bundles the Rust binary as `osprey`.
EXT_DIR         := vscode-extension
EXT_ID          := nimblesite.osprey
VSCODE_STORAGE  := $(HOME)/Library/Application Support/Code/User/globalStorage/storage.json

# C runtime compile flag profiles (mirror the original hardened recipes).
A    := -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Werror -Wall -Wextra -ftrapv -fPIE -D_GNU_SOURCE
B    := $(A) -std=c11
OSSL := -DOPENSSL_SUPPRESS_DEPRECATED -DOPENSSL_API_COMPAT=30000 -Wno-deprecated-declarations
# Object lists for the archives (paths relative to compiler/, where `ar` runs).
FIB_OBJ  := bin/fiber_runtime.o bin/system_runtime.o bin/effects_runtime.o bin/string_runtime.o bin/string_runtime_list.o bin/list_runtime.o bin/map_runtime.o bin/map_runtime_hamt.o bin/json_runtime.o bin/ffi_runtime.o bin/term_runtime.o
HTTP_OBJ := bin/http_shared.o bin/http_client_runtime.o bin/http_server_runtime.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o $(FIB_OBJ)

# =============================================================================
# Public targets
# =============================================================================

## build: Build the C runtime archives + the Rust workspace (release).
build: _runtime
	@echo "==> cargo build --release --workspace"
	cargo build --release --workspace

## tui: Build, then launch the interactive TUI demo (live GitHub API browser).
##      Runs in the current terminal so the raw-mode key reader gets a real stdin.
tui: build
	@echo "==> launching TUI: $(TUI_DEMO)"
	./$(BIN) $(TUI_DEMO) --run

## run: Compile and run an Osprey file (usage: make run FILE=<path>).
run: build
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	./$(BIN) $(FILE) --run

## test: cargo tests + the differential golden harness (osprey-rs --run vs
##       each example's .expectedoutput, byte-for-byte).
test: build
	@echo "==> cargo test --workspace"
	cargo test --workspace
	@echo "==> differential golden tests"
	@out=$$(zsh crates/diff_examples.sh); echo "$$out"; \
	  echo "$$out" | grep -Eq 'FAIL=0 '  || { echo 'FAIL: differential mismatch'; exit 1; }; \
	  echo "$$out" | grep -Eq 'NOEXP=0 ' || { echo 'FAIL: example missing .expectedoutput'; exit 1; }

## lint: Format check + clippy at maximum strictness (warnings are errors).
lint:
	@echo "==> cargo fmt --check + clippy -D warnings"
	cargo fmt --all --check
	cargo clippy --workspace --all-targets -- -D warnings

## fmt: Format all Rust code in-place.
fmt:
	cargo fmt --all

## clean: Remove all build artifacts (Rust target + C runtime objects/archives).
clean:
	cargo clean
	$(RM) $(RTB) compiler/lib outputs

## ci: lint + test + build (full local CI simulation).
ci: lint test build

## setup: Install the Rust components the build needs.
setup:
	rustup component add rustfmt clippy
	cargo fetch
	@echo "==> setup complete. Run 'make ci' to validate."

## install: Install osprey-rs + the runtime archives system-wide (cargo + /usr/local/lib).
install: build
	cargo install --path crates/osprey-cli --force
	sudo mkdir -p /usr/local/lib
	sudo cp $(RTB)/libfiber_runtime.a $(RTB)/libhttp_runtime.a /usr/local/lib/
	@echo "==> installed osprey-rs and runtime archives."

## vsix: Build, bundle the Rust binary as the extension's compiler, package and
##       install the VSIX into every VSCode profile. macOS only.
vsix: build _vsix_clean _vsix_build _vsix_bundle _vsix_package _vsix_install

# =============================================================================
# Internal helpers — not public targets.
# =============================================================================

# Build the pure-C runtime archives that osprey-rs links at `--run` time. One
# shell so `cd` persists; faithful port of the original hardened C recipes.
_runtime:
	@echo "==> building C runtime archives ($(RTB)/lib*_runtime.a)"
	@cd compiler && set -e && mkdir -p bin lib && \
	  $(CC) -c -fPIC -O2 -Werror -Wall -Wextra -Wpedantic -std=c11 -D_GNU_SOURCE runtime/fiber_runtime.c -o bin/fiber_runtime.o && \
	  $(CC) $(A) runtime/system_runtime.c       -o bin/system_runtime.o && \
	  $(CC) $(A) runtime/effects_runtime.c      -o bin/effects_runtime.o && \
	  $(CC) $(A) runtime/string_runtime.c       -o bin/string_runtime.o && \
	  $(CC) $(A) runtime/string_runtime_list.c  -o bin/string_runtime_list.o && \
	  $(CC) $(B) runtime/list_runtime.c         -o bin/list_runtime.o && \
	  $(CC) $(B) runtime/map_runtime.c          -o bin/map_runtime.o && \
	  $(CC) $(B) runtime/map_runtime_hamt.c     -o bin/map_runtime_hamt.o && \
	  $(CC) $(B) runtime/json_runtime.c         -o bin/json_runtime.o && \
	  $(CC) $(B) runtime/ffi_runtime.c          -o bin/ffi_runtime.o && \
	  $(CC) $(B) runtime/term_runtime.c         -o bin/term_runtime.o && \
	  $(CC) -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Werror -Wall -Wextra \
	        -Wformat -Werror=format-security -Werror=implicit-function-declaration \
	        -Werror=incompatible-pointer-types -Werror=int-conversion -Warray-bounds -ftrapv \
	        -fno-delete-null-pointer-checks -fno-strict-overflow -fno-strict-aliasing -fPIE \
	        -DWITH_OPENSSL $(OSSL) `pkg-config --cflags openssl 2>/dev/null || echo ""` \
	        runtime/http_shared.c -o bin/http_shared.o && \
	  $(CC) $(A) $(OSSL) `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/http_client_runtime.c      -o bin/http_client_runtime.o && \
	  $(CC) $(A) $(OSSL) `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/http_server_runtime.c      -o bin/http_server_runtime.o && \
	  $(CC) $(A) $(OSSL) `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/websocket_client_runtime.c -o bin/websocket_client_runtime.o && \
	  $(CC) $(A) $(OSSL) `pkg-config --cflags openssl 2>/dev/null || echo ""` runtime/websocket_server_runtime.c -o bin/websocket_server_runtime.o && \
	  $(AR) rcs bin/libfiber_runtime.a $(FIB_OBJ) && \
	  $(AR) rcs bin/libhttp_runtime.a  $(HTTP_OBJ) && \
	  cp bin/libfiber_runtime.a bin/libhttp_runtime.a lib/

# --- vsix sub-steps --------------------------------------------------------
_vsix_clean:
	cd $(EXT_DIR) && $(RM) out dist *.vsix

_vsix_build:
	cd $(EXT_DIR) && npm run compile

# Stage the freshly-built Rust binary where the extension expects its bundled
# compiler (bin/<os>-<arch>/osprey), so the VSIX runs against THIS build.
_vsix_bundle:
	@OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	case "$$OS" in darwin) OS=darwin;; linux) OS=linux;; *) OS=win32;; esac; \
	ARCH=$$(uname -m); case "$$ARCH" in arm64|aarch64) ARCH=arm64;; *) ARCH=x64;; esac; \
	DEST="$(EXT_DIR)/bin/$$OS-$$ARCH"; mkdir -p "$$DEST"; \
	cp $(BIN) "$$DEST/osprey"; \
	echo "  bundled $(BIN) -> $$DEST/osprey"

_vsix_package:
	cd $(EXT_DIR) && npm run package

# Install the newest VSIX into the default profile + every named profile.
_vsix_install:
	@VSIX=$$(ls -t $(EXT_DIR)/*.vsix 2>/dev/null | head -1); \
	if [ -z "$$VSIX" ]; then echo "FAIL: no .vsix in $(EXT_DIR)/"; exit 1; fi; \
	echo "  vsix: $$VSIX"; \
	code --install-extension "$$VSIX" --force && echo "  [default] installed"; \
	jq -r '.userDataProfiles[]?.name' "$(VSCODE_STORAGE)" 2>/dev/null | while IFS= read -r prof; do \
	  [ -z "$$prof" ] && continue; \
	  code --profile "$$prof" --install-extension "$$VSIX" --force && echo "  [$$prof] installed"; \
	done
