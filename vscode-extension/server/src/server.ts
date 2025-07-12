import {
    CodeActionKind,
    CompletionItem,
    CompletionItemKind,
    createConnection,
    Definition,
    Diagnostic,
    DiagnosticSeverity,
    DidChangeConfigurationNotification,
    DocumentSymbol,
    Hover,
    InitializeParams,
    InitializeResult,
    Location,
    MarkupKind,
    ParameterInformation,
    Position,
    ProposedFeatures,
    SignatureHelp,
    SignatureInformation,
    SymbolKind,
    TextDocumentPositionParams,
    TextDocuments,
    TextDocumentSyncKind,
    CodeAction,
    TextEdit,
    Range,
    DocumentHighlight,
    FoldingRange,
    WorkspaceSymbol,
    RenameParams,
    PrepareRenameParams,
    ImplementationParams,
    DocumentFormattingParams,
    DocumentRangeFormattingParams
} from 'vscode-languageserver/node';

import { execFile } from 'child_process';
import * as fs from 'fs';
import * as path from 'path';
import { TextDocument } from 'vscode-languageserver-textdocument';

// Enhanced symbol and type information tracking
interface OspreySymbol {
  name: string;
  type: string;
  kind: 'function' | 'variable' | 'type' | 'parameter' | 'effect' | 'operation' | 'module' | 'import';
  location: Location;
  documentation?: string;
  signature?: string;
  parameters?: OspreyParameter[];
  returnType?: string;
  effectSet?: string; // For functions with effects
  isEffectOperation?: boolean; // For effect operations
  parentEffect?: string; // For effect operations
  isExported?: boolean; // For module exports
  isImported?: boolean; // For imports
}

interface OspreyParameter {
  name: string;
  type: string;
  documentation?: string;
  isOptional?: boolean;
}

interface DocumentAnalysis {
  symbols: OspreySymbol[];
  errors: Diagnostic[];
  uri: string;
  imports: string[]; // Track imported modules
  effects: string[]; // Track effect declarations
  handlers: string[]; // Track handler expressions
}

// Symbol reference tracking
interface SymbolReference {
  symbol: string;
  location: Location;
  kind: 'definition' | 'usage' | 'implementation' | 'effect_operation' | 'handler';
  context?: string; // Additional context like effect name for operations
}

// Global symbol table (workspace-wide)
const workspaceSymbols: Map<string, OspreySymbol[]> = new Map();
const documentAnalyses: Map<string, DocumentAnalysis> = new Map();
const symbolReferences: Map<string, SymbolReference[]> = new Map(); // Track all references
const effectOperations: Map<string, OspreySymbol[]> = new Map(); // Track effect operations
const implementations: Map<string, Location[]> = new Map(); // Track implementations

// Create a connection for the server
const connection = createConnection(ProposedFeatures.all);

// Log server startup
connection.console.log('🚀 Osprey Language Server starting up...');

// Create a text document manager
const documents: TextDocuments<TextDocument> = new TextDocuments(TextDocument);

// Add error handlers
process.on('uncaughtException', (error) => {
  connection.console.error(`💥 Uncaught exception: ${error.message}`);
  connection.console.error(`Stack: ${error.stack}`);
});

process.on('unhandledRejection', (reason, promise) => {
  connection.console.error(`💥 Unhandled promise rejection: ${reason}`);
  connection.console.error(`Promise: ${promise}`);
});

let hasConfigurationCapability = false;
let hasWorkspaceFolderCapability = false;
let hasDiagnosticRelatedInformationCapability = false;

connection.onInitialize((params: InitializeParams) => {
  connection.console.log('🎯 Initialize request received');
  connection.console.log(`Client capabilities: ${JSON.stringify(params.capabilities, null, 2)}`);
  
  try {
    const capabilities = params.capabilities;

    // Check client capabilities
    hasConfigurationCapability = !!(
      capabilities.workspace && !!capabilities.workspace.configuration
    );
    hasWorkspaceFolderCapability = !!(
      capabilities.workspace && !!capabilities.workspace.workspaceFolders
    );
    hasDiagnosticRelatedInformationCapability = !!(
      capabilities.textDocument &&
      capabilities.textDocument.publishDiagnostics &&
      capabilities.textDocument.publishDiagnostics.relatedInformation
    );

    connection.console.log(`✅ Configuration capability: ${hasConfigurationCapability}`);
    connection.console.log(`✅ Workspace folder capability: ${hasWorkspaceFolderCapability}`);
    connection.console.log(`✅ Diagnostic related info capability: ${hasDiagnosticRelatedInformationCapability}`);

    const result: InitializeResult = {
      capabilities: {
        textDocumentSync: TextDocumentSyncKind.Incremental,
        // Auto-completion
        completionProvider: {
          resolveProvider: true,
          triggerCharacters: ['.', ':', '$', '(', '|', '!', '{']
        },
        // Hover information (TYPE INFO!)
        hoverProvider: true,
        // Go to definition
        definitionProvider: true,
        // Go to implementation
        implementationProvider: true,
        // Find references
        referencesProvider: true,
        // Document symbols (outline)
        documentSymbolProvider: true,
        // Workspace symbols (global search)
        workspaceSymbolProvider: true,
        // Signature help (function parameters)
        signatureHelpProvider: {
          triggerCharacters: ['(', ',']
        },
        // Rename symbols
        renameProvider: {
          prepareProvider: true
        },
        // Code actions (quick fixes)
        codeActionProvider: {
          codeActionKinds: [
            CodeActionKind.QuickFix,
            CodeActionKind.Refactor,
            CodeActionKind.Source,
            CodeActionKind.SourceOrganizeImports
          ]
        },
        // Document formatting
        documentFormattingProvider: true,
        // Range formatting
        documentRangeFormattingProvider: true,
        // Document highlights (same symbol highlighting)
        documentHighlightProvider: true,
        // Folding ranges
        foldingRangeProvider: true
      }
    };

    if (hasWorkspaceFolderCapability) {
      result.capabilities.workspace = {
        workspaceFolders: {
          supported: true
        }
      };
    }

    connection.console.log(`🎉 Server capabilities configured: ${JSON.stringify(result.capabilities, null, 2)}`);
    return result;
  } catch (error) {
    connection.console.error(`💥 Error during initialization: ${error}`);
    throw error;
  }
});

connection.onInitialized(() => {
  connection.console.log('🎊 Server initialization completed');
  
  try {
    if (hasConfigurationCapability) {
      connection.console.log('📋 Registering configuration change notifications');
      connection.client.register(DidChangeConfigurationNotification.type, undefined);
    }
    
    if (hasWorkspaceFolderCapability) {
      connection.console.log('📁 Setting up workspace folder change handlers');
      connection.workspace.onDidChangeWorkspaceFolders((_event: any) => {
        connection.console.log('📂 Workspace folder change event received.');
      });
    }
    
    connection.console.log('🚀 Osprey Language Server is ready!');
    connection.console.log('🎯 ALL LSP FEATURES REGISTERED AND READY!');
  } catch (error) {
    connection.console.error(`💥 Error during post-initialization: ${error}`);
  }
});

