import * as path from "path";
import {
  workspace,
  ExtensionContext,
  window,
  commands,
  debug,
  languages,
} from "vscode";
import { execFile } from "child_process";
import * as fs from "fs";
import {
  CloseAction,
  ErrorAction,
  Executable,
  LanguageClient,
  LanguageClientOptions,
  RevealOutputChannelOn,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

// @nimblesite/shipwright-vscode is ESM-only; this extension is CommonJS, so it
// is loaded via dynamic import() (never a static require) inside activate().

let client: LanguageClient;

// shipwrightPlatform maps the Node platform/arch to the Shipwright platform id
// (e.g. darwin-arm64, win32-x64) used in the bundled binary path. Exported for
// unit testing of the platform-string mapping.
export function shipwrightPlatform(): string {
  const arch = process.arch === "arm64" ? "arm64" : "x64";
  const os =
    process.platform === "win32"
      ? "win32"
      : process.platform === "darwin"
        ? "darwin"
        : "linux";
  return `${os}-${arch}`;
}

// resolveBundledCompiler returns the absolute path to the version-matched
// osprey binary bundled in this VSIX for the current platform, or undefined
// when running unbundled (e.g. a local dev install). The release pipeline
// stages it at bin/<platform>/osprey[.exe]. [SWR-VERSION-MANIFEST] Exported so
// both the bundled-present and unbundled branches can be unit tested.
export function resolveBundledCompiler(
  context: ExtensionContext,
): string | undefined {
  const exe = process.platform === "win32" ? ".exe" : "";
  const bundled = context.asAbsolutePath(
    path.join("bin", shipwrightPlatform(), `osprey${exe}`),
  );
  return fs.existsSync(bundled) ? bundled : undefined;
}

// looksLikePath reports whether a configured compiler value is a filesystem
// path (absolute or relative) rather than a bare command name resolved on PATH.
// Only path-like values are existence-checked; a bare `osprey` is left for the
// OS to resolve at spawn time.
export function looksLikePath(value: string): boolean {
  return value.includes("/") || value.includes("\\");
}

// resolveServerCommand picks the osprey binary that backs the language server:
// an explicit user setting, then the version-matched bundled compiler, then a
// plain `osprey` on PATH. The server is launched as `<command> lsp` over stdio.
// A configured path that points at a MISSING file would make the language
// client fail to spawn (ENOENT) and silently kill every feature — hover,
// diagnostics, go-to-definition. Rather than die, fall back to the bundled/PATH
// compiler and warn. `warn` is injectable so the fallback branch is unit
// testable; it defaults to a no-op. Exported so each branch is unit tested
// independently of a single live activation.
export function resolveServerCommand(
  context: ExtensionContext,
  warn: (message: string) => void = () => undefined,
): string {
  const config = workspace.getConfiguration("osprey");
  const userPath =
    config.get<string>("server.compilerPath") ||
    config.get<string>("server.path");
  if (userPath) {
    if (looksLikePath(userPath) && !fs.existsSync(userPath)) {
      const fallback = resolveBundledCompiler(context) ?? "osprey";
      warn(
        `osprey.server.compilerPath "${userPath}" does not exist; ` +
          `falling back to "${fallback}". Run \`make build\` to produce it.`,
      );
      return fallback;
    }
    return userPath;
  }
  return resolveBundledCompiler(context) ?? "osprey";
}

// makeClientFailureHandling builds the language client's failure callbacks: the
// one-shot initialization-failed handler and the runtime error/closed handlers
// that keep the server alive (Continue) or restart it (Restart). These fire only
// on real LSP transport failures, which an integration test cannot reliably
// induce — so they are extracted here and the side effects (`log`, `showError`)
// are injected, letting each callback be unit-tested directly. Behaviour is
// identical to the previous inline handlers.
export function makeClientFailureHandling(
  log: (message: string) => void,
  showError: (message: string) => void,
): Pick<LanguageClientOptions, "initializationFailedHandler" | "errorHandler"> {
  return {
    initializationFailedHandler: (error) => {
      log(`Initialization failed: ${error}`);
      showError(`Osprey language server initialization failed: ${error}`);
      return false;
    },
    errorHandler: {
      error: (error, message, count) => {
        log(
          `Language server error: ${error}, message: ${message}, count: ${count}`,
        );
        return { action: ErrorAction.Continue };
      },
      closed: () => {
        log("Language server connection closed; restarting");
        return { action: CloseAction.Restart };
      },
    },
  };
}

// A minimal stand-in for the active editor the debug provider reads — just the
// document fields the synthesis needs.
export interface ActiveEditorLike {
  document: { languageId: string; fileName: string };
}

// applyDefaultOspreyDebugConfig fills an otherwise-empty launch config from the
// active osprey editor, so pressing Run with no `.vscode/launch.json` still
// works ([EDITOR-VSCODE]). It mutates and returns `config`: synthesis happens
// only when type/request/name are all absent AND an osprey document is focused;
// any already-populated config is returned untouched. Pure (no VS Code globals)
// so the debug provider's branches are unit-testable without a debug session.
export function applyDefaultOspreyDebugConfig(
  config: any,
  activeEditor: ActiveEditorLike | undefined,
): any {
  if (!config.type && !config.request && !config.name) {
    if (activeEditor && activeEditor.document.languageId === "osprey") {
      config.type = "osprey";
      config.name = "Run Osprey File";
      config.request = "launch";
      config.program = activeEditor.document.fileName;
    }
  }
  return config;
}

export function activate(context: ExtensionContext) {
  console.log("Osprey extension is now active!");

  // Create output channel for diagnostics
  const outputChannel = window.createOutputChannel("Osprey Debug");
  outputChannel.appendLine("=== Osprey Extension Activation ===");
  outputChannel.show();

  // Check if Osprey server is enabled
  const config = workspace.getConfiguration("osprey");
  if (!config.get("server.enabled", true)) {
    outputChannel.appendLine("Language server is disabled in configuration");
    return;
  }

  // Shipwright: verify the bundled osprey compiler matches the version this
  // extension expects before we launch it for diagnostics. On mismatch the
  // host surfaces a prompt-reinstall message (hosts.vscode.onMismatch).
  // [SWR-VERSION-HANDSHAKE] Best-effort: never block activation on it.
  const manifestPath = context.asAbsolutePath("shipwright.json");
  if (fs.existsSync(manifestPath)) {
    // Adapter normalizing VS Code's Thenable-returning API to the Promise-typed
    // shape the library expects (VscodeApiLike).
    const vscodeApi = {
      workspace: {
        getConfiguration: (s?: string) => workspace.getConfiguration(s),
      },
      window: {
        showErrorMessage: (
          m: string,
          o: { modal: boolean },
          ...items: string[]
        ) => Promise.resolve(window.showErrorMessage(m, o, ...items)),
        showWarningMessage: (
          m: string,
          o: { modal: boolean },
          ...items: string[]
        ) => Promise.resolve(window.showWarningMessage(m, o, ...items)),
      },
    };
    void (async () => {
      try {
        const sw = await import("@nimblesite/shipwright-vscode");
        const r = await sw.activateShipwright(context, {
          vscode: vscodeApi,
          manifestPath,
          showMessages: true,
        });
        outputChannel.appendLine(
          `Shipwright activation: ok=${r.ok}, diagnostics=${r.diagnostics.length}`,
        );
      } catch (e) {
        outputChannel.appendLine(`Shipwright activation error: ${e}`);
      }
    })();
  }

  // The language server is the Rust `osprey lsp` subcommand (the osprey-lsp
  // crate, built on the published lspkit crates), spoken over stdio. Resolve
  // the binary: explicit user setting first, then the version-matched bundled
  // compiler, then `osprey` on PATH.
  const ospreyCommand = resolveServerCommand(context, (m) => {
    outputChannel.appendLine(m);
    window.showWarningMessage(m);
  });
  outputChannel.appendLine(`Language server command: ${ospreyCommand} lsp`);

  const serverExecutable: Executable = {
    command: ospreyCommand,
    args: ["lsp"],
    transport: TransportKind.stdio,
  };
  const serverOptions: ServerOptions = {
    run: serverExecutable,
    debug: serverExecutable,
  };

  // Client options. The server analyzes document text (not the filesystem), so
  // unsaved `untitled:` buffers are supported alongside on-disk files.
  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: "file", language: "osprey" },
      { scheme: "untitled", language: "osprey" },
    ],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("**/*.osp"),
    },
    outputChannelName: "Osprey Language Server",
    revealOutputChannelOn: RevealOutputChannelOn.Error,
    ...makeClientFailureHandling(
      (message) => outputChannel.appendLine(message),
      (message) => {
        window.showErrorMessage(message);
      },
    ),
  };

  // Create and start the language client
  client = new LanguageClient(
    "ospreyLanguageServer",
    "Osprey Language Server",
    serverOptions,
    clientOptions,
  );

  outputChannel.appendLine("Starting language client...");

  // Start the client and server
  client
    .start()
    .then(() => {
      outputChannel.appendLine(
        "SUCCESS: Osprey language server started successfully",
      );
      console.log("Osprey language server started successfully");
    })
    .catch((error: any) => {
      const errorMsg = `Failed to start Osprey language server: ${error.message || error}`;
      outputChannel.appendLine(`ERROR: ${errorMsg}`);
      outputChannel.appendLine(
        `Error stack: ${error.stack || "No stack trace"}`,
      );
      console.error("Failed to start Osprey language server:", error);
      window.showErrorMessage(errorMsg);
    });

  // Add status bar item
  const statusBar = window.createStatusBarItem();
  statusBar.text = "$(check) Osprey";
  statusBar.tooltip = "Osprey Language Server is running";
  statusBar.show();
  context.subscriptions.push(statusBar);

  // Register debug adapter
  const provider = debug.registerDebugAdapterDescriptorFactory("osprey", {
    createDebugAdapterDescriptor(_session: any) {
      // Return null to use inline debug adapter
      return null;
    },
  });

  context.subscriptions.push(provider);

  // Register debug configuration provider
  context.subscriptions.push(
    debug.registerDebugConfigurationProvider("osprey", {
      resolveDebugConfiguration(folder: any, config: any, token: any) {
        // If no config is provided, synthesize one from the active osprey editor.
        config = applyDefaultOspreyDebugConfig(config, window.activeTextEditor);

        if (!config.program) {
          return window
            .showInformationMessage("Cannot find a program to run")
            .then((_) => {
              return undefined;
            });
        }

        // Actually run the Osprey program instead of debugging
        compileAndRunCurrentFile(resolveServerCommand(context));
        return undefined; // Cancel the debug session
      },
    }),
  );

  // Auto-detect and force language association for .osp files
  workspace.onDidOpenTextDocument((document) => {
    outputChannel.appendLine(`📁 Document opened: ${document.fileName}`);
    if (
      document.fileName.endsWith(".osp") &&
      document.languageId !== "osprey"
    ) {
      outputChannel.appendLine(
        `🔧 Forcing language association for ${document.fileName} (was: ${document.languageId})`,
      );
      // Use the proper API to set language
      languages.setTextDocumentLanguage(document, "osprey").then(
        () => {
          outputChannel.appendLine(
            `✅ Successfully set language to osprey for ${document.fileName}`,
          );
        },
        (error: any) => {
          outputChannel.appendLine(`❌ Failed to set language: ${error}`);
        },
      );
    }
  });

  // Check already open documents
  workspace.textDocuments.forEach((document) => {
    if (
      document.fileName.endsWith(".osp") &&
      document.languageId !== "osprey"
    ) {
      outputChannel.appendLine(
        `🔧 Forcing language association for already open file: ${document.fileName}`,
      );
      languages.setTextDocumentLanguage(document, "osprey");
    }
  });

  // Register commands
  context.subscriptions.push(
    commands.registerCommand("osprey.compile", () => {
      compileCurrentFile(resolveServerCommand(context));
    }),
    commands.registerCommand("osprey.run", () => {
      compileAndRunCurrentFile(resolveServerCommand(context));
    }),
    commands.registerCommand("osprey.setLanguage", () => {
      const activeEditor = window.activeTextEditor;
      if (activeEditor) {
        languages.setTextDocumentLanguage(activeEditor.document, "osprey");
        window.showInformationMessage("Set language to Osprey");
      }
    }),
    workspace.onDidChangeConfiguration((event: any) => {
      if (event.affectsConfiguration("osprey")) {
        window.showInformationMessage(
          "Osprey configuration changed. Restart required.",
        );
      }
    }),
  );
}

