# agent-pmo:74cf183
# =============================================================================
# Standard Makefile — osprey
# Cross-platform: Linux, macOS, Windows (via GNU Make)
# Primary language: Go (compiler/), with TypeScript sub-projects
# =============================================================================

.PHONY: build test lint fmt clean ci setup

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
  RM = rm -rf
  MKDIR = mkdir -p
endif

# ---------------------------------------------------------------------------
# Coverage — single source of truth is coverage-thresholds.json
# See REPO-STANDARDS-SPEC [COVERAGE-THRESHOLDS-JSON].
# ---------------------------------------------------------------------------
COVERAGE_THRESHOLDS_FILE := coverage-thresholds.json

# =============================================================================
# Standard Targets
# =============================================================================

## build: Compile all artifacts (Go compiler + C runtimes + TypeScript extension)
build:
	@echo "==> Building..."
	cd compiler && $(MAKE) build
	cd vscode-extension && npm run compile

## test: Fail-fast tests + coverage + threshold enforcement.
##       See REPO-STANDARDS-SPEC [TEST-RULES] and [COVERAGE-THRESHOLDS-JSON].
test:
	@echo "==> Testing (fail-fast + coverage + threshold)..."
	cd compiler && $(MAKE) _test
	$(MAKE) _coverage_check

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

# Implements [TEST-RULES] — fail-fast Go tests with coverage
_test:
	cd compiler && go test -failfast -covermode=atomic -coverprofile=coverage.out ./... -p 1
	cd compiler && go tool cover -func=coverage.out

# Implements [COVERAGE-THRESHOLDS-JSON] — reads coverage-thresholds.json, fails below threshold
_coverage_check:
	@if [ ! -f "$(COVERAGE_THRESHOLDS_FILE)" ]; then echo "FAIL: $(COVERAGE_THRESHOLDS_FILE) not found"; exit 1; fi; \
	THRESHOLD=$$(jq -r '.default_threshold' "$(COVERAGE_THRESHOLDS_FILE)"); \
	PCT=$$(cd compiler && go tool cover -func=coverage.out | awk '/^total:/{print $$3}' | tr -d '%'); \
	PCT_INT=$$(echo "$$PCT" | awk '{printf "%d", $$1}'); \
	echo "Line coverage: $${PCT}% (threshold: $${THRESHOLD}%)"; \
	if [ "$$PCT_INT" -lt "$${THRESHOLD}" ]; then \
	  echo "FAIL: $${PCT}% < $${THRESHOLD}%"; exit 1; \
	else \
	  echo "OK: $${PCT}% >= $${THRESHOLD}%"; \
	fi

# =============================================================================
# Repo-Specific Targets
# =============================================================================

.PHONY: install uninstall regenerate-parser run website-dev website-build

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
