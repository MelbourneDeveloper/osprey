package codegen

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

// WebSocket bridge functions exported to C runtime
// These functions are called by the C WebSocket runtime when events occur

//export osprey_handle_websocket_connection
func osprey_handle_websocket_connection(serverID C.int, connectionID C.int, clientIP *C.char) C.int {
	// Convert C strings to Go strings
	clientIPStr := C.GoString(clientIP)

	fmt.Printf("🔌 WebSocket connection established - Server: %d, Connection: %d, Client: %s\n",
		int(serverID), int(connectionID), clientIPStr)

	// TODO: Call into Osprey user code to handle connection event
	// For now, just return success
	return 0
}

//export osprey_handle_websocket_message
func osprey_handle_websocket_message(
	serverID C.int,
	connectionID C.int,
	message *C.char,
	messageLength C.size_t,
) (*C.char, C.size_t) {
	// Convert C string to Go string
	messageData := C.GoStringN(message, C.int(messageLength))

	fmt.Printf("💬 WebSocket message received - Server: %d, Connection: %d, Message: %s\n",
		int(serverID), int(connectionID), messageData)

	// TODO: Call into Osprey user code to handle message
	// For now, echo the message back
	response := "Echo: " + messageData

	// Allocate C string for response
	cResponse := C.CString(response)
	responseLength := C.size_t(len(response))

	return cResponse, responseLength
}

//export osprey_handle_websocket_disconnect
func osprey_handle_websocket_disconnect(serverID C.int, connectionID C.int) C.int {
	fmt.Printf("🔌 WebSocket connection closed - Server: %d, Connection: %d\n",
		int(serverID), int(connectionID))

	// TODO: Call into Osprey user code to handle disconnect event
	// For now, just return success
	return 0
}

//export osprey_free_websocket_response
func osprey_free_websocket_response(response *C.char) {
	// Free the allocated response string
	if response != nil {
		C.free(unsafe.Pointer(response))
	}
}
