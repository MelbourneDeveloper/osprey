import Foundation
import CryptoKit
import Darwin

public struct FileFingerprint: Codable, Equatable {
    public let path: String
    public let sha256: String
    public let bytes: Int64
    public let device: Int32
    public let inode: UInt64
    public let observed: Date
}
// Implements [SECURITY-COLLECTORS-FILES]: bounded local regular files only; no cloud downloads or link traversal.
public func fingerprint(_ url: URL) throws -> FileFingerprint {
    let properties = try url.resourceValues(forKeys: [.isUbiquitousItemKey])
    guard properties.isUbiquitousItem != true else { throw CollectionError.invalid("Cloud-managed files are excluded from hashing") }
    let descriptor = open(url.path, O_RDONLY | O_NOFOLLOW | O_NONBLOCK)
    guard descriptor >= 0 else { throw CollectionError.invalid("Cannot open selected file without following links") }
    defer { close(descriptor) }
    var before = stat(), after = stat()
    guard fstat(descriptor, &before) == 0, before.st_mode & S_IFMT == S_IFREG, before.st_size <= 32_000_000 else {
        throw CollectionError.invalid("Choose a local regular file no larger than 32 MB")
    }
    let digest = try digestFile(descriptor)
    guard fstat(descriptor, &after) == 0, sameFileVersion(before, after) else { throw CollectionError.invalid("File changed while hashing; retry inspection") }
    return FileFingerprint(path: url.path, sha256: digest, bytes: after.st_size, device: after.st_dev, inode: after.st_ino, observed: Date())
}
private func digestFile(_ descriptor: Int32) throws -> String {
    var hasher = SHA256(), buffer = [UInt8](repeating: 0, count: 65536), total = 0
    while true {
        try Task.checkCancellation()
        let count = read(descriptor, &buffer, buffer.count)
        guard count >= 0 else { throw CollectionError.invalid("File read failed") }
        if count == 0 { break }
        total += count
        guard total <= 32_000_000 else { throw CollectionError.invalid("File grew beyond hashing limit") }
        hasher.update(data: Data(buffer.prefix(count)))
    }
    return hasher.finalize().map { String(format: "%02x", $0) }.joined()
}
private func sameFileVersion(_ a: stat, _ b: stat) -> Bool {
    a.st_dev == b.st_dev && a.st_ino == b.st_ino && a.st_size == b.st_size &&
        a.st_mtimespec.tv_sec == b.st_mtimespec.tv_sec && a.st_mtimespec.tv_nsec == b.st_mtimespec.tv_nsec &&
        a.st_ctimespec.tv_sec == b.st_ctimespec.tv_sec && a.st_ctimespec.tv_nsec == b.st_ctimespec.tv_nsec
}
public struct BaselineStore {
    public let file: URL
    public init(directory: URL) { file = directory.appendingPathComponent("baselines.json") }
    public func load() throws -> [FileFingerprint] {
        guard FileManager.default.fileExists(atPath: file.path) else { return [] }
        return try JSONDecoder().decode([FileFingerprint].self, from: boundedData(file))
    }
    public func record(_ value: FileFingerprint) throws {
        let values = try load().filter { $0.path != value.path } + [value]
        guard values.count <= 1000 else { throw CollectionError.invalid("Baseline limit reached") }
        try privateWrite(values, to: file)
    }
}
