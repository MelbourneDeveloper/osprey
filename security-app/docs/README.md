# Security application: research and delivery dossier

**Status: proposed; researched 22 October 2026.** This is the starting specification for a separate security product repository. The research documents do not deliver a security application, complete Osprey UI framework, Apple entitlement approval, or production AI integration. A native design adapter and gallery are now available separately. Vigil is the provisional visual identity; final product naming remains open.

This is a professional anti-spyware product for individuals and businesses, with prevention, detection, investigation and recovery. It explains what runs on a person's device, what access it has, what it communicates with, and what evidence warrants attention. Osprey owns evidence correlation, detection rules, policy, application state and workflows. Platform collectors provide observations through explicit effects; deterministic handlers make those decisions testable. Expert recommendations are an external validation goal, never a claim made before an assessment.

## Decisions recorded from the brief

- Production and App Store readiness are required quality standards. Full Mac protection is planned through a Developer ID signed, hardened and notarized application where necessary privileges cannot ship through the store. Eligible native editions target actual App Store review; their capability profile must retain useful, honestly described protection. A store decision cannot silently remove required protections or imply parity with privileged Mac monitoring.
- iOS is a companion with explicitly narrower collection capabilities. Alternative distribution does not remove iOS sandbox restrictions. Android is a planned additional UI target, with independently researched collectors.
- AI is an optional extra, disabled initially. The full non-AI product works without an account, AI credits, model availability or a NIMBLESITE connection. Threat-feed updates are separate from AI.
- **NIMBLESITE.ai is the required integration when AI is enabled.** Users select the supported provider and model; a model is never silently substituted. Local production inference is not currently established by NIMBLESITE's published contract.
- Build an ergonomic, strongly typed Osprey application framework with a public model independent of SwiftUI, evolving the storage app's pattern. Its contracts include lifecycle, typed state/actions/bindings, routing, composable screens, forms, commands, accessibility, themes, persistence, errors, diagnostics and tests. SwiftUI is the first backend for native macOS/iOS; AppKit/UIKit supply platform integration. Flutter, a web shell and a widget-kind switch presented as a framework are excluded.
- Reuse and evolve the storage app's domain/host separation, native renderer, collection techniques, C ABI and tests after a security review. Extract shared infrastructure once; do not fork its renderer or copy the compiler into the new repository. A future Android renderer may share the Osprey framework contract; it does not dictate Apple UI choices.

## Documents

| Document | Owns |
| --- | --- |
| [Design system and native gallery](design-system.md) | Logo, native visual components, semantic tokens, motion and verification boundaries |
| [Research and source register](research.md) | Platform/distribution feasibility, authoritative quotations, technique selection and unresolved evidence |
| [Product and security specification](specification.md) | Capability matrix, findings, collectors, actions, evidence, privacy and acceptance requirements |
| [Application and UI framework](ui-framework.md) | Backend-independent Osprey API, native host boundaries, reusable packages, interaction and quality contracts |
| [Optional NIMBLESITE agent](ai-agent.md) | Model choice, authentication, evidence disclosure, tools, costs and AI-specific verification |
| [Implementation and extraction plan](plan.md) | Dependency-ordered milestones, ownership, tests, feasibility gates and release criteria |

These documents belong together and move into the product repository at extraction. They are product proposals rather than additions to the implemented Osprey language specification. Requirement IDs are descriptive slugs; future implementation and tests cite them with `Implements [ID]` and `Verifies [ID]` comments.

## First delivery

Prove the platform and memory boundaries, then deliver a read-only Mac inventory with permission visibility limits, signed-app identity, background activity, persistence and scoped file inspection. Show the evidence and collection gaps for every result. Add vulnerability matching and network observation after those foundations, then reversible network rules. Deliver the optional AI investigation alongside verified evidence, never as the detection engine.

The framework must prove itself with two consumers: a storage-browser feature and a security-inventory feature. A handful of mapped controls, a polished screenshot, or a compiled application is insufficient acceptance.