// Global settings
interface OspreySettings {
  compilerPath: string;
  enableDiagnostics: boolean;
}

const defaultSettings: OspreySettings = {
  compilerPath: 'osprey',
  enableDiagnostics: true
};

let globalSettings: OspreySettings = defaultSettings;

// Cache the settings of all open documents
const documentSettings: Map<string, Promise<OspreySettings>> = new Map();

connection.onDidChangeConfiguration((change: any) => {
  if (hasConfigurationCapability) {
    // Reset all cached document settings
    documentSettings.clear();
  } else {
    globalSettings = <OspreySettings>(
      (change.settings.osprey || defaultSettings)
    );
  }

  // Revalidate all open text documents
  documents.all().forEach(validateTextDocument);
});

function getDocumentSettings(resource: string): Promise<OspreySettings> {
  if (!hasConfigurationCapability) {
    return Promise.resolve(globalSettings);
  }
  let result = documentSettings.get(resource);
  if (!result) {
    result = connection.workspace.getConfiguration({
      scopeUri: resource,
      section: 'osprey'
    }).then((config: any) => {
      return {
        compilerPath: config.server?.compilerPath || config.server?.path || 'osprey',
        enableDiagnostics: config.diagnostics?.enabled !== false
      };
    });
    documentSettings.set(resource, result);
  }
  return result;
}

// Only keep settings for open documents
documents.onDidClose((e: any) => {
  documentSettings.delete(e.document.uri);
  // Also clean up references
  symbolReferences.delete(e.document.uri);
});

// The content of a text document has changed
documents.onDidChangeContent((change: any) => {
  connection.console.log(`📄 Document content changed: ${change.document.uri}`);
  validateTextDocument(change.document);
});

// Document opened
documents.onDidOpen((event: any) => {
  connection.console.log(`📂 Document opened: ${event.document.uri}`);
  validateTextDocument(event.document);
});

async function validateTextDocument(textDocument: TextDocument): Promise<void> {
  connection.console.log(`🔍 VALIDATING DOCUMENT: ${textDocument.uri}`);
  
  const settings = await getDocumentSettings(textDocument.uri);
  connection.console.log(`⚙️ Settings: enableDiagnostics=${settings.enableDiagnostics}, compilerPath=${settings.compilerPath}`);
  
  // Force diagnostics to be enabled for comprehensive analysis
  settings.enableDiagnostics = true;
  
  if (!settings.enableDiagnostics) {
    connection.console.log(`⚠️ Diagnostics disabled for ${textDocument.uri}`);
    return;
  }

  connection.console.log(`🏃 Running analysis for ${textDocument.uri}`);
  
  // Parse symbols for IDE features using the compiler
  connection.console.log(`🔍 PARSING SYMBOLS for ${textDocument.uri}`);
  const symbols = await getSymbolsFromCompiler(textDocument.getText(), textDocument.uri);
  connection.console.log(`📝 FOUND ${symbols.length} SYMBOLS:`);
  if (symbols.length === 0) {
    connection.console.log(`⚠️ NO SYMBOLS FOUND - checking document content...`);
    const text = textDocument.getText();
    connection.console.log(`📄 Document has ${text.length} characters`);
    connection.console.log(`📄 First 200 chars: ${text.substring(0, 200)}`);
  } else {
    symbols.forEach((sym, i) => {
      connection.console.log(`  ${i+1}. ${sym.kind} ${sym.name}: ${sym.type} at line ${sym.location.range.start.line + 1}`);
    });
  }
  
  // Find all symbol references in the document
  const references = findAllSymbolReferences(textDocument);
  symbolReferences.set(textDocument.uri, references);
  
  // Run diagnostics
  const diagnostics: Diagnostic[] = await analyzeDocument(textDocument, settings);
  
  // Store analysis results
  const analysis: DocumentAnalysis = {
    symbols,
    errors: diagnostics,
    uri: textDocument.uri,
    imports: [], // Initialize
    effects: [], // Initialize
    handlers: [] // Initialize
  };
  documentAnalyses.set(textDocument.uri, analysis);
  workspaceSymbols.set(textDocument.uri, symbols);
  
  connection.console.log(`📊 Found ${diagnostics.length} diagnostics for ${textDocument.uri}`);
  diagnostics.forEach((diag, i) => {
    connection.console.log(`  ${i+1}. [${diag.severity}] Line ${diag.range.start.line + 1}: ${diag.message}`);
  });
  
  connection.sendDiagnostics({ uri: textDocument.uri, diagnostics });
  connection.console.log(`✅ Sent ${diagnostics.length} diagnostics to client`);
}

async function analyzeDocument(textDocument: TextDocument, settings: OspreySettings): Promise<Diagnostic[]> {
  connection.console.log(`📝 Analyzing document content (${textDocument.getText().length} chars)`);
  const diagnostics: Diagnostic[] = [];
  const text = textDocument.getText();
  
  try {
    // Write content to temporary file
    const outputsDir = path.join(process.cwd(), 'outputs');
    if (!fs.existsSync(outputsDir)) {
      fs.mkdirSync(outputsDir, { recursive: true });
    }
    const tempFile = path.join(outputsDir, `temp_${Date.now()}.osp`);
    connection.console.log(`💾 Writing temp file: ${tempFile}`);
    fs.writeFileSync(tempFile, text);

    // Call Osprey compiler to check for errors
    connection.console.log(`🔨 Running osprey compiler: ${settings.compilerPath}`);
    const result = await runOspreyCompiler(tempFile, settings.compilerPath);
    
    connection.console.log(`📤 Compiler stdout: ${result.stdout.substring(0, 200)}${result.stdout.length > 200 ? '...' : ''}`);
    connection.console.log(`📤 Compiler stderr: ${result.stderr.substring(0, 200)}${result.stderr.length > 200 ? '...' : ''}`);
    
    // Clean up temp file
    fs.unlinkSync(tempFile);
    connection.console.log(`🗑️ Cleaned up temp file`);

    // Parse compiler output for errors
    const errors = parseCompilerErrors(result.stderr, result.stdout, result.error, text);
    connection.console.log(`🔍 Parsed ${errors.length} errors from compiler output`);
    diagnostics.push(...errors);

  } catch (error) {
    connection.console.log(`💥 Exception during analysis: ${error}`);
    // If we can't compile, add a generic error
    diagnostics.push({
      severity: DiagnosticSeverity.Error,
      range: {
        start: { line: 0, character: 0 },
        end: { line: 0, character: 100 }
      },
      message: `Osprey compiler error: ${error}`,
      source: 'osprey'
    });
  }

  return diagnostics;
}

function runOspreyCompiler(filePath: string, compilerPath: string): Promise<{stdout: string, stderr: string, error?: Error}> {
  return new Promise((resolve) => {
    execFile('osprey', [filePath], (error: any, stdout: any, stderr: any) => {
      // Don't treat non-zero exit codes as errors - they might just be syntax errors
      resolve({ stdout, stderr, error: undefined });
    });
  });
}

