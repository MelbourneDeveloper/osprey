import Foundation

// Implements [IOS-TARGET-MACOS]: assertions against the real Mac. The boundary
// and volume checks hold on every Mac. Index checks need Spotlight indexing;
// run.sh requires them whenever `mdutil` reports it enabled.
enum SmokeCheck {
    static func run(requireIndex: Bool) -> Bool {
        let checks = boundaryChecks() + volumeChecks() + indexChecks(required: requireIndex)
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

    /// /Applications is indexed on every Mac with Spotlight on, and each app
    /// bundle carries its whole size in the index.
    private static func indexChecks(required: Bool) -> [(String, Bool)] {
        let root = "/Applications"
        let report = ospreyString(osprey_storage_report(root, 5))
        let rendered = ("report renders", report.hasPrefix("Storage report for \(root)\nVolume: "))
        guard required else { return [rendered] }
        let count = storageScan(root)
        let entries = (0..<count).map { (path: ospreyString(storageEntryPath($0)), bytes: storageEntryBytes($0)) }
        let largest = entries.first.map { URL(fileURLWithPath: $0.path).lastPathComponent } ?? "MISSING"
        return [
            rendered,
            ("indexed entries", count > 0),
            ("largest first", zip(entries, entries.dropFirst()).allSatisfy { $0.bytes >= $1.bytes }),
            ("entries under root", entries.allSatisfy { $0.path.hasPrefix(root + "/") && $0.bytes > 0 }),
            ("report ranks the largest", report.contains("%  \(largest)\n") || report.hasSuffix("%  \(largest)")),
            ("out-of-range entry", storageEntryBytes(count) == -1),
        ]
    }
}
