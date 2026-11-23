import Foundation

struct ChainClient {
    let url: URL
    var session: URLSession = .shared
    func rpc(_ method: String, _ params: [Any]) async throws -> Any {
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.timeoutInterval = 30
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.httpBody = try JSONSerialization.data(withJSONObject: ["jsonrpc": "2.0", "id": 1, "method": method, "params": params])
        let (data, response) = try await session.data(for: request)
        guard let response = response as? HTTPURLResponse, (200..<300).contains(response.statusCode) else { throw PVoteError("RPC HTTP request failed") }
        guard let body = try JSONSerialization.jsonObject(with: data) as? [String: Any] else { throw PVoteError("Invalid RPC response") }
        if let error = body["error"] as? [String: Any] { throw PVoteError(error["message"] as? String ?? "RPC request failed") }
        return body["result"] ?? NSNull()
    }
    static func quantity(_ value: Any) throws -> UInt64 {
        guard let text = value as? String, text.hasPrefix("0x"), let n = UInt64(text.dropFirst(2), radix: 16) else {
            throw PVoteError("Invalid RPC quantity")
        }
        return n
    }
    func receipt(_ hash: String) async throws -> [String: Any]? {
        try await rpc("eth_getTransactionReceipt", [hash]) as? [String: Any]
    }
    static func accepted(_ receipt: [String: Any], contract: String, topic: String) -> Bool {
        guard (try? quantity(receipt["status"] ?? "")) == 1, let logs = receipt["logs"] as? [[String: Any]] else { return false }
        return logs.contains { log in
            guard (log["address"] as? String)?.lowercased() == contract.lowercased(),
                  let topics = log["topics"] as? [String], topics.count == 2 else { return false }
            return topics[0].lowercased() == topic.lowercased()
        }
    }
}
