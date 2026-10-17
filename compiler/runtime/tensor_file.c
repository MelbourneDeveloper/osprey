// The safetensors reader (tensor_file.h) and the host loader built on it,
// `loadTensor` [BUILTIN-FILE-TENSOR]: one named tensor of a safetensors file —
// the Hugging Face interchange format, an 8-byte little-endian header length,
// a JSON header mapping each name to its dtype, shape and data_offsets, then
// the raw tensor bytes — straight into a List<float>. No GpuBuffer exists on
// this path; gpu_host.c reads through the same functions into its own buffers.
//
// Every stored float widens EXACTLY to the double a List<float> or
// GpuBuffer<float> holds. Every failure path records WHY on the io_error
// channel; codegen turns a NULL result into `Error { message }`
// [BUILTIN-FILE-ERRMSG].

#include "collection_runtime.h"
#include "io_error.h"
#include "memory_hooks.h"
#include "tensor_file.h"

#include <errno.h>
#include <limits.h>
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

// A safetensors header is JSON text; one this large is not a header.
#define SAFETENSORS_HEADER_MAX ((int64_t)100 * 1024 * 1024)
#define SAFETENSORS_PREFIX 8

// Elements per fread while widening a tensor.
#define LOAD_CHUNK 4096

static const char OP_LOAD_TENSOR[] = "loadTensor";

// --- stored float formats ------------------------------------------------------

