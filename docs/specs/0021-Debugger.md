# Debugger

> **Engineering spec** (tooling), not part of the `0001`-`0019` language
> reference. It defines how Osprey programs are built, launched, and inspected
> by debuggers.

The Osprey debugger is a source-level debugging system for native Osprey
programs. It is integrated with editors through the same extension surface as
the language server, but it uses a different protocol: LSP is the static
analysis plane; DAP is the runtime control plane. The implementation plan is
[Plan 0012](../plans/0012-osprey-debugger.md).

## Status

| Capability                          | State                                                                          |
| ----------------------------------- | ------------------------------------------------------------------------------ |
| Debug build mode (`osprey --debug`) | In progress. Emits LLVM/DWARF metadata for native builds.                      |
| VS Code DAP launch                  | In progress. Compiles `.osp` to a native debug binary and launches `lldb-dap`. |
| Variables / value rendering         | Planned. Requires local-variable metadata and Osprey runtime inspectors.       |
| Fibers / effects inspection         | Planned. Requires runtime debug APIs.                                          |
| Replay / time travel                | Planned. Requires deterministic runtime event recording.                       |

## Protocol Split `[DEBUGGER-PROTOCOLS]`

Osprey uses two editor protocols:

- **LSP** (`osprey lsp`) owns editor-time analysis: diagnostics, hover,
  symbols, definition, completion, and source position normalization.
- **DAP** owns runtime control: launch, breakpoints, stepping, stack traces,
  scopes, variables, evaluate, pause, and terminate.

The debugger MUST NOT fake a debug session by canceling DAP and running
`osprey --run`. The `osprey.run` command is a run command; F5 is a debugger
launch.

Both planes MUST agree on source identity and positions. AST/source positions
used by LSP are also the provenance for emitted debug metadata.

## Debug Build Contract `[DEBUGGER-BUILD]`

`osprey --debug --compile` builds a native executable suitable for source-level
debugging.

Required behavior:

- `--debug` is accepted by `--llvm`, `--compile`, and `--run`.
- Native debug builds emit LLVM debug metadata that lowers to DWARF.
- Native debug builds pass debugger-friendly driver flags (`-g`, no omitted
  frame pointer where supported).
- Native debug builds default to no optimization (`-O0`) unless an explicit
  debug optimization override is supplied.
- Non-debug builds keep their release-oriented defaults.
- `--debug --target=wasm32` is rejected until WebAssembly debug information is
  specified and tested.

Minimum emitted metadata:

- `source_filename`.
- `!llvm.dbg.cu`.
- `!llvm.module.flags` including debug-info version and DWARF version.
- `!DIFile`.
- `!DICompileUnit`.
- `!DISubprogram` for user functions and generated `main`.
- `!DILocation` on instructions derived from executable source statements.

## Source Mapping `[DEBUGGER-SOURCE-MAP]`

The parser and lowerers must preserve source positions for executable
statements and declarations.

Rules:

- Osprey AST positions use 1-based lines and 0-based columns.
- DAP/source debugger positions exposed to users use 1-based lines and columns.
- Compiler-generated code may be associated with the nearest source statement
  only when doing so improves stepping/breakpoint behavior.
- Generated helper frames should be hidden from normal stepping once smart
  stepping exists.

## Editor Launch `[DEBUGGER-EDITOR-LAUNCH]`

For VS Code:

1. The debug provider resolves the `.osp` source file from the active editor or
   launch configuration.
2. Dirty documents are saved or the debug launch is rejected.
3. The provider runs the version-matched compiler:

   ```text
   osprey <source.osp> --debug --compile -o <debug-binary>
   ```

4. The provider launches a real DAP adapter, initially `lldb-dap`, against the
   compiled native binary.
5. DAP handles breakpoints, stepping, stack, scopes, and variables.

The extension may let users configure:

- Osprey compiler path.
- LLDB-DAP path.
- Debug output path.
- Program args.
- Working directory.
- Environment variables.
- Stop-on-entry.

## Reusable Debugger Helpers `[DEBUGGER-REUSE]`

Generic debugger utilities MUST live outside Osprey-specific compiler modules.
The `osprey-debug` crate holds small editor/compiler-neutral primitives such as
debug source identity and native debug-build policy. It intentionally avoids
Osprey parser, type-checker, codegen, and editor dependencies so useful pieces
can later move into `lspkit`.

Osprey-specific lowering remains in `osprey-codegen`; editor-specific launch
logic remains in the VS Code extension.

## Future Runtime Inspection `[DEBUGGER-RUNTIME]`

The finished debugger must inspect Osprey values, not just native pointers.

Required future support:

- Local variables and parameters via `DILocalVariable` and `dbg.value` /
  `dbg.declare`.
- Records, unions, `Result`, strings, lists, maps, closures, fibers, channels,
  and effect handlers rendered as Osprey values.
- Safe runtime inspection helpers for opaque handles.
- Fiber and effect runtime debug ids.
- Replayable scheduler/effect/IO event streams.

These features are not allowed to guess raw memory layouts ad hoc from the
editor. Stable runtime inspection APIs are required.

## Conformance

A change is conformant only if:

1. `osprey --debug --llvm` emits the minimum debug metadata in
   `[DEBUGGER-BUILD]`.
2. `osprey --debug --compile` produces a native executable that a supported DAP
   adapter can launch.
3. The VS Code debugger contribution starts a DAP session; it does not proxy to
   `osprey --run`.
4. LSP and debugger source positions follow `[DEBUGGER-SOURCE-MAP]`.
5. Generic debugger utilities remain isolated per `[DEBUGGER-REUSE]`.
