# =============================================================================
# Standard Targets
# =============================================================================

## build: C runtime archives + Rust workspace (release) + VSCode extension
export CARGO_BUILD_JOBS ?= 4

build: _runtime $(EXT_NODE_DEPS)
	@echo "==> Building..."
	cargo build --release --workspace
	cd $(EXT_DIR) && npm run compile

## test: Fail-fast tests + coverage + per-project threshold enforcement.
##       Projects listed in coverage-thresholds.json are each tested + checked.
test: build
	@echo "==> Testing (fail-fast + coverage + per-project thresholds)..."
	$(MAKE) _test_runtime_incremental
	$(MAKE) _test_rust
	$(MAKE) _coverage_check_rust
	$(MAKE) _test_c_runtime
	$(MAKE) _coverage_check_c_runtime
# _test_language_corpus is deliberately NOT here. It runs every program in
# tests/ through the TAP subcommand, and _test_goldens then runs every one of
# them AGAIN — `run_test_corpus.sh` checks the in-language assertions AND the
# byte-exact golden in a single pass (see its header), so the golden run is
# strictly the stronger of the two observations. Executing the same assertions
# twice buys nothing and doubles the slowest stage in the build. The `osprey
# test` subcommand itself stays covered by bank-test and the CLI e2e suite.
	$(MAKE) _test_goldens
	$(MAKE) _conformance-gc
	$(MAKE) _conformance-arc
	$(MAKE) _test_profiler
	$(MAKE) _test_vscode_extension
	$(MAKE) _coverage_check_vscode_extension
	$(MAKE) _test_bench_tools

## bank: Rebuild and run the Talon Bank showcase for manual testing.
##       Opens http://127.0.0.1:18790 (dashboard) / /api/accounts (JSON API).
##       The hold marker keeps the server up; Ctrl-C removes it and exits.
bank: bank-web
	@echo "==> Talon Bank live on http://127.0.0.1:18790  (Ctrl-C to stop)"
	@touch /tmp/talon_bank.hold
	@trap 'rm -f /tmp/talon_bank.hold' EXIT INT TERM; \
	  (if command -v open >/dev/null && command -v curl >/dev/null; then \
	    attempts=0; \
	    until accounts="$$(curl -fsS http://127.0.0.1:18790/api/accounts 2>/dev/null)" && [[ "$$accounts" == *'"Priya Sharma"'* ]]; do \
	      attempts=$$((attempts + 1)); [ "$$attempts" -ge 200 ] && exit 0; sleep 0.1; \
	    done; \
	    open http://127.0.0.1:18790; \
	  fi) & \
	  ./$(BIN) examples/projects/modules --run

## bank-web: Regenerate the embedded React host + Osprey WebAssembly client.
##           Requires Node and a WASI sysroot; the generated Osprey Bundle is
##           committed so ordinary native/CI builds do not need either tool.
bank-web: build _runtime_wasm $(BANK_WEB_NODE_DEPS)
	@echo "==> Building Talon Bank browser application..."
	cd examples/projects/modules/web && npm run build

## bank-test: Native Osprey unit tests for the Talon Bank pure domain layer,
##            run through the built-in `osprey test` harness (TAP output).
bank-test: build
	$(MAKE) _bank_test

# The same suite against an ALREADY-BUILT binary. No prerequisites, so the
# staged pipeline can run it on the downloaded `osprey-build` artifact. Going
# through `bank-test` there rebuilt the whole workspace and then ran `tsc -b`
# in vscode-extension, whose node_modules that job never installs — three
# TS2688 errors and a dead stage, for work the build job had already done.
_bank_test:
	@echo "==> Bank native tests (osprey test)..."
	./$(BIN) test examples/projects/modules/test

# _language-test: Run the assertion-driven Default + ML core language corpus.
_language-test: build
	$(MAKE) _test_language_corpus

