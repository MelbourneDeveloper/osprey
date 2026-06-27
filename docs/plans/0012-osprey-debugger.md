# Plan 0012 - Modern Osprey Debugger

## Summary

Osprey needs a real debugger, not a run button hidden behind VS Code's
debugger UI. The debugger must support source breakpoints, stepping, stack
inspection, variables, Osprey value rendering, expression evaluation, fibers,
effects, replay, and machine-checkable debug-info quality. It should integrate
with existing debugger infrastructure first, then add Osprey-specific semantics
where generic native debuggers cannot know the language.

The architecture is layered:

1. Emit standards-based source debug information from the compiler.
2. Compile debug builds as native executables that LLDB/GDB can understand.
3. Use DAP for editor integration and headless automation.
4. Add an Osprey debug shim for language-specific values, fibers, effects, and
   expression evaluation.
5. Add deterministic replay/time-travel after the runtime can record the needed
   events.

The first implementation slice is intentionally small: line tables,
`--debug`, and VS Code launch through LLDB-DAP. This plan is not small. It is
the complete target system.

## Research basis

Only authoritative sources are accepted here: official standards/docs, ACM
published papers, and author/publisher paper pages. Short quotes are included
as anchors; the implementation decisions below are paraphrased requirements,
not copied text.

| Source                                                                                                                                                               | Authority               | Relevant anchor                             |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- | ------------------------------------------- |
| DWARF v5 standard, DWARF Committee, 2017 - https://dwarfstd.org/doc/DWARF5.pdf                                                                                       | Debug-info standard     | "accurate picture of the source program"    |
| LLVM Source Level Debugging docs - https://llvm.org/docs/SourceLevelDebugging.html                                                                                   | LLVM producer docs      | "source-language's AST map onto LLVM code"  |
| Debug Adapter Protocol - https://microsoft.github.io/debug-adapter-protocol/                                                                                         | DAP spec                | "defines the abstract protocol"             |
| VS Code debugger extension guide - https://code.visualstudio.com/api/extension-guides/debugger-extension                                                             | VS Code API docs        | "debug adapter which is a separate process" |
| LLDB-DAP docs - https://lldb.llvm.org/resources/lldbdap.html                                                                                                         | LLDB project docs       | "`lldb-dap` exposes LLDB's functionality"   |
| Kell and Stinnett, Onward! 2024, "Source-Level Debugging of Compiler-Optimised Code: Ill-Posed, but Not Impossible" - https://dl.acm.org/doi/10.1145/3689492.3690047 | ACM peer-reviewed paper | "loss of state"                             |
| Huang, Liang, Su, Zhang, PLDI 2025, "Robustifying Debug Information Updates in LLVM via Control-Flow Conformance Analysis" - https://dl.acm.org/doi/10.1145/3729267  | ACM peer-reviewed paper | "debug location updates"                    |
| Stinnett and Kell, OOPSLA 2026, "Debugging Debugging Information using Dynamic Call Trees" - https://dl.acm.org/doi/10.1145/3798213                                  | ACM peer-reviewed paper | "observing a running source program"        |
| Lu, Liu, Wang, Zhang, FSE 2024, "DTD: Comprehensive and Scalable Testing for Debuggers" - https://dl.acm.org/doi/10.1145/3643779                                     | ACM peer-reviewed paper | debugger conformance testing                |

Implications for Osprey:

- Use DWARF as the native interchange format on Unix-like platforms. It is the
  standard consumer contract for LLDB/GDB and already models source languages,
  scopes, functions, variables, line tables, and locations. Emit the version
  the target toolchain best supports: DWARF 4 on macOS (Apple `dsymutil`/LLDB
  lag on v5), DWARF 5 elsewhere.
