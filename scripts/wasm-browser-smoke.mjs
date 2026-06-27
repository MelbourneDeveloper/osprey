// Browser-path smoke test for an Osprey-compiled WebAssembly module. [WASM-TARGET]
//
// Usage: node scripts/wasm-browser-smoke.mjs <module.wasm> [expected-stdout-file]
//
// Complements scripts/wasm-smoke.mjs (which runs under Node's WASI host) by
// exercising the EXACT inline WASI shim the browser uses — examples/wasm/
// wasi-shim.mjs — so a regression in the browser loader is caught in CI without
// launching a browser. Exits non-zero on trap or stdout mismatch.

import { readFile } from "node:fs/promises";
import { runModule } from "../examples/wasm/wasi-shim.mjs";

const [, , wasmPath, expectedPath] = process.argv;
if (!wasmPath) {
  console.error("usage: node wasm-browser-smoke.mjs <module.wasm> [expected-stdout-file]");
  process.exit(2);
}

const bytes = await readFile(wasmPath);
if (!WebAssembly.validate(bytes)) {
  console.error(`FAIL: ${wasmPath} is not a valid WebAssembly module`);
  process.exit(1);
}

let captured = "";
try {
  await runModule(bytes, (text) => {
    captured += text;
  });
} catch (err) {
  console.error(`FAIL: module trapped under the browser shim: ${err?.message ?? err}`);
  process.exit(1);
}
process.stdout.write(captured);

if (expectedPath) {
  const expected = (await readFile(expectedPath, "utf8")).trim();
  if (captured.trim() !== expected) {
    console.error("FAIL: browser-shim stdout mismatch");
    console.error(`  expected: ${JSON.stringify(expected)}`);
    console.error(`  actual:   ${JSON.stringify(captured.trim())}`);
    process.exit(1);
  }
}

console.error(`OK: ${wasmPath} ran cleanly under the browser WASI shim`);
