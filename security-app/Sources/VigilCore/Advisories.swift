import Foundation

private struct KEVCatalog: Decodable {
    let catalogVersion: String
    let dateReleased: String
    let count: Int
    let vulnerabilities: [KEVEntry]
}
private struct KEVEntry: Decodable {
    let cveID: String
    let vendorProject: String
    let product: String
    let vulnerabilityName: String
    let shortDescription: String
    let requiredAction: String
    let dateAdded: String
}
// Implements [SECURITY-COLLECTORS-ADVISORIES]: an imported primary-source catalog is never a device match.
public func decodeKEV(_ data: Data) throws -> ScanSection {
    guard data.count <= 8_000_000 else { throw CollectionError.invalid("Advisory catalog exceeds 8 MB") }
    let catalog = try JSONDecoder().decode(KEVCatalog.self, from: data)
    guard catalog.count == catalog.vulnerabilities.count, catalog.count <= 20_000,
          Set(catalog.vulnerabilities.map(\.cveID)).count == catalog.count else { throw CollectionError.invalid("Invalid advisory catalog count or identities") }
    return ScanSection(category: .advisories, state: .partial,
        reason: "Imported CISA KEV-format catalog \(catalog.catalogVersion), released \(catalog.dateReleased). File provenance is unverified. Search for research; installed versions and applicability are NOT matched. No matching record is not proof of safety.",
        records: catalog.vulnerabilities.map { advisoryEvidence($0, version: catalog.catalogVersion) })
}
private func advisoryEvidence(_ entry: KEVEntry, version: String) -> Evidence {
    Evidence(id: entry.cveID, title: "\(entry.cveID) · \(entry.vendorProject) \(entry.product)",
        detail: "\(entry.vulnerabilityName). \(entry.shortDescription)", source: "Imported CISA KEV-format file (unverified origin)",
        facts: ["applicability": "not assessed", "requiredAction": entry.requiredAction, "added": entry.dateAdded,
                "catalogVersion": version, "reference": "https://www.cisa.gov/known-exploited-vulnerabilities-catalog"])
}

public func fetchKEV() async throws -> ScanSection {
    guard let url = URL(string: "https://raw.githubusercontent.com/cisagov/kev-data/develop/known_exploited_vulnerabilities.json") else {
        throw CollectionError.invalid("Invalid CISA catalog endpoint")
    }
    var section = try decodeKEV(await boundedResponse(URLRequest(url: url, timeoutInterval: 30), limit: 8_000_000))
    section.reason = section.reason.replacingOccurrences(of: "File provenance is unverified.", with: "Retrieved over HTTPS from CISA's official GitHub mirror; no independent artifact signature or installed-version mapping.")
    section.records = section.records.map { record in
        var value = record; value.source = url.absoluteString; return value
    }
    return section
}

public func retainAdvisories(_ section: ScanSection, previous: ScanSection?) -> ScanSection {
    guard section.category == .advisories, section.records.isEmpty, var cached = previous, !cached.records.isEmpty else { return section }
    cached.state = .stale; cached.reason = "Cached catalog; refresh explicitly. Current collection: " + section.reason + " " + cached.reason
    return cached
}
