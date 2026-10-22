import XCTest
@testable import VigilCore

final class CollectorTests: XCTestCase {
    func testLaunchAgentReadsProgramWithoutExecutingIt() throws {
        let data = try PropertyListSerialization.data(fromPropertyList: [
            "Label": "test.agent", "ProgramArguments": ["/tmp/.hidden/agent", "--token", "secret"], "RunAtLoad": true
        ], format: .xml, options: 0)
        let record = try decodeLaunchItem(data, path: "/tmp/test.plist")
        XCTAssertEqual(record.path, "/tmp/.hidden/agent")
        XCTAssertEqual(record.facts["hidden"], "true")
        XCTAssertFalse(record.detail.contains("secret"))
    }
    func testUnrecognizedEntitlementOutputNeverMeansSandboxAbsent() {
        XCTAssertThrowsError(try decodeEntitlements("unrecognized signing output"))
    }
    func testMalformedLaunchItemIsAnError() {
        XCTAssertThrowsError(try decodeLaunchItem(Data("garbage".utf8), path: "bad"))
    }
    func testMissingExecutableIsNotReportedAsSuccessfulInventory() throws {
        let data = try PropertyListSerialization.data(fromPropertyList: ["Label": "bad"], format: .xml, options: 0)
        XCTAssertThrowsError(try decodeLaunchItem(data, path: "bad"))
    }
    func testCommandNonzeroAndTimeoutRemainErrors() throws {
        XCTAssertThrowsError(try runCommand("/usr/bin/false", []))
        XCTAssertThrowsError(try runCommand("/bin/sleep", ["5"], timeout: 0.1))
    }
    func testCommandPreservesArgumentsWithoutShellExpansion() throws {
        XCTAssertEqual(try runCommand("/bin/echo", ["$(id); secret"]), "$(id); secret\n")
    }
    func testProcessParserKeepsExecutableSpacesAndRejectsInvalidRows() throws {
        let records = try parseProcesses(" 42 1 Thu Oct 22 09:10:11 2026 /Applications/My App.app/Contents/MacOS/My App\n")
        XCTAssertEqual(records.count, 1)
        XCTAssertEqual(records[0].path, "/Applications/My App.app/Contents/MacOS/My App")
        XCTAssertEqual(records[0].facts["pid"], "42")
        XCTAssertThrowsError(try parseProcesses("not a process"))
    }
    func testAppInventoryRetainsReadableAppsAfterMalformedNeighbor() throws {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
        defer { try? FileManager.default.removeItem(at: root) }
        for name in ["Broken", "Valid"] {
            try FileManager.default.createDirectory(at: root.appendingPathComponent("\(name).app/Contents"), withIntermediateDirectories: true)
        }
        let data = try PropertyListSerialization.data(fromPropertyList: ["CFBundleIdentifier": "test.valid", "CFBundleShortVersionString": "1"], format: .xml, options: 0)
        try data.write(to: root.appendingPathComponent("Valid.app/Contents/Info.plist"))
        let section = try appInventory(roots: [root.path])
        XCTAssertEqual(section.records.map(\.title), ["Valid"])
        XCTAssertEqual(section.state, .partial)
        XCTAssertTrue(section.reason.contains("Broken.app"))
    }
    func testExportRoundTripPreservesUnknownCoverage() throws {
        let report = ScanReport(sections: [ScanSection(category: .network, state: .unsupported, reason: "Provider not installed", records: [])])
        let decoded = try JSONDecoder().decode(ScanReport.self, from: JSONEncoder().encode(report))
        XCTAssertEqual(decoded.sections[0].state, .unsupported)
        XCTAssertFalse(decoded.sections[0].reason.isEmpty)
    }
}
