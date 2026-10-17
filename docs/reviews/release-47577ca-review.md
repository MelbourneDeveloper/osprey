# Review of work since Osprey release 47577ca

Date: 17 October 2026.

## Scope and outcome

Reviewed the net changes from `47577ca1a9e5264e8642b36bf97747f53bb25fb4` through `98d2fa6e3cd01b408504616e8050995ff0ce355b`: 116 reachable commits and 397 changed files. This includes compiler/type/code generation changes, the Python frontend and interop package, runtime extraction, GPT-2 handlers, build/release configuration, editor changes, website changes, and test changes. Much of the diff is repository/specification reorganization.

The review began at `991bbe809120a07081a6e796d2ffe54c2e529ee7`. Commit `98d2fa6e` arrived during the review and was included. Uncommitted work continued changing in the shared checkout; it is not part of the committed review boundary. Findings were reproduced using compilers built from separate, exported commit snapshots.

**Two defects in the newly added Python frontend remain at the reviewed head. A high-impact GPU miscompile also remains, but reproduces at the release baseline and must not be called a new regression. No weakening of existing CI gates, coverage thresholds, or corpus floors was found.**

| ID | Priority | Finding | Classification |
| --- | --- | --- | --- |
| R1 | P2 | Invalid Python dedents are accepted and executable statements move outside their apparent block | New frontend defect |
| R2 | P2 | Unicode escapes producing adjacent braces lose characters in f-strings | New frontend defect |
| B1 | P1 | Device kernels resolve a named helper's global constant to a caller's shadowing parameter | Pre-existing; reproduced at both release and reviewed head |
| H1 | P2 | `loadTensor` was rejected inside a static handler | Found at initial head; fixed by `98d2fa6e` |

P1 denotes a high-impact correctness issue; P2 denotes a correctness issue requiring a fix in the affected feature. No implementation fixes were applied as part of this review.

## R1 — [P2] Validate dedents independently of tree-sitter recovery

