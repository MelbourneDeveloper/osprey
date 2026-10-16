# GPT-2 in Python syntax

The [model](model.py) performs the existing `Tensor` effect. The handler selects host lists or the GPU surface; the forward mathematics stays in Python-flavor source.

```python
def matmul(x, w, rows, inner, cols):
    return perform(Tensor.matMul(x, w, rows, inner, cols))
```

From the repository root:

```sh
make gpt2 FLAVOR=python GPU=host PROMPT="The capital of France is" TOKENS=1
make gpt2 FLAVOR=python BACKEND=lists PROMPT="The capital of France is" TOKENS=1
```

The GPU-surface host run produces `The capital of France is the`, matching the independent reference. `GPU=auto`, `GPU=metal` and `GPU=opencl` select the existing device adapters. A missing adapter stops execution with an error.

The [sampler](generate/generate.py) recomputes the forward pass for each token. The existing Default example uses a KV cache, so its timings measure different work. Tokenization and weight loading reuse the existing Osprey modules; the independent [Python reference](../reference.py) stays separate.

The ordinary compiler test compares every element of 24 deterministic tensors against that reference, at the real embedding width and sequence lengths 1/2/3, under both handlers and default/GC/ARC memory. This includes all intermediate values from a transformer block. Run it directly with an existing compiler:

```sh
python3 scripts/check-python-gpt2.py --compiler target/debug/osprey
```

Pass `--device metal`, `--device opencl` or `--device auto` for a strict hardware comparison. That mode requires an adapter identity and completed device dispatches. Host arithmetic compares exactly; only values downstream of `exp`/`tanh` use the existing GPU precision bound.

The [port notes](../../../docs/plans/python-flavor-port/README.md) show the concrete changes from reference Python.
