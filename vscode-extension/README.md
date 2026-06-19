<p align="center">
  <img src="https://raw.githubusercontent.com/MelbourneDeveloper/osprey/main/website/src/assets/images/logo.png" alt="Osprey logo" width="128" />
</p>

# Osprey for VS Code

> **Preview.** Osprey is pre-production and evolving fast. Expect rough edges.

Language support for [Osprey](https://ospreylang.dev) — a functional programming
language with algebraic effects, fiber-based concurrency, pattern matching, and
strong compile-time safety.

Powered by a Rust language server (`osprey lsp`, built on
[lspkit](https://github.com/Nimblesite/lspkit)) that runs the compiler front-end
in-process — the same engine targeted at Neovim and Zed next.

## Features

- **Syntax highlighting** for `.osp` files — keywords, types, string
  interpolation (`"Hello ${name}!"`), operators, and comments.
- **Live diagnostics** — errors and warnings from the Osprey compiler as you
  type, inline in the editor.
- **Hover, go-to-definition, find-references, document symbols, signature help,
  and completion** — driven by the compiler's own parser and type checker.
- **Compile & run** from the editor:
  - `Osprey: Compile Osprey File` (`Ctrl/Cmd+Shift+B`)
  - `Osprey: Compile and Run Osprey File` (`F5`)
- **Bracket matching, auto-closing, and comment toggling.**

## Requirements

The extension bundles a version-matched Osprey compiler for your platform and
verifies it at startup, so syntax checking works out of the box.

To **compile and run** programs, Osprey invokes LLVM and a C toolchain, so install:

- **LLVM** (provides `llc`) — `brew install llvm` / `scoop install llvm`
- A C compiler — `clang` (macOS/Linux) or MinGW `gcc` (`scoop install gcc`)

Or install the full toolchain via a package manager (this also puts `osprey` on
your `PATH`):

```bash
brew install nimblesite/tap/osprey            # macOS / Linux
scoop bucket add nimblesite https://github.com/Nimblesite/scoop-bucket && scoop install osprey   # Windows
```

## Settings

| Setting | Default | Description |
|---------|---------|-------------|
| `osprey.server.enabled` | `true` | Enable/disable the language server. |
| `osprey.diagnostics.enabled` | `true` | Enable/disable inline diagnostics. |
| `osprey.server.compilerPath` | `""` | Path to an Osprey compiler. **Leave empty** to use the version-matched compiler bundled with this extension (falling back to `osprey` on `PATH`). |

## Links

- Website & docs: <https://ospreylang.dev>
- Source & issues: <https://github.com/MelbourneDeveloper/osprey>

## License

See [LICENSE](LICENSE).
