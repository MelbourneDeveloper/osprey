import Foundation

private struct StorageCommand: Decodable {
    let kind: String
    let path: String
}

private struct StorageEnvelope: Decodable {
    let model: String
    let ui: NativeNode
    let commands: [StorageCommand]
}

// Implements [MOBILE-REACTIVE-HOST]: keep Osprey's opaque browser, show the
// screen it returns, and do the host work it asks for. Nothing here decides
// what the screen holds.
@MainActor
final class StorageStore: ObservableObject {
    @Published private(set) var ui: NativeNode?
    @Published private(set) var failure: String?
    @Published var picking = false
    private var model = ""
    private var walk: Task<Void, Never>?
    private var scoped: URL?

    func start() {
        guard ui == nil, failure == nil else { return }
        let status = osprey_main()
        guard status == 0 else { return fail("Osprey initialization failed: \(status)") }
        accept(osprey_storage_start(canonical(NSHomeDirectory())))
    }

    func send(_ event: [String: Any]) {
        guard let data = try? JSONSerialization.data(withJSONObject: event),
            let json = String(data: data, encoding: .utf8)
        else { return fail("An event could not be encoded for Osprey") }
        accept(osprey_storage_dispatch(model, json))
    }

    /// A folder chosen in the system picker. On iPhone it stays readable only
    /// while its security scope is held, so the scope lasts until the next pick.
    func picked(_ result: Result<URL, Error>) {
        guard case let .success(url) = result else { return }
        scoped?.stopAccessingSecurityScopedResource()
        scoped = url.startAccessingSecurityScopedResource() ? url : nil
        send(["type": "picked", "value": canonical(url.path)])
    }

    private func accept(_ pointer: UnsafePointer<CChar>?) {
        guard let pointer,
            let envelope = try? JSONDecoder().decode(StorageEnvelope.self, from: Data(String(cString: pointer).utf8))
        else { return fail("Osprey returned a screen that could not be read") }
        model = envelope.model
        ui = envelope.ui
        envelope.commands.forEach(run)
    }

    private func run(_ command: StorageCommand) {
        switch command.kind {
        case "measure": measure(command.path)
        case "pick": picking = true
        default: fail("Osprey asked for unknown host work: \(command.kind)")
        }
    }

    /// Walks off the main thread; every result goes back to Osprey as an event.
    private func measure(_ root: String) {
        walk?.cancel()
        walk = Task.detached(priority: .userInitiated) { [weak self] in
            let tree = DiskTree.measure(root) { items, bytes in
                Task { @MainActor in self?.send(["type": "progress", "path": root, "items": items, "bytes": bytes]) }
            }
            await self?.finish(root, tree)
        }
    }

    private func finish(_ root: String, _ tree: DiskTree?) {
        guard !Task.isCancelled else { return }
        if let tree { DiskTree.current = tree }
        send(["type": tree == nil ? "failed" : "measured", "path": root])
    }

    private func fail(_ message: String) {
        failure = message
        walk?.cancel()
    }
}
