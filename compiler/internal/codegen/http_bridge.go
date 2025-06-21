package codegen

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"unsafe"
)

// HTTP status codes
const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusNotFound            = 404
	StatusMethodNotAllowed    = 405
	StatusInternalServerError = 500
)

// Content type headers
const (
	ContentTypeJSON = "Content-Type: application/json\r\n"
	ContentTypeText = "Content-Type: text/plain\r\n"
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
			*responseStatus = StatusInternalServerError
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
func handleRawHTTPRequest(_ int, _ string, _ string, _ string, _ []byte) (int, string, []byte) {
	// TODO: implement this properly!!!

	// THIS MUST NOT HARDCODE API ENDPOINTS!!!

	// THIS MUST JUST FORWARD THE REQUEST TO THE OSPREY CODE VIA A CALLBACK (FUNCTION COMPOSITION)

	// Temporary return to fix compilation error
	// This will be replaced with proper callback forwarding
	return StatusNotFound, "", []byte("Not implemented yet")
}
