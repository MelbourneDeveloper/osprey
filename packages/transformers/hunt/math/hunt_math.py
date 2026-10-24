# osprey: flavor=python
"""Hunt driver: transformers' shared numerics — softmax, log_softmax, argmax
and multinomial — over four symbolic logits (NaN/inf/subnormal/huge/finite) and
a symbolic uniform draw. Integer index arithmetic is discharged by a local
value-mode `Arith` handler."""
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
    x2 = unwrap_or(parseFloat(input()), 0.0)
    x3 = unwrap_or(parseFloat(input()), 0.0)
    u = unwrap_or(parseFloat(input()), 0.5)
    logits = [x0, x1, x2, x3]
    probs = softmax(logits)
    logprobs = log_softmax(logits)
    best = argmax(logits)
    drawn = multinomial(probs, u)
    return print(len(probs) + len(logprobs) + best + drawn)


def main():
    return arith_total(run)
