import XCTest
@testable import VigilCore

final class ReportStoreTests: XCTestCase {
    func testSnapshotRoundTripMarksRestoredEvidenceStale() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        defer { try? FileManager.default.removeItem(at: root) }
        let store = ReportStore(directory: root)
        try store.save(ScanReport(sections: [ScanSection(category: .apps, state: .partial, reason: "Limited scope", records: [])]))
        let restored = try XCTUnwrap(store.load())
        XCTAssertEqual(restored.sections.first?.state, .stale)
        XCTAssertTrue(restored.sections.first?.reason.contains("Limited scope") == true)
        let attributes = try FileManager.default.attributesOfItem(atPath: root.appendingPathComponent("latest.json").path)
        XCTAssertEqual((attributes[.posixPermissions] as? NSNumber)?.intValue, 0o600)
    }
    func testCorruptReportDoesNotBecomeAnEmptySuccessfulScan() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        try FileManager.default.createDirectory(at: root, withIntermediateDirectories: false)
        defer { try? FileManager.default.removeItem(at: root) }
        try Data("broken".utf8).write(to: root.appendingPathComponent("latest.json"))
        XCTAssertThrowsError(try ReportStore(directory: root).load())
    }
    func testNewSchemaIsRejected() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        defer { try? FileManager.default.removeItem(at: root) }
        var report = ScanReport(sections: []); report.schema = 999
        XCTAssertThrowsError(try ReportStore(directory: root).save(report))
    }
    func testImportedDuplicateIdentitiesAreRejected() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        try FileManager.default.createDirectory(at: root, withIntermediateDirectories: false)
        defer { try? FileManager.default.removeItem(at: root) }
        let section = ScanSection(category: .apps, state: .partial, reason: "test", records: [])
        let url = root.appendingPathComponent("duplicate.json")
        try JSONEncoder().encode(ScanReport(sections: [section, section])).write(to: url)
        XCTAssertThrowsError(try readReport(url))
    }
    func testTimestampChangesDoNotProduceEvidenceChanges() {
        let original = Evidence(id: "id", title: "app", detail: "same", source: "test")
        var newer = original; newer.observed = Date().addingTimeInterval(10)
        XCTAssertEqual(evidenceChanges(previous: [original], current: [newer]), [])
        newer.facts["sandbox"] = "absent"
        XCTAssertEqual(evidenceChanges(previous: [original], current: [newer]).map(\.kind), [.changed])
    }
    func testAddedAndRemovedEvidenceHaveStableIdentity() {
        let old = Evidence(id: "old", title: "same name", detail: "", source: "test")
        let new = Evidence(id: "new", title: "same name", detail: "", source: "test")
        XCTAssertEqual(Set(evidenceChanges(previous: [old], current: [new]).map(\.kind)), [.added, .removed])
    }
}
