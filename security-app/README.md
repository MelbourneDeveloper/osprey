# Vigil design system

Native visual foundation for the security product. **Vigil is a working identity**, pending name and trademark review. This package delivers the SwiftUI design adapter and an interactive design studio with synthetic fixtures. It does not implement monitoring, detection, the Osprey application framework, or the NIMBLESITE agent.

```sh
swift run --package-path security-app VigilGallery
swift test --package-path security-app
```

Open `Package.swift` in Xcode to consume `VigilDesign` from a native macOS or iOS host. Minimum library deployment: macOS 14 / iOS 17. The gallery executable is intended for Mac; iOS library compilation is checked independently.

## Included

- Original generated aperture icon and system-font wordmark.
- Typed semantic colours, spacing, system typography and motion policy.
- Reusable surfaces, coverage/severity badges, evidence fields, coverage strip and native button styling.
- Overview with search, selection, evidence inspector and explanation sheet.
- Interactive component catalog, light/dark switching and Reduce Motion preview.
- Offline native rendering for desktop and compact layout inspection.

The design contract, component behaviour, motion timing and acceptance boundaries are in [the design specification](docs/design-system.md). `Sources/VigilDesign/Theme.swift` is the current executable token source. Do not create a second handwritten token palette for another backend; extract a shared token source when the Osprey framework is implemented.

## Native previews

```sh
swift run --package-path security-app VigilGallery --render security-app/.build/previews
```

This renders the same content through `NSHostingView`, including native fields, into local PNG files. It does not capture the user's screen, collect device data, or replace real VoiceOver/keyboard/device testing. Previews omit the window chrome/sidebar; the running gallery uses native `NavigationSplitView`.

## Verification

- macOS gallery builds; semantic text and primary-button contrast tests pass in both appearances.
- iOS library build succeeds with `xcodebuild -scheme VigilDesign -destination 'generic/platform=iOS' CODE_SIGNING_ALLOWED=NO build` from this directory.
- Native desktop light/dark and 390-point compact content renderings were inspected.
- Live keyboard, VoiceOver, animation and device interaction checks remain pending: the Mac was locked during this session. No App Store readiness or full accessibility certification is claimed.

## Asset provenance

`Sources/VigilDesign/Resources/Brand.xcassets/VigilIcon.imageset/VigilIcon.png` was created with the built-in Image Generation tool, not fetched from an icon vendor. The exact prompt is in [logo-prompt.txt](logo-prompt.txt). Direction: “A professional anti-spyware macOS app icon: full-bleed deep graphite #172521, a large ivory aperture made from two offset angular brackets around a small muted mint square, precise optical balance, restrained bevel, no text, shield, padlock or eye illustration.” The packaged raster is the source artwork; final Icon Composer/vector, small-size optical variants and store export production remain release work. SF Symbols supply interface icons on Apple platforms and must not be exported as the Android icon library.
