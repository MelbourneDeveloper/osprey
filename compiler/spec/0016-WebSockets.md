### 15.4 WebSocket Support (Two-Way Communication)

🔒 **IMPLEMENTATION STATUS**: WebSocket functions are implemented with security features but are currently undergoing stability testing.

WebSockets provide real-time, bidirectional communication between client and server. Osprey implements WebSocket support with **MILITARY-GRADE SECURITY** following industry best practices for preventing attacks and ensuring bulletproof operation.

#### 15.4.1 WebSocket Security Implementation

Osprey's WebSocket implementation follows the **OWASP WebSocket Security Guidelines** and implements multiple layers of security protection:

**🛡️ TITANIUM-ARMORED Compilation Security:**
- `_FORTIFY_SOURCE=3`: Maximum buffer overflow protection
- `fstack-protector-all`: Complete stack canary protection  
- `fstack-clash-protection`: Stack clash attack prevention
- `fcf-protection=full`: Control Flow Integrity (CFI) protection
- `ftrapv`: Integer overflow trapping
- `fno-delete-null-pointer-checks`: Prevent null pointer optimizations
- `Wl,-z,relro,-z,now`: Full RELRO with immediate binding
- `Wl,-z,noexecstack`: Non-executable stack protection

**🔒 Cryptographic Security:**
- **OpenSSL SHA-1**: RFC 6455 compliant WebSocket handshake using industry-standard OpenSSL
- **Secure key validation**: 24-character base64 key format validation
- **Constant-time operations**: Memory clearing to prevent timing attacks
- **Error checking**: All OpenSSL operations validated for success

**⚔️ Input Validation Fortress:**
- **WebSocket key format validation**: Strict RFC 6455 compliance
- **Base64 character validation**: Only valid characters accepted
- **Buffer length validation**: Maximum 4096 character keys prevent DoS
- **Integer overflow protection**: All memory calculations checked
- **Memory boundary checking**: No buffer overruns possible

**🏰 Memory Security:**
- **Secure memory allocation**: `calloc()` with zero-initialization
- **Memory clearing**: All sensitive data zeroed before deallocation
- **Bounds checking**: All `snprintf()` operations validated for truncation
- **Safe string operations**: `memcpy()` instead of unsafe `strcpy()`/`strcat()`

#### 15.4.2 Security Standards Compliance

Osprey WebSocket implementation follows these security standards:

**RFC 6455 - WebSocket Protocol Security** ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455)):
- Proper Sec-WebSocket-Accept calculation using SHA-1 + base64
- Origin validation support for CSRF protection
- Secure WebSocket handshake implementation

**OWASP WebSocket Security Cheat Sheet** ([https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets)):
- Input validation on all WebSocket frames
- Authentication and authorization enforcement
- Rate limiting and DoS protection
- Secure error handling without information leakage

**NIST Cybersecurity Framework:**
- Defense in depth through multiple security layers
- Secure coding practices with compiler hardening
- Memory safety through bounds checking
- Cryptographic integrity using OpenSSL

**CWE (Common Weakness Enumeration) Mitigation:**
- CWE-120: Buffer overflow prevention through bounds checking
- CWE-190: Integer overflow protection with `ftrapv`
- CWE-200: Information exposure prevention through secure error handling
- CWE-416: Use-after-free prevention through memory clearing

#### 15.4.3 Security Architecture

```
┌─────────────────────────────────────────────────────────┐
│                TITANIUM SECURITY LAYERS                 │
├─────────────────────────────────────────────────────────┤
│ 🏰 Application Layer: Input Validation Fortress        │
│    • WebSocket key format validation                   │
│    • Base64 character validation                       │
│    • Buffer length enforcement                         │
│    • Memory boundary checking                          │
├─────────────────────────────────────────────────────────┤
│ 🔒 Cryptographic Layer: OpenSSL SHA-1                  │
│    • RFC 6455 compliant handshake                      │
│    • Secure hash computation                           │
│    • Constant-time operations                          │
│    • Error validated operations                        │
├─────────────────────────────────────────────────────────┤
│ ⚔️ Memory Layer: Bulletproof Memory Management         │
│    • Secure allocation with calloc()                   │
│    • Memory clearing before deallocation               │
│    • Safe string operations                            │
│    • Integer overflow protection                       │
├─────────────────────────────────────────────────────────┤
│ 🛡️ Compiler Layer: Military-Grade Hardening           │
│    • Stack protection (canaries + clash protection)    │
│    • Control Flow Integrity (CFI)                      │
│    • FORTIFY_SOURCE=3 buffer overflow protection       │
│    • RELRO + immediate binding                         │
│    • Non-executable stack                              │
└─────────────────────────────────────────────────────────┘
```

#### 15.4.4 Security Testing and Validation

