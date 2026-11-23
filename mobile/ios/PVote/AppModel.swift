import SwiftUI
import UIKit
import CryptoKit

@MainActor
final class AppModel: ObservableObject {
    @Published var fixtures: [Fixture] = []
    @Published var selectedFixtureID = "n4-l2-t3"
    @Published var scores: [Int] = [5, 5]
    @Published var election: Election?
    @Published var address = ""
    @Published var status = "Loading local engine…"
    @Published var busy = false
    @Published var progress: Double = 0
    @Published var lastMeasurement: [String: Any]?
    @Published var benchmarkRows: [[String: Any]] = []
    @Published var exportURLs: [URL] = []
    @Published var pending: [String: Any]?
    @Published var repetitions = 30
    @Published var warmups = 3
    @Published var benchmarkSuite = "paper3"
    @Published var proofProtocol = "context-v1"
    @Published var verifyEach = false
    private let engine = NativeEngine()
    private var task: Task<Void, Never>?
    private var metadata: [String: Any] = [:]
    private var session: [String: Any] = [:]
    private var runID = ""
    private var votingCalldata: String?
    private var acceptanceTopic: String?
    private var activeBenchmark = false
    private var loaded = false
    private var benchmarkOrder = "ascending"

    var fixture: Fixture? { fixtures.first { $0.id == selectedFixtureID } }
    var parameters: PublicParameters? { election?.parameters ?? fixture?.parameters }
    var canSend: Bool { election != nil && votingCalldata != nil && pending == nil && !busy }
    private var documents: URL { FileManager.default.urls(for: .documentDirectory, in: .userDomainMask)[0] }

    func start() async {
        guard !loaded else { return }; loaded = true
        do {
            guard let url = Bundle.main.url(forResource: "fixtures", withExtension: "json") else { throw PVoteError("Public fixtures missing") }
            fixtures = try JSONDecoder().decode([Fixture].self, from: Data(contentsOf: url))
            metadata = try await engine.call(["action": "version"])
            metadata.merge(deviceMetadata()) { _, new in new }
            UIDevice.current.isBatteryMonitoringEnabled = true
            if let argument = ProcessInfo.processInfo.arguments.first(where: { $0.hasPrefix("--benchmark-comparison=") }) {
                await runScoreComparison(String(argument.dropFirst("--benchmark-comparison=".count)))
                return
            }
            if let argument = ProcessInfo.processInfo.arguments.first(where: { $0.hasPrefix("--benchmark-recheck=") }) {
                try await recheckCurve(String(argument.dropFirst("--benchmark-recheck=".count)))
                return
            }
            if let argument = ProcessInfo.processInfo.arguments.first(where: { $0.hasPrefix("--export-http=") }) {
                let endpoint = String(argument.dropFirst("--export-http=".count))
                try await exportLatestReportHTTP(endpoint)
                return
            }
            if ProcessInfo.processInfo.arguments.contains("--dump-latest-report") {
                let reports = try FileManager.default.contentsOfDirectory(at: documents, includingPropertiesForKeys: nil)
                    .filter { $0.lastPathComponent.hasPrefix("benchmark-") && $0.pathExtension == "json" }
                    .sorted { $0.lastPathComponent > $1.lastPathComponent }
                guard let file = reports.first else { throw PVoteError("No device report") }
                let encoded = try await engine.exportReport(Data(contentsOf: file))
                guard let digest = encoded["sha256"] as? String, let payload = encoded["gzipBase64"] as? String else { throw PVoteError("Report export failed") }
                print("PVOTE_EXPORT_BEGIN \(digest) \(payload.count)")
                for offset in stride(from: 0, to: payload.count, by: 512) {
                    let start = payload.index(payload.startIndex, offsetBy: offset)
                    let end = payload.index(start, offsetBy: min(512, payload.count-offset))
                    print("PVOTE_DATA \(payload[start..<end])")
                }
                print("PVOTE_EXPORT_END")
                fflush(stdout)
                status = "Existing report exported"
                return
            }
            if let key = try Wallet.load() {
                address = try await engine.call(["action": "address", "privateKey": key])["address"] as? String ?? ""
            } else {
                let wallet = try await engine.call(["action": "wallet"])
                guard let key = wallet["privateKey"] as? String else { throw PVoteError("Account creation failed") }
                try Wallet.save(key)
                address = wallet["address"] as? String ?? ""
            }
            if let data = try? Data(contentsOf: documents.appendingPathComponent("election.json")) {
                try installElection(data)
            }
            if let data = try? Data(contentsOf: documents.appendingPathComponent("pending.json")) {
                pending = try JSONSerialization.jsonObject(with: data) as? [String: Any]
            }
            refreshExports()
            status = "Ready for offline voting and benchmarks"
            // Explicit opt-in for reproducible device runs launched from Xcode.
            if ProcessInfo.processInfo.arguments.contains("--benchmark-curve") {
                proofProtocol = "context-v1"
                benchmarkSuite = "curve"
                runBenchmark()
            }
        } catch { status = error.localizedDescription }
    }

