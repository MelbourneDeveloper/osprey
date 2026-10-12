// Assertion-driven tests for gpu_runtime.c and gpu_host.c — the
// [GPU-BACKEND-HOST] staging buffer behind toGpu/fromGpu/gpuMap/gpuFold
// ([GPU-BUFFER], [GPU-BUFFER-ELEM], [GPU-BUFFER-LENGTH], [GPU-GET],
// [GPU-FILTER]) and the host operations a trained model needs: the dense
// product [GPU-MATMUL] and the two weight loaders [GPU-BUFFER-LOAD],
// [GPU-BUFFER-TENSOR] (docs/specs/0034). Linked with gpu_host.c,
// file_runtime.c (the io_error channel) and memory_runtime.c (the default
// backend) by the Makefile's _test_c_runtime.
//
// The load-bearing contract is TOTALITY: a negative length, an overflowing
// length, and a failed data allocation must all yield an observably EMPTY
// buffer — never a trap, never a partial buffer, never an out-of-bounds
// write. Every accessor must be NULL-safe and bounds-checked, because these
// are the no-trap backstops behind codegen's counted loops. The host
// operations extend it: every failure is NULL plus a recorded reason, never a
// partial result.
#ifndef _GNU_SOURCE
#define _GNU_SOURCE
#endif
#include <assert.h>
#include <errno.h>
#include <math.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

void *osprey_gpu_alloc(int64_t length);
int64_t osprey_gpu_len(void *buffer);
int64_t osprey_gpu_get(void *buffer, int64_t index);
void osprey_gpu_set(void *buffer, int64_t index, int64_t word);
int32_t osprey_gpu_in_bounds(void *buffer, int64_t index);
void osprey_gpu_take(void *buffer, int64_t count);
void *osprey_gpu_matmul(void *left, void *right, int64_t rows, int64_t inner, int64_t cols);
void *osprey_gpu_load_f32(const char *path, int64_t offset, int64_t count);
void *osprey_gpu_load_tensor(const char *path, const char *name);
const char *osp_io_error(void);
void osp_io_error_clear(void);

static long g_checks = 0;
#define CHECK(c)                                                               \
  do {                                                                         \
    g_checks++;                                                                \
    assert(c);                                                                 \
  } while (0)

// Elements are 8-byte words; this mirrors gpu_runtime.c's GPU_MAX_LEN.
#define GPU_MAX_LEN (INT64_MAX / 8)
#define ROUND_LEN ((int64_t)4096)
#define SMALL_LEN ((int64_t)16)

// A fresh buffer reports its exact length and reads zero at every index.
static void t_alloc_zero_filled(void) {
  void *b = osprey_gpu_alloc(SMALL_LEN);
  CHECK(b != NULL);
  CHECK(osprey_gpu_len(b) == SMALL_LEN);
  for (int64_t i = 0; i < SMALL_LEN; i++) {
    CHECK(osprey_gpu_get(b, i) == 0);
    CHECK(osprey_gpu_in_bounds(b, i) == 1);
  }
  CHECK(osprey_gpu_in_bounds(b, SMALL_LEN) == 0);
  CHECK(osprey_gpu_in_bounds(b, -1) == 0);
}

// Out-of-range lengths clamp to the EMPTY buffer — allocation is total.
static void t_alloc_clamps_to_empty(void) {
  int64_t bad[] = {-1, INT64_MIN, GPU_MAX_LEN + 1, INT64_MAX};
  for (unsigned i = 0; i < sizeof(bad) / sizeof(bad[0]); i++) {
    void *b = osprey_gpu_alloc(bad[i]);
    CHECK(b != NULL);
    CHECK(osprey_gpu_len(b) == 0);
    CHECK(osprey_gpu_in_bounds(b, 0) == 0);
    CHECK(osprey_gpu_get(b, 0) == 0);
    osprey_gpu_set(b, 0, 99); // ignored, not a wild store
    CHECK(osprey_gpu_get(b, 0) == 0);
  }
  void *zero = osprey_gpu_alloc(0);
  CHECK(zero != NULL && osprey_gpu_len(zero) == 0);
}

