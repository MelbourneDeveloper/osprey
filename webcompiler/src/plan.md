# Osprey Web Compiler Implementation Plan

## 🎯 GOAL: Implement Working Web API

The web API must return clean JSON responses with:
1. `compilerOutput` - Compiler messages (errors only, success messages suppressed)
2. `programOutput` - Actual program output (what the Osprey program prints)

### 🔥 REMAINING TODOs IN SERVER.OSP

1. **JSON Request Parsing**: Need to extract `code` field from POST body
2. **Proper File I/O**: Write code to temp file and execute
3. **Output Separation**: run the compiler with quiet mode
4. **JSON Response Formatting**: Escape strings properly for JSON

## 📋 IMPLEMENTATION PLAN

### Phase 1: Fix JSON Request Parsing
```osprey
// CURRENT BROKEN CODE:
let actualCode = match extractResult {
    Success { value } => {
        //TODO: get the actual code from the request  ← NEEDS IMPLEMENTATION
    }
    Error { message } => "print(\"Failed to extract code from JSON!\")"
}

// WHAT IT SHOULD BE:
let actualCode = match extractResult {
    Success { value } => {
        // Read the extracted code from jq output
        let readResult = readFile("/tmp/osprey_code.txt")
        match readResult {
            Success { value } => value
            Error { message } => "print(\"JSON parsing failed!\")"
        }
    }
    Error { message } => "print(\"Failed to extract code from JSON!\")"
}
```

### Phase 2: Fix Compiler Execution
```osprey
// CURRENT CODE:
let programOutput = spawnProcess("../compiler/bin/osprey /tmp/osprey_web_temp.osp --run --quiet", processEventHandler)

// IMPROVEMENTS NEEDED:
1. Use absolute path or PATH lookup for osprey binary
2. Better error handling for process failures
3. Timeout protection (prevent processes from running indefinitely)
```

### Phase 3: Fix JSON Response Generation
```osprey
// NEED TO ADD BACK:
fn escapeJsonString(input: string) -> string = {
    // Replace " with \"
    // Replace \n with \\n  
    // Replace \r with \\r
    // Replace \t with \\t
    // Replace \ with \\
    input  // TODO: Implement proper escaping
}
```

## 🐛 KNOWN COMPILER BUGS & LIMITATIONS

### 1. String Manipulation Functions Missing
- **Problem**: No built-in string replace/escape functions
- **Impact**: Can't properly escape JSON strings
- **Workaround**: Basic escaping or rely on client-side handling
- **Fix Required**: Add string manipulation to standard library

### 2. Process Output Capture Issues
- **Problem**: Effects system not capturing all output properly
- **Impact**: Program output gets mixed with process exit messages
- **Status**: Partially fixed with --quiet flag
- **Remaining**: Need better separation of stdout/stderr in effects

### 3. File Path Resolution
- **Problem**: Relative paths break when running from different directories
- **Impact**: Server can't find osprey binary
- **Fix**: Use absolute paths or proper PATH resolution

### 4. Error Propagation
- **Problem**: Process failures not properly reported
- **Impact**: Client gets generic error instead of specific compiler error
- **Fix**: Better error message extraction from spawnProcess

## 🔧 REQUIRED COMPILER FEATURES

### Immediate Needs (for web API):
1. **String escaping functions**: `stringReplace`, `stringEscape`
2. **Better process error reporting**: Extract stderr from spawnProcess
3. **Path resolution**: `resolvePath`, `getAbsolutePath`
4. **Timeout handling**: `spawnProcessWithTimeout`

### Future Enhancements:
1. **JSON parsing/generation**: Native JSON support
2. **HTTP status codes**: Better HTTP response control
3. **Streaming responses**: For large output

## 🚀 IMPLEMENTATION STRATEGY

### Step 1: Manual JSON Parsing (IMMEDIATE)
```osprey
// Extract code from JSON request manually:
// 1. Write JSON to file
// 2. Use jq to extract .code field  
// 3. Read extracted code
// 4. Write to .osp file
// 5. Compile with --quiet flag
```

### Step 2: Fix Process Execution
```osprey
// Proper error handling:
match spawnProcess("osprey file.osp --run --quiet") {
    Success { value } => {
        // Process succeeded, get output
        let output = perform ProcessOutput.getAccumulatedOutput()
        "{\"success\": true, \"programOutput\": \"" + escapeJsonString(output) + "\"}"
    }
    Error { message } => {
        // Process failed, return error
        let output = perform ProcessOutput.getAccumulatedOutput() 
        "{\"success\": false, \"compilerOutput\": \"" + escapeJsonString(message) + "\", \"programOutput\": \"" + escapeJsonString(output) + "\"}"
    }
}
```

### Step 3: Production Hardening
1. **Security**: Validate input code for malicious content
2. **Resource limits**: CPU/memory/time limits per request
3. **Cleanup**: Remove temp files after execution
4. **Logging**: Proper error logging for debugging

## 🧪 TESTING STRATEGY

### Unit Tests:
- JSON parsing edge cases
- Error message formatting
- Output escaping

### Integration Tests:
- End-to-end API calls
- Compiler error scenarios
- Large output handling
- Timeout scenarios

### Performance Tests:
- Concurrent request handling
- Memory usage under load
- Response time benchmarks

## 📦 DEPLOYMENT CONSIDERATIONS

### Docker Requirements:
- Osprey compiler binary in PATH
- jq for JSON parsing
- Proper file permissions for /tmp
- Process cleanup on container shutdown

### Security:
- Sandbox mode for user code execution
- Input validation and sanitization
- Resource limits (CPU, memory, disk)
- Network isolation

## 🎯 SUCCESS CRITERIA

### API Must:
1. ✅ Accept POST requests with `{"code": "osprey source code"}`
2. ✅ Return `{"success": true/false, "compilerOutput": "...", "programOutput": "..."}`
3. ✅ Suppress compiler success messages (use --quiet)
4. ✅ Always show compiler error messages
5. ✅ Return clean program output without pollution
6. ✅ Handle malformed JSON gracefully
7. ✅ Clean up temp files
8. ✅ Respond within reasonable time limits

### Test Case:
```bash
curl -X POST http://localhost:3001/api/run \
  -H 'Content-Type: application/json' \
  -d '{"code":"print(\"Hello World!\")"}'

# Expected Response:
{"success": true, "compilerOutput": "", "programOutput": "Hello World!\n"}
```

## 🔥 PRIORITY ORDER

1. **HIGH**: Fix JSON parsing TODO (get actual code from request)
2. **HIGH**: Add back JSON escaping function
3. **MEDIUM**: Better error handling and reporting
4. **MEDIUM**: Production hardening (timeouts, cleanup)
5. **LOW**: Performance optimizations
6. **LOW**: Advanced features (streaming, etc.)

---

**Ready for implementation! 🚀**
