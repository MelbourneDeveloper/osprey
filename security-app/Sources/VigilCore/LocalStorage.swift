import Foundation

// Implements [SECURITY-PRIVACY]: bounded local records, no network synchronization.
public func privateWrite<T: Encodable>(_ value: T, to file: URL) throws {
    let directory = file.deletingLastPathComponent()
    try FileManager.default.createDirectory(at: directory, withIntermediateDirectories: true, attributes: [.posixPermissions: 0o700])
    try FileManager.default.setAttributes([.posixPermissions: 0o700], ofItemAtPath: directory.path)
    let data = try JSONEncoder().encode(value)
    guard data.count <= 8_000_000 else { throw CollectionError.invalid("Local record exceeds 8 MB") }
    #if os(iOS)
    try data.write(to: file, options: [.atomic, .completeFileProtectionUnlessOpen])
    #else
    try data.write(to: file, options: .atomic)
    #endif
    try FileManager.default.setAttributes([.posixPermissions: 0o600], ofItemAtPath: file.path)
}
public func boundedData(_ file: URL, limit: Int = 8_000_000) throws -> Data {
    let handle = try FileHandle(forReadingFrom: file)
    defer { try? handle.close() }
    let data = try handle.read(upToCount: limit + 1) ?? Data()
    guard data.count <= limit else { throw CollectionError.invalid("Input exceeds \(limit) bytes") }
    return data
}
public struct HistoryEntry: Identifiable {
    public let file: URL
    public let created: Date
    public var id: String { file.lastPathComponent }
}
public struct HistoryStore {
    public let directory: URL
    public init(directory: URL) { self.directory = directory }
    public func entries() throws -> [HistoryEntry] {
        guard FileManager.default.fileExists(atPath: directory.path) else { return [] }
        return try FileManager.default.contentsOfDirectory(at: directory, includingPropertiesForKeys: nil)
            .filter { $0.pathExtension == "json" }.map { HistoryEntry(file: $0, created: try readReport($0).created) }
            .sorted { $0.created > $1.created }
    }
    public func save(_ report: ScanReport, now: Date = Date(), limit: Int = 30) throws {
        guard limit > 0, limit <= 100 else { throw CollectionError.invalid("Invalid history retention limit") }
        try privateWrite(report, to: directory.appendingPathComponent(UUID().uuidString + ".json"))
        for (index, entry) in try entries().enumerated() where index >= limit || now.timeIntervalSince(entry.created) > 30 * 86400 {
            try FileManager.default.removeItem(at: entry.file)
        }
    }
}
public struct FindingStore {
    public let file: URL
    public init(directory: URL) { file = directory.appendingPathComponent("findings.json") }
    public func load(now: Date = Date()) throws -> [Finding] {
        guard FileManager.default.fileExists(atPath: file.path) else { return [] }
        return try JSONDecoder().decode([Finding].self, from: boundedData(file)).filter { now.timeIntervalSince($0.lastSeen) <= 90 * 86400 }.map { expireSuppression($0, now: now) }
    }
    public func save(_ findings: [Finding], now: Date = Date()) throws {
        let retained = findings.filter { now.timeIntervalSince($0.lastSeen) <= 90 * 86400 }
        guard retained.count <= 10_000 else { throw CollectionError.invalid("Finding retention limit reached") }
        try privateWrite(retained, to: file)
    }
}
