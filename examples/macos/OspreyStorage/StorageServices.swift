import Foundation

// Implements [IOS-HOST-IMPORTS] and [IOS-TARGET-MACOS]. Osprey calls these C
// functions synchronously on the host thread: volume capacities come from
// Foundation, and folder sizes from measuring the tree, hidden files included.

/// Volume capacity in bytes, or -1. `important` adds the space macOS can purge
/// on demand (caches, local snapshots), which is what System Settings shows as
/// available; `available` excludes it.
@_cdecl("storage_volume")
func storageVolume(_ path: UnsafePointer<CChar>?, _ measure: UnsafePointer<CChar>?) -> Int64 {
    guard let path, let measure else { return -1 }
    let keys: Set<URLResourceKey> = [
        .volumeTotalCapacityKey, .volumeAvailableCapacityKey, .volumeAvailableCapacityForImportantUsageKey,
    ]
    guard let values = try? URL(fileURLWithPath: String(cString: path)).resourceValues(forKeys: keys) else {
        return -1
    }
    switch String(cString: measure) {
    case "total": return values.volumeTotalCapacity.map(Int64.init) ?? -1
    case "available": return values.volumeAvailableCapacity.map(Int64.init) ?? -1
    case "important": return values.volumeAvailableCapacityForImportantUsage ?? -1
    default: return -1
    }
}

/// The entries and unreadable-folder count of the most recent scan.
nonisolated(unsafe) private var scanned = DiskWalk.Ranking(entries: [], unreadable: 0)

/// Ranks the direct children of `root` by the bytes stored under them, hidden
/// files included. Returns the entry count, or -1 when `root` cannot be listed.
@_cdecl("storage_scan")
func storageScan(_ root: UnsafePointer<CChar>?) -> Int64 {
    guard let root, let ranking = DiskWalk.ranking(under: String(cString: root)) else {
        scanned = DiskWalk.Ranking(entries: [], unreadable: 0)
        return -1
    }
    scanned = ranking
    return Int64(ranking.entries.count)
}

/// Osprey may keep a returned string for the rest of the process, so each call
/// hands over its own heap copy and Swift never frees it. [IOS-HOST-IMPORTS]
@_cdecl("storage_entry_path")
func storageEntryPath(_ index: Int64) -> UnsafePointer<CChar>? {
    UnsafePointer(strdup(scanned.entries.indices.contains(Int(index)) ? scanned.entries[Int(index)].path : ""))
}

@_cdecl("storage_entry_bytes")
func storageEntryBytes(_ index: Int64) -> Int64 {
    scanned.entries.indices.contains(Int(index)) ? scanned.entries[Int(index)].bytes : -1
}

/// Folders the last scan could not read, so their contents are not counted.
@_cdecl("storage_unreadable")
func storageUnreadable() -> Int64 {
    scanned.unreadable
}

/// macOS keeps no record of a folder's size, and Spotlight skips hidden files
/// entirely, so the host measures the tree: Foundation's enumerator lists each
/// directory in bulk, and every file contributes its allocated size. Hidden
/// files and package contents count; a hard-linked file counts once; mounted
/// volumes below the root are not entered.
enum DiskWalk {
    typealias Entry = (path: String, bytes: Int64)

    struct Ranking {
        var entries: [Entry]
        var unreadable: Int64
    }

    /// Per-directory rankings from the last walk: the root and each child.
    nonisolated(unsafe) private static var rankings: [String: Ranking] = [:]

    private static let keys: [URLResourceKey] = [
        .isDirectoryKey, .isVolumeKey, .totalFileAllocatedSizeKey, .linkCountKey, .fileResourceIdentifierKey,
    ]

    /// A drill-down into a child of the previous root is answered from that
    /// walk, once; any other root is walked afresh.
    static func ranking(under root: String) -> Ranking? {
        let base = root.hasSuffix("/") ? root : root + "/"
        if let cached = rankings.removeValue(forKey: base) {
            rankings = [:]
            return cached
        }
        guard let walked = walk(base) else { return nil }
        rankings = walked.filter { $0.key != base }
        return walked[base]
    }

    /// Sums each file into its direct child of `base` and that child's own
    /// child, and attributes each unreadable folder the same way. Every child
    /// gets a ranking, empty for a file.
    private static func walk(_ base: String) -> [String: Ranking]? {
        let root = URL(fileURLWithPath: base, isDirectory: true)
        var tally = Tally(base: base, resolved: canonical(base))
        guard (try? root.checkResourceIsReachable()) == true,
            let walker = FileManager.default.enumerator(
                at: root, includingPropertiesForKeys: keys, options: [],
                errorHandler: { url, _ in tally.unreadable(url.path) })
        else { return nil }
        var linked = Set<NSObject>()
        for case let url as URL in walker {
            guard let values = try? url.resourceValues(forKeys: Set(keys)) else { continue }
            if values.isDirectory == true {
                if values.isVolume == true { walker.skipDescendants() }
                continue
            }
            if (values.linkCount ?? 1) > 1, let id = values.fileResourceIdentifier as? NSObject,
                !linked.insert(id).inserted
            {
                continue
            }
            tally.add(url.path, Int64(values.totalFileAllocatedSize ?? 0))
        }
        return tally.rankings()
    }
}

/// The path with every symlink resolved, as the enumerator reports it.
/// `URL.resolvingSymlinksInPath` would strip `/private` back off `/var`.
private func canonical(_ path: String) -> String {
    guard let resolved = realpath(path, nil) else { return path }
    defer { free(resolved) }
    return String(cString: resolved)
}

/// Running totals for one walk, keyed by directory with a trailing slash in the
/// caller's spelling. The enumerator reports resolved paths (`/var` arrives as
/// `/private/var`), so entries are matched against the resolved root.
private struct Tally {
    let base: String
    let resolved: String
    private var bytes: [String: [String: Int64]] = [:]
    private var blocked: [String: Int64] = [:]

    init(base: String, resolved: String) {
        self.base = base
        self.resolved = resolved.hasSuffix("/") ? resolved : resolved + "/"
        bytes[base] = [:]
    }

    /// The child of `base` holding `path`, and that child's child, if any.
    private func levels(_ path: String) -> (child: String, inside: String?)? {
        let prefix = path.hasPrefix(resolved) ? resolved : base
        guard path.hasPrefix(prefix) else { return nil }
        let parts = path.dropFirst(prefix.count).split(separator: "/", maxSplits: 2).map(String.init)
        guard let first = parts.first else { return nil }
        let child = base + first
        return (child, parts.count > 1 ? child + "/" + parts[1] : nil)
    }

    mutating func add(_ path: String, _ size: Int64) {
        guard let level = levels(path) else { return }
        bytes[base, default: [:]][level.child, default: 0] += size
        if bytes[level.child + "/"] == nil { bytes[level.child + "/"] = [:] }
        if let inside = level.inside { bytes[level.child + "/", default: [:]][inside, default: 0] += size }
    }

    /// Records a folder the walk could not read, against the root and against
    /// the child holding it, and keeps walking.
    mutating func unreadable(_ path: String) -> Bool {
        blocked[base, default: 0] += 1
        if let level = levels(path) { blocked[level.child + "/", default: 0] += 1 }
        return true
    }

    func rankings() -> [String: DiskWalk.Ranking] {
        bytes.reduce(into: [:]) { result, level in
            let entries = level.value.map { (path: $0.key, bytes: $0.value) }.sorted { $0.bytes > $1.bytes }
            result[level.key] = DiskWalk.Ranking(entries: entries, unreadable: blocked[level.key] ?? 0)
        }
    }
}
