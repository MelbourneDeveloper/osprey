// Host-backend operations on dense buffers beyond the buffer primitive itself
// (docs/specs/0034-GPUComputation.md): the row-major dense product behind
// `gpuMatMul` [GPU-MATMUL], and the two ways a trained model's weights enter a
// program without passing through a List — raw little-endian binary32 words
// at a byte offset (`gpuLoad`, [GPU-BUFFER-LOAD]: the layout of llm.c's
// checkpoints and of any raw dump) and one named tensor out of a safetensors
// file (`gpuLoadTensor`, [GPU-BUFFER-TENSOR]: the Hugging Face interchange
// format — an 8-byte little-endian header length, a JSON header mapping each
// name to its dtype, shape and data_offsets, then the raw tensor bytes).
//
// Every stored float widens EXACTLY to the double word a GpuBuffer<float>
// holds. Every failure path records WHY on the io_error channel and returns
// NULL; codegen turns that pair into `Error { message }` [BUILTIN-FILE-ERRMSG].
// Like gpu_runtime.c this unit is backend-neutral: buffers come from
// osprey_gpu_alloc, scratch memory never outlives the call, and a buffer
// allocated for a result that then fails is released before the failure is
// reported, so a NULL result owes the memory backend nothing.

#include "io_error.h"
#include "memory_hooks.h"
#include "gpu_device.h"

#include <errno.h>
#include <limits.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// Seek to a 64-bit offset on every target. Windows spells it _fseeki64; WASI's
// libc declares only fseek(long), and a wasm32 long is 32 bits, so an offset
// beyond it is refused there rather than truncated into the wrong tensor.
#if defined(_WIN32)
#define osp_fseek _fseeki64
#elif defined(__wasi__)
static int osp_fseek(FILE *f, int64_t offset, int whence) {
  if (offset > (int64_t)LONG_MAX) {
    errno = EOVERFLOW;
    return -1;
  }
  return fseek(f, (long)offset, whence);
}
#else
#define osp_fseek fseeko
#endif

void *osprey_gpu_alloc(int64_t length);
int64_t osprey_gpu_len(void *buffer);
int64_t *osprey_gpu_data(void *buffer);

// Elements are 8-byte words; mirrors gpu_runtime.c's GPU_MAX_LEN.
#define GPU_MAX_LEN (INT64_MAX / 8)

// A safetensors header is JSON text; one this large is not a header.
#define SAFETENSORS_HEADER_MAX ((int64_t)100 * 1024 * 1024)
#define SAFETENSORS_PREFIX 8

// Elements per fread while widening a tensor into its buffer.
#define LOAD_CHUNK 4096

// The NULL of a failed operation, with the buffer it had allocated released.
// The ARC corpus lane found two live objects per failed load before this
// existed: tests/core/gpu/transformer pins it, memory_golden_tests.c counts it.
static void *abandoned(void *buffer) {
  osp_release(buffer);
  return NULL;
}

static const char OP_MATMUL[] = "gpuMatMul";
static const char OP_LOAD[] = "gpuLoad";
static const char OP_TENSOR[] = "gpuLoadTensor";

// Broken-code quarantine: the GPT-2 example used to turn failed tensors and
// products into empty buffers, then generate/train with fabricated zeros.
void *osprey_gpu_panic(const char *message) {
  (void)fprintf(stderr, "GPU execution cannot continue: %s\n", message ? message : "unknown GPU failure");
  abort();
}

static double word_double(int64_t word) {
  double d;
  memcpy(&d, &word, sizeof d);
  return d;
}

static int64_t double_word(double d) {
  int64_t word;
  memcpy(&word, &d, sizeof word);
  return word;
}

// --- the dense product -------------------------------------------------------

// rows × cols as an element count, or -1 when negative or overflowing.
static int64_t area(int64_t rows, int64_t cols) {
  if (rows < 0 || cols < 0 || (rows != 0 && cols > GPU_MAX_LEN / rows)) {
    return -1;
  }
  return rows * cols;
}

