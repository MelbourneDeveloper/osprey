---
layout: page.njk
title: "GPT-2, Written Once, Run Two Ways"
excerpt: "The trained GPT-2 model is one Osprey program written against a static Tensor effect. One handler runs it over plain lists, another on the GPU surface, and the compiler discharges the choice before the program starts. Same tokens, same loss, same 148 gradient sums."
description: "How Osprey's static effects let the published GPT-2 124M weights run unchanged under two interpretations — host lists and the GPU surface's typed buffers and pure kernels — with byte-identical output against an independent reference, and what the port fixed in the compiler."
tags: ["blog", "gpu", "gpt2", "machine-learning", "algebraic-effects", "compiler"]
author: "Christian Findlay"
readingTime: 9
---

Every numeric codebase that runs on more than one kind of hardware has the same shape: a reference implementation that is slow and trusted, a fast implementation that is trusted because it agrees with the slow one, and a growing list of places where they disagree. The two are written twice. They drift. The test that compares them is the one that gets skipped when the fast path changes.

Osprey's answer is to write the program once and make the choice of implementation an algebraic effect. This post is about doing that to GPT-2, the 124M-parameter model exactly as Hugging Face publishes it: load the weights, generate text, run a training step. One source. Two handlers. The output is identical to the byte.

```bash
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=gpu
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=lists
```

```text
gpt2 (124M) on host — 50000 merges, 50257 vocabulary rows, 12 layers
prompt: The capital of France is
tokens: 464 3139 286 4881 318
generated 20 tokens
The capital of France is the capital of the French Republic, and the capital of the French Republic is the capital of the French
```

The second run prints `gpt2 (124M) on lists` on its first line and the same four lines after it. The first run downloads `model.safetensors` and `merges.txt` into `examples/gpt2/weights/`. There is no conversion step and no Python at run time.

## The effect is the interface

A tensor program needs nine operations. Build a buffer from an index function. Read an element. Measure a buffer. Combine two element-wise. Reduce one. Multiply two matrices. Lift a host list. Load a tensor from a checkpoint. Say who is answering. That list is the effect, and `B` is whatever the handler wants a tensor to be:

```osprey
export static effect Tensor<B> {
    build: fn(int, fn(int) -> float) -> B
    get: fn(B, int) -> float
    length: fn(B) -> int
    zipWith: fn(B, B, fn(float, float) -> float) -> B
    reduce: fn(B, float, fn(float, float) -> float) -> float
    matMul: fn(B, B, int, int, int) -> B
    fromList: fn(List<float>) -> B
    load: fn(string, string) -> B
    device: fn() -> string
}
```

The model performs these and nothing else. LayerNorm over the `t` rows of a `[t, width]` tensor is three `build`s:

```osprey
fn layerNorm(x, t, g, b) = {
    let means = rowMeans(x, t)
    let variances = rowVariances(x, t, means)
    perform Backends::Tensor.build(t * width, fn(i) => normalised(x, means, variances, g, b, i))
}
```

A linear layer is one `matMul` and one bias gather. The weights are a record whose twelve fields are `B`. Attention, GELU, softmax, the backward pass, the SGD step: all of it is `perform`, index arithmetic and `fold`. Nowhere does the model say "list" or "buffer". Hindley-Milner inference carries `B` through every record and every helper, so there is nothing to annotate and no second copy to keep in step.

## Two handlers

The first handler answers with plain lists. It is the slow, obvious implementation, the one you would write to be sure.

```osprey
export fn onLists(work) = {
    handle static Tensor<List<float>> {
        build n f => range(0, n) |> map(f) |> toList()
        get b i => listGet(b, i) ?: 0.0
        length b => listLength(b)
        zipWith a b f => range(0, listLength(a)) |> map(fn(i) => f(listGet(a, i) ?: 0.0, listGet(b, i) ?: 0.0)) |> toList()
        reduce b seed f => range(0, listLength(b)) |> fold(seed, fn(a, i) => f(a, listGet(b, i) ?: 0.0))
        matMul a b rows inner cols => range(0, rows * cols) |> map(fn(i) => range(0, inner) |> fold(0.0, fn(acc, k) => acc + (listGet(a, intDiv(i, cols) * inner + k) ?: 0.0) * (listGet(b, k * cols + i % cols) ?: 0.0))) |> toList()
        fromList xs => xs
        load path name => fromGpu(gpuLoadTensor(path, name) ?: toGpu([]))
        device => "lists"
    }
    work()
}
```

The second answers with the GPU surface: `GpuBuffer<float>`, kernels that are pure functions, `gpuMatMul`, `gpuLoadTensor`.

```osprey
export fn onGpu(work) = {
    handle static Tensor<GpuBuffer<float>> {
        build n f => gpuIota(n) |> gpuMap(f)
        get b i => gpuGet(b, i) ?: 0.0
        length b => gpuLength(b)
        zipWith a b f => gpuZipWith(a, b, f)
        reduce b seed f => b |> gpuFold(seed, f)
        matMul a b rows inner cols => gpuMatMul(a, b, rows, inner, cols) ?: toGpu([])
        fromList xs => toGpu(xs)
        load path name => gpuLoadTensor(path, name) ?: toGpu([])
        device => gpuDevice()
    }
    work()
}
```

