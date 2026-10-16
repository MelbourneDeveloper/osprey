# osprey: flavor=python
"""Check short deterministic sequences at GPT-2's actual tensor dimensions."""
from gpt2 import Model as Base, Backends, PythonModel as M
from gpt2.Backends import Tensor


def vector(size, scale):
    return perform(Tensor.build(size, lambda i: toFloat(i % 7) * scale))


def line(name, values):
    size = perform(Tensor.length(values))
    text = join(toList(map(range(size), lambda i: toString(perform(Tensor.get(values, i))))), ",")
    return name + ":" + text + "\n"


def layer():
    return Base.Layer(
        ln1g=vector(Base.width, 0.1), ln1b=vector(Base.width, -0.05),
        qkvW=vector(Base.width * Base.qkvWidth, 0.01), qkvB=vector(Base.qkvWidth, -0.02),
        projW=vector(Base.width * Base.width, 0.02), projB=vector(Base.width, 0.03),
        ln2g=vector(Base.width, 0.2), ln2b=vector(Base.width, -0.01),
        fcW=vector(Base.width * Base.mlpWidth, 0.01), fcB=vector(Base.mlpWidth, 0.02),
        outW=vector(Base.mlpWidth * Base.width, -0.01), outB=vector(Base.width, 0.01),
    )


def sequence(t):
    x = vector(t * Base.width, 0.25)
    g = vector(Base.width, 0.1)
    b = vector(Base.width, -0.05)
    qkv = vector(t * Base.qkvWidth, 0.125)
    attn = M.attention(qkv, t)
    return line(f"norm{t}", M.layer_norm(x, t, g, b)) + line(f"probs{t}", attn.probs) + line(f"attention{t}", attn.out)


def block_lines():
    acts = M.block(vector(3 * Base.width, 0.25), 3, layer())
    return line("block.x", acts.x) + line("block.h1", acts.h1) + line("block.qkv", acts.qkv) + line("block.probs", acts.probs) + line("block.a", acts.a) + line("block.x2", acts.x2) + line("block.h2", acts.h2) + line("block.f", acts.f) + line("block.u", acts.u) + line("block.y", acts.y)


def run():
    a = vector(6, 0.25)
    b = vector(6, -0.125)
    bias = vector(2, 0.1)
    gelus = perform(Tensor.build(7, lambda i: M.gelu(toFloat(i - 3) / 2.0)))
    small = line("transpose", M.transpose(a, 2, 3)) + line("matmul", M.matmul(a, b, 2, 3, 2)) + line("linear", M.linear(a, 2, 3, 2, b, bias)) + line("residual", M.residual(a, b)) + line("gelu", gelus)
    print("@device:" + Base.device() + "\n" + small + sequence(1) + sequence(2) + sequence(3) + block_lines())


def main():
    backend = input()
    match Backends.strictly(lambda: Backends.under(backend, run)):
        case Success(value=value):
            return value
        case Error(message=message):
            _ = Backends.orAbort(Error(message=message))
            return None
