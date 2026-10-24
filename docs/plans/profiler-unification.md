# Shared profiler integration

Osprey keeps its runtime and compiler layer: optimized native builds, DWARF symbolization, the platform sampler, the live-attach endpoint every native program serves ([PROF-ATTACH-ENDPOINT] in [0028-Profiler](../specs/tooling/0028-Profiler.md)), compiler invocation, cancellation and source classification (`isOspreySource`). lspkit owns everything above that: the `lsp-profiler` service mounted by `osprey lsp` through `crates/lsp-profiler-osprey`, and the `@nimblesite/lsp-profiler-ui` package whose views, palette, debugger links and VSIX contracts both Osprey and SharpLsp mount. The behaviour is specified in lspkit's [LSP-PROFILER-SPEC](../../lspkit/docs/specs/LSP-PROFILER-SPEC.md): `[PROF-UI-SAMPLED]`, `[PROF-UI-PALETTE]`, `[PROF-UI-DEBUG-LINK]`, `[PROF-UI-CONTRACT]` and `[PROF-UI-CAPABILITIES]`.

| Capability | Where it lives | Status |
| --- | --- | --- |
| Optimized native builds, DWARF, platform sampler | Osprey compiler and runtime | Done |
| Fiber IDs, running/waiting samples, dropped samples | Osprey summary, shared header and chips | Done |
| Left-heavy/time-order canvas, zoom/pan/search, slab budget | lspkit `flame-model`, `flame-time-order`, `flame-script` | Done |
| Hot self/total functions, line heat, source jumps | lspkit results panel and decorations; Osprey heat annotations | Done |
| Speedscope, V8, folded exports | Osprey CLI artifacts, opened from the shared panel | Done |
| Cancellable compilation/run, test explorer profiling | Osprey commands and test integration | Done |
| Process discovery, live CPU attach, snapshot, stop, live progress | `crates/lsp-profiler-osprey` over the shared protocol | Done on Unix; Windows has no endpoint (`attach_none.rs`) |
| Heap snapshots, diffs, leaks, retention | No Osprey memory backend | Absent; `memory: false` hides every memory surface, never faked |
| Host palettes over canvas, panels, memory canvases and editor heat | lspkit `palette.ts`; Osprey `OSPREY_PALETTE` | Done |
| Watch ↔ profiler navigation | lspkit `debug-link`; session-scoped, type queries only | Done; type retention needs a memory backend |
| VSIX regression contract | lspkit `testing/webview-contract`, `testing/cpu-capture-contract` | Runs in both real extension hosts |

## Workspace integration

Osprey builds the shared UI from its `lspkit` submodule (`file:../lspkit/packages/lsp-profiler-ui`; CI runs `make _profiler_ui` before `npm ci`, which installs a packed copy). The submodule pins a commit on lspkit's `feature/shared-profiler-ui` branch, which is pushed but not merged into lspkit `main`, and has no pull request. SharpLsp's own `lspkit` submodule pins the older `4566b47` on the same branch.

## Review of 2026-10-24

Fixed, each with a test written first where the defect was observable:

