import Foundation

// Implements [SECURITY-COLLECTORS-NETWORK]: metadata snapshot, no traffic interception or PID-based authority.
public func decodeSockets(_ text: String) throws -> [Evidence] {
    var records: [Evidence] = [], process: [String: String] = [:], socket: [String: String] = [:]
    for line in text.split(separator: "\n") {
        guard let key = line.first else { continue }
        let value = String(line.dropFirst())
        if key == "p" || key == "f" {
            if !socket.isEmpty { records.append(try socketEvidence(process, socket)); socket = [:] }
        }
        switch key {
        case "p": guard Int(value) != nil else { throw CollectionError.invalid("Invalid socket PID") }; process = ["pid": value]
        case "c": process["command"] = value
        case "f", "t", "P", "n": socket[String(key)] = value
        case "T": if value.hasPrefix("ST=") { socket["state"] = String(value.dropFirst(3)) }
        default: break
        }
    }
    if !socket.isEmpty { records.append(try socketEvidence(process, socket)) }
    return records
}
private func socketEvidence(_ process: [String: String], _ socket: [String: String]) throws -> Evidence {
    guard let pid = process["pid"], let descriptor = socket["f"], let endpoint = socket["n"] else {
        throw CollectionError.invalid("Incomplete socket metadata; no clean-network conclusion is available")
    }
    return Evidence(id: "\(pid):\(descriptor):\(endpoint)", title: process["command"] ?? "Process \(pid)",
        detail: "\(socket["P"] ?? "Unknown protocol") · \(endpoint) · \(socket["state"] ?? "State unavailable")",
        source: "lsof numeric socket metadata", facts: ["pid": pid, "endpoint": endpoint,
        "protocol": socket["P"] ?? "unknown", "state": socket["state"] ?? "unknown", "attribution": "unverified PID snapshot"])
}
#if os(macOS)
public func networkSnapshot(pid: Int32? = nil) throws -> ScanSection {
    let scope = pid.map { ["-a", "-p", String($0)] } ?? []
    let text = try runCommand("/usr/sbin/lsof", ["-nP", "-i", "-FpcfnPtT"] + scope, timeout: 15)
    return ScanSection(category: .network, state: .partial,
        reason: "Readable open IP sockets at one moment. Numeric endpoints only; no payload, byte counts, domains, continuous coverage or blocking. PID attribution can race process exit/reuse. Listeners and legitimate connections are not spyware evidence.",
        records: try decodeSockets(text))
}
#endif
