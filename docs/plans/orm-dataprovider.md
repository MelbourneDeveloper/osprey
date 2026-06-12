# Plan: DataProvider YAML → Typed ORM (L3)

Parent: [`backend-framework.md`](backend-framework.md). Runtime dep: [`ffi-sqlite.md`](ffi-sqlite.md).
Ergonomics dep: closures — **landed** (closure cells, `crates/osprey-codegen/src/closure.rs`). New `0020-ORM.md` lands with Phase 1.

## Goal

Generate typed Osprey data-access from [DataProvider](https://github.com/Nimblesite/DataProvider) YAML —
records + CRUD + named SQL queries — the [sqlc](https://docs.sqlc.dev/) way: SQL is the source of truth,
codegen only synthesizes the param signature + result record. **No expression trees, no runtime builder.**

## Three tiers, one foundation (raw SQL is tier 1)

All bottom out at the `Database` effect's parameterized `query`/`exec`:

1. **Raw `query(conn, sql, params)`** — total flexibility (joins, CTEs, windows, `RETURNING`); you write SQL.
2. **Named typed SQL** (sqlc-style) — annotated SQL → typed function + `<Query>Row` record:
   ```sql
   -- name: activeUsersByOrg :many
   SELECT u.id, u.name FROM users u JOIN orgs o ON o.id = u.org_id WHERE o.id = $1 ORDER BY u.name;
   ```
   → `fn activeUsersByOrg(orgId: Uuid) -> Result<[ActiveUsersByOrgRow], Error> !Database`
3. **Generated CRUD** — zero-ceremony sugar over the same primitive; never a ceiling.

## Type mapping (cited)

`isNullable: true` → `Option<T>` ([Diesel](https://diesel.rs/guides/getting-started)/[sqlc](https://docs.sqlc.dev/en/latest/reference/datatypes.html));
one record per table; projections get a `<Query>Row` record. PKs are single-column `Guid` (DataProvider convention).

| DataProvider | non-null | nullable |  | DataProvider | non-null | nullable |
|---|---|---|---|---|---|---|
| `Text` | `String` | `Option<String>` | | `Guid` | `Uuid` | `Option<Uuid>` |
| `Integer` | `Int` | `Option<Int>` | | `DateTime` | `Timestamp` | `Option<Timestamp>` |
| `Real` | `Float` | `Option<Float>` | | `Blob` | `Bytes` | `Option<Bytes>` |
| `Boolean` | `Bool` | `Option<Bool>` | | `vector`(PG) | `[Float]` | `Option<[Float]>` |

## Generated CRUD (parameterized at gen time; binds values at runtime)

| Function | sqlc | SQL | Return |
|---|---|---|---|
| `findById(id)` | `:one` | `SELECT <cols> FROM <t> WHERE <pk>=$1 LIMIT 1` | `Result<Option<E>, Error>` |
| `all()` | `:many` | `SELECT <cols> FROM <t>` | `Result<[E], Error>` |
| `insert(rec)` | `:one` | `INSERT INTO <t>(<cols>) VALUES($1..$n) RETURNING <cols>` | `Result<E, Error>` |
| `update(rec)` | `:exec` | `UPDATE <t> SET <col=$i…> WHERE <pk>=$k` | `Result<Unit, Error>` |
| `delete(id)` | `:exec` | `DELETE FROM <t> WHERE <pk>=$1` | `Result<Unit, Error>` |

## Decisions

- **Generator = `osprey gen-data <schema.yaml>`** emitting `.osp` (least resistance: no .NET SDK dep; pure
  testable transform; type-map lives next to the HM inferer). Commit generated `.osp` beside the YAML.
- **Migrations: delegate** to `DataProviderMigrate migrate --schema … --provider sqlite|postgres` (declarative,
  idempotent, forward-only) — don't fork the diff engine.
- **LQL deferred** (no Osprey parser yet); its pipe syntax is already Osprey's idiom — natural future surface.
- **Upstream:** [Nimblesite/DataProvider#66](https://github.com/Nimblesite/DataProvider/issues/66) (AOT migrate + Osprey target).

## TODO

### Phase 1 — YAML → records + CRUD
- [ ] 1.1 `osprey gen-data` subcommand: parse + validate DataProvider YAML (Go yaml) → schema IR
- [ ] 1.2 Type-map each column (table above); emit one record per table
- [ ] 1.3 Emit `findById/all/insert/update/delete` (parameterized SQL) over the `Database` effect
- [ ] 1.4 Golden `examples/tested/db/orm_generated.osp` (checked-in generated `.osp` + round-trip test)

### Phase 2 — named typed SQL
- [ ] 2.1 Parse `-- name: … :one|:many|:exec` annotations; infer `<Query>Row` from selected columns
- [ ] 2.2 Compile-time column/type check against the schema (sqlc `verify` analog)

### Phase 3 — migrations + relations
- [ ] 3.1 `osprey migrate` shells out to `DataProviderMigrate`; optional shadow-DB drift check
- [ ] 3.2 FK-aware helpers
- [ ] 3.3 Spec `0020-ORM.md`

### Phase 4 — LQL (gated on parser story)
- [ ] 4.1 LQL → parameterized SQL transpilation

## Authorities

[sqlc](https://docs.sqlc.dev/) ([annotations](https://docs.sqlc.dev/en/latest/reference/query-annotations.html)),
[Prisma Migrate](https://www.prisma.io/docs/orm/prisma-migrate/getting-started),
[Ent](https://entgo.io/docs/crud/), [Diesel](https://diesel.rs/guides/getting-started),
[DataProvider](https://github.com/Nimblesite/DataProvider) +
[Migration README](https://raw.githubusercontent.com/Nimblesite/DataProvider/main/Migration/README.md).
