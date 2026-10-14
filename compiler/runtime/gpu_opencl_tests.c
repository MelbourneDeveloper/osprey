// Driver lifecycle/error tests without GPU hardware. Hardware execution and
// numerical parity are independently gated by gpu_device_tests.c.
#include "gpu_opencl_api.h"
#include "io_error.h"
#include <assert.h>
#include <stdio.h>
#include <string.h>

static unsigned calls, fail_at, live, next_handle;
static int status(void) { return ++calls == fail_at ? -7 : 0; }
static void *handle(int *error) {
  *error = status();
  if (*error) return NULL;
  live++;
  return (void *)(uintptr_t)++next_handle;
}
static int OSP_CL_CALL platforms(unsigned count, void **out, unsigned *n) {
  if (n) *n = 1;
  if (count && out) *out = (void *)(uintptr_t)42;
  return status();
}
static int OSP_CL_CALL devices(void *platform, uint64_t type, unsigned count, void **out, unsigned *n) {
  assert(platform && type == 4 && count == 1 && !n);
  int error = status();
  if (!error) *out = (void *)(uintptr_t)43;
  return error;
}
static int OSP_CL_CALL info(void *device, unsigned key, size_t size, void *out, size_t *n) {
  assert(device && key == 0x102B && size >= 5 && !n);
  memcpy(out, "mock", 5);
  return status();
}
static void *OSP_CL_CALL context(const intptr_t *p, unsigned count, void *const *d, void *notify, void *user, int *error) {
  assert(!p && count == 1 && *d && !notify && !user);
  return handle(error);
}
static void *OSP_CL_CALL queue(void *ctx, void *device, uint64_t properties, int *error) {
  assert(ctx && device && !properties);
  return handle(error);
}
static void *OSP_CL_CALL program(void *ctx, unsigned count, const char **source, const size_t *size, int *error) {
  assert(ctx && count == 1 && *size && *source);
  return handle(error);
}
static int OSP_CL_CALL build(void *p, unsigned count, void *const *device, const char *options, void *notify, void *user) {
  assert(p && count == 1 && *device && strstr(options, "CL1.2") && !notify && !user);
  return status();
}
static int OSP_CL_CALL build_info(void *p, void *d, unsigned key, size_t size, void *out, size_t *n) {
  assert(p && d && key == 0x1183 && size > 5 && !n);
  memcpy(out, "build", 6);
  return 0;
}
static void *OSP_CL_CALL kernel(void *p, const char *name, int *error) {
  assert(p && strcmp(name, "osprey_matmul64") == 0);
  return handle(error);
}
static void *OSP_CL_CALL buffer(void *ctx, uint64_t flags, size_t bytes, const void *source, int *error) {
  assert(ctx && bytes == 8 && (flags == 36 ? source != NULL : source == NULL));
  return handle(error);
}
static int OSP_CL_CALL arg(void *k, unsigned i, size_t bytes, const void *value) {
  assert(k && i < 6 && bytes == 8 && value);
  return status();
}
static int OSP_CL_CALL launch(void *q, void *k, unsigned rank, const size_t *offset, const size_t *global,
                            const size_t *local, unsigned count, const void *events, void *event) {
  assert(q && k && rank == 1 && !offset && *global == 1 && !local && !count && !events && !event);
  return status();
}
static int OSP_CL_CALL read_buffer(void *q, void *b, unsigned blocking, size_t offset, size_t size,
                                  void *out, unsigned count, const void *events, void *event) {
  assert(q && b && blocking == 1 && !offset && size == 8 && !count && !events && !event);
  memset(out, 0, size);
  return status();
}
static int OSP_CL_CALL release(void *object) { assert(object && live); live--; return 0; }

static const OspreyOpenCLApi mock = {platforms, devices, info, context, queue, program, build,
    build_info, kernel, buffer, arg, launch, read_buffer, release, release, release, release, release};

static void failure_case(unsigned point) {
  osprey_gpu_opencl_test_driver(&mock);
  assert(live == 0);
  calls = 0; fail_at = point;
  int64_t left = 0, right = 0, out = 99;
  osp_io_error_clear();
  int result = osprey_gpu_opencl_matmul(&left, &right, &out, 1, 1, 1);
  assert(point ? result == -1 : result == 1);
  assert(point ? osp_io_error() != NULL : (out == 0 && osprey_gpu_opencl_completed() == 1));
  assert(live <= 4); // Failed initialization releases everything; launches retain only the context.
}

int main(void) {
  // Exercise the real dynamic loader even on a machine without a compute GPU.
  assert(osprey_gpu_opencl_name() != NULL);
  for (unsigned point = 0; point <= 20; point++) failure_case(point);
  osprey_gpu_opencl_test_driver(&mock);
  assert(live == 0);
  puts("[ok] OpenCL driver: success plus every initialization/transfer/launch/readback failure releases resources");
  return 0;
}
