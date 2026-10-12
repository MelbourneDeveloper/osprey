# GPT-2 in Osprey

The trained GPT-2 model, exactly as Hugging Face publishes it, running on
Osprey's GPU surface. No training from scratch, no conversion step, no Python at
run time: `gpuLoadTensor` reads the published `model.safetensors`, the byte-level
BPE tokenizer is rebuilt from the published `merges.txt`, and the transformer is
written as `gpu*` kernels and `gpuMatMul` products.

```
make gpt2 PROMPT="The capital of France is" TOKENS=20
make gpt2-train PROMPT="The capital of France is"
```

The first run downloads the two files (about 550 MB) into
`examples/gpt2/weights/`, which is ignored by git. `make gpt2` prints the
device, the prompt's token ids, and the greedy continuation:

```
gpt2 (124M) on host — 50000 merges, 50257 vocabulary rows, 12 layers
prompt: The capital of France is
tokens: 464 3139 286 4881 318
generated 20 tokens
The capital of France is the capital of the French Republic, and the capital of the French Republic is the capital of the French
```

`make gpt2-train` runs one training step on the same downloaded weights: the
next-token cross-entropy of the prompt, the gradient of every parameter through
the twelve blocks, one SGD step, and the loss again. It prints the loss, the sum
of every gradient tensor (one line per block, in the checkpoint's tensor order),
and the loss after the step:

```
gpt2 (124M) on host — one training step over 5 tokens, learning rate 0.0001
prompt: The capital of France is
tokens: 464 3139 286 4881 318
loss 4.5031892154473265
grad wte 3.580909711548702e-12 wpe 4.3062341886779265e-15 ln_f 0.20932690619903307 0.5128505627737475
h.0 -3.3628631026292846 -7.137586885391019 3.8554208861589476 …
…
loss after one step 3.8539085307403833
```

## Layout

| Path | What it is |
| --- | --- |
| `model/tokenizer.osp` | `gpt2::Tokenizer` — the byte-level BPE tokenizer from `merges.txt` alone |
| `model/model.osp` | `gpt2::Model` — the configuration, the weights, LayerNorm, attention, the blocks and the forward pass, keeping every block's activations |
| `generate/` | A one-file project: greedy sampling over the forward pass |
| `train/` | A one-file project: the loss, the backward pass and one SGD step over the same forward pass |
| `reference.py` | The independent pure-Python oracle for both programs |

A standalone script cannot import across files, so each program is its own
project naming `model/` as a second source root
([modules](../../docs/specs/0025-ModulesAndNamespaces.md)), the layout the
graphics examples use. `osprey examples/gpt2/generate --run` runs a project
directory; the Makefile targets add `--memory=arc`.

## What runs where

Everything runs on the CPU today. `gpuDevice()` reports `host`, and both
programs say so on their first line. The GPU surface is the shape the model is
written in, so the same source runs unchanged when a device backend lands
([plan 0023](../../docs/plans/0023-gpu-computation.md)); nothing here is
GPU-accelerated.

| Measured 2026-10-12, AMD Ryzen 9 3900X, Linux 6.18 (WSL2), release build | |
| --- | --- |
| Loading 124M parameters as float buffers plus the transposed embedding | 9.3 s |
| One greedy token, 5 to 25 tokens of context, no KV cache | about 1.25 s |
| 20 tokens, default memory backend (nothing reclaimed) | 34 s, 10.0 GB peak |
| 20 tokens, `--memory=arc` | 42 s, 1.3 GB peak |
| `reference.py`, the same 20 tokens in pure Python | 16 min 25 s |
| One training step over the 5-token prompt (load, forward, loss, backward, SGD, forward again), `--memory=arc` | 42.2 s, 3.6 GB peak |
| `reference.py --train`, the same step in pure Python | 5 min 49 s, 16.2 GB peak |

The Makefile targets run under `--memory=arc` because of the memory rows.

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
the step print identically in both implementations.

The same kernels — LayerNorm, GELU, softmax, causal attention, `linear`, the
block and the argmax — are pinned at width 4 against values the Python oracle
computed in
[`tests/core/gpu/transformer.test.osp`](../../tests/core/gpu/transformer.test.osp),
which the differential harness runs under every memory backend, both kernel
lowerings and wasm32 on every `make test`.

## How it is written

- **Weights** come straight from the safetensors header: `tensor("h.3.ln_2.bias")`
  is a `GpuBuffer<float>`, or an `Error` naming the tensor and the reason.
  Shapes are stated at each `gpuMatMul`, because buffers are rank-1 today.
- **Tokenizer** — GPT-2 encodes bytes, and merge *i* of `merges.txt` is token
  `256 + i`, so the merge list alone is the vocabulary; `vocab.json` is never
  read. Pre-tokenization approximates GPT-2's regular expression for ASCII
  text and treats every non-ASCII byte as a letter.
- **Attention, LayerNorm, GELU and softmax** are `gpuIota |> gpuMap` kernels
  over flat index arithmetic, with `exp`, `sqrt` and `tanh` called inside
  kernels and `fold` over `range` for the inner reductions.
- **Products** are `gpuMatMul`, whose host implementation accumulates each
  element in k order from 0.0 — the reference order a device backend must match.
- **Sampling** is greedy: a `gpuFold` for the peak logit, a `gpuFilter` for
  its index.
- **The backward pass** is the same kinds of kernel. A `linear`'s gradients are
  two products against transposed operands plus a column sum; LayerNorm's
  recomputes the normalised rows from its input; attention's is the softmax
  Jacobian per row and then the query, key and value columns of each head.
  The forward pass keeps each block's activations in an `Acts` record so the
  backward pass reads each exactly once.

## Known limits

- One prompt line (and, for `generate`, an optional token count line) on
  stdin; no KV cache, so every generated token recomputes the whole sequence.
- The pre-tokenizer is an ASCII approximation of GPT-2's Unicode regular
  expression; prompts in other scripts tokenize correctly only where the
  letter/number/punctuation classes agree with the ASCII rules.
- Computation is double precision throughout (the fp32 weights widen exactly),
  so logits differ from an fp32 implementation in the last digits; the greedy
  choice has agreed on every token compared so far.
- Training is one plain SGD step at a fixed learning rate: no optimiser state,
  no batching, no weight decay, and the updated weights are not written back.