async function getSymbolsFromCompiler(sourceCode: string, uri: string): Promise<OspreySymbol[]> {
  try {
    // Write source code to a temporary file
    const fs = require('fs');
    const path = require('path');
    const os = require('os');
    
    const tempDir = os.tmpdir();
    const tempFile = path.join(tempDir, `osprey_temp_${Date.now()}.osp`);
    
    fs.writeFileSync(tempFile, sourceCode);
    
    // Get compiler path from settings
    const settings = await getDocumentSettings('');
    const compilerPath = settings.compilerPath || 'osprey';
    
    // Run compiler with --symbols flag
    const result = await new Promise<{stdout: string, stderr: string, error?: Error}>((resolve) => {
      execFile(compilerPath, [tempFile, '--symbols'], (error: any, stdout: any, stderr: any) => {
        resolve({ stdout, stderr, error: error || undefined });
      });
    });
    
    // Clean up temp file
    try {
      fs.unlinkSync(tempFile);
    } catch (e) {
      // Ignore cleanup errors
    }
    
    if (result.error) {
      connection.console.log(`❌ Compiler error: ${result.error.message}`);
      return [];
    }
    
    if (result.stderr) {
      connection.console.log(`⚠️ Compiler stderr: ${result.stderr}`);
    }
    
    if (!result.stdout.trim()) {
      connection.console.log(`❌ No compiler output`);
      return [];
    }
    
    // Parse JSON output
    const symbolsData = JSON.parse(result.stdout);
    
    // Convert to OspreySymbol format
    const symbols: OspreySymbol[] = symbolsData.map((sym: any) => ({
      name: sym.name,
      type: sym.type,
      kind: sym.kind as 'function' | 'variable' | 'type' | 'parameter' | 'effect' | 'operation' | 'module' | 'import',
      location: {
        uri: uri,
        range: {
          start: { line: sym.line - 1, character: sym.column - 1 },
          end: { line: sym.line - 1, character: sym.column - 1 + sym.name.length }
        }
      },
      documentation: sym.documentation,
      signature: sym.signature,
      parameters: sym.parameters,
      returnType: sym.returnType,
      effectSet: sym.effectSet,
      isEffectOperation: sym.isEffectOperation,
      parentEffect: sym.parentEffect,
      isExported: sym.isExported,
      isImported: sym.isImported
    }));
    
    connection.console.log(`✅ Parsed ${symbols.length} symbols from compiler`);
    return symbols;
    
  } catch (error) {
    connection.console.log(`❌ Error getting symbols from compiler: ${error}`);
    return [];
  }
}

async function getBuiltinHoverDocumentation(elementName: string): Promise<string | null> {
  try {
    // Get compiler path from settings
    const settings = await getDocumentSettings('');
    const compilerPath = settings.compilerPath || 'osprey';
    
    // Call compiler with hover flag
    const result = await new Promise<{stdout: string, stderr: string, error?: Error}>((resolve) => {
      execFile(compilerPath, ['--hover', elementName], (error: any, stdout: any, stderr: any) => {
        resolve({ stdout, stderr, error: error || undefined });
      });
    });
    
    if (result.error || result.stderr) {
      connection.console.log(`❌ Error getting hover documentation for ${elementName}: ${result.stderr || result.error}`);
      return null;
    }
    
    return result.stdout.trim() || null;
  } catch (error) {
    connection.console.log(`❌ Exception getting hover documentation for ${elementName}: ${error}`);
    return null;
  }
}

function parseCompilerErrors(stderr: string, stdout: string, error: Error | undefined, documentText?: string): Diagnostic[] {
  const diagnostics: Diagnostic[] = [];
  
  // Combine stderr and stdout to look for errors
  const allOutput = `${stderr}\n${stdout}`;
  const lines = allOutput.split('\n');

  for (const line of lines) {
    if (!line.trim()) continue;
    
    // AGGRESSIVE LLVM IR DETECTION - Skip anything that looks like LLVM IR
    const trimmedLine = line.trim();
    
    // Skip any line that contains LLVM IR patterns
    if (
      // LLVM IR variable/global assignments
      trimmedLine.match(/^[@%]\w+\s*=/) ||
      
      // LLVM IR function definitions and declarations
      trimmedLine.match(/^define\s+/) ||
      trimmedLine.match(/^declare\s+/) ||
      
      // LLVM IR instructions
      trimmedLine.match(/^\s*(br|ret|call|load|store|alloca|getelementptr|icmp|add|sub|mul|div|and|or|xor)\s/) ||
      
      // LLVM IR types and casts
      trimmedLine.match(/\bi(1|8|16|32|64)\b/) ||
      trimmedLine.match(/\b(zext|sext|trunc|bitcast|inttoptr|ptrtoint)\b/) ||
      
      // LLVM IR string constants (the main culprit)
      trimmedLine.includes('x i8]') ||
      trimmedLine.includes('\\00"') ||
      trimmedLine.includes('c"') ||
      
      // LLVM IR control flow
      trimmedLine.match(/^\s*\w+:$/) || // labels
      trimmedLine.match(/^entry:$/) ||
      
      // LLVM IR metadata and attributes
      trimmedLine.match(/^target\s+(datalayout|triple)/) ||
      trimmedLine.includes('!dbg') ||
      trimmedLine.includes('attributes') ||
      
      // LLVM IR brackets and braces (common in IR)
      trimmedLine.match(/^\s*[{}]\s*$/) ||
      
      // Catch any line that has LLVM-style syntax
      trimmedLine.includes(' = global ') ||
      trimmedLine.includes(' = constant ') ||
      trimmedLine.includes('getelementptr') ||
      trimmedLine.includes('noundef') ||
      
      // Catch string literals in LLVM IR (main issue)
      (trimmedLine.includes('"') && (
        trimmedLine.includes('c"') ||
        trimmedLine.includes('= global') ||
        trimmedLine.includes('x i8]')
      ))
    ) {
      // This is LLVM IR, skip it completely
      continue;
    }
    
    // Try different error patterns (only for non-LLVM IR lines)
    let match;
    
    // Pattern: "line X:Y message" (Osprey compiler format)
    match = line.match(/line (\d+):(\d+)\s+(.+)/i);
    if (match) {
      const lineNum = Math.max(0, parseInt(match[1]) - 1); // Convert to 0-based
      const charNum = Math.max(0, parseInt(match[2]) - 1);
      const message = match[3];

      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: {
          start: { line: lineNum, character: charNum },
          end: { line: lineNum, character: Number.MAX_VALUE }
        },
        message: message,
        source: 'osprey',
        code: 'syntax-error'
      });
      continue;
    }
    
    // Pattern: "filename:line:col: message"
    match = line.match(/([^:]+):(\d+):(\d+):\s*(.+)/);
    if (match) {
      const lineNum = Math.max(0, parseInt(match[2]) - 1);
      const charNum = Math.max(0, parseInt(match[3]) - 1);
      const message = match[4];

      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: {
          start: { line: lineNum, character: charNum },
          end: { line: lineNum, character: charNum + 10 }
        },
        message: message,
        source: 'osprey'
      });
      continue;
    }
    
    // Pattern: "Error generating LLVM IR: validation error:" (Osprey specific validation errors)
    match = line.match(/Error generating LLVM IR:\s*validation error:\s*(.+)/);
    if (match) {
      const message = match[1];
      
      // Try to extract function name from validation error
      const functionMatch = message.match(/Function '(\w+)'/);
      const parameterMatch = message.match(/Parameter '(\w+)'/);
      
      // Default error range
      let errorRange = { start: { line: 0, character: 0 }, end: { line: 0, character: Number.MAX_VALUE } };
      
      // If we have the document text, try to find the exact location
      if (documentText && functionMatch) {
        const functionName = functionMatch[1];
        const lines = documentText.split('\n');
        
        for (let i = 0; i < lines.length; i++) {
          const currentLine = lines[i];
          // Look for function declaration pattern: fn functionName(
          const fnPattern = new RegExp(`\\bfn\\s+${functionName}\\s*\\(`);
          const fnMatch = currentLine.match(fnPattern);
          
          if (fnMatch) {
            const startChar = currentLine.indexOf('fn');
            const endChar = currentLine.indexOf('(') + 1;
            
            errorRange = {
              start: { line: i, character: startChar },
              end: { line: i, character: endChar }
            };
            break;
          }
        }
      }
      
      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: errorRange,
        message: message,
        source: 'osprey',
        code: 'type-inference-error'
      });
      continue;
    }
    
    // Pattern: "Error generating LLVM IR: parse errors:" (Osprey specific)
    if (line.includes('Error generating LLVM IR:') || line.includes('parse errors:')) {
      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: {
          start: { line: 0, character: 0 },
          end: { line: 0, character: Number.MAX_VALUE }
        },
        message: line.trim(),
        source: 'osprey',
        code: 'compilation-error'
      });
      continue;
    }
    
    // Pattern: "token recognition error" (ANTLR specific)
    if (line.includes('token recognition error')) {
      // Extract line number if available
      const tokenMatch = line.match(/line (\d+):(\d+)/);
      let lineNum = 0;
      let charNum = 0;
      if (tokenMatch) {
        lineNum = Math.max(0, parseInt(tokenMatch[1]) - 1);
        charNum = Math.max(0, parseInt(tokenMatch[2]) - 1);
      }
      
      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: {
          start: { line: lineNum, character: charNum },
          end: { line: lineNum, character: charNum + 5 }
        },
        message: line.trim(),
        source: 'osprey',
        code: 'token-error'
      });
      continue;
    }
    
    // Pattern: Actual error messages (but NOT LLVM IR containing "error" in strings)
    // Only match lines that start with "Error:" or "error:" (actual error messages)
    if (line.match(/^(Error|error):/i)) {
      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: {
          start: { line: 0, character: 0 },
          end: { line: 0, character: 100 }
        },
        message: line.trim(),
        source: 'osprey'
      });
      continue;
    }
    
    // Pattern: "Syntax error at line" (our error listener format)
    match = line.match(/Syntax error at line (\d+):(\d+)\s*-\s*(.+)/i);
    if (match) {
      const lineNum = Math.max(0, parseInt(match[1]) - 1);
      const charNum = Math.max(0, parseInt(match[2]) - 1);
      const message = match[3];

      diagnostics.push({
        severity: DiagnosticSeverity.Error,
        range: {
          start: { line: lineNum, character: charNum },
          end: { line: lineNum, character: charNum + 10 }
        },
        message: message,
        source: 'osprey'
      });
      continue;
    }
  }

  return diagnostics;
}

