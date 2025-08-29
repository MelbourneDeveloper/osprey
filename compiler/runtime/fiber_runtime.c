#include <pthread.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

// Fiber runtime implementation in C for linking with LLVM-generated code

typedef struct Fiber {
  int64_t id;
  int64_t (*function)(void);
  int64_t result;
  bool completed;
  pthread_t thread;
  pthread_mutex_t mutex;
  pthread_cond_t cond;
  bool uses_thread; // Track if this fiber uses threading
} Fiber;

typedef struct Channel {
  int64_t id;
  int capacity;
  int64_t *buffer;
  int head;
  int tail;
  int count;
  pthread_mutex_t mutex;
  pthread_cond_t not_empty;
  pthread_cond_t not_full;
} Channel;

// Global runtime state
static Fiber *fibers[1000];
static Channel *channels[1000];
static int64_t next_id = 1;
static pthread_mutex_t runtime_mutex = PTHREAD_MUTEX_INITIALIZER;

// Deterministic execution mode
static bool deterministic_mode = false;
static int64_t execution_queue[1000];
static int64_t queue_size = 0;

// Enable/disable deterministic fiber execution
void fiber_set_deterministic_mode(bool enabled) {
  pthread_mutex_lock(&runtime_mutex);
  deterministic_mode = enabled;
  if (enabled) {
    queue_size = 0; // Reset queue when enabling
  }
  pthread_mutex_unlock(&runtime_mutex);
}

// Execute a fiber directly (for deterministic mode)
static void execute_fiber_directly(Fiber *fiber) {
  fiber->result = fiber->function();
  fiber->completed = true;
}

// Thread function for executing fibers
static void *fiber_thread_func(void *arg) {
  Fiber *fiber = (Fiber *)arg;

  // Execute the fiber function
  fiber->result = fiber->function();

  // Mark as completed and signal
  pthread_mutex_lock(&fiber->mutex);
  fiber->completed = true;
  pthread_cond_signal(&fiber->cond);
  pthread_mutex_unlock(&fiber->mutex);

  return NULL;
}

// Create and schedule a fiber
int64_t fiber_spawn(int64_t (*fn)(void)) {
  if (!fn) {
    return -1; // Invalid function pointer
  }

  pthread_mutex_lock(&runtime_mutex);

  int64_t id = next_id++;

  // Check if we've exceeded the fiber array bounds
  if (id >= 1000) {
    pthread_mutex_unlock(&runtime_mutex);
    return -4; // Fiber array full
  }

  Fiber *fiber = malloc(sizeof(Fiber));
  if (!fiber) {
    pthread_mutex_unlock(&runtime_mutex);
    return -2; // Memory allocation failed
  }

  fiber->id = id;
  fiber->function = fn;
  fiber->completed = false;
  fiber->uses_thread = false;
  
  if (!deterministic_mode) {
    // Normal concurrent mode - use threads
    pthread_mutex_init(&fiber->mutex, NULL);
    pthread_cond_init(&fiber->cond, NULL);
    fiber->uses_thread = true;

    fibers[id] = fiber;

    // Create thread to execute fiber
    int result = pthread_create(&fiber->thread, NULL, fiber_thread_func, fiber);
    if (result != 0) {
      // Thread creation failed, clean up
      fibers[id] = NULL;
      pthread_mutex_destroy(&fiber->mutex);
      pthread_cond_destroy(&fiber->cond);
      free(fiber);
      pthread_mutex_unlock(&runtime_mutex);
      return -3; // Thread creation failed
    }
  } else {
    // Deterministic mode - queue for sequential execution
    fibers[id] = fiber;
    execution_queue[queue_size++] = id;
  }

  pthread_mutex_unlock(&runtime_mutex);

  return id;
}

// Wait for fiber completion
int64_t fiber_await(int64_t fiber_id) {
  // Check bounds first to prevent buffer overflow
  if (fiber_id < 1 || fiber_id >= 1000) {
    return -1;
  }

  pthread_mutex_lock(&runtime_mutex);
  Fiber *fiber = fibers[fiber_id];
  bool is_deterministic = deterministic_mode;
  pthread_mutex_unlock(&runtime_mutex);

  if (!fiber)
    return -1;

  if (is_deterministic) {
    // Deterministic mode - execute fibers in queue order up to the requested one
    pthread_mutex_lock(&runtime_mutex);
    for (int64_t i = 0; i < queue_size; i++) {
      int64_t current_id = execution_queue[i];
      Fiber *current_fiber = fibers[current_id];
      if (current_fiber && !current_fiber->completed) {
        execute_fiber_directly(current_fiber);
      }
      if (current_id == fiber_id) {
        break; // Stop once we've executed the requested fiber
      }
    }
    int64_t result = fiber->result;
    pthread_mutex_unlock(&runtime_mutex);
    return result;
  } else {
    // Normal concurrent mode - wait for thread completion
    pthread_mutex_lock(&fiber->mutex);
    while (!fiber->completed) {
      pthread_cond_wait(&fiber->cond, &fiber->mutex);
    }
    int64_t result = fiber->result;
    pthread_mutex_unlock(&fiber->mutex);

    // Join thread
    if (fiber->uses_thread) {
      pthread_join(fiber->thread, NULL);
    }

    return result;
  }
}

