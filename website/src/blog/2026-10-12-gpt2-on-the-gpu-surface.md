---
layout: page.njk
title: "GPT-2, Written Once, Run Two Ways"
excerpt: "The trained GPT-2 model is one Osprey program written against a static Tensor effect. One handler runs it over plain lists, another on the GPU surface, whose kernels and matrix products execute on Metal or OpenCL. The compiler discharges the choice before the program starts."
description: "How Osprey's static effects let the published GPT-2 124M weights run unchanged under two interpretations — host lists and the GPU surface's typed buffers and pure kernels, on the CPU or on a Metal or OpenCL device — what agrees to the byte, what agrees to twelve digits, and what the port fixed in the compiler."
modified: 2026-10-16
tags: ["blog", "gpu", "gpt2", "machine-learning", "algebraic-effects", "compiler"]
author: "Christian Findlay"
readingTime: 11
---

Every numeric codebase that runs on more than one kind of hardware has the same shape: a reference implementation that is slow and trusted, a fast implementation that is trusted because it agrees with the slow one, and a growing list of places where they disagree. The two are written twice. They drift. The test that compares them is the one that gets skipped when the fast path changes.

Osprey's answer is to write the program once and make the choice of implementation an algebraic effect. This post is about doing that to GPT-2, the 124M-parameter model exactly as Hugging Face publishes it: load the weights, generate text, run a training step. One source. Two handlers. On the CPU the two print identical output, byte for byte, and so does an independent Python reference. On a GPU, through Metal or OpenCL, the same source generates the same text, and a training step runs about three times faster than on the CPU.

*Updated 2026-10-16. When this post first went up, nothing in it ran on a GPU. The GPU-surface handler now runs every kernel and every matrix product on Metal or OpenCL, so the sections on running it, the numbers and the limits are rewritten, and every figure is re-measured.*

```bash
make gpt2 PROMPT="The capital of France is" TOKENS=20                # GPU=auto: a device when there is one
make gpt2 PROMPT="The capital of France is" TOKENS=20 GPU=host       # the same handler on the CPU
make gpt2 PROMPT="The capital of France is" TOKENS=20 BACKEND=lists  # the list handler
```

```text
gpt2 (124M) on metal:0 (Apple M4 Max) — 50000 merges, 50257 vocabulary rows, 12 layers
prompt: The capital of France is
tokens: 464 3139 286 4881 318
generated 20 tokens
The capital of France is the capital of the French Republic, and the capital of the French Republic is the capital of the French
```

The first line names who answered. `GPU=opencl` prints `on opencl:0 (Apple M4 Max)`, `GPU=host` prints `on host` and `BACKEND=lists` prints `on lists`; all four print the same four lines after it. The first run downloads `model.safetensors` and `merges.txt` into `examples/gpt2/weights/`. There is no conversion step and no Python at run time.

## The effect is the interface

A tensor program needs ten operations. Build a buffer from an index function. Read an element. Measure a buffer. Combine two element-wise. Reduce one. Find the first index holding a value. Multiply two matrices. Lift a host list. Load a tensor from a checkpoint. Say who is answering. That list is the effect, and `B` is whatever the handler wants a tensor to be:

```osprey
export type Kernel = fn(int) -> float
export type Combine = fn(float, float) -> float

export static effect Tensor<B> {
    build: fn(int, Kernel) -> B
    get: fn(B, int) -> float
    length: fn(B) -> int
    zipWith: fn(B, B, Combine) -> B
    reduce: fn(B, float, Combine) -> float
    indexOf: fn(B, float) -> int
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
        indexOf b x => range(0, listLength(b)) |> map(fn(i) => hit(listGet(b, i) ?: 0.0, x, i, listLength(b))) |> fold(listLength(b), minI)
        matMul a b rows inner cols => range(0, rows * cols) |> map(fn(i) => range(0, inner) |> fold(0.0, fn(acc, k) => acc + (listGet(a, intDiv(i, cols) * inner + k) ?: 0.0) * (listGet(b, k * cols + i % cols) ?: 0.0))) |> toList()
        fromList xs => xs
        load path name => fromGpu(orAbort(gpuLoadTensor(path, name)))
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
        indexOf b x => gpuIota(gpuLength(b)) |> gpuMap(fn(i) => hit(gpuGet(b, i) ?: 0.0, x, i, gpuLength(b))) |> gpuFold(gpuLength(b), minI)
        matMul a b rows inner cols => orAbort(gpuMatMul(a, b, rows, inner, cols))
        fromList xs => toGpu(xs)
        load path name => orAbort(gpuLoadTensor(path, name))
        device => gpuDevice()
    }
    work()
}
```

Notice what the arms do with `f`. The model's callback is forwarded straight into `gpuMap`, `gpuZipWith` and `gpuFold`. The compiler proves that callback pure before it accepts the program; a kernel that performs a runtime effect is rejected at compile time with the operation's name in the message. Once it is proven pure, the compiler lifts it into a kernel and compiles it to the device dialect that Metal and OpenCL both run. `indexOf`, the step that picks the next token, is itself a map and a fold, so even the sampler's argmax runs on the device.

`indexOf`'s two helpers name an element's index when it holds the value, and keep the smaller of two indices. `orAbort` turns a failed load or product into an error that ends the run and names the cause. An earlier version substituted an empty tensor, and the model generated text from fabricated zeros.

The choice is made in one place, `Backends::under`, from the name `main` reads:

```osprey
export fn under(name, work) = match name {
    "gpu" => onGpu(work)
    "lists" => onLists(work)
    _ => print("unknown backend `${name}`: choose gpu or lists")
}
```

A mistyped name runs nothing rather than silently landing on either handler. `main` wraps the call in `Backends::strictly`, which handles `Arith`: an integer overflow or a division by zero in the model ends the run with a message, as a device kernel's fault does, instead of becoming a value.

## Discharged, not dispatched

`static` is the word that makes this different from a trait, an interface or a function table. A static effect is answered at compile time. The compiler specialises `run` and everything it reaches once per handler region and inlines the arms at every `perform`. The compiled program contains two copies of the model and no handler at all: `Tensor` does not exist at run time, and there is no indirect call where it used to be. The `match` in `under` is choosing between two programs that were both checked, both specialised and both complete before the first byte of the weights was read.

That staging is also what makes the second handler honest. A static handler must answer every operation, so adding an eleventh operation to `Tensor` fails to compile in whichever handler forgot it. `Tensor<List<float>>` and `Tensor<GpuBuffer<float>>` are different effects to the handler search, so a request can never be answered by the wrong representation. And the declared types travel with the rewrite: a kernel the model writes as `fn(g, x) => g` is checked against `fn(float, float) -> float` and keeps that type after discharge, which is how the backend knows to lower it over doubles.

## The same numbers, three ways

`examples/gpt2/reference.py` is an independent pure-Python forward and backward pass with no numpy, over the same two files, using the same left-to-right accumulation order the Osprey handlers do. That order is the contract: `matMul` sums over `k` from `0.0`, in double, each product formed before it is added, no fused multiply-add. Both handlers follow it, which is why, on the CPU, they agree with each other and with the reference to the last bit.

```bash
make gpt2-train PROMPT="The capital of France is" GPU=host
```

```text
gpt2 (124M) on host — one training step over 5 tokens, learning rate 0.0001
prompt: The capital of France is
tokens: 464 3139 286 4881 318
loss 4.503189215447316
grad wte -6.369071936518367e-13 wpe -1.4510961876545991e-15 ln_f 0.20932690619905386 0.512850562773751
h.0 -3.3628631026293108 -7.1375868853910704 3.8554208861589614 …
…
loss after one step 3.8539085307403895
```

