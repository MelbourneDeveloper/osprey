#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

// Process spawning function
int64_t spawn_process(char *command) {
  if (!command) {
    return -1;
  }

  // Execute the command and return the exit code
  int result = system(command);

  // Extract the actual exit code from system() result
  int exit_code = 0;
  if (WIFEXITED(result)) {
    exit_code = WEXITSTATUS(result);
  } else {
    exit_code = -1; // Command didn't exit normally
  }
  return (int64_t)exit_code;
}

// Write file function - returns 0 for success, negative for error
int64_t write_file(char *filename, char *content) {
  if (!filename || !content) {
    return -1;
  }

  FILE *file = fopen(filename, "w");
  if (!file) {
    return -2;
  }

  size_t written = fwrite(content, 1, strlen(content), file);
  fclose(file);

  return (int64_t)written;
}

// Read file function - returns content or NULL on error
char *read_file(char *filename) {
  if (!filename) {
    return NULL;
  }

  FILE *file = fopen(filename, "r");
  if (!file) {
    return NULL;
  }

  // Get file size
  fseek(file, 0, SEEK_END);
  long size = ftell(file);
  fseek(file, 0, SEEK_SET);

  // Allocate buffer and read content
  char *content = malloc(size + 1);
  if (!content) {
    fclose(file);
    return NULL;
  }

  size_t read_size = fread(content, 1, size, file);
  content[read_size] = '\0';
  fclose(file);

  return content;
}

// Simple JSON parsing - extract "code" field
char *parse_json(char *json_string) {
  if (!json_string) {
    return NULL;
  }

  // For now, just return the input
  // TODO: Implement proper JSON parsing
  return strdup(json_string);
}

// Extract arbitrary field from JSON {"field": "value"}
char *extract_json_field(char *json_string, char *field_name) {
  if (!json_string || !field_name) {
    return NULL;
  }

  // Create the search pattern: "field_name":
  char *pattern = malloc(strlen(field_name) + 4); // "field_name":
  sprintf(pattern, "\"%s\":", field_name);

  char *field_start = strstr(json_string, pattern);
  free(pattern);

  if (!field_start) {
    return NULL;
  }

  // Skip past "field_name":
  field_start += strlen(field_name) + 3;

  // Skip whitespace and quotes
  while (*field_start == ' ' || *field_start == '\t' || *field_start == '"') {
    field_start++;
  }

  // Find the end quote
  char *field_end = strchr(field_start, '"');
  if (!field_end) {
    return NULL;
  }

  // Extract the field value
  size_t field_len = field_end - field_start;
  char *extracted_value = malloc(field_len + 1);
  strncpy(extracted_value, field_start, field_len);
  extracted_value[field_len] = '\0';

  return extracted_value;
}

// Extract code from JSON {"code": "..."} - backward compatibility
char *extract_code(char *json_string) {
  return extract_json_field(json_string, "code");
}
