import Foundation

private final class NoRedirect: NSObject, URLSessionTaskDelegate, @unchecked Sendable {
    func urlSession(_ session: URLSession, task: URLSessionTask, willPerformHTTPRedirection response: HTTPURLResponse,
                    newRequest request: URLRequest, completionHandler: @escaping (URLRequest?) -> Void) {
        completionHandler(nil)
    }
}
// No redirects, persistent cookies, disk cache, automatic application retries or unbounded bodies.
func boundedResponse(_ request: URLRequest, limit: Int) async throws -> Data {
    let configuration = URLSessionConfiguration.ephemeral
    configuration.httpCookieStorage = nil; configuration.urlCache = nil; configuration.timeoutIntervalForResource = 90
    let session = URLSession(configuration: configuration, delegate: NoRedirect(), delegateQueue: nil)
    defer { session.invalidateAndCancel() }
    let (bytes, response) = try await session.bytes(for: request)
    guard let response = response as? HTTPURLResponse, response.statusCode == 200 else {
        throw CollectionError.invalid("Service returned HTTP \((response as? HTTPURLResponse)?.statusCode ?? 0). No automatic retry; check connectivity, authentication and service status.")
    }
    var data = Data()
    for try await byte in bytes {
        try Task.checkCancellation()
        guard data.count < limit else { throw CollectionError.invalid("Response exceeded size limit") }
        data.append(byte)
    }
    return data
}
