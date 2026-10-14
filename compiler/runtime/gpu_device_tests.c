// Hardware parity gate [GPU-MATMUL] [GPU-BACKEND-PORTABILITY]. Run explicitly
// on a GPU runner: make _test_gpu_device. Host and device use identical words.
#include "gpu_device.h"
#include "io_error.h"
#include <assert.h>
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void *osprey_gpu_alloc(int64_t length);
void *osprey_gpu_matmul(void *left, void *right, int64_t rows, int64_t inner, int64_t cols);
int64_t *osprey_gpu_data(void *buffer);
static uint64_t state = 0x3f27d489aeb5106c;
static unsigned checks;
static const char *selected;
static char selected_storage[16];

static uint64_t random_word(void) {
  state ^= state << 13;
  state ^= state >> 7;
  state ^= state << 17;
  return state;
}

static void *random_buffer(int64_t count, int low, int range) {
  void *buffer = osprey_gpu_alloc(count);
  int64_t *words = osprey_gpu_data(buffer);
  assert(words);
  for (int64_t i = 0; i < count; i++) {
    uint64_t r = random_word();
    uint64_t e = (uint64_t)low + ((r >> 52) % (uint64_t)range);
    uint64_t word = (r & UINT64_C(0x800fffffffffffff)) | (e << 52);
    memcpy(words + i, &word, sizeof word);
  }
  return buffer;
}

static int nan_word(int64_t word) {
  double value;
  memcpy(&value, &word, sizeof value);
  return isnan(value);
}

static void compare(void *left, void *right, int64_t rows, int64_t inner, int64_t cols) {
  assert(setenv("OSPREY_GPU", "host", 1) == 0);
  osp_io_error_clear();
  void *host = osprey_gpu_matmul(left, right, rows, inner, cols);
  assert(host && !osp_io_error());
  assert(setenv("OSPREY_GPU", selected, 1) == 0);
  void *device = osprey_gpu_matmul(left, right, rows, inner, cols);
  if (!device) fprintf(stderr, "%s\n", osp_io_error());
  assert(device && !osp_io_error());
  const int64_t *a = osprey_gpu_data(host), *b = osprey_gpu_data(device);
  for (int64_t i = 0; i < rows * cols; i++) {
    if (a[i] != b[i] && !(nan_word(a[i]) && nan_word(b[i])))
      fprintf(stderr, "parity at %lld: host %016llx device %016llx\n",
          (long long)i, (unsigned long long)a[i], (unsigned long long)b[i]);
    assert(a[i] == b[i] || (nan_word(a[i]) && nan_word(b[i])));
    checks++;
  }
}

static void extremes(void) {
  uint64_t words[] = {0, UINT64_C(0x8000000000000000), 1, UINT64_C(0x000fffffffffffff),
      UINT64_C(0x0010000000000000), UINT64_C(0x3ff0000000000000),
      UINT64_C(0x3ff0000000000001), UINT64_C(0x7fefffffffffffff),
      UINT64_C(0x7ff0000000000000), UINT64_C(0xfff0000000000000)};
  int64_t count = (int64_t)(sizeof words / sizeof words[0]);
  void *left = osprey_gpu_alloc(count), *right = osprey_gpu_alloc(count);
  memcpy(osprey_gpu_data(left), words, sizeof words);
  memcpy(osprey_gpu_data(right), words, sizeof words);
  compare(left, right, count, 1, count);
}

int main(void) {
  selected = getenv("OSPREY_GPU");
  (void)snprintf(selected_storage, sizeof selected_storage, "%s", selected ? selected : "auto");
  selected = selected_storage;
  extremes();
  compare(random_buffer(64, 0, 2047), random_buffer(64, 0, 2047), 64, 1, 64);
  compare(random_buffer(32 * 511, 700, 600), random_buffer(511 * 16, 700, 600), 32, 511, 16);
  compare(random_buffer(32 * 511, 0, 100), random_buffer(511 * 16, 0, 2047), 32, 511, 16);
  compare(random_buffer(32 * 511, 1020, 7), random_buffer(511 * 16, 1020, 7), 32, 511, 16);
#if defined(OSPREY_GPU_METAL)
  if (strcmp(selected, "opencl") != 0) assert(osprey_gpu_metal_completed() == 5);
#endif
  printf("[ok] device binary64 parity: %u results, five completed GPU dispatches\n", checks);
  return 0;
}
