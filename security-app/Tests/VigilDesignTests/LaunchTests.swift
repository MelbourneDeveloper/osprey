#if os(macOS)
import AppKit
import XCTest

final class LaunchTests: XCTestCase {
    func testCommandLineLaunchBecomesRegularApplication() throws {
        let root = URL(fileURLWithPath: #filePath).deletingLastPathComponent()
            .deletingLastPathComponent().deletingLastPathComponent()
        let process = Process()
        process.executableURL = root.appendingPathComponent(".build/out/Products/Debug/VigilGallery")
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
