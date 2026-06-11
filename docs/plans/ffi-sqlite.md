# Plan: Generic C Interop → SQLite Driver (L1)

Parent: [`backend-framework.md`](backend-framework.md). Gated by the existing **FFI** capability
(`--no-ffi` / `PermissionFFI`) — see [`0016`](../specs/0016-SecurityAndSandboxing.md). New `0019-Database.md`
(the Osprey-level `Database` effect) lands with Phase 3.

## The decision (one mechanism, not hardcoded DB functions)

The ORM needs ONE of: a C network layer to speak Postgres, **or** a generic C interop layer to call SQLite.
**Choose the generic interop layer + SQLite** — least resistance (embedded; no socket/auth/wire protocol).

**There are NO `db*` compiler builtins.** SQLite is reached the same way any C library is: through generic
**`extern fn`** declarations bound to `libsqlite3`. The DB driver is an **Osprey library**, not compiler
code. This also subsumes Postgres later (bind `libpq` via the same FFI) — a from-scratch network layer is
only needed for a *pure-Osprey* Postgres client, which is optional/future.

## What the generic FFI layer needs (reusable for any C lib)

`extern fn` marshals `int` (i64), `string` (i8*/char*), struct pointers, and i64 opaque handles
(extern lowering in [`crates/osprey-codegen/src/lower.rs`](../../crates/osprey-codegen/src/lower.rs) /
[`extern_call.rs`](../../crates/osprey-codegen/src/extern_call.rs)). Before this plan, linking was
hardcoded to the bundled runtime archives. Gaps closed, all generic:

1. **Arbitrary library linking from source** — a way to say "link `sqlite3`" (a CLI `-l`/link directive or
   an `extern`/annotation form), resolving **system libraries** (`-lsqlite3`), not just bundled archives.
2. **Opaque pointers + out-parameters** — `sqlite3_open_v2(path, &db, …)` and `sqlite3_prepare_v2(db, sql,
   n, &stmt, …)` take pointer-to-pointer out-params. Need a generic opaque `Ptr` (i8*) value Osprey can hold
   and pass, plus a convention for `&out` results (return the pointer, or a tiny generic alloc/deref).
3. **Param binding marshalling** — pass Osprey `int`/`string` into `sqlite3_bind_int64`/`sqlite3_bind_text`
   (bound params only — [OWASP](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html)).

These are language-level FFI improvements — every future C binding (compression, crypto, libpq) reuses them.

## SQLite driver = Osprey `extern fn` library

```osprey
extern fn sqlite3_open(path: string, ppDb: Ptr) -> int
extern fn sqlite3_prepare_v2(db: Ptr, sql: string, n: int, ppStmt: Ptr, tail: Ptr) -> int
extern fn sqlite3_bind_text(stmt: Ptr, idx: int, val: string, n: int, destructor: Ptr) -> int
extern fn sqlite3_step(stmt: Ptr) -> int
extern fn sqlite3_column_text(stmt: Ptr, col: int) -> string
extern fn sqlite3_finalize(stmt: Ptr) -> int
extern fn sqlite3_close(db: Ptr) -> int
```
A thin Osprey module wraps these into `Result`-returning, bound-parameter helpers. The **`Database` effect**
(L1 Phase 3) is an Osprey-level abstraction over that module — still no compiler DB code.

## Status — ✅ COMPLETE (SQLite end-to-end; Postgres binding verified)

All phases landed and verified. Spec: [`0019-Database.md`](../specs/0019-Database.md). Two tested examples
(`examples/tested/db/`) + a Postgres bind smoke demo. Build, clippy, the C runtime's `-Werror` build, crate
unit tests, and the DB examples in the differential harness all green.

## TODO

### Phase 1 — generic FFI buildout  ✅ DONE
- [x] 1.1 Generic library linking via `// @link: <lib>` → `-l<lib>`, plus `// @linkdir: <path>` → `-L<path>`
  (for keg-only/custom libs). The directives are read from source and handed to clang at link time by
  [`crates/osprey-cli/src/main.rs`](../../crates/osprey-cli/src/main.rs) — shared by `--run` and
  `--compile`. Names/paths validated (no injection).
- [x] 1.2 Opaque `Ptr` type (→ i8*) across extern + user-fn signatures; C out-params via generic pointer
  cells in [ffi_runtime.c](../../compiler/runtime/ffi_runtime.c)
  (`osprey_ffi_cell`/`deref`/`free`/`null`/`transient`), archived into both runtimes.
- [x] 1.3 Proven against SQLite (a real system lib).

### Phase 2 — SQLite binding (pure Osprey `extern fn`)  ✅ DONE
- [x] 2.1 `extern fn` decls for the SQLite C API.
- [x] 2.3 Golden [examples/tested/db/sqlite_basics.osp](../../compiler/examples/tested/db/sqlite_basics.osp) (+ `.expectedoutput`).
- [x] 2.4 Wired into the differential harness ([`crates/diff_examples.sh`](../../crates/diff_examples.sh),
  run by `make test`); CI installs `libsqlite3-dev` ([ci.yml](../../.github/workflows/ci.yml)). Passes locally.
- [~] 2.2 Result-typed ergonomics: the `Database` effect's `open` returns `Result<Ptr, Error>`. A fully
  `Result`-everywhere reusable module is deferred — needs (a) an Osprey module/import story and (b) the
  mixed-Result-payload codegen fix below.

### Phase 3 — Database effect + Postgres  ✅ DONE
- [x] 3.1 [`Database` effect](../../compiler/examples/tested/db/database_effect.osp) over the SQLite bindings
  (capability-gated `!Database`, swappable handler, bound params, typed rows); spec
  [`0019-Database.md`](../specs/0019-Database.md). Verified.
- [x] 3.2 Postgres binds through the **same** FFI: [examples/db_postgres/pg_smoke.osp](../../compiler/examples/db_postgres/pg_smoke.osp)
  (`// @link: pq` + `// @linkdir:`, `PQlibVersion()` → "pg linked ok"). A full connect+query round-trip binds
  `PQconnectdb`/`PQexecParams`/`PQgetvalue`/`PQfinish` identically and needs a running server (not auto-tested).

## Known limitation (discovered; future fix)

Mixing two `Result<T, Error>` payload shapes in one program — e.g. `Result<Ptr, Error>` (`{i8*, i8}`) and
`Result<int, Error>` (`{i64, i8}`) — makes `llc` reject the IR (`instruction forward referenced`). The
`Database` effect works around it by having only `open` return `Result` (others return raw rc/handles). Root
cause is the generic-`Result` payload lowering, related to the error-payload work in
[`production-primitives.md`](production-primitives.md); fixing it would let every fallible op return `Result`.

## Authorities

[SQLite C API](https://sqlite.org/cintro.html) ([bind](https://www.sqlite.org/c3ref/bind_blob.html),
[column](https://www.sqlite.org/c3ref/column_blob.html), [open_v2](https://www.sqlite.org/c3ref/open.html)),
[OWASP Query Parameterization](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html),
[libpq](https://www.postgresql.org/docs/current/libpq-exec.html).
