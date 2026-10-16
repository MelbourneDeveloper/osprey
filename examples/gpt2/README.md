# GPT-2 in Osprey

The trained GPT-2 model, exactly as Hugging Face publishes it, written once
against a static `Tensor` effect and run by two handlers: one over plain host
lists, one on Osprey's GPU surface. No training from scratch, no conversion
step, no Python at run time: the `load` operation reads the published
`model.safetensors`, the byte-level BPE tokenizer is rebuilt from the published
`merges.txt`, and the transformer performs `Tensor` operations that the handler
a program installs decides how to execute. On the CPU both handlers print the
same tokens, the same loss and the same 148 gradient sums; on a GPU the text is
the same and the numbers agree to twelve digits.

```
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=gpu
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=lists
make gpt2-train PROMPT="The capital of France is" BACKEND=gpu
make gpt2-train PROMPT="The capital of France is" BACKEND=lists
```

The first run downloads the two files (about 550 MB) into
`examples/gpt2/weights/`, which is ignored by git. `BACKEND` defaults to `gpu`.
`make gpt2` prints which handler answered, the prompt's token ids, and the
greedy continuation:

```
gpt2 (124M) on host — 50000 merges, 50257 vocabulary rows, 12 layers
prompt: The capital of France is
tokens: 464 3139 286 4881 318
generated 20 tokens
The capital of France is the capital of the French Republic, and the capital of the French Republic is the capital of the French
```

That first line is from `GPU=host`. The default `GPU=auto` names the adapter instead, such as `gpt2 (124M) on metal:0 (Apple M4 Max)`, and prints the same four lines after it; `BACKEND=lists` prints `gpt2 (124M) on lists`.

`make gpt2-train` runs one training step on the same downloaded weights: the
next-token cross-entropy of the prompt, the gradient of every parameter through
the twelve blocks, one SGD step, and the loss again. It prints the loss, a
checksum of every gradient tensor (one line per block, in the checkpoint's
tensor order: a vector's elements added in order; a matrix's rows each added in
order, then the rows), and the loss after the step:

```
gpt2 (124M) on host — one training step over 5 tokens, learning rate 0.0001
prompt: The capital of France is
tokens: 464 3139 286 4881 318
loss 4.503189215447316
grad wte -6.369071936518367e-13 wpe -1.4510961876545991e-15 ln_f 0.20932690619905386 0.512850562773751
h.0 -3.3628631026293108 -7.1375868853910704 3.8554208861589614 …
…
loss after one step 3.8539085307403895
```

## Layout

| Path | What it is |
| --- | --- |
| `model/backends.osp` | `gpt2::Backends` — the `Tensor<B>` effect (build, get, length, zipWith, reduce, indexOf, matMul, fromList, load, device) and its two handlers, `onLists` over `List<float>` and `onGpu` over `GpuBuffer<float>` |
| `model/tokenizer.osp` | `gpt2::Tokenizer` — the byte-level BPE tokenizer from `merges.txt` alone |
| `model/model.osp` | `gpt2::Model` — the configuration, the weights, LayerNorm, attention, the blocks, the forward pass keeping every block's activations, and the KV-cached step; every tensor is `B` |
| `generate/` | A one-file project: greedy sampling over the forward pass and the cached steps; `main` picks the handler |
| `train/` | A one-file project: the loss, the backward pass and one SGD step over the same forward pass; `main` picks the handler |
| `reference.py` | The independent pure-Python oracle for both programs |

A standalone script cannot import across files, so each program is its own
project naming `model/` as a second source root
([modules](../../docs/specs/language/0025-ModulesAndNamespaces.md)), the layout the
graphics examples use. `osprey examples/gpt2/generate --run` runs a project
directory; the Makefile targets add `--memory=arc`.

## One model, two handlers

`gpt2::Model` never names a list or a buffer. Its records are generic in the
tensor representation — `Layer<B>`, `Weights<B>`, `Acts<B>` — and every
operation is a `perform`:

```osprey
export fn layerNorm(x, t, g, b) = {
    let means = rowMeans(x, t)
    let variances = rowVariances(x, t, means)
    perform Backends::Tensor.build(t * width, fn(i) => normalised(x, means, variances, g, b, i))
}
```

`onLists` answers `build` with `range |> map |> toList` and `matMul` with a
`fold` over `k`; `onGpu` answers them with `gpuIota |> gpuMap` and `gpuMatMul`,
forwarding the model's callbacks to the kernels. `main` reads the backend name
from stdin and `Backends::under` installs the handler it names around the whole
run; any other name runs nothing:

```osprey
Backends::under(input(), fn() => run(prompt, count))
```

