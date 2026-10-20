# Osprey Storage: where your disk space goes, on the Mac and iPhone

One SwiftUI app that shows what uses your disk, folder by folder. It opens on your home folder with hidden files shown, ranks everything inside by size, and lets you click into any folder. The eye button in the toolbar (⇧⌘.) leaves hidden files out of every total; ⌘↑ goes to the enclosing folder, ⌘R measures again and ⌘O measures another folder or drive. It never shows a privacy prompt: folders macOS guards wait until you choose them.

Osprey owns the app. [`storage/src`](storage/src) holds the browser's state and events, the ranking, the hidden-file rule, the symbols, the formatting and every screen, which Osprey returns as a tree of native nodes. Swift does two things: it draws that tree with the platform's own grouped list, toolbar, toggle and capacity bar, and it measures the disk.

## Build and run

Use an Apple silicon Mac with Xcode selected by `xcode-select`. From the repository root:

```sh
make macos                          # Mac runtime, compiler and every macOS example
examples/macos/run.sh               # build and open the Mac app
examples/macos/run.sh ios-sim       # the same app in the iPhone simulator
make macos-test                     # C ABI, corpus, domain suite and smoke
```

`run.sh --smoke` builds a folder with known contents, including a hidden `.git` inside a visible folder, a hidden `.cache` and an unreadable folder, then drives Osprey through the events the app sends: it requires the walk to count them, the screen to show and explain them, the toggle to leave the hidden ones out, and a folder to open. `run.sh ios-sim-smoke` runs the same check in the simulator. `OSPREY_MACOS_SKIP_BUILD=1` reuses the existing Mac build, and `OSPREY_BIN` selects another compiler.

## One app, two platforms

The iPhone apps and this one build the same way: [`examples/ios/build.sh`](../ios/build.sh) compiles the Osprey project into `libOspreyApp.a` and its generated `libOspreyApp.h`, and `xcodebuild` links them into the app. [`OspreyStorage.xcconfig`](OspreyStorage.xcconfig) declares one target for `macosx`, `iphoneos` and `iphonesimulator` and picks the matching archive (`--target=macos`, `ios` or `ios-sim`). The Swift sources, the renderer and the Osprey code are the same on both. See [the Mac slice](../../docs/specs/compiler/iOSCodegen.md#the-mac-slice-ios-target-macos) for the boundary.

The renderer is the one [Issue Inbox](../mobile/ios/IssueInbox/NativeRenderer.swift) uses. This app adds node kinds that map straight to stock controls: `navigation` (a titled screen of grouped sections with a toolbar), `entry`, `gauge`, `stack`, `pathbar`, `toggle`, `progress` and `empty`. Osprey chooses every word, symbol, colour, share and key equivalent.

On iPhone, home is the app's own container. Choose Folder opens the Files picker, and the app holds the chosen folder's security scope until the next pick.

## The boundary

[`storage/src/main.ospml`](storage/src/main.ospml) exports two functions, and the [store](OspreyStorage/StorageStore.swift) calls nothing else:

```ospml
start home = envelope (Browse::start home)
dispatch (model, event) = envelope (Browse::dispatch (Browse::decode model) event)
```

Each returns the browser as opaque text, the screen, and any host work: `measure` a folder (with the folders not to enter, and whether to add it to the walk already made) or `pick` one. Swift does that work and reports back as events (`progress`, `measured`, `failed`, `picked`), and Osprey ignores a report about a folder it has since left.

The screen asks the platform through an effect, so [`View`](storage/src/view.ospml) never touches the file system and the [domain suite](storage/test/main.ospml) renders it from fixed data:

```ospml
export effect Facts
    volume : string => Volume
    listing : string => Listing
```

`main.ospml` answers it with C functions the Swift host defines with `@_cdecl` in [`StorageServices.swift`](OspreyStorage/StorageServices.swift). All calls are synchronous on the main thread. Swift strings passed to Osprey are borrowed for the call. A string a Swift import returns must stay valid for the rest of the process, so each listing's paths are copied once per walk and never freed. Calling Swift crosses Osprey's memory-safety boundary.

Ranking loops over a folder's entries and keeps only the 120 largest, so a folder of any width keeps the stack flat. The iPhone main thread has a 1 MB stack, and Osprey does not turn tail calls into loops.

## Where the numbers come from

- **Volume capacity:** Foundation's `volumeTotalCapacity`, `volumeAvailableCapacity` and `volumeAvailableCapacityForImportantUsage`. The last includes space the system can purge on demand, such as caches and local snapshots, which is what Settings shows as available.
- **Folder sizes:** one walk of the chosen folder with Foundation's directory enumerator. Each folder records the allocated size of everything below it, and the size of everything below it outside hidden items, so the hidden-files toggle needs no second walk. Hidden files, caches and package contents count. A hard-linked file counts once, a mounted volume below the folder is not entered, and APFS clones count in full, as in Finder. An iCloud Drive or File Provider folder whose contents were evicted is not entered either: it holds nothing on this disk, and listing it would make the provider fetch it from the network. "Hidden" means what Finder hides: names starting with a dot and items flagged hidden, such as `~/Library`.

The system keeps no record of folder sizes that an app can read. APFS directory statistics return `ATTR_DIR_ALLOCSIZE` as 0 for ordinary folders, and the Storage pane in System Settings sizes folders through the private SpaceAttribution framework. Spotlight records file sizes but never indexes hidden files or caches, so measuring is the only complete source. On an Apple silicon Mac, a 330 GB home folder of several million files takes about three and a half minutes, with live progress.

## Permissions

The app never asks for access. It reads folder metadata, the same information Finder shows, and never asks for Full Disk Access.

macOS asks the person before any app reads Desktop, Documents, Downloads, iCloud Drive, cloud storage folders or another app's container, and an app that opens one waits until the question is answered. Osprey lists those folders ([`Browse`](storage/src/browse.ospml)) and sends them with every walk, and the walk does not enter them. They appear under "Protected by macOS" instead. Choosing one opens the system folder picker already on that folder; choosing it there is the person's consent, so macOS shows no prompt, and its walk is added into the totals already measured.

Other apps' private data, such as Mail and Messages, is closed to every app. The walk counts those folders in the "could not be read" footer instead of silently shrinking a total.
