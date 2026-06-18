import * as path from 'path';
import { workspace, ExtensionContext, window, ConfigurationChangeEvent, commands, Uri, debug, languages } from 'vscode';
import { execFile } from 'child_process';
import * as fs from 'fs';
import {
  CloseAction,
  ErrorAction,
  Executable,
  LanguageClient,
  LanguageClientOptions,
  RevealOutputChannelOn,
  ServerOptions,
  TransportKind
} from 'vscode-languageclient/node';

// @nimblesite/shipwright-vscode is ESM-only; this extension is CommonJS, so it
// is loaded via dynamic import() (never a static require) inside activate().

let client: LanguageClient;

// shipwrightPlatform maps the Node platform/arch to the Shipwright platform id
// (e.g. darwin-arm64, win32-x64) used in the bundled binary path.
function shipwrightPlatform(): string {
  const arch = process.arch === 'arm64' ? 'arm64' : 'x64';
  const os = process.platform === 'win32' ? 'win32'
    : process.platform === 'darwin' ? 'darwin' : 'linux';
  return `${os}-${arch}`;
}

// resolveBundledCompiler returns the absolute path to the version-matched
// osprey binary bundled in this VSIX for the current platform, or undefined
// when running unbundled (e.g. a local dev install). The release pipeline
// stages it at bin/<platform>/osprey[.exe]. [SWR-VERSION-MANIFEST]
function resolveBundledCompiler(context: ExtensionContext): string | undefined {
  const exe = process.platform === 'win32' ? '.exe' : '';
  const bundled = context.asAbsolutePath(path.join('bin', shipwrightPlatform(), `osprey${exe}`));
  return fs.existsSync(bundled) ? bundled : undefined;
}

// resolveServerCommand picks the osprey binary that backs the language server:
// an explicit user setting, then the version-matched bundled compiler, then a
// plain `osprey` on PATH. The server is launched as `<command> lsp` over stdio.
function resolveServerCommand(context: ExtensionContext): string {
  const config = workspace.getConfiguration('osprey');
  const userPath = config.get<string>('server.compilerPath') || config.get<string>('server.path');
  if (userPath) {
    return userPath;
  }
  return resolveBundledCompiler(context) ?? 'osprey';
}

