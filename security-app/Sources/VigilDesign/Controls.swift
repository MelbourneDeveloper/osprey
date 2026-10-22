import SwiftUI

// Implements [DESIGN-MOTION]. No looping animation, scale, parallax or scan theatre.
public struct VigilButtonStyle: ButtonStyle {
    @Environment(\.colorScheme) private var scheme
    @ReducedMotion private var reduced: Bool
    @Environment(\.isEnabled) private var enabled
    let prominent: Bool
    public init(prominent: Bool = false) { self.prominent = prominent }
    public func makeBody(configuration: Configuration) -> some View {
        configuration.label.font(.callout.weight(.semibold))
            .padding(.horizontal, 16).frame(minHeight: 44)
            .foregroundStyle(prominent ? Ink.canvas.color(scheme) : Ink.primary.color(scheme))
            .background(prominent ? Ink.accent.color(scheme) : Ink.inset.color(scheme), in: RoundedRectangle(cornerRadius: 10))
            .opacity(enabled ? (configuration.isPressed ? 0.72 : 1) : 0.45)
            .animation(Motion.feedback(reduced: reduced), value: configuration.isPressed)
    }
}

public struct Brand: View {
    public init() {}
    public var body: some View {
        HStack(spacing: 10) {
            Image("VigilIcon", bundle: .module).resizable().interpolation(.high)
                .frame(width: 36, height: 36).clipShape(RoundedRectangle(cornerRadius: 9))
                .accessibilityHidden(true)
            Text("vigil").font(.system(size: 27, weight: .semibold, design: .rounded)).tracking(-1)
        }.accessibilityLabel("Vigil")
    }
}

public struct CoverageStrip: View {
    @Environment(\.colorScheme) private var scheme
    let states: [Coverage]
    public init(_ states: [Coverage]) { self.states = states }
    public var body: some View {
        HStack(spacing: 5) {
            ForEach(Array(states.enumerated()), id: \.offset) { _, state in
                RoundedRectangle(cornerRadius: 3).fill(state.ink.color(scheme))
                    .frame(height: 5)
            }
        }.accessibilityLabel("\(states.filter { $0 == .available }.count) of \(states.count) collection areas available")
    }
}