- Emit debug metadata through LLVM IR metadata, not a sidecar map that only our
  extension understands. Osprey's backend already emits textual LLVM IR, so the
  right boundary is `!DICompileUnit`, `!DISubprogram`, `!DILocation`,
  `!DILocalVariable`, and value-location records. The long-term target is
  LLVM's `#dbg_value` / `#dbg_declare` **debug records**; the current textual IR
  backend may emit the older `@llvm.dbg.*` intrinsic compatibility form only
  while the pinned local LLVM toolchain accepts and lowers it correctly.
- Use DAP for IDEs, but do not define a new DAP dialect unless required.
  Reuse `lldb-dap` first; add Osprey-only requests only behind the standard
  adapter boundary.
- Optimized-code debugging is **ill-posed, not impossible** (Kell & Stinnett,
  Onward! 2024): there is no single faithful source view of optimized state, so
  the answer is an explicit choice of view plus residual state/computation on
  demand — not "give up" and not `-Og`-style de-optimization. Pragmatically,
  debug builds start at `-O0`; the correctness criterion for emitted local/line
  info is the paper's **per-variable history oracle** (every local's observed
  value-change history matches the source program's, each change at the correct
  file and line; invariant under DCE/CSE/code motion). Never fabricate a value:
  report it unavailable, and do not conflate "optimized out" with a debug-info
  emission bug.
- When Osprey gains its own IR transforms, debug-location updates must obey the
  **control-flow conformance** invariant (Huang et al., PLDI 2025): an optimized
  path's location set must be a _subset_ of the corresponding source path's. Per
  moved/cloned/merged instruction choose Preserve / Merge / Drop accordingly;
  never attach a source location to a path that did not already carry it, and
  prefer Drop over a guessed location. (Today Osprey optimizes at the clang
  level, so this binds clang/LLVM; it binds us the moment we add passes.)
- Validate stack/scope correctness against the **source-level dynamic call
  tree**, comparing an optimized run to its unoptimized reference (Stinnett &
  Kell, OOPSLA 2026). This is strictly stronger than the weak "backtrace
  invariant," which reversed or empty backtraces can satisfy. Classify each
  divergence as a compiler bug vs a debug-info _format_ limitation, use DWARF
  `inlined_subroutine` records and location views, and mark synthesized frames
  artificial so debuggers can hide them.
- Test the debugger itself by **cross-debugger differential testing** over the
  _same_ binary (DTD; Lu et al., FSE 2024): drive two independent adapters,
  instruction-step, and diff full program state including registers, treating
  optimized-out / error / not-found as _distinct_ states. Manual "looks OK in
  VS Code" is not acceptance.

## User-facing requirements

The finished debugger must support these workflows:

- Press F5 in VS Code on an `.osp` file and hit source breakpoints.
- Run `osprey file.osp --debug --compile` to produce a native debug binary.
- Run `osprey file.osp --debug --run` and preserve debug info for crash/stack
  analysis.
- Launch an Osprey program with args, cwd, env, stdin/stdout/stderr terminal
  configuration, and optional stop-on-entry.
- Attach to an existing Osprey process when symbols are present.
- Load a core dump when the platform/debugger supports it.
- Use breakpoints: source line, conditional, hit count, logpoint, function,
  exception/panic-like runtime trap, and data watchpoint where supported.
- Step over, step into, step out, continue, pause, restart, and stop.
- Show call stacks in Osprey names, not only lowered helper symbols.
- Show locals, parameters, captured variables, module globals, and runtime
  handles.
- Render Osprey values: ints, floats, bools, strings, Unit, records, unions,
  Result, lists, maps, fibers, channels, effects, closures, and Ptr/FFI values.
- Evaluate Osprey expressions in the current frame, with side-effect controls.
- Inspect fibers as logical tasks even when the OS sees one thread.
- Inspect active effect handlers and resumptions.
- Record/replay deterministic executions for concurrency/effects bugs.
- Work headlessly in CI so debug behavior is tested without a GUI.

## Non-goals

- Do not build a native-code debugger from scratch.
- Do not invent a proprietary source map format as the primary contract.
- Do not make VS Code the only frontend.
- Do not promise reliable optimized variable inspection before explicit
  optimized-debug validation exists.
