import Foundation

// Implements [IOS-TARGET-ENTRY] and [IOS-TARGET-MACOS]. Swift owns the process
// and the platform APIs; Osprey owns the report. Usage:
//   osprey-storage [folder] [limit]
//   osprey-storage --smoke [--require-index]
let arguments = Array(CommandLine.arguments.dropFirst())
let status = osprey_main()
guard status == 0 else {
    FileHandle.standardError.write(Data("Osprey initialization failed: \(status)\n".utf8))
    exit(1)
}
if arguments.first == "--smoke" {
    exit(SmokeCheck.run(requireIndex: arguments.contains("--require-index")) ? 0 : 1)
}
let root = arguments.first ?? NSHomeDirectory()
let limit = arguments.dropFirst().first.flatMap { Int64($0) } ?? 8
FileHandle.standardError.write(Data("Reading Spotlight's index for \(root)…\n".utf8))
print(ospreyString(osprey_storage_report(root, limit)))

/// Copies a string Osprey returned before the next call into the library.
func ospreyString(_ pointer: UnsafePointer<CChar>?) -> String {
    pointer.map { String(cString: $0) } ?? "Osprey returned a null string"
}
