// Portable direct-driver adapter [GPU-BACKEND-PORTABILITY] [GPU-MATMUL].
#include "gpu_opencl_api.h"
#include "io_error.h"
#include "../bin/gpu_shader_source.h"
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#if defined(_WIN32)
#include <windows.h>
#else
#include <dlfcn.h>
#endif

static OspreyOpenCLApi cl;
static void *cl_device, *cl_context, *cl_queue, *cl_program, *cl_kernel;
static pthread_mutex_t cl_lock = PTHREAD_MUTEX_INITIALIZER;
static uint64_t cl_completed;
static int cl_ready;
static char cl_name[256];
static char cl_failure[1024];

#define CL_FUNCTIONS(X) \
  X(platforms, "clGetPlatformIDs") X(devices, "clGetDeviceIDs") \
  X(info, "clGetDeviceInfo") X(context, "clCreateContext") \
  X(queue, "clCreateCommandQueue") X(program, "clCreateProgramWithSource") \
  X(build, "clBuildProgram") X(build_info, "clGetProgramBuildInfo") \
  X(kernel, "clCreateKernel") X(buffer, "clCreateBuffer") \
  X(arg, "clSetKernelArg") X(launch, "clEnqueueNDRangeKernel") \
  X(read, "clEnqueueReadBuffer") X(release_buffer, "clReleaseMemObject") \
  X(release_kernel, "clReleaseKernel") X(release_program, "clReleaseProgram") \
  X(release_queue, "clReleaseCommandQueue") X(release_context, "clReleaseContext")

#if defined(OSPREY_GPU_TESTING)
static const OspreyOpenCLApi *test_driver;
#endif

static int failed(const char *operation, int status) {
  (void)snprintf(cl_failure, sizeof cl_failure, "%s failed (OpenCL status %d)", operation, status);
  return -1;
}

static void *load_driver(void) {
#if defined(OSPREY_GPU_TESTING)
  if (test_driver) return (void *)(uintptr_t)1;
#endif
#if defined(_WIN32)
  return (void *)LoadLibraryA("OpenCL.dll");
#elif defined(__APPLE__)
  return dlopen("/System/Library/Frameworks/OpenCL.framework/OpenCL", RTLD_NOW | RTLD_LOCAL);
#else
  return dlopen("libOpenCL.so.1", RTLD_NOW | RTLD_LOCAL);
#endif
}

#if defined(OSPREY_GPU_TESTING)
static int test_symbol(const char *name, void *slot, size_t size) {
#define MATCH(field, text) if (strcmp(name, text) == 0 && test_driver->field && size == sizeof test_driver->field) { \
    memcpy(slot, &test_driver->field, size); return 0; }
  CL_FUNCTIONS(MATCH)
#undef MATCH
  return failed(name, -1);
}
#endif

static int symbol(void *library, const char *name, void *slot, size_t size) {
#if defined(OSPREY_GPU_TESTING)
  if (test_driver) return test_symbol(name, slot, size);
#endif
#if defined(_WIN32)
  FARPROC address = GetProcAddress((HMODULE)library, name);
#else
  void *address = dlsym(library, name);
#endif
  if (!address || size != sizeof address) return failed(name, -1);
  memcpy(slot, &address, size);
  return 0;
}

#define LOAD(field, name) if (symbol(library, name, &cl.field, sizeof cl.field)) return -1;
static int load_api(void) {
  void *library = load_driver();
  if (!library) return failed("OpenCL driver unavailable", -1);
  CL_FUNCTIONS(LOAD)
  return 0;
}
#undef LOAD

static void release_context(void) {
  if (cl_kernel) (void)cl.release_kernel(cl_kernel);
  if (cl_program) (void)cl.release_program(cl_program);
  if (cl_queue) (void)cl.release_queue(cl_queue);
  if (cl_context) (void)cl.release_context(cl_context);
  cl_kernel = cl_program = cl_queue = cl_context = NULL;
}

static int find_device(void) {
  unsigned count = 0;
  int status = cl.platforms(0, NULL, &count);
  if (status || !count) return failed("GPU platform discovery", status);
  void **platforms = calloc(count, sizeof *platforms);
  if (!platforms) return failed("platform allocation", -1);
  status = cl.platforms(count, platforms, NULL);
  for (unsigned i = 0; !status && i < count; i++) {
    if (cl.devices(platforms[i], UINT64_C(4), 1, &cl_device, NULL) == 0 && cl_device) break;
  }
  free(platforms);
  if (status || !cl_device) return failed("compute GPU discovery", status);
  char name[192] = {0};
  status = cl.info(cl_device, 0x102B, sizeof name, name, NULL);
  if (status) return failed("device name", status);
  (void)snprintf(cl_name, sizeof cl_name, "opencl:0 (%s; matrix kernels)", name);
  return 0;
}

static int create_context(void) {
  int status = 0;
  cl_context = cl.context(NULL, 1, &cl_device, NULL, NULL, &status);
  if (status || !cl_context) return failed("context creation", status);
  cl_queue = cl.queue(cl_context, cl_device, 0, &status);
  return status || !cl_queue ? failed("queue creation", status) : 0;
}

