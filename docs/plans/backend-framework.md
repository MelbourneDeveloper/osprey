# Plan Index: Industrial Backend Framework

Thin index only — each layer's detail + TODOs live in its own plan. (Sibling index style:
[`production-primitives.md`](production-primitives.md).)

The goal: stand up a typed, effect-safe HTTP+DB service in a few lines. Two ground rules:

- **Raw SQL is first-class.** Type safety layers *on top* of SQL ([sqlc](https://github.com/sqlc-dev/sqlc)
  model) — no expression trees, no runtime query builder.
- **Capability-gated effects.** DB/Net are algebraic effects behind permission flags; unhandled effects
  are a compile error.

Not greenfield: we beef the existing `httpListen` server; DB is new (today only `readFile`/`writeFile`).

## Layers

| Layer | Adds | Plan | Compiler work |
|---|---|---|---|
| **L0** | General closure capture (gates composable middleware) | [`closures.md`](closures.md) | Large (codegen) |
| **L1** | `Database` effect + SQLite driver, parameterized, `Result`-typed | [`database-effect.md`](database-effect.md) | Small (C shim + 1 permission, **no codegen**) |
| **L2** | shelf-style framework: `Handler`/`Middleware`, `Request`/`Response`, router | [`http-framework.md`](http-framework.md) | None for v0; L0 for v1 |
| **L3** | DataProvider YAML → typed records + CRUD + named SQL; migrations | [`orm-dataprovider.md`](orm-dataprovider.md) | Medium (Go generator) |

## Two load-bearing constraints (this session's probes)

1. **Closure capture is broken language-wide** (function values are bare pointers, no env slot —
   [function_signatures.go:1272](../../compiler/internal/codegen/function_signatures.go#L1272)). Gates
   `Middleware = fn(Handler) -> Handler`. → [`closures.md`](closures.md).
2. **A handler can't close over a `let`** (top-level `let`s run inside `main`); but top-level `fn`s are
   mutually visible. → L2 v0 router is a top-level dispatch `fn` (zero compiler change). → [`http-framework.md`](http-framework.md).

## Build order (least resistance)

1. **L1 SQLite** — self-contained, no codegen, demoable now. Unblocks the ORM runtime. **Start here.**
2. **L0 closures** — gates elegant middleware + record ergonomics.
3. **L2 framework** — v0 today → v1 middleware after L0.
4. **L3 ORM generator** — ties HTTP+DB over real schema.

## Schema source of truth

[DataProvider](https://github.com/Nimblesite/DataProvider) (Nimblesite, .NET) YAML + `DataProviderMigrate`
+ LQL. Upstream ask filed: [Nimblesite/DataProvider#66](https://github.com/Nimblesite/DataProvider/issues/66)
(AOT single-binary migrate + Osprey codegen target). Detail in [`orm-dataprovider.md`](orm-dataprovider.md).

## Master TODO

- [ ] **L1** [`database-effect.md`](database-effect.md): SQLite shim + `Database` effect + `PermissionDatabase`.
- [ ] **L0** [`closures.md`](closures.md): Phase 2 capture + Phase 5 UFCS/field-call disambiguation.
- [ ] **L2** [`http-framework.md`](http-framework.md): v0 (router + Request/Response) → v1 (middleware).
- [ ] **L3** [`orm-dataprovider.md`](orm-dataprovider.md): YAML → records + CRUD → named SQL queries.
- [ ] **L1+** [`database-effect.md`](database-effect.md): Postgres via libpq.
- [ ] Spec: `0016` Database capability + `0012` Database builtins (done); add `0018`/`0019`/`0020` as layers land.
- [ ] One end-to-end golden example `examples/tested/http/todo_service.osp` (HTTP + ORM + SQLite).
