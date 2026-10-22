#if os(macOS)
import AppKit
import XCTest

final class LaunchTests: XCTestCase {
    private var packagedApp: URL {
        URL(fileURLWithPath: #filePath).deletingLastPathComponent().deletingLastPathComponent().deletingLastPathComponent().appendingPathComponent(".build/Vigil.app")
    }
    func testPackagedApplicationHasValidResourceSeal() throws {
        let process = Process()
        process.executableURL = URL(fileURLWithPath: "/usr/bin/codesign")
        process.arguments = ["--verify", "--deep", "--strict", packagedApp.path]
        try process.run(); process.waitUntilExit()
        XCTAssertEqual(process.terminationStatus, 0, "The packaged app must be signed after its resources and Info.plist are assembled")
    }
    func testPackagedApplicationLaunches() throws {
        try launch(packagedApp.appendingPathComponent("Contents/MacOS/Vigil"))
    }
    func testCommandLineLaunchBecomesRegularApplication() throws {
        let products = Bundle(for: Self.self).bundleURL.deletingLastPathComponent()
        try launch(products.appendingPathComponent("Vigil"))
    }
    private func launch(_ executable: URL) throws {
        let process = Process()
        process.executableURL = executable
        try process.run()
        defer { process.terminate(); process.waitUntilExit() }
        let deadline = Date().addingTimeInterval(8)
        while Date() < deadline {
            if NSRunningApplication(processIdentifier: process.processIdentifier)?.activationPolicy == .regular { return }
            Thread.sleep(forTimeInterval: 0.1)
        }
        XCTFail("Vigil must launch as a regular macOS application, with a Dock presence and foreground window")
    }
}
#endif