// An IN-range length whose data allocation fails (no machine has 2^63-aligned
// bytes) still yields the empty buffer, never NULL-deref or a partial length.
static void t_alloc_failure_is_empty(void) {
  void *b = osprey_gpu_alloc(GPU_MAX_LEN);
  CHECK(b != NULL);
  CHECK(osprey_gpu_len(b) == 0);
  CHECK(osprey_gpu_in_bounds(b, 0) == 0);
  osprey_gpu_set(b, 0, 7);
  CHECK(osprey_gpu_get(b, 0) == 0);
}

// Every accessor tolerates NULL — the empty-buffer degenerate case.
static void t_null_safety(void) {
  CHECK(osprey_gpu_len(NULL) == 0);
  CHECK(osprey_gpu_get(NULL, 0) == 0);
  CHECK(osprey_gpu_in_bounds(NULL, 0) == 0);
  osprey_gpu_set(NULL, 0, 1);
  osprey_gpu_take(NULL, 3);
  CHECK(osprey_gpu_len(NULL) == 0);
}

// Full write/read round trip, and out-of-bounds stores land NOWHERE: after
// poking every invalid index the whole payload is still byte-identical.
static void t_set_get_roundtrip(void) {
  void *b = osprey_gpu_alloc(ROUND_LEN);
  CHECK(b != NULL && osprey_gpu_len(b) == ROUND_LEN);
  for (int64_t i = 0; i < ROUND_LEN; i++) {
    osprey_gpu_set(b, i, i * i - 3);
  }
  osprey_gpu_set(b, ROUND_LEN, 0x7A7A);     // one past the end
  osprey_gpu_set(b, -1, 0x7A7A);            // below zero
  osprey_gpu_set(b, INT64_MAX, 0x7A7A);     // absurdly far
  osprey_gpu_set(b, INT64_MIN, 0x7A7A);
  for (int64_t i = 0; i < ROUND_LEN; i++) {
    CHECK(osprey_gpu_get(b, i) == i * i - 3);
  }
}

// Extreme word values survive storage exactly (raw 8-byte words, no coercion).
static void t_boundary_words(void) {
  void *b = osprey_gpu_alloc(2);
  CHECK(b != NULL);
  osprey_gpu_set(b, 0, INT64_MIN);
  osprey_gpu_set(b, 1, INT64_MAX);
  CHECK(osprey_gpu_get(b, 0) == INT64_MIN);
  CHECK(osprey_gpu_get(b, 1) == INT64_MAX);
}

// take publishes a compaction prefix: it only ever SHRINKS, the surviving
// prefix is untouched, and the amputated tail becomes unaddressable.
static void t_take_shrinks_only(void) {
  void *b = osprey_gpu_alloc(10);
  CHECK(b != NULL);
  for (int64_t i = 0; i < 10; i++) {
    osprey_gpu_set(b, i, 100 + i);
  }
  osprey_gpu_take(b, 20); // larger than length: no growth
  CHECK(osprey_gpu_len(b) == 10);
  osprey_gpu_take(b, 7);
  CHECK(osprey_gpu_len(b) == 7);
  for (int64_t i = 0; i < 7; i++) {
    CHECK(osprey_gpu_get(b, i) == 100 + i);
  }
  CHECK(osprey_gpu_in_bounds(b, 7) == 0);
  CHECK(osprey_gpu_get(b, 7) == 0);
  osprey_gpu_set(b, 7, 55); // past the published length: ignored
  osprey_gpu_take(b, 7);    // idempotent at the same count
  CHECK(osprey_gpu_len(b) == 7);
  osprey_gpu_take(b, -3);   // negative clamps to zero: fully compacted away
  CHECK(osprey_gpu_len(b) == 0);
  CHECK(osprey_gpu_in_bounds(b, 0) == 0);
}

// --- float words and the io_error channel --------------------------------------

