# WebAssembly Target

Osprey compiles to WebAssembly so programs run in the browser. The backend
reuses the existing LLVM-IR pipeline unchanged and adds a target selector
(`osprey --target=wasm32`), a wasm-portable runtime archive, and a `wasm-ld`
link step. The output is a `wasm32-wasip1` command module that runs under any
WASI host — `wasmtime`, Node's `node:wasi`, or a browser WASI shim. [WASM-TARGET]

## Status

Implemented for the portable language core. `osprey --target=wasm32 --compile`
emits a validated `.wasm` that prints correctly under wasmtime, Node's WASI, and
in the browser (`examples/wasm/`). Of the tested example suite, **30/48 run on
wasm with byte-identical stdout**; the other 18 use a non-portable feature and
skip (see below). The CI `wasm` job gates on both the browser-loadable example
(`wasm-validate` + Node-WASI stdout) and the full golden suite
(`crates/diff_wasm_examples.sh`, FAIL=0).

Not yet ported (link-time `undefined symbol`, by design — see Limitations):
fibers/`spawn` (pthreads), HTTP/WebSocket (sockets/OpenSSL), FFI (`dlopen`),
and the `random`/`input` builtins (CSPRNG/stdin syscalls).

## Design

### Target triple: `wasm32-wasip1` [WASM-TARGET-TRIPLE]

`wasm32-wasip1` (the modern spelling of `wasm32-wasi`) is chosen over
`wasm32-unknown-unknown`. Osprey's runtime needs libc — `malloc`, `sprintf`,
`snprintf`, `strcmp`, `memcpy` — and a `stdout`. wasi-libc supplies all of them,
so the portable C runtime compiles **unchanged** and `print` works via WASI
`fd_write` with no hand-written libc. The browser still runs the module through
a small WASI shim (`examples/wasm/index.html`), so "needs WASI" does not mean
"can't run in a browser".

### Codegen is unchanged — the IR was already portable [WASM-TARGET-IR]

The textual LLVM IR carries no target triple or datalayout (clang supplies the
target's), `int` is `i64`, and pointers round-trip through `i64` (the uniform
machine-word boxing). On wasm32 a pointer is 32 bits, so that round-trip
zero-extends and truncates losslessly — addresses fit in 32 bits. No byte
offsets are hard-coded; LLVM computes struct layout per the target datalayout.
The backend needed **no** pointer-width refactor.

### Three width fixes (correctness on ILP32) [WASM-TARGET-WIDTH]

The IR baked in host type widths that differ on wasm32 (ILP32), each invisible
on LP64 (where `i8*` and `i64` are both 8 bytes) but wrong on wasm32:

1. **`size_t` (string length).** Codegen called libc `strlen`, declared `i64`,
   but `size_t` is 32-bit on wasm32 — a `wasm-ld` signature mismatch that traps.
   Fixed by routing the `length` builtin and string concatenation through a
   runtime shim `osp_strlen(const char*) -> int64_t`, so the `size_t -> int64`
   cast lives in C (correct per target) and the IR stays `i64` everywhere.
2. **`long` (integer formatting).** Integer→string used `sprintf("%ld", i64)`;
   on wasm32 `long` is 32-bit, truncating large values. Fixed by `"%lld"`
   (`long long` is 64-bit on every target; identical to `%ld` on LP64).
3. **`Result` success-slot type.** The `Result<T, E>` block is `{ T, i8 disc,
   i8* errmsg }`. An `Error { message }` constructor typed its success slot from
   the *message* (`i8*`), but a function declared `-> Result<int, _>` is read
   back with an `i64` slot. On LP64 both are 8 bytes so it worked by accident;
   on wasm32 the 4-byte `i8*` slot shifts the `disc`/`errmsg` offsets, flipping
   `Error` to `Success` with a garbage payload. Fixed by re-laying a returned
   `Result` to the declared success-slot type (`result::repack_to_inner`, called
   from `coerce_return`) so producer and reader agree on the layout.

All three are width-stable improvements that leave native output byte-identical
(the differential golden suite is unchanged at 48/48).

### Entry point: a `__main_void` thunk [WASM-ENTRY]

wasi-libc's `crt1-command.o` `_start` calls `__main_void`. The wasm driver
appends a thunk `define i32 @__main_void() { call i32 @main() ... }` to the IR,
sidestepping the clang/wasi-libc `main`-mangling skew. The result is a command
module (`_start` → `__wasm_call_ctors` → `__main_void` → `@main`) that runs
identically under wasmtime, Node's WASI, and a browser shim.

### Linking: drive `wasm-ld` directly [WASM-TARGET-LINK]

clang lowers the IR to a wasm object (`-c` only) and the driver then calls
`wasm-ld` itself: `crt1-command.o` + object + `libosprey_runtime_wasm.a`
(on-demand) + `-lc`. Going straight to `wasm-ld` avoids clang's auto-added
`libclang_rt.builtins-wasm32.a`, which stock Homebrew/apt LLVM doesn't ship and
which the portable core doesn't need (wasm has native `i64`). Sysroot and tool
locations are discovered from `OSPREY_WASI_SYSROOT` / `OSPREY_WASM_LD` or the
conventional Homebrew / wasi-sdk / Linux paths.

### Runtime subset [WASM-TARGET-RUNTIME]

`make _runtime_wasm` cross-compiles the portable C units — allocator, strings,
list/map containers, JSON, effects — to `libosprey_runtime_wasm.a`. The
allocator is the same `malloc` passthrough (`@osp_alloc`); reclamation stays
unobservable [MEM-OPAQUE]. Non-portable units (fibers, HTTP/WebSocket, system,
terminal, FFI, CSPRNG) are excluded; because archives link on demand, a program
that does not reference their symbols links cleanly, and one that does fails at
link with a clear `undefined symbol`.

## Limitations

- **No fibers/HTTP/WebSocket/FFI/`random`/`input`.** These depend on
  pthreads / sockets / OpenSSL / `dlopen` / syscalls absent under
  `wasm32-wasip1`. A program using them fails at link, not silently.
- **WASI in the browser** needs a shim (`examples/wasm/index.html` ships a
  ~60-line inline one mapping `fd_write` to the page/console). A future
  `wasm32-unknown-unknown` mode could import I/O from JS directly.
- **No reclamation on wasm** beyond the optimizer's static frees, same as
  native default [MEM-OPAQUE].

## Verification

- `osprey examples/wasm/hello.osp --target=wasm32 --compile -o hello.wasm`
- `wasm-validate hello.wasm` — structural well-formedness
- `node scripts/wasm-smoke.mjs hello.wasm examples/wasm/hello.expectedoutput`
  — runs under Node's WASI and asserts stdout
- `examples/wasm/index.html` — loads and runs it in the browser, output to page
- `zsh crates/diff_wasm_examples.sh` — the golden suite: compile every tested
  example to wasm, run under Node's WASI, diff stdout; non-portable examples
  (undefined symbol) are SKIPped. Reports `PASS=30 FAIL=0 SKIP=18`.
- CI `wasm` job runs the validate + Node-WASI smoke **and** the golden suite on
  every PR.
