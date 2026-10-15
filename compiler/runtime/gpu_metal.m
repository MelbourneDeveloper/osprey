// Direct Metal compute adapter [GPU-BACKEND-PORTABILITY] [GPU-MATMUL].
#import <Foundation/Foundation.h>
#import <Metal/Metal.h>
#include "gpu_device.h"
#include "io_error.h"
#include "../bin/gpu_shader_source.h"
#include <pthread.h>
#include <stdatomic.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static id<MTLDevice> gpu_device;
static id<MTLCommandQueue> gpu_queue;
static id<MTLComputePipelineState> gpu_pipeline;
static pthread_once_t gpu_once = PTHREAD_ONCE_INIT;
static _Atomic uint64_t gpu_completed = 0;
static char gpu_failure[512];
static char gpu_name[256];
typedef struct { uint64_t rows, inner, cols; } GpuShape;

static void init_pipeline(void) {
  NSError *error = nil;
  NSString *source = [[NSString alloc] initWithBytes:osprey_gpu_shader
      length:osprey_gpu_shader_len encoding:NSUTF8StringEncoding];
  id<MTLLibrary> library = [gpu_device newLibraryWithSource:source options:nil error:&error];
  id<MTLFunction> function = [library newFunctionWithName:@"osprey_matmul64"];
  if (function) gpu_pipeline = [gpu_device newComputePipelineStateWithFunction:function error:&error];
  if (!gpu_pipeline) (void)snprintf(gpu_failure, sizeof gpu_failure, "%s",
      error ? error.localizedDescription.UTF8String : "Metal kernel compilation failed");
}

static void init_device(void) {
  @autoreleasepool {
    gpu_device = MTLCreateSystemDefaultDevice();
    gpu_queue = [gpu_device newCommandQueue];
    if (!gpu_device || !gpu_queue) {
      (void)snprintf(gpu_failure, sizeof gpu_failure, "no Metal compute device or command queue");
      return;
    }
    (void)snprintf(gpu_name, sizeof gpu_name, "metal:0 (%s; matrix kernels)", gpu_device.name.UTF8String);
    init_pipeline();
  }
}

const char *osprey_gpu_metal_name(void) {
  (void)pthread_once(&gpu_once, init_device);
  return gpu_pipeline ? gpu_name : "metal:unavailable";
}

uint64_t osprey_gpu_metal_completed(void) { return atomic_load(&gpu_completed); }

static id<MTLBuffer> input_buffer(const int64_t *words, uint64_t count) {
  return [gpu_device newBufferWithBytes:words length:(NSUInteger)(count * 8)
      options:MTLResourceStorageModeShared];
}

static int completed(id<MTLCommandBuffer> command, const GpuShape *shape) {
  [command commit];
  [command waitUntilCompleted];
  if (command.status != MTLCommandBufferStatusCompleted) {
    osp_io_error_text("gpuMatMul", "Metal", command.error.localizedDescription.UTF8String);
    return -1;
  }
  uint64_t count = atomic_fetch_add(&gpu_completed, 1) + 1;
  osprey_gpu_report(gpu_name, count, (const uint64_t[3]){shape->rows, shape->inner, shape->cols});
  return 1;
}

static void bind_product(id<MTLComputeCommandEncoder> encoder, id<MTLBuffer> left,
                         id<MTLBuffer> right, id<MTLBuffer> out, const GpuShape *shape) {
  [encoder setComputePipelineState:gpu_pipeline];
  [encoder setBuffer:left offset:0 atIndex:0];
  [encoder setBuffer:right offset:0 atIndex:1];
  [encoder setBuffer:out offset:0 atIndex:2];
  [encoder setBytes:&shape->rows length:sizeof shape->rows atIndex:3];
  [encoder setBytes:&shape->inner length:sizeof shape->inner atIndex:4];
  [encoder setBytes:&shape->cols length:sizeof shape->cols atIndex:5];
}

static int encode_product(id<MTLBuffer> left, id<MTLBuffer> right, id<MTLBuffer> out,
                          const GpuShape *shape) {
  id<MTLCommandBuffer> command = [gpu_queue commandBuffer];
  id<MTLComputeCommandEncoder> encoder = [command computeCommandEncoder];
  if (!command || !encoder) { osp_io_error_text("gpuMatMul", "Metal", "command allocation failed"); return -1; }
  bind_product(encoder, left, right, out, shape);
  NSUInteger count = (NSUInteger)(shape->rows * shape->cols);
  NSUInteger group = MIN(gpu_pipeline.maxTotalThreadsPerThreadgroup, (NSUInteger)256);
  [encoder dispatchThreads:MTLSizeMake(count, 1, 1) threadsPerThreadgroup:MTLSizeMake(group, 1, 1)];
  [encoder endEncoding];
  return completed(command, shape);
}

static int product(const int64_t *left, const int64_t *right, int64_t *out, const GpuShape *shape) {
  uint64_t bytes = shape->rows * shape->cols * 8;
  if (shape->rows * shape->cols > UINT32_MAX || bytes > gpu_device.maxBufferLength ||
      shape->rows * shape->inner * 8 > gpu_device.maxBufferLength ||
      shape->inner * shape->cols * 8 > gpu_device.maxBufferLength) {
    osp_io_error_text("gpuMatMul", "Metal", "matrix exceeds device buffer or dispatch limits"); return -1;
  }
  id<MTLBuffer> a = input_buffer(left, shape->rows * shape->inner);
  id<MTLBuffer> b = input_buffer(right, shape->inner * shape->cols);
  id<MTLBuffer> c = [gpu_device newBufferWithLength:(NSUInteger)bytes options:MTLResourceStorageModeShared];
  if (!a || !b || !c) { osp_io_error_text("gpuMatMul", "Metal", "buffer allocation failed"); return -1; }
  int status = encode_product(a, b, c, shape);
  if (status == 1) memcpy(out, c.contents, (size_t)bytes);
  return status;
}

int osprey_gpu_metal_matmul(const int64_t *left, const int64_t *right, int64_t *out,
                            int64_t rows, int64_t inner, int64_t cols) {
  (void)pthread_once(&gpu_once, init_device);
  if (!gpu_pipeline) { osp_io_error_text("gpuMatMul", "Metal", gpu_failure); return -1; }
  if (!rows || !inner || !cols) return 1;
  GpuShape shape = {(uint64_t)rows, (uint64_t)inner, (uint64_t)cols};
  @autoreleasepool { return product(left, right, out, &shape); }
}