// Byte width of one stored element.
static size_t dtype_width(OspDtype t) {
  return t == OSP_DT_BF16 ? 2 : (size_t)t;
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

// Widen the little-endian element at `at` to a double.
static double widen(OspDtype t, const unsigned char *at) {
  switch (t) {
  case OSP_DT_F64: {
    double d;
    memcpy(&d, at, sizeof d);
    return d;
  }
  case OSP_DT_F32:
    return from_f32_bits((uint32_t)at[0] | ((uint32_t)at[1] << 8) | ((uint32_t)at[2] << 16) |
                         ((uint32_t)at[3] << 24));
  case OSP_DT_F16:
    return from_f16_bits((uint16_t)(at[0] | (at[1] << 8)));
  case OSP_DT_BF16:
    return from_f32_bits(((uint32_t)at[0] << 16) | ((uint32_t)at[1] << 24));
  }
  return 0.0;
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

static int read_span(FILE *f, const char *op, const char *path, const OspTensorSpan *span,
                     OspTensorSink sink, void *ctx) {
  unsigned char chunk[LOAD_CHUNK * 8];
  double values[LOAD_CHUNK];
  size_t width = dtype_width(span->dtype);
  for (int64_t done = 0; done < span->count;) {
    size_t want = (size_t)(span->count - done < LOAD_CHUNK ? span->count - done : LOAD_CHUNK);
    if (fread(chunk, width, want, f) != want) {
      osp_io_error_text(op, path, ferror(f) ? "read failed" : "file ends before the tensor does");
      return 0;
    }
    for (size_t i = 0; i < want; i++) {
      values[i] = widen(span->dtype, chunk + i * width);
    }
    sink(ctx, values, want);
    done += (int64_t)want;
  }
  return 1;
}

int osp_tensor_read(const char *op, const char *path, const OspTensorSpan *span, OspTensorSink sink,
                    void *ctx) {
  if (span->count < 0) {
    osp_io_error_text(op, path, "a negative element count");
    return 0;
  }
  FILE *f = open_at(op, path, span->offset);
  if (!f) {
    return 0;
  }
  int read = read_span(f, op, path, span, sink, ctx);
  fclose(f);
  return read;
}

// --- safetensors ---------------------------------------------------------------

// The header text, NUL-terminated, or NULL with the reason recorded. The
// 8-byte prefix is the header length in little-endian.
static char *read_header(FILE *f, const char *op, const char *path, int64_t *data_start) {
  unsigned char prefix[SAFETENSORS_PREFIX];
  if (fread(prefix, 1, sizeof prefix, f) != sizeof prefix) {
    osp_io_error_text(op, path, "shorter than a safetensors header prefix");
    return NULL;
  }
  uint64_t n = 0;
  for (int i = SAFETENSORS_PREFIX - 1; i >= 0; i--) {
    n = (n << 8) | prefix[i];
  }
  if (n > (uint64_t)SAFETENSORS_HEADER_MAX) {
    osp_io_error_text(op, path, "header length is not plausible for safetensors");
    return NULL;
  }
  char *header = malloc((size_t)n + 1);
  if (!header) {
    osp_io_error_set(op, path, ENOMEM);
    return NULL;
  }
  if (fread(header, 1, (size_t)n, f) != n) {
    osp_io_error_text(op, path, "file ends inside the header");
    free(header);
    return NULL;
  }
  header[n] = '\0';
  *data_start = SAFETENSORS_PREFIX + (int64_t)n;
  return header;
}

// The dtype named in one tensor's header entry, or 0 when it is not a float
// format this reader widens.
static OspDtype dtype_of(const char *entry) {
  const char *at = strstr(entry, "\"dtype\":\"");
  if (!at) {
    return 0;
  }
  at += strlen("\"dtype\":\"");
  if (strncmp(at, "F64\"", 4) == 0) return OSP_DT_F64;
  if (strncmp(at, "F32\"", 4) == 0) return OSP_DT_F32;
  if (strncmp(at, "F16\"", 4) == 0) return OSP_DT_F16;
  if (strncmp(at, "BF16\"", 5) == 0) return OSP_DT_BF16;
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

// The header of `path`, NULL with the reason recorded.
static char *header_of(const char *op, const char *path, int64_t *data_start) {
  FILE *f = open_at(op, path, 0);
  if (!f) {
    return NULL;
  }
  char *header = read_header(f, op, path, data_start);
  fclose(f);
  return header;
}

int osp_tensor_locate(const char *op, const char *path, const char *name, OspTensorSpan *span) {
  int64_t data_start = 0;
  char *header = header_of(op, path, &data_start);
  if (!header) {
    return 0;
  }
  const char *entry = entry_of(header, name ? name : "");
  OspDtype t = entry ? dtype_of(entry) : 0;
  int64_t begin = 0;
  int64_t end = 0;
  int located = entry && t && offsets_of(entry, &begin, &end);
  free(header);
  if (!located) {
    osp_io_error_text(op, name, !entry ? "no tensor of that name in the header"
                                : !t   ? "dtype is not F64, F32, F16 or BF16"
                                       : "data_offsets are missing or inverted");
    return 0;
  }
  int64_t width = (int64_t)dtype_width(t);
  if ((end - begin) % width != 0) {
    osp_io_error_text(op, name, "byte span is not a whole number of elements");
    return 0;
  }
  *span = (OspTensorSpan){data_start + begin, t, (end - begin) / width};
  return 1;
}

// --- loadTensor ----------------------------------------------------------------

// Push each value as the IEEE-754 word a List<float> element holds.
static void into_list(void *ctx, const double *values, size_t n) {
  for (size_t i = 0; i < n; i++) {
    int64_t word;
    memcpy(&word, values + i, sizeof word);
    osprey_list_builder_push_of(ctx, word, 0);
  }
}

// One tensor out of a safetensors file as a List<float>, or NULL with the
// reason recorded. A list built before a failed read is released, so a NULL
// result owes the memory backend nothing.
OspreyList *osprey_load_tensor(const char *path, const char *name) {
  OspTensorSpan span;
  if (!osp_tensor_locate(OP_LOAD_TENSOR, path, name, &span)) {
    return NULL;
  }
  OspreyListBuilder *b = osprey_list_builder_new_of(0);
  int read = osp_tensor_read(OP_LOAD_TENSOR, path, &span, into_list, b);
  OspreyList *list = osprey_list_builder_seal(b);
  if (!read) {
    osp_release(list);
    return NULL;
  }
  return list;
}

// Broken-code quarantine: a host interpretation whose weights fail to load
// would otherwise run on whatever its caller substituted. It aborts naming why.
void *osprey_tensor_panic(const char *message) {
  (void)fprintf(stderr, "tensor load cannot continue: %s\n", message ? message : "unknown failure");
  abort();
}
