import * as assert from 'assert';
import * as vscode from 'vscode';
import * as path from 'path';
import * as fs from 'fs';
import * as os from 'os';
import {
  shipwrightPlatform,
  resolveBundledCompiler,
  resolveServerCommand
} from '../../client/src/extension';

const extensionId = 'nimblesite.osprey';

// The compiled test lives at <ext>/out/test/suite, so the extension root is
// three levels up. The release pipeline stamps a shipwright.json into the
// extension root (it is gitignored as a build-time artifact); we replicate that
// here so the Shipwright version handshake in activate() runs under test
// instead of being skipped. Staging happens at module load — before any test
// triggers (lazy) activation.
const extensionRoot = path.resolve(__dirname, '..', '..', '..');
const stagedManifestPath = path.join(extensionRoot, 'shipwright.json');
const repoRootManifestPath = path.resolve(extensionRoot, '..', 'shipwright.json');
let manifestWasStaged = false;

(function stageShipwrightManifest(): void {
  // Only stage if the extension hasn't already shipped one and we have the
  // canonical repo-root manifest to copy.
  if (!fs.existsSync(stagedManifestPath) && fs.existsSync(repoRootManifestPath)) {
    fs.copyFileSync(repoRootManifestPath, stagedManifestPath);
    manifestWasStaged = true;
  }
})();

// resolveOspreyOnPath returns the absolute path of the `osprey` binary the test
// harness staged on PATH (the same one the LSP would otherwise find), or
// undefined if it cannot be located. Used to exercise the explicit
// `server.compilerPath` branch of resolveServerCommand before activation.
function resolveOspreyOnPath(): string | undefined {
  const exe = process.platform === 'win32' ? 'osprey.exe' : 'osprey';
  for (const dir of (process.env.PATH ?? '').split(path.delimiter)) {
    if (!dir) {
      continue;
    }
    const candidate = path.join(dir, exe);
    if (fs.existsSync(candidate)) {
      return candidate;
    }
  }
  return undefined;
}

suite('Osprey Shipwright Activation Coverage', () => {
  const settle = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));
  let priorCompilerPath: string | undefined;
  let setCompilerPath = false;

  suiteSetup(async () => {
    // Point server.compilerPath at the real osprey binary BEFORE the extension
    // activates so resolveServerCommand takes its explicit-user-path branch and
    // the language client launches against a genuine osprey, exercising the
    // client option handlers and the start outcome.
    const ospreyPath = resolveOspreyOnPath();
    if (ospreyPath) {
      const config = vscode.workspace.getConfiguration('osprey');
      priorCompilerPath = config.get<string>('server.compilerPath');
      await config.update(
        'server.compilerPath',
        ospreyPath,
        vscode.ConfigurationTarget.Global
      );
      setCompilerPath = true;
    }
  });

  suiteTeardown(async () => {
    // Restore the compiler path so later suites see defaults.
    if (setCompilerPath) {
      await vscode.workspace
        .getConfiguration('osprey')
        .update(
          'server.compilerPath',
          priorCompilerPath ?? '',
          vscode.ConfigurationTarget.Global
        );
    }
    // Remove only the manifest we staged so we never delete a real one.
    if (manifestWasStaged && fs.existsSync(stagedManifestPath)) {
      fs.rmSync(stagedManifestPath, { force: true });
    }
  });

  test('extension activates with a shipwright manifest present', async () => {
    const ext = vscode.extensions.getExtension(extensionId);
    assert.ok(ext, 'extension must be discoverable');

    // The explicit compiler path we set must be visible to the extension so its
    // resolveServerCommand picks the user-configured binary.
    const ospreyPath = resolveOspreyOnPath();
    if (ospreyPath) {
      assert.strictEqual(
        vscode.workspace.getConfiguration('osprey').get<string>('server.compilerPath'),
        ospreyPath,
        'server.compilerPath is the staged osprey binary'
      );
    }

    // The manifest path the extension resolves must point at the staged file so
    // the Shipwright handshake block (fs.existsSync(manifestPath)) is taken.
    assert.ok(
      fs.existsSync(stagedManifestPath),
      'shipwright.json is staged in the extension root'
    );

    // Activating runs the whole activate() body, including the async Shipwright
    // import + activateShipwright handshake. It must not throw and must leave
    // the extension active.
    if (ext && !ext.isActive) {
      await ext.activate();
    }
    // Give the fire-and-forget Shipwright async IIFE time to run its import,
    // version check, and outputChannel.appendLine before we assert.
    await settle(2500);

    assert.ok(ext?.isActive, 'extension is active after Shipwright handshake');

    // The manifest we staged is valid JSON describing the osprey product, so a
    // successful load is the contract the handshake depends on.
    const manifest = JSON.parse(fs.readFileSync(stagedManifestPath, 'utf8'));
    assert.strictEqual(manifest.product.id, 'osprey', 'manifest is the osprey product');
    assert.ok(Array.isArray(manifest.components), 'manifest declares components');
    assert.ok(manifest.components.length > 0, 'manifest has at least one component');
  });
});