Source: [`crates/osprey-syntax/src/python/mod.rs:77–91`](https://github.com/Nimblesite/osprey/blob/98d2fa6e3cd01b408504616e8050995ff0ce355b/crates/osprey-syntax/src/python/mod.rs#L77-L91).

The indentation pass checks tabs but does not verify that each dedent returns to an existing indentation level. The Python tree-sitter grammar accepts the following text without an error/missing node, so the later CST error scan does not catch it:

```python
if False:
    print("hidden")
  print("leaked")
```

Reproduction, with those exact spaces saved as `/tmp/bad-dedent.py`:

```sh
osprey /tmp/bad-dedent.py --run --quiet
```

Observed at the reviewed head: exit status **0**, stdout **`leaked\n`**. The two-space statement is lowered outside the conditional. The declared Python-flavor contract requires inconsistent dedents to be errors; CPython rejects this same text with `IndentationError`.

This can turn a layout mistake into an unconditional operation instead of a compilation failure.

**Test gap:** [`python_reports_malformed_indentation_and_incomplete_cst_nodes`](https://github.com/Nimblesite/osprey/blob/98d2fa6e3cd01b408504616e8050995ff0ce355b/crates/osprey-syntax/tests/python.rs#L405-L414) only asserts that *some* error exists. Its inconsistent-dedent fixture puts `return` at the recovered outer level, where `return` is unsupported anyway. That fixture can therefore pass without proving indentation validation works.

Recommended regression: reject the example above during parsing and assert an indentation-specific diagnostic at line 3. Also cover nested blocks and preserve valid multiline strings and continuations.

## R2 — [P2] Collapse source brace escapes before decoding character escapes

Source: [`crates/osprey-syntax/src/python/strings.rs:41–52`](https://github.com/Nimblesite/osprey/blob/98d2fa6e3cd01b408504616e8050995ff0ce355b/crates/osprey-syntax/src/python/strings.rs#L41-L52).

`content()` decodes backslash escapes first, then replaces every decoded `{{` or `}}` with one brace. This conflates two escaped brace characters with the f-string syntax for one literal brace.

Save as `/tmp/escaped-braces.py`:

```python
print(f"\u007b\u007b")
```

```sh
osprey /tmp/escaped-braces.py --run --quiet
```

Observed at the reviewed head: exit status **0**, stdout **`{\n`** (`7b 0a`). Expected: **`{{\n`** (`7b 7b 0a`), as produced by the same supported Unicode escapes in a normal string and by Python's f-string syntax. The corresponding closing-brace escapes have the same problem.

The compiler silently changes string data. This affects generated text, templates, and protocols that use brace sequences.

Recommended regression: assert literal output for adjacent `\u007b`, `\u007d`, `\x7b`, and source `{{`/`}}`, both beside interpolations and without them. Brace syntax must be distinguished in the source/CST before Unicode decoding creates new brace characters.

## B1 — [P1, pre-existing] Named GPU helpers acquire caller bindings

Affected code: [`crates/osprey-codegen/src/gpu_device/emit.rs:274–295`](https://github.com/Nimblesite/osprey/blob/98d2fa6e3cd01b408504616e8050995ff0ce355b/crates/osprey-codegen/src/gpu_device/emit.rs#L274-L295) and [`343–361`](https://github.com/Nimblesite/osprey/blob/98d2fa6e3cd01b408504616e8050995ff0ce355b/crates/osprey-codegen/src/gpu_device/emit.rs#L343-L361).

```osprey
let bias = 10
fn add(x) = wrapAdd(x, bias)
fn apply(f, x) = f(x)
fn main() = {
    let output = gpuMap(toGpu([1, 2]), fn(bias) => apply(add, bias))
    print("${gpuGet(output, 0) ?: 0}:${gpuGet(output, 1) ?: 0}")
}
```

Save as `/tmp/gpu-scope.osp` and run:

```sh
OSPREY_GPU=host osprey /tmp/gpu-scope.osp --run --quiet
OSPREY_GPU=metal osprey /tmp/gpu-scope.osp --run --quiet
OSPREY_GPU=opencl osprey /tmp/gpu-scope.osp --run --quiet
```

| Compiler snapshot | Backend | Output | Exit status |
| --- | --- | --- | --- |
| Release `47577ca1` | Host | `11:12` | 0 |
| Release `47577ca1` | Metal | `2:4` | 0 |
| Reviewed head `98d2fa6e` | Host | `11:12` | 0 |
| Reviewed head `98d2fa6e` | Metal | `2:4` | 0 |
| Reviewed head `98d2fa6e` | OpenCL | `2:4` | 0 |

Hardware runs reported completed kernel dispatches on an Apple M4 Max. These are actual device results, not host fallback.

`add` must read the top-level `bias = 10`. Device helper emission instead makes caller scopes available to the helper; `passed()` also records the current scopes for a named top-level function. The kernel parameter `bias` then supplies the helper's free variable, producing `x + x`. Named functions need their definition's lexical environment; inline lambdas need their actual capture environment.

**This is not a regression introduced since the release.** It is included because the range modifies this compiler path and adds higher-order device support while the existing correctness hole remains.

**Test gap:** the new [`kernels_reducing_through_higher_order_helpers_keep_their_device_form`](https://github.com/Nimblesite/osprey/blob/98d2fa6e3cd01b408504616e8050995ff0ce355b/crates/osprey-cli/tests/cross_flavor_ir_equiv/device.rs#L169-L203) checks launch presence, a generated loop, and equal device source across flavors. It does not execute that fixture or assert its numeric results. Equal generated code across flavors cannot detect a bug shared by all flavors.

Recommended regression: compare numeric host/device results when a kernel parameter shadows a helper's global constant, and when helper parameters shadow lambda captures. Follow the repository's broken-code quarantine policy until lexical binding identity is preserved.

## H1 — Static tensor loading was fixed during the review

At initial head `991bbe80`, `loadTensor` was registered inside `runtime_group(files)`. A static handler using the new host loader failed with:

```text
static handler arm `Read.load` requires runtime builtins: loadTensor
```

The reproducer was:

```osprey
static effect Read { load: fn(string, string) -> Result<List<float>, Error> }
fn main() = {
    handle static Read { load p n => loadTensor(p, n) }
    let result = perform Read.load("m", "w")
    print(0)
}
```

Commit `98d2fa6e` moves the loader into a separate declaration group and switches the GPT-2 CPU handlers to the host loader and host abort function. The same source passes `--check` with a compiler built from that commit. This issue is **resolved within the final reviewed range**, not an outstanding finding.

## Test weakening and coverage assessment

- Existing per-project coverage thresholds were not reduced. Two new C runtime entries, `tensor_file` and `tensor_list`, have 90% floors.
- The duplication ceiling tightened from 6.1% to 3.03%.
- Native, wasm, mobile, and alternate-GPU-lowering corpus floors were not lowered. Native worker isolation still preserves output and exit status for evaluation.
- The removed inline GPU loader assertions were moved into `tensor_file_tests.h` and invoked through the parameterized loader suite. They were not simply deleted. The C runtime suites passed, including the GPU and host-list tensor loaders.
- The TextMate corpus sweep was extracted into its own test/harness and gained additional invariants; the existing keyword/operator checks remain.
- The Makefile split preserves the substantive test and coverage recipes. Release classification now includes `scripts/make/**`, and the Node dependency guard follows included fragments.
- No newly narrowed existing required job or newly weakened existing assertion was identified. The weaknesses above are gaps in new/retained tests, not evidence that an existing failing assertion was removed to obtain a pass.
- The new full-checkpoint GPT-2 parity test is ignored by default. The ordinary differential corpus discovers only `*.test.osp` and `*.test.ospml`; its 223 passing goldens do not establish equivalent corpus coverage for Python-flavor sources.

## Verification and limits

| Check | Result |
| --- | --- |
| `cargo test --workspace --exclude osprey-cli --quiet` | Passed |
| CLI integration targets `python_interop`, `cross_flavor_equiv`, `operation_modes` | Passed: 17 tests; one checkpoint-dependent test ignored |
| `make _test_c_runtime` | Passed, 30 suites |
| Native default-memory differential corpus | Passed: 223 programs/goldens, 20 alternate GPU lowerings, 6 doctests |
| GC differential corpus | Passed: 223 programs/goldens, 20 alternate GPU lowerings, 6 doctests |
| ARC differential corpus | Passed: 223 programs/goldens, 20 alternate GPU lowerings, 6 doctests; zero leaky programs |
| CI scope classifier tests | Passed: 16 tests |
| Node dependency guard | Passed: 10 guarded recipes |
| Mobile/build-tool script tests | Passed: 8 tests |
| Isolated compiler builds for release, initial head, and final reviewed head | Passed; used to distinguish new, pre-existing, and resolved defects |
| Direct frontend and real-device reproducers | Confirmed the findings above |

Broad suites ran in the shared checkout while unrelated edits were occurring. The isolated snapshot builds provide the commit-specific evidence for the findings and for H1's resolution. The GPU reproducer uses runtime facilities whose implementation did not change in this range, apart from a documentation reference.

This was a source review with targeted execution, not a full CI certification. The full checkpoint test, wasm/iOS/Android execution, installed-VSIX/browser suites, and complete lint/coverage gates were not run as part of this review.
