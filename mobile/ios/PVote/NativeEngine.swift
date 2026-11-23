import Foundation

struct PVoteError: LocalizedError {
    let message: String
    var errorDescription: String? { message }
    init(_ message: String) { self.message = message }
}

// Actor confines Go calls to one operation at a time, away from the main actor.
actor NativeEngine {
    func exportReport(_ data: Data) throws -> [String: Any] {
        guard let text = String(data: data, encoding: .utf8) else { throw PVoteError("Invalid report encoding") }
        let pointer = text.withCString { PVoteExportReport(UnsafeMutablePointer(mutating: $0)) }
        guard let pointer else { throw PVoteError("Report export failed") }
        defer { PVoteFree(pointer) }
        guard let result = try JSONSerialization.jsonObject(with: Data(String(cString: pointer).utf8)) as? [String: Any],
              result["error"] == nil else { throw PVoteError("Report export failed") }
        return result
    }

    func call(_ request: [String: Any]) throws -> [String: Any] {
        let bytes = try JSONSerialization.data(withJSONObject: request, options: [.sortedKeys])
        guard let input = String(data: bytes, encoding: .utf8) else { throw PVoteError("Request encoding failed") }
        let output = input.withCString { PVoteCall(UnsafeMutablePointer(mutating: $0)) }
        guard let output else { throw PVoteError("Local engine did not respond") }
        defer { PVoteFree(output) }
        let response = Data(String(cString: output).utf8)
        guard let object = try JSONSerialization.jsonObject(with: response) as? [String: Any] else {
            throw PVoteError("Invalid engine response")
        }
        if let error = object["error"] as? String { throw PVoteError(error) }
        return object
    }
}

struct PublicParameters: Codable, Hashable {
    let n: Int
    let t: Int
    let l: Int
    let min: Int
    let max: Int
    let g0: String
    let h0: String
    let g1: String
    let pki: String
    let pk: [String]
    let sigma: [String]
    var object: [String: Any] { (try? JSONSerialization.jsonObject(with: JSONEncoder().encode(self))) as? [String: Any] ?? [:] }
}

struct Fixture: Codable, Identifiable {
    let id: String
    let parameters: PublicParameters
    var label: String { "n=\(parameters.n), l=\(parameters.l), t=\(parameters.t)" }
}

struct Election: Codable {
    let schemaVersion: Int
    let electionId: String
    let `protocol`: String
    let rpcURL: String
    let contract: String
    let chainId: String
    let voterAddress: String
    let voterIndex: UInt64
    let labels: [String]
    let parameters: PublicParameters
}

// Sampling covers the native-call window, including encoding and optional DVerify.
// Sampled peaks are lower bounds; they are not exact peaks of Share alone.
final class FootprintSampler: @unchecked Sendable {
    private let lock = NSLock()
    private var maximum: UInt64 = 0
    private var timer: DispatchSourceTimer?
    let before: UInt64
    init() {
        before = PVoteFootprintBytes()
        maximum = before
        let source = DispatchSource.makeTimerSource(queue: DispatchQueue(label: "PVote.memory", qos: .utility))
        source.schedule(deadline: .now(), repeating: .milliseconds(2))
        source.setEventHandler { [weak self] in self?.sample() }
        timer = source
        source.resume()
    }
    private func sample() {
        let current = PVoteFootprintBytes()
        lock.lock(); maximum = max(maximum, current); lock.unlock()
    }
    func finish() -> [String: Any] {
        timer?.cancel(); timer = nil
        sample()
        lock.lock(); let peak = maximum; lock.unlock()
        return ["nativeCallFootprintBeforeBytes": before, "nativeCallSampledPeakBytes": peak,
                "nativeCallFootprintAfterBytes": PVoteFootprintBytes(), "footprintSamplingIntervalMs": 2]
    }
    deinit { timer?.cancel() }
}
