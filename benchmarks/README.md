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

> **Heads-up on RAM.** Osprey currently allocates heap memory per operation and
> does not reclaim it during a run (see *Findings*), so the heavier cases peak at
> **2–3 GB** each. Run on a machine with ≥8 GB free, or use `BENCH_FILTER`.

Results are written to `benchmarks/results/` (gitignored):

| file | contents |
|------|----------|
| `results.md`   | human-readable CPU + memory tables and Osprey-relative slowdown |
| `results.json` | the same data, structured, for tracking over time |
| `hf/*.json`    | raw [hyperfine](https://github.com/sharkdp/hyperfine) exports per case |

## The benchmarks (15)

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
| Osprey  | `osprey <f>.osp --compile` |
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

## Findings (illustrative)

On the author's machine (Apple Silicon, macOS), geometric mean across the suite:

```
CPU:    Osprey ~10–15× slower than Rust / C / OCaml / Haskell — a stable gap,
        no single benchmark blows up disproportionately.
Memory: Osprey uses ~100–1000× more peak RSS than the others.
```

The memory result is the headline. Peak RSS scales with the **number of
operations**, not live data: `print("${1+1}")` ≈ 1.4 MB, `fib(25)` ≈ 13 MB,
`fib(35)` (≈120× more calls) ≈ 1.4 GB. That is the signature of the runtime
heap-allocating per checked operation / call and not freeing during a run (no GC,
no scope-based free) — the single highest-leverage thing to fix, now tracked
directly by this suite.

## Not yet benchmarked (and why)

Blocked on language features Osprey doesn't expose today (left out, not faked):

| Benchmark | Blocked on |
|-----------|-----------|
| mandelbrot, n-body, spectral-norm | no `sqrt`/trig stdlib, no `int`↔`float` conversion, float-formatting differs across languages (no exact integer oracle) |
| n-queens, quicksort, mergesort, fannkuch | no list literals / cons / list pattern-matching / mutable arrays |
| sieve of Eratosthenes, matrix-multiply, n-sieve | no mutable arrays |
| collatz, digit-sum, integer-sqrt | no integer division or bit-shift operator (only float `/`) |
| pidigits | no arbitrary-precision integers (i64 only) |

As Osprey grows a math stdlib, numeric conversions, integer division, and
arrays, these become expressible and should be added here.

## Adding a benchmark

1. `mkdir benchmarks/cases/<name>/`.
2. Write `<name>.osp`, verify: `osprey <name>.osp --run`.
3. Add `expected.txt` (exact stdout) and `bench.json` (`{"name","description"}`).
4. Add `<name>.rs`, `<name>.c`, `<name>.ml`, `<name>.hs` — identical algorithm
   and parameters, each printing only the integer result.
5. `BENCH_FILTER=<name> make bench` and confirm every language reports `ok`.