export function activate(context: ExtensionContext) {
  console.log('Osprey extension is now active!');
  
  // Create output channel for diagnostics
  const outputChannel = window.createOutputChannel('Osprey Debug');
  outputChannel.appendLine('=== Osprey Extension Activation ===');
  outputChannel.show();

  // Check if Osprey server is enabled
  const config = workspace.getConfiguration('osprey');
  if (!config.get('server.enabled', true)) {
    outputChannel.appendLine('Language server is disabled in configuration');
    return;
  }

  // Shipwright: verify the bundled osprey compiler matches the version this
  // extension expects before we launch it for diagnostics. On mismatch the
  // host surfaces a prompt-reinstall message (hosts.vscode.onMismatch).
  // [SWR-VERSION-HANDSHAKE] Best-effort: never block activation on it.
  const manifestPath = context.asAbsolutePath('shipwright.json');
  if (fs.existsSync(manifestPath)) {
    // Adapter normalizing VS Code's Thenable-returning API to the Promise-typed
    // shape the library expects (VscodeApiLike).
    const vscodeApi = {
      workspace: { getConfiguration: (s?: string) => workspace.getConfiguration(s) },
      window: {
        showErrorMessage: (m: string, o: { modal: boolean }, ...items: string[]) =>
          Promise.resolve(window.showErrorMessage(m, o, ...items)),
        showWarningMessage: (m: string, o: { modal: boolean }, ...items: string[]) =>
          Promise.resolve(window.showWarningMessage(m, o, ...items))
      }
    };
    void (async () => {
      try {
        const sw = await import('@nimblesite/shipwright-vscode');
        const r = await sw.activateShipwright(context, { vscode: vscodeApi, manifestPath, showMessages: true });
        outputChannel.appendLine(`Shipwright activation: ok=${r.ok}, diagnostics=${r.diagnostics.length}`);
      } catch (e) {
        outputChannel.appendLine(`Shipwright activation error: ${e}`);
      }
    })();
  }

  // The language server is the Rust `osprey lsp` subcommand (the osprey-lsp
  // crate, built on the published lspkit crates), spoken over stdio. Resolve
  // the binary: explicit user setting first, then the version-matched bundled
  // compiler, then `osprey` on PATH.
  const ospreyCommand = resolveServerCommand(context);
  outputChannel.appendLine(`Language server command: ${ospreyCommand} lsp`);

  const serverExecutable: Executable = {
    command: ospreyCommand,
    args: ['lsp'],
    transport: TransportKind.stdio
  };
  const serverOptions: ServerOptions = {
    run: serverExecutable,
    debug: serverExecutable
  };

  // Client options. The server analyzes document text (not the filesystem), so
  // unsaved `untitled:` buffers are supported alongside on-disk files.
  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: 'file', language: 'osprey' },
      { scheme: 'untitled', language: 'osprey' }
    ],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher('**/*.osp')
    },
    outputChannelName: 'Osprey Language Server',
    revealOutputChannelOn: RevealOutputChannelOn.Error,
    initializationFailedHandler: (error) => {
      outputChannel.appendLine(`Initialization failed: ${error}`);
      window.showErrorMessage(`Osprey language server initialization failed: ${error}`);
      return false;
    },
    errorHandler: {
      error: (error, message, count) => {
        outputChannel.appendLine(`Language server error: ${error}, message: ${message}, count: ${count}`);
        return { action: ErrorAction.Continue };
      },
      closed: () => {
        outputChannel.appendLine('Language server connection closed; restarting');
        return { action: CloseAction.Restart };
      }
    }
  };

  // Create and start the language client
  client = new LanguageClient(
    'ospreyLanguageServer',
    'Osprey Language Server',
    serverOptions,
    clientOptions
  );

  outputChannel.appendLine('Starting language client...');

  // Start the client and server
  client.start().then(() => {
    outputChannel.appendLine('SUCCESS: Osprey language server started successfully');
    console.log('Osprey language server started successfully');
  }).catch((error: any) => {
    const errorMsg = `Failed to start Osprey language server: ${error.message || error}`;
    outputChannel.appendLine(`ERROR: ${errorMsg}`);
    outputChannel.appendLine(`Error stack: ${error.stack || 'No stack trace'}`);
    console.error('Failed to start Osprey language server:', error);
    window.showErrorMessage(errorMsg);
  });

  // Add status bar item
  const statusBar = window.createStatusBarItem();
  statusBar.text = "$(check) Osprey";
  statusBar.tooltip = "Osprey Language Server is running";
  statusBar.show();
  context.subscriptions.push(statusBar);

  // Register debug adapter
  const provider = debug.registerDebugAdapterDescriptorFactory('osprey', {
    createDebugAdapterDescriptor(_session: any) {
      // Return null to use inline debug adapter
      return null;
    }
  });

  context.subscriptions.push(provider);

  // Register debug configuration provider
  context.subscriptions.push(debug.registerDebugConfigurationProvider('osprey', {
    resolveDebugConfiguration(folder: any, config: any, token: any) {
      // If no config is provided, create a default one
      if (!config.type && !config.request && !config.name) {
        const editor = window.activeTextEditor;
        if (editor && editor.document.languageId === 'osprey') {
          config.type = 'osprey';
          config.name = 'Run Osprey File';
          config.request = 'launch';
          config.program = editor.document.fileName;
        }
      }

      if (!config.program) {
        return window.showInformationMessage("Cannot find a program to run").then(_ => {
          return undefined;
        });
      }

      // Actually run the Osprey program instead of debugging
      compileAndRunCurrentFile();
      return undefined; // Cancel the debug session
    }
  }));

  // Auto-detect and force language association for .osp files
  workspace.onDidOpenTextDocument((document) => {
    outputChannel.appendLine(`📁 Document opened: ${document.fileName}`);
    if (document.fileName.endsWith('.osp') && document.languageId !== 'osprey') {
      outputChannel.appendLine(`🔧 Forcing language association for ${document.fileName} (was: ${document.languageId})`);
      // Use the proper API to set language
      languages.setTextDocumentLanguage(document, 'osprey').then(() => {
        outputChannel.appendLine(`✅ Successfully set language to osprey for ${document.fileName}`);
      }, (error: any) => {
        outputChannel.appendLine(`❌ Failed to set language: ${error}`);
      });
    }
  });

  // Check already open documents
  workspace.textDocuments.forEach((document) => {
    if (document.fileName.endsWith('.osp') && document.languageId !== 'osprey') {
      outputChannel.appendLine(`🔧 Forcing language association for already open file: ${document.fileName}`);
      languages.setTextDocumentLanguage(document, 'osprey');
    }
  });

  // Register commands
  context.subscriptions.push(
    commands.registerCommand('osprey.compile', () => {
      compileCurrentFile();
    }),
    commands.registerCommand('osprey.run', () => {
      compileAndRunCurrentFile();
    }),
    commands.registerCommand('osprey.setLanguage', () => {
      const activeEditor = window.activeTextEditor;
      if (activeEditor) {
        languages.setTextDocumentLanguage(activeEditor.document, 'osprey');
        window.showInformationMessage('Set language to Osprey');
      }
    }),
    workspace.onDidChangeConfiguration((event: any) => {
      if (event.affectsConfiguration('osprey')) {
        window.showInformationMessage('Osprey configuration changed. Restart required.');
      }
    })
  );
}

