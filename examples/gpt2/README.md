# GPT-2 in Osprey

The trained GPT-2 model, exactly as Hugging Face publishes it, running on
Osprey's GPU surface. No training, no conversion step, no Python at run time:
`gpuLoadTensor` reads the published `model.safetensors`, the byte-level BPE
tokenizer is rebuilt from the published `merges.txt`, and the transformer is
written as `gpu*` kernels and `gpuMatMul` products.

```
make gpt2 PROMPT="The capital of France is" TOKENS=20
```

The first run downloads the two files (about 550 MB) into
`examples/gpt2/weights/`, which is ignored by git. Every run prints the
device, the prompt's token ids, and the greedy continuation:

```
gpt2 (124M) on host — 50000 merges, 50257 vocabulary rows, 12 layers
prompt: The capital of France is
tokens: 464 3139 286 4881 318
generated 20 tokens
The capital of France is the capital of the French Republic, and the capital of the French Republic is the capital of the French
```

## What runs where

Everything runs on the CPU today. `gpuDevice()` reports `host`, and the
program says so on its first line. The GPU surface is the shape the model is
written in, so the same source runs unchanged when a device backend lands
([plan 0023](../../docs/plans/0023-gpu-computation.md)); nothing here is
GPU-accelerated yet.

| Measured 2026-10-12, AMD Ryzen 9 3900X, Linux 6.18 (WSL2), release build | |
| --- | --- |
| Loading 124M parameters as float buffers plus the transposed embedding | 9.3 s |
| One greedy token, 5 to 25 tokens of context, no KV cache | about 1.25 s |
| 20 tokens, default memory backend (nothing reclaimed) | 34 s, 10.0 GB peak |
| 20 tokens, `--memory=arc` | 42 s, 1.3 GB peak |
| `reference.py`, the same 20 tokens in pure Python | 16 min 25 s |

`make gpt2` runs under `--memory=arc` because of that memory row.

## Checking the output

[`reference.py`](reference.py) is an independent pure-Python implementation
(no numpy) that reads the same two files and runs the same forward pass with
the same left-to-right accumulation order. Run both on one prompt and compare
the token ids and the text:

```
python3 examples/gpt2/reference.py examples/gpt2/weights "The capital of France is" 20
```

On the prompts tried so far — the default, `Hello, I'm a language model,`,
`  two  spaces 123abc don't` and `naïve café — 日本` — the tokenizer produces
identical ids in both implementations, and the greedy continuations compared
so far are identical: all 20 tokens of the default prompt, and the two of
`Hello, I'm a language model,` (`not a`).

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

## Known limits

- One prompt line and an optional token count line on stdin; no KV cache, so
  every step recomputes the whole sequence.
- The pre-tokenizer is an ASCII approximation of GPT-2's Unicode regular
  expression; prompts in other scripts tokenize correctly only where the
  letter/number/punctuation classes agree with the ASCII rules.
- Computation is double precision throughout (the fp32 weights widen exactly),
  so logits differ from an fp32 implementation in the last digits; the greedy
  choice has agreed on every token compared so far.