    func resetScores() {
        guard let p = parameters else { return }
        scores = Array(repeating: (p.min + p.max) / 2, count: p.l)
        lastMeasurement = nil; votingCalldata = nil
    }

    private func recheckCurve(_ destination: String) async throws {
        guard var components = URLComponents(string: destination),
              let token = components.fragment, token.count >= 16 else { throw PVoteError("Recheck export token missing") }
        components.fragment = nil
        guard let endpoint = components.url?.absoluteString else { throw PVoteError("Invalid recheck export URL") }
        proofProtocol = "context-v1"; benchmarkSuite = "curve"
        repetitions = 30; warmups = 3; verifyEach = false
        UIApplication.shared.isIdleTimerDisabled = true
        defer { UIApplication.shared.isIdleTimerDisabled = false }
        for (round, order) in ["ascending", "descending", "rotated"].enumerated() {
            try Task.checkCancellation()
            if round > 0 {
                status = "Cooling for 60 s between runs…"
                try await Task.sleep(for: .seconds(60))
            }
            var coolingSeconds = 0
            while ProcessInfo.processInfo.thermalState != .nominal {
                guard coolingSeconds < 600 else { throw PVoteError("Cooling timed out before the next run") }
                status = "Cooling before round \(round+1) (\(coolingSeconds) s)…"
                if coolingSeconds % 30 == 0 {
                    var progressURL = components
                    progressURL.path = "/"; progressURL.query = "phase=cooling&round=\(round+1)&seconds=\(coolingSeconds)&thermal=\(ProcessInfo.processInfo.thermalState.rawValue)"
                    if let url = progressURL.url {
                        var request = URLRequest(url: url); request.timeoutInterval = 5
                        _ = try? await URLSession.shared.data(for: request)
                    }
                }
                try await Task.sleep(for: .seconds(5))
                coolingSeconds += 5
            }
            guard UIApplication.shared.applicationState == .active else { throw PVoteError("Recheck stopped in background") }
            benchmarkOrder = order
            runBenchmark()
            session["recheckRound"] = round + 1
            session["preRunNominalThermalRequired"] = true
            session["interRoundRestSeconds"] = 60
            session["additionalCoolingWaitSeconds"] = coolingSeconds
            await task?.value
            UIApplication.shared.isIdleTimerDisabled = true
            guard session["status"] as? String == "complete" else {
                throw PVoteError("Recheck interrupted: \(session["interruptionReason"] as? String ?? "unknown reason"). Raw report saved")
            }
            try await exportLatestReportHTTP(endpoint, token: token)
        }
        benchmarkOrder = "ascending"
        status = "Three runs exported: 2,070 measurements"
    }

