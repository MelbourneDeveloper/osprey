# Package-system fixtures

Declarative inputs for the suites in `registry/tests/src`. Every file is data; nothing here executes. Each JSON file names the spec section it pins in a `spec` field.

| File | Pins |
| --- | --- |
| `names.json` | `[PACKAGE-NAMES]` accepted keys and each rejection's reason code. |
| `release-ids.json` | `[PACKAGE-IDENTITY]` display forms, including malformed digests, epochs and ordinals. |
| `paths.json` | `[PACKAGE-SOURCE-CANONICAL]` path grammar, case-fold collisions, excluded control paths, forbidden entry kinds, the LFS pointer shape and submission bounds. |
| `source/*.json` | Whole trees with each file's SHA-256, the exact envelope bytes and `SourceDigest`; `minimal-tree.json` also carries a digest-substitution variant. |
| `manifests/ledger/` | The spec 0029 example manifest, its RFC 8785 canonical JSON and `H("osprey.manifest.v1", …)`. |
| `digests.json` | `H(d, x)` over small payloads. |

Goldens were computed once with an independent implementation (Python `hashlib` and `json`, which agrees with RFC 8785 for ASCII keys and integer-only values) and are asserted byte for byte by Osprey code that never saw that implementation. Change a fixture only together with the spec change that justifies it.