static void build_failure(int status) {
  (void)failed("kernel compilation", status);
  char log[900] = {0};
  if (cl.build_info(cl_program, cl_device, 0x1183, sizeof log, log, NULL) == 0 && log[0])
    (void)snprintf(cl_failure, sizeof cl_failure, "kernel compilation (status %d): %.899s", status, log);
}

static int create_program(void) {
  int status = 0;
  const char *source = (const char *)osprey_gpu_shader;
  size_t size = osprey_gpu_shader_len;
  cl_program = cl.program(cl_context, 1, &source, &size, &status);
  if (status || !cl_program) return failed("program creation", status);
  status = cl.build(cl_program, 1, &cl_device, "-cl-std=CL1.2", NULL, NULL);
  if (status) {
    build_failure(status);
    return -1;
  }
  cl_kernel = cl.kernel(cl_program, "osprey_matmul64", &status);
  return status || !cl_kernel ? failed("kernel creation", status) : 0;
}

static int initialize(void) {
  if (cl_ready) return cl_ready;
  if ((!cl.platforms && load_api()) || find_device() || create_context() || create_program()) {
    release_context();
    cl_ready = -1;
  } else {
    cl_ready = 1;
  }
  return cl_ready;
}

const char *osprey_gpu_opencl_name(void) {
  (void)pthread_mutex_lock(&cl_lock);
  int ready = initialize();
  (void)pthread_mutex_unlock(&cl_lock);
  return ready > 0 ? cl_name : "opencl:unavailable";
}

uint64_t osprey_gpu_opencl_completed(void) {
  (void)pthread_mutex_lock(&cl_lock);
  uint64_t count = cl_completed;
  (void)pthread_mutex_unlock(&cl_lock);
  return count;
}

static int set_arguments(void *const buffers[3], const uint64_t dimensions[3]) {
  for (unsigned i = 0; i < 6; i++) {
    size_t size = i < 3 ? sizeof buffers[i] : sizeof dimensions[i - 3];
    const void *value = i < 3 ? (const void *)&buffers[i] : (const void *)&dimensions[i - 3];
    int status = cl.arg(cl_kernel, i, size, value);
    if (status) return failed("kernel argument", status);
  }
  return 0;
}

static int dispatch(void *const buffers[3], const uint64_t dimensions[3], int64_t *out) {
  if (set_arguments(buffers, dimensions)) return -1;
  size_t count = (size_t)(dimensions[0] * dimensions[2]);
  int status = cl.launch(cl_queue, cl_kernel, 1, NULL, &count, NULL, 0, NULL, NULL);
  if (status) return failed("GPU dispatch", status);
  status = cl.read(cl_queue, buffers[2], 1, 0, count * 8, out, 0, NULL, NULL);
  if (status) return failed("GPU completion/readback", status);
  cl_completed++;
  if (cl_completed == 1 || getenv("OSPREY_GPU_TRACE"))
    (void)fprintf(stderr, "[gpu] %s: completed dispatch %llu\n", cl_name, (unsigned long long)cl_completed);
  return 1;
}

static int product(const int64_t *left, const int64_t *right, int64_t *out, const uint64_t dims[3]) {
  int status = 0;
  void *buffers[3] = {NULL, NULL, NULL};
  const void *sources[3] = {left, right, NULL};
  uint64_t sizes[3] = {dims[0] * dims[1] * 8, dims[1] * dims[2] * 8, dims[0] * dims[2] * 8};
  for (unsigned i = 0; i < 3; i++) {
    buffers[i] = cl.buffer(cl_context, i < 2 ? UINT64_C(36) : UINT64_C(2), (size_t)sizes[i], sources[i], &status);
    if (status || !buffers[i]) { status = failed("GPU buffer allocation", status); break; }
  }
  if (!status) status = dispatch(buffers, dims, out);
  for (unsigned i = 0; i < 3; i++) if (buffers[i]) (void)cl.release_buffer(buffers[i]);
  return status;
}

int osprey_gpu_opencl_matmul(const int64_t *left, const int64_t *right, int64_t *out,
                             int64_t rows, int64_t inner, int64_t cols) {
  (void)pthread_mutex_lock(&cl_lock);
  uint64_t dimensions[3] = {(uint64_t)rows, (uint64_t)inner, (uint64_t)cols};
  int status = initialize() < 0 ? -1 : (!rows || !inner || !cols ? 1 : product(left, right, out, dimensions));
  if (status < 0) osp_io_error_text("gpuMatMul", cl_name[0] ? cl_name : "OpenCL", cl_failure);
  (void)pthread_mutex_unlock(&cl_lock);
  return status;
}

#if defined(OSPREY_GPU_TESTING)
void osprey_gpu_opencl_test_driver(const OspreyOpenCLApi *api) {
  release_context();
  test_driver = api;
  memset(&cl, 0, sizeof cl);
  cl_device = NULL;
  cl_ready = 0;
  cl_completed = 0;
}
#endif
