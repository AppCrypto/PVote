import SwiftUI
import UniformTypeIdentifiers

struct ContentView: View {
    @ObservedObject var model: AppModel
    @State private var importing = false
    private let accent = Color(red: 0.07, green: 0.43, blue: 0.43)
    var body: some View {
        TabView {
            NavigationStack { voting.navigationTitle("PVote") }
                .tabItem { Label("Voting", systemImage: "checkmark.seal") }
            NavigationStack { benchmark.navigationTitle("Voter benchmark") }
                .tabItem { Label("Benchmark", systemImage: "speedometer") }
            NavigationStack { results.navigationTitle("Results") }
                .tabItem { Label("Results", systemImage: "chart.xyaxis.line") }
        }
        .tint(accent)
        .safeAreaInset(edge: .bottom) {
            Text(model.status).font(.caption).frame(maxWidth: .infinity, alignment: .leading)
                .padding(10).background(.regularMaterial)
        }
    }

    private var voting: some View {
        Form {
            Section("Local voting account") {
                Text(model.address.isEmpty ? "Loading…" : model.address)
                    .font(.caption.monospaced()).textSelection(.enabled)
                Button { UIPasteboard.general.string = model.address } label: {
                    Label("Copy address", systemImage: "doc.on.doc")
                }
            }
            Section("Election") {
                if let election = model.election {
                    LabeledContent("Mode", value: "On-chain experiment")
                    Text(election.electionId).font(.caption).textSelection(.enabled)
                    Text(election.contract).font(.caption.monospaced())
                    Button("Use offline demo") { model.useOffline() }
                } else {
                    Picker("Public parameters", selection: $model.selectedFixtureID) {
                        ForEach(model.fixtures) { Text($0.label).tag($0.id) }
                    }
                    .onChange(of: model.selectedFixtureID) { _, _ in model.resetScores() }
                    Picker("Proof version", selection: $model.proofProtocol) {
                        Text("Context-bound").tag("context-v1")
                        Text("Legacy verifier").tag("legacy-solidity-v1")
                    }
                    .onChange(of: model.proofProtocol) { _, _ in model.invalidateBallot() }
                }
                Button { importing = true } label: { Label("Import election JSON", systemImage: "square.and.arrow.down") }
                if let p = model.parameters { Text("\(p.n) talliers · threshold \(p.t) · \(p.l) options").font(.caption).foregroundStyle(.secondary) }
            }
            Section("Scores") {
                if let p = model.parameters {
                    ForEach(model.scores.indices, id: \.self) { index in
                        Stepper(value: $model.scores[index], in: p.min...p.max) {
                            LabeledContent(optionLabel(index), value: String(model.scores[index]))
                        }
                    }
                    .onChange(of: model.scores) { _, _ in model.invalidateBallot() }
                }
                Button { model.generateBallot() } label: {
                    Label("Generate and verify ballot", systemImage: "lock.shield")
                }.buttonStyle(.borderedProminent)
            }
            if let result = model.lastMeasurement {
                Section("Local computation") {
                    metric("RB-PVSS.Share", result["shareMs"], unit: "ms")
                    metric("Encoding", result["encodeMs"], unit: "ms")
                    metric("Local verification", result["verifyMs"], unit: "ms")
                    metric("Ballot", result["transcriptBytes"], unit: "bytes", digits: 0)
                    metric("Calldata", result["calldataBytes"], unit: "bytes", digits: 0)
                }
            }
            if model.election != nil {
                Section("Submission and confirmation") {
                    Button("Sign and submit ballot") { model.submitBallot() }.disabled(!model.canSend)
                    if let hash = model.pending?["transactionHash"] as? String {
                        Text(hash).font(.caption.monospaced()).textSelection(.enabled)
                        Button("Check transaction") { model.checkPending() }
                    }
                    Text("Keep the transaction hash after connection loss. Do not create a replacement automatically. No tallying participation is needed after acceptance.")
                        .font(.caption).foregroundStyle(.secondary)
                }
            }
        }
        .disabled(model.busy)
        .overlay(alignment: .topTrailing) { if model.busy { ProgressView().padding() } }
        .fileImporter(isPresented: $importing, allowedContentTypes: [.json]) { result in
            switch result { case .success(let url): model.importElection(url); case .failure(let error): model.status = error.localizedDescription }
        }
    }

