// The safetensors reader every tensor loader shares: `loadTensor` builds a
// host List<float> from it [BUILTIN-FILE-TENSOR], `gpuLoadTensor` and `gpuLoad`
// a GpuBuffer<float> [GPU-BUFFER-TENSOR] [GPU-BUFFER-LOAD]. It knows no
// container: it locates a tensor's bytes and hands the widened doubles, in
// file order, to the caller's sink, so neither surface depends on the other.
#ifndef OSPREY_TENSOR_FILE_H
#define OSPREY_TENSOR_FILE_H

#include <stddef.h>
#include <stdint.h>

// Stored float formats, valued by byte width (BF16 is 2 bytes, tagged 3).
typedef enum { OSP_DT_F64 = 8, OSP_DT_F32 = 4, OSP_DT_F16 = 2, OSP_DT_BF16 = 3 } OspDtype;

// Where a run of stored elements lives: the byte offset of the first, the
// stored format, and how many there are.
typedef struct {
  int64_t offset;
  OspDtype dtype;
  int64_t count;
} OspTensorSpan;

// Receives the next `n` widened values, in file order.
typedef void (*OspTensorSink)(void *ctx, const double *values, size_t n);

// The span of tensor `name` in safetensors file `path`: 1, or 0 with the
// reason recorded on the io_error channel under `op`.
int osp_tensor_locate(const char *op, const char *path, const char *name, OspTensorSpan *span);

// Widen every element of `span` exactly to a double and hand them to `sink`:
// 1, or 0 with the reason recorded. A short read is an error, never a shorter
// tensor: a truncated weight file must not load as a smaller model.
int osp_tensor_read(const char *op, const char *path, const OspTensorSpan *span, OspTensorSink sink,
                    void *ctx);

#endif
