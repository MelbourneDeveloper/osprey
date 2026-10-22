#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
compiler="${OSPREY_COMPILER:-$(command -v osprey || true)}"
if [[ -z "$compiler" ]]; then
    echo 'Set OSPREY_COMPILER to the Osprey compiler executable.' >&2
    exit 1
fi
mkdir -p .build/engine
"$compiler" engine --compile --memory=arc -o .build/engine/vigil-engine
swift build
products="$(swift build --show-bin-path)"
cp .build/engine/vigil-engine "$products/vigil-engine"
app="$PWD/.build/Vigil.app"
mkdir -p "$app/Contents/MacOS" "$app/Contents/Resources"
copy_executable() {
    cp "$1" "$2.new"
    mv -f "$2.new" "$2"
}
copy_executable "$products/Vigil" "$app/Contents/MacOS/Vigil"
copy_executable .build/engine/vigil-engine "$app/Contents/MacOS/vigil-engine"
cp -R "$products/Vigil_VigilDesign.bundle" "$app/Contents/Resources/"
cat > "$app/Contents/Info.plist" <<'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
<key>CFBundleIdentifier</key><string>local.vigil.security</string>
<key>CFBundleName</key><string>Vigil</string>
<key>CFBundleExecutable</key><string>Vigil</string>
<key>CFBundlePackageType</key><string>APPL</string>
<key>CFBundleShortVersionString</key><string>0.1.0</string>
<key>CFBundleVersion</key><string>1</string>
<key>LSMinimumSystemVersion</key><string>14.0</string>
<key>NSHighResolutionCapable</key><true/>
</dict></plist>
PLIST
# Seal the assembled development bundle; this is not Developer ID signing or notarization.
/usr/bin/codesign --force --sign - "$app/Contents/MacOS/vigil-engine"
/usr/bin/codesign --force --sign - "$app"
/usr/bin/codesign --verify --deep --strict "$app"
printf 'Built %s\n' "$app"
