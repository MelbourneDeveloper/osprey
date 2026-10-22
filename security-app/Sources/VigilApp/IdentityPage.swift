import SwiftUI
import VigilDesign

struct IdentityPage: View {
    var snapshot = false
    @Environment(\.colorScheme) private var scheme
    @ReducedMotion private var reduced: Bool
    @State private var revealed = false
    var body: some View {
        Group {
            if snapshot { content } else { ScrollView { content } }
        }.navigationTitle("Identity & motion")
    }
    var content: some View {
            VStack(alignment: .leading, spacing: 28) {
                hero
                palette
                motion
            }.padding(32).frame(maxWidth: 1000, alignment: .leading).frame(maxWidth: .infinity)
    }
    private var hero: some View {
        Surface {
            VStack(alignment: .leading, spacing: 24) {
                Brand().scaleEffect(1.4, anchor: .leading).padding(.vertical, 16)
                Eyebrow("Vigil / Working identity")
                Text("Nothing hidden.\nNothing overstated.").font(.system(.largeTitle, design: .serif))
                Text("An aperture for attention. A quiet centre. Space to see clearly.").font(.title3).foregroundStyle(Ink.secondary.color(scheme))
                Text("Provisional name and generated mark; trademark and final icon production review remain open.").font(.caption).foregroundStyle(Ink.secondary.color(scheme))
            }
        }
    }
    private var palette: some View {
        LazyVGrid(columns: [GridItem(.adaptive(minimum: 140))], spacing: 16) {
            ForEach(Array(Ink.allCases.enumerated()), id: \.offset) { _, ink in
                VStack(alignment: .leading, spacing: 8) {
                    RoundedRectangle(cornerRadius: 12).fill(ink.color(scheme)).frame(height: 64)
                        .overlay(RoundedRectangle(cornerRadius: 12).stroke(Ink.line.color(scheme)))
                    Text(String(describing: ink).capitalized).font(.caption.weight(.semibold))
                    Text(String(format: "#%06X", ink.hex(dark: scheme == .dark))).font(.system(.caption, design: .monospaced)).foregroundStyle(Ink.secondary.color(scheme))
                }
            }
        }
    }
    private var motion: some View {
        Surface {
            VStack(alignment: .leading, spacing: 18) {
                Eyebrow("Motion with a purpose")
                Text("120 ms to acknowledge. 220 ms to settle.").font(.title3.weight(.semibold))
                Text("No pulsing alarms, rotating radars or perpetual animation. Reduced Motion removes custom transitions entirely.").foregroundStyle(Ink.secondary.color(scheme))
                Button(revealed ? "Close detail" : "Reveal detail") {
                    withAnimation(Motion.settle(reduced: reduced)) { revealed.toggle() }
                }.buttonStyle(VigilButtonStyle(prominent: true))
                if revealed {
                    EvidenceLine("Collection updated", "The source and collection time remain attached to the observation.")
                        .transition(.opacity)
                }
            }
        }
    }
}