- Do not let debug mode change Osprey language semantics.

## Architecture

### Layer 1: compiler source model

The AST must carry complete source spans, not just optional declaration
positions.

Required model:

- `SourceId`: canonical source file identity.
- `SourceSpan`: file, start line/column, end line/column, byte offsets.
- Every `Stmt`, every `Expr`, every pattern, every type expression, every
  parameter, every field, every handler arm, and every generated wrapper must
  either have a real source span or be marked synthetic.
- Lowering rewrites such as pipe desugaring, method-call normalization,
  lambda lifting, generic specialization, effect handler lowering, and closure
  construction must preserve provenance.
- Synthetic compiler/runtime code must be marked as artificial and skipped by
  source stepping unless the user disables smart stepping.

Acceptance:

- Parser tests assert spans for representative declarations, statements,
  nested block expressions, lambdas, match arms, handler arms, and interpolation
  holes.
- LSP and debugger share the same line/column convention: Osprey AST uses
  1-based lines and 0-based columns; DAP uses 1-based lines and columns.
  Emitted `!DILocation` columns are 1-based, so convert the 0-based AST column
  with `column + 1` — LLVM reserves column `0` as the "no column" sentinel, and
  emitting a raw 0-based column silently collides with it.

### Layer 2: LLVM/DWARF debug info

The backend must emit LLVM debug metadata for native debuggers.

Required module metadata:

- `source_filename` with the canonical source path.
- `!llvm.dbg.cu` with one `DICompileUnit`.
- `!llvm.module.flags` with `"Debug Info Version"` and DWARF version.
- `DIFile` for every source file involved.
- Osprey producer string and language tag. `DW_LANG_C` is the pragmatic interim
  code, but it is **not neutral**: LLDB/GDB drive expression-evaluation language,
  value/type formatting, name demangling, and default array lower bounds off
  `DW_AT_language`, so claiming C imports C semantics that the Osprey-aware
  evaluator (Layer 6) must then override. `DW_LANG_lo_user`..`hi_user` is honest
  but unrecognized by every debugger (and ≥ `0x8000` forces `DW_FORM_data2`).
  The correct long-term path is to register `DW_LNAME_Osprey` with dwarfstd.org
  and emit the DWARF 6 `DW_AT_language_name` + `DW_AT_language_version` pair
  (permitted from a DWARF 5 producer) while dual-emitting legacy `DW_AT_language`
  for older consumers — the same path Rust and Swift took rather than reusing C.

Required function metadata:

- `DISubprogram` for user functions and `main`.
- Distinct linkage names for lowered symbols, but display names should be
  Osprey names.
- `DISubroutineType` for function signatures.
- Lexical scopes for blocks, lambdas, handler arms, match arms, and generated
  closures.
- Artificial/generated functions marked artificial where LLVM/DWARF supports
  it, so stepping favors source code.

Required instruction metadata:

- Every executable instruction derived from user code must carry `!dbg`.
- Statement-frontier lines receive stable locations for breakpoints.
- Pure bookkeeping instructions may inherit the nearest source location only
  when needed to keep breakpoints/stepping coherent; otherwise they are
  artificial/no-location.
- Tail code generated for `return`, block completion, and implicit `Unit`
  should map to the last user expression, not to random generated lines.

Required variable metadata:

- Emit variable locations as LLVM value-location records. Prefer
  `#dbg_declare` / `#dbg_value` **debug records** on LLVM 19+ / RemoveDIs
  toolchains. The current Osprey textual backend is allowed to use the older
  `@llvm.dbg.value` compatibility intrinsic while tested LLDB/DWARF output
  remains correct; this is a bridge, not the desired final IR spelling.
- Function parameters: `DILocalVariable` plus a value-location record or
  compatibility `dbg.value`.
- Immutable locals: value-location records for SSA values; `dbg.declare`/
  `#dbg_declare` for stack/heap slots when addressable.
