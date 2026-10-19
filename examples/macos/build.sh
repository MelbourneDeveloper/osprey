#!/usr/bin/env bash
# Implements [IOS-TARGET-MACOS]: compile the Osprey library for the Mac slice,
# then link it into the Swift host through the generated header.
set -euo pipefail

example_dir=$(cd "$(dirname "$0")" && pwd)
repo_dir=$(cd "$example_dir/../.." && pwd)
compiler=${OSPREY_BIN:-"$repo_dir/target/release/osprey"}
build="$example_dir/build"
if [[ ! -x "$compiler" ]]; then
    echo "Osprey compiler not found at $compiler. Run make macos from the repository root." >&2
    exit 1
fi
mkdir -p "$build"
"$compiler" "$example_dir/storage" --compile --target=macos -o "$build/libOspreyStorage.a"
xcrun --sdk macosx swiftc -O -target arm64-apple-macos12.0 \
    -import-objc-header "$example_dir/OspreyStorage/Bridge.h" -I "$build" \
    "$example_dir"/OspreyStorage/*.swift -L "$build" -lOspreyStorage \
    -o "$build/osprey-storage"
echo "Built $build/osprey-storage"
