import Foundation

public struct RunningAppFact: Sendable {
    public let pid: Int32
    public let launched: Date?
    public let active: Bool
    public let hidden: Bool
    public let policy: String
    public init(pid: Int32, launched: Date?, active: Bool, hidden: Bool, policy: String) {
        self.pid = pid; self.launched = launched; self.active = active; self.hidden = hidden; self.policy = policy
    }
}
public func processStartDate(_ value: String) -> Date? {
    let formatter = DateFormatter()
    formatter.locale = Locale(identifier: "en_US_POSIX")
    formatter.dateFormat = "EEE MMM d HH:mm:ss yyyy"
    return formatter.date(from: value)
}
public func annotateRunningApps(_ records: [Evidence], applications: [RunningAppFact]) -> [Evidence] {
    records.map { record in
        guard let pidText = record.facts["pid"], let pid = Int32(pidText),
              let app = applications.first(where: { $0.pid == pid }), let launch = app.launched,
              let started = record.facts["started"].flatMap(processStartDate), abs(launch.timeIntervalSince(started)) < 1 else { return record }
        var result = record
        result.facts["applicationState"] = app.active ? "Frontmost at observation time" : "Background at observation time"
        result.facts["hiddenWindows"] = String(app.hidden)
        result.facts["activationPolicy"] = app.policy
        result.source += "; NSWorkspace application state matched by PID and launch time"
        return result
    }
}
