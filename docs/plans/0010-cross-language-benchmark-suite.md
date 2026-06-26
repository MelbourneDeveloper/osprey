# Plan 0010 — Cross-Language Benchmark Suite

**Subsystem:** `benchmarks/` (harness + cases), `Makefile` (`make bench`),
`.devcontainer` (comparison toolchains)
**Status:** Suite shipped (15 cases × 5 languages); feature-blocked classics pending
**Spec ID:** `[BENCH-SUITE]`

## Summary

An **accurate, reproducible** way to see where Osprey sits on CPU time and peak
memory against **Rust, C, OCaml, and Haskell**. Every benchmark is implemented
in all five languages with the *same naive algorithm and parameters*, compiled
to a native binary, checked byte-for-byte against an integer oracle
(`expected.txt`), then timed with `hyperfine` (CPU) and `/usr/bin/time` (peak
RSS). All source lives **in-tree and version-controlled** under
`benchmarks/cases/<name>/` — `<name>.{osp,rs,c,ml,hs}` + `expected.txt` +
`bench.json`. Only build/run *output* (`benchmarks/results/`) is gitignored.

## What works today (15 cases)

**Recursion-bound:** `fib`, `ackermann`, `tak`, `hanoi`, `pascal`, `coins`, `mutual`
**Iteration / number theory:** `primes`, `gcdsum`, `nestedloop`, `factorial`, `powmod`, `josephus`, `coprime`
**Allocation / memory:** `binarytrees`

Harness: [benchmarks/run.sh](../../benchmarks/run.sh) (toolchain detection,
build-once, correctness oracle, CPU + memory measurement) →
[benchmarks/report.py](../../benchmarks/report.py) (CPU table, relative-speed
table, peak-memory table, geomean Osprey-slowdown). `make bench` /
`BENCH_FILTER=<name> make bench`. Dev container installs `ghc ocaml time` +
hyperfine.

## Key finding (the headline)

Osprey peak RSS scales with **operation count**, not live data — the runtime
heap-allocates per checked op/call and never frees during a run (no GC, no
scope-based free): `1+1` ≈ 1.4 MB, `fib(25)` ≈ 13 MB, `fib(35)` ≈ 1.4 GB. This
suite now tracks that regression directly. See
[spec 0018 — Memory Management](../specs/0018-MemoryManagement.md).

## Blocked classics — need Osprey language features (not faked)

Each row is a benchmark we *cannot* express today. Ordered by leverage: the
feature that unblocks the most classics, with the least scope, first.

| Missing feature | Unblocks | Scope |
|-----------------|----------|-------|
| **Integer division `/` (+ already-present `%`)** | collatz, digit-sum, integer-sqrt, sieve-of-eratosthenes (index math), radix benchmarks | **Low** — new typed op, no new types |
| Mutable arrays / fixed-size buffers | sieve, matrix-multiply, quicksort, mergesort, fannkuch, n-queens | High — new aggregate type + codegen |
| `int`↔`float` conversion + `sqrt`/trig stdlib | mandelbrot, n-body, spectral-norm | Medium — math runtime + exact float oracle |
| List literals / cons / list pattern-matching ops | list-heavy classics already partly covered by `[h, ...t]` | Medium |
| Arbitrary-precision integers | pidigits, big-factorial (exact) | High |

**Decision:** integer division is the next feature to add — lowest scope,
unblocks the most classic integer benchmarks, and `%` already exists so the
codegen/runtime path is half-built.

## Implementation plan (next feature: integer division)

1. Find how `%` is lowered (`crates/osprey-codegen`) and the float `/` path; add
   integer `/` as a typed, overflow/zero-checked op returning
   `Result<int, MathError>` (divide-by-zero → `MathError`), mirroring `%`.
2. Type rule in `crates/osprey-types`: `int / int -> Result<int, MathError>`;
   keep `float / float -> float` unchanged (type-directed dispatch).
3. Whole-body arithmetic helpers already auto-unwrap `Result` — confirm `/`
   composes (e.g. `fn idiv(a: int, b: int) -> int = a / b`).
4. `find-similar` before adding any runtime helper; reuse the `%` path.
5. Add `collatz`, `digit-sum`, `integer-sqrt` benchmark cases (all 5 languages).
6. `make ci` green; refresh `benchmarks/results/`.

## TODO

- [x] Harness: build-once, correctness oracle, CPU (hyperfine) + peak RSS.
- [x] `report.py`: CPU + relative-speed + peak-memory tables, geomean summary.
- [x] 15 cases × 5 languages, all source version-controlled under `cases/`.
- [x] `make bench` target + `BENCH_FILTER`; `.gitignore` tracks source, ignores `results/`.
- [x] Dev container: `ghc`, `ocaml`, `time`, hyperfine.
- [x] README documents methodology, fairness caveats, the memory finding.
- [ ] Run full suite end-to-end → publish `results.md` numbers in README findings.
- [ ] **Add integer division `/`** (typed, checked) — codegen + types + tests.
- [ ] Add `collatz`, `digit-sum`, `integer-sqrt` cases once `/` lands.
- [ ] (Later) mutable arrays → sieve, matrix-multiply, sort/fannkuch/n-queens.
- [ ] (Later) `int`↔`float` + `sqrt` → mandelbrot, n-body, spectral-norm.
