# Osprey cross-language benchmarks
Languages compared: Osprey, Rust, C, OCaml, Haskell.
Lower is better. CPU = hyperfine mean ± stddev; memory = peak resident set size.

## CPU time
| Benchmark | Osprey | Rust | C | OCaml | Haskell |
|---|---|---|---|---|---|
| tak | 488.8 ms ± 7.3 | 32.9 ms ± 0.5 | 33.5 ms ± 2.2 | 44.4 ms ± 0.7 | 63.3 ms ± 2.0 |

## Relative speed (× the fastest language per benchmark)
| Benchmark | Osprey | Rust | C | OCaml | Haskell |
|---|---|---|---|---|---|
| tak | 14.9× | 1.0× | 1.0× | 1.3× | 1.9× |

## Peak memory
| Benchmark | Osprey | Rust | C | OCaml | Haskell |
|---|---|---|---|---|---|
| tak | 1162.0 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |

## Where Osprey sits
**CPU** (geometric mean across benchmarks):
- 14.9× slower than Rust
- 14.6× slower than C
- 11.0× slower than OCaml
- 7.7× slower than Haskell

**Peak memory** (geometric mean across benchmarks):
- 751× more than Rust
- 836× more than C
- 527× more than OCaml
- 105× more than Haskell
