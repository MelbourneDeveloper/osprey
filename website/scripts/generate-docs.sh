#!/bin/bash

# Generate reference documentation from Osprey compiler examples
# This script extracts symbols from example files and generates markdown documentation

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEBSITE_DIR="$(dirname "$SCRIPT_DIR")"
COMPILER_DIR="$WEBSITE_DIR/../compiler"
DOCS_DIR="$WEBSITE_DIR/src/docs"

echo "Generating Osprey reference documentation..."

# Always rebuild the compiler for local platform
echo "Building Osprey compiler for local platform..."
cd "$COMPILER_DIR"
make clean
make build

# Create reference documentation directory
mkdir -p "$DOCS_DIR"

# Generate documentation from key example files
cd "$COMPILER_DIR"

echo "Generating standard library documentation..."

# Create API reference documentation
cat > "$DOCS_DIR/stdlib.md" << 'EOF'
---
title: "API Reference - Osprey Programming Language"
description: "Complete API reference for built-in functions, types, operators, and language constructs"
---

EOF

# Generate API reference documentation from the compiler
echo "Generating API reference from compiler..."
"$COMPILER_DIR/bin/osprey" --docs --docs-dir "$DOCS_DIR"

# The compiler now generates docs directly to $DOCS_DIR, so no copying needed
if [ -f "$DOCS_DIR/index.md" ]; then
    # Create a main API reference page that includes the generated content
    cat > "$DOCS_DIR/stdlib.md" << 'EOF'
---
title: "API Reference - Osprey Programming Language"
description: "Complete API reference for built-in functions, types, operators, and language constructs"
---

EOF
    
    # Append the generated README content if it exists
    if [ -f "$DOCS_DIR/README.md" ]; then
        tail -n +1 "$DOCS_DIR/README.md" >> "$DOCS_DIR/stdlib.md"
    fi
    
    echo "Documentation generation complete!"
    echo "Generated files:"
    echo "  - $DOCS_DIR/stdlib.md (Main API Reference)"
    echo "  - $DOCS_DIR/functions/ (Individual function docs)"
    echo "  - $DOCS_DIR/types/ (Individual type docs)"
    echo "  - $DOCS_DIR/operators/ (Individual operator docs)"
    echo "  - $DOCS_DIR/keywords/ (Individual keyword docs)"
else
    echo "Error: Documentation generation failed - no docs generated to $DOCS_DIR"
    exit 1
fi 