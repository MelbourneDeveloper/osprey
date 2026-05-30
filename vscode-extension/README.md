# Osprey for VS Code

Language support for [Osprey](https://ospreylang.dev) — a functional programming
language with algebraic effects, fiber-based concurrency, pattern matching, and
strong compile-time safety.

## Features

- **Syntax highlighting** for `.osp` files — keywords, types, string
  interpolation (`"Hello ${name}!"`), operators, and comments.
- **Live diagnostics** — errors and warnings from the Osprey compiler as you
  type, inline in the editor.
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