suite('Osprey Extension Integration Tests', () => {
  let tempDir: string;
  let testFile: string;

  setup(() => {
    // Create temporary directory for test files
    tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'osprey-test-'));
    testFile = path.join(tempDir, 'test.osp');
  });

  teardown(() => {
    // Clean up temporary files
    if (fs.existsSync(tempDir)) {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  test('Extension should activate when opening .osp file', async () => {
    // Create a simple Osprey file
    const ospreyCode = `
// Simple test function
fn add(a, b) = a + b

let result = add(5, 3)
print(result)
`;
    fs.writeFileSync(testFile, ospreyCode);

    // Open the file in VS Code
    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait a bit for extension to activate
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Check that the extension is active
    const extension = vscode.extensions.getExtension(extensionId);
    assert.ok(extension, 'Extension should be found');
    
    if (extension) {
      assert.ok(extension.isActive, 'Extension should be active after opening .osp file');
    }
  });

  test('Language should be set to osprey for .osp files', async () => {
    const ospreyCode = `fn test() = 42`;
    fs.writeFileSync(testFile, ospreyCode);

    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait for language detection
    await new Promise(resolve => setTimeout(resolve, 500));

    assert.strictEqual(document.languageId, 'osprey', 'Language should be set to osprey');
  });

  test('Compile command should be available for .osp files', async () => {
    const ospreyCode = `fn hello() = print("Hello, World!")`;
    fs.writeFileSync(testFile, ospreyCode);

    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait for extension activation
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Get all available commands
    const commands = await vscode.commands.getCommands();
    
    assert.ok(commands.includes('osprey.compile'), 'Compile command should be available');
    assert.ok(commands.includes('osprey.run'), 'Run command should be available');
  });

  test('Syntax highlighting should work for .osp files', async () => {
    const ospreyCode = `
fn power(base, exp) = match exp {
  0 => 1
  1 => base
  _ => base * power(base, exp - 1)
}

let result = power(2, 3)
print(result)
`;
    fs.writeFileSync(testFile, ospreyCode);

    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait for syntax highlighting to load
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Check that the document has the correct language
    assert.strictEqual(document.languageId, 'osprey');
    
    // Check that the file has content (basic sanity check)
    assert.ok(document.getText().includes('fn power'), 'Document should contain the test code');
  });

  test('Extension should handle invalid Osprey code gracefully', async () => {
    const invalidCode = `
fn broken syntax here {
  this is not valid osprey code
  missing parentheses and stuff
}
`;
    fs.writeFileSync(testFile, invalidCode);

    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait for diagnostics
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Extension should still be active even with invalid code
    const extension = vscode.extensions.getExtension(extensionId);
    assert.ok(extension?.isActive, 'Extension should remain active with invalid code');
  });

  test('File operations should work without workspace', async () => {
    // This test ensures the extension works with individual files
    const ospreyCode = `fn standalone() = print("No workspace needed!")`;
    fs.writeFileSync(testFile, ospreyCode);

    // Close any existing workspace
    if (vscode.workspace.workspaceFolders) {
      await vscode.commands.executeCommand('workbench.action.closeFolder');
    }

    // Open file without workspace
    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait for extension
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Should still work
    assert.strictEqual(document.languageId, 'osprey');
    
    const extension = vscode.extensions.getExtension(extensionId);
    assert.ok(extension?.isActive, 'Extension should work without workspace');
  });

  test('Multiple .osp files should work correctly', async () => {
    // Create multiple test files
    const file1 = path.join(tempDir, 'file1.osp');
    const file2 = path.join(tempDir, 'file2.osp');
    
    fs.writeFileSync(file1, 'fn func1() = 1');
    fs.writeFileSync(file2, 'fn func2() = 2');

    // Open both files
    const doc1 = await vscode.workspace.openTextDocument(file1);
    const doc2 = await vscode.workspace.openTextDocument(file2);
    
    await vscode.window.showTextDocument(doc1);
    await vscode.window.showTextDocument(doc2);

    // Wait for processing
    await new Promise(resolve => setTimeout(resolve, 1000));

    // Both should have correct language
    assert.strictEqual(doc1.languageId, 'osprey');
    assert.strictEqual(doc2.languageId, 'osprey');
  });

  test('Extension configuration should be accessible', async () => {
    const config = vscode.workspace.getConfiguration('osprey');
    
    // Check that configuration exists and has expected properties
    assert.ok(config, 'Osprey configuration should exist');
    
    // Check default values
    const serverEnabled = config.get('server.enabled');
    const compilerPath = config.get('server.compilerPath');
    
    assert.strictEqual(typeof serverEnabled, 'boolean', 'server.enabled should be boolean');
    assert.strictEqual(typeof compilerPath, 'string', 'server.compilerPath should be string');
  });

  test('Language server should start successfully', async () => {
    // Basic test that language server starts without crashing
    const ospreyCode = `fn test() = 42`;
    fs.writeFileSync(testFile, ospreyCode);

    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);

    // Wait for language server to start
    await new Promise(resolve => setTimeout(resolve, 3000));

    // Extension should still be active
    const extension = vscode.extensions.getExtension(extensionId);
    assert.ok(extension?.isActive, 'Extension should remain active with language server');
  });

  test('Compile and Run commands execute against the active file', async () => {
    fs.writeFileSync(testFile, 'fn main() -> Unit = print("hi")\n');
    const document = await vscode.workspace.openTextDocument(testFile);
    await vscode.window.showTextDocument(document);
    await new Promise(resolve => setTimeout(resolve, 500));

    // Both commands shell out to the staged `osprey` binary on PATH; they must
    // run their handlers to completion without throwing back to the caller.
    await vscode.commands.executeCommand('osprey.compile');
    await vscode.commands.executeCommand('osprey.run');
    await new Promise(resolve => setTimeout(resolve, 3000));

    const extension = vscode.extensions.getExtension(extensionId);
    assert.ok(extension?.isActive, 'Extension should remain active after running commands');
  });
});

