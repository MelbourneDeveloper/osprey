// Single source of truth for the Osprey mark: website/src/assets/images/logo.png.
//
// The VS Code extension needs a physical icon file inside the package — vsce
// cannot reference a URL for the Marketplace icon or the `.osp` file icon — so we
// GENERATE it from the canonical logo at build time rather than committing a copy.
// Wired into `precompile` and `vscode:prepublish`; the output (icon.png) is
// gitignored. This keeps exactly one logo in version control.
import { copyFileSync, existsSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, resolve } from "node:path";

const here = dirname(fileURLToPath(import.meta.url));
const SOURCE = resolve(here, "../../website/src/assets/images/logo.png");
const DEST = resolve(here, "../icon.png");

if (!existsSync(SOURCE)) {
  console.error(`sync-icon: canonical logo not found at ${SOURCE}`);
  process.exit(1);
}

try {
  copyFileSync(SOURCE, DEST);
  console.log(`sync-icon: ${SOURCE} -> ${DEST}`);
} catch (err) {
  console.error(`sync-icon: failed to copy icon: ${err.message}`);
  process.exit(1);
}
