import SwiftUI

// Implements [DESIGN-STATUS]. Coverage and severity are separate domains.
public enum Coverage: String, CaseIterable, Identifiable {
    case available, partial, denied, unsupported, failed, stale
    public var id: Self { self }
    public var title: String { rawValue.capitalized }
    public var symbol: String {
        switch self {
        case .available: return "checkmark.circle.fill"
        case .partial: return "circle.lefthalf.filled"
        case .denied: return "hand.raised.fill"
        case .unsupported: return "minus.circle"
        case .failed: return "exclamationmark.circle.fill"
        case .stale: return "clock.badge.exclamationmark"
        }
    }
    public var ink: Ink {
        switch self {
        case .available: return .good
        case .partial, .stale: return .caution
        case .failed: return .danger
        case .denied, .unsupported: return .secondary
        }
    }
}

public enum Severity: String, CaseIterable, Identifiable {
    case information, review, high, critical
    public var id: Self { self }
    public var title: String { rawValue.capitalized }
    public var ink: Ink {
        switch self {
        case .information: return .secondary
        case .review: return .caution
        case .high, .critical: return .danger
        }
    }
    public var symbol: String {
        switch self {
        case .information: return "info.circle"
        case .review: return "magnifyingglass"
        case .high: return "exclamationmark.triangle"
        case .critical: return "exclamationmark.octagon.fill"
        }
    }
}

public struct StatusBadge: View {
    @Environment(\.colorScheme) private var scheme
    private let title: String
    private let symbol: String
    private let ink: Ink
    public init(_ state: Coverage) { title = state.title; symbol = state.symbol; ink = state.ink }
    public init(_ severity: Severity) { title = severity.title; symbol = severity.symbol; ink = severity.ink }
    public var body: some View {
        Label(title, systemImage: symbol).font(.caption.weight(.semibold))
            .foregroundStyle(ink.color(scheme)).padding(.horizontal, 9).padding(.vertical, 6)
            .background(Ink.inset.color(scheme), in: Capsule())
            .accessibilityElement(children: .combine)
    }
}

public struct EvidenceLine: View {
    @Environment(\.colorScheme) private var scheme
    let label: String
    let value: String
    public init(_ label: String, _ value: String) { self.label = label; self.value = value }
    public var body: some View {
        VStack(alignment: .leading, spacing: 5) {
            Text(label).font(.caption).foregroundStyle(Ink.secondary.color(scheme))
            Text(value).font(.callout).textSelection(.enabled).fixedSize(horizontal: false, vertical: true)
        }.frame(maxWidth: .infinity, alignment: .leading).accessibilityElement(children: .combine)
    }
}
