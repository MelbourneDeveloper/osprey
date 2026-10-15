// Device selection and the compiled-kernel launch contract against mock
// adapters [GPU-DEVICE] [GPU-BACKEND-DEVICE] [GPU-DEVICE-FAULT]: the host
// oracle runs when the host is selected, a selected device never falls back
// to it, and a fault or a kernel without a device form aborts by name.
#include "gpu_device.h"
#include "gpu_opencl_api.h"
#include "io_error.h"
#include "test_death.h"
#include <assert.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/wait.h>
#include <unistd.h>

void *osprey_gpu_alloc(int64_t length);
int64_t *osprey_gpu_data(void *buffer);
int32_t osprey_gpu_launch(const char *kernel, const char *source, const char *kinds, const int64_t *uniforms,
                          void *in0, void *in1, int64_t *out, int64_t out_len, int64_t count, int32_t sequential);
int32_t osprey_gpu_declined(const char *what);
int osprey_gpu_buffer_tests(int argc, char **argv);

static int device_error;
static unsigned dispatches, last_args;
static uint64_t last_threads;

#if defined(OSPREY_GPU_METAL) || defined(OSPREY_GPU_OPENCL)
static int64_t injected_fault;
// The mock device "runs" a kernel by stamping word 0 of the result and
// recording the injected fault, if any, in the fault word (argument 3).
static int mock_dispatch(const OspreyGpuDispatch *d, int64_t mark) {
  assert(d->kernel && d->count >= 1 && d->threads >= 1);
  dispatches++;
  last_args = d->count;
  last_threads = d->threads;
  for (unsigned i = 0; i < d->count; i++) {
    if (d->args[i].kind == OSPREY_GPU_OUTPUT && d->args[i].count) d->args[i].out[0] = i == 3 ? injected_fault : mark;
  }
  return device_error ? -1 : 1;
}
#endif
#if defined(OSPREY_GPU_METAL)
const char *osprey_gpu_metal_name(void) { return "metal:mock"; }
int osprey_gpu_metal_dispatch(const OspreyGpuDispatch *d) { return mock_dispatch(d, 42); }
#endif
#if defined(OSPREY_GPU_OPENCL)
const char *osprey_gpu_opencl_name(void) { return "opencl:mock"; }
int osprey_gpu_opencl_dispatch(const OspreyGpuDispatch *d) { return mock_dispatch(d, 43); }
#endif

static void select_mode(const char *mode) {
  if (mode) assert(setenv("OSPREY_GPU", mode, 1) == 0);
  else assert(unsetenv("OSPREY_GPU") == 0);
}

static int32_t launch(int64_t *out, const char *kinds, const int64_t *uniforms) {
  void *in = osprey_gpu_alloc(3);
  return osprey_gpu_launch("osprey_kernel_1", "kernel void osprey_kernel_1() {}", kinds, uniforms, in, NULL, out, 1, 3, 0);
}

static void host_case(const char *mode) {
  select_mode(mode);
  int64_t a = 1, out = 99;
  assert(strcmp(osprey_gpu_device_name(), "host") == 0);
  assert(osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1) == 0 && out == 99);
  assert(launch(&out, "", NULL) == 0 && out == 99);
  assert(osprey_gpu_declined("gpuScan") == 0);
}

static void device_case(const char *mode, int available, const char *name) {
  select_mode(mode);
  int64_t a = 1, out = 99;
  osp_io_error_clear();
  assert(strcmp(osprey_gpu_device_name(), name) == 0);
  int status = osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1);
  assert(available ? status == 1 && last_args == 6 && last_threads == 1 : status == -1);
  if (!available) {
    assert(strstr(osp_io_error(), "no GPU adapter") != NULL);
    return;
  }
  device_error = 1;
  assert(osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1) == -1);
  device_error = 0;
  unsigned before = dispatches;
  assert(osprey_gpu_device_matmul(&a, &a, &out, 0, 1, 1) == 1 && dispatches == before);
  void *table = osprey_gpu_alloc(4);
  int64_t uniforms[2] = {7, (int64_t)(intptr_t)table};
  assert(launch(&out, "sb", uniforms) == 1 && out == (name[0] == 'm' ? 42 : 43));
  assert(last_args == 5 + 1 + 2 && last_threads == 3);
}

