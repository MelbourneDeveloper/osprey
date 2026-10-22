#if os(macOS)
import SwiftUI
import VigilCore
import VigilDesign

struct SecurityWorkspace: View {
    @StateObject private var store = InvestigationStore()
    var body: some View {
        TabView {
            InvestigationView(store: store).tabItem { Label("Investigation", systemImage: "magnifyingglass") }
            FindingsView(store: store).tabItem { Label("Review", systemImage: "checklist") }
            HistoryView(store: store).tabItem { Label("History", systemImage: "clock") }
            IntegrityView(storage: store.baselineStorage).tabItem { Label("File integrity", systemImage: "doc.badge.gearshape") }
            AgentView().tabItem { Label("Optional assistant", systemImage: "bubble.left.and.bubble.right") }
        }.task { store.scan() }
    }
}
struct FindingsView: View {
    @ObservedObject var store: InvestigationStore
    @State private var selection: String?
    @State private var reason = ""
    private var finding: Finding? { store.findings.first { $0.id == selection } }
    var body: some View {
        NavigationSplitView {
            List(store.findings, selection: $selection) { finding in
                VStack(alignment: .leading) {
                    Text(finding.evidence.title).font(.headline)
                    Text("\(finding.ruleID) · \(finding.status.rawValue)").font(.caption)
                }.tag(finding.id)
            }.navigationTitle("Review queue")
        } detail: {
            if let finding { findingDetail(finding) }
            else { ContentUnavailableView("Review observations", systemImage: "checklist", description: Text("Review items appear after scanning or signature inspection. An empty queue does not establish that the device is safe.")) }
        }.onChange(of: selection) { _, _ in reason = "" }
        .alert("Review could not be saved", isPresented: Binding(get: { store.error != nil }, set: { if !$0 { store.error = nil } })) {
            Button("OK") { store.error = nil }
        } message: { Text(store.error ?? "") }
    }
    private func findingDetail(_ finding: Finding) -> some View {
        Form {
            SwiftUI.Section("Evidence") {
                Text(finding.evidence.title).font(.title)
                Text(finding.evidence.detail).textSelection(.enabled)
                LabeledContent("Rule", value: finding.ruleID)
                LabeledContent("First seen", value: finding.firstSeen.formatted())
                LabeledContent("Last observed", value: finding.lastSeen.formatted())
                LabeledContent("Status", value: finding.status.rawValue)
                Text("Review context, not a malware verdict. Resolution requires a relevant follow-up observation; missing data cannot resolve a finding.").font(.caption)
            }
            SwiftUI.Section("Investigation") { reviewControls(finding) }
        }.formStyle(.grouped)
    }
    private func reviewControls(_ finding: Finding) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            TextField("Investigation note / suppression reason", text: $reason)
            HStack {
                Button("Acknowledge") { store.updateFinding(finding, status: .acknowledged, reason: reason) }
                Button("Investigate") { store.updateFinding(finding, status: .investigating, reason: reason) }
                Button("Suppress for 7 days") { store.updateFinding(finding, status: .suppressed, reason: reason) }.disabled(reason.trimmingCharacters(in: .whitespaces).isEmpty)
                Button("Reopen") { store.updateFinding(finding, status: .new, reason: reason) }
            }.disabled(store.busy)
            Text(finding.reason).foregroundStyle(.secondary)
            if let expiry = finding.expires { Text("Suppression expires \(expiry.formatted()). Changes in evidence reopen review.").font(.caption) }
        }
    }
}
struct HistoryView: View {
    @ObservedObject var store: InvestigationStore
    @State private var selected: String?
    @State private var report: ScanReport?
    @State private var failure: String?
    var body: some View {
        NavigationSplitView {
            List(store.history, selection: $selected) { item in Text(item.created.formatted()).tag(item.id) }
                .navigationTitle("Snapshot history")
        } detail: {
            if let report {
                List(report.sections) { section in
                    VStack(alignment: .leading) {
                        Text("\(section.category.title): \(section.records.count) entries · \(section.state.rawValue)").font(.headline)
                        Text(section.reason).font(.caption)
                        ForEach(section.records.prefix(100)) { Text("\($0.title) — \($0.detail)").font(.caption).textSelection(.enabled) }
                        if section.records.count > 100 { Text("Showing the first 100 entries; full evidence is retained in this snapshot.").font(.caption) }
                    }
                }
            } else { ContentUnavailableView("Local snapshot history", systemImage: "clock", description: Text(failure ?? "Select a snapshot. Up to 30 snapshots are kept for 30 days. Historical observations do not describe current state.")) }
        }.onChange(of: selected) { _, id in load(id) }
    }
    private func load(_ id: String?) {
        report = nil; failure = nil
        guard let entry = store.history.first(where: { $0.id == id }) else { return }
        do { report = try readReport(entry.file) }
        catch { failure = error.localizedDescription }
    }
}
#endif