// The double a buffer word holds [GPU-BUFFER-ELEM].
static double at(void *b, int64_t i) {
  int64_t word = osprey_gpu_get(b, i);
  double d;
  memcpy(&d, &word, sizeof d);
  return d;
}

// A float buffer holding `xs` verbatim.
static void *floats(const double *xs, int64_t n) {
  void *b = osprey_gpu_alloc(n);
  CHECK(b != NULL && osprey_gpu_len(b) == n);
  for (int64_t i = 0; i < n; i++) {
    int64_t word;
    memcpy(&word, &xs[i], sizeof word);
    osprey_gpu_set(b, i, word);
  }
  return b;
}

// Every failure records exactly one sentence, "op: subject: reason".
static void expect_text(const char *message) {
  const char *actual = osp_io_error();
  if (!actual || strcmp(actual, message) != 0) {
    fprintf(stderr, "expected: %s\n  actual: %s\n", message, actual ? actual : "(none)");
  }
  CHECK(actual != NULL && strcmp(actual, message) == 0);
}

// An errno failure records "op: subject: " followed by the OS's own wording.
static void expect_errno(const char *prefix, int err) {
  const char *message = osp_io_error();
  CHECK(message != NULL);
  CHECK(strncmp(message, prefix, strlen(prefix)) == 0);
  CHECK(strstr(message, strerror(err)) != NULL);
}

// --- the dense product ---------------------------------------------------------

// A 2×3 · 3×2 product lands exactly, row-major, at rows × cols.
static void t_matmul_exact(void) {
  double l[] = {1, 2, 3, 4, 5, 6};
  double r[] = {7, 8, 9, 10, 11, 12};
  double expected[] = {58, 64, 139, 154};
  osp_io_error_clear();
  void *out = osprey_gpu_matmul(floats(l, 6), floats(r, 6), 2, 3, 2);
  CHECK(out != NULL && osp_io_error() == NULL);
  CHECK(osprey_gpu_len(out) == 4);
  for (int64_t i = 0; i < 4; i++) {
    CHECK(at(out, i) == expected[i]);
  }
}

// The sum runs in k order from 0.0: (2^53 + 1) absorbs the 1 before -2^53
// cancels, so the result is 0.0; the mirrored row gives 1.0. A backend that
// reassociates (pairing ±2^53 first) would answer 1.0 for both.
static void t_matmul_accumulates_in_k_order(void) {
  double big = 9007199254740992.0;
  double forward[] = {big, 1.0, -big};
  double mirrored[] = {-big, 1.0, big};
  double ones[] = {1.0, 1.0, 1.0};
  void *a = osprey_gpu_matmul(floats(forward, 3), floats(ones, 3), 1, 3, 1);
  void *b = osprey_gpu_matmul(floats(mirrored, 3), floats(ones, 3), 1, 3, 1);
  CHECK(a != NULL && osprey_gpu_len(a) == 1 && at(a, 0) == 0.0);
  CHECK(b != NULL && osprey_gpu_len(b) == 1 && at(b, 0) == 1.0);
}

// Any zero dimension with consistent lengths is the zero-filled rows × cols
// result; a zero inner dimension means empty (even NULL) operands.
static void t_matmul_zero_dimensions(void) {
  double two[] = {1.0, 2.0};
  osp_io_error_clear();
  void *inner0 = osprey_gpu_matmul(NULL, NULL, 2, 0, 2);
  CHECK(inner0 != NULL && osprey_gpu_len(inner0) == 4 && osp_io_error() == NULL);
  for (int64_t i = 0; i < 4; i++) {
    CHECK(at(inner0, i) == 0.0);
  }
  void *rows0 = osprey_gpu_matmul(osprey_gpu_alloc(0), floats(two, 2), 0, 2, 1);
  CHECK(rows0 != NULL && osprey_gpu_len(rows0) == 0);
  void *cols0 = osprey_gpu_matmul(floats(two, 2), osprey_gpu_alloc(0), 1, 2, 0);
  CHECK(cols0 != NULL && osprey_gpu_len(cols0) == 0);
}