The effect is `static`, so the compiler discharges it before the program runs:
each arm of the `match` in `under` is its own specialised copy of the model with the
handler's arms inlined at every `perform`, and no `Tensor` handler exists at run
time. [One model, two backends](../../website/src/docs/one-model-two-backends.md)
walks through the effect, both handlers and what the compiler proves about them.

## What runs where

`BACKEND=gpu` (the default) runs the model under `onGpu`, and `GPU` picks where its kernels and matrix products execute. `make gpt2` and `make gpt2-train` request the GPU by default (`GPU=auto`). On macOS, they dispatch directly to Metal; elsewhere they dispatch to an OpenCL compute GPU through its installed driver. `GPU=metal` or `GPU=opencl` requests a specific adapter, and `GPU=host` keeps every kernel and product on the CPU as the oracle. `BACKEND=lists` runs the same model over host lists, on the CPU, whatever `GPU` says. The same model source runs on every path, without CUDA. `gpuDevice()` names the resolved adapter and labels matrix-kernel execution, and both programs print it on their first line; a completed-command message on stderr proves a matrix reached the GPU. Adapter failures abort the model rather than continuing with empty buffers.

```sh
make gpt2 GPU=host      # every kernel and product on the CPU
make gpt2 GPU=metal     # every kernel and product on the Apple GPU
make gpt2 BACKEND=lists # host lists, CPU only
```

To run the compiler interactively from the repository root:

```sh
OSPREY_GPU=auto osprey examples/gpt2/generate --run --memory=arc
```

Enter the prompt, the token count and the backend (`gpu` or `lists`) on three lines. Set `OSPREY_GPU_TRACE=1` to report every completed dispatch. Under `BACKEND=gpu` with a device selected, every kernel the model performs — embeddings, LayerNorm, attention scores and softmax, GELU, the residual adds, every reduction, the next token's index, the training loss — and every matrix product executes on that device, compiled from the same Osprey source the host runs ([GPU-BACKEND-DEVICE](../../docs/specs/runtime/0034-GPUComputation.md#device-backends--gpu-backend-device)). Both adapters use the same software binary64 arithmetic, preserving the host's accumulation order without float32 substitution. What stays on the CPU: tokenization, building and loading buffers, and one download of each value the program reads back — the logits a token's index is taken from, and the loss. Every buffer is resident on the device from the first kernel that touches it ([GPU-BUFFER-RESIDENT](../../docs/specs/runtime/0034-GPUComputation.md#device-resident-buffers--gpu-buffer-resident)): the weights upload once, a layer's intermediate buffers never leave the device, and a buffer's device copy is released with the buffer. A kernel the device dialect cannot express aborts the run naming itself rather than running on the CPU; `GPU=host` runs the same program on the CPU and `BACKEND=lists` runs it over plain lists without touching the GPU surface at all.

The device's `exp`, `ln` and `tanh` are correctly rounded; a host libm is not on every argument. The generated text is identical on every path. The training loss agrees to fifteen significant digits and differs in the last (`4.503189215447316` on the host, `4.503189215447317` on Metal). Across all 150 numbers the step prints, the largest difference from the host is 2.4 × 10⁻¹³ and no number larger than 10⁻⁶ differs by more than 4 parts in 10¹³, so a gradient sum that cancels to about `1e-15`, such as `wpe`, differs in every digit. Metal and OpenCL print identical bytes apart from the device name.

| Measured 2026-10-13, AMD Ryzen 9 3900X, Linux 6.18 (WSL2), release build, `--memory=arc`, `GPU=host` | `BACKEND=gpu` | `BACKEND=lists` |
| --- | --- | --- |
| 20 tokens from the 5-token prompt: load, one full pass, 19 KV-cached steps | 19.6 s, 1.4 GB peak | 7 min 18 s, 1.7 GB peak |
| One training step over the 5-token prompt: load, forward, loss, backward, SGD, forward again | 52 s, 3.6 GB peak | 9 min, 4.6 GB peak |
| `reference.py`, the same 20 tokens in pure Python | 16 min 25 s | |

| Measured 2026-10-16, Apple M4 Max, macOS, release build, `--memory=arc`, `BACKEND=gpu`, median of three | `GPU=host` | `GPU=metal` | `GPU=opencl` |
| --- | --- | --- | --- |
| 20 tokens from the 5-token prompt, every kernel and product on the selected device, buffers resident on it | 6.5 s, 1.7 GB peak | 6.9 s, 3.1 GB peak | 6.8 s, 1.7 GB peak |
| One training step over the 5-token prompt | 25.2 s, 4.0 GB peak | 7.6 s, 5.9 GB peak | 7.8 s, 3.1 GB peak |
| `reference.py --train`, the same step in pure Python (2026-10-15) | 5 min 49 s, 16.2 GB peak | | |

