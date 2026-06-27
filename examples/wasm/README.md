# Osprey → WebAssembly

A minimal end-to-end example: an Osprey program compiled to `wasm32-wasip1` and
run under `wasmtime`, Node's WASI, and in the browser. See
[`docs/specs/0022-WebAssemblyTarget.md`](../../docs/specs/0022-WebAssemblyTarget.md)
for the design.

## Prerequisites

- `clang` with the wasm32 backend (any recent LLVM)
- `wasm-ld` — from LLVM's `lld` (`brew install lld`, or `apt-get install lld`)
- A WASI sysroot — `brew install wasi-libc`, or the
  [wasi-sdk](https://github.com/WebAssembly/wasi-sdk). Override the autodetected
  location with `OSPREY_WASI_SYSROOT=/path/to/wasi-sysroot`.
- `node` (for the smoke test) and, optionally, `wasmtime` (for `--run`).

## Build & run

From the repo root:

```sh
# One target builds everything ready to go: the wasm runtime archive, the
# compiled example, and validation + smoke-runs (Node's WASI and the browser shim).
make wasm
```

Or drive the compiler directly:

```sh
# compile to WebAssembly
osprey examples/wasm/hello.osp --target=wasm32 --compile -o examples/wasm/build/hello.wasm

# run it under a WASI host
wasmtime examples/wasm/build/hello.wasm
osprey examples/wasm/hello.osp --target=wasm32 --run     # uses wasmtime under the hood
node scripts/wasm-smoke.mjs         examples/wasm/build/hello.wasm examples/wasm/hello.expectedoutput  # Node's WASI host
node scripts/wasm-browser-smoke.mjs examples/wasm/build/hello.wasm examples/wasm/hello.expectedoutput  # the browser's WASI shim
```

Expected output ([`hello.expectedoutput`](hello.expectedoutput)):

```
Osprey on WebAssembly
length = 21
2 * 1500000000 = 3000000000
```

## In the browser

`index.html` ships a tiny inline WASI shim (`fd_write` → page + console), so no
bundler or npm packages are needed. Serve this directory over HTTP and open it:

```sh
cd examples/wasm && python3 -m http.server 8080
# then visit http://localhost:8080/  and click “Run hello.wasm”
```

From the devtools console you can drive the module directly:

```js
await osprey.run()   // (re)instantiate and run _start
osprey.exports       // raw wasm exports: memory, _start
```

## What works / what doesn't

The target is `wasm32-wasip1`; libc (malloc, `sprintf`, `printf`/`puts` → WASI
`fd_write`) comes from wasi-libc, so the portable runtime subset — allocator,
strings, lists, maps, JSON, effects — runs unchanged. **Not** supported on wasm
yet: fibers/`spawn` (pthreads), HTTP/WebSocket (sockets/OpenSSL), and FFI
(`dlopen`). A program using those fails at link with an undefined-symbol error.
