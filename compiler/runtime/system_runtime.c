#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/wait.h>

// Forward declarations for Osprey functions
extern char* processCompileRequest(char* jsonBody);
extern char* processRunRequest(char* jsonBody);

// Process spawning function
int64_t spawn_process(char* command) {
    if (!command) {
        return -1;
    }
    
    printf("🚀 Spawning process: %s\n", command);
    
    // Execute the command and return the exit code
    int result = system(command);
    
    printf("✅ Process completed with result: %d\n", result);
    return (int64_t)result;
}

// Write file function
int64_t write_file(char* filename, char* content) {
    if (!filename || !content) {
        return -1;
    }
    
    printf("📝 Writing file: %s\n", filename);
    
    FILE* file = fopen(filename, "w");
    if (!file) {
        printf("❌ Failed to open file for writing: %s\n", filename);
        return -2;
    }
    
    size_t written = fwrite(content, 1, strlen(content), file);
    fclose(file);
    
    printf("✅ Wrote %zu bytes to %s\n", written, filename);
    return (int64_t)written;
}

// Read file function
char* read_file(char* filename) {
    if (!filename) {
        return strdup("");
    }
    
    printf("📖 Reading file: %s\n", filename);
    
    FILE* file = fopen(filename, "r");
    if (!file) {
        printf("❌ Failed to open file for reading: %s\n", filename);
        return strdup("");
    }
    
    // Get file size
    fseek(file, 0, SEEK_END);
    long size = ftell(file);
    fseek(file, 0, SEEK_SET);
    
    // Allocate buffer and read content
    char* content = malloc(size + 1);
    if (!content) {
        fclose(file);
        return strdup("");
    }
    
    size_t read_size = fread(content, 1, size, file);
    content[read_size] = '\0';
    fclose(file);
    
    printf("✅ Read %zu bytes from %s\n", read_size, filename);
    return content;
}

// Simple JSON parsing - extract "code" field
char* parse_json(char* json_string) {
    if (!json_string) {
        return strdup("");
    }
    
    printf("🔍 Parsing JSON: %s\n", json_string);
    
    // For now, just return the input
    // TODO: Implement proper JSON parsing
    return strdup(json_string);
}

// Extract code from JSON {"code": "..."}
char* extract_code(char* json_string) {
    if (!json_string) {
        return strdup("");
    }
    
    printf("🎯 Extracting code from JSON: %s\n", json_string);
    
    // Simple extraction - look for "code":"..." pattern
    char* code_start = strstr(json_string, "\"code\":");
    if (!code_start) {
        printf("❌ No 'code' field found in JSON\n");
        return strdup("print(\"No code found\")");
    }
    
    // Skip past "code":
    code_start += 7;
    
    // Skip whitespace and quotes
    while (*code_start == ' ' || *code_start == '\t' || *code_start == '"') {
        code_start++;
    }
    
    // Find the end quote
    char* code_end = strchr(code_start, '"');
    if (!code_end) {
        printf("❌ Malformed JSON - no closing quote\n");
        return strdup("print(\"Malformed JSON\")");
    }
    
    // Extract the code
    size_t code_len = code_end - code_start;
    char* extracted_code = malloc(code_len + 1);
    strncpy(extracted_code, code_start, code_len);
    extracted_code[code_len] = '\0';
    
    // Handle escape sequences (simple unescape)
    char* unescaped = malloc(strlen(extracted_code) + 1);
    char* src = extracted_code;
    char* dst = unescaped;
    
    while (*src) {
        if (*src == '\\' && *(src + 1)) {
            src++; // Skip backslash
            switch (*src) {
                case 'n': *dst++ = '\n'; break;
                case 't': *dst++ = '\t'; break;
                case 'r': *dst++ = '\r'; break;
                case '"': *dst++ = '"'; break;
                case '\\': *dst++ = '\\'; break;
                default: *dst++ = *src; break;
            }
            src++;
        } else {
            *dst++ = *src++;
        }
    }
    *dst = '\0';
    
    free(extracted_code);
    
    printf("✅ Extracted code: %s\n", unescaped);
    return unescaped;
}

// Bridge functions for HTTP server to call Osprey functions
char* process_compile_request(char* json_body) {
    if (!json_body) {
        return strdup("{\"error\": \"No request body provided\"}");
    }
    
    printf("🔗 Bridge: calling processCompileRequest with: %s\n", json_body);
    char* result = processCompileRequest(json_body);
    if (!result) {
        return strdup("{\"error\": \"Compilation function returned null\"}");
    }
    
    return strdup(result);
}

char* process_run_request(char* json_body) {
    if (!json_body) {
        return strdup("{\"error\": \"No request body provided\"}");
    }
    
    printf("🔗 Bridge: calling processRunRequest with: %s\n", json_body);
    char* result = processRunRequest(json_body);
    if (!result) {
        return strdup("{\"error\": \"Run function returned null\"}");
    }
    
    return strdup(result);
} 