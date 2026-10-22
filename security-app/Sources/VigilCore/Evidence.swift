import Foundation

// Implements [SECURITY-EVIDENCE] and [SECURITY-COVERAGE].
public enum Area: String, Codable, CaseIterable, Identifiable, Sendable {
    case apps, processes, startup, files, protection, packages, network, permissions, advisories
    public var id: Self { self }
    public var title: String { rawValue.capitalized }
}
public enum CollectionState: String, Codable, Sendable {
    case available, partial, denied, unsupported, failed, stale
}
public struct Evidence: Identifiable, Codable, Hashable, Sendable {
    public var id: String
    public var title: String
    public var detail: String
    public var path: String
    public var source: String
    public var observed = Date()
    public var facts: [String: String] = [:]
    public var assessment = "information"
    public init(id: String, title: String, detail: String, path: String = "", source: String, facts: [String: String] = [:]) {
        self.id = id; self.title = title; self.detail = detail; self.path = path; self.source = source; self.facts = facts
    }
}
public struct ScanSection: Identifiable, Codable, Sendable {
    public var id: Area { category }
    public var category: Area
    public var state: CollectionState
    public var reason: String
    public var records: [Evidence]
    public var observed = Date()
    public init(category: Area, state: CollectionState, reason: String, records: [Evidence]) {
        self.category = category; self.state = state; self.reason = reason; self.records = records
    }
}
public struct ScanReport: Codable, Sendable {
    public var schema = 1
    public var created = Date()
    public var sections: [ScanSection]
    public init(sections: [ScanSection]) { self.sections = sections }
}
public enum CollectionError: Error, LocalizedError {
    case invalid(String)
    public var errorDescription: String? {
        switch self { case .invalid(let reason): return reason }
    }
}
