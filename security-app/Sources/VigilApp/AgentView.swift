import SwiftUI
import UniformTypeIdentifiers
import VigilCore

struct AgentView: View {
    @State private var enabled = false
    @State private var approved = false
    @State private var token = ""
    @State private var question = ""
    @State private var reply = ""
    @State private var failure = ""
    @State private var choices: [AgentConfiguration] = []
    @State private var selected: UUID?
    @State private var importing = false
    @State private var request: Task<Void, Never>?
    private var model: AgentConfiguration? { choices.first { $0.id == selected } }
    var body: some View {
        Form {
            SwiftUI.Section("Optional NIMBLESITE assistant") {
                Toggle("Enable optional AI", isOn: $enabled)
                Text("All local investigation features work with AI off. This assistant receives only the question you type; it has no access to your device evidence or local tools.")
            }
            if enabled { configuration; disclosure; answer }
        }.formStyle(.grouped)
            .fileImporter(isPresented: $importing, allowedContentTypes: [.json]) { result in importModels(result) }
            .onChange(of: enabled) { _, value in if !value { stop(); token = ""; reply = ""; approved = false } }
            .onChange(of: selected) { _, _ in stop(); reply = ""; approved = false }
            .onChange(of: token) { _, _ in stop(); approved = false; reply = "" }
            .onChange(of: question) { _, _ in approved = false }
            .onDisappear { stop(); token = "" }
    }
    private var configuration: some View {
        SwiftUI.Section("Your deployed model configurations") {
            Button("Import model configurations") { importing = true }
            Picker("Model", selection: $selected) {
                Text("Select a configuration").tag(UUID?.none)
                ForEach(choices) { Text("\($0.name) · \($0.provider)/\($0.model)").tag(Optional($0.id)) }
            }
            SecureField("End-user JWT (session only)", text: $token)
            Text("Import configurations supplied by your deployment administrator. Model labels are supplied by that file; verify the actual server configuration. Never enter an operator API key. Tokens are discarded on leaving this view.").font(.caption)
        }
    }
    private var disclosure: some View {
        SwiftUI.Section("Disclosure preview — exactly what will be sent") {
            TextEditor(text: $question).frame(minHeight: 100)
            Text("Destination: api.nimblesite.ai and the configured model provider. NIMBLESITE retains conversations and tenant operators may access them. Removing text here does not delete server history. A request may incur a charge; it is never automatically retried.").font(.caption)
            Toggle("I approve sending this question to the selected configuration", isOn: $approved)
            HStack {
                Button("Send question") { send() }.disabled(!approved || model == nil || token.isEmpty || request != nil)
                if request != nil { ProgressView().controlSize(.small); Button("Cancel") { stop() } }
            }
        }
    }
    private var answer: some View {
        SwiftUI.Section("Assistant response — unverified explanation") {
            Text(reply.isEmpty ? "No request has been sent." : reply).textSelection(.enabled)
            if !failure.isEmpty { Text(failure).foregroundStyle(.red) }
            Text("Responses cannot change findings or run actions. Each question starts a new conversation; switching models never forwards prior conversations.").font(.caption)
        }
    }
    private func send() {
        guard let model else { return }
        do {
            let outbound = try agentRequest(configuration: model, token: token, message: question, consent: enabled && approved)
            failure = ""; reply = ""; approved = false
            request = Task {
                do { let result = try await askAgent(outbound); try Task.checkCancellation(); reply = result }
                catch { if !Task.isCancelled { failure = error.localizedDescription } }
                if !Task.isCancelled { request = nil }
            }
        } catch { failure = error.localizedDescription }
    }
    private func stop() { request?.cancel(); request = nil }
    private func importModels(_ result: Result<URL, Error>) {
        do {
            let url = try result.get(), access = url.startAccessingSecurityScopedResource()
            defer { if access { url.stopAccessingSecurityScopedResource() } }
            let models = try JSONDecoder().decode([AgentConfiguration].self, from: boundedData(url, limit: 100_000))
            guard !models.isEmpty, models.count <= 100, Set(models.map(\.id)).count == models.count else { throw CollectionError.invalid("Invalid model catalog") }
            stop(); choices = models; selected = models.first?.id; approved = false; reply = ""
        } catch { failure = error.localizedDescription }
    }
}
