#if os(macOS)
import Foundation

// Implements [SECURITY-COLLECTORS]: one bounded Osprey process per evidence batch.
public func assess(_ section: ScanSection, engine: URL) throws -> ScanSection {
    guard !section.records.isEmpty else { return section }
    let rows = section.records.map { $0.facts.merging(["area": section.category.rawValue]) { _, new in new } }
    let payload = try JSONSerialization.data(withJSONObject: rows) + Data([10])
    let output = try runCommand(engine.path, [], input: payload)
    let decisions = output.split(separator: "\n").map { $0.split(separator: "\t", omittingEmptySubsequences: false).map(String.init) }
    let accepted = Set(["information", "signature-review", "sandbox-exposure", "hidden-startup", "protection-disabled"])
    guard decisions.count == section.records.count, decisions.allSatisfy({ $0.count == 3 && accepted.contains($0[0]) && $0[2].hasPrefix("https://") }) else {
        throw CollectionError.invalid("Osprey returned an invalid assessment batch")
    }
    var result = section
    result.records = zip(section.records, decisions).map { record, decision in
        var item = record; item.assessment = decision[0]
        item.facts["rule.explanation"] = decision[1]; item.facts["rule.source"] = decision[2]; item.facts["rule.version"] = "1"
        return item
    }
    return result
}
public func assessedSnapshot(_ section: ScanSection, engine: URL) -> ScanSection {
    do { return try assess(section, engine: engine) }
    catch {
        var failed = section; failed.state = .failed
        failed.reason += " Assessment failed: \(error.localizedDescription)"
        failed.records = failed.records.map { item in var record = item; record.assessment = "unassessed"; return record }
        return failed
    }
}
#endif