## bank-e2e: Browser end-to-end tests for the Talon Bank modules showcase
##           (examples/projects/modules) — real Chromium via Playwright drives
##           the compiled osprey binary serving its HTTP API and web UI.
##           The server's hold marker must not outlive the run: a leaked one
##           makes every later bank-*-test refuse to start.
bank-e2e: bank-web $(BANK_E2E_NODE_DEPS)
	@echo "==> Bank e2e (Playwright)..."
	cd examples/projects/modules/e2e && npx playwright install chromium && npx playwright test
	@test ! -e /tmp/talon_bank.hold || { echo "FAIL: bank e2e left /tmp/talon_bank.hold behind; its webServer was stopped without running its trap"; exit 1; }

## docs-html-test: Browser acceptance for the HTML documentation export. Real
##                 Chromium over file:// — the strictest case, and the one a
##                 reader hits first — covering every theme, custom stylesheet
##                 ordering, user guides, search with no server behind it, the
##                 mobile disclosure and keyboard navigation. Reuses the bank
##                 e2e project's Playwright install rather than adding another.
docs-html-test: build $(BANK_E2E_NODE_DEPS)
	@echo "==> Documentation HTML browser acceptance..."
	cd examples/projects/modules/e2e && npx playwright install chromium
	node scripts/verify-docs-html.mjs $(BIN)

## lint: Run all linters/analyzers (read-only). Checks formatting but does
## NOT rewrite it — `make fmt` does that. The fmt check lives HERE because
## `lint` is what the required CI job runs; a format gate only in an optional
## job cannot block a merge.
lint: _deslop _lint

_lint: $(EXT_NODE_DEPS)
	@echo "==> Linting..."
	node scripts/verify-node-deps-guard.mjs
	node scripts/verify-no-dead-code.mjs
	node scripts/verify-release-gates.mjs
	python3 scripts/test-release-gate.py
	cargo fmt --all --check
	cargo clippy --workspace --all-targets -- -D warnings
	cd $(EXT_DIR) && npm run lint

# _deslop: Code-duplication gate. Fails the build when measured
# duplication exceeds the ceiling in .deslop.toml (exit 3). Exclusions and the
# threshold live in that committed config — the single source of truth.
# The installer resolves the latest release for both local and hosted CI.

_deslop:
	@echo "==> Duplication gate (deslop)..."
	bash scripts/install-deslop.sh
	deslop . --nohtml --nojson --output $(CURDIR)/target/deslop-report --log-to-console --log-level error --no-color

## hawk: Dead-code gate (astral-sh/hawk). Fails the build when any `pub`
## declaration is unreachable from the osprey binary (hawk::dead_public).
##
## hawk counts the workspace's TEST binaries as reachability roots, so an item
## whose only callers are its own `#[cfg(test)]` module passes this gate — that
## is how `osprey_debug::DebugBuild` and `osprey_syntax::dependency_sets` lived
## in the tree. `scripts/verify-no-dead-code.mjs` (run by `_lint`) answers the
## narrower question "does PRODUCT code name this?" and catches that class.
## Neither gate subsumes the other; both run. Scoped
## to dead_public ONLY — unnecessary_public / restricted-visibility findings are
## over-exposure, not dead code, and several are irreducibly public for the
## integration tests under this workspace's `dead_code = "deny"` policy, so they
## must NOT fail the gate. hawk needs rustc_private (RUSTC_BOOTSTRAP=1) and its
## prebuilt driver is built against one exact rustc and refuses any other, so
## the gate runs hawk under that toolchain, not under whatever stable is today:
## CI installing the latest hawk over the latest stable failed the day stable
## moved to 1.99.1 while hawk 0.1.15 was built for 1.99.0. HAWK_VERSION and
## HAWK_RUSTC move together; `make hawk-install` installs both. An absent
## cargo-hawk FAILS this
## target: a gate that cannot run must not report success.
## Install: make hawk-install (https://github.com/astral-sh/hawk)
HAWK_VERSION ?= 0.1.15
HAWK_RUSTC   ?= 1.99.0
hawk-install:
	rustup toolchain install $(HAWK_RUSTC) --profile minimal
	curl --proto '=https' --tlsv1.2 -LsSf https://github.com/astral-sh/hawk/releases/download/$(HAWK_VERSION)/cargo-hawk-installer.sh | sh

