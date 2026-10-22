import SwiftUI
import VigilDesign

struct Overview: View {
    var snapshot = false
    @Environment(\.colorScheme) private var scheme
    @ReducedMotion private var reduced: Bool
    @State private var selected: Finding? = findings.first
    @State private var query = ""
    private var filtered: [Finding] {
        findings.filter { query.isEmpty || "\($0.title) \($0.app)".localizedCaseInsensitiveContains(query) }
    }
    var body: some View {
        Group {
            if snapshot { content } else { ScrollView { content } }
        }.navigationTitle("Overview")
            .animation(Motion.settle(reduced: reduced), value: selected?.id)
    }
    var content: some View {
            VStack(alignment: .leading, spacing: 28) {
                heading
                coverage
                ViewThatFits(in: .horizontal) {
                    HStack(alignment: .top, spacing: 20) {
                        findingList.frame(minWidth: 370)
                        inspector.frame(width: 290)
                    }
                    VStack(spacing: 20) { findingList; inspector }
                }
                footer
            }.padding(32).frame(maxWidth: 1240, alignment: .leading).frame(maxWidth: .infinity)
    }
    private var heading: some View {
        VStack(alignment: .leading, spacing: 10) {
            Eyebrow("Your device. In plain sight.")
            Text("Clarity is protection.").font(.system(.largeTitle, design: .serif).weight(.medium))
            Text("Review what changed. Understand what matters.").font(.title3).foregroundStyle(Ink.secondary.color(scheme))
        }
    }
    private var coverage: some View {
        Surface {
            VStack(alignment: .leading, spacing: 18) {
                HStack(alignment: .top) {
                    VStack(alignment: .leading, spacing: 7) {
                        Eyebrow("Collection coverage")
                        Text("4 of 6 areas available").font(.title2.weight(.semibold))
                    }
                    Spacer()
                    StatusBadge(Coverage.partial)
                }
                CoverageStrip([.available, .available, .available, .available, .denied, .unsupported])
                Text("File access is denied. Network inspection is unavailable in this fixture.")
                    .font(.callout).foregroundStyle(Ink.secondary.color(scheme))
            }
        }
    }
    private var findingList: some View {
        VStack(alignment: .leading, spacing: 16) {
            HStack { Text("Worth your attention").font(.title3.weight(.semibold)); Spacer(); Text("3 findings").foregroundStyle(Ink.secondary.color(scheme)) }
            TextField("Filter findings", text: $query).textFieldStyle(.roundedBorder)
                .accessibilityLabel("Filter findings by title or application")
            ForEach(filtered) { finding in findingButton(finding) }
            if filtered.isEmpty { ContentUnavailableView.search(text: query) }
        }
    }
    private func findingButton(_ finding: Finding) -> some View {
        Button { selected = finding } label: {
            FindingRow(finding: finding, selected: selected?.id == finding.id)
        }.buttonStyle(.plain).accessibilityAddTraits(selected?.id == finding.id ? .isSelected : [])
    }
    private var inspector: some View {
        Surface {
            if let finding = selected {
                FindingInspector(finding: finding)
            } else {
                Text("Select a finding to inspect its evidence.").foregroundStyle(Ink.secondary.color(scheme))
            }
        }
    }
    private var footer: some View {
        HStack(spacing: 8) {
            Image(systemName: "internaldrive")
            Text("Local by default. AI is optional.")
            Spacer()
            Text("DEMONSTRATION DATA").font(.caption2.weight(.semibold)).tracking(1)
        }.font(.caption).foregroundStyle(Ink.secondary.color(scheme))
    }
}

struct FindingRow: View {
    @Environment(\.colorScheme) private var scheme
    let finding: Finding
    let selected: Bool
    var body: some View {
        HStack(alignment: .top, spacing: 14) {
            Image(systemName: finding.symbol).font(.title2).frame(width: 30).padding(.top, 3)
            VStack(alignment: .leading, spacing: 9) {
                Text(finding.app).font(.caption).foregroundStyle(Ink.secondary.color(scheme))
                Text(finding.title).font(.headline).fixedSize(horizontal: false, vertical: true)
                StatusBadge(finding.severity)
            }
            Spacer(minLength: 0)
            Image(systemName: "chevron.right").font(.caption).foregroundStyle(Ink.secondary.color(scheme))
        }.padding(18).frame(maxWidth: .infinity, alignment: .leading)
            .background(selected ? Ink.inset.color(scheme) : Ink.surface.color(scheme), in: RoundedRectangle(cornerRadius: 14))
            .overlay(RoundedRectangle(cornerRadius: 14).strokeBorder(selected ? Ink.accent.color(scheme) : Ink.line.color(scheme), lineWidth: selected ? 1.5 : 0.5))
            .contentShape(RoundedRectangle(cornerRadius: 14))
    }
}

struct FindingInspector: View {
    @Environment(\.colorScheme) private var scheme
    let finding: Finding
    @State private var showingExplanation = false
    var body: some View {
        VStack(alignment: .leading, spacing: 20) {
            Eyebrow("Evidence inspector")
            Image(systemName: finding.symbol).font(.system(size: 32, weight: .light)).accessibilityHidden(true)
            Text(finding.app).font(.title2.weight(.semibold))
            Text(finding.detail).font(.callout).lineSpacing(4)
            Divider()
            EvidenceLine("Source", finding.evidence)
            EvidenceLine("Confidence", "Illustrative · not an assessment")
            EvidenceLine("Next step", "Verify the app’s purpose and expected access.")
            Button("How to interpret this", systemImage: "info.circle") { showingExplanation = true }
                .buttonStyle(VigilButtonStyle()).sheet(isPresented: $showingExplanation) { explanation }
        }
    }
    private var explanation: some View {
        VStack(alignment: .leading, spacing: 20) {
            Text("Evidence before action").font(.title2.weight(.semibold))
            Text("These are synthetic observations, not results from your Mac. In the product, every finding must identify its source, collection time, uncertainty and a reversible next step. A permission or background process alone does not prove spyware.")
            Button("Done") { showingExplanation = false }.buttonStyle(VigilButtonStyle(prominent: true))
        }.padding(32).frame(idealWidth: 460)
    }
}
