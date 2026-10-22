#if os(macOS)
import SwiftUI
import AppKit
import VigilCore
import VigilDesign

struct InvestigationView: View {
    @ObservedObject var store: InvestigationStore
    @State private var area: Area? = .apps
    @State private var query = ""
    @State private var reviewOnly = false
    @Environment(\.colorScheme) private var scheme
    private var section: ScanSection? { store.report.sections.first { $0.category == area } }
    private var records: [Evidence] {
        (section?.records ?? []).filter {
            (!reviewOnly || $0.assessment != "information") && (query.isEmpty || "\($0.title) \($0.path) \($0.detail)".localizedCaseInsensitiveContains(query))
        }
    }
    var body: some View {
        NavigationSplitView { sidebar } content: { inventory } detail: { detail }
            .tint(Ink.accent.color(scheme))
            .toolbar { actions }
            .alert("Collection needs attention", isPresented: Binding(get: { store.error != nil }, set: { if !$0 { store.error = nil } })) {
                Button("OK") { store.error = nil }
            } message: { Text(store.error ?? "") }

            .onChange(of: area) { _, _ in store.selected = nil; query = "" }
    }
    private var sidebar: some View {
        VStack(alignment: .leading) {
            Brand().padding(20)
            List(Area.allCases, selection: $area) { item in
                Label(item.title, systemImage: symbol(item)).tag(item)
            }.listStyle(.sidebar)
            Text("Local investigation\nRead-only snapshots · AI off").font(.caption).foregroundStyle(.secondary).padding(20)
        }.navigationSplitViewColumnWidth(min: 180, ideal: 210)
    }
    private var inventory: some View {
        VStack(alignment: .leading, spacing: 12) {
            coverage.padding(.horizontal)
            Toggle("Only items to review", isOn: $reviewOnly).padding(.horizontal)
            List(records, selection: $store.selected) { record in
                VStack(alignment: .leading, spacing: 5) {
                    Text(record.title).font(.headline)
                    Text(record.detail).font(.caption).foregroundStyle(.secondary).lineLimit(2)
                    if record.assessment != "information" { Text(record.assessment.replacingOccurrences(of: "-", with: " ")).font(.caption).foregroundStyle(Ink.caution.color(scheme)) }
                }.padding(.vertical, 5).tag(record)
            }.overlay { if records.isEmpty { ContentUnavailableView("No entries", systemImage: "tray", description: Text(section?.reason ?? "Collection has not completed.")) } }
        }.navigationTitle(area?.title ?? "Investigation")
            .searchable(text: $query, prompt: "Filter local evidence")
            .navigationSplitViewColumnWidth(min: 300, ideal: 430)
    }
    private var coverage: some View {
        VStack(alignment: .leading, spacing: 8) {
            HStack {
                if store.busy { ProgressView().controlSize(.small) }
                Text(store.progress).font(.caption).foregroundStyle(.secondary)
                if !store.changes.isEmpty { Text("\(store.changes.count) startup/package snapshot differences").font(.caption) }
            }
            if let section {
                HStack { StatusBadge(Coverage(rawValue: section.state.rawValue) ?? .failed); Spacer(); Text("\(records.count) entries").font(.caption) }
                DisclosureGroup("Scope and collection limits") {
                    Text(section.reason).font(.caption).foregroundStyle(.secondary).textSelection(.enabled)
                }.font(.caption)
            }
        }.padding(.top, 14)
    }
    @ViewBuilder private var detail: some View {
        if let record = store.selected {
            ScrollView {
                VStack(alignment: .leading, spacing: 20) {
                    Eyebrow("Observed evidence")
                    Text(record.title).font(.largeTitle)
                    Text(record.detail).textSelection(.enabled)
                    EvidenceLine("Source", record.source)
                    EvidenceLine("Observed", record.observed.formatted())
                    EvidenceLine("Path", record.path.isEmpty ? "Not applicable" : record.path)
                    ForEach(record.facts.keys.sorted(), id: \.self) { key in EvidenceLine(key, record.facts[key] ?? "Unknown") }
                    controls(record)
                    Text("This evidence is not a spyware verdict. Snapshots can become stale; refresh before taking action.").font(.caption).foregroundStyle(.secondary)
                }.padding(28).frame(maxWidth: .infinity, alignment: .leading)
            }
        } else { ContentUnavailableView("Inspect the evidence", systemImage: "doc.text.magnifyingglass", description: Text("Select an entry to see its source, timestamp and collection limits.")) }
    }
    private func controls(_ record: Evidence) -> some View {
        HStack {
            if area == .permissions {
                Button("Open System Settings") { NSWorkspace.shared.openApplication(at: URL(fileURLWithPath: "/System/Applications/System Settings.app"), configuration: .init()) }
            }
            if area == .apps { Button("Verify signature & sandbox") { store.inspect(record) }.disabled(store.busy) }
            if !record.path.isEmpty { Button("Reveal in Finder") { NSWorkspace.shared.activateFileViewerSelecting([URL(fileURLWithPath: record.path)]) } }
        }.buttonStyle(VigilButtonStyle())
    }
    @ToolbarContentBuilder private var actions: some ToolbarContent {
        ToolbarItemGroup {
            Button("Refresh", systemImage: "arrow.clockwise") { store.scan() }.keyboardShortcut("r").disabled(store.busy)
            if store.busy { Button("Stop", systemImage: "stop.circle") { store.cancel() } }
            if area == .files { Button("Choose folder", systemImage: "folder.badge.plus") { store.chooseFolder() }.disabled(store.busy) }
            if area == .network { Button("Observe open sockets", systemImage: "network") { store.observeNetwork() }.disabled(store.busy) }
            if area == .advisories {
                Button("Download CISA catalog", systemImage: "arrow.down.circle") { store.downloadAdvisories() }.disabled(store.busy)
                Button("Import KEV catalog", systemImage: "square.and.arrow.down") { store.importAdvisories() }.disabled(store.busy) }
            Button("Export evidence", systemImage: "square.and.arrow.up") { store.exportReport() }.disabled(store.report.sections.isEmpty || store.busy)
        }
    }
    private func symbol(_ area: Area) -> String {
        switch area {
        case .apps: return "app.badge"
        case .processes: return "waveform.path.ecg"
        case .startup: return "power"
        case .files: return "folder"
        case .protection: return "lock.shield"
        case .packages: return "shippingbox"
        case .network: return "network"
        case .permissions: return "hand.raised"
        case .advisories: return "exclamationmark.bubble"
        }
    }
}
#endif
