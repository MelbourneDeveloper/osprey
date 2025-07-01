#include <assert.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

// Include the system runtime header (we'll define the interface)
extern int64_t spawn_process_with_handler(char *command,
                                          void (*handler)(int64_t, int64_t,
                                                          char *));
extern int64_t await_process(int64_t process_id);
extern void cleanup_process(int64_t process_id);

// Test event handler data
typedef struct {
  int stdout_count;
  int stderr_count;
  int exit_count;
  char last_stdout[1024];
  char last_stderr[1024];
  int64_t last_exit_code;
  pthread_mutex_t mutex;
} TestEventData;

// Test event handler
static void test_event_handler(int64_t process_id, int64_t event_type,
                               char *data) {
  // This would be passed in from test, but for now we'll use a global
  static TestEventData test_data = {0};
  static int initialized = 0;

  if (!initialized) {
    pthread_mutex_init(&test_data.mutex, NULL);
    initialized = 1;
  }

  pthread_mutex_lock(&test_data.mutex);

  switch (event_type) {
  case 1: // PROCESS_STDOUT_DATA
    test_data.stdout_count++;
    strncpy(test_data.last_stdout, data, sizeof(test_data.last_stdout) - 1);
    test_data.last_stdout[sizeof(test_data.last_stdout) - 1] = '\0';
    printf("TEST: Got stdout from process %lld: %s", (long long)process_id,
           data);
    break;
  case 2: // PROCESS_STDERR_DATA
    test_data.stderr_count++;
    strncpy(test_data.last_stderr, data, sizeof(test_data.last_stderr) - 1);
    test_data.last_stderr[sizeof(test_data.last_stderr) - 1] = '\0';
    printf("TEST: Got stderr from process %lld: %s", (long long)process_id,
           data);
    break;
  case 3: // PROCESS_EXIT
    test_data.exit_count++;
    test_data.last_exit_code = atoll(data);
    printf("TEST: Process %lld exited with code: %lld\n", (long long)process_id,
           (long long)test_data.last_exit_code);
    break;
  default:
    printf("TEST: Unknown event type %d from process %lld\n", event_type,
           (long long)process_id);
    break;
  }

  pthread_mutex_unlock(&test_data.mutex);
}

void test_basic_process_spawn() {
  printf("=== Testing Basic Process Spawn ===\n");

  int64_t process_id = spawn_process_with_handler(
      "echo 'Hello from test process'", test_event_handler);

  assert(process_id > 0);
  printf("Process spawned with ID: %lld\n", (long long)process_id);

  // Wait for completion
  int64_t exit_code = await_process(process_id);
  printf("Process completed with exit code: %lld\n", (long long)exit_code);

  assert(exit_code == 0);

  // Clean up
  cleanup_process(process_id);
  printf("Process cleaned up\n");

  printf("=== Basic Process Spawn Test PASSED ===\n\n");
}

void test_multiple_processes() {
  printf("=== Testing Multiple Processes ===\n");

  int64_t pid1 =
      spawn_process_with_handler("echo 'Process 1'", test_event_handler);
  int64_t pid2 =
      spawn_process_with_handler("echo 'Process 2'", test_event_handler);
  int64_t pid3 =
      spawn_process_with_handler("echo 'Process 3'", test_event_handler);

  assert(pid1 > 0);
  assert(pid2 > 0);
  assert(pid3 > 0);
  assert(pid1 != pid2);
  assert(pid2 != pid3);

  printf("Spawned processes: %lld, %lld, %lld\n", (long long)pid1,
         (long long)pid2, (long long)pid3);

  // Wait for all to complete
  int64_t exit1 = await_process(pid1);
  int64_t exit2 = await_process(pid2);
  int64_t exit3 = await_process(pid3);

  assert(exit1 == 0);
  assert(exit2 == 0);
  assert(exit3 == 0);

  // Clean up all
  cleanup_process(pid1);
  cleanup_process(pid2);
  cleanup_process(pid3);

  printf("=== Multiple Processes Test PASSED ===\n\n");
}

void test_process_with_error() {
  printf("=== Testing Process With Error ===\n");

  int64_t process_id = spawn_process_with_handler(
      "false", test_event_handler); // 'false' command returns exit code 1

  assert(process_id > 0);

  // Wait for completion
  int64_t exit_code = await_process(process_id);
  printf("Error process completed with exit code: %lld\n",
         (long long)exit_code);

  assert(exit_code == 1); // false command should return 1

  cleanup_process(process_id);

  printf("=== Process With Error Test PASSED ===\n\n");
}

void test_process_with_stderr() {
  printf("=== Testing Process With Stderr ===\n");

  int64_t process_id = spawn_process_with_handler(
      "sh -c 'echo \"error message\" >&2'", test_event_handler);

  assert(process_id > 0);

  // Wait for completion
  int64_t exit_code = await_process(process_id);

  assert(exit_code == 0);

  cleanup_process(process_id);

  printf("=== Process With Stderr Test PASSED ===\n\n");
}

void test_long_running_process() {
  printf("=== Testing Long Running Process ===\n");

  int64_t process_id = spawn_process_with_handler(
      "sh -c 'for i in 1 2 3; do echo \"Line $i\"; sleep 0.1; done'",
      test_event_handler);

  assert(process_id > 0);

  // Wait for completion
  int64_t exit_code = await_process(process_id);

  assert(exit_code == 0);

  cleanup_process(process_id);

  printf("=== Long Running Process Test PASSED ===\n\n");
}

void test_invalid_command() {
  printf("=== Testing Invalid Command ===\n");

  int64_t process_id = spawn_process_with_handler("nonexistent_command_12345",
                                                  test_event_handler);

  // Should still get a process ID (the failure happens in the child process)
  assert(process_id > 0);

  // Wait for completion - should get exit code 127 (command not found)
  int64_t exit_code = await_process(process_id);
  printf("Invalid command exit code: %lld\n", (long long)exit_code);

  assert(exit_code == 127); // Standard exit code for command not found

  cleanup_process(process_id);

  printf("=== Invalid Command Test PASSED ===\n\n");
}

int main() {
  printf("Running System Runtime Tests...\n\n");

  test_basic_process_spawn();
  test_multiple_processes();
  test_process_with_error();
  test_process_with_stderr();
  test_long_running_process();
  test_invalid_command();

  printf("=== ALL SYSTEM RUNTIME TESTS PASSED ===\n");
  return 0;
}