// A shape the operands do not have is NULL plus the reason naming the operand;
// a negative or overflowing result shape is refused before any operand is read.
static void t_matmul_shape_errors(void) {
  double six[] = {1, 2, 3, 4, 5, 6};
  void *l = floats(six, 6);
  void *r = floats(six, 6);
  osp_io_error_clear();
  CHECK(osprey_gpu_matmul(l, r, -1, 3, 2) == NULL);
  expect_text("gpuMatMul: result: rows * cols is negative or overflows");
  osp_io_error_clear();
  CHECK(osprey_gpu_matmul(l, r, GPU_MAX_LEN, 1, 2) == NULL);
  expect_text("gpuMatMul: result: rows * cols is negative or overflows");
  osp_io_error_clear();
  CHECK(osprey_gpu_matmul(l, r, 2, 4, 2) == NULL);
  expect_text("gpuMatMul: left operand: has 6 elements, the stated shape needs 8");
  osp_io_error_clear();
  CHECK(osprey_gpu_matmul(l, r, 2, 3, 3) == NULL);
  expect_text("gpuMatMul: right operand: has 6 elements, the stated shape needs 9");
  osp_io_error_clear();
  CHECK(osprey_gpu_matmul(NULL, r, 1, 1, 6) == NULL);
  expect_text("gpuMatMul: left operand: has 0 elements, the stated shape needs 1");
  osp_io_error_clear();
  CHECK(osprey_gpu_matmul(l, r, 2, -3, 2) == NULL);
  CHECK(strncmp(osp_io_error(), "gpuMatMul: left operand:", 24) == 0);
}

// --- files on disk ---------------------------------------------------------------

static void put_u16(unsigned char *at, uint16_t v) {
  at[0] = (unsigned char)(v & 0xff);
  at[1] = (unsigned char)(v >> 8);
}

static void put_u32(unsigned char *at, uint32_t v) {
  put_u16(at, (uint16_t)(v & 0xffff));
  put_u16(at + 2, (uint16_t)(v >> 16));
}

static void put_u64(unsigned char *at, uint64_t v) {
  put_u32(at, (uint32_t)(v & 0xffffffffu));
  put_u32(at + 4, (uint32_t)(v >> 32));
}

static void put_f32(unsigned char *at, float f) {
  uint32_t bits;
  memcpy(&bits, &f, sizeof bits);
  put_u32(at, bits);
}

static void put_f64(unsigned char *at, double d) {
  uint64_t bits;
  memcpy(&bits, &d, sizeof bits);
  put_u64(at, bits);
}

// `bytes` written to a fresh temporary file whose path is returned (static).
static const char *temp_file(const unsigned char *bytes, size_t n) {
  static char path[64];
  (void)snprintf(path, sizeof path, "/tmp/osprey_gpu_host_XXXXXX");
  int fd = mkstemp(path);
  CHECK(fd >= 0);
  FILE *f = fdopen(fd, "wb");
  CHECK(f != NULL);
  CHECK(fwrite(bytes, 1, n, f) == n);
  CHECK(fclose(f) == 0);
  return path;
}