Training is about three times faster on the device: half its 1 499 dispatches run over 3 072 threads or more. Generation is not: after the prompt each token is a one-row pass through the KV cache, 4 338 of its 5 919 dispatches run over 768 threads or fewer, and without kernel fusion each is its own launch.

On the CPU the GPU-surface handler beats the list handler because its buffers
are dense doubles and its kernels are lifted into flat functions; the list handler boxes every element and calls
the callback through a closure. The Makefile targets run under `--memory=arc`
because of the memory column.

## Checking the output

[`reference.py`](reference.py) is an independent pure-Python implementation
(no numpy) that reads the same two files and runs the same forward and backward
passes with the same left-to-right accumulation order. Python's built-in `sum`
has used compensated summation for floats since 3.12, so the script adds
sequentially from `0.0` itself. Run both on one prompt and compare the output:

```
python3 examples/gpt2/reference.py examples/gpt2/weights "The capital of France is" 20
python3 examples/gpt2/reference.py examples/gpt2/weights --train "The capital of France is"
```

On the prompts tried so far — the default, `Hello, I'm a language model,`,
`  two  spaces 123abc don't` and `naïve café — 日本` — the tokenizer produces
identical ids in both implementations, and the greedy continuations compared
so far are identical: all 20 tokens of the default prompt, and the two of
`Hello, I'm a language model,` (`not a`). For the training step on the default
prompt, the loss, every one of the 148 gradient-tensor sums and the loss after
the step print identically in both implementations. Both Osprey handlers print
the same bytes as each other, because both accumulate in the same order.

The same model at width 4 — one block written against the same `Tensor`
effect, run under both handlers — is pinned against values the Python oracle
computed in
[`tests/effects/staged/tensor_backends.test.osp`](../../tests/effects/staged/tensor_backends.test.osp)
and its ML twin, which must emit identical LLVM IR; the kernels themselves are
pinned in
[`tests/core/gpu/transformer.test.osp`](../../tests/core/gpu/transformer.test.osp).
The differential harness runs both under every memory backend, both kernel
lowerings and wasm32 on every `make test`.

## How it is written

- **Weights** come through the `load` operation: `tensor("h.3.ln_2.bias")` is a `B`. A failed load aborts with its error instead of substituting an empty tensor. Under `onGpu`, `load` uses `gpuLoadTensor`; under `onLists`, the same loader reads back into a list. Shapes are stated at each `matMul`, because tensors are rank-1.
- **Tokenizer** — GPT-2 encodes bytes, and merge *i* of `merges.txt` is token
  `256 + i`, so the merge list alone is the vocabulary; `vocab.json` is never
  read. Pre-tokenization approximates GPT-2's regular expression for ASCII
  text and treats every non-ASCII byte as a letter.
- **Attention, LayerNorm, GELU and softmax** are `build` operations over flat
  index arithmetic, with `exp`, `sqrt` and `tanh` called inside the index
  function and `fold` over `range` for the inner reductions.
- **Products** are `matMul`. Both handlers accumulate each element in k order
  from 0.0 — the reference order a device backend must match.
- **The KV cache** — the prompt takes one full forward pass; every later token
  runs `forwardStep`, one position through the blocks over each block's cached
  qkv rows, so a token costs one position rather than the whole sequence.
- **Sampling** is greedy: a `reduce` for the peak logit and an `indexOf` for
  its first index, both kernels.
- **The backward pass** is the same kinds of operation. A `linear`'s gradients
  are two products against transposed operands plus a column sum; LayerNorm's
  recomputes the normalised rows from its input; attention's is the softmax
  Jacobian per row and then the query, key and value columns of each head.
  The forward pass keeps each block's activations in an `Acts<B>` record so the
  backward pass reads each exactly once.

## Known limits

- Three stdin lines for `generate` (prompt, token count, backend) and two for
  `train` (prompt, backend); an unknown backend name selects the GPU surface.
- The handler is chosen where it is written: `match input()` selects between
  two regions the compiler has already specialised. A handler computed at run
  time would need a dynamic effect and a runtime handler.
- The pre-tokenizer is an ASCII approximation of GPT-2's Unicode regular
  expression; prompts in other scripts tokenize correctly only where the
  letter/number/punctuation classes agree with the ASCII rules.
- Computation is double precision throughout (the fp32 weights widen exactly),
  so logits differ from an fp32 implementation in the last digits; the greedy
  choice has agreed on every token compared so far.
- Training is one plain SGD step at a fixed learning rate: no optimiser state,
  no batching, no weight decay, and the updated weights are not written back.
