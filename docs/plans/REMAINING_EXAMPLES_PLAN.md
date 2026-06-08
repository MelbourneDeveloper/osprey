# Remaining `examples/tested` — Fix Plan

Status at time of writing: **24 / 35** byte-exact (after whole-string `TrimSpace`)
against the sibling `.expectedoutput`, verified by `crates/diff_examples.sh`.
The Rust workspace is fully green: `cargo build --release`, `cargo clippy
--workspace --all-targets` (max strictness), `cargo fmt --check`, and
`cargo test --workspace` all pass.

The 11 remaining failures fall into **4 subsystems**. None is a one-line tweak;
each needs a genuinely new capability. They are ordered below by ROI
(examples-unlocked ÷ effort), most attractive first.

---

## 1. Generic monomorphization — unlocks 3

**Examples:** `types/any_type_comprehensive`, `types/pure_hindley_milner_test`,
`types/type_equality_comprehensive`

**Symptom:** all three abort codegen with
`invalid program: expected an integer, found a string/handle`.

**Root cause.** A generic record field (`type Generic<T> = { data: T }`) has its
field type written as the type variable `T`, which `types::ltype_of_name` maps to
`LType::Ptr`. So `gen1 = makeGenericInt(42)` stores `42` as `inttoptr 42` and
`gen3 = makeGenericString("x")` stores the string pointer — both as `i8*`. At a
use site (`${gen1.data}`, `gen1.data * 2`) the backend cannot tell an int payload
from a string payload, so `as_i64`/`toString` on the `Ptr` fails or prints
garbage. There is **no per-expression type information**: `osprey_types::ProgramTypes`
publishes only `functions: HashMap<name,(params,ret)>`, `ctors`, and `unions` —
frozen, substitution-free tables — so the backend cannot recover that *this*
`gen1.data` is `int` while *that* `gen3.data` is `string`.

**Fix approach (pick one):**

- **(A) Publish a per-expression type map.** Extend `osprey_types::ProgramTypes`
  with `expr_types: HashMap<ExprId, Type>` (or key by `Position`). During
  inference (`crates/osprey-types/src/expr.rs::infer_expr`), record the resolved
  type of every node after the final substitution is applied. The backend then
  asks `prog.expr_type(&expr)` in `gen_field_access` / `gen_object` /
  `gen_index` to pick the concrete `LType` instead of defaulting `T → Ptr`.
  Requires giving AST nodes stable ids (add an `id: u32` to `Expr`, or thread a
  `Position`-keyed side table). **Largest, cleanest, most reusable.**

- **(B) Monomorphize by call site.** When lowering `let gen1 = makeGenericInt(42)`,
  read `makeGenericInt`'s *concrete resolved return type* (`Generic<int>`) from
  `prog.return_type("makeGenericInt")` and attach the resolved field layout to
  the binding's `Value` via a per-variable layout override (a new
  `HashMap<owner_instance, Vec<(field,LType)>>` on `Codegen`). Cheaper than (A)
  but only covers values that flow straight from a typed function return; it will
  not handle a generic field read inside a polymorphic body.

- **(C) Uniform boxed representation + runtime type tag.** Store every generic
  field as an `{ i8 tag, i64 payload }` and have `toString`/`as_i64` branch on the
  tag at runtime. Avoids the inference plumbing but adds a tag word to every
  generic value and a runtime dispatch — and changes the ABI for these types.

**Recommendation:** (A). It is the only approach that fixes the *general* case
(`fn getFirst(p) = p.first` already works via the by-name fallback in
`aggregate::find_field_owner`, but generic *value* typing needs real per-expr
types). Effort: ~1–2 days incl. inference changes + AST ids.

**Verify:** the three `types/*` examples; watch for regressions in every example
that uses records/generics (`pure_hindley_milner` exercises Box/Pair2/Generic).

---

## 2. Runtime FFI symbols — unlocks `files` (1), foundation for `db`/`http`

**Example:** `files/file_io_json_workflow`

**Symptom:** clang link error — `Undefined symbols for architecture arm64:
"_readFile", "_writeFile"` (and likely `parseJson`/JSON helpers).

**Root cause.** The program calls builtins `readFile`/`writeFile` (and JSON ops);
the backend emits `call @readFile(...)` but the linked static runtimes
(`compiler/bin/libfiber_runtime.a`, `libhttp_runtime.a`) export no such symbol.
The backend already declares unknown callees as externs (`call_with_values`), so
codegen succeeds — it is purely a **missing runtime symbol** at link time.

