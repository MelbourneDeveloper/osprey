// Mock-driver lifecycle coverage; gpu_device_tests separately proves hardware.
#import <Foundation/Foundation.h>
#import <Metal/Metal.h>
#include "gpu_device.h"
#include "io_error.h"
#include <assert.h>
#include <sys/wait.h>
#include <unistd.h>

#define MOCK_MARK 0x2a
static unsigned failure, libraries, functions, bound;
static NSUInteger expected_threads = 1;
static NSMutableArray *created;

@interface OspreyMockMetal : NSObject
@property(nonatomic, strong) NSMutableData *data;
@end
@implementation OspreyMockMetal
- (NSString *)name { return @"mock Metal"; }
- (NSUInteger)maxBufferLength { return failure == 10 ? 1 : 1 << 20; }
- (NSUInteger)maxTotalThreadsPerThreadgroup { return 256; }
- (void *)contents { return self.data.mutableBytes; }
- (id)newCommandQueue { return failure == 2 ? nil : [OspreyMockMetal new]; }
- (id)newLibraryWithSource:(NSString *)source options:(id)options error:(NSError **)error {
  assert([source containsString:@"osprey_matmul64"] && !options); (void)error;
  libraries++;
  return failure == 3 ? nil : [OspreyMockMetal new];
}
- (id)newFunctionWithName:(NSString *)name {
  assert(name.length);
  functions++;
  return failure == 4 ? nil : [OspreyMockMetal new];
}
- (id)newComputePipelineStateWithFunction:(id)function error:(NSError **)error {
  assert(function); (void)error;
  return failure == 5 ? nil : [OspreyMockMetal new];
}
- (id)newBufferWithBytes:(const void *)bytes length:(NSUInteger)length options:(MTLResourceOptions)options {
  assert(bytes && length >= 8 && options == MTLResourceStorageModeShared);
  if (failure == 6) return nil;
  OspreyMockMetal *buffer = [OspreyMockMetal new];
  buffer.data = [NSMutableData dataWithBytes:bytes length:length];
  [created addObject:buffer];
  return buffer;
}
- (id)commandBuffer { return failure == 7 ? nil : [OspreyMockMetal new]; }
- (id)computeCommandEncoder { return failure == 8 ? nil : [OspreyMockMetal new]; }
- (void)setComputePipelineState:(id)pipeline { assert(pipeline); }
- (void)setBuffer:(id)buffer offset:(NSUInteger)offset atIndex:(NSUInteger)index {
  assert(buffer && !offset && index < OSPREY_GPU_MAX_ARGS);
  bound++;
}
- (void)setBytes:(const void *)bytes length:(NSUInteger)length atIndex:(NSUInteger)index {
  assert(bytes && length == 8 && index < OSPREY_GPU_MAX_ARGS);
  bound++;
}
- (void)dispatchThreads:(MTLSize)size threadsPerThreadgroup:(MTLSize)group {
  assert(size.width == expected_threads && size.height == 1 && size.depth == 1 && group.width == 256);
}
- (void)endEncoding {}
// The mock "runs" every kernel by stamping word 0 of each buffer it was given.
- (void)commit {
  for (OspreyMockMetal *buffer in created) ((int64_t *)buffer.data.mutableBytes)[0] = MOCK_MARK;
  [created removeAllObjects];
}
- (void)waitUntilCompleted {}
- (MTLCommandBufferStatus)status { return failure == 9 ? MTLCommandBufferStatusError : MTLCommandBufferStatusCompleted; }
- (NSError *)error { return [NSError errorWithDomain:@"mock" code:7 userInfo:nil]; }
@end

id<MTLDevice> MTLCreateSystemDefaultDevice(void) {
  return failure == 1 ? nil : (id)[OspreyMockMetal new];
}

static int product(int64_t *out) {
  int64_t left = 0, right = 0;
  OspreyGpuDispatch d = {"osprey_matmul64", NULL, {{0, NULL, NULL, 0}}, 0, 1};
  osprey_gpu_arg(&d, OSPREY_GPU_INPUT, &left, NULL, 1);
  osprey_gpu_arg(&d, OSPREY_GPU_INPUT, &right, NULL, 1);
  osprey_gpu_arg(&d, OSPREY_GPU_OUTPUT, out, out, 1);
  for (uint64_t i = 0; i < 3; i++) osprey_gpu_arg(&d, OSPREY_GPU_SCALAR, NULL, NULL, 1);
  expected_threads = 1;
  return osprey_gpu_metal_dispatch(&d);
}

static int compiled(const char *source, int64_t *out, uint64_t threads) {
  int64_t in = 7, fault = 0;
  OspreyGpuDispatch d = {"osprey_kernel_1", source, {{0, NULL, NULL, 0}}, 0, threads};
  osprey_gpu_arg(&d, OSPREY_GPU_INPUT, &in, NULL, 1);
  osprey_gpu_arg(&d, OSPREY_GPU_INPUT, &in, NULL, 0);
  osprey_gpu_arg(&d, OSPREY_GPU_OUTPUT, out, out, 2);
  osprey_gpu_arg(&d, OSPREY_GPU_OUTPUT, &fault, &fault, 1);
  osprey_gpu_arg(&d, OSPREY_GPU_SCALAR, NULL, NULL, threads);
  expected_threads = (NSUInteger)threads;
  int status = osprey_gpu_metal_dispatch(&d);
  assert(status != 1 || fault == MOCK_MARK);
  return status;
}

// The success path: the product, then a program's kernels, which recompile
// the library once and reuse their pipeline afterwards.
static void success_case(void) {
  static const char source[] = "kernel void osprey_kernel_1() {}";
  int64_t out[2] = {99, 99};
  assert(product(out) == 1 && out[0] == MOCK_MARK && osprey_gpu_metal_completed() == 1);
  assert(libraries == 1 && functions == 1 && bound == 6);
  assert(compiled(source, out, 3) == 1 && libraries == 2 && functions == 2 && out[1] == 99);
  assert(compiled(source, out, 5) == 1 && libraries == 2 && functions == 2);
  assert(product(out) == 1 && libraries == 2 && functions == 3 && osprey_gpu_metal_completed() == 4);
  assert(strcmp(osprey_gpu_metal_name(), "metal:0 (mock Metal)") == 0);
}

static void test_case(unsigned mode) {
  pid_t child = fork();
  assert(child >= 0);
  if (!child) {
    failure = mode;
    created = [NSMutableArray new];
    int64_t out = 99;
    if (!mode) success_case();
    else {
      assert(product(&out) == -1 && osp_io_error() != NULL && out == 99);
      assert(strcmp(osprey_gpu_metal_name(), mode <= 3 ? "metal:unavailable" : "metal:0 (mock Metal)") == 0);
    }
    exit(0); // Flush gcov from each fresh initialization state.
  }
  int status = 0;
  assert(waitpid(child, &status, 0) == child);
  assert(WIFEXITED(status) && WEXITSTATUS(status) == 0);
}

int main(void) {
  for (unsigned i = 0; i <= 10; i++) test_case(i);
  puts("[ok] Metal driver: product and compiled-kernel dispatch, library recompilation, pipeline reuse, every initialization/allocation/launch/readback error");
  return 0;
}
