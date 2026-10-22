#if os(macOS)
import AppKit

@MainActor
public func presentPanel(_ panel: NSSavePanel, for window: NSWindow, completion: @escaping (NSApplication.ModalResponse) -> Void) {
    panel.beginSheetModal(for: window, completionHandler: completion)
}
#endif
