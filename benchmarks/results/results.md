# Osprey cross-language benchmarks
Languages compared: Osprey, Rust, C, OCaml, Haskell.
Lower is better. CPU = hyperfine mean ± stddev; memory = peak resident set size.

## Tuning priorities — where Osprey needs work (worst gap first)
Ranked by Osprey's CPU slowdown against the *fastest other* language on each benchmark. Start at the top. Memory overhead (vs the lightest other language) is shown alongside — in this suite it is the bigger problem and tracks operation count, not data size.

| Rank | Benchmark | Osprey CPU | × slower vs fastest | Osprey peak RSS | × heavier |
|---|---|---|---|---|---|
| 1 | binarytrees | 420.3 ms | 25× | 905.0 MiB | 531× |
| 2 | ackermann | 127.9 ms | 2× | 1.6 MiB | 1× |
| 3 | coins | 71.7 ms | 1× | 1.4 MiB | 1× |
| 4 | collatz | 12.3 ms | 1× | 1.4 MiB | 1× |
| 5 | isqrt | 13.4 ms | 1× | 1.4 MiB | 1× |
| 6 | fib | 21.1 ms | 1× | 1.4 MiB | 1× |
| 7 | coprime | 62.4 ms | 1× | 1.4 MiB | 1× |
| 8 | powmod | 23.4 ms | 1× | 1.4 MiB | 1× |
| 9 | mutual | 13.4 ms | 1× | 1.4 MiB | 1× |
| 10 | gcdsum | 81.1 ms | 1× | 1.4 MiB | 1× |
| 11 | nestedloop | 44.8 ms | 1× | 1.4 MiB | 1× |
| 12 | pascal | 27.7 ms | 1× | 1.4 MiB | 1× |
| 13 | tak | 32.7 ms | 1× | 1.4 MiB | 1× |
| 14 | hanoi | 38.4 ms | 1× | 1.4 MiB | 1× |
| 15 | josephus | 32.8 ms | 1× | 1.4 MiB | 1× |
| 16 | factorial | 33.5 ms | 1× | 1.4 MiB | 1× |
| 17 | primes | 6.3 ms | 1× | 1.4 MiB | 1× |
| 18 | digitsum | 4.9 ms | 1× | 1.4 MiB | 1× |

## CPU time
| Benchmark | Osprey | Rust | C | OCaml | Haskell |
|---|---|---|---|---|---|
| ackermann | 127.9 ms ± 1.3 | 132.5 ms ± 2.8 | 128.0 ms ± 2.0 | 113.6 ms ± 2.7 | 65.6 ms ± 3.7 |
| binarytrees | 420.3 ms ± 6.8 | 713.2 ms ± 5.7 | 348.9 ms ± 3.1 | 50.7 ms ± 0.7 | 16.7 ms ± 1.0 |
| coins | 71.7 ms ± 1.6 | 76.2 ms ± 2.5 | 71.1 ms ± 2.7 | 93.6 ms ± 2.5 | 52.3 ms ± 1.6 |
| collatz | 12.3 ms ± 0.4 | 11.4 ms ± 0.6 | 9.4 ms ± 0.4 | 54.2 ms ± 1.0 | 39.3 ms ± 4.3 |
| coprime | 62.4 ms ± 1.0 | 60.9 ms ± 2.0 | 58.5 ms ± 0.6 | 88.1 ms ± 1.0 | 100.0 ms ± 2.6 |
| digitsum | 4.9 ms ± 0.2 | 5.2 ms ± 0.2 | 5.3 ms ± 0.2 | 19.0 ms ± 2.6 | 29.1 ms ± 2.6 |
| factorial | 33.5 ms ± 1.2 | 34.8 ms ± 0.8 | 34.6 ms ± 0.9 | 50.2 ms ± 6.0 | 53.9 ms ± 3.6 |
| fib | 21.1 ms ± 0.5 | 17.8 ms ± 2.0 | 18.6 ms ± 1.7 | 24.3 ms ± 1.4 | 49.8 ms ± 4.4 |
| gcdsum | 81.1 ms ± 5.9 | 79.8 ms ± 2.0 | 79.3 ms ± 1.4 | 101.6 ms ± 1.4 | 103.3 ms ± 3.9 |
| hanoi | 38.4 ms ± 1.5 | 38.8 ms ± 1.6 | 39.3 ms ± 1.7 | 61.6 ms ± 1.0 | 55.8 ms ± 5.1 |
| isqrt | 13.4 ms ± 0.5 | 11.2 ms ± 0.9 | 10.6 ms ± 0.4 | 20.9 ms ± 1.3 | 40.5 ms ± 1.4 |
| josephus | 32.8 ms ± 0.6 | 33.5 ms ± 0.9 | 33.4 ms ± 1.6 | 41.2 ms ± 0.5 | 44.5 ms ± 4.8 |
| mutual | 13.4 ms ± 1.1 | 13.0 ms ± 0.9 | 12.8 ms ± 0.5 | 28.9 ms ± 1.7 | 40.5 ms ± 2.1 |
| nestedloop | 44.8 ms ± 0.8 | 46.4 ms ± 3.2 | 44.5 ms ± 0.9 | 57.1 ms ± 0.4 | 63.9 ms ± 1.7 |
| pascal | 27.7 ms ± 1.0 | 27.7 ms ± 1.0 | 27.6 ms ± 1.2 | 44.6 ms ± 0.8 | 62.3 ms ± 4.3 |
| powmod | 23.4 ms ± 0.9 | 22.7 ms ± 0.9 | 22.4 ms ± 1.0 | 59.6 ms ± 0.9 | 57.4 ms ± 4.6 |
| primes | 6.3 ms ± 0.3 | 6.7 ms ± 0.9 | 6.6 ms ± 0.6 | 8.8 ms ± 1.0 | 15.8 ms ± 1.3 |
| tak | 32.7 ms ± 0.7 | 32.8 ms ± 0.5 | 32.9 ms ± 0.5 | 45.1 ms ± 2.2 | 64.4 ms ± 1.7 |

