import Foundation

// Implements [IOS-TARGET-MACOS]: assertions against the real device. The
// boundary and volume checks read this machine; the browser checks walk a
// folder whose contents are known and drive Osprey through the same events
// the app sends, so hidden files must be measured and filtered for them to pass.
enum SmokeCheck {
    static func run() -> Bool {
        let checks = boundaryChecks() + volumeChecks() + browserChecks()
        let failures = checks.compactMap { $0.1 ? nil : $0.0 }
        let verdict = failures.isEmpty ? "OSPREY_MACOS_STORAGE_SMOKE_OK" : "FAIL: \(failures.joined(separator: ", "))"
        print(verdict)
        #if os(iOS)
        let documents = FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0]
        try? verdict.write(to: documents.appendingPathComponent("storage-smoke-result.txt"), atomically: true, encoding: .utf8)
        #endif
        return failures.isEmpty
    }

    /// The exports and imports cross the C boundary the iPhone apps use.
    private static func boundaryChecks() -> [(String, Bool)] {
        [
            ("initialize once", osprey_main() == 0),
            ("string export", text(osprey_storage_Disk_formatBytes(4_143_181_824)) == "4.1 GB"),
            ("grouped count", text(osprey_storage_Disk_formatCount(1_234_005)) == "1,234,005"),
            ("integer export", osprey_storage_Disk_percent(450, 750) == 60),
            ("zero divisor", osprey_storage_Disk_percent(1, 0) == 0),
        ]
    }

    private static func volumeChecks() -> [(String, Bool)] {
        let home = NSHomeDirectory()
        let total = storageVolume(home, "total")
        let available = storageVolume(home, "available")
        return [
            ("volume total", total > 0),
            ("free space fits", available >= 0 && available <= total),
            ("unknown measure", storageVolume(home, "bogus") == -1),
            ("volume name", !text(storageVolumeName(home)).isEmpty),
        ]
    }

    /// `Projects` holds 10 KB in sight and 300 KB inside a hidden `.git`; a
    /// hidden `.cache` holds 200 KB; `Locked` cannot be read.
    private static func browserChecks() -> [(String, Bool)] {
        guard let root = try? Fixture.make() else { return [("fixture", false)] }
        defer { Fixture.remove(root) }
        let started = envelope(osprey_storage_start(root))
        DiskTree.current = DiskTree.measure(root) { _, _ in }
        let ready = dispatch(started, ["type": "measured", "path": root])
        let top = entries(storageScan(root))
        let unreadable = storageUnreadable()
        let projects = top.first { $0.name == "Projects" }
        let sighted = dispatch(ready, ["type": "hidden", "value": "false"])
        let inside = dispatch(ready, ["type": "open", "value": root + "/Projects"])
        return [
            ("asks for the walk", started.text.contains("\"commands\":[{\"kind\":\"measure\",\"path\":\"\(root)\"}]")),
            ("ranks every entry", Set(top.map(\.name)) == ["Projects", ".cache", "notes.txt", "Locked"]),
            ("hidden files inside a folder counted", (projects?.bytes ?? 0) >= 310_000),
            ("hidden files inside a folder can be left out", (projects?.visible ?? 0) >= 10_000 && (projects?.visible ?? 0) < 300_000),
            ("hidden folder flagged", top.first { $0.name == ".cache" }?.hidden == true),
            ("unreadable folder reported", unreadable == 1),
            ("screen shows hidden entries", ready.text.contains("\"text\":\".cache\"")),
            ("screen explains the unreadable folder", ready.text.contains("1 folder here could not be read")),
            ("hidden files can be turned off", !sighted.text.contains("\"text\":\".cache\"") && sighted.text.contains("\"text\":\"Projects\"")),
            ("folders open", inside.text.contains("\"text\":\".git\"") && inside.text.contains("Enclosing Folder")),
            ("out-of-range entry", storageEntryBytes(Int64(top.count)) == -1),
        ]
    }

    private struct Envelope {
        let text: String
        let model: String
    }

    private static func envelope(_ pointer: UnsafePointer<CChar>?) -> Envelope {
        let raw = text(pointer)
        let json = (try? JSONSerialization.jsonObject(with: Data(raw.utf8))) as? [String: Any]
        return Envelope(text: raw, model: json?["model"] as? String ?? "")
    }

    private static func dispatch(_ from: Envelope, _ event: [String: Any]) -> Envelope {
        let data = (try? JSONSerialization.data(withJSONObject: event)) ?? Data()
        return envelope(osprey_storage_dispatch(from.model, String(decoding: data, as: UTF8.self)))
    }

    private static func text(_ pointer: UnsafePointer<CChar>?) -> String { pointer.map { String(cString: $0) } ?? "" }

    private static func entries(_ count: Int64) -> [(name: String, bytes: Int64, visible: Int64, hidden: Bool)] {
        (0..<max(count, 0)).map {
            (URL(fileURLWithPath: text(storageEntryPath($0))).lastPathComponent, storageEntryBytes($0),
             storageEntryVisible($0), storageEntryHidden($0))
        }
    }
}

/// A throwaway folder with known contents, hidden files included.
private enum Fixture {
    static let files: [(String, Int)] = [
        ("Projects/app.swift", 10_000), ("Projects/.git/pack", 300_000), (".cache/blob", 200_000),
        ("notes.txt", 1_000), ("Locked/secret", 50_000),
    ]

    static func make() throws -> String {
        let root = FileManager.default.temporaryDirectory.appendingPathComponent("osprey-storage-\(UUID().uuidString)")
        for (path, size) in files {
            let file = root.appendingPathComponent(path)
            try FileManager.default.createDirectory(at: file.deletingLastPathComponent(), withIntermediateDirectories: true)
            try Data((0..<size).map { UInt8($0 % 251) }).write(to: file)
        }
        try FileManager.default.setAttributes([.posixPermissions: 0o000], ofItemAtPath: root.appendingPathComponent("Locked").path)
        return canonical(root.path)
    }

    static func remove(_ root: String) {
        try? FileManager.default.setAttributes([.posixPermissions: 0o755], ofItemAtPath: root + "/Locked")
        try? FileManager.default.removeItem(atPath: root)
    }
}
