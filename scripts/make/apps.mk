## graphics: Build the platform graphics host and run the animated GPU demo
#            Osprey drives the scene; the GPU shades every pixel. macOS renders
#            it through Metal in a Swift host linked to the Mac slice, Windows
#            through a Direct3D 12 DLL linked to the native executable — and the
#            Osprey sources are the same files on both, which is the entire point.
#            SCENE=kali|opal|character selects one of the three scenes; every
#            one is a one-file project sharing examples/graphics/base/base.osp
#            and one named fragment entry in the platform shader library.
#            GFX_WARN overrides the C warning set. See $(GFX_DIR)/README.md.
GFX_DIR := examples/graphics
GFX_WARN ?= -Wall -Werror
SCENE ?= kali

ifeq ($(OS),Windows_NT)
# Direct3D 12. Built with the same MinGW UCRT64 toolchain that builds the C
# runtime on Windows, so `// @link(windows): ospgfx` in base/base.osp resolves
# the import library below.
#
# UNVERIFIED. This branch, the two C files it compiles and base.hlsl were all
# written on a macOS machine with no Windows toolchain and have never been
# built or run. $(GFX_DIR)/README.md says what would establish that they work.
GFX_SHADER := $(GFX_DIR)/base.hlsl
GFX_SRC    := $(GFX_DIR)/ospgfx_d3d12.c $(GFX_DIR)/ospgfx_d3d12_setup.c
GFX_DLL    := $(GFX_DIR)/ospgfx.dll
GFX_LIB    := $(GFX_DIR)/libospgfx.dll.a
GFX_SYSLIB := -ld3d12 -ldxgi -ld3dcompiler -ldxguid -luser32
# Every entry point the run-time compiler will be asked for, with its target
# profile. `make _graphics-shader` compiles all four so a typo in one scene's
# fragment cannot hide behind another scene running fine.
GFX_ENTRIES := osp_vertex:vs_5_0 osp_fragment:ps_5_0 \
	osp_fragment_opal:ps_5_0 osp_fragment_character:ps_5_0

ifeq ($(MSYSTEM),)
# Native PowerShell has neither bash for the recipes below nor a MinGW cc.
graphics _graphics-shader:
	@Write-Host "make graphics needs the MSYS2 UCRT64 shell -- see $(GFX_DIR)/README.md"
else
$(GFX_LIB): $(GFX_SRC) $(GFX_DIR)/ospgfx_d3d12.h $(GFX_DIR)/ospgfx_abi.h
	@echo "==> building Osprey -> Direct3D 12 graphics bridge"
	$(CC) -shared -O2 $(GFX_WARN) -o $(GFX_DLL) $(GFX_SRC) \
		-Wl,--out-implib,$(GFX_LIB) $(GFX_SYSLIB)

# The bridge compiles the shader at run time, so an error in it would otherwise
# only surface as a window that refuses to open. fxc is the same compiler
# D3DCompile is, so checking here checks exactly what the bridge will do; the
# Windows SDK ships it, but not on PATH, so the check skips when it is absent.
_graphics-shader:
	@if command -v fxc.exe >/dev/null 2>&1; then \
		for entry in $(GFX_ENTRIES); do \
			echo "==> checking $(GFX_SHADER) $${entry%%:*}"; \
			fxc.exe -nologo -T $${entry##*:} -E $${entry%%:*} \
				-Fo $(GFX_DIR)/.shadercheck.cso $(GFX_SHADER) || exit 1; \
		done; \
		$(RM) $(GFX_DIR)/.shadercheck.cso; \
	else \
		echo "==> skipping shader check (no fxc.exe on PATH)"; \
	fi

# The DLL lives beside its source rather than beside the compiled scene, so the
# loader is pointed at it for the length of the run.
graphics: build $(GFX_LIB) _graphics-shader
	@echo "==> launching Osprey graphics demo: $(SCENE) (close the window to exit)"
	PATH="$(CURDIR)/$(GFX_DIR):$$PATH" ./$(BIN) $(GFX_DIR)/$(SCENE) --run
endif

else
GFX_SHADER := $(GFX_DIR)/base.metal
GFX_SCENES := kali opal character
GFX_SMOKE_FRAMES ?= 60

# The bridge compiles the shader at run time, so a syntax error in it would
# otherwise only surface as a window that refuses to open. Check it up front
# when the Metal toolchain is installed, and skip quietly when it is not. Xcode
# 26 and later ship `metal` as a stub until the toolchain component downloads,
# so availability is asked of the compiler itself rather than of `xcrun --find`.
_graphics-shader:
	@if xcrun -sdk macosx metal --version >/dev/null 2>&1; then \
		echo "==> checking $(GFX_SHADER)"; \
		xcrun -sdk macosx metal -c $(GFX_SHADER) -o /dev/null; \
	else \
		echo "==> skipping shader check (no Metal toolchain; xcodebuild -downloadComponent MetalToolchain installs it)"; \
	fi

# The scene compiles for the Mac slice and links into the Swift Metal host in
# $(GFX_DIR)/macos, the way every macOS example joins Osprey and Swift.
# [IOS-TARGET-MACOS]
graphics: build _runtime_macos _graphics-shader
	@echo "==> launching Osprey graphics demo: $(SCENE) (close the window to exit)"
	OSPREY_BIN="$(CURDIR)/$(BIN)" bash $(GFX_DIR)/macos/run.sh $(SCENE)

## graphics-smoke: Render every scene for GFX_SMOKE_FRAMES frames and require
##                 each to report them. Needs a desktop session and a Metal GPU.
graphics-smoke: build _runtime_macos _graphics-shader
	@for scene in $(GFX_SCENES); do \
		report=$$(OSPREY_GFX_FRAMES=$(GFX_SMOKE_FRAMES) OSPREY_BIN="$(CURDIR)/$(BIN)" \
			bash $(GFX_DIR)/macos/run.sh $$scene) || exit 1; \
		echo "$$report"; \
		case "$$report" in *": $(GFX_SMOKE_FRAMES) frames in "*) ;; \
			*) echo "graphics-smoke: $$scene did not draw $(GFX_SMOKE_FRAMES) frames" >&2; exit 1 ;; esac; \
	done
