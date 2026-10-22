# Product and security specification [SECURITY-PRODUCT]

**Status: proposed contract, not implemented.** [Research](research.md) owns source evidence and platform constraints. [The framework](ui-framework.md) owns application composition; [the agent contract](ai-agent.md) owns optional AI. Initial implementation targets a directly distributed native Mac app and a constrained native iOS companion.

## Anti-spyware mission and threat model [SECURITY-SPYWARE]

This is a production endpoint-protection application, not a language sample or a disk-cleanup utility. Guard against covert collection, persistence and outbound transfer through evidenced detection, supported preventive controls, investigation and recovery. [CIS Control 10](https://www.cisecurity.org/controls/malware-defenses) calls for controlling malicious installation, spread and execution; [NIST incident-response guidance](https://csrc.nist.gov/pubs/sp/800/61/r3/final) informs the response lifecycle. Mapping requirements to these sources does not establish compliance, certification or endorsement.

| Spyware behavior | Detection/prevention design | Required counterexample and coverage limit |
| --- | --- | --- |
| Credential/input capture ([T1056.001](https://attack.mitre.org/techniques/T1056/001/)) | Correlate available Accessibility/Input Monitoring authorization evidence, executable reputation, persistence and process/network behavior; guide revocation through supported controls. | Legitimate assistive tools and password managers. Having input access is not proof of keylogging; unrestricted live capture visibility is not established. |
| Screen capture ([T1113](https://attack.mitre.org/techniques/T1113/)) | Correlate screen-recording authorization evidence and observable execution/activity with known indicators. | Conferencing/recording tools. Do not install our own screen recorder to claim detection or assume every capture API produces observable process events. |
| Microphone/camera surveillance ([T1123](https://attack.mitre.org/techniques/T1123/)) | Correlate supported permission/device-activity evidence, signer identity and unexpected background/persistence/network context. | Legitimate calls. Per-app live sensor attribution must be proven per OS; unavailable activity remains unknown. |
| Hidden persistent implant ([T1543.001](https://attack.mitre.org/techniques/T1543/001/), [T1564.001](https://attack.mitre.org/techniques/T1564/001/)) | Observe persistence definitions and executable identity changes, suspicious writable/hidden targets and confirmed indicators. | Legitimate launch agents, developer tools and caches. Hidden or persistent alone is insufficient. |
| Command/control and outbound collection ([T1041](https://attack.mitre.org/techniques/T1041/)) | Correlate attributable connections with identity, destination intelligence and independent host evidence; enforce explicit network rules through an approved provider. | Backups, sync and large business transfers. Encryption and incomplete attribution limit payload conclusions. |
| Defense tampering | Verify our component signatures, service health, policy/update integrity, event gaps and allowed protection-state observations; surface degraded protection locally and to authorized administrators. | Legitimate OS/app updates and approved administrator changes. A hostile kernel/root administrator may defeat local observations. |

Before a protective-execution feature ships, prototype supported ES authorization events and a bounded local policy decision path. Use prevalidated rules and identity evidence; never wait for AI, remote reputation or UI interaction in an authorization callback. Define timeout/failure behavior, cache invalidation and operational recovery, then test them. Read-only inventory is a development milestone; the production anti-spyware claim requires measured, tested protective capabilities for the advertised scope.

Assets include business documents, credentials, communications, sensor access, evidence integrity and availability. Adversaries include malicious applications, scripts, compromised packages, deceptive signed software and persistence installed under a user account. Treat hostile evidence, compromised feeds, IPC peers, optional AI output and update infrastructure as attack surfaces. Remote kernel exploits, already-compromised boot chains and targeted mercenary spyware are not covered by an unqualified guarantee. [Apple's threat-notification guidance](https://support.apple.com/en-us/102174) establishes the need for specialist escalation; do not claim that a clean scan clears a targeted iPhone.

The product itself must not become employee surveillance software: do not record keystrokes, messages, screen/audio/video content or browsing histories to implement ordinary protection. Managed collection is transparent and policy-bound, with minimum necessary endpoint evidence. Confirmed/high-risk incidents offer evidence preservation and expert-response guidance; avoid blind cleanup that destroys investigation material.

## Production and business contract [SECURITY-BUSINESS]

Native macOS/iOS app quality, accurate store metadata and privacy disclosures, supportability and signed release artifacts are mandatory. App Store acceptance is evaluated per capability profile; privileges that cannot ship in that channel require the separately specified direct Mac edition. Do not label research documents or unsigned examples “App Store ready.”

The business edition requires authenticated device enrollment, policy version/signature verification, explicit administrative roles, tenant/device isolation, scoped audit export, update rings and documented MDM integration where supported. Enrollment must not grant administrators unrestricted file contents or access to personal AI conversations. Define user/admin conflict resolution, policy expiry/offline behavior, secure de-enrollment and lost-device token revocation. A management outage cannot stop the local detection engine or erase known policy without a recorded reason.

Separate endpoint security telemetry from employee productivity monitoring. Document collection categories, visibility, retention and administrator access; keep customer-controlled evidence export and deletion workflows. Require a customer-facing threat/coverage document, independent security assessment, responsible disclosure process, incident-response runbook and tested support/recovery procedures. Recommendations by cybersecurity experts are a validation objective, not marketing copy invented by the implementation team.

## Capability and coverage contract [SECURITY-COVERAGE]

Every collector reports one of `available`, `partial`, `denied`, `unsupported`, `failed` or `stale`, with observation time, scope and reason. An empty result under denied access never means a clean scan. An overall summary shows checks completed, checks incomplete and the age of each data source; there is no unsupported universal “device is safe” verdict.

| Requested capability | Direct macOS | Consumer iOS | First delivery and limits |
| --- | --- | --- | --- |
| Unnecessary permissions | Signed entitlements; selected privacy observations where supported; guided review | Guidance and user-supplied evidence; no other-app grant enumeration | Distinguish requested capability, granted permission and observed use. Recommend review with a reason; lack of observed use does not prove unnecessary access. |
| Apps with security/privacy history | Versioned installed identity matched to reviewed advisories and incident records | User-selected app/report identity; no assumed complete inventory | Separate current applicable vulnerability, historical incident, disputed claim, remediation and unknown version. |
| Suspicious processes | Process snapshots and identity initially; ES events after entitlement approval | System-wide inspection unsupported | Combine executable identity, ancestry, persistence, permissions and network observations. No single name, CPU spike or unsigned build establishes maliciousness. |
| Background apps | Running apps, foreground state, agents, daemons and configured login items as distinct observations | Other-app background inventory unsupported | A background app is information, not automatically a finding. Configuration does not prove execution. |
| Hidden/problematic files | Scoped directory metadata and selective hashing | Container and explicitly selected files only | Dotfiles/hidden flags are inventory facts. Correlate writable persistence, executable identity and known indicators. Never traverse denied boundaries silently. |
| Modified system files | Signature checks, trusted baselines for selected writable configuration, supported protection-state observations | General system verification unavailable | Separate SSV-protected OS content from writable files. Unknown baseline or permission failure means unknown integrity. |
| Apps without sandboxing | Verify signed entitlement evidence and running context where available | Third-party apps use OS sandbox; no installed-app audit | Show declared sandbox entitlement as such; missing/invalid evidence is unknown. Non-sandboxed software is exposure context, not a malware verdict. |
| Suspicious outbound data / blocking | Network Extension feasibility, metadata attribution and observe mode before filtering | Provider-specific, OS/eligibility-dependent scope only | Metadata supports destination/behavior findings; no payload-theft claims from encrypted traffic. Blocking needs a separately verified provider. |
| Problematic Homebrew packages | Installed receipts, provenance, services and supported advisory matches | Not applicable | Old, disabled, unmaintained, externally tapped and known-vulnerable are distinct. Package identity/version mapping must be evidenced. |

An optional Mac App Store build is a separate capability profile of this contract, not a promise of direct-Mac parity. iOS does not gain Mac privileges by moving outside the App Store. Managed-device capabilities must be labeled and cannot be used to market a consumer feature.

## Domain and finding model [SECURITY-EVIDENCE]

Osprey owns immutable `Asset`, `Observation`, `Coverage`, `Finding`, `Rule`, `Advisory`, `ActionPlan` and `ActionOutcome` concepts. Their required data is:

| Concept | Required information |
| --- | --- |
| Asset | Local ID, asset kind, observed identity and identity confidence. Apps include bundle ID, Team ID, signing identifier, designated requirement/hash where meaningful, version and executable path. Processes include PID plus start time/audit identity; PID alone is insufficient. |
| Observation | ID, collector and version, device-local ID, scan/session, wall and monotonic times, scope, result/error, redacted summary and evidence digest/reference. Keep original and normalized values distinguishable. |
| Finding | Stable rule/asset identity, category, severity, confidence, status, rule version, evidence IDs, counterevidence, source links, first/last seen, applicable version range and suggested action. |
| Advisory | Publisher, authoritative URL, publication/update/retrieval dates, affected identity/version range, applicability, aliases, withdrawal/remediation status, license and freshness. |
| ActionPlan | Local target identity, exact operation, prerequisites, expected change, risk, rollback, expiry and evidence digest. |
| ActionOutcome | Approval identity/time, preflight identity, operation/result, postcondition, rollback state and audit ID. |

Severity expresses potential impact; confidence expresses evidence strength. An observed signature failure can be high-confidence evidence without proving its cause. A historical privacy incident remains historical even if a current version belongs to the same vendor. Missing version mapping cannot become a definite CVE match through name similarity or AI inference.

Finding lifecycle is `new`, `acknowledged`, `investigating`, `resolved` or `suppressed`. Resolution requires a relevant follow-up check, not dismissal. Suppressions bind to a rule and asset identity, record a reason/expiry, and are re-evaluated when signing, version or evidence changes. Preserve contrary evidence and corrections to reputational records.

## Collection contracts [SECURITY-COLLECTORS]

Native services gather facts; Osprey handlers request collection and interpret normalized results. Initial effect responsibilities are inventory, process snapshots, persistence definitions, scoped file metadata, signing, permission observations, advisory snapshots, persistence storage and network events. Asynchronous work completes as correlated events; no suspended dynamic control effect is assumed at the Apple C ABI.

The host must enforce bounded batches, explicit cancellation, scan generations, progress, maximum input size and deadlines. Ignore stale completions after scope changes; report partial data on event loss or failures. Catch-up/rescan after sleep, collector restart and ES/network event loss is required before continuous-monitoring claims. Never block ES response deadlines or network callbacks on AI, a database scan or the UI.

### App, process and persistence inspection [SECURITY-COLLECTORS-IDENTITY]

Inventory only readable app locations and record inaccessible roots. Validate signatures with Apple's Code Signing Services; distinguish on-disk validation from running-code evidence. Enumerate application state through NSWorkspace and supplement it with a verified process API adapter, then ES execution events when authorized. Observe app-list changes rather than relying solely on launch notifications, which omit some background applications. Normalize launch-agent/daemon definitions without executing their contents. Correlate referenced executables, owner, permissions, identity and observed execution. Detect process exit/reuse races and label unavailable parent identity.

### Permission review [SECURITY-COLLECTORS-PERMISSIONS]

Extract entitlements from signing evidence, not untrusted Info.plist assertions. Usage-description strings are requests/explanations, not proof of grants. Supported permission APIs describe our app's access; they are not an all-app permission inventory. A future read-only TCC adapter is separately gated by explicit consent, actual authorization, known schema and OS tests; unknown schemas return unsupported. Never write TCC databases, revoke other apps' grants through private APIs, or claim that Full Disk Access overrides every restriction. Guide people to supported Settings controls for review/revocation.

### File and integrity inspection [SECURITY-COLLECTORS-FILES]

Reuse the storage walk's cancellation, unreadable-item accounting, hidden metadata and hard-link lessons. Constrain access to authorized roots; default to metadata, hash selectively, and do not upload contents. Record symlink identity and resolved target separately; do not follow symlinks outside scope. Detect mount boundaries, APFS identity, cloud placeholders and changes during hashing; never trigger cloud downloads by accident. Use file identity plus a fresh open for sensitive operations, not string-prefix containment alone.

A trusted integrity baseline records origin, OS/build or app version, signature, hash algorithm and creation conditions. A first scan is a baseline observation, not proof that the device began clean. OS updates invalidate incompatible baselines. Do not offer automatic “repair” of sealed OS files or require disabling SIP/SSV.

### Advisory and Homebrew inspection [SECURITY-COLLECTORS-ADVISORIES]

Use reviewed vendor advisories, CISA KEV, suitable OSV records and official regulator findings with explicit provenance. Record affected ranges and fixed versions; preserve withdrawn advisories and feed age. Validate and authenticate distributed rule/feed snapshots, with a pinned trust root, rollback protection and an emergency revocation path. Feed/network failure retains the last valid snapshot with a stale marker. Absence from a feed means no matching record found, not no vulnerability.

Read installed Homebrew receipts/metadata as data first. Never evaluate formula/cask Ruby, run tap hooks or execute an unknown `brew` binary as part of scanning. If a user explicitly permits CLI collection, verify the executable/provenance, constrain environment and arguments, disable automatic updates/analytics for the invocation, bound runtime/output and record that code executed. Preserve formula/cask distinction, tap, installed revision and upstream mapping; do not assume the Homebrew token is an OSV ecosystem identifier. Unsupported mappings remain unknown. No automatic `brew upgrade`, uninstall or service stop.

### Network observation and filtering [SECURITY-COLLECTORS-NETWORK]

Start in observation mode. Collect the minimum available process/audit identity, destination address/port, domain evidence, protocol, byte counts and times, with attribution confidence. DNS observations, SNI and hostname availability are conditional; encrypted DNS, encrypted client hello, QUIC, VPNs and shared IP/CDN hosting limit inference. A foreign address, large transfer or tracker label alone is not confirmed exfiltration. Link findings to current source evidence and observed behavior.

Rules bind to validated app identity and scoped destination criteria, not display name or PID alone. Define precedence as explicit user exception, explicit deny, then configured default allow; an exception must identify the deny it overrides. Expiry, rule conflicts, already-open connections, redirects, child helpers, captive portals, DNS and VPN coexistence require tests. The monitor must say when traffic bypasses the provider. For the first release, provider failure defaults to allowing connectivity with a visible protection gap; explicitly selected stricter behavior is separate future scope. Undo and disabling filtering must remain available without cloud connectivity or AI.

## Response and privilege boundaries [SECURITY-ACTIONS]

The UI, agent transport, collectors and privileged services are separate trust boundaries. Start read-only and request elevated services only for demonstrated operations. Validate IPC callers using signing/audit identity and authorize every operation independently; an unprivileged UI parameter is never authority. Normalize and recheck targets at action time to prevent path swaps, PID reuse and stale findings.

Allowed initial actions are supported Settings navigation, reveal/export, suppression, and approved reversible network-rule changes once filtering is verified. Do not delete hidden files, kill arbitrary processes, alter system protections or silently quarantine business-critical data. Any later destructive/remediation feature needs its own evidence-preservation, authorization and recovery contract. AI can propose a plan but cannot approve one.

## Local data and privacy [SECURITY-PRIVACY]

Scanning runs locally. Store minimum evidence in an app-private versioned database with explicit migrations, restricted file permissions and platform protection; Keychain stores keys/tokens. Define a 30-day default for detailed observations, 90 days for local finding/action history, and 7 days for flow metadata, each configurable and bounded by storage quotas. These are proposed product defaults, not external requirements. User-requested case preservation is explicit and shows its retention. Purging includes caches, backups under our control and exports where applicable; do not promise remote deletion from a local purge.

No telemetry or document uploads by default. Threat-feed retrieval, software updates, optional AI and local diagnostics are separate settings with visible destinations. Usernames, home paths, command arguments, URLs and tokens are sensitive; redact them from logs and exports unless deliberately included. Render evidence as text, constrain links, reject control sequences and never execute scripts from findings or agent replies. Export includes scope, time, rule/feed versions, evidence references and redaction choices. Hashes/audit chains can expose accidental changes but do not make local evidence invulnerable to a compromised administrator.

## Verification contract [SECURITY-VERIFICATION]

Each collector/rule needs positive, benign-negative, denied-access, malformed-input, partial-result and stale-result fixtures before its implementation. Use synthetic local artifacts, not live malware. Cross-check deterministic Osprey findings under the supported backends; verify native ABI and platform permission behavior on current supported devices. Examples that must remain informational include `.git`, a legitimate updater launch agent, a developer's unsigned tool, a normal background sync and a non-sandboxed signed application.

Release requires validated evidence for all nine capability rows, including explicit unsupported states, zero nonconsensual writes, action rollback tests, malformed feed/IPC resilience and the memory/lifecycle gate in [the plan](plan.md). Measure false positives on a labeled benign/attack-fixture corpus and publish the corpus scope and rates; do not invent detection efficacy from unit-test pass counts.
