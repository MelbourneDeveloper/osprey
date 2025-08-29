// Package cli provides security configuration for the Osprey compiler.
package cli

import (
	"strings"
)

// SecurityConfig holds the security configuration for compilation.
type SecurityConfig struct {
	// Network access controls
	AllowHTTP      bool `json:"allowHttp"`
	AllowWebSocket bool `json:"allowWebSocket"`

	// File system access controls
	AllowFileRead  bool `json:"allowFileRead"`
	AllowFileWrite bool `json:"allowFileWrite"`

	// External function interface controls
	AllowFFI bool `json:"allowFfi"`

	// Process execution controls (for future use)
	AllowProcessExecution bool `json:"allowProcessExecution"`

	// Sandbox mode (disables all risky operations)
	SandboxMode bool `json:"sandboxMode"`
}

// NewDefaultSecurityConfig creates a security config with permissive defaults.
func NewDefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		AllowHTTP:             true,
		AllowWebSocket:        true,
		AllowFileRead:         true,
		AllowFileWrite:        true,
		AllowFFI:              true,
		AllowProcessExecution: true,
		SandboxMode:           false,
	}
}

// NewSandboxSecurityConfig creates a security config for sandbox mode (web compiler).
func NewSandboxSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		AllowHTTP:             false,
		AllowWebSocket:        false,
		AllowFileRead:         false,
		AllowFileWrite:        false,
		AllowFFI:              false,
		AllowProcessExecution: false,
		SandboxMode:           true,
	}
}

// ApplySandboxMode applies sandbox restrictions to the security config.
func (sc *SecurityConfig) ApplySandboxMode() {
	sc.SandboxMode = true
	sc.AllowHTTP = false
	sc.AllowWebSocket = false
	sc.AllowFileRead = false
	sc.AllowFileWrite = false
	sc.AllowFFI = false
	sc.AllowProcessExecution = false
}

// GetBlockedFunctions returns a list of function names that are not available in current security mode.
func (sc *SecurityConfig) GetBlockedFunctions() []string {
	var blocked []string

	if !sc.AllowHTTP {
		blocked = append(blocked,
			"httpCreateServer", "httpListen", "httpStopServer",
			"httpCreateClient", "httpGet", "httpPost", "httpPut",
			"httpDelete", "httpRequest", "httpCloseClient")
	}

	if !sc.AllowWebSocket {
		blocked = append(blocked,
			"websocketConnect", "websocketSend", "websocketClose",
			"websocketCreateServer", "websocketServerListen",
			"websocketServerSend", "websocketServerBroadcast",
			"websocketStopServer")
	}

	// File I/O functions would be added here when implemented
	if !sc.AllowFileRead {
		blocked = append(blocked, "readFile", "openFile")
	}

	if !sc.AllowFileWrite {
		blocked = append(blocked, "writeFile", "createFile", "deleteFile")
	}

	// FFI functions would be added here when extern validation is implemented
	if !sc.AllowFFI {
		blocked = append(blocked, "extern")
	}

	return blocked
}

// GetSecuritySummary returns a human-readable summary of the security configuration.
func (sc *SecurityConfig) GetSecuritySummary() string {
	if sc.SandboxMode {
		return "Security: SANDBOX MODE - Only safe functions available"
	}

	var allowed, blocked []string

	if sc.AllowHTTP {
		allowed = append(allowed, "HTTP")
	} else {
		blocked = append(blocked, "HTTP")
	}

	if sc.AllowWebSocket {
		allowed = append(allowed, "WebSocket")
	} else {
		blocked = append(blocked, "WebSocket")
	}

	if sc.AllowFileRead {
		allowed = append(allowed, "FileRead")
	} else {
		blocked = append(blocked, "FileRead")
	}

	if sc.AllowFileWrite {
		allowed = append(allowed, "FileWrite")
	} else {
		blocked = append(blocked, "FileWrite")
	}

	if sc.AllowFFI {
		allowed = append(allowed, "FFI")
	} else {
		blocked = append(blocked, "FFI")
	}

	summary := "Security:"
	if len(allowed) > 0 {
		summary += " Available=[" + strings.Join(allowed, ",") + "]"
	}

	if len(blocked) > 0 {
		summary += " Unavailable=[" + strings.Join(blocked, ",") + "]"
	}

	return summary
}
