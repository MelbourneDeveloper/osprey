# Implementation, extraction and release plan

**Status: research/design phase.** This plan builds a professional native anti-spyware product in a separate repository. The application must not ship as an Osprey demonstration. Contracts are in the [product specification](specification.md), [native framework](ui-framework.md) and [optional agent](ai-agent.md); [research](research.md) owns primary-source evidence. No implementation milestone below is complete merely because its document exists.

## Ownership and repository extraction [DELIVERY-EXTRACTION]

Keep this dossier together until the product repository is provisioned, then move it under that repository's `docs/` and update inbound links. The new repo owns security domain/features, native app hosts, extension targets, framework packages, deterministic fixtures, integration tests, release tooling and deployment documentation. Compiler/runtime defects stay in Osprey with regression tests; consume pinned toolchain artifacts rather than copying compiler sources or depending on a developer's checkout.

Extract generic native-app infrastructure once into a versioned library boundary. Both the storage feature and security app consume it. Until an Osprey package manager is available, use explicit local source roots or a pinned verified source artifact; Swift components use a Swift package. Avoid parallel framework implementations in the two repos. Preserve upstream MIT notices and inventory third-party licenses. The new product's license, name and ownership metadata are release decisions; do not silently assign the entire product Osprey's license.

| Existing infrastructure | Reuse decision and required work |
| --- | --- |
| `examples/macos/storage/src/{main,browse,view}.ospml` | Preserve pure update, effect handlers, UI ownership and stale-scan handling; extract general lifecycle/composition rather than copy storage domain logic. |
| `examples/mobile/ios/IssueInbox/NativeRenderer.swift` | Evolve in place into the shared native backend; split by responsibility, introduce typed protocol validation and move generic code once. Both apps use it. |
| `StorageStore.swift` | Extract serialized dispatch, command correlation, cancellation and scope management; eliminate shared mutable globals where an app/session store is required. |
| `StorageServices.swift` | Reuse reviewed enumeration and accounting techniques. Harden symlink/identity races, denied-access reporting and cloud/mount boundaries. Do not reuse its TCC-file-open heuristic as an authoritative permission audit. |
| Apple build scripts, C ABI and xcconfig | Reuse target selection and generated headers; replace sample bundle IDs, permission strings and signing with production-specific settings. Current examples are not store-ready artifacts. |
| Storage and mobile domain/smoke suites | Preserve regression assertions and adapt reusable framework scenarios; add threat fixtures and real-device permission tests. |

## Milestone order and acceptance [DELIVERY-MILESTONES]

### Gate A: platform, distribution and runtime feasibility

- [ ] Record supported OS/device baselines; start Mac Apple silicon, add other architectures only with verified toolchain support. Match iOS minimum versions to the tested native framework APIs.
- [ ] Request/prove ES entitlement eligibility and a signed install on a stock Mac; document user approval, privacy access and supported event coverage without disabling SIP.
- [ ] Prototype the chosen Network Extension provider, prove process attribution and install/recovery; independently assess App Store eligibility and iOS URL-filter/tunnel constraints against current Apple tables.
- [ ] Produce a store-ready capability profile and a direct-distribution profile. Do not defer this decision until marketing or silently remove anti-spyware protections to pass review.
- [ ] Pin failing tests for runtime ownership/lifetime limitations, then implement a supported bounded-memory embedding contract. Verify Unicode, scalar ABI, copied inputs/outputs, disposal and repeated dispatch through Mac/iOS libraries.
- [ ] Deliver reproducible evidence: signed install/uninstall logs, permission-denied results, ABI suite and an allocation soak report. Unavailable ES means a clearly labeled read-only beta, not a continuous-protection launch.

This gate is first because a native framework cannot make unavailable privileges or unsafe long-running memory behavior disappear.

### Gate B: framework and ergonomic authoring

- [ ] Define typed Osprey actions, bindings, routes, components, command/effect requests and a bounded versioned host protocol. Prove the proposed generic APIs compile using current language features; add negative type-checking fixtures before implementation.
- [ ] Extract the storage-app pattern into reusable core/UI/test libraries and a SwiftUI native backend. Implement headless interaction semantics for backend-independence tests.
- [ ] Deliver navigation, collections, forms, command menus, scoped state/subscriptions, themes/localization, accessibility and lifecycle as [framework contracts](ui-framework.md), not per-screen host branches.
- [ ] Run the storage and security consumers through the same packages. Prove that adding a feature does not require editing renderer internals.
- [ ] Pass keyboard/VoiceOver, IME, live-update edit preservation, multiwindow, iPhone/iPad adaptation and large-list tests. Measure authoring boilerplate and revise the API until shared mechanisms replace it.

### Gate C: local read-only Mac investigation

- [ ] Build inventory/signing, running app/process, launch persistence and scoped metadata collectors with explicit coverage states.
- [ ] Implement permission distinctions, sandbox evidence, hidden-file context and trusted writable-configuration baselines.
- [ ] Implement local evidence storage, retention/migrations, deterministic findings, counterevidence, suppression and export.
- [ ] Deliver synthetic spyware/persistence scenarios and a substantial benign developer/business-workstation corpus. Each detection rule must have an independently asserted positive and negative case.
- [ ] Demonstrate every initial workflow without AI, cloud authentication or network availability; cached feed age remains visible.

