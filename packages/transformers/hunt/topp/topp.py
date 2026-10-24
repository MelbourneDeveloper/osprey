# osprey: flavor=python
"""Hunt driver: transformers' sort-and-cumsum warpers — top_p (with a symbolic
top_p), typical and eta — over a list mixing symbolic logits with the concrete
special values NaN, +inf, -inf, huge and subnormal."""
from transformers.Math import neg_inf
from transformers.LogitsProcess import (
    EtaLogitsWarper, TopPLogitsWarper, TypicalLogitsWarper, apply,
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
    p = unwrap_or(parseFloat(input()), 0.9)
    inf = exp(1000.0)
    ninf = 0.0 - inf
    nan = inf - inf
    logits = [x0, x1, inf, ninf, nan]
    b = apply(TopPLogitsWarper(p, neg_inf(), 1), [], logits)
    d = apply(TypicalLogitsWarper(0.9, neg_inf(), 1), [], logits)
    f = apply(EtaLogitsWarper(0.1, neg_inf(), 1), [], logits)
    return print(len(b) + len(d) + len(f))


def main():
    return arith_total(run)
