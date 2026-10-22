#if os(macOS)
import Foundation

// No shell, bounded runtime/output, private temporary files, no inherited command-line secrets.
public func runCommand(_ executable: String, _ arguments: [String], timeout: TimeInterval = 10, input: Data? = nil) throws -> String {
    let directory = FileManager.default.temporaryDirectory.appendingPathComponent(UUID().uuidString)
    try FileManager.default.createDirectory(at: directory, withIntermediateDirectories: false, attributes: [.posixPermissions: 0o700])
    defer { try? FileManager.default.removeItem(at: directory) }
    let output = directory.appendingPathComponent("output")
    let incoming = directory.appendingPathComponent("input")
    try Data().write(to: output); try (input ?? Data()).write(to: incoming)
    let writer = try FileHandle(forWritingTo: output), reader = try FileHandle(forReadingFrom: incoming)
    defer { try? writer.close(); try? reader.close() }
    let process = Process()
    process.executableURL = URL(fileURLWithPath: executable); process.arguments = arguments
    process.standardOutput = writer; process.standardError = writer; process.standardInput = reader
    process.environment = ["PATH": "/usr/bin:/bin:/usr/sbin:/sbin", "LC_ALL": "C", "HOME": NSHomeDirectory()]
    try process.run()
    defer { if process.isRunning { kill(process.processIdentifier, SIGKILL); process.waitUntilExit() } }
    try awaitExit(process, output: output, timeout: timeout)
    return try commandResult(output, status: process.terminationStatus, executable: executable)
}
private func commandResult(_ output: URL, status: Int32, executable: String) throws -> String {
    guard (try output.resourceValues(forKeys: [.fileSizeKey])).fileSize ?? 0 <= 8_000_000 else {
        throw CollectionError.invalid("Command exceeded 8 MB")
    }
    let data = try Data(contentsOf: output)
    guard let text = String(data: data, encoding: .utf8) else { throw CollectionError.invalid("Command returned invalid UTF-8") }
    guard status == 0 else { throw CollectionError.invalid("\(URL(fileURLWithPath: executable).lastPathComponent) exited \(status): \(text.prefix(500))") }
    return text
}
private func awaitExit(_ process: Process, output: URL, timeout: TimeInterval) throws {
    let deadline = Date().addingTimeInterval(timeout)
    while process.isRunning {
        let size = (try output.resourceValues(forKeys: [.fileSizeKey])).fileSize ?? 0
        if Date() > deadline || size > 8_000_000 || Task.isCancelled {
            process.terminate()
            if process.isRunning { kill(process.processIdentifier, SIGKILL) }
            process.waitUntilExit()
            throw CollectionError.invalid("Collection cancelled, timed out, or exceeded 8 MB")
        }
        Thread.sleep(forTimeInterval: 0.02)
    }
    process.waitUntilExit()
}
#endif
