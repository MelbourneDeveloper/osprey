# @osprey/testing

Assertions for the Python-flavor package ports ([spec 0043](../../docs/specs/ecosystem/0043-PythonFlavorPorts.md)). Each check returns an `Outcome` value; `Check.report` prints one `ok`/`FAIL` line per outcome and a `N passed, M failed` total that a golden transcript pins.

| Function | Original | Notes |
| --- | --- | --- |
| `Check.near(got, want, rtol, atol)` | `numpy.isclose(got, want, rtol, atol, equal_nan=True)` | Equal infinities are equal, opposite infinities are not, NaN equals NaN |
| `Check.close`, `Check.close_all` | `numpy.isclose`, `numpy.allclose` | `close_all` names the first differing index |
| `Check.equal`, `Check.ints_equal`, `Check.truth` | `assert a == b` | Exact |
| `Check.strictly(work)` | — | Runs a suite under a policy that stops at the first integer overflow or zero division |

The suite in `tests/` compares each verdict with real `numpy.isclose`; `scripts/oracles/testing.py` regenerates `tests/expected.py` from it.
