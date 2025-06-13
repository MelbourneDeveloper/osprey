#!/bin/bash

# Osprey Homebrew Tap Setup Script
# Based on: https://docs.brew.sh/How-to-Create-and-Maintain-a-Tap

set -e

echo "🍺 Setting up Osprey Homebrew Tap..."

# Configuration
TAP_NAME="homebrew-osprey"
GITHUB_USER="melbournedeveloper"  # Change this to your GitHub username
REPO_URL="https://github.com/${GITHUB_USER}/${TAP_NAME}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_step() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if we have a release file first
if [ ! -f "../compiler/osprey" ]; then
    print_error "No osprey binary found in ../compiler/osprey"
    print_warning "Please build your project first with: make build"
    exit 1
fi

# Step 1: Create GitHub repository (manual step)
echo "📋 Manual Setup Required:"
echo "1. Go to GitHub and create a repository named: ${TAP_NAME}"
echo "2. Make it public"
echo "3. Don't initialize with README (we'll add our own)"
echo "4. Press Enter when done..."
read -p ""

# Step 2: Create local tap structure
print_step "Creating local tap directory structure..."

# Clean up any existing directory
if [ -d "${TAP_NAME}" ]; then
    rm -rf "${TAP_NAME}"
fi

# Create the tap directory structure
mkdir -p "${TAP_NAME}/Formula"
cd "${TAP_NAME}"

# Step 3: Initialize Git repository
print_step "Initializing Git repository..."
git init
git branch -M main

# Step 4: Create README
print_step "Creating README.md..."
cat > README.md << EOF
# Osprey Homebrew Tap

Personal Homebrew tap for the Osprey programming language.

## Installation

\`\`\`bash
# Add the tap
brew tap ${GITHUB_USER}/osprey

# Install Osprey
brew install osprey
\`\`\`

## Direct Installation

\`\`\`bash
# Install directly without adding tap first
brew install ${GITHUB_USER}/osprey/osprey
\`\`\`

## About Osprey

Osprey is a modern functional programming language designed for clarity, safety, and expressiveness.

- Homepage: https://www.ospreylang.dev
- Documentation: https://www.ospreylang.dev/docs
- Source: https://github.com/melbournedeveloper/osprey

## Issues

If you have issues with this formula, please report them at the [main Osprey repository](https://github.com/melbournedeveloper/osprey/issues).
EOF

# Step 5: Copy and update the formula
print_step "Creating Osprey formula..."
cp ../osprey.rb Formula/osprey.rb

# Step 6: Create release automation script
print_step "Creating release automation..."
cat > update-formula.sh << 'EOF'
#!/bin/bash

# Script to update the Osprey formula with new releases
# Usage: ./update-formula.sh <version> <sha256>

set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <version> <sha256>"
    echo "Example: $0 1.0.0 abc123..."
    exit 1
fi

VERSION=$1
SHA256=$2

print_step() {
    echo -e "\033[0;32m✓\033[0m $1"
}

print_step "Updating Osprey formula to version ${VERSION}..."

# Update the formula
sed -i.bak \
    -e "s/version \".*\"/version \"${VERSION}\"/" \
    -e "s/sha256 \".*\"/sha256 \"${SHA256}\"/" \
    -e "s/v[0-9]\+\.[0-9]\+\.[0-9]\+/v${VERSION}/g" \
    Formula/osprey.rb

# Remove backup file
rm Formula/osprey.rb.bak

print_step "Updated formula:"
grep -E "(version|sha256|url)" Formula/osprey.rb

print_step "Committing changes..."
git add Formula/osprey.rb
git commit -m "osprey ${VERSION}"
git push origin main

print_step "Formula updated successfully!"
echo "Users can now install with: brew upgrade osprey"
EOF

chmod +x update-formula.sh

# Step 7: Create GitHub Actions for CI (optional but recommended)
print_step "Creating GitHub Actions workflow..."
mkdir -p .github/workflows

cat > .github/workflows/test.yml << 'EOF'
name: Test Formula

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: macos-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Homebrew
      id: set-up-homebrew
      uses: Homebrew/actions/setup-homebrew@master
    
    - name: Test formula
      run: |
        brew test-bot --only-formulae --only-json-tab --skip-dependents Formula/osprey.rb
    
    - name: Audit formula
      run: |
        brew audit --strict --online Formula/osprey.rb
EOF

# Step 8: Create .gitignore
print_step "Creating .gitignore..."
cat > .gitignore << 'EOF'
# macOS
.DS_Store

# Homebrew
*.bottle.tar.gz
EOF

# Step 9: Add everything to git
print_step "Adding files to Git..."
git add .
git commit -m "Initial tap setup for Osprey"

# Step 10: Add remote and push
print_step "Setting up remote repository..."
git remote add origin "${REPO_URL}.git"

echo ""
print_warning "About to push to ${REPO_URL}"
echo "Make sure the GitHub repository exists and is empty!"
read -p "Press Enter to continue..."

git push -u origin main

# Step 11: Create installation instructions
print_step "Creating installation instructions..."
cd ..

cat > INSTALLATION.md << EOF
# Osprey Installation Instructions

## From Homebrew Tap

### Method 1: Add tap first, then install
\`\`\`bash
brew tap ${GITHUB_USER}/osprey
brew install osprey
\`\`\`

### Method 2: Direct install
\`\`\`bash
brew install ${GITHUB_USER}/osprey/osprey
\`\`\`

## Verify Installation

\`\`\`bash
osprey --version
osprey --help
\`\`\`

## Update Osprey

\`\`\`bash
brew update
brew upgrade osprey
\`\`\`

## Uninstall

\`\`\`bash
brew uninstall osprey
brew untap ${GITHUB_USER}/osprey  # Optional: remove the tap
\`\`\`

## Tap Repository

The tap is hosted at: ${REPO_URL}

## Issues

Report issues at: https://github.com/melbournedeveloper/osprey/issues
EOF

echo ""
echo "🎉 Tap setup complete!"
echo ""
print_step "Tap repository: ${REPO_URL}"
print_step "Installation command: brew install ${GITHUB_USER}/osprey/osprey"
echo ""
echo "📝 Next steps:"
echo "1. Users can now install Osprey with the commands above"
echo "2. When you release a new version, use: ./${TAP_NAME}/update-formula.sh <version> <sha256>"
echo "3. Check the INSTALLATION.md file for user instructions"
echo ""
print_warning "Remember to create actual GitHub releases with binaries for the URLs in your formula!" 