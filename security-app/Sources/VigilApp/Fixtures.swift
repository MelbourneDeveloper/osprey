import SwiftUI
import VigilDesign

// Synthetic observations only; this gallery never inspects the user's device.
struct StudioFinding: Identifiable, Hashable {
    let id: String
    let title: String
    let app: String
    let symbol: String
    let severity: Severity
    let detail: String
    let evidence: String
}

let findings = [
    StudioFinding(id: "capture", title: "Screen capture needs a closer look", app: "Atlas Helper",
            symbol: "rectangle.dashed.badge.record", severity: .high,
            detail: "A background helper appears in the screen-capture fixture. Confirm that this access is expected before deciding what to change.",
            evidence: "Synthetic capture observation · 10:42:18"),
    StudioFinding(id: "startup", title: "A new item starts at login", app: "Northstar Sync",
            symbol: "arrow.clockwise.circle", severity: .review,
            detail: "A launch agent was added after the previous inventory. Background activity is common in sync tools; check the publisher and purpose.",
            evidence: "Synthetic launch-agent change · 10:38:04"),
    StudioFinding(id: "network", title: "An unfamiliar connection", app: "Atlas Helper",
            symbol: "network", severity: .review,
            detail: "This fixture connects a background process to an unfamiliar destination. A connection alone does not establish data theft.",
            evidence: "collector.example.invalid:443 · synthetic flow")
]

enum Section: String, CaseIterable, Identifiable {
    case overview = "Overview", components = "Design library", identity = "Identity & motion"
    var id: Self { self }
    var symbol: String {
        switch self {
        case .overview: return "square.grid.2x2"
        case .components: return "square.stack.3d.up"
        case .identity: return "circle.hexagongrid"
        }
    }
}