// Reject an operand whose length is not the stated shape, naming which one.
static int shaped(void *buffer, int64_t count, const char *operand) {
  if (count >= 0 && osprey_gpu_len(buffer) == count) {
    return 1;
  }
  char reason[96];
  (void)snprintf(reason, sizeof reason, "has %lld elements, the stated shape needs %lld",
                 (long long)osprey_gpu_len(buffer), (long long)count);
  osp_io_error_text(OP_MATMUL, operand, reason);
  return 0;
}

// One rank-1 update, out[r][*] += left[r][k] * right[k][*]: the product and the
// sum are separate statements so no target contracts them into a fused
// multiply-add, and every output element still accumulates in k order.
static void axpy_row(int64_t *out, const int64_t *right, double x, int64_t cols) {
  for (int64_t c = 0; c < cols; c++) {
    double product = x * word_double(right[c]);
    out[c] = double_word(word_double(out[c]) + product);
  }
}

// out[r*cols + c] = sum over k, in order from 0.0, of left[r*inner + k] *
// right[k*cols + c], in double. NULL with the reason recorded when either
// operand's length disagrees with the shape or the shape overflows.
static void *matrix_result(void *left, void *right, int64_t rows, int64_t inner, int64_t cols) {
  int64_t n = area(rows, cols);
  if (n < 0) {
    osp_io_error_text(OP_MATMUL, "result", "rows * cols is negative or overflows");
    return NULL;
  }
  if (!shaped(left, area(rows, inner), "left operand") ||
      !shaped(right, area(inner, cols), "right operand")) {
    return NULL;
  }
  void *out = osprey_gpu_alloc(n);
  if (!out || osprey_gpu_len(out) != n) {
    osp_io_error_set(OP_MATMUL, "result", ENOMEM);
    return abandoned(out);
  }
  return out;
}

static void host_product(int64_t *o, const int64_t *a, const int64_t *b, int64_t rows, int64_t inner, int64_t cols) {
  for (int64_t r = 0; o && a && b && r < rows; r++) {
    for (int64_t k = 0; k < inner; k++)
      axpy_row(o + r * cols, b + k * cols, word_double(a[r * inner + k]), cols);
  }
}

void *osprey_gpu_matmul(void *left, void *right, int64_t rows, int64_t inner, int64_t cols) {
  void *out = matrix_result(left, right, rows, inner, cols);
  if (!out) return NULL;
  int64_t *o = osprey_gpu_data(out);
  const int64_t *a = osprey_gpu_data(left);
  const int64_t *b = osprey_gpu_data(right);
  int device = osprey_gpu_device_matmul(a, b, o, rows, inner, cols);
  if (device != 0) return device > 0 ? out : abandoned(out);
  host_product(o, a, b, rows, inner, cols);
  return out;
}

// --- stored float formats ------------------------------------------------------

typedef enum { DT_F64 = 8, DT_F32 = 4, DT_F16 = 2, DT_BF16 = 3 } Dtype;

// Byte width of one stored element.
static size_t dtype_width(Dtype t) {
  return t == DT_BF16 ? 2 : (size_t)t;
}

static double from_f32_bits(uint32_t bits) {
  float f;
  memcpy(&f, &bits, sizeof f);
  return (double)f;
}

// IEEE-754 binary16 → double, exactly: normals, subnormals, zeros, infinities
// and NaN each widen to the double with the same value.
static double from_f16_bits(uint16_t bits) {
  int sign = (bits >> 15) & 1;
  int exponent = (bits >> 10) & 0x1f;
  uint32_t mantissa = bits & 0x3ff;
  double magnitude;
  if (exponent == 0x1f) {
    magnitude = mantissa ? __builtin_nan("") : __builtin_inf();
  } else if (exponent == 0) {
    magnitude = (double)mantissa * 0x1p-24;
  } else {
    magnitude = (double)(mantissa | 0x400) * __builtin_ldexp(0x1p-10, exponent - 15);
  }
  return sign ? -magnitude : magnitude;
}

