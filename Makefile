# agent-pmo:b636503
# =============================================================================
# Standard Makefile — osprey
# Cross-platform: Linux, macOS, Windows (via GNU Make)
# Primary language: Rust (crates/ workspace → the osprey-rs compiler), with a
# pure-C runtime (compiler/runtime → lib*_runtime.a, linked by `osprey-rs
# --run`) and TypeScript sub-projects (vscode-extension, webcompiler, website).
# =============================================================================

.PHONY: build test lint fmt clean ci setup tui run install uninstall website-dev website-build rebuild-install-vsix

# ---------------------------------------------------------------------------
# OS Detection
# ---------------------------------------------------------------------------
ifeq ($(OS),Windows_NT)
  SHELL := powershell.exe
  .SHELLFLAGS := -NoProfile -Command
  RM = Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
  MKDIR = New-Item -ItemType Directory -Force
  HOME ?= $(USERPROFILE)
else
  # bash needed for `pipefail` in tee'd test recipes; Ubuntu's /bin/sh is dash.
  SHELL := /bin/bash
  RM = rm -rf
  MKDIR = mkdir -p
endif

# ---------------------------------------------------------------------------
# Variables. NOTE: `?=` (not `:=`) on purpose — the VSCode Makefile-Tools panel
# lists `:=` assignments as if they were targets; `?=` keeps the panel clean.
# ---------------------------------------------------------------------------
# Coverage — single source of truth is coverage-thresholds.json.
# See REPO-STANDARDS-SPEC [COVERAGE-THRESHOLDS-JSON].
COVERAGE_THRESHOLDS_FILE ?= coverage-thresholds.json

# Toolchain / paths. BIN: the built CLI. RTB: C-runtime archive output dir
# (osprey-rs searches compiler/bin at --run time).
CC  ?= cc
AR  ?= ar
BIN ?= target/release/osprey-rs
RTB ?= compiler/bin

# VSIX (VSCode extension) — macOS only. Bundles the Rust binary as `osprey`.
EXT_DIR        ?= vscode-extension
EXT_ID         ?= nimblesite.osprey
VSCODE_STORAGE ?= $(HOME)/Library/Application Support/Code/User/globalStorage/storage.json

# C runtime compile flag profiles (hardened; mirror the original recipes).
A    ?= -c -fPIC -O2 -D_FORTIFY_SOURCE=2 -fstack-protector-strong -Werror -Wall -Wextra -ftrapv -fPIE -D_GNU_SOURCE
B    ?= $(A) -std=c11
OSSL ?= -DOPENSSL_SUPPRESS_DEPRECATED -DOPENSSL_API_COMPAT=30000 -Wno-deprecated-declarations
# Object lists for the archives (paths relative to compiler/, where `ar` runs).
FIB_OBJ  ?= bin/fiber_runtime.o bin/system_runtime.o bin/effects_runtime.o bin/string_runtime.o bin/string_runtime_list.o bin/list_runtime.o bin/map_runtime.o bin/map_runtime_hamt.o bin/json_runtime.o bin/ffi_runtime.o bin/term_runtime.o
HTTP_OBJ ?= bin/http_shared.o bin/http_client_runtime.o bin/http_server_runtime.o bin/websocket_client_runtime.o bin/websocket_server_runtime.o $(FIB_OBJ)

# =============================================================================
# Standard Targets
# =============================================================================

## build: C runtime archives + Rust workspace (release) + VSCode extension
build: _runtime
	@echo "==> Building..."
	cargo build --release --workspace
	cd $(EXT_DIR) && npm run compile

## test: Fail-fast tests + coverage + per-project threshold enforcement.
##       See REPO-STANDARDS-SPEC [TEST-RULES] and [COVERAGE-THRESHOLDS-JSON].
##       Projects listed in coverage-thresholds.json are each tested + checked.
test: build
	@echo "==> Testing (fail-fast + coverage + per-project thresholds)..."
	$(MAKE) _test_rust
	$(MAKE) _coverage_check_rust
	$(MAKE) _test_differential
	$(MAKE) _test_vscode_extension
	$(MAKE) _coverage_check_vscode_extension

## lint: Run all linters/analyzers (read-only). Does NOT format.
lint:
	@echo "==> Linting..."
	cargo clippy --workspace --all-targets -- -D warnings
	cd $(EXT_DIR) && npm run lint

