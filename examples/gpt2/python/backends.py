# osprey: flavor=python
# osprey: namespace=gpt2
# osprey: module=Backends
"""The tensor dialect GPT-2 is written in.

`gpt2.Model` performs `Tensor` operations and never names a buffer type;
`Cpu.on_lists` (cpu.py) answers every operation with host list code and
`Gpu.on_gpu` (gpu.py) with the GPU surface. The handler a program installs decides
where the arithmetic runs, and both must agree with examples/gpt2/reference.py bit
for bit, because every arm accumulates in the same k order from 0.0.
Implements [STAGE-BACKEND], [GPU-DEVICE], [FLAVOR-PYTHON-GPU].
"""
from typing import Callable

from gpt2 import Cpu, Gpu

__all__ = ["Kernel", "Combine", "Tensor", "hit", "min_i", "strictly", "under"]


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
            return Gpu.on_gpu(work)
        case "lists":
            return Cpu.on_lists(work)
        case _:
            return print(f"unknown backend `{name}`: choose gpu or lists")
