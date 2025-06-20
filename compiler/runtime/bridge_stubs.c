#include <stdlib.h>
#include <string.h>

// Stub implementations for bridge functions expected by HTTP server runtime
// These would normally be implemented by the Go compiler when linking

char *process_compile_request(char *json_body) {
  // Stub implementation - return a simple JSON response
  char *response = malloc(256);
  if (response) {
    strcpy(response,
           "{\"status\":\"error\",\"message\":\"Compiler not linked\"}");
  }
  return response;
}

char *process_run_request(char *json_body) {
  // Stub implementation - return a simple JSON response
  char *response = malloc(256);
  if (response) {
    strcpy(response,
           "{\"status\":\"error\",\"message\":\"Runtime not linked\"}");
  }
  return response;
}