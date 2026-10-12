---
layout: page.njk
title: "GPT-2 on the GPU Surface"
excerpt: "The trained GPT-2 model loads its published weights, generates text and runs a training step as Osprey GPU kernels. Every one of them executes on the CPU today, and that is the point."
description: "How Osprey's GPU surface — typed buffers, pure kernels, gpuMatMul and gpuLoadTensor — runs the published GPT-2 124M weights on the host backend, byte-for-byte against an independent reference, before any device backend exists."
tags: ["blog", "gpu", "gpt2", "machine-learning", "compiler"]
author: "Christian Findlay"
readingTime: 7
---

Osprey has a GPU surface: a `GpuBuffer<T>` type, kernels that are ordinary pure functions, and a dozen built-ins such as `gpuMap`, `gpuFold` and `gpuMatMul` that apply them. Nothing in it runs on a GPU yet. Every kernel executes as a loop on the CPU, and `gpuDevice()` reports `host`.

So before building a device backend, we asked a harder question: is the surface right? The test was to run a real model on it, unchanged, and compare every number against an implementation we did not write in Osprey. The model is GPT-2, the 124M-parameter one, exactly as Hugging Face publishes it.

```bash
make gpt2 PROMPT="The capital of France is" TOKENS=20
```

```text
gpt2 (124M) on host — 50000 merges, 50257 vocabulary rows, 12 layers
prompt: The capital of France is
tokens: 464 3139 286 4881 318
generated 20 tokens
The capital of France is the capital of the French Republic, and the capital of the French Republic is the capital of the French
```

The first run downloads `model.safetensors` and `merges.txt` into `examples/gpt2/weights/`. There is no conversion step and no Python at run time.

## What the surface looks like

A kernel is a function. `gpuIota(n)` is the buffer `0 .. n-1`; `gpuMap` applies a function to each element; `gpuGet` reads one element as a `Result`, so an out-of-range index is a value, not a crash. LayerNorm over the `t` rows of a `[t, width]` buffer is three kernels:

```osprey
fn layerNorm(x, t, g, b) = {
    let means = rowMeans(x, t)
    let variances = rowVariances(x, t, means)
    gpuIota(t * width) |> gpuMap(fn(i) => {
        let r = intDiv(i, width)
        let c = i % width
        (at(x, i) - at(means, r)) / sqrt(at(variances, r) + epsilon) * at(g, c) + at(b, c)
    })
}
```

`at` is `gpuGet(m, i) ?: 0.0`. Buffers are rank-1, so every kernel in the model is flat index arithmetic: `intDiv(i, width)` for the row, `i % width` for the column, `(h * t + r) * t + c` for an attention score. A matrix is a buffer plus the shape the program states:

```osprey
fn linear(x, rows, inner, cols, w, bias) = {
    let value = product(x, w, rows, inner, cols)
    gpuIota(rows * cols) |> gpuMap(fn(i) => at(value, i) + at(bias, i % cols))
}
```

`product` is `gpuMatMul` with its `Error` reported and turned into an empty buffer. The weights come straight from the safetensors header: `gpuLoadTensor(path, "h.3.ln_2.bias")` is a `GpuBuffer<float>`, or an `Error` naming the tensor and the reason. Four storage types widen exactly: F64, F32, F16 and BF16.

The whole model is twelve of these blocks, a byte-level BPE tokenizer rebuilt from the merge list alone, and an argmax. It is a shared model module and two one-file projects under `examples/gpt2/`, written in the same `fn`, `match`, `fold` and record syntax as any other Osprey program. The compiler proves each kernel pure before it accepts it: a kernel that performs an unhandled effect is a compile error, not a runtime surprise.

## The numbers are the same, not close