// Widen the little-endian element at `at` to a double word.
static int64_t widen(Dtype t, const unsigned char *at) {
  switch (t) {
  case DT_F64: {
    double d;
    memcpy(&d, at, sizeof d);
    return double_word(d);
  }
  case DT_F32:
    return double_word(from_f32_bits((uint32_t)at[0] | ((uint32_t)at[1] << 8) |
                                     ((uint32_t)at[2] << 16) | ((uint32_t)at[3] << 24)));
  case DT_F16:
    return double_word(from_f16_bits((uint16_t)(at[0] | (at[1] << 8))));
  case DT_BF16:
    return double_word(from_f32_bits(((uint32_t)at[0] << 16) | ((uint32_t)at[1] << 24)));
  }
  return 0;
}

// --- reading -------------------------------------------------------------------

static FILE *open_at(const char *op, const char *path, int64_t offset) {
  FILE *f = path ? fopen(path, "rb") : NULL;
  if (!f) {
    osp_io_error_set(op, path, path ? errno : EINVAL);
    return NULL;
  }
  if (offset < 0 || osp_fseek(f, offset, SEEK_SET) != 0) {
    osp_io_error_set(op, path, offset < 0 ? EINVAL : errno);
    fclose(f);
    return NULL;
  }
  return f;
}

// Read `count` stored elements of type `t` from the stream position into a
// fresh buffer of as many double words. A short read is an error, not a
// shorter buffer: a truncated weight file must never load as a smaller model.
static void *read_words(FILE *f, const char *op, const char *path, Dtype t, int64_t count) {
  unsigned char chunk[LOAD_CHUNK * 8];
  size_t width = dtype_width(t);
  void *out = count > GPU_MAX_LEN ? NULL : osprey_gpu_alloc(count);
  int64_t *words = osprey_gpu_data(out);
  if (!out || osprey_gpu_len(out) != count || (count > 0 && !words)) {
    osp_io_error_set(op, path, ENOMEM);
    return abandoned(out);
  }
  for (int64_t done = 0; done < count;) {
    size_t want = (size_t)(count - done < LOAD_CHUNK ? count - done : LOAD_CHUNK);
    if (fread(chunk, width, want, f) != want) {
      osp_io_error_text(op, path, ferror(f) ? "read failed" : "file ends before the tensor does");
      return abandoned(out);
    }
    for (size_t i = 0; i < want; i++, done++) {
      words[done] = widen(t, chunk + i * width);
    }
  }
  return out;
}

static void *load(const char *op, const char *path, int64_t offset, Dtype t, int64_t count) {
  if (count < 0) {
    osp_io_error_text(op, path, "a negative element count");
    return NULL;
  }
  FILE *f = open_at(op, path, offset);
  if (!f) {
    return NULL;
  }
  void *out = read_words(f, op, path, t, count);
  fclose(f);
  return out;
}

// `count` little-endian binary32 values starting `offset` bytes into `path`.
void *osprey_gpu_load_f32(const char *path, int64_t offset, int64_t count) {
  return load(OP_LOAD, path, offset, DT_F32, count);
}

// --- safetensors ---------------------------------------------------------------

