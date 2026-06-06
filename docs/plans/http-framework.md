# Plan: HTTP Framework (L2, shelf-style)

Parent: [`backend-framework.md`](backend-framework.md). Builds on [`0014-HTTP.md`](../specs/0014-HTTP.md);
new `0018-HTTPFramework.md` lands with Phase 2. v1 gated on [`closures.md`](closures.md).

## Goal

Beef the bare `httpListen` callback into a [shelf](https://pub.dev/packages/shelf)-style framework. Target:

```osprey
fn main() = serve(port: 8080, handler:
    pipeline([logRequests, cors, recover]) |> router([
        get("/health",     fn(req) => text(200, "ok")),
        get("/users/<id>", fn(req) => json(200, lookup(param(req, "id")))),
        post("/users",     fn(req) => created(insertUser(body(req)))),
    ]))
```

## Key decisions (each cited in [Authorities](#authorities))

- **Core model adopted literally:** `Handler = fn(Request) -> Response`, `Middleware = fn(Handler) -> Handler` (shelf).
- **Drop `FutureOr`:** IO/async is an **effect row**, not a wrapper type — the checker sees "this handler does IO".
- **Immutable `Request`/`Response`** records + copy-with helpers (shelf `change()`); reject Go/Express mutable sink.
- **Typed effects replace shelf's stringly `context` bag** for per-request state (auth/req-id/DB handle);
  middleware installs a handler, downstream `perform`s it — and this is the Koa onion mechanism.
- **USP:** effect-typed middleware lets the checker enforce "auth route only reachable via auth middleware".
- **Router:** shelf `<id>`/`<id|regex>` params; per-method, most-specific-wins; **404 vs 405+`Allow`**
  ([RFC 9110](https://www.rfc-editor.org/rfc/rfc9110)); linear match v0 → radix trie v2; `mount` + `Cascade` later.
- **Parsing split:** C populates headers + splits `?query` (one tiny runtime change at
  [http_server_runtime.c:69-70](../../compiler/runtime/http_server_runtime.c#L69-L70)); routing + `<id>` +
  header/query parse in Osprey via existing `split`/`indexOf`/`substring` builtins.

## v0 vs v1

- **v0 (no compiler change):** router is a **top-level dispatch `fn`** (handler can't close over a `let`);
  Request/Response records, helpers, `<id>` params. Middleware = hand-composed top-level `fn`s. Ships now.
- **v1 (after [`closures.md`](closures.md) Phase 2):** real `Pipeline` (`foldr` of middleware), `Cascade`,
  data-valued `router([...])`, shippable `recover`/`logging`/`cors` middleware.

## Records (target shape)

```osprey
type Request  = { method: string, path: string, query: string, headers: string, body: string, params: Map<string,string> }
type Response = { status: int, contentType: string, headers: string, body: string }
```
Helpers: `ok/text/json/created/noContent/badRequest/unauthorized/forbidden/notFound/methodNotAllowed/internalError`.

## TODO

### Phase 1 — v0 (no compiler change)
- [ ] 1.1 `Request`/`Response` records + response helpers (extend, don't fork, `http_server_example.osp`)
- [ ] 1.2 C-side: populate `raw_headers`; split `?query` off `path` ([http_server_runtime.c:55-70](../../compiler/runtime/http_server_runtime.c#L55-L70))
- [ ] 1.3 Top-level `fn` router: `<id>` params, header/query parse in Osprey, 404 vs **405 + `Allow`**
- [ ] 1.4 Golden: `examples/tested/http/framework_basic.osp` (server + client round-trip) + expected output

### Phase 2 — v1 (after L0 closures)
- [ ] 2.1 `Handler`/`Middleware` type aliases; `pipeline([...])` as `foldr`; `Cascade` (404/405 fallthrough)
- [ ] 2.2 Data-valued `router([...])`; ship `recover`→500, `logRequests`+latency, `cors` preflight
- [ ] 2.3 Spec `0018-HTTPFramework.md`; extend golden example with middleware

### Phase 3 — scale
- [ ] 3.1 Radix-trie router; `mount` sub-routers; typed positional params
- [ ] 3.2 Streaming bodies (lazy byte streams ⇒ backpressure) + `hijack` for WS/SSE

### Phase 4 — production hardening (adapter)
- [ ] 4.1 Graceful shutdown (drain) + Read/Write/Idle timeouts; body-size limit → 413
- [ ] 4.2 Content negotiation → 406/`Vary`; gzip + `Vary`

## Authorities

[shelf](https://pub.dev/packages/shelf) + [shelf_router](https://pub.dev/documentation/shelf_router/latest/shelf_router/Router-class.html),
[PEP 3333](https://peps.python.org/pep-3333/), [Rack](https://github.com/rack/rack/blob/main/SPEC.rdoc),
[Express](https://expressjs.com/en/guide/using-middleware.html)/[Koa](https://github.com/koajs/koa/blob/master/docs/guide.md),
[chi](https://github.com/go-chi/chi)/[httprouter](https://github.com/julienschmidt/httprouter),
[Go Server.Shutdown](https://pkg.go.dev/net/http#Server.Shutdown), [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110),
[WHATWG Fetch CORS](https://fetch.spec.whatwg.org/#http-cors-protocol).
