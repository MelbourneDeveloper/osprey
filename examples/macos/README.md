# Osprey on the Mac: which folders use your disk

A command-line Mac app that reports which folders take the most disk space, hidden files included. Osprey owns the report. A Swift host answers its questions: volume capacity from Foundation, and folder sizes by measuring the tree.

Output of `examples/macos/run.sh ~ 3`:

```text
Storage report for /Users/ada
Volume: 994.7 GB total, 958.2 GB used, 36.5 GB free, 189.9 GB available for important data.
Measured 566.8 GB here, hidden files included (59% of used space).
Folders that could not be read: 157. macOS protects them; grant the host Full Disk Access to count them.

Largest items
 245.6 GB  ████████░░░░░░░░░░░░  43%  Documents
 151.5 GB  █████░░░░░░░░░░░░░░░  27%  Library
  83.8 GB  ███░░░░░░░░░░░░░░░░░  15%  .colima

Inside Documents
 207.2 GB  ████████████████░░░░  84%  Stuff
  37.0 GB  ███░░░░░░░░░░░░░░░░░  15%  Code
   1.2 GB  ░░░░░░░░░░░░░░░░░░░░   0%  Desktops
```

## The effect at the boundary

[`storage/src/disk.ospml`](storage/src/disk.ospml) declares what the report needs from the platform and performs it while building the report:

```ospml
export effect Facts
    volume : string => Volume
    largest : string => Listing

insideLargest limit folders = match folders
    [] => []
    [largest, ..._rest] => top limit (perform Facts.largest largest.path).entries
```

[`storage/src/main.ospml`](storage/src/main.ospml) answers that effect with C functions the Swift host defines:

```ospml
extern storage_scan (root : string) -> int
extern storage_unreadable -> int

report (root, limit) =
    handle Disk::Facts
        volume path => Volume(total = storage_volume path "total", available = storage_volume path "available", important = storage_volume path "important")
        largest folder =>
            found = entries folder
            Listing(entries = found, blocked = storage_unreadable ())
    Disk::render (Disk::survey root limit)
```

[`OspreyStorage/StorageServices.swift`](OspreyStorage/StorageServices.swift) implements each import with `@_cdecl`, and [`main.swift`](OspreyStorage/main.swift) calls the generated `osprey_storage_report`. The [domain suite](storage/test/main.ospml) runs the same survey and report under a fixed handler, natively and without Swift.

## Where the numbers come from

- **Volume capacity:** Foundation's `volumeTotalCapacity`, `volumeAvailableCapacity` and `volumeAvailableCapacityForImportantUsage`. The last includes space macOS can purge on demand, such as caches and local snapshots, which is what System Settings shows as available.
- **Folder sizes:** a walk of the folder with Foundation's directory enumerator, which lists each directory in bulk. Every file adds its allocated size to the child of the folder that holds it and to that child's child, so one walk also answers the drill-down. Hidden files, caches and package contents count. A hard-linked file counts once, a mounted volume below the folder is not entered, and APFS clones count in full, as in Finder.

macOS keeps no record of folder sizes that an app can read. APFS directory statistics return `ATTR_DIR_ALLOCSIZE` as 0 for ordinary folders, and the Storage pane in System Settings sizes folders through the private SpaceAttribution framework. Spotlight records file sizes but never indexes hidden files or caches: on the machine above it accounted for 17.5 GB of a 236 GB `Documents` folder and 95 GB of the 567 GB home folder, and `.colima` was invisible to it. Measuring is the only complete source, and a home folder of that size takes about a minute and a half.

## Permissions

Hidden files need no permission. "Hidden" is a Finder display convention, and the walk reads them like any other file.

What macOS protects is other apps' data and a few user folders. A folder the walk cannot read is counted in the report's "could not be read" line instead of silently shrinking a total. macOS asks the person automatically the first time an app reads Desktop, Documents or Downloads, or another app's container; a bundled app explains why in its `NSDesktopFolderUsageDescription`, `NSDocumentsFolderUsageDescription` and `NSDownloadsFolderUsageDescription` strings. A command-line host borrows the permissions of the terminal that runs it. An app can also ask for one folder by letting the person choose it in `NSOpenPanel`, which grants that folder's whole tree.

Mail, Messages, Safari and similar data need Full Disk Access, which no API can request. An app can only open the Full Disk Access pane (`x-apple.systempreferences:com.apple.preference.security?Privacy_AllFiles`) and ask the person to add it. Apple's [2 October 2026 notice](https://developer.apple.com/news/?id=p6zjojqw) says granting it will soon require even more explicit action.

## Build and run

Use an Apple silicon Mac with Xcode selected by `xcode-select`. From the repository root:

```sh
make macos                          # Mac runtime, compiler and every macOS example
examples/macos/run.sh               # report on your home folder
examples/macos/run.sh /Applications 5
make macos-test                     # C ABI, corpus, domain suite and smoke
```

`run.sh --smoke` builds a folder with known contents, including a hidden `.git` inside a visible folder, a hidden `.cache` and an unreadable folder, and requires the report to count and name them. `OSPREY_MACOS_SKIP_BUILD=1` reuses the existing build, and `OSPREY_BIN` selects another compiler.

## One boundary for Mac and iPhone

[`scripts/macos-host.sh`](../../scripts/macos-host.sh) compiles the project with `--target=macos` into `libOspreyApp.a` and `libOspreyApp.h`, which [`Bridge.h`](OspreyStorage/Bridge.h) imports into Swift as the iPhone apps import theirs, then links the Swift sources. The [graphics scenes](../graphics/README.md) build the same way. This is the [iPhone target](../ios) built for the Mac: the same header, scalar types, entry, thunks and target rules. A Swift `@_cdecl` function written for one platform links against the other unchanged; only the Apple API inside it has to exist on both. See [the Mac slice](../../docs/specs/compiler/iOSCodegen.md#the-mac-slice-ios-target-macos) for the contract.

Calls are synchronous on one host thread. Swift strings passed to Osprey are borrowed for the call; Swift copies strings Osprey returns with `String(cString:)`. A string a Swift import returns must stay valid for the rest of the process, so `storage_entry_path` hands over a `strdup` copy that Swift never frees. Calling Swift crosses Osprey's memory-safety boundary.
