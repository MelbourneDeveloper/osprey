# Osprey VSCode Extension Features

## 🎨 Syntax Highlighting

Based on the Osprey tree-sitter grammar, provides rich syntax highlighting for:

### Keywords
- **Control flow**: `match`, `if`, `else`, `loop`, `select`
- **Declarations**: `fn`, `let`, `mut`, `type`, `import`, `extern`, `module`
- **Fiber operations**: `spawn`, `await`, `yield`, `channel`, `send`, `recv`
- **Other**: `return`

### Literals
- **Numbers**: `42`, `123`, `123.45`
- **Booleans**: `true`, `false`
- **Strings**: `"hello world"`
- **Interpolated strings**: `"Hello ${name}!"`

### Functions
- **Declarations**: `fn add(x, y) = x + y`
- **Calls**: `add(x: 1, y: 2)`
- **Built-ins**: `print(value)`

### Types
- **Definitions**: `type Result = Ok { value: String } | Error { message: String }`
- **Variants**: `Ok`, `Error`, `Some`, `None`
- **Pattern matching**: `match expr { Ok => "success" }`

### Operators
- **Arithmetic**: `+`, `-`, `*`, `/`
- **Assignment**: `=`
- **Arrows**: `=>`, `|->`
- **Pipe**: `|>`
- **Punctuation**: `()`, `[]`, `{}`, `,`, `;`, `:`

## 🔍 Error Diagnostics

Real-time compilation error detection:

- **Parse errors**: Syntax mistakes highlighted immediately
- **Type errors**: Function call mismatches (e.g., missing named arguments)
- **Semantic errors**: Undefined variables, incorrect patterns
- **Line-precise**: Errors shown exactly where they occur

### Example Error Detection
```osprey
fn add(x, y) = x + y
let result = add(5, 10)  // ❌ Error: named arguments required
```

## 💡 Code Completion

Intelligent autocompletion for:

### Keywords
- `fn` → `fn ${1:name}(${2:params}) = ${3:body}`
- `let` → `let ${1:name} = ${2:value}`
- `match` → `match ${1:expr} { ${2:pattern} => ${3:result} }`
- `type` → `type ${1:Name} = ${2:Variant} | ${3:Variant}`

### Built-in Functions
- `print` → `print(${1:value})`

### Trigger Characters
- `.` for method chaining
- `:` for named arguments
- `$` for string interpolation
- `(` for function parameters
- `|` for pipe operations

## 🔬 Advanced Language Features

### Hover Information
- **Type information**: Hover over variables to see their types
- **Function signatures**: Detailed parameter and return type info
- **Documentation**: Built-in function documentation
- **Pipe operator help**: Comprehensive `|>` operator documentation

### Signature Help
- **Function parameters**: Shows parameter names and types as you type
- **Named arguments**: Helps with required named parameter syntax
- **Trigger on**: `(` and `,` characters

### Document Symbols
- **Function outline**: Quick navigation to functions
- **Type definitions**: Jump to type declarations
- **Symbol hierarchy**: Organized code structure view

## 🛠️ Language Features

### Bracket Matching
- Auto-closing: `()`, `[]`, `{}`, `""`
- Bracket highlighting
- Smart indentation

### Comments
- Line comments: `// comment`
- Comment toggling with Ctrl+/
- Syntax highlighting in comments

### String Interpolation
- Syntax highlighting inside `${...}`
- Nested expression support
- Proper escaping

## ⚙️ Configuration

Customizable via VSCode settings:

```json
{
  "osprey.server.enabled": true,
  "osprey.server.path": "/custom/path/to/osprey",
  "osprey.diagnostics.enabled": true,
  "osprey.server.compilerPath": "osprey"
}
```

## 🚀 Performance

- **Lightweight**: TypeScript handles UI, the Rust compiler handles computation
- **Fast**: Incremental compilation on document changes
- **Responsive**: Non-blocking error checking
- **Memory efficient**: Temporary files cleaned up immediately

## 🔧 Development Features

### Commands
- **Compile**: `Ctrl+Shift+B` / `Cmd+Shift+B` - Compile current file
- **Run**: `F5` - Compile and run current file
- **Set Language**: Force language association for `.osp` files

### Status Bar
- Shows "✅ Osprey" when language server is running
- Click for server information

### Output Panel
- "Osprey Language Server" channel for debugging
- Compilation errors and warnings
- Server startup/shutdown logs

