import CoreServices
import Foundation

// Implements [IOS-HOST-IMPORTS] and [IOS-TARGET-MACOS]. Osprey calls these C
// functions synchronously on the host thread. Every number comes from records
// macOS already keeps: Foundation's volume capacities and the Spotlight index.

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

/// The entries of the most recent scan, largest first.
nonisolated(unsafe) private var scanned: [SpotlightIndex.Entry] = []

/// Ranks the direct children of `root` by the bytes Spotlight attributes to
/// them. Returns the entry count, or -1 when the index cannot be queried.
@_cdecl("storage_scan")
func storageScan(_ root: UnsafePointer<CChar>?) -> Int64 {
    guard let root, let entries = SpotlightIndex.entries(under: String(cString: root)) else {
        scanned = []
        return -1
    }
    scanned = entries
    return Int64(entries.count)
}

/// Osprey may keep a returned string for the rest of the process, so each call
/// hands over its own heap copy and Swift never frees it. [IOS-HOST-IMPORTS]
@_cdecl("storage_entry_path")
func storageEntryPath(_ index: Int64) -> UnsafePointer<CChar>? {
    UnsafePointer(strdup(scanned.indices.contains(Int(index)) ? scanned[Int(index)].path : ""))
}

@_cdecl("storage_entry_bytes")
func storageEntryBytes(_ index: Int64) -> Int64 {
    scanned.indices.contains(Int(index)) ? scanned[Int(index)].bytes : -1
}

/// Spotlight already records each indexed item's size, including whole app and
/// document bundles (it does not index inside packages, so nothing is counted
/// twice). A query reads those records; no directory is opened or walked.
enum SpotlightIndex {
    typealias Entry = (path: String, bytes: Int64)

    /// Per-directory rankings from the last query: the root and each child.
    nonisolated(unsafe) private static var rankings: [String: [Entry]] = [:]

    /// A drill-down into a child of the previous root is answered from that
    /// query, once; any other root is a fresh query.
    static func entries(under root: String) -> [Entry]? {
        let base = root.hasSuffix("/") ? root : root + "/"
        if let cached = rankings.removeValue(forKey: base) {
            rankings = [:]
            return cached
        }
        guard let items = query(base) else { return nil }
        let ranked = rank(items, under: base)
        rankings = ranked.filter { $0.key != base }
        return ranked[base] ?? []
    }

    /// Sums each item into its direct child of `base` and into that child's own
    /// child. Every child gets a ranking, empty for a file or a package.
    private static func rank(_ items: [Entry], under base: String) -> [String: [Entry]] {
        var totals: [String: [String: Int64]] = [base: [:]]
        for item in items where item.path.hasPrefix(base) {
            let parts = item.path.dropFirst(base.count).split(separator: "/", maxSplits: 2).map(String.init)
            guard let first = parts.first else { continue }
            let child = base + first
            let inside = child + "/"
            totals[base, default: [:]][child, default: 0] += item.bytes
            if totals[inside] == nil { totals[inside] = [:] }
            if parts.count > 1 { totals[inside, default: [:]][inside + parts[1], default: 0] += item.bytes }
        }
        return totals.mapValues { sizes in
            sizes.map { (path: $0.key, bytes: $0.value) }.sorted { $0.bytes > $1.bytes }
        }
    }

    /// Every indexed item under `base` with its size on disk: the physical size
    /// when Spotlight records it (APFS compression and clones included),
    /// otherwise its logical size.
    private static func query(_ base: String) -> [Entry]? {
        let sizes = ["kMDItemPhysicalSize", "kMDItemFSSize"] as CFArray
        guard let query = MDQueryCreate(kCFAllocatorDefault, "kMDItemPhysicalSize > 0 || kMDItemFSSize > 0" as CFString, sizes, nil) else {
            return nil
        }
        MDQuerySetSearchScope(query, [base] as CFArray, 0)
        guard MDQueryExecute(query, CFOptionFlags(kMDQuerySynchronous.rawValue)) else { return nil }
        return (0..<MDQueryGetResultCount(query)).compactMap { index in
            guard let raw = MDQueryGetResultAtIndex(query, index),
                let path = MDItemCopyAttribute(Unmanaged<MDItem>.fromOpaque(raw).takeUnretainedValue(), kMDItemPath) as? String
            else { return nil }
            return (path: path, bytes: size(query, index, "kMDItemPhysicalSize") ?? size(query, index, "kMDItemFSSize") ?? 0)
        }
    }

    private static func size(_ query: MDQuery, _ index: CFIndex, _ attribute: String) -> Int64? {
        MDQueryGetAttributeValueOfResultAtIndex(query, attribute as CFString, index).map {
            (Unmanaged<CFNumber>.fromOpaque($0).takeUnretainedValue() as NSNumber).int64Value
        }
    }
}
