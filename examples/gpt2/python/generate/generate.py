# osprey: flavor=python
"""Greedy GPT-2 generation in Osprey Python flavor."""

from gpt2 import Backends, Model, PythonModel, Tokenizer
from gpt2.Backends import Tensor


class Sampling:
    ids: list[int]
    done: bool


def stop(message):
    _ = Backends.orAbort(Error(message=message))


def peak(a, b):
    return b if b > a else a


def next_token(w, ids):
    logits = PythonModel.forward(w, ids)
    first = perform(Tensor.get(logits, 0))
    best = perform(Tensor.reduce(logits, first, peak))
    return perform(Tensor.indexOf(logits, best))


def step(w, state):
    if state.done or len(state.ids) >= Model.positions:
        return Sampling(ids=state.ids, done=True)
    else:
        token = next_token(w, state.ids)
        if token == Model.endOfText:
            return Sampling(ids=state.ids, done=True)
        else:
            return Sampling(ids=listAppend(state.ids, token), done=False)


def sample(w, ids, count):
    initial = Sampling(ids=ids, done=False)
    return fold(range(count), initial, lambda state, _i: step(w, state)).ids


def run(prompt, count):
    bpe = Tokenizer.loadBpe(Model.weightsDir + "/merges.txt")
    ids = Tokenizer.encode(bpe, prompt)
    if len(ids) == 0 or len(ids) > Model.positions:
        stop(f"prompt must contain 1..{Model.positions} tokens")
    w = Model.loadWeights()
    out = sample(w, ids, count if count > 0 else 0)
    print(f"gpt2 (124M) on {Model.device()}\nprompt: {prompt}\ntokens:{Tokenizer.idsText(ids)}\ngenerated {len(out) - len(ids)} tokens\n{Tokenizer.decode(bpe, out)}")


def token_count(line):
    match parseInt(line):
        case Success(value=value):
            return value
        case Error():
            return 20


def main():
    prompt_line = input()
    prompt = prompt_line if prompt_line != "" else "The capital of France is"
    count = token_count(input())
    backend = input()
    if backend != "gpu" and backend != "lists":
        stop(f"unknown backend {backend}: choose gpu or lists")
    match Backends.strictly(lambda: Backends.under(backend, lambda: run(prompt, count))):
        case Success(value=value):
            return value
        case Error(message=message):
            return stop(f"generation stopped: {message}")