### File Association
- Automatic `.osp` file recognition
- Proper language mode activation
- Icon association (if configured)

## 📊 Implementation Status

### ✅ **FIXED AND WORKING**

#### Core Extension Features
- **Extension Activation**: Extension properly activates when opening .osp files
- **Language Detection**: Files are correctly detected as Osprey language  
- **Syntax Highlighting**: Works correctly with the TextMate grammar
- **Commands**: Compile and run commands are available and working
- **Configuration**: Extension configuration is accessible and functional

#### Language Server Infrastructure  
- **Server Startup**: Language server starts successfully
- **Document Management**: Text documents are properly tracked
- **Diagnostics**: Syntax errors are detected and reported
- **Hover Information**: Works for built-in functions and language constructs
- **Pipe Operator Documentation**: Comprehensive hover documentation for `|>`
- **Signature Help**: Function signatures are provided for built-in functions
- **Code Completion**: Basic completion for keywords and built-in functions

### ⚠️ **PARTIALLY WORKING**

#### Symbol Information
- **Compiler Integration**: Osprey compiler `--symbols` flag works correctly
- **Symbol Parsing**: `getSymbolsFromCompiler()` function correctly parses JSON output
- **Symbol Tracking**: `findAllSymbolReferences()` function correctly identifies definitions and usages

### ❌ **NOT WORKING (Main Issues)**

#### Core Language Features
- **Go to Definition**: VSCode not calling our `onDefinition` handler
- **Find All References**: VSCode not calling our `onReferences` handler  
- **Document Symbols**: Limited symbol information in outline view

### 🔧 **ROOT CAUSE ANALYSIS**

#### Issue: Language Server Protocol Integration
The main problem is that VSCode isn't routing language feature requests to our language server handlers, despite:
1. ✅ Language server registering definition/references capabilities
2. ✅ Extension properly starting the language server
3. ✅ Server responding to other LSP requests (hover, completion, diagnostics)

#### Possible Causes
1. **Request Routing**: VSCode may not be sending `textDocument/definition` requests to our server
2. **Capability Registration**: Definition provider capability may not be properly registered
3. **Document URI Handling**: URI format mismatches between client and server
4. **Timing Issues**: Requests arriving before symbol parsing is complete

## 🎯 Future Roadmap

### Next Priority Fixes
1. **Go-to Definition**: Debug LSP communication for definition requests
2. **Find References**: Fix VSCode routing for reference requests
3. **Document Symbols**: Enhanced symbol information in outline view

### Planned Features
1. **Rename Symbol**: Intelligent renaming across files
2. **Semantic Highlighting**: Advanced syntax coloring
3. **Code Formatting**: Auto-format Osprey code
4. **Debugging Support**: Integrated debugger
5. **REPL Integration**: Interactive Osprey shell

### Advanced Features
1. **Call Hierarchy**: Function call trees
2. **Code Actions**: Quick fixes and refactoring
3. **Workspace Symbols**: Project-wide symbol search
4. **Multi-file Analysis**: Cross-file type checking

## 📋 Requirements

- **VSCode**: 1.96.0 or higher
- **Node.js**: 20.19.2 (exact version required)
- **Osprey Compiler**: Rust-based compiler (`osprey`) in PATH or bundled with the extension
- **Operating System**: Windows, macOS, Linux

## 🐛 Known Issues

1. **Large Files**: May be slow on files >1000 lines (delegated to the compiler)
2. **Complex Types**: Advanced type inference not yet implemented
3. **Multi-file**: Cross-file analysis limited
4. **Go to Definition**: VSCode LSP integration needs debugging
5. **Find References**: Handler not being called by VSCode

## 🔍 Debugging Information

### Compiler Integration Working
```bash
$ echo 'fn double(x) = x * 2\nlet result = double(5)' | osprey --symbols
[
  {
    "name": "double",
    "kind": "function", 
    "type": "Function(x: Any) -> Any",
    "line": 1,
    "column": 1,
    ...
  }
]
```

### Extension Infrastructure Working
- Language server starts: ✅
- Documents are tracked: ✅  
- Other LSP features work: ✅
- Only definition/references broken: ❌

## 💬 Feedback

Report issues and feature requests at the Osprey repository. The extension is designed to grow with the language!

**Status**: 🟡 Core features working, advanced navigation features need LSP debugging 