static void invalid_case(void) {
  select_mode("invalid");
  int64_t a = 1, out = 99;
  assert(strcmp(osprey_gpu_device_name(), "invalid:device") == 0);
  assert(osprey_gpu_device_matmul(&a, &a, &out, 1, 1, 1) == -1 && out == 99);
  assert(strstr(osp_io_error(), "unknown backend") != NULL);
}

// Runs `scenario` in a child that must abort with `expected` on stderr.
static void aborts(void (*scenario)(void), const char *expected) {
  int pipes[2];
  assert(pipe(pipes) == 0);
  pid_t child = fork();
  assert(child >= 0);
  if (!child) {
    // The guard's own lines count: dump the coverage counters before the abort.
    (void)signal(SIGABRT, osp_death_dump_and_reraise);
    assert(dup2(pipes[1], 2) >= 0);
    scenario();
    _exit(0);
  }
  close(pipes[1]);
  char text[1024] = {0};
  ssize_t n = read(pipes[0], text, sizeof text - 1);
  close(pipes[0]);
  int status = 0;
  assert(waitpid(child, &status, 0) == child && n > 0);
  assert(WIFSIGNALED(status) && strstr(text, expected) != NULL);
}

static void declined_on_device(void) {
  select_mode("auto");
  (void)osprey_gpu_declined("gpuScan over a kernel the device cannot run");
}
#if defined(OSPREY_GPU_METAL)
static void faulting_launch(void) {
  select_mode("auto");
  injected_fault = OSPREY_GPU_FAULT_OVERFLOW;
  int64_t out = 0;
  (void)launch(&out, "", NULL);
}
static void too_many_uniforms(void) {
  select_mode("auto");
  int64_t out = 0, uniforms[OSPREY_GPU_MAX_ARGS] = {0};
  char kinds[OSPREY_GPU_MAX_ARGS + 1];
  memset(kinds, 's', OSPREY_GPU_MAX_ARGS);
  kinds[OSPREY_GPU_MAX_ARGS] = '\0';
  (void)launch(&out, kinds, uniforms);
}
static void device_failure_on_launch(void) {
  select_mode("auto");
  device_error = 1;
  int64_t out = 0;
  (void)launch(&out, "", NULL);
}
#endif

int main(void) {
  host_case(NULL); host_case(""); host_case("host");
  assert(osprey_gpu_buffer_tests(0, NULL) == 0);
#if defined(OSPREY_GPU_METAL)
  device_case("auto", 1, "metal:mock"); device_case("metal", 1, "metal:mock");
  aborts(faulting_launch, "kernel osprey_kernel_1: integer overflow");
  aborts(declined_on_device, "gpuScan over a kernel the device cannot run; a selected device cannot run it");
  aborts(too_many_uniforms, "too many kernel arguments");
  aborts(device_failure_on_launch, "kernel osprey_kernel_1:");
#else
  device_case("auto", 0, "unavailable:device"); device_case("metal", 0, "unavailable:device");
  aborts(declined_on_device, "a selected device cannot run it");
#endif
#if defined(OSPREY_GPU_OPENCL)
  device_case("opencl", 1, "opencl:mock");
#else
  device_case("opencl", 0, "unavailable:device");
#endif
  invalid_case();
  assert(unsetenv("OSPREY_GPU_TRACE") == 0 && osprey_gpu_traced(1) && !osprey_gpu_traced(2));
  assert(setenv("OSPREY_GPU_TRACE", "0", 1) == 0 && !osprey_gpu_traced(2));
  assert(setenv("OSPREY_GPU_TRACE", "1", 1) == 0 && osprey_gpu_traced(2));
  puts("[ok] GPU selection: host, auto, Metal, OpenCL, missing adapter, compiled-kernel launch, faults and declines abort by name, trace switch and no silent CPU fallback");
  return 0;
}