// TODO: Implement proper fiber yielding with context switching
int64_t fiber_yield(int64_t value) {
  // NOTE: Current implementation is incomplete and needs proper context
  // Don't ignore this. FIX IT!
  return value;
}

// Create a channel
int64_t channel_create(int64_t capacity) {
  pthread_mutex_lock(&runtime_mutex);

  int64_t id = next_id++;
  Channel *channel = malloc(sizeof(Channel));
  channel->id = id;
  channel->capacity = (int)capacity;
  channel->buffer = malloc((size_t)capacity * sizeof(int64_t));
  channel->head = 0;
  channel->tail = 0;
  channel->count = 0;
  pthread_mutex_init(&channel->mutex, NULL);
  pthread_cond_init(&channel->not_empty, NULL);
  pthread_cond_init(&channel->not_full, NULL);

  channels[id] = channel;

  pthread_mutex_unlock(&runtime_mutex);

  return id;
}

// Send value to channel
int64_t channel_send(int64_t channel_id, int64_t value) {
  // Check bounds first to prevent buffer overflow
  if (channel_id < 1 || channel_id >= 1000) {
    return 0;
  }

  pthread_mutex_lock(&runtime_mutex);
  Channel *channel = channels[channel_id];
  pthread_mutex_unlock(&runtime_mutex);

  if (!channel)
    return 0;

  pthread_mutex_lock(&channel->mutex);

  // Wait while channel is full
  while (channel->count == channel->capacity) {
    pthread_cond_wait(&channel->not_full, &channel->mutex);
  }

  // Add value to buffer
  channel->buffer[channel->tail] = value;
  channel->tail = (channel->tail + 1) % channel->capacity;
  channel->count++;

  // Signal that channel is not empty
  pthread_cond_signal(&channel->not_empty);

  pthread_mutex_unlock(&channel->mutex);

  return 1; // Success
}

// Receive from channel
int64_t channel_recv(int64_t channel_id) {
  // Check bounds first to prevent buffer overflow
  if (channel_id < 1 || channel_id >= 1000) {
    return -1;
  }

  pthread_mutex_lock(&runtime_mutex);
  Channel *channel = channels[channel_id];
  pthread_mutex_unlock(&runtime_mutex);

  if (!channel)
    return -1;

  pthread_mutex_lock(&channel->mutex);

  // Wait while channel is empty
  while (channel->count == 0) {
    pthread_cond_wait(&channel->not_empty, &channel->mutex);
  }

  // Get value from buffer
  int64_t value = channel->buffer[channel->head];
  channel->head = (channel->head + 1) % channel->capacity;
  channel->count--;

  // Signal that channel is not full
  pthread_cond_signal(&channel->not_full);

  pthread_mutex_unlock(&channel->mutex);

  return value;
}

// Sleep for specified milliseconds
int64_t fiber_sleep(int64_t milliseconds) {
  usleep((unsigned int)(milliseconds * 1000)); // Convert milliseconds to microseconds
  return 0;
}

// FIBER-BASED PROCESS SPAWNING FUNCTIONS
// These functions integrate process spawning with the fiber runtime

// External process functions from system_runtime.c
extern int64_t spawn_process_with_handler(char *command,
                                          void (*handler)(int64_t, int64_t,
                                                          char *));
extern int64_t await_process(int64_t process_id);
extern void cleanup_process(int64_t process_id);

// Process event types (matching system_runtime.c)
#define PROCESS_STDOUT_DATA 1
#define PROCESS_STDERR_DATA 2
#define PROCESS_EXIT 3

// Simple process event handler that just prints output (demo implementation)
static void default_process_event_handler(int64_t process_id,
                                          int64_t event_type, char *data) {
  switch (event_type) {
  case PROCESS_STDOUT_DATA:
    printf("Process %lld stdout: %s", (long long)process_id, data);
    break;
  case PROCESS_STDERR_DATA:
    printf("Process %lld stderr: %s", (long long)process_id, data);
    break;
  case PROCESS_EXIT:
    printf("Process %lld exited with code: %s\n", (long long)process_id, data);
    break;
  default:
    printf("Process %lld unknown event %lld: %s\n", (long long)process_id,
           (long long)event_type, data);
    break;
  }
}

// Spawn a process with event handler - returns process ID
int64_t fiber_spawn_process(char *command) {
  if (!command) {
    return -1;
  }

  // Spawn process with default event handler
  return spawn_process_with_handler(command, default_process_event_handler);
}

// Spawn a process with custom handler - for advanced use cases
int64_t fiber_spawn_process_with_handler(char *command,
                                         void (*handler)(int64_t, int64_t,
                                                         char *)) {
  if (!command || !handler) {
    return -1;
  }

  return spawn_process_with_handler(command, handler);
}

// Await process completion in fiber context
int64_t fiber_await_process(int64_t process_id) {
  return await_process(process_id);
}

// Await process completion with stdout callback in fiber context
int64_t fiber_await_process_with_callback(int64_t process_id,
                                          void (*stdout_callback)(char *)) {
  if (!stdout_callback) {
    return await_process(process_id);
  }

  return fiber_await_process(process_id);
}

// Clean up process resources
void fiber_cleanup_process(int64_t process_id) { cleanup_process(process_id); }
