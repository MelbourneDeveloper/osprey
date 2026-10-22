import SwiftUI

// Implements [DESIGN-FOUNDATIONS]. Native adapter; domain policy stays in Osprey.
public enum Ink: CaseIterable {
    case canvas, surface, inset, primary, secondary, accent, good, caution, danger, line

    public func hex(dark: Bool) -> UInt32 {
        switch self {
        case .canvas: return dark ? 0x101715 : 0xF5F5F0
        case .surface: return dark ? 0x19221F : 0xFFFFFF
        case .inset: return dark ? 0x202C27 : 0xEBEEE7
        case .primary: return dark ? 0xEDF2E9 : 0x1A2922
        case .secondary: return dark ? 0xAFBDB2 : 0x58665D
        case .accent: return dark ? 0xB7D9A4 : 0x365D39
        case .good: return dark ? 0xA8D4B2 : 0x306342
        case .caution: return dark ? 0xEAC28E : 0x855211
        case .danger: return dark ? 0xF4AAA0 : 0xA13D34
        case .line: return dark ? 0x3E5045 : 0xC6CFC4
        }
    }

    public func color(_ scheme: ColorScheme) -> Color {
        let value = hex(dark: scheme == .dark)
        return Color(red: Double(value >> 16) / 255,
                     green: Double((value >> 8) & 255) / 255,
                     blue: Double(value & 255) / 255)
    }
}

public enum Space {
    public static let micro: CGFloat = 4
    public static let tight: CGFloat = 8
    public static let regular: CGFloat = 12
    public static let section: CGFloat = 20
    public static let generous: CGFloat = 28
    public static let page: CGFloat = 36
}

public enum Motion {
    public static func settle(reduced: Bool) -> Animation? {
        reduced ? nil : .easeInOut(duration: 0.22)
    }
    public static func feedback(reduced: Bool) -> Animation? {
        reduced ? nil : .easeOut(duration: 0.12)
    }
}

public struct Surface<Content: View>: View {
    @Environment(\.colorScheme) private var scheme
    private let content: Content
    public init(@ViewBuilder content: () -> Content) { self.content = content() }
    public var body: some View {
        content.padding(Space.section).frame(maxWidth: .infinity, alignment: .leading)
            .background(Ink.surface.color(scheme), in: RoundedRectangle(cornerRadius: 16))
            .overlay(RoundedRectangle(cornerRadius: 16).strokeBorder(Ink.line.color(scheme), lineWidth: 0.5))
    }
}

public struct Eyebrow: View {
    @Environment(\.colorScheme) private var scheme
    let text: String
    public init(_ text: String) { self.text = text }
    public var body: some View {
        Text(text.uppercased()).font(.caption.weight(.semibold)).tracking(1.5)
            .foregroundStyle(Ink.secondary.color(scheme))
    }
}

private struct MotionPreviewKey: EnvironmentKey {
    static let defaultValue = false
}

public extension EnvironmentValues {
    var vigilReduceMotion: Bool {
        get { self[MotionPreviewKey.self] }
        set { self[MotionPreviewKey.self] = newValue }
    }
}

@propertyWrapper public struct ReducedMotion: DynamicProperty {
    @Environment(\.accessibilityReduceMotion) private var system
    @Environment(\.vigilReduceMotion) private var preview
    public init() {}
    public var wrappedValue: Bool { system || preview }
}
