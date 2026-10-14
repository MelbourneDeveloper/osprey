// Mock-driver lifecycle coverage; gpu_device_tests separately proves hardware.
#import <Foundation/Foundation.h>
#import <Metal/Metal.h>
#include "gpu_device.h"
#include "io_error.h"
#include <assert.h>
#include <sys/wait.h>
#include <unistd.h>

static unsigned failure;
@interface OspreyMockMetal : NSObject
@property(nonatomic, strong) NSMutableData *data;
@end
@implementation OspreyMockMetal
- (NSString *)name { return @"mock Metal"; }
- (NSUInteger)maxBufferLength { return failure == 11 ? 1 : 1024; }
- (NSUInteger)maxTotalThreadsPerThreadgroup { return 256; }
- (void *)contents { return self.data.mutableBytes; }
- (id)newCommandQueue { return failure == 2 ? nil : [OspreyMockMetal new]; }
- (id)newLibraryWithSource:(NSString *)source options:(id)options error:(NSError **)error {
  assert(source && !options); (void)error;
  return failure == 3 ? nil : [OspreyMockMetal new];
}
- (id)newFunctionWithName:(NSString *)name {
  assert([name isEqualToString:@"osprey_matmul64"]);
  return failure == 4 ? nil : [OspreyMockMetal new];
}
- (id)newComputePipelineStateWithFunction:(id)function error:(NSError **)error {
  assert(function); (void)error;
  return failure == 5 ? nil : [OspreyMockMetal new];
}
- (id)newBufferWithBytes:(const void *)bytes length:(NSUInteger)length options:(MTLResourceOptions)options {
  assert(bytes && length == 8 && options == MTLResourceStorageModeShared);
  if (failure == 6) return nil;
  OspreyMockMetal *buffer = [OspreyMockMetal new];
  buffer.data = [NSMutableData dataWithBytes:bytes length:length];
  return buffer;
}
- (id)newBufferWithLength:(NSUInteger)length options:(MTLResourceOptions)options {
  assert(length == 8 && options == MTLResourceStorageModeShared);
  if (failure == 7) return nil;
  OspreyMockMetal *buffer = [OspreyMockMetal new];
  buffer.data = [NSMutableData dataWithLength:length];
  return buffer;
}
- (id)commandBuffer { return failure == 8 ? nil : [OspreyMockMetal new]; }
- (id)computeCommandEncoder { return failure == 9 ? nil : [OspreyMockMetal new]; }
- (void)setComputePipelineState:(id)pipeline { assert(pipeline); }
- (void)setBuffer:(id)buffer offset:(NSUInteger)offset atIndex:(NSUInteger)index {
  assert(buffer && !offset && index < 3);
}
- (void)setBytes:(const void *)bytes length:(NSUInteger)length atIndex:(NSUInteger)index {
  assert(bytes && length == 8 && index >= 3 && index <= 5);
}
- (void)dispatchThreads:(MTLSize)size threadsPerThreadgroup:(MTLSize)group {
  assert(size.width == 1 && size.height == 1 && size.depth == 1 && group.width == 256);
}
- (void)endEncoding {}
- (void)commit {}
- (void)waitUntilCompleted {}
- (MTLCommandBufferStatus)status { return failure == 10 ? MTLCommandBufferStatusError : MTLCommandBufferStatusCompleted; }
- (NSError *)error { return [NSError errorWithDomain:@"mock" code:7 userInfo:nil]; }
@end

id<MTLDevice> MTLCreateSystemDefaultDevice(void) {
  return failure == 1 ? nil : (id)[OspreyMockMetal new];
}

static void test_case(unsigned mode) {
  pid_t child = fork();
  assert(child >= 0);
  if (!child) {
    failure = mode;
    int64_t left = 0, right = 0, out = 99;
    int status = osprey_gpu_metal_matmul(&left, &right, &out, 1, 1, 1);
    assert(mode ? status == -1 : status == 1);
    assert(mode ? osp_io_error() != NULL : (out == 0 && osprey_gpu_metal_completed() == 1));
    assert(osprey_gpu_metal_name() != NULL);
    if (!mode) assert(osprey_gpu_metal_matmul(NULL, NULL, NULL, 0, 0, 0) == 1);
    exit(0); // Flush gcov from each fresh initialization state.
  }
  int status = 0;
  assert(waitpid(child, &status, 0) == child);
  assert(WIFEXITED(status) && WEXITSTATUS(status) == 0);
}

int main(void) {
  for (unsigned i = 0; i <= 11; i++) test_case(i);
  puts("[ok] Metal driver: completed dispatch plus every initialization/allocation/launch/readback error");
  return 0;
}
