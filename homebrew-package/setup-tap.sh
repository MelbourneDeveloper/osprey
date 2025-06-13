#!/bin/bash

set -e

echo "🍺 Setting up Osprey Homebrew Tap..."

# Check dependencies
if [ ! -f "../compiler/osprey" ]; then
    echo "✗ No osprey binary found in ../compiler/osprey"
    echo "Build first: make build"
    exit 1
fi

# Step 2: Create local tap structure
echo "✓ Creating local tap structure..."

# Clean up any existing directory
if [ -d "homebrew-osprey" ]; then
    rm -rf "homebrew-osprey"
fi

# Create the tap directory structure
mkdir -p "homebrew-osprey/Formula"
cd "homebrew-osprey"

# Step 3: Initialize Git repository
echo "✓ Initializing Git repository..."
git init
git branch -M main

# Step 4: Create README
echo "✓ Creating README.md..."
cat > README.md << 'EOF'
# Osprey Homebrew Tap

Homebrew tap for the Osprey programming language.

## Installation

```bash
# Add the tap
brew tap melbournedeveloper/osprey

# Install Osprey
brew install osprey
```

## Direct Installation

```bash
# Install directly without adding tap first
brew install melbournedeveloper/osprey/osprey
```

## About Osprey

Osprey is a modern functional programming language designed for clarity, safety, and expressiveness.

- Homepage: https://www.ospreylang.dev
- Documentation: https://www.ospreylang.dev/docs
- Source: https://github.com/melbournedeveloper/osprey

## Issues

Report issues at the [main Osprey repository](https://github.com/melbournedeveloper/osprey/issues).
EOF

# Step 5: Copy formula
echo "✓ Creating Osprey formula..."
cp ../osprey.rb Formula/osprey.rb

# Step 6: Create release automation script
echo "✓ Creating release automation..."
cat > update-formula.sh << 'EOF'
#!/bin/bash
set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <version> <sha256>"
    exit 1
fi

echo "✓ Updating Osprey formula to version $1..."

sed -i.bak \
    -e "s/version \".*\"/version \"$1\"/" \
    -e "s/sha256 \".*\"/sha256 \"$2\"/" \
    -e "s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v$1/g" \
    Formula/osprey.rb

rm Formula/osprey.rb.bak

echo "✓ Committing changes..."
git add Formula/osprey.rb
git commit -m "osprey $1"
git push origin main

echo "✓ Updated successfully!"
EOF

chmod +x update-formula.sh

# Step 7: Create .gitignore
echo "✓ Creating .gitignore..."
cat > .gitignore << 'EOF'
.DS_Store
*.bottle.tar.gz
EOF

# Step 8: Add everything to git
echo "✓ Adding files to Git..."
git add .
git commit -m "Initial tap setup for Osprey"

# Step 9: Add remote and push
echo "✓ Pushing to remote repository..."
git remote add origin "https://github.com/melbournedeveloper/homebrew-osprey.git"
git push -u origin main

echo ""
echo "🎉 Tap setup complete!"
echo "✓ Repository: https://github.com/melbournedeveloper/homebrew-osprey"
echo "✓ Install with: brew install melbournedeveloper/osprey/osprey"
echo "✓ Update releases with: ./homebrew-osprey/update-formula.sh <version> <sha256>" 