    private func runScoreComparison(_ destination: String) async {
        guard !busy else { return }
        busy = true; activeBenchmark = true; UIApplication.shared.isIdleTimerDisabled = true
        task = Task {
            defer { busy = false; activeBenchmark = false; UIApplication.shared.isIdleTimerDisabled = false; refreshExports() }
            do { try await scoreComparisonSuite(destination) }
            catch {
                if !runID.isEmpty && session["status"] as? String != "complete" {
                    session["status"] = "interrupted"
                    session["interruptionReason"] = error.localizedDescription
                    try? persistBenchmark()
                }
                status = "Comparison stopped: \(error.localizedDescription). Data saved"
            }
        }
        await task?.value
    }

    private func comparisonCooling(_ endpoint: URLComponents, label: String) async throws {
        var waited = 0
        while ProcessInfo.processInfo.thermalState != .nominal {
            try Task.checkCancellation()
            guard waited < 600 else { throw PVoteError("Comparison cooling timed out") }
            status = "Cooling · \(label) · \(waited) s"
            if waited % 30 == 0 {
                var progressURL = endpoint; progressURL.path = "/"
                progressURL.queryItems = [URLQueryItem(name:"phase",value:"comparison-cooling"),URLQueryItem(name:"case",value:label),URLQueryItem(name:"seconds",value:String(waited))]
                if let url = progressURL.url { var r = URLRequest(url:url); r.timeoutInterval = 5; _ = try? await URLSession.shared.data(for:r) }
            }
            try await Task.sleep(for:.seconds(5)); waited += 5
        }
        guard UIApplication.shared.applicationState == .active else { throw PVoteError("App is in background") }
    }

    private func scoreComparisonSuite(_ destination: String) async throws {
        guard var endpoint = URLComponents(string:destination), let token = endpoint.fragment, token.count >= 16 else { throw PVoteError("Comparison export token missing") }
        endpoint.fragment = nil
        guard let url = endpoint.url?.absoluteString else { throw PVoteError("Invalid comparison export URL") }
        var setupRecords: [[String:Any]] = []
        for kappa in 1...8 {
            status = "Preparing comparison parameters · κ=\(kappa) (untimed)"
            setupRecords.append(try await engine.call(["action":"prepareComparison","kappa":kappa]))
        }
        let schemes = ["PVote","PriScore","3PVSS-bit"]
        let ascending = (1...8).flatMap { k in schemes.map { (k,$0) } }
        for round in 0..<3 {
            if round > 0 { status = "Cooling for 60 s between runs"; try await Task.sleep(for:.seconds(60)) }
            let cases = round == 1 ? Array(ascending.reversed()) : round == 2 ? Array(ascending.dropFirst(12))+Array(ascending.prefix(12)) : ascending
            runID = ISO8601DateFormatter().string(from:Date()).replacingOccurrences(of:":",with:"-")+"-cmp-"+UUID().uuidString.prefix(8)
            benchmarkRows = []; progress = 0
            session = ["schemaVersion":1,"runId":runID,"metadata":metadata,"experiment":"score_domain_generation",
                       "runsPerConfiguration":10,"warmupsPerConfiguration":3,"recheckRound":round+1,
                       "configurationOrderPolicy":["ascending","descending","rotated"][round],
                       "fixtureIds":cases.map { "\($0.1)-k\($0.0)" },"setupRecords":setupRecords,
                       "scorePolicy":"same cyclic distribution of B=2^kappa-1 across five coordinates for all schemes",
                       "setupAndParameterLoadingTimed":false,"verifyEachMeasuredBallot":false,
                       "interRoundRestSeconds":60,"preCaseNominalThermalRequired":true,
                       "startedAt":ISO8601DateFormatter().string(from:Date()),"status":"running"]
            var excluded: [[String:Any]] = []
            var verified: [String] = []
            for (caseIndex,item) in cases.enumerated() {
                let (kappa,scheme) = item; let id = "\(scheme)-k\(kappa)"
                try await comparisonCooling(endpoint,label:"round \(round+1) \(id)")
                let base: [String:Any] = ["action":"generateComparison","scheme":scheme,"kappa":kappa,"electionId":"comparison-\(runID)","voterIndex":1]
                for warm in 0..<3 {
                    try Task.checkCancellation(); status = "Warmup · round \(round+1) \(id) · \(warm+1)/3"
                    var request = base; request["verify"] = warm == 0
                    _ = try await engine.call(request)
                }
                verified.append(id); session["verifiedWarmupCases"] = verified
                var sample = 1
                while sample <= 10 {
                    try Task.checkCancellation()
                    if ProcessInfo.processInfo.thermalState == .serious || ProcessInfo.processInfo.thermalState == .critical {
                        try await comparisonCooling(endpoint,label:id)
                    }
                    guard !ProcessInfo.processInfo.isLowPowerModeEnabled else { throw PVoteError("Disable Low Power Mode before measuring") }
                    var request = base; request["verify"] = false
                    status = "Measuring · round \(round+1) \(id) · \(sample)/10"
                    var row = try await measuredCall(request)
                    row["fixtureId"] = id; row["sampleIndex"] = sample
                    let after = row["environmentAfter"] as? [String:Any] ?? [:]
                    if ["serious","critical"].contains(after["thermalState"] as? String ?? "") || after["lowPowerMode"] as? Bool == true {
                        row["exclusionReason"] = "environment outside predefined thermal/power limits"
                        excluded.append(row); session["excludedSamples"] = excluded; try persistBenchmark()
                        try await comparisonCooling(endpoint,label:id); continue
                    }
                    benchmarkRows.append(row)
                    progress = Double(caseIndex*10+sample)/240
                    try persistBenchmark(); sample += 1
                }
            }
            session["status"] = "complete"; session["completedAt"] = ISO8601DateFormatter().string(from:Date())
            try persistBenchmark(); try await exportLatestReportHTTP(url,token:token)
        }
        status = "Comparison exported: 720 measurements"
    }