// Raw binary32 words widen exactly — sign of zero, infinities, NaN and a
// subnormal included — from a byte offset, a short read is an error and count
// zero is the empty buffer.
static void t_load_f32(void) {
  float stored[] = {1.5f, -0.0f, INFINITY, NAN, 1e-45f, -2.5e-3f};
  unsigned char bytes[4 + sizeof stored];
  memset(bytes, 0x5a, 4);
  for (size_t i = 0; i < 6; i++) {
    put_f32(bytes + 4 + i * 4, stored[i]);
  }
  const char *path = temp_file(bytes, sizeof bytes);
  osp_io_error_clear();
  void *b = osprey_gpu_load_f32(path, 4, 6);
  CHECK(b != NULL && osprey_gpu_len(b) == 6 && osp_io_error() == NULL);
  CHECK(at(b, 0) == 1.5 && at(b, 1) == 0.0 && signbit(at(b, 1)));
  CHECK(isinf(at(b, 2)) && isnan(at(b, 3)));
  CHECK(at(b, 4) == (double)1e-45f && at(b, 5) == (double)-2.5e-3f);
  void *none = osprey_gpu_load_f32(path, 4, 0);
  CHECK(none != NULL && osprey_gpu_len(none) == 0);
  char prefix[96];
  (void)snprintf(prefix, sizeof prefix, "gpuLoad: %s: ", path);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_f32(path, 4, 7) == NULL);
  CHECK(strncmp(osp_io_error(), prefix, strlen(prefix)) == 0);
  CHECK(strstr(osp_io_error(), "file ends before the tensor does") != NULL);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_f32(path, 0, -1) == NULL);
  CHECK(strstr(osp_io_error(), "a negative element count") != NULL);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_f32(path, -1, 1) == NULL);
  expect_errno(prefix, EINVAL);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_f32("/nonexistent/osprey/weights.bin", 0, 1) == NULL);
  expect_errno("gpuLoad: /nonexistent/osprey/weights.bin: ", ENOENT);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_f32(NULL, 0, 1) == NULL);
  expect_errno("gpuLoad: ", EINVAL);
  (void)unlink(path);
}

// The safetensors header this suite reads: `__metadata__` LAST, a name that
// is a prefix of another (`h.1.w` inside `h.11.w`), every admitted dtype, a
// refused one, and three malformed entries. Data offsets are relative to the
// data section. Returns the complete file image's length.
static size_t safetensors_image(unsigned char *out, size_t cap) {
  const char *header =
      "{\"a\":{\"dtype\":\"F32\",\"shape\":[2,3],\"data_offsets\":[0,24]},"
      "\"h.1.w\":{\"dtype\":\"F16\",\"shape\":[5],\"data_offsets\":[24,34]},"
      "\"h.11.w\":{\"dtype\":\"BF16\",\"shape\":[2],\"data_offsets\":[34,38]},"
      "\"d\":{\"dtype\":\"F64\",\"shape\":[2],\"data_offsets\":[38,54]},"
      "\"i\":{\"dtype\":\"I64\",\"shape\":[1],\"data_offsets\":[54,62]},"
      "\"inv\":{\"dtype\":\"F32\",\"shape\":[1],\"data_offsets\":[8,4]},"
      "\"odd\":{\"dtype\":\"F32\",\"shape\":[1],\"data_offsets\":[0,6]},"
      "\"beyond\":{\"dtype\":\"F32\",\"shape\":[250],\"data_offsets\":[0,1000]},"
      "\"__metadata__\":{\"format\":\"pt\"}}";
  size_t n = strlen(header);
  CHECK(8 + n + 62 <= cap);
  put_u64(out, (uint64_t)n);
  memcpy(out + 8, header, n);
  unsigned char *data = out + 8 + n;
  float a[] = {0.5f, -1.0f, 2.0f, 3.0f, -4.5f, 6.0f};
  for (size_t i = 0; i < 6; i++) {
    put_f32(data + i * 4, a[i]);
  }
  uint16_t f16[] = {0x3C00, 0xC100, 0x0001, 0x7C00, 0x7E00}; // 1, -2.5, 2^-24, inf, nan
  for (size_t i = 0; i < 5; i++) {
    put_u16(data + 24 + i * 2, f16[i]);
  }
  put_u16(data + 34, 0x3F80); // bf16 1.0
  put_u16(data + 36, 0xC040); // bf16 -3.0
  put_f64(data + 38, 3.25);
  put_f64(data + 46, -1e300);
  put_u64(data + 54, 7);
  return 8 + n + 62;
}

