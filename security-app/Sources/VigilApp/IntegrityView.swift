#if os(macOS)
import SwiftUI
import VigilCore
import VigilDesign

struct IntegrityView: View {
    let storage: BaselineStore
    @State private var baselines: [FileFingerprint] = []
    @State private var observed: FileFingerprint?
    @State private var message = "Choose a local configuration file to calculate SHA-256. Recording the first observation does not prove the file is clean."
    @State private var busy = false
    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("File integrity").font(.largeTitle)
            Text(message).textSelection(.enabled)
            HStack {
                Button("Choose file to inspect") { choose() }.disabled(busy)
                if busy { ProgressView().controlSize(.small) }
                if let observed { Button("Record this observation as baseline") { record(observed) }.disabled(busy) }
            }
            if let observed { Text("SHA-256: \(observed.sha256)\n\(observed.path)").font(.system(.caption, design: .monospaced)).textSelection(.enabled) }
            List(baselines, id: \.path) { baseline in
                VStack(alignment: .leading) { Text(baseline.path); Text("Recorded \(baseline.observed.formatted()) · \(baseline.sha256)").font(.caption) }
            }
            Text("Local regular files up to 32 MB. Symlink targets and cloud-managed files are excluded. OS files are never repaired or modified.").font(.caption).foregroundStyle(.secondary)
        }.padding(28).task { reload() }
    }
    private func choose() {
        let panel = NSOpenPanel(); panel.canChooseFiles = true; panel.canChooseDirectories = false
        guard let window = NSApp.keyWindow else { message = "Activate the window to choose a file"; return }
        presentPanel(panel, for: window) { response in
            guard response == .OK, let file = panel.url else { return }
            inspect(file)
        }
    }
    private func inspect(_ file: URL) {
        busy = true; observed = nil
        let scoped = file.startAccessingSecurityScopedResource()
        Task {
            defer { busy = false; if scoped { file.stopAccessingSecurityScopedResource() } }
            do {
                let result = try await Task.detached { try fingerprint(file) }.value
                observed = result
                let prior = baselines.first { $0.path == result.path }
                message = prior.map { $0.sha256 == result.sha256 ? "Content matches the recorded baseline observation." : "Content changed since the baseline. Review the change; updates may be legitimate." } ?? "No prior baseline. This observation does not establish a trusted clean state."
            } catch { message = error.localizedDescription }
        }
    }
    private func record(_ fingerprint: FileFingerprint) {
        do { try storage.record(fingerprint); reload(); message = "Baseline recorded locally. This is a comparison point, not a certificate of safety." }
        catch { message = error.localizedDescription }
    }
    private func reload() {
        do { baselines = try storage.load() }
        catch { message = error.localizedDescription }
    }
}
#endif
