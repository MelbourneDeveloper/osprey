# osprey: flavor=python
# osprey: namespace=gpt2
# osprey: module=Backends
"""The tensor dialect GPT-2 is written in, and the two interpretations that run it.

`gpt2.Model` performs `Tensor` operations and never names a buffer type; `on_lists`
answers every operation with host list code and `on_gpu` with the GPU surface. The
handler a program installs decides where the arithmetic runs, and both must agree
with examples/gpt2/reference.py bit for bit, because every arm accumulates in the
same k order from 0.0.
Implements [STAGE-BACKEND], [GPU-DEVICE], [FLAVOR-PYTHON-GPU].
"""
from typing import Callable

__all__ = ["Kernel", "Combine", "Tensor", "or_abort", "on_lists", "on_gpu", "strictly", "under"]


@extern
def osprey_gpu_panic(message: str) -> GpuBuffer[float]: ...


def or_abort(result):
    """Broken-code quarantine: a failed load or product used to become an empty
    tensor, and the model generated from fabricated zeros. It aborts instead."""
    match result:
        case Success(value=value):
            return value
        case Error(message=message):
            return osprey_gpu_panic(message)


# The kernels a model hands its handler: an index to its element, and two
# elements to one.
type Kernel = Callable[[int], float]
type Combine = Callable[[float, float], float]


@static_effect
class Tensor[B]:
    """The operations a model needs: build a buffer from an index function, read
    and measure one, combine two, reduce one, find the first index holding a
    value, multiply two matrices, lift a host list, load a tensor from a
    safetensors file, and say which interpretation is answering. `B` is the
    representation the handler picks."""
    def build(size: int, kernel: Kernel) -> B: ...
    def get(buffer: B, index: int) -> float: ...
    def length(buffer: B) -> int: ...
    def zip_with(a: B, b: B, combine: Combine) -> B: ...
    def reduce(buffer: B, seed: float, combine: Combine) -> float: ...
    def index_of(buffer: B, value: float) -> int: ...
    def mat_mul(a: B, b: B, rows: int, inner: int, cols: int) -> B: ...
    def from_list(values: list[float]) -> B: ...
    def load(path: str, name: str) -> B: ...
    def device() -> str: ...


def hit(v, x, i, n):
    """`index_of`'s first step under both handlers: element i names itself when it
    holds x and the length n otherwise."""
    return i if v == x else n


def min_i(a, i):
    """`index_of`'s second step: the smallest name wins, so the answer is the first
    index holding x or n when none does."""
    return i if i < a else a


def element(b, i):
    return unwrap_or(listGet(b, i), 0.0)


def dot(a, b, i, inner, cols):
    return fold(range(inner), 0.0, lambda acc, k: acc + element(a, i // cols * inner + k) * element(b, k * cols + i % cols))


@static_handler(Tensor[list[float]])
class on_lists:
    """Host lists, loaded through the GPU surface's safetensors reader."""
    def build(n, f):
        return toList(map(range(n), f))

    def get(b, i):
        return element(b, i)

    def length(b):
        return len(b)

    def zip_with(a, b, f):
        return toList(map(range(len(a)), lambda i: f(element(a, i), element(b, i))))

    def reduce(b, seed, f):
        return fold(range(len(b)), seed, lambda a, i: f(a, element(b, i)))

    def index_of(b, x):
        return fold(map(range(len(b)), lambda i: hit(element(b, i), x, i, len(b))), len(b), min_i)

    def mat_mul(a, b, rows, inner, cols):
        return toList(map(range(rows * cols), lambda i: dot(a, b, i, inner, cols)))

    def from_list(xs):
        return xs

    def load(path, name):
        return fromGpu(or_abort(gpuLoadTensor(path, name)))

    def device():
        return "lists"


def fetched(b, i):
    return unwrap_or(gpuGet(b, i), 0.0)


@static_handler(Tensor[GpuBuffer[float]])
class on_gpu:
    """GPU buffers: each operation is a kernel, a product or a loader of the GPU
    surface, and `device` reports which backend executes them. Kernels and
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


@effect
class Fault:
    """An arithmetic fault ends the run: its arm answers the whole region."""
    @control
    def halt(message: str) -> int: ...


@handler(Fault)
class halting:
    def halt(fault):
        return Error(message=fault)


@handler(Arith)
class arithmetic:
    def overflow(op, l, r, _):
        return perform(Fault.halt(f"integer overflow: {l} {op} {r}"))

    def divideByZero(op, l):
        return toFloat(perform(Fault.halt(f"division by zero: {l} {op} 0")))

    def remainderByZero(l):
        return perform(Fault.halt(f"remainder by zero: {l} % 0"))


def strictly(work) -> Result[None, str]:
    """Runs `work` under GPT-2's arithmetic policy: Success with its value, or Error
    naming the first fault. A fault is a bug in the model, never a value to
    substitute (a fabricated 0.0 is indistinguishable from a real logit), so it
    stops the run, as a device kernel's fault does [GPU-DEVICE-FAULT]. The
    annotation pins what inference leaves open: each `with` suite is a closure
    whose Result type must be concrete to compile."""
    with halting:
        with arithmetic:
            return Success(value=work())


def under(name, work):
    """Runs `work` under the handler `name` selects. An unknown name runs nothing,
    so a mistyped backend never silently lands on the GPU surface."""
    match name:
        case "gpu":
            return on_gpu(work)
        case "lists":
            return on_lists(work)
        case _:
            return print(f"unknown backend `{name}`: choose gpu or lists")
