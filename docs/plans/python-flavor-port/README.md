# Python-flavor GPT-2 port notes

The authored sources now live in [`examples/gpt2/python`](../../../examples/gpt2/python/README.md). The model and sampler compile through the Python frontend and execute the real 124M checkpoint on the host. [Plan 0031](../0031-python-flavor.md) tracks the complete flavor, device and training acceptance.

| Reference Python | Osprey Python flavor |
| --- | --- |
| `class Tokenizer` and `class Safetensors`, stdlib parsing/loading | Import the existing tokenizer and checkpoint loader; keep the independent reference unchanged |
| Runtime dictionaries for activations/parameters | Existing `Model` records and local `Statistics`/`Projection` records |
| Tensor lists, comprehensions and indexed mutation | Explicit `Tensor.build`, `get`, `zipWith` and `matMul`; the handler selects representation/execution |
| `add_all(values)` | Ordered `fold` from `0.0` over an index function |
| Mutable trace and sampling loops | Immutable record accumulators and `fold` |
| `math.sqrt`, `math.exp`, `math.tanh` | Existing `sqrt`, `exp`, `tanh` |
| Dictionary keys | Named record fields, preserving checkpoint tensor ordering |
| Host token IDs used in embedding | Existing `Model.embed` performs host orchestration and passes scalar IDs to kernels |
| Slicing the final row and dotting each vocabulary row | Build the last-row tensor and multiply by the transposed tied embedding |
| `sys.argv`/`__main__` | Three-line stdin contract and Osprey `main` |

`Tensor.build` appears directly at each kernel construction site so static specialization receives the actual lambda. The model ports the forward mathematics; it does not call `Model.forward`. It shares tokenization, checkpoint data, records and handlers. The sampler recomputes the full pass each token; the Default example's KV-cache timings do not measure this port.

Native compiler validation compares every element of 24 tensors to the unchanged reference: three sequence lengths for LayerNorm and causal attention, small matrix/transpose/linear/residual/GELU cases, and every intermediate activation in one block. Both handlers pass under default, GC and ARC; ARC exits retain no live objects. Real-checkpoint host runs match the reference's first token and decoded output for `The capital of France is` and `Melbourne is the capital of`; the France case also matches under the list handler. These are host execution results; strict hardware comparisons and the backward/training port remain tracked in the plan.