function compileCurrentFile(compilerCommand: string) {
  const activeEditor = window.activeTextEditor;
  if (!activeEditor) {
    window.showErrorMessage("No active Osprey file found");
    return;
  }

  const document = activeEditor.document;
  if (!document.fileName.endsWith(".osp")) {
    window.showErrorMessage("Please open a .osp file to compile");
    return;
  }

  // Save the file first
  document.save().then(() => {
    const outputChannel = window.createOutputChannel("Osprey Compiler");
    outputChannel.show();
    outputChannel.appendLine(`Compiling ${document.fileName}...`);

    // Get the directory containing the file (no workspace required)
    const fileDir = path.dirname(document.fileName);

    // Use the resolved osprey compiler (user setting → version-matched bundled
    // binary → `osprey` on PATH) — same resolution the language server uses.
    execFile(
      compilerCommand,
      [document.fileName],
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
          outputChannel.appendLine(`Exit code: ${error.code || "unknown"}`);
          outputChannel.appendLine(`Signal: ${error.signal || "none"}`);
          outputChannel.appendLine(`Error message: ${error.message}`);
          window.showErrorMessage(
            "Compilation failed. Check output for details.",
          );
        } else {
          outputChannel.appendLine("=== COMPILATION SUCCESS ===");
          window.showInformationMessage("Osprey file compiled successfully!");
        }

        outputChannel.appendLine(`=== END OUTPUT ===`);
      },
    );
  });
}