**Fix approach:**

1. Grep the C runtime for the real symbol names
   (`grep -rn "readFile\|read_file\|osp.*file" compiler/runtime/`). They likely
   exist under a different spelling (e.g. `osp_read_file`).
2. Either (a) add a name-mapping in the backend so `readFile` → the real C symbol
   (a small table in `runtime.rs`/`call_with_values`), or (b) add thin C shims
   named `readFile`/`writeFile` to a runtime `.c` and rebuild the static libs, or
   (c) implement them in `crates/osprey-runtime-sys` if the Rust runtime is the
   intended home.
3. Confirm `osprey-cli`'s `link_args` includes the runtime that exports them.

**Risk:** file ops are deterministic if the example reads/writes fixed paths;
check the `.expectedoutput` does not embed absolute paths or timestamps. Effort:
~half a day once the symbol names are located.

---

## 3. Function-pointers-as-values + process runtime — unlocks `processes` (1)

**Example:** `processes/async_process_management`

**Symptom:** `unknown name processEventHandler`.

**Root cause (two parts):**

- **Higher-order values.** `spawnProcess("echo …", processEventHandler)` passes a
  *named function* as an argument value. `gen_expr(Identifier("processEventHandler"))`
  finds no local binding and no constructor → `unknown name`. The backend only
  supports inline/let-bound lambdas (beta-reduced) and direct calls; it has no
  notion of a first-class function pointer.
- **Process runtime.** Even once the callback is passed, `spawnProcess` must exist
  in the runtime, run the command, and invoke the callback — deterministically.

**Fix approach:**

1. In `gen_expr`'s `Identifier` arm, if the name is a known top-level function
   (`cg.fn_params.contains_key(name)` / inference has a `Fun` type for it), emit a
   function pointer: `Value` carrying `bitcast <sig>* @name to i8*` (a new `Ptr`
   value tagged as a code pointer). Pass it as `i8*`.
2. Teach `call_with_values` / a new indirect-call path to accept an `i8*` callee
   and `bitcast` it back to the right fn-ptr type before `call` — mirror what
   `effects::gen_perform` already does for handler pointers.
3. Ensure `spawnProcess` (and the callback ABI) exists in the runtime and is
   deterministic. The example uses `echo`/`false`; output ordering across spawned
   processes must match Go's.

**Risk:** medium-high — touches the call ABI. Do (1)+(2) behind the existing
`gen_user_call` so direct calls are unaffected; the `effects` module is the
template. Effort: ~1–2 days incl. runtime.

---

## 4. HTTP server/client runtime + `HttpResponse` builtin — unlocks 3

**Examples:** `http/http_response_handle`, `http/http_server_example`,
`http/tui_repo_table`

**Symptom:** `invalid program: missing field 'body' for 'HttpResponse'` when
lowering `HttpResponse { status, headers, contentType, streamFd, isComplete,
partialBody }`.

**Root cause.** `HttpResponse` is a *builtin* type (not declared in the `.osp`).
The registered ctor layout (in `osprey-types` builtins / the C runtime contract)
expects a `body` field, but the literal supplies `partialBody` — the two
definitions are out of sync. Beyond that, these examples need a **live HTTP
server + client** (`httpCreateServer`, `httpListen`, `httpGetResponse`, handle
lifecycle) producing byte-exact output (`status=200`, `body=hello body`, header
lookups, double-free → `Error`).

**Fix approach:**

