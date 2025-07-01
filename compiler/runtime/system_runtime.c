#include <errno.h>
#include <fcntl.h>
#include <pthread.h>
#include <signal.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

// Process event handler function type - Osprey provides this callback
typedef void (*ProcessEventHandler)(int64_t process_id, int64_t event_type,
                                    char *data);

// Event types for process callbacks
#define PROCESS_STDOUT_DATA 1
#define PROCESS_STDERR_DATA 2
#define PROCESS_EXIT 3

// Process result structure
typedef struct {
  int64_t process_id;          // Process ID for tracking
  char *command;               // Command being executed
  int64_t exit_code;           // Process exit code
  bool is_running;             // Process status
  pthread_t monitor_thread;    // Thread monitoring the process
  pthread_mutex_t mutex;       // Mutex for thread safety
  int stdout_pipe[2];          // Pipes for capturing stdout
  int stderr_pipe[2];          // Pipes for capturing stderr
  pid_t pid;                   // Actual process PID
  ProcessEventHandler handler; // Callback for events
} ProcessResult;

// Global process tracking
#define MAX_PROCESSES 1000
static ProcessResult *processes[MAX_PROCESSES];
static int64_t next_process_id = 1;
static pthread_mutex_t process_mutex = PTHREAD_MUTEX_INITIALIZER;

// Thread function to monitor process and send callbacks
static void *process_monitor_thread(void *arg) {
  ProcessResult *proc = (ProcessResult *)arg;

  // Close write ends in parent
  close(proc->stdout_pipe[1]);
  close(proc->stderr_pipe[1]);

  // Make pipes non-blocking
  fcntl(proc->stdout_pipe[0], F_SETFL, O_NONBLOCK);
  fcntl(proc->stderr_pipe[0], F_SETFL, O_NONBLOCK);

  char buffer[1024];
  fd_set read_fds;
  struct timeval timeout;

  // Monitor process and send callbacks for output
  while (proc->is_running) {
    FD_ZERO(&read_fds);
    FD_SET(proc->stdout_pipe[0], &read_fds);
    FD_SET(proc->stderr_pipe[0], &read_fds);

    timeout.tv_sec = 0;
    timeout.tv_usec = 100000; // 100ms timeout

    int max_fd = (proc->stdout_pipe[0] > proc->stderr_pipe[0])
                     ? proc->stdout_pipe[0]
                     : proc->stderr_pipe[0];

    int ready = select(max_fd + 1, &read_fds, NULL, NULL, &timeout);

    if (ready > 0) {
      // Read stdout and send callback
      if (FD_ISSET(proc->stdout_pipe[0], &read_fds)) {
        ssize_t bytes = read(proc->stdout_pipe[0], buffer, sizeof(buffer) - 1);
        if (bytes > 0) {
          buffer[bytes] = '\0';

          // Send stdout data to Osprey via callback
          if (proc->handler) {
            proc->handler(proc->process_id, PROCESS_STDOUT_DATA, buffer);
          }
        }
      }

      // Read stderr and send callback
      if (FD_ISSET(proc->stderr_pipe[0], &read_fds)) {
        ssize_t bytes = read(proc->stderr_pipe[0], buffer, sizeof(buffer) - 1);
        if (bytes > 0) {
          buffer[bytes] = '\0';

          // Send stderr data to Osprey via callback
          if (proc->handler) {
            proc->handler(proc->process_id, PROCESS_STDERR_DATA, buffer);
          }
        }
      }
    }

    // Check if process is still running
    int status;
    pid_t result = waitpid(proc->pid, &status, WNOHANG);
    if (result > 0) {
      // Process finished
      pthread_mutex_lock(&proc->mutex);
      proc->is_running = false;
      if (WIFEXITED(status)) {
        proc->exit_code = WEXITSTATUS(status);
      } else if (WIFSIGNALED(status)) {
        proc->exit_code = -1; // Terminated by signal
      }
      pthread_mutex_unlock(&proc->mutex);

      // Send exit event to Osprey
      if (proc->handler) {
        char exit_code_str[32];
        snprintf(exit_code_str, sizeof(exit_code_str), "%lld",
                 (long long)proc->exit_code);
        proc->handler(proc->process_id, PROCESS_EXIT, exit_code_str);
      }
      break;
    } else if (result < 0 && errno != ECHILD) {
      // Error in waitpid
      pthread_mutex_lock(&proc->mutex);
      proc->is_running = false;
      proc->exit_code = -1;
      pthread_mutex_unlock(&proc->mutex);

      // Send error exit event
      if (proc->handler) {
        proc->handler(proc->process_id, PROCESS_EXIT, "-1");
      }
      break;
    }
  }

  // Clean up pipes
  close(proc->stdout_pipe[0]);
  close(proc->stderr_pipe[0]);

  return NULL;
}

