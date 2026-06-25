# Releasing Osprey

Osprey ships through a single **tag-triggered** pipeline
([`.github/workflows/release.yml`](../.github/workflows/release.yml)) built on the
[Shipwright](https://github.com/Nimblesite/Shipwright) supply-chain contract.

## The three trigger rules

| Event | What runs |
|-------|-----------|
| **Push a tag `v*`** | The full release: build → GitHub Release → Homebrew + Scoop + Marketplace → website. |
| **Open a PR to `main`** | CI only ([`ci.yml`](../.github/workflows/ci.yml) + [`ci-windows.yml`](../.github/workflows/ci-windows.yml)). |
| **Merge to `main`** | **Nothing.** No build, no deploy. |

## Cutting a release

```bash
# from an up-to-date main
git tag v1.2.3
git push origin v1.2.3
```

That's it. The pipeline does the rest:

1. **Resolve version** — `1.2.3` is derived from the tag and validated against
   `shipwright.json`.
2. **Build** the compiler for `darwin-arm64`, `darwin-x64`, `linux-x64`, and
   `win32-x64`, verify `osprey --version` prints `osprey 1.2.3`, and package each
   as `osprey-1.2.3-<platform>.tar.gz` (binary + runtime libs) + `.sha256`.
3. **GitHub Release** on `Nimblesite/osprey` with all tarballs.
4. **Homebrew** — writes `Formula/osprey.rb` to
   [`Nimblesite/homebrew-tap`](https://github.com/Nimblesite/homebrew-tap).
5. **Scoop** — writes `bucket/osprey.json` to
   [`Nimblesite/scoop-bucket`](https://github.com/Nimblesite/scoop-bucket).
6. **VS Code Marketplace** — publishes the per-platform VSIX as
   `nimblesite.osprey` (binary bundled and version-checked at activation).
7. **Website + web compiler** deploy — only after every step above succeeds.

## Versioning — never hard-code it ([SWR-VERSION-BUILD-STAMPING])

Source-controlled version fields MUST stay at the placeholder **`0.0.0-dev`**:

- `Cargo.toml` (`[workspace.package] version`) and the CLI fallback in
  `crates/osprey-cli/src/main.rs`
- `vscode-extension/package.json` (`version`)
- `shipwright.json` (`product.version` + each component `expectedVersion`)

The real version is stamped from the tag at build time — the `osprey` binary
via the `OSPREY_VERSION` environment variable at `cargo build` time,
`package.json`/`shipwright.json` via the release job. **A PR that changes a placeholder to a real version is a defect and
must be rejected in review.**

The compiler honors the version contract ([SWR-VERSION-CLI-OUTPUT]):

```text
$ osprey --version
osprey 1.2.3
$ osprey --version --json
{"manifestVersion":1,"name":"osprey","version":"1.2.3","kind":"cli","product":"osprey"}
```

## Required secrets / variables

Configure these on `Nimblesite/osprey`:

| Secret | Used by | Purpose |
|--------|---------|---------|
| `VSCE_PAT` | `vsix` job | VS Code Marketplace PAT for publisher `nimblesite`. |
| `TAP_TOKEN` | `brew` job | PAT with push to `Nimblesite/homebrew-tap`. |
| `SCOOP_BUCKET_TOKEN` | `scoop` job | PAT with push to `Nimblesite/scoop-bucket`. |
| `FLY_API_TOKEN` | web compiler | Fly.io deploy (existing). |

The GitHub Release uses the built-in `GITHUB_TOKEN` (`contents: write`).
Set repo **variable** `SKIP_VSCE_PUBLISH=true` to dry-run the VSIX build without
publishing to the Marketplace.

## Installing released builds

```bash
# Homebrew (macOS / Linux)
brew install nimblesite/tap/osprey

# Scoop (Windows)
scoop bucket add nimblesite https://github.com/Nimblesite/scoop-bucket
scoop install osprey
```

Osprey shells out to LLVM (`llc`) and a C compiler (`clang`/`gcc`) at compile
time, so those are package-manager dependencies (`llvm` for brew; `llvm` + `gcc`
for scoop).

## Windows support status

The Windows build is delivered in phases (see the `[WINDOWS-PORT-*]` markers and
the C runtime under `compiler/runtime/`):

- **Phase 1 (shipped):** core language — collections, strings, fibers (via
  winpthreads), effects, pattern matching. Built under MSYS2 UCRT64.
- **Phase 2 (in progress):** HTTP / WebSocket via Winsock2.
- **Phase 3 (planned):** process spawning via the Win32 process APIs
  (`CreateProcess`); currently stubbed on Windows.