1. Reconcile the `HttpResponse` field set: make the registered builtin layout
   (wherever `HttpResponse`'s ctor is registered for the Rust backend) match the
   field names the examples use (`status, headers, contentType, streamFd,
   isComplete, partialBody`) — or map them. Cross-check `compiler/runtime/`'s
   HTTP struct.
2. Link `libhttp_runtime.a` and wire the `http*` builtins to its symbols
   (`httpCreateServer`, `httpListen`, `httpGetResponse`, `httpResponseStatus/
   Body/Header/Free`, `httpCreateClient`, `httpCloseClient`, `httpStopServer`).
3. Drive the lifecycle deterministically: the examples bind to `127.0.0.1:18095`
   and expect a fixed status/body, so the runtime must serve the registered
   handler (`handleRequest`) — which itself depends on **function-pointers-as-
   values** (§3, `httpListen(serverId, handleRequest)`). **Do §3 first.**

**Risk:** highest — real sockets, handler dispatch, port binding, and
deterministic ordering. Effort: ~2–3 days; depends on §3.

---

## 5. Go parser precedence bug — `comprehensive_math` (1), do last / maybe never

**Example:** `math/comprehensive_math`

**Symptom:** one line differs — `complex = 18` (ours, arithmetically correct) vs
`complex = 2` (expected).

**Root cause.** `fn complex(a, b) = match (a*2) + (b*3) - 1 { Success{value} =>
value … }` with `a=5, b=3`. The correct value is `10 + 9 - 1 = 18`; Go's
`.expectedoutput` is `2`, produced by a **precedence/grouping bug in the Go
ANTLR parser**. Our tree-sitter grammar parses it correctly.

**Fix approach.** Two unpleasant options:

- **(A) Replicate the bug** in the tree-sitter grammar / lowering for this exact
  shape. **Strongly discouraged:** it would corrupt *correct* arithmetic
  precedence everywhere and likely break several of the 24 passing examples.
- **(B) Regenerate the oracle.** If the project's intent is "match the *new* Rust
  compiler," update `comprehensive_math.osp.expectedoutput` to `complex = 18`
  (the correct value) and treat the Go file as the stale artifact. This requires
  an explicit product decision: are the `.expectedoutput` files frozen Go output
  (bug-for-bug) or the spec of correct behavior?

**Recommendation:** raise with the team. Do **not** silently change either side.
If the contract is "bug-for-bug Go parity," this example stays red until (A) is
deemed worth the global risk (it almost certainly is not).

---

## Dependency order

```
§2 files (independent, easy)
§1 types×3 (independent, medium-large)
§3 processes  ──┐
                ├─►  §4 http×3   (http handler dispatch needs first-class fn ptrs)
§3 fn-ptrs  ────┘
§5 comprehensive_math  (product decision, do last)
§? db×2 (sqlite)  — same shape as §2: locate/link the sqlite FFI symbols;
                    sqlite_basics currently produces empty output (runtime no-op).
```

---

## TODO

- [ ] **types×3 — generic monomorphization.** Add a per-expression resolved-type
      map to `osprey_types::ProgramTypes` (give `Expr` stable ids or key by
      `Position`); record types in `infer_expr` after final substitution; consume
      in `gen_field_access`/`gen_object`/`gen_index` to pick the concrete `LType`
      instead of `T → Ptr`. Verify `any_type_comprehensive`,
      `pure_hindley_milner_test`, `type_equality_comprehensive`; watch
      record/generic regressions.
- [ ] **files — runtime symbols.** Locate real `readFile`/`writeFile`/JSON
      symbols in `compiler/runtime/`; add a backend name-map *or* C shims; ensure
      `osprey-cli` links the exporting lib. Verify `file_io_json_workflow` (check
      `.expectedoutput` has no absolute paths/timestamps).
- [ ] **db×2 — sqlite FFI.** Same recipe as files: find/link the sqlite runtime
      symbols so `sqlite_basics` stops producing empty output; then
      `database_effect`. Verify both.
- [ ] **processes — first-class function pointers.** Emit a code-pointer `Value`
      for a named-function identifier; add an indirect-call path
      (`bitcast i8* → sig* → call`) modeled on `effects::gen_perform`; keep
      `gen_user_call` direct calls unchanged. Then wire/verify the `spawnProcess`
      runtime + callback ABI. Verify `async_process_management`.
- [ ] **http×3 — HttpResponse layout + HTTP runtime.** Reconcile the builtin
      `HttpResponse` field set (`…partialBody`) with the examples; link
      `libhttp_runtime.a` and wire the `http*` builtins; rely on first-class fn
      ptrs for `httpListen(serverId, handleRequest)`. Verify
      `http_response_handle`, `http_server_example`, `tui_repo_table`.
- [ ] **comprehensive_math — product decision.** Decide whether
      `.expectedoutput` is bug-for-bug Go parity or correct-behavior spec. If the
      latter, update the expected file to `complex = 18`. Do **not** replicate the
      Go precedence bug globally. Verify `comprehensive_math`.
- [ ] **Regression gate.** After each item: `cargo build --release`,
      `cargo clippy --workspace --all-targets`, `cargo fmt --check`,
      `cargo test --workspace`, and `crates/diff_examples.sh` must all stay green
      and the pass count must not drop.
```
