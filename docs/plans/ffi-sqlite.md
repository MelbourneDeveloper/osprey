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

Today `extern fn` already marshals `int` (i64), `string` (i8*/char*), struct pointers, and i64 opaque
handles ([statement_generation.go:101-141](../../compiler/internal/codegen/statement_generation.go#L101-L141)),
but linking is **hardcoded** to bundled `lib<name>.a` for `rust_utils`
([compilation.go:129-163](../../compiler/internal/codegen/compilation.go#L129-L163),
[:369](../../compiler/internal/codegen/compilation.go#L369)). Gaps to close, all generic:

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

## Status

**Phase 1 + 2 landed and verified** via `./bin/osprey examples/tested/db/sqlite_basics.osp --run`: a typed,
bound-parameter SQLite round-trip (`:memory:` → DDL → two bound inserts → ordered `SELECT` → typed int+text
columns → close) from pure Osprey `extern fn` — zero hardcoded DB builtins. Build/lint (Go + C) + codegen
unit tests green.

## TODO

### Phase 1 — generic FFI buildout  ✅ DONE
- [x] 1.1 Generic library linking via the `// @link: <lib>` source directive → `-l<lib>`, carried in the IR
  as a `; osprey-link:` marker so it survives both `--compile` and JIT. New
  [ffi_linking.go](../../compiler/internal/codegen/ffi_linking.go); wired into
  [compilation.go](../../compiler/internal/codegen/compilation.go) +
  [jit_executor.go](../../compiler/internal/codegen/jit_executor.go). Lib names validated (no flag/shell injection).
- [x] 1.2 Opaque `Ptr` type (→ i8*) across extern + user-fn signatures (`TypePtr` in
  [constants.go](../../compiler/internal/codegen/constants.go),
  [statement_generation.go](../../compiler/internal/codegen/statement_generation.go),
  [function_signatures.go](../../compiler/internal/codegen/function_signatures.go)); C out-params via generic
  pointer cells in [ffi_runtime.c](../../compiler/runtime/ffi_runtime.c)
  (`osprey_ffi_cell`/`deref`/`free`/`null`/`transient`), archived into both runtimes (Makefile).
- [x] 1.3 Proven against SQLite (a real system lib): `// @link: sqlite3` + `sqlite3_libversion()`.

### Phase 2 — SQLite binding (pure Osprey `extern fn`)
- [x] 2.1 `extern fn` decls: `sqlite3_open`/`exec`/`prepare_v2`/`bind_int`/`bind_text`/`step`/`column_int`/`column_text`/`finalize`/`close`.
- [x] 2.3 Golden [examples/tested/db/sqlite_basics.osp](../../compiler/examples/tested/db/sqlite_basics.osp)
  (+ `.expectedoutput`) — **verified locally**.
- [ ] 2.2 Factor extern decls + helpers into a reusable `Result`-returning module (blocked on a module/import
  story; inline in the example for now).
- [ ] 2.4 Wire the example into a test suite — **needs `libsqlite3-dev` on CI runners** (harness forbids skips);
  gated on a CI-dependency decision so CI stays green.

### Phase 3 — Database effect + Postgres
- [ ] 3.1 Osprey `Database` effect wrapping the SQLite module; spec `0019-Database.md`.
- [ ] 3.2 Postgres: bind `libpq` (`PQconnectdb`/`PQexecParams`/`PQgetvalue`/`PQfinish`) via the **same** FFI
  (`// @link: pq`), same `Database` effect — engine chosen by connection string. (Pure-Osprey wire client = future.)

## Authorities

[SQLite C API](https://sqlite.org/cintro.html) ([bind](https://www.sqlite.org/c3ref/bind_blob.html),
[column](https://www.sqlite.org/c3ref/column_blob.html), [open_v2](https://www.sqlite.org/c3ref/open.html)),
[OWASP Query Parameterization](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html),
[libpq](https://www.postgresql.org/docs/current/libpq-exec.html).