hawk:
	@echo "==> Dead-code gate (hawk)..."
	@if ! command -v cargo-hawk >/dev/null 2>&1; then \
		echo "FAIL: cargo-hawk is not installed, so the dead-code gate cannot run."; \
		echo "      A gate that cannot run must not report success — this used to"; \
		echo "      print a warning and exit 0, so 'make hawk' was green on every"; \
		echo "      machine without the binary, the CI runner included if its"; \
		echo "      install step ever stopped landing cargo-hawk on PATH."; \
		echo "      Install: make hawk-install"; \
		exit 1; \
	fi
	RUSTC_BOOTSTRAP=1 cargo +$(HAWK_RUSTC) hawk check --only dead-public -D hawk::dead_public --target-dir $(CURDIR)/target/hawk

## fmt: Format all code in-place. Pass CHECK=1 for read-only check (CI use).
fmt: $(EXT_NODE_DEPS)
	@echo "==> Formatting$(if $(CHECK), (check mode),)..."
	cargo fmt --all$(if $(CHECK), --check,)
	cd $(EXT_DIR) && npx prettier$(if $(CHECK), --check, --write) .

## clean: Remove all build artifacts
clean:
	@echo "==> Cleaning..."
	cargo clean
	$(RM) $(RTB) compiler/lib outputs lcov.info test.log
	cd $(EXT_DIR) && $(RM) out dist coverage test.log

## ci: lint + hawk + test + bank-test + bank-e2e + build (full CI simulation)
ci: lint hawk test bank-test bank-e2e docs-html-test mobile-domain-test bank-mobile-domain-test build

## wasm: Build everything for the WebAssembly target, ready to go — the wasm
## runtime archive (compiler/bin/libosprey_runtime_wasm.a), the hello example,
## and Osprey Data Studio in BOTH flavors (studio.{osp,ospml} -> one byte-
## identical manifest that drives the SQLite dashboard in examples/wasm/
## index.html) — then validate them and smoke-run under Node's WASI, the browser
## WASI shim, with committed expected output for hello and both Studio flavors.
## Requires clang (wasm32 backend),
## wasm-ld and a WASI sysroot —
## `brew install lld wasi-libc` (macOS) or the wasi-sdk. See
## docs/specs/compiler/0022-WebAssemblyTarget.md.
wasm: build _runtime_wasm
	@echo "==> compiling the wasm example -> examples/wasm/build/"
	@$(MKDIR) examples/wasm/build
	$(BIN) examples/wasm/hello.osp --target=wasm32 --compile -o examples/wasm/build/hello.wasm
	@echo "==> validating + smoke-running examples/wasm/build/hello.wasm"
	$(call wasm_validate,examples/wasm/build/hello.wasm)
	node scripts/wasm-smoke.mjs         examples/wasm/build/hello.wasm examples/wasm/hello.expectedoutput
	node scripts/wasm-browser-smoke.mjs examples/wasm/build/hello.wasm examples/wasm/hello.expectedoutput
	@echo "==> compiling Osprey Data Studio (BOTH flavors) -> examples/wasm/build/"
	$(BIN) examples/wasm/studio.osp   --target=wasm32 --compile -o examples/wasm/build/studio.osp.wasm
	$(BIN) examples/wasm/studio.ospml --target=wasm32 --compile -o examples/wasm/build/studio.ospml.wasm
	$(call wasm_validate,examples/wasm/build/studio.osp.wasm examples/wasm/build/studio.ospml.wasm)
	@echo "==> both Studio flavors must emit the SAME manifest (byte-identical golden)"
	node scripts/wasm-smoke.mjs         examples/wasm/build/studio.osp.wasm   examples/wasm/studio.expectedoutput
	node scripts/wasm-browser-smoke.mjs examples/wasm/build/studio.osp.wasm   examples/wasm/studio.expectedoutput
	node scripts/wasm-smoke.mjs         examples/wasm/build/studio.ospml.wasm examples/wasm/studio.expectedoutput
	node scripts/wasm-browser-smoke.mjs examples/wasm/build/studio.ospml.wasm examples/wasm/studio.expectedoutput
	$(MAKE) _test_wasm_goldens
	@echo "==> wasm ready: built + validated + WASI/browser smoke green"

