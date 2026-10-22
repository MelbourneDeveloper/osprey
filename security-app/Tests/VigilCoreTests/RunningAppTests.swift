import XCTest
@testable import VigilCore

final class RunningAppTests: XCTestCase {
    func testReusedPIDCannotInheritApplicationState() throws {
        let records = try parseProcesses("42 1 Thu Oct 22 09:10:11 2026 /Applications/Test.app/Contents/MacOS/Test\n")
        let fact = RunningAppFact(pid: 42, launched: Date(timeIntervalSince1970: 0), active: true, hidden: false, policy: "regular")
        XCTAssertNil(annotateRunningApps(records, applications: [fact])[0].facts["applicationState"])
    }
    func testBackgroundAndHiddenAreDifferentObservations() throws {
        let records = try parseProcesses("42 1 Thu Oct 22 09:10:11 2026 /Applications/Test.app/Contents/MacOS/Test\n")
        let launched = try XCTUnwrap(processStartDate("Thu Oct 22 09:10:11 2026"))
        let fact = RunningAppFact(pid: 42, launched: launched, active: false, hidden: false, policy: "regular")
        let result = annotateRunningApps(records, applications: [fact])[0]
        XCTAssertEqual(result.facts["applicationState"], "Background at observation time")
        XCTAssertEqual(result.facts["hiddenWindows"], "false")
        XCTAssertEqual(result.assessment, "information")
    }
}
