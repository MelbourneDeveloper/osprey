# osprey: flavor=python
"""Hunt driver: transformers' shared numerics — softmax, log_softmax, argmax
and multinomial — over lists mixing two symbolic logits with the concrete
special values NaN, +inf, -inf, huge (exp 700) and subnormal (exp -745), plus a
symbolic uniform draw. The special values guarantee the non-finite paths are
explored on every run."""
from transformers.Math import argmax, log_softmax, multinomial, softmax


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
    u = unwrap_or(parseFloat(input()), 0.5)
    inf = exp(1000.0)
    ninf = 0.0 - inf
    nan = inf - inf
    huge = exp(700.0)
    subnormal = exp(-745.0)
    logits = [x0, x1, inf, ninf]
    logits2 = [x0, nan, huge, subnormal]
    probs = softmax(logits)
    logprobs = log_softmax(logits)
    best = argmax(logits)
    drawn = multinomial(probs, u)
    p2 = softmax(logits2)
    drawn2 = multinomial(p2, u)
    return print(len(probs) + len(logprobs) + best + drawn + len(p2) + drawn2)


def main():
    return arith_total(run)
