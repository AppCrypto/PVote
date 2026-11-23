import Foundation

final class RPCStub: URLProtocol {
    override class func canInit(with request: URLRequest) -> Bool { true }
    override class func canonicalRequest(for request: URLRequest) -> URLRequest { request }
    override func stopLoading() {}
    override func startLoading() {
        var bytes = request.httpBody ?? Data()
        if bytes.isEmpty, let stream = request.httpBodyStream {
            stream.open()
            var buffer = [UInt8](repeating: 0, count: 4096)
            while stream.hasBytesAvailable {
                let count = stream.read(&buffer, maxLength: buffer.count)
                if count <= 0 { break }
                bytes.append(contentsOf: buffer.prefix(count))
            }
            stream.close()
        }
        let body = (try? JSONSerialization.jsonObject(with: bytes)) as? [String: Any]
        let method = body?["method"] as? String
        let response: [String: Any]
        if method == "eth_chainId" { response = ["jsonrpc": "2.0", "id": 1, "result": "0x539"] }
        else if method == "eth_getTransactionReceipt" { response = ["jsonrpc": "2.0", "id": 1, "result": NSNull()] }
        else { response = ["jsonrpc": "2.0", "id": 1, "error": ["code": -32000, "message": "test rejection"]] }
        let data = try! JSONSerialization.data(withJSONObject: response)
        client?.urlProtocol(self, didReceive: HTTPURLResponse(url: request.url!, statusCode: 200, httpVersion: nil, headerFields: nil)!, cacheStoragePolicy: .notAllowed)
        client?.urlProtocol(self, didLoad: data)
        client?.urlProtocolDidFinishLoading(self)
    }
}

@main
struct BridgeSmoke {
    static func main() async throws {
        let engine = NativeEngine()
        let version = try await engine.call(["action": "version"])
        precondition((version["sourceSHA256"] as? String)?.count == 64)
        let data = try Data(contentsOf: URL(fileURLWithPath: CommandLine.arguments[1]))
        guard let fixtures = try JSONSerialization.jsonObject(with: data) as? [[String: Any]] else { fatalError("fixtures") }
        let ids = Set(["n4-l2-t3", "n10-l5-t7", "n100-l19-t59"])
        for item in fixtures where ids.contains(item["id"] as? String ?? "") {
            _ = try await engine.call(["action": "prepare", "parameters": item["parameters"]!])
            let p = item["parameters"] as! [String: Any]
            let l = p["l"] as! Int
            for proof in ["context-v1", "legacy-solidity-v1"] {
                let result = try await engine.call(["action": "generate", "scores": Array(repeating: 5, count: l),
                    "protocol": proof, "electionId": "swift-bridge-smoke", "voterIndex": 1, "verify": true])
                precondition((result["shareMs"] as? Double ?? 0) > 0)
                precondition(result["verifyMs"] != nil)
                precondition(result["calldata"] == nil)
            }
        }
        let wallet = try await engine.call(["action": "wallet"])
        for kappa in [1,8] {
            _ = try await engine.call(["action":"prepareComparison","kappa":kappa])
            for scheme in ["PVote","PriScore","3PVSS-bit"] {
                let m = try await engine.call(["action":"generateComparison","kappa":kappa,"scheme":scheme,"electionId":"bridge-comparison","voterIndex":1,"verify":true])
                precondition((m["generationMs"] as? Double ?? 0) > 0 && m["verifyMs"] != nil)
            }
        }
        let address = try await engine.call(["action": "address", "privateKey": wallet["privateKey"]!])
        precondition(address["address"] as? String == wallet["address"] as? String)
        do {
            _ = try await engine.call(["action": "generate", "scores": [999999], "protocol": "context-v1"])
            fatalError("invalid ballot accepted")
        } catch is PVoteError { }
        let topic = "0xaccept"
        let receipt: [String: Any] = ["status": "0x1", "logs": [["address": "0xabc", "topics": [topic, "0x1"]]]]
        precondition(ChainClient.accepted(receipt, contract: "0xABC", topic: topic))
        precondition(!ChainClient.accepted(receipt, contract: "0xDEF", topic: topic))
        precondition(!ChainClient.accepted(["status": "0x0", "logs": receipt["logs"]!], contract: "0xabc", topic: topic))
        precondition(PVoteFootprintBytes() > 0)
        let configuration = URLSessionConfiguration.ephemeral
        configuration.protocolClasses = [RPCStub.self]
        let session = URLSession(configuration: configuration)
        defer { session.invalidateAndCancel() }
        let chain = ChainClient(url: URL(string: "https://unit-test.invalid")!, session: session)
        let id = try ChainClient.quantity(try await chain.rpc("eth_chainId", []))
        precondition(id == 1337)
        let pending = try await chain.receipt("0xhash")
        precondition(pending == nil)
        do { _ = try await chain.rpc("rejected", []); fatalError("RPC error ignored") }
        catch let error as PVoteError { precondition(error.message == "test rejection") }
        print("PASS: macOS Swift/Go bridge, six real RB-PVSS proofs, wallet, JSON-RPC, receipt validation, and memory probe. These are not iPhone measurements.")
    }
}