function compileCurrentFile() {
  const activeEditor = window.activeTextEditor;
  if (!activeEditor) {
    window.showErrorMessage('No active Osprey file found');
    return;
  }

  const document = activeEditor.document;
  if (!document.fileName.endsWith('.osp')) {
    window.showErrorMessage('Please open a .osp file to compile');
    return;
  }

  // Save the file first
  document.save().then(() => {
    const outputChannel = window.createOutputChannel('Osprey Compiler');
    outputChannel.show();
    outputChannel.appendLine(`Compiling ${document.fileName}...`);

    // Get the directory containing the file (no workspace required)
    const fileDir = path.dirname(document.fileName);
    
    // Use the installed osprey compiler
    execFile('osprey', [document.fileName], 
      { cwd: fileDir }, 
      (error: any, stdout: any, stderr: any) => {
        outputChannel.appendLine(`=== COMPILATION OUTPUT ===`);
        
        if (stdout) {
          outputChannel.appendLine(`STDOUT:`);
          outputChannel.appendLine(stdout);
        }
        
        if (stderr) {
          outputChannel.appendLine(`STDERR:`);
          outputChannel.appendLine(stderr);
        }
        
        if (error) {
          outputChannel.appendLine(`ERROR:`);
          outputChannel.appendLine(`Exit code: ${error.code || 'unknown'}`);
          outputChannel.appendLine(`Signal: ${error.signal || 'none'}`);
          outputChannel.appendLine(`Error message: ${error.message}`);
          window.showErrorMessage('Compilation failed. Check output for details.');
        } else {
          outputChannel.appendLine('=== COMPILATION SUCCESS ===');
          window.showInformationMessage('Osprey file compiled successfully!');
        }
        
        outputChannel.appendLine(`=== END OUTPUT ===`);
      }
    );
  });
}

function compileAndRunCurrentFile() {
  const activeEditor = window.activeTextEditor;
  if (!activeEditor) {
    window.showErrorMessage('No active Osprey file found');
    return;
  }

  const document = activeEditor.document;
  if (!document.fileName.endsWith('.osp')) {
    window.showErrorMessage('Please open a .osp file to run');
    return;
  }

  // Save the file first
  document.save().then(() => {
    const outputChannel = window.createOutputChannel('Osprey Runner');
    outputChannel.show();
    outputChannel.appendLine(`Compiling and running ${document.fileName}...`);

    // Get the directory containing the file (no workspace required)
    const fileDir = path.dirname(document.fileName);
    
    // Use the installed osprey compiler with --run flag
    execFile('osprey', [document.fileName, '--run'], 
      { cwd: fileDir }, 
      (error: any, stdout: any, stderr: any) => {
        outputChannel.appendLine(`=== COMPILE AND RUN OUTPUT ===`);
        
        if (stdout) {
          outputChannel.appendLine(`STDOUT:`);
          outputChannel.appendLine(stdout);
        }
        
        if (stderr) {
          outputChannel.appendLine(`STDERR:`);
          outputChannel.appendLine(stderr);
        }
        
        if (error) {
          outputChannel.appendLine(`ERROR:`);
          outputChannel.appendLine(`Exit code: ${error.code || 'unknown'}`);
          outputChannel.appendLine(`Signal: ${error.signal || 'none'}`);
          outputChannel.appendLine(`Error message: ${error.message}`);
          window.showErrorMessage('Compilation or execution failed. Check output for details.');
        } else {
          outputChannel.appendLine('=== SUCCESS ===');
          window.showInformationMessage('Osprey program executed successfully!');
        }
        
        outputChannel.appendLine(`=== END OUTPUT ===`);
      }
    );
  });
}

export function deactivate(): Promise<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
} 