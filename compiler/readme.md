# Osprey Compiler Assets

The compiler itself is the Rust workspace at the repository root
([crates/](../crates/), binary `osprey`). This directory holds the assets
compiled programs and the test suites depend on:

- [runtime/](runtime/) - The C runtime (fibers, effects, HTTP/WebSocket,
  strings, lists, maps, JSON, terminal). Built into `lib/libfiber_runtime.a` /
  `lib/libhttp_runtime.a` by `make _runtime`; `osprey` links these at
  `--compile`/`--run` time.
- [examples/tested/](examples/tested/) - Golden examples; each `.osp` runs in
  CI and must match its `.expectedoutput` byte-for-byte
  ([crates/diff_examples.sh](../crates/diff_examples.sh)).
- [examples/failscompilation/](examples/failscompilation/) - Ill-formed
  programs the compiler must reject (the must-reject ratchet).
- [examples/tui/](examples/tui/) - The interactive TUI showcase (`make tui`).

## Related Projects

- [VSCode Extension](../vscode-extension/) - Language support for VS Code
- [Web Compiler](../webcompiler/) - Browser-based playground
- [Website](../website/) - Official documentation site

## Development Notes

- Don't put .osp files in the root of this folder
- This is for runtime/examples only - no loose project files
