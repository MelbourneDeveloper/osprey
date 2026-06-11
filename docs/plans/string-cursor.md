# Plan: O(1) String Cursor (`byteAt`, `codePointAt`, friends)

Spec: [`0012-Built-InFunctions.md` — Cursor Access](../specs/0012-Built-InFunctions.md#cursor-access-total-o1).

Parent: [`production-primitives.md`](production-primitives.md).

## Problem

Every existing string primitive allocates. `take(s, 1)`, `drop(s, 1)`, `substring(s, i, i+1)` each call `malloc` and copy bytes. Writing a parser this way is O(n²) at best — every advance copies the remaining input — and produces O(n log n) garbage that the allocator has to wade through.

No primitive today gives you "the byte at position i" or "the codepoint at position i" without allocating. A user-written JSON parser, query-string parser, CSV reader, or markdown reader is therefore quadratic or worse, before any actual parsing logic.

A parser writer needs three things:

1. A way to look at one byte (or codepoint) at a known offset in O(1).
2. A way to know how far to advance after consuming a codepoint.
3. A way to build a small string (e.g., a JSON string-literal value) without a chain of `+`.

This plan adds the first two as builtins. The third is mostly handled today by `+` on `string` (already infallible) but with the cursor primitives in hand, parsers can buffer bytes into a `List<int>` and call `fromCodePoint` per char if they want.

## Scope

Five builtins, all spec'd in [`0012-Built-InFunctions.md`](../specs/0012-Built-InFunctions.md#cursor-access-total-o1):

| Builtin | Signature | Behaviour |
|---|---|---|
| `byteLength` | `(string) -> int` | UTF-8 byte length. O(1). |
| `byteAt` | `(string, int) -> Result<int, StringError>` | Byte at index `i` ∈ [0, 255]. O(1). |
| `codePointAt` | `(string, int) -> Result<int, StringError>` | UTF-8 decode at byte offset. O(1). |
| `codePointWidth` | `(int) -> Result<int, StringError>` | UTF-8 byte width of a codepoint. O(1). |
| `fromCodePoint` | `(int) -> Result<string, StringError>` | Encode one codepoint to a single-char string. O(1). |

## Why builtins (not pure Osprey)

The spec ([`0012` — Cursor Access](../specs/0012-Built-InFunctions.md#cursor-access-total-o1)) says: "they are the lowest-level string operations in the language; everything above is implementable in pure Osprey on top of them." That's literally true — you cannot implement `byteAt` in Osprey without already having a way to address bytes. So these five are the irreducible primitive layer; *everything else* (parsing, validation, encoding) is up for grabs as a user library.

This is the opposite of the JSON-as-builtin shortcut. Builtins here are the *minimum* C surface, chosen so that user code can do *more* in Osprey, not less.

## Phase 1 — C runtime

All five live in [`compiler/runtime/string_runtime.c`](../../compiler/runtime/string_runtime.c) and [`string_runtime.h`](../../compiler/runtime/string_runtime.h). The file is currently 279 LOC; adding five small functions keeps it well under the 500-LOC cap.

- [ ] **1.1** `int64_t osp_string_byte_length(const char *s)` — `s == NULL` → 0; otherwise `strlen(s)`. Wraps to `i64`.
- [ ] **1.2** `osp_result_int osp_string_byte_at(const char *s, int64_t i)` — bounds-check `i ∈ [0, byteLen)`; on success return `(unsigned char)s[i]` as `int64_t`; on failure populate the result's error message per [`error-payloads.md`](error-payloads.md) Phase 1.
- [ ] **1.3** `osp_result_int osp_string_codepoint_at(const char *s, int64_t byte_index)` — decode UTF-8 starting at `byte_index`. Validate: first byte's high bits select the width (0xxx → 1B, 110x → 2B, 1110 → 3B, 11110 → 4B); each continuation byte must be `10xx`; `byte_index + width <= byteLen`. Return the decoded scalar value. Error on any failure with a specific message (`"codePointAt: invalid UTF-8 lead byte"`, `"codePointAt: truncated codepoint"`, `"codePointAt: invalid continuation byte"`, `"codePointAt: index out of range"`).
- [ ] **1.4** `osp_result_int osp_string_codepoint_width(int64_t cp)` — return 1 if `cp <= 0x7F`; 2 if `cp <= 0x7FF`; 3 if `cp <= 0xFFFF` (excluding surrogates `0xD800..0xDFFF`); 4 if `cp <= 0x10FFFF`; otherwise Error.
- [ ] **1.5** `char *osp_string_from_codepoint(int64_t cp)` — encode `cp` to a new heap-allocated, null-terminated UTF-8 string (1–4 bytes plus terminator). Returns NULL on invalid scalar value (caller wraps as Error). Verify the existing memory model frees this (every other `osp_string_*` allocator returns malloc'd memory; assume same convention).
- [ ] **1.6** Add hard-assertion C unit tests in [`runtime/string_runtime_tests.c`](../../compiler/runtime/string_runtime_tests.c):
  - Empty string: `byteLength("") == 0`; `byteAt("", 0)` errors.
  - ASCII: every position works; index `-1` and `len` error.
  - Multi-byte: `byteAt("héllo", 1) == 0xC3` (first byte of `é`); `codePointAt("héllo", 1).value == 0xE9`; `codePointAt("héllo", 2)` errors (mid-codepoint).
  - 4-byte codepoint: `"𝄞"` (U+1D11E) — `byteLength == 4`; `codePointAt(s, 0).value == 0x1D11E`.
  - `codePointWidth(0x7F) == 1`; `codePointWidth(0x80) == 2`; `codePointWidth(0x10000) == 4`; `codePointWidth(0x110000)` errors; `codePointWidth(0xD800)` errors (surrogate).
  - `fromCodePoint(0x1F600)` (😀) returns a 4-byte string; `fromCodePoint(0x110000)` errors.
  - Round-trip: `codePointAt(fromCodePoint(cp).value, 0).value == cp` for `cp ∈ {0, 0x7F, 0x80, 0x7FF, 0x800, 0xFFFF, 0x10000, 0x10FFFF}`.
- [ ] **1.7** Run the C runtime test suite ([`compiler/runtime/run_tests.sh`](../../compiler/runtime/run_tests.sh)) and rebuild via `make _runtime`. The strict flags (`-Werror` plus the hardening set) must stay clean.

## Phase 2 — Builtin registry + LLVM codegen

- [ ] **2.1** Register the five signatures in the checker's builtin table ([`crates/osprey-types/src/builtins.rs`](../../crates/osprey-types/src/builtins.rs)). `byteLength` returns plain `int`, the others return `Result<_, StringError>` (but see Phase 3 caveat below).
- [ ] **2.2** Add the C-ABI lowerings in [`crates/osprey-codegen/src/extern_call.rs`](../../crates/osprey-codegen/src/extern_call.rs) — the same one-line `sig(...)` registry pattern as `length` and `contains`. Result-returning entries use the `ResultStr`-style return kinds already in that file.
- [ ] **2.3** No new module needed — each builtin is one registry row plus its C function; only add a module if a lowering ever needs more than the shared call path.

## Phase 3 — Result payload caveat

Today, the registry's "Result<T, StringError>" is mostly a lie — see [`error-payloads.md`](error-payloads.md). The error message slot is hardcoded `"Error occurred"` regardless of which builtin failed.

- [ ] **3.1** Implement the five generators against the **post-[`error-payloads.md`](error-payloads.md)-Phase-2** runtime contract from day one. If `error-payloads.md` hasn't shipped yet, sequence: ship Phase 1 of error-payloads first (runtime contract), then string-cursor Phase 1+2 on top of it. The cursor builtins are the first new fallible builtins to land; they should not perpetuate the hardcoded-message hack.

## Phase 4 — End-to-end tests

- [ ] **4.1** `examples/tested/basics/strings/cursor_basic.osp`:
  ```osprey
  let s = "héllo"
  print(byteLength(s))                          // 6
  match byteAt(s, 1) { Success { value } => print(value) }    // 195 (0xC3)
  match codePointAt(s, 1) { Success { value } => print(value) }  // 233 (0xE9)
  match codePointWidth(0x1F600) { Success { value } => print(value) }  // 4
  match fromCodePoint(0x1F600) { Success { value } => print(value) }   // 😀
  ```
  Pin output verbatim in `.expectedoutput`.
- [ ] **4.2** `examples/tested/basics/strings/cursor_advance.osp`: walk every codepoint of `"héllo, 世界"` using `codePointAt` + `codePointWidth` in a recursive fn (no closures needed); print each scalar value. Asserts that the cursor really does run in O(n) of input bytes.
- [ ] **4.3** Negative example `examples/failscompilation/cursor_wrong_arity.ospo` — `byteAt("a")` (missing index) — expected error names the registered signature.

## Phase 5 — The canary

Once all of `error-payloads.md`, `closures.md`, `recursive-union-payloads.md`, and this plan ship, the JSON-parser canary in [`production-primitives.md`](production-primitives.md) is the cross-cutting acceptance test.

- [ ] **5.1** Land `examples/tested/json/json_parser.osp` written in pure Osprey:
  - Uses `JsonValue` (recursive-union-payloads) as its output type.
  - Uses `codePointAt`/`codePointWidth` (this plan) to walk the input.
  - Uses closures (closures.md) for parser-combinator style helpers OR plain top-level functions (your call; both must work).
  - Uses real error messages (`error-payloads.md`) reporting `"unexpected '}' at byte 47"` style failures.
  - Round-trips at least `{"name": "alice", "age": 30, "tags": [1, true, null, "x"]}`.
  - **Asserts**: the parser file is under 500 LOC. If it isn't, we missed a primitive — open a follow-up issue rather than splitting the file.

## Out of scope

- A `Char` type (deferred).
- Grapheme-cluster iteration.
- Unicode normalisation forms.
- Regex.
- `String.format` / printf-style.

## TODO checklist

### Phase 1 — C runtime
- [ ] 1.1 `osp_string_byte_length`
- [ ] 1.2 `osp_string_byte_at`
- [ ] 1.3 `osp_string_codepoint_at` with UTF-8 validation
- [ ] 1.4 `osp_string_codepoint_width` with scalar-validity check
- [ ] 1.5 `osp_string_from_codepoint`
- [ ] 1.6 C unit tests covering ASCII, multi-byte, 4-byte, surrogates, boundaries
- [ ] 1.7 `runtime/run_tests.sh` + `make _runtime` clean under strict flags

### Phase 2 — Builtin registry + codegen
- [ ] 2.1 Checker signatures in `crates/osprey-types/src/builtins.rs`
- [ ] 2.2 C-ABI lowerings in `crates/osprey-codegen/src/extern_call.rs`

### Phase 3 — Coordinate with error-payloads
- [ ] 3.1 Land on top of `error-payloads.md` Phase 1 runtime contract — do not perpetuate the static-message hack

### Phase 4 — E2E tests
- [ ] 4.1 `cursor_basic.osp` happy-path
- [ ] 4.2 `cursor_advance.osp` codepoint walk
- [ ] 4.3 `cursor_wrong_arity.ospo` negative

### Phase 5 — Canary
- [ ] 5.1 `examples/tested/json/json_parser.osp` (pure Osprey) parses `{"name":"alice","age":30,"tags":[1,true,null,"x"]}` and asserts the parsed tree; ≤500 LOC

### Acceptance
- [ ] User-written JSON parser passes a 20-input round-trip corpus.
- [ ] No O(n²) string copying in the JSON parser source (no `take`/`drop`/`substring` on the input cursor).