// Find all symbol references (definitions and usages) in a document
function findAllSymbolReferences(document: TextDocument): SymbolReference[] {
  connection.console.log(`🔍 Finding symbol references in ${document.uri}`);
  const references: SymbolReference[] = [];
  const text = document.getText();
  const lines = text.split('\n');
  
  connection.console.log(`📄 Document has ${lines.length} lines, ${text.length} characters`);
  
  // Regular expressions for different symbol patterns
  const patterns = {
    // Function definitions: fn name(
    functionDef: /\bfn\s+(\w+)\s*\(/g,
    // Variable definitions: let/mut name =
    variableDef: /\b(let|mut)\s+(\w+)\s*=/g,
    // Type definitions: type Name =
    typeDef: /\btype\s+(\w+)\s*=/g,
    // Function calls: name(
    functionCall: /\b(\w+)\s*\(/g,
    // Variable/type references: any word not in definition context
    symbolRef: /\b(\w+)\b/g
  };
  
  lines.forEach((line, lineIndex) => {
    connection.console.log(`📄 Processing line ${lineIndex + 1}: "${line}"`);
    
    // Track function definitions
    let match;
    const functionDefs = new Set<string>();
    const variableDefs = new Set<string>();
    const typeDefs = new Set<string>();
    
    // Find function definitions
    patterns.functionDef.lastIndex = 0;
    while ((match = patterns.functionDef.exec(line)) !== null) {
      const name = match[1];
      functionDefs.add(name);
      connection.console.log(`  🔧 Found function definition: ${name}`);
      references.push({
        symbol: name,
        location: {
          uri: document.uri,
          range: {
            start: { line: lineIndex, character: match.index + match[0].indexOf(name) },
            end: { line: lineIndex, character: match.index + match[0].indexOf(name) + name.length }
          }
        },
        kind: 'definition'
      });
    }
    
    // Find variable definitions
    patterns.variableDef.lastIndex = 0;
    while ((match = patterns.variableDef.exec(line)) !== null) {
      const name = match[2];
      variableDefs.add(name);
      connection.console.log(`  📦 Found variable definition: ${name}`);
      references.push({
        symbol: name,
        location: {
          uri: document.uri,
          range: {
            start: { line: lineIndex, character: match.index + match[0].indexOf(name) },
            end: { line: lineIndex, character: match.index + match[0].indexOf(name) + name.length }
          }
        },
        kind: 'definition'
      });
    }
    
    // Find type definitions
    patterns.typeDef.lastIndex = 0;
    while ((match = patterns.typeDef.exec(line)) !== null) {
      const name = match[1];
      typeDefs.add(name);
      connection.console.log(`  🏷️ Found type definition: ${name}`);
      references.push({
        symbol: name,
        location: {
          uri: document.uri,
          range: {
            start: { line: lineIndex, character: match.index + match[0].indexOf(name) },
            end: { line: lineIndex, character: match.index + match[0].indexOf(name) + name.length }
          }
        },
        kind: 'definition'
      });
    }
    
    // Find all symbol references (usages)
    patterns.symbolRef.lastIndex = 0;
    while ((match = patterns.symbolRef.exec(line)) !== null) {
      const name = match[1];
      const startChar = match.index;
      
      // Skip if this is a keyword
      const keywords = ['fn', 'let', 'mut', 'type', 'match', 'if', 'then', 'else', 'case', 'of', 'import', 'extern'];
      if (keywords.includes(name)) continue;
      
      // Skip if this is part of a definition we already tracked
      if (functionDefs.has(name) || variableDefs.has(name) || typeDefs.has(name)) {
        // Check if this occurrence is the definition itself
        const defRef = references.find(r => 
          r.symbol === name && 
          r.kind === 'definition' && 
          r.location.range.start.line === lineIndex &&
          r.location.range.start.character === startChar
        );
        if (defRef) continue;
      }
      
      // This is a usage
      connection.console.log(`  👁️ Found symbol usage: ${name} at character ${startChar}`);
      references.push({
        symbol: name,
        location: {
          uri: document.uri,
          range: {
            start: { line: lineIndex, character: startChar },
            end: { line: lineIndex, character: startChar + name.length }
          }
        },
        kind: 'usage'
      });
    }
  });
  
  connection.console.log(`✅ Found ${references.length} total symbol references`);
  return references;
}

// Provide completion items
connection.onCompletion((_textDocumentPosition: TextDocumentPositionParams): CompletionItem[] => {
  return [
    {
      label: 'fn',
      kind: CompletionItemKind.Keyword,
      data: 1,
      detail: 'Function declaration',
      insertText: 'fn ${1:name}(${2:params}) = ${3:body}'
    },
    {
      label: 'let',
      kind: CompletionItemKind.Keyword,
      data: 2,
      detail: 'Variable declaration',
      insertText: 'let ${1:name} = ${2:value}'
    },
    {
      label: 'mut',
      kind: CompletionItemKind.Keyword,
      data: 3,
      detail: 'Mutable variable declaration',
      insertText: 'mut ${1:name} = ${2:value}'
    },
    {
      label: 'match',
      kind: CompletionItemKind.Keyword,
      data: 4,
      detail: 'Pattern matching',
      insertText: 'match ${1:expr} {\n\t${2:pattern} => ${3:result}\n}'
    },
    {
      label: 'type',
      kind: CompletionItemKind.Keyword,
      data: 5,
      detail: 'Type declaration',
      insertText: 'type ${1:Name} = ${2:Variant} | ${3:Variant}'
    },
    {
      label: 'print',
      kind: CompletionItemKind.Function,
      data: 6,
      detail: 'Print function',
      insertText: 'print(${1:value})'
    }
  ];
});

connection.onCompletionResolve((item: CompletionItem): CompletionItem => {
  return item;
});

// HOVER PROVIDER - SHOW TYPE INFORMATION
connection.onHover(async (params): Promise<Hover | null> => {
  connection.console.log(`🎯 HOVER REQUEST RECEIVED at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) {
    connection.console.log(`❌ Document not found for hover`);
    return null;
  }
  
  const text = document.getText();
  const lines = text.split('\n');
  const currentLine = lines[params.position.line];
  if (!currentLine) {
    connection.console.log(`❌ No line found at position ${params.position.line}`);
    return null;
  }
  
  connection.console.log(`📄 Current line: "${currentLine}"`);
  connection.console.log(`📍 Character position: ${params.position.character}`);
  
  // Get word at cursor position
  const wordAtPosition = getWordAtPosition(currentLine, params.position.character);
  if (!wordAtPosition) {
    connection.console.log(`❌ No word found at cursor position`);
    return null;
  }
  
  connection.console.log(`🔍 Word at position: "${wordAtPosition.word}" (${wordAtPosition.start}-${wordAtPosition.end})`);
  
  // Get symbol information from the compiler
  try {
    const symbols = await getSymbolsFromCompiler(text, params.textDocument.uri);
    connection.console.log(`📊 Got ${symbols.length} symbols from compiler`);
    
    // First, try to find exact symbol match
    for (const symbol of symbols) {
      if (symbol.name === wordAtPosition.word) {
        connection.console.log(`✅ Found exact symbol match: ${symbol.name} (${symbol.type})`);
        return createHoverContent(symbol, params.position, currentLine);
      }
    }
    
    // If no exact match, check if we're hovering over a function call or variable reference
    const hoverInfo = await analyzeWordContext(wordAtPosition.word, currentLine, params.position.character, symbols);
    if (hoverInfo) {
      connection.console.log(`✅ Found contextual hover info for: ${wordAtPosition.word}`);
      return hoverInfo;
    }
    
    connection.console.log(`❌ No symbol information found for: ${wordAtPosition.word}`);
    return null;
  } catch (error) {
    connection.console.log(`❌ Error getting symbols from compiler: ${error}`);
    return null;
  }
});

function getWordAtPosition(line: string, character: number): { word: string; start: number; end: number } | null {
  // Find word boundaries around the character position
  let start = character;
  let end = character;
  
  // Move start backwards to find word start
  while (start > 0 && /[a-zA-Z0-9_]/.test(line[start - 1])) {
    start--;
  }
  
  // Move end forwards to find word end
  while (end < line.length && /[a-zA-Z0-9_]/.test(line[end])) {
    end++;
  }
  
  if (start === end) {
    return null;
  }
  
  const word = line.substring(start, end);
  return { word, start, end };
}

function createHoverContent(symbol: OspreySymbol, position: Position, currentLine: string): Hover {
  let content = `**${symbol.name}**\n\n`;
  
  // Add kind badge
  const kindBadge = symbol.kind === 'function' ? '🔧' : symbol.kind === 'variable' ? '📦' : '🏷️';
  content += `${kindBadge} *${symbol.kind.charAt(0).toUpperCase() + symbol.kind.slice(1)}*\n\n`;
  
  // Add type information
  content += `**Type:** \`${symbol.type}\`\n\n`;
  
  // Add signature for functions
  if (symbol.signature) {
    content += `**Signature:**\n\`\`\`osprey\n${symbol.signature}\n\`\`\`\n\n`;
  }
  
  // Add parameters for functions
  if (symbol.parameters && symbol.parameters.length > 0) {
    content += `**Parameters:**\n`;
    symbol.parameters.forEach(param => {
      content += `- **\`${param.name}\`**: \`${param.type}\``;
      if (param.documentation) {
        content += ` - ${param.documentation}`;
      }
      content += '\n';
    });
    content += '\n';
  }
  
  // Add return type for functions
  if (symbol.returnType && symbol.kind === 'function') {
    content += `**Returns:** \`${symbol.returnType}\`\n\n`;
  }

  // Add effect set for functions
  if (symbol.effectSet) {
    content += `**Effect Set:** \`${symbol.effectSet}\`\n\n`;
  }

  // Add documentation
  if (symbol.documentation) {
    content += `**Description:**\n${symbol.documentation}\n\n`;
  }
  
  // Add context information
  content += `---\n*Line ${position.line + 1}, Column ${position.character + 1}*`;
  
  return {
    contents: {
      kind: MarkupKind.Markdown,
      value: content
    }
  };
}

async function analyzeWordContext(word: string, line: string, character: number, symbols: OspreySymbol[]): Promise<Hover | null> {
  // Check if this is a function call (word followed by parentheses)
  const afterWord = line.substring(character);
  const beforeWord = line.substring(0, character - word.length);
  
  // Function call pattern
  if (afterWord.match(/^\s*\(/)) {
    // Find function definition
    const functionSymbol = symbols.find(s => s.name === word && s.kind === 'function');
    if (functionSymbol) {
      return createHoverContent(functionSymbol, { line: 0, character: 0 }, line);
    }
  }
  
  // Check for pipe operator
  if (line.includes('|>') && character >= line.indexOf('|>') && character <= line.indexOf('|>') + 2) {
    const hoverDoc = await getBuiltinHoverDocumentation('|>');
    if (hoverDoc) {
      return {
        contents: {
          kind: MarkupKind.Markdown,
          value: hoverDoc
        }
      };
    }
  }
  
  // Variable reference
  const variableSymbol = symbols.find(s => s.name === word && s.kind === 'variable');
  if (variableSymbol) {
    return createHoverContent(variableSymbol, { line: 0, character: 0 }, line);
  }
  
  // Type reference
  const typeSymbol = symbols.find(s => s.name === word && s.kind === 'type');
  if (typeSymbol) {
    return createHoverContent(typeSymbol, { line: 0, character: 0 }, line);
  }
  
  // Check for built-in language elements (keywords, operators, types, functions)
  const hoverDoc = await getBuiltinHoverDocumentation(word);
  if (hoverDoc) {
    return {
      contents: {
        kind: MarkupKind.Markdown,
        value: hoverDoc
      }
    };
  }
  
  return null;
}

// GO TO DEFINITION
connection.onDefinition((params): Definition | null => {
  connection.console.log(`🎯 GO TO DEFINITION REQUEST RECEIVED at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) {
    connection.console.log(`❌ Document not found: ${params.textDocument.uri}`);
    return null;
  }
  
  // Get word at position using the existing function
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) {
    connection.console.log(`❌ No line at position ${params.position.line}`);
    return null;
  }
  
  connection.console.log(`📄 Line content: "${line}"`);
  
  const wordAtPosition = getWordAtPosition(line, params.position.character);
  if (!wordAtPosition) {
    connection.console.log(`❌ No word at cursor position ${params.position.character}`);
    return null;
  }
  
  const targetWord = wordAtPosition.word;
  connection.console.log(`🔍 Looking for definition of: "${targetWord}"`);
  
  // Debug: log all current references
  const currentDocRefs = symbolReferences.get(params.textDocument.uri) || [];
  connection.console.log(`📊 Current document has ${currentDocRefs.length} symbol references`);
  currentDocRefs.forEach((ref, i) => {
    connection.console.log(`  ${i+1}. ${ref.kind} "${ref.symbol}" at line ${ref.location.range.start.line + 1}`);
  });
  
  // First, check if we're already on a definition
  const currentRef = currentDocRefs.find(ref => 
    ref.symbol === targetWord &&
    ref.location.range.start.line === params.position.line &&
    params.position.character >= ref.location.range.start.character &&
    params.position.character <= ref.location.range.end.character
  );
  
  if (currentRef && currentRef.kind === 'definition') {
    connection.console.log(`✋ Already on definition of ${targetWord}`);
    return null; // Already on the definition
  }
  
  // Search for the definition in the current document first
  const definitionInCurrentDoc = currentDocRefs.find(ref => 
    ref.symbol === targetWord && ref.kind === 'definition'
  );
  
  if (definitionInCurrentDoc) {
    connection.console.log(`✅ Found definition in current document at line ${definitionInCurrentDoc.location.range.start.line + 1}`);
    return definitionInCurrentDoc.location;
  }
  
  // Search in workspace symbols (includes imports and external symbols)
  connection.console.log(`🔍 Searching workspace symbols...`);
  let foundInWorkspace = false;
  for (const [uri, symbols] of workspaceSymbols) {
    connection.console.log(`  📁 Checking ${uri} with ${symbols.length} symbols`);
    for (const symbol of symbols) {
      if (symbol.name === targetWord) {
        connection.console.log(`✅ Found definition in ${uri} at line ${symbol.location.range.start.line + 1}`);
        foundInWorkspace = true;
        return symbol.location;
      }
    }
  }
  
  if (!foundInWorkspace) {
    connection.console.log(`🔍 No workspace symbols found, searching all references...`);
  }
  
  // Search in all document references as fallback
  for (const [uri, refs] of symbolReferences) {
    const def = refs.find(ref => ref.symbol === targetWord && ref.kind === 'definition');
    if (def) {
      connection.console.log(`✅ Found definition via references in ${uri}`);
      return def.location;
    }
  }
  
  connection.console.log(`❌ No definition found for "${targetWord}"`);
  return null;
});

// GO TO IMPLEMENTATION
connection.onImplementation((params: ImplementationParams): Location[] | null => {
  connection.console.log(`🎯 GO TO IMPLEMENTATION REQUEST RECEIVED at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) {
    connection.console.log(`❌ Document not found: ${params.textDocument.uri}`);
    return null;
  }
  
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) {
    connection.console.log(`❌ No line at position ${params.position.line}`);
    return null;
  }
  
  connection.console.log(`📄 Line content: "${line}"`);
  
  const wordAtPosition = getWordAtPosition(line, params.position.character);
  if (!wordAtPosition) {
    connection.console.log(`❌ No word at cursor position ${params.position.character}`);
    return null;
  }
  
  const targetWord = wordAtPosition.word;
  connection.console.log(`🔍 Looking for implementation of: "${targetWord}"`);
  
  // Search in all document references
  for (const [uri, refs] of symbolReferences) {
    for (const ref of refs) {
      if (ref.symbol === targetWord && ref.kind === 'implementation') {
        connection.console.log(`✅ Found implementation at ${uri}:${ref.location.range.start.line + 1}`);
        return [ref.location];
      }
    }
  }
  
  // Search in workspace symbols
  for (const [uri, symbols] of workspaceSymbols) {
    for (const symbol of symbols) {
      if (symbol.name === targetWord && symbol.kind === 'implementation') {
        connection.console.log(`✅ Found implementation at ${uri}:${symbol.location.range.start.line + 1}`);
        return [symbol.location];
      }
    }
  }
  
  connection.console.log(`❌ No implementation found for "${targetWord}"`);
  return null;
});

// DOCUMENT SYMBOLS (Outline)
connection.onDocumentSymbol((params): DocumentSymbol[] => {
  connection.console.log(`📋 Document symbols request for ${params.textDocument.uri}`);
  
  const analysis = documentAnalyses.get(params.textDocument.uri);
  if (!analysis) {
    connection.console.log(`❌ No analysis found for document symbols`);
    return [];
  }
  
  const documentSymbols: DocumentSymbol[] = analysis.symbols.map(symbol => {
    let symbolKind: SymbolKind;
    switch (symbol.kind) {
      case 'function':
        symbolKind = SymbolKind.Function;
        break;
      case 'variable':
        symbolKind = SymbolKind.Variable;
        break;
      case 'type':
        symbolKind = SymbolKind.Class;
        break;
      case 'parameter':
        symbolKind = SymbolKind.Variable;
        break;
      case 'effect':
        symbolKind = SymbolKind.Enum;
        break;
      case 'operation':
        symbolKind = SymbolKind.Function;
        break;
      case 'module':
        symbolKind = SymbolKind.Module;
        break;
      case 'import':
        symbolKind = SymbolKind.Import;
        break;
      default:
        symbolKind = SymbolKind.Variable;
    }
    
    return {
      name: symbol.name,
      detail: symbol.type,
      kind: symbolKind,
      range: symbol.location.range,
      selectionRange: symbol.location.range
    };
  });
  
  connection.console.log(`✅ Returning ${documentSymbols.length} document symbols`);
  return documentSymbols;
});

// SIGNATURE HELP
connection.onSignatureHelp((params): SignatureHelp | null => {
  connection.console.log(`✍️ Signature help request at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return null;
  
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) return null;
  
  // Find function call at current position
  const beforeCursor = line.substring(0, params.position.character);
  const functionMatch = beforeCursor.match(/(\w+)\s*\(\s*([^)]*)$/);
  
  if (functionMatch) {
    const [, functionName] = functionMatch;
    connection.console.log(`🔍 Looking for signature of: ${functionName}`);
    
    // Find function in workspace symbols (includes built-ins from compiler)
    for (const [uri, symbols] of workspaceSymbols) {
      for (const symbol of symbols) {
        if (symbol.name === functionName && symbol.kind === 'function') {
          connection.console.log(`✅ Found function signature for ${functionName}`);
          
          const parameters: ParameterInformation[] = symbol.parameters?.map(param => ({
            label: `${param.name}: ${param.type}`,
            documentation: param.documentation || `Parameter ${param.name} of type ${param.type}`
          })) || [];
          
          const signature: SignatureInformation = {
            label: symbol.signature || `${functionName}(${parameters.map(p => p.label).join(', ')})`,
            documentation: symbol.documentation,
            parameters
          };
          
          // Calculate active parameter based on comma count
          const parameterText = functionMatch[2] || '';
          const commaCount = (parameterText.match(/,/g) || []).length;
          const activeParameter = Math.min(commaCount, parameters.length - 1);
          
          return {
            signatures: [signature],
            activeSignature: 0,
            activeParameter: Math.max(0, activeParameter)
          };
        }
      }
    }
  }
  
  connection.console.log(`❌ No signature help found`);
  return null;
});

// CODE ACTIONS
connection.onCodeAction((params): CodeAction[] => {
  connection.console.log(`🔧 Code action request for ${params.textDocument.uri}`);
  
  const actions: CodeAction[] = [];
  const diagnostics = params.context.diagnostics;
  
  for (const diagnostic of diagnostics) {
    // Quick fix for common errors
    if (diagnostic.message.includes('undefined variable')) {
      const variableName = diagnostic.message.match(/undefined variable '([^']+)'/)?.[1];
      if (variableName) {
        actions.push({
          title: `Add 'let ${variableName}' declaration`,
          kind: CodeActionKind.QuickFix,
          diagnostics: [diagnostic],
          edit: {
            changes: {
              [params.textDocument.uri]: [{
                range: diagnostic.range,
                newText: `let ${variableName} = \n`
              }]
            }
          }
        });
      }
    }
    
    // Quick fix for missing semicolons
    if (diagnostic.message.includes('expected')) {
      actions.push({
        title: 'Add missing semicolon',
        kind: CodeActionKind.QuickFix,
        diagnostics: [diagnostic],
        edit: {
          changes: {
            [params.textDocument.uri]: [{
              range: diagnostic.range,
              newText: ';'
            }]
          }
        }
      });
    }
  }
  
  // Add organize imports action
  actions.push({
    title: 'Organize Imports',
    kind: CodeActionKind.SourceOrganizeImports,
    edit: {
      changes: {
        [params.textDocument.uri]: [] // Would implement import organization logic
      }
    }
  });
  
  connection.console.log(`✅ Returning ${actions.length} code actions`);
  return actions;
});

// DOCUMENT FORMATTING
connection.onDocumentFormatting((params: DocumentFormattingParams): TextEdit[] => {
  connection.console.log(`📝 Document formatting request for ${params.textDocument.uri}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return [];
  
  const text = document.getText();
  const lines = text.split('\n');
  const edits: TextEdit[] = [];
  
  // Basic formatting: consistent indentation
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    const trimmed = line.trim();
    
    if (trimmed.length > 0) {
      // Calculate proper indentation based on context
      let indentLevel = 0;
      
      // Increase indent for blocks
      if (trimmed.includes('{') && !trimmed.includes('}')) {
        indentLevel = 1;
      }
      
      // Decrease indent for closing braces
      if (trimmed.includes('}') && !trimmed.includes('{')) {
        indentLevel = -1;
      }
      
      const properIndent = '  '.repeat(Math.max(0, indentLevel));
      const formattedLine = properIndent + trimmed;
      
      if (formattedLine !== line) {
        edits.push({
          range: {
            start: { line: i, character: 0 },
            end: { line: i, character: line.length }
          },
          newText: formattedLine
        });
      }
    }
  }
  
  connection.console.log(`✅ Applied ${edits.length} formatting edits`);
  return edits;
});

// DOCUMENT RANGE FORMATTING
connection.onDocumentRangeFormatting((params: DocumentRangeFormattingParams): TextEdit[] => {
  connection.console.log(`📝 Range formatting request for ${params.textDocument.uri}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return [];
  
  const text = document.getText();
  const lines = text.split('\n');
  const edits: TextEdit[] = [];
  
  const startLine = params.range.start.line;
  const endLine = params.range.end.line;
  
  // Format only the specified range
  for (let i = startLine; i <= endLine; i++) {
    const line = lines[i];
    const trimmed = line.trim();
    
    if (trimmed.length > 0) {
      const properIndent = '  '; // Simple 2-space indentation
      const formattedLine = properIndent + trimmed;
      
      if (formattedLine !== line) {
        edits.push({
          range: {
            start: { line: i, character: 0 },
            end: { line: i, character: line.length }
          },
          newText: formattedLine
        });
      }
    }
  }
  
  connection.console.log(`✅ Applied ${edits.length} range formatting edits`);
  return edits;
});

// DOCUMENT HIGHLIGHTS
connection.onDocumentHighlight((params): DocumentHighlight[] => {
  connection.console.log(`🔍 Document highlight request at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return [];
  
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) return [];
  
  const wordAtPosition = getWordAtPosition(line, params.position.character);
  if (!wordAtPosition) return [];
  
  const targetWord = wordAtPosition.word;
  const highlights: DocumentHighlight[] = [];
  
  // Find all occurrences of the word in the document
  for (let i = 0; i < lines.length; i++) {
    const currentLine = lines[i];
    let startIndex = 0;
    
    while (true) {
      const index = currentLine.indexOf(targetWord, startIndex);
      if (index === -1) break;
      
      // Check if it's a complete word match
      const before = index > 0 ? currentLine[index - 1] : ' ';
      const after = index + targetWord.length < currentLine.length ? currentLine[index + targetWord.length] : ' ';
      
      if (!/[a-zA-Z0-9_]/.test(before) && !/[a-zA-Z0-9_]/.test(after)) {
        highlights.push({
          range: {
            start: { line: i, character: index },
            end: { line: i, character: index + targetWord.length }
          },
          kind: 1 // Read
        });
      }
      
      startIndex = index + 1;
    }
  }
  
  connection.console.log(`✅ Found ${highlights.length} highlights for "${targetWord}"`);
  return highlights;
});

// FOLDING RANGES
connection.onFoldingRanges((params): FoldingRange[] => {
  connection.console.log(`📁 Folding ranges request for ${params.textDocument.uri}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return [];
  
  const text = document.getText();
  const lines = text.split('\n');
  const ranges: FoldingRange[] = [];
  
  let braceStack: number[] = [];
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    
    // Track opening braces
    for (let j = 0; j < line.length; j++) {
      if (line[j] === '{') {
        braceStack.push(i);
      } else if (line[j] === '}') {
        if (braceStack.length > 0) {
          const startLine = braceStack.pop()!;
          if (i > startLine) {
            ranges.push({
              startLine,
              endLine: i,
              kind: 'region'
            });
          }
        }
      }
    }
  }
  
  // Add function folding
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line.trim().startsWith('fn ')) {
      // Find the end of the function
      let braceCount = 0;
      let foundStart = false;
      
      for (let j = i; j < lines.length; j++) {
        const currentLine = lines[j];
        
        for (let k = 0; k < currentLine.length; k++) {
          if (currentLine[k] === '{') {
            braceCount++;
            foundStart = true;
          } else if (currentLine[k] === '}') {
            braceCount--;
            if (foundStart && braceCount === 0) {
              if (j > i) {
                ranges.push({
                  startLine: i,
                  endLine: j,
                  kind: 'function'
                });
              }
              break;
            }
          }
        }
        
        if (foundStart && braceCount === 0) break;
      }
    }
  }
  
  connection.console.log(`✅ Found ${ranges.length} folding ranges`);
  return ranges;
});

// WORKSPACE SYMBOLS
connection.onWorkspaceSymbol((params): WorkspaceSymbol[] => {
  connection.console.log(`🔍 Workspace symbol search for: "${params.query}"`);
  
  const symbols: WorkspaceSymbol[] = [];
  
  for (const [uri, docSymbols] of workspaceSymbols) {
    for (const symbol of docSymbols) {
      if (symbol.name.toLowerCase().includes(params.query.toLowerCase())) {
        let symbolKind: SymbolKind;
        switch (symbol.kind) {
          case 'function':
            symbolKind = SymbolKind.Function;
            break;
          case 'variable':
            symbolKind = SymbolKind.Variable;
            break;
          case 'type':
            symbolKind = SymbolKind.Class;
            break;
          case 'effect':
            symbolKind = SymbolKind.Enum;
            break;
          case 'operation':
            symbolKind = SymbolKind.Function;
            break;
          case 'module':
            symbolKind = SymbolKind.Module;
            break;
          default:
            symbolKind = SymbolKind.Variable;
        }
        
        symbols.push({
          name: symbol.name,
          kind: symbolKind,
          location: symbol.location,
          containerName: path.basename(uri)
        });
      }
    }
  }
  
  connection.console.log(`✅ Found ${symbols.length} workspace symbols`);
  return symbols;
});

// PREPARE RENAME
connection.onPrepareRename((params: PrepareRenameParams): Range | null => {
  connection.console.log(`✏️ Prepare rename request at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return null;
  
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) return null;
  
  const wordAtPosition = getWordAtPosition(line, params.position.character);
  if (!wordAtPosition) return null;
  
  return {
    start: { line: params.position.line, character: wordAtPosition.start },
    end: { line: params.position.line, character: wordAtPosition.end }
  };
});

// RENAME
connection.onRename((params: RenameParams): WorkspaceEdit | null => {
  connection.console.log(`✏️ Rename request for "${params.newName}" at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) return null;
  
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) return null;
  
  const wordAtPosition = getWordAtPosition(line, params.position.character);
  if (!wordAtPosition) return null;
  
  const oldName = wordAtPosition.word;
  const changes: { [uri: string]: TextEdit[] } = {};
  
  // Find all references and rename them
  for (const [uri, refs] of symbolReferences) {
    const uriChanges: TextEdit[] = [];
    
    for (const ref of refs) {
      if (ref.symbol === oldName) {
        uriChanges.push({
          range: ref.location.range,
          newText: params.newName
        });
      }
    }
    
    if (uriChanges.length > 0) {
      changes[uri] = uriChanges;
    }
  }
  
  connection.console.log(`✅ Renamed "${oldName}" to "${params.newName}" in ${Object.keys(changes).length} files`);
  return { changes };
});

// FIND REFERENCES
connection.onReferences((params) => {
  connection.console.log(`🔗 Find references request at ${params.textDocument.uri}:${params.position.line}:${params.position.character}`);
  
  const document = documents.get(params.textDocument.uri);
  if (!document) {
    connection.console.log(`❌ Document not found: ${params.textDocument.uri}`);
    return [];
  }
  
  // Get word at position using the existing function
  const text = document.getText();
  const lines = text.split('\n');
  const line = lines[params.position.line];
  if (!line) {
    connection.console.log(`❌ No line at position ${params.position.line}`);
    return [];
  }
  
  const wordAtPosition = getWordAtPosition(line, params.position.character);
  if (!wordAtPosition) {
    connection.console.log(`❌ No word at cursor position`);
    return [];
  }
  
  const targetWord = wordAtPosition.word;
  connection.console.log(`🔍 Finding all references for: "${targetWord}"`);
  
  const references: Location[] = [];
  const includeDeclaration = params.context.includeDeclaration;
  
  // Search in all document references
  for (const [uri, docRefs] of symbolReferences) {
    for (const ref of docRefs) {
      if (ref.symbol === targetWord) {
        // Include or exclude declarations based on context
        if (includeDeclaration || ref.kind === 'usage') {
          references.push(ref.location);
          connection.console.log(`  📍 Found ${ref.kind} at ${uri}:${ref.location.range.start.line + 1}:${ref.location.range.start.character + 1}`);
        }
      }
    }
  }
  
  // Also check workspace symbols for any we might have missed
  if (includeDeclaration) {
    for (const [uri, symbols] of workspaceSymbols) {
      for (const symbol of symbols) {
        if (symbol.name === targetWord) {
          // Check if we already have this location
          const alreadyIncluded = references.some(loc => 
            loc.uri === symbol.location.uri &&
            loc.range.start.line === symbol.location.range.start.line &&
            loc.range.start.character === symbol.location.range.start.character
          );
          
          if (!alreadyIncluded) {
            references.push(symbol.location);
            connection.console.log(`  📍 Found symbol definition at ${uri}:${symbol.location.range.start.line + 1}`);
          }
        }
      }
    }
  }
  
  connection.console.log(`✅ Found ${references.length} references for "${targetWord}"`);
  return references;
});

// Make the text document manager listen on the connection
documents.listen(connection);

// Listen on the connection
connection.listen(); 