endif

## gpt2: Run GPT-2 (usage: make gpt2 FLAVOR=default|python PROMPT="..." TOKENS=20 BACKEND=gpu|lists)
## BACKEND=gpu runs the model on the GPU surface: every kernel and product on GPU=auto|metal|opencl, or on the CPU with GPU=host; BACKEND=lists runs the same model over host lists.
## The published weights and merges download once into examples/gpt2/weights (ignored by git).
GPT2_DIR    := examples/gpt2
GPT2_SOURCE := https://huggingface.co/openai-community/gpt2/resolve/main
PROMPT      ?= The capital of France is
TOKENS      ?= 20
BACKEND     ?= gpu
GPU         ?= auto
FLAVOR      ?= default
GPT2_GENERATOR := $(GPT2_DIR)/$(if $(filter python,$(FLAVOR)),python/generate,generate)

$(GPT2_DIR)/weights/%:
	@$(MKDIR) $(GPT2_DIR)/weights
	@echo "==> downloading $* from openai-community/gpt2"
	curl -fsSL --retry 3 -o $@ $(GPT2_SOURCE)/$*

gpt2: build $(GPT2_DIR)/weights/model.safetensors $(GPT2_DIR)/weights/merges.txt
	@case "$(FLAVOR)" in default|python) ;; *) echo "gpt2 FLAVOR must be default or python" >&2; exit 1 ;; esac
	@printf '%s\n%s\n%s\n' "$(PROMPT)" "$(TOKENS)" "$(BACKEND)" | OSPREY_GPU=$(GPU) ./$(BIN) $(GPT2_GENERATOR) --run --memory=arc

## gpt2-train: One training step of GPT-2 on the prompt: loss, gradients, SGD, loss again (usage: make gpt2-train PROMPT="..." BACKEND=gpu|lists)
gpt2-train: build $(GPT2_DIR)/weights/model.safetensors $(GPT2_DIR)/weights/merges.txt
	@test "$(FLAVOR)" = default || { echo "gpt2-train FLAVOR must be default" >&2; exit 1; }
	@printf '%s\n%s\n' "$(PROMPT)" "$(BACKEND)" | OSPREY_GPU=$(GPU) ./$(BIN) $(GPT2_DIR)/train --run --memory=arc

## run: Compile and run an Osprey file (usage: make run FILE=<path>)
run: build
	@if [ -z "$(FILE)" ]; then echo "Usage: make run FILE=<path>"; exit 1; fi
	./$(BIN) $(FILE) --run

## install: Install osprey + runtime archives system-wide
install: build
	cargo install --path crates/osprey-cli --force
	sudo $(MKDIR) /usr/local/lib
	sudo cp $(RTB)/libfiber_runtime.a $(RTB)/libhttp_runtime.a /usr/local/lib/
	@echo "==> installed osprey and runtime archives."

# _uninstall: Remove osprey + runtime archives from the system
_uninstall:
	cargo uninstall osprey-cli 2>/dev/null || true
	sudo rm -f /usr/local/lib/libfiber_runtime.a /usr/local/lib/libhttp_runtime.a
	@echo "==> uninstalled."

# _website-dev: Start local website development server
_website-dev: $(WEBSITE_NODE_DEPS)
	cd website && npm run dev

# _website-build: Build static site
_website-build: $(WEBSITE_NODE_DEPS)
	cd website && npm run build

# _test_bench_tools: Unit-test the benchmark result merger and the README figure
# generator. Every benchmark run stages and merges validated replacements, so
# the checks guard the recorded data set against failed or incomplete runs;
# `update_readme` is what keeps every
# quoted binarytrees figure in step with that data.
_test_bench_tools:
	@echo "==> Bench tooling tests..."
	python3 benchmarks/test_merge_results.py

## bench: Build, then run the cross-language performance benchmark suite
##        (Osprey vs Rust/C/OCaml/Haskell — CPU via hyperfine, peak memory via
##        /usr/bin/time). Absent toolchains are skipped. Informational only:
##        NOT part of `make ci`/`make test` (perf is noisy on shared runners).
##        Pass a name filter via BENCH_FILTER=<substr>; results in
##        benchmarks/results/results.html. Every run preserves unmeasured
##        values and publishes only validated replacements. See benchmarks/README.md.
bench: build
	@zsh benchmarks/run.sh $(BENCH_FILTER)

## bench-osprey: Re-time ONLY the Osprey columns — default, `--memory=arc`,
##        `--memory=gc` and wasm32 — and merge them into the results already on
##        disk. No other language is compiled, run or timed, so this is the fast
##        loop for a compiler/runtime change. Every non-Osprey measurement in
##        benchmarks/results/ survives byte-for-byte at the data-record level;
##        the Osprey re-run happens in a temp dir and is published only if every
##        case succeeded, so a broken build can never erase the recorded data.
##        Needs one prior `make bench` for the other languages' baseline.
##        Honours BENCH_FILTER=<substr>. See benchmarks/README.md.
bench-osprey: build
	@BENCH_PARTIAL=1 zsh benchmarks/run.sh $(BENCH_FILTER)
