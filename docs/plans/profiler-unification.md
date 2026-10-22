# Shared profiler integration

Preserve Osprey's runtime and compiler layer. Move reusable views and behavior into lspkit; mount the same views and exercise the same VSIX contract from both extensions. No git operations are part of this work.

| Capability | Osprey today | lspkit / SharpLsp today | Required destination |
| --- | --- | --- | --- |
| Optimized native builds, DWARF, platform sampler | Implemented | Runtime-specific .NET sampler | Keep Osprey implementation |
| Fiber IDs, running/waiting samples, dropped samples | Implemented | No fiber-state view | Shared UI accepts runtime metadata |
| Left-heavy/time-order canvas, zoom/pan/search | Implemented | SVG preview/external interactive viewer | Move engine into lspkit |
| Hot self/total functions, line heat, source jumps | Implemented | Implemented | Shared results, preserve Osprey annotations |
| Speedscope, V8, folded exports | Implemented | Implemented | Preserve artifacts |
| Cancellable compilation/run, test explorer profiling | Implemented | Launch/attach session flow | Keep Osprey commands and test integration |
| Process explorer, live sessions | Absent | Implemented | Runtime adapter required |
| Heap snapshots/diffs/leaks/retention | Absent | .NET backend implemented | Native heap capture required; never fake support |
| Configurable product palettes | Absent | Fixed orange/purple | Shared palette with host defaults |
| Watch ↔ profiler navigation | Absent | Type picker, no watch bridge | Shared bridge, exact session/type provenance |
| VSIX regression contract | Osprey-specific suites | Shared UI + .NET integration suites | Common suite runs in both real extensions |

## Verification checklist

- [x] Existing Osprey model, HTML, panel, heat, cancellation and export tests pass.
- [x] Both hosts use the same interactive CPU renderer and test contract.
- [x] Product palettes cover canvas, panels and editor heat.
- [x] Debugger links preserve session identity and do not guess object identity from a type.
- [x] Real Osprey runs reach the shared profiler; .NET regression runs still pass.
- [x] Native runtime capabilities and gaps are explicitly verified.
- [x] Desktop, 520px/360px panels, light/dark themes and live canvas palette repaint checked.
- [x] Reduced-motion behavior checked in a rendered panel.

Verification on 2026-10-22: all 438 Osprey extension tests passed before the final refinements; the final focused Osprey profiler run passed 91 tests, SharpLsp passed 51 shared profiler/CPU/decoration tests, and the final broader Osprey profiler/debugger run passed 163 tests including the real variable-to-profiler bridge. Both suites import the same lspkit webview and CPU capture contracts. These run in real VS Code extension hosts; they do not establish installed-package validation or native heap/attach parity.

## Workspace integration

The unpublished shared-package changes live in the existing `../SharpLsp/lspkit` checkout. Osprey currently resolves that local source package. A coordinated lspkit release or repository dependency pin is required before this can build from an isolated Osprey checkout; no git or publication was authorized in this task. Do not describe this as shipping-ready or full backend parity.

## Outstanding runtime work

- [ ] Implement and validate native process discovery and live CPU attach/snapshot/stop through the shared protocol.
- [ ] Implement native heap snapshots, allocation sites, diffs, leak evidence and retention graphs for supported Osprey memory backends.
- [ ] Run the common process/memory/debugger runtime contracts in Osprey, not only the UI contracts.
- [ ] Resolve existing compiler CI failures: missing Default twins for `ml_name_positions` and `verdict_declared`, and ML parse errors in `finalizers_masks_instances`.
