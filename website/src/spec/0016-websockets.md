---
layout: page
title: "WebSocket Functions"
description: "Osprey Language Specification: WebSocket Functions"
date: 2025-06-26
tags: ["specification", "reference", "documentation"]
author: "Christian Findlay"
permalink: "/spec/0016-websockets/"
---

## 16. WebSocket Functions

🔒 **IMPLEMENTATION STATUS**: WebSocket functions are implemented with security features but have critical design flaws that violate Osprey's functional programming principles.

### Table of Contents
- [16. WebSocket Functions](#16-websocket-functions)
  - [Table of Contents](#table-of-contents)
  - [16.1 WebSocket Core Types](#161-websocket-core-types)
  - [16.2 WebSocket Security Implementation](#162-websocket-security-implementation)
  - [16.3 WebSocket Client Functions](#163-websocket-client-functions)
    - [`websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`](#websocketconnecturl-string-messagehandler-fnstring---resultsuccess-string---resultwebsocketid-string)
    - [`websocketSend(wsID: Int, message: String) -> Result<Success, String>`](#websocketsendwsid-int-message-string---resultsuccess-string)
    - [`websocketClose(wsID: Int) -> Result<Success, String>`](#websocketclosewsid-int---resultsuccess-string)
  - [16.4 WebSocket Server Functions](#164-websocket-server-functions)
    - [`websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>`](#websocketcreateserverport-int-address-string-path-string---resultserverid-string)
    - [`websocketServerListen(serverID: Int) -> Result<Success, String>`](#websocketserverlistenserverid-int---resultsuccess-string)
    - [`websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>`](#websocketserverbroadcastserverid-int-message-string---resultsuccess-string)
    - [`websocketStopServer(serverID: Int) -> Result<Success, String>`](#websocketstopserverserverid-int---resultsuccess-string)
  - [16.5 Complete WebSocket Example](#165-complete-websocket-example)

WebSockets provide real-time, bidirectional communication between client and server. Osprey implements WebSocket support following functional programming principles with comprehensive error handling through Result types.

All WebSocket functions comply with:
- **RFC 6455**: The WebSocket Protocol ([https://tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455))
- **Result types** for all operations that can fail
- **Immutable message handling** with functional callbacks
- **Type safety** through structured error handling

### 16.1 WebSocket Core Types

```osprey
type WebSocketID = Int
type ServerID = Int

type WebSocketMessage = {
    type: String,
    data: String,
    timestamp: Int
}

type WebSocketConnection = {
    id: WebSocketID,
    url: String,
    isConnected: Bool
}
```

### 16.2 WebSocket Security Implementation

Osprey's WebSocket implementation follows **OWASP WebSocket Security Guidelines** with multiple security layers:

**🛡️ Cryptographic Security:**
- **OpenSSL SHA-1**: RFC 6455 compliant WebSocket handshake
- **Secure key validation**: 24-character base64 key format validation
- **Memory clearing**: All sensitive data zeroed before deallocation

**⚔️ Input Validation:**
- **WebSocket key format validation**: Strict RFC 6455 compliance
- **Buffer length validation**: Maximum limits prevent DoS attacks
- **Memory boundary checking**: No buffer overruns possible

### 16.3 WebSocket Client Functions

#### `websocketConnect(url: String, messageHandler: fn(String) -> Result<Success, String>) -> Result<WebSocketID, String>`

Establishes a WebSocket connection with a message handler callback.

**Parameters:**
- `url`: WebSocket URL (e.g., "ws://localhost:8080/chat")
- `messageHandler`: Callback function to handle incoming messages

**Returns:**
- `Success(wsID)`: WebSocket connection identifier
- `Err(message)`: Connection error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` instead of `Result<WebSocketID, String>` and takes string handler instead of function pointer

**Example:**
```osprey
fn handleMessage(message: String) -> Result<Success, String> = {
    print("Received: ${message}")
    Success()
}

let wsResult = websocketConnect(url: "ws://localhost:8080/chat", messageHandler: handleMessage)
match wsResult {
    Success wsId => {
        print("Connected with ID: ${wsId}")
        // Use the connection
    }
    Err message => print("Failed to connect: ${message}")
}
```

#### `websocketSend(wsID: Int, message: String) -> Result<Success, String>`

Sends a message through the WebSocket connection.

**Parameters:**
- `wsID`: WebSocket identifier from `websocketConnect`
- `message`: Message to send

**Returns:**
- `Success()`: Message sent successfully
- `Err(message)`: Send error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` instead of `Result<Success, String>`

**Example:**
```osprey
let sendResult = websocketSend(wsID: wsId, message: "Hello, WebSocket!")
match sendResult {
    Success _ => print("Message sent successfully")
    Err message => print("Failed to send: ${message}")
}
```

#### `websocketClose(wsID: Int) -> Result<Success, String>`

Closes the WebSocket connection and cleans up resources.

**Parameters:**
- `wsID`: WebSocket identifier to close

**Returns:**
- `Success()`: Connection closed successfully
- `Err(message)`: Close error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` instead of `Result<Success, String>`

**Example:**
```osprey
let closeResult = websocketClose(wsID: wsId)
match closeResult {
    Success _ => print("Connection closed")
    Err message => print("Failed to close: ${message}")
}
```

### 16.4 WebSocket Server Functions

#### `websocketCreateServer(port: Int, address: String, path: String) -> Result<ServerID, String>`

Creates a WebSocket server bound to the specified port, address, and path.

**Parameters:**
- `port`: Port number (1-65535)
- `address`: IP address to bind to (e.g., "127.0.0.1", "0.0.0.0")
- `path`: WebSocket endpoint path (e.g., "/chat", "/live")

**Returns:**
- `Success(serverID)`: Unique WebSocket server identifier
- `Err(message)`: Server creation error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` instead of `Result<ServerID, String>`. Also has critical runtime issues with port binding failures.

**Example:**
```osprey
let serverResult = websocketCreateServer(port: 8080, address: "127.0.0.1", path: "/chat")
match serverResult {
    Success serverId => print("WebSocket server created with ID: ${serverId}")
    Err message => print("Failed to create server: ${message}")
}
```

#### `websocketServerListen(serverID: Int) -> Result<Success, String>`

Starts the WebSocket server listening for connections.

**Parameters:**
- `serverID`: Server identifier from `websocketCreateServer`

**Returns:**
- `Success()`: Server started listening successfully
- `Err(message)`: Listen error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` instead of `Result<Success, String>`. Currently returns `-4` (bind failed) due to port binding issues.

**Example:**
```osprey
let listenResult = websocketServerListen(serverID: serverId)
match listenResult {
    Success _ => print("Server listening on ws://127.0.0.1:8080/chat")
    Err message => print("Failed to start listening: ${message}")
}
```

#### `websocketServerBroadcast(serverID: Int, message: String) -> Result<Success, String>`

Broadcasts a message to all connected WebSocket clients.

**Parameters:**
- `serverID`: Server identifier
- `message`: Message to broadcast to all clients

**Returns:**
- `Success()`: Message broadcasted successfully
- `Err(message)`: Broadcast error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` (number of clients sent to) instead of `Result<Success, String>`

**Example:**
```osprey
let broadcastResult = websocketServerBroadcast(serverID: serverId, message: "Welcome to Osprey Chat!")
match broadcastResult {
    Success _ => print("Message broadcasted to all clients")
    Err message => print("Failed to broadcast: ${message}")
}
```

#### `websocketStopServer(serverID: Int) -> Result<Success, String>`

Stops the WebSocket server and closes all connections.

**Parameters:**
- `serverID`: Server identifier to stop

**Returns:**
- `Success()`: Server stopped successfully
- `Err(message)`: Stop error description

**Implementation Status:** ⚠️ **INCORRECT** - Current C runtime returns raw `int64_t` instead of `Result<Success, String>`

**Example:**
```osprey
let stopResult = websocketStopServer(serverID: serverId)
match stopResult {
    Success _ => print("Server stopped successfully")
    Err message => print("Failed to stop server: ${message}")
}
```

### 16.5 Complete WebSocket Example

A practical WebSocket server and client implementation demonstrating real-time communication:

```osprey
fn main() -> Int = {
    print("Starting WebSocket chat server...")
    
    // Create WebSocket server
    let serverResult = websocketCreateServer(port: 8080, address: "127.0.0.1", path: "/chat")
    
    match serverResult {
        Success serverId => {
            print("WebSocket server created with ID: ${serverId}")
            
            // Start listening for connections
            let listenResult = websocketServerListen(serverID: serverId)
            
            match listenResult {
                Success _ => {
                    print("Server listening on ws://127.0.0.1:8080/chat")
                    
                    // Broadcast welcome message
                    let welcomeResult = websocketServerBroadcast(
                        serverID: serverId, 
                        message: "Welcome to Osprey Chat!"
                    )
                    
                    match welcomeResult {
                        Success _ => print("Welcome message broadcasted")
                        Err message => print("Failed to broadcast welcome: ${message}")
                    }
                    
                    // Keep server alive for 10 seconds
                    sleep(10000)
                    
                    // Clean shutdown
                    let stopResult = websocketStopServer(serverID: serverId)
                    match stopResult {
                        Success _ => print("Server stopped successfully")
                        Err message => print("Failed to stop server: ${message}")
                    }
                }
                Err message => print("Failed to start listening: ${message}")
            }
        }
        Err message => print("Failed to create server: ${message}")
    }
    
    0
}

// WebSocket client example
fn connectToChat() -> Int = {
    fn chatMessageHandler(message: String) -> Result<Success, String> = {
        print("Chat: ${message}")
        
        // Process different message types
        match message {
            "ping" => {
                print("Received ping, responding with pong")
                Success()
            }
            _ => {
                print("Received message: ${message}")
                Success()
            }
        }
    }
    
    let wsResult = websocketConnect(
        url: "ws://127.0.0.1:8080/chat", 
        messageHandler: chatMessageHandler
    )
    
    match wsResult {
        Success wsId => {
            print("Connected to chat server with ID: ${wsId}")
            
            // Send some messages
            let sendResult1 = websocketSend(wsID: wsId, message: "Hello from Osprey!")
            match sendResult1 {
                Success _ => print("First message sent")
                Err message => print("Failed to send first message: ${message}")
            }
            
            let sendResult2 = websocketSend(wsID: wsId, message: "How is everyone?")
            match sendResult2 {
                Success _ => print("Second message sent")
                Err message => print("Failed to send second message: ${message}")
            }
            
            let pingResult = websocketSend(wsID: wsId, message: "ping")
            match pingResult {
                Success _ => print("Ping sent")
                Err message => print("Failed to send ping: ${message}")
            }
            
            // Wait a bit before closing
            sleep(2000)
            
            // Close connection
            let closeResult = websocketClose(wsID: wsId)
            match closeResult {
                Success _ => print("Disconnected from chat")
                Err message => print("Failed to disconnect: ${message}")
            }
        }
        Err message => print("Failed to connect: ${message}")
    }
    
    0
}

// Chat message processing with pattern matching
fn processMessage(message: String) -> Result<String, String> = 
    match message {
        "{\"type\":\"join\"}" => Success("User joined the chat")
        "{\"type\":\"leave\"}" => Success("User left the chat")
        "{\"type\":\"message\"}" => Success("Regular chat message")
        _ => Success("Unknown message type")
    }
```

This example demonstrates:

- **Result type error handling** for all WebSocket operations
- **Functional message handlers** with proper callback signatures
- **Pattern matching** for different message types
- **Resource management** with proper connection cleanup
- **Real-time communication** patterns for chat applications
- **Comprehensive error handling** following Osprey's functional principles
- **Type-safe WebSocket operations** with structured error messages

**Key Differences from Current Implementation:**
1. **All functions return Result types** instead of raw integers
2. **Message handlers are proper function pointers** instead of string identifiers
3. **Comprehensive error handling** with descriptive error messages
4. **Functional programming patterns** with immutable message processing
5. **Type safety** through structured return types and pattern matching