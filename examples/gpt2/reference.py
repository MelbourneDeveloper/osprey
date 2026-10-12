#!/usr/bin/env python3
"""Independent GPT-2 (124M) reference in plain Python — no numpy, no torch.

Reads the same `model.safetensors` and `merges.txt` the Osprey program reads,
runs the same fp32-widened-to-double forward pass with the same left-to-right
accumulation order, and prints the greedy continuation. It exists only to check
examples/gpt2/gpt2.osp: run both on one prompt and compare token ids and text.

    python3 examples/gpt2/reference.py examples/gpt2/weights "The capital of France is" 20
"""
import json
import math
import operator
import re
import struct
import sys

WIDTH, HEADS, HEAD_WIDTH, LAYERS, VOCAB, END = 768, 12, 64, 12, 50257, 50256
PATTERN = re.compile(
    r"""'s|'t|'re|'ve|'m|'ll|'d| ?[A-Za-z\u0080-￿]+| ?[0-9]+| ?[^\sA-Za-z0-9\u0080-￿]+|\s+(?!\S)|\s+"""
)


def bytes_to_unicode():
    bs = list(range(33, 127)) + list(range(161, 173)) + list(range(174, 256))
    cs = bs[:]
    n = 0
    for b in range(256):
        if b not in bs:
            bs.append(b)
            cs.append(256 + n)
            n += 1
    return dict(zip(bs, map(chr, cs)))


class Tokenizer:
    def __init__(self, merges_path):
        lines = open(merges_path, encoding="utf-8").read().split("\n")[1:]
        self.merges = [tuple(m.split(" ")) for m in lines if m]
        self.ranks = {m: i for i, m in enumerate(self.merges)}
        self.b2u = bytes_to_unicode()
        self.u2b = {u: b for b, u in self.b2u.items()}
        base = [self.b2u[b] for b in sorted(self.b2u, key=lambda b: ord(self.b2u[b]) if self.b2u[b] >= "Ā" else 0)]
        self.ids = {t: i for i, t in enumerate(self._base_order())}
        for i, (a, b) in enumerate(self.merges):
            self.ids[a + b] = 256 + i
        self.tokens = {i: t for t, i in self.ids.items()}

    def _base_order(self):
        order = list(range(33, 127)) + list(range(161, 173)) + list(range(174, 256))
        order += [b for b in range(256) if b not in order]
        return [self.b2u[b] for b in order]

    def bpe(self, word):
        word = tuple(word)
        while len(word) > 1:
            rank, at = min((self.ranks.get(p, 1 << 62), i) for i, p in enumerate(zip(word, word[1:])))
            if rank == 1 << 62:
                break
            a, b = word[at], word[at + 1]
            out, j = [], 0
            while j < len(word):
                if j < len(word) - 1 and word[j] == a and word[j + 1] == b:
                    out.append(a + b)
                    j += 2
                else:
                    out.append(word[j])
                    j += 1
            word = tuple(out)
        return word

    def encode(self, text):
        ids = []
        for piece in PATTERN.findall(text):
            stand_in = "".join(self.b2u[b] for b in piece.encode("utf-8"))
            ids.extend(self.ids[w] for w in self.bpe(stand_in))
        return ids

    def decode(self, ids):
        raw = bytes(self.u2b[c] for i in ids for c in self.tokens[i])
        return raw.decode("utf-8", errors="replace")


