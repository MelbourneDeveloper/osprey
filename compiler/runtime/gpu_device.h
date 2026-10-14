#ifndef OSPREY_GPU_DEVICE_H
#define OSPREY_GPU_DEVICE_H
#include <stdint.h>

// Shared builtin dispatch contract [GPU-BACKEND-PORTABILITY]. A zero return
// selects the host oracle, one means completed device work, -1 is an error.
int osprey_gpu_device_matmul(const int64_t *left, const int64_t *right, int64_t *out,
                             int64_t rows, int64_t inner, int64_t cols);
const char *osprey_gpu_device_name(void);
void *osprey_gpu_panic(const char *message);

#if defined(OSPREY_GPU_METAL)
int osprey_gpu_metal_matmul(const int64_t *left, const int64_t *right, int64_t *out,
                            int64_t rows, int64_t inner, int64_t cols);
const char *osprey_gpu_metal_name(void);
uint64_t osprey_gpu_metal_completed(void);
#endif
#endif