### Gate D: production anti-spyware monitoring and intelligence

- [ ] Add approved ES collection with loss accounting, resume/restart reconciliation, bounded queues and process identity validation.
- [ ] Add authenticated, versioned advisory/rule snapshots and verified package/version matching, including Homebrew provenance and services.
- [ ] Build process/persistence/signing/file/network correlation for the threat scenarios in [the specification](specification.md). Distinguish observed evidence from inferred intent.
- [ ] Tune on labeled benign/attack fixtures, record false positives and missed scenarios, and require independent technical review of detection logic and coverage claims.

### Gate E: network protection and response

- [ ] Ship observation mode first; validate destinations and attribution under VPNs, QUIC, encrypted DNS, proxies, shared hosting, sleep and fast process churn.
- [ ] Implement explicit scoped rules, precedence, expiry, local approval, identity revalidation, rollback and collector-health UI.
- [ ] Test failure policy, captive portals, business-critical app exclusions, uninstall/recovery and restart behavior. Changes must leave an auditable outcome and verified postcondition.
- [ ] Independently assess any additional containment/quarantine functionality before adding it; no blanket process-killing or hidden-file cleanup feature qualifies as spyware remediation.

### Gate F: optional NIMBLESITE agent

- [ ] Retrieve/pin the deployed API schema; verify native HTTP auth, end-user isolation, contract hashes, errors and cancellation with synthetic evidence.
- [ ] Provision approved model configurations and a user-selectable model catalog without client operator keys or shared mutable model settings.
- [ ] Implement disclosure previews/redaction, scoped tools, citations, cost controls and reviewable action proposals. Resolve provider terms, retention and deletion/export requirements for business customers.
- [ ] Pass AI-off parity, prompt-injection, denied disclosure, wrong-evidence, replay, model-switch and provider-outage tests. No privileged operation depends on AI availability.

### Gate G: iOS and business deployment

- [ ] Deliver native iOS investigation, report import where verified, local explanations, supported prevention capabilities and explicit OS coverage limits. Test on physical supported devices, not only simulators.
- [ ] Add opt-in paired Mac findings with authenticated pairing, device approval, replay protection, minimum synchronized fields and revocation. The iPhone does not acquire the Mac's permissions.
- [ ] Specify and verify MDM-assisted deployment/approval where Apple supports it, policy distribution/signature checks, organization roles, tenant/device separation, update rings, reporting and recovery. Personal use remains independent of organizational enrollment.
- [ ] Pilot with consenting business administrators and independent security practitioners. Publish what the test cohort establishes; do not claim that experts recommend the product before they have assessed it.

### Gate H: production release

- [ ] Pass threat-model review, independent penetration test, dependency/supply-chain review and findings remediation; retain a documented residual-risk register.
- [ ] Complete accessibility/usability, bounded-memory/performance/battery, service restart, migration, recovery, support and deployment tests against release artifacts.
- [ ] Generate software bills of materials, signed artifacts/provenance, vulnerability disclosure policy, incident-response runbook and signed update/rollback safeguards.
- [ ] Complete store metadata, privacy labels/manifests, entitlement justification, reviewer instructions, usable offline mode and evidence-backed protection claims for every store edition. Obtain actual approval before saying approved.
- [ ] Complete Developer ID signing, notarization/stapling, extension lifecycle and managed installation/uninstall for the full Mac edition. Verify from a clean machine and preserve protections/availability during updates.
- [ ] Establish a published support/patch policy, monitored feed/update health, release notes and correction process for erroneous reputational findings.

## Threat and quality evidence [DELIVERY-TESTS]

Use test-first domain/native work. Required scenario classes include legitimate background agents, disguised persistence referencing a writable hidden executable, invalid signatures, changed configuration with a known baseline, stale/withdrawn advisories, denied permission collection, attributed suspicious destinations, process identity swaps, malicious input aimed at the optional agent, and interrupted remediation with rollback. Keep benign controls and errors in the same suites; do not shrink scope to obtain green checks.

Measure UI latency, scan duration, CPU, memory, energy, disk writes and event loss against fixed small/large fixtures. Set documented release budgets from the first prototype measurements, then enforce them against release builds. A bounded-memory soak includes at least 24 hours of synthetic activity, repeated scans/navigation/edits, cancellation, sleep/resume and service restarts; extend duration in business pilots. This is a proposed minimum validation exercise, not a claim of safety after 24 hours.

No compiler-wide benchmark run or existing test execution is needed for this documentation phase. Implementation will run the relevant Osprey domain/corpus, native framework and platform suites; product CI gets its own explicit gates. Never bypass existing compiler CI or treat signing/build success as application acceptance.

## Future Android adapter [DELIVERY-ANDROID]

- [ ] Reuse the same Osprey feature/core/UI modules and semantic protocol through the existing Android C ABI/JNI pattern.
- [ ] Implement native Kotlin presentation for the accepted component contracts, preserving interaction/accessibility behavior.
- [ ] Research Android permissions, package visibility, VPN exclusivity and Play distribution separately; test actual device capabilities and label differences.

Android is a portability proof and later product target. It does not introduce a nonnative Apple UI layer or delay Mac/iOS release.