class Safetensors:
    def __init__(self, path):
        self.f = open(path, "rb")
        n = struct.unpack("<Q", self.f.read(8))[0]
        self.header = json.loads(self.f.read(n))
        self.start = 8 + n

    def tensor(self, name):
        e = self.header[name]
        assert e["dtype"] == "F32", name
        begin, end = e["data_offsets"]
        self.f.seek(self.start + begin)
        return list(struct.unpack("<%df" % ((end - begin) // 4), self.f.read(end - begin)))


def dot(a, b):
    return sum(map(operator.mul, a, b))


def matmul(x, w, rows, inner, cols):
    # Same per-element k order as osprey_gpu_matmul: out[r][c] = Σ_k x[r][k] w[k][c].
    wcols = [w[c::cols] for c in range(cols)]
    return [dot(x[r * inner:(r + 1) * inner], wcols[c]) for r in range(rows) for c in range(cols)]


def linear(x, rows, inner, cols, w, b):
    y = matmul(x, w, rows, inner, cols)
    return [y[i] + b[i % cols] for i in range(rows * cols)]


def layer_norm(x, t, g, b):
    out = []
    for r in range(t):
        row = x[r * WIDTH:(r + 1) * WIDTH]
        mean = sum(row) / 768.0
        var = sum((v - mean) * (v - mean) for v in row) / 768.0
        s = math.sqrt(var + 1e-5)
        out.extend((v - mean) / s * g[c] + b[c] for c, v in enumerate(row))
    return out


def gelu(x):
    return 0.5 * x * (1.0 + math.tanh(0.7978845608028654 * (x + 0.044715 * x * x * x)))


def attention(qkv, t):
    out = [0.0] * (t * WIDTH)
    for h in range(HEADS):
        for r in range(t):
            q = qkv[r * 2304 + h * 64:r * 2304 + h * 64 + 64]
            scores = [dot(q, qkv[c * 2304 + 768 + h * 64:c * 2304 + 768 + h * 64 + 64]) * 0.125 for c in range(r + 1)]
            peak = max(scores)
            weights = [math.exp(s - peak) for s in scores]
            total = sum(weights)
            probs = [w / total for w in weights]
            for d in range(64):
                out[r * WIDTH + h * 64 + d] = sum(probs[c] * qkv[c * 2304 + 1536 + h * 64 + d] for c in range(r + 1))
    return out


def forward(st, layers, ids):
    t = len(ids)
    wte, wpe = st["wte"], st["wpe"]
    x = [wte[ids[i // WIDTH] * WIDTH + i % WIDTH] + wpe[i] for i in range(t * WIDTH)]
    for l in layers:
        h1 = layer_norm(x, t, l["ln_1.weight"], l["ln_1.bias"])
        qkv = linear(h1, t, WIDTH, 2304, l["attn.c_attn.weight"], l["attn.c_attn.bias"])
        a = linear(attention(qkv, t), t, WIDTH, WIDTH, l["attn.c_proj.weight"], l["attn.c_proj.bias"])
        x = [u + v for u, v in zip(x, a)]
        h2 = layer_norm(x, t, l["ln_2.weight"], l["ln_2.bias"])
        hidden = [gelu(v) for v in linear(h2, t, WIDTH, 3072, l["mlp.c_fc.weight"], l["mlp.c_fc.bias"])]
        m = linear(hidden, t, 3072, WIDTH, l["mlp.c_proj.weight"], l["mlp.c_proj.bias"])
        x = [u + v for u, v in zip(x, m)]
    xf = layer_norm(x, t, st["ln_f.weight"], st["ln_f.bias"])
    last = xf[(t - 1) * WIDTH:]
    logits = [dot(last, wte[v * WIDTH:(v + 1) * WIDTH]) for v in range(VOCAB)]
    return logits


def main():
    weights, prompt, count = sys.argv[1], sys.argv[2], int(sys.argv[3])
    tok = Tokenizer(weights + "/merges.txt")
    sf = Safetensors(weights + "/model.safetensors")
    st = {k: sf.tensor(k + ".weight" if k in ("wte", "wpe") else k) for k in ("wte", "wpe", "ln_f.weight", "ln_f.bias")}
    names = ["ln_1.weight", "ln_1.bias", "attn.c_attn.weight", "attn.c_attn.bias", "attn.c_proj.weight",
             "attn.c_proj.bias", "ln_2.weight", "ln_2.bias", "mlp.c_fc.weight", "mlp.c_fc.bias",
             "mlp.c_proj.weight", "mlp.c_proj.bias"]
    layers = [{n: sf.tensor("h.%d.%s" % (i, n)) for n in names} for i in range(LAYERS)]
    ids = tok.encode(prompt)
    print("tokens:", " ".join(map(str, ids)))
    for _ in range(count):
        logits = forward(st, layers, ids)
        nxt = max(range(VOCAB), key=lambda v: (logits[v], -v))
        print("next", nxt, repr(tok.decode([nxt])), "logit", repr(logits[nxt]), flush=True)
        if nxt == END:
            break
        ids.append(nxt)
    print(tok.decode(ids))


if __name__ == "__main__":
    main()
