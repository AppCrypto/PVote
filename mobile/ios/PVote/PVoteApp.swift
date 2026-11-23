import SwiftUI

@main
struct PVoteApp: App {
    @StateObject private var model = AppModel()
    @Environment(\.scenePhase) private var scenePhase
    var body: some Scene {
        WindowGroup {
            ContentView(model: model)
                .task { await model.start() }
                .onChange(of: scenePhase) { _, phase in
                    if phase != .active { model.cancelBenchmark() }
                }
        }
    }
}
