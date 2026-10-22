import XCTest
import Darwin
@testable import VigilCore

final class IntelligenceTests: XCTestCase {
    func testSocketMetadataPreservesNumericEndpointsWithoutTrustingPID() throws {
        let data = "p42\ncExample\nf12\ntIPv6\nPTCP\nn[::1]:123->[::1]:456\nTST=ESTABLISHED\nf13\ntIPv4\nPUDP\nn*:5353\n"
        let records = try decodeSockets(data)
        XCTAssertEqual(records.count, 2)
        XCTAssertEqual(records[0].facts["endpoint"], "[::1]:123->[::1]:456")
        XCTAssertEqual(records[0].facts["attribution"], "unverified PID snapshot")
        XCTAssertEqual(records[1].assessment, "information")
        XCTAssertThrowsError(try decodeSockets("f1\nn1.2.3.4:80\n"))
    }
    func testLiveSocketCollectorSeesTestListener() throws {
        let descriptor = try listener()
        defer { close(descriptor) }
        let section = try networkSnapshot(pid: ProcessInfo.processInfo.processIdentifier)
        XCTAssertEqual(section.state, .partial)
        XCTAssertTrue(section.records.contains { $0.facts["pid"] == String(ProcessInfo.processInfo.processIdentifier) && $0.facts["state"] == "LISTEN" })
    }
    private func listener() throws -> Int32 {
        let descriptor = socket(AF_INET, SOCK_STREAM, 0)
        guard descriptor >= 0 else { throw CollectionError.invalid("Cannot create test socket") }
        var address = sockaddr_in()
        address.sin_len = UInt8(MemoryLayout<sockaddr_in>.size); address.sin_family = sa_family_t(AF_INET)
        address.sin_addr.s_addr = inet_addr("127.0.0.1")
        let bound = withUnsafePointer(to: &address) { pointer in
            pointer.withMemoryRebound(to: sockaddr.self, capacity: 1) { Darwin.bind(descriptor, $0, socklen_t(MemoryLayout<sockaddr_in>.size)) }
        }
        guard bound == 0, listen(descriptor, 1) == 0 else { close(descriptor); throw CollectionError.invalid("Cannot bind test listener") }
        return descriptor
    }
    func testUnavailableFeedRetainsCacheAsStale() {
        let evidence = Evidence(id: "CVE-test", title: "Test", detail: "Research", source: "fixture")
        let cached = ScanSection(category: .advisories, state: .partial, reason: "Prior catalog", records: [evidence])
        let failed = ScanSection(category: .advisories, state: .failed, reason: "Offline", records: [])
        let retained = retainAdvisories(failed, previous: cached)
        XCTAssertEqual(retained.state, .stale)
        XCTAssertEqual(retained.records.first?.id, "CVE-test")
        XCTAssertTrue(retained.reason.contains("Offline"))
    }
    func testAdvisoryCatalogIsResearchNotInstalledVersionMatch() throws {
        let json = Data(#"{"catalogVersion":"2026.10.22","dateReleased":"2026-10-22T00:00:00Z","count":1,"vulnerabilities":[{"cveID":"CVE-2026-1234","vendorProject":"Example","product":"Tool","vulnerabilityName":"Example issue","shortDescription":"Description","requiredAction":"Update as directed by vendor","dateAdded":"2026-10-22"}]}"#.utf8)
        let section = try decodeKEV(json)
        XCTAssertEqual(section.records.count, 1)
        XCTAssertEqual(section.records[0].facts["applicability"], "not assessed")
        XCTAssertEqual(section.records[0].assessment, "information")
        XCTAssertEqual(section.state, .partial)
        XCTAssertThrowsError(try decodeKEV(Data("{}".utf8)))
        XCTAssertThrowsError(try decodeKEV(Data(String(decoding: json, as: UTF8.self).replacingOccurrences(of: "\"count\":1", with: "\"count\":2").utf8)))
    }
}