## Relative speed (× the fastest language per benchmark)
| Benchmark | Osprey | Rust | C | OCaml | Haskell |
|---|---|---|---|---|---|
| ackermann | 1.9× | 2.0× | 2.0× | 1.7× | 1.0× |
| binarytrees | 25.2× | 42.8× | 21.0× | 3.0× | 1.0× |
| coins | 1.4× | 1.5× | 1.4× | 1.8× | 1.0× |
| collatz | 1.3× | 1.2× | 1.0× | 5.8× | 4.2× |
| coprime | 1.1× | 1.0× | 1.0× | 1.5× | 1.7× |
| digitsum | 1.0× | 1.1× | 1.1× | 3.9× | 5.9× |
| factorial | 1.0× | 1.0× | 1.0× | 1.5× | 1.6× |
| fib | 1.2× | 1.0× | 1.0× | 1.4× | 2.8× |
| gcdsum | 1.0× | 1.0× | 1.0× | 1.3× | 1.3× |
| hanoi | 1.0× | 1.0× | 1.0× | 1.6× | 1.5× |
| isqrt | 1.3× | 1.1× | 1.0× | 2.0× | 3.8× |
| josephus | 1.0× | 1.0× | 1.0× | 1.3× | 1.4× |
| mutual | 1.0× | 1.0× | 1.0× | 2.2× | 3.1× |
| nestedloop | 1.0× | 1.0× | 1.0× | 1.3× | 1.4× |
| pascal | 1.0× | 1.0× | 1.0× | 1.6× | 2.3× |
| powmod | 1.0× | 1.0× | 1.0× | 2.7× | 2.6× |
| primes | 1.0× | 1.1× | 1.0× | 1.4× | 2.5× |
| tak | 1.0× | 1.0× | 1.0× | 1.4× | 2.0× |

## Peak memory
| Benchmark | Osprey | Rust | C | OCaml | Haskell |
|---|---|---|---|---|---|
| ackermann | 1.6 MiB | 1.7 MiB | 1.6 MiB | 2.6 MiB | 15.1 MiB |
| binarytrees | 905.0 MiB | 2.2 MiB | 1.7 MiB | 5.1 MiB | 11.0 MiB |
| coins | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| collatz | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| coprime | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| digitsum | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| factorial | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| fib | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| gcdsum | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| hanoi | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| isqrt | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| josephus | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| mutual | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| nestedloop | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.0 MiB |
| pascal | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| powmod | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| primes | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |
| tak | 1.4 MiB | 1.5 MiB | 1.4 MiB | 2.2 MiB | 11.1 MiB |

## Where Osprey sits
**CPU** (geometric mean across benchmarks):
- 1.0× slower than Rust
- 1.0× slower than C
- 0.7× slower than OCaml
- 0.7× slower than Haskell

**Peak memory** (geometric mean across benchmarks):
- 1× more than Rust
- 1× more than C
- 1× more than OCaml
- 0× more than Haskell