- Mutable locals: a user variable and its storage cell must be represented so
  the user sees the value, not only the cell pointer.
- Captures: closure environment fields mapped back to capture names.
- Match bindings and handler parameters: lexical-block scoped locals.
- Optimized-away values: reported honestly as unavailable; never fabricated.

Required type metadata:

- Primitive Osprey types mapped to base types.
- Strings as Osprey strings, not raw `i8*` by default.
- Records/objects as composite types with named fields.
- Unions as tagged variants.
- Result as tagged success/error shape with value/message rendering.
- Lists/maps/channels/fibers as runtime handle types with pretty-printer
  expansion, because their concrete C layout is runtime-owned.
- Function values/closures as callable handles with captures.
- Ptr/FFI types as opaque pointers unless an extern binding supplies richer
  metadata later.

Acceptance:

- `osprey --llvm --debug file.osp` contains `DICompileUnit`,
  `DISubprogram`, and `DILocation`.
- Compiled binaries contain a source line table visible to `llvm-dwarfdump` or
  platform equivalent.
- LLDB can set a breakpoint by `.osp` file and line, then stop there.
- Stepping through a simple program visits Osprey source lines in expected
  order.

### Layer 3: compiler CLI and build modes

New CLI surface:

- `--debug`: emit debug info and build for debugging.
- `--debug-info=dwarf|none`: explicit debug-info choice, default `dwarf` when
  `--debug` is present.
- `--debug-opt=none|limited|optimized`: optimization/debuggability policy.
  Initial implementation supports `none` only.
- `--debug-out <path>`: optional path for the debug executable.
- `--debug-preserve-ir`: keep the `.ll` file for inspection.
- `--debug-preserve-symbols`: keep platform symbol bundles such as `.dSYM`.

Debug build behavior:

- `--debug` implies `-O0` unless `--debug-opt` requests otherwise.
- Pass `-g` to the C compiler driver.
- Prefer `-fno-omit-frame-pointer` where supported.
- On macOS, run or document `dsymutil` behavior so line tables and symbols are
  discoverable by LLDB.
- Do not strip debug symbols.
- Runtime archives should have debug variants once runtime stepping is needed.
  First slice may debug user code only.

Acceptance:

- `osprey file.osp --debug --compile` writes a binary and leaves enough debug
  info for LLDB to resolve Osprey lines.
- `osprey file.osp --debug --run` executes normally.
- Existing non-debug `--compile`/`--run` behavior and performance defaults are
  unchanged.

### Layer 4: DAP and editor integration

Initial adapter strategy:

- Reuse LLDB-DAP as the real adapter.
- The Osprey VS Code extension resolves the Osprey compiler, compiles the
  active `.osp` file with `--debug --compile`, then launches `lldb-dap` against
  the produced binary.
- The extension contributes a real `osprey` debugger configuration, not a
  provider that cancels the session and runs the file.

Debug configuration fields:

- `program`: `.osp` source path.
- `args`: program args.
- `cwd`: working directory.
- `env`: environment map.
- `stopOnEntry`: bool.
- `compilerPath`: override compiler path.
- `lldbDapPath`: override LLDB-DAP path.
- `debugOutput`: optional binary path.
- `preserveArtifacts`: bool.
- `console`: integrated terminal, external terminal, or internal console.
- `preLaunchTask` support through VS Code's native mechanisms.

Resolution:

- Find `lldb-dap` from explicit config, extension setting, common LLVM toolchain
  paths, or PATH.
- If missing, show a precise error explaining the required tool and the setting
  to configure.
- On systems where only `lldb-vscode` exists, detect it as a compatibility
  fallback if it supports DAP.

Future Osprey-aware adapter:

- A small `osprey-dap` shim may sit in front of LLDB-DAP when native DAP cannot
  express Osprey semantics.
- It should proxy standard requests and only intercept Osprey-specific
  variables, evaluate, fibers, effects, and replay requests.
