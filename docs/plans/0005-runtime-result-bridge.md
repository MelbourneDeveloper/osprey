# Plan 0005 — HTTP / WebSocket `Result` Bridge Alignment

**Subsystem:** `compiler/runtime` (C), `crates/osprey-codegen` (extern call ABI),
`crates/osprey-types` (signatures)
**Status:** Partially implemented
**Spec:** [0014-HTTP.md](../specs/0014-HTTP.md), [0015-WebSockets.md](../specs/0015-WebSockets.md), [0013-ErrorHandling.md](../specs/0013-ErrorHandling.md)

## Summary

HTTP and WebSocket servers/clients work end-to-end in examples, but the C runtime
returns raw `int64_t` (error codes / handles) while the type system expects
`Result<T, string>`. The two specs say so explicitly in their `## Status`
sections — "the bridge is being aligned." This is one defect class shared by both
subsystems, so it is planned together to avoid duplicate work.

## Evidence

- HTTP status: *"The current C runtime returns raw `int64_t` for the
  create/listen/stop and request functions; the type system expects
  `Result<T, string>` … The handler bridge in `httpListen` currently expects a raw
  string return rather than the `HttpResponse` record."* —
  [0014-HTTP.md](../specs/0014-HTTP.md) §Status.
- WebSocket status: *"The current C runtime returns raw `int64_t` for several of
  these functions; the type system expects `Result<T, string>` … WebSocket server
  `listen` currently fails to bind in some environments."* —
  [0015-WebSockets.md](../specs/0015-WebSockets.md) §Status.
- Runtime: [compiler/runtime/http_server_runtime.c](../../compiler/runtime/http_server_runtime.c),
  [compiler/runtime/websocket_server_runtime.c](../../compiler/runtime/websocket_server_runtime.c),
  [compiler/runtime/http_shared.h](../../compiler/runtime/http_shared.h).
- ABI: [crates/osprey-codegen/src/extern_call.rs](../../crates/osprey-codegen/src/extern_call.rs).

## What works today

- `httpCreateServer` / `httpListen` / `httpStopServer`, client requests, and the
  `HttpResponse` record round-trip — see
  [compiler/examples/tested/http/](../../compiler/examples/tested/http/).
- WebSocket create/listen/send/close/broadcast paths exist in the runtime.

## Where it's misaligned

- Runtime functions return bare `int64_t` codes, so the `Result<T, string>`
  surface that the rest of the language relies on is synthesized loosely (or not
  at all) at the ABI boundary, and error *messages* are lost (only codes survive).
- `httpListen`'s handler bridge expects a raw `string` return rather than the
  spec'd `HttpResponse` record.
- WebSocket `listen` `bind()` fails in some environments (likely missing
  `SO_REUSEADDR` / IPv4-vs-IPv6 dual-stack handling).

## Implementation plan

1. **Define one canonical C `Result` shape** for these intrinsics: a struct
   carrying `{ ok: i64/handle, err: char* }` (reuse the existing Result block ABI
   the codegen already understands — see
   [crates/osprey-codegen/src/result.rs](../../crates/osprey-codegen/src/result.rs)).
   Do not invent a second representation.
2. **Convert HTTP functions** (`httpCreateServer`, `httpListen`, `httpStopServer`,
   client request) to populate that struct: success → handle/value, failure →
   `errno`/getaddrinfo/parse message string. Update the matching `extern_call.rs`
   return descriptors from `Ret::Int` to the Result form.
3. **Fix the `httpListen` handler bridge** to accept the `HttpResponse` record
   (already laid out in [http_shared.h](../../compiler/runtime/http_shared.h))
   instead of a raw string.
4. **Convert WebSocket functions** identically; fix `listen` binding by setting
   `SO_REUSEADDR` and handling address-family selection explicitly.
5. **Thread real error messages** end-to-end so a failed `listen`/`connect`
   surfaces a human-readable string in the `Error` arm ([ERR-PAYLOAD]).
6. **Deduplicate**: HTTP and WebSocket share socket setup — factor the bind/listen
   helper once (run `find-similar` first).

## Testing

- Extend [compiler/examples/tested/http/](../../compiler/examples/tested/http/)
  to `match` on the `Result` of `httpListen` and on a deliberately failing bind
  (port in use) to assert the `Error` message threads through.
- Re-enable / add a WebSocket server `listen` example now that binding is fixed
  (see [compiler/examples/websocketserver/](../../compiler/examples/websocketserver/)).

## Risks / considerations

- Struct layout must match exactly between C and the codegen Result ABI — verify
  field order/sizes against [aggregate.rs](../../crates/osprey-codegen/src/aggregate.rs).
- The bind fix is environment-sensitive; test on macOS (dev) and Linux (CI).

## TODO

- [ ] Settle on the shared C `Result` struct mirroring the codegen Result ABI.
- [ ] Convert HTTP create/listen/stop + client to return `Result<T, string>`;
      update `extern_call.rs` descriptors.
- [ ] Switch `httpListen`'s handler bridge to the `HttpResponse` record.
- [ ] Convert WebSocket functions to `Result<T, string>`.
- [ ] Fix WebSocket `listen` bind (`SO_REUSEADDR`, address family).
- [ ] Thread real error messages into `Error` payloads ([ERR-PAYLOAD]).
- [ ] Factor shared socket setup once (`find-similar` first).
- [ ] Extend `tested/http` (+ websocket) examples with `Result` matching and a
      failing-bind case; refresh `.expectedoutput`.
- [ ] Update the `## Status` sections of 0014 and 0015 once aligned.
- [ ] `make ci` green on macOS and Linux.
