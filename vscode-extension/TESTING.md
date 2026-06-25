# Testing Guide for Osprey VS Code Extension

This document explains how to run and develop tests for the Osprey VS Code extension.

## Test Environment

The extension uses the VS Code Extension Testing framework which runs tests inside a special VS Code instance called the **Extension Development Host**. This provides full access to the VS Code API.

## Quick Start

```bash
# Run all tests
npm test

# Compile and run tests
npm run compile && npm test

# Use the test launcher for more options
npm run test:launcher help
```

## Test Structure

```
test/
├── runTest.ts              # Test runner entry point
├── suite/
│   ├── index.ts           # Test suite loader
│   └── extension.test.ts  # All tests (integration + language features)
└── mocha.opts            # Mocha configuration
```

## Test Categories

### Integration Tests

- Tests that require VS Code API and extension functionality
- Located in `test/suite/extension.test.ts`
- Test extension activation, commands, language features, etc.
- Includes two test suites:
  - **Osprey Extension Integration Tests**: Core extension functionality (working)
  - **Osprey Language Features Tests**: Advanced LSP features (some expected to fail until LSP issues are resolved)

## Running Tests

### Method 1: VS Code Test CLI (Recommended)

```bash
# Install dependencies first
npm install

# Run all tests
npm test

# The tests will:
# 1. Download VS Code if not present
# 2. Launch VS Code Extension Host
# 3. Load the extension
# 4. Run tests with full VS Code API access
```

### Method 2: Test Launcher Script

```bash
# Show all available commands
npm run test:launcher help

# Run tests
npm run test:launcher test

# Compile TypeScript
npm run test:launcher compile

# Clean output
npm run test:launcher clean
```

### Method 3: VS Code Debugger

1. Open the extension in VS Code
2. Go to Run and Debug (Ctrl+Shift+D)
3. Select "Extension Tests" configuration
4. Press F5 to run tests with debugging

## Test Configuration

### VS Code Test CLI Configuration (.vscode-test.js)

```javascript
const { defineConfig } = require("@vscode/test-cli");

module.exports = defineConfig({
  files: "out/test/suite/**/*.test.js",
  version: "stable",
  mocha: {
    ui: "tdd",
    timeout: 10000,
    color: true,
  },
  launchArgs: ["--disable-extensions", "--disable-workspace-trust"],
});
```

### Debug Configuration (.vscode/launch.json)

Two configurations are available:

- **Run Extension**: Launch extension for manual testing
- **Extension Tests**: Run tests with debugging support

## Writing Tests

### Basic Test Structure

```typescript
import * as assert from "assert";
import * as vscode from "vscode";

suite("My Test Suite", () => {
  test("My test", async () => {
    // Test code here
    assert.strictEqual(1 + 1, 2);
  });
});
```

### Testing Extension Functionality

```typescript
test("Extension should activate", async () => {
  const extension = vscode.extensions.getExtension(
    "christianfindlay.osprey-language-support",
  );
  assert.ok(extension);
  assert.ok(extension.isActive);
});
```

### Testing with Files

```typescript
test("Language detection", async () => {
  const tempFile = path.join(os.tmpdir(), "test.osp");
  fs.writeFileSync(tempFile, "fn test() = 42");

  const document = await vscode.workspace.openTextDocument(tempFile);
  assert.strictEqual(document.languageId, "osprey");

  // Cleanup
  fs.unlinkSync(tempFile);
});
```

### Testing Language Features (Expected Failures)

Some language feature tests are expected to fail until LSP integration issues are resolved:

```typescript
test("Go to Definition - Function (Expected to fail until LSP fixed)", async () => {
  // Test implementation with try/catch to handle expected failures
  try {
    // Test go to definition functionality
  } catch (error) {
    console.log("Go to Definition failed as expected:", error);
  }
});
```

## Test Status

### ✅ Working Tests

- Extension activation
- Language detection for `.osp` files
- Command availability
- Syntax highlighting
- Configuration access
- File operations
- Multiple file handling
- Basic diagnostics

### ⚠️ Expected Failures (LSP Issues)

- Go to Definition
- Find All References
- Advanced hover information
- Document symbols
- Some signature help features

## Common Issues

### "Cannot find module 'vscode'" Error

This means tests are running in regular Node.js instead of VS Code Extension Host.

- ✅ **Fixed**: Use `vscode-test` CLI or proper test runner
- ❌ **Wrong**: Running tests with regular `mocha` command

### Tests Timing Out

- Increase timeout in `.vscode-test.js` or test files
- Use `await` for async operations
- Add delays for VS Code operations: `await new Promise(resolve => setTimeout(resolve, 1000))`

### Extension Not Activating

- Check that test files have `.osp` extension
- Ensure extension is properly configured in `package.json`
- Use `--disable-extensions` flag to avoid conflicts

## Test Output

Successful test run should show:

```
✔ Validated version: 1.100.2
✔ Found existing install in .vscode-test/vscode-darwin-arm64-1.100.2
Loading development extension at /path/to/extension
✔ Extension should activate when opening .osp file
✔ Language should be set to osprey for .osp files
...
Integration Tests: 10 passing
Language Feature Tests: 3 passing, 3 expected failures
```

## Continuous Integration

For CI/CD pipelines, use:

```bash
# Install dependencies
npm ci

# Run tests in headless mode
npm test
```

The VS Code Test CLI automatically handles downloading and running VS Code in CI environments.

## References

- [VS Code Extension Testing Guide](https://code.visualstudio.com/api/working-with-extensions/testing-extension)
- [@vscode/test-cli Documentation](https://www.npmjs.com/package/@vscode/test-cli)
- [@vscode/test-electron Documentation](https://www.npmjs.com/package/@vscode/test-electron)
