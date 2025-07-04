#include "unity.h"
#include <pthread.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

// Forward declarations for system runtime functions
extern int64_t spawn_process_with_handler(const char *command,
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

// Global test data
static TestEventData test_data = {0};
static int initialized = 0;

// Test event handler
static void test_event_handler(int64_t process_id, int64_t event_type,
                               char *data) {
  (void)process_id; // Suppress unused parameter warning

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
    break;
  case 2: // PROCESS_STDERR_DATA
    test_data.stderr_count++;
    strncpy(test_data.last_stderr, data, sizeof(test_data.last_stderr) - 1);
    test_data.last_stderr[sizeof(test_data.last_stderr) - 1] = '\0';
    break;
  case 3: // PROCESS_EXIT
    test_data.exit_count++;
    test_data.last_exit_code = atoll(data);
    break;
  default:
    // Unknown event type, do nothing
    break;
  }

  pthread_mutex_unlock(&test_data.mutex);
}

// Unity Test Functions
static void test_basic_process_spawn(void) {
  int64_t process_id = spawn_process_with_handler(
      "echo 'Hello from test process'", test_event_handler);

  TEST_ASSERT_TRUE(process_id > 0);

  int64_t exit_code = await_process(process_id);
  TEST_ASSERT_EQUAL(0, exit_code);

  cleanup_process(process_id);
}

static void test_multiple_processes(void) {
  int64_t pid1 =
      spawn_process_with_handler("echo 'Process 1'", test_event_handler);
  int64_t pid2 =
      spawn_process_with_handler("echo 'Process 2'", test_event_handler);
  int64_t pid3 =
      spawn_process_with_handler("echo 'Process 3'", test_event_handler);

  TEST_ASSERT_TRUE(pid1 > 0);
  TEST_ASSERT_TRUE(pid2 > 0);
  TEST_ASSERT_TRUE(pid3 > 0);
  TEST_ASSERT_TRUE(pid1 != pid2);
  TEST_ASSERT_TRUE(pid2 != pid3);

  int64_t exit1 = await_process(pid1);
  int64_t exit2 = await_process(pid2);
  int64_t exit3 = await_process(pid3);

  TEST_ASSERT_EQUAL(0, exit1);
  TEST_ASSERT_EQUAL(0, exit2);
  TEST_ASSERT_EQUAL(0, exit3);

  cleanup_process(pid1);
  cleanup_process(pid2);
  cleanup_process(pid3);
}

static void test_process_with_error(void) {
  int64_t process_id = spawn_process_with_handler("false", test_event_handler);

  TEST_ASSERT_TRUE(process_id > 0);

  int64_t exit_code = await_process(process_id);
  TEST_ASSERT_EQUAL(1, exit_code);

  cleanup_process(process_id);
}

static void test_process_with_stderr(void) {
  int64_t process_id = spawn_process_with_handler(
      "sh -c 'echo \"error message\" >&2'", test_event_handler);

  TEST_ASSERT_TRUE(process_id > 0);

  int64_t exit_code = await_process(process_id);
  TEST_ASSERT_EQUAL(0, exit_code);

  cleanup_process(process_id);
}

static void test_long_running_process(void) {
  int64_t process_id = spawn_process_with_handler(
      "sh -c 'for i in 1 2 3; do echo \"Line $i\"; sleep 0.1; done'",
      test_event_handler);

  TEST_ASSERT_TRUE(process_id > 0);

  int64_t exit_code = await_process(process_id);
  TEST_ASSERT_EQUAL(0, exit_code);

  cleanup_process(process_id);
}

static void test_invalid_command(void) {
  int64_t process_id = spawn_process_with_handler("nonexistent_command_12345",
                                                  test_event_handler);

  TEST_ASSERT_TRUE(process_id > 0);

  int64_t exit_code = await_process(process_id);
  TEST_ASSERT_EQUAL(127, exit_code); // Command not found

  cleanup_process(process_id);
}

static void test_null_command(void) {
  int64_t process_id = spawn_process_with_handler(NULL, test_event_handler);
  TEST_ASSERT_EQUAL(-1, process_id); // Should handle null gracefully
}

static void test_null_handler(void) {
  int64_t process_id = spawn_process_with_handler("echo 'test'", NULL);
  TEST_ASSERT_EQUAL(-1, process_id); // Should handle null handler gracefully
}

// Main test runner
int main(void) {
  UNITY_BEGIN();

  RUN_TEST(test_basic_process_spawn);
  RUN_TEST(test_multiple_processes);
  RUN_TEST(test_process_with_error);
  RUN_TEST(test_process_with_stderr);
  RUN_TEST(test_long_running_process);
  RUN_TEST(test_invalid_command);
  RUN_TEST(test_null_command);
  RUN_TEST(test_null_handler);

  UNITY_END();
}
