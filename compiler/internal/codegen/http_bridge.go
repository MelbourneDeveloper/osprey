package codegen

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"unsafe"
)

// HTTP response status codes
const (
	HTTPStatusOK                  = 200
	HTTPStatusCreated             = 201
	HTTPStatusNotFound            = 404
	HTTPStatusMethodNotAllowed    = 405
	HTTPStatusInternalServerError = 500
)

// HTTP header constants
const (
	ContentTypeJSON      = "Content-Type: application/json\r\n"
	ContentTypePlainText = "Content-Type: text/plain\r\n"
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
			*responseStatus = HTTPStatusInternalServerError
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
// TODO: This should call into an Osprey function..
func handleRawHTTPRequest(_ int, method, fullURL, _ string, body []byte) (int, string, []byte) {
	// NOTE: WE CANNOT HARD CODE FUNCTIONALITY FOR SPECIFIC SERVICE HERE
	// WE JUST NEED TO CALL THE FUNCTION AND PASS THE REQUEST DATA ALONG

	// Temporary implementation until proper Osprey function calling is implemented
	// This should be replaced with a call to the user-defined Osprey HTTP handler function
	_ = method
	_ = fullURL
	_ = body

	// Return a basic "not implemented" response
	responseBody := []byte(`{"error": "HTTP handler not implemented"}`)
	return HTTPStatusInternalServerError, ContentTypeJSON, responseBody
}