    private var benchmark: some View {
        Form {
            Section("Benchmark settings") {
                Picker("Parameter set", selection: $model.benchmarkSuite) {
                    Text("Three paper configurations").tag("paper3")
                    Text("Full n/l curve (23 configurations)").tag("curve")
                    Text("Current parameters").tag("current")
                }
                Picker("Proof version", selection: $model.proofProtocol) {
                    Text("Context-bound").tag("context-v1")
                    Text("Legacy verifier").tag("legacy-solidity-v1")
                }
                Stepper("Warmups: \(model.warmups) per configuration", value: $model.warmups, in: 1...10)
                Stepper("Samples: \(model.repetitions) per configuration", value: $model.repetitions, in: 3...200)
                Toggle("Verify each measured ballot", isOn: $model.verifyEach)
                Text("Share and encoding are timed separately. Setup and loading are excluded. The first warmup is verified.")
                    .font(.caption).foregroundStyle(.secondary)
            }.disabled(model.busy)
            Section {
                if model.busy {
                    ProgressView(value: model.progress)
                    Button("Stop and save", role: .destructive) { model.cancelBenchmark() }
                } else {
                    Button { model.runBenchmark() } label: { Label("Run benchmark", systemImage: "play.fill") }
                        .buttonStyle(.borderedProminent)
                }
            }
            Section("Measurements") {
                Label("Share wall time, CPU time, and encoding", systemImage: "clock")
                Label("Ballot size, allocations, and GC", systemImage: "doc.text")
                Label("Sampled app footprint (2 ms)", systemImage: "memorychip")
                Label("Device, OS, build, temperature, and battery", systemImage: "iphone")
                Text("Keep the app in the foreground. Backgrounding or high temperature stops measurement. Battery level is not energy usage.")
                    .font(.caption).foregroundStyle(.secondary)
            }
        }
    }

    private var results: some View {
        List {
            Section("Generation summary") {
                let grouped = Dictionary(grouping: model.benchmarkRows) { $0["fixtureId"] as? String ?? "" }
                if grouped.isEmpty { Text("Run a benchmark to view mean, SD, and P95").foregroundStyle(.secondary) }
                ForEach(grouped.keys.sorted(), id: \.self) { key in
                    let samples = (grouped[key] ?? []).compactMap { ($0["generationMs"] ?? $0["shareMs"]) as? Double }.sorted()
                    if !samples.isEmpty {
                        let mean = samples.reduce(0,+)/Double(samples.count)
                        let sd = samples.count > 1 ? sqrt(samples.reduce(0) { $0 + pow($1-mean,2) }/Double(samples.count-1)) : 0
                        let p95 = samples[max(0, Int(ceil(Double(samples.count)*0.95))-1)]
                        VStack(alignment: .leading, spacing: 5) {
                            Text(key).font(.subheadline.monospaced())
                            Text(String(format: "%.2f ± %.2f ms · P95 %.2f ms · n=%d", mean, sd, p95, samples.count))
                                .font(.caption).foregroundStyle(.secondary)
                        }
                    }
                }
            }
            Section("Raw reports") {
                if model.exportURLs.isEmpty { Text("No exports").foregroundStyle(.secondary) }
                ForEach(model.exportURLs, id: \.self) { url in
                    ShareLink(item: url) { Label(url.lastPathComponent, systemImage: "square.and.arrow.up").font(.caption) }
                }
                Text("Use complete physical-device Release reports for analysis. Simulator and interrupted runs are labeled.")
                    .font(.caption).foregroundStyle(.secondary)
            }
        }
    }

    private func optionLabel(_ index: Int) -> String {
        if let labels = model.election?.labels, labels.indices.contains(index) { return labels[index] }
        return "Option \(index+1)"
    }
    private func metric(_ title: String, _ value: Any?, unit: String, digits: Int = 2) -> some View {
        let number = (value as? NSNumber)?.doubleValue ?? 0
        return LabeledContent(title, value: String(format: "%.*f %@", digits, number, unit))
    }
}
