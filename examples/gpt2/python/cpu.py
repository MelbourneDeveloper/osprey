# osprey: flavor=python
# osprey: namespace=gpt2
# osprey: module=Cpu
"""The CPU interpretation of `gpt2.Backends.Tensor`: every operation answered with
host list code, accumulating in the same k order from 0.0 as
examples/gpt2/reference.py. Nothing here names the GPU surface: weights load
through the host `loadTensor` and a failed load aborts through the host runtime.
Implements [STAGE-BACKEND], [FLAVOR-PYTHON-GPU], [BUILTIN-FILE-TENSOR].
"""
from gpt2.Backends import Tensor, hit, min_i

__all__ = ["on_lists"]


@extern
def osprey_tensor_panic(message: str) -> list[float]: ...


def or_abort(result):
    """Broken-code quarantine: a failed load or product used to become an empty
    tensor, and the model generated from fabricated zeros. It aborts instead,
    through this handler's own runtime."""
    match result:
        case Success(value=value):
            return value
        case Error(message=message):
            return osprey_tensor_panic(message)


@static_handler(Tensor[list[float]])
class on_lists:
    """Host lists, loaded straight into lists from the safetensors file."""
    def build(n, f):
        return toList(map(range(n), f))

    def get(b, i):
        return unwrap_or(listGet(b, i), 0.0)

    def length(b):
        return len(b)

    def zip_with(a, b, f):
        return toList(map(range(len(a)), lambda i: f(unwrap_or(listGet(a, i), 0.0), unwrap_or(listGet(b, i), 0.0))))

    def reduce(b, seed, f):
        return fold(range(len(b)), seed, lambda a, i: f(a, unwrap_or(listGet(b, i), 0.0)))

    def index_of(b, x):
        return fold(map(range(len(b)), lambda i: hit(unwrap_or(listGet(b, i), 0.0), x, i, len(b))), len(b), min_i)

    def mat_mul(a, b, rows, inner, cols):
        return toList(map(range(rows * cols), lambda i: fold(range(inner), 0.0, lambda acc, k: acc + unwrap_or(listGet(a, i // cols * inner + k), 0.0) * unwrap_or(listGet(b, k * cols + i % cols), 0.0))))

    def from_list(xs):
        return xs

    def load(path, name):
        return or_abort(loadTensor(path, name))

    def device():
        return "lists"