    private func exportLatestReportHTTP(_ endpoint: String, token: String? = nil) async throws {
        guard var components = URLComponents(string: endpoint), ["http", "https"].contains(components.scheme ?? "") else { throw PVoteError("Invalid HTTP export URL") }
        let authorization = token ?? components.fragment
        components.fragment = nil
        guard let url = components.url else { throw PVoteError("Invalid HTTP export URL") }
        let reports = try FileManager.default.contentsOfDirectory(at: documents, includingPropertiesForKeys: nil)
            .filter { $0.lastPathComponent.hasPrefix("benchmark-") && $0.pathExtension == "json" }
            .sorted { $0.lastPathComponent > $1.lastPathComponent }
        guard let file = reports.first else { throw PVoteError("No benchmark report") }
        let bytes = try Data(contentsOf: file)
        let digest = SHA256.hash(data: bytes).map { String(format: "%02x", $0) }.joined()
        var request = URLRequest(url: url)
        request.httpMethod = "POST"; request.timeoutInterval = 30
        request.setValue("application/json", forHTTPHeaderField: "Content-Type")
        request.setValue(digest, forHTTPHeaderField: "X-PVote-SHA256")
        if let authorization { request.setValue("Bearer \(authorization)", forHTTPHeaderField: "Authorization") }
        status = "Exporting completed report over HTTP…"
        let (response, info) = try await URLSession.shared.upload(for: request, from: bytes)
        guard let info = info as? HTTPURLResponse, info.statusCode == 200,
              let ack = try JSONSerialization.jsonObject(with: response) as? [String: Any],
              ack["sha256"] as? String == digest else { throw PVoteError("HTTP export checksum mismatch") }
        print("PVOTE_HTTP_EXPORTED \(digest) \(bytes.count)")
        status = "HTTP export verified"
    }

    func invalidateBallot() { lastMeasurement = nil; votingCalldata = nil }

