# Implementation status — 22 October 2026

The application, documents and design assets are consolidated under the standalone `security-app/` root. This is a working read-only investigation build. The full production implementation is **not complete**; research, specification and delivery gates remain the acceptance contract.

| Contract | Implemented and checked | Remaining work |
| --- | --- | --- |
| [SECURITY-COVERAGE] | Typed coverage states, explicit partial scopes, failed collection and stale restored evidence | Continuous collector-health reconciliation, loss accounting and age policy |
| [SECURITY-COLLECTORS-IDENTITY] | Standard app roots, PID/start-time process snapshots without argument collection, NSWorkspace foreground/background observations matched to launch time, launch plist parsing, on-demand codesign validation | Exhaustive inventory, ES events and trusted process identity revalidation for actions |
| [SECURITY-COLLECTORS-PERMISSIONS] | Declared usage strings and validated sandbox entitlements, labeled separately from grants | Actual privacy grant observations; guidance now covers six sensitive categories and opens System Settings |
| [SECURITY-EVIDENCE] | Osprey review rules; source/time/scope; native inspector/filter; private latest report, bounded snapshot history, finding lifecycle/expiry/reopening, explicit export and startup/package difference counts | Independent benign/attack corpus, full audit trail/database migrations and approved privileged actions |
| Scoped files | Immediate hidden metadata; selected regular-file SHA-256, baseline recording and comparison, bounded reads and descriptor race checks | Trusted origin attestation, OS-version-bound baselines and stronger adversarial filesystem identity guarantees |
| System integrity | SIP, authenticated-root and Gatekeeper configuration commands | Verification of protected content and trusted writable-file baselines; these commands alone do not establish system integrity |
| Homebrew | Standard Cellar name/version plus bounded receipt/tap/request/bottle metadata without invoking package code | Authenticated provenance, casks, nonstandard prefixes, services and installed-version advisory matching |
| iOS | Native JSON report import and review, schema/duplicate/size validation; compilation | Physical-device tests, signing/provisioning, report authentication, pairing and permitted protection providers |
| [AGENT-INTEGRATION] | AI off by default; optional native model picker, session JWT, question disclosure and bounded end-user HTTPS client; no local tools | Authenticated deployment tests, verified model catalog, sign-in/token refresh, scoped evidence/tool workflow, remote lifecycle and adversarial integration tests |
| Framework | Shared native design components and presentation helpers | Typed backend-independent Osprey composition/state/effects framework and shared storage-app consumption |

## External prerequisites and release gates

Endpoint Security and Network Extension prevention are not implemented. Their entitlement eligibility, signing identity, installation/approval and recovery must be proved on stock devices before enabling a protection claim. The NIMBLESITE native client is implemented against the inspected local vendor schema but has no verified live authenticated tenant/configuration; no operator credential belongs in the app. These prerequisites do not prevent local read-only collection.

Enterprise enrollment, MDM, tenant isolation, policy signing, advisory distribution, secure updates, notarization, App Store submission, independent security assessment, false-positive/false-negative measurement and the 24-hour lifecycle soak remain open. No business-security endorsement, detection efficacy or production readiness is claimed.

## Verification evidence

Tests cover launch activation, contrast, native panel ownership, malformed launch definitions, command argument handling/timeouts/output limits, process parsing and live PID observation, scoped hidden entries, unsupported collection, valid/unsigned signing evidence, malformed entitlement output, real Osprey batch integration, invalid engine output, local report permissions, stale reloads, corruption, schema/identity rejection and snapshot differences. Osprey rules have benign controls for hidden files/background processes and review cases for signing, sandbox exposure, hidden startup and disabled protections.

Live Mac interaction checks include inventory navigation, filtering, selecting evidence and Calculator signature/sandbox validation. Tests establish these specific behaviors, not anti-spyware effectiveness. Run `make ci` from the application root with `OSPREY_COMPILER` set to reproduce automated checks.

Apple requires signature validation before trusting returned signing information: “you must successfully call” the validity check before extracting trusted metadata. See [SecCodeCopySigningInformation](https://developer.apple.com/documentation/security/seccodecopysigninginformation(_:_:_:)). The adapter runs strict static validation first; failures remain unverified, not evidence of absent sandboxing or malware.

The application-local `make ci` passed on 22 October 2026: 30 Swift XCTest cases, nine Osprey rule cases, and unsigned iOS compilation. The launch, output-size, entitlement-parser and panel-ownership regression checks were observed failing before their fixes. The Mac locked during final interactive testing; the corrected sheet picker, export dialog, final cancellation behavior and physical iOS import are not claimed as interactively verified. Static and automated checks do not substitute for those remaining checks.

## Expanded application verification

The workspace now includes native Investigation, Review, History, File integrity and Optional assistant tabs. Network and Advisories expose explicit collection/import actions in Investigation. CISA catalog records remain research material with unknown installed applicability. The same optional assistant is compiled into the iOS companion; no account is required for local report review.

Additional tests cover suppression expiry at reconciliation and reload, evidence changes, unknown-signing resolution prevention, recurrence, partial-scan retention, history quotas/permissions, receipt provenance, exact SHA-256 output, symbolic-link/directory rejection, permission guidance, socket parsing and a real loopback listener, malformed advisory counts, stale-cache preservation, AI consent, end-user-only request headers, contract validation, token header injection and rejected agent tools. Osprey tests now check explanatory uncertainty and the MITRE source reference as well as classifications.

New request transport was not exercised against a billable authenticated tenant. The native desktop remained locked during this pass; no new interactive UI acceptance is claimed. The application-local Deslop CLI scanned zero supported files (Swift/Osprey unsupported), so its empty report is not evidence of no duplication; shared persistence/HTTP helpers and source were reviewed manually.

The final packaging regression reproduced an invalid resource seal on the assembled `.app`. Build tooling now replaces executable files atomically and ad-hoc signs the helper and assembled bundle after writing resources/Info.plist. Dedicated tests verify the seal and launch the packaged executable, separately from the Swift build product. Ad-hoc signing supports local development only; Developer ID/notarization remains a release gate.

Final automated result for this pass: **52 Swift XCTest cases, 11 Osprey cases, packaged-bundle signature/launch checks, and unsigned iOS compilation passed** through application-local `make ci`. Offscreen native rendering also succeeded and was inspected with real saved inventory (marked stale); it does not replace unlocked interaction testing. Verification logs and rendered views are under `.build/verification/`.
