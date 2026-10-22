// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "VigilDesign",
    platforms: [.macOS(.v14), .iOS(.v17)],
    products: [
        .library(name: "VigilDesign", targets: ["VigilDesign"]),
        .executable(name: "VigilGallery", targets: ["VigilGallery"])
    ],
    targets: [
        .target(name: "VigilDesign", resources: [.process("Resources")]),
        .executableTarget(name: "VigilGallery", dependencies: ["VigilDesign"]),
        .testTarget(name: "VigilDesignTests", dependencies: ["VigilDesign"])
    ]
)
