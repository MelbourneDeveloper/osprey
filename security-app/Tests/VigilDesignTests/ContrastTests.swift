import XCTest
@testable import VigilDesign

final class ContrastTests: XCTestCase {
    // Readable status labels are required on every shipped surface in both appearances.
    func testSemanticTextMeetsNormalTextContrast() {
        let text: [Ink] = [.primary, .secondary, .accent, .good, .caution, .danger]
        for dark in [false, true] {
            for foreground in text {
                for background in [Ink.canvas, .surface, .inset] {
                    let ratio = contrast(foreground.hex(dark: dark), background.hex(dark: dark))
                    XCTAssertGreaterThanOrEqual(ratio, 4.5, "\(foreground) on \(background), dark=\(dark)")
                }
            }
        }
    }

    func testPrimaryButtonTextMeetsContrast() {
        for dark in [false, true] {
            XCTAssertGreaterThanOrEqual(contrast(Ink.canvas.hex(dark: dark), Ink.accent.hex(dark: dark)), 4.5)
        }
    }

    private func contrast(_ first: UInt32, _ second: UInt32) -> Double {
        let values = [luminance(first), luminance(second)].sorted()
        return (values[1] + 0.05) / (values[0] + 0.05)
    }

    private func luminance(_ hex: UInt32) -> Double {
        let channels = [Double(hex >> 16), Double((hex >> 8) & 255), Double(hex & 255)]
        let linear = channels.map { value -> Double in
            let normalized = value / 255
            return normalized <= 0.04045 ? normalized / 12.92 : pow((normalized + 0.055) / 1.055, 2.4)
        }
        return zip(linear, [0.2126, 0.7152, 0.0722]).map(*).reduce(0, +)
    }
}
