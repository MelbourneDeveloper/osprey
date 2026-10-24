# osprey: flavor=python
"""Hunt driver: transformers' cheap float warpers — temperature (with a
symbolic temperature, so 0.0 reaches the divide-by-zero policy), min_p and
epsilon — over a list mixing symbolic logits with the concrete special values
NaN, +inf, -inf, huge and subnormal."""
from transformers.Math import neg_inf
from transformers.LogitsProcess import (
    EpsilonLogitsWarper, MinPLogitsWarper, TemperatureLogitsWarper, apply,
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
    t = unwrap_or(parseFloat(input()), 1.0)
    inf = exp(1000.0)
    ninf = 0.0 - inf
    nan = inf - inf
    logits = [x0, x1, inf, ninf, nan]
    a = apply(TemperatureLogitsWarper(t), [], logits)
    c = apply(MinPLogitsWarper(0.1, neg_inf(), 1), [], logits)
    e = apply(EpsilonLogitsWarper(0.1, neg_inf(), 1), [], logits)
    return print(len(a) + len(c) + len(e))


def main():
    return arith_total(run)