    func importElection(_ url: URL) {
        guard !busy else { return }
        let accessed = url.startAccessingSecurityScopedResource()
        defer { if accessed { url.stopAccessingSecurityScopedResource() } }
        do {
            let data = try Data(contentsOf: url)
            try installElection(data)
            try data.write(to: documents.appendingPathComponent("election.json"), options: .atomic)
            status = "Election imported. Generation and signing stay on device"
        } catch { status = error.localizedDescription }
    }

    private func installElection(_ data: Data) throws {
        guard data.count <= 1_048_576 else { throw PVoteError("Election file too large") }
        let item = try JSONDecoder().decode(Election.self, from: data)
        guard item.schemaVersion == 1, item.protocol == "legacy-solidity-v1", !item.electionId.isEmpty,
              item.voterIndex > 0, (UInt64(item.chainId) ?? 0) > 0,
              item.contract.range(of: "^0x[0-9a-fA-F]{40}$", options: .regularExpression) != nil,
              item.voterAddress.range(of: "^0x[0-9a-fA-F]{40}$", options: .regularExpression) != nil,
              let url = URL(string: item.rpcURL), ["http", "https"].contains(url.scheme ?? ""), url.host != nil,
              item.parameters.n > 0, item.parameters.n <= 200, item.parameters.l > 0,
              item.parameters.l <= 64, item.parameters.l < item.parameters.t, item.parameters.t <= item.parameters.n,
              item.parameters.min >= -1_000_000, item.parameters.max <= 1_000_000,
              item.parameters.min <= item.parameters.max,
              item.parameters.max - item.parameters.min <= 255 else { throw PVoteError("Invalid election file") }
        election = item; resetScores()
    }

    func useOffline() {
        guard !busy else { return }
        election = nil
        try? FileManager.default.removeItem(at: documents.appendingPathComponent("election.json"))
        resetScores(); status = "Offline mode"
    }

    func generateBallot() {
        guard !busy, let p = parameters else { return }
        busy = true
        task = Task {
            defer { busy = false }
            do {
                _ = try await engine.call(["action": "prepare", "parameters": p.object])
                let request: [String: Any] = ["action": "generate", "scores": scores,
                    "protocol": election?.protocol ?? proofProtocol, "electionId": election?.electionId ?? "offline-demo",
                    "voterIndex": election?.voterIndex ?? 1, "includeCalldata": true, "verify": true]
                var result = try await measuredCall(request)
                votingCalldata = result.removeValue(forKey: "calldata") as? String
                acceptanceTopic = result["acceptanceTopic"] as? String
                lastMeasurement = result
                status = "Ballot generated and verified"
            } catch { status = error.localizedDescription }
        }
    }

    private func measuredCall(_ request: [String: Any]) async throws -> [String: Any] {
        let sampler = FootprintSampler()
        let start = ProcessInfo.processInfo.systemUptime
        let before = environment()
        do {
            var result = try await engine.call(request)
            result.merge(sampler.finish()) { _, new in new }
            result["nativeCallWallMs"] = (ProcessInfo.processInfo.systemUptime - start) * 1000
            result["environmentBefore"] = before
            result["environmentAfter"] = environment()
            result["timestamp"] = ISO8601DateFormatter().string(from: Date())
            return result
        } catch { _ = sampler.finish(); throw error }
    }

