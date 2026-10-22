#if os(macOS)
import Foundation

// Implements [SECURITY-COLLECTORS-IDENTITY]. All snapshots are read-only and bounded.
public func collect(_ area: Area, folder: URL? = nil) -> ScanSection {
    do {
        switch area {
        case .apps: return try appInventory()
        case .processes: return ScanSection(category: area, state: .partial, reason: "Point-in-time process snapshot; no execution monitoring or behavior attribution.", records: try parseProcesses(runCommand("/bin/ps", ["-axo", "pid=,ppid=,lstart=,comm="])))
        case .startup: return try startupInventory()
        case .files: return try hiddenInventory(folder)
        case .packages: return try packageInventory()
        case .protection: return protectionInventory()
        case .network: return unavailable(area, "Network Extension is not installed. Traffic inspection and blocking are unavailable.")
        case .permissions: return permissionGuidance()
        case .advisories: return unavailable(area, "No verified advisory feed is configured. No security-history or vulnerability verdict is available.")
        }
    } catch {
        return ScanSection(category: area, state: .failed, reason: error.localizedDescription, records: [])
    }
}
private func unavailable(_ area: Area, _ reason: String) -> ScanSection {
    ScanSection(category: area, state: .unsupported, reason: reason, records: [])
}
private func children(_ path: String) throws -> [URL] {
    try FileManager.default.contentsOfDirectory(at: URL(fileURLWithPath: path), includingPropertiesForKeys: [.isDirectoryKey, .isSymbolicLinkKey, .isHiddenKey])
}
public func appInventory(roots: [String] = ["/Applications", "/System/Applications", NSHomeDirectory() + "/Applications"]) throws -> ScanSection {
    var records: [Evidence] = [], issues: [String] = []
    for root in roots where FileManager.default.fileExists(atPath: root) && records.count < 10_000 {
        do {
            for item in try children(root).sorted(by: { $0.path < $1.path }).prefix(1000) {
                try Task.checkCancellation()
                do {
                    for app in try appCandidates(item).prefix(10_000 - records.count) {
                        do { records.append(try appRecord(app)) }
                        catch { issues.append("\(app.path): \(error.localizedDescription)") }
                    }
                } catch { issues.append("\(item.path): \(error.localizedDescription)") }
            }
        } catch { issues.append("\(root): \(error.localizedDescription)") }
    }
    return ScanSection(category: .apps, state: .partial, reason: "Readable standard application folders, two levels, at most 1,000 entries per folder and 10,000 apps total; symbolic links excluded. Inspect an app to validate its signature. \(issues.prefix(100).joined(separator: "; "))", records: records.sorted { $0.title < $1.title })
}
private func appCandidates(_ item: URL) throws -> [URL] {
    let values = try item.resourceValues(forKeys: [.isDirectoryKey, .isSymbolicLinkKey])
    guard values.isSymbolicLink != true else { throw CollectionError.invalid("Symbolic link excluded") }
    if item.pathExtension == "app" { return [item] }
    guard values.isDirectory == true else { return [] }
    return try children(item.path).filter { $0.pathExtension == "app" }.prefix(1000).filter {
        try $0.resourceValues(forKeys: [.isSymbolicLinkKey]).isSymbolicLink != true
    }
}
private func appRecord(_ url: URL) throws -> Evidence {
    let plist = url.appendingPathComponent("Contents/Info.plist")
    guard let info = try PropertyListSerialization.propertyList(from: boundedData(plist, limit: 1_000_000), format: nil) as? [String: Any] else {
        throw CollectionError.invalid("App Info.plist is not a dictionary")
    }
    let usage = info.keys.filter { $0.hasSuffix("UsageDescription") }.sorted().joined(separator: ", ")
    return Evidence(id: url.path, title: url.deletingPathExtension().lastPathComponent,
                    detail: "Version \(info["CFBundleShortVersionString"] as? String ?? "unknown") · \(info["CFBundleIdentifier"] as? String ?? "unknown identity")",
                    path: url.path, source: "On-disk Info.plist (unverified until signature inspection)", facts: ["signing": "unknown", "sandbox": "unknown", "declaredUsage": usage.isEmpty ? "None declared; this does not establish permission status" : usage])
}
private func startupInventory() throws -> ScanSection {
    let roots = [NSHomeDirectory() + "/Library/LaunchAgents", "/Library/LaunchAgents", "/Library/LaunchDaemons"]
    var records: [Evidence] = [], issues: [String] = []
    for root in roots where FileManager.default.fileExists(atPath: root) {
        do {
            for item in try children(root).filter({ $0.pathExtension == "plist" }).prefix(2000) {
                do { records.append(try decodeLaunchItem(boundedData(item, limit: 1_000_000), path: item.path)) }
                catch { issues.append(error.localizedDescription) }
            }
        } catch { issues.append(error.localizedDescription) }
    }
    return ScanSection(category: .startup, state: .partial, reason: "User and third-party launch plists only; modern login-item registrations and system daemons are outside this snapshot. \(issues.joined(separator: "; "))", records: records)
}
private func hiddenInventory(_ folder: URL?) throws -> ScanSection {
    guard let folder else { return ScanSection(category: .files, state: .partial, reason: "Choose a folder to inspect its immediate hidden entries. No recursive scan or file contents are collected.", records: []) }
    let entries = try children(folder.path)
    let records = try entries.prefix(2000).filter { try $0.resourceValues(forKeys: [.isHiddenKey]).isHidden == true || $0.lastPathComponent.hasPrefix(".") }.map {
        Evidence(id: $0.path, title: $0.lastPathComponent, detail: "Hidden metadata only. Hidden files are common and do not by themselves indicate spyware.", path: $0.path, source: "FileManager metadata", facts: ["hidden": "true"])
    }
    return ScanSection(category: .files, state: .partial, reason: "Immediate entries of \(folder.path); limit 2,000. No traversal through subdirectories or symlinks.", records: records)
}
private func packageInventory() throws -> ScanSection {
    var records: [Evidence] = []
    for root in ["/opt/homebrew/Cellar", "/usr/local/Cellar"] where FileManager.default.fileExists(atPath: root) {
        for package in try children(root).prefix(2000) {
            for version in try children(package.path).prefix(100) {
                do { records.append(try packageReceipt(boundedData(version.appendingPathComponent("INSTALL_RECEIPT.json"), limit: 1_000_000), path: version.path)) }
                catch { records.append(Evidence(id: version.path, title: package.lastPathComponent, detail: "Receipt unavailable: \(error.localizedDescription)", path: version.path, source: "Homebrew Cellar directory", facts: ["version": version.lastPathComponent, "tap": "unknown"])) }
            }
        }
    }
    return ScanSection(category: .packages, state: .partial, reason: "Standard Homebrew Cellar locations only. Custom prefixes, casks, outdated/disabled status and vulnerability matching are not covered.", records: records)
}
private func protectionInventory() -> ScanSection {
    let checks = [("System Integrity Protection", "/usr/bin/csrutil", ["status"]), ("Signed system volume", "/usr/bin/csrutil", ["authenticated-root", "status"]), ("Gatekeeper", "/usr/sbin/spctl", ["--status"])]
    var records: [Evidence] = [], issues: [String] = []
    for (name, path, arguments) in checks {
        do {
            let text = try runCommand(path, arguments)
            records.append(Evidence(id: name, title: name, detail: text.trimmingCharacters(in: .whitespacesAndNewlines), source: path, facts: ["disabled": String(text.lowercased().contains("disabled"))]))
        } catch { issues.append("\(name): \(error.localizedDescription)") }
    }
    return ScanSection(category: .protection, state: .partial, reason: "Protection configuration only; not a forensic verification of all system files. \(issues.joined(separator: "; "))", records: records)
}
#endif
