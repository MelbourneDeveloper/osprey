import SwiftUI
import VigilDesign

struct ComponentCatalog: View {
    var snapshot = false
    @Environment(\.colorScheme) private var scheme
    @State private var selection: Coverage = .partial
    @State private var notes = ""
    @State private var alerts = true
    @State private var saved = false
    var body: some View {
        Group {
            if snapshot { content } else { ScrollView { content } }
        }.navigationTitle("Design library")
    }
    var content: some View {
            VStack(alignment: .leading, spacing: 24) {
                Text("A language of clarity.").font(.largeTitle.weight(.medium))
                Text("Native controls. Explicit states. One visual grammar.").foregroundStyle(Ink.secondary.color(scheme))
                Surface { states }
                Surface { interaction }
                Surface { typography }
            }.padding(32).frame(maxWidth: 1000, alignment: .leading).frame(maxWidth: .infinity)
    }
    private var states: some View {
        VStack(alignment: .leading, spacing: 18) {
            Eyebrow("01 / Status is never just a colour")
            LazyVGrid(columns: [GridItem(.adaptive(minimum: 150), alignment: .leading)], alignment: .leading, spacing: 16) {
                ForEach(Coverage.allCases) { StatusBadge($0) }
                ForEach(Severity.allCases) { StatusBadge($0) }
            }
            Divider()
            Picker("Preview coverage", selection: $selection) {
                ForEach(Coverage.allCases) { Text($0.title).tag($0) }
            }.frame(maxWidth: 360)
            HStack { Text("Selected collection state"); Spacer(); StatusBadge(selection) }
        }
    }
    private var interaction: some View {
        VStack(alignment: .leading, spacing: 16) {
            Eyebrow("02 / Controls & feedback")
            Toggle("Notify me when collection becomes unavailable", isOn: $alerts)
            TextField("Investigation note", text: $notes).textFieldStyle(.roundedBorder)
                .onChange(of: notes) { _, _ in saved = false }
            HStack {
                Button(saved ? "Saved in this session" : "Save note") { saved = true }
                    .buttonStyle(VigilButtonStyle(prominent: true)).disabled(notes.trimmingCharacters(in: .whitespaces).isEmpty)
                Button("Clear") { notes = "" }.buttonStyle(VigilButtonStyle())
            }
            Text("Gallery settings are session-only and never change device protection.").font(.caption).foregroundStyle(Ink.secondary.color(scheme))
        }
    }
    private var typography: some View {
        VStack(alignment: .leading, spacing: 14) {
            Eyebrow("03 / System typography")
            Text("A quieter kind of confidence.").font(.system(.largeTitle, design: .serif))
            Text("The important details, beautifully legible.").font(.title3.weight(.semibold))
            Text("System fonts respect the platform, scale with accessibility settings, and keep evidence readable.").font(.body)
            Text("SHA-256  8a2c…f019").font(.system(.callout, design: .monospaced)).textSelection(.enabled)
        }
    }
}
