# osprey: flavor=python
# osprey: namespace=gpt2
# osprey: module=Gpu
"""The GPU interpretation of `gpt2.Backends.Tensor`: every operation is a kernel, a
product or a loader of the GPU surface.
Implements [STAGE-BACKEND], [GPU-DEVICE], [FLAVOR-PYTHON-GPU].
"""
from gpt2.Backends import Tensor, hit, min_i

__all__ = ["on_gpu"]


@extern
def osprey_gpu_panic(message: str) -> GpuBuffer[float]: ...


def or_abort(result):
    """Broken-code quarantine: a failed load or product used to become an empty
    tensor, and the model generated from fabricated zeros. It aborts instead,
    through this handler's own runtime."""
    match result:
        case Success(value=value):
            return value
        case Error(message=message):
            return osprey_gpu_panic(message)


def fetched(b, i):
    return unwrap_or(gpuGet(b, i), 0.0)


@static_handler(Tensor[GpuBuffer[float]])
class on_gpu:
    """GPU buffers, and `device` reports which backend executes them. Kernels and
    products run on Metal or OpenCL when `OSPREY_GPU` selects a device."""
    def build(n, f):
        return gpuMap(gpuIota(n), f)

    def get(b, i):
        return fetched(b, i)

    def length(b):
        return gpuLength(b)

    def zip_with(a, b, f):
        return gpuZipWith(a, b, f)

    def reduce(b, seed, f):
        return gpuFold(b, seed, f)

    def index_of(b, x):
        return gpuFold(gpuMap(gpuIota(gpuLength(b)), lambda i: hit(fetched(b, i), x, i, gpuLength(b))), gpuLength(b), min_i)

    def mat_mul(a, b, rows, inner, cols):
        return or_abort(gpuMatMul(a, b, rows, inner, cols))

    def from_list(xs):
        return toGpu(xs)

    def load(path, name):
        return or_abort(gpuLoadTensor(path, name))

    def device():
        return gpuDevice()