// The header text, NUL-terminated, or NULL with the reason recorded. The
// 8-byte prefix is the header length in little-endian.
static char *read_header(FILE *f, const char *path, int64_t *data_start) {
  unsigned char prefix[SAFETENSORS_PREFIX];
  if (fread(prefix, 1, sizeof prefix, f) != sizeof prefix) {
    osp_io_error_text(OP_TENSOR, path, "shorter than a safetensors header prefix");
    return NULL;
  }
  uint64_t n = 0;
  for (int i = SAFETENSORS_PREFIX - 1; i >= 0; i--) {
    n = (n << 8) | prefix[i];
  }
  if (n > (uint64_t)SAFETENSORS_HEADER_MAX) {
    osp_io_error_text(OP_TENSOR, path, "header length is not plausible for safetensors");
    return NULL;
  }
  char *header = malloc((size_t)n + 1);
  if (!header) {
    osp_io_error_set(OP_TENSOR, path, ENOMEM);
    return NULL;
  }
  if (fread(header, 1, (size_t)n, f) != n) {
    osp_io_error_text(OP_TENSOR, path, "file ends inside the header");
    free(header);
    return NULL;
  }
  header[n] = '\0';
  *data_start = SAFETENSORS_PREFIX + (int64_t)n;
  return header;
}

// The dtype named in one tensor's header entry, or 0 when it is not a float
// format this loader widens.
static Dtype dtype_of(const char *entry) {
  const char *at = strstr(entry, "\"dtype\":\"");
  if (!at) {
    return 0;
  }
  at += strlen("\"dtype\":\"");
  if (strncmp(at, "F64\"", 4) == 0) return DT_F64;
  if (strncmp(at, "F32\"", 4) == 0) return DT_F32;
  if (strncmp(at, "F16\"", 4) == 0) return DT_F16;
  if (strncmp(at, "BF16\"", 5) == 0) return DT_BF16;
  return 0;
}

// `[begin,end]` byte offsets relative to the data section, or 0 when absent.
static int offsets_of(const char *entry, int64_t *begin, int64_t *end) {
  const char *at = strstr(entry, "\"data_offsets\":[");
  if (!at) {
    return 0;
  }
  char *after = NULL;
  *begin = (int64_t)strtoll(at + strlen("\"data_offsets\":["), &after, 10);
  if (!after || *after != ',') {
    return 0;
  }
  *end = (int64_t)strtoll(after + 1, &after, 10);
  return after && *after == ']' && *begin >= 0 && *end >= *begin;
}

// The header entry for `name` — the text between its `{` and `}` — isolated
// in place by terminating the header at the closing brace. Names are matched
// quoted and followed by a colon, so `h.1.x` never matches inside `h.11.x`.
static char *entry_of(char *header, const char *name) {
  size_t len = strlen(name) + 4;
  char *key = malloc(len);
  if (!key) {
    return NULL;
  }
  (void)snprintf(key, len, "\"%s\":", name);
  char *at = strstr(header, key);
  free(key);
  char *close = at ? strchr(at, '}') : NULL;
  if (!close) {
    return NULL;
  }
  *close = '\0';
  return at;
}

// One tensor out of a safetensors file, widened to a float buffer.
void *osprey_gpu_load_tensor(const char *path, const char *name) {
  FILE *f = open_at(OP_TENSOR, path, 0);
  if (!f) {
    return NULL;
  }
  int64_t data_start = 0;
  char *header = read_header(f, path, &data_start);
  fclose(f);
  if (!header) {
    return NULL;
  }
  const char *entry = entry_of(header, name ? name : "");
  Dtype t = entry ? dtype_of(entry) : 0;
  int64_t begin = 0;
  int64_t end = 0;
  int located = entry && t && offsets_of(entry, &begin, &end);
  free(header);
  if (!located) {
    osp_io_error_text(OP_TENSOR, name, !entry ? "no tensor of that name in the header"
                                       : !t ? "dtype is not F64, F32, F16 or BF16"
                                            : "data_offsets are missing or inverted");
    return NULL;
  }
  int64_t width = (int64_t)dtype_width(t);
  if ((end - begin) % width != 0) {
    osp_io_error_text(OP_TENSOR, name, "byte span is not a whole number of elements");
    return NULL;
  }
  return load(OP_TENSOR, path, data_start + begin, t, (end - begin) / width);
}
