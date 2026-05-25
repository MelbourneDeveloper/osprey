# Plan Index: Production-App Primitives

## Why this exists

The goal: a developer can write a real production app in Osprey. The honest test of that claim is **whether they can write a JSON parser in Osprey itself** — not call out to a C builtin. If they can't, then every future user library (CSV, TOML, query strings, markdown, log parsing, configuration formats) hits the same wall, and no amount of bolted-on C functions makes the language self-hosting.

A probe in this session ([/tmp/probe_recursive_union.osp](../../tmp/probe_recursive_union.osp), [/tmp/probe_closure_capture.osp](../../tmp/probe_closure_capture.osp), [/tmp/probe_lambda_simple.osp](../../tmp/probe_lambda_simple.osp)) confirmed four primitives are broken or missing today:

| # | Primitive | Status | Plan |
|---|---|---|---|
| 1 | Lambdas at all (`fn(x) => x + 1`) | ❌ codegen broken: `undefined variable: x` | [`closures.md`](closures.md) |
| 2 | Recursive unions with `List<Self>` / `Map<K,Self>` payload | ✅ shipped (PR #67): all union-payload constructors compile; `List<Self>`, `Map<K,Self>`, mutual recursion, and `Tree`-style self-recursion all round-trip end-to-end under strict llc per [TYPE-UNION-REC] | — |
| 3 | Error message payload threading through `Result<T, E>` | ✅ Phase 1+2 shipped: every builtin returns its real per-call message via the new `err_msg` slot. Phase 5 (user-constructed `Error { message: ... }` for non-string success types) still open. | [`error-payloads.md`](error-payloads.md) |
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

1. ~~`error-payloads.md`~~ — Phase 1+2 shipped; every builtin's Error branch now reports its real per-call message. Phase 5 (user-defined `Error { message: ... }` with non-string success types) deferred — captured at the bottom of `error-payloads.md`.
2. **`closures.md`** second — small surface, but every higher-order primitive (`map`, `filter`, `fold`, parser combinators) depends on it.
3. ~~`recursive-union-payloads.md`~~ — shipped in PR #67; `JsonValue`, tree types, and ASTs now compile and round-trip per [TYPE-UNION-REC].
4. **`string-cursor.md`** fourth — adds the C primitives that the JSON parser will sit on top of.
5. **`list-patterns.md`** fifth — wraps the parser in idiomatic recursive descent (`[head, ...tail]`).

After all four land, the proof point is one tested example: `examples/tested/json/json_parser.osp` — a JSON parser **written in Osprey** that consumes a real input and produces a `JsonValue`. That example becomes the canary for every future regression in these four areas.

## Master TODO (across all five plans)

- [x] Land `error-payloads.md` Phase 1 (codegen-side threading) and Phase 2 (codegen rewrite). Phase 5 (user-constructed Error with non-string success types) still open.
- [ ] Land `closures.md` Phase 1 (lambda param binding fix) and Phase 2 (capture).
- [x] Land recursive-union-payloads Phase 1–5 (layout, codegen, mutual recursion, deep-tree stress, strict-llc round-trip fixes); see PR #67. Plan file deleted now that all sub-items are checked.
- [ ] Land `string-cursor.md` Phase 1 (C runtime) and Phase 2 (builtins + registry).
- [ ] Land `list-patterns.md` Phase 1–3 (grammar, codegen, `osprey_list_drop` runtime).
- [ ] Land `examples/tested/json/json_parser.osp` written in pure Osprey, using only the above primitives plus existing `List`/`Map`/`match`. Must round-trip RFC 8259 conforming inputs.
- [ ] Update `coverage-to-90-percent.md` ratchet plan once the JSON example's tests are exercising the new branches.
