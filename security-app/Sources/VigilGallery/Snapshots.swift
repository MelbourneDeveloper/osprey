#if os(macOS)
import AppKit
import SwiftUI
import VigilDesign

// Offline design review; no device observations or screen capture permissions.
@MainActor func renderSnapshotsIfRequested() {
    guard let index = CommandLine.arguments.firstIndex(of: "--render"),
          CommandLine.arguments.indices.contains(index + 1) else { return }
    let directory = URL(fileURLWithPath: CommandLine.arguments[index + 1])
    do {
        try FileManager.default.createDirectory(at: directory, withIntermediateDirectories: true)
        try render(Overview(snapshot: true), name: "overview-light", directory: directory, width: 1100, dark: false)
        try render(Overview(snapshot: true), name: "overview-dark", directory: directory, width: 1100, dark: true)
        try render(Overview(snapshot: true), name: "overview-compact", directory: directory, width: 390, dark: false)
        try render(ComponentCatalog(snapshot: true), name: "components", directory: directory, width: 1000, dark: false)
        try render(IdentityPage(snapshot: true), name: "identity", directory: directory, width: 1000, dark: true)
        exit(EXIT_SUCCESS)
    } catch { fputs("Snapshot failed: \(error)\n", stderr); exit(EXIT_FAILURE) }
}

@MainActor private func render<V: View>(_ view: V, name: String, directory: URL, width: CGFloat, dark: Bool) throws {
    let scheme: ColorScheme = dark ? .dark : .light
    let content = view.frame(width: width, height: width < 500 ? 1700 : 1000, alignment: .top)
        .background(Ink.canvas.color(scheme)).foregroundStyle(Ink.primary.color(scheme))
        .environment(\.colorScheme, scheme).environment(\.vigilReduceMotion, true)
    let host = NSHostingView(rootView: content)
    host.frame = NSRect(x: 0, y: 0, width: width, height: width < 500 ? 1700 : 1000)
    host.layoutSubtreeIfNeeded()
    try saveSnapshot(host, to: directory.appendingPathComponent("\(name).png"))
}

@MainActor private func saveSnapshot(_ host: NSView, to url: URL) throws {
    guard let bitmap = host.bitmapImageRepForCachingDisplay(in: host.bounds) else {
        throw CocoaError(.fileWriteUnknown)
    }
    host.cacheDisplay(in: host.bounds, to: bitmap)
    guard let data = bitmap.representation(using: .png, properties: [:]) else {
        throw CocoaError(.fileWriteUnknown)
    }
    try data.write(to: url)
}
#endif
