// Assertion-driven tests for tensor_list.c, `loadTensor` [BUILTIN-FILE-TENSOR]:
// the host loader must satisfy the same safetensors contract as
// `gpuLoadTensor` (tensor_file_tests.h) while building a List<float> and
// linking no GPU unit at all — this suite links tensor_list.c, tensor_file.c,
// list_runtime.c, file_runtime.c (the io_error channel) and memory_runtime.c,
// and nothing named gpu_*.
#include <assert.h>

#include "collection_runtime.h"

OspreyList *osprey_load_tensor(const char *path, const char *name);
const char *osp_io_error(void);
void osp_io_error_clear(void);

static long g_checks = 0;
#define CHECK(c)                                                               \
  do {                                                                         \
    g_checks++;                                                                \
    assert(c);                                                                 \
  } while (0)

#include "tensor_file_tests.h"

static void *load(const char *path, const char *name) {
  return osprey_load_tensor(path, name);
}

static int64_t length(void *list) {
  return osprey_list_length(list);
}

// The double a List<float> element word holds.
static double at(void *list, int64_t i) {
  CHECK(osprey_list_in_bounds(list, i));
  int64_t word = osprey_list_get(list, i);
  double d;
  memcpy(&d, &word, sizeof d);
  return d;
}

int main(void) {
  t_tensor_loader(&(TensorLoader){"loadTensor", load, length, at});
  printf("[ok] tensor_list: %ld assertions\n", g_checks);
  return 0;
}
