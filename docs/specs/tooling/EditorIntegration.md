# Editor integration and distribution

The corresponding LanguageServerAndEditors contract is [here](0020-LanguageServerAndEditors.md).

## Editor integrations

The VS Code integration is a thin client over `[LSP-TRANSPORT]`.

### VS Code `[EDITOR-VSCODE]`

- Extension id `nimblesite.osprey`; client in
  [`vscode-extension/client/src/extension.ts`](../../../vscode-extension/client/src/extension.ts)
  spawns `osprey lsp` over stdio.
- Packaged as a **per-platform VSIX** (`darwin-arm64`, `linux-x64`,
  `win32-x64`). Each VSIX **bundles** a version-matched `osprey` binary + runtime
  libs + a stamped `shipwright.json`, verified present inside the package at build
  time.
- Client resolves the server command in priority order: user setting
  (`osprey.server.compilerPath`) → bundled binary → `PATH` (per the Shipwright
  `sources` list in [`shipwright.json`](../../../shipwright.json)).
- The extension's native DAP integration is specified in
  [Debugger](0021-Debugger.md); it is separate from the LSP request path.
- Marketplace publication uses **OIDC** (no PAT) — see `[EDITOR-VERSIONING]` and
  the [release workflow](../../../.github/workflows/release.yml). Open VSX
  publication uses the same VSIX artifacts and an independent optional-token
  job, so either registry can succeed alone.

## Versioning & supply chain `[EDITOR-VERSIONING]`

The VS Code distributions obey the
[Shipwright](https://github.com/Nimblesite/Shipwright) version contract: the
extension and bundled binary MUST be version-matched.

- The binary is the source of truth: `osprey --version` → `osprey X.Y.Z`;
  `osprey --version --json` → the version manifest (`[SWR-VERSION-CLI-OUTPUT]`).
- Components are declared in [`shipwright.json`](../../../shipwright.json):
  `osprey` (the CLI — which _is_ the language server, via the `lsp`
  subcommand) and `osprey-vscode`. The component id **must** equal the name the
  binary reports from `osprey --version` (Shipwright matches the probed name
  against the component id), so the CLI component is `osprey`, not
  `osprey-compiler`. The LSP is **not** a separate component; it is
  the same binary, so no separate version surface exists to drift.
- Source version fields stay at `0.0.0-dev`; the release version is stamped
  from the git tag at build time (`[SWR-VERSION-BUILD-STAMPING]`). Hard-coding
  a version is a defect.
- VS Code activation verifies the bundled compiler against the manifest and
  prompts to reinstall on mismatch (`hosts.vscode.onMismatch`). PATH/registry
  sources are verified at startup (`verifyStartup`).
- Marketplace publishing uses GitHub OIDC and Microsoft Entra
  workload-identity federation, with no stored PAT.
