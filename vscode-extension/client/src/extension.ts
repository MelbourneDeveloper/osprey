import * as path from 'path';
import * as vscode from 'vscode';
import { workspace, ExtensionContext, window, ConfigurationChangeEvent, commands, Uri, debug, languages } from 'vscode';
import { execFile } from 'child_process';
import * as fs from 'fs';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind
} from 'vscode-languageclient/node';
import { activateShipwright, type ActivationResult } from '@nimblesite/shipwright-vscode';

let client: LanguageClient;
// Resolved by Shipwright on activation. All subsequent execFile calls use this
// path; falling back to PATH is forbidden by SWR-IDE-* in the VSIX bundle path.
let ospreyBinaryPath = 'osprey';

export async function activate(context: ExtensionContext): Promise<void> {
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

  // Resolve the bundled `osprey` binary via Shipwright before starting the LSP.
  // SWR-IDE rules forbid PATH/global-install fallback for VS Code native-binary
  // extensions, so the per-platform bundle inside the VSIX is the only normal
  // startup source. User settings (osprey.binaries.*) override the bundled
  // binary; a mismatch fails activation instead of being silently masked.
  const shipwright = await resolveOspreyBinary(context, outputChannel);
  if (!shipwright.ok) {
    outputChannel.appendLine('Shipwright resolution failed; aborting activation.');
    return;
  }

  // Server options - use the TypeScript language server.
  // Try both layouts: top-level tsc -b (out/server/src/server.js) and the
  // server/tsconfig.json layout (server/out/src/server.js).
  const candidatePaths = [
    path.join('out', 'server', 'src', 'server.js'),
    path.join('server', 'out', 'src', 'server.js')
  ];
  let serverModule = '';
  for (const candidate of candidatePaths) {
    const abs = context.asAbsolutePath(candidate);
    if (fs.existsSync(abs)) {
      serverModule = abs;
      break;
    }
  }
  outputChannel.appendLine(`Server module path: ${serverModule || '(not found)'}`);

  if (!serverModule) {
    const errorMsg = `Server module not found in any of: ${candidatePaths.join(', ')}`;
    outputChannel.appendLine(`ERROR: ${errorMsg}`);
    window.showErrorMessage(errorMsg);
    return;
  }
  
  outputChannel.appendLine('Server module exists, proceeding with setup...');
  
  const debugOptions = { execArgv: ['--nolazy', '--inspect=6009'] };
  
  const serverOptions: ServerOptions = {
    run: { module: serverModule, transport: TransportKind.ipc },
    debug: {
      module: serverModule,
      transport: TransportKind.ipc,
      options: debugOptions
    }
  };

  // Client options
  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: 'file', language: 'osprey' }],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher('**/*.osp')
    },
    // The LSP server reads ospreyBinaryPath via params.initializationOptions
    // and uses it for every execFile call. Resolution happens once in the
    // extension host via Shipwright; the server never falls back to PATH.
    initializationOptions: {
      ospreyBinaryPath
    },
    outputChannelName: 'Osprey Language Server',
    revealOutputChannelOn: 4, // Error
    initializationFailedHandler: (error) => {
      outputChannel.appendLine(`Initialization failed: ${error}`);
      window.showErrorMessage(`Osprey language server initialization failed: ${error}`);
      return false;
    },
    errorHandler: {
      error: (error, message, count) => {
        outputChannel.appendLine(`Language server error: ${error}, message: ${message}, count: ${count}`);
        return { action: 1 }; // Continue
      },
      closed: () => {
        outputChannel.appendLine('Language server connection closed');
        return { action: 1 }; // Restart
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
    
    // Use the Shipwright-resolved osprey compiler (bundled-or-override).
    execFile(ospreyBinaryPath, [document.fileName],
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
    
    // Use the Shipwright-resolved osprey compiler with --run flag.
    execFile(ospreyBinaryPath, [document.fileName, '--run'],
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