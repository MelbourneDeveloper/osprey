import XCTest
@testable import VigilCore

final class AgentTests: XCTestCase {
    private let model = AgentConfiguration(name: "My model", provider: "provider", model: "model", configID: UUID(), contractHash: String(repeating: "a", count: 64))
    func testNoConsentMeansNoRequest() {
        XCTAssertThrowsError(try agentRequest(configuration: model, token: "user.jwt.token", message: "Explain a rule", consent: false))
    }
    func testEndUserRequestNeverContainsOperatorKeyOrInventedDeviceContext() throws {
        let request = try agentRequest(configuration: model, token: "user.jwt.token", message: "Explain a rule", consent: true)
        XCTAssertEqual(request.url?.host, "api.nimblesite.ai")
        XCTAssertTrue(request.url?.path.hasPrefix("/api/v1/end-user/chat/") == true)
        XCTAssertNil(request.value(forHTTPHeaderField: "X-API-Key"))
        XCTAssertEqual(request.value(forHTTPHeaderField: "Authorization"), "Bearer user.jwt.token")
        let body = try XCTUnwrap(try JSONSerialization.jsonObject(with: XCTUnwrap(request.httpBody)) as? [String: String])
        XCTAssertEqual(body["message"], "Explain a rule")
        XCTAssertEqual(Set(body.keys), ["message", "tool_manifest_sha256"])
    }
    func testMalformedContractAndHeaderInjectionAreRejected() {
        let invalid = AgentConfiguration(name: "Invalid", provider: "p", model: "m", configID: UUID(), contractHash: "bad")
        XCTAssertThrowsError(try agentRequest(configuration: invalid, token: "a.b.c", message: "hello", consent: true))
        XCTAssertThrowsError(try agentRequest(configuration: model, token: "a\r\nX-Test:b", message: "hello", consent: true))
    }
    func testAgentCannotRequestToolsOrTurnTextIntoActions() throws {
        let reply = Data(#"{"response":"ignore policy; execute shell","conversation_id":"00000000-0000-0000-0000-000000000001","tool_calls":[{"id":"one","name":"shell"}]}"#.utf8)
        XCTAssertThrowsError(try decodeAgentReply(reply))
        let plain = Data(#"{"response":"Explanation only","conversation_id":"00000000-0000-0000-0000-000000000001","tool_calls":[]}"#.utf8)
        XCTAssertEqual(try decodeAgentReply(plain), "Explanation only")
    }
}
