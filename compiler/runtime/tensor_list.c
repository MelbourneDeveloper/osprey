// `loadTensor` [BUILTIN-FILE-TENSOR]: one safetensors tensor read by the
// shared reader (tensor_file.h) straight into a host List<float>. No GpuBuffer
// and no osprey_gpu_* symbol exists on this path, so a program that loads its
// weights here links nothing of the GPU surface.

#include "collection_runtime.h"
#include "memory_hooks.h"
#include "tensor_file.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static const char OP_LOAD_TENSOR[] = "loadTensor";

// Push each value as the IEEE-754 word a List<float> element holds.
static void into_list(void *ctx, const double *values, size_t n) {
  for (size_t i = 0; i < n; i++) {
    int64_t word;
    memcpy(&word, values + i, sizeof word);
    osprey_list_builder_push_of(ctx, word, 0);
  }
}

// One tensor out of a safetensors file as a List<float>, or NULL with the
// reason recorded. A list built before a failed read is released, so a NULL
// result owes the memory backend nothing.
OspreyList *osprey_load_tensor(const char *path, const char *name) {
  OspTensorSpan span;
  if (!osp_tensor_locate(OP_LOAD_TENSOR, path, name, &span)) {
    return NULL;
  }
  OspreyListBuilder *b = osprey_list_builder_new_of(0);
  int read = osp_tensor_read(OP_LOAD_TENSOR, path, &span, into_list, b);
  OspreyList *list = osprey_list_builder_seal(b);
  if (!read) {
    osp_release(list);
    return NULL;
  }
  return list;
}

// Broken-code quarantine: a host interpretation whose weights fail to load
// would otherwise run on whatever its caller substituted. It aborts naming why.
void *osprey_tensor_panic(const char *message) {
  (void)fprintf(stderr, "tensor load cannot continue: %s\n", message ? message : "unknown failure");
  abort();
}
