# Plan: Colored TUI That Calls Ad-Hoc APIs

Parent: [`production-primitives.md`](production-primitives.md).

Related specs:
- [`0014-HTTP.md`](../../compiler/spec/0014-HTTP.md) — current HTTP client surface
- [`0012-Built-InFunctions.md`](../../compiler/spec/0012-Built-InFunctions.md) — builtin registry
- [`0017-AlgebraicEffects.md`](../../compiler/spec/0017-AlgebraicEffects.md) — effect handler shape used by the TUI abstraction

## Why this exists

The pitch: *a developer can sit down and write an amazing colored TUI in Osprey that calls arbitrary HTTP APIs and renders the responses.* Today that pitch is false. This plan enumerates the exact gaps and lays out the runtime + builtin work to close them.

A survey of the current state (`runtime/`, `internal/codegen/`, `examples/tested/`, `osprey.g4`) found **five hard blockers**, listed below. The first one alone makes "call an API and show the result" impossible.

## Current state — what's stopping us

### Blocker 1 — HTTP returns only a status code, not the response body
[`http_client_runtime.c:4-38`](../../compiler/runtime/http_client_runtime.c#L4-L38). Every HTTP method (`httpGet`, `httpPost`, `httpPut`, `httpDelete`, `httpRequest`) returns `int` — the HTTP status. The response body is read into a local buffer and discarded. There is no `httpResponseBody` or equivalent. You can ping endpoints but you cannot read what they return.

### Blocker 2 — JSON parsing is a stub
[`system_runtime.c:380-417`](../../compiler/runtime/system_runtime.c#L380-L417). `extract_json_field` is a naive substring search: no nesting, no arrays, no escape sequences, no number-vs-string discrimination, no nulls. The file itself contains a `TODO: Implement proper JSON parsing` comment. Without a real parser, responses are unreadable even if Blocker 1 is fixed.

### Blocker 3 — No raw terminal / single-keystroke input
`input()` reads one whole line and blocks until newline. There is no `termios` raw mode, no escape-sequence reading, no arrow keys, no Ctrl+C trap. So: no menus driven by arrow keys, no live filters, no scrollable lists, no `q` to quit.

### Blocker 4 — No ANSI / terminal helpers
There is no `clearScreen`, `moveCursor`, `getTerminalSize`, `hideCursor`. You can embed `"\x1b[31m"` directly in `print()` strings, but with no `stringCellWidth()` (length of a string after stripping ANSI) you cannot lay out colored text reliably, and with no terminal-size query you cannot redraw on resize.

### Blocker 5 — Grammar has no loop construct
Verified in `osprey.g4` — no `for`, `while`, `until`, `loop`. A render loop has to be expressed as recursion (`fn loop() = { render(); handleInput(); loop() }`). This actually works because of TCO/fiber yield, but it's worth naming because every TUI tutorial in the world assumes `while true`. We do **not** propose adding a loop keyword here — recursion is fine — but the docs need to teach the pattern.

## What already works that we're building on

- ✅ String toolkit (`substring`, `split`, `padStart/End`, `${}` interpolation, `indexOf`, `trim`) — enough surface for a hand-written JSON parser once recursive-union payloads land.
- ✅ Fibers + channels — render fiber + input-reader fiber pattern is viable once raw stdin exists.
- ✅ Pattern matching on union types — perfect for `Event = KeyPress(s) | Tick | ApiResponse(r) | Quit`.
- ✅ Effects system — `effect Terminal { render, readKey, size }` + `handle ... in` is the right shape for testable TUI code.
- ✅ `print()` writes to stdout and ANSI escape sequences pass through unchanged.

## Dependencies on other plans

This plan **does not duplicate** the JSON-in-Osprey work tracked by [`production-primitives.md`](production-primitives.md). That plan's end state is a JSON parser *written in Osprey itself*, sitting on closures, recursive union payloads, list patterns, and the string cursor.

For the TUI plan we have two options for JSON:

- **Option A (recommended for v1):** ship a C-side JSON parser as a builtin (`jsonParse`, `jsonGet`, `jsonType`, `jsonKeys`, `jsonFree`) so the TUI demo can ship before `production-primitives.md` is fully done. ~200 LOC of C using a vetted single-file library (e.g. jsmn or cJSON).
- **Option B:** wait for the Osprey-native parser to land via `production-primitives.md` and use that. Slower, but no FFI burden and dogfoods the language.

The TUI plan **assumes Option A** for v1 and notes Option B as a follow-up cleanup. If `production-primitives.md` lands first, we delete the C JSON builtins and switch the demo over.

## Design decisions worth flagging before code is written

### Decision 1 — HTTP response shape

Two viable shapes for the response handle:

```
// Shape A — opaque handle, accessors
let r = httpGet(client, "/users/1")  // -> Result<HttpResponse, Error>
let body = httpResponseBody(r)
let status = httpResponseStatus(r)
httpResponseFree(r)
```

```
// Shape B — record
let r = httpGet(client, "/users/1")  // -> Result<{status: int, body: string, headers: Map<string, string>}, Error>
```

**Recommendation: Shape A** for v1. Records over FFI need `Map<string, string>` codegen for headers and force every caller to allocate the headers map even if they don't read it. Opaque handle + accessors lets the body stream into a Osprey-owned `string` on first call to `httpResponseBody(r)` and lazily build the headers map on first call to `httpResponseHeader(r, name)`. Shape B is a v2 ergonomic improvement once we have struct-return-over-FFI worked out.

### Decision 2 — TUI primitive vs widget library

**Recommendation: primitives only for v1.** Ship `termRawMode`, `termReadKey`, `termSize`, `termClear`, `termMoveCursor`, `termHideCursor`, `termShowCursor`, plus a pure-Osprey `ansi` module for colors. Build *one* widget (a colored arrow-key menu) inside the demo app. Decide on a widget library only after we see what the demo wants.

### Decision 3 — Sync vs async HTTP

**Recommendation: sync only for v1.** Wrap the HTTP call in a fiber if you want non-blocking (`spawn fn() => httpGet(...)`). A real `httpGetAsync` returning `Fiber<HttpResponse>` is Phase 6.

## Sequencing

The phases are ordered so each one ends with a runnable, testable artifact. **Phase 1 alone unblocks "call an API and show the body"** — that's the bare minimum.

### Phase 1 — Unblock API response bodies

Smallest, most contained change. Without this, nothing else matters.

- [ ] Refactor [`http_client_runtime.c`](../../compiler/runtime/http_client_runtime.c) so the response body and status are retained in a heap-allocated `HttpResponse` struct, not discarded.
- [ ] Change `httpGet`/`httpPost`/`httpPut`/`httpDelete`/`httpRequest` to return an opaque `int` handle (or `-1` on error) instead of the status code.
- [ ] Add `httpResponseStatus(handle: int) -> Result<int, Error>`.
- [ ] Add `httpResponseBody(handle: int) -> Result<string, Error>`.
- [ ] Add `httpResponseHeader(handle: int, name: string) -> Result<string, Error>` (returns first matching header, or error if absent).
- [ ] Add `httpResponseFree(handle: int) -> Result<bool, Error>`. Double-free must error, not segfault.
- [ ] Register the new builtins in `internal/codegen/` with correct `Result<T, Error>` type signatures.
- [ ] Spec update: amend [`0014-HTTP.md`](../../compiler/spec/0014-HTTP.md) to document the new response surface; mark the old "returns status code" behavior as removed.
- [ ] Migrate every existing [`examples/tested/http`](../../compiler/examples/tested/) example to the new API; delete any that only checked status.
- [ ] Add `examples/tested/http/http_body_roundtrip.osp` — `GET https://httpbin.org/json`, print the body, assert it contains `"slideshow"`. (Pick a stable test endpoint or use the test HTTP server already in the runtime tests.)
- [ ] Add C runtime test covering the handle lifecycle: alloc → status read → body read → free → double-free errors cleanly.

### Phase 2 — JSON parsing (C builtin, v1 shortcut)

- [ ] Add `runtime/json_runtime.c` wrapping a vetted single-file JSON library. Surface: `jsonParse(s: string) -> Result<int, Error>` (handle), `jsonType(h, path: string) -> string` (returns `"null"`/`"bool"`/`"number"`/`"string"`/`"array"`/`"object"`), `jsonGet(h, path) -> Result<string, Error>` (path syntax: `"a.b[0].c"`), `jsonKeys(h, path) -> Result<List<string>, Error>`, `jsonLength(h, path) -> Result<int, Error>` (for arrays), `jsonFree(h) -> Result<bool, Error>`.
- [ ] Replace the [`extract_json_field` stub](../../compiler/runtime/system_runtime.c#L380-L417) — either delete it or make it forward to `jsonParse`+`jsonGet`.
- [ ] Register the builtins in `internal/codegen/`.
- [ ] Spec update: amend [`0012-Built-InFunctions.md`](../../compiler/spec/0012-Built-InFunctions.md) with a `JSON` section. Mark it as "v1 builtin — slated for replacement by Osprey-native parser per `production-primitives.md`."
- [ ] Add `examples/tested/http/github_api.osp` — `GET https://api.github.com/users/<known-stable-user>`, parse, print `login` and `public_repos`. Use a recorded fixture if hitting the live API in CI is too flaky.
- [ ] Document the path syntax (`a.b[0].c`) in the spec, including escape rules for keys containing `.` or `[`.

### Phase 3 — Terminal raw mode and ANSI control

- [ ] Add `runtime/term_runtime.c`:
  - [ ] `termRawMode(on: bool) -> Result<bool, Error>` using `termios` (POSIX). Save original state on first enable. Register `atexit` handler that restores cooked mode + shows cursor — non-negotiable, every uncaught panic must restore the terminal.
  - [ ] `termSize() -> Result<{cols: int, rows: int}, Error>` via `ioctl(TIOCGWINSZ)`. If records-over-FFI is not yet supported, return two builtins `termCols()` / `termRows()` for v1.
  - [ ] `termReadKey() -> Result<string, Error>` — reads one keystroke in raw mode and translates ANSI escape sequences to human-readable strings (`"Up"`, `"Down"`, `"Left"`, `"Right"`, `"Enter"`, `"Esc"`, `"Tab"`, `"Backspace"`, `"Ctrl-C"`, otherwise the literal character). Returns error on EOF.
  - [ ] `termClear()`, `termMoveCursor(row: int, col: int)`, `termHideCursor()`, `termShowCursor()` — write the corresponding ANSI escape and flush.
- [ ] Register all of the above in `internal/codegen/`.
- [ ] Spec update: new section in [`0012-Built-InFunctions.md`](../../compiler/spec/0012-Built-InFunctions.md) called "Terminal Control", spec ID `[BUILTIN-TERM]`.
- [ ] Add pure-Osprey module `examples/lib/ansi.osp` (or wherever stdlib lives) with `red(s)`, `green(s)`, `blue(s)`, `yellow(s)`, `cyan(s)`, `magenta(s)`, `bold(s)`, `dim(s)`, `bgRed(s)`, … — all just string-wrapping helpers. No new builtins.
- [ ] Add `stringCellWidth(s: string) -> int` builtin that strips ANSI escape sequences (CSI / SGR) before counting cells. Needed so `padEnd(coloredString, 20)` aligns visibly. C implementation: state machine over the bytes, skip `\x1b[...m` runs.
- [ ] Add a runtime test for `termReadKey` driven by a pty (or skip in CI and have a manual test recipe documented).

### Phase 4 — Demo TUI app (the proof point)

This is the deliverable that says "yes, you can write a TUI in Osprey."

- [ ] Add `examples/tested/tui/github_repo_browser.osp`:
  - On launch, prompts (line-mode) for a GitHub username.
  - Fetches `https://api.github.com/users/<name>/repos` via Phase 1 HTTP, parses with Phase 2 JSON.
  - Enters raw mode via Phase 3, draws a colored menu of repos (name + stars + language).
  - Arrow-key navigation, Enter to view detail (description + URL), Esc/`q` to go back.
  - On exit, restores terminal cleanly even if the user hits Ctrl-C.
  - Use `effect Terminal { render, readKey, size }` so a fake terminal handler can drive the app in tests.
- [ ] Add a headless test that wires up a `handle Terminal in` block returning a scripted sequence of keys and asserts on the rendered output. This is what makes the TUI testable.
- [ ] Document the pattern in [`website/src/docs`](../../website/src/docs/) so users can copy it.
- [ ] Screenshot or asciinema recording in the website docs.

### Phase 5 — Nice to have (out of scope for v1, listed so we don't pretend they're free)

- [ ] Mouse support (`termMouseMode`).
- [ ] Bracketed paste.
- [ ] Widget library (`Menu`, `Table`, `TextInput`, `Spinner`).
- [ ] `httpGetAsync` returning `Fiber<HttpResponse>` — true non-blocking IO.
- [ ] Replace C JSON builtins with the Osprey-native parser from [`production-primitives.md`](production-primitives.md). When this lands, the C builtins from Phase 2 get deleted.
- [ ] Streaming HTTP response bodies (for large payloads or SSE) — currently the full body is buffered.

## Effort estimate

| Phase | Estimate |
|---|---|
| 1 (HTTP body) | half a day |
| 2 (JSON builtin) | one day |
| 3 (terminal control) | one to two days |
| 4 (demo + tests + docs) | one day |
| **Total to "credible TUI demo"** | **3–4 focused days** |

Phase 5 items are open-ended and not on the critical path.

## Master TODO

- [ ] Land Phase 1 — HTTP response handles + body/status/header accessors, spec update, example migration.
- [ ] Land Phase 2 — C JSON builtins, spec update, GitHub API example.
- [ ] Land Phase 3 — `termios` raw mode, ANSI helpers, `stringCellWidth`, spec update.
- [ ] Land Phase 4 — `github_repo_browser.osp` demo with effect-handler-driven tests + website docs.
- [ ] Decide whether to revisit Phase 5 items based on demo feedback.
- [ ] After [`production-primitives.md`](production-primitives.md) JSON parser lands: delete C JSON builtins and migrate the demo to the Osprey-native parser.
