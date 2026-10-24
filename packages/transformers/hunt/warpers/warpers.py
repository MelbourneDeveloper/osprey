# osprey: flavor=python
"""Hunt driver: transformers' float logits warpers — temperature (with a
symbolic temperature, so 0.0 reaches the divide-by-zero policy), top_p, min_p,
typical, epsilon and eta — over four symbolic logits (NaN/inf/subnormal/huge)."""
from transformers.Math import neg_inf
from transformers.LogitsProcess import (
    EpsilonLogitsWarper, EtaLogitsWarper, MinPLogitsWarper, TemperatureLogitsWarper,
    TopPLogitsWarper, TypicalLogitsWarper, apply,
)


@handler(Arith)
class arith_total:
    def overflow(op, l, r, wrapped):
        return wrapped

    def divideByZero(op, l):
        return 0.0

    def remainderByZero(l):
        return 0


def run():
    x0 = unwrap_or(parseFloat(input()), 0.0)
    x1 = unwrap_or(parseFloat(input()), 0.0)
    x2 = unwrap_or(parseFloat(input()), 0.0)
    x3 = unwrap_or(parseFloat(input()), 0.0)
    t = unwrap_or(parseFloat(input()), 1.0)
    p = unwrap_or(parseFloat(input()), 0.9)
    logits = [x0, x1, x2, x3]
    a = apply(TemperatureLogitsWarper(t), [], logits)
    b = apply(TopPLogitsWarper(p, neg_inf(), 1), [], logits)
    c = apply(MinPLogitsWarper(0.1, neg_inf(), 1), [], logits)
    d = apply(TypicalLogitsWarper(0.9, neg_inf(), 1), [], logits)
    e = apply(EpsilonLogitsWarper(0.1, neg_inf(), 1), [], logits)
    f = apply(EtaLogitsWarper(0.1, neg_inf(), 1), [], logits)
    return print(len(a) + len(b) + len(c) + len(d) + len(e) + len(f))


def main():
    return arith_total(run)
