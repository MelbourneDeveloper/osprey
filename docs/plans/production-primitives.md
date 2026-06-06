# Plan Index: Production-App Primitives

## Why this exists

The goal: a developer can write a real production app in Osprey. The honest test of that claim is **whether they can write a JSON parser in Osprey itself** — not call out to a C builtin. If they can't, then every future user library (CSV, TOML, query strings, markdown, log parsing, configuration formats) hits the same wall, and no amount of bolted-on C functions makes the language self-hosting.

A probe in this session ([/tmp/probe_recursive_union.osp](../../tmp/probe_recursive_union.osp), [/tmp/probe_closure_capture.osp](../../tmp/probe_closure_capture.osp), [/tmp/probe_lambda_simple.osp](../../tmp/probe_lambda_simple.osp)) confirmed four primitives are broken or missing today:

| # | Primitive | Status | Plan |
|---|---|---|---|
| 1 | Lambdas at all (`fn(x) => x + 1`) | ✅ no-capture lambdas land; ❌ **capture** still broken (`use of undefined value '%n'`) | [`closures.md`](closures.md) |
| 2 | Recursive unions with `List<Self>` / `Map<K,Self>` payload | ❌ codegen panics: `store operands are not compatible: src=i8*; dst=i1*` at [expression_generation.go:1631](../../compiler/internal/codegen/expression_generation.go#L1631) | [`recursive-union-payloads.md`](recursive-union-payloads.md) |
| 3 | Error message payload threading through `Result<T, E>` | ❌ hardcoded `"Error occurred"` global at [llvm.go:2305](../../compiler/internal/codegen/llvm.go#L2305) | [`error-payloads.md`](error-payloads.md) |
| 4 | O(1) codepoint/byte cursor over `string` | ❌ no `byteAt` / `codePointAt`; every existing op (`take`/`drop`/`substring`) allocates | [`string-cursor.md`](string-cursor.md) |
| 5 | List patterns (`[head, ...tail]`) | ❌ spec'd at [TYPE-LIST-PATTERNS] but no AST node / no codegen — escalated to critical-path | [`list-patterns.md`](list-patterns.md) |

What works today that we are **building on**:
- Plain self-recursive unions (`Tree = Leaf | Node { left: Tree, right: Tree }`) compile and run.
- `List<T>` and `Map<K, V>` persistent runtime is in place ([`collection_codegen.go`](../../compiler/internal/codegen/collection_codegen.go), 15 C runtime tests pass).
- String functions (`split`, `indexOf`, `startsWith`, `trim`, `parseInt`, …) work — but every one of them allocates.
- Pattern matching on union variants with bound fields works.

## Spec changes that landed alongside these plans

- `0004-TypeSystem.md` — new `Closures` subsection under Function Types ([TYPE-FN-CLOSURE]); new top-level `Union Types` section with explicit recursive-payload requirement ([TYPE-UNION-REC]).
- `0012-Built-InFunctions.md` — new `Cursor Access` section: `byteLength`, `byteAt`, `codePointAt`, `codePointWidth`, `fromCodePoint`.
- `0013-ErrorHandling.md` — new `Error Payload Propagation` requirement ([ERR-PAYLOAD]) making the hardcoded-message implementation explicitly non-conforming.

## Sequencing

The five plans are **independent at the implementation level** but have a natural priority order if a single agent picks them up:

1. **`error-payloads.md`** first — smallest, most contained, unblocks meaningful error messages everywhere immediately.
2. **`closures.md`** second — small surface, but every higher-order primitive (`map`, `filter`, `fold`, parser combinators) depends on it.
3. **`recursive-union-payloads.md`** third — unblocks `JsonValue`, tree types, ASTs.
4. **`string-cursor.md`** fourth — adds the C primitives that the JSON parser will sit on top of.
5. **`list-patterns.md`** fifth — wraps the parser in idiomatic recursive descent (`[head, ...tail]`).

After all four land, the proof point is one tested example: `examples/tested/json/json_parser.osp` — a JSON parser **written in Osprey** that consumes a real input and produces a `JsonValue`. That example becomes the canary for every future regression in these four areas.

## Sibling plans

[`backend-framework.md`](backend-framework.md) — the industrial HTTP framework + typed DB/ORM ecosystem
that sits **on top** of these primitives (its composable middleware is gated on [`closures.md`](closures.md)
Phase 2).

[`tui-http-app.md`](tui-http-app.md) — colored TUI that calls ad-hoc HTTP APIs. Independent runtime/builtin work (HTTP response bodies, terminal raw mode, ANSI helpers) but shares the "needs JSON" requirement. The TUI plan ships a C-side JSON builtin as a v1 shortcut and deletes it once this plan's Osprey-native parser lands.

## Master TODO (across all five plans)

- [ ] Land `error-payloads.md` Phase 1 (runtime threading) and Phase 2 (codegen rewrite).
- [ ] Land `closures.md` Phase 1 (lambda param binding fix) and Phase 2 (capture).
- [ ] Land `recursive-union-payloads.md` Phase 1 (layout) and Phase 2 (codegen).
- [ ] Land `string-cursor.md` Phase 1 (C runtime) and Phase 2 (builtins + registry).
- [ ] Land `list-patterns.md` Phase 1–3 (grammar, codegen, `osprey_list_drop` runtime).
- [ ] Land `examples/tested/json/json_parser.osp` written in pure Osprey, using only the above primitives plus existing `List`/`Map`/`match`. Must round-trip RFC 8259 conforming inputs.
- [ ] Update `coverage-to-90-percent.md` ratchet plan once the JSON example's tests are exercising the new branches.