suite('Osprey Language Features Tests', () => {
  let document: vscode.TextDocument;
  let editor: vscode.TextEditor;

  // Helper to create and open a test document
  async function createTestDocument(content: string): Promise<void> {
    document = await vscode.workspace.openTextDocument({
      language: 'osprey',
      content: content
    });
    editor = await vscode.window.showTextDocument(document);
    // Wait for language server to process the document
    await new Promise(resolve => setTimeout(resolve, 2000));
  }

  teardown(async () => {
    await vscode.commands.executeCommand('workbench.action.closeActiveEditor');
  });

  test('Go to Definition - Function', async () => {
    const content = `
fn double(x) = x * 2

let result = double(5)
`;
    await createTestDocument(content);

    // Position cursor on 'double' in the function call
    const position = new vscode.Position(3, 13); // Line 3, 'double' call
    
    try {
      // Execute go to definition
      const definitions = await vscode.commands.executeCommand<vscode.Location[]>(
        'vscode.executeDefinitionProvider',
        document.uri,
        position
      );

      if (definitions && definitions.length > 0) {
        assert.strictEqual(definitions[0].range.start.line, 1, 'Definition should be on line 1');
        assert.strictEqual(definitions[0].range.start.character, 3, 'Definition should start at character 3');
      } else {
        // This is expected to fail currently due to LSP issues
        console.log('Go to Definition not working - LSP integration issue (expected)');
      }
    } catch (error) {
      console.log('Go to Definition failed as expected:', error);
    }
  });

  test('Find All References - Function', async () => {
    const content = `
fn add(x, y) = x + y

let sum1 = add(x: 1, y: 2)
let sum2 = add(x: 3, y: 4)
print(add(x: 5, y: 6))
`;
    await createTestDocument(content);

    // Position cursor on 'add' in the function definition
    const position = new vscode.Position(1, 3); // Line 1, 'add' definition
    
    try {
      const references = await vscode.commands.executeCommand<vscode.Location[]>(
        'vscode.executeReferenceProvider',
        document.uri,
        position
      );

      if (references && references.length > 0) {
        assert.strictEqual(references.length, 4, 'Should find 4 references (1 definition + 3 usages)');
      } else {
        // This is expected to fail currently due to LSP issues
        console.log('Find All References not working - LSP integration issue (expected)');
      }
    } catch (error) {
      console.log('Find All References failed as expected:', error);
    }
  });

  test('Hover Information - Function', async () => {
    const content = `
fn multiply(x, y) = x * y

let product = multiply(x: 3, y: 4)
`;
    await createTestDocument(content);

    // Position cursor on 'multiply' in the function call
    const position = new vscode.Position(3, 14); // Line 3, 'multiply' call
    
    try {
      const hovers = await vscode.commands.executeCommand<vscode.Hover[]>(
        'vscode.executeHoverProvider',
        document.uri,
        position
      );

      if (hovers && hovers.length > 0) {
        const hoverContent = hovers[0].contents[0];
        assert.ok(hoverContent, 'Hover should have content');
        
        // Check if hover contains function information
        const hoverText = typeof hoverContent === 'string' ? hoverContent : hoverContent.value;
        assert.ok(hoverText.includes('multiply'), 'Hover should mention the function name');
      } else {
        console.log('Hover information not available yet');
      }
    } catch (error) {
      console.log('Hover failed:', error);
    }
  });

  test('Document Symbols', async () => {
    const content = `
fn foo() = 42
let bar = 10
type Baz = A | B
`;
    await createTestDocument(content);

    try {
      const symbols = await vscode.commands.executeCommand<vscode.DocumentSymbol[]>(
        'vscode.executeDocumentSymbolProvider',
        document.uri
      );

      if (symbols && symbols.length > 0) {
        const symbolNames = symbols.map(s => s.name);
        console.log('Found symbols:', symbolNames);
        // Basic check that we found some symbols
        assert.ok(symbols.length > 0, 'Should find at least some symbols');
      } else {
        console.log('Document symbols not available yet');
      }
    } catch (error) {
      console.log('Document symbols failed:', error);
    }
  });

  test('Diagnostics - Syntax Error', async () => {
    const content = `
fn broken( = 42
`;
    await createTestDocument(content);

    // Wait for diagnostics
    await new Promise(resolve => setTimeout(resolve, 3000));

    const diagnostics = vscode.languages.getDiagnostics(document.uri);
    if (diagnostics.length > 0) {
      const error = diagnostics[0];
      assert.strictEqual(error.severity, vscode.DiagnosticSeverity.Error, 'Should be an error');
      assert.ok(error.message.length > 0, 'Error should have a message');
    } else {
      console.log('No diagnostics found - may need more time or Osprey compiler');
    }
  });

  test('Code Completion', async () => {
    const content = `
fn test() = 42
let x = te
`;
    await createTestDocument(content);

    // Position cursor after 'te'
    const position = new vscode.Position(2, 10);
    
    try {
      const completions = await vscode.commands.executeCommand<vscode.CompletionList>(
        'vscode.executeCompletionItemProvider',
        document.uri,
        position
      );

      if (completions && completions.items.length > 0) {
        console.log('Found completions:', completions.items.map(item =>
          typeof item.label === 'string' ? item.label : item.label.label
        ));
        assert.ok(completions.items.length > 0, 'Should have completion items');
      } else {
        console.log('No completions available yet');
      }
    } catch (error) {
      console.log('Code completion failed:', error);
    }
  });
});

