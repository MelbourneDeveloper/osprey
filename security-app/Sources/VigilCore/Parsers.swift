import Foundation

public func decodeLaunchItem(_ data: Data, path: String) throws -> Evidence {
    guard data.count <= 1_000_000,
          let item = try PropertyListSerialization.propertyList(from: data, format: nil) as? [String: Any],
          let label = item["Label"] as? String,
          let program = (item["Program"] as? String) ?? (item["ProgramArguments"] as? [String])?.first,
          !program.isEmpty else { throw CollectionError.invalid("Invalid startup definition: \(path)") }
    return Evidence(id: path, title: label, detail: "Configured startup item; this does not establish execution. Definition: \(path)",
                    path: program, source: "launchd property list", facts: [
                        "hidden": String(program.split(separator: "/").contains { $0.hasPrefix(".") }),
                        "runAtLoad": String(item["RunAtLoad"] as? Bool ?? false), "definition": path
                    ])
}
public func parseProcesses(_ text: String) throws -> [Evidence] {
    try text.split(separator: "\n").map { line in
        let fields = line.split(maxSplits: 7, omittingEmptySubsequences: true, whereSeparator: { $0 == " " || $0 == "\t" })
        guard fields.count == 8, let pid = Int(fields[0]), Int(fields[1]) != nil else {
            throw CollectionError.invalid("Unrecognized process snapshot format")
        }
        let started = fields[2...6].joined(separator: " "), path = String(fields[7])
        return Evidence(id: "\(pid):\(started)", title: URL(fileURLWithPath: path).lastPathComponent,
                        detail: "PID \(pid) · parent \(fields[1]) · started \(started). Snapshot; process may have exited.",
                        path: path, source: "ps: pid, ppid, lstart, comm (no arguments)",
                        facts: ["pid": String(pid), "parent": String(fields[1]), "started": started])
    }
}
