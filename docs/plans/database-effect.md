# Plan: Database Effect + Driver (L1)

Parent: [`backend-framework.md`](backend-framework.md). Spec: [`0016`](../specs/0016-SecurityAndSandboxing.md)
Database capability + [`0012`](../specs/0012-Built-InFunctions.md) Database builtins (added); new `0019-Database.md`
lands with the effect. **Start here** — self-contained, no codegen.

## Goal

A typed, capability-gated `Database` effect with full SQL flexibility and **mandatory bound parameters**
(injection structurally impossible). SQLite first (embedded = least resistance); Postgres later via libpq,
behind the **same** effect. Today there is no DB layer — only `readFile`/`writeFile`.

## Key decisions

- **SQLite via C shim, not raw `extern fn`.** The `json`/`http` builtin-shim pattern gives `Result`-typed
  returns with **zero new codegen**: `generateRuntimeBuiltinCall`
  ([runtime_builtins_generation.go:25-73](../../compiler/internal/codegen/runtime_builtins_generation.go#L25-L73))
  already wraps `int64` handles (neg → `Error`) and `char*`/NULL (NULL → `Error`).
- **Never hand-roll Postgres wire protocol / SCRAM** ([RFC 7677](https://datatracker.ietf.org/doc/html/rfc7677)) —
  `-lpq` gives protocol v3 + auth + TLS for free.
- **Parameter-only API** — no interpolated-SQL function exists ([OWASP](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html)).
- **Row model:** typed column accessors primary (SQLite `sqlite3_column_type`/`_name`; PG `PQftype`/`PQfname`/
  `PQgetisnull`); rows-as-JSON `queryJson` as the dynamic escape hatch.

## Effect surface (engine-agnostic)

```osprey
effect Database {
  open  : fn(string) -> Result<int, Error>                 // conn str; sqlite://… | postgres://…
  exec  : fn(int, string, [Param]) -> Result<int, Error>   // DML/DDL → rows affected
  query : fn(int, string, [Param]) -> Result<int, Error>   // SELECT → result-set handle
  close : fn(int) -> Result<int, Error>
}
```
`Param = PInt | PFloat | PText | PBlob | PNull`. Engine selected by conn string, not signatures.

## TODO

### Phase 1 — SQLite (no codegen)
- [ ] 1.1 `compiler/runtime/database_runtime.c` mirroring [json_runtime.c](../../compiler/runtime/json_runtime.c):
  mutex-guarded handle tables; `db_open/db_exec/db_query/db_row/db_column/db_free_result/db_close` over
  `sqlite3_open_v2`/`prepare_v2`/`bind_*`(`SQLITE_TRANSIENT`)/`step`/`column_*`/`finalize`; `strdup` strings.
- [ ] 1.2 `Makefile`: compile `database_runtime.c` in both runtime targets ([:76](../../compiler/Makefile#L76),
  [:122](../../compiler/Makefile#L122)); archive `.o` ([:78](../../compiler/Makefile#L78), [:124](../../compiler/Makefile#L124));
  `-lsqlite3` after [compilation.go:369](../../compiler/internal/codegen/compilation.go#L369) (+ JIT
  [jit_executor.go:117](../../compiler/internal/codegen/jit_executor.go#L117)) with Homebrew `-L` fallback.
- [ ] 1.3 `constants.go` `DB*Osprey`/`DB*Func` block beside JSON ([:220-235](../../compiler/internal/codegen/constants.go#L220-L235)).
- [ ] 1.4 `builtin_registry.go` Database group after JSON ([:1484](../../compiler/internal/codegen/builtin_registry.go#L1484))
  via `reg`/`str`/`intp`; return-type strings exactly `"Result<int, Error>"`/`"Result<string, Error>"`;
  `SecurityFlag: PermissionDatabase`; add `CategoryDatabase`.
- [ ] 1.5 `PermissionDatabase` in the enum ([:72-90](../../compiler/internal/codegen/builtin_registry.go#L72-L90));
  thread `AllowDatabase` through `checkSecurityPermission` ([llvm.go:894-913](../../compiler/internal/codegen/llvm.go#L894-L913)),
  `SecurityConfig` ([generator.go:82-90](../../compiler/internal/codegen/generator.go#L82-L90)),
  CLI + `--no-db` ([security.go:10-61](../../compiler/internal/cli/security.go#L10-L61)).
- [ ] 1.6 C runtime test for the handle table (mirror fiber/json test style); run `make c-test`.
- [ ] 1.7 Golden `examples/tested/db/sqlite_basics.osp`: open → DDL → param insert → typed query → match rows → close.

### Phase 2 — Postgres (same effect)
- [ ] 2.1 `database_pg_runtime.c` (~8 fns): `PQconnectdb`/`PQexecParams`(text)/`PQntuples`/`PQfname`/`PQgetvalue`/`PQfinish`; `-lpq`.
- [ ] 2.2 Prefer `PQexecParams` over `PQprepare` (survives transaction-pooled [PgBouncer](https://www.pgbouncer.org/features.html)).
- [ ] 2.3 Example switches only the connection string.

### Phase 3 — ergonomics
- [ ] 3.1 Typed column-accessor row API + `queryJson` escape hatch
- [ ] 3.2 Transaction helpers (`begin`/`commit`/`rollback`)

## Authorities

[SQLite C API](https://sqlite.org/cintro.html) ([bind](https://www.sqlite.org/c3ref/bind_blob.html),
[column](https://www.sqlite.org/c3ref/column_blob.html), [threading](https://www.sqlite.org/threadsafe.html)),
[OWASP](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html),
[libpq](https://www.postgresql.org/docs/current/libpq-exec.html),
[PG protocol](https://www.postgresql.org/docs/current/protocol.html),
[RFC 7677](https://datatracker.ietf.org/doc/html/rfc7677).
