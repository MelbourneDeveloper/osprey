# GPT-2 in Osprey

This is GPT-2: OpenAI's 124M-parameter language model, running the trained weights exactly as [Hugging Face publishes them](https://huggingface.co/openai-community/gpt2). Nothing is retrained, converted or delegated to a Python runtime. The `load` operation reads the published `model.safetensors`, the byte-level BPE tokenizer is rebuilt from the published `merges.txt`, and the transformer, the KV-cached generation loop and the greedy sampler are all Osprey source.

The model is written in both the Default flavor and the Python flavor:

- **Default flavor**: [`model/`](model) and [`generate/`](generate), plus a [training step](train).
- **Python flavor**: [`python/`](python/README.md), Osprey source that looks like typed Python, with no `.osp` file anywhere in the project. Each file starts with `# osprey: flavor=python` and is compiled by `osprey`; CPython never runs it. The [Python flavor guide](https://www.ospreylang.dev/docs/python-flavor/) explains how it differs from Python.

Given the same prompt, the Default and Python versions print exactly the same bytes on the CPU, on Metal and on OpenCL. The test `python_gpt2_generator_output_is_byte_identical_to_the_default_generator` compares them on each of those devices.

## CPU or GPU is an algebraic effect

The model never names a device or a buffer type. Every tensor operation it needs (build, read, combine, reduce, matrix product, load) is an operation of one static effect, `Tensor<B>`. The model *performs* those operations, and whichever handler the program installs around it decides what `B` is and where the arithmetic runs:

| Handler | Default | Python | Runs on | `B` |
| --- | --- | --- | --- | --- |
| `onLists` / `on_lists` | [`model/cpu.osp`](model/cpu.osp) | [`python/cpu.py`](python/cpu.py) | the CPU | `List<float>` |
| `onGpu` / `on_gpu` | [`model/gpu.osp`](model/gpu.osp) | [`python/gpu.py`](python/gpu.py) | Metal or OpenCL | `GpuBuffer<float>` |

`BACKEND=lists` or `BACKEND=gpu` picks the handler. The model source is identical in both runs, and on the CPU both handlers print the same tokens, the same loss and the same 148 gradient sums; on a GPU the text is the same and the numbers agree to twelve digits.

Why this matters:

- **Where the arithmetic runs is chosen by the caller.** In a typical framework, device placement is a property of each tensor (`x.to("cuda")`) that the model code must thread through every layer. Here it is chosen once, by installing a handler around the computation; the transformer cannot even mention it.
- **A new backend is a new handler, not a new model.** A backend is one handler answering ten operations. Adding one changes no model code.
- **The compiler checks it.** A program that performs a `Tensor` operation without a handler is rejected at compile time. Because `Tensor` is a `static` effect, the handler is selected when the program is compiled, so the effect itself adds no run-time dispatch.
- **The flavor doesn't change any of this.** Effects belong to Osprey's shared AST, so the Python flavor declares the same effect and the same two handlers (as `@static_effect` and `@static_handler` classes) and gets the same checking and the same output.

## Run it

```
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=gpu
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=lists
make gpt2 FLAVOR=python PROMPT="The capital of France is" TOKENS=20 BACKEND=gpu
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
| `model/backends.osp` | `gpt2::Backends` — the `Tensor<B>` effect (build, get, length, zipWith, reduce, indexOf, matMul, fromList, load, device), the arithmetic policy `strictly` and the handler selector `under` |
| `model/cpu.osp` | `gpt2::Cpu` — the CPU handler `onLists`, answering every `Tensor` operation over `List<float>` |
| `model/gpu.osp` | `gpt2::Gpu` — the GPU handler `onGpu`, answering every `Tensor` operation over `GpuBuffer<float>` on Metal or OpenCL |
| `model/tokenizer.osp` | `gpt2::Tokenizer` — the byte-level BPE tokenizer from `merges.txt` alone |
| `model/model.osp` | `gpt2::Model` — the configuration, the weights, LayerNorm, attention, the blocks, the forward pass keeping every block's activations, and the KV-cached step; every tensor is `B` |
| `generate/` | A one-file project: greedy sampling over the forward pass and the cached steps; `main` picks the handler |
| `train/` | A one-file project: the loss, the backward pass and one SGD step over the same forward pass; `main` picks the handler |
| [`python/`](python/README.md) | The same program in the Python flavor (Osprey, not Python), self-contained: the `Tensor` effect and both handlers, the tokenizer, the loader, the model and the sampler as `.py` modules, printing the Default generator's output byte for byte |
| `reference.py` | Real Python, run by CPython: the independent oracle both Osprey programs are checked against |

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

`onLists` in `cpu.osp` answers `build` with `range |> map |> toList` and `matMul` with a
`fold` over `k`; `onGpu` in `gpu.osp` answers them with `gpuIota |> gpuMap` and `gpuMatMul`,
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

| Measured 2026-10-17, Apple M4 Max, macOS, release build, `--memory=arc`; `BACKEND=gpu` median of three, `BACKEND=lists` one run | `BACKEND=lists`, `OSPREY_GPU=invalid` | `BACKEND=gpu`, `GPU=host` | `BACKEND=gpu`, `GPU=metal` |
| --- | --- | --- | --- |
| 20 tokens from the 5-token prompt, Default flavor | 2 min 23 s, 1.8 GB peak | 8.2 s, 1.6 GB peak | 6.1 s, 2.8 GB peak |
| The same 20 tokens, Python flavor | 2 min 19 s, 1.8 GB peak | | 7.3 s, 2.8 GB peak |
| One training step over the 5-token prompt | 3 min 19 s, 3.9 GB peak | | 6.7 s, 3.9 GB peak |

The list handler ran under an `OSPREY_GPU` no backend accepts, so any GPU call would have aborted the run; it printed the same text as Metal in both flavors. Most of its gap is the representation rather than the device: the GPU-surface handler on the CPU (`GPU=host`) is 17 times faster than the list handler, and Metal is 1.3 times faster again for generation.

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

- **Weights** come through the `load` operation: `tensor("h.3.ln_2.bias")` is a `B`. A failed load aborts with its error instead of substituting an empty tensor. Under `onGpu`, `load` uses `gpuLoadTensor`; under `onLists`, it uses `loadTensor`, which reads the same file with the same reader straight into a list, so the list handler never touches the GPU surface. Shapes are stated at each `matMul`, because tensors are rank-1.
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
  `train` (prompt, backend); an unknown backend name prints an error and runs nothing.
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