function compileAndRunCurrentFile(compilerCommand: string) {
  const activeEditor = window.activeTextEditor;
  if (!activeEditor) {
    window.showErrorMessage("No active Osprey file found");
    return;
  }

  const document = activeEditor.document;
  if (!document.fileName.endsWith(".osp")) {
    window.showErrorMessage("Please open a .osp file to run");
    return;
  }

  // Save the file first
  document.save().then(() => {
    const outputChannel = window.createOutputChannel("Osprey Runner");
    outputChannel.show();
    outputChannel.appendLine(`Compiling and running ${document.fileName}...`);

    // Get the directory containing the file (no workspace required)
    const fileDir = path.dirname(document.fileName);

    // Use the resolved osprey compiler with --run (user setting → version-matched
    // bundled binary → `osprey` on PATH) — same resolution the language server uses.
    execFile(
      compilerCommand,
      [document.fileName, "--run"],
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
          outputChannel.appendLine(`Exit code: ${error.code || "unknown"}`);
          outputChannel.appendLine(`Signal: ${error.signal || "none"}`);
          outputChannel.appendLine(`Error message: ${error.message}`);
          window.showErrorMessage(
            "Compilation or execution failed. Check output for details.",
          );
        } else {
          outputChannel.appendLine("=== SUCCESS ===");
          window.showInformationMessage(
            "Osprey program executed successfully!",
          );
        }

        outputChannel.appendLine(`=== END OUTPUT ===`);
      },
    );
  });
}

export function deactivate(): Promise<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