## fmt: Format all code in-place. Pass CHECK=1 for read-only check (CI use).
fmt:
	@echo "==> Formatting$(if $(CHECK), (check mode),)..."
	cargo fmt --all$(if $(CHECK), --check,)
	cd $(EXT_DIR) && npx prettier$(if $(CHECK), --check, --write) .

## clean: Remove all build artifacts
clean:
	@echo "==> Cleaning..."
	cargo clean
	$(RM) $(RTB) compiler/lib outputs lcov.info test.log
	cd $(EXT_DIR) && $(RM) out dist coverage test.log

## ci: lint + test + build (full CI simulation)
ci: lint test build

## setup: Post-create dev environment setup (used by devcontainer)
setup:
	@echo "==> Setting up development environment..."
	rustup component add rustfmt clippy llvm-tools-preview
	command -v cargo-llvm-cov >/dev/null 2>&1 || cargo install cargo-llvm-cov
	cd $(EXT_DIR) && npm ci
	cd webcompiler && npm ci
	cd website && npm ci
	@echo "==> Setup complete. Run 'make ci' to validate."

# ---------------------------------------------------------------------------
# Internal helpers — NOT public targets, NOT in .PHONY
# ---------------------------------------------------------------------------

# Build the pure-C runtime archives osprey-rs links at `--run` time. One shell
# so `cd` persists; faithful port of the original hardened C recipes.
_runtime:
	@echo "==> building C runtime archives ($(RTB)/lib*_runtime.a)"
	@cd compiler && set -e && $(MKDIR) bin lib && \
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

# --- rust (crates/) ---------------------------------------------------------
# Implements [TEST-RULES] — cargo test is fail-fast at the binary level by
# default (a failing test binary aborts the run); coverage via cargo-llvm-cov.
# `--profile ci` is the workspace's fast-compile profile (see root Cargo.toml).
_test_rust:
	@echo "==> [rust] running tests with coverage..."
	set -o pipefail && cargo llvm-cov --workspace --profile ci --lcov --output-path lcov.info 2>&1 | tee test.log

