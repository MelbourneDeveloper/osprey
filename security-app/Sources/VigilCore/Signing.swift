#if os(macOS)
import Foundation

// Implements [SECURITY-COLLECTORS-IDENTITY]. Validate before trusting entitlement data.
public func inspectSignature(_ record: Evidence) throws -> Evidence {
    var result = record
    do { _ = try runCommand("/usr/bin/codesign", ["--verify", "--deep", "--strict", record.path], timeout: 30) }
    catch {
        result.facts["signing"] = "unverified"
        result.facts["sandbox"] = "unknown"
        result.detail += "\nSignature verification did not succeed: \(error.localizedDescription). This is not a malware verdict."
        return result
    }
    let text = try runCommand("/usr/bin/codesign", ["-d", "--entitlements", ":-", record.path])
    let entitlements = try decodeEntitlements(text)
    result.facts["signing"] = "valid"
    result.facts["sandbox"] = (entitlements["com.apple.security.app-sandbox"] as? Bool == true) ? "present" : "absent"
    result.facts["entitlements"] = entitlements.keys.sorted().joined(separator: "\n")
    result.source = "codesign --verify --deep --strict; validated on-disk entitlement evidence"
    result.observed = Date()
    return result
}
func decodeEntitlements(_ text: String) throws -> [String: Any] {
    guard let start = text.range(of: "<?xml"), let end = text.range(of: "</plist>") else {
        let lines = text.split(separator: "\n")
        guard !lines.isEmpty, lines.allSatisfy({ $0.hasPrefix("Executable=") || $0.hasPrefix("warning: Specifying ':'") }) else {
            throw CollectionError.invalid("Unrecognized entitlement output; sandbox status is unknown")
        }
        return [:]
    }
    let data = Data(text[start.lowerBound..<end.upperBound].utf8)
    guard let dictionary = try PropertyListSerialization.propertyList(from: data, format: nil) as? [String: Any] else {
        throw CollectionError.invalid("Entitlements are not a dictionary")
    }
    return dictionary
}
#endif
