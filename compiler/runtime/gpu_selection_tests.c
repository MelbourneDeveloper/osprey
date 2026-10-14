#include "gpu_device.h"
#include "gpu_opencl_api.h"
#include "io_error.h"
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static int device_error;
int osprey_gpu_buffer_tests(int argc, char **argv);
#if defined(OSPREY_GPU_METAL)
const char *osprey_gpu_metal_name(void) { return "metal:mock"; }
int osprey_gpu_metal_matmul(const int64_t *a, const int64_t *b, int64_t *out, int64_t r, int64_t k, int64_t c) {
  assert(a && b && r == 1 && k == 1 && c == 1);
  *out = 42;
  return device_error ? -1 : 1;
}
#endif
#if defined(OSPREY_GPU_OPENCL)
const char *osprey_gpu_opencl_name(void) { return "opencl:mock"; }
int osprey_gpu_opencl_matmul(const int64_t *a, const int64_t *b, int64_t *out, int64_t r, int64_t k, int64_t c) {
  assert(a && b && r == 1 && k == 1 && c == 1);
  *out = 43;
  return device_error ? -1 : 1;
}
#endif

static void host_case(const char *mode) {
  if (mode) assert(setenv("OSPREY_GPU", mode, 1) == 0);
  else assert(unsetenv("OSPREY_GPU") == 0);
  int64_t a = 1, out = 99;
  assert(strcmp(osprey_gpu_device_name(), "host") == 0);
  assert(osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1) == 0 && out == 99);
}

static void device_case(const char *mode, int available, const char *name) {
  assert(setenv("OSPREY_GPU", mode, 1) == 0);
  int64_t a = 1, out = 99;
  osp_io_error_clear();
  assert(strcmp(osprey_gpu_device_name(), name) == 0);
  int status = osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1);
  assert(available ? status == 1 : status == -1);
  if (available) {
    device_error = 1;
    assert(osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1) == -1);
    device_error = 0;
  } else assert(strstr(osp_io_error(), "no GPU adapter") != NULL);
}

int main(void) {
  host_case(NULL); host_case(""); host_case("host");
  assert(osprey_gpu_buffer_tests(0, NULL) == 0);
#if defined(OSPREY_GPU_METAL)
  device_case("auto", 1, "metal:mock"); device_case("metal", 1, "metal:mock");
#else
  device_case("auto", 0, "unavailable:device"); device_case("metal", 0, "unavailable:device");
#endif
#if defined(OSPREY_GPU_OPENCL)
  device_case("opencl", 1, "opencl:mock");
#else
  device_case("opencl", 0, "unavailable:device");
#endif
  assert(setenv("OSPREY_GPU", "invalid", 1) == 0);
  assert(strcmp(osprey_gpu_device_name(), "invalid:device") == 0);
  int64_t a = 1, out = 99;
  assert(osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1) == -1 && out == 99);
  assert(strstr(osp_io_error(), "unknown backend") != NULL);
  puts("[ok] GPU selection: host, auto, Metal, OpenCL, missing adapter and no silent CPU fallback");
  return 0;
}
