import Foundation

// Implements [IOS-TARGET-MACOS]: assertions against the real Mac. The boundary
// and volume checks read this machine; the walk checks build a folder whose
// sizes are known, so hidden files must be counted for them to pass.
enum SmokeCheck {
    static func run() -> Bool {
        let checks = boundaryChecks() + volumeChecks() + walkChecks()
        let failures = checks.compactMap { $0.1 ? nil : $0.0 }
        print(failures.isEmpty ? "OSPREY_MACOS_STORAGE_SMOKE_OK" : "FAIL: \(failures.joined(separator: ", "))")
        return failures.isEmpty
    }

    /// The exports and imports cross the C boundary the iPhone apps use.
    private static func boundaryChecks() -> [(String, Bool)] {
        [
            ("initialize once", osprey_main() == 0),
            ("string export", ospreyString(osprey_storage_Disk_formatBytes(4_143_181_824)) == "4.1 GB"),
            ("unknown size", ospreyString(osprey_storage_Disk_formatBytes(-1)) == "unknown"),
            ("integer export", osprey_storage_Disk_percent(450, 750) == 60),
            ("zero divisor", osprey_storage_Disk_percent(1, 0) == 0),
        ]
    }

    private static func volumeChecks() -> [(String, Bool)] {
        let home = NSHomeDirectory()
        let total = storageVolume(home, "total")
        let available = storageVolume(home, "available")
        let important = storageVolume(home, "important")
        return [
            ("volume total", total > 0),
            ("free space fits", available >= 0 && available <= total),
            ("important space fits", important >= 0 && important <= total),
            ("unknown measure", storageVolume(home, "bogus") == -1),
        ]
    }

    /// `Projects` holds 10 KB in sight and 300 KB inside a hidden `.git`; a
    /// hidden `.cache` holds 200 KB; `Locked` cannot be read.
    private static func walkChecks() -> [(String, Bool)] {
        guard let root = try? Fixture.make() else { return [("fixture", false)] }
        defer { Fixture.remove(root) }
        let report = ospreyString(osprey_storage_report(root, 5))
        let top = entries(storageScan(root))
        let unreadable = storageUnreadable()
        let inside = entries(storageScan(root + "/Projects"))
        return [
            ("ranks every entry", top.map(\.name) == ["Projects", ".cache", "notes.txt"]),
            ("hidden folder counted", top.first { $0.name == ".cache" }.map { $0.bytes >= 200_000 } ?? false),
            ("hidden files inside a folder counted", top.first.map { $0.bytes >= 310_000 } ?? false),
            ("unreadable folder reported", unreadable == 1),
            ("drill-down shows the hidden folder", inside.first?.name == ".git" && inside.first.map { $0.bytes >= 300_000 } ?? false),
            ("report includes hidden entries", report.contains("%  .cache\n") && report.contains("%  .git")),
            ("report names unreadable folders", report.contains("Folders that could not be read: 1.")),
            ("out-of-range entry", storageEntryBytes(Int64(inside.count)) == -1),
        ]
    }

    private static func entries(_ count: Int64) -> [(name: String, bytes: Int64)] {
        (0..<max(count, 0)).map {
            (name: URL(fileURLWithPath: ospreyString(storageEntryPath($0))).lastPathComponent, bytes: storageEntryBytes($0))
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
        return root.path
    }

    static func remove(_ root: String) {
        try? FileManager.default.setAttributes([.posixPermissions: 0o755], ofItemAtPath: root + "/Locked")
        try? FileManager.default.removeItem(atPath: root)
    }
}
