// Shared test-environment resolution for the Osprey VSIX suites: locate the
// freshly-built osprey compiler and the lldb-dap adapter. Defined once and
// imported by every suite so no binary-resolution logic is duplicated.

import * as assert from "assert";
import * as fs from "fs";
import * as path from "path";
import { resolveLldbDapExecutable } from "../../client/src/extension";

// The compiled test lives at <ext>/out/test/suite, so the extension root is
// three levels up.
export const extensionRoot = path.resolve(__dirname, "..", "..", "..");

/** Absolute path of the `osprey` binary staged on PATH, or undefined. */
export function resolveOspreyOnPath(): string | undefined {
  const exe = process.platform === "win32" ? "osprey.exe" : "osprey";
  for (const dir of (process.env.PATH ?? "").split(path.delimiter)) {
    if (!dir) {
      continue;
    }
    const candidate = path.join(dir, exe);
    if (fs.existsSync(candidate)) {
      return candidate;
    }
  }
  return undefined;
}

/**
 * The osprey binary tests should launch: PREFER the freshly-built dev compiler
 * under the repo (target/release/osprey — the `make build` output that speaks
 * the current protocol) over any older `osprey` on PATH. Falls back to the PATH
 * binary, then undefined.
 */
export function resolveBuiltOsprey(): string | undefined {
  const exe = process.platform === "win32" ? "osprey.exe" : "osprey";
  const built = path.resolve(extensionRoot, "..", "target", "release", exe);
  return fs.existsSync(built) ? built : resolveOspreyOnPath();
}

/** Resolve lldb-dap for the debugger E2E, failing the test if it is absent. */
export function resolveRequiredLldbDap(): string {
  const command = resolveLldbDapExecutable({});
  if (command && fs.existsSync(command)) {
    return command;
  }
  assert.fail(
    `lldb-dap is required for the Osprey VSIX debugger E2E; resolved "${command}" but it was not executable`,
  );
}
