#if os(macOS)
import XCTest
@testable import VigilCore

final class IntegrationTests: XCTestCase {
    private var engine: URL {
        URL(fileURLWithPath: #filePath).deletingLastPathComponent().deletingLastPathComponent()
            .deletingLastPathComponent().appendingPathComponent(".build/engine/vigil-engine")
    }
    func testRealOspreyEngineDistinguishesHiddenFileFromHiddenStartup() throws {
        let record = Evidence(id: "fixture", title: "Fixture", detail: "Test only", source: "test", facts: ["hidden": "true"])
        for (area, expected) in [(Area.files, "information"), (.startup, "hidden-startup")] {
            let result = try assess(ScanSection(category: area, state: .partial, reason: "Test", records: [record]), engine: engine)
            XCTAssertEqual(result.records.first?.assessment, expected)
            XCTAssertEqual(result.state, .partial)
        }
    }
    func testMissingEngineIsNotSilentSuccess() {
        let record = Evidence(id: "fixture", title: "Fixture", detail: "Test", source: "test")
        XCTAssertThrowsError(try assess(ScanSection(category: .apps, state: .partial, reason: "Test", records: [record]), engine: URL(fileURLWithPath: "/nonexistent/vigil-engine")))
    }
    func testInvalidEngineOutputIsRejected() {
        let record = Evidence(id: "fixture", title: "Fixture", detail: "Test", source: "test")
        XCTAssertThrowsError(try assess(ScanSection(category: .apps, state: .partial, reason: "Test", records: [record]), engine: URL(fileURLWithPath: "/usr/bin/true")))
    }
    func testLiveProcessSnapshotContainsThisTestProcess() {
        let section = collect(.processes)
        XCTAssertEqual(section.state, .partial)
        XCTAssertTrue(section.records.contains { $0.facts["pid"] == String(ProcessInfo.processInfo.processIdentifier) }, section.reason)
    }
    func testHiddenScanStaysInSelectedDirectory() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        try FileManager.default.createDirectory(at: root, withIntermediateDirectories: false)
        defer { try? FileManager.default.removeItem(at: root) }
        try Data().write(to: root.appendingPathComponent(".hidden"))
        try Data().write(to: root.appendingPathComponent("visible"))
        try FileManager.default.createSymbolicLink(at: root.appendingPathComponent("outside"), withDestinationURL: URL(fileURLWithPath: "/"))
        let section = collect(.files, folder: root)
        XCTAssertEqual(section.records.map(\.title), [".hidden"])
        XCTAssertEqual(section.state, .partial)
    }
    func testMissingFolderNeverLooksClean() {
        let section = collect(.files, folder: URL(fileURLWithPath: "/nonexistent/vigil-fixture"))
        XCTAssertEqual(section.state, .failed)
        XCTAssertFalse(section.reason.isEmpty)
    }
    func testUnavailableCapabilitiesNeverClaimCoverage() {
        for area in [Area.network, .permissions, .advisories] {
            let section = collect(area)
            XCTAssertEqual(section.state, .unsupported)
            XCTAssertFalse(section.reason.isEmpty)
        }
    }
    func testUnsignedExecutableNeverClaimsSandboxAbsence() throws {
        let url = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        defer { try? FileManager.default.removeItem(at: url) }
        try Data("#!/bin/sh\nexit 0\n".utf8).write(to: url)
        let result = try inspectSignature(Evidence(id: "unsigned", title: "Test", detail: "", path: url.path, source: "test"))
        XCTAssertEqual(result.facts["signing"], "unverified")
        XCTAssertEqual(result.facts["sandbox"], "unknown")
    }
    func testVerifiedSystemAppExposesActualSandboxEntitlement() throws {
        let result = try inspectSignature(Evidence(id: "calculator", title: "Calculator", detail: "", path: "/System/Applications/Calculator.app", source: "test"))
        XCTAssertEqual(result.facts["signing"], "valid")
        XCTAssertEqual(result.facts["sandbox"], "present")
    }
    func testFastOversizedOutputIsRejected() {
        XCTAssertThrowsError(try runCommand("/usr/bin/head", ["-c", "9000000", "/dev/zero"]))
    }
}
#endif
