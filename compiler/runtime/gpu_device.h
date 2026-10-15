#ifndef OSPREY_GPU_DEVICE_H
#define OSPREY_GPU_DEVICE_H
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Shared builtin dispatch contract [GPU-BACKEND-PORTABILITY]. A zero return
// selects the host oracle, one means completed device work, -1 is an error.
int osprey_gpu_device_matmul(const int64_t *left, const int64_t *right, int64_t *out,
                             int64_t rows, int64_t inner, int64_t cols);
const char *osprey_gpu_device_name(void);
void *osprey_gpu_panic(const char *message);

// The first completed dispatch is always reported; `OSPREY_GPU_TRACE=1`
// reports every one.
static inline int osprey_gpu_traced(uint64_t count) {
  const char *trace = getenv("OSPREY_GPU_TRACE");
  return count == 1 || (trace != NULL && strcmp(trace, "1") == 0);
}

// `dims` is the product's rows, inner and cols.
static inline void osprey_gpu_report(const char *device, uint64_t count, const uint64_t dims[3]) {
  if (osprey_gpu_traced(count))
    (void)fprintf(stderr, "[gpu] %s: completed dispatch %llu (%llu x %llu x %llu)\n", device,
                  (unsigned long long)count, (unsigned long long)dims[0],
                  (unsigned long long)dims[1], (unsigned long long)dims[2]);
}

#if defined(OSPREY_GPU_METAL)
int osprey_gpu_metal_matmul(const int64_t *left, const int64_t *right, int64_t *out,
                            int64_t rows, int64_t inner, int64_t cols);
const char *osprey_gpu_metal_name(void);
uint64_t osprey_gpu_metal_completed(void);
#endif
#endif
