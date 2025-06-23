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
  pthread_mutex_init(&fiber->mutex, NULL);
  pthread_cond_init(&fiber->cond, NULL);

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
  pthread_mutex_unlock(&runtime_mutex);

  if (!fiber)
    return -1;

  // Wait for fiber to complete
  pthread_mutex_lock(&fiber->mutex);
  while (!fiber->completed) {
    pthread_cond_wait(&fiber->cond, &fiber->mutex);
  }
  int64_t result = fiber->result;
  pthread_mutex_unlock(&fiber->mutex);

  // Join thread
  pthread_join(fiber->thread, NULL);

  return result;
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
  channel->capacity = capacity;
  channel->buffer = malloc(capacity * sizeof(int64_t));
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
  usleep(milliseconds * 1000); // Convert milliseconds to microseconds
  return 0;
}