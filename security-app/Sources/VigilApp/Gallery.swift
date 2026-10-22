import SwiftUI
import VigilDesign

@main
struct VigilApp: App {
    #if os(macOS)
    @NSApplicationDelegateAdaptor(VigilApplicationDelegate.self) private var delegate
    #endif
    init() {
        #if os(macOS)
        renderSnapshotsIfRequested()
        #endif
    }
    var body: some Scene {
        WindowGroup("Vigil") {
            #if os(macOS)
            Group {
                if CommandLine.arguments.contains("--studio") { Gallery() }
                else { SecurityWorkspace() }
            }.frame(minWidth: 1000, minHeight: 640)
            #else
            TabView {
                PhoneInvestigation().tabItem { Label("Investigation", systemImage: "magnifyingglass") }
                AgentView().tabItem { Label("Optional assistant", systemImage: "bubble.left") }
            }
            #endif
        }
        #if os(macOS)
        .defaultSize(width: 1360, height: 900)
        #endif
    }
}

struct Gallery: View {
    @Environment(\.colorScheme) private var systemScheme
    @State private var section: Section? = .overview
    @State private var dark = false
    @State private var reduced = false
    var body: some View {
        NavigationSplitView {
            sidebar.navigationSplitViewColumnWidth(min: 210, ideal: 230, max: 270)
        } detail: {
            page.background(Ink.canvas.color(dark ? .dark : .light))
                .foregroundStyle(Ink.primary.color(dark ? .dark : .light))
                .toolbar { toolbar }
        }
        .tint(Ink.accent.color(dark ? .dark : .light)).preferredColorScheme(dark ? .dark : .light)
        .environment(\.vigilReduceMotion, reduced)
        .onAppear { dark = systemScheme == .dark }
    }
    private var sidebar: some View {
        VStack(alignment: .leading, spacing: 28) {
            Brand().padding(.horizontal, 20).padding(.top, 24)
            List(Section.allCases, selection: $section) { item in
                Label(item.rawValue, systemImage: item.symbol).padding(.vertical, 7).tag(item)
            }.listStyle(.sidebar)
            VStack(alignment: .leading, spacing: 8) {
                Label("Design studio", systemImage: "square.dashed").font(.callout.weight(.medium))
                Text("Synthetic data\nNo device monitoring").font(.caption).foregroundStyle(Ink.secondary.color(dark ? .dark : .light))
            }.padding(20)
        }.background(Ink.inset.color(dark ? .dark : .light))
    }
    @ViewBuilder private var page: some View {
        switch section ?? .overview {
        case .overview: Overview()
        case .components: ComponentCatalog()
        case .identity: IdentityPage()
        }
    }
    @ToolbarContentBuilder private var toolbar: some ToolbarContent {
        ToolbarItemGroup {
            Toggle(isOn: $reduced) { Label("Reduce motion", systemImage: "figure.stand") }
                .help("Preview without motion; system Reduce Motion always takes priority")
            Toggle(isOn: $dark) { Label("Dark appearance", systemImage: "moon") }
                .help("Switch light and dark appearance")
        }
    }
}

#if os(macOS)
final class VigilApplicationDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        NSApplication.shared.setActivationPolicy(.regular)
        NSApplication.shared.activate(ignoringOtherApps: true)
        for window in NSApplication.shared.windows { window.makeKeyAndOrderFront(nil) }
    }
}
#endif
