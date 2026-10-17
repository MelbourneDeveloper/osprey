# osprey: flavor=python
# osprey: namespace=gpt2
# osprey: module=Cpu
"""The CPU interpretation of `gpt2.Backends.Tensor`: every operation answered with
host list code, accumulating in the same k order from 0.0 as
examples/gpt2/reference.py.
Implements [STAGE-BACKEND], [FLAVOR-PYTHON-GPU].
"""
from gpt2.Backends import Tensor, or_abort, hit, min_i

__all__ = ["on_lists"]


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
