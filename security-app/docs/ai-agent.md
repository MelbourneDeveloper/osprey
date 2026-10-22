# Optional NIMBLESITE.ai investigation agent [AGENT-INTEGRATION]

**Status: proposed integration contract.** AI is disabled by default and never required for spyware detection, protection rules, evidence explanations, investigation, recovery, export or managed policy. When enabled, the agent uses NIMBLESITE.ai and a user-selected supported model. Native UI uses the [Osprey framework](ui-framework.md); NIMBLESITE's Flutter SDK is not an application dependency.

## User and administrator controls [AGENT-CHOICE]

Offer AI in Settings and finding details as an explicit additional investigation tool. Show provider, exact model, costs, account/credit state, data destination and retention disclosure before activation. Users can change the model and disable AI; organizations may restrict the approved model set or prohibit cloud processing. Discover currently eligible models from a verified catalog rather than hardcoded marketing examples. Never silently switch provider/model after an error, unavailable model or insufficient credit.

Model selection applies per user/investigation, not by patching one shared configuration for everyone. The proposed deployment publishes separately provisioned configurations for approved model choices with identical permitted tool contracts; a trusted configuration service binds the selected config to its provider/model. Prove end-user invocation authorization and catalog authenticity against the deployed API. If the service cannot enforce this safely, model choice is a release blocker, not permission to embed a tenant management key.

A model switch starts a new conversation unless the person explicitly agrees to disclose the existing selected context to the new provider. Display the model used for each answer. Never resend a prior transcript to another provider merely because Settings changed.

## Verified vendor boundary [AGENT-VENDOR]

[NIMBLESITE's model contract](https://www.nimblesite.ai/docs/models/) describes provider/model configuration and live discovery. [Its end-user authentication contract](https://www.nimblesite.ai/docs/multi-tenancy/) provides IdP JWT access and warns against shipping operator keys. [Tool documentation](https://www.nimblesite.ai/docs/tools/) defines client-side calls with typed arguments/results and a server-issued contract hash. Use Swift URLSession transport over that HTTP contract; Osprey owns agent state, disclosure policy, tool authorization and results interpretation.

The authenticated OpenAPI schema must be retrieved, pinned and checked before implementation; public retrieval returned 403 during research. Operator examples using `/api/v1/chat/{config_id}` and an API key are not the distributed-client authentication design. Exact end-user paths, streaming/cancellation semantics and error shapes must come from the deployed schema. No live tenant configuration or billable request has been made for this research.

Provision tenant/config/IdP settings through a separate operator tool. The app uses short-lived end-user tokens stored with Keychain protection and normal expiry/revocation. A conversation ID is a reference, not authorization. Test cross-user, cross-tenant and wrong-config access. Ordinary scanning requires none of these credentials.

## Disclosure and retention [AGENT-PRIVACY]

[NIMBLESITE's privacy policy](https://www.nimblesite.ai/privacy/) discloses hosted model-provider processing, persisted conversations, tenant-operator visibility and limited current deletion facilities. Its hosted offering does not establish production on-device inference or zero-retention processing. Revalidate these terms at release and when providers change. Do not claim data remains local because the tool ran on the device.

Before disclosure, show the exact evidence categories and redacted preview, selected provider/model and retention notice. Default to rule explanations and minimized finding summaries with pseudonymous asset IDs. Keep the local mapping private. File contents, secrets, credentials, full command lines, customer data, browsing URLs and complete process/file inventories are excluded by default. Redaction must run before transport and before diagnostic logging. Typed prompts supplied by users can themselves contain sensitive information; disclose that too.

Each tool request rechecks disclosure consent and managed policy for the requested scope. A consented investigation may share only its chosen evidence set; it cannot silently expand to the entire device. New evidence categories, unredacted values or a different provider require fresh consent. Denial returns a typed tool error and leaves local investigation available. A local delete does not promise deletion of provider or NIMBLESITE history. Ship a tested remote deletion/export process and accurately describe any manual operator/support step before enabling the feature for business data.

## Local tools and action authority [AGENT-TOOLS]

These are proposed application tools, not pre-existing NIMBLESITE tools:

| Tool | Allowed scope |
| --- | --- |
| `get_findings` | Filtered, paginated summaries for the consented device/investigation |
| `get_evidence` | Explicit evidence IDs from that scope, with mandatory redaction and provenance |
| `get_asset_context` | Authorized identity and correlated observations, with coverage gaps |
| `get_rule_explanation` | Versioned deterministic rule rationale, source links and counterexamples |
| `lookup_advisory` | Reviewed local advisory records; an approved feed updater handles network retrieval |
| `request_collection` | A bounded allowlisted read-only collector against an authorized asset; separate disclosure consent governs its results |
| `propose_action` | Creates a local reviewable plan; performs no system change |

There is no arbitrary shell, executable installer, unrestricted file reader, credential reader, generic URL fetcher or direct privileged-service tool. Native services enforce capability scope even when the agent or UI requests more. Tool arguments and results have strict schemas, bounded sizes and validated asset references. Register client execution; a hosted sandbox has no device access and does not become a privileged deputy.

Copy the vendor-issued tool contract hash into requests and reject mismatches. Execute only authenticated structured tool calls, never instructions embedded in answer text. Correlate every call/result and handle duplicate/replayed calls idempotently. Keep local records of completed work so recovering a conversation never reruns an action just to redraw its card. Do not retry a billable turn automatically after an ambiguous timeout.

Use a proposed default limit of eight tool rounds and a bounded evidence/token budget per investigation; expose cancellation and configurable cost limits. Validate real vendor metering and available spend controls before promising a monetary ceiling. Empty balance, timeouts, authentication failures, throttling and schema drift are ordinary AI errors; they cannot stop local protection.

## Evidence-based answers [AGENT-ANSWERS]

Answers separate observed facts, source-backed applicability and hypotheses. Every factual device claim cites an authorized observation ID; every advisory claim cites a reviewed source. Validate cited IDs locally. Unsupported claims are labeled and cannot become deterministic findings or change severity automatically. Markdown is rendered with safe native components; remote images and arbitrary deep links do not load or execute implicitly.

A proposed response is a typed explanation plus evidence references, uncertainty, counterevidence and optional action-plan IDs. A native action card is created only from a valid local plan. Approval occurs in trusted native UI, bound to target identity, operation, expiry and current evidence. The agent cannot mint approval or broaden an approved action. Independent local policy and postcondition checks still apply.

## Adversarial and failure tests [AGENT-VERIFICATION]

[OWASP prompt-injection guidance](https://genai.owasp.org/llmrisk/llm01-prompt-injection/) and [excessive-agency guidance](https://genai.owasp.org/llmrisk/llm062025-excessive-agency/) inform fixtures for malicious filenames, plist values, logs, advisory text and pasted documents. Prove these cannot override tool scope, disclose secrets, approve remediation or hide collection gaps. Delimiters and a system prompt are not the security boundary.

Required tests cover AI-off zero AI traffic; deterministic finding equivalence with AI on/off; user model isolation; provider-change consent; denied disclosure; token expiry/revocation; cross-user access; schema/contract mismatch; hallucinated evidence; replayed tools; exhausted credit; ambiguous network errors; interrupted streams; cancellation; and bounded tool loops. A signed-in user must be able to disable AI and continue a local scan immediately.