    func runBenchmark() {
        guard !busy else { return }
        let paperIDs: Set<String> = ["n10-l5-t7", "n50-l10-t30", "n100-l19-t59"]
        let ordered = fixtures.filter { item in
            if benchmarkSuite == "paper3" { return paperIDs.contains(item.id) }
            if benchmarkSuite == "curve" { return item.parameters.n % 10 == 0 && [1,10,19].contains(item.parameters.l) && item.parameters.l <= item.parameters.n / 3 }
            return item.id == selectedFixtureID
        }.sorted { ($0.parameters.l, $0.parameters.n) < ($1.parameters.l, $1.parameters.n) }
        let selected: [Fixture]
        if benchmarkOrder == "descending" { selected = Array(ordered.reversed()) }
        else if benchmarkOrder == "rotated" {
            let offset = ordered.count / 2
            selected = Array(ordered.dropFirst(offset)) + Array(ordered.prefix(offset))
        } else { selected = ordered }
        guard !selected.isEmpty else { return }
        busy = true; activeBenchmark = true; UIApplication.shared.isIdleTimerDisabled = true
        benchmarkRows = []; progress = 0
        runID = ISO8601DateFormatter().string(from: Date()).replacingOccurrences(of: ":", with: "-") + "-" + UUID().uuidString.prefix(8)
        let runs = repetitions, warm = warmups, protocolName = proofProtocol, verify = verifyEach
        session = ["schemaVersion": 1, "runId": runID, "metadata": metadata,
                   "warmupsPerConfiguration": warm, "runsPerConfiguration": runs,
                   "protocol": protocolName, "verifyEachMeasuredBallot": verify,
                   "setupAndParameterLoadingTimed": false,
                   "footprintScope": "whole app during native call; sampled lower bound, not Share-only",
                   "scorePolicy": "all coordinates are interval midpoint; fresh cryptographic randomness each run",
                   "configurationOrderPolicy": benchmarkOrder,
                   "fixtureIds": selected.map(\.id), "startedAt": ISO8601DateFormatter().string(from: Date()),
                   "status": "running"]
        task = Task {
            defer { busy = false; activeBenchmark = false; UIApplication.shared.isIdleTimerDisabled = false; refreshExports() }
            do {
                for (configIndex, item) in selected.enumerated() {
                    try Task.checkCancellation()
                    let p = item.parameters
                    status = "Loading \(item.label)"
                    _ = try await engine.call(["action": "prepare", "parameters": p.object])
                    let base: [String: Any] = ["action": "generate", "scores": Array(repeating: (p.min+p.max)/2, count: p.l),
                        "protocol": protocolName, "electionId": "benchmark-\(runID)-\(item.id)", "voterIndex": 1]
                    for index in 0..<warm {
                        try Task.checkCancellation()
                        var request = base; request["verify"] = index == 0
                        _ = try await engine.call(request)
                    }
                    for index in 0..<runs {
                        try Task.checkCancellation()
                        guard ProcessInfo.processInfo.thermalState != .serious && ProcessInfo.processInfo.thermalState != .critical else {
                            throw PVoteError("Device too warm. Data saved")
                        }
                        var request = base; request["verify"] = verify
                        var row = try await measuredCall(request)
                        row["sampleIndex"] = index + 1; row["fixtureId"] = item.id
                        benchmarkRows.append(row)
                        progress = Double(configIndex*runs+index+1) / Double(selected.count*runs)
                        status = "\(item.label) · \(index+1)/\(runs)"
                        try persistBenchmark()
                    }
                }
                session["status"] = "complete"
                session["completedAt"] = ISO8601DateFormatter().string(from: Date())
                status = "Completed \(benchmarkRows.count) measurements. JSON and CSV ready"
            } catch {
                session["status"] = "interrupted"
                session["interruptionReason"] = error is CancellationError ? "cancelled or app left foreground" : error.localizedDescription
                status = error is CancellationError ? "Benchmark stopped. Data saved" : error.localizedDescription
            }
            try? persistBenchmark()
        }
    }

    func cancelBenchmark() { if activeBenchmark { task?.cancel() } }

