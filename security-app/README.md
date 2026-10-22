# Vigil

A standalone native security application under development. All application code, Osprey rules, design assets, tests, research, specifications and plans live in this directory. It can be moved into its own repository without importing source from the enclosing Osprey checkout. The Osprey compiler is an external build dependency.

## Build and run

From this directory, with Xcode installed:

```sh
export OSPREY_COMPILER=/absolute/path/to/osprey
./build.sh
open -n .build/Vigil.app
```

The app opens the real local investigation workspace. `./build.sh` builds and bundles the Osprey engine; running Swift alone after a clean checkout does not build that engine. For this checkout the existing compiler is `/Users/christianfindlay/Documents/Code/osprey/target/release/osprey`.

```sh
make ci
```

This builds the app, runs Osprey rule tests and Swift tests, and compiles the iOS target without signing. macOS 14 and iOS 17 are the deployment baselines. The bundled Mac app is ad-hoc signed for local development, not a Developer ID/notarized release.

## Implemented

- **Investigation:** installed apps, processes and background-app observations, launch definitions, Homebrew receipt provenance, hidden entries and system-protection settings. App signature/sandbox inspection runs on demand.
- **Review:** Osprey assessment, explanations and primary-source references; persisted findings with acknowledge/investigate/reopen and reasoned, expiring suppression. Relevant follow-up evidence resolves a finding; unknown or missing evidence cannot.
- **History:** private latest report plus up to 30 snapshots retained for 30 days; finding records for 90 days, startup/package comparisons and JSON export. Exports contain local identifiers and paths.
- **File integrity:** explicitly selected local regular files, bounded SHA-256 hashing, race checks and recorded baseline comparison. No cloud download, automatic repair or claim that a first baseline is clean.
- **Network:** explicit open-IP-socket observation without name resolution, payload inspection or blocking. PID attribution is provisional.
- **Permissions:** guidance for reviewing six sensitive permission categories in System Settings. Actual other-app grants remain unknown.
- **Advisories:** explicit download from CISA's official KEV mirror or local JSON import, searchable research and visible provenance. Installed-version applicability is not inferred from names. Cached catalogs survive collection failure as stale evidence.
- **Optional NIMBLESITE:** off by default; native model-configuration picker, session-only end-user JWT, exact question disclosure, cancellation and bounded HTTPS transport. Local tools and automatic device-evidence disclosure are disabled. Deployment setup is in [the agent contract](docs/ai-agent.md#implemented-client-and-deployment-contract-agent-client).
- **iOS:** native evidence-report import and optional assistant. Unsigned compilation is verified; physical-device interaction and distribution are not yet verified.

The app is a read-only investigation build. Active spyware prevention, continuous Endpoint Security monitoring, Network Extension blocking, complete TCC visibility, installed-version advisory matching, enterprise management, signed releases and the full backend-independent Osprey UI framework remain unfinished. [Implementation status](docs/implementation-status.md) tracks the remaining contract.

Osprey owns assessment and explanatory policy in `engine/`; native collection, storage and presentation are Swift. A short-lived engine process bounds its lifetime. This does not solve the separate long-running C ABI ownership and framework extraction requirements.

## Design system

The same executable retains the design studio explicitly:

```sh
.build/Vigil.app/Contents/MacOS/Vigil --studio
.build/Vigil.app/Contents/MacOS/Vigil --render .build/previews
```

The [design contract](docs/design-system.md), [research](docs/research.md), [product specification](docs/specification.md), [framework contract](docs/ui-framework.md), [agent contract](docs/ai-agent.md) and [delivery plan](docs/plan.md) belong to this application. `Sources/VigilDesign/Theme.swift` owns the native design tokens; resources include the logo. The studio uses synthetic fixtures, clearly labeled separately from the default live workspace.

## Asset provenance

`Sources/VigilDesign/Resources/Brand.xcassets/VigilIcon.imageset/VigilIcon.png` was created with the built-in Image Generation tool, not fetched from an icon vendor. The exact prompt is in [logo-prompt.txt](logo-prompt.txt). Direction: “A professional anti-spyware macOS app icon: full-bleed deep graphite #172521, a large ivory aperture made from two offset angular brackets around a small muted mint square, precise optical balance, restrained bevel, no text, shield, padlock or eye illustration.” The packaged raster is the source artwork; final Icon Composer/vector, small-size optical variants and store export production remain release work. SF Symbols supply interface icons on Apple platforms and must not be exported as the Android icon library.
