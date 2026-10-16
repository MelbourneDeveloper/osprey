# osprey: flavor=python
# osprey: namespace=gpt2
# osprey: module=PythonModel
"""GPT-2 forward mathematics in Osprey Python flavor.

This is Python concrete syntax with Osprey semantics, not a CPython module.
Forward mathematics follows examples/gpt2/reference.py. Existing Osprey
modules supply checkpoint data, representation-neutral records and handlers.
"""

from gpt2 import Model as Base
from gpt2.Backends import Tensor

WIDTH = Base.width
HEADS = Base.heads
HEAD_WIDTH = Base.headWidth
LAYERS = Base.layerCount
VOCAB = Base.vocabSize
QKV = Base.qkvWidth
MLP = Base.mlpWidth
SCALE = Base.scoreScale
EPS = Base.epsilon
GELU_SCALE = Base.geluScale
GELU_CUBIC = Base.geluCubic


def at(buffer, index):
    return perform(Tensor.get(buffer, index))


def add_all(size, value_at):
    """Preserve the reference's left-to-right addition from 0.0."""
    return fold(range(size), 0.0, lambda acc, i: acc + value_at(i))


@export
def matmul(x, w, rows, inner, cols):
    return perform(Tensor.matMul(x, w, rows, inner, cols))


@export
def transpose(m, rows, cols):
    return perform(Tensor.build(rows * cols, lambda i: at(m, (i % rows) * cols + i // rows)))


@export
def linear(x, rows, inner, cols, w, b):
    y = matmul(x, w, rows, inner, cols)
    return perform(Tensor.build(rows * cols, lambda i: at(y, i) + at(b, i % cols)))


def row_mean(x, row):
    return add_all(WIDTH, lambda c: at(x, row * WIDTH + c)) / toFloat(WIDTH)


def row_variance(x, means, row):
    mean = at(means, row)
    return add_all(WIDTH, lambda c: square(at(x, row * WIDTH + c) - mean)) / toFloat(WIDTH)


def square(value):
    return value * value


def row_stats(x, t):
    means = perform(Tensor.build(t, lambda row: row_mean(x, row)))
    variances = perform(Tensor.build(t, lambda row: row_variance(x, means, row)))
    return Statistics(means=means, variances=variances)


class Statistics[B]:
    means: B
    variances: B


def normalized(x, means, variances, g, b, i):
    row = i // WIDTH
    col = i % WIDTH
    scale = sqrt(at(variances, row) + EPS)
    return (at(x, i) - at(means, row)) / scale * at(g, col) + at(b, col)


@export
def layer_norm(x, t, g, b):
    stats = row_stats(x, t)
    means = stats.means
    variances = stats.variances
    return perform(Tensor.build(t * WIDTH, lambda i: normalized(x, means, variances, g, b, i)))


@export
def gelu(x):
    return 0.5 * x * (1.0 + tanh(GELU_SCALE * (x + GELU_CUBIC * x * x * x)))


def q_at(qkv, row, head, d):
    return at(qkv, row * QKV + head * HEAD_WIDTH + d)


def k_at(qkv, row, head, d):
    return at(qkv, row * QKV + WIDTH + head * HEAD_WIDTH + d)


def v_at(qkv, row, head, d):
    return at(qkv, row * QKV + 2 * WIDTH + head * HEAD_WIDTH + d)


def score(qkv, t, i):
    head = i // (t * t)
    row = (i // t) % t
    col = i % t
    return add_all(HEAD_WIDTH, lambda d: q_at(qkv, row, head, d) * k_at(qkv, col, head, d)) * SCALE if col <= row else 0.0


def max_value(a, b):
    return b if b > a else a


def row_peak(scores, t, row):
    return fold(range(row % t + 1), at(scores, row * t), lambda peak, c: max_value(peak, at(scores, row * t + c)))


def causal_weight(scores, peaks, t, i):
    return exp(at(scores, i) - at(peaks, i // t)) if i % t <= (i // t) % t else 0.0


def probabilities(qkv, t):
    scores = perform(Tensor.build(HEADS * t * t, lambda i: score(qkv, t, i)))
    peaks = perform(Tensor.build(HEADS * t, lambda row: row_peak(scores, t, row)))
    weights = perform(Tensor.build(HEADS * t * t, lambda i: causal_weight(scores, peaks, t, i)))
    totals = perform(Tensor.build(HEADS * t, lambda row: add_all(row % t + 1, lambda c: at(weights, row * t + c))))
    return perform(Tensor.build(HEADS * t * t, lambda i: at(weights, i) / at(totals, i // t)))


def attended(probs, qkv, t, i):
    row = i // WIDTH
    head = (i % WIDTH) // HEAD_WIDTH
    d = i % HEAD_WIDTH
    return add_all(row + 1, lambda col: at(probs, (head * t + row) * t + col) * v_at(qkv, col, head, d))


@export
def attention(qkv, t):
    probs = probabilities(qkv, t)
    return Base.Attn(probs=probs, out=perform(Tensor.build(t * WIDTH, lambda i: attended(probs, qkv, t, i))))


@export
def residual(x, part):
    return perform(Tensor.zipWith(x, part, lambda u, v: u + v))


class Projection[B]:
    normalized: B
    values: B


def normalized_projection(x, t, cols, g, b, w, bias):
    h = layer_norm(x, t, g, b)
    return Projection(normalized=h, values=linear(h, t, WIDTH, cols, w, bias))


@export
def block(x, t, layer):
    first = normalized_projection(x, t, QKV, layer.ln1g, layer.ln1b, layer.qkvW, layer.qkvB)
    h1 = first.normalized
    qkv = first.values
    attn = attention(qkv, t)
    x2 = residual(x, linear(attn.out, t, WIDTH, WIDTH, layer.projW, layer.projB))
    second = normalized_projection(x2, t, MLP, layer.ln2g, layer.ln2b, layer.fcW, layer.fcB)
    h2 = second.normalized
    f = second.values
    u = perform(Tensor.build(t * MLP, lambda i: gelu(at(f, i))))
    y = residual(x2, linear(u, t, MLP, WIDTH, layer.outW, layer.outB))
    return Base.Acts(x=x, h1=h1, qkv=qkv, probs=attn.probs, a=attn.out, x2=x2, h2=h2, f=f, u=u, y=y)


def advance_block(w, t, trace, index):
    acts = block(trace.x, t, Base.layer(w, index))
    return Base.Trace(acts=listAppend(trace.acts, acts), x=acts.y)


@export
def forward_trace(w, ids):
    t = len(ids)
    initial = Base.Trace(acts=[], x=Base.embed(w, ids, t))
    return fold(range(LAYERS), initial, lambda trace, i: advance_block(w, t, trace, i))


@export
def final_norm(w, ids):
    return layer_norm(forward_trace(w, ids).x, len(ids), w.lnfG, w.lnfB)


@export
def forward(w, ids):
    xf = final_norm(w, ids)
    start = (len(ids) - 1) * WIDTH
    last = perform(Tensor.build(WIDTH, lambda c: at(xf, start + c)))
    return matmul(last, w.wteT, 1, WIDTH, VOCAB)
