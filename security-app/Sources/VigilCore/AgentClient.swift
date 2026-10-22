import Foundation

// Implements [AGENT-VENDOR], [AGENT-CHOICE], [AGENT-PRIVACY]: end-user-only, explicit one-turn disclosure.
public struct AgentConfiguration: Codable, Identifiable, Hashable {
    public var id: UUID { configID }
    public let name: String
    public let provider: String
    public let model: String
    public let configID: UUID
    public let contractHash: String
    public init(name: String, provider: String, model: String, configID: UUID, contractHash: String) {
        self.name = name; self.provider = provider; self.model = model; self.configID = configID; self.contractHash = contractHash
    }
}
public func agentRequest(configuration: AgentConfiguration, token: String, message: String, consent: Bool) throws -> URLRequest {
    guard consent else { throw CollectionError.invalid("AI disclosure was not approved") }
    guard !message.trimmingCharacters(in: .whitespacesAndNewlines).isEmpty, message.utf8.count <= 16_000 else { throw CollectionError.invalid("Enter a question up to 16 KB") }
    guard token.split(separator: ".").count == 3, token.utf8.count <= 16_000,
          token.unicodeScalars.allSatisfy({ CharacterSet(charactersIn: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_.").contains($0) }) else {
        throw CollectionError.invalid("Use a valid end-user JWT; operator API keys are not supported")
    }
    guard configuration.contractHash.count == 64, configuration.contractHash.allSatisfy({ $0.isHexDigit }),
          let url = URL(string: "https://api.nimblesite.ai/api/v1/end-user/chat/\(configuration.configID.uuidString.lowercased())") else {
        throw CollectionError.invalid("Invalid deployed configuration contract hash")
    }
    var request = URLRequest(url: url, cachePolicy: .reloadIgnoringLocalCacheData, timeoutInterval: 60)
    request.httpMethod = "POST"; request.setValue("application/json", forHTTPHeaderField: "Content-Type")
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
    request.httpBody = try JSONEncoder().encode(["message": message, "tool_manifest_sha256": configuration.contractHash])
    return request
}
public func decodeAgentReply(_ data: Data) throws -> String {
    guard data.count <= 1_000_000, let reply = try JSONSerialization.jsonObject(with: data) as? [String: Any],
          let calls = reply["tool_calls"] as? [Any], calls.isEmpty,
          let id = reply["conversation_id"] as? String, UUID(uuidString: id) != nil,
          let text = reply["response"] as? String else {
        throw CollectionError.invalid("Unsupported agent response or tool request. No local tool was executed.")
    }
    return text
}
public func askAgent(_ request: URLRequest) async throws -> String {
    guard request.url?.host == "api.nimblesite.ai", request.url?.scheme == "https" else { throw CollectionError.invalid("Invalid agent destination") }
    return try decodeAgentReply(await boundedResponse(request, limit: 1_000_000))
}
