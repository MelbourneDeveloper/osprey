import Foundation

// Implements [SECURITY-COLLECTORS-ADVISORIES]: local receipts are claims, not authenticated package identity.
public func packageReceipt(_ data: Data, path: String) throws -> Evidence {
    guard data.count <= 1_000_000,
          let receipt = try JSONSerialization.jsonObject(with: data) as? [String: Any],
          let source = receipt["source"] as? [String: Any], let tap = source["tap"] as? String else {
        throw CollectionError.invalid("Missing or malformed Homebrew receipt provenance")
    }
    let url = URL(fileURLWithPath: path), version = url.lastPathComponent
    return Evidence(id: path, title: url.deletingLastPathComponent().lastPathComponent,
        detail: "Installed \(version) · \(tap). Local receipt; support status and vulnerability applicability remain unknown.",
        path: path, source: "Homebrew INSTALL_RECEIPT.json", facts: ["version": version, "tap": tap,
        "installedOnRequest": receipt["installed_on_request"].map { String(describing: $0) } ?? "unknown",
        "bottle": receipt["poured_from_bottle"].map { String(describing: $0) } ?? "unknown", "kind": "formula"])
}