    private func persistBenchmark() throws {
        var report = session; report["samples"] = benchmarkRows
        let file = documents.appendingPathComponent("benchmark-\(runID).json")
        try JSONSerialization.data(withJSONObject: report, options: [.sortedKeys, .prettyPrinted]).write(to: file, options: .atomic)
        let fields = session["experiment"] as? String == "score_domain_generation" ?
                     ["fixtureId","sampleIndex","scheme","phase","kappa","budget","n","l","t","generationMs","generationCpuMs","serializationMs","transcriptBytes","transcriptHash","parameterFingerprint","timestamp"] :
                     ["fixtureId","sampleIndex","protocol","n","l","t","shareMs","shareCpuMs","encodeMs","verifyMs",
                      "transcriptBytes","calldataBytes","shareAllocatedBytes","shareMallocs","shareGCs",
                      "nativeCallSampledPeakBytes","nativeCallFootprintBeforeBytes","nativeCallFootprintAfterBytes","timestamp"]
        func escape(_ value: Any?) -> String { "\"" + (value.map { String(describing: $0) } ?? "").replacingOccurrences(of: "\"", with: "\"\"") + "\"" }
        let csv = ([fields.joined(separator: ",")] + benchmarkRows.map { row in fields.map { escape(row[$0]) }.joined(separator: ",") }).joined(separator: "\n") + "\n"
        try Data(csv.utf8).write(to: file.deletingPathExtension().appendingPathExtension("csv"), options: .atomic)
    }

    private func refreshExports() {
        exportURLs = ((try? FileManager.default.contentsOfDirectory(at: documents, includingPropertiesForKeys: nil)) ?? [])
            .filter { $0.lastPathComponent.hasPrefix("benchmark-") || $0.lastPathComponent.hasPrefix("vote-") }
            .sorted { $0.lastPathComponent > $1.lastPathComponent }
    }

    func submitBallot() {
        guard !busy, let election, let calldata = votingCalldata, let topic = acceptanceTopic else { return }
        guard pending == nil else { status = "Check the pending transaction first"; return }
        busy = true
        task = Task {
            defer { busy = false; refreshExports() }
            do {
                guard address.lowercased() == election.voterAddress.lowercased() else { throw PVoteError("Account not registered for this election") }
                guard let key = try Wallet.load(), let url = URL(string: election.rpcURL) else { throw PVoteError("Account or RPC URL unavailable") }
                let chain = ChainClient(url: url)
                let chainID = try ChainClient.quantity(try await chain.rpc("eth_chainId", []))
                guard String(chainID) == election.chainId else { throw PVoteError("RPC chain ID mismatch") }
                guard let code = try await chain.rpc("eth_getCode", [election.contract, "latest"]) as? String, code != "0x" else { throw PVoteError("No contract at destination") }
                let started = ProcessInfo.processInfo.systemUptime
                let tx: [String: Any] = ["from": address, "to": election.contract, "data": calldata]
                let gas = try ChainClient.quantity(try await chain.rpc("eth_estimateGas", [tx]))
                let nonce = try ChainClient.quantity(try await chain.rpc("eth_getTransactionCount", [address, "pending"]))
                let price = try ChainClient.quantity(try await chain.rpc("eth_gasPrice", []))
                let signStart = ProcessInfo.processInfo.systemUptime
                let signed = try await engine.call(["action": "sign", "privateKey": key, "transaction": ["chainId": election.chainId,
                    "nonce": nonce, "gas": gas + gas/5, "gasPrice": String(price), "to": election.contract, "data": calldata]])
                guard let hash = signed["transactionHash"] as? String, let raw = signed["rawTransaction"] as? String else { throw PVoteError("Signing failed") }
                let signingMs = (ProcessInfo.processInfo.systemUptime-signStart)*1000
                // Persist before sending: a network timeout must not cause a new transaction.
                pending = ["transactionHash": hash, "rpcURL": election.rpcURL, "contract": election.contract, "acceptanceTopic": topic,
                    "submittedAt": Date().timeIntervalSince1970, "signingMs": signingMs, "estimatedGas": gas,
                    "preBroadcastMs": (ProcessInfo.processInfo.systemUptime-started)*1000,
                    "generation": lastMeasurement ?? [:], "metadata": metadata]
                try savePending()
                status = "Broadcasting transaction…"
                guard let returned = try await chain.rpc("eth_sendRawTransaction", [raw]) as? String,
                      returned.lowercased() == hash.lowercased() else { throw PVoteError("Transaction hash mismatch. Check the saved transaction") }
                try await pollPending()
            } catch {
                if pending != nil {
                    pending?["lastRPCError"] = error.localizedDescription
                    pending?["lastRPCErrorAt"] = ISO8601DateFormatter().string(from: Date())
                    try? savePending()
                }
                status = error.localizedDescription + (pending == nil ? "" : ". Transaction hash saved for polling")
            }
        }
    }