- Any non-standard DAP extension must be prefixed and documented.

Shared editor components (no per-language duplication):

- Osprey, Basilisk, and SharpLsp each ship a VS Code debugger. The editor-side
  glue is language-neutral and MUST be shared via a package under the
  [LspKit](https://github.com/Nimblesite/lspkit) umbrella, not forked into each
  extension. See spec `[DEBUGGER-REUSE]`.
- LspKit is currently Rust-only and carries no debugger code, so this needs a
  new shared TypeScript package (VS Code DAP glue + DAP test harness) alongside
  a future `lspkit-debug` Rust crate for native debug-build policy.
- Shared TypeScript surface: adapter resolution (setting → toolchain paths →
  PATH + precise missing-tool error; cf. SharpLsp `findNetcoredbg`), debug-config
  synthesis (cf. Basilisk `applyDebugConfigDefaults`, SharpLsp
  `resolveDebugConfiguration`), save-dirty-or-reject, the pre-launch build hook,
  and the DAP test client/poll/UI-stub harness.
- Language-specific only: debug `type`, adapter binary (`lldb-dap`), build
  command, toolchain paths. The Osprey extension contributes its generic glue
  upstream rather than maintaining a private copy.

Acceptance:

- Pressing F5 on a focused `.osp` file starts a debug session.
- The editor DAP glue and DAP test harness are imported from the shared LspKit
  package, not duplicated in the Osprey/Basilisk/SharpLsp extensions.
- The VS Code debug UI can set/hit a source breakpoint.
- The debug session does not use `osprey.run`; it launches a debugger.
- A headless DAP integration test can initialize, launch, set breakpoints,
  continue, and observe a stopped event.

### Layer 5: Osprey value model and pretty-printers

Generic LLDB display will expose raw pointers and structs. A modern Osprey
debugger must render language values.

Required value renderers:

- `int`, `float`, `bool`, `string`, `Unit`.
- Records and anonymous objects with field names and nested values.
- Union variants with tag and payload fields.
- `Result<T, E>` as `Success(value)` or `Error(message/value)`.
- Lists with length and paged children.
- Maps with length and paged key/value entries.
- Strings with UTF-8 validation and byte/codepoint display policy matching the
  language spec.
- Closures with function name and captures.
- Fibers with state, result type, and join/await status.
- Channels with open/closed state and buffer/queue info if available.
- Effect handlers with active operation, captured env, and continuation state.
- FFI pointers with address and optional extern binding name.

Implementation options:

- LLDB Python synthetic providers and summaries for native LLDB.
- DAP variableReference expansion through an Osprey shim.
- Runtime helper functions for safe inspection of opaque handles. Helpers must
  be side-effect-free and safe when the program is paused.

Acceptance:

- Variables pane shows Osprey names and values for primitive locals.
- Records/unions expand by field/variant name.
- Lists/maps are paged and do not dump unbounded memory.
- Runtime handles never crash the debug session when null/corrupt; they render
  as invalid/unavailable.

### Layer 6: expression evaluation

Expression evaluation must be Osprey-aware. LLDB C expression evaluation is not
enough for Osprey syntax or type rules.

Required modes:

- `hover`: display variable value for symbols in scope.
- `watch`: evaluate pure Osprey expressions.
- `repl`: evaluate expressions and optional debugger commands.
- `call`: explicitly allow side-effecting calls only when the user requests it.

Implementation:

- Parse the expression with `osprey-syntax`.
- Resolve names from the selected frame's debug scope.
- Type-check against the frame environment.
- Lower either to an LLDB expression using known native symbols or to a small
  JIT/evaluation thunk when safe.
- Block or prompt for side-effecting operations, FFI, IO, spawning, HTTP, and
  effect operations unless explicitly enabled.

Acceptance:

- Evaluate local arithmetic/string expressions.
- Evaluate field access and match-safe value inspection.
- Reject unavailable optimized values with a clear message.
- Reject side effects by default.

### Layer 7: stack traces and source navigation

Stack traces must hide implementation noise by default but preserve access for
advanced debugging.

Required behavior:

- Frames display Osprey function/module/lambda names.
- Closure/lambda frames show declaration line and capture summary.
- Handler/resume frames show effect and operation.
- Runtime frames are hidden under smart stepping by default.
- Users can disable smart stepping to step into generated/runtime code.
- Stack traces from crashes include file/line/function.

Acceptance:

- `print(backtrace)` equivalent runtime support or LLDB backtrace shows Osprey
  file/line for user code.
- Stepping into `print` or runtime helpers is skipped by default.
- A crash in user code reports the Osprey call stack.

### Layer 8: fibers, channels, and effects

Osprey's runtime has logical concurrency that native OS-thread debuggers do not
fully model.

Required fiber support:

- Each fiber has a stable debug id.
- Debugger can list fibers with state: runnable, running, waiting, completed,
  failed.
- Stack for suspended fibers is available when runtime representation allows it.
- Breakpoints can stop all fibers or only the triggering fiber.
- Scheduler stepping can advance one fiber or run until next breakpoint.

Required channel/select support:

- Inspect channel state.
- Show waiting senders/receivers if runtime tracks them.
- Break on send/receive/select operations.

Required effect support:

- Show active handler stack.
- Show performed effect operation and arguments.
- Step into handler arm.
- Step over `resume`.
- Break on unhandled effect or invalid resume.

Acceptance:

- Debugger lists fibers in a concurrent test program.
- Breakpoint inside spawned work stops with a meaningful logical fiber id.
- Effect handler frames are visible and named.

### Layer 9: deterministic replay and time travel

Replay is not the first debugger, but a modern debugger plan must include it
because concurrency/effects bugs are otherwise hard to reproduce.

Required record stream:

- Program args, cwd, env, Osprey version, compiler version, target triple.
- External inputs: stdin, files when sandbox/record mode permits, network when
  enabled, clock/random/process results.
- Scheduler decisions for fibers/channels/select.
- Effect operation/resume ordering.
- FFI boundary calls when recordable; opaque calls marked non-replayable.
- Runtime allocation ids for stable handle identity where needed.

Required replay behavior:

- `osprey file.osp --debug --record trace.ospreytrace --run`.
- `osprey debug replay trace.ospreytrace`.
- Reverse-continue and reverse-step once snapshots/checkpoints exist.
- Deterministic failure if the run depends on non-recorded external state.

Acceptance:

- A fiber scheduling test replays exactly.
- A trace records metadata sufficient to reject replay on incompatible binary or
  compiler version.

### Layer 10: optimized debugging

Optimized debugging must be explicit and tested. It cannot be hand-waved.

Policy:

- `--debug-opt=none` is the default and the first supported mode.
- `--debug-opt=limited` may enable optimizations proven not to break line/var
  expectations for covered cases.
- `--debug-opt=optimized` is allowed only after validation infrastructure is in
  place; it must be honest about unavailable variables.

Required validation:

- Per-variable history oracle (Kell & Stinnett): each local's observed
  value-change history must match the source program's, with each change at the
  correct file and line. This factors out temporal imprecision and flags only
  outright wrong or missing variable info; it is the automatable correctness
  test for emitted local/line metadata.
- Control-flow conformance for location updates (Huang et al.): for every
  instruction a pass moves/clones/merges, assert the optimized path's location
  set is a subset of the source path's, and that the Preserve/Merge/Drop choice
  follows from that containment. Applies once Osprey runs its own IR passes.
- Source-level dynamic call-tree differential (Stinnett & Kell): compare the
  call tree recovered from an optimized run against its unoptimized reference;
  attribute divergences to compiler vs debug-info format. Do not settle for the
  weak backtrace invariant (reversed/empty backtraces satisfy it).
- Cross-debugger differential testing (DTD, Lu et al.): diff full program state
  between two adapters over the same binary, treating optimized-out / error /
  not-found as distinct states.
- Track coverage of debug metadata: functions, lines, variables, scopes, and
  value locations.

Acceptance:

- Optimized debug mode has its own test suite and documented limitations.
- Any optimized-away variable is reported as unavailable, not stale.

### Layer 11: security and safety

Debugger integration runs programs and can inspect memory; it needs explicit
boundaries.

Requirements:

- Never launch a program without user action in the IDE.
- Preserve workspace trust expectations.
- Quote paths safely and avoid shell execution where `execFile`/spawn APIs are
  available.
- Do not expose remote attach by default.
- Redact sensitive environment values from logs.
- Keep debug adapter logs opt-in.
- Side-effecting expression evaluation is off by default.
- FFI and network replay require explicit user opt-in.
- Debug artifacts should not be silently published in release packages unless
  intended.

Acceptance:

- All compiler/debugger subprocess launches use argument arrays, not shell
  interpolation.
- Missing tool errors do not leak env secrets.

### Layer 12: packaging and cross-platform

Platforms:

- macOS: LLDB/LLDB-DAP first, DWARF/dSYM support.
- Linux: LLDB-DAP first, GDB/DAP fallback only if needed.
- Windows: MinGW DWARF via LLDB/GDB first; native MSVC CodeView/PDB is a later
  track if the toolchain moves that way.

Packaging:

- The VS Code extension should not bundle a large LLDB toolchain initially.
- It should detect and explain missing LLDB-DAP.
- Release notes must state debugger prerequisites.
- Future binary releases may bundle platform-specific `lldb-dap` if licensing,
  size, and update cadence are acceptable.

Acceptance:

- Debug prerequisites documented for macOS/Linux/Windows.
- Extension fails cleanly when prerequisites are missing.

## Phased implementation

### Phase 0 - Spec and safety rails

- [x] Write this plan.
- [ ] Add `docs/specs/0021-Debugger.md` once phase 1 behavior is real enough to
      become language/tooling contract.
- [ ] Add issue checklist and labels for debugger phases.

### Phase 1 - Source line debugging

- [ ] Add complete statement source positions, including bare expression
      statements.
- [ ] Add `osprey_codegen::compile_program_debug`.
- [ ] Emit `DICompileUnit`, `DIFile`, `DISubprogram`, `DILocation`, and module
      flags.
- [ ] Add `--debug` and debug clang flags.
- [ ] Preserve debug artifacts for tests.
- [ ] Add golden tests for debug IR metadata.
- [ ] Add a headless LLDB smoke test for source breakpoint and stepping.

### Phase 2 - VS Code real debug launch

- [ ] Replace the current fake debug provider that cancels the session and runs
      the file.
- [ ] Compile the selected `.osp` file with `--debug --compile`.
- [ ] Resolve `lldb-dap`.
- [ ] Launch LLDB-DAP with the compiled binary.
- [ ] Add extension tests for configuration synthesis and missing-tool errors.
- [ ] Add a DAP smoke test that launches and hits a source breakpoint.

### Phase 3 - Variables and lexical scopes

- [ ] Emit `DILocalVariable` for parameters and lets.
- [ ] Emit value-location records for locals, mut cells, captures, match
      bindings, and handler parameters; migrate the textual spelling from
      `@llvm.dbg.*` compatibility intrinsics to `#dbg_*` debug records once the
      supported LLVM floor makes that path portable.
- [ ] Add lexical scopes for blocks, lambdas, match arms, and handlers.
- [ ] Validate primitive local inspection in LLDB.
- [ ] Add unavailable-value reporting tests.

### Phase 4 - Osprey type and value rendering

- [ ] Emit type metadata for primitives, records, unions, Result, and closures.
- [ ] Add LLDB summaries/synthetic providers or an Osprey DAP shim.
- [ ] Add runtime inspection helpers for list/map/string/fiber/channel/effect
      handles.
- [ ] Page large collections.
- [ ] Add corruption/null safety tests for runtime handle display.

### Phase 5 - Expression evaluation

- [ ] Parse/type-check watch expressions in Osprey syntax.
- [ ] Resolve names against current frame scopes.
- [ ] Implement pure expression evaluation.
- [ ] Gate side effects behind explicit user action.
- [ ] Add REPL/watch/hover tests.

### Phase 6 - Fibers, channels, effects

- [ ] Add stable runtime debug ids for fibers and channels.
- [ ] Expose runtime debug inspection APIs.
- [ ] Add fiber list and logical stack display.
- [ ] Add breakpoints on spawn/await/send/recv/select.
- [ ] Expose effect handler stack and resume events.
- [ ] Add concurrent/effects debugger tests.

### Phase 7 - Replay/time travel

- [ ] Define `.ospreytrace`.
- [ ] Record scheduler, IO, time/random, process, network, and effect events.
- [ ] Implement deterministic replay runner.
- [ ] Add checkpointing for reverse operations.
- [ ] Add replay tests for fibers/effects.

### Phase 8 - Optimized debug mode

- [ ] Implement debug metadata coverage metrics.
- [ ] Add unoptimized-vs-optimized stepping comparisons.
- [ ] Add control-flow conformance validation for debug locations.
- [ ] Add dynamic call-tree validation for representative programs.
- [ ] Enable `--debug-opt=limited`, then `--debug-opt=optimized` only after
      tests define the contract.

## Testing matrix

Compiler tests:

- Parse span tests.
- Debug IR golden tests.
- Type metadata tests.
- Variable metadata tests.
- Non-debug IR unchanged tests.

Native debugger tests:

- LLDB command-script tests.
- `llvm-dwarfdump` line/function/variable checks.
- macOS dSYM checks.
- Linux ELF DWARF checks.
- Windows MinGW DWARF checks when CI is available.

DAP tests:

- Initialize/launch/configurationDone.
- SetBreakpoints.
- Continue/stopped.
- StackTrace.
- Scopes/variables.
- Evaluate.
- Disconnect/terminate.

Osprey semantic tests:

- Functions, top-level statements, blocks.
- Lambdas and closures.
- Records/unions/Result/list/map.
- Match and pattern bindings.
- Fibers/channels/select.
- Effects/handlers/resume.
- FFI safe/opaque boundary.

Regression tests:

- Breakpoint on every executable source line in selected examples.
- Step sequence is stable for `--debug-opt=none`.
- Variables are correct or explicitly unavailable.
- Runtime/generated frames hidden by default.

## Risks

- LLVM textual debug metadata is easy to get syntactically valid but semantically
  weak. Mitigation: validate with LLDB and `llvm-dwarfdump`, not just string
  checks.
- Osprey currently optimizes at the clang invocation level. Debug mode must
  override that without changing release behavior.
- Runtime handles are opaque. Pretty-printers need stable runtime inspection
  APIs, not ad hoc memory guessing.
- Fibers/effects may require runtime changes to expose logical stacks.
- macOS debug symbols may require dSYM handling even when ELF-like tests pass
  elsewhere.
- Windows native debug format support will lag unless the project adopts an
  MSVC/CodeView path.

## Definition of done

The Osprey debugger is complete when:

- A user can debug normal Osprey programs from VS Code and from a headless DAP
  test runner.
- Breakpoints, stepping, stacks, scopes, variables, and expression evaluation
  work for the language constructs in `examples/tested`.
- Osprey runtime values render as language values, not raw implementation
  pointers.
- Fibers and effects are inspectable at the language level.
- Replay can reproduce deterministic concurrency/effects executions.
- Debug-info quality is continuously tested, including optimized-debug modes
  once those modes are enabled.
- `make ci` includes the non-GUI debugger checks; GUI/VS Code checks remain in
  the extension test lane.
