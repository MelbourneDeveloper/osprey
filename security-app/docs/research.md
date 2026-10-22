# Research and authoritative sources

**Research date: 22 October 2026.** Quotations below are short excerpts from the linked primary sources. Design decisions are ours; a reference establishes the cited technique or platform constraint, not certification or a guarantee of detection. Recheck Apple deployment tables, entitlements, SDK support and vendor contracts at each release.

## Distribution and operating-system boundaries

### macOS

[Apple's App Review Guidelines](https://developer.apple.com/app-store/review/guidelines/), §2.4.5, require Mac App Store applications to be “appropriately sandboxed” and prohibit requesting “escalation to root privileges”. Sections 2.5 and 5.4 also govern public APIs, background behavior and VPN services. These constraints conflict with parts of a comprehensive host-monitoring design. Decision: prioritize direct distribution, keep the UI unprivileged, and give each separately signed service only its required capabilities.

[Apple's System Extensions guidance](https://developer.apple.com/system-extensions/) says EndpointSecurity use requires requesting an entitlement. Direct distribution does not remove that dependency, user approval, privacy permissions or SIP protections. [Apple DTS stated in 2019](https://developer.apple.com/forums/thread/125508) that “EndpointSecurity system extensions cannot be deployed via the Mac App Store.” This is historical Apple staff guidance, not a current blanket store ruling. Treat ES as direct-distribution scope and reconfirm current eligibility with Apple before promising otherwise.

[TN3134](https://developer.apple.com/documentation/Technotes/tn3134-network-extension-provider-deployment) is the deployment authority for Network Extension provider types. It describes content-filter system extensions on macOS and different requirements for packet tunnels and newer URL filters. A network extension and an ES extension are distinct targets with distinct privileges. Prove signing, installation, user approval and recovery on stock devices with SIP enabled.

### iOS and iPadOS

[Apple Platform Security](https://support.apple.com/en-au/guide/security/sec15bfe098e/web) states: “All third-party apps are sandboxed.” This prevents a consumer app from providing unrestricted app inventory, other apps' permission databases, system-wide process inspection, filesystem scanning or general system-file verification. Present unavailable coverage explicitly, including when the installation channel changes.

TN3134 documents supervised-device, managed-app and Screen Time child-device restrictions for iOS content filters. It also documents URL-filter providers introduced with iOS/macOS 26; these are a separate feasibility investigation, not evidence of a general consumer per-app firewall. Packet tunnels require a real networking design and conflict/coexistence tests with existing VPNs. Do not route traffic through a fake VPN merely to claim monitoring.

[Apple's personal-safety guide](https://support.apple.com/guide/personal-safety/app-privacy-features-in-apple-products-ipsd4dac4a2c/1.0/web/1.0) describes App Privacy Report, including sensor/data access and network activity. Initial iOS scope is guidance and analysis of explicitly supplied reports/files; report export format and completeness must be verified on supported devices. There is no assumed live API into that report.

### Android

[Android package visibility](https://developer.android.com/training/package-visibility) limits installed-app queries. [VpnService](https://developer.android.com/reference/android/net/VpnService) requires user authorization and permits only one granted VPN application at a time. Shared UI and Osprey domain logic are feasible; collection parity and Google Play permission eligibility are separate work. Android is not a way to infer unavailable Apple observations.

## Detection sources and their use

| Authority and source | Short quotation or documented fact | Technique and product implication |
| --- | --- | --- |
| [NIST SP 800-137](https://csrc.nist.gov/pubs/sp/800/137/final) | “visibility into organizational assets” | Inventory, change tracking, collection health and explicit scope precede a protection claim. Adapt the methodology to personal devices; do not imply federal compliance. |
| [NIST SP 800-61 Rev. 3](https://csrc.nist.gov/pubs/sp/800/61/r3/final) | “Incident Response Recommendations and Considerations” | Separate detection, investigation, response and recovery; record action outcomes and preserve evidence before changes. |
| [MITRE ATT&CK T1543.001](https://attack.mitre.org/techniques/T1543/001/) | “Adversaries may create or modify launch agents” | Inspect launch-agent definitions and executable identity; correlate new persistence with executable location, signing and observed behavior. Legitimate updaters are negative fixtures. |
| [MITRE ATT&CK T1564.001](https://attack.mitre.org/techniques/T1564/001/) | “Many applications create these hidden files and folders” | Inspect dotfiles and hidden flags, but require additional evidence before a suspicious finding. `.git`, `.ssh` and normal caches are negative controls. |
| [MITRE ATT&CK T1041](https://attack.mitre.org/techniques/T1041/) | “exfiltrating it over an existing command and control channel” | Correlate process identity, destination, timing and transfer metadata. TLS metadata does not prove payload contents or theft. |
| [CISA KEV catalog](https://www.cisa.gov/known-exploited-vulnerabilities-catalog) | “an input to their vulnerability management prioritization framework” | Prioritize a confirmed installed-version vulnerability with known exploitation; presence in KEV does not prove this device is compromised. |
| [OSV API](https://google.github.io/osv.dev/api/) | Queries accept project commit hashes or versions. | Match exact ecosystem/package identity and affected versions where covered; retain aliases and withdrawn records. Missing ecosystem coverage is unknown, not safe. |
| [Apple Code Signing Services](https://developer.apple.com/documentation/security/code-signing-services/) | “Performs dynamic validation of signed code.” | Use static/dynamic validation and signing information appropriately. Preserve error codes; distinguish unsigned, invalid, ad hoc and unavailable. A valid signature is provenance evidence, not innocence. |
| [Apple signed system volume](https://support.apple.com/guide/security/signed-system-volume-security-secd698747c9/1/web/1) | “verifies the integrity of the system content at runtime” | Report SSV context separately from writable-data-volume changes. Compare only compatible OS-build baselines; do not label normal updates tampering. |
| [Apple SIP](https://support.apple.com/en-my/guide/security/secb7ea06b49/web) | Describes protected system locations and restrictions. | Collect supported protection-state evidence where available; never require disabling protection to install or scan. |
| [Apple malware protections](https://support.apple.com/en-nz/guide/security/sec469d47bd8/web) | Describes Gatekeeper, notarization and XProtect. | Complement built-in controls. Report observed checks precisely; do not claim that a passed notarization check is an antivirus verdict. |
| [Homebrew manual](https://docs.brew.sh/Manpage) | Describes JSON inventory, outdated packages, taps and services. | Capture installed receipts and versions, provenance, service state and available updates. `brew audit` checks packaging conventions; `brew doctor` checks environment health. Neither establishes vulnerability absence. |
| [FTC Zoom settlement announcement](https://search.ftc.gov/news-events/news/press-releases/2020/11/ftc-requires-zoom-enhance-its-security-practices-part-settlement) | “to settle allegations” | Example of authoritative, dated privacy/security history. Label allegations and orders accurately, preserve remediation context and affected dates; do not classify today's installed app as malicious from a historical incident. |

ATT&CK IDs describe adversary behavior. They are neither signatures nor proof that the observed behavior is malicious. New persistence or hidden files alone merit inventory/change information; corroboration determines escalation.

### Anti-spyware and business protection sources

[CIS Control 10](https://www.cisecurity.org/controls/malware-defenses) says: “Prevent or control the installation, spread, and execution of malicious applications, code, or scripts on enterprise assets.” This establishes a protection objective beyond inventory. The product must verify its preventive capabilities and operational recovery before a professional anti-spyware launch; mapping a control is not compliance certification or a CIS recommendation of this product.

[MITRE keylogging T1056.001](https://attack.mitre.org/techniques/T1056/001/) describes intercepting keystrokes and credentials; [screen capture T1113](https://attack.mitre.org/techniques/T1113/) describes desktop collection; [audio capture T1123](https://attack.mitre.org/techniques/T1123/) describes microphone/application recording. These establish threat scenarios. They do not establish that Apple exposes a universal sensor-use monitoring API to this app. Each advertised observation requires a supported, tested collector and legitimate-use negative fixtures.

[Apple's mercenary-spyware guidance](https://support.apple.com/en-us/102174) notes that targeted attacks are “much harder to detect and prevent.” Provide high-risk escalation and available system-protection guidance rather than promising that consumer iOS inspection detects or excludes every spyware infection.

## Permissions and background activity

[App Sandbox documentation](https://developer.apple.com/documentation/security/protecting-user-data-with-app-sandbox) describes restrictions on other apps' preferences, terminating other apps and protected resources. An entitlement is a declared capability, a privacy grant is an authorization decision, and actual use is an observation: collect and display these separately. No documented general public API for auditing all other apps' TCC grants was established in this research. A consented direct-Mac TCC database adapter is therefore experimental, read-only, version-specific and optional; never edit that database.

[NSWorkspace.runningApplications](https://developer.apple.com/documentation/appkit/nsworkspace/runningapplications) supplies running application observations, not a complete process inventory. [Apple's launch-notification documentation](https://developer.apple.com/documentation/appkit/nsworkspace/didlaunchapplicationnotification) notes omissions for background apps and recommends observing `runningApplications` for that use. Combine app observations with verified process collection; do not build background detection from launch notifications alone. Requested, configured, running, hidden, inactive and suspended are different states.

## Framework research and decision

| Candidate | Evidence | Assessment for this product |
| --- | --- | --- |
| Osprey framework with SwiftUI host | [SwiftUI](https://developer.apple.com/swiftui/), [model data](https://developer.apple.com/documentation/swiftui/model-data), [navigation](https://developer.apple.com/documentation/swiftui/navigation), [scenes](https://developer.apple.com/documentation/swiftui/scenes), [forms](https://developer.apple.com/documentation/swiftui/form) | Selected: evolve the existing storage-app pattern into a reusable native application framework. Osprey owns composition and application semantics; SwiftUI owns native control behavior, with AppKit/UIKit integration where required. |
| Future Android native host | Existing [Kotlin renderer](https://github.com/Nimblesite/osprey/blob/main/examples/mobile/android/app/src/main/java/org/ospreylang/issueinbox/NativeRenderer.kt) and [Android target contract](https://github.com/Nimblesite/osprey/blob/main/docs/specs/platforms/0039-AndroidTarget.md) | Bonus target: share Osprey state, rules, component semantics and events through a Kotlin native host. Android collector capabilities remain independent. |
| Flutter / React Native / web UI | Evaluated only as alternatives. | Excluded by the product requirement for a native Mac/iOS application following the storage app. No Dart layer or JavaScript runtime is planned. |

The storage app already supplies the Osprey/C ABI/SwiftUI pattern. It does not establish a complete reusable application framework. [Framework acceptance](ui-framework.md#acceptance-contract-ui-verification) must prove composition, identity, bindings, navigation, lifecycle and the complete interaction model rather than merely recognizing control names.

## NIMBLESITE research

[NIMBLESITE's model documentation](https://www.nimblesite.ai/docs/models/) says “Pick the model in `model_config`.” It documents provider/model discovery and model configuration changes. [Its integration introduction](https://www.nimblesite.ai/docs/) links an OpenAPI contract and a Flutter SDK. [Tool documentation](https://www.nimblesite.ai/docs/tools/) describes typed client-executed tools, contract hashes and correlated results. These support an optional agent whose actual device capabilities stay local.

[The authentication documentation](https://www.nimblesite.ai/docs/multi-tenancy/) states: “Keys are for servers and operators. Never ship one in an app.” Its end-user JWT integration is the intended consumer boundary. [The privacy policy](https://www.nimblesite.ai/privacy/) states: “There is no "don't store" mode today, and no automatic expiry.” It also discloses third-party model processing, tenant-operator access, and currently limited deletion facilities. Provider and retention disclosures are required before sending evidence; client-side tool execution does not make returned evidence private to the device.

The privacy page lists more providers than the examples on the models page. Discover availability at runtime and verify the deployed contract; do not copy either page into a permanent model picker. Production local inference and zero-retention operation are not established capabilities. The public `/openapi.json` request returned HTTP 403 from this research environment; authenticated schema retrieval and SDK/API consistency are explicit feasibility work, not verified integration.

[OWASP LLM01](https://genai.owasp.org/llmrisk/llm01-prompt-injection/) explains that retrieved content can redirect model behavior. [OWASP LLM06](https://genai.owasp.org/llmrisk/llm062025-excessive-agency/) addresses excessive tool authority. Treat scanned filenames, metadata, logs, reports, websites and agent replies as untrusted input. Enforcement belongs in deterministic local policy and privileged services, outside the model.

## Evidence still required

- Current ES distribution approval, permission behavior, supported event coverage and event-loss recovery on the intended macOS releases.
- macOS network attribution and filtering under VPNs, QUIC, encrypted DNS, shared hosting and process restarts; iOS consumer eligibility for each provider type.
- Versioned TCC/permission visibility and lawful App Privacy Report import formats on physical devices.
- Authenticated NIMBLESITE schemas, per-user model configuration authorization, deletion/export workflow, provider terms, quotas and spend limits.
- Bounded native memory and string lifetime for long-running Osprey embedding; native framework ABI tests on Mac and iOS, with Android as a later adapter.

These are implementation gates with named deliverables in [the plan](plan.md). None authorizes a claim of complete device coverage.

## Verified implementation decisions — 22 October 2026 [RESEARCH-IMPLEMENTATION]

| Decision | Primary evidence and technique | Implementation / qualification |
| --- | --- | --- |
| Launch persistence | MITRE [T1543.001](https://attack.mitre.org/techniques/T1543/001/) says adversaries “may create or modify launch agents”. It also documents legitimate updater and developer uses. | Read `launchd` plists as data; correlate hidden target context, retain source and never infer execution or malicious intent from configuration alone. |
| Privacy review | Apple's [Privacy & Security guide](https://support.apple.com/guide/mac-help/change-privacy-security-settings-on-mac-mchl211c911f/mac) documents system controls for sensitive access. | Guidance links the user to System Settings. Usage descriptions and signing entitlements are distinct from privacy grants; no private TCC write is attempted. |
| Package provenance | Homebrew's [Formula Cookbook](https://docs.brew.sh/Formula-Cookbook) defines a tab as “information about a keg” and identifies `INSTALL_RECEIPT.json`. | Parse bounded receipt JSON without executing Ruby or `brew`. Tap and version are local claims, not authenticated upstream identity or a CVE match. |
| Integrity | Apple says SIP “restricts the root user account” in [its platform explanation](https://support.apple.com/en-us/102149). | Protection settings and selected writable-file comparisons are separate observations. SHA-256 detects byte changes relative to a recorded observation; it cannot establish an initially clean machine or validate the sealed OS. |
| Advisory research | [CISA's official KEV mirror](https://github.com/cisagov/kev-data) calls itself “Mirror of cisa.gov/kev data files” and publishes JSON/schema. | Explicit download uses that fixed HTTPS origin, without device upload. The catalog is research evidence. No name-based installed-version matching, signed rule distribution or complete vulnerability coverage is claimed. |
| Network inspection | Apple's [Network Extension entitlement](https://developer.apple.com/documentation/bundleresources/entitlements/com.apple.developer.networking.networkextension) describes provider capabilities and Developer ID provisioning. The installed macOS `lsof(8)` documents field output and numeric address/port switches. | The current collector reads numeric open-socket metadata using `lsof -nP -i -FpcfnPtT`. It is not a provider, exfiltration detector or prevention mechanism. A loopback-listener test checks the actual adapter. |
| Agent credentials | [NIMBLESITE end-user authentication](https://www.nimblesite.ai/docs/multi-tenancy/) defines BYO-IdP JWT calls and prohibits client operator keys. | Native transport uses `/api/v1/end-user/chat/{config_id}` with a session-only Bearer token. Imported config labels do not independently attest server model selection. |
| Agent privacy | [NIMBLESITE privacy](https://www.nimblesite.ai/privacy/) documents provider processing, persisted history and tenant access. | User-authored question only, approval per question/model, no auto-sharing of files or evidence, no local tool execution, no automatic retry after ambiguous failure. Remote deletion remains a deployment/support obligation. |

The local vendor SDK OpenAPI artifact was inspected in addition to the public docs. SHA-256: `7011397cf3445e59050f44362737db450469dc512f7dc05882c03235f25fb5d5`. It establishes the implemented request/response shape, not proof that an authenticated production tenant has been provisioned or that every deployment serves the identical schema. Live public schema retrieval remained unavailable. End-user auth, cross-tenant isolation, actual model/config binding and billable behavior must be tested against the chosen deployment before enabling production AI.

Mac baseline: macOS 14+, Apple silicon, direct local development distribution. iOS baseline: iOS 17+, native report/assistant companion. Endpoint Security eligibility, provider provisioning and store approval are release evidence to obtain from Apple, not technical claims inferred from having an SDK. The [delivery plan](plan.md) preserves these acceptance gates.
