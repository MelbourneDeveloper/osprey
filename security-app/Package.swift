// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "Vigil",
    platforms: [.macOS(.v14), .iOS(.v17)],
    products: [
        .library(name: "VigilDesign", targets: ["VigilDesign"]),
        .executable(name: "Vigil", targets: ["VigilApp"])
    ],
    targets: [
        .target(name: "VigilCore"),
        .testTarget(name: "VigilCoreTests", dependencies: ["VigilCore"]),
        .target(name: "VigilDesign", resources: [.process("Resources")]),
        .executableTarget(name: "VigilApp", dependencies: ["VigilDesign", "VigilCore"]),
        .testTarget(name: "VigilDesignTests", dependencies: ["VigilDesign"])
    ]
)