    func checkPending() {
        guard !busy, pending != nil else { return }; busy = true
        task = Task { defer { busy = false; refreshExports() }; do { try await pollPending() } catch { status = error.localizedDescription } }
    }

    private func savePending() throws {
        guard let pending else { return }
        try JSONSerialization.data(withJSONObject: pending, options: [.sortedKeys]).write(to: documents.appendingPathComponent("pending.json"), options: .atomic)
    }

    private func pollPending() async throws {
        guard let record = pending, let hash = record["transactionHash"] as? String,
              let rpc = record["rpcURL"] as? String, let url = URL(string: rpc),
              let contract = record["contract"] as? String, let topic = record["acceptanceTopic"] as? String else { return }
        let chain = ChainClient(url: url)
        for _ in 0..<120 {
            try Task.checkCancellation()
            if let receipt = try await chain.receipt(hash) {
                var report = record
                report["receipt"] = receipt
                report["accepted"] = ChainClient.accepted(receipt, contract: contract, topic: topic)
                report["receiptObservedAt"] = Date().timeIntervalSince1970
                report["submissionToReceiptObservationMs"] = (Date().timeIntervalSince1970-(record["submittedAt"] as? Double ?? Date().timeIntervalSince1970))*1000
                report["pollIntervalMs"] = 1000
                let filename = documents.appendingPathComponent("vote-\(hash).json")
                try JSONSerialization.data(withJSONObject: report, options: [.prettyPrinted, .sortedKeys]).write(to: filename, options: .atomic)
                pending = nil
                try? FileManager.default.removeItem(at: documents.appendingPathComponent("pending.json"))
                votingCalldata = nil
                guard report["accepted"] as? Bool == true else { throw PVoteError("Transaction mined without ballot acceptance. See receipt") }
                status = "Ballot accepted · gas \((try? ChainClient.quantity(receipt["gasUsed"] ?? "")) ?? 0)"
                return
            }
            status = "Transaction pending: \(hash.prefix(14))…"
            try await Task.sleep(for: .seconds(1))
        }
        status = "Still pending. Transaction hash saved"
    }

    private func environment() -> [String: Any] {
        let p = ProcessInfo.processInfo
        let names = ["nominal", "fair", "serious", "critical"]
        return ["thermalState": names[min(p.thermalState.rawValue, 3)], "lowPowerMode": p.isLowPowerModeEnabled,
                "batteryLevel": UIDevice.current.batteryLevel, "batteryState": UIDevice.current.batteryState.rawValue]
    }

    private func deviceMetadata() -> [String: Any] {
        var info = utsname(); uname(&info)
        let machineSize = MemoryLayout.size(ofValue: info.machine)
        let machine = withUnsafePointer(to: &info.machine) { pointer in
            pointer.withMemoryRebound(to: CChar.self, capacity: machineSize) { String(cString: $0) }
        }
        #if targetEnvironment(simulator)
        let simulator = true
        #else
        let simulator = false
        #endif
        #if DEBUG
        let configuration = "Debug"
        #else
        let configuration = "Release"
        #endif
        return ["hardwareIdentifier": machine, "deviceModel": UIDevice.current.model, "iOSVersion": UIDevice.current.systemVersion,
                "isSimulator": simulator, "buildConfiguration": configuration,
                "appVersion": Bundle.main.object(forInfoDictionaryKey: "CFBundleShortVersionString") as? String ?? "",
                "processorCount": ProcessInfo.processInfo.processorCount, "physicalMemoryBytes": ProcessInfo.processInfo.physicalMemory,
                "osVersion": ProcessInfo.processInfo.operatingSystemVersionString]
    }
}
