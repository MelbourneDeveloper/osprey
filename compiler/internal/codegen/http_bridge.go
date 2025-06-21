package codegen

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

//export osprey_handle_http_request
func osprey_handle_http_request(
	serverID C.int,
	method *C.char,
	fullURL *C.char,
	rawHeaders *C.char,
	body *C.char,
	bodyLength C.size_t,
	responseStatus *C.int,
	responseHeaders **C.char,
	responseBody **C.char,
	responseBodyLength *C.size_t,
) C.int {
	// SAFETY: Check for null pointers first
	if method == nil || fullURL == nil || rawHeaders == nil {
		if responseStatus != nil {
			*responseStatus = 500
		}
		return -1
	}

	// Convert C strings to Go strings
	goMethod := C.GoString(method)
	goURL := C.GoString(fullURL)
	goHeaders := C.GoString(rawHeaders)

	// Handle body data (may be binary)
	var goBody []byte
	if body != nil && bodyLength > 0 {
		goBody = C.GoBytes(unsafe.Pointer(body), C.int(bodyLength))
	}

	// Handle the request in Osprey - NO ROUTING IN GO!
	// This should call actual Osprey pattern matching code
	status, headers, respBody := handleRawHTTPRequest(int(serverID), goMethod, goURL, goHeaders, goBody)

	// Set response status
	*responseStatus = C.int(status)

	// Allocate and set response headers
	if headers != "" {
		*responseHeaders = C.CString(headers)
	} else {
		*responseHeaders = nil
	}

	// Allocate and set response body
	if len(respBody) > 0 {
		// Allocate C memory for response body
		*responseBody = (*C.char)(C.malloc(C.size_t(len(respBody))))
		C.memcpy(unsafe.Pointer(*responseBody), unsafe.Pointer(&respBody[0]), C.size_t(len(respBody)))
		*responseBodyLength = C.size_t(len(respBody))
	} else {
		*responseBody = nil
		*responseBodyLength = 0
	}

	return 0 // Success
}

// handleRawHTTPRequest handles the raw HTTP request data
// TODO: This should call into actual Osprey pattern matching code
func handleRawHTTPRequest(serverID int, method, fullURL, rawHeaders string, body []byte) (int, string, []byte) {
	// Extract path from full URL (remove query parameters for routing)
	path := fullURL
	if queryIndex := strings.Index(fullURL, "?"); queryIndex != -1 {
		path = fullURL[:queryIndex]
	}

	// Simple routing - this should be replaced with Osprey pattern matching
	switch method {
	case "GET":
		switch path {
		case "/api/users":
			headers := "Content-Type: application/json\r\n"
			body := []byte(`[{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}]`)
			return 200, headers, body
		case "/api/health":
			headers := "Content-Type: application/json\r\n"
			body := []byte(`{"status": "healthy"}`)
			return 200, headers, body
		default:
			headers := "Content-Type: text/plain\r\n"
			body := []byte("Not Found")
			return 404, headers, body
		}
	case "POST":
		switch path {
		case "/api/users":
			headers := "Content-Type: application/json\r\n"
			responseBody := fmt.Sprintf(`{"id": 3, "name": "New User", "message": "User created", "received_body_length": %d}`, len(body))
			return 201, headers, []byte(responseBody)
		default:
			headers := "Content-Type: text/plain\r\n"
			body := []byte("Endpoint not found")
			return 404, headers, body
		}
	default:
		headers := "Content-Type: text/plain\r\n"
		body := []byte("Method not allowed")
		return 405, headers, body
	}
}
