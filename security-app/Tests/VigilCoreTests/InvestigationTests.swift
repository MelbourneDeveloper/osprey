import XCTest
@testable import VigilCore

final class InvestigationTests: XCTestCase {
    private func record(_ value: String = "one") -> Evidence {
        var result = Evidence(id: "agent", title: "Agent", detail: "Definition", path: "/tmp/agent", source: "fixture", facts: ["version": value])
        result.assessment = "hidden-startup"
        return result
    }
    func testSuppressionExpiresAndCannotHideChangedEvidence() throws {
        let now = Date(timeIntervalSince1970: 100)
        let finding = Finding(area: .startup, evidence: record(), now: now)
        let suppressed = try transition(finding, to: .suppressed, reason: "Approved developer tool", until: now.addingTimeInterval(60), now: now)
        XCTAssertEqual(reconcileFinding(suppressed, evidence: record(), now: now).status, .suppressed)
        XCTAssertEqual(reconcileFinding(suppressed, evidence: record("two"), now: now).status, .new)
        XCTAssertEqual(reconcileFinding(suppressed, evidence: record(), now: now.addingTimeInterval(61)).status, .new)
    }
    func testResolutionRequiresFollowUpEvidenceAndSuppressionNeedsReason() throws {
        let finding = Finding(area: .startup, evidence: record(), now: Date())
        XCTAssertThrowsError(try transition(finding, to: .resolved, reason: "Dismiss", until: nil, now: Date()))
        XCTAssertThrowsError(try transition(finding, to: .suppressed, reason: " ", until: Date().addingTimeInterval(60), now: Date()))
        var clean = record(); clean.assessment = "information"; clean.facts["hidden"] = "false"
        XCTAssertEqual(reconcileFinding(finding, evidence: clean, now: Date()).status, .resolved)
    }
    func testMissingEvidenceInPartialScanDoesNotResolveFinding() {
        let finding = Finding(area: .startup, evidence: record(), now: Date())
        let partial = ScanSection(category: .startup, state: .partial, reason: "Unreadable", records: [])
        let updated = reconcileFindings([finding], section: partial, now: Date())
        XCTAssertEqual(updated.first?.status, .new)
        XCTAssertEqual(updated.first?.id, finding.id)
    }
    func testUnknownSigningDoesNotResolveSignatureFinding() {
        var invalid = record(); invalid.assessment = "signature-review"; invalid.facts["signing"] = "unverified"
        let finding = Finding(area: .apps, evidence: invalid, now: Date())
        var unknown = invalid; unknown.assessment = "information"; unknown.facts["signing"] = "unknown"
        XCTAssertNotEqual(reconcileFinding(finding, evidence: unknown, now: Date()).status, .resolved)
    }
    func testResolvedFindingReopensOnRecurrenceWithoutLosingRuleIdentity() {
        let finding = Finding(area: .startup, evidence: record(), now: Date())
        var clear = record(); clear.assessment = "information"; clear.facts["hidden"] = "false"
        let resolved = reconcileFinding(finding, evidence: clear, now: Date())
        XCTAssertEqual(resolved.status, .resolved)
        XCTAssertEqual(reconcileFinding(resolved, evidence: record(), now: Date()).status, .new)
    }
    func testExpiredSuppressionReopensOnReload() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        defer { try? FileManager.default.removeItem(at: root) }
        let now = Date(), store = FindingStore(directory: root)
        let finding = Finding(area: .startup, evidence: record(), now: now)
        let suppressed = try transition(finding, to: .suppressed, reason: "Temporary exception", until: now.addingTimeInterval(60), now: now)
        try store.save([suppressed], now: now)
        XCTAssertEqual(try store.load(now: now.addingTimeInterval(61)).first?.status, .new)
    }
    func testClearSigningResolvesSigningEvenWhenSandboxStillNeedsReview() {
        var invalid = record(); invalid.assessment = "signature-review"; invalid.facts["signing"] = "unverified"
        let finding = Finding(area: .apps, evidence: invalid, now: Date())
        var valid = invalid; valid.assessment = "sandbox-exposure"; valid.facts["signing"] = "valid"; valid.facts["sandbox"] = "absent"
        XCTAssertEqual(reconcileFinding(finding, evidence: valid, now: Date()).status, .resolved)
        let updated = reconcileFindings([finding], section: ScanSection(category: .apps, state: .partial, reason: "Test", records: [valid]), now: Date())
        XCTAssertEqual(updated.count, 2)
        XCTAssertTrue(updated.contains { $0.ruleID == "sandbox-exposure" && $0.status == .new })
    }
    func testHistoryRetentionAndPrivateStorage() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        defer { try? FileManager.default.removeItem(at: root) }
        let store = HistoryStore(directory: root)
        for seconds in [0.0, 1.0, 2.0] {
            var report = ScanReport(sections: []); report.created = Date(timeIntervalSince1970: seconds)
            try store.save(report, now: Date(timeIntervalSince1970: 3), limit: 2)
        }
        XCTAssertEqual(try store.entries().count, 2)
        XCTAssertEqual(try store.entries().first?.created, Date(timeIntervalSince1970: 2))
        XCTAssertEqual(try FileManager.default.attributesOfItem(atPath: root.path)[.posixPermissions] as? Int, 0o700)
    }
    func testReceiptPreservesProvenanceWithoutClaimingVulnerability() throws {
        let data = Data(#"{"source":{"tap":"vendor/tools","versions":{"stable":"1.2"}},"installed_on_request":true,"poured_from_bottle":false}"#.utf8)
        let evidence = try packageReceipt(data, path: "/opt/homebrew/Cellar/tool/1.2")
        XCTAssertEqual(evidence.facts["tap"], "vendor/tools")
        XCTAssertEqual(evidence.facts["version"], "1.2")
        XCTAssertEqual(evidence.assessment, "information")
        XCTAssertThrowsError(try packageReceipt(Data("[]".utf8), path: "/tmp/tool"))
    }
    func testBaselineDetectsChangesAndRefusesSymbolicLinks() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        try FileManager.default.createDirectory(at: root, withIntermediateDirectories: false)
        defer { try? FileManager.default.removeItem(at: root) }
        let file = root.appendingPathComponent("config"), link = root.appendingPathComponent("link")
        try Data("original".utf8).write(to: file)
        let baseline = try fingerprint(file)
        XCTAssertEqual(baseline.sha256, "0682c5f2076f099c34cfdd15a9e063849ed437a49677e6fcc5b4198c76575be5")
        try Data("changed".utf8).write(to: file)
        XCTAssertNotEqual(try fingerprint(file).sha256, baseline.sha256)
        try FileManager.default.createSymbolicLink(at: link, withDestinationURL: file)
        XCTAssertThrowsError(try fingerprint(link))
        XCTAssertThrowsError(try fingerprint(root))
    }
    func testPermissionGuidanceDoesNotInventGrants() {
        let section = permissionGuidance()
        XCTAssertEqual(section.state, .unsupported)
        XCTAssertTrue(section.records.count >= 5)
        XCTAssertTrue(section.records.allSatisfy { $0.facts["grant"] == "unknown" })
    }
}
