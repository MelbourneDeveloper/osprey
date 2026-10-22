#if os(iOS)
import SwiftUI
import UniformTypeIdentifiers
import VigilCore
import VigilDesign

struct PhoneInvestigation: View {
    @State private var picking = false
    @State private var report: ScanReport?
    @State private var failure: String?
    var body: some View {
        NavigationStack {
            List {
                SwiftUI.Section {
                    Brand()
                    Text("iPhone investigation").font(.title2)
                    Text("iOS does not expose other apps’ processes, permissions or files. Import a Mac evidence report to review its observations. No live protection or synchronization is active.")
                    Button("Import evidence report") { picking = true }
                }
                if let report {
                    ForEach(report.sections) { section in reportSection(section) }
                }
                if let failure { Text(failure).foregroundStyle(.red) }
            }.navigationTitle("Vigil")
                .fileImporter(isPresented: $picking, allowedContentTypes: [.json]) { load($0) }
        }
    }
    private func reportSection(_ section: ScanSection) -> some View {
        SwiftUI.Section(section.category.title) {
            Text("Imported snapshot · \(section.observed.formatted())").font(.caption)
            Text(section.reason).font(.caption)
            ForEach(section.records) { record in
                NavigationLink(record.title) {
                    ScrollView {
                        VStack(alignment: .leading, spacing: 20) {
                            Text(record.detail)
                            EvidenceLine("Source", record.source)
                            EvidenceLine("Observed", record.observed.formatted())
                            EvidenceLine("Path", record.path)
                        }.padding()
                    }.navigationTitle(record.title)
                }
            }
        }
    }
    private func load(_ result: Result<URL, Error>) {
        do {
            let url = try result.get(), scoped = url.startAccessingSecurityScopedResource()
            defer { if scoped { url.stopAccessingSecurityScopedResource() } }
            report = try readReport(url); failure = nil
        } catch { failure = error.localizedDescription }
    }
}
#endif
