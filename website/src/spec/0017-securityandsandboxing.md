---
layout: page
title: "Security and Sandboxing"
description: "Osprey Language Specification: Security and Sandboxing"
date: 2025-06-26
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0017-securityandsandboxing/"
---

# 17. Security and Sandboxing

- [Security Flags](#171-security-flags)
- [Security Policies](#172-security-policies)
- [Blocked Functions by Category](#173-blocked-functions-by-category)
- [Function Availability](#174-function-availability)
- [Programming Best Practices](#175-programming-best-practices)
- [Implementation Details](#176-implementation-details)

## 17. Security and Sandboxing

The Osprey compiler includes built-in security controls to restrict access to potentially dangerous functionality like network operations and file system access. This is essential for safe code execution in environments like web compilers where untrusted code may be executed.

## 17.1 Security Flags

#### `--sandbox`
Enables sandbox mode, which disables all potentially risky operations:
- HTTP/HTTPS operations (httpCreateServer, httpGet, httpPost, etc.)
- WebSocket operations (websocketConnect, websocketSend, etc.)
- File system access (when implemented)
- Foreign Function Interface (FFI)
- Process execution

**Example:**
```bash
osprey program.osp --sandbox --llvm
```

#### Granular Security Flags

For more granular control, you can disable specific categories of operations:

- `--no-http`: Disable HTTP client and server functions
- `--no-websocket`: Disable WebSocket client and server functions  
- `--no-fs`: Disable file system read/write operations
- `--no-ffi`: Disable foreign function interface

**Examples:**
```bash
# Disable only HTTP operations
osprey program.osp --no-http --compile

# Disable HTTP and WebSocket operations
osprey program.osp --no-http --no-websocket --run

# Disable file system access only
osprey program.osp --no-fs --llvm
```

## 17.2 Security Policies

#### Default Security (Permissive)
By default, all operations are allowed for backward compatibility and normal development use.

#### Sandbox Security (Restrictive)
When `--sandbox` is used, all potentially dangerous functions are unavailable. This is recommended for:
- Web-based code execution
- Untrusted code evaluation
- Educational environments
- Code review systems

## 17.3 Blocked Functions by Category

#### HTTP Functions
When HTTP access is disabled (`--no-http` or `--sandbox`), these functions are unavailable:
- `httpCreateServer` - Create HTTP server
- `httpListen` - Start HTTP server listening
- `httpStopServer` - Stop HTTP server
- `httpCreateClient` - Create HTTP client
- `httpGet` - HTTP GET request
- `httpPost` - HTTP POST request
- `httpPut` - HTTP PUT request
- `httpDelete` - HTTP DELETE request
- `httpRequest` - Generic HTTP request
- `httpCloseClient` - Close HTTP client

#### WebSocket Functions
When WebSocket access is disabled (`--no-websocket` or `--sandbox`), these functions are unavailable:
- `websocketConnect` - Connect to WebSocket server
- `websocketSend` - Send WebSocket message
- `websocketClose` - Close WebSocket connection
- `websocketCreateServer` - Create WebSocket server
- `websocketServerListen` - Start WebSocket server
- `websocketServerSend` - Send message to specific client
- `websocketServerBroadcast` - Broadcast message to all clients
- `websocketStopServer` - Stop WebSocket server

#### File System Functions (Future)
When file system access is disabled (`--no-fs` or `--sandbox`), these functions will be unavailable:
- `readFile` - Read file contents
- `writeFile` - Write file contents
- `deleteFile` - Delete file
- `createDirectory` - Create directory
- `listDirectory` - List directory contents

## 17.4 Function Availability

In different security modes, certain functions are simply not available in the language:

**Sandbox Mode**: Only safe functions like `print`, `toString`, `range`, etc. are available. Dangerous functions like `httpCreateServer` or `websocketConnect` result in "undefined function" compile errors.

**Partial Restrictions**: When specific categories are disabled (e.g., `--no-http`), those functions are unavailable while others remain accessible.

**Default Mode**: All functions are available.
- A human-readable explanation

## 17.5 Programming Best Practices

#### For Safe Code
Write code that doesn't use security-sensitive functions:
```osprey
// Safe operations - work in all security modes
let x = 42
let y = 24
let sum = x + y
print("Sum: ")
print(sum)
```

#### For Network Code
When writing network code, be aware that it may be restricted:
```osprey
// This will fail in sandbox mode or with --no-http
let serverID = httpCreateServer(port: 8080, address: "127.0.0.1")
```

## 17.6 Implementation Details

#### Security Configuration
Security settings are configured at compilation time and cannot be bypassed by the compiled program. The security checks happen during the LLVM IR generation phase, preventing security-sensitive functions from being included in the generated code.

#### Performance Impact
Security checks add minimal overhead during compilation and no runtime overhead, as restricted functions are simply not compiled into the final program.

#### Backward Compatibility
All existing code continues to work with default settings. Security restrictions are opt-in and don't affect normal development workflows.

#### Integration with Web Compiler
The security features are designed specifically for web compiler integration:

```javascript
// Example web compiler usage
const result = await compileOsprey(sourceCode, {
    mode: 'sandbox',  // Enable sandbox mode
    outputFormat: 'llvm'
});
```

#### Security Summary
When using security restrictions, the compiler will display a security summary:

```bash
# Sandbox mode
Security: SANDBOX MODE - All risky operations disabled

# Partial restrictions
Security: Allowed=[FileRead,FileWrite,FFI] Blocked=[HTTP,WebSocket]
```

---

## Summary

Osprey is a modern functional programming language with:

- **Type Safety**: No runtime panics, all errors handled explicitly via Result types
- **Named Arguments**: Multi-parameter functions require named arguments for clarity
- **Functional Programming**: Powerful iterator functions with pipe operator
- **Lightweight Fibers**: Zero-cost concurrency with Rust-like async/await
- **Fiber-Isolated Modules**: No global state, each fiber gets its own module instances
- **Rust Interoperability**: Seamless integration with Rust libraries
- **Memory Safety**: No shared mutable state between fibers

**Key Innovation**: The fiber-isolated module system eliminates data races by design while maintaining clean encapsulation through accessor patterns.