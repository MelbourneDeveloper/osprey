# Database Access

Osprey has **no built-in database functions**. Databases are reached through the generic
foreign-function interface ([FFI](0005-FunctionCalls.md), `extern fn`): a driver is a set of `extern`
declarations bound to a C client library, linked with a `// @link:` directive. SQLite and PostgreSQL are
both reached this way — the language core stays database-agnostic.

Plan: [`ffi-sqlite.md`](../plans/ffi-sqlite.md). Examples:
[`sqlite_basics.osp`](../../compiler/examples/tested/db/sqlite_basics.osp) (raw bindings) and
[`database_effect.osp`](../../compiler/examples/tested/db/database_effect.osp) (the effect).

## Status

Implemented and verified for SQLite: generic library linking, the opaque `Ptr` type, FFI out-parameter
cells, the SQLite bindings, and a `Database` effect over them. PostgreSQL binds through the same mechanism
(`// @link: pq`); a full server round-trip requires a running PostgreSQL instance.

## Linking a C library — `// @link:` [DB-LINK]

A source-level directive requests that a third-party C library be linked:

```osprey
// @link: sqlite3
extern fn sqlite3_libversion() -> string
```

The directive maps to a `-l<lib>` linker flag and is carried through the IR so it applies to both
`--compile` and JIT (`--run`). Library names are validated (no linker-flag or shell metacharacters). A
companion `// @linkdir: <path>` directive adds a `-L<path>` search directory for libraries installed
outside the default search path (e.g. a Homebrew keg-only `libpq`). Linking is gated by the FFI capability
(`--no-ffi` / `--sandbox` disable it) — see [0016-SecurityAndSandboxing.md](0016-SecurityAndSandboxing.md).

## The opaque `Ptr` type [DB-PTR]

`Ptr` is the opaque foreign-pointer type (a C `void*`). It is used for connection handles, statement
handles, and out-parameter cells. C out-parameters (e.g. `sqlite3_open(const char*, sqlite3**)`) are
handled with generic FFI pointer cells from the runtime — no per-library C code:

| Function | Purpose |
|----------|---------|
| `osprey_ffi_cell() -> Ptr` | allocate a zeroed pointer cell for an out-parameter |
| `osprey_ffi_deref(cell: Ptr) -> Ptr` | read the pointer written into a cell |
| `osprey_ffi_free(cell: Ptr) -> int` | release a cell |
| `osprey_ffi_null() -> Ptr` | a `NULL` pointer for optional C arguments |
| `osprey_ffi_transient() -> Ptr` | SQLite's `SQLITE_TRANSIENT` sentinel (copy-the-buffer) |

## Bound parameters are mandatory [DB-PARAMS]

Values are always **bound** to placeholders (`?`, `$1`), never interpolated into SQL text, making
injection structurally impossible (per the [OWASP Query Parameterization Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Query_Parameterization_Cheat_Sheet.html)).
A driver binds with the engine's bind calls (e.g. `sqlite3_bind_int` / `sqlite3_bind_text`).

## The `Database` effect [DB-EFFECT]

Database access SHOULD be modeled as an [algebraic effect](0017-AlgebraicEffects.md) so that it is a
capability: a function that touches the database is `!Database`, an unhandled `Database` effect is a
compile error, and the handler is swappable (a real driver in production, a mock in tests). Connection and
statement handles are passed through operation arguments so handler arms hold no captured state. Fallible
operations return `Result<_, Error>`.

```osprey
effect Database {
    open: fn(string) -> Result<Ptr, Error>
    exec: fn(Ptr, string) -> int
    prepare: fn(Ptr, string) -> Ptr
    bindInt: fn(Ptr, int, int) -> int
    bindText: fn(Ptr, int, string) -> int
    step: fn(Ptr) -> int
    columnInt: fn(Ptr, int) -> int
    columnText: fn(Ptr, int) -> string
    finalize: fn(Ptr) -> int
    close: fn(Ptr) -> int
}
```

The handler is the only code that references the C bindings; application logic performs operations and
never sees SQLite or the FFI. Selecting PostgreSQL instead of SQLite is a matter of swapping the handler
(binding `libpq`'s `PQ*` functions) — application code is unchanged.
