# Osprey benchmarks

A cross-language performance harness that measures where Osprey sits relative to
**Rust, C, OCaml, and Haskell** on classic compute benchmarks — both **CPU time**
and **peak memory**.

Every benchmark is implemented identically in all five languages, compiled to a
native binary, checked for correct output, then timed. **All source is in this
folder** under `cases/<name>/` so you can read and compare every line:

```
benchmarks/cases/<name>/
  <name>.osp   <name>.rs   <name>.c   <name>.ml   <name>.hs   expected.txt   bench.json
```

```bash
make bench                       # build everything, run the whole suite
BENCH_FILTER=fib make bench      # only cases whose name contains "fib"
zsh benchmarks/run.sh            # run directly (assumes `make build` already ran)
zsh benchmarks/run.sh primes     # direct, single case
```

> **Heads-up on RAM.** With the optimized build (below) every case except
> `binarytrees` peaks at ~1.4 MB — on par with C. `binarytrees` still peaks near
> **900 MB** because its tree nodes genuinely escape and the default allocator
> does not reclaim them yet (see *Findings*); run it on a machine with a few GB
> free, or skip it with `BENCH_FILTER`.

Results are written to `benchmarks/results/` (gitignored):

| file | contents |
|------|----------|
| `results.md`   | human-readable CPU + memory tables and Osprey-relative slowdown |
| `results.json` | the same data, structured, for tracking over time |
| `hf/*.json`    | raw [hyperfine](https://github.com/sharkdp/hyperfine) exports per case |

## The benchmarks (18)

All have a single deterministic **integer** result, so output is byte-comparable
across languages (a broken implementation is caught and excluded from timing).

**Recursion-bound**

| Case | What it stresses | Workload |
|------|------------------|----------|
| `fib`       | function-call + recursion overhead   | naive recursive `fib(35)` |
| `ackermann` | deep non-tail recursion              | `ack(3, 10)` |
| `tak`       | heavy ternary self-recursion         | `tak(32, 16, 8)` |
| `hanoi`     | exponential double recursion         | Towers of Hanoi move count, n=25 |
| `pascal`    | un-memoised binomial recursion       | `C(27, 13)` via `C(n-1,k-1)+C(n-1,k)` |
| `coins`     | combinatorial tree recursion (SICP)  | ways to make 600 from `[1,5,10,25,50]` |
| `mutual`    | mutual recursion                     | `isEven`/`isOdd` across a range |

**Iteration / number theory**

| Case | What it stresses | Workload |
|------|------------------|----------|
| `primes`     | integer `%` in a hot loop          | count primes below 200000 (trial division) |
| `gcdsum`     | Euclidean recursion, modulo        | sum `gcd(i, 1234567)`, i in 1..1,999,999 |
| `nestedloop` | nested iteration + arithmetic      | triple loop `250³`, accumulate `(i*j*k) mod 1e9+7` |
| `factorial`  | multiplication-heavy fold          | product `1..10,000,000` mod 1e9+7 |
| `powmod`     | naive modular exponentiation       | sum of `i^20 mod 1e9+7`, i in 1..1,000,000 |
| `josephus`   | modular iteration                  | Josephus survivor, n=10,000,000, k=7 |
| `coprime`    | nested iteration + gcd             | count coprime pairs in a 2000×2000 grid |
| `collatz`    | integer `/` in deep recursion      | sum of Collatz (3n+1) stopping times over 1..100,000 |
| `digitsum`   | integer `/` + `%` in recursion     | sum of decimal digit-sums over 1..2,000,000 |
| `isqrt`      | Newton's method, integer `/`       | sum of integer square roots over 1..1,000,000 |

**Allocation / memory**

| Case | What it stresses | Workload |
|------|------------------|----------|
| `binarytrees` | allocation / GC / memory pressure | build & checksum 1200 trees of depth 13 |

`binarytrees` is the primary **memory** benchmark: it churns millions of small
heap nodes, so peak RSS reveals each language's allocation/GC strategy (Rust
`Box`, C `malloc`/`free`, OCaml/Haskell GC, Osprey's runtime).

## Methodology

1. **Build once, time the binary.** `osprey … --compile` emits a persistent
   native executable; we time *that*, never `--run` (which bundles compile+link
   into the measurement). Comparison languages use their standard optimizing
   release flags (below).
2. **Correctness oracle.** Each binary is run once and its output compared to the
   case's `expected.txt`. A mismatch or build failure is reported and excluded
   from timing — we never publish a number for a program that computed the wrong
   thing.
3. **CPU.** `hyperfine -N --warmup 3 --min-runs 10` per case across all available
   languages → statistical mean ± stddev.
4. **Memory.** `/usr/bin/time` peak resident set size, max over a few runs
   (`-l` on macOS, `-v` on Linux).
5. **Missing toolchains are skipped,** not fatal.

Compile commands (source of truth: [`run.sh`](run.sh)):

| Lang | Command |
|------|---------|
| Osprey  | `osprey <f>.osp --compile` (emits LLVM IR, compiled by clang at `-O2`; override with `OSPREY_OPT`) |
| Rust    | `rustc -C opt-level=3 -C overflow-checks=off -o <bin> <f>.rs` |
| C       | `cc -O2 -o <bin> <f>.c` |
| OCaml   | `ocamlopt -O3 -unsafe -o <bin> <f>.ml` |
| Haskell | `ghc -O2 -o <bin> <f>.hs` |

## Reading the numbers fairly

- **Osprey does checked arithmetic on every `+ - * %`** (each returns
  `Result<int, MathError>`, overflow-checked). The others do not by default — we
  even pass `-C overflow-checks=off` to Rust to match its release profile. Part
  of Osprey's CPU gap is the cost of that safety, a real language semantic.
- **Same algorithm everywhere.** Identical *naive* algorithm and parameters in
  every language — no memoization, closed forms, SIMD, or parallelism. We measure
  the language/compiler/runtime, not who is cleverest. Ranges match Osprey's
  half-open `range(a,b)` = `[a,b)` exactly.
- **Osprey loops via `range |> fold`,** not deep linear recursion, because Osprey
  has no tail-call optimization yet (a 1e6-deep recursion overflows the stack).
  The work is identical; only the iteration mechanism differs.
- **OCaml is built without flambda** (stock `ocamlopt`), so its numbers are
  conservative versus an flambda build.
- **Single machine, wall clock.** Treat ratios as indicative; re-run locally.

## Findings

On the author's machine (Apple Silicon, macOS), geometric mean across the 18
benchmarks (`results/results.md` has the live numbers and a ranked
**tuning-priorities** table):

```
CPU:    Osprey ≈ 1.0× Rust, 1.1× C, 0.7× OCaml, 0.7× Haskell (geomean).
        At parity with the fastest systems languages; faster than OCaml/Haskell.
Memory: ≈ C on 17 of 18 cases (~1.4 MB). One outlier: binarytrees.
```

**Osprey is the fastest of all five languages outright on several cases**
(varies run-to-run, since Osprey/Rust/C are now within measurement noise of each
other) — e.g. `digitsum`, `pascal`, `tak`, `powmod`, `mutual`, `josephus`,
`primes`, `gcdsum`. On the rest it ties C/Rust.

**How it got here.** The first version of this suite measured Osprey 12–89×
slower and using up to 2244× more memory. The cause was **not** the language: the
compiler was handing its LLVM IR to clang with **no optimization flag (`-O0`)**,
so every per-operation `Result` block stayed a live `malloc`. Compiling the IR at
`-O2` lets LLVM prove those allocations non-escaping and delete them outright
(heap → registers): `fib(35)` went from **0.52 s / 1.37 GB to 0.01 s / 1.4 MB**.
See [plan 0010](../docs/plans/0010-cross-language-benchmark-suite.md).

**The one remaining gap: `binarytrees`** (~900 MB, ~24× the fastest). Its tree
nodes genuinely *escape* — they are built, held, then checksummed — so the
optimizer cannot statically free them, and Osprey's default allocator
(`compiler/runtime/memory_runtime.c`, a `malloc` passthrough) does not reclaim
during a run. Allocation now funnels through one swappable `@osp_alloc` hook
([MEM-BACKENDS](../docs/specs/0018-MemoryManagement.md)), so a reclaiming backend
(ARC / arena / tracing GC) can be linked in to close this last gap without
touching the language.

## Not yet benchmarked (and why)

Blocked on language features Osprey doesn't expose today (left out, not faked):

| Benchmark | Blocked on |
|-----------|-----------|
| mandelbrot, n-body, spectral-norm | no `sqrt`/trig stdlib, no `int`↔`float` conversion, float-formatting differs across languages (no exact integer oracle) |
| n-queens, quicksort, mergesort, fannkuch | no list literals / cons / list pattern-matching / mutable arrays |
| sieve of Eratosthenes, matrix-multiply, n-sieve | no mutable arrays |
| pidigits | no arbitrary-precision integers (i64 only) |

`collatz`, `digitsum`, and `isqrt` were unblocked by adding the `intDiv` builtin
(truncating, divide-by-zero-checked integer division — the `/` operator stays
float-only per spec; see [BUILTIN-INTDIV](../docs/specs/0012-Built-InFunctions.md)).
As Osprey grows a math stdlib, numeric conversions, and mutable arrays, the
remaining classics above become expressible and should be added here.

## Adding a benchmark

1. `mkdir benchmarks/cases/<name>/`.
2. Write `<name>.osp`, verify: `osprey <name>.osp --run`.
3. Add `expected.txt` (exact stdout) and `bench.json` (`{"name","description"}`).
4. Add `<name>.rs`, `<name>.c`, `<name>.ml`, `<name>.hs` — identical algorithm
   and parameters, each printing only the integer result.
5. `BENCH_FILTER=<name> make bench` and confirm every language reports `ok`.
