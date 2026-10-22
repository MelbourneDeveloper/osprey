import Foundation

// Implements [SECURITY-COLLECTORS-PERMISSIONS]: guidance is not observed authorization.
public func permissionGuidance() -> ScanSection {
    let reviews = [
        ("Accessibility", "Assistive tools may control other apps. Review unfamiliar entries; legitimate accessibility tools need this access."),
        ("Input Monitoring", "Review apps permitted to monitor keyboard input. Permission alone does not establish keylogging."),
        ("Screen & System Audio Recording", "Review recording tools against the work you expect them to perform."),
        ("Full Disk Access", "Review broad file access for tools you no longer use. Backup and security tools may require it."),
        ("Camera", "Review camera access in Privacy & Security. A usage-description string does not prove access was granted."),
        ("Microphone", "Review microphone access. Legitimate calling and recording apps commonly require it.")
    ]
    return ScanSection(category: .permissions, state: .unsupported,
        reason: "Guided review only; other apps’ actual grants and sensor use are not enumerated. Open System Settings → Privacy & Security.",
        records: reviews.map { Evidence(id: $0.0, title: $0.0, detail: $0.1, source: "Apple Mac User Guide: Privacy & Security", facts: ["grant": "unknown", "kind": "guidance"]) })
}