Notice what the arms do with `f`. The model's callback is forwarded straight into `gpuMap`, `gpuZipWith` and `gpuFold`. The compiler proves that callback pure before it accepts the program; a kernel that performs a runtime effect is rejected at compile time with the operation's name in the message.

And `main` is the one place the choice is made:

```osprey
match input() {
    "lists" => Backends::onLists(fn() => run(prompt, count))
    _ => Backends::onGpu(fn() => run(prompt, count))
}
```

## Discharged, not dispatched

`static` is the word that makes this different from a trait, an interface or a function table. A static effect is answered at compile time. The compiler specialises `run` and everything it reaches once per handler region and inlines the arms at every `perform`. The compiled program contains two copies of the model and no handler at all: `Tensor` does not exist at run time, and there is no indirect call where it used to be. The `match` in `main` is choosing between two programs that were both checked, both specialised and both complete before the first byte of the weights was read.

That staging is also what makes the second handler honest. A static handler must answer every operation, so adding a tenth operation to `Tensor` fails to compile in whichever handler forgot it. `Tensor<List<float>>` and `Tensor<GpuBuffer<float>>` are different effects to the handler search, so a request can never be answered by the wrong representation. And the declared types travel with the rewrite: a kernel the model writes as `fn(g, x) => g` is checked against `fn(float, float) -> float` and keeps that type after discharge, which is how the backend knows to lower it over doubles.

## The same numbers, three ways

`examples/gpt2/reference.py` is an independent pure-Python forward and backward pass with no numpy, over the same two files, using the same left-to-right accumulation order the Osprey handlers do. That order is the contract: `matMul` sums over `k` from `0.0`, in double, each product formed before it is added, no fused multiply-add. Both handlers follow it, which is why they agree with each other and with the reference to the last bit.

```bash
make gpt2-train PROMPT="The capital of France is" BACKEND=gpu
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

`BACKEND=lists` prints the same loss, the same sum for every one of the 148 gradient tensors and the same loss after the step. Several of those sums are near `1e-16`, where any difference in accumulation order shows. The reference prints them too.

| Measured 2026-10-13, AMD Ryzen 9 3900X, Linux 6.18 (WSL2), release build, `--memory=arc` | `BACKEND=gpu` | `BACKEND=lists` |
| --- | --- | --- |
| 20 greedy tokens from the 5-token prompt, with the KV cache | 19.6 s, 1.4 GB peak | 7 min 18 s, 1.7 GB peak |
| One training step over the 5-token prompt | 52 s, 3.6 GB peak | 9 min, 4.6 GB peak |

The GPU surface is faster because its buffers are dense doubles and its kernels are lifted into flat functions; the list handler boxes every element and calls through a closure. Neither number is a GPU result. `gpuDevice()` reports `host`, the model prints it, and this post does not describe anything here as accelerated.

## What the port fixed

A model is a better review of a compiler than a test suite is. Landing the two-handler version found three defects, each pinned by a failing test before it was fixed.

**Hygiene renamed a payload binder.** The static rewrite renames bound names so inlined arms cannot capture the model's variables. It also renamed the binder behind `?:`, which the checker types by name, so `gpuMatMul(…) ?: toGpu([])` lost its element type and a float buffer was read back as integers. Result field labels now keep their spelling and shadow by scope.

**Inference tables collided across specialisations.** The checker publishes each lambda's and `let`'s resolved type keyed by source position. Two specialised copies of one function share every position, so the list copy's types overwrote the GPU copy's. The tables are now keyed by the enclosing function as well, and code generation names the function it is lowering, including while inlining a generic body.

**Discharge dropped the declared types.** The rewrite bound each operation argument with a plain `let`, so a kernel inside a generic helper was re-inferred as `fn(a, b) -> a` and lowered over raw words. The training step's `ln_2` and `mlp.c_fc` gradients came out as exact zeros. The bound arguments now carry the operation's declared parameter types, and the fix is a specification rule with a corpus case in both syntax flavours.

Earlier ports found the codegen defects the previous version of this post listed: a `fold` seeded with a list literal rebuilt the list every iteration, a string list was copied on every `listGet`, and `join` crashed on runtime-built lists. Those are fixed with regression programs in the corpus.

## What this is not

Nothing here runs on a GPU. The GPU surface is the shape one of the two handlers is written in, so the same model runs unchanged when a device backend lands; that backend is stage 4 of [plan 0023](https://github.com/Nimblesite/osprey/blob/main/docs/plans/0023-gpu-computation.md), and its gate is this program against `llm.c`'s CUDA path on the same GPU, within a factor of two. When it lands, `onGpu` is the one file that notices, and `gpuDevice()` changes its answer.

The handler is chosen where it is written. `match input()` selects between two regions the compiler has already specialised; a handler computed at run time would need a dynamic effect and a runtime handler. The pre-tokenizer approximates GPT-2's regular expression for ASCII text. One SGD step at a fixed learning rate is not an optimiser, and the updated weights are not written back.

To run it, `make gpt2` after [installing Osprey](/docs/installation/). [One model, two backends](/docs/one-model-two-backends/) walks through the effect, both handlers and what the compiler proves about them, with the same code in the ML flavour. The [example's README](https://github.com/Nimblesite/osprey/blob/main/examples/gpt2/README.md) has the measurements, the comparison procedure and the source layout, and the [GPU specification](/spec/0034-gpucomputation/) and [effects specification](/spec/0017-algebraiceffects/) hold the contracts every number above is held to.