// Each admitted dtype widens exactly, the scanner matches whole quoted names
// regardless of entry order, and every malformed entry is a named error.
static void t_load_tensor(void) {
  unsigned char image[1024];
  const char *path = temp_file(image, safetensors_image(image, sizeof image));
  osp_io_error_clear();
  void *a = osprey_gpu_load_tensor(path, "a");
  CHECK(a != NULL && osprey_gpu_len(a) == 6 && osp_io_error() == NULL);
  CHECK(at(a, 0) == 0.5 && at(a, 1) == -1.0 && at(a, 4) == -4.5 && at(a, 5) == 6.0);
  void *h1 = osprey_gpu_load_tensor(path, "h.1.w");
  CHECK(h1 != NULL && osprey_gpu_len(h1) == 5);
  CHECK(at(h1, 0) == 1.0 && at(h1, 1) == -2.5 && at(h1, 2) == ldexp(1.0, -24));
  CHECK(isinf(at(h1, 3)) && at(h1, 3) > 0 && isnan(at(h1, 4)));
  void *h11 = osprey_gpu_load_tensor(path, "h.11.w");
  CHECK(h11 != NULL && osprey_gpu_len(h11) == 2);
  CHECK(at(h11, 0) == 1.0 && at(h11, 1) == -3.0);
  void *d = osprey_gpu_load_tensor(path, "d");
  CHECK(d != NULL && osprey_gpu_len(d) == 2 && at(d, 0) == 3.25 && at(d, 1) == -1e300);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, "i") == NULL);
  expect_text("gpuLoadTensor: i: dtype is not F64, F32, F16 or BF16");
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, "__metadata__") == NULL);
  expect_text("gpuLoadTensor: __metadata__: dtype is not F64, F32, F16 or BF16");
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, "h.1") == NULL);
  expect_text("gpuLoadTensor: h.1: no tensor of that name in the header");
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, "inv") == NULL);
  expect_text("gpuLoadTensor: inv: data_offsets are missing or inverted");
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, "odd") == NULL);
  expect_text("gpuLoadTensor: odd: byte span is not a whole number of elements");
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, "beyond") == NULL);
  CHECK(strstr(osp_io_error(), "file ends before the tensor does") != NULL);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(path, NULL) == NULL);
  expect_text("gpuLoadTensor: <null path>: no tensor of that name in the header");
  (void)unlink(path);
}

// A file that is not a safetensors file is refused at the first byte that
// proves it: too short for the prefix, an implausible header length, or a
// header the file ends inside of.
static void t_load_tensor_malformed_files(void) {
  unsigned char stub[16] = {0};
  const char *tiny = temp_file(stub, 4);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(tiny, "a") == NULL);
  CHECK(strstr(osp_io_error(), "shorter than a safetensors header prefix") != NULL);
  (void)unlink(tiny);
  put_u64(stub, (uint64_t)200 * 1024 * 1024);
  const char *huge = temp_file(stub, 16);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(huge, "a") == NULL);
  CHECK(strstr(osp_io_error(), "header length is not plausible for safetensors") != NULL);
  (void)unlink(huge);
  put_u64(stub, 100);
  const char *cut = temp_file(stub, 16);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor(cut, "a") == NULL);
  CHECK(strstr(osp_io_error(), "file ends inside the header") != NULL);
  (void)unlink(cut);
  osp_io_error_clear();
  CHECK(osprey_gpu_load_tensor("/nonexistent/osprey/model.safetensors", "a") == NULL);
  expect_errno("gpuLoadTensor: /nonexistent/osprey/model.safetensors: ", ENOENT);
}

int main(void) {
  t_alloc_zero_filled();
  t_alloc_clamps_to_empty();
  t_alloc_failure_is_empty();
  t_null_safety();
  t_set_get_roundtrip();
  t_boundary_words();
  t_take_shrinks_only();
  t_matmul_exact();
  t_matmul_accumulates_in_k_order();
  t_matmul_zero_dimensions();
  t_matmul_shape_errors();
  t_load_f32();
  t_load_tensor();
  t_load_tensor_malformed_files();
  printf("[ok] gpu_runtime: %ld assertions\n", g_checks);
  return 0;
}