`examples/gpt2/reference.py` is an independent pure-Python forward pass, with no numpy, over the same two files. It uses the same left-to-right accumulation order the Osprey runtime does, and that order is the contract: `gpuMatMul` sums over `k` from `0.0`, in double, with each product formed before it is added, no fused multiply-add. The specification says so, and a device backend may reorder only where it reproduces the host's bytes.

On every prompt tried, the token ids agree exactly, and so do the greedy continuations. The reference takes 16 minutes 25 seconds for the twenty tokens above. Osprey's host backend takes 34 seconds, or 42 seconds under reference counting with a 1.3 GB peak instead of 10 GB. That is not a GPU result and we do not describe it as one. It is a CPU loop generated from a functional program, compared against pure Python.

## A training step

Generation only exercises the forward pass. The backward pass is where a surface shows its gaps, so `examples/gpt2/train` runs one step on the downloaded weights: the next-token cross-entropy of the prompt, the gradient of every parameter through the twelve blocks, one SGD step, and the loss again.

```bash
make gpt2-train PROMPT="The capital of France is"
```

```text
gpt2 (124M) on host — one training step over 5 tokens, learning rate 0.0001
prompt: The capital of France is
tokens: 464 3139 286 4881 318
loss 4.5031892154473265
grad wte 3.580909711548702e-12 wpe 4.3062341886779265e-15 ln_f 0.20932690619903307 0.5128505627737475
h.0 -3.3628631026292846 -7.137586885391019 3.8554208861589476 …
…
loss after one step 3.8539085307403833
```

The backward pass is written as the same kinds of kernel as the forward one. The gradient of `linear` is two products against transposed operands and a column sum; attention's is the softmax Jacobian per row and then the query, key and value columns of each head. The reference script gained the same step in the same order, and prints the same loss, the same sum for every one of the 148 gradient tensors, and the same loss after the step. Those sums are a stricter check than they look: several are near `1e-16`, where any difference in accumulation order shows.

The step runs in 42 seconds at a 3.6 GB peak under reference counting. It keeps every block's activations in a record, because the backward pass reads each of them once, and the forward pass that generation uses is the same function.

## What the port told us

Writing a model is a better review of a surface than writing tests for it. Three things came out.

**Codegen defects.** A `fold` whose seed is a list literal rebuilt the list on every iteration. A string list bound with `let` was copied on every `listGet`. `join` crashed on any list the runtime had built rather than one written as a literal. Each is now fixed, with a regression program in the corpus and an IR test that pins the lowering.

**Surface gaps.** Rank-1 indexing is the largest readability cost: every kernel indexes a `[heads, t, t]` tensor by hand. Nothing fuses: one block allocates nineteen buffers per step, and every `linear` is a product followed by a separate bias gather. Host-ordered reductions are harmless here, because the per-row folds are exactly the order the reference uses. Transfer cost is invisible on the host backend, which is precisely the problem a device backend will expose first.

**New built-ins.** `gpuMatMul`, `gpuLoad`, `gpuLoadTensor` and `sqrt`, `exp`, `ln`, `tanh` exist because the model needed them. Each landed with a specification section, C tests for every error sentence, IR tests and rejection fixtures in both syntax flavors.

## What this is not

Nothing here is GPU-accelerated. The GPU surface is the shape the model is written in, so the same source runs unchanged when a device backend lands; that backend is stage 4 of [plan 0023](https://github.com/Nimblesite/osprey/blob/main/docs/plans/0023-gpu-computation.md), and its gate is this program against `llm.c`'s CUDA path on the same GPU, within a factor of two. The pre-tokenizer approximates GPT-2's regular expression for ASCII text. There is no KV cache, so every generated token recomputes the whole sequence, and one SGD step is not an optimiser.

To run it, `make gpt2` after [installing Osprey](/docs/installation/). The [example's README](https://github.com/Nimblesite/osprey/blob/main/examples/gpt2/README.md) has the measurements, the comparison procedure and the source layout, and the [GPU specification](/spec/0034-gpucomputation/) has the contract every number above is held to.