// Spawn process with event handler - similar to HTTP server pattern
int64_t spawn_process_with_handler(char *command, ProcessEventHandler handler) {
  if (!command || !handler) {
    return -1;
  }

  pthread_mutex_lock(&process_mutex);

  int64_t process_id = next_process_id++;
  if (process_id >= MAX_PROCESSES) {
    pthread_mutex_unlock(&process_mutex);
    return -2; // Too many processes
  }

  ProcessResult *proc = malloc(sizeof(ProcessResult));
  if (!proc) {
    pthread_mutex_unlock(&process_mutex);
    return -3; // Memory allocation failed
  }

  // Initialize process structure
  proc->process_id = process_id;
  proc->command = strdup(command);
  proc->exit_code = -999; // Not finished yet
  proc->is_running = true;
  proc->handler = handler;
  pthread_mutex_init(&proc->mutex, NULL);

  // Create pipes for stdout and stderr
  if (pipe(proc->stdout_pipe) != 0 || pipe(proc->stderr_pipe) != 0) {
    free(proc->command);
    free(proc);
    pthread_mutex_unlock(&process_mutex);
    return -4; // Pipe creation failed
  }

  // Fork the process
  proc->pid = fork();
  if (proc->pid == 0) {
    // Child process
    close(proc->stdout_pipe[0]); // Close read end
    close(proc->stderr_pipe[0]);

    // Redirect stdout and stderr to pipes
    dup2(proc->stdout_pipe[1], STDOUT_FILENO);
    dup2(proc->stderr_pipe[1], STDERR_FILENO);

    close(proc->stdout_pipe[1]);
    close(proc->stderr_pipe[1]);

    // Execute the command
    execl("/bin/sh", "sh", "-c", command, (char *)NULL);
    _exit(127); // If execl fails
  } else if (proc->pid > 0) {
    // Parent process
    processes[process_id] = proc;

    // Create monitoring thread
    if (pthread_create(&proc->monitor_thread, NULL, process_monitor_thread,
                       proc) != 0) {
      // Thread creation failed, clean up
      close(proc->stdout_pipe[0]);
      close(proc->stdout_pipe[1]);
      close(proc->stderr_pipe[0]);
      close(proc->stderr_pipe[1]);
      kill(proc->pid, SIGTERM);
      waitpid(proc->pid, NULL, 0);
      free(proc->command);
      free(proc);
      processes[process_id] = NULL;
      pthread_mutex_unlock(&process_mutex);
      return -5; // Thread creation failed
    }

    pthread_mutex_unlock(&process_mutex);
    return process_id;
  } else {
    // Fork failed
    close(proc->stdout_pipe[0]);
    close(proc->stdout_pipe[1]);
    close(proc->stderr_pipe[0]);
    close(proc->stderr_pipe[1]);
    free(proc->command);
    free(proc);
    pthread_mutex_unlock(&process_mutex);
    return -6; // Fork failed
  }
}

// Wait for process completion - blocks until process finishes
int64_t await_process(int64_t process_id) {
  if (process_id < 1 || process_id >= MAX_PROCESSES) {
    return -1;
  }

  pthread_mutex_lock(&process_mutex);
  ProcessResult *proc = processes[process_id];
  pthread_mutex_unlock(&process_mutex);

  if (!proc) {
    return -1;
  }

  // Wait for monitor thread to complete
  pthread_join(proc->monitor_thread, NULL);

  return proc->exit_code;
}

// Clean up process resources
void cleanup_process(int64_t process_id) {
  if (process_id < 1 || process_id >= MAX_PROCESSES) {
    return;
  }

  pthread_mutex_lock(&process_mutex);
  ProcessResult *proc = processes[process_id];
  if (proc) {
    processes[process_id] = NULL;

    if (proc->command)
      free(proc->command);
    pthread_mutex_destroy(&proc->mutex);
    free(proc);
  }
  pthread_mutex_unlock(&process_mutex);
}

// Legacy spawn_process function for backward compatibility - now blocking
char *spawn_process(char *command) {
  if (!command) {
    return NULL;
  }

  // Use popen for simple blocking behavior (legacy support)
  FILE *pipe = popen(command, "r");
  if (!pipe) {
    return NULL;
  }

  // Read all output
  char *output = malloc(4096);
  if (!output) {
    pclose(pipe);
    return NULL;
  }

  size_t total_read = 0;
  size_t buffer_size = 4096;
  char buffer[256];

  while (fgets(buffer, sizeof(buffer), pipe) != NULL) {
    size_t len = strlen(buffer);

    // Resize if needed
    if (total_read + len >= buffer_size) {
      buffer_size *= 2;
      output = realloc(output, buffer_size);
      if (!output) {
        pclose(pipe);
        return NULL;
      }
    }

    strcpy(output + total_read, buffer);
    total_read += len;
  }

  output[total_read] = '\0';
  pclose(pipe);

  return output;
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
