import Foundation

// Implements [SECURITY-EVIDENCE]: private local snapshot, explicit staleness, version validation.
public struct ReportStore {
    public let directory: URL
    public init(directory: URL) { self.directory = directory }
    private var file: URL { directory.appendingPathComponent("latest.json") }
    public func save(_ report: ScanReport) throws {
        guard report.schema == 1 else { throw CollectionError.invalid("Unsupported evidence schema") }
        try privateWrite(report, to: file)
    }
    public func load() throws -> ScanReport? {
        guard FileManager.default.fileExists(atPath: file.path) else { return nil }
        var report = try readReport(file)
        report.sections = report.sections.map { section in
            var result = section
            result.state = .stale; result.reason = "Saved snapshot; refresh before use. " + section.reason
            return result
        }
        return report
    }
}
public struct EvidenceChange: Equatable, Identifiable {
    public enum Kind: String, Hashable { case added, removed, changed }
    public let kind: Kind
    public let evidence: Evidence
    public var id: String { "\(kind.rawValue):\(evidence.id)" }
}
public func evidenceChanges(previous: [Evidence], current: [Evidence]) -> [EvidenceChange] {
    let old = Dictionary(previous.map { ($0.id, $0) }, uniquingKeysWith: { first, _ in first })
    let new = Dictionary(current.map { ($0.id, $0) }, uniquingKeysWith: { first, _ in first })
    let updated = current.compactMap { item -> EvidenceChange? in
        guard let prior = old[item.id] else { return EvidenceChange(kind: .added, evidence: item) }
        let differs = prior.path != item.path || prior.title != item.title || prior.detail != item.detail || prior.facts != item.facts
        return differs ? EvidenceChange(kind: .changed, evidence: item) : nil
    }
    return updated + previous.filter { new[$0.id] == nil }.map { EvidenceChange(kind: .removed, evidence: $0) }
}

public func readReport(_ file: URL) throws -> ScanReport {
    guard (try file.resourceValues(forKeys: [.fileSizeKey])).fileSize ?? 0 <= 8_000_000 else { throw CollectionError.invalid("Evidence exceeds storage limit") }
    let report = try JSONDecoder().decode(ScanReport.self, from: boundedData(file))
    guard report.schema == 1 else { throw CollectionError.invalid("Unsupported evidence schema") }
    guard Set(report.sections.map(\.category)).count == report.sections.count,
          report.sections.allSatisfy({ Set($0.records.map(\.id)).count == $0.records.count }) else {
        throw CollectionError.invalid("Evidence contains duplicate identities")
    }
    return report
}
