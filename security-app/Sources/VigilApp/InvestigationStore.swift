#if os(macOS)
import AppKit
import VigilCore
import VigilDesign

@MainActor
final class InvestigationStore: ObservableObject {
    @Published private(set) var report = ScanReport(sections: [])
    @Published private(set) var busy = false
    @Published private(set) var progress = "Ready to scan"
    @Published var error: String?
    @Published var selected: Evidence?
    @Published private(set) var changes: [EvidenceChange] = []
    @Published var findings: [Finding] = []
    @Published var history: [HistoryEntry] = []
    private var previous = ScanReport(sections: [])
    private let storage = ReportStore(directory: FileManager.default.homeDirectoryForCurrentUser.appendingPathComponent("Library/Application Support/Vigil"))
    private var worker: Task<ScanSection, Error>?
    private var generation = UUID()
    private let engine = Bundle.main.executableURL?.deletingLastPathComponent().appendingPathComponent("vigil-engine")

    var historyStorage: HistoryStore { HistoryStore(directory: storage.directory.appendingPathComponent("history")) }
    var findingStorage: FindingStore { FindingStore(directory: storage.directory) }
    var baselineStorage: BaselineStore { BaselineStore(directory: storage.directory) }
    init() {
        do { if let saved = try storage.load() { report = saved }; findings = try findingStorage.load(); history = try historyStorage.entries() }
        catch { self.error = "Saved evidence could not be loaded: \(error.localizedDescription)" }
    }
    func scan() {
        previous = report; selected = nil; changes = []
        cancel(); report = ScanReport(sections: []); busy = true
        let run = generation
        Task {
            for area in Area.allCases {
                guard generation == run else { return }
                progress = "Collecting \(area.title.lowercased())…"
                await refresh(area, run: run)
            }
            guard generation == run else { return }
            busy = false; progress = "Snapshot completed at \(Date().formatted(date: .omitted, time: .standard))"
            self.save(); self.archive()
        }
    }
    func cancel() {
        generation = UUID(); worker?.cancel(); worker = nil; busy = false
        progress = "Stopped; collected snapshots are retained"
    }
    private func refresh(_ area: Area, folder: URL? = nil, run: UUID) async {
        let engine = engine
        let task = Task.detached(priority: .userInitiated) {
            let section = collect(area, folder: folder)
            guard let engine else { throw CollectionError.invalid("The Osprey engine is missing; build with ./build.sh") }
            return assessedSnapshot(section, engine: engine)
        }
        worker = task
        do {
            var section = try await task.value
            guard generation == run else { return }
            if area == .processes { section.records = annotateRunningApps(section.records, applications: runningApplications()) }
            accept(retainAdvisories(section, previous: previous.sections.first { $0.category == area }))
            if area == .startup || area == .packages,
               let prior = previous.sections.first(where: { $0.category == area }) {
                changes += evidenceChanges(previous: prior.records, current: section.records)
            }
        } catch { if generation == run { self.error = error.localizedDescription } }
    }
    private func accept(_ section: ScanSection) {
        findings = reconcileFindings(findings, section: section, now: Date())
        report.sections.removeAll { $0.category == section.category }; report.sections.append(section)
    }
    func chooseFolder() {
        let panel = NSOpenPanel()
        panel.canChooseDirectories = true; panel.canChooseFiles = false; panel.allowsMultipleSelection = false
        guard let window = NSApp.keyWindow else { error = "No active window for the file dialog"; return }
        presentPanel(panel, for: window) { [weak self] response in
            guard response == .OK, let url = panel.url, let self else { return }
            let scoped = url.startAccessingSecurityScopedResource()
            self.cancel(); self.busy = true; self.selected = nil
            let run = self.generation
            Task {
                defer { if scoped { url.stopAccessingSecurityScopedResource() } }
                await self.refresh(.files, folder: url, run: run)
                if self.generation == run { self.busy = false; self.progress = "Folder inspection completed"; self.save() }
            }
        }
    }
    func inspect(_ record: Evidence) {
        guard !busy else { return }
        busy = true; progress = "Verifying \(record.title)…"
        let run = generation, engine = engine
        worker = Task.detached {
            let inspected = try inspectSignature(record)
            guard let engine else { throw CollectionError.invalid("Osprey engine is missing") }
            return try assess(ScanSection(category: .apps, state: .partial, reason: "Selected app", records: [inspected]), engine: engine)
        }
        Task { await finishInspection(run) }
    }
    private func finishInspection(_ run: UUID) async {
        do {
            guard let section = try await worker?.value, generation == run, let record = section.records.first else { return }
            if let index = report.sections.firstIndex(where: { $0.category == .apps }),
               let row = report.sections[index].records.firstIndex(where: { $0.id == record.id }) { report.sections[index].records[row] = record }
            findings = reconcileFindings(findings, section: section, now: Date())
            selected = record
        } catch { if generation == run { self.error = error.localizedDescription } }
        if generation == run { busy = false; progress = "Inspection completed"; save() }
    }
    private func runningApplications() -> [RunningAppFact] {
        NSWorkspace.shared.runningApplications.map {
            RunningAppFact(pid: $0.processIdentifier, launched: $0.launchDate, active: $0.isActive,
                           hidden: $0.isHidden, policy: String(describing: $0.activationPolicy))
        }
    }
    private func save() {
        do { try storage.save(report); try findingStorage.save(findings) }
        catch { self.error = "Evidence could not be saved: \(error.localizedDescription)" }
    }
    func archive() {
        do { try historyStorage.save(report); history = try historyStorage.entries() }
        catch { self.error = error.localizedDescription }
    }
    func updateFinding(_ finding: Finding, status: FindingStatus, reason: String, days: Int = 7) {
        do {
            let updated = try transition(finding, to: status, reason: reason, until: Date().addingTimeInterval(Double(days) * 86400), now: Date())
            findings = findings.map { $0.id == finding.id ? updated : $0 }; try findingStorage.save(findings)
        } catch { self.error = error.localizedDescription }
    }
    func observeNetwork() { runSupplement("Observing open IP sockets…") { try networkSnapshot() } }
    func downloadAdvisories() { runSupplement("Downloading public CISA catalog; no device data is sent…") { try await fetchKEV() } }
    private func runSupplement(_ label: String, operation: @escaping @Sendable () async throws -> ScanSection) {
        guard !busy else { return }
        busy = true; progress = label
        let run = generation
        worker = Task.detached { try await operation() }
        Task {
            do { if let section = try await worker?.value, generation == run { replace(section) } }
            catch { if generation == run { self.error = error.localizedDescription } }
            if generation == run { busy = false; progress = "Collection finished" }
        }
    }
    func importAdvisories() {
        let panel = NSOpenPanel(); panel.allowedContentTypes = [.json]
        guard let window = NSApp.keyWindow else { error = "Activate the window to choose a catalog"; return }
        presentPanel(panel, for: window) { [weak self] response in
            guard response == .OK, let file = panel.url, let self else { return }
            do { self.replace(try decodeKEV(boundedData(file))) }
            catch { self.error = error.localizedDescription }
        }
    }
    private func replace(_ section: ScanSection) {
        accept(section); save()
    }
    func exportReport() {
        let panel = NSSavePanel(); panel.nameFieldStringValue = "vigil-evidence.json"
        guard let window = NSApp.keyWindow else { error = "No active window for the file dialog"; return }
        presentPanel(panel, for: window) { [weak self] response in
            guard response == .OK, let url = panel.url, let self else { return }
            do {
                let encoder = JSONEncoder(); encoder.outputFormatting = [.prettyPrinted, .sortedKeys]
                try encoder.encode(self.report).write(to: url, options: .atomic)
            } catch { self.error = error.localizedDescription }
        }
    }
}
#endif