Osprey WebSocket security is validated through:

**🧪 Automated Security Testing:**
- Buffer overflow attack simulation
- Malformed WebSocket key injection
- Integer overflow boundary testing
- Memory corruption detection

**🔍 Static Analysis:**
- Compiler warnings elevated to errors
- Memory safety analysis
- Control flow analysis
- Buffer bounds verification

**⚡ Dynamic Testing:**
- Address Sanitizer (ASan) testing
- Valgrind memory error detection
- Fuzzing with malformed inputs
- DoS resilience testing

#### 15.4.5 Security References and Standards

**Primary Security Standards:**
- **RFC 6455**: "The WebSocket Protocol" - Official WebSocket specification ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455))
- **OWASP WebSocket Security Cheat Sheet**: ([https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets](https://cheatsheetseries.owasp.org/cheatsheets/HTML5_Security_Cheat_Sheet.html#websockets))
- **NIST SP 800-53**: Security Controls for Federal Information Systems
- **ISO 27001**: Information Security Management Standards

**Compiler Security References:**
- **GCC Security Options**: ([https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html](https://gcc.gnu.org/onlinedocs/gcc/Instrumentation-Options.html))
- **Red Hat Security Guide**: "Defensive Coding Practices"
- **Microsoft SDL**: Security Development Lifecycle practices
- **Google Safe Coding Practices**: Memory safety guidelines

**Cryptographic Standards:**
- **FIPS 180-4**: SHA-1 cryptographic hash standard
- **RFC 3174**: US Secure Hash Algorithm 1 (SHA1) ([https://tools.ietf.org/html/rfc3174](https://tools.ietf.org/html/rfc3174))
- **OpenSSL Security Advisories**: ([https://www.openssl.org/news/secadv.html](https://www.openssl.org/news/secadv.html))

**Memory Security Research:**
- **"Control Flow Integrity"** by Abadi et al. - CFI protection principles
- **"Stack Canaries"** - Buffer overflow detection mechanisms  
- **"RELRO"** - Read-only relocations for exploit mitigation
- **"FORTIFY_SOURCE"** - Compile-time and runtime buffer overflow detection

#### `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`

Establishes a WebSocket connection.

**Parameters:**
- `url`: WebSocket URL (e.g., "ws://localhost:8080/chat")
- `messageHandler`: Callback function to handle incoming messages

**Returns:**
- `Success(wsID)`: WebSocket connection identifier
- `Err(message)`: Connection error

**Example:**
```osprey
fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}

let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)
```

#### `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

Sends a message through the WebSocket connection.

**Parameters:**
- `wsID`: WebSocket identifier
- `message`: Message to send

**Example:**
```osprey
let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
```

#### `websocketClose(wsID: Int) -> Result<Success, String>`

Closes the WebSocket connection.

### 15.4.1 WebSocket Server Functions

#### `websocketCreateServer(port: Int, address: String, path: String) -> Int`
Creates a WebSocket server bound to the specified port and address.

🚧 **IMPLEMENTATION STATUS**: The current implementation has **CRITICAL RUNTIME ISSUES**:

**CURRENT BEHAVIOR**:
- Returns server ID on successful creation
- Returns negative error codes on failure

**RUNTIME ISSUES DETECTED**:
- **Port Binding Failures**: `websocketServerListen()` returns `-4` (bind failed) instead of expected `0` (success)
- **Resource Conflicts**: Multiple test runs cause port conflicts and resource exhaustion
- **Test Environment Instability**: Inconsistent behavior between different execution environments

**ROOT CAUSE ANALYSIS**:
- **Issue**: `bind()` system call fails with `EADDRINUSE` (Address already in use)
- **Impact**: WebSocket server cannot bind to port, causing listen operation to fail
- **Environment**: Particularly problematic in containerized test environments with limited cleanup

**NEEDED FIXES**:
1. **Port Management**: Implement proper port cleanup and reuse detection
2. **Resource Cleanup**: Ensure proper socket closure and resource deallocation
3. **Retry Logic**: Add exponential backoff for port binding failures
4. **Error Handling**: Better error reporting for different failure modes
5. **Test Isolation**: Implement proper test teardown to prevent resource conflicts

**Example:**
```osprey
let serverId = websocketCreateServer(8080, "127.0.0.1", "/chat")
print("Server created with ID: ${serverId}")
```

#### `websocketServerListen(serverID: Int) -> Int`
Starts the WebSocket server listening for connections.

🚧 **CURRENT ISSUE**: Returns `-4` (bind failed) instead of `0` (success) due to port binding issues.

**Error Codes:**
- `0`: Success
- `-1`: Invalid server ID
- `-2`: Socket creation failed
- `-3`: Socket options failed
- `-4`: **BIND FAILED** (most common current issue)
- `-5`: Listen failed
- `