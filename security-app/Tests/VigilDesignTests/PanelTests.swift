#if os(macOS)
import AppKit
import XCTest
@testable import VigilDesign

final class PanelTests: XCTestCase {
    @MainActor
    func testFileDialogIsAttachedToItsOwningWindow() async throws {
        _ = NSApplication.shared
        let window = NSWindow(contentRect: NSRect(x: 0, y: 0, width: 800, height: 600), styleMask: [.titled, .closable], backing: .buffered, defer: false)
        window.isReleasedWhenClosed = false
        let panel = NSOpenPanel()
        defer { panel.cancel(nil); window.close() }
        window.makeKeyAndOrderFront(nil)
        presentPanel(panel, for: window) { _ in }
        try await Task.sleep(nanoseconds: 300_000_000)
        XCTAssertTrue(panel.sheetParent === window, "File dialogs must remain attached to the owning app window")
    }
}
#endif
