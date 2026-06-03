# agent-pmo:74cf183
# =============================================================================
# Standard Makefile — osprey
# Cross-platform: Linux, macOS, Windows (via GNU Make)
# Primary language: Go (compiler/), with TypeScript sub-projects
# =============================================================================

.PHONY: build tui test lint fmt clean ci setup ratchet

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
# Coverage — single source of truth is coverage-thresholds.json
# See REPO-STANDARDS-SPEC [COVERAGE-THRESHOLDS-JSON].
# ---------------------------------------------------------------------------
COVERAGE_THRESHOLDS_FILE := coverage-thresholds.json

# ---------------------------------------------------------------------------
# VSIX (VSCode extension) — macOS only.
# Profile auto-discovery reads names from globalStorage/storage.json so
# `make vsix` installs into the default profile AND every named profile.
# ---------------------------------------------------------------------------
EXT_DIR         := vscode-extension
EXT_ID          := nimblesite.osprey
VSCODE_USER_DIR := $(HOME)/Library/Application Support/Code/User
VSCODE_STORAGE  := $(VSCODE_USER_DIR)/globalStorage/storage.json

# =============================================================================
# Standard Targets
# =============================================================================

## build: Compile all artifacts (Go compiler + C runtimes + TypeScript extension)
build:
	@echo "==> Building..."
	cd compiler && $(MAKE) build
	cd vscode-extension && npm run compile

## tui: Rebuild the compiler and launch the interactive TUI demo (live GitHub
##      API browser). Runs in the current terminal so the raw-mode key reader
##      gets a real stdin.
tui:
	@echo "==> Building compiler + launching TUI..."
	cd compiler && $(MAKE) tui

## test: Fail-fast tests + coverage + per-project threshold enforcement.
##       See REPO-STANDARDS-SPEC [TEST-RULES] and [COVERAGE-THRESHOLDS-JSON].
##       Each project listed in coverage-thresholds.json is tested separately.
test:
	@echo "==> Testing (fail-fast + coverage + per-project thresholds)..."
	$(MAKE) _test_compiler
	$(MAKE) _coverage_check_compiler
	$(MAKE) _test_vscode_extension
	$(MAKE) _coverage_check_vscode_extension

## ratchet: Update each project's coverage threshold in coverage-thresholds.json
##          to (measured - 1) so the next run requires at least the current level.
##          Run after improving coverage; commit the resulting JSON change.
ratchet:
	@echo "==> Ratcheting thresholds to (measured - 1)..."
	$(MAKE) _ratchet_compiler
	$(MAKE) _ratchet_vscode_extension
	@echo "==> Updated $(COVERAGE_THRESHOLDS_FILE). Review and commit."

## lint: Run all linters/analyzers (read-only). Does NOT format.
lint:
	@echo "==> Linting..."
	cd compiler && golangci-lint run --config .golangci.yml
	cd vscode-extension && npm run lint

## fmt: Format all code in-place. Pass CHECK=1 for read-only check (CI use).
fmt:
	@echo "==> Formatting$(if $(CHECK), (check mode),)..."
	gofmt$(if $(CHECK), -l compiler/... | grep . && exit 1 || true, -w compiler/...)
	cd vscode-extension && npx prettier$(if $(CHECK), --check, --write) .

## clean: Remove all build artifacts
clean:
	@echo "==> Cleaning..."
	cd compiler && $(MAKE) clean
	cd vscode-extension && $(RM) out dist

## ci: lint + test + build (full CI simulation)
ci: lint test build

## setup: Post-create dev environment setup (used by devcontainer)
setup:
	@echo "==> Setting up development environment..."
	cd compiler && go mod download
	cd compiler && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.1.6
	cd vscode-extension && npm ci
	cd webcompiler && npm ci
	cd website && npm ci
	@echo "==> Setup complete. Run 'make ci' to validate."

# ---------------------------------------------------------------------------
# Internal helpers — NOT public targets, NOT in .PHONY
# ---------------------------------------------------------------------------

