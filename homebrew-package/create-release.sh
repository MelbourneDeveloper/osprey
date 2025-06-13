#!/bin/bash
set -e

VERSION=${1:-"0.1.0"}
echo "Creating release v$VERSION..."

# Build runtime libraries
echo "✓ Building runtime libraries..."
cd ../compiler/runtime
gcc -c fiber_runtime.c -o fiber_runtime.o
ar rcs libfiber_runtime.a fiber_runtime.o

gcc -c http_server_runtime.c -o http_runtime.o
gcc -c http_shared.c -o http_shared.o
ar rcs libhttp_runtime.a http_runtime.o http_shared.o

# Create release package
echo "✓ Creating release package..."
cd ..
mkdir -p release-package
cp osprey release-package/
cp runtime/libfiber_runtime.a release-package/
cp runtime/libhttp_runtime.a release-package/

# Create tarball
cd release-package
tar -czf "../osprey-darwin-amd64.tar.gz" *
cd ..

# Calculate SHA256
SHA256=$(sha256sum osprey-darwin-amd64.tar.gz | cut -d' ' -f1)
echo "✓ SHA256: $SHA256"

# Create GitHub release
echo "✓ Creating GitHub release..."
gh release create "v$VERSION" \
    --repo "MelbourneDeveloper/osprey" \
    --title "Osprey v$VERSION" \
    --notes "Release v$VERSION of the Osprey programming language" \
    osprey-darwin-amd64.tar.gz

# Update homebrew formula
echo "✓ Updating homebrew formula..."
cd ../homebrew-package/homebrew-osprey
./update-formula.sh "$VERSION" "$SHA256"

echo "✓ Release v$VERSION created successfully!"
echo "✓ Users can now install with: brew install melbournedeveloper/osprey/osprey"

# Cleanup
cd ../../compiler
rm -rf release-package osprey-darwin-amd64.tar.gz
rm runtime/*.o runtime/*.a 