`BACKEND=lists` prints the same loss, the same checksum for every one of the 148 gradient tensors and the same loss after the step. A matrix's checksum adds each row in order and then adds the row sums, the same two-level sum under either handler. Several of those sums are near `1e-16`, where any difference in accumulation order shows. The reference prints them too.

| Measured 2026-10-13, AMD Ryzen 9 3900X, Linux 6.18 (WSL2), release build, `--memory=arc`, `GPU=host` | `BACKEND=gpu` | `BACKEND=lists` |
| --- | --- | --- |
| 20 greedy tokens from the 5-token prompt, with the KV cache | 19.6 s, 1.4 GB peak | 7 min 18 s, 1.7 GB peak |
| One training step over the 5-token prompt | 52 s, 3.6 GB peak | 9 min, 4.6 GB peak |

Both columns ran on the CPU. The GPU-surface handler is faster there because its buffers are dense doubles and its kernels are lifted into flat functions; the list handler boxes every element and calls through a closure.

## On the device

`GPU=metal` and `GPU=opencl` run the same handler with every kernel and every matrix product on the GPU: embeddings, LayerNorm, attention scores and softmax, GELU, the residual adds, every reduction, the next token's index and the loss. The CPU tokenizes, builds and loads buffers, sequences the launches, and reads back only what the program reads on the host: the next token's index, the loss and the printed checksums. A buffer is uploaded on its first launch and stays on the device after it, so the weights go up once and a layer's intermediate buffers never come back. Under `--memory=arc`, which the make targets use, dropping a buffer releases its device copy. A kernel the device cannot compile stops the run and names itself; it never falls back to the CPU.

```text
gpt2 (124M) on metal:0 (Apple M4 Max) — one training step over 5 tokens, learning rate 0.0001
prompt: The capital of France is
tokens: 464 3139 286 4881 318
loss 4.503189215447317
grad wte -7.653704059418232e-13 wpe 5.182486384480711e-17 ln_f 0.20932690619897562 0.5128505627737499
h.0 -3.362863102629283 -7.137586885391083 3.8554208861588597 …
…
loss after one step 3.8539085307403917
```

That is not the host's output, and the difference has one cause. Metal has no double-precision type, so the device library does binary64 arithmetic in software, and the OpenCL adapter runs the same library. Its add, subtract, multiply, divide, square root, conversions and comparisons match the host's IEEE results bit for bit; `gpu_shader_tests.c` checks that over twelve million random words, subnormals and special values. `exp`, `ln` and `tanh` are different. The device evaluates them in 128-bit fixed point and rounds once, so it returns the correctly rounded double. A host libm does not on every argument: Apple's `tanh` is one ulp off on about one sampled argument in 1 400, and glibc's by up to two. On aarch64 Linux the device matches glibc's binary128 `long double` functions on all 800 000 sampled arguments of each. GELU calls `tanh`, softmax calls `exp` and the loss calls `ln`, so every number downstream inherits a last-digit difference.

The generated text is identical on every path. The loss agrees to fifteen significant digits. Across all 150 numbers the step prints, the largest difference from the host is 2.4 × 10⁻¹³, and no number larger than 10⁻⁶ differs by more than 4 parts in 10¹³. A sum that cancels to around `1e-15`, such as `wpe`, is made of that noise and differs in every digit. Metal and OpenCL print identical bytes apart from the device name.

| Measured 2026-10-16, Apple M4 Max, macOS, release build, `--memory=arc`, `BACKEND=gpu`, median of three | `GPU=host` | `GPU=metal` | `GPU=opencl` |
| --- | --- | --- | --- |
| 20 greedy tokens from the 5-token prompt, with the KV cache | 6.5 s, 1.7 GB peak | 6.9 s, 3.1 GB peak | 6.8 s, 1.7 GB peak |
| One training step over the 5-token prompt | 25.2 s, 4.0 GB peak | 7.6 s, 5.9 GB peak | 7.8 s, 3.1 GB peak |

These are whole-program wall-clock times, including reading 548 MB of weights. Training is about three times faster on the GPU. Its 1 499 dispatches are mostly large: half run over 3 072 threads or more, and a tenth over more than half a million. Generation is not faster. After the prompt, the KV cache makes each new token a one-row pass, so 4 338 of its 5 919 dispatches run over 768 threads or fewer. With no kernel fusion, each of those small kernels is a separate launch.

## What the port fixed

A model is a better review of a compiler than a test suite is. Landing the two-handler version found three defects, each pinned by a failing test before it was fixed.

**Hygiene renamed a payload binder.** The static rewrite renames bound names so inlined arms cannot capture the model's variables. It also renamed the binder behind `?:`, which the checker types by name, so `gpuMatMul(…) ?: toGpu([])` lost its element type and a float buffer was read back as integers. Result field labels now keep their spelling and shadow by scope.

**Inference tables collided across specialisations.** The checker publishes each lambda's and `let`'s resolved type keyed by source position. Two specialised copies of one function share every position, so the list copy's types overwrote the GPU copy's. The tables are now keyed by the enclosing function as well, and code generation names the function it is lowering, including while inlining a generic body.

**Discharge dropped the declared types.** The rewrite bound each operation argument with a plain `let`, so a kernel inside a generic helper was re-inferred as `fn(a, b) -> a` and lowered over raw words. The training step's `ln_2` and `mlp.c_fc` gradients came out as exact zeros. The bound arguments now carry the operation's declared parameter types, and the fix is a specification rule with a corpus case in both syntax flavours.

Earlier ports found the codegen defects the previous version of this post listed: a `fold` seeded with a list literal rebuilt the list every iteration, a string list was copied on every `listGet`, and `join` crashed on runtime-built lists. Those are fixed with regression programs in the corpus.

The device port found two more. The first checked-multiply helper used signed 64-bit division, which crashed Metal's shader compiler (`XPC_ERROR_CONNECTION_INTERRUPTED`) and made OpenCL refuse to create the kernel. The library now divides unsigned magnitudes only. The second was in the gradient checksums. A device `reduce` runs on one thread in element order, which is what keeps its result equal to the host's. Before buffers stayed on the device, a training step took 68.5 s on Metal, and 56 s of it went to checksums that folded up to 38.6 million elements on that one thread. Each checksum is now a per-row kernel and a `reduce` over the rows.

## What this is not

This is not a fast GPT-2. Every operation on the device is software binary64, so that the device computes the same values as the host; a float32 handler would be a third interpretation, checked against this one. There is no kernel fusion yet, so every kernel is its own dispatch, and a large reduction has to be written as a per-row kernel and a reduction over the rows. The gate for this stage in [plan 0023](https://github.com/Nimblesite/osprey/blob/main/docs/plans/0023-gpu-computation.md) is this program within a factor of two of a native GPT-2 on the same GPU at the same precision, and it has not been measured. The OpenCL adapter is vendor-neutral, but it has only run on Apple's OpenCL. AMD, Intel and NVIDIA hardware are untested.

The handler is chosen where it is written. `under`'s `match` selects between two regions the compiler has already specialised; a handler computed at run time would need a dynamic effect and a runtime handler. The pre-tokenizer approximates GPT-2's regular expression for ASCII text. One SGD step at a fixed learning rate is not an optimiser, and the updated weights are not written back.

To run it, `make gpt2` after [installing Osprey](/docs/installation/); add `GPU=host` to keep it on the CPU. [One model, two backends](/docs/one-model-two-backends/) walks through the effect, both handlers and what the compiler proves about them, with the same code in the ML flavour. The [example's README](https://github.com/Nimblesite/osprey/blob/main/examples/gpt2/README.md) has the measurements, the comparison procedure and the source layout, and the [GPU specification](/spec/0034-gpucomputation/) and [effects specification](/spec/0017-algebraiceffects/) hold the contracts every number above is held to.
