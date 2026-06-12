# Plan Index: Production-App Primitives

## Why this exists

The goal: a developer can write a real production app in Osprey. The honest test of that claim is **whether they can write a JSON parser in Osprey itself** — not call out to a C builtin. If they can't, then every future user library (CSV, TOML, query strings, markdown, log parsing, configuration formats) hits the same wall, and no amount of bolted-on C functions makes the language self-hosting.

Probes (re-run against the Rust compiler, `target/release/osprey`) put the five primitives here today:

| # | Primitive | Status | Plan |
|---|---|---|---|
| 1 | Lambdas at all (`fn(x) => x + 1`) | ✅ **landed** — one closure model: every function value is a closure cell `{ fnptr, captures… }` ([`crates/osprey-codegen/src/closure.rs`](../../crates/osprey-codegen/src/closure.rs)); escaping closures (`makeAdder`), curried middleware, nested capture, capturing lambdas in record fields and iterator callbacks all work; `spawn` shares the same cells (per-instance, re-entrant). Golden coverage in `function_composition_test.osp`. Remaining: UFCS field-call disambiguation (`obj.fnField(x)`) — see below | — |
| 2 | Recursive unions with `List<Self>` / `Map<K,Self>` payload | ✅ landed — tagged heap layout with pointer-indirected payloads ([`crates/osprey-codegen/src/aggregate.rs`](../../crates/osprey-codegen/src/aggregate.rs)); golden coverage in `recursive_unions.osp` | [`recursive-union-payloads.md`](recursive-union-payloads.md) |
| 3 | Error message payload threading through `Result<T, E>` | ❌ the `Error { message }` arm binds the zeroed payload slot — a division-by-zero probe interpolates `message` as `0.0` instead of a reason ([`crates/osprey-codegen/src/result.rs`](../../crates/osprey-codegen/src/result.rs)) | [`error-payloads.md`](error-payloads.md) |
| 4 | O(1) codepoint/byte cursor over `string` | ❌ no `byteAt` / `codePointAt`; every existing op (`take`/`drop`/`substring`) allocates | [`string-cursor.md`](string-cursor.md) |
| 5 | List patterns (`[head, ...tail]`) | ❌ spec'd at [TYPE-LIST-PATTERNS] but no AST node / no codegen — escalated to critical-path | [`list-patterns.md`](list-patterns.md) |

What works today that we are **building on**:
- Plain self-recursive unions (`Tree = Leaf | Node { left: Tree, right: Tree }`) compile and run.
- `List<T>` and `Map<K, V>` persistent runtime is in place ([`crates/osprey-codegen/src/collections.rs`](../../crates/osprey-codegen/src/collections.rs) over [`compiler/runtime/`](../../compiler/runtime/), 15 C runtime tests pass).
- String functions (`split`, `indexOf`, `startsWith`, `trim`, `parseInt`, …) work — but every one of them allocates.
- Pattern matching on union variants with bound fields works.

## Spec changes that landed alongside these plans

- `0004-TypeSystem.md` — new `Closures` subsection under Function Types ([TYPE-FN-CLOSURE]); new top-level `Union Types` section with explicit recursive-payload requirement ([TYPE-UNION-REC]).
- `0012-Built-InFunctions.md` — new `Cursor Access` section: `byteLength`, `byteAt`, `codePointAt`, `codePointWidth`, `fromCodePoint`.
- `0013-ErrorHandling.md` — new `Error Payload Propagation` requirement ([ERR-PAYLOAD]) making the hardcoded-message implementation explicitly non-conforming.

## Sequencing

The five plans are **independent at the implementation level** but have a natural priority order if a single agent picks them up:

1. **`error-payloads.md`** first — smallest, most contained, unblocks meaningful error messages everywhere immediately.
2. **`string-cursor.md`** next — adds the C primitives that the JSON parser will sit on top of.
3. **`list-patterns.md`** last — wraps the parser in idiomatic recursive descent (`[head, ...tail]`).

Closures and recursive unions are done. One closure follow-up survives:

- **UFCS field-call disambiguation** — `obj.fnField(args)` must call the
  function-valued field instead of UFCS-rewriting to `fnField(obj, args)`
  (the rule is already specified at
  [0012-Built-InFunctions.md](../specs/0012-Built-InFunctions.md): "field
  access wins; UFCS is the fallback"). The rewrite happens pre-inference in
  [`crates/osprey-syntax/src/expr.rs`](../../crates/osprey-syntax/src/expr.rs)
  (`lower_call`), so either parse to a neutral `MethodCall` node resolved
  during inference, or have the checker re-interpret `f(x, …)` when `x`'s
  record type declares a function-typed field `f`. Until then the working
  spelling is `let f = obj.fnField` then `f(args)` (golden-covered;
  `failscompilation/function_typed_record_field.ospo` pins the rejection).

After all four land, the proof point is one tested example: `examples/tested/json/json_parser.osp` — a JSON parser **written in Osprey** that consumes a real input and produces a `JsonValue`. That example becomes the canary for every future regression in these four areas.

## Sibling plans

[`backend-framework.md`](backend-framework.md) — the industrial HTTP framework + typed DB/ORM ecosystem
that sits **on top** of these primitives (its composable middleware gate — escaping closures — has landed).

[`tui-http-app.md`](tui-http-app.md) — colored TUI that calls ad-hoc HTTP APIs. Independent runtime/builtin work (HTTP response bodies, terminal raw mode, ANSI helpers) but shares the "needs JSON" requirement. The TUI plan ships a C-side JSON builtin as a v1 shortcut and deletes it once this plan's Osprey-native parser lands.

## Master TODO (across all five plans)

- [ ] Land `error-payloads.md` Phase 1 (runtime threading) and Phase 2 (codegen rewrite).
- [x] Land closures (escaping capture, one closure-cell model; plan completed and deleted). Follow-up: UFCS field-call disambiguation (see Sequencing).
- [x] Land `recursive-union-payloads.md` Phase 1 (layout) and Phase 2 (codegen).
- [ ] Land `string-cursor.md` Phase 1 (C runtime) and Phase 2 (builtins + registry).
- [ ] Land `list-patterns.md` Phase 1–3 (grammar, codegen, `osprey_list_drop` runtime).
- [ ] Land `examples/tested/json/json_parser.osp` written in pure Osprey, using only the above primitives plus existing `List`/`Map`/`match`. Must round-trip RFC 8259 conforming inputs.
