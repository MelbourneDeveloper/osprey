# Osprey on the Mac: which folders use your disk

A command-line Mac app that reports which folders take the most disk space. Osprey owns the report. A Swift host answers its questions from records macOS already keeps, Foundation's volume capacities and the sizes in the Spotlight index, so no directory is opened or walked.

Output of `examples/macos/run.sh ~ 3`:

```text
Storage report for /Users/ada
Volume: 994.7 GB total, 897.6 GB used, 97.0 GB free, 188.3 GB available for important data.
Spotlight attributes 95.4 GB here (11% of used space).

Largest items
  71.4 GB  ███████████████░░░░░  75%  Library
  17.5 GB  ███░░░░░░░░░░░░░░░░░  18%  Documents
   3.3 GB  ░░░░░░░░░░░░░░░░░░░░   3%  Downloads

Inside Library
  44.2 GB  ████████████░░░░░░░░  62%  Application Support
  24.4 GB  ██████░░░░░░░░░░░░░░  34%  Android
 957.6 MB  ░░░░░░░░░░░░░░░░░░░░   1%  Python
```

## The effect at the boundary

[`storage/src/disk.ospml`](storage/src/disk.ospml) declares what the report needs from the platform and performs it while building the report:

```ospml
export effect Facts
    volume : string => Volume
    largest : string => List<Folder>

insideLargest limit folders = match folders
    [] => []
    [largest, ..._rest] => top limit (perform Facts.largest largest.path)
```

[`storage/src/main.ospml`](storage/src/main.ospml) answers that effect with C functions the Swift host defines:

```ospml
extern storage_scan (root : string) -> int
extern storage_entry_path (index : int) -> string

report (root, limit) =
    handle Disk::Facts
        volume path => Volume(total = storage_volume path "total", available = storage_volume path "available", important = storage_volume path "important")
        largest folder => entries folder
    Disk::render (Disk::survey root limit)
```

[`OspreyStorage/StorageServices.swift`](OspreyStorage/StorageServices.swift) implements each import with `@_cdecl`, and [`main.swift`](OspreyStorage/main.swift) calls the generated `osprey_storage_report`. The [domain suite](storage/test/main.ospml) runs the same survey and report under a fixed handler, natively and without Swift.

## Where the numbers come from

- **Volume capacity:** Foundation's `volumeTotalCapacity`, `volumeAvailableCapacity` and `volumeAvailableCapacityForImportantUsage`. The last includes space macOS can purge on demand, such as caches and local snapshots, which is what System Settings shows as available.
- **Folder sizes:** the Spotlight index, queried with `MDQuery`. Spotlight already records each indexed item's size, `kMDItemPhysicalSize` (bytes on disk, after APFS compression) when present and `kMDItemFSSize` otherwise. An app or document bundle is one item carrying its whole size, and Spotlight does not index inside packages, so nothing is counted twice. One query sums every indexed item into the folder's children and their children, which also answers the drill-down.

macOS has no public API that returns a folder's total size. APFS can keep directory statistics, but `getattrlist` returns `ATTR_DIR_ALLOCSIZE` as 0 for ordinary folders. The Storage pane in System Settings sizes folders through the private SpaceAttribution framework, which apps cannot use. Apple's [2 October 2026 notice](https://developer.apple.com/news/?id=p6zjojqw) tightens Full Disk Access without adding a storage API, so the Spotlight index is the public record to build on.

Spotlight only knows what it indexes. It skips caches such as `~/Library/Caches`, hidden folders, locations excluded in Spotlight's privacy settings, and volumes with indexing off. The report states how much of the volume's used space it could attribute. Check indexing with `mdutil -s /`.

## Build and run

Use an Apple silicon Mac with Xcode selected by `xcode-select`. From the repository root:

```sh
make macos                          # Mac runtime, compiler and this app
examples/macos/run.sh               # report on your home folder
examples/macos/run.sh /Applications 5
make macos-test                     # C ABI, corpus, domain suite and smoke
```

A report on a home folder with several hundred thousand indexed items takes tens of seconds while Spotlight answers. `OSPREY_MACOS_SKIP_BUILD=1` reuses the existing build, and `OSPREY_BIN` selects another compiler.

## One boundary for Mac and iPhone

`osprey storage --compile --target=macos -o build/libOspreyStorage.a` emits the archive and `libOspreyStorage.h`, which [`Bridge.h`](OspreyStorage/Bridge.h) imports into Swift. This is the [iPhone target](../ios) built for the Mac: the same header, scalar types, entry, thunks and target rules. A Swift `@_cdecl` function written for one platform links against the other unchanged; only the Apple API inside it has to exist on both. See [the Mac slice](../../docs/specs/compiler/iOSCodegen.md#the-mac-slice-ios-target-macos) for the contract.

Calls are synchronous on one host thread. Swift strings passed to Osprey are borrowed for the call; Swift copies strings Osprey returns with `String(cString:)`. A string a Swift import returns must stay valid for the rest of the process, so `storage_entry_path` hands over a `strdup` copy that Swift never frees. Calling Swift crosses Osprey's memory-safety boundary.