wasm wasm-site _runtime_wasm bank-web _test_wasm_goldens: export PATH := $(WASM_PATH_PREFIX)$(PATH)

## wasm-site: Build only the WebAssembly artifacts published by the website.
##      Used by GitHub Pages before `npm run build`; does not rely on checked-in
##      wasm binaries. Requires clang, wasm-ld, a WASI sysroot, and node.
wasm-site: _runtime_wasm
	@echo "==> building osprey compiler for the website wasm demo"
	cargo build --release -p osprey-cli
	@echo "==> compiling Osprey Data Studio website assets -> examples/wasm/build/"
	@$(MKDIR) examples/wasm/build
	$(BIN) examples/wasm/studio.osp   --target=wasm32 --compile -o examples/wasm/build/studio.osp.wasm
	$(BIN) examples/wasm/studio.ospml --target=wasm32 --compile -o examples/wasm/build/studio.ospml.wasm
	$(call wasm_validate,examples/wasm/build/studio.osp.wasm examples/wasm/build/studio.ospml.wasm)
	node scripts/wasm-smoke.mjs         examples/wasm/build/studio.osp.wasm   examples/wasm/studio.expectedoutput
	node scripts/wasm-browser-smoke.mjs examples/wasm/build/studio.osp.wasm   examples/wasm/studio.expectedoutput
	node scripts/wasm-smoke.mjs         examples/wasm/build/studio.ospml.wasm examples/wasm/studio.expectedoutput
	node scripts/wasm-browser-smoke.mjs examples/wasm/build/studio.ospml.wasm examples/wasm/studio.expectedoutput
	@echo "==> website wasm demo ready"

## wasm-serve: Build the wasm target (full `make wasm`), then static-host
##      $(WASM_SERVE_DIR) at http://localhost:$(WASM_SERVE_PORT)/ and open it in
##      your browser. Long-running dev server — Ctrl-C to stop. Override the port
##      with WASM_SERVE_PORT=<n>. (`make wasm` itself stays headless for CI.)
wasm-serve: wasm
	@URL="http://localhost:$(WASM_SERVE_PORT)/"; \
	  command -v python3 >/dev/null 2>&1 || { echo "FAIL: python3 not found (needed for the dev server)"; exit 1; }; \
	  echo "==> serving $(WASM_SERVE_DIR)/ at $$URL — opening browser (Ctrl-C to stop)"; \
	  OPENER=$$(command -v open || command -v xdg-open || true); \
	  if [ -n "$$OPENER" ]; then ( sleep 1; "$$OPENER" "$$URL" >/dev/null 2>&1 || true ) & \
	  else echo "  (no 'open'/'xdg-open' found — browse to $$URL manually)"; fi; \
	  cd $(WASM_SERVE_DIR) && exec python3 -m http.server $(WASM_SERVE_PORT)

## setup: Post-create dev environment setup (used by devcontainer)
setup: $(EXT_NODE_DEPS) $(WEBCOMPILER_NODE_DEPS) $(WEBSITE_NODE_DEPS)
	@echo "==> Setting up development environment..."
	rustup component add rustfmt clippy llvm-tools-preview
	command -v cargo-llvm-cov >/dev/null 2>&1 || cargo install cargo-llvm-cov
	bash scripts/install-deslop.sh
	@echo "==> Setup complete. Run 'make ci' to validate."
