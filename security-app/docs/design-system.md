# Vigil: native design system [DESIGN-SYSTEM]

**Delivered:** a [native SwiftUI design package and interactive Mac gallery](../README.md), source logo artwork, semantic foundations and illustrative security screens. **Proposed:** production integration through the Osprey framework, final store icon production and complete device/accessibility validation. Vigil is a provisional name, not a cleared trademark.

## Identity [DESIGN-IDENTITY]

The visual idea is an aperture: making concealed behaviour visible while keeping attention on the evidence. Two ivory architectural brackets frame a mint centre on graphite. The mark is deliberately distinct from the Osprey language identity. The wordmark uses lowercase system rounded type with restrained tracking; product typography uses system text and a serif display accent.

Keep one mark-width of clear space around standalone branding. The raster artwork is square and unmasked; native icon tools apply platform shapes. Never put severity badges inside the logo, animate the mark continuously, or use the mint centre to imply that the device is safe. The app icon is identity, not a monitoring indicator. Produce small-size optical variants before store submission; the current raster is suitable for the gallery and identity review.

## Foundations [DESIGN-FOUNDATIONS]

The executable source is [Theme.swift](../Sources/VigilDesign/Theme.swift). Public roles are `canvas`, `surface`, `inset`, `primary`, `secondary`, `accent`, `good`, `caution`, `danger` and `line`. Resolve roles against appearance; do not let app features select raw RGB values. The current native package is the first backend implementation. A shared Osprey token contract must consume an extracted single source when the broader framework ships, rather than duplicating colour constants.

The palette pairs warm ivory with deep green graphite. Muted botanical accents identify selection and available collection; amber calls for review; brick red identifies high severity. Low elevation, thin borders and generous whitespace create hierarchy. Text carries the meaning, with symbols as reinforcement. Disabled control opacity is never used for ordinary explanatory text.

| Foundation | Rule |
| --- | --- |
| Typography | System body, callout, headline and caption roles; serif large title only for overview/identity; monospaced evidence identifiers. Never shrink evidence to make it fit. |
| Spacing | 4, 8, 12, 20, 28, 36 points for reusable primitives; native layouts may resolve their own platform margins. |
| Shapes | 16-point surfaces, 14-point findings, 10-point buttons; capsules only for short status labels. |
| Elevation | Solid opaque surfaces. Borders establish grouping without requiring transparency or shadow. |
| Hit areas | Custom action buttons have a minimum 44-point height. Native controls keep platform behaviour and focus treatment. |
| Selection | Surface tint plus visible outline and accessibility selected trait. Never colour alone. |
| Layout | Mac sidebar/content/inspector; findings and inspector stack when the width cannot accommodate both. Text wraps; the whole screen never requires horizontal scrolling. |

Contrast tests check all six semantic text roles against all three surfaces in both appearances at 4.5:1 or better, plus primary button text. These numeric checks complement visual/device tests; they do not certify every native system control or all accessibility settings.

## Evidence as the interaction model [DESIGN-EVIDENCE]

The core interaction connects an observation to its source, time, uncertainty and next step. Selecting a finding changes the adjacent evidence inspector. The user should not have to open a modal merely to understand why a finding exists. Narrow layouts stack the detail; production iPhone routing will use the framework's native navigation contract.

Each finding carries severity independently of confidence and independently of collection coverage. The illustrative gallery labels every screen as demonstration data, makes no observations about the user's device and performs no protection changes. Its session-only note field exists to exercise controls, not to store an investigation.

## Status grammar [DESIGN-STATUS]

| Domain | States | Interpretation |
| --- | --- | --- |
| Collection coverage | Available, partial, denied, unsupported, failed, stale | Whether observations can support a conclusion; available does not mean threat-free. |
| Finding severity | Information, review, high, critical | Consequence and priority; not certainty. |
| Confidence | Label with explanation and source | Must be provided by the evidence workflow; never inferred from the severity colour. |
| Task lifecycle | Planned: idle, collecting, cancelling, completed, failed | Real progress only; no invented scan completion percentages. |

Coverage segments represent individual collection areas, not a security score. Visible text explains missing areas. A production coverage view must identify each area and its age, with an accessible equivalent. The gallery demonstrates the visual grouping and aggregate label; a generic percentage of “protection” is forbidden.

## Component contracts [DESIGN-COMPONENTS]

| Component | Behaviour |
| --- | --- |
| `Surface` | Composable native content, consistent padding and border, no implicit clipping of long text. |
| `StatusBadge` | Separate typed initializers for coverage and severity; combined accessible text and icon. |
| `EvidenceLine` | Wrapping, selectable value and source label; identifiers never silently truncated. |
| `CoverageStrip` | Bounded collection-area segments with textual aggregate accessibility label. |
| `VigilButtonStyle` | Native button action/focus semantics, disabled treatment, 44-point height and restrained press feedback. |
| `Brand` | Resource-backed mark plus readable wordmark; one accessibility name. |
| Finding row | Full-row native button, stable identity, visible selected state; opens evidence without performing a security action. |
| Gallery form | Native text field, dirty-state feedback, empty-note save disabled, session-only state. |

Production additions must cover loading, empty, denied, stale and failure content slots, confirmation sheets for approved actions, and cancellable work. The gallery's component subset does not satisfy the entire application-framework acceptance contract. App state, real evidence, authorization and service effects remain Osprey responsibilities; this package owns native visual interpretation only.

## Motion [DESIGN-MOTION]

Motion explains a transition and then stops. Press feedback uses 120 ms ease-out; selection/detail changes use 220 ms ease-in-out. Detail demonstrations crossfade. There is no radar animation, perpetual logo pulse, confetti after a scan, or artificial activity noise. Critical observations remain readable rather than flashing.

The system Reduce Motion preference always wins. The gallery can additionally force reduced motion; it can never turn system reduction off. Under reduction, custom animations resolve to `nil`, with immediate state changes. Native OS components retain their own accessibility-aware behaviour. Animation is never required to discover a finding, confirm an action or understand a coverage gap.

## Native accessibility and release checks [DESIGN-VERIFICATION]

The system follows Apple's [motion guidance](https://developer.apple.com/design/human-interface-guidelines/motion) and [accessibility guidance](https://developer.apple.com/design/human-interface-guidelines/accessibility). Keep system text styles, native editing/selection, keyboard semantics and labelled controls. Decorative icons are hidden from accessibility; evidence labels and values are combined. Rendering is opaque by default, so Reduce Transparency needs no custom glass substitute.

Before production integration, verify keyboard-only navigation and focus visibility, VoiceOver order and announcements, Dynamic Type including accessibility sizes on physical iPhone/iPad, increased contrast, RTL and long localized strings, multiwindow state, IME editing during observation updates, and repeated open/close/cancel interactions. Audit press/selection motion on a live device. Screenshots and contrast calculations alone cannot prove these behaviours.

The current verification covers a Mac build, iOS library compilation, contrast tests and native rendered layout inspection at desktop and 390-point widths. Live interaction/accessibility testing was blocked by the locked Mac. The generated logo requires final production exports and identity clearance. No security efficacy, store approval, or expert endorsement follows from completing this design system.