_coverage_check_rust:
	@if [ ! -f "$(COVERAGE_THRESHOLDS_FILE)" ]; then echo "FAIL: $(COVERAGE_THRESHOLDS_FILE) not found"; exit 1; fi; \
	THRESHOLD=$$(jq -r '.projects.crates.threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	LH=$$(grep '^LH:' lcov.info | awk -F: '{sum+=$$2} END{print sum+0}'); \
	LF=$$(grep '^LF:' lcov.info | awk -F: '{sum+=$$2} END{print sum+0}'); \
	if [ "$$LF" -eq 0 ]; then echo "[rust] FAIL: no lines in lcov.info"; exit 1; fi; \
	PCT=$$(awk "BEGIN{printf \"%.1f\", $$LH/$$LF*100}"); \
	PCT_INT=$$(awk "BEGIN{printf \"%d\", $$LH/$$LF*100}"); \
	echo "[rust] coverage: $${PCT}% (threshold: $${THRESHOLD}%)"; \
	if [ "$$PCT_INT" -lt "$$THRESHOLD" ]; then echo "[rust] FAIL: $${PCT}% < $${THRESHOLD}%"; exit 1; fi; \
	echo "[rust] OK: $${PCT}% >= $${THRESHOLD}%"

# Differential golden harness: every examples/tested/*.osp run through
# `osprey-rs --run` must match its .expectedoutput byte-for-byte, and the
# must-reject suite (examples/failscompilation) must stay within the
# FC_EXPECTED_ESCAPES ratchet declared in the harness.
_test_differential:
	@echo "==> [differential] osprey-rs --run vs .expectedoutput..."
	@out=$$(zsh crates/diff_examples.sh); echo "$$out"; \
	  echo "$$out" | grep -Eq 'FAIL=0 '  || { echo 'FAIL: differential mismatch'; exit 1; }; \
	  echo "$$out" | grep -Eq 'NOEXP=0 ' || { echo 'FAIL: example missing .expectedoutput'; exit 1; }; \
	  echo "$$out" | grep -q  'FC_OK'    || { echo 'FAIL: must-reject ratchet exceeded'; exit 1; }

# --- vscode-extension -------------------------------------------------------
# The extension's LSP server spawns the `osprey` binary at runtime, so the
# integration tests need a real compiler on PATH: the Rust binary is staged as
# `osprey`. vscode-test runs with V8 coverage; c8 merges the profiles into
# coverage/coverage-summary.json.
_test_vscode_extension:
	@echo "==> [vscode-extension] staging osprey-rs as 'osprey' for LSP integration..."
	$(MKDIR) target/path-bin
	cp $(BIN) target/path-bin/osprey
	@echo "==> [vscode-extension] running tests with V8 coverage..."
	$(RM) $(EXT_DIR)/coverage
	cd $(EXT_DIR) && set -o pipefail && \
	  PATH="$(CURDIR)/target/path-bin:$$PATH" \
	  npm run pretest 2>&1 | tee test.log && \
	  PATH="$(CURDIR)/target/path-bin:$$PATH" \
	  ./node_modules/.bin/vscode-test --coverage --coverage-output coverage \
	    --coverage-reporter text-summary --coverage-reporter json-summary --coverage-reporter html 2>&1 | tee -a test.log

_coverage_check_vscode_extension:
	@if [ ! -f "$(COVERAGE_THRESHOLDS_FILE)" ]; then echo "FAIL: $(COVERAGE_THRESHOLDS_FILE) not found"; exit 1; fi; \
	THRESHOLD=$$(jq -r '.projects["vscode-extension"].threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	if [ ! -f "$(EXT_DIR)/coverage/coverage-summary.json" ]; then \
	  echo "[vscode-extension] FAIL: coverage-summary.json not produced"; exit 1; \
	fi; \
	PCT=$$(jq -r '.total.lines.pct' "$(EXT_DIR)/coverage/coverage-summary.json"); \
	PCT_INT=$$(echo "$$PCT" | awk '{printf "%d", $$1}'); \
	echo "[vscode-extension] coverage: $${PCT}% (threshold: $${THRESHOLD}%)"; \
	if [ "$$PCT_INT" -lt "$$THRESHOLD" ]; then echo "[vscode-extension] FAIL: $${PCT}% < $${THRESHOLD}%"; exit 1; fi; \
	echo "[vscode-extension] OK: $${PCT}% >= $${THRESHOLD}%"

# =============================================================================
# Repo-Specific Targets
# =============================================================================

## tui: Build, then launch the interactive TUI demo (live GitHub API browser).
##      Runs in the current terminal so the raw-mode key reader gets real stdin.
tui: build
	@echo "==> launching TUI demo (live GitHub API browser)"
	./$(BIN) compiler/examples/tui/api_browser.osp --run

## run: Compile and run an Osprey file (usage: make run FILE=<path>)
run: build
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	./$(BIN) $(FILE) --run

## install: Install osprey-rs + runtime archives system-wide
install: build
	cargo install --path crates/osprey-cli --force
	sudo $(MKDIR) /usr/local/lib
	sudo cp $(RTB)/libfiber_runtime.a $(RTB)/libhttp_runtime.a /usr/local/lib/
	@echo "==> installed osprey-rs and runtime archives."

## uninstall: Remove osprey-rs + runtime archives from the system
uninstall:
	cargo uninstall osprey-cli 2>/dev/null || true
	sudo rm -f /usr/local/lib/libfiber_runtime.a /usr/local/lib/libhttp_runtime.a
	@echo "==> uninstalled."

## website-dev: Start local website development server
website-dev:
	cd website && npm run dev

## website-build: Build static site
website-build:
	cd website && npm run build

## rebuild-install-vsix: Uninstall → clean → rebuild → package → install the
##      VSCode extension into every VSCode profile, bundling the freshly-built
##      Rust compiler as `osprey`. macOS only. See [MAKE-IDE-EXT].
rebuild-install-vsix: build _vsix_uninstall _vsix_clean _vsix_build _vsix_bundle _vsix_package _vsix_install

# --- vsix sub-steps ---------------------------------------------------------
# Uninstall from default profile + every named profile in storage.json.
# `code --uninstall-extension` exits non-zero when not installed; swallowed so
# uninstall-before-install stays idempotent.
_vsix_uninstall:
	-@code --uninstall-extension $(EXT_ID) >/dev/null 2>&1 && echo "  [default] uninstalled" || echo "  [default] not installed"
	@jq -r '.userDataProfiles[]?.name' "$(VSCODE_STORAGE)" 2>/dev/null | while IFS= read -r prof; do \
	  [ -z "$$prof" ] && continue; \
	  code --profile "$$prof" --uninstall-extension $(EXT_ID) >/dev/null 2>&1 \
	    && echo "  [$$prof] uninstalled" || echo "  [$$prof] not installed"; \
	done

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
	DEST="$(EXT_DIR)/bin/$$OS-$$ARCH"; $(MKDIR) "$$DEST"; \
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