- [x] A speedscope export with a non-object profile threw from `parseSpeedscope`, and a `null`, negative or string weight was accepted and turned every slab into `NaN`. Both are now classified errors.
- [x] Search's "% of total" compared every matched slab with every counted one, O(n²) over a 150,000-slab budget. It is now one sort and sweep (`coveredSpan`).
- [x] A label cached for an 8px width bucket was fitted to the widest width in the bucket and could overflow its slab by up to 7px. It is now fitted to the narrowest.
- [x] The frame ranking was implemented twice (`frameColors` and `paletteFrames`), the search twice (`searchMatches` and the webview's own loop), the palette's CSS variable names twice (`paletteCss` and the live repaint) and the palette's default colours twice (`DEFAULT_PALETTE` and `PROFILER_CSS_VARS`). Each now has one definition, and the webview runs the tested functions through `Function.prototype.toString`.
- [x] The memory timeline, the retention graph's target and cycle colours, the leak badges and the high-confidence leak decoration were fixed purple and red whatever the host palette said. They now read the palette.
- [x] No test pinned the live progress's top function or the hot-function tie order. The Rust tests `progress_names_the_function_with_the_most_self_samples`, `progress_breaks_equal_self_samples_by_name` and `hot_functions_break_equal_totals_by_self_samples_then_name` do now.
- [x] The UI package's `npm run lint` ran ESLint with no config and always failed. It now type-checks with `tsc --noEmit`, as Osprey's extension does.
- [x] Osprey's `flame-html`, `flame-script` and `flame-script-helpers` were one-line re-exports carrying a comment that did not describe them. They are deleted, and their importers use the package directly. `OSPREY_PALETTE` is typed as `ProfilerPalette`.
- [x] Osprey's flame panel re-implemented lspkit's `SingletonWebviewPanel` and its message routing, with a dead `navigateToSource` rewrite and no `openFlamegraphSvg` route. Its untested copies dropped the extension's branch coverage under the 95% gate. The panel is now a `SingletonWebviewPanel` routed through lspkit's exported `handleFlamegraphMessage`, titled after the profiled program.
- [x] A posted message that was not an object with a string `type` (or had a non-string `file` or non-numeric `line`) reached every profiler panel's handler, which threw on `msg.type`. The panel boundary now drops it (`webviewMessage`), pinned by a webview-contract test that runs in every host.
- [x] Stale comments: fixed hex colours in the decoration headers, an orphaned `MAX_INLINE_SVG_BYTES` comment, the `tracksMemory` contract, and imports placed mid-file.

Verification: lspkit `make lint` (clippy, spec citations, skill links), `cargo fmt --check`, the three new Rust tests and the UI package's 172 VS Code tests pass. Osprey's extension lints and compiles against the rebuilt package. Osprey's profiler, flame, debugger and palette suites (`OSPREY_TEST_GREP="[Pp]rofil|[Ff]lame|[Dd]ebug|[Pp]alette|[Ss]ampled"`) passed 269 of 272 in one seven-minute run. The three failures all passed when re-run on their own: two `runInTerminal` launcher timeouts in the Test Explorer debug suite, and `a debug-launched program is profiled over the wire`, whose process listing timed out. After the panel was moved onto the shared code, the whole extension suite (518 tests, plus 166 against the installed VSIX) passed with branch coverage at 95.75%, over the 95% gate that the copies had dropped it under.

## Todo

Delivery:

- [ ] Open the lspkit pull request for `feature/shared-profiler-ui`, merge it into `main`, and repin Osprey's submodule to the merged commit.
- [ ] Bump SharpLsp's `lspkit` submodule to the same commit and run its shared profiler, CPU and decoration suites. SharpLsp has not been run against these fixes.
- [ ] Publish `@nimblesite/lsp-profiler-ui` so neither host depends on a `file:` path (lspkit plan, Phase 4).

Runtime:

- [ ] Windows live attach: `attach_none.rs` leaves Osprey programs on Windows undiscoverable and unprofilable from the Processes view.
- [ ] A native heap backend for Osprey's memory modes (ARC and GC): snapshots, allocation sites, diffs, leak evidence and retention graphs. Then drop `memory: false` from `PROFILER_HOST` and run the shared memory contracts in Osprey.
- [ ] Debugger type retention for Osprey variables, which needs the heap backend.

Shared UI quality debt, which predates this work unless marked:

- [ ] `profiler.ts` is 542 lines, over the 500-line limit.
- [ ] Functions over 20 lines in `profiler.ts`, `process-explorer.ts`, `process-explorer-rows.ts`, `memory-dashboard-html.ts`, `memory-status.ts`, `profiler-decorations.ts`, `profiler-flamegraph-html.ts`, `profiler-webview.ts` and `scripts/contributions.mjs`. `createProfilerStore` (28 lines) and `registerProfilerUi` (25 lines) are shorter than before this work but still over.
- [ ] Canvas labels (`rgba(15, 15, 15, 0.9)`) and the active toggle's text (`#171b16`) assume a light palette. A host with dark CPU colours gets unreadable labels.
- [ ] The UI package's own suite runs against a mock language server, which lspkit's `CLAUDE.md` forbids. Real-server coverage is the host suites. Port the suite or record the exception in the spec.
- [ ] An ESLint config with SharpLsp's reactive-screen and no-legacy rules (lspkit plan, Phase 4).

Test reliability:

- [x] The live profiler e2e tests recorded for fixed seconds, so a starved sampler failed them: the snapshot found no samples, and the progress test stopped before any progress was sent. They now record until `cpu.list` shows the samples they need, the progress test waits for a progress notification covering them, and the duration test waits for its session to leave the listing. Each wait has a two-minute deadline. All 13 passed three runs in a row with 28 busy loops on 14 cores.
- [ ] `profiler-cpu-launch-e2e` (`waitForListing`) and `test-debug.e2e` (integrated-terminal console selection) fail under a long, loaded run and pass alone. Find what they race on and remove the timing dependence. Do not raise the timeouts.

Compiler CI on PR #270, owned by the effects work rather than the profiler:

- [ ] `tests/effects/scope/finalizers_masks_instances` needs spec 0017's `mask` and `instance`.