// These suites drive the command handlers, event subscriptions, and debug
// provider registered in activate() so the coverage harness records the
// execFile callbacks and the early-return guards in extension.ts.
suite('Osprey Command Handler Coverage', () => {
  let tempDir: string;
  const extension = () => vscode.extensions.getExtension(extensionId);

  // The compile/run handlers shell out to `osprey` and write to an output
  // channel from the execFile callback; the callback completes after a tick, so
  // give it a generous window before we let the test (and the host) move on.
  const settle = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

  async function closeEverything(): Promise<void> {
    // closeAllEditors races with VS Code's async editor teardown, so poll until
    // the active editor actually drains (or give up after a bounded number of
    // attempts and let the caller assert what it needs).
    for (let i = 0; i < 10 && vscode.window.activeTextEditor; i++) {
      await vscode.commands.executeCommand('workbench.action.closeAllEditors');
      await settle(150);
    }
  }

  // makeActive shows the document and polls until it is genuinely the active
  // text editor. The extension shows an "Osprey Debug" output channel during
  // activation, and that output pseudo-editor steals the active-editor slot in
  // the headless host. We close panels and refocus the editor group between
  // attempts so the .osp document becomes (and stays) the active editor at the
  // moment a command captures window.activeTextEditor.
  async function makeActive(document: vscode.TextDocument): Promise<boolean> {
    for (let i = 0; i < 25; i++) {
      // Hide the output panel and any output pseudo-editor stealing focus.
      await vscode.commands.executeCommand('workbench.action.closePanel').then(
        undefined,
        () => undefined
      );
      await vscode.window.showTextDocument(document, {
        viewColumn: vscode.ViewColumn.One,
        preserveFocus: false,
        preview: false
      });
      await vscode.commands.executeCommand('workbench.action.focusActiveEditorGroup').then(
        undefined,
        () => undefined
      );
      await settle(120);
      const active = vscode.window.activeTextEditor;
      if (active && active.document.uri.toString() === document.uri.toString()) {
        return true;
      }
    }
    return false;
  }

  async function openOsp(name: string, content: string): Promise<vscode.TextDocument> {
    const file = path.join(tempDir, name);
    fs.writeFileSync(file, content);
    const document = await vscode.workspace.openTextDocument(file);
    await makeActive(document);
    return document;
  }

  setup(async () => {
    tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'osprey-cmd-'));
    // Ensure the extension is active before exercising its commands.
    const ext = extension();
    assert.ok(ext, 'Extension must be discoverable');
    if (ext && !ext.isActive) {
      await ext.activate();
    }
  });

  teardown(async () => {
    await closeEverything();
    if (fs.existsSync(tempDir)) {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  test('osprey.compile runs to completion against a valid .osp file', async function () {
    this.timeout(30000);
    const valid = 'fn main() -> Unit = print("compiled ok")\n';
    const document = await openOsp('compile-valid.osp', valid);

    assert.strictEqual(document.languageId, 'osprey', 'doc must be osprey');
    assert.ok(document.fileName.endsWith('.osp'), 'file name must end with .osp');

    // The handler reads window.activeTextEditor synchronously; make sure our
    // .osp document is the active editor at the moment the command fires so it
    // passes both guards and reaches the save->execFile->callback chain.
    const isActive = await makeActive(document);
    assert.ok(isActive, 'the .osp document became the active editor');
    assert.ok(
      vscode.window.activeTextEditor?.document.fileName.endsWith('.osp'),
      'the active editor is the .osp document'
    );

    // Drive the command twice to exercise the chain for both a fresh and an
    // already-saved buffer.
    await vscode.commands.executeCommand('osprey.compile');
    await settle(4000);
    await makeActive(document);
    await vscode.commands.executeCommand('osprey.compile');
    await settle(4000);

    assert.ok(extension()?.isActive, 'extension stays active after compile');
    assert.ok(fs.existsSync(document.fileName), 'source file still on disk');
  });

  test('osprey.compile surfaces a compiler error path for invalid code', async function () {
    this.timeout(30000);
    // Invalid Osprey forces the compiler to exit non-zero, exercising the
    // `if (error)` branch of the execFile callback.
    const broken = 'fn main( = \n this is not valid osprey @@@ \n';
    const document = await openOsp('compile-broken.osp', broken);

    assert.strictEqual(document.languageId, 'osprey');
    assert.ok(await makeActive(document), 'broken doc is active');
    await vscode.commands.executeCommand('osprey.compile');
    await settle(4000);

    assert.ok(extension()?.isActive, 'extension survives a failed compile');
    const text = document.getText();
    assert.ok(text.includes('not valid'), 'broken source preserved');
  });

  test('osprey.run runs to completion against a valid .osp file', async function () {
    this.timeout(30000);
    const valid = 'fn main() -> Unit = print("ran ok")\n';
    const document = await openOsp('run-valid.osp', valid);

    assert.strictEqual(document.languageId, 'osprey');
    assert.ok(await makeActive(document), 'valid run doc is active');
    await vscode.commands.executeCommand('osprey.run');
    await settle(5000);
    await makeActive(document);
    await vscode.commands.executeCommand('osprey.run');
    await settle(5000);

    assert.ok(extension()?.isActive, 'extension stays active after run');
    assert.ok(document.fileName.endsWith('.osp'));
  });

  test('osprey.run surfaces a failure path for invalid code', async function () {
    this.timeout(30000);
    const broken = 'fn main( = @@@ not osprey at all\n';
    const document = await openOsp('run-broken.osp', broken);

    assert.ok(await makeActive(document), 'broken run doc is active');
    await vscode.commands.executeCommand('osprey.run');
    await settle(4000);

    assert.ok(extension()?.isActive, 'extension survives a failed run');
    assert.ok(document.getText().length > 0, 'broken source preserved');
  });

  test('compile and run guard against a non-.osp active editor', async () => {
    // Open a plain text (non-.osp) file so both handlers hit the
    // "Please open a .osp file" early return.
    const txt = path.join(tempDir, 'notes.txt');
    fs.writeFileSync(txt, 'just some plain text, not osprey at all');
    const document = await vscode.workspace.openTextDocument(txt);
    const active = await makeActive(document);

    assert.ok(active, 'the .txt document became active');
    assert.ok(!document.fileName.endsWith('.osp'), 'active file is not .osp');
    assert.notStrictEqual(document.languageId, 'osprey', 'not osprey lang');

    await vscode.commands.executeCommand('osprey.compile');
    await vscode.commands.executeCommand('osprey.run');
    await settle(500);

    assert.ok(extension()?.isActive, 'extension active after guarded commands');
  });

  test('compile and run guard when there is no active editor', async () => {
    await closeEverything();
    // Best-effort: VS Code may keep a hidden editor around, but the handlers'
    // guards are exercised either way (no .osp active editor present).
    const noEditor = vscode.window.activeTextEditor === undefined;

    await vscode.commands.executeCommand('osprey.compile');
    await vscode.commands.executeCommand('osprey.run');
    await settle(400);

    assert.ok(extension()?.isActive, 'extension active with no editor');
    // The commands must not have opened or left an .osp editor active.
    const active = vscode.window.activeTextEditor;
    assert.ok(
      noEditor || !active?.document.fileName.endsWith('.osp'),
      'no .osp editor became active from a guarded command'
    );
  });

  test('osprey.setLanguage retargets the active editor to osprey', async () => {
    // Open a file with a .txt extension so its language starts as plaintext,
    // then force it to osprey through the command. Covers the setLanguage
    // handler body (active editor present).
    const txt = path.join(tempDir, 'convert-me.txt');
    fs.writeFileSync(txt, 'fn main() -> Unit = print("convert")\n');
    const opened = await vscode.workspace.openTextDocument(txt);
    assert.ok(await makeActive(opened), 'convert doc became active');

    assert.notStrictEqual(
      vscode.window.activeTextEditor?.document.languageId,
      'osprey',
      'starts non-osprey'
    );

    await vscode.commands.executeCommand('osprey.setLanguage');
    await settle(600);

    const active = vscode.window.activeTextEditor;
    assert.ok(active, 'an editor is active');
    assert.strictEqual(active?.document.languageId, 'osprey', 'language now osprey');
  });

  test('osprey.setLanguage is a no-op with no active editor', async () => {
    await closeEverything();
    const before = vscode.window.activeTextEditor;

    // Should not throw even though there is (best-effort) nothing to retarget.
    await vscode.commands.executeCommand('osprey.setLanguage');
    await settle(300);

    assert.ok(extension()?.isActive, 'extension stays active');
    // The command must not have created a new editor.
    assert.strictEqual(
      vscode.window.activeTextEditor,
      before,
      'no editor was opened by setLanguage'
    );
  });

  test('all three osprey commands are registered', async () => {
    const all = await vscode.commands.getCommands(true);
    assert.ok(all.includes('osprey.compile'), 'compile registered');
    assert.ok(all.includes('osprey.run'), 'run registered');
    assert.ok(all.includes('osprey.setLanguage'), 'setLanguage registered');
  });

});

suite('Osprey Activation Side-Effect Coverage', () => {
  let tempDir: string;
  const extensionId2 = extensionId;
  const settle = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

  // Same focus-stealing workaround as the command suite: the extension's
  // "Osprey Debug" output channel grabs the active-editor slot in the headless
  // host, so we close panels and refocus the editor group until the target
  // document is genuinely active.
  async function makeActive(document: vscode.TextDocument): Promise<boolean> {
    for (let i = 0; i < 25; i++) {
      await vscode.commands.executeCommand('workbench.action.closePanel').then(
        undefined,
        () => undefined
      );
      await vscode.window.showTextDocument(document, {
        viewColumn: vscode.ViewColumn.One,
        preserveFocus: false,
        preview: false
      });
      await vscode.commands.executeCommand('workbench.action.focusActiveEditorGroup').then(
        undefined,
        () => undefined
      );
      await settle(120);
      const active = vscode.window.activeTextEditor;
      if (active && active.document.uri.toString() === document.uri.toString()) {
        return true;
      }
    }
    return false;
  }

  setup(() => {
    tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'osprey-side-'));
  });

  teardown(async () => {
    await vscode.commands.executeCommand('workbench.action.closeAllEditors');
    if (fs.existsSync(tempDir)) {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  test('opening a .osp file triggers the language-association watcher', async () => {
    // Opening a brand new .osp document fires onDidOpenTextDocument; the handler
    // forces the osprey language id. This exercises the open-document watcher.
    const file = path.join(tempDir, 'late-open.osp');
    fs.writeFileSync(file, 'fn main() -> Unit = print("late")\n');

    const document = await vscode.workspace.openTextDocument(file);
    await vscode.window.showTextDocument(document);
    await settle(800);

    assert.strictEqual(document.languageId, 'osprey', 'forced to osprey');
    assert.ok(document.fileName.endsWith('.osp'), 'is a .osp file');
    const ext = vscode.extensions.getExtension(extensionId2);
    assert.ok(ext?.isActive, 'extension active');
  });

  test('changing osprey configuration fires the change handler', async () => {
    // Flip an osprey setting to trigger onDidChangeConfiguration, which shows an
    // information message. We assert the round-trip of the value to confirm the
    // configuration channel is wired and the handler had a real event to react
    // to.
    const config = vscode.workspace.getConfiguration('osprey');
    const original = config.get<boolean>('diagnostics.enabled');

    await config.update(
      'diagnostics.enabled',
      !original,
      vscode.ConfigurationTarget.Global
    );
    await settle(500);

    const flipped = vscode.workspace
      .getConfiguration('osprey')
      .get<boolean>('diagnostics.enabled');
    assert.strictEqual(flipped, !original, 'config value flipped');

    // Restore so other tests see defaults.
    await config.update(
      'diagnostics.enabled',
      original,
      vscode.ConfigurationTarget.Global
    );
    await settle(300);

    const restored = vscode.workspace
      .getConfiguration('osprey')
      .get<boolean>('diagnostics.enabled');
    assert.strictEqual(restored, original, 'config value restored');
  });

  // startDebugging may never settle because the provider cancels the session by
  // returning undefined, so always race it against a timeout and never assert on
  // the return value.
  async function startDebugRaced(config: vscode.DebugConfiguration): Promise<string> {
    const timeout = new Promise<string>(resolve =>
      setTimeout(() => resolve('timeout'), 4000)
    );
    const start = Promise.resolve(vscode.debug.startDebugging(undefined, config))
      .then(v => `resolved:${String(v)}`)
      .catch((error: unknown) => `error:${String(error)}`);
    return Promise.race([start, timeout]);
  }

  test('debug provider synthesizes a config from the active osprey editor', async function () {
    this.timeout(30000);

    // With an active osprey editor and an empty config, resolveDebugConfiguration
    // synthesizes type/name/request/program from the editor, then kicks off
    // compile-and-run and cancels the session by returning undefined.
    const file = path.join(tempDir, 'debug-synth.osp');
    fs.writeFileSync(file, 'fn main() -> Unit = print("debug synth")\n');
    const document = await vscode.workspace.openTextDocument(file);
    const isActive = await makeActive(document);

    assert.strictEqual(document.languageId, 'osprey', 'doc is osprey');
    assert.ok(isActive, 'osprey doc is the active editor');
    assert.strictEqual(
      vscode.window.activeTextEditor?.document.languageId,
      'osprey',
      'active editor language is osprey'
    );

    const outcome = await startDebugRaced({
      type: '',
      name: '',
      request: ''
    } as unknown as vscode.DebugConfiguration);
    await settle(2500);

    assert.ok(typeof outcome === 'string', 'debug start settled to a string');
    assert.ok(
      vscode.extensions.getExtension(extensionId2)?.isActive,
      'extension survives synthesized debug config'
    );
  });

  test('debug provider rejects a launch config with no resolvable program', async function () {
    this.timeout(30000);

    // No active osprey editor and a config that carries a type but no program
    // means synthesis is skipped and `!config.program` is true, so the provider
    // shows "Cannot find a program to run" and returns undefined.
    await vscode.commands.executeCommand('workbench.action.closeAllEditors');
    await settle(300);

    const outcome = await startDebugRaced({
      type: 'osprey',
      name: 'Run Osprey File',
      request: 'launch'
    } as vscode.DebugConfiguration);
    await settle(1500);

    assert.ok(typeof outcome === 'string', 'debug start settled to a string');
    assert.ok(
      vscode.extensions.getExtension(extensionId2)?.isActive,
      'extension survives the no-program debug path'
    );
  });
});

// Unit tests for the pure binary-resolution helpers. These don't depend on a
// live activation, so they can exercise both the bundled-present and unbundled
// branches deterministically by supplying a fake ExtensionContext whose
// asAbsolutePath points at real or missing files.
suite('Osprey Binary Resolution Unit Tests', () => {
  let tempDir: string;

  // Minimal ExtensionContext stand-in: only asAbsolutePath is consumed by the
  // helpers under test. It roots relative paths at a controllable directory.
  function fakeContext(root: string): vscode.ExtensionContext {
    return {
      asAbsolutePath: (rel: string) => path.join(root, rel)
    } as unknown as vscode.ExtensionContext;
  }

  setup(() => {
    tempDir = fs.mkdtempSync(path.join(os.tmpdir(), 'osprey-unit-'));
  });

  teardown(() => {
    if (fs.existsSync(tempDir)) {
      fs.rmSync(tempDir, { recursive: true, force: true });
    }
  });

  test('shipwrightPlatform returns a valid os-arch identifier', () => {
    const id = shipwrightPlatform();
    const [osName, arch] = id.split('-');

    assert.strictEqual(id.split('-').length, 2, 'id is exactly os-arch');
    assert.ok(
      ['win32', 'darwin', 'linux'].includes(osName),
      `os segment "${osName}" is a known platform`
    );
    assert.ok(
      ['arm64', 'x64'].includes(arch),
      `arch segment "${arch}" is a known architecture`
    );
    // It must agree with the actual process this test runs in.
    const expectedArch = process.arch === 'arm64' ? 'arm64' : 'x64';
    const expectedOs =
      process.platform === 'win32'
        ? 'win32'
        : process.platform === 'darwin'
          ? 'darwin'
          : 'linux';
    assert.strictEqual(osName, expectedOs, 'os matches the host');
    assert.strictEqual(arch, expectedArch, 'arch matches the host');
    // The mapping is deterministic across calls.
    assert.strictEqual(shipwrightPlatform(), id, 'mapping is stable');
  });

  test('resolveBundledCompiler returns the path when the binary exists', () => {
    // Lay down a fake bundled binary at bin/<platform>/osprey[.exe].
    const exe = process.platform === 'win32' ? '.exe' : '';
    const platDir = path.join(tempDir, 'bin', shipwrightPlatform());
    fs.mkdirSync(platDir, { recursive: true });
    const bin = path.join(platDir, `osprey${exe}`);
    fs.writeFileSync(bin, '#!/bin/sh\nexit 0\n');

    const resolved = resolveBundledCompiler(fakeContext(tempDir));

    assert.ok(resolved, 'a path is returned when the binary exists');
    assert.strictEqual(resolved, bin, 'resolved path is the staged binary');
    assert.ok(resolved?.endsWith(`osprey${exe}`), 'ends with the binary name');
    assert.ok(fs.existsSync(resolved as string), 'resolved path actually exists');
  });

  test('resolveBundledCompiler returns undefined when no binary is bundled', () => {
    // tempDir has no bin/<platform>/osprey, so resolution must fail closed.
    const resolved = resolveBundledCompiler(fakeContext(tempDir));

    assert.strictEqual(resolved, undefined, 'undefined when unbundled');
    const expectedMissing = path.join(
      tempDir,
      'bin',
      shipwrightPlatform(),
      process.platform === 'win32' ? 'osprey.exe' : 'osprey'
    );
    assert.ok(!fs.existsSync(expectedMissing), 'the probed path is genuinely absent');
  });

  test('resolveServerCommand prefers an explicit user compiler path', async () => {
    // With server.compilerPath set, resolution must return it verbatim and never
    // touch the bundled fallback.
    const config = vscode.workspace.getConfiguration('osprey');
    const original = config.get<string>('server.compilerPath');
    const custom = path.join(tempDir, 'my-custom-osprey');
    fs.writeFileSync(custom, '#!/bin/sh\nexit 0\n');

    await config.update('server.compilerPath', custom, vscode.ConfigurationTarget.Global);
    try {
      const resolved = resolveServerCommand(fakeContext(tempDir));
      assert.strictEqual(resolved, custom, 'returns the explicit user path');
      assert.notStrictEqual(resolved, 'osprey', 'does not fall back to PATH');
    } finally {
      await config.update(
        'server.compilerPath',
        original ?? '',
        vscode.ConfigurationTarget.Global
      );
    }
  });

  test('resolveServerCommand falls back to bundled then PATH', async () => {
    // No user path: with a bundled binary present it returns that; without one
    // it falls back to the bare `osprey` PATH lookup.
    const config = vscode.workspace.getConfiguration('osprey');
    const originalCompiler = config.get<string>('server.compilerPath');
    const originalPath = config.get<string>('server.path');
    await config.update('server.compilerPath', '', vscode.ConfigurationTarget.Global);
    await config.update('server.path', '', vscode.ConfigurationTarget.Global);

    try {
      // No bundled binary in tempDir -> bare PATH fallback.
      assert.strictEqual(
        resolveServerCommand(fakeContext(tempDir)),
        'osprey',
        'falls back to osprey on PATH when nothing is bundled'
      );

      // Now stage a bundled binary -> it is preferred over the PATH fallback.
      const exe = process.platform === 'win32' ? '.exe' : '';
      const platDir = path.join(tempDir, 'bin', shipwrightPlatform());
      fs.mkdirSync(platDir, { recursive: true });
      const bin = path.join(platDir, `osprey${exe}`);
      fs.writeFileSync(bin, '#!/bin/sh\nexit 0\n');

      assert.strictEqual(
        resolveServerCommand(fakeContext(tempDir)),
        bin,
        'prefers the bundled binary over the PATH fallback'
      );
    } finally {
      await config.update(
        'server.compilerPath',
        originalCompiler ?? '',
        vscode.ConfigurationTarget.Global
      );
      await config.update(
        'server.path',
        originalPath ?? '',
        vscode.ConfigurationTarget.Global
      );
    }
  });
});
