import Foundation

// Implements [SECURITY-EVIDENCE]: a review decision never establishes malware.
public enum FindingStatus: String, Codable, CaseIterable { case new, acknowledged, investigating, resolved, suppressed }
public struct Finding: Codable, Identifiable {
    public let id: String
    public let area: Area
    public let ruleID: String
    public var evidence: Evidence
    public var status = FindingStatus.new
    public let firstSeen: Date
    public var lastSeen: Date
    public var reason = ""
    public var expires: Date?
    public var ruleVersion = "1"
    public init(area: Area, evidence: Evidence, now: Date) {
        id = "\(area.rawValue):\(evidence.id):\(evidence.assessment)"
        ruleID = evidence.assessment
        self.area = area; self.evidence = evidence; firstSeen = now; lastSeen = now
    }
}
public func transition(_ finding: Finding, to status: FindingStatus, reason: String, until: Date?, now: Date) throws -> Finding {
    guard status != .resolved else { throw CollectionError.invalid("Resolution requires a follow-up observation of the same asset.") }
    if status == .suppressed {
        guard !reason.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty,
              let until, until > now, until.timeIntervalSince(now) <= 90 * 86400 else {
            throw CollectionError.invalid("Suppression needs a reason and expiry within 90 days.")
        }
    }
    var result = finding
    result.status = status; result.reason = reason; result.expires = status == .suppressed ? until : nil
    return result
}
public func reconcileFinding(_ finding: Finding, evidence: Evidence, now: Date) -> Finding {
    var result = finding
    let changed = finding.evidence.path != evidence.path || finding.evidence.facts != evidence.facts
    guard evidence.assessment == finding.ruleID || resolves(finding.ruleID, evidence) else { return finding }
    result.evidence = evidence; result.lastSeen = now
    if resolves(finding.ruleID, evidence) { result.status = .resolved; result.expires = nil }
    else if changed || finding.status == .resolved || (finding.status == .suppressed && (finding.expires ?? .distantPast) <= now) {
        result.status = .new; result.reason = "Evidence changed, review resumed or suppression expired"; result.expires = nil
    }
    return result
}
public func reconcileFindings(_ existing: [Finding], section: ScanSection, now: Date) -> [Finding] {
    guard section.state == .available || section.state == .partial else { return existing }
    let observed = Dictionary(section.records.map { ($0.id, $0) }, uniquingKeysWith: { first, _ in first })
    let updated = existing.map { finding in
        guard finding.area == section.category, let evidence = observed[finding.evidence.id], evidence.assessment != "unassessed" else { return finding }
        return reconcileFinding(finding, evidence: evidence, now: now)
    }
    let added = section.records.filter { $0.assessment != "information" && $0.assessment != "unassessed" }.map {
        Finding(area: section.category, evidence: $0, now: now)
    }.filter { candidate in !updated.contains { $0.id == candidate.id } }
    return updated + added
}

private func resolves(_ rule: String, _ evidence: Evidence) -> Bool {
    guard evidence.assessment != "unassessed" else { return false }
    switch rule {
    case "signature-review": return evidence.facts["signing"] == "valid"
    case "sandbox-exposure": return evidence.facts["signing"] == "valid" && evidence.facts["sandbox"] == "present"
    case "hidden-startup": return evidence.facts["hidden"] == "false"
    case "protection-disabled": return evidence.facts["disabled"] == "false"
    default: return false
    }
}

public func expireSuppression(_ finding: Finding, now: Date) -> Finding {
    guard finding.status == .suppressed, (finding.expires ?? .distantPast) <= now else { return finding }
    var result = finding
    result.status = .new; result.expires = nil; result.reason = "Suppression expired; review the observation again"
    return result
}
