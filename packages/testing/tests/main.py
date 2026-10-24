# osprey: flavor=python
"""@osprey/testing's own suite: `Check.near` against numpy.isclose's verdicts
(scripts/oracles/testing.py), then the reporting helpers."""
from testing import Check
from testing_expected import Expected

RTOL = 0.00001
ATOL = 0.00000001


def inf():
    return exp(1000.0)


def value_of(text):
    """The oracle spells non-finite operands as Python prints them."""
    match text:
        case "inf":
            return inf()
        case "-inf":
            return 0.0 - inf()
        case "nan":
            return inf() - inf()
        case _:
            return unwrap_or(parseFloat(text), 0.0)


def verdict(name, pair, want):
    a = unwrap_or(listGet(pair, 0), 0.0)
    b = unwrap_or(listGet(pair, 1), 0.0)
    return Check.equal(f"{name}({a}, {b})", Check.near(a, b, RTOL, ATOL), want)


def special(i):
    pair = unwrap_or(listGet(Expected.SPECIAL, i), [])
    values = [value_of(unwrap_or(listGet(pair, 0), "")), value_of(unwrap_or(listGet(pair, 1), ""))]
    return verdict("near.special", values, unwrap_or(listGet(Expected.SPECIAL_VERDICTS, i), False))


def oracle_cases():
    finite = toList(map(range(len(Expected.PAIRS)), lambda i: verdict("near", unwrap_or(listGet(Expected.PAIRS, i), []), unwrap_or(listGet(Expected.VERDICTS, i), False))))
    return listConcat(finite, toList(map(range(len(Expected.SPECIAL)), special)))


def helper_cases():
    return [
        Check.close_all("close_all.lists", [1.0, 2.0, 3.0], [1.0, 2.0, 3.0000001], RTOL, ATOL),
        Check.truth("close_all.length", not Check.close_all("x", [1.0], [1.0, 2.0], RTOL, ATOL).ok, "lengths differ"),
        Check.equal("close_all.detail", Check.close_all("x", [1.0, 5.0], [1.0, 2.0], RTOL, ATOL).detail, "index 1: got 5.0, want 2.0"),
        Check.ints_equal("ints_equal", [1, 2, 3], [1, 2, 3]),
        Check.equal("ints_equal.detail", Check.ints_equal("x", [1, 2, 4], [1, 2, 3]).detail, "index 2: got [1, 2, 4], want [1, 2, 3]"),
        Check.equal("ints_text.empty", Check.ints_text([]), "[]"),
        Check.equal("summary", Check.summary([Check.passed("a"), Check.failed("b", "c")]), "1 passed, 1 failed"),
    ]


def main():
    match Check.strictly(lambda: Check.report(listConcat(oracle_cases(), helper_cases()))):
        case Success(value=value):
            return value
        case Error(message=message):
            return print(f"suite stopped: {message}")