# --- compiler -------------------------------------------------------------
# Implements [TEST-RULES] — fail-fast Go tests with coverage.
# -coverpkg=./... instruments ALL packages so integration tests
# (which call codegen via Go API) contribute to coverage.
# Generated code (ANTLR parser) is excluded from coverage.out before threshold check.
_test_compiler:
	@echo "==> [compiler] running tests..."
	cd compiler && set -o pipefail && go test -failfast -covermode=atomic -coverpkg=./... -coverprofile=coverage.out.raw ./... -p 1 2>&1 | tee test.log
	cd compiler && grep -v '/parser/osprey_' coverage.out.raw > coverage.out
	cd compiler && go tool cover -func=coverage.out | tail -1

_coverage_check_compiler:
	@if [ ! -f "$(COVERAGE_THRESHOLDS_FILE)" ]; then echo "FAIL: $(COVERAGE_THRESHOLDS_FILE) not found"; exit 1; fi; \
	THRESHOLD=$$(jq -r '.projects.compiler.threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	PCT=$$(cd compiler && go tool cover -func=coverage.out | awk '/^total:/{print $$3}' | tr -d '%'); \
	PCT_INT=$$(echo "$$PCT" | awk '{printf "%d", $$1}'); \
	echo "[compiler] coverage: $${PCT}% (threshold: $${THRESHOLD}%)"; \
	if [ "$$PCT_INT" -lt "$${THRESHOLD}" ]; then \
	  echo "[compiler] FAIL: $${PCT}% < $${THRESHOLD}%"; exit 1; \
	fi; \
	echo "[compiler] OK: $${PCT}% >= $${THRESHOLD}%"; \
	NEW_THRESHOLD=$$((PCT_INT - 1)); \
	if [ "$$NEW_THRESHOLD" -lt 0 ]; then NEW_THRESHOLD=0; fi; \
	if [ "$$NEW_THRESHOLD" -gt "$$THRESHOLD" ]; then \
	  jq ".projects.compiler.threshold = $$NEW_THRESHOLD" "$(COVERAGE_THRESHOLDS_FILE)" > "$(COVERAGE_THRESHOLDS_FILE).tmp" && mv "$(COVERAGE_THRESHOLDS_FILE).tmp" "$(COVERAGE_THRESHOLDS_FILE)"; \
	  echo "[compiler] auto-ratchet: threshold $${THRESHOLD} -> $${NEW_THRESHOLD} (measured $${PCT}%)"; \
	fi

_ratchet_compiler:
	@if [ ! -f "compiler/coverage.out" ]; then echo "Run 'make test' first to produce coverage.out"; exit 1; fi; \
	PCT=$$(cd compiler && go tool cover -func=coverage.out | awk '/^total:/{print $$3}' | tr -d '%'); \
	PCT_INT=$$(echo "$$PCT" | awk '{printf "%d", $$1}'); \
	NEW_THRESHOLD=$$((PCT_INT - 1)); \
	if [ "$$NEW_THRESHOLD" -lt 0 ]; then NEW_THRESHOLD=0; fi; \
	OLD=$$(jq -r '.projects.compiler.threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	if [ "$$NEW_THRESHOLD" -le "$$OLD" ]; then \
	  echo "[compiler] threshold unchanged: $${OLD} (measured $${PCT}%, ratchet would set $${NEW_THRESHOLD})"; \
	else \
	  jq ".projects.compiler.threshold = $$NEW_THRESHOLD" "$(COVERAGE_THRESHOLDS_FILE)" > "$(COVERAGE_THRESHOLDS_FILE).tmp" && mv "$(COVERAGE_THRESHOLDS_FILE).tmp" "$(COVERAGE_THRESHOLDS_FILE)"; \
	  echo "[compiler] threshold $${OLD} -> $${NEW_THRESHOLD} (measured $${PCT}%)"; \
	fi

# --- vscode-extension -----------------------------------------------------
# The extension's LSP server spawns the `osprey` binary at runtime, so the
# integration tests need the real compiler on PATH. Build it first, then run
# vscode-test with PATH augmented to include compiler/bin and
# NODE_V8_COVERAGE set so the Electron Extension Host writes V8 coverage
# profiles for the extension code (client + server). After the run, c8
# merges those profiles into coverage/coverage-summary.json.
_test_vscode_extension:
	@echo "==> [vscode-extension] building compiler for LSP integration..."
	cd compiler && $(MAKE) build
	@echo "==> [vscode-extension] running tests with real compiler + V8 coverage..."
	rm -rf vscode-extension/coverage
	cd vscode-extension && set -o pipefail && \
	  PATH="$(CURDIR)/compiler/bin:$$PATH" \
	  npm run pretest 2>&1 | tee test.log && \
	  PATH="$(CURDIR)/compiler/bin:$$PATH" \
	  ./node_modules/.bin/vscode-test --coverage --coverage-output coverage \
	    --coverage-reporter text-summary --coverage-reporter json-summary --coverage-reporter html 2>&1 | tee -a test.log

_coverage_check_vscode_extension:
	@if [ ! -f "$(COVERAGE_THRESHOLDS_FILE)" ]; then echo "FAIL: $(COVERAGE_THRESHOLDS_FILE) not found"; exit 1; fi; \
	THRESHOLD=$$(jq -r '.projects["vscode-extension"].threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	if [ ! -f "vscode-extension/coverage/coverage-summary.json" ]; then \
	  echo "[vscode-extension] FAIL: coverage-summary.json not produced — c8 report failed"; exit 1; \
	fi; \
	PCT=$$(jq -r '.total.lines.pct' "vscode-extension/coverage/coverage-summary.json"); \
	PCT_INT=$$(echo "$$PCT" | awk '{printf "%d", $$1}'); \
	echo "[vscode-extension] coverage: $${PCT}% (threshold: $${THRESHOLD}%)"; \
	if [ "$$PCT_INT" -lt "$${THRESHOLD}" ]; then \
	  echo "[vscode-extension] FAIL: $${PCT}% < $${THRESHOLD}%"; exit 1; \
	fi; \
	echo "[vscode-extension] OK: $${PCT}% >= $${THRESHOLD}%"; \
	NEW_THRESHOLD=$$((PCT_INT - 1)); \
	if [ "$$NEW_THRESHOLD" -lt 0 ]; then NEW_THRESHOLD=0; fi; \
	if [ "$$NEW_THRESHOLD" -gt "$$THRESHOLD" ]; then \
	  jq ".projects[\"vscode-extension\"].threshold = $$NEW_THRESHOLD" "$(COVERAGE_THRESHOLDS_FILE)" > "$(COVERAGE_THRESHOLDS_FILE).tmp" && mv "$(COVERAGE_THRESHOLDS_FILE).tmp" "$(COVERAGE_THRESHOLDS_FILE)"; \
	  echo "[vscode-extension] auto-ratchet: threshold $${THRESHOLD} -> $${NEW_THRESHOLD} (measured $${PCT}%)"; \
	fi

_ratchet_vscode_extension:
	@if [ ! -f "vscode-extension/coverage/coverage-summary.json" ]; then \
	  echo "[vscode-extension] no coverage report — skipping ratchet"; \
	else \
	  PCT=$$(jq -r '.total.lines.pct' "vscode-extension/coverage/coverage-summary.json"); \
	  PCT_INT=$$(echo "$$PCT" | awk '{printf "%d", $$1}'); \
	  NEW_THRESHOLD=$$((PCT_INT - 1)); \
	  if [ "$$NEW_THRESHOLD" -lt 0 ]; then NEW_THRESHOLD=0; fi; \
	  OLD=$$(jq -r '.projects["vscode-extension"].threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	  if [ "$$NEW_THRESHOLD" -le "$$OLD" ]; then \
	    echo "[vscode-extension] threshold unchanged: $${OLD} (measured $${PCT}%)"; \
	  else \
	    jq ".projects[\"vscode-extension\"].threshold = $$NEW_THRESHOLD" "$(COVERAGE_THRESHOLDS_FILE)" > "$(COVERAGE_THRESHOLDS_FILE).tmp" && mv "$(COVERAGE_THRESHOLDS_FILE).tmp" "$(COVERAGE_THRESHOLDS_FILE)"; \
	    echo "[vscode-extension] threshold $${OLD} -> $${NEW_THRESHOLD} (measured $${PCT}%)"; \
	  fi; \
	fi

# =============================================================================
# Repo-Specific Targets
# =============================================================================

.PHONY: install uninstall regenerate-parser run website-dev website-build vsix

## install: Install osprey compiler globally
install:
	cd compiler && $(MAKE) install

## uninstall: Remove osprey compiler from system
uninstall:
	cd compiler && $(MAKE) uninstall

## regenerate-parser: Regenerate ANTLR parser from grammar
regenerate-parser:
	cd compiler && $(MAKE) regenerate-parser

## run: Run compiler on a specific file (usage: make run FILE=<path>)
run:
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	cd compiler && go run cmd/osprey/main.go $(FILE)

## website-dev: Start local website development server
website-dev:
	cd website && npm run dev

## website-build: Build static site
website-build:
	cd website && npm run build

## vsix: Clean → uninstall → build → package → install (every VSCode profile)
##       Compiler is rebuilt cleanly first so the LSP-bundled binary is in sync.
##       Discovers profiles from $(VSCODE_STORAGE) and uses `code --profile`
##       to install into each one. macOS only.
vsix:
	@echo "==> [vsix] clean"
	$(MAKE) _vsix_clean
	@echo "==> [vsix] compiler build (fresh binary + libs)"
	cd compiler && $(MAKE) build
	@echo "==> [vsix] uninstall (all profiles)"
	$(MAKE) _vsix_uninstall
	@echo "==> [vsix] extension build"
	$(MAKE) _vsix_build
	@echo "==> [vsix] bundle freshly-built compiler"
	$(MAKE) _vsix_bundle
	@echo "==> [vsix] package"
	$(MAKE) _vsix_package
	@echo "==> [vsix] install (all profiles)"
	$(MAKE) _vsix_install

# --- vsix sub-targets -----------------------------------------------------
_vsix_clean:
	cd $(EXT_DIR) && $(RM) out dist *.vsix

# Uninstall from default profile + every named profile in storage.json.
# `code --uninstall-extension` exits non-zero when the extension isn't
# present; we swallow that since uninstall-before-install must be idempotent.
_vsix_uninstall:
	-@code --uninstall-extension $(EXT_ID) >/dev/null 2>&1 && echo "  [default] uninstalled" || echo "  [default] not installed"
	@jq -r '.userDataProfiles[]?.name' "$(VSCODE_STORAGE)" 2>/dev/null | while IFS= read -r prof; do \
	  [ -z "$$prof" ] && continue; \
	  if code --profile "$$prof" --uninstall-extension $(EXT_ID) >/dev/null 2>&1; then \
	    echo "  [$$prof] uninstalled"; \
	  else \
	    echo "  [$$prof] not installed"; \
	  fi; \
	done

_vsix_build:
	cd $(EXT_DIR) && npm run compile

# Stage the freshly-built compiler where the extension's resolveBundledCompiler
# looks for it (bin/<os>-<arch>/osprey), so the packaged VSIX runs LSP
# diagnostics against THIS compiler instead of a stale global `osprey`. The
# platform string mirrors shipwrightPlatform() in client/src/extension.ts.
# bin/ is not in .vscodeignore, so vsce includes it in the package.
_vsix_bundle:
	@OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	case "$$OS" in darwin) OS=darwin;; linux) OS=linux;; *) OS=win32;; esac; \
	ARCH=$$(uname -m); case "$$ARCH" in arm64|aarch64) ARCH=arm64;; *) ARCH=x64;; esac; \
	DEST="$(EXT_DIR)/bin/$$OS-$$ARCH"; \
	mkdir -p "$$DEST"; \
	cp compiler/bin/osprey "$$DEST/osprey"; \
	echo "  bundled compiler -> $$DEST/osprey ($$($$DEST/osprey --version 2>/dev/null | head -1))"

_vsix_package:
	cd $(EXT_DIR) && npm run package

# Install latest *.vsix to default profile + every named profile.
# --force lets `code` reinstall over an existing version without prompting.
_vsix_install:
	@VSIX=$$(ls -t $(EXT_DIR)/*.vsix 2>/dev/null | head -1); \
	if [ -z "$$VSIX" ]; then echo "FAIL: no .vsix in $(EXT_DIR)/ — did _vsix_package run?"; exit 1; fi; \
	echo "  vsix: $$VSIX"; \
	code --install-extension "$$VSIX" --force && echo "  [default] installed"; \
	jq -r '.userDataProfiles[]?.name' "$(VSCODE_STORAGE)" 2>/dev/null | while IFS= read -r prof; do \
	  [ -z "$$prof" ] && continue; \
	  code --profile "$$prof" --install-extension "$$VSIX" --force && echo "  [$$prof] installed"; \
	done
