# Osprey package system

The package manager and registry from [plan 0020](../docs/plans/0020-package-manager.md) and specs [0029](../docs/specs/ecosystem/0029-PackageManagement.md)–[0032](../docs/specs/ecosystem/0032-PackageEvidenceAndDiscovery.md), written entirely in Osprey and built by the existing compiler (`[PACKAGE-IMPLEMENTATION-LANGUAGE]`). **All Osprey, and a working API that serves working packages comes first.** The active work is plan 0020 Milestone 0: the database-free folder registry of spec 0029 `[PACKAGE-FOLDER-REGISTRY]`, served by `api/` and fetched by `cli/`. The sections below say exactly what runs today.

## Layout

| Directory | Namespace | Contents |
| --- | --- | --- |
| `core/` | `registry` | Pure rules shared by every frontend: package names, release identities, digest envelopes and `H(d, x)`, the canonical source-tree grammar and `SourceDigest`, a pure-Osprey SHA-256, byte encodings, canonical ordering, the arithmetic policy, the language-conformance gate and the blocked-capability inventory. No I/O. |
| `host/` | `registry_host` | The `Host` effects (`Files`, `Tree`, `Web`, `Process`, `Console`) and their native handlers built on Osprey builtins. Handlers carry no policy. |
| `api/` | `registry_api` | The folder-registry HTTP server: configuration, catalog loading, the pure `Router`, and the listener. |
| `cli/` | `registry_cli` | Command routing and the commands themselves (`help`, `doctor`, `fetch`). `src/main.osp` is the executable entry; `lib/` holds everything tests assemble. |
| `sample-registry/` | — | A registry root with four real releases of three packages; `@demo/report` uses the other two. |
| `manifest/` | `registry` | The TOML parser and manifest validators, kept out of `core/` until the type checker handles the parser in seconds. |
| `resolver/` | `registry` | The SAT/CDCL resolver core and its own test project. |
| `tests/` | `registry_tests` | Every suite, assembled with core, host and `cli/lib` into one program. `fixtures/` holds declarative JSON fixtures whose digest goldens were computed with an independent implementation (Python `hashlib`/`json`) and are asserted byte for byte. |
| `blockers/` | `registry_blockers` | Probes for capabilities Osprey lacks; red by design until each lands. |
| `../schemas/package-registry.td` | — | The canonical typeDiagram model. Spec 0031 renders the same bytes and a test keeps them identical. |

Projects share source through `source_roots` (for example `cli` lists `../core/src`), because Osprey has no cross-project dependency model yet. `osprey test` compiles each `*.test.osp` in isolation and cannot import project modules, so suites are ordinary modules with an `export fn run()` called from `tests/src/main.osp`, run with `osprey registry/tests --run`.

## Run a registry

A registry is a folder. Each release is a directory `@scope/name/<epoch>.<ordinal>/` holding the package's source tree, as in `sample-registry/`. Nothing else is needed: no database, no account, no hosted service.

```sh
make registry-serve          # serves registry/sample-registry on http://127.0.0.1:8787
curl http://127.0.0.1:8787/v1/packages
curl http://127.0.0.1:8787/v1/packages/demo/maths
```

To serve your own folder on your network, edit `api/registry-server.json`:

```json
{ "root": "/srv/osprey-packages", "address": "0.0.0.0", "port": 8787 }
```

Then run `osprey . --run` in `registry/api`. At startup the server walks the root, computes each release's `SourceDigest`, prints every release it serves and every directory it refuses with the reason, and answers the endpoints of spec 0029 `[PACKAGE-FOLDER-REGISTRY]`.

Fetch from it with the package frontend, which verifies the digest before writing anything:

```sh
echo "fetch @demo/report --registry http://127.0.0.1:8787" | ./target/osprey-pkg
```

The release lands in `.osprey/packages/@demo/report/<digest>/`.

## Running

```sh
make registry-test                                   # all suites, TAP output
make registry-blockers                               # capability probes, expected red
./target/release/osprey build registry/cli -o target/osprey-pkg
echo "doctor" | ./target/osprey-pkg                  # one command line on stdin
```

The frontend reads its command line from standard input and prints `status: N` as its last line because Osprey has no argv or exit-status builtin yet. Both are blocked capabilities below, not design.

## Host capability inventory (plan 0020, Phase 0)

Measured against the compiler at `target/release/osprey` on macOS. `registry/core/src/capabilities.osp` is the single source for the blocked list; `osprey-pkg doctor` prints it.

| Capability | State | Evidence |
| --- | --- | --- |
| Read and write whole text files with OS failure reasons | available | `readFile`/`writeFile`; `tests/src/dispatch_tests.osp` through `Host::Files` |
| Run a shell command, capture stdout/stderr and exit code | available | `spawnProcess` + redirection in `host/src/native.osp` |
| Standard input, one line at a time | available | `input()` in `cli/src/main.osp` |
| HTTP/WebSocket client and server, JSON documents | available | runtime builtins; not yet used here |
| Byte-level string access, 64-bit integers, lists, maps | available | `core/src/bytes.osp`, `core/src/order.osp` |
| Command-line arguments | blocked `process.argv` | no builtin; not probeable from inside a program |
| Process exit status | blocked `process.exit-status` | `main`'s value is ignored |
| Environment variables | blocked `process.environment` | no builtin |
| Binary-safe file reads | blocked `files.binary-read` | `blockers/src/main.osp`: a 5-byte file with a NUL reads as 1 byte |
| Directory listing, stat, mkdir, rename, unlink, fsync | blocked `files.directory` | no builtins; the host shells out, which Windows and atomicity rules will not accept |
| UTC clock | blocked `time.clock` | no builtin |
| Bitwise operators | blocked `int.bitwise` | SHA-256 emulates them with arithmetic and lookup tables |

A blocked capability is resolved only by Osprey gaining it; no component is written in another language meanwhile (plan 0020, fixed risk controls).

## Measured

- `SourceDigest` of the fixture trees equals the independently computed goldens, including the digest-substitution case.
- SHA-256 in Osprey hashes 64 KiB in about 1.6 s on an Apple Silicon laptop: bitwise operations are nibble-table lookups because the language has none. Correctness is complete; throughput is the `int.bitwise` blocker's cost.

## Defects found while building (reported, not worked around in the compiler)

| Probe in `registry/blockers` | Behavior today | Needed |
| --- | --- | --- |
| `files.binary-read` | `readFile` returns the bytes before the first NUL only. | Whole-file reads so declared PNG/JSON assets can be hashed. |
| `runtime.default-operator-pipeline` | `byteAt(s, i) ?: 0` inside a `map` or `fold` over a million elements crashes (signal 11); the `match` form and a pipeline without `?:` both succeed. | `?:` must behave like the `match` it abbreviates at any iteration count. |

## Known compiler limits that shape the code

- A recursive function with generic parameter types is instantiated once per program (issue #277). `Order::sortBy` therefore recurses over integer index lists only, and a second use of it at a different element type inside the `registry` namespace still fails to type-check (issue #278), so `Arithmetic::percentile` uses the builtin `listSort`.
- Two records sharing a field name make an untracked receiver a codegen error (issue #273); functions that read such fields spell out the record type and say so in a comment.
- `==` on lists compares identity, and lists are not assertion operands; tests compare rendered strings.
- Imports live at file level, never inside a module, and constructors from another module are qualified (`Names::Key { ... }`).
