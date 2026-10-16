---
layout: page.njk
title: "GPT-2 Runs on Osprey’s GPU Backend"
excerpt: "A real language model, written in Osprey, generating text and running a training step on the GPU. The same model runs through Metal or OpenCL, with a measured 3.3× training speedup on an Apple M4 Max."
description: "How Osprey runs the published GPT-2 model on CPUs and GPUs, what Metal and OpenCL support today, and where performance is improving."
modified: 2026-10-16
tags: ["blog", "gpu", "gpt2", "machine-learning"]
author: "Christian Findlay"
readingTime: 3
---

Osprey can run **GPT-2, the 124-million-parameter language model, on a GPU**. It loads the published weights, generates text and runs a training step. The model itself is Osprey code: attention, matrix multiplication, predictions and weight updates.

On an Apple M4 Max, that training step takes **7.6 seconds on Metal versus 25.2 seconds on the CPU**: about **3.3× faster**. Metal and OpenCL run the same model source.

Osprey is running a real AI workload, with a CPU implementation to check its answers.

## Write the model once

Osprey provides typed GPU buffers and operations that apply ordinary Osprey functions to their contents. The compiler turns supported functions into GPU kernels—the small programs a GPU executes. You write the calculation in Osprey; you do not need a separate CUDA or Metal version of the model.

GPT-2 describes its calculations through a small set of array operations. One implementation uses CPU lists; another uses GPU buffers, executed through Metal or OpenCL. The model’s calculations stay the same. [The model and both implementations are available to read](https://github.com/Nimblesite/osprey/tree/main/examples/gpt2/model).

## What actually runs on the GPU?

The CPU turns your prompt into tokens and loads the weights. The GPU handles the model’s layers, attention calculations, predictions and training calculations.

Weights upload once. Intermediate results stay on the GPU for the next operation, avoiding a trip back to the CPU after every calculation. The CPU coordinates the work, handles text and reads back the results the program needs.

Generation reuses earlier calculations as it adds tokens. The training example updates the weights once and checks the prediction error again.

An independent Python implementation checks the calculations. Recorded CPU and GPU runs generated the same text, with closely agreeing training results. The [verification details](https://github.com/Nimblesite/osprey/blob/main/examples/gpt2/README.md) explain the small numerical differences.

## Metal, OpenCL and CUDA

Osprey’s GPU interface has no CUDA dependency.

| Execution path | Current support |
| --- | --- |
| CPU | GPU-buffer operations on the host, plus a list implementation for comparison |
| Metal | Implemented for Apple GPUs |
| OpenCL | Implemented; provides a route to compatible AMD, Intel and NVIDIA devices with suitable drivers |
| CUDA | No dedicated CUDA backend today |

These measurements used an Apple M4 Max through Metal and OpenCL, not every supported GPU vendor.

## Where acceleration helps today

| Workload | CPU | Metal | OpenCL |
| --- | ---: | ---: | ---: |
| One training step | 25.2 s | **7.6 s** | **7.8 s** |
| Generate 20 tokens | **6.5 s** | 6.9 s | 6.8 s |

These recorded measurements are from 16 October 2026: Apple M4 Max, release build, ARC memory mode, a five-token prompt and the median of three runs. Times include loading the weights.

Training gives the GPU larger pieces of work and gets a clear speedup. Short text generation still pays for many small GPU launches. Combining those jobs and adding faster numerical formats are further work; this version uses carefully checked 64-bit arithmetic.

## Run it

From the repository root, after building Osprey:

```bash
# Choose the backend available on your machine: metal, opencl or host.
make gpt2 GPU=metal PROMPT="The capital of France is" TOKENS=20
make gpt2-train GPU=metal PROMPT="The capital of France is"
```

The first run downloads the weights and tokenizer data. The program prints its actual backend and reports unsupported GPU work without silently substituting CPU execution.

[Try the complete example](https://github.com/Nimblesite/osprey/tree/main/examples/gpt2), and compare its CPU and GPU runs. The model and its GPU calculations are Osprey code you can